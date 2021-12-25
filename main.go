package main

import (
	"ZBProxy/config"
	"ZBProxy/console"
	"ZBProxy/service"
	"fmt"
	"github.com/fatih/color"
	"log"
)

var onlineConnections = 0

const (
	ServerAddr             = "mc.hypixel.net"
	ServerPort      uint16 = 25565 // this must be uint16 (unsigned short) to be compatible with the protocol
	LocalPort       uint16 = 25565
	MotdDescription        = ""
)

func main() {
	log.SetOutput(color.Output)
	console.SetTitle(fmt.Sprintf("ZBProxy %v | Loading...", Version))
	console.Println(color.HiRedString(` ______  _____   _____   _____    _____  __    __ __    __
|___  / |  _  \ |  _  \ |  _  \  /  _  \ \ \  / / \ \  / /
   / /  | |_| | | |_| | | |_| |  | | | |  \ \/ /   \ \/ /`), color.HiWhiteString(`
  / /   |  _  { |  ___/ |  _  /  | | | |   }  {     \  /
 / /__  | |_| | | |     | | \ \  | |_| |  / /\ \    / /
/_____| |_____/ |_|     |_|  \_\ \_____/ /_/  \_\  /_/`))
	color.HiGreen("Welcome to ZBProxy %s!\n\n", Version)
	go CheckUpdate()

	config.LoadConfig()

	for i := 0; i < len(config.Config.Services); i++ {
		go service.StartNewService(&config.Config.Services[i])
	}
}
