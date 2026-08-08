//go:build windows

package util

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	cfUnicodeText = 13
	gmemMoveable = 0x0002
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	openClipboard     = user32.NewProc("OpenClipboard")
	closeClipboard    = user32.NewProc("CloseClipboard")
	emptyClipboard    = user32.NewProc("EmptyClipboard")
	setClipboardData  = user32.NewProc("SetClipboardData")
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	globalAlloc       = kernel32.NewProc("GlobalAlloc")
	globalLock        = kernel32.NewProc("GlobalLock")
	globalUnlock      = kernel32.NewProc("GlobalUnlock")
)

// CopyToClipboard 将文本复制到系统剪贴板（Windows，无控制台窗口）
func CopyToClipboard(text string) error {
	utf16Text, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}

	r, _, err := openClipboard.Call(0)
	if r == 0 {
		return fmt.Errorf("OpenClipboard failed: %v", err)
	}
	defer closeClipboard.Call()

	emptyClipboard.Call()

	size := uintptr(len(utf16Text) * 2)
	hMem, _, err := globalAlloc.Call(gmemMoveable, size)
	if hMem == 0 {
		return fmt.Errorf("GlobalAlloc failed: %v", err)
	}

	ptr, _, err := globalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("GlobalLock failed: %v", err)
	}
	copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:size:size], (*[1 << 30]byte)(unsafe.Pointer(&utf16Text[0]))[:size:size])
	globalUnlock.Call(hMem)

	r, _, err = setClipboardData.Call(cfUnicodeText, hMem)
	if r == 0 {
		return fmt.Errorf("SetClipboardData failed: %v", err)
	}
	return nil
}
