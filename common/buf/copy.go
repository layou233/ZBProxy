package buf

import (
	"io"
	"syscall"
)

func Copy(dst io.Writer, src *ReaderV) (n int64, err error) {
	var writerV *WriterV
	useWriterV := false
	if sysConn, ok := dst.(syscall.Conn); ok {
		if rawConn, err := sysConn.SyscallConn(); err == nil {
			writerV, useWriterV = NewWriterV(dst, rawConn)
		}
	}

	for {
		buffers, err := src.ReadVectorized()
		if err != nil {
			PutMulti(buffers)
			return n, err
		}

		var written int64
		if useWriterV {
			written, err = writerV.WriteVectorized(buffers)
		} else {
			written, err = buffers.WriteTo(dst)
		}

		PutMulti(buffers)
		if err != nil {
			return n, err
		}
		n += written
	}
}
