package network

import "github.com/layou233/zbproxy/v3/common/jsonx"

type InboundSocketOptions struct {
	SendThrough     string
	KeepAlivePeriod jsonx.Duration `json:",omitempty"`
	Mark            int            `json:",omitempty"`
	TCPCongestion   string         `json:",omitempty"`
	TCPFastOpen     bool           `json:",omitempty"`
	MultiPathTCP    bool           `json:",omitempty"`
}

type OutboundSocketOptions struct {
	SendThrough     string         `json:",omitempty"`
	KeepAlivePeriod jsonx.Duration `json:",omitempty"`
	Mark            int            `json:",omitempty"`
	Interface       string         `json:",omitempty"`
	TCPCongestion   string         `json:",omitempty"`
	TCPFastOpen     bool           `json:",omitempty"`
	MultiPathTCP    bool           `json:",omitempty"`
}

func ConvertLegacyOutboundOptions(inbound *InboundSocketOptions) *OutboundSocketOptions {
	if inbound == nil {
		return nil
	}
	return &OutboundSocketOptions{
		SendThrough:     inbound.SendThrough,
		KeepAlivePeriod: inbound.KeepAlivePeriod,
		Mark:            inbound.Mark,
		TCPCongestion:   inbound.TCPCongestion,
		TCPFastOpen:     inbound.TCPFastOpen,
		MultiPathTCP:    inbound.MultiPathTCP,
	}
}
