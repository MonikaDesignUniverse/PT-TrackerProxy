package main

import (
	"strings"
	"io/ioutil"
	"net/http"
)

func Fetch(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log("Fetch", "请求时发生了错误 (Part 1): %s", true, err.Error())
		return nil
	}
	req.Header.Set("User-Agent", userAgent)
	response, err := httpClient.Do(req)
	if err != nil {
		Log("Fetch", "请求时发生了错误 (Part 2): %s", true, err.Error())
		return nil
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		Log("Fetch", "读取时发生了错误: %s", true, err.Error())
		return nil
	}

	if response.StatusCode == 403 {
		Log("Fetch", "请求时发生了错误: 认证失败 %s", true, responseBody)
		return nil
	}

	if response.StatusCode == 404 {
		Log("Fetch", "请求时发生了错误: 资源不存在", true)
		return nil
	}

	return responseBody
}
func Submit(url string, postdata string) []byte {
	req, err := http.NewRequest("POST", url, strings.NewReader(postdata))
	if err != nil {
		Log("Submit", "请求时发生了错误 (Part 1): %s", true, err.Error())
		return nil
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := httpClient.Do(req)
	if err != nil {
		Log("Submit", "请求时发生了错误 (Part 2): %s", true, err.Error())
		return nil
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		Log("Submit", "读取时发生了错误: %s", true, err.Error())
		return nil
	}

	if response.StatusCode == 403 {
		Log("Submit", "请求时发生了错误: 认证失败 %s", true, responseBody)
		return nil
	}

	if response.StatusCode == 404 {
		Log("Submit", "请求时发生了错误: 资源不存在", true)
		return nil
	}

	return responseBody
}
