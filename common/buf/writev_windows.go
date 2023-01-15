//go:build windows

package buf

import (
	"net"
	"syscall"
)

type windowsWriter struct {
	bufs []syscall.WSABuf
}

func (w *windowsWriter) Write(fd uintptr, buffers net.Buffers) int64 {
	if w.bufs == nil {
		w.bufs = make([]syscall.WSABuf, 0, len(buffers))
	} else {
		w.bufs = w.bufs[:0]
	}
	for _, buffer := range buffers {
		w.bufs = append(w.bufs, syscall.WSABuf{
			Len: uint32(len(buffer)),
			Buf: &buffer[0],
		})
	}
	var n uint32
	if syscall.WSASend(syscall.Handle(fd), &w.bufs[0], uint32(len(w.bufs)), &n, 0, nil, nil) == nil {
		return int64(n)
	}
	return -1
}

func newVectorizedWriter() (vectorizedWriter, bool) {
	return new(windowsWriter), true
}
