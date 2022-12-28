package service

import (
	"fmt"
	"log"
	"net"

	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/minecraft"
	"github.com/layou233/ZBProxy/service/tls"
	"github.com/layou233/ZBProxy/service/transfer"

	"github.com/fatih/color"
)

var (
	GreenPlus = color.HiGreenString("[+]")
	RedMinus  = color.HiRedString("[-]")
)

func newConnReceiver(s *config.ConfigProxyService,
	conn *net.TCPConn,
	options *transfer.Options,
) {
	ctx := new(transfer.ConnContext).Init()
	log.Println("Service", s.Name, ":", ctx.ColoredID, GreenPlus, conn.RemoteAddr().String())
	defer log.Println("Service", s.Name, ":", ctx.ColoredID, RedMinus, conn.RemoteAddr().String(), ctx)
	var err error // avoid scoop problems
	var remote net.Conn

	if options.IsTLSHandleNeeded {
		remote, err = tls.NewConnHandler(s, conn, options.Out)
		if err != nil {
			conn.Close()
			return
		}
	} else if options.IsMinecraftHandleNeeded {
		remote, err = minecraft.NewConnHandler(s, ctx, conn, options)
		if err != nil {
			conn.Close()
			return
		}
	}

	if remote == nil {
		remote, err = options.Out.Dial("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
		if err != nil {
			log.Printf("Service %s: Failed to dial to target server: %v", s.Name, err.Error())
			conn.Close()
			return
		}
	}
	options.OnlineCount.Add(1)
	defer options.OnlineCount.Add(-1)
	transfer.SimpleTransfer(conn, remote, options.FlowType)
}
