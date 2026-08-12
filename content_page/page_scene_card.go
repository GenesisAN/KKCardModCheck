package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	g "github.com/AllenDang/giu"
	ic "github.com/GenesisAN/illusionsCard"
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
)

var (
	sceneCardPath        string
	sceneOutputDir       string
	sceneUniqueOnly      = true
	sceneStatus          string
	sceneSummaryText     string
	scenePluginsText     string
	sceneCharactersText  string
	sceneMissingModsText string
	sceneUsedModsText    string
	sceneExportLog       string
	parsedScene          *KK.KKSceneCard
	parsedScenePath      string
)

// SceneCardToolsTab displays KK Studio scene dependencies and extracts the
// embedded characters as standalone KK character-card PNGs.
func SceneCardToolsTab() *g.TabItemWidget {
	if sceneCardPath == "" && strings.EqualFold(filepath.Ext(config.Instance.CardPath), ".png") {
		sceneCardPath = config.Instance.CardPath
	}
	if parsedScenePath != "" && cleanScenePath(sceneCardPath) != parsedScenePath {
		resetSceneResult()
	}

	layout := []g.Widget{
		g.Label("解析 KK 工作室场景卡的插件、Zipmod 和内嵌角色"),
		g.Separator(),
		util.FilePicker("场景卡文件", &sceneCardPath).ForFiles().Filters("*.png"),
		util.FilePicker("角色输出目录", &sceneOutputDir).ForDirectories(),
		g.Checkbox("按角色去重（推荐；取消后导出全部 Studio 角色实例）", &sceneUniqueOnly),
		g.Row(
			g.Button("解析场景卡").Size(160, 48).OnClick(parseSelectedScene),
			g.Button("拆分角色卡").Size(160, 48).OnClick(extractSelectedSceneCharacters),
			g.Button("打开输出目录").Size(140, 48).OnClick(func() {
				if sceneOutputDir != "" {
					util.OpenFile(sceneOutputDir)
				}
			}),
			g.Button("清除结果").Size(120, 48).OnClick(func() {
				resetSceneResult()
				sceneExportLog = ""
				g.Update()
			}),
		),
		g.Separator(),
	}

	if strings.EqualFold(filepath.Ext(cleanScenePath(sceneCardPath)), ".png") {
		layout = append(layout,
			g.Align(g.AlignCenter).To(g.ImageWithFile(cleanScenePath(sceneCardPath)).Size(320, 180)),
			g.Separator(),
		)
	}
	if sceneStatus != "" {
		layout = append(layout, g.Label(sceneStatus))
	}
	if parsedScene != nil {
		layout = append(layout,
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

	return g.TabItem(" 场景卡解析 / 角色拆分 ").Layout(layout...)
}

func sceneTextTab(label string, text *string) *g.TabItemWidget {
	return g.TabItem(label).Layout(
		g.InputTextMultiline(text).Size(-1, -1).
			Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsAllowTabInput | g.InputTextFlagsNoHorizontalScroll),
	)
}

func parseSelectedScene() {
	path := cleanScenePath(sceneCardPath)
	if path == "" {
		sceneStatus = "请选择 KK 工作室场景卡"
		g.Update()
		return
	}
	card, err := ic.ReadCardFromPath(path)
	if err != nil {
		resetSceneResult()
		sceneStatus = fmt.Sprintf("场景卡读取失败：%v", err)
		g.Update()
		return
	}
	scene, ok := card.(*KK.KKSceneCard)
	if !ok {
		resetSceneResult()
		sceneStatus = fmt.Sprintf("所选文件不是 KK 工作室场景卡（识别为 %T）", card)
		g.Update()
		return
	}

	parsedScene = scene
	parsedScenePath = path
	if sceneOutputDir == "" {
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		sceneOutputDir = filepath.Join(filepath.Dir(path), base+"_characters")
	}
	buildSceneTexts(scene)
	missingCount, usedCount, modErr := buildSceneModTexts(scene)
	sceneStatus = fmt.Sprintf("解析完成：场景插件 %d 项，角色实例 %d 个，去重后 %d 个；综合使用 MOD %d 个",
		len(scene.SceneExtendedList), len(scene.CharaCards), len(scene.UniqueCharaCardIndexes()), usedCount)
	if modErr != nil {
		sceneStatus += fmt.Sprintf("；缺失 MOD 核对失败：%v", modErr)
	} else {
		sceneStatus += fmt.Sprintf("，综合缺失 %d 个", missingCount)
	}
	g.Update()
}

func buildSceneModTexts(scene Base.CardInterface) (missingCount, usedCount int, err error) {
	_, used := summarizeCardMods([]Base.CardInterface{scene}, nil)
	usedCount = len(used)
	sceneUsedModsText = fmt.Sprintf("综合使用 MOD（场景本体 + 全部内嵌角色，共 %d 项；包含缺失项）：\n%s",
		usedCount, formatStringList(used))

	localGUIDs, err := util.LoadModGUIDsFromJSON(config.Instance.ModInfoPath)
	if err != nil {
		sceneMissingModsText = fmt.Sprintf("无法读取本地 MOD 数据库，暂时不能核对缺失项：\n%v", err)
		return 0, usedCount, err
	}
	missing, _ := summarizeCardMods([]Base.CardInterface{scene}, localGUIDs)
	missingCount = len(missing)
	sceneMissingModsText = fmt.Sprintf("综合缺失 MOD（场景本体 + 全部内嵌角色，共 %d 项）：\n%s",
		missingCount, formatStringList(missing))
	return missingCount, usedCount, nil
}

func extractSelectedSceneCharacters() {
	path := cleanScenePath(sceneCardPath)
	if parsedScene == nil || parsedScenePath != path {
		parseSelectedScene()
		if parsedScene == nil {
			return
		}
	}
	output := cleanScenePath(sceneOutputDir)
	if output == "" {
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		output = filepath.Join(filepath.Dir(path), base+"_characters")
		sceneOutputDir = output
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		sceneStatus = fmt.Sprintf("创建输出目录失败：%v", err)
		g.Update()
		return
	}

	indexes := sceneExtractionIndexes(parsedScene, sceneUniqueOnly)
	modeName := "全部实例"
	if sceneUniqueOnly {
		modeName = "按角色去重"
	}

	var log strings.Builder
	written := 0
	for order, sourceIndex := range indexes {
		data, err := parsedScene.ExtractCharaCard(sourceIndex)
		if err != nil {
			fmt.Fprintf(&log, "[%d] 提取失败：%v\n", sourceIndex+1, err)
			continue
		}
		name := sceneCharacterName(parsedScene.CharaCards[sourceIndex])
		filename := fmt.Sprintf("%02d_%s.png", order+1, safeSceneFilePart(name))
		outputPath := nextAvailableScenePath(filepath.Join(output, filename))
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Fprintf(&log, "[%d] 写入失败：%v\n", sourceIndex+1, err)
			continue
		}
		written++
		fmt.Fprintf(&log, "[%d] %s\n", sourceIndex+1, outputPath)
	}
	sceneExportLog = log.String()
	sceneStatus = fmt.Sprintf("拆分完成（%s）：成功写入 %d/%d 张角色卡", modeName, written, len(indexes))
	g.Update()
}

