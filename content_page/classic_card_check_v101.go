package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	g "github.com/AllenDang/giu"
)

type classicCheckSnapshotV101 struct {
	tips           string
	missing        []string
	used           []string
	missingText    string
	usedText       string
	checked        bool
	checking       bool
	progressParsed int
	progressScenes int
	progressErrors int
}

func takeClassicCheckSnapshotV101() classicCheckSnapshotV101 {
	checkMu.Lock()
	defer checkMu.Unlock()
	return classicCheckSnapshotV101{
		tips:           tips,
		missing:        append([]string(nil), missingMods...),
		used:           append([]string(nil), usedMods...),
		missingText:    missingModsText,
		usedText:       usedModsText,
		checked:        checked,
		checking:       isChecking,
		progressParsed: progressParsed,
		progressScenes: progressScenes,
		progressErrors: progressErrors,
	}
}

func syncClassicCardSelectionV101() {
	currentPath := config.Instance.CardPath
	checkMu.Lock()
	defer checkMu.Unlock()
	if lastCardPath == "" {
		lastCardPath = currentPath
		return
	}
	if currentPath == lastCardPath {
		return
	}
	if isChecking {
		cancelRequested.Store(true)
		tips = "路径变化，正在取消当前检查..."
	} else {
		resetCheckStateLocked()
	}
	lastCardPath = currentPath
}

func ClassicCardCheckPageV101() g.Layout {
	syncClassicCardSelectionV101()
	snapshot := takeClassicCheckSnapshotV101()
	return g.Layout{
		g.Custom(func() {
			width, _ := g.GetAvailableRegion()
			gap := float32(8)
			left := (width - gap) * 0.58
			g.Row(classicInputPanel(left), classicPreviewPanelV101(width-gap-left, snapshot)).Build()
		}),
		classicResultsPanelV101(snapshot),
	}
}

func classicPreviewPanelV101(width float32, snapshot classicCheckSnapshotV101) g.Widget {
	path := strings.Trim(strings.TrimSpace(config.Instance.CardPath), "\"")
	layout := g.Layout{g.Label("预览 / 操作"), g.Dummy(0, 2)}
	if isPreviewableCard(path) {
		layout = append(layout,
			g.Row(
				g.ImageWithFile(path).Size(100, 142),
				g.Layout{g.Label(filepath.Base(path)), g.Label("PNG 角色卡 / KK 场景卡"), g.Label("状态：已载入")},
			),
		)
	} else if path != "" {
		layout = append(layout, g.Label("目录模式"), g.Label(filepath.Base(path)), g.Label("状态：已载入"))
	} else {
		layout = append(layout, g.Label("尚未选择卡片或目录"))
	}
	layout = append(layout, g.Custom(func() {
		_, remaining := g.GetAvailableRegion()
		// Reserve the action height and the gap preceding it.
		if remaining > classicButtonHeight+8 {
			g.Dummy(0, remaining-classicButtonHeight-8).Build()
		}
	}))
	if snapshot.checking {
		layout = append(layout,
			classicButton("取消检测").Size(-1, classicButtonHeight).OnClick(cancelClassicCheckV101),
		)
	} else {
		layout = append(layout,
			g.Style().SetDisabled(path == "").To(classicButton("开始检测").Size(-1, classicButtonHeight).OnClick(SingleCardCheck)),
		)
	}
	return g.Child().ID(g.ID("classic-preview-v101")).Size(width, 280).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(layout...)
}

