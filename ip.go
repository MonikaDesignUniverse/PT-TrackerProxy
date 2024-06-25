package main

var maxFailedCount = 3

var ipv4FailedCount = 0
var ipv4APIPos = 0
var ipv4DetectAPIs = []string { "https://api-ipv4.ip.sb/ip", "https://ipv4.ip.mir6.com" }

var ipv6FailedCount = 0
var ipv6APIPos = 0
var ipv6DetectAPIs = []string { "https://api-ipv6.ip.sb/ip", "https://ipv6.ip.mir6.com" }

func RefreshCurrentIPv4() bool {
	ipv4ResponseBody := Fetch(ipv4DetectAPIs[ipv4APIPos], true)
	if ipv4ResponseBody == nil {
		Log("RefreshCurrentIPv4", "获取 IPv4 地址时发生了错误 (Error 1)")
		return false
	}

	ipv4ResponseBodyStr := StrTrim(string(ipv4ResponseBody))
	if ipv4ResponseBodyStr == "ERR_NOROUTE" {
		Log("Debug-RefreshCurrentIPv4", "获取 IPv4 地址失败.")
		return true
	}
	if CheckPublicIP(ipv4ResponseBodyStr) != 4 {
		Log("RefreshCurrentIPv4", "获取 IPv4 地址时发生了错误 (Error 2)")
		return false
	}

	if currentIPv4 != ipv4ResponseBodyStr {
		currentIPv4 = ipv4ResponseBodyStr
		Log("RefreshCurrentIPv4", "获取当前 IPv4 地址: %s", ipv4ResponseBodyStr)
	}

	return true
}
func RefreshCurrentIPv6() bool {
	ipv6ResponseBody := Fetch(ipv6DetectAPIs[ipv6APIPos], true)
	if ipv6ResponseBody == nil {
		Log("RefreshCurrentIPv6", "获取 IPv6 地址时发生了错误 (Error 1)")
		return false
	}

	ipv6ResponseBodyStr := StrTrim(string(ipv6ResponseBody))
	if ipv6ResponseBodyStr == "ERR_NOROUTE" {
		Log("Debug-RefreshCurrentIPv6", "获取 IPv6 地址失败.")
		return true
	}
	if CheckPublicIP(ipv6ResponseBodyStr) != 6 {
		Log("RefreshCurrentIPv6", "获取 IPv6 地址时发生了错误 (Error 2)")
		return false
	}

	if currentIPv6 != ipv6ResponseBodyStr {
		currentIPv6 = ipv6ResponseBodyStr
		Log("RefreshCurrentIPv6", "获取当前 IPv6 地址: %s", ipv6ResponseBodyStr)
	}

	return true
}
func RefreshCurrentIP() {
	if !RefreshCurrentIPv4() {
		ipv4FailedCount++

		if ipv4FailedCount >= maxFailedCount {
			ipv4FailedCount = 0
			
			if (len(ipv4DetectAPIs) - 1) > ipv4APIPos {
				ipv4APIPos++
			} else {
				ipv4APIPos = 0
			}

			Log("RefreshCurrentIP", "获取 IPv4 地址错误次数过多, 已更换 %d 号 API", ipv4APIPos)
		}
	}

	if !RefreshCurrentIPv6() {
		ipv6FailedCount++

		if ipv6FailedCount >= maxFailedCount {
			ipv6FailedCount = 0

			if (len(ipv6DetectAPIs) - 1) > ipv6APIPos {
				ipv6APIPos++
			} else {
				ipv6APIPos = 0
			}

			Log("RefreshCurrentIP", "获取 IPv6 地址错误次数过多, 已更换 %d 号 API", ipv6APIPos)
		}
	}
}
