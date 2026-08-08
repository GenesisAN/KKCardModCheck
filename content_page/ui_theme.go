package content_page

import (
	"image/color"

	g "github.com/AllenDang/giu"
)

var (
	uiBackground  = color.RGBA{R: 0x0F, G: 0x0F, B: 0x0F, A: 0xFF}
	uiSidebar     = color.RGBA{R: 0x19, G: 0x19, B: 0x19, A: 0xFF}
	uiPanel       = color.RGBA{R: 0x19, G: 0x19, B: 0x19, A: 0xFF}
	uiPanelAlt    = color.RGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xFF}
	uiInput       = color.RGBA{R: 0x29, G: 0x4A, B: 0x7A, A: 0x8A}
	uiBorder      = color.RGBA{R: 0x6E, G: 0x6E, B: 0x80, A: 0x80}
	uiText        = color.RGBA{R: 0xF2, G: 0xF2, B: 0xF2, A: 0xFF}
	uiMuted       = color.RGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
	uiPrimary     = color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0xFF}
	uiPrimaryHi   = color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0xFF}
	uiPrimaryLo   = color.RGBA{R: 0x33, G: 0x69, B: 0xAD, A: 0xFF}
	uiSuccess     = color.RGBA{R: 0x57, G: 0xC7, B: 0x85, A: 0xFF}
	uiWarning     = color.RGBA{R: 0xF0, G: 0xB8, B: 0x49, A: 0xFF}
	uiDanger      = color.RGBA{R: 0xF0, G: 0x69, B: 0x69, A: 0xFF}
	uiTransparent = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	uiFonts       = map[float32]*g.FontInfo{}
)

func InitializeUIFonts() {
	for _, size := range []float32{15, 16, 18, 20, 24, 26, 28, 36, 46, 48} {
		if font := g.Context.FontAtlas.AddFont("MSYH", size); font != nil {
			uiFonts[size] = font
		}
	}
}

func appTheme(widgets ...g.Widget) *g.StyleSetter {
	return g.Style().
		SetColor(g.StyleColorText, uiText).
		SetColor(g.StyleColorTextDisabled, uiMuted).
		SetColor(g.StyleColorWindowBg, uiBackground).
		SetColor(g.StyleColorChildBg, uiPanel).
		SetColor(g.StyleColorPopupBg, uiPanel).
		SetColor(g.StyleColorBorder, uiBorder).
		SetColor(g.StyleColorBorderShadow, uiTransparent).
		SetColor(g.StyleColorFrameBg, uiInput).
		SetColor(g.StyleColorFrameBgHovered, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0x66}).
		SetColor(g.StyleColorFrameBgActive, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0xAA}).
		SetColor(g.StyleColorButton, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0x66}).
		SetColor(g.StyleColorButtonHovered, uiPrimary).
		SetColor(g.StyleColorButtonActive, uiPrimaryLo).
		SetColor(g.StyleColorHeader, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0x50}).
		SetColor(g.StyleColorHeaderHovered, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0xCC}).
		SetColor(g.StyleColorHeaderActive, uiPrimary).
		SetColor(g.StyleColorSeparator, uiBorder).
		SetColor(g.StyleColorTab, uiPanel).
		SetColor(g.StyleColorTabHovered, uiPrimary).
		SetColor(g.StyleColorTabActive, uiPrimaryLo).
		SetColor(g.StyleColorTableHeaderBg, uiPanelAlt).
		SetColor(g.StyleColorTableBorderStrong, uiBorder).
		SetColor(g.StyleColorTableBorderLight, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 8, 8).
		SetStyle(g.StyleVarFramePadding, 6, 4).
		SetStyle(g.StyleVarItemSpacing, 8, 6).
		SetStyle(g.StyleVarItemInnerSpacing, 6, 4).
		SetStyleFloat(g.StyleVarWindowRounding, 0).
		SetStyleFloat(g.StyleVarChildRounding, 0).
		SetStyleFloat(g.StyleVarFrameRounding, 0).
		SetStyleFloat(g.StyleVarPopupRounding, 0).
		SetStyleFloat(g.StyleVarScrollbarRounding, 0).
		SetStyleFloat(g.StyleVarGrabRounding, 0).
		SetStyleFloat(g.StyleVarTabRounding, 0).
		To(widgets...)
}

func panelStyle(id string, width, height float32, widgets ...g.Widget) *g.StyleSetter {
	return g.Style().SetColor(g.StyleColorChildBg, uiPanel).SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 8, 8).SetStyleFloat(g.StyleVarChildRounding, 0).SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID(id)).Size(width, height).Border(true).Layout(widgets...))
}

func insetStyle(id string, width, height float32, widgets ...g.Widget) *g.StyleSetter {
	return g.Style().SetColor(g.StyleColorChildBg, uiPanelAlt).SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 8, 8).SetStyleFloat(g.StyleVarChildRounding, 0).SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID(id)).Size(width, height).Border(true).Layout(widgets...))
}

func primaryButton(label string, width, height float32, onClick func()) *g.StyleSetter {
	return g.Style().SetColor(g.StyleColorText, uiText).
		SetColor(g.StyleColorButton, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0xA0}).
		SetColor(g.StyleColorButtonHovered, uiPrimary).SetColor(g.StyleColorButtonActive, uiPrimaryLo).
		SetStyleFloat(g.StyleVarFrameRounding, 0).To(g.Button(label).Size(width, height).OnClick(onClick))
}

func secondaryButton(label string, width, height float32, onClick func()) *g.StyleSetter {
	return g.Style().SetColor(g.StyleColorButton, color.RGBA{R: 0x42, G: 0x96, B: 0xFA, A: 0x66}).
		SetColor(g.StyleColorButtonHovered, uiPrimary).SetColor(g.StyleColorButtonActive, uiPrimaryLo).
		SetColor(g.StyleColorBorder, uiBorder).SetStyleFloat(g.StyleVarFrameBorderSize, 1).SetStyleFloat(g.StyleVarFrameRounding, 0).
		To(g.Button(label).Size(width, height).OnClick(onClick))
}

func textLabel(text string, size float32, col color.Color) *g.StyleSetter {
	style := g.Style().SetColor(g.StyleColorText, col)
	if font := uiFonts[size]; font != nil {
		style.SetFont(font)
	}
	return style.To(g.Label(text))
}

func titleLabel(text string) *g.StyleSetter   { return textLabel(text, 24, uiText) }
func sectionLabel(text string) *g.StyleSetter { return textLabel(text, 18, uiText) }
func mutedLabel(text string) *g.StyleSetter   { return textLabel(text, 16, uiMuted) }

func statusLabel(text string, col color.Color) *g.StyleSetter {
	return g.Style().SetColor(g.StyleColorText, col).To(g.Label(text))
}
