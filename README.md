# PT-TrackerProxy
一个简单的 PT Tracker 本地代理, 目前仅适用于 MDU, 且正处于公测阶段.

它通过本机反向代理 Tracker 方式, 使 Tracker 请求可以通过本工具中转, 从而达到使用系统代理的效果.

它依赖以下环境变量: ```HTTP_PROXY``` 和 ```HTTPS_PROXY```.

支持的 PT Tracker (只支持 HTTPS) 可通过 HTTP Header 来接收版本号 (```X-PTTP-Version```)/IPv4 地址 (```X-PTTP-IP4```)/IPv6 地址(```X-PTTP-IP6```).

使用方法: 始终启动此工具, 设置代理所需环境变量 (HTTP_PROXY 和 HTTPS_PROXY), 关闭客户端设置中的“服务器端请求伪造 (SSRF)”保护, 最后将原有站点 Tracker 中的 https:// 批量替换为 http://127.0.0.1:7887/;
