package transfer

import (
	"ZBProxy/common/zerocopy"
	"github.com/fatih/color"
	"github.com/xtls/xray-core/common/buf"
	"io"
	"log"
	"net"
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

func SimpleTransfer(a, b *net.TCPConn, flow int) {
	switch flow {
	case FLOW_ORIGIN:
		go io.Copy(writerOnly{b}, a)
		io.Copy(writerOnly{a}, b)

	case FLOW_LINUX_ZEROCOPY:
		log.Print(color.HiRedString("The `linux-zerocopy` is deprecated, please use `zerocopy` instead."))
		fallthrough

	case FLOW_ZEROCOPY:
		fallthrough

	case FLOW_AUTO:
		go zerocopy.CopyTCP(b, a)
		zerocopy.CopyTCP(a, b)
		return // TODO: Use MULTIPLE when fail to sendfile or splice

	case FLOW_MULTIPLE:
		aReader := buf.NewReader(a)
		bReader := buf.NewReader(b)
		aWriter := buf.NewWriter(a)
		bWriter := buf.NewWriter(b)
		go buf.Copy(bReader, aWriter)
		buf.Copy(aReader, bWriter)
	}
}
