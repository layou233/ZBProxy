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
	Debug         bool

	IPAccess  access    `json:",omitempty"`
	Minecraft minecraft `json:",omitempty"`
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

	EnableMojangCapeRequirement bool `json:",omitempty"`
}

type configAnyDest struct {
	WildcardRootDomainName string `json:",omitempty"`
}
