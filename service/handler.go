package service

import (
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/minecraft"
	"github.com/layou233/ZBProxy/service/transfer"
	"log"
	"net"
)

func newConnReceiver(s *config.ConfigProxyService,
	conn *net.TCPConn,
	isMinecraftHandleNeeded bool,
	flowType int,
	remoteAddr *net.TCPAddr,
	//mcNameLists []*set.StringSet,
	mcNameMode int) {

	log.Println("Service", s.Name, ": A new connection request sent by", conn.RemoteAddr().String(), "is received.")
	defer log.Println("Service", s.Name, ": A connection with", conn.RemoteAddr().String(), "is closed.")
	var err error // in order to avoid scoop problems
	var remote *net.TCPConn = nil

	if isMinecraftHandleNeeded {
		if s.Debug {
			log.Println("Using Minecraft Autofix packs mode")
		}
		remote, err = minecraft.NewConnHandler(s, conn, remoteAddr, mcNameMode)
		if err != nil {
			return
		}
	}

	if remote == nil {
		if s.Debug {
			log.Println("Using normal mode")
		}
		remote, err = net.DialTCP("tcp", nil, remoteAddr)
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	transfer.SimpleTransfer(conn, remote, flowType)
}
