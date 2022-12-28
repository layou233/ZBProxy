package transfer

import (
	"github.com/layou233/ZBProxy/common"
	"io"
	"log"
	"net"
	"runtime"
	"syscall"

	"github.com/layou233/ZBProxy/common/buf"

	"github.com/fatih/color"
)

const (
	FLOW_ORIGIN = iota
	FLOW_LINUX_ZEROCOPY
	FLOW_ZEROCOPY
	FLOW_MULTIPLE
	FLOW_AUTO

	osSupportSplice = runtime.GOOS == "linux" || runtime.GOOS == "android"
)

type writerOnly struct {
	io.Writer
}

func SimpleTransfer(a, b net.Conn, flow int) {
	//nolint:errcheck
	switch flow {
	case FLOW_ORIGIN:
		go func() {
			buffer := buf.Get(32 * 1024)
			io.CopyBuffer(writerOnly{b}, a, buffer)
			buf.Put(buffer)
			a.Close()
			b.Close()
		}()
		buffer := buf.Get(32 * 1024)
		io.CopyBuffer(writerOnly{a}, b, buffer)
		buf.Put(buffer)
		a.Close()
		b.Close()

	case FLOW_ZEROCOPY:
		fallthrough

	case FLOW_LINUX_ZEROCOPY:
		if !osSupportSplice {
			log.Panic(color.HiRedString("Only Linux based systems support Linux ZeroCopy, please set your flow to origin or auto."))
		}
		fallthrough

	case FLOW_AUTO:
		if osSupportSplice {
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
		aReader := buf.NewReaderV(a, common.Must(a.(syscall.Conn).SyscallConn()))
		bReader := buf.NewReaderV(b, common.Must(b.(syscall.Conn).SyscallConn()))
		//aWriter := buf.NewWriter(a)
		//bWriter := buf.NewWriter(b)

		go func() {
			buf.Copy(a, bReader)
			a.Close()
			b.Close()
		}()
		buf.Copy(b, aReader)
		a.Close()
		b.Close()
	}
}
