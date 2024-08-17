package network

import (
	"time"

	"github.com/layou233/zbproxy/v3/common/jsonx"
)

type InboundSocketOptions struct {
	keepAliveOptions
	Mark          int    `json:",omitempty"`
	SendThrough   string `json:",omitempty"`
	TCPCongestion string `json:",omitempty"`
	TCPFastOpen   bool   `json:",omitempty"`
	MultiPathTCP  bool   `json:",omitempty"`
}

type OutboundSocketOptions struct {
	keepAliveOptions
	Mark          int    `json:",omitempty"`
	Interface     string `json:",omitempty"`
	SendThrough   string `json:",omitempty"`
	TCPCongestion string `json:",omitempty"`
	TCPFastOpen   bool   `json:",omitempty"`
	MultiPathTCP  bool   `json:",omitempty"`
}

type keepAliveOptions struct {
	KeepAliveIdle     jsonx.Duration `json:",omitempty"`
	KeepAliveInterval jsonx.Duration `json:",omitempty"`
	KeepAliveCount    int            `json:",omitempty"`
}

func (o *keepAliveOptions) KeepAliveConfig() (c KeepAliveConfig) {
	c = KeepAliveConfig{
		Idle:     time.Duration(o.KeepAliveIdle),
		Interval: time.Duration(o.KeepAliveInterval),
		Count:    o.KeepAliveCount,
	}
	if c.Idle > 0 || c.Interval > 0 || c.Count > 0 {
		c.Enable = true
	}
	return
}

func ConvertLegacyOutboundOptions(inbound *InboundSocketOptions) *OutboundSocketOptions {
	if inbound == nil {
		return nil
	}
	return &OutboundSocketOptions{
		keepAliveOptions: keepAliveOptions{
			KeepAliveIdle:     inbound.KeepAliveIdle,
			KeepAliveInterval: inbound.KeepAliveInterval,
			KeepAliveCount:    inbound.KeepAliveCount,
		},
		Mark:          inbound.Mark,
		SendThrough:   inbound.SendThrough,
		TCPCongestion: inbound.TCPCongestion,
		TCPFastOpen:   inbound.TCPFastOpen,
		MultiPathTCP:  inbound.MultiPathTCP,
	}
}
