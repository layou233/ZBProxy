package protocol

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/network"
	"github.com/layou233/zbproxy/v3/common/network/socks"
	"github.com/layou233/zbproxy/v3/common/proxyprotocol"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/protocol/minecraft"

	"github.com/phuslu/log"
)

func NewOutbound(logger *log.Logger, newConfig *config.Outbound) (adapter.Outbound, error) {
	if newConfig == nil {
		return nil, os.ErrInvalid
	}
	switch {
	case newConfig.Minecraft != nil:
		return minecraft.NewOutbound(logger, newConfig)
	}
	return &Plain{
		logger: logger,
		config: newConfig,
	}, nil
}

type Plain struct {
	logger *log.Logger
	config *config.Outbound
	router adapter.Router
	dialer network.Dialer
}

var (
	_ adapter.Outbound         = (*Plain)(nil)
	_ adapter.MetadataOutbound = (*Plain)(nil)
	_ network.Dialer           = (*Plain)(nil)
)

func (o *Plain) Name() string {
	if o.config != nil {
		return o.config.Name
	}
	return ""
}

func (o *Plain) PostInitialize(router adapter.Router, provider adapter.RouteResourceProvider) error {
	var err error
	if o.config.Dialer != "" {
		if o.config.SocketOptions != nil {
			return errors.New("socket options are not available when dialer is specified")
		}
		o.dialer, err = provider.FindOutboundByName(o.config.Dialer)
		if err != nil {
			return err
		}
	} else {
		o.dialer = network.NewSystemDialer(o.config.SocketOptions)
	}
	switch o.config.ProxyProtocolVersion {
	case proxyprotocol.VersionUnspecified,
		proxyprotocol.Version1,
		proxyprotocol.Version2:
	default:
		return fmt.Errorf("invalid proxy protocol version: %v", o.config.ProxyProtocolVersion)
	}
	switch o.config.ProxyOptions.Type {
	case "socks", "socks5", "socks4a", "socks4":
		o.dialer = &socks.Client{
			Dialer:  o.dialer,
			Version: o.config.ProxyOptions.Type,
			Network: o.config.ProxyOptions.Network,
			Address: o.config.ProxyOptions.Address,
		}
	}
	o.router = router
	return nil
}

func (o *Plain) Reload(options adapter.OutboundReloadOptions) error {
	o.config = options.Config
	return o.PostInitialize(o.router, &options)
}

func (o *Plain) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	return o.dialer.DialContext(ctx, network, address)
}

func (o *Plain) DialContextWithMetadata(ctx context.Context, network string, address string, metadata *adapter.Metadata) (net.Conn, error) {
	conn, err := adapter.DialContextWithMetadata(o.dialer, ctx, network, address, metadata)
	if err != nil {
		return nil, err
	}
	if o.config.ProxyProtocolVersion != proxyprotocol.VersionUnspecified {
		var localAddress netip.AddrPort
		localAddress, err = netip.ParseAddrPort(conn.LocalAddr().String())
		if err != nil {
			conn.Close()
			return nil, common.Cause("failed to parse local address: ", err)
		}
		err = (&proxyprotocol.Header{
			Version:           uint8(o.config.ProxyProtocolVersion),
			Command:           proxyprotocol.CommandProxy,
			TransportProtocol: proxyprotocol.TransportProtocolByNetwork(network) | proxyprotocol.AddressFamilyByAddr(metadata.SourceAddress.Addr()),
			SourceAddress:     metadata.SourceAddress,
		}).WriteHeader(conn, localAddress)
		if err != nil {
			conn.Close()
			return nil, common.Cause("failed to write PROXY protocol header: ", err)
		}
	}
	return conn, nil
}
