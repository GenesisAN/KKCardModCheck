package util

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"unicode/utf16"
)

type ModXml struct {
	GUID        string `xml:"guid"`
	Version     string `xml:"version"`
	Name        string `xml:"name"`
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Website     string `xml:"website"`
	Game        string `xml:"game"`
	Path        string
	Upload      bool
}

// ReadZip 打开ZIP找到 manifest.xml 并尝试解析 ModXml
func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{Path: src}

	//检查文件大小,如果文件大小为0,则返回错误
	if fi, err := os.Stat(src); err != nil {
		mod.Name = fmt.Sprintf("获取文件信息失败:%s", err)
		return mod, fmt.Errorf("获取文件信息失败:%s", err)
	} else if fi.Size() == 0 {
		mod.Name = "文件大小为0"
		return mod, fmt.Errorf("文件大小为0")
	}
	zr, err := zip.OpenReader(src)

	if err != nil {
		mod.Name = fmt.Sprintf("打开zip失败:%s", err)
		return mod, fmt.Errorf("打开zip失败:%s", err)
	}
	defer zr.Close()

	for _, v := range zr.File {
		if v.FileInfo().Name() == "manifest.xml" {
			fr, err := v.Open()
			if err != nil {
				mod.Name = fmt.Sprintf("打开manifest.xml失败:%s", err)
				return mod, fmt.Errorf("打开manifest.xml失败:%s", err)
			}
			defer fr.Close()

			data, err := io.ReadAll(fr)
			if err != nil {
				mod.Name = fmt.Sprintf("读取manifest.xml失败:%s", err)
				return mod, fmt.Errorf("读取manifest.xml失败:%s", err)
			}

			// 如果可能是 UTF-16，则先转为 UTF-8
			if isLikelyUTF16(data) {
				data = utf16ToUtf8(data)
			}

			// 第一次解析
			if err := xml.Unmarshal(data, &mod); err != nil {
				// 如果失败了，尝试“移除多余闭合标签”再解析
				fixed := removeUnmatchedClosingTags(string(data))
				if err2 := xml.Unmarshal([]byte(fixed), &mod); err2 != nil {
					mod.Name = fmt.Sprintf("解析manifest.xml失败: %v; 修补后仍失败: %v", err, err2)
					return mod, fmt.Errorf("解析manifest.xml失败: %v; 修补后仍失败: %v", err, err2)
				}
			}
			break // 找到manifest.xml后停止查找
		}
	}
	return mod, nil
}

// 判断是否为UTF-16
func isLikelyUTF16(data []byte) bool {
	return len(data) >= 2 && (data[0] == 0xFF && data[1] == 0xFE || data[0] == 0xFE && data[1] == 0xFF)
}

// 将UTF-16数据转换为UTF-8
func utf16ToUtf8(data []byte) []byte {
	u16 := make([]uint16, len(data)/2)
	for i := range u16 {
		u16[i] = uint16(data[2*i]) | uint16(data[2*i+1])<<8
	}
	runes := utf16.Decode(u16)
	return []byte(string(runes)) // 转换为字符串后再转换为[]byte
}
