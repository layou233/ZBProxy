package service

import (
	"ZBProxy/config"
	"ZBProxy/service/minecraft"
	"ZBProxy/service/transfer"
	"log"
	"net"
)

func newConnReceiver(s *config.ConfigProxyService,
	conn *net.TCPConn,
	isMinecraftHandleNeeded bool,
	flowType int,
	remoteAddr *net.TCPAddr) {

	log.Println("Service", s.Name, ": A new connection request sent by", conn.RemoteAddr().String(), "is received.")
	defer log.Println("Service", s.Name, ": A connection with", conn.RemoteAddr().String(), "is closed.")
	var err error // in order to avoid scoop problems
	var remote *net.TCPConn = nil

	if isMinecraftHandleNeeded {
		remote, err = minecraft.NewConnHandler(s, conn, remoteAddr)
		if err != nil {
			return
		}
	}
	log.Println(remote)

	if remote == nil {
		remote, err = net.DialTCP("tcp", nil, remoteAddr)
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	transfer.SimpleTransfer(conn, remote, flowType)
}
