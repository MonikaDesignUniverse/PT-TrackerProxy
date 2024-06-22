//go:build (!darwin && !windows && !linux)
package main

import (
	"net/url"
	"net/http"
)
var httpProxyURL *url.URL
var httpsProxyURL *url.URL

func GetProxy(r *http.Request) (*url.URL, error) {
	return http.ProxyFromEnvironment(r)
}
