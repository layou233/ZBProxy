package outbound

import (
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	if option == nil {
		return nil
	}
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if option.Mark != 0 {
				err = syscall.SetsockoptInt(fdInt, syscall.SOL_SOCKET, syscall.SO_USER_COOKIE, option.Mark)
				if err != nil {
					return
				}
			}

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, unix.TCP_FASTOPEN, 1)
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
