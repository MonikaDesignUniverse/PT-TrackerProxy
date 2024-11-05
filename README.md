# PT-TrackerProxy
一个简单的 PT Tracker 本地代理, 但不仅限于 Tracker.

它通过本机反向代理, 使 RSS/Tracker 请求可以通过本工具中转, 从而达到使用系统代理的效果. 自 1.2 版本才支持 RSS 请求.

它依赖网络设置或以下环境变量 ```HTTP_PROXY```/```HTTPS_PROXY```.

使用方法:
1. 根据对应平台下载本工具;
2. 设置代理所需网络设置或环境变量;
3. 在启动 BT 客户端期间始终启动本工具;
4. 关闭客户端设置中的 "服务器端请求伪造 (SSRF) 缓解";
5. 将原有 PT 站点 RSS/Tracker 中的 https:// 批量替换为 http://127.0.0.1:7887/;

常规版本:
<details>
<summary>查看 常见平台下载版本 对照表</summary>

| 操作系统 | 处理器架构 | 处理器位数 | 下载版本      | 说明 |
| -------- | ---------- | ---------- | ------------- | ----------------- |
| macOS    | ARM64      | 64 位      | darwin-arm64  | 常见于 Apple M 系列 |
| macOS    | AMD64      | 64 位      | darwin-amd64  | 常见于 Intel 系列 |
| Windows  | AMD64      | 64 位      | windows-amd64 | 常见于大部分现代 PC |
| Windows  | i386       | 32 位      | windows-386   | 少见于部分老式 PC |
| Windows  | ARM64      | 64 位      | windows-arm64 | 常见于新型平台, 应用于部分平板/笔记本/少数特殊硬件 |
| Windows  | ARMv7      | 32 位      | windows-arm   | 少见于罕见平台, 应用于部分上古硬件, 如 Surface RT 等 |
| Linux    | AMD64      | 64 位      | linux-amd64   | 常见于大部分 NAS 及服务器 |
| Linux    | i386       | 32 位      | linux-386     | 少见于部分老式 NAS 及服务器 |
| Linux    | ARM64      | 64 位      | linux-arm64   | 常见于部分服务器及开发板, 如 Oracle 或 Raspberry Pi 等 |
| Linux    | ARMv*      | 32 位      | linux-armv*   | 少见于部分老式服务器及开发板, 查看 /proc/cpuinfo 或 从高到底试哪个能跑 |

其它版本的 Linux/NetBSD/FreeBSD/OpenBSD/Solaris 可以此类推, 并在列表中选择适合自己的.
</details>

Docker 版本: 于 [Docker Hub](https://hub.docker.com/r/monikadesignuniverse/pt-trackerproxy) 提供.

--------------------

支持的 PT 站点列表
| PT 站点 | RSS 请求 | RSS 下载 | 站点下载 | Tracker 请求 | 局域网支持 (自定义监听) |
| ------ | ------ | ------ | ------ | ------ | ------ |
| MDU | ✅ | ✅️ | ✅ | ✅ | ✅ | ✅ |

✅: PT 站点支持并允许该功能.  
❌: PT 站点不支持或不允许该功能.

RSS 请求: 可将 RSS 请求通过本工具中转.  
RSS 下载: 可将 RSS 下载通过本工具中转, 且下载 Torrent 的链接及 Tracker 支持自动适配 PTTP.  
站点下载: 可通过站点下载 Torrent,  且下载 Torrent 的 Tracker 支持根据用户设置适配 PTTP.  
Tracker 请求: 可将 Tracker 请求通过本工具中转, 且客户端 IP 地址支持由 PTTP 上报.  
局域网支持 (自定义监听): 可由 PTTP 上报自定义监听地址和监听端口以用于自动适配, 且若支持站点下载, 则用户设置也须支持自定义监听地址及监听端口.

为确保用户安全, PTTP 支持要求 PT 站点须处于上述支持列表内, 列入上述支持列表的 PT 站点同时会被列入工具白名单.

--------------------

支持的 PT 站点可通过 HTTP Header (```X-PTTP-*```) 来接收信息, 样例可见 server_PTTPHelper.php.  
不支持但兼容的 PT 站点可通过 HTTP Header (```X-Forwarded-For```) 来接收信息, 须打开 XFF 兼容模式, XFF 兼容模式下仅支持上报一个 IP 地址.  
请注意: 支持的 PT 站点不应再支持 XFF 兼容模式, 不支持但兼容的 PT 站点可通过 XFF 兼容模式来进行过渡.

自 1.0 版本:  
版本号 (```X-PTTP-Version```)  
IPv4 地址 (```X-PTTP-IP4```)  
IPv6 地址 (```X-PTTP-IP6```)

自 1.2 版本:  
监听地址 (```X-PTTP-ListenAddr```)  
监听端口 (```X-PTTP-ListenPort```)
