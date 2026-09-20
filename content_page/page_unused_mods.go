package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/internal/modextract"
	"KKCardModCheck/util"
	"fmt"
	g "github.com/AllenDang/giu"
	"path/filepath"
	"strings"
)

var unusedSource, unusedDestination, unusedLog string
var unusedPlan *modextract.Plan
var unusedBusy bool
var unusedGeneration uint64
var unusedPreviewSource, unusedPreviewDestination string
var unusedResults = make(chan unusedResult, 1)

type unusedResult struct {
	plan  *modextract.Plan
	log   string
	moved int
}

func unusedScanSnapshot() (bool, uint64, []string) {
	checkMu.Lock()
	defer checkMu.Unlock()
	return checked && !isChecking && progressParsed > 0 && progressErrors == 0 && completedCardPath == config.Instance.CardPath,
		checkGeneration, append([]string(nil), usedMods...)
}

func unusedModsTab() *g.TabItemWidget {
	select {
	case result := <-unusedResults:
		unusedBusy = false
		unusedPlan, unusedLog = result.plan, result.log
		if result.moved > 0 {
			config.ModInfoChage = true
		}
	default:
	}
	if unusedSource == "" {
		unusedSource = config.Instance.ModPath
	}
	if unusedSource == "" && config.Instance.GamePath != "" {
		unusedSource = filepath.Join(config.Instance.GamePath, "mods")
	}
	if unusedDestination == "" {
		unusedDestination = config.Instance.UnusedModOutputPath
	}
	ready, generation, _ := unusedScanSnapshot()
	valid := ready && unusedPlan != nil && generation == unusedGeneration && unusedSource == unusedPreviewSource && unusedDestination == unusedPreviewDestination
	empty := unusedPlan == nil
	if unusedPlan != nil {
		empty = len(unusedPlan.Entries) == 0
	}
	return g.TabItem("未使用 MOD 批量提取").Layout(
		g.Label("仅以本次扫描的角色卡 / 场景卡为准，建议先扫描完整卡片目录。"),
		g.Label("保留使用中的 MOD 和无法识别的文件；移动时保留子目录，不覆盖同名文件。"),
		g.Style().SetDisabled(unusedBusy).To(
			util.FilePicker("MOD 源目录", &unusedSource).ForDirectories().LabelWidth(120).ButtonWidth(144),
			util.FilePicker("提取后存放目录", &unusedDestination).ForDirectories().LabelWidth(120).ButtonWidth(144),
		),
		g.Row(
			classicButton("预览未使用 MOD").Disabled(unusedBusy || !ready).OnClick(previewUnusedMods),
			classicButton("确认批量移动").Disabled(unusedBusy || !valid || empty).OnClick(moveUnusedMods),
			classicButton("打开存放目录").Disabled(unusedDestination == "").OnClick(func() { util.OpenFile(unusedDestination) }),
		),
		g.Label(func() string {
			if unusedBusy {
				return "正在处理，请稍候..."
			}
			if !ready {
				return "请先完成卡片检测；取消、读取错误或没有成功解析的卡片时不能提取。"
			}
			if unusedPlan != nil && !valid {
				return "扫描结果或目录已变化，请重新预览。"
			}
			return "先预览清单，确认后将文件从 MOD 源目录移动到指定目录。"
		}()),
		g.InputTextMultiline(&unusedLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly),
	)
}

func previewUnusedMods() {
	ready, generation, used := unusedScanSnapshot()
	if unusedBusy || !ready {
		return
	}
	unusedGeneration = generation
	unusedPreviewSource, unusedPreviewDestination = unusedSource, unusedDestination
	source, destination := unusedSource, unusedDestination
	unusedBusy, unusedPlan, unusedLog = true, nil, "正在扫描 MOD 目录..."
	go func() {
		result := unusedResult{}
		defer func() {
			if r := recover(); r != nil {
				result = unusedResult{log: fmt.Sprintf("预览失败: %v", r)}
			}
			unusedResults <- result
			g.Update()
		}()
		plan, err := modextract.Preview(source, destination, used, func(path string) (string, error) {
			mod, err := util.ReadZip("", path, 0)
			return mod.GUID, err
		})
		if err != nil {
			result.log = "预览失败: " + err.Error()
			return
		}
		var log strings.Builder
		fmt.Fprintf(&log, "待移动 %d 个 MOD；保留无法识别 / 链接文件 %d 个\n源目录：%s\n存放目录：%s\n\n", len(plan.Entries), len(plan.Skipped), plan.Source, plan.Destination)
		for _, entry := range plan.Entries {
			fmt.Fprintf(&log, "%s  [%s]\n", entry.Relative, entry.GUID)
		}
		for _, skip := range plan.Skipped {
			fmt.Fprintf(&log, "保留：%s\n", skip)
		}
		result.plan, result.log = plan, log.String()
	}()
}

func moveUnusedMods() {
	ready, generation, _ := unusedScanSnapshot()
	if unusedBusy || !ready || unusedPlan == nil || generation != unusedGeneration || unusedSource != unusedPreviewSource || unusedDestination != unusedPreviewDestination {
		return
	}
	plan := unusedPlan
	unusedBusy, unusedPlan = true, nil
	config.Instance.UnusedModOutputPath = unusedDestination
	config.Save()
	go func() {
		var log strings.Builder
		moved, failed := plan.Move(func(line string) { fmt.Fprintln(&log, line) })
		unusedResults <- unusedResult{moved: moved, log: fmt.Sprintf("完成：已移动 %d 个，失败 %d 个。\n存放目录：%s\n移动后请重新构建 MOD 信息数据库。\n\n%s", moved, failed, plan.Destination, log.String())}
		g.Update()
	}()
}
