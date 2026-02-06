//go:build !windows

package util

import (
	"os/exec"
	"runtime"
)

// CopyToClipboard 将文本复制到系统剪贴板（非 Windows）
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		// linux: 尝试 xclip（常见）
		// 如果 xclip 不存在，返回错误让调用方决定
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, _ = stdin.Write([]byte(text))
	_ = stdin.Close()
	return cmd.Wait()
}
