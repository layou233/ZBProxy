package service

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/minecraft"
	"github.com/layou233/ZBProxy/service/transfer"
	"github.com/layou233/ZBProxy/version"
	"log"
	"net"
	"strconv"
	"strings"
)

var ListenerArray = make([]net.Listener, 1)

func StartNewService(s *config.ConfigProxyService) {
	// Check Settings
	var isMinecraftHandleNeeded = s.EnableHostnameRewrite ||
		s.EnableAnyDest ||
		s.EnableWhiteList ||
		s.EnableMojangCapeRequirement ||
		s.MotdDescription != "" ||
		s.MotdFavicon != ""
	flowType := getFlowType(s.Flow)
	if flowType == -1 {
		log.Panic(color.HiRedString("Service %s: Unknown flow type '%s'.", s.Name, s.Flow))
	}
	if s.MotdFavicon == "{DEFAULT_MOTD}" {
		s.MotdFavicon = minecraft.DefaultMotd
	}
	s.MotdDescription = strings.NewReplacer(
		"{INFO}", "ZBProxy "+version.Version,
		"{NAME}", s.Name,
		"{HOST}", s.TargetAddress,
		"{PORT}", strconv.Itoa(int(s.TargetPort)),
	).Replace(s.MotdDescription)
	if s.EnableHostnameRewrite && s.RewrittenHostname == "" {
		s.RewrittenHostname = s.TargetAddress
	}
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   nil, // listens on all available IP addresses of the local system
		Port: int(s.Listen),
	})
	if err != nil {
		log.Panic(color.HiRedString("Service %s: Can't start listening on port %v: %v", s.Name, s.Listen, err.Error()))
	}
	ListenerArray = append(ListenerArray, listen) // add to ListenerArray
	remoteAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
	for {
		conn, err := listen.AcceptTCP()
		if err == nil {
			go newConnReceiver(s, conn, isMinecraftHandleNeeded, flowType, remoteAddr)
		}
	}
}

func getFlowType(flow string) int {
	switch flow {
	case "origin":
		return transfer.FLOW_ORIGIN
	case "linux-zerocopy":
		return transfer.FLOW_LINUX_ZEROCOPY
	case "zerocopy":
		return transfer.FLOW_ZEROCOPY
	case "multiple":
		return transfer.FLOW_MULTIPLE
	case "auto":
		return transfer.FLOW_AUTO
	default:
		return -1
	}
}
