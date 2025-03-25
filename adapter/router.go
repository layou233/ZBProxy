package adapter

import (
	"net"

	"github.com/layou233/zbproxy/v3/common/set"
)

type Router interface {
	RouteResourceProvider
	HandleConnection(conn net.Conn, metadata *Metadata)
}

type RouteResourceProvider interface {
	FindOutboundByName(name string) (Outbound, error)
	FindListsByTag(tags []string) ([]set.StringSet, error)
}
