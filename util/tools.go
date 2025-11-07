package util

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgryski/go-farm"
)

// GetAllFiles 遍历目录获取指定后缀文件
func GetAllFiles(root, ext string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ext {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// GetFileHash 使用 farmhash 的 Hash32 计算文件哈希
func GetFileHash32(filePath string) (uint32, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 读取文件内容并计算哈希
	buf := make([]byte, 4096) // 每次读取 4KB
	h := uint32(0)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}
		h = farm.Hash32(buf[:n])
	}

	return h, nil
}

// GetFileHash64 使用 farmhash 的 Hash64 计算文件哈希
func GetFileHash64(filePath string) (uint64, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 读取文件内容并计算哈希
	buf := make([]byte, 4096) // 每次读取 4KB
	h := uint64(0)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}
		h = farm.Hash64(buf[:n])
	}

	return h, nil
}

// IsExist 文件/路径存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsNotExist 文件/路径不存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}
func IsWin() bool {
	return runtime.GOOS == "windows"
}
func RemoveNewlines(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
}

// DedupeStrings 对字符串切片去重，保留首次出现的顺序
func DedupeStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, s := range items {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
