package transfer

import (
	"acln.ro/zerocopy"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"runtime"
)

const (
	FLOW_ORIGIN = iota
	FLOW_LINUX_ZEROCOPY
	FLOW_AUTO
)

func SimpleTransfer(a, b net.Conn, flow int) {
	switch flow {
	case FLOW_ORIGIN:
		io.Copy(a, b)
		go io.Copy(b, a)

	case FLOW_LINUX_ZEROCOPY:
		if runtime.GOOS != "linux" {
			log.Panic(color.HiRedString("Only Linux based systems support Linux ZeroCopy, please set your flow to origin or auto."))
		}
		fallthrough

	case FLOW_AUTO:
		zerocopy.Transfer(a, b)
		go zerocopy.Transfer(b, a)
	}
}
