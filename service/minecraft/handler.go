package minecraft

import (
	"ZBProxy/config"
	mcnet "github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/net/packet"
	"net"
)

func NewConnHandler(s *config.ConfigProxyService, c *net.Conn) error {
	conn := mcnet.WrapConn(c)
	var p packet.Packet
	err := conn.ReadPacket(&p)
	if err != nil {
		return err
	}

	var (
		protocol  packet.VarInt
		hostname  packet.String
		port      packet.UnsignedShort
		nextState packet.VarInt
	)
	err := p.Scan(&protocol, &hostname, &port, &nextState)
	if err != nil {
		return err
	}
	if nextState == 1 { // status

	}

	// login
}
