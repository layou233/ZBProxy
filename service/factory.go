package service

import (
	"github.com/fatih/color"
	"github.com/layou233/ZBProxy/common"
	"github.com/layou233/ZBProxy/common/set"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/outbound"
	"github.com/layou233/ZBProxy/outbound/socks"
	"github.com/layou233/ZBProxy/service/access"
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
	var (
		isTLSHandleNeeded = s.TLSSniffing.RejectNonTLS ||
			s.TLSSniffing.RejectIfNonMatch ||
			len(s.TLSSniffing.SNIAllowListTags) != 0
		isMinecraftHandleNeeded = s.Minecraft.EnableHostnameRewrite ||
			s.Minecraft.EnableAnyDest ||
			s.Minecraft.MotdDescription != "" ||
			s.Minecraft.MotdFavicon != ""
	)
	if isTLSHandleNeeded && isMinecraftHandleNeeded {
		log.Panic(color.HiRedString("Service %s: The current version can't handle TLS and Minecraft at the same time.", s.Name))
	}
	flowType := getFlowType(s.Flow)
	if flowType == -1 {
		log.Panic(color.HiRedString("Service %s: Unknown flow type '%s'.", s.Name, s.Flow))
	}
	if s.Minecraft.MotdFavicon == "{DEFAULT_MOTD}" {
		s.Minecraft.MotdFavicon = minecraft.DefaultMotd
	}
	s.Minecraft.MotdDescription = strings.NewReplacer(
		"{INFO}", "ZBProxy "+version.Version,
		"{NAME}", s.Name,
		"{HOST}", s.TargetAddress,
		"{PORT}", strconv.Itoa(int(s.TargetPort)),
	).Replace(s.Minecraft.MotdDescription)
	if s.Minecraft.EnableHostnameRewrite && s.Minecraft.RewrittenHostname == "" {
		s.Minecraft.RewrittenHostname = s.TargetAddress
	}
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   nil, // listens on all available IP addresses of the local system
		Port: int(s.Listen),
	})
	if err != nil {
		log.Panic(color.HiRedString("Service %s: Can't start listening on port %v: %v", s.Name, s.Listen, err.Error()))
	}
	ListenerArray = append(ListenerArray, listen) // add to ListenerArray

	// load access lists
	ipAccessMode := access.ParseAccessMode(s.IPAccess.Mode)
	if ipAccessMode != access.DefaultMode { // IP access control enabled
		if s.IPAccess.ListTags == nil {
			log.Panic(color.HiRedString("Service %s: ListTags can't be null when access control enabled.", s.Name))
		}
		for _, tag := range s.IPAccess.ListTags {
			if common.GetSecond[error](access.GetTargetList(tag)) != nil {
				log.Panic(color.HiRedString("Service %s: %s", s.Name, err.Error()))
			}
		}
	}

	// load Minecraft player name access lists
	mcNameAccessMode := access.ParseAccessMode(s.Minecraft.NameAccess.Mode)
	if isMinecraftHandleNeeded && mcNameAccessMode != access.DefaultMode { // IP access control enabled
		if s.Minecraft.NameAccess.ListTags == nil {
			log.Panic(color.HiRedString("Service %s: ListTags can't be null when access control enabled.", s.Name))
		}
		for _, tag := range s.Minecraft.NameAccess.ListTags {
			if common.GetSecond[error](access.GetTargetList(tag)) != nil {
				log.Panic(color.HiRedString("Service %s: %s", s.Name, err.Error()))
			}
		}
	}

	out := outbound.SystemOutbound
	switch s.Outbound.Type {
	case "socks", "socks5", "socks4a", "socks4":
		out = &socks.Client{
			Version: s.Outbound.Type,
			Network: s.Outbound.Network,
			Address: s.Outbound.Address,
		}
	}

	options := &transfer.Options{
		Out:                     out,
		IsTLSHandleNeeded:       isTLSHandleNeeded,
		IsMinecraftHandleNeeded: isMinecraftHandleNeeded,
		FlowType:                flowType,
		McNameMode:              mcNameAccessMode,
	}
	for {
		conn, err := listen.AcceptTCP()
		if err == nil {
			if ipAccessMode != access.DefaultMode {
				// https://stackoverflow.com/questions/29687102/how-do-i-get-a-network-clients-ip-converted-to-a-string-in-golang
				ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
				hit := false
				for _, list := range s.IPAccess.ListTags {
					if hit = common.Must[*set.StringSet](access.GetTargetList(list)).Has(ip); hit {
						break
					}
				}
				switch ipAccessMode {
				case access.AllowMode:
					if !hit {
						forciblyCloseTCP(conn)
						continue
					}
				case access.BlockMode:
					if hit {
						forciblyCloseTCP(conn)
						continue
					}
				}
			}
			go newConnReceiver(s, conn, options)
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

func forciblyCloseTCP(conn *net.TCPConn) {
	conn.SetLinger(0) // let Close send RST to forcibly close the connection
	conn.Close()      // forcibly close
}
