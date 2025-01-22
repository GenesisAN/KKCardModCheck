package util

import (
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func OKMsg(pages *tview.Pages, text string, page string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage(page)
				}
			}),
		true,
		false)
	pages.SwitchToPage("弹窗")
	//app.Draw()
}

func MsgWeb(pages *tview.Pages, text, page, button, url string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text).
			AddButtons([]string{"OK", button}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage(page)
				} else {
					exec.Command("cmd", "/C", "start", url).Start()
				}
			}),
		true,
		false)
	pages.SwitchToPage("弹窗")
}

func MsgTips(pages *tview.Pages, text string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text),
		true,
		false)
	pages.SwitchToPage("弹窗")
}
