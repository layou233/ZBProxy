package minecraft

import (
	"errors"
	"log"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/layou233/ZBProxy/common"
	"github.com/layou233/ZBProxy/common/buf"
	"github.com/layou233/ZBProxy/common/mcprotocol"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/service/access"
	"github.com/layou233/ZBProxy/service/transfer"
)

var (
	// ErrSuccessfullyHandledMOTDRequest means the Minecraft client requested for MOTD
	// and has been correctly handled by program. This used to skip the data forward
	// process and directly go to the end of this connection.
	ErrSuccessfullyHandledMOTDRequest         = errors.New("handled MOTD")
	ErrRejectedLoginAccessControl             = errors.New("interrupted by access control")
	ErrRejectedLoginPlayerNumberLimitExceeded = errors.New("rejected due to player number limit exceeded")
	ErrBadPlayerName                          = errors.New("rejected due to bad player name")
)

func NewConnHandler(s *config.ConfigProxyService,
	ctx *transfer.ConnContext,
	c net.Conn,
	options *transfer.Options,
) (net.Conn, error) {
	buffer := buf.NewSize(256)
	defer buffer.Release()
	buffer.Reset(mcprotocol.MaxVarIntLen)

	conn := mcprotocol.StreamConn(c)
	err := conn.ReadLimitedPacket(buffer, 250)
	if err != nil {
		return nil, err
	}

	var packetID mcprotocol.VarInt
	// Server bound : Handshake
	var (
		protocol  mcprotocol.VarInt
		hostname  string
		port      uint16
		nextState byte
	)
	err = mcprotocol.Scan(buffer, &packetID, &protocol, &hostname, &port, &nextState)
	if err != nil {
		return nil, err
	}
	if s.Minecraft.HostnameAccess.Mode != access.DefaultMode {
		hit := false
		hostnameClean := strings.TrimSuffix(hostname, "\x00FML\x00")
		for _, list := range s.Minecraft.HostnameAccess.ListTags {
			if hit = common.Must(access.GetTargetList(list)).Has(hostnameClean); hit {
				break
			}
		}
		switch s.Minecraft.HostnameAccess.Mode {
		case access.AllowMode:
			if !hit {
				c.(*net.TCPConn).SetLinger(0)
				c.Close()
				return nil, ErrRejectedLoginAccessControl
			}
		case access.BlockMode:
			if hit {
				c.(*net.TCPConn).SetLinger(0)
				c.Close()
				return nil, ErrRejectedLoginAccessControl
			}
		}
	}
	if nextState == 1 { // status
		if s.Minecraft.MotdDescription == "" && s.Minecraft.MotdFavicon == "" {
			// directly proxy MOTD from server
			remote, err := options.Out.Dial("tcp", net.JoinHostPort(s.TargetAddress, strconv.FormatInt(int64(s.TargetPort), 10)))
			if err != nil {
				return nil, err
			}

			buffer.Rewind(mcprotocol.MaxVarIntLen)
			err = mcprotocol.StreamConn(remote).WritePacket(buffer) // Server bound : Handshake
			if err != nil {
				return nil, err
			}

			return remote, nil
		} else {
			// Server bound : Status Request
			// Must read, but not used (and also nothing included in it)
			//buffer.Reset(mcprotocol.MaxVarIntLen)
			err = conn.ReadLimitedPacket(buffer, 1)
			if err != nil {
				return nil, err
			}

			// send custom MOTD
			motd := generateMOTD(int(protocol), s, options)
			motdLen := len(motd)

			buffer.Reset(mcprotocol.MaxVarIntLen)
			common.Must0(mcprotocol.WriteToPacket(buffer,
				byte(0x00), // Client bound : Status Response
				mcprotocol.VarInt(motdLen),
			))
			err = conn.WriteVectorizedPacket(buffer, motd)
			if err != nil {
				return nil, err
			}

			// handle ping request
			buffer.Reset(mcprotocol.MaxVarIntLen)
			switch s.Minecraft.PingMode {
			case pingModeDisconnect:
			case pingMode0ms:
				err = mcprotocol.WriteToPacket(buffer,
					byte(0x01),           // Client bound : Ping Response
					int64(math.MaxInt64), // this makes no sense but only a number
				)
				if err != nil {
					return nil, err
				}
				err = conn.WritePacket(buffer)
				if err != nil {
					return nil, err
				}
			default:
				err = conn.ReadLimitedPacket(buffer, 9)
				if err != nil {
					return nil, err
				}
				err = conn.WritePacket(buffer)
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
	// We only read its packet length and the player name, ignoring the rest part.
	// Unread part would be sent to target during the copy stage.
	// The reason for doing this is that this packet format has been modified many times in the history,
	// so it would take a lot of code to make it all compatible. So why not just forward it?
	// Get player name and check the profile
	buffer.Reset(mcprotocol.MaxVarIntLen)
	loginStartLen, _, err := mcprotocol.ReadVarIntFrom(c)
	if err != nil {
		return nil, err
	}
	_, _, err = mcprotocol.ReadVarIntFrom(c) // skip packet ID
	if err != nil {
		return nil, err
	}
	var playerName string
	{
		playerNameLen, _, err := mcprotocol.ReadVarIntFrom(c)
		if err != nil {
			return nil, err
		}
		if playerNameLen > 16 || playerNameLen <= 0 {
			return nil, ErrBadPlayerName
		}
		_, err = buffer.ReadFullFrom(c, int(playerNameLen))
		if err != nil {
			return nil, err
		}
		playerName = string(buffer.Bytes())
	}

	if s.Minecraft.OnlineCount.EnableMaxLimit && s.Minecraft.OnlineCount.Max <= int(options.OnlineCount.Load()) {
		log.Printf("Service %s : %s Rejected a new Minecraft player login request due to online player number limit: %s", s.Name, ctx.ColoredID, playerName)
		msg, err := generatePlayerNumberLimitExceededMessage(s, playerName).MarshalJSON()
		if err != nil {
			return nil, err
		}

		buffer.Reset(mcprotocol.MaxVarIntLen)
		common.Must0(mcprotocol.WriteToPacket(buffer,
			byte(0x00), // Client bound : Disconnect (login)
			mcprotocol.VarInt(len(msg)),
		))
		err = conn.WriteVectorizedPacket(buffer, msg)
		if err != nil {
			return nil, err
		}

		c.(*net.TCPConn).SetLinger(10) //nolint:errcheck
		c.Close()
		return nil, ErrRejectedLoginPlayerNumberLimitExceeded
	}

	accessibility := "DEFAULT"
	if s.Minecraft.NameAccess.Mode != access.DefaultMode {
		hit := false
		for _, list := range s.Minecraft.NameAccess.ListTags {
			if hit = common.Must(access.GetTargetList(list)).Has(playerName); hit {
				break
			}
		}
		switch s.Minecraft.NameAccess.Mode {
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
	log.Printf("Service %s : %s New Minecraft player logged in: %s [%s]", s.Name, ctx.ColoredID, playerName, accessibility)
	ctx.AttachInfo("PlayerName=" + playerName)
	if accessibility == "DENY" || accessibility == "REJECT" {
		msg, err := generateKickMessage(s, playerName).MarshalJSON()
		if err != nil {
			return nil, err
		}

		buffer.Reset(mcprotocol.MaxVarIntLen)
		common.Must0(mcprotocol.WriteToPacket(buffer,
			byte(0x00), // Client bound : Disconnect (login)
			mcprotocol.VarInt(len(msg)),
		))
		err = conn.WriteVectorizedPacket(buffer, msg)
		if err != nil {
			return nil, err
		}

		c.(*net.TCPConn).SetLinger(10) //nolint:errcheck
		c.Close()
		return nil, ErrRejectedLoginAccessControl
	}

	remote, err := options.Out.Dial("tcp", net.JoinHostPort(s.TargetAddress, strconv.FormatInt(int64(s.TargetPort), 10)))
	if err != nil {
		conn.Close()
		return nil, common.Cause("failed to dial to target server: ", err)
	}
	remoteMC := mcprotocol.StreamConn(remote)

	// Hostname rewritten
	buffer.Reset(mcprotocol.MaxVarIntLen)
	if s.Minecraft.EnableHostnameRewrite {
		err = mcprotocol.WriteToPacket(buffer,
			byte(0x00), // Server bound : Handshake
			protocol,
			func() string {
				if !s.Minecraft.IgnoreFMLSuffix &&
					strings.HasSuffix(hostname, "\x00FML\x00") {
					return s.Minecraft.RewrittenHostname + "\x00FML\x00"
				}
				return s.Minecraft.RewrittenHostname
			}(),
			s.TargetPort,
			byte(2),
		)
	} else {
		err = mcprotocol.WriteToPacket(buffer,
			byte(0x00), // Server bound : Handshake
			protocol,
			hostname,
			port,
			byte(2),
		)
	}
	if err != nil {
		return nil, err
	}
	err = remoteMC.WritePacket(buffer)
	if err != nil {
		return nil, err
	}

	// Server bound : Login Start
	buffer.Reset(mcprotocol.MaxVarIntLen)
	err = mcprotocol.WriteToPacket(buffer,
		byte(0x00),
		playerName,
	)
	if err != nil {
		return nil, err
	}
	mcprotocol.AppendPacketLength(buffer, int(loginStartLen))
	_, err = remote.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return remote, nil
}
