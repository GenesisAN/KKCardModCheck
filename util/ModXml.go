package util

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/GenesisAN/illusionsCard/util"
	"io"
	"strings"
)

type ModXml struct {
	GUID        string `xml:"guid"`
	Version     string `xml:"version"`
	Name        string `xml:"name"`
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Website     string `xml:"website"`
	Game        string `xml:"game"`
	BDMD5       string `gorm:"uniqueIndex"`
	Path        string
	Upload      bool
}

func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{}
	//===============BDMD5计算===================
	bdmd5, err := util.GetFileBDMD5(src)
	if err != nil {
		return mod, err
	}

	//==========================================

	zr, err := zip.OpenReader(src) //打开modzip
	if err != nil {
		return mod, errors.New(fmt.Sprintf("打开zip失败:%s", err.Error()))
	}

	for _, v := range zr.File { //遍历里面的文件
		if v.FileInfo().Name() == "manifest.xml" { //找到manifest.xml
			fr, _ := v.Open()         //打开它
			data, _ := io.ReadAll(fr) //读取里面的所有内容
			//构造XML的结构体
			err = xml.Unmarshal(data, &mod) //内容反序列化，按指定格式拿到里面的数据
			if err != nil {
				return mod, errors.New(fmt.Sprintf("读取zipmod信息失败:%s", err.Error()))
			}

			if err != nil {
				fmt.Errorf(err.Error())
				return mod, err
			}
			mod.BDMD5 = bdmd5
			mod.Path = strings.Replace(src, "mods\\", "", -1)
		}
	}
	return mod, err
}
