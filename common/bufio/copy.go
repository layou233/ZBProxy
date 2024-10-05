package bufio

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime"

	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
)

func CopyConn(remote net.Conn, local net.Conn) error {
	done := make(chan struct{})
	var errLocal, errRemote error
	go func() {
		_, errRemote = CopyBuffer(local, remote, nil)
		local.Close()
		close(done)
	}()
	_, errLocal = CopyBuffer(remote, local, nil)
	remote.Close()
	<-done
	if errLocal != nil || errRemote != nil {
		return fmt.Errorf("relay connnections: download: %w | upload: %w", errRemote, errLocal)
	}
	return nil
}

func Copy(destination io.Writer, source io.Reader) (written int64, err error) {
	return CopyBuffer(destination, source, nil)
}

func CopyBuffer(destination io.Writer, source io.Reader, buffer *buf.Buffer) (written int64, err error) {
	for {
		destination, source = common.UnwrapWriter(destination), common.UnwrapReader(source)
		if cachedConn, isSourceCachedConn := source.(*CachedConn); isSourceCachedConn {
			if cachedConn.cache == nil || cachedConn.cache.Len() == 0 {
				source = cachedConn.Conn
				cachedConn.Release()
				continue
			}
			written, err = cachedConn.cache.WriteTo(destination)
			if err != nil {
				return
			}
			continue
		}
		break
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "android" { // Linux optimizations
		if destinationTCPConn, isDestinationTCP := destination.(*net.TCPConn); isDestinationTCP {
			switch typedSource := source.(type) {
			case *net.TCPConn, *net.UnixConn, *os.File:
				println("real copy!!!")
				written, err = io.Copy(destinationTCPConn, typedSource)
				switch common.Unwrap(err) {
				case io.EOF, net.ErrClosed:
					err = nil
				}
				return
			}
		}
	}
	if buffer == nil {
		buffer = buf.NewSize(16 * 1024)
		written, err = CopyBuffer(destination, source, buffer)
		defer buffer.Release()
	}
	for {
		buffer.Reset(0) // TODO: headroom support
		var nRead, nWrite int64
		var errWrite error
		nRead, err = buffer.ReadOnceFrom(source)
		if nRead > 0 {
			nWrite, errWrite = buffer.WriteTo(destination)
			written += nWrite
			if errWrite != nil {
				return written, errWrite
			}
			if nWrite != nRead {
				return written, io.ErrShortWrite
			}
		}
		if err != nil {
			switch common.Unwrap(err) {
			case io.EOF, net.ErrClosed:
				err = nil
			}
			break
		}
	}
	return
}
