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
