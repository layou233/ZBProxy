//go:build windows

package buf

import (
	"net"
	"syscall"
)

type windowsWriter struct{}

func (w windowsWriter) Write(fd uintptr, buffers net.Buffers) int64 {
	bufs := make([]syscall.WSABuf, 0, len(buffers))
	for _, buffer := range buffers {
		bufs = append(bufs, syscall.WSABuf{
			Len: uint32(len(buffer)),
			Buf: &buffer[0],
		})
	}
	var n uint32
	if syscall.WSASend(syscall.Handle(fd), &bufs[0], uint32(len(bufs)), &n, 0, nil, nil) == nil {
		return int64(n)
	}
	return -1
}

func newVectorizedWriter() (vectorizedWriter, bool) {
	return new(windowsWriter), true
}
