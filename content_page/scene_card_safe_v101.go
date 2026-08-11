package content_page

import (
	"fmt"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
)

func parseSelectedSceneSafelyV101() {
	defer func() {
		if recovered := recover(); recovered != nil {
			resetSceneResult()
			sceneStatus = fmt.Sprintf("场景卡解析失败：场景数据异常（%v）", recovered)
			g.Update()
		}
	}()
	parseSelectedScene()
}

func ClassicSceneCardToolsPageV101() g.Layout {
	if sceneCardPath == "" && strings.EqualFold(filepath.Ext(configCardPathForScene()), ".png") {
		sceneCardPath = configCardPathForScene()
	}
	if parsedScenePath != "" && cleanScenePath(sceneCardPath) != parsedScenePath {
		resetSceneResult()
	}

	widgets := g.Layout{
		g.Label("场景卡文件与输出目录"),
		g.Separator(),
		classicSceneFilePicker(),
		classicSceneOutputPicker(),
		g.Checkbox("按角色去重（推荐；取消后导出全部 Studio 角色实例）", &sceneUniqueOnly),
		g.Row(
			g.Button("解析场景卡").OnClick(parseSelectedSceneSafelyV101),
			g.Style().SetDisabled(!sceneReadyForExtraction()).To(
				g.Button("拆分角色卡").OnClick(extractSelectedSceneCharactersSafe),
			),
			g.Style().SetDisabled(sceneOutputDir == "").To(
				g.Button("打开输出目录").OnClick(func() {
					if sceneOutputDir != "" {
						openSceneOutputDir()
					}
				}),
			),
			g.Button("清除结果").OnClick(func() {
				resetSceneResult()
				sceneExportLog = ""
				g.Update()
			}),
		),
		g.Separator(),
	}
	if sceneStatus != "" {
		widgets = append(widgets, g.Label(sceneStatus))
	}
	if parsedScene != nil {
		widgets = append(widgets,
			g.TabBar().TabItems(
				sceneTextTab("场景概要", &sceneSummaryText),
				sceneTextTab("综合缺失 MOD", &sceneMissingModsText),
				sceneTextTab("综合使用 MOD", &sceneUsedModsText),
				sceneTextTab("场景插件 / MOD", &scenePluginsText),
				sceneTextTab("内嵌角色 / MOD", &sceneCharactersText),
				sceneTextTab("拆分日志", &sceneExportLog),
			),
		)
	}
	return widgets
}
