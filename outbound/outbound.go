// Package outbound contains all protocols for remote connection used by ZBProxy.
//
// To implement an outbound protocol, one needs to do the following:
// 1. Implement the interface(s) below.
package outbound

import (
	"io"
	"net"
)

type Outbound interface {
	Dial(network, address string) (net.Conn, error)
	DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error)
	Handshake(r io.Reader, w io.Writer, network, address string) error
}
