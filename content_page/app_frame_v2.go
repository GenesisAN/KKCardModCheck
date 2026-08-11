package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"strconv"

	g "github.com/AllenDang/giu"
)

// InitialConfigWindowV2 是修正导航尺寸后的新版应用入口。
func InitialConfigWindowV2() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{appTheme(Welcome()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return g.Layout{appTheme(
			g.Style().
				SetColor(g.StyleColorChildBg, uiBackground).
				SetStyle(g.StyleVarWindowPadding, 24, 20).
				To(g.Child().ID(g.ID("first-build-v2")).Size(-1, -1).Border(false).Layout(
					titleLabel("构建 MOD 信息数据库"),
					mutedLabel("首次使用前扫描游戏目录，之后即可进行卡片依赖检查"),
					g.Dummy(0, 8),
					ModInfoBuild(),
				)),
		)}
	}
	return mainContentV2()
}

func mainContentV2() g.Layout {
	return g.Layout{appTheme(
		g.Row(
			sidebarV2(),
			g.Style().
				SetColor(g.StyleColorChildBg, uiBackground).
				SetStyle(g.StyleVarWindowPadding, 24, 20).
				To(g.Child().ID(g.ID("main-content-v2")).Size(-1, -1).Border(false).Layout(activePageLayoutV2()...)),
		),
	)}
}

func sidebarV2() g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, uiSidebar).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 18, 18).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("sidebar-v2")).Size(236, -1).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			textLabel("KKCardModCheck", 20, uiText),
			mutedLabel("角色卡 MOD 工具"),
			g.Dummy(0, 18),
			navItemV2("卡片检测", pageCardCheck),
			navItemV2("场景卡解析", pageSceneCard),
			navItemV2("卡片与 Mod 配对", pageModMatch),
			navItemV2("卡片使用 Mod", pageCardUse),
			navItemV2("Mod 信息构建", pageModBuild),
			navItemV2("设置", pageSettings),
			g.Custom(func() {
				_, availableHeight := g.GetAvailableRegion()
				if availableHeight > 132 {
					g.Dummy(0, availableHeight-132).Build()
				}
			}),
			insetStyle("database-status-v2", -1, 112,
				g.Row(textLabel("DB", 18, uiMuted), sectionLabel("MOD 数据库")),
				textLabel(modDatabaseSummary(), 15, databaseStatusColor()),
				mutedLabel("用于核对本地依赖"),
			),
		))
}

func navItemV2(label string, page appPage) g.Widget {
	selected := activePage == page
	visibleLabel := "    " + label + "##nav-v2-" + strconv.Itoa(int(page))
	textColor := uiMuted
	if selected {
		visibleLabel = ">   " + label + "##nav-v2-" + strconv.Itoa(int(page))
		textColor = uiPrimary
	}

	return g.Style().
		SetColor(g.StyleColorText, textColor).
		SetColor(g.StyleColorHeader, uiPrimaryLo).
		SetColor(g.StyleColorHeaderHovered, uiPanelAlt).
		SetColor(g.StyleColorHeaderActive, uiPrimaryLo).
		SetStyle(g.StyleVarFramePadding, 10, 10).
		To(g.Selectable(visibleLabel).Selected(selected).Size(196, 46).OnClick(func() {
			activePage = page
			g.Update()
		}))
}

func activePageLayoutV2() g.Layout {
	switch activePage {
	case pageSceneCard:
		return pageShellV2("场景卡解析 / 角色拆分", "解析 KK 工作室场景依赖，并导出内嵌角色卡", SceneCardToolsPage())
	case pageModMatch:
		return pageShellV2("卡片与 MOD 配对", "检查指定 Zipmod 是否被角色卡实际引用", ModeMatchPage())
	case pageCardUse:
		return pageShellV2("卡片使用 MOD", "根据角色卡生成便于分享的借物表", ModeCardUsePage())
	case pageModBuild:
		return pageShellV2("MOD 信息构建", "重新扫描游戏目录并更新本地 MOD 数据库", ModInfoBuild())
	case pageSettings:
		return settingsPage()
	default:
		return pageShellV2("卡片 MOD 缺失检测", "支持角色卡与 KK 场景卡，并综合检查场景及内嵌角色依赖", CardModVerificationPage())
	}
}

func pageShellV2(title, subtitle string, body g.Layout) g.Layout {
	layout := g.Layout{
		g.Row(
			titleLabel(title),
			g.Align(g.AlignRight).To(
				secondaryButton("使用说明", 110, 36, func() {
					util.OpenURL("https://github.com/GenesisAN/KKCardModCheck#使用说明")
				}),
			),
		),
		mutedLabel(subtitle),
		g.Dummy(0, 10),
	}
	return append(layout, body...)
}
