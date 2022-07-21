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

	IPAccess  access    `json:",omitempty"`
	Minecraft minecraft `json:",omitempty"`
	Outbound  outbound  `json:",omitempty"`
}

type access struct {
	Mode     string   // 'accept' or 'deny' or empty
	ListTags []string `json:",omitempty"`
}

type minecraft struct {
	EnableHostnameRewrite bool
	RewrittenHostname     string `json:",omitempty"`

	IgnoreFMLSuffix bool

	NameAccess access `json:",omitempty"`

	EnableAnyDest   bool          `json:",omitempty"`
	AnyDestSettings configAnyDest `json:",omitempty"`

	MotdFavicon     string
	MotdDescription string
}

type configAnyDest struct {
	WildcardRootDomainName string `json:",omitempty"`
}

type outbound struct {
	Type           string
	Network        string
	Address        string
	DomainStrategy string
}
