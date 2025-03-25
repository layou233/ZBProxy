package route

import (
	"context"
	"net"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common/bufio"
)

type rejectOutbound struct{}

var (
	_ adapter.Outbound       = rejectOutbound{}
	_ adapter.InjectOutbound = rejectOutbound{}
)

func (r rejectOutbound) Name() string {
	return "REJECT"
}

func (r rejectOutbound) PostInitialize(adapter.Router, adapter.RouteResourceProvider) error {
	return nil
}

func (r rejectOutbound) Reload(adapter.OutboundReloadOptions) error {
	return nil
}

func (r rejectOutbound) InjectConnection(ctx context.Context, conn *bufio.CachedConn, metadata *adapter.Metadata) error {
	return conn.Close()
}

func (r rejectOutbound) DialContext(context.Context, string, string) (net.Conn, error) {
	return nil, adapter.ErrInjectionRequired
}

type resetOutbound struct{}

var (
	_ adapter.Outbound       = resetOutbound{}
	_ adapter.InjectOutbound = resetOutbound{}
)

func (r resetOutbound) Name() string {
	return "RESET"
}

func (r resetOutbound) PostInitialize(adapter.Router, adapter.RouteResourceProvider) error {
	return nil
}

func (r resetOutbound) Reload(adapter.OutboundReloadOptions) error {
	return nil
}

func (r resetOutbound) InjectConnection(ctx context.Context, conn *bufio.CachedConn, metadata *adapter.Metadata) error {
	if tcpConn, isTCPConn := conn.Conn.(*net.TCPConn); isTCPConn {
		tcpConn.SetLinger(0)
	}
	return conn.Close()
}

func (r resetOutbound) DialContext(context.Context, string, string) (net.Conn, error) {
	return nil, adapter.ErrInjectionRequired
}
