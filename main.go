package main

import (
	"io"
	"os"
	"os/signal"
	"log"
	"fmt"
	"time"
	"strings"
	"strconv"
	"syscall"
	"context"
	"net"
	"net/url"
	"net/http"
	"net/http/httputil"
)

var isServerRunning = true
var currentIPv4 = ""
var currentIPv6 = ""
var intervalTicker = time.NewTicker(time.Duration(900) * time.Second)
var logwriter = LogWriter {}

var reserveProxy = &httputil.ReverseProxy {
	Transport: httpTransport,
	Rewrite: func(r *httputil.ProxyRequest) {
		ctx := r.In.Context()
		ctx_parsedReserveURL := ctx.Value("parsedReserveURL")
		if ctx_parsedReserveURL == nil {
			return
		}

		parsedReserveURL := ctx_parsedReserveURL.(*url.URL)
		r.SetURL(parsedReserveURL)

		r.Out.Header.Set("X-PTTP-Version", programVersion)
		if currentIPv4 != "" {
			r.Out.Header.Set("X-PTTP-IP4", currentIPv4)
		}
		if currentIPv6 != "" {
			r.Out.Header.Set("X-PTTP-IP6", currentIPv6)
		}
	},
}

type LogWriter struct {
    w io.Writer
}

func (w LogWriter) Write(p []byte) (n int, err error) {
	Log("LogWriter", string(p))
	return len(p), nil
}
func Log(module string, str string, args ...interface {}) {
	if !debug && strings.HasPrefix(module, "Debug") {
		return
	}
	if module == "LogWriter" {
		str = StrTrim(str)
		if str == "http: proxy error: EOF" {
			return
		}
	}
	logStr := fmt.Sprintf("[" + GetDateTime(true) + "][" + module + "] " + str + ".\n", args...)
	fmt.Print(logStr)
}
func StartProxy() {
	http.HandleFunc("/", ProcessRequest)

	domainArr := make([]string, 0, len(domain_whitelist))
	for k := range domain_whitelist {
		domainArr = append(domainArr, k)
	}

	listenStr := (config.ListenAddr + ":" + strconv.Itoa(listenPort))
	Log("StartProxy", "监听于: %s, 支持以下 Tracker: %s", listenStr, strings.Join(domainArr, " | "))
	reserveServer.Addr = listenStr
	listen, err := net.Listen("tcp4", listenStr)
	if err != nil {
		Log("StartProxy", "监听端口时发生错误: %s", err.Error())
		return
	}

	for isServerRunning {
		if err := reserveServer.Serve(listen); err != http.ErrServerClosed {
			Log("StartProxy", "处理请求时发生错误: %s", err.Error())
		}
	}
}
func WriteResponse(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(s + "\n"))
}
func ProcessRequest(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		WriteResponse(w, programUserAgent + ".")
		return
	}

	targetHostSplit := strings.SplitN(r.URL.Path, "/", 3)
	targetHostSplitLen := len(targetHostSplit)
	if targetHostSplitLen < 2 {
		WriteResponse(w, "坏目标.")
		return
	}

	targetHost := targetHostSplit[1]
	if _, exist := domain_whitelist[targetHost]; !exist {
		WriteResponse(w, "坏域名.")
		return
	}

	if targetHostSplitLen < 3 || !strings.HasPrefix(targetHostSplit[2], "announce") {
		WriteResponse(w, "坏路径.")
		return
	}

	reserveURL := ("https://" + targetHost)
	parsedReserveURL, err := url.Parse(reserveURL)
	if err != nil {
		WriteResponse(w, "解析目标 URL 时发生错误: " + err.Error())
		return
	}

	ctx := context.WithValue(r.Context(), "parsedReserveURL", parsedReserveURL)
	r = r.WithContext(ctx)
	r.URL.Path = targetHostSplit[2]
	r.URL.RawPath = ""
	r.RequestURI = r.URL.RequestURI()

	reserveProxy.ServeHTTP(w, r)
}
func CatchSignal() {
	exitSignal := make(chan os.Signal, 2)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	isServerRunning = false
	intervalTicker.Stop()
	reserveServer.Close()
	os.Exit(0)
}
func main() {
	LoadConfig()
	ShowVersion()
	log.SetFlags(0)
	log.SetOutput(logwriter)
	go CatchSignal()
	go RefreshCurrentIP()
	StartProxy()
}
