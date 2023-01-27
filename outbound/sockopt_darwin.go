package outbound

import (
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	// TCP_FASTOPEN_CLIENT is the value to enable TCP fast open on darwin for client connections.
	TCP_FASTOPEN_CLIENT = 0x02 //nolint: revive,stylecheck
)

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, unix.TCP_FASTOPEN, TCP_FASTOPEN_CLIENT)
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