func classicResultsPanelV101(snapshot classicCheckSnapshotV101) g.Widget {
	status := "状态：等待检测"
	if snapshot.checking {
		status = fmt.Sprintf("状态：正在检测，已解析 %d（场景卡 %d），错误 %d", snapshot.progressParsed, snapshot.progressScenes, snapshot.progressErrors)
	} else if snapshot.checked {
		status = "状态：检测完成"
	} else if snapshot.tips != "" {
		status = "状态：" + snapshot.tips
	}

	return g.Child().ID(g.ID("classic-results-v101")).Size(-1, -1).Border(true).Layout(
		g.Row(
			g.Label("检测结果"),
			g.Dummy(12, 0),
			g.Label(status),
			g.Spacing(),
			g.Style().SetDisabled(snapshot.checking).To(classicButton("清除").OnClick(clearClassicCheckV101)),
		),
		g.Row(
			g.Label(fmt.Sprintf("解析卡片: %d", snapshot.progressParsed)),
			g.Dummy(12, 0),
			g.Style().SetColor(g.StyleColorText, uiDanger).To(g.Label(fmt.Sprintf("综合缺失 MOD: %d", len(snapshot.missing)))),
			g.Dummy(12, 0),
			g.Style().SetColor(g.StyleColorText, uiSuccess).To(g.Label(fmt.Sprintf("综合使用 MOD: %d", len(snapshot.used)))),
		),
		g.TabBar().TabItems(
			g.TabItem(fmt.Sprintf("综合缺失的 MOD (%d)", len(snapshot.missing))).Layout(
				g.Row(g.Style().SetDisabled(snapshot.missingText == "").To(classicButton("复制列表").OnClick(func() {
					copyClassicCheckTextV101(snapshot.missingText)
				}))),
				classicResultTableV101(snapshot.missing, "缺失", snapshot.checked),
			),
			g.TabItem(fmt.Sprintf("综合使用的 MOD (%d)", len(snapshot.used))).Layout(
				g.Row(g.Style().SetDisabled(snapshot.usedText == "").To(classicButton("复制列表").OnClick(func() {
					copyClassicCheckTextV101(snapshot.usedText)
				}))),
				classicResultTableV101(snapshot.used, "使用", snapshot.checked),
			),
			unusedModsTab(),
		),
	)
}

func classicResultTableV101(entries []string, status string, checked bool) g.Widget {
	rows := make([]*g.TableRowWidget, 0, len(entries)+1)
	if len(entries) == 0 {
		message := "尚无结果"
		if checked {
			message = "没有条目"
		}
		rows = append(rows, g.TableRow(g.Label(message), g.Label("-"), g.Label("-")).MinHeight(28))
	} else {
		for index, entry := range entries {
			value := entry
			rows = append(rows, g.TableRow(
				g.Label(value),
				g.Label(status),
				classicButton("复制##classic-v101-"+status+strconv.Itoa(index)).Size(classicInlineButtonWidth, classicInlineButtonHeight).OnClick(func() { _ = util.CopyToClipboard(value) }),
			).MinHeight(30))
		}
	}
	return g.Table().Columns(
		g.TableColumn("MOD GUID").Flags(g.TableColumnFlagsWidthStretch).InnerWidthOrWeight(0.72),
		g.TableColumn("状态").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(90),
		g.TableColumn("操作").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(80),
	).Rows(rows...).Size(-1, -1).Flags(g.TableFlagsRowBg | g.TableFlagsBordersInnerH | g.TableFlagsResizable | g.TableFlagsScrollY | g.TableFlagsSizingStretchProp)
}

func cancelClassicCheckV101() {
	cancelRequested.Store(true)
	checkMu.Lock()
	tips = "正在取消，请稍候..."
	checkMu.Unlock()
	g.Update()
}

func clearClassicCheckV101() {
	checkMu.Lock()
	if isChecking {
		cancelRequested.Store(true)
		tips = "正在取消，请稍候..."
	} else {
		resetCheckStateLocked()
	}
	checkMu.Unlock()
	g.Update()
}

func copyClassicCheckTextV101(text string) {
	if text == "" {
		return
	}
	if err := util.CopyToClipboard(text); err != nil {
		checkMu.Lock()
		tips = "复制到剪贴板失败"
		checkMu.Unlock()
	} else {
		checkMu.Lock()
		tips = "已复制到剪贴板"
		checkMu.Unlock()
	}
	g.Update()
}
