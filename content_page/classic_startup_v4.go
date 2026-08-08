package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"

	g "github.com/AllenDang/giu"
)

func InitialConfigWindowClassicV4() g.Layout {
	if !config.Instance.InitConfig {
		return g.Layout{classicTheme(classicWelcomeV4()...)}
	}
	if !config.Instance.HideModeInfoBuild {
		return g.Layout{classicTheme(classicModBuildV4()...)}
	}
	return g.Layout{classicTheme(
		g.Row(
			classicSidebar(),
			g.Child().ID(g.ID("classic-workspace-v4")).Size(-1, -1).Border(false).Layout(classicActivePageV4()...),
		),
	)}
}

func classicWelcomeV4() g.Layout {
	return g.Layout{
		g.Label("KKCardModCheck - 初次运行配置"),
		g.Separator(),
		g.Label("在开始检测角色卡之前，需要设置游戏目录并构建本地 MOD 信息数据库。").Wrapped(true),
		g.Dummy(0, 8),
		g.Child().ID(g.ID("classic-welcome-v4")).Size(-1, 180).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			g.Label("步骤 1 / 2：设置游戏目录"),
			g.Separator(),
			util.FilePicker("游戏目录", &config.Instance.GamePath).ForDirectories().LabelWidth(100).ButtonWidth(100),
			g.Label("请选择 Koikatsu / Koikatsu Sunshine 安装目录，也可以直接选择 Koikatu_Data。").Wrapped(true),
			classicWelcomeStatusV4(),
		),
		g.Dummy(0, 6),
		g.Row(
			g.Style().SetDisabled(config.Instance.GamePath == "").To(g.Button("保存并进入下一步").OnClick(func() {
				config.Instance.InitConfig = true
				config.Instance.HideModeInfoBuild = false
				config.Save()
				g.Update()
			})),
			g.Button("BUG 反馈").OnClick(func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/issues") }),
			g.Button("GitHub").OnClick(func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck") }),
		),
	}
}

func classicWelcomeStatusV4() g.Widget {
	if config.Instance.GamePath == "" {
		return g.Label("状态：尚未选择目录")
	}
	if isKoikatsuGameDirectory(config.Instance.GamePath) {
		return g.Style().SetColor(g.StyleColorText, uiSuccess).To(g.Label("状态：已识别游戏目录"))
	}
	return g.Style().SetColor(g.StyleColorText, uiWarning).To(g.Label("状态：未发现 Koikatu_Data，将使用当前目录扫描 MOD"))
}

func classicModBuildV4() g.Layout {
	actions := []g.Widget{
		g.Button(func() string {
			if modBuildInProgress {
				return "正在构建..."
			}
			if config.Instance.ModInfoBuild {
				return "重新构建"
			}
			return "开始构建"
		}()).Disabled(modBuildInProgress).OnClick(func() {
			if !modBuildInProgress {
				startModBuild()
			}
		}),
		g.Button("返回上一步").Disabled(modBuildInProgress).OnClick(func() {
			config.Instance.InitConfig = false
			config.Instance.HideModeInfoBuild = false
			config.Save()
			g.Update()
		}),
	}
	if len(modBuildFailedFiles) > 0 {
		actions = append(actions, g.Button("移动损坏 MOD").Disabled(movingBadMods).OnClick(func() {
			if !movingBadMods {
				moveFailedMods()
			}
		}))
	}
	if config.Instance.ModInfoBuild {
		actions = append(actions, g.Button("进入主界面").OnClick(func() {
			config.Instance.HideModeInfoBuild = true
			config.Save()
			g.Update()
		}))
	}

	return g.Layout{
		g.Label("KKCardModCheck - 初次运行配置"),
		g.Separator(),
		g.Label("步骤 2 / 2：构建 MOD 信息数据库"),
		g.Label("扫描游戏目录中的 .zipmod 文件并生成本地索引。添加新 MOD 后可以在主界面中重新构建。").Wrapped(true),
		g.Dummy(0, 6),
		g.Row(actions...),
		g.ProgressBar(modBuildProgress).Size(-1, 0),
		g.Label(func() string {
			if modBuildStatus == "" {
				return "状态：等待构建"
			}
			return "状态：" + modBuildStatus
		}()),
		g.InputTextMultiline(&modBuildLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsAllowTabInput | g.InputTextFlagsNoHorizontalScroll),
	}
}

func classicActivePageV4() g.Layout {
	switch activePage {
	case pageSceneCard:
		return classicPageShell("场景卡解析 / 角色拆分", SceneCardToolsPage())
	case pageModMatch:
		return classicPageShell("卡片与 MOD 配对", ModeMatchPage())
	case pageCardUse:
		return classicPageShell("卡片使用 MOD", ModeCardUsePage())
	case pageModBuild:
		return classicPageShell("MOD 信息构建", classicModBuildV4())
	case pageSettings:
		return classicSettingsV4()
	default:
		return classicPageShell("卡片 MOD 缺失检测", ClassicCardCheckPage())
	}
}

func classicSettingsV4() g.Layout {
	cfgPath := config.GetConfigPath()
	return classicPageShell("设置", g.Layout{
		g.Label("路径"),
		g.Separator(),
		util.FilePicker("游戏目录", &config.Instance.GamePath).ForDirectories().LabelWidth(100).ButtonWidth(100),
		g.Row(g.Label("配置文件"), g.InputText(&cfgPath).Flags(g.InputTextFlagsReadOnly)),
		g.Dummy(0, 6),
		g.Label("操作"),
		g.Separator(),
		g.Row(
			g.Button("打开配置文件").OnClick(func() { util.OpenFile(config.GetConfigPath()) }),
			g.Button("GitHub").OnClick(func() { util.OpenURL("https://github.com/GenesisAN/KKCardModCheck") }),
			g.Button("重置首次运行流程").OnClick(func() {
				config.Instance.InitConfig = false
				config.Instance.HideModeInfoBuild = false
				config.Save()
				g.Update()
			}),
		),
		g.Label(fmt.Sprintf("数据库：%s", modDatabaseSummary())),
	})
}
