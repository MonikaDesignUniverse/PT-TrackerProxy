package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"net"
	"net/url"
	"net/http"
	"net/http/httputil"
)

type Config struct {
	Listen string
}

var debug = false
var configFilename = "config.json"
var userAgent = "PT-TrackerProxy/1.0"
var reserveURL = "http://t.acg.rip:6699/announce"
var config = Config {
	Listen: "127.0.0.1:8765",
}

var isServerRunning = true
var currentIPv4 = ""
var currentIPv6 = ""
var intervalTicker *time.Ticker
var httpTransport = &http.Transport {
	DisableKeepAlives: true,
	ForceAttemptHTTP2: false,
}
var httpClient = http.Client {
	Timeout:   6 * time.Second,
	Transport: httpTransport,
}
var reserveServer = &http.Server {
}

func StrTrim(str string) string {
	return strings.Trim(str, " \n\r")
}
func GetDateTime(withTime bool) string {
	formatStr := "2006-01-02"
	if withTime {
		formatStr += " 15:04:05"
	}
	return time.Now().Format(formatStr)
}
func CheckIP(ip string) uint32 {
	ipParsed := net.ParseIP(ip)
	if ipParsed == nil {
		return 0
	}

	if strings.Count(ip, ":") < 2 {
		return 4
	}

	return 6
}
func LoadConfig() bool {
	_, err := os.Stat(configFilename)
	if err != nil {
		Log("Debug-LoadConfig", "读取配置文件元数据时发生了错误: %s", false, err.Error())
		return false
	}
	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		Log("Debug-LoadConfig", "读取配置文件时发生了错误: %s", false, err.Error())
		return false
	}
	if err := json.Unmarshal(configFile, &config); err != nil {
		Log("Debug-LoadConfig", "解析配置文件时发生了错误: %s", false, err.Error())
		return false
	}
	Log("LoadConfig", "读取配置文件成功", true)
	return true
}
func Log(module string, str string, args ...interface {}) {
	if !debug && strings.HasPrefix(module, "Debug") {
		return
	}
	logStr := fmt.Sprintf("[" + GetDateTime(true) + "][" + module + "] " + str + ".\n", args...)
	fmt.Print(logStr)
}
func RefreshCurrentIPv4() bool {
	ipv4ResponseBody := Fetch("https://api-ipv4.ip.sb/ip")
	if ipv4ResponseBody == nil {
		Log("RefreshCurrentIPv4", "获取 IPv4 地址时发生了错误 (Error 1)")
		return false
	}
	ipv4ResponseBodyStr := StrTrim(string(ipv4ResponseBody))
	if CheckIP(ipv4ResponseBodyStr) != 4 {
		Log("RefreshCurrentIPv4", "获取 IPv4 地址时发生了错误 (Error 2)")
		return false
	}
	currentIPv4 = ipv4ResponseBodyStr
	Log("RefreshCurrentIPv4", "获取当前 IPv4 地址: %s", ipv4ResponseBodyStr)
	return true
}
func RefreshCurrentIPv6() bool {
	ipv6ResponseBody := Fetch("https://api-ipv6.ip.sb/ip")
	if ipv6ResponseBody == nil {
		Log("RefreshCurrentIPv6", "获取 IPv6 地址时发生了错误 (Error 1)")
		return false
	}
	ipv6ResponseBodyStr := StrTrim(string(ipv6ResponseBody))
	if CheckIP(ipv6ResponseBodyStr) != 6 {
		Log("RefreshCurrentIPv6", "获取 IPv6 地址时发生了错误 (Error 2)")
		return false
	}
	currentIPv6 = ipv6ResponseBodyStr
	Log("RefreshCurrentIPv6", "获取当前 IPv6 地址: %s", ipv6ResponseBodyStr)
	return true
}
func RefreshCurrentIP() {
	intervalTicker = time.NewTicker(time.Duration(900) * time.Second)
	for ; true; <-intervalTicker.C {
		if !RefreshCurrentIPv4() {
			currentIPv4 = ""
		}
		if !RefreshCurrentIPv6() {
			currentIPv6 = ""
		}
	}
}
func StartProxy() {
	parsedReserveURL, err := url.Parse(reserveURL)
	if err != nil {
		Log("StartProxy", "监听本地 URL 时发生错误: %s", err.Error())
		return
	}
	reserveProxy := &httputil.ReverseProxy {
        Rewrite: func(r *httputil.ProxyRequest) {
            r.SetURL(parsedReserveURL)
            if strings.Contains(r.Out.URL.RawQuery, "&") {
	            if currentIPv4 != "" {
	            	r.Out.URL.RawQuery += ("&pttp_ip4=" + currentIPv4)
	            }
	            if currentIPv6 != "" {
	            	r.Out.URL.RawQuery += ("&pttp_ip6=" + currentIPv6)
	            }
	        }
        },
    }
	http.Handle("/announce", reserveProxy)
	http.Handle("/announce/", reserveProxy)
	Log("StartProxy", "监听于: %s, 服务于: %s", config.Listen, parsedReserveURL)
	reserveServer.Addr = config.Listen
	for isServerRunning {
		if err := reserveServer.ListenAndServe(); err != http.ErrServerClosed {
			Log("StartProxy", "处理请求时发生错误: %s", err.Error())
		}
	}
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
	go CatchSignal()
	go RefreshCurrentIP()
	StartProxy()
}
