package transfer

import (
	"fmt"

	"github.com/layou233/ZBProxy/console"

	"github.com/fatih/color"
	"github.com/zhangyunhao116/fastrand"
)

type ConnContext struct {
	ColoredID      string
	AdditionalInfo []string
	Err            error
}

func (c *ConnContext) AttachInfo(info string) {
	c.AdditionalInfo = append(c.AdditionalInfo, info)
}

func (c *ConnContext) Init() *ConnContext {
	id := fastrand.Int31()
	idColor := fastrand.Intn(len(console.ColorList))
	c.ColoredID = color.New(console.ColorList[idColor]).Sprint("[", id, "]")

	c.AdditionalInfo = make([]string, 0, 1)

	return c
}

func (c *ConnContext) String() (info string) {
	if len(c.AdditionalInfo) != 0 {
		info = fmt.Sprint(c.AdditionalInfo)
	}
	if c.Err == nil {
		info += ": √"
	} else {
		info += ": " + c.Err.Error()
	}
	return
}
