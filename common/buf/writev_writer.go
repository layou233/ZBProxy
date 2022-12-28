package buf

import (
	"io"
	"net"
	"syscall"
)

// WriterV is a writer that supports writing buffers
type WriterV struct {
	io.Writer
	rawConn syscall.RawConn
	mw      vectorizedWriter
}

type vectorizedWriter interface {
	Write(fd uintptr, buffers net.Buffers) int64
}

func NewWriterV(writer io.Writer, rawConn syscall.RawConn) (*WriterV, bool) {
	if mw, ok := newVectorizedWriter(); ok {
		return &WriterV{
			Writer:  writer,
			rawConn: rawConn,
			mw:      mw,
		}, true
	}
	return nil, false
}

func (w *WriterV) WriteVectorized(buffers net.Buffers) (n int64, err error) {
	err = w.rawConn.Write(func(fd uintptr) bool {
		n = w.mw.Write(fd, buffers)
		if n < 0 {
			n = 0
			return false
		}
		return true
	})
	return
}
