package content_page

import (
	"KKCardModCheck/util"
	"fmt"
	"strings"

	g "github.com/AllenDang/giu"
)

var classicVersionState struct {
	current  string
	latest   string
	outdated bool
	checking bool
	failed   bool
	retry    func()
}

func SetClassicVersionState(current, latest string, outdated, checking, failed bool, retry func()) {
	classicVersionState.current = current
	classicVersionState.latest = latest
	classicVersionState.outdated = outdated
	classicVersionState.checking = checking
	classicVersionState.failed = failed
	classicVersionState.retry = retry
}

func classicSidebarWithUpdate() g.Widget {
	return g.Child().ID(g.ID("classic-sidebar-update")).Size(205, -1).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
		g.Label("KKCardModCheck"),
		g.Label("Version "+classicVersionState.current),
		g.Dummy(0, 6),
		classicNav("卡片检测", pageCardCheck),
		classicNav("场景卡解析", pageSceneCard),
		classicNav("卡片与 MOD 配对", pageModMatch),
		classicNav("卡片使用 MOD", pageCardUse),
		classicNav("MOD 信息构建", pageModBuild),
		classicNav("设置", pageSettings),
		g.Custom(func() {
			_, height := g.GetAvailableRegion()
			if height > 180 {
				g.Dummy(0, height-180).Build()
			}
		}),
		classicUpdateStatusWidget(),
		g.Separator(),
		g.Label("MOD 数据库"),
		g.Style().SetColor(g.StyleColorText, databaseStatusColor()).To(g.Label(modDatabaseSummary())),
	)
}

func classicUpdateStatusWidget() g.Widget {
	widgets := g.Layout{g.Label("版本检查")}
	current := strings.ToLower(strings.TrimSpace(classicVersionState.current))
	switch {
	case current == "" || current == "dev" || current == "development" || current == "unknown":
		widgets = append(widgets, g.Label("开发构建，不检查更新"))
	case classicVersionState.checking:
		widgets = append(widgets, g.Label("正在检查..."))
	case classicVersionState.outdated:
		widgets = append(widgets,
			g.Style().SetColor(g.StyleColorText, uiWarning).To(
				g.Label(fmt.Sprintf("发现新版本: %s", classicVersionState.latest)),
			),
			classicButton("打开发布页").OnClick(func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/releases")
			}),
		)
	case classicVersionState.failed:
		widgets = append(widgets,
			g.Style().SetColor(g.StyleColorText, uiDanger).To(g.Label("检查失败")),
			classicButton("重新检查").OnClick(runClassicVersionRetry),
		)
	case classicVersionState.latest != "":
		widgets = append(widgets,
			g.Label("当前已是最新版本"),
			classicButton("再次检查").OnClick(runClassicVersionRetry),
		)
	default:
		widgets = append(widgets, classicButton("检查更新").OnClick(runClassicVersionRetry))
	}
	return g.Child().ID(g.ID("classic-update-status")).Size(-1, 102).Border(false).Flags(g.WindowFlagsNoScrollbar).Layout(widgets...)
}

func runClassicVersionRetry() {
	if classicVersionState.retry != nil {
		classicVersionState.retry()
	}
}
