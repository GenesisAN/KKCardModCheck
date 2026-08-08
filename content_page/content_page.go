package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	g "github.com/AllenDang/giu"
)

// 软件所在目录
var appDir string = "."

type appPage int

const (
	pageCardCheck appPage = iota
	pageSceneCard
	pageModMatch
	pageCardUse
	pageModBuild
	pageSettings
)

var (
	activePage       = pageCardCheck
	dbSummaryPath    string
	dbSummaryModTime time.Time
	dbSummaryText    = "尚未构建"
)

func init() {
	appDir = filepath.Dir(os.Args[0])
}

// InitialConfigWindow 显示初始配置流程或主界面。
func InitialConfigWindow() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{appTheme(Welcome()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return g.Layout{appTheme(
			g.Style().
				SetColor(g.StyleColorChildBg, uiBackground).
				SetStyle(g.StyleVarWindowPadding, 24, 20).
				To(g.Child().ID(g.ID("first-build")).Size(-1, -1).Border(false).Layout(
					titleLabel("构建 MOD 信息数据库"),
					mutedLabel("首次使用前扫描游戏目录，之后即可进行卡片依赖检查"),
					g.Dummy(0, 8),
					ModInfoBuild(),
				)),
		)}
	}
	return MainContent()
}

// MainContent 是左侧导航加工作区的统一应用框架。
func MainContent() g.Layout {
	return g.Layout{
		appTheme(
			g.Row(
				sidebar(),
				g.Style().
					SetColor(g.StyleColorChildBg, uiBackground).
					SetStyle(g.StyleVarWindowPadding, 24, 20).
					To(g.Child().ID(g.ID("main-content")).Size(-1, -1).Border(false).Layout(activePageLayout()...)),
			),
		),
	}
}

func sidebar() g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, uiSidebar).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 18, 18).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("sidebar")).Size(236, -1).Border(true).Layout(
			g.Row(
				textLabel("◆", 24, uiPrimary),
				textLabel("KKCardModCheck", 20, uiText),
			),
			mutedLabel("角色卡 MOD 工具"),
			g.Dummy(0, 18),
			navItem("卡片检测", pageCardCheck),
			navItem("场景卡解析", pageSceneCard),
			navItem("卡片与 Mod 配对", pageModMatch),
			navItem("卡片使用 Mod", pageCardUse),
			navItem("Mod 信息构建", pageModBuild),
			navItem("设置", pageSettings),
			g.Custom(func() {
				_, availableHeight := g.GetAvailableRegion()
				if availableHeight > 132 {
					g.Dummy(0, availableHeight-132).Build()
				}
			}),
			insetStyle("database-status", -1, 112,
				g.Row(textLabel("▣", 24, uiMuted), sectionLabel("MOD 数据库")),
				textLabel(modDatabaseSummary(), 15, databaseStatusColor()),
				mutedLabel("用于核对本地依赖"),
			),
		))
}

func navItem(label string, page appPage) g.Widget {
	selected := activePage == page
	visibleLabel := "     " + label + "##nav-" + strconv.Itoa(int(page))
	textColor := uiMuted
	if selected {
		visibleLabel = "●   " + label + "##nav-" + strconv.Itoa(int(page))
		textColor = uiPrimary
	}

	return g.Style().
		SetColor(g.StyleColorText, textColor).
		SetColor(g.StyleColorHeader, uiPrimaryLo).
		SetColor(g.StyleColorHeaderHovered, uiPanelAlt).
		SetColor(g.StyleColorHeaderActive, uiPrimaryLo).
		SetStyle(g.StyleVarFramePadding, 10, 10).
		SetStyleFloat(g.StyleVarFrameRounding, 6).
		To(g.Selectable(visibleLabel).Selected(selected).Size(-1, 46).OnClick(func() {
			activePage = page
			g.Update()
		}))
}

func activePageLayout() g.Layout {
	switch activePage {
	case pageSceneCard:
		return pageShell("场景卡解析 / 角色拆分", "解析 KK 工作室场景依赖，并导出内嵌角色卡", SceneCardToolsPage())
	case pageModMatch:
		return pageShell("卡片与 MOD 配对", "检查指定 Zipmod 是否被角色卡实际引用", ModeMatchPage())
	case pageCardUse:
		return pageShell("卡片使用 MOD", "根据角色卡生成便于分享的借物表", ModeCardUsePage())
	case pageModBuild:
		return pageShell("MOD 信息构建", "重新扫描游戏目录并更新本地 MOD 数据库", ModInfoBuild())
	case pageSettings:
		return settingsPage()
	default:
		return pageShell("卡片 MOD 缺失检测", "选择角色卡或卡片目录，快速检查缺失依赖", CardModVerificationPage())
	}
}

