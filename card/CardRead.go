package card

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

type BlockHeaderInfo struct {
	Name    string `json:"name"`    //名称
	Version string `json:"version"` //版本
	Pos     int64  `json:"pos"`     //坐标
	Size    int64  `json:"size"`    //大小
}

var Parameter = "Parameter"
var KKEx = "KKEx"

// ReadCardKK 读取KK的卡片,传入卡7片路径
func ReadCardKK(path string) (*KoiCard, error) {
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
	pngend := bytes.Index(f, PngEndChunk) + len(PngEndChunk)
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
		_, err = BufRead(buffer, 64, "读取失败：数据缺失")
		if err != nil {
			return nil, err
		}
	} else if fb == 0x64 {
		_, err = BufRead(buffer, 3, "读取失败：数据缺失")
		if err != nil {
			return nil, err
		}
	}
	typelen, err := buffer.ReadByte()
	if err != nil {
		return nil, errors.New("读取失败：卡片类型无法识别")
	}
	cardtypebyte, err := BufRead(buffer, int(typelen), "读取失败：卡片类型无法识别")
	if err != nil {
		return nil, err
	}
	cardtype := string(cardtypebyte)
	//fmt.Println("卡片类型:", cardtype)
	//fmt.Println("封面大小:", pngend)
	card, err := ParseKoiChara(buffer)
	if err != nil {
		return nil, err
	}
	//card.PrintCardInfo()
	card.Image = outpng
	card.Path = path
	card.CardType = cardtype
	return &card, nil
	//版本号
}

func BufRead(buffer *bytes.Buffer, n int, errMsg string) ([]byte, error) {
	if buffer.Len() < n {
		return nil, errors.New(errMsg)
	}
	return buffer.Next(n), nil
}

// 加载卡片头部
func LoadHead() {

}
