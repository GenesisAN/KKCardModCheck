package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	g "github.com/AllenDang/giu"
)

var (
	classicBg       = color.RGBA{R: 15, G: 15, B: 15, A: 255}
	classicPanel    = color.RGBA{R: 25, G: 25, B: 25, A: 255}
	classicFrame    = color.RGBA{R: 41, G: 74, B: 122, A: 138}
	classicFrameHot = color.RGBA{R: 66, G: 150, B: 250, A: 102}
	classicBlue     = color.RGBA{R: 66, G: 150, B: 250, A: 255}
	classicText     = color.RGBA{R: 242, G: 242, B: 242, A: 255}
	classicMuted    = color.RGBA{R: 153, G: 153, B: 153, A: 255}
	classicBorder   = color.RGBA{R: 110, G: 110, B: 128, A: 128}
)

func classicTheme(widgets ...g.Widget) g.Widget {
	return g.Style().
		SetColor(g.StyleColorText, classicText).
		SetColor(g.StyleColorTextDisabled, classicMuted).
		SetColor(g.StyleColorWindowBg, classicBg).
		SetColor(g.StyleColorChildBg, classicPanel).
		SetColor(g.StyleColorPopupBg, classicPanel).
		SetColor(g.StyleColorBorder, classicBorder).
		SetColor(g.StyleColorBorderShadow, uiTransparent).
		SetColor(g.StyleColorFrameBg, classicFrame).
		SetColor(g.StyleColorFrameBgHovered, classicFrameHot).
		SetColor(g.StyleColorFrameBgActive, classicFrameHot).
		SetColor(g.StyleColorTitleBg, classicPanel).
		SetColor(g.StyleColorTitleBgActive, classicPanel).
		SetColor(g.StyleColorMenuBarBg, classicPanel).
		SetColor(g.StyleColorScrollbarBg, classicBg).
		SetColor(g.StyleColorScrollbarGrab, color.RGBA{R: 79, G: 79, B: 79, A: 255}).
		SetColor(g.StyleColorCheckMark, classicBlue).
		SetColor(g.StyleColorButton, color.RGBA{R: 66, G: 150, B: 250, A: 102}).
		SetColor(g.StyleColorButtonHovered, color.RGBA{R: 66, G: 150, B: 250, A: 255}).
		SetColor(g.StyleColorButtonActive, color.RGBA{R: 15, G: 135, B: 250, A: 255}).
		SetColor(g.StyleColorHeader, color.RGBA{R: 66, G: 150, B: 250, A: 80}).
		SetColor(g.StyleColorHeaderHovered, color.RGBA{R: 66, G: 150, B: 250, A: 204}).
		SetColor(g.StyleColorHeaderActive, color.RGBA{R: 66, G: 150, B: 250, A: 255}).
		SetColor(g.StyleColorSeparator, classicBorder).
		SetColor(g.StyleColorTab, classicPanel).
		SetColor(g.StyleColorTabHovered, classicBlue).
		SetColor(g.StyleColorTabActive, color.RGBA{R: 51, G: 105, B: 173, A: 255}).
		SetColor(g.StyleColorTableHeaderBg, color.RGBA{R: 48, G: 48, B: 48, A: 255}).
		SetColor(g.StyleColorTableBorderStrong, classicBorder).
		SetColor(g.StyleColorTableBorderLight, color.RGBA{R: 59, G: 59, B: 64, A: 255}).
		SetColor(g.StyleColorTableRowBg, uiTransparent).
		SetColor(g.StyleColorTableRowBgAlt, color.RGBA{R: 255, G: 255, B: 255, A: 15}).
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

func InitialConfigWindowClassic() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{classicTheme(Welcome()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return g.Layout{classicTheme(
			g.Child().ID(g.ID("classic-first-build")).Size(-1, -1).Border(false).Layout(
				g.Label("构建 MOD 信息数据库"),
				g.Separator(),
				ModInfoBuild(),
			),
		)}
	}

	return g.Layout{classicTheme(
		g.Row(
			classicSidebar(),
			g.Child().ID(g.ID("classic-workspace")).Size(-1, -1).Border(false).Layout(classicActivePage()...),
		),
	)}
}

func classicSidebar() g.Widget {
	return g.Child().ID(g.ID("classic-sidebar")).Size(205, -1).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
		g.Label("KKCardModCheck"),
		g.Label("Version 0.1.10"),
		g.Separator(),
		classicNav("卡片检测", pageCardCheck),
		classicNav("场景卡解析", pageSceneCard),
		classicNav("卡片与 MOD 配对", pageModMatch),
		classicNav("卡片使用 MOD", pageCardUse),
		classicNav("MOD 信息构建", pageModBuild),
		classicNav("设置", pageSettings),
		g.Custom(func() {
			_, h := g.GetAvailableRegion()
			if h > 72 {
				g.Dummy(0, h-72).Build()
			}
		}),
		g.Separator(),
		g.Label("MOD 数据库"),
		g.Style().SetColor(g.StyleColorText, databaseStatusColor()).To(g.Label(modDatabaseSummary())),
	)
}

func classicNav(label string, page appPage) g.Widget {
	selected := activePage == page
	return g.Selectable(label+"##classic-nav-"+strconv.Itoa(int(page))).Selected(selected).Size(187, 34).OnClick(func() {
		activePage = page
		g.Update()
	})
}

func classicActivePage() g.Layout {
	switch activePage {
	case pageSceneCard:
		return classicPageShell("场景卡解析 / 角色拆分", SceneCardToolsPage())
	case pageModMatch:
		return classicPageShell("卡片与 MOD 配对", ModeMatchPage())
	case pageCardUse:
		return classicPageShell("卡片使用 MOD", ModeCardUsePage())
	case pageModBuild:
		return classicPageShell("MOD 信息构建", ModInfoBuild())
	case pageSettings:
		return settingsPage()
	default:
		return classicPageShell("卡片 MOD 缺失检测", ClassicCardCheckPage())
	}
}

