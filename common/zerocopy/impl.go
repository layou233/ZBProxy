package zerocopy

import "net"

func CopyTCP(dst, src *net.TCPConn) (int64, error) {
	return dst.ReadFrom(src)
}
