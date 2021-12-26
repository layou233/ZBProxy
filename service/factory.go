package service

import (
	"ZBProxy/config"
	"ZBProxy/service/minecraft"
	"ZBProxy/service/transfer"
	"fmt"
	mcnet "github.com/Tnze/go-mc/net"
	"github.com/fatih/color"
	"log"
	"net"
	"sync"
)

func StartNewService(s *config.ConfigProxyService, wg *sync.WaitGroup) {
	defer wg.Done()
	var isMinecraftHandleNeeded = s.EnableHostnameRewrite ||
		s.EnableAnyDest ||
		s.EnableWhiteList ||
		s.EnableMojangCapeRequirement ||
		s.MotdDescription != "" ||
		s.MotdFavicon != ""
	var flowType = getFlowType(s.Flow)
	if flowType == -1 {
		log.Panic(color.HiRedString("Service %s: Unknown flow type '%s'.", s.Name, s.Flow))
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Listen))
	if err != nil {
		log.Panic(color.HiRedString("Service %s: Can't start listening on port %v: %v", s.Name, s.Listen, err.Error()))
	}

	for {
		conn, err := listen.Accept()
		if err == nil {
			var remote *mcnet.Conn = nil
			if isMinecraftHandleNeeded {
				remote, err = minecraft.NewConnHandler(s, &conn)
			}
			if err != nil {
				continue
			}
			if remote == nil {
				remote, err = mcnet.DialMC(fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
				if err != nil {
					log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
					continue
				}
			}
			transfer.SimpleTransfer(conn, remote.Socket, flowType)
		}
	}
}

func getFlowType(flow string) int {
	switch flow {
	case "origin":
		return transfer.FLOW_ORIGIN
	case "linux-zerocopy":
		return transfer.FLOW_LINUX_ZEROCOPY
	case "auto":
		return transfer.FLOW_AUTO
	default:
		return -1
	}
}
