package transfer

import (
	"acln.ro/zerocopy"
	"github.com/fatih/color"
	"github.com/xtls/xray-core/common/buf"
	"io"
	"log"
	"net"
	"runtime"
)

const (
	FLOW_ORIGIN = iota
	FLOW_LINUX_ZEROCOPY
	FLOW_MULTIPLE
	FLOW_AUTO
)

func SimpleTransfer(a, b net.Conn, flow int) {
	switch flow {
	case FLOW_ORIGIN:
		go io.Copy(b, a)
		io.Copy(a, b)

	case FLOW_LINUX_ZEROCOPY:
		if runtime.GOOS != "linux" {
			log.Panic(color.HiRedString("Only Linux based systems support Linux ZeroCopy, please set your flow to origin or auto."))
		}
		fallthrough

	case FLOW_AUTO:
		if runtime.GOOS == "linux" {
			go zerocopy.Transfer(b, a)
			zerocopy.Transfer(a, b)
			return
		}
		fallthrough

	case FLOW_MULTIPLE:
		aReader := buf.NewReader(a)
		bReader := buf.NewReader(b)
		aWriter := buf.NewWriter(a)
		bWriter := buf.NewWriter(b)
		go buf.Copy(bReader, aWriter)
		buf.Copy(aReader, bWriter)
	}
}
