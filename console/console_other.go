//go:build !windows
// +build !windows

package console

import "fmt"

func SetTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
