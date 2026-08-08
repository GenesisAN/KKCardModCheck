package content_page

import g "github.com/AllenDang/giu"

// fixedPanelStyle 用于不应出现内部滚动条的固定高度卡片。
func fixedPanelStyle(id string, width, height float32, widgets ...g.Widget) *g.StyleSetter {
	return g.Style().
		SetColor(g.StyleColorChildBg, uiPanel).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 18, 16).
		SetStyleFloat(g.StyleVarChildRounding, 8).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID(id)).Size(width, height).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(widgets...))
}
