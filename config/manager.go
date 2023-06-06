package config

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/layou233/ZBProxy/common/set"
	"github.com/layou233/ZBProxy/version"

	"github.com/fatih/color"
	"github.com/zhangyunhao116/fastrand"
)

const DefaultMotd = `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAMAAACdt4HsAAAAAXNSR0IB2cksfwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAeZQTFRF/+IAzrcHKyodHx8fR0IZ89cC2MAFKikdw64IIyMeXFQX/eAAxa8IIiIfrpwLISAfgHQRnYwOuKQKaF4VJyYecGUUj4AQ99sBUUsYMS4d38YE+98AWVEXrJkM7dMCQDsaPzsb9NgCTkgYuqYJ3sYEUEoYODQcJiUe0LkHW1MX+94B178GMC4dcmgTwq0I5s0DpJINoZAN69shjLi/i7i/ZZfEU4fHVIfHeqnB0NFOZqr/f7PWY1sVtcd7n7+g+NwBm76nv8tq7tMCRkEagLPU4Nc06M4DPjob4McE49cvgLTUxMxjnL6lpMCYi6B7s8d/ssZ/a4iJQmWPZaj7NUdcXnNvYJ7sMkVcu8pxW5XdLjtMNjMcssaBVovNKTI9qMORUYG+JCgtnr6hTHiu+N8MlbuxR26f7twci7fBQ2SPfXES5dgsgbTRPVuA29U8ebLgxMxdiHoQ0tFLb63w0dFMOFFverHe3NU6PVl9hLTN5tkpQWKLj7m78d0XRWqZmb2q++EGSnOnICIkTny0rsSHqpgMJSoxUoTCuch1KTNAV43QZHlxSXKkcY6IdrDl+N4EZqn8lLuy2tQ+t8h40tJMmLyr8N0ZtqIKerHdd7DjloYPkbq3dGkT59ooq8OMxM1ixs1gzbYH1nu7OAAAA1JJREFUeJyFl/dfE0EQxZesRkPU0CIlGCACIiAWrBhFQQELotjF3sGK2HvD3huK2P5TL7nL7bwZPnvvR/bNl83eu5k9pajyQpppylQwhKdxg0JNj3BDfhQMM2baAWpWjBMK7AYOUIUcoAutBgFQBRwQK7IZJEAVc0JJHqxH4wCYLQClZZxQXgGGRCUF6DmCkKzihOokGGpSFBCbKwi15ZxQh3FI1hOAntcgCHklnDC/EQxNIQLQzWFBKBJxWICGlggB6PqFgiDjwE57EQXwxGck4qAXT2pwATzxS1pbly5bzrRiZasn1xSnAJb4VW02rXZNiTQFYOLXWAFrPVd7igIg8eusgPU5W0cnAUDiN1gBG31fVzcB0MRb69vIVnuaCcAkfpO1fjM9rQYK8BO/xQrYinmArHiJ77UCtlkAXuL7iLbv6O+H+p27WiwAnnhHid1Qv2evDtVaADzxSu2D+v0HHM9AjwXA+8tBqD90OGvq7jCGOk7A/nIE6o8e80ypGt+RrOYE2l+OQ/0JY6pM+J4K0QAHurylk6eg/jQ1pU0DkAOz0/2JZ85C/SCahsw+xcDUqXbnz+fOQ/0FbiJPXE7U4qi6eAnqL6e5hz5x2UKHr4zQ8pGrbvsARUgkRQu9Bv/++g3ltg9UqMkH4MDU+ibU37qdNTntg6nejDwYmPoO1N+955l6mjmhrNQn0IF5H+ofPPRNDeKJV5pjMAPzEdT30mTLJ04WvYH5+AnUP1WgUQvAHZjPnkP9C8U0bAFkBubLV1D/OsoB0WILwInDG6gf7H/bx8RvQIwfxw1IveOBYoD3HwIAHx1T1wAF4I/8FFDf1pdxhUmg1DAAPgcBvmRtJFBKj1LA1yDAN9c3FjMA+lKq7wH1P3LGcQOAORG0gZ++c8IAyJz4FQT47QMa8w3AzIk/QYC/ZrPe6+tOJjInXPHEaz2uuDq6DYDOiax44h2NCUK2v+RmY5q9NfLOX/JPEDJx8IfrEFuULbS8VhCc/mKmM78ZQOKzqpIX4ln0lsZvBmHRQrMDCzVOAJEWtihbKD+pjMgymROuxsTI4yeFADonvA1ywCR3KFgmc8LVhCCIOxQu80DlEm8kvtHYepwdU42Yyfwbja+zb20v8VT4jfYfTXskk4+wbR0AAAAASUVORK5CYII=`

