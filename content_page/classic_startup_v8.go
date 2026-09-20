package content_page

import (
	"KKCardModCheck/config"

	g "github.com/AllenDang/giu"
)

func InitialConfigWindowClassicV8() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{classicTheme(classicWelcomeV4()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return g.Layout{classicTheme(classicModBuildV4()...)}
	}
	return g.Layout{classicTheme(
		g.Row(
			classicSidebarWithUpdate(),
			g.Style().SetStyle(g.StyleVarWindowPadding, 0, 0).To(
				g.Child().ID(g.ID("classic-workspace-v8")).Size(-1, -1).Border(false).Layout(
					g.Style().SetStyle(g.StyleVarWindowPadding, 12, 12).To(classicActivePageV8()...),
				),
			),
		),
	)}
}

func classicActivePageV8() g.Layout {
	switch activePage {
	case pageSceneCard:
		return classicPageShell("场景卡解析 / 角色拆分", ClassicSceneCardToolsPageV101())
	case pageModMatch:
		return classicPageShell("卡片与 MOD 配对", ModeMatchPage())
	case pageCardUse:
		return classicPageShell("卡片使用 MOD", ModeCardUsePage())
	case pageModBuild:
		return classicPageShell("MOD 信息构建", classicModBuildV4())
	case pageSettings:
		return classicSettingsV4()
	default:
		return classicPageShell("卡片 MOD 缺失检测", ClassicCardCheckPageV101())
	}
}
