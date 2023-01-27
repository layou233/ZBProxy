//go:build !linux && !windows && !darwin && !freebsd

package outbound

import "syscall"

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	return emptyDialerControl
}

func emptyDialerControl(_ string, _ string, _ syscall.RawConn) error {
	return nil
}
