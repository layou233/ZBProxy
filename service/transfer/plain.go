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
		go io.Copy(b, a)
		io.Copy(a, b)

	case FLOW_LINUX_ZEROCOPY:
		if runtime.GOOS != "linux" {
			log.Panic(color.HiRedString("Only Linux based systems support Linux ZeroCopy, please set your flow to origin or auto."))
		}
		fallthrough

	case FLOW_AUTO:
		go zerocopy.Transfer(b, a)
		zerocopy.Transfer(a, b)
	}
}
