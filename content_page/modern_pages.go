package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
)

// SceneCardToolsPage 是场景卡工具在新侧栏框架中的页面布局。
func SceneCardToolsPage() g.Layout {
	if sceneCardPath == "" && strings.EqualFold(filepath.Ext(config.Instance.CardPath), ".png") {
		sceneCardPath = config.Instance.CardPath
	}
	if parsedScenePath != "" && cleanScenePath(sceneCardPath) != parsedScenePath {
		resetSceneResult()
	}

	widgets := []g.Widget{
		sectionLabel("选择场景卡"),
		mutedLabel("支持 KK 工作室 PNG 场景卡，可查看插件、Zipmod 与内嵌角色"),
		g.Dummy(0, 4),
		util.FilePicker("场景卡文件", &sceneCardPath).ForFiles().Filters("*.png").LabelWidth(120).ButtonWidth(130),
		util.FilePicker("角色输出目录", &sceneOutputDir).ForDirectories().LabelWidth(120).ButtonWidth(130),
		g.Checkbox("按角色去重（推荐；取消后导出全部 Studio 角色实例）", &sceneUniqueOnly),
		g.Row(
			primaryButton("解析场景卡", 150, 44, parseSelectedScene),
			secondaryButton("拆分角色卡", 150, 44, extractSelectedSceneCharacters),
			secondaryButton("打开输出目录", 150, 44, func() {
				if sceneOutputDir != "" {
					util.OpenFile(sceneOutputDir)
				}
			}),
			secondaryButton("清除结果", 120, 44, func() {
				resetSceneResult()
				sceneExportLog = ""
				g.Update()
			}),
		),
	}

	if strings.EqualFold(filepath.Ext(cleanScenePath(sceneCardPath)), ".png") {
		widgets = append(widgets,
			g.Align(g.AlignCenter).To(g.ImageWithFile(cleanScenePath(sceneCardPath)).Size(256, 144)),
		)
	}
	if sceneStatus != "" {
		widgets = append(widgets, textLabel(sceneStatus, 16, uiSuccess))
	}
	if parsedScene != nil {
		widgets = append(widgets,
			g.TabBar().TabItems(
				sceneTextTab(" 场景概要 ", &sceneSummaryText),
				sceneTextTab(" 综合缺失 Mod ", &sceneMissingModsText),
				sceneTextTab(" 综合使用 Mod ", &sceneUsedModsText),
				sceneTextTab(" 场景插件 / Mod ", &scenePluginsText),
				sceneTextTab(" 内嵌角色 / Mod ", &sceneCharactersText),
				sceneTextTab(" 拆分日志 ", &sceneExportLog),
			),
		)
	}

	return g.Layout{panelStyle("scene-card-page", -1, -1, widgets...)}
}

// ModeMatchPage 是卡片与 Mod 配对功能在新侧栏框架中的页面布局。
func ModeMatchPage() g.Layout {
	if mmCardPath == "" {
		mmCardPath = config.Instance.CardPath
	}
	status := "等待操作"
	statusColor := uiMuted
	if mmStatus != "" {
		status = mmStatus
		statusColor = uiSuccess
		if strings.Contains(mmStatus, "错误") {
			statusColor = uiDanger
		}
	}

	return g.Layout{panelStyle("mod-match-page", -1, -1,
		sectionLabel("选择 MOD 与角色卡"),
		mutedLabel("可拖入文件到窗口，或使用右侧按钮选择文件"),
		g.Dummy(0, 4),
		util.FilePicker("MOD 文件", &mmModPath).ForFiles().Filters("*.zipmod").LabelWidth(120).ButtonWidth(130),
		util.FilePicker("卡片文件", &mmCardPath).ForFiles().Filters("*.png").LabelWidth(120).ButtonWidth(130),
		g.Row(
			primaryButton(func() string {
				if mmInProgress {
					return "正在比对..."
				}
				return "开始比对"
			}(), 170, 48, func() {
				if !mmInProgress {
					startModMatch()
				}
			}),
			secondaryButton("清除", 110, 48, func() {
				mmModPath = ""
				mmCardPath = ""
				mmStatus = ""
				mmLog = ""
				g.Update()
			}),
		),
		g.Row(g.Label("状态"), textLabel(status, 16, statusColor)),
		g.InputTextMultiline(&mmLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
	)}
}

// ModeCardUsePage 是借物表生成功能在新侧栏框架中的页面布局。
func ModeCardUsePage() g.Layout {
	if cuCardPath == "" {
		cuCardPath = config.Instance.CardPath
	}
	status := "等待操作"
	statusColor := uiMuted
	if cuStatus != "" {
		status = cuStatus
		statusColor = uiSuccess
		if strings.Contains(cuStatus, "错误") {
			statusColor = uiDanger
		}
	}

	return g.Layout{panelStyle("card-use-page", -1, -1,
		sectionLabel("生成角色卡借物表"),
		mutedLabel("读取角色卡依赖并生成 card_use_mod.txt"),
		g.Dummy(0, 4),
		util.FilePicker("卡片文件", &cuCardPath).ForFiles().Filters("*.png").LabelWidth(120).ButtonWidth(130),
		g.Row(
			primaryButton(func() string {
				if cuInProgress {
					return "正在生成..."
				}
				return "生成借物表"
			}(), 170, 48, func() {
				if !cuInProgress {
					startGenerateCardMods()
				}
			}),
			secondaryButton("打开借物表", 150, 48, func() {
				if cuOutputPath != "" {
					util.OpenFile(cuOutputPath)
				}
			}),
			secondaryButton("清除", 110, 48, func() {
				cuCardPath = ""
				cuStatus = ""
				cuLog = ""
				cuOutputPath = ""
				g.Update()
			}),
		),
		g.Row(g.Label("状态"), textLabel(status, 16, statusColor)),
		g.InputTextMultiline(&cuLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
	)}
}
