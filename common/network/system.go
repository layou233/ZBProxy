package network

import (
	"context"
	"net"
	"syscall"
)

var SystemDialer Dialer = &systemOutbound{}

func NewSystemDialer(options *OutboundSocketOptions) Dialer {
	if options == nil {
		return SystemDialer
	}

	out := &systemOutbound{
		Dialer: net.Dialer{
			Control: NewDialerControlFromOptions(options),
		},
	}
	SetDialerTCPKeepAlive(&out.Dialer, options.KeepAliveConfig())
	if options.SendThrough != "" {
		out.Dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP(options.SendThrough)}
	}
	if options.MultiPathTCP {
		SetDialerMultiPathTCP(&out.Dialer, true)
	}

	return out
}

type ControlFunc = func(network string, address string, c syscall.RawConn) error

type systemOutbound struct {
	net.Dialer
}

func (o *systemOutbound) DialTCPContext(ctx context.Context, network, address string) (*net.TCPConn, error) {
	conn, err := o.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}
