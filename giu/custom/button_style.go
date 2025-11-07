package custom

import (
	g "github.com/AllenDang/giu"
	colorful "github.com/lucasb-eyer/go-colorful"
)

// 此文件保留自定义按钮样式的代码

// 将 ButtonStyleBlue 暴露为包级全局变量，供其它包复用
var ButtonStyleBlue = g.Style()
var ButtonStyleRed = g.Style()

func init() {
	blueColor, _ := colorful.Hex("#007ACC") // VS Code 蓝
	redColor, _ := colorful.Hex("#FF3C00")  // VS Code 红
	ButtonStyleBlue.SetColor(g.StyleColorButton, blueColor)
	ButtonStyleRed.SetColor(g.StyleColorButton, redColor)
}
