//go:build (!darwin && !windows && !linux)
package main

import (
	"net/url"
	"net/http"
)

func GetProxy(r *http.Request) (*url.URL, error) {
	return http.ProxyFromEnvironment(r)
}
