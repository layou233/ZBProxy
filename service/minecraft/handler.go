package minecraft

import (
	"ZBProxy/config"
	"fmt"
	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
	"net"
)

func NewConnHandler(s *config.ConfigProxyService, c *net.Conn) (*mcnet.Conn, error) {
	conn := mcnet.WrapConn(*c)
	var p packet.Packet
	err := conn.ReadPacket(&p)
	if err != nil {
		return nil, err
	}

	var ( // Server bound : Handshake
		protocol  packet.VarInt
		hostname  packet.String
		port      packet.UnsignedShort
		nextState packet.VarInt
	)
	err = p.Scan(&protocol, &hostname, &port, &nextState)
	if err != nil {
		return nil, err
	}
	if nextState == 1 { // status
		if s.MotdDescription != "" ||
			s.MotdFavicon != "" {
			// Server bound : Status Request
			// Must read, but not used (and also nothing included in it)
			conn.ReadPacket(&p)

			// send custom MOTD
			conn.WritePacket(generateMotdPacket(
				int(protocol),
				s.MotdFavicon, s.MotdDescription))

			// handle for ping request
			conn.ReadPacket(&p)
			conn.WritePacket(p)

			conn.Close()
			return conn, nil
		} else {
			// directly proxy MOTD from server
			remote, err := mcnet.DialMC(fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
			if err != nil {
				return nil, err
			}

			remote.WritePacket(p)      // Server bound : Handshake
			remote.Write([]byte{1, 0}) // Server bound : Status Request
			return remote, nil
		}
	}
	// else: login

	// Server bound : Login Start
	// Get player name and check the profile
	conn.ReadPacket(&p)
	var (
		playerName packet.String
	)
	p.Scan(&playerName)
	// TODO PlayerName handle

	remote, err := mcnet.DialMC(fmt.Sprintf("%v:%v", s.TargetAddress, s.TargetPort))
	if err != nil {
		return nil, err
	}

	// Hostname rewritten
	if s.EnableHostnameRewrite {
		if s.RewrittenHostname == "" {
			s.RewrittenHostname = s.TargetAddress
		}
		remote.WritePacket(packet.Marshal(
			0x0, // Server bound : Handshake
			packet.String(s.RewrittenHostname),
			packet.UnsignedShort(s.TargetPort),
			packet.VarInt(2),
		))
	} else {
		remote.WritePacket(packet.Marshal(
			0x0, // Server bound : Handshake
			hostname,
			port,
			packet.VarInt(2),
		))
	}
	remote.WritePacket(p)
	return remote, nil
}
