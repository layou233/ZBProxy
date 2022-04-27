package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

var Config configMain

func LoadConfig() {
	configFile, err := os.ReadFile("ZBProxy.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Configuration file is not exists. Generating a new one...")
			generateDefaultConfig()
			goto success
		} else {
			log.Panic(color.HiRedString("Unexpected error when loading config: ", err.Error()))
		}
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Panic(color.HiRedString("Config format error: ", err.Error()))
	}

success:
	log.Println(color.HiYellowString("Successfully loaded config from file."))
}

func generateDefaultConfig() {
	file, err := os.Create("ZBProxy.json")
	defer file.Close()
	if err != nil {
		log.Panic("Failed to create configuration file:", err.Error())
	}
	Config = configMain{
		Services: []*ConfigProxyService{
			{
				Name:                  "HypixelDefault",
				TargetAddress:         "mc.hypixel.net",
				TargetPort:            25565,
				Listen:                25565,
				Flow:                  "auto",
				EnableHostnameRewrite: true,
				MotdFavicon:           "{DEFAULT_MOTD}",
				MotdDescription:       "§d{NAME}§e service is working on §a§o{INFO}§r\n§c§lProxy for §6§n{HOST}:{PORT}§r",
			},
		},
	}
	newConfig, _ :=
		json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	if err != nil {
		log.Panic("Failed to save configuration file:", err.Error())
	}
}
