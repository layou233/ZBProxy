package console

import (
	"fmt"

	"github.com/fatih/color"
)

func Println(a ...interface{}) {
	fmt.Fprintln(color.Output, a...)
}

func Printf(format string, a ...interface{}) {
	fmt.Fprintf(color.Output, format, a...)
}

var ColorList = [...]color.Attribute{
	color.FgRed, color.FgGreen, color.FgYellow,
	color.FgBlue, color.FgMagenta, color.FgCyan, color.FgWhite,

	color.FgHiRed, color.FgHiGreen, color.FgHiYellow,
	color.FgHiBlue, color.FgHiMagenta, color.FgHiCyan, color.FgHiWhite,
}
