package transfer

import (
	"sync/atomic"

	"github.com/layou233/ZBProxy/outbound"
)

type Options struct {
	Out                     outbound.Outbound
	IsTLSHandleNeeded       bool
	IsMinecraftHandleNeeded bool
	FlowType                int
	OnlineCount             atomic.Int32
}
