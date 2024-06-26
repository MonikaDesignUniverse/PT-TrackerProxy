package main

import (
	"os"
	"time"
	"encoding/json"
	"net/url"
	"net/http"
	"github.com/tidwall/jsonc"
)

type Config struct {
	Debug      bool
	ListenAddr string
	ListenPort int
}

var programName = "PT-TrackerProxy"
var programVersion = "Unknown"
var programUserAgent = programName + "/" + programVersion

var domain_whitelist = map[string]*url.URL {
	"monikadesign.uk": ParseURL("https://monikadesign.uk"),
	"tracker.monikadesign.uk": ParseURL("https://tracker.monikadesign.uk"),
	"daisuki.monikadesign.uk": ParseURL("https://daisuki.monikadesign.uk"),
	"daikirai.monikadesign.uk": ParseURL("https://daikirai.monikadesign.uk"),
}

var configFilename = "config.json"
var config = Config {
	Debug:      false,
	ListenAddr: "127.0.0.1",
	ListenPort: 7887,
}

var httpTransport = &http.Transport {
	DisableKeepAlives:     false,
	ForceAttemptHTTP2:     false,
	MaxIdleConns:          200,
	MaxConnsPerHost:       200,
	MaxIdleConnsPerHost:   200,
	IdleConnTimeout:       60 * time.Second,
	TLSHandshakeTimeout:   12 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	Proxy:                 GetProxy,
}
var httpTransportWithoutProxy = &http.Transport {
	DisableKeepAlives:     true,
	ForceAttemptHTTP2:     false,
	MaxIdleConns:          200,
	MaxConnsPerHost:       200,
	MaxIdleConnsPerHost:   200,
	IdleConnTimeout:       60 * time.Second,
	TLSHandshakeTimeout:   12 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	Proxy:                 nil,
}
var httpClient = &http.Client {
	Timeout:   6 * time.Second,
	Transport: httpTransportWithoutProxy,
}
var reserveServer = &http.Server {
	ReadTimeout:  60 * time.Second,
	WriteTimeout: 60 * time.Second,
	IdleTimeout:  60 * time.Second,
}

func ShowVersion() {
	Log("ShowVersion", "%s %s", programName, programVersion)
}
func LoadConfig() bool {
	_, err := os.Stat(configFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			Log("LoadConfig", "读取配置文件元数据时发生了错误: %s", err.Error())
		}
		return false
	}
	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		Log("LoadConfig", "读取配置文件时发生了错误: %s", err.Error())
		return false
	}
	if err := json.Unmarshal(jsonc.ToJSON(configFile), &config); err != nil {
		Log("LoadConfig", "解析配置文件时发生了错误: %s", err.Error())
		return false
	}
	Log("LoadConfig", "读取配置文件成功")

	return true
}
