package transfer

import (
	"fmt"
	"math/rand"

	"github.com/layou233/ZBProxy/console"

	"github.com/fatih/color"
)

type ConnContext struct {
	ColoredID      string
	AdditionalInfo []string
}

func (c *ConnContext) AttachInfo(info string) {
	c.AdditionalInfo = append(c.AdditionalInfo, info)
}

func (c *ConnContext) Init() *ConnContext {
	id := rand.Int31()
	idColor := rand.Intn(console.ColorListN)
	c.ColoredID = color.New(console.ColorList[idColor]).Sprint("[", id, "]")

	c.AdditionalInfo = make([]string, 0, 1)

	return c
}

func (c *ConnContext) String() string {
	if len(c.AdditionalInfo) == 0 {
		return ""
	}
	return fmt.Sprint(c.AdditionalInfo)
}
