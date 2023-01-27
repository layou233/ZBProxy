package outbound

import (
	"io"
	"net"
	"syscall"
)

var SystemOutbound Outbound = &systemOutbound{}

func NewSystemOutbound(control DialerControl) Outbound {
	return &systemOutbound{
		Dialer: net.Dialer{Control: control},
	}
}

type DialerControl = func(network string, address string, c syscall.RawConn) error

type systemOutbound struct {
	net.Dialer
}

func (o *systemOutbound) DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	return net.DialTCP(network, laddr, raddr)
}

func (o *systemOutbound) Handshake(_ io.Reader, _ io.Writer, _, _ string) error {
	return nil
}