func sceneExtractionIndexes(scene *KK.KKSceneCard, unique bool) []int {
	if scene == nil {
		return nil
	}
	if unique {
		return scene.UniqueCharaCardIndexes()
	}
	indexes := make([]int, len(scene.CharaCards))
	for index := range indexes {
		indexes[index] = index
	}
	return indexes
}

func buildSceneTexts(scene *KK.KKSceneCard) {
	sceneZipmods := dedupeSorted(append(
		collectScenePluginZipmods(scene.SceneExtendedList),
		sceneItemZipmodDependencies(scene)...,
	))
	allZipmods := dedupeSorted(append(scene.GetZipmodsDependencies(), sceneItemZipmodDependencies(scene)...))
	uniqueIndexes := scene.UniqueCharaCardIndexes()

	var summary strings.Builder
	fmt.Fprintf(&summary, "文件：%s\n", scene.SourceFileName)
	fmt.Fprintf(&summary, "场景版本：%s\n", scene.Version)
	fmt.Fprintf(&summary, "ExtendedSave 版本：%d\n", scene.ExtendedVersion)
	fmt.Fprintf(&summary, "场景插件数据：%d 项\n", len(scene.SceneExtendedList))
	fmt.Fprintf(&summary, "场景直接 Zipmod：%d 项\n", len(sceneZipmods))
	fmt.Fprintf(&summary, "场景与角色合计 Zipmod：%d 项\n", len(allZipmods))
	fmt.Fprintf(&summary, "角色实例：%d 个\n", len(scene.CharaCards))
	fmt.Fprintf(&summary, "不同角色：%d 个\n", len(uniqueIndexes))
	fmt.Fprintf(&summary, "\n全部 Zipmod 依赖：\n%s", formatStringList(allZipmods))
	sceneSummaryText = summary.String()

	pluginNames := make([]string, 0, len(scene.SceneExtendedList))
	for name := range scene.SceneExtendedList {
		pluginNames = append(pluginNames, name)
	}
	sort.Strings(pluginNames)
	var plugins strings.Builder
	fmt.Fprintf(&plugins, "场景插件数据（%d 项）：\n", len(pluginNames))
	for _, name := range pluginNames {
		plugin := scene.SceneExtendedList[name]
		fmt.Fprintf(&plugins, "\n- %s (v%d)\n", name, plugin.Version)
		zipmods := resolveInfoGUIDs(plugin.RequiredZipmodGUIDs)
		if len(zipmods) > 0 {
			fmt.Fprintf(&plugins, "  Zipmod: %s\n", strings.Join(zipmods, ", "))
		}
		if len(plugin.RequiredPluginGUIDs) > 0 {
			fmt.Fprintf(&plugins, "  插件依赖: %s\n", strings.Join(dedupeSorted(plugin.RequiredPluginGUIDs), ", "))
		}
	}
	fmt.Fprintf(&plugins, "\n场景直接 Zipmod（%d 项）：\n%s", len(sceneZipmods), formatStringList(sceneZipmods))
	scenePluginsText = plugins.String()

	uniqueSet := make(map[int]struct{}, len(uniqueIndexes))
	for _, index := range uniqueIndexes {
		uniqueSet[index] = struct{}{}
	}
	var characters strings.Builder
	fmt.Fprintf(&characters, "角色实例 %d 个，不同角色 %d 个：\n", len(scene.CharaCards), len(uniqueIndexes))
	for index, chara := range scene.CharaCards {
		kind := "重复角色实例"
		if _, ok := uniqueSet[index]; ok {
			kind = "不同角色（保留项）"
		}
		deps := dedupeSorted(chara.GetZipmodsDependencies())
		fmt.Fprintf(&characters, "\n[%d] %s [%s]\n", index+1, sceneCharacterName(chara), kind)
		fmt.Fprintf(&characters, "    插件数据：%d 项；Zipmod：%d 项\n", len(chara.ExtendedList), len(deps))
		for _, guid := range deps {
			fmt.Fprintf(&characters, "    - %s\n", guid)
		}
	}
	sceneCharactersText = characters.String()
}

