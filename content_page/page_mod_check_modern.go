package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
)

var resultViewMissing = true

// CardModVerificationPage 是卡片检测功能的现代主页面。
func CardModVerificationPage() g.Layout {
	syncCardSelection()

	return g.Layout{
		g.Custom(func() {
			availableWidth, _ := g.GetAvailableRegion()
			gap := float32(14)
			leftWidth := (availableWidth - gap) * 0.58
			rightWidth := availableWidth - gap - leftWidth
			g.Row(cardImportPanel(leftWidth), cardPreviewPanel(rightWidth)).Build()
		}),
		g.Dummy(0, 4),
		cardResultsPanel(),
	}
}

func syncCardSelection() {
	if lastCardPath == "" {
		lastCardPath = config.Instance.CardPath
		return
	}
	if config.Instance.CardPath == lastCardPath {
		return
	}

	missingMods = nil
	missingModsText = ""
	usedMods = nil
	usedModsText = ""
	tips = ""
	checked = false
	lastCardPath = config.Instance.CardPath
}

func cardImportPanel(width float32) g.Widget {
	return fixedPanelStyle("card-import-panel", width, 326,
		sectionLabel("选择卡片"),
		insetStyle("card-drop-zone", -1, 184,
			g.Align(g.AlignCenter).To(textLabel("PNG", 24, uiMuted)),
			g.Align(g.AlignCenter).To(sectionLabel("拖入角色卡或文件夹")),
			g.Align(g.AlignCenter).To(mutedLabel("支持 PNG 角色卡，也可批量检测目录")),
			g.Align(g.AlignCenter).To(g.Row(
				primaryButton("选择文件", 130, 40, chooseCardFile),
				secondaryButton("选择文件夹", 130, 40, chooseCardDirectory),
			)),
		),
		g.Label("当前路径"),
		g.InputText(&config.Instance.CardPath).Size(-1).Flags(g.InputTextFlagsReadOnly),
	)
}

func chooseCardFile() {
	selected, err := util.BrowseFile("选择角色卡", config.Instance.CardPath, "*.png")
	if err == nil && selected != "" {
		config.Instance.CardPath = selected
		g.Update()
	}
}

func chooseCardDirectory() {
	selected, err := util.BrowseDirectory("选择卡片目录", config.Instance.CardPath)
	if err == nil && selected != "" {
		config.Instance.CardPath = selected
		g.Update()
	}
}

func cardPreviewPanel(width float32) g.Widget {
	path := strings.Trim(strings.TrimSpace(config.Instance.CardPath), "\"")
	widgets := []g.Widget{sectionLabel("卡片预览")}

	if isPreviewableCard(path) {
		widgets = append(widgets,
			g.Row(
				g.ImageWithFile(path).Size(130, 184),
				g.Layout{
					g.Dummy(0, 18),
					textLabel(filepath.Base(path), 18, uiText),
					mutedLabel("PNG 角色卡"),
					statusLabel("已识别", uiSuccess),
				},
			),
		)
	} else if path != "" {
		widgets = append(widgets,
			g.Dummy(0, 32),
			g.Align(g.AlignCenter).To(textLabel("DIR", 28, uiMuted)),
			g.Align(g.AlignCenter).To(sectionLabel(filepath.Base(path))),
			g.Align(g.AlignCenter).To(mutedLabel("目录模式 · 将递归检测 PNG 卡片")),
			g.Dummy(0, 44),
		)
	} else {
		widgets = append(widgets,
			g.Dummy(0, 42),
			g.Align(g.AlignCenter).To(textLabel("[  ]", 28, uiBorder)),
			g.Align(g.AlignCenter).To(mutedLabel("选择卡片后显示预览")),
			g.Dummy(0, 58),
		)
	}

	widgets = append(widgets,
		g.Style().SetDisabled(path == "").To(primaryButton("开始检测", -1, 48, SingleCardCheck)),
	)
	return fixedPanelStyle("card-preview-panel", width, 326, widgets...)
}

