//go:build !linux && !windows && !darwin && !freebsd

package outbound

func NewDialerControlFromOptions(option *SocketOptions) DialerControl {
	return nil
}
