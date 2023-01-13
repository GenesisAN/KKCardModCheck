//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// RunningByDoubleClick 检查是否通过双击直接运行
func RunningByDoubleClick() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var ids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&ids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}

// NoMoreDoubleClick 提示用户不要双击运行，并生成安全启动脚本
func NoMoreDoubleClick() error {
	r := boxW(0, "请勿通过双击直接运行本程序\n这会因CMD字符格式问题导致界面渲染出现问题\n点击确认将释出启动脚本，点击取消则关闭程序", "警告", 0x00000030|0x00000001)
	if r == 2 {
		return nil
	}
	r = boxW(0, "点击确认将覆盖start.bat，点击取消则关闭程序", "警告", 0x00000030|0x00000001)
	if r == 2 {
		return nil
	}
	f, err := os.OpenFile("start.bat", os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Printf("打开start.bat失败: %v", err)
		return err
	}
	_ = f.Truncate(0)

	ex, _ := os.Executable()
	exPath := filepath.Base(ex)
	_, err = f.WriteString("%Created by KKCardModCheck. DO NOT EDIT ME!%\nchcp 65001 \n\"" + exPath + "\" -s")
	if err != nil {
		fmt.Printf("写入启动.bat失败: %v", err)
		return err
	}
	f.Close()
	boxW(0, "启动脚本已生成，请双击start.bat启动", "提示", 0x00000040|0x00000000)
	return nil
}

// BoxW of Win32 API. Check https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw for more detail.
func boxW(hwnd uintptr, caption, title string, flags uint) int {
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(captionPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags))

	return int(ret)
}
