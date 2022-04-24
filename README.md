# ZBProxy
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_small)
[![Go Reference](https://pkg.go.dev/badge/github.com/layou233/ZBProxy.svg)](https://pkg.go.dev/github.com/layou233/ZBProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/layou233/ZBProxy)](https://goreportcard.com/report/github.com/layou233/ZBProxy)  
**NEWS: NEW ZBPROXY 3.0 WITH CONFIGURATION SUPPORT IS COMING SOON AND ALREADY AVAILABLE IN THE ACTIONS.**  

🚀Easy proxy your data at the best privacy, giving you better experience enjoying Minecraft.  
Support multiplatform because of Golang\'s attributes.  
一键搭建Minecraft加速IP软件，作者[B站@贴吧蜡油](https://space.bilibili.com/404017926 "点我前往空间")。  
本项目遵守国际化理念，所以代码内注释以及README都将以英语编写，可借助翻译工具辅助阅读。  

#### **[Join Official Telegram Group](https://t.me/launium)**  
#### **[Document (W.I.P)](https://launium.com/doc/ZBProxy)**

## What can it do?
In many situations you can use Nginx ```proxy_pass``` to easy proxy your Minecraft data.  
The complete code is as follows:

```
stream {
    server {
        listen 25565;
        proxy_pass TARGET_SERVER_ADDRESS;
    }
}
```
But start from 2020, Hypixel set up an authentication of the player login address.  
If you do not log in from their official address as known as ```mc.hypixel.net:25565```, you will not be able to join the game.  
The original method is to cheat the server by modifying the ```hosts``` file.  
But that\'s too complicated for people who don\'t know the principle.  
We studied its working principle, and successfully bypassed the detection by modifying the data sent by client at the technical level.  
The product of the research is what you see now as ZBProxy.  
For players, just enter the address of your proxy server, you can join the game **directly** as usual.

### Is it safe?
There is no need to worry about privacy at all, because the connection to any Minecraft server which requires online verification is fully **encrypted**.  
Our code is completely open source, so you can freely check whether there is a backdoor.

## How to use it?
1. Download the compiled executable file at [releases page](https://github.com/layou233/ZBProxy/releases/ "releases").  
2. Run it, and your data proxy service is now established!  
For Linux system, you may need to give permissions to the executable file in order to solve problems that cannot run or run blocked. Just enter the following command:
```bash
chmod 777 PATH_OF_THE_FILE
```
3. Ensure the port **25565** is fully open on the server.
4. Enter your proxy server IP into your Minecraft client, and join it for game!  
    (Since the listening port is **25565**, you don\'t need to input the port number in the client, and the client will complete it automatically)  

Since ZBProxy 3.0, users are allowed to set the listening port and forwarding destination through the automatically generated JSON configuration file, including choosing whether to enable the hostname rewriting function and a series of surprise functions.  
At the first startup, **a JSON configuration file** is **automatically generated**, which contains a preset Hypixel forwarding configuration, so users can still build forwarding services in only one step like old versions, but reasonable use of the configuration file can help users explore more possibilities of ZBProxy. Including but not limited to quickly setting up an ordinary and efficient reverse proxy.  
If you are just new in here, you can view the **[ZBProxy Document](https://launium.com/doc/ZBProxy)** to learn about how to unlock the power of ZBProxy through configuration.

## Are there any other ways to improve my games?
Generally speaking, Linux-based operating environments have more room for optimization.  
ZBProxy supports **Zero Copy** technology on Linux, which can *reduce memory usage **by one time**, save **a lot of** CPU processing, and reduce **network latency***. When users set `Flow` to `auto` or `linux-zerocopy` on their `Service` configuration, this technology will be automatically adopted in due course.  
If you are running ZBProxy on a Linux-based system, you can go to **[ZBProxy Document](https://launium.com/doc/ZBProxy)** to view **tips for optimizing network settings**.

## TODO List
1. Some functions are still not implemented.

## Sponsor
[![JetBrains logo](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=ZBProxy)  
JetBrains for open source support development license.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Flayou233%2FZBProxy.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Flayou233%2FZBProxy?ref=badge_large)