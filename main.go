package main

import (
	"ZBProxy/windows"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var onlineConnections = 0

const (
	version = "1.0-SNAPSHOT"

	serverAddr        = "mc.hypixel.net"
	serverPort uint16 = 25565 // this must be uint16 (unsigned short) to be compatible with the protocol
	localPort  uint16 = 25565
)

func main() {
	windows.SetTitle(fmt.Sprintf("ZBProxy %v | Loading...", version))
	fmt.Printf("Welcome to ZBProxy %s!\n\n", version)
	fmt.Println("Your current settings:")
	fmt.Println("serverAddr=" + serverAddr)
	fmt.Printf("serverPort=%d\n", serverPort)
	fmt.Printf("localPort=%d\n", localPort)

	log.Printf("Starting listening on local port %d...", localPort)
	listen, err1 := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err1 != nil {
		log.Printf("Unable to listen on port %d.\n", localPort)
		log.Printf("Caution: %v\n", err1.Error())
		log.Printf("The program will exit in 5 seconds.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	defer listen.Close()
	for {
		time.Sleep(time.Second)
		windows.SetTitle(
			fmt.Sprintf("ZBProxy %v | Online Connections: %v", version, onlineConnections/2))
		fromConn, err2 := listen.Accept()
		if err2 != nil {
			continue
		}
		serverIp, err3 := net.ResolveIPAddr("ip4", serverAddr)
		if err3 != nil {
			log.Printf("Can't resolve hostname: %v", serverAddr)
			continue
		}
		go forDial(fromConn, fmt.Sprintf("%s:%d", serverIp.String(), serverPort))
	}
}
