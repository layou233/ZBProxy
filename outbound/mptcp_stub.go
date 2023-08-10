//go:build !go1.21

package outbound

import "net"

func SetMultiPathTCP(*net.Dialer, bool) {
	panic("MultiPath TCP requires go1.21, please recompile your binary.")
}
