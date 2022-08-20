package transfer

import (
	"github.com/layou233/ZBProxy/outbound"
	"sync/atomic"
)

type Options struct {
	Out                     outbound.Outbound
	IsTLSHandleNeeded       bool
	IsMinecraftHandleNeeded bool
	FlowType                int
	McNameMode              int
	onlineCount             int32
}

func (receiver Options) AddCount(n int32) {
	atomic.AddInt32(&receiver.onlineCount, n)
}

func (receiver Options) GetCount() int32 {
	return receiver.onlineCount
}
