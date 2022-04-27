package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/console"
	"github.com/layou233/ZBProxy/service"
	"github.com/layou233/ZBProxy/version"
)

func main() {
	log.SetOutput(color.Output)
	console.SetTitle(fmt.Sprintf("Running %s ", version.Version))
	color.HiBlack("Build Information: %s, %s-%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	sum := 0
	for sum <= 1000 {
		sum += 1
		// Load
		for _, s := range config.Config.Services {
			go service.StartNewService(s)
		}
		// Loaded times
		if sum < 1 {
			fmt.Printf("First Loaded ")
		} else {
			fmt.Printf("Loaded %d times \n", sum)
		}
		// keyboard Listener
		keysEvents, err := keyboard.GetKeys(69)
		if err != nil {
			panic(err)
		}

		for {
			event := <-keysEvents
			if event.Err != nil {
				panic(event.Err)
			}

			if event.Rune == 'r' {
				fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
				break
			}

			if event.Key == keyboard.KeyEsc {
				osSignals := make(chan os.Signal, 1)
				signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
				<-osSignals // wait for exits

				// sometimes after the program exits on Windows, the ports are still occupied and "listening".
				// so manually closes these listeners when the program exits.
				for _, listener := range service.ListenerArray {
					if listener != nil { // avoid null pointers
						listener.Close()
					}
				}
			}
		}
	}
}
