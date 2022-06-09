package access

import "log"

const (
	DefaultMode = iota
	AllowMode
	BlockMode
)

func GetAccessMode(mode string) int {
	switch mode {
	case "allow", "whitelist":
		return AllowMode
	case "block", "blacklist":
		return BlockMode
	case "":
		return DefaultMode
	default:
		log.Panicf("Unknown access control mode: %q", mode)
		return 0
	}
}
