package windows

import (
	"syscall"
	"unsafe"
)

func SetTitle(title string) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	kernel32.NewProc("SetConsoleTitleW").
		Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
}
