//go:build illumos

package buf

import (
	"net"

	"golang.org/x/sys/unix"
)

type unixReader struct {
	iovs net.Buffers
}

func (r *unixReader) Init(bs net.Buffers) {
	iovs := r.iovs
	if iovs == nil {
		iovs = make(net.Buffers, 0, len(bs))
	}
	for _, b := range bs {
		iovs = append(iovs, b)
	}
	r.iovs = iovs
}

func (r *unixReader) Read(fd uintptr) int32 {
	n, e := unix.Readv(int(fd), r.iovs)
	if e != nil {
		return -1
	}
	return int32(n)
}

func (r *unixReader) Clear() {
	r.iovs = r.iovs[:0]
}

func newVectorizedReader() vectorizedReader {
	return &unixReader{}
}
