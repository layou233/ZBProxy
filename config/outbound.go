package config

import "github.com/layou233/zbproxy/v3/common/network"

type Outbound struct {
	Name                 string                         `json:",omitempty"`
	Dialer               string                         `json:",omitempty"`
	TargetAddress        string                         `json:",omitempty"`
	TargetPort           uint16                         `json:",omitempty"`
	Minecraft            *MinecraftService              `json:",omitempty"`
	SocketOptions        *network.OutboundSocketOptions `json:",omitempty"`
	ProxyProtocolVersion int8                           `json:",omitempty"`
	ProxyOptions         proxyOptions                   `json:",omitempty"`
}
