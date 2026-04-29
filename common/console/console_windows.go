//go:build windows
// +build windows

package console

import (
	"syscall"
	"unsafe"
)

var SetConsoleTitleW = syscall.MustLoadDLL("kernel32.dll").MustFindProc("SetConsoleTitleW")

func SetTitle(title string) {
	SetConsoleTitleW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title)))) //nolint:errcheck,staticcheck
}
