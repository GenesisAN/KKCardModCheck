package util

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type ModXml struct {
	GUID        string `xml:"guid"`
	Version     string `xml:"version"`
	Name        string `xml:"name"`
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Website     string `xml:"website"`
	Game        string `xml:"game"`
	//BDMD5       string `gorm:"uniqueIndex"`
	Path   string
	Upload bool
}

func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{}
	//===============BDMD5 Build ===================
	//bdmd5, err := util.GetFileBDMD5(src)
	//if err != nil {
	//	return mod, fail
	//}

	//==========================================
	mod.Path = src
	zr, err := zip.OpenReader(src) //open modzip
	if err != nil {
		mod.GUID = fmt.Sprintf("打开zip失败:%s", err.Error())
		return mod, errors.New(fmt.Sprintf("打开zip失败:%s", err.Error()))
	}

	for _, v := range zr.File {
		if v.FileInfo().Name() == "manifest.xml" { //find manifest.xml
			fr, _ := v.Open()
			data, _ := io.ReadAll(fr) //read manifest.xml
			//构造XML的结构体
			err = xml.Unmarshal(data, &mod) //unmarshal xml to struct
			if err != nil {
				mod.GUID = fmt.Sprintf("打开zip失败:%s", err.Error())
				return mod, errors.New(fmt.Sprintf("读取zipmod信息失败:%s", err.Error()))
			}
			//mod.BDMD5 = bdmd5
			//mod.Path = strings.Replace(src, "mods\\", "", -1)
		}
	}
	return mod, err
}
