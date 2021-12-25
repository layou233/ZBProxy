package main

import (
	"ZBProxy/config"
	"ZBProxy/console"
	"ZBProxy/service"
	"ZBProxy/version"
	"fmt"
	"github.com/fatih/color"
	"log"
	"sync"
)

func main() {
	log.SetOutput(color.Output)
	console.SetTitle(fmt.Sprintf("ZBProxy %v | Running...", version.Version))
	console.Println(color.HiRedString(` ______  _____   _____   _____    _____  __    __ __    __
|___  / |  _  \ |  _  \ |  _  \  /  _  \ \ \  / / \ \  / /
   / /  | |_| | | |_| | | |_| |  | | | |  \ \/ /   \ \/ /`), color.HiWhiteString(`
  / /   |  _  { |  ___/ |  _  /  | | | |   }  {     \  /
 / /__  | |_| | | |     | | \ \  | |_| |  / /\ \    / /
/_____| |_____/ |_|     |_|  \_\ \_____/ /_/  \_\  /_/`))
	color.HiGreen("Welcome to ZBProxy %s!\n\n", version.Version)
	go version.CheckUpdate()

	config.LoadConfig()

	group := sync.WaitGroup{}
	for i := 0; i < len(config.Config.Services); i++ {
		group.Add(1)
		go service.StartNewService(&config.Config.Services[i], &group)
	}
	(&group).Wait()
}
