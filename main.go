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

func loop() {
	latestVersion, outdated, checking, failed := versionSnapshot()
	content_page.SetClassicVersionState(version, latestVersion, outdated, checking, failed, startVersionCheck)
	g.SingleWindow().Layout(content_page.InitialConfigWindowClassicV8())
}

var wnd *g.MasterWindow

func main() {
	defer func() {
		log.Println("Window has closed. Saving configuration...")
		config.Save()
	}()

	wnd = g.NewMasterWindow("KKCardModCheck V"+version, 1200, 800, 0)
	content_page.InitializeUIFonts()
	wnd.SetDropCallback(content_page.HandleDroppedFiles)
	wnd.SetSizeLimits(980, 680, 3840, 2160)
	wnd.Run(loop)
}
