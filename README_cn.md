# ZBProxy
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_small)
[![Go Reference](https://pkg.go.dev/badge/github.com/layou233/ZBProxy.svg)](https://pkg.go.dev/github.com/layou233/ZBProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/layou233/ZBProxy)](https://goreportcard.com/report/github.com/layou233/ZBProxy)

[**English**](README.md) | **简体中文**

**新闻：ZBProxy-3.0版本已经推出，请前往[**Actions**](https://github.com/layou233/ZBProxy/actions)下载最新版本**

🚀快速搭建Minecraft服务器加速IP，给您最好的体验.
使用go语言编写，支持多平台.
一键搭建Minecraft加速IP软件，作者[B站@贴吧蜡油](https://space.bilibili.com/404017926 "点我前往空间").

#### **[加入tg群](https://t.me/launium)** 
#### **[文档](https://launium.com/doc/ZBProxy)**

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

## 赞助商
[![JetBrains logo](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=ZBProxy)  
JetBrains for open source support development license.

## 许可证
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_large)