func isPreviewableCard(path string) bool {
	if !strings.EqualFold(filepath.Ext(path), ".png") {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func cardResultsPanel() g.Widget {
	statusText := "等待检测"
	statusColor := uiMuted
	if checked {
		statusText = "检测完成"
		statusColor = uiSuccess
	} else if tips != "" {
		statusText = "提示"
		statusColor = uiWarning
	}

	return panelStyle("card-results-panel", -1, -1,
		g.Row(sectionLabel("检测结果"), statusLabel(statusText, statusColor), mutedLabel(tips)),
		resultStats(),
		resultToolbar(),
		resultTable(),
	)
}

func resultStats() g.Widget {
	return g.Custom(func() {
		availableWidth, _ := g.GetAvailableRegion()
		gap := float32(12)
		cardWidth := (availableWidth - gap*2) / 3
		g.Row(
			statCard("解析卡片", parsedCardCount(), uiPrimary, cardWidth),
			statCard("缺失 MOD", len(missingMods), uiDanger, cardWidth),
			statCard("使用 MOD", len(usedMods), uiSuccess, cardWidth),
		).Build()
	})
}

func statCard(label string, value int, accent color.Color, width float32) g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, uiPanelAlt).
		SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 16, 12).
		SetStyleFloat(g.StyleVarChildRounding, 6).
		SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().Size(width, 82).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			mutedLabel(label),
			textLabel(fmt.Sprintf("%d", value), 26, accent),
		))
}

func parsedCardCount() int {
	if !checked {
		return 0
	}
	var count int
	if _, err := fmt.Sscanf(tips, "检索完成，共解析 %d 张卡片", &count); err == nil {
		return count
	}
	if strings.EqualFold(filepath.Ext(config.Instance.CardPath), ".png") {
		return 1
	}
	return 0
}

func resultToolbar() g.Widget {
	return g.Row(
		resultTabButton(fmt.Sprintf("缺失的 MOD  %d", len(missingMods)), true),
		resultTabButton(fmt.Sprintf("使用的 MOD  %d", len(usedMods)), false),
		g.Spacing(),
		secondaryButton("复制列表", 110, 36, copyCurrentResult),
		secondaryButton("清除结果", 110, 36, clearCardResults),
	)
}

func resultTabButton(label string, missing bool) g.Widget {
	selected := resultViewMissing == missing
	buttonColor := uiInput
	textColor := uiMuted
	if selected {
		buttonColor = uiPrimaryLo
		textColor = uiPrimary
	}
	return g.Style().
		SetColor(g.StyleColorText, textColor).
		SetColor(g.StyleColorButton, buttonColor).
		SetColor(g.StyleColorButtonHovered, uiPanelAlt).
		SetColor(g.StyleColorButtonActive, uiPrimaryLo).
		To(g.Button(label).Size(160, 38).OnClick(func() {
			resultViewMissing = missing
			g.Update()
		}))
}

func resultTable() g.Widget {
	entries := usedMods
	if resultViewMissing {
		entries = missingMods
	}
	rows := make([]*g.TableRowWidget, 0, len(entries)+1)
	if len(entries) == 0 {
		emptyText := "尚未检测"
		if checked {
			if resultViewMissing {
				emptyText = "没有缺失的 MOD"
			} else {
				emptyText = "未发现 MOD 依赖"
			}
		}
		rows = append(rows, g.TableRow(mutedLabel(emptyText), g.Label("—"), g.Label("—")).MinHeight(42))
	} else {
		for _, entry := range entries {
			value := entry
			stateText := "已使用"
			stateColor := uiSuccess
			if resultViewMissing {
				stateText = "缺失"
				stateColor = uiDanger
			}
			rows = append(rows, g.TableRow(
				g.Label(value),
				textLabel(stateText, 15, stateColor),
				g.Style().SetColor(g.StyleColorText, uiPrimary).To(
					g.Button("复制##"+value).OnClick(func() { _ = util.CopyToClipboard(value) }),
				),
			).MinHeight(42))
		}
	}

	return g.Table().
		Columns(
			g.TableColumn("MOD GUID").Flags(g.TableColumnFlagsWidthStretch).InnerWidthOrWeight(0.68),
			g.TableColumn("状态").Flags(g.TableColumnFlagsWidthStretch).InnerWidthOrWeight(0.16),
			g.TableColumn("操作").Flags(g.TableColumnFlagsWidthStretch).InnerWidthOrWeight(0.16),
		).
		Rows(rows...).Size(-1, -1).
		Flags(g.TableFlagsRowBg | g.TableFlagsBordersInnerH | g.TableFlagsScrollY | g.TableFlagsSizingStretchProp)
}

func copyCurrentResult() {
	text := usedModsText
	if resultViewMissing {
		text = missingModsText
	}
	if text == "" {
		return
	}
	if err := util.CopyToClipboard(text); err != nil {
		tips = "复制到剪贴板失败"
	} else {
		tips = "已复制到剪贴板"
	}
	g.Update()
}

func clearCardResults() {
	missingMods = nil
	missingModsText = ""
	usedMods = nil
	usedModsText = ""
	tips = ""
	checked = false
	g.Update()
}