var (
	Config     configMain
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
		Lists: map[string]set.StringSet{
			//"test": {"foo", "bar"},
		},
	}
	newConfig, _ := json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	file.Close()
	if err != nil {
		log.Panic("Failed to save configuration file:", err.Error())
	}
}

func LoadLists(isReload bool) bool {
	reloadLock.Lock()
	defer reloadLock.Unlock()
	var config configMain
	if isReload {
		configFile, err := os.ReadFile("ZBProxy.json")
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(color.HiRedString("Fail to reload : Configuration file is not exists."))
			} else {
				log.Println(color.HiRedString("Unexpected error when reloading config: %s", err.Error()))
			}
			return false
		}

		err = json.Unmarshal(configFile, &config)
		if err != nil {
			log.Println(color.HiRedString("Fail to reload : Config format error: %s", err.Error()))
			return false
		}
	} else {
		config = Config
	}

	for _, s := range config.Services {
		if s.Minecraft.MotdFavicon == "{DEFAULT_MOTD}" {
			s.Minecraft.MotdFavicon = DefaultMotd
		}
		s.Minecraft.MotdDescription = strings.NewReplacer(
			"{INFO}", "ZBProxy "+version.Version,
			"{NAME}", s.Name,
			"{HOST}", s.TargetAddress,
			"{PORT}", strconv.Itoa(int(s.TargetPort)),
		).Replace(s.Minecraft.MotdDescription)

		if samples := s.Minecraft.OnlineCount.Sample; samples != nil {
			var convertedSamples []Sample
			switch samples := samples.(type) {
			case map[string]any:
				convertedSamples = make([]Sample, 0, len(samples))
				for uuid, name := range samples {
					convertedSamples = append(convertedSamples, Sample{
						Name: name.(string),
						ID:   uuid,
					})
				}

			case []any:
				convertedSamples = make([]Sample, 0, len(samples))
				var u [16]byte
				var dst [36]byte
				for i, sample := range samples {
					// generate random UUID with ZBProxy signature
					fastrand.Read(u[:])
					u[0] = byte(i)
					u[1] = '$'
					u[2] = 'Z'
					u[3] = 'B'
					u[4] = '$'

					// marshal UUID string
					const hexTable = "0123456789abcdef"
					dst[8] = '-'
					dst[13] = '-'
					dst[18] = '-'
					dst[23] = '-'
					for i, x := range [16]byte{
						0, 2, 4, 6,
						9, 11,
						14, 16,
						19, 21,
						24, 26, 28, 30, 32, 34,
					} {
						c := u[i]
						dst[x] = hexTable[c>>4]
						dst[x+1] = hexTable[c&0x0F]
					}

					convertedSamples = append(convertedSamples, Sample{
						Name: sample.(string),
						ID:   string(dst[:]),
					})
				}

			default:
				log.Println(color.HiMagentaString(
					"Config Reload : Failed to reload samples: unknown samples input type: %T", samples))
				return false
			}
			s.Minecraft.OnlineCount.Sample = convertedSamples
		}
	}

	Config = config
	debug.FreeOSMemory()
	return true
}
