package outbound

import (
	"strings"
	"syscall"
)

const (
	// TCP_FASTOPEN_CONNECT for out-going connections.
	TCP_FASTOPEN_CONNECT = 30 //nolint: revive,stylecheck
)

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	if option == nil {
		return nil
	}
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if option.Mark != 0 {
				err = syscall.SetsockoptInt(fdInt, syscall.SOL_SOCKET, syscall.SO_MARK, option.Mark)
				if err != nil {
					return
				}
			}

			if option.Interface != "" {
				err = syscall.BindToDevice(fdInt, option.Interface)
				if err != nil {
					return
				}
			}

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.SOL_TCP, TCP_FASTOPEN_CONNECT, 1)
					if err != nil {
						return
					}
				}

				if option.TCPCongestion != "" {
					err = syscall.SetsockoptString(fdInt, syscall.SOL_TCP, syscall.TCP_CONGESTION, option.TCPCongestion)
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
