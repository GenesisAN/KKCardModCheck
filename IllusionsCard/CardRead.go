package IllusionsCard

import (
	"KKCardModCheck/IllusionsCard/Base"
	"KKCardModCheck/IllusionsCard/KK"
	"KKCardModCheck/IllusionsCard/Tools"
	"bytes"
	"errors"
	"os"
	"strings"
)

var Parameter = "Parameter"
var KKEx = "KKEx"

// ReadCardKK 读取KK的卡片,传入卡7片路径
func ReadCardKK(path string) (*KK.KoiCard, error) {
	path = strings.Replace(path, "\\", "/", -1)
	//fmt.Println(path)
	//提取文件名
	//读取图片
	f, err := os.ReadFile(path)
	buffer := bytes.NewBuffer(f)
	if err != nil {
		return nil, err
	}
	//value := get_png(f)
	//fmt.Println(value)
	//切割图片
	pngend := bytes.Index(f, Base.PngEndChunk) + len(Base.PngEndChunk)
	if pngend == -1 {
		return nil, errors.New("未找到PNG文件块尾")
	}
	outpng := buffer.Next(pngend)
	//os.WriteFile("Out.png", outpng, 0776)
	fb, err := buffer.ReadByte()
	if err != nil {
		return nil, errors.New("读取失败：卡片数据头无法识别")
	}
	if fb == 0x7 {
		_, err = Tools.BufRead(buffer, 64, "读取失败：数据缺失")
		if err != nil {
			return nil, err
		}
	} else if fb == 0x64 {
		_, err = Tools.BufRead(buffer, 3, "读取失败：数据缺失")
		if err != nil {
			return nil, err
		}
	}
	typelen, err := buffer.ReadByte()
	if err != nil {
		return nil, errors.New("读取失败：卡片类型无法识别")
	}
	cardtypebyte, err := Tools.BufRead(buffer, int(typelen), "读取失败：卡片类型无法识别")
	if err != nil {
		return nil, err
	}
	cardtype := string(cardtypebyte)
	//fmt.Println("卡片类型:", cardtype)
	//fmt.Println("封面大小:", pngend)
	card, err := KK.ParseKoiChara(buffer)
	if err != nil {
		return nil, err
	}
	//IllusionsCard.PrintCardInfo()
	card.Image = outpng
	card.Path = path
	card.CardType = cardtype
	return &card, nil
	//版本号
}
