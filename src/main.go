package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// MessageBox of Win32 API.
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {

	ret, _, _ := windows.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

// MessageBoxPlain of Win32 API.
func MessageBoxPlain(title, caption string) int {
	const (
		NULL  = 0
		MB_OK = 0
	)
	return MessageBox(NULL, caption, title, MB_OK)
}

func main() {
	MessageBoxPlain("Test", "Hello World!")
}
