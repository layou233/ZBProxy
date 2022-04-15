package zerocopy

import "net"

func CopyTCP(dst, src *net.TCPConn) (int64, error) {
	defer dst.Close()
	return dst.ReadFrom(src)
}
