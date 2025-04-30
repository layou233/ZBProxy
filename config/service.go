package config

import "github.com/layou233/zbproxy/v3/common/network"

type Service struct {
	Name          string
	TargetAddress string `json:",omitempty"`
	TargetPort    uint16 `json:",omitempty"`
	Listen        uint16

	EnableProxyProtocol bool                          `json:",omitempty"`
	IPAccess            access                        `json:",omitempty"`
	Minecraft           *MinecraftService             `json:",omitempty"`
	TLSSniffing         *tlsSniffing                  `json:",omitempty"`
	SocketOptions       *network.InboundSocketOptions `json:",omitempty"`
	Outbound            outbound                      `json:",omitempty"`
}

type access struct {
	Mode      string   // 'accept' or 'deny' or empty
	ListTags  []string `json:",omitempty"`
	LowerCase bool     `json:",omitempty"`
}

type MinecraftService struct {
	EnableHostnameRewrite bool
	RewrittenHostname     string `json:",omitempty"`

	OnlineCount onlineCount

	IgnoreFMLSuffix   bool `json:",omitempty"`
	IgnoreSRVRedirect bool `json:",omitempty"`

	HostnameAccess access `json:",omitempty"`
	NameAccess     access `json:",omitempty"`

	PingMode        string
	MotdFavicon     string
	MotdDescription string
}

type onlineCount struct {
	Max            int32
	Online         int32
	EnableMaxLimit bool
	Sample         any `json:",omitempty"`
}

type tlsSniffing struct {
	RejectNonTLS     bool
	RejectIfNonMatch bool     `json:",omitempty"`
	SNIAllowListTags []string `json:",omitempty"`
}

type outbound struct {
	Type    string
	Network string `json:",omitempty"`
	Address string `json:",omitempty"`
}
