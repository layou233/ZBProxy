package adapter

import (
	"context"
	"errors"
	"net"

	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/network"
)

type Outbound interface {
	Name() string
	PostInitialize(router Router, provider RouteResourceProvider) error
	Reload(options OutboundReloadOptions) error
	DialContext(ctx context.Context, network string, address string) (net.Conn, error)
}

type InjectOutbound interface {
	InjectConnection(ctx context.Context, conn *bufio.CachedConn, metadata *Metadata) error
}

type MetadataOutbound interface {
	DialContextWithMetadata(ctx context.Context, network string, address string, metadata *Metadata) (net.Conn, error)
}

type SRVOutbound interface {
	DialContextWithSRV(ctx context.Context, network string, address string, serviceName string) (net.Conn, error)
}

func DialContextWithMetadata(dialer network.Dialer, ctx context.Context, network, addr string, metadata *Metadata) (net.Conn, error) {
	if metadataOutbound, isMetadataOutbound := dialer.(MetadataOutbound); isMetadataOutbound {
		return metadataOutbound.DialContextWithMetadata(ctx, network, addr, metadata)
	} else if metadata.SRV != nil {
		if srvOutbound, isSRVOutbound := dialer.(SRVOutbound); isSRVOutbound {
			return srvOutbound.DialContextWithSRV(ctx, network, addr, metadata.SRV.ServiceName)
		}
	}
	return dialer.DialContext(ctx, network, addr)
}

var ErrInjectionRequired = errors.New("injection required")
