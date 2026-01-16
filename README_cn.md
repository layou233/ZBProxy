# ZBProxy
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_small)
[![Go Reference](https://pkg.go.dev/badge/github.com/layou233/zbproxy/v3.svg)](https://pkg.go.dev/github.com/layou233/zbproxy/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/layou233/zbproxy/v3)](https://goreportcard.com/report/github.com/layou233/zbproxy/v3)

[**English**](README.md) | **简体中文**

🚀 一个简单、快速、高性能的多用途 TCP 中继，主要为搭建 Hypixel 加速 IP 而开发。

一键搭建Minecraft加速IP软件，作者[B站@贴吧蜡油](https://space.bilibili.com/404017926 "点我前往空间")。

## Feature Highlights

- [x] ☝ 一键部署
- [x] 📋 高可自定义的配置
- [x] 🔌 在 Linux 上使用 `splice(2)` 进行零拷贝转发, 以及其它两种转发模式
- [x] 👮 在 IP 和 Minecraft 玩家名 上启用黑/白名单 (访问控制)
- [x] 🔄 配置文件热重载 列表 和 Minecraft MOTD
- [x] 📦 定制的轻量高性能 Minecraft 网络协议框架
- [x] 💻 干净且多彩的日志输出，易于跟踪每一个连接
- [x] 🔮 多平台和 CPU 架构支持
- 以及更多...

#### **[加入 Telegram 群](https://t.me/launium)** 
#### **[文档 (开发中)](https://launium.com/doc/ZBProxy)**

## 本程序可以做什么？
在大多数情况下，你可以使用Nginx的```proxy_pass```来中转Minecraft数据。 
完整代码如下:
```
stream {
    server {
        listen 25565;
        proxy_pass TARGET_SERVER_ADDRESS;
    }
}
```
但从2020年开始，Hypixel会验证玩家的登录地址.
如果你没有从Hypixel官方地址```mc.hypixel.net:25565```登录, 你将无法加入游戏.
最初的方法是通过修改```hosts```文件来欺骗服务器.  
但这对于很多玩家来说太复杂了. 
我们研究了它的工作原理, 在技术层面通过修改客户端发送的数据, 成功地绕过了检测.
这项研究的成果就是你现在看到的 ZBProxy.  
对于玩家来说,**直接**输入代理服务器地址便可以加入游戏.

**在最新版本，你甚至可以修改加速IP的图标和MOTD**

## 它安全吗?
完全不需要担心隐私问题，我们的代码是完全开源的，所以你可以自由检查是否有后门。

## 如何使用？
完整的文档已迁移至
https://launium.com/doc/ZBProxy

## 鸣谢
[![JetBrains logo.](https://resources.jetbrains.com/storage/products/company/brand/logos/jetbrains.svg)](https://jb.gg/OpenSource)  
使用 JetBrains IDE 开发。

## 许可证
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_large)