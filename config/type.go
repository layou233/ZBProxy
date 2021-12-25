package config

type configMain struct {
	Services []ConfigProxyService
}

type ConfigProxyService struct {
	Name          string
	TargetAddress string
	TargetPort    uint16
	Listen        uint16
	Flow          string

	EnableHostnameRewrite bool
	RewrittenHostname     string

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
