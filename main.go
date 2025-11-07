package main

import (
	"KKCardModCheck/config"
	"KKCardModCheck/content_page"
	_ "embed"
	"log"

	g "github.com/AllenDang/giu"
)

//go:embed style.css
var cssStyle []byte

// 主循环
func loop() {
	g.SingleWindowWithMenuBar().Layout(
		AppMenuBar(),
		content_page.InitialConfigWindow(),
	)
}

var wnd *g.MasterWindow

func main() {
	defer func() {
		log.Println("Window has closed. Saving configuration...")
		config.Save()
	}()
	wnd = g.NewMasterWindow("KKCardModCheck V"+version, 1200, 900, 0)
	if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
		panic(err)
	}
	wnd.Run(loop)
}
