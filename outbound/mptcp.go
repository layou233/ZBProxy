//go:build go1.21

package outbound

import "net"

func SetMultiPathTCP(dialer *net.Dialer, use bool) {
	dialer.SetMultipathTCP(use)
}
