package transfer

import (
	"io"
	"log"
	"net"
	"runtime"

	"github.com/fatih/color"
	"github.com/xtls/xray-core/common/buf"
)

const (
	FLOW_ORIGIN = iota
	FLOW_LINUX_ZEROCOPY
	FLOW_ZEROCOPY
	FLOW_MULTIPLE
	FLOW_AUTO
)

type writerOnly struct {
	io.Writer
}

func SimpleTransfer(a, b net.Conn, flow int) {
	//nolint:errcheck
	switch flow {
	case FLOW_ORIGIN:
		go func() {
			io.Copy(writerOnly{b}, a)
			a.Close()
			b.Close()
		}()
		io.Copy(writerOnly{a}, b)
		a.Close()
		b.Close()

	case FLOW_ZEROCOPY:
		fallthrough

	case FLOW_LINUX_ZEROCOPY:
		if runtime.GOOS != "linux" {
			log.Panic(color.HiRedString("Only Linux based systems support Linux ZeroCopy, please set your flow to origin or auto."))
		}
		fallthrough

	case FLOW_AUTO:
		if runtime.GOOS == "linux" {
			go func() {
				io.Copy(b, a)
				a.Close()
				b.Close()
			}()
			io.Copy(a, b)
			a.Close()
			b.Close()
			return // TODO: Use MULTIPLE when fail to sendfile or splice
		}
		fallthrough

	case FLOW_MULTIPLE:
		aReader := buf.NewReader(a)
		bReader := buf.NewReader(b)
		aWriter := buf.NewWriter(a)
		bWriter := buf.NewWriter(b)

		go func() {
			buf.Copy(bReader, aWriter)
			a.Close()
			b.Close()
		}()
		buf.Copy(aReader, bWriter)
		a.Close()
		b.Close()
	}
}
