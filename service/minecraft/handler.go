package minecraft

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/layou233/ZBProxy/common"
	"github.com/layou233/ZBProxy/common/set"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/access"
	"github.com/layou233/ZBProxy/service/transfer"

	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
	"github.com/fatih/color"
)

// ErrSuccessfullyHandledMOTDRequest means the Minecraft client requested for MOTD
// and has been correctly handled by program. This used to skip the data forward
// process and directly go to the end of this connection.
var ErrSuccessfullyHandledMOTDRequest = errors.New("")

var ErrRejectedLogin = ErrSuccessfullyHandledMOTDRequest // don't cry baby

func badPacketPanicRecover(s *config.ConfigProxyService) {
	// Non-Minecraft packet which uses `go-mc` packet scan method may cause panic.
	// So a panic handler is needed.
	if err := recover(); err != nil {
		log.Print(color.HiRedString("Service %s : Bad Minecraft packet was received: %v", s.Name, err))
	}
}

func NewConnHandler(s *config.ConfigProxyService,
	c net.Conn,
	options *transfer.Options,
) (net.Conn, error) {
	defer badPacketPanicRecover(s)

	conn := mcnet.WrapConn(c)
	var p packet.Packet
	err := conn.ReadPacket(&p)
	if err != nil {
		return nil, err
	}

	var ( // Server bound : Handshake
		protocol  packet.VarInt
		hostname  packet.String
		port      packet.UnsignedShort
		nextState packet.Byte
	)
	err = p.Scan(&protocol, &hostname, &port, &nextState)
	if err != nil {
		return nil, err
	}
	if nextState == 1 { // status
		if s.Minecraft.MotdDescription == "" && s.Minecraft.MotdFavicon == "" {
			// directly proxy MOTD from server

			remote, err := options.Out.Dial("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
			if err != nil {
				return nil, err
			}
			remoteMC := mcnet.WrapConn(remote)

			err = remoteMC.WritePacket(p) // Server bound : Handshake
			if err != nil {
				return nil, err
			}

			_, err = remote.Write([]byte{1, 0}) // Server bound : Status Request
			if err != nil {
				return nil, err
			}

			return remote, nil
		} else {
			// Server bound : Status Request
			// Must read, but not used (and also nothing included in it)
			err = conn.ReadPacket(&p)
			if err != nil {
				return nil, err
			}

			// send custom MOTD
			err = conn.WritePacket(generateMotdPacket(
				int(protocol),
				s, options))
			if err != nil {
				return nil, err
			}

			// handle for ping request
			switch s.Minecraft.PingMode {
			case pingModeDisconnect:
			case pingMode0ms:
				err = conn.WritePacket(
					packet.Marshal(0x01, packet.Long(1<<63-1))) // max int64 value
				if err != nil {
					return nil, err
				}
			default:
				err = conn.ReadPacket(&p)
				if err != nil {
					return nil, err
				}
				err = conn.WritePacket(p)
				if err != nil {
					return nil, err
				}
			}

			conn.Close()
			return nil, ErrSuccessfullyHandledMOTDRequest
		}
	}
	// else: login

	// Server bound : Login Start
	// Get player name and check the profile
	err = conn.ReadPacket(&p)
	if err != nil {
		return nil, err
	}
	var playerName packet.String
	err = p.Scan(&playerName)
	if err != nil {
		return nil, err
	}

	if s.Minecraft.OnlineCount.EnableMaxLimit && s.Minecraft.OnlineCount.Max <= int(options.OnlineCount.Load()) {
		log.Printf("Service %s : Rejected a new Minecraft player login request due to online player number limit: %s", s.Name, playerName)
		err := conn.WritePacket(packet.Marshal(
			0x00, // Client bound : Disconnect (login)
			generatePlayerNumberLimitExceededMessage(s, playerName),
		))
		if err != nil {
			return nil, err
		}

		c.(*net.TCPConn).SetLinger(10) //nolint:errcheck
		c.Close()
		return nil, ErrRejectedLogin
	}

	accessibility := "DEFAULT"
	if options.McNameMode != access.DefaultMode {
		hit := false
		for _, list := range s.Minecraft.NameAccess.ListTags {
			if hit = common.Must[*set.StringSet](access.GetTargetList(list)).Has(string(playerName)); hit {
				break
			}
		}
		switch options.McNameMode {
		case access.AllowMode:
			if hit {
				accessibility = "ALLOW"
			} else {
				accessibility = "DENY"
			}
		case access.BlockMode:
			if hit {
				accessibility = "REJECT"
			} else {
				accessibility = "PASS"
			}
		}
	}
	log.Printf("Service %s : A new Minecraft player requested a login: %s [%s]", s.Name, playerName, accessibility)
	if accessibility == "DENY" || accessibility == "REJECT" {
		err = conn.WritePacket(packet.Marshal(
			0x00, // Client bound : Disconnect (login)
			generateKickMessage(s, playerName),
		))
		if err != nil {
			return nil, err
		}

		c.(*net.TCPConn).SetLinger(10) //nolint:errcheck
		c.Close()
		return nil, ErrRejectedLogin
	}

	remote, err := options.Out.Dial("tcp", fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
	if err != nil {
		log.Printf("Service %s : Failed to dial to target server: %v", s.Name, err.Error())
		conn.Close()
		return nil, err
	}
	remoteMC := mcnet.WrapConn(remote)

	// Hostname rewritten
	if s.Minecraft.EnableHostnameRewrite {
		err = remoteMC.WritePacket(packet.Marshal(
			0x00, // Server bound : Handshake
			protocol,
			packet.String(func() string {
				if !s.Minecraft.IgnoreFMLSuffix &&
					strings.HasSuffix(string(hostname), "\x00FML\x00") {
					return s.Minecraft.RewrittenHostname + "\x00FML\x00"
				}
				return s.Minecraft.RewrittenHostname
			}()),
			packet.UnsignedShort(s.TargetPort),
			packet.Byte(2),
		))
	} else {
		err = remoteMC.WritePacket(packet.Marshal(
			0x00, // Server bound : Handshake
			protocol,
			hostname,
			port,
			packet.Byte(2),
		))
	}
	if err != nil {
		return nil, err
	}

	// Server bound : Login Start
	err = remoteMC.WritePacket(p)
	if err != nil {
		return nil, err
	}
	return remote, nil
}
