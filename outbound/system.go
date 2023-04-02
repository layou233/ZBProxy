package outbound

import (
	"io"
	"net"
	"syscall"
)

var SystemOutbound Outbound = &systemOutbound{}

func NewSystemOutbound(control DialerControl) Outbound {
	if control == nil {
		return SystemOutbound
	}
	return &systemOutbound{
		Dialer: net.Dialer{Control: control},
	}
}

type DialerControl = func(network string, address string, c syscall.RawConn) error

type systemOutbound struct {
	net.Dialer
}

func (o *systemOutbound) DialTCP(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	if o.Dialer.Control == nil {
		return net.DialTCP(network, laddr, raddr)
	}

	conn, err := (&net.Dialer{
		LocalAddr: laddr,
		Control:   o.Dialer.Control,
	}).Dial(network, raddr.String())
	if err == nil {
		return conn.(*net.TCPConn), nil
	}
	return nil, err
}

func (o *systemOutbound) Handshake(_ io.Reader, _ io.Writer, _, _ string) error {
	return nil
}