func classicPageShell(title string, body g.Layout) g.Layout {
	layout := g.Layout{
		g.Row(
			textLabel(title, 24, classicText),
			g.Spacing(),
			g.Button("使用说明").OnClick(func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck#使用说明")
			}),
		),
		g.Separator(),
	}
	return append(layout, body...)
}

func ClassicCardCheckPage() g.Layout {
	syncCardSelection()
	return g.Layout{
		g.Custom(func() {
			width, _ := g.GetAvailableRegion()
			gap := float32(8)
			left := (width - gap) * 0.58
			g.Row(classicInputPanel(left), classicPreviewPanel(width-gap-left)).Build()
		}),
		classicResultsPanel(),
	}
}

func classicInputPanel(width float32) g.Widget {
	return g.Child().ID(g.ID("classic-input-panel")).Size(width, 245).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
		g.Label("输入"),
		g.Separator(),
		g.Label("角色卡、KK 场景卡或目录"),
		g.InputText(&config.Instance.CardPath).Size(-1).Flags(g.InputTextFlagsReadOnly),
		g.Row(
			g.Button("选择 PNG 文件").OnClick(chooseCardFile),
			g.Button("选择目录").OnClick(chooseCardDirectory),
		),
		g.Dummy(0, 8),
		g.Label("说明"),
		g.BulletText("可将 PNG 角色卡、KK 场景卡或目录直接拖入窗口"),
		g.BulletText("目录模式会递归扫描所有 PNG 文件"),
		g.BulletText("场景卡结果会综合场景本体与全部内嵌角色依赖"),
		g.BulletText("检测依赖前需要先构建本地 MOD 数据库"),
	)
}

func classicPreviewPanel(width float32) g.Widget {
	path := strings.Trim(strings.TrimSpace(config.Instance.CardPath), "\"")
	layout := g.Layout{g.Label("预览 / 操作"), g.Separator()}
	if isPreviewableCard(path) {
		layout = append(layout,
			g.Row(
				g.ImageWithFile(path).Size(100, 142),
				g.Layout{
					g.Label(filepath.Base(path)),
					g.Label("PNG 角色卡 / KK 场景卡"),
					g.Label("状态：已载入"),
				},
			),
		)
	} else if path != "" {
		layout = append(layout, g.Label("目录模式"), g.Label(filepath.Base(path)), g.Label("状态：已载入"), g.Dummy(0, 52))
	} else {
		layout = append(layout, g.Label("尚未选择卡片或目录"), g.Dummy(0, 94))
	}
	layout = append(layout, g.Style().SetDisabled(path == "").To(g.Button("开始检测").Size(-1, 34).OnClick(SingleCardCheck)))
	return g.Child().ID(g.ID("classic-preview-panel")).Size(width, 245).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(layout...)
}

func classicResultsPanel() g.Widget {
	return g.Child().ID(g.ID("classic-results-panel")).Size(-1, -1).Border(true).Layout(
		g.Row(
			g.Label("检测结果"),
			g.Separator(),
			g.Label(classicStatusText()),
			g.Spacing(),
			g.Button("复制当前列表").OnClick(copyCurrentResult),
			g.Button("清除").OnClick(clearCardResults),
		),
		g.Row(
			g.Label(fmt.Sprintf("解析卡片: %d", parsedCardCount())),
			g.Separator(),
			g.Style().SetColor(g.StyleColorText, uiDanger).To(g.Label(fmt.Sprintf("综合缺失 MOD: %d", len(missingMods)))),
			g.Separator(),
			g.Style().SetColor(g.StyleColorText, uiSuccess).To(g.Label(fmt.Sprintf("综合使用 MOD: %d", len(usedMods)))),
		),
		g.TabBar().TabItems(
			g.TabItem(fmt.Sprintf("综合缺失的 MOD (%d)", len(missingMods))).Layout(classicResultTable(true)),
			g.TabItem(fmt.Sprintf("综合使用的 MOD (%d)", len(usedMods))).Layout(classicResultTable(false)),
		),
	)
}

func classicStatusText() string {
	if checked {
		return "状态：检测完成"
	}
	if tips != "" {
		return "状态：" + tips
	}
	return "状态：等待检测"
}

func classicResultTable(showMissing bool) g.Widget {
	entries := usedMods
	if showMissing {
		entries = missingMods
	}
	rows := make([]*g.TableRowWidget, 0, len(entries)+1)
	if len(entries) == 0 {
		msg := "尚无结果"
		if checked {
			msg = "没有条目"
		}
		rows = append(rows, g.TableRow(g.Label(msg), g.Label("-"), g.Label("-")).MinHeight(28))
	} else {
		for i, entry := range entries {
			value := entry
			status := "使用"
			if showMissing {
				status = "缺失"
			}
			rows = append(rows, g.TableRow(
				g.Label(value),
				g.Label(status),
				g.Button("复制##classic-copy-"+strconv.Itoa(i)).OnClick(func() { _ = util.CopyToClipboard(value) }),
			).MinHeight(30))
		}
	}
	return g.Table().Columns(
		g.TableColumn("MOD GUID").Flags(g.TableColumnFlagsWidthStretch).InnerWidthOrWeight(0.72),
		g.TableColumn("状态").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(90),
		g.TableColumn("操作").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(80),
	).Rows(rows...).Size(-1, -1).Flags(g.TableFlagsRowBg | g.TableFlagsBorders | g.TableFlagsResizable | g.TableFlagsScrollY | g.TableFlagsSizingStretchProp)
}

func _classicFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
