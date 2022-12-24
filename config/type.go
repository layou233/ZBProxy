package config

type configMain struct {
	Services []*ConfigProxyService
	Lists    map[string][]string
}

type ConfigProxyService struct {
	Name          string
	TargetAddress string
	TargetPort    uint16
	Listen        uint16
	Flow          string

	IPAccess    access      `json:",omitempty"`
	Minecraft   minecraft   `json:",omitempty"`
	TLSSniffing tlsSniffing `json:",omitempty"`
	Outbound    outbound    `json:",omitempty"`
}

type access struct {
	Mode     string   // 'accept' or 'deny' or empty
	ListTags []string `json:",omitempty"`
}

type minecraft struct {
	EnableHostnameRewrite bool
	RewrittenHostname     string `json:",omitempty"`

	OnlineCount onlineCount

	IgnoreFMLSuffix bool `json:",omitempty"`

	NameAccess access `json:",omitempty"`

	EnableAnyDest   bool          `json:",omitempty"`
	AnyDestSettings configAnyDest `json:",omitempty"`

	MotdFavicon     string
	MotdDescription string
}

type onlineCount struct {
	Max            int
	Online         int32
	EnableMaxLimit bool
}

type configAnyDest struct {
	WildcardRootDomainName string `json:",omitempty"`
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

type mojangUuidSingle struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type hypixelConfig struct {
	State bool         `json:"success,omitempty"`
	Guild hypixelGuild `json:"guild,omitempty"`
}

type hypixelGuild struct {
	Members []hypixelGuildMemeber `json:"members,omitempty"`
}

type hypixelGuildMemeber struct {
	Uuid string `json:"uuid,omitempty"`
}
