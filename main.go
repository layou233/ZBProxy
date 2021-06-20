package main

import (
	"ZBProxy/windows"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var onlinePlayers = 0

const (
	version = "1.0-SNAPSHOT"

	serverAddr = "mc.remiaft.com"
	serverPort = 25565
	localPort  = 25565
)

func main() {
	windows.SetTitle(fmt.Sprintf("ZBProxy %v | Loading...", version))
	fmt.Printf("Welcome to ZBProxy %s!\n\n", version)
	fmt.Println("Your current settings:")
	fmt.Println("serverAddr=" + serverAddr)
	fmt.Printf("serverPort=%d\n", serverPort)
	fmt.Printf("localPort=%d\n", localPort)

	log.Printf("Starting listening on local port %d...", localPort)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		log.Printf("Unable to listen on port %d.\n", localPort)
		log.Printf("Caution: %d\n", err.Error())
		log.Printf("The program will exit in 5 seconds.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	defer listen.Close()
	for {
		time.Sleep(time.Second)
		windows.SetTitle(
			fmt.Sprintf("ZBProxy %v | Online Connections: %v", version, onlinePlayers/2))
		fromConn, err := listen.Accept()
		if err != nil {
			continue
		}
		serverIp, err := net.ResolveIPAddr("ip4", serverAddr)
		if err != nil {
			log.Printf("Can't resolve hostname: %v", serverAddr)
			continue
		}
		go forDial(fromConn, fmt.Sprintf("%s:%d", serverIp.String(), serverPort))
	}
}
