package util

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
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

func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{Path: src}

	zr, err := zip.OpenReader(src)
	if err != nil {
		return mod, fmt.Errorf("打开zip失败:%s", err)
	}
	defer zr.Close()

	for _, v := range zr.File {
		if v.FileInfo().Name() == "manifest.xml" {
			fr, err := v.Open()
			if err != nil {
				return mod, fmt.Errorf("打开manifest.xml失败:%s", err)
			}
			defer fr.Close()

			data, err := io.ReadAll(fr)
			if err != nil {
				return mod, fmt.Errorf("读取manifest.xml失败:%s", err)
			}

			// 首次尝试解码 XML
			if err = xml.Unmarshal(data, &mod); err != nil {
				if isLikelyUTF16(data) { // 如果失败了，尝试将utf-16 转换为 UTF-8 后再解码
					utf8Data := utf16ToUtf8(data)
					if err = xml.Unmarshal(utf8Data, &mod); err != nil {
						return mod, fmt.Errorf("UTF-16 转换后读取失败:%s", err)
					}
				} else {
					return mod, fmt.Errorf("读取zipmod信息失败:%s", err)
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
