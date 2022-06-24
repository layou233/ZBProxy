package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/console"
	"github.com/layou233/ZBProxy/service"
	"github.com/layou233/ZBProxy/version"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	color.HiGreen("Welcome to ZBProxy %s!\n", version.Version)
	color.HiBlack("Build Information: %s, %s-%s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
	go version.CheckUpdate()

	config.LoadConfig()

	var accessLists []*config.AccessLists
	num := 0
	accessLists = make([]*config.AccessLists, 1024)
	for _, s := range config.Config.Services {
		accessLists[num] = &config.AccessLists{IpAccessLists: nil, McNameAccessLists: nil}
		go service.StartNewService(s, accessLists[num])
		num += 1
	}

	osSignals := make(chan os.Signal, 1)

	for {
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
		// capture signals to execute hot reload
		if <-osSignals == syscall.SIGHUP {
			// reload whitelist or blacklist
			config.LoadConfig()
			numNew := 0
			var accessListsNew []*config.AccessLists
			for _, s := range config.Config.Services {
				accessLists[numNew] = &config.AccessLists{IpAccessLists: nil, McNameAccessLists: nil}
				service.ParseAccessLists(s, accessListsNew[numNew], true)
				if accessLists[numNew].IPAccessMode != accessListsNew[numNew].IPAccessMode {
					// ipAccessMode can't be changed, need check
					log.Println(color.HiRedString("IPAccessMode can't be changed, please check and reload"))
					goto End
				}
				numNew += 1
			}
			if numNew != num {
				log.Println(color.HiRedString("service configuration number was changed, please check and reload"))
				goto End
			}
			accessLists = accessListsNew
		End:
		} else {
			// sometimes after the program exits on Windows, the ports are still occupied and "listening".
			// so manually closes these listeners when the program exits.
			for _, listener := range service.ListenerArray {
				if listener != nil { // avoid null pointers
					listener.Close()
				}
			}
			break
		}
	}
}
