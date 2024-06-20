package main

import (
	"time"
	"strings"
	"net"
	"net/url"
)

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
func ParseURL(targetURL string) *url.URL {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		Log("ParseURL", "解析目标 URL 时发生错误: %s", err.Error())
		return nil
	}

	return parsedURL
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
