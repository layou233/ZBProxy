package buf

import (
	"io"
	"net"
)

func Copy(dst io.Writer, src *ReaderV) (n int64, err error) {
	for {
		var buffers net.Buffers
		buffers, err = src.ReadVectorized()
		if err != nil {
			PutMulti(buffers)
			return
		}

		var written int64
		written, err = buffers.WriteTo(dst)

		PutMulti(buffers)
		if err != nil {
			return
		}
		n += written
	}
}