func collectScenePluginZipmods(plugins map[string]*Base.PluginDataEx) []string {
	var values []string
	for _, plugin := range plugins {
		for _, info := range plugin.RequiredZipmodGUIDs {
			if info.GUID != "" {
				values = append(values, info.GUID)
			}
		}
	}
	return dedupeSorted(values)
}

func resolveInfoGUIDs(infos []Base.ResolveInfo) []string {
	values := make([]string, 0, len(infos))
	for _, info := range infos {
		if info.GUID != "" {
			values = append(values, info.GUID)
		}
	}
	return dedupeSorted(values)
}

func dedupeSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func formatStringList(values []string) string {
	if len(values) == 0 {
		return "（无）"
	}
	return "- " + strings.Join(values, "\n- ")
}

func sceneCharacterName(card *KK.KKCharaCard) string {
	if card == nil || card.CharInfo == nil {
		return "character"
	}
	name := strings.TrimSpace(card.CharInfo.Lastname + " " + card.CharInfo.Firstname)
	if name == "" {
		name = strings.TrimSpace(card.CharInfo.Nickname)
	}
	if name == "" {
		return "character"
	}
	return name
}

func safeSceneFilePart(value string) string {
	replacer := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "_", "/", "_",
		"\\", "_", "|", "_", "?", "_", "*", "_", " ", "_",
	)
	value = strings.Trim(replacer.Replace(value), " ._")
	if value == "" {
		return "character"
	}
	return value
}

func nextAvailableScenePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for suffix := 2; ; suffix++ {
		candidate := fmt.Sprintf("%s_%d%s", base, suffix, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func cleanScenePath(path string) string {
	return strings.Trim(strings.TrimSpace(path), "\"")
}

func resetSceneResult() {
	parsedScene = nil
	parsedScenePath = ""
	sceneSummaryText = ""
	scenePluginsText = ""
	sceneCharactersText = ""
	sceneMissingModsText = ""
	sceneUsedModsText = ""
	sceneStatus = ""
}
