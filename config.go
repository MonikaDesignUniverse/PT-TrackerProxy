package main

import (
	"os"
	"time"
	"encoding/json"
	"net/http"
)

var debug = false

var programName = "PT-TrackerProxy"
var programVersion = "Unknown"
var programUserAgent = programName + "/" + programVersion

var domain_whitelist = map[string]bool {
	"monikadesign.uk": true,
	"tracker.monikadesign.uk": true,
	"daisuki.monikadesign.uk": true,
	"daikirai.monikadesign.uk": true,
}

var configFilename = "config.json"
var config = Config {
	ListenAddr: "127.0.0.1",
}
var listenPort = 7887

var httpTransport = &http.Transport {
	DisableKeepAlives:     false,
	ForceAttemptHTTP2:     false,
	MaxIdleConns:          200,
	MaxConnsPerHost:       200,
	MaxIdleConnsPerHost:   200,
	IdleConnTimeout:       60 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	Proxy:                 http.ProxyFromEnvironment,
}
var httpTransportWithoutProxy = &http.Transport {
	DisableKeepAlives:     true,
	ForceAttemptHTTP2:     false,
	MaxIdleConns:          200,
	MaxConnsPerHost:       200,
	MaxIdleConnsPerHost:   200,
	IdleConnTimeout:       60 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	Proxy:                 nil,
}
var httpClient = http.Client {
	Timeout:   6 * time.Second,
	Transport: httpTransportWithoutProxy,
}
var reserveServer = &http.Server {
	ReadTimeout:  60 * time.Second,
	WriteTimeout: 60 * time.Second,
	IdleTimeout:  60 * time.Second,
}

type Config struct {
	ListenAddr string
}

func ShowVersion() {
	Log("ShowVersion", "%s %s", programName, programVersion)
}
func LoadConfig() bool {
	_, err := os.Stat(configFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			Log("Debug-LoadConfig", "读取配置文件元数据时发生了错误: %s", err.Error())
		}
		return false
	}
	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		Log("Debug-LoadConfig", "读取配置文件时发生了错误: %s", err.Error())
		return false
	}
	if err := json.Unmarshal(configFile, &config); err != nil {
		Log("Debug-LoadConfig", "解析配置文件时发生了错误: %s", err.Error())
		return false
	}
	Log("LoadConfig", "读取配置文件成功")

	return true
}
