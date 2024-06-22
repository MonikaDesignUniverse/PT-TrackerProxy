//go:build (darwin || windows || linux)
package main

import (
	"net/url"
	"net/http"
	"github.com/bdwyertech/go-get-proxied/proxy"
)

var httpProxyURL *url.URL
var httpsProxyURL *url.URL

func GetProxy(r *http.Request) (*url.URL, error) {
	if httpProxyURL == nil || httpsProxyURL == nil {
		proxyProvider := proxy.NewProvider("")

		if httpProxyURL == nil {
			httpProxy := proxyProvider.GetHTTPProxy("")
			if httpProxy == nil {
				return nil, nil
			}

			Log("GetProxy", "发现 HTTP 代理: %s", httpProxy.Src(), httpProxy.String())
			httpProxyURL = httpProxy.URL()
		}

		if httpsProxyURL == nil {
			httpsProxy := proxyProvider.GetHTTPSProxy("")
			if httpsProxy == nil {
				return nil, nil
			}

			Log("GetProxy", "发现 HTTPS 代理: %s", httpsProxy.Src(), httpsProxy.String())
			httpsProxyURL = httpsProxy.URL()
		}
	}

	if r != nil {
		if r.URL.Scheme == "https" {
			return httpsProxyURL, nil
		} else if r.URL.Scheme == "http" {
			return httpProxyURL, nil
		}
	}

	return nil, nil
}
