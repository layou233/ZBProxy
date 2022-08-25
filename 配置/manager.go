package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/layou233/ZBProxy/common/set"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	Config     configMain
	Lists      map[string]*set.StringSet
	reloadLock sync.Mutex
)

func LoadConfig() {
	configFile, err := os.ReadFile("ZBProxy.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Configuration file is not exists. Generating a new one...")
			generateDefaultConfig()
			goto success
		} else {
			log.Panic(color.HiRedString("Unexpected error when loading config: %s", err.Error()))
		}
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Panic(color.HiRedString("Config format error: %s", err.Error()))
	}

success:
	LoadLists(false)
	log.Println(color.HiYellowString("Successfully loaded config from file."))
}

func generateDefaultConfig() {
	file, err := os.Create("ZBProxy.json")
	if err != nil {
		log.Panic("Failed to create configuration file:", err.Error())
	}
	Config = configMain{
		Services: []*ConfigProxyService{
			{
				Name:          "HypixelDefault",
				TargetAddress: "mc.hypixel.net",
				TargetPort:    25565,
				Listen:        25565,
				Flow:          "auto",
				Minecraft: minecraft{
					EnableHostnameRewrite: true,
					IgnoreFMLSuffix:       true,
					OnlineCount: onlineCount{
						Max:            114514,
						Online:         -1,
						EnableMaxLimit: false,
					},
					MotdFavicon:     "{DEFAULT_MOTD}",
					MotdDescription: "§d{NAME}§e service is working on §a§o{INFO}§r\n§c§lProxy for §6§n{HOST}:{PORT}§r",
				},
			},
		},
		Lists: map[string][]string{
			//"test": {"foo", "bar"},
		},
	}
	newConfig, _ :=
		json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	file.Close()
	if err != nil {
		log.Panic("Failed to save configuration file:", err.Error())
	}
}

func LoadLists(isReload bool) bool {
	reloadLock.Lock()
	if isReload {
		configFile, err := os.ReadFile("ZBProxy.json")
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(color.HiRedString("Fail to reload : Configuration file is not exists."))
			} else {
				log.Println(color.HiRedString("Unexpected error when reloading config: %s", err.Error()))
			}
			reloadLock.Unlock()
			return false
		}

		err = json.Unmarshal(configFile, &Config)
		if err != nil {
			log.Println(color.HiRedString("Fail to reload : Config format error: %s", err.Error()))
			reloadLock.Unlock()
			return false
		}
	}
	//log.Println("Lists:", Config.Lists)
	if l := len(Config.Lists); l == 0 { // if nothing in Lists
		Lists = map[string]*set.StringSet{} // empty map
	} else {
		Lists = make(map[string]*set.StringSet, l) // map size init
		for k, v := range Config.Lists {
			//log.Println("List: Loading", k, "value:", v)
			set := set.NewStringSetFromSlice(v)
			Lists[k] = &set
		}
	}
	Config.Lists = nil // free memory
	reloadLock.Unlock()
	runtime.GC()
	return true
}

func MonitorConfig(watcher *fsnotify.Watcher) error {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write { // config reload
					log.Println(color.HiMagentaString("Config Reload : file change detected. Reloading..."))
					if LoadLists(true) { // reload success
						log.Println(color.HiMagentaString("Config Reload : Successfully reloaded Lists."))
					} else {
						log.Println(color.HiMagentaString("Config Reload : Failed to reload Lists."))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				log.Println(color.HiRedString("Config Reload Error : ", err))
			}
		}
	}()

	return watcher.Add("config.json")
}
