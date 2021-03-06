package service

import (
	"fmt"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/outbound"
	"github.com/layou233/ZBProxy/service/minecraft"
	"github.com/layou233/ZBProxy/service/transfer"
	"log"
	"net"
)

func newConnReceiver(s *config.ConfigProxyService,
	conn *net.TCPConn,
	out outbound.Outbound,
	isMinecraftHandleNeeded bool,
	flowType int,
	mcNameMode int) {

	log.Println("Service", s.Name, ": A new connection request sent by", conn.RemoteAddr().String(), "is received.")
	defer log.Println("Service", s.Name, ": A connection with", conn.RemoteAddr().String(), "is closed.")
	var err error // in order to avoid scoop problems
	var remote net.Conn = nil

	if isMinecraftHandleNeeded {
		remote, err = minecraft.NewConnHandler(s, conn, out, mcNameMode)
		if err != nil {
			return
		}
	}

	if remote == nil {
		remote, err = out.Dial("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	transfer.SimpleTransfer(conn, remote, flowType)
}
