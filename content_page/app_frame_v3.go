package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"image/color"
	"strconv"

	g "github.com/AllenDang/giu"
)

var uiBoldFontsV3 = map[float32]*g.FontInfo{}

func InitializeUIFontsV3() {
	for _, size := range []float32{18, 20, 22, 26, 30, 32} {
		if font := g.Context.FontAtlas.AddFont("MSYHBD", size); font != nil {
			uiBoldFontsV3[size] = font
		}
	}
}

func boldTextV3(text string, size float32, col color.Color) g.Widget {
	s := g.Style().SetColor(g.StyleColorText, col)
	if f := uiBoldFontsV3[size]; f != nil {
		s.SetFont(f)
	}
	return s.To(g.Label(text))
}

func InitialConfigWindowV3() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{appTheme(Welcome()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return InitialConfigWindowV2()
	}
	return g.Layout{appTheme(
		brandBarV3(),
		g.Row(
			sidebarV3(),
			g.Style().SetColor(g.StyleColorChildBg, uiBackground).SetStyle(g.StyleVarWindowPadding, 26, 22).
				To(g.Child().ID(g.ID("workspace-v3")).Size(-1, -1).Border(false).Flags(g.WindowFlagsNoScrollbar).Layout(activePageLayoutV3()...)),
		),
	)}
}

func brandBarV3() g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, color.RGBA{R: 0x08, G: 0x11, B: 0x20, A: 0xFF}).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 22, 10).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("brand-bar-v3")).Size(-1, 58).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			g.Row(
				iconV3("brand-logo-v3", "logo", 34, uiPrimary),
				boldTextV3("KKCardModCheck", 22, uiText),
				textLabel("v0.1.10", 16, uiMuted),
			),
		))
}

func sidebarV3() g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, uiSidebar).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 0, 28).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("sidebar-v3")).Size(284, -1).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			navItemV3("document", "卡片检测", pageCardCheck),
			navItemV3("scene", "场景卡解析", pageSceneCard),
			navItemV3("mods", "卡片与Mod配对", pageModMatch),
			navItemV3("list", "卡片使用Mod", pageCardUse),
			navItemV3("database", "Mod信息构建", pageModBuild),
			navItemV3("settings", "设置", pageSettings),
			g.Custom(func() {
				_, h := g.GetAvailableRegion()
				if h > 132 {
					g.Dummy(0, h-132).Build()
				}
			}),
			databaseCardV3(),
		))
}

func navItemV3(kind, label string, page appPage) g.Widget {
	selected := activePage == page
	bg := uiSidebar
	col := uiMuted
	if selected {
		bg = color.RGBA{R: 0x18, G: 0x2A, B: 0x41, A: 0xFF}
		col = uiPrimary
	}
	activate := func() { activePage = page; g.Update() }
	return g.Style().SetColor(g.StyleColorChildBg, bg).SetStyle(g.StyleVarWindowPadding, 24, 10).
		To(g.Child().ID(g.ID("nav-v3-"+strconv.Itoa(int(page)))).Size(282, 64).Border(false).Flags(g.WindowFlagsNoScrollbar).Layout(
			g.Row(
				clickableIconV3("nav-icon-v3-"+strconv.Itoa(int(page)), kind, 36, col, activate),
				g.Style().SetColor(g.StyleColorText, col).SetColor(g.StyleColorHeader, uiTransparent).SetColor(g.StyleColorHeaderHovered, uiPanelAlt).
					SetColor(g.StyleColorHeaderActive, uiPrimaryLo).SetStyle(g.StyleVarFramePadding, 6, 9).
					To(g.Selectable(label+"##nav-label-v3-"+strconv.Itoa(int(page))).Selected(false).Size(184, 40).OnClick(activate)),
			),
		))
}

func databaseCardV3() g.Widget {
	return g.Style().SetColor(g.StyleColorChildBg, uiPanelAlt).SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 16, 14).SetStyleFloat(g.StyleVarChildRounding, 7).SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("database-card-v3")).Size(240, 108).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			g.Row(
				iconV3("database-card-icon-v3", "database", 38, uiMuted),
				g.Layout{sectionLabel("MOD 数据库"), textLabel(modDatabaseSummary(), 15, databaseStatusColor())},
			),
		))
}

func activePageLayoutV3() g.Layout {
	if activePage != pageCardCheck {
		return activePageLayoutV2()
	}
	return pageShellV3("卡片 MOD 缺失检测", "支持角色卡与 KK 场景卡，并综合检查场景及内嵌角色依赖", CardModVerificationPageV3())
}

func pageShellV3(title, subtitle string, body g.Layout) g.Layout {
	result := g.Layout{
		g.Row(
			g.Layout{boldTextV3(title, 32, uiText), textLabel(subtitle, 18, uiMuted)},
			g.Spacing(),
			g.Style().SetColor(g.StyleColorText, uiPrimary).SetColor(g.StyleColorButton, uiTransparent).SetColor(g.StyleColorButtonHovered, uiPanelAlt).
				SetColor(g.StyleColorBorder, uiPrimary).SetStyleFloat(g.StyleVarFrameBorderSize, 1).SetStyleFloat(g.StyleVarFrameRounding, 7).
				To(g.Button("使用说明").Size(128, 46).OnClick(func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck#使用说明") })),
		),
		g.Dummy(0, 8),
	}
	return append(result, body...)
}
