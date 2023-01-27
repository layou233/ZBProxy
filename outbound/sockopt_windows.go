package outbound

import (
	"strings"
	"syscall"
)

const (
	TCP_FASTOPEN = 15 //nolint: revive,stylecheck
)

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			handle := syscall.Handle(fd)

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(handle, syscall.IPPROTO_TCP, TCP_FASTOPEN, 1)
					if err != nil {
						return
					}
				}
			}
		})
		if err != nil {
			return err
		}
		return err_
	}
}
