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

			httpProxyURL = httpProxy.URL()
			if httpProxyURL.Scheme == "" {
				httpProxyURL.Scheme = "http"
			}

			Log("GetProxy", "发现 HTTP 代理: %s (来源: %s)", httpProxyURL.String(), httpProxy.Src())
		}

		if httpsProxyURL == nil {
			httpsProxy := proxyProvider.GetHTTPSProxy("")
			if httpsProxy == nil {
				return nil, nil
			}

			httpsProxyURL = httpsProxy.URL()
			if httpsProxyURL.Scheme == "" {
				httpsProxyURL.Scheme = "http"
			}

			Log("GetProxy", "发现 HTTPS 代理: %s (来源: %s)", httpsProxyURL.String(), httpsProxy.Src())
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
