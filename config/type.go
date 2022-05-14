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

	IPAccess  access `json:"omitempty"`
	Minecraft minecraft
}

type access struct {
	Mode     string `json:"omitempty"` // 'accept' or 'deny' or empty
	ListTags []string
}

type minecraft struct {
	EnableHostnameRewrite bool `json:"omitempty"`
	RewrittenHostname     string

	NameAccess access

	EnableAnyDest   bool
	AnyDestSettings configAnyDest

	MotdFavicon     string
	MotdDescription string

	EnableWhiteList             bool
	EnableMojangCapeRequirement bool
}

type configAnyDest struct {
	WildcardRootDomainName string
}
