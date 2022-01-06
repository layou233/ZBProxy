package service

import (
	"ZBProxy/config"
	"ZBProxy/service/minecraft"
	"ZBProxy/service/transfer"
	"fmt"
	mcnet "github.com/Tnze/go-mc/net"
	"log"
	"net"
)

func newConnReceiver(s *config.ConfigProxyService,
	conn net.Conn,
	isMinecraftHandleNeeded bool,
	flowType int) {

	log.Println("Service", s.Name, ": A new connection request sent by", conn.RemoteAddr().String(), "is received.")
	defer log.Println("Service", s.Name, ": A connection with", conn.RemoteAddr().String(), "is closed.")
	var err error // in order to avoid scoop problems
	var remote *mcnet.Conn = nil

	if isMinecraftHandleNeeded {
		remote, err = minecraft.NewConnHandler(s, &conn)
	}
	if err != nil {
		return
	}

	if remote == nil {
		remote, err = mcnet.DialMC(fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	transfer.SimpleTransfer(conn, remote.Socket, flowType)
}