func pageShell(title, subtitle string, body g.Layout) g.Layout {
	layout := g.Layout{
		g.Row(
			titleLabel(title),
			g.Spacing(),
			secondaryButton("使用说明", 110, 36, func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck#使用说明")
			}),
		),
		mutedLabel(subtitle),
		g.Dummy(0, 10),
	}
	return append(layout, body...)
}

func settingsPage() g.Layout {
	cfgPath := config.GetConfigPath()
	return pageShell("设置", "管理游戏目录、配置文件和支持入口", g.Layout{
		panelStyle("settings-main", -1, -1,
			sectionLabel("路径设置"),
			mutedLabel("设置 Koikatsu / Koikatsu Sunshine 游戏目录"),
			g.Dummy(0, 4),
			util.FilePicker("游戏目录", &config.Instance.GamePath).ForDirectories().LabelWidth(110).ButtonWidth(130),
			g.Dummy(0, 12),
			sectionLabel("软件与配置"),
			g.Row(g.Label("软件目录"), g.InputText(&appDir).Flags(g.InputTextFlagsReadOnly)),
			g.Row(g.Label("配置文件"), g.InputText(&cfgPath).Flags(g.InputTextFlagsReadOnly)),
			g.Row(
				secondaryButton("打开软件目录", 150, 38, func() { util.OpenURL(".") }),
				secondaryButton("打开配置文件", 150, 38, func() { util.OpenFile(config.GetConfigPath()) }),
				g.Style().SetColor(g.StyleColorButton, uiDanger).To(
					g.Button("重置配置").Size(120, 38).OnClick(config.ResetToDefault),
				),
			),
			g.Dummy(0, 12),
			sectionLabel("帮助与反馈"),
			g.Row(
				secondaryButton("BUG 反馈", 130, 38, func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/issues") }),
				secondaryButton("讨论区", 110, 38, func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/discussions") }),
				secondaryButton("GitHub 仓库", 140, 38, func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck") }),
				secondaryButton("QQ群", 100, 38, func() { util.OpenURL("https://qm.qq.com/q/rza05u1550") }),
			),
		),
	})
}

// HandleDroppedFiles 接收从 Windows 资源管理器拖入窗口的文件或目录。
func HandleDroppedFiles(names []string) {
	for _, name := range names {
		name = strings.Trim(strings.TrimSpace(name), "\"")
		info, err := os.Stat(name)
		if err != nil {
			continue
		}
		if info.IsDir() || strings.EqualFold(filepath.Ext(name), ".png") {
			config.Instance.CardPath = name
			activePage = pageCardCheck
			tips = "已载入：" + filepath.Base(name)
			g.Update()
			return
		}
	}
	tips = "拖入内容无效，请选择 PNG 角色卡或卡片目录"
	g.Update()
}

func modDatabaseSummary() string {
	path := config.Instance.ModInfoPath
	if path == "" {
		dbSummaryText = "尚未构建"
		return dbSummaryText
	}
	info, err := os.Stat(path)
	if err != nil {
		dbSummaryText = "数据库不可用"
		return dbSummaryText
	}
	if path == dbSummaryPath && info.ModTime().Equal(dbSummaryModTime) {
		return dbSummaryText
	}

	dbSummaryPath = path
	dbSummaryModTime = info.ModTime()
	data, err := os.ReadFile(path)
	if err != nil {
		dbSummaryText = "数据库不可用"
		return dbSummaryText
	}
	var entries []json.RawMessage
	if err := json.Unmarshal(data, &entries); err != nil {
		dbSummaryText = "数据库格式错误"
		return dbSummaryText
	}
	dbSummaryText = fmt.Sprintf("已就绪 · %s 项", formatThousands(len(entries)))
	return dbSummaryText
}

func databaseStatusColor() color.Color {
	if strings.HasPrefix(modDatabaseSummary(), "已就绪") {
		return uiSuccess
	}
	return uiWarning
}

func formatThousands(value int) string {
	s := strconv.Itoa(value)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
