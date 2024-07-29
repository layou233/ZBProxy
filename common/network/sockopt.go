package network

import "github.com/layou233/zbproxy/v3/common/jsonx"

type InboundSocketOptions struct {
	KeepAlivePeriod jsonx.Duration `json:",omitempty"`
	Mark            int            `json:",omitempty"`
	SendThrough     string         `json:",omitempty"`
	TCPCongestion   string         `json:",omitempty"`
	TCPFastOpen     bool           `json:",omitempty"`
	MultiPathTCP    bool           `json:",omitempty"`
}

type OutboundSocketOptions struct {
	KeepAlivePeriod jsonx.Duration `json:",omitempty"`
	Mark            int            `json:",omitempty"`
	Interface       string         `json:",omitempty"`
	SendThrough     string         `json:",omitempty"`
	TCPCongestion   string         `json:",omitempty"`
	TCPFastOpen     bool           `json:",omitempty"`
	MultiPathTCP    bool           `json:",omitempty"`
}

func ConvertLegacyOutboundOptions(inbound *InboundSocketOptions) *OutboundSocketOptions {
	if inbound == nil {
		return nil
	}
	return &OutboundSocketOptions{
		KeepAlivePeriod: inbound.KeepAlivePeriod,
		Mark:            inbound.Mark,
		SendThrough:     inbound.SendThrough,
		TCPCongestion:   inbound.TCPCongestion,
		TCPFastOpen:     inbound.TCPFastOpen,
		MultiPathTCP:    inbound.MultiPathTCP,
	}
}
