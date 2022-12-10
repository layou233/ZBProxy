//go:build windows
// +build windows

package console

import (
	"syscall"
	"unsafe"
)

var SetConsoleTileW = syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleTitleW")

func SetTitle(title string) {
	SetConsoleTileW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title)))) //nolint:errcheck,staticcheck
}
