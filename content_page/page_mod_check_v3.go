package content_page

import (
	"KKCardModCheck/config"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
)

func CardModVerificationPageV3() g.Layout {
	syncCardSelection()
	return g.Layout{
		g.Custom(func() {
			w, _ := g.GetAvailableRegion()
			gap := float32(16)
			left := (w - gap) * 0.52
			g.Row(cardImportPanelV3(left), cardPreviewPanelV3(w-gap-left)).Build()
		}),
		g.Dummy(0, 4),
		cardResultsPanelV3(),
	}
}

func cardImportPanelV3(width float32) g.Widget {
	return fixedPanelStyle("card-import-panel-v3", width, 396,
		boldTextV3("选择卡片", 20, uiText),
		dropZoneV3("drop-zone-v3", 246,
			g.Align(g.AlignCenter).To(iconV3("drop-image-v3", "image", 66, uiMuted)),
			g.Align(g.AlignCenter).To(boldTextV3("拖入角色卡或文件夹", 22, uiText)),
			g.Align(g.AlignCenter).To(textLabel("支持 PNG 角色卡，也可批量检测目录", 17, uiMuted)),
			g.Align(g.AlignCenter).To(g.Row(
				primaryButton("选择文件", 160, 46, chooseCardFile),
				secondaryButton("选择文件夹", 160, 46, chooseCardDirectory),
			)),
		),
		textLabel("当前路径", 16, uiText),
		g.InputText(&config.Instance.CardPath).Size(-1).Flags(g.InputTextFlagsReadOnly),
	)
}

func cardPreviewPanelV3(width float32) g.Widget {
	path := strings.Trim(strings.TrimSpace(config.Instance.CardPath), "\"")
	var content g.Widget
	if isPreviewableCard(path) {
		content = g.Row(
			g.ImageWithFile(path).Size(178, 252),
			g.Layout{
				g.Dummy(0, 30),
				boldTextV3(filepath.Base(path), 18, uiText),
				textLabel("PNG 角色卡", 17, uiMuted),
				statusLabel("已识别", uiSuccess),
			},
		)
	} else {
		content = g.Row(
			placeholderCardV3(160, 238),
			g.Layout{
				g.Dummy(0, 42),
				boldTextV3("等待选择卡片", 18, uiText),
				textLabel("PNG · 角色卡预览", 16, uiMuted),
				statusLabel("未载入", uiMuted),
			},
		)
	}
	return fixedPanelStyle("card-preview-panel-v3", width, 396,
		boldTextV3("卡片预览", 20, uiText),
		content,
		g.Style().SetDisabled(path == "").To(primaryButton("开始检测", -1, 54, SingleCardCheck)),
	)
}

func cardResultsPanelV3() g.Widget {
	status := "等待检测"
	col := uiMuted
	if checked {
		status = "检测完成"
		col = uiSuccess
	} else if tips != "" {
		status = "提示"
		col = uiWarning
	}
	return panelStyle("card-results-panel-v3", -1, -1,
		g.Row(boldTextV3("检测结果", 22, uiText), statusLabel(status, col), textLabel(tips, 15, uiMuted)),
		resultStatsV3(),
		resultToolbar(),
		resultTable(),
		g.Align(g.AlignRight).To(textLabel(resultFooterV3(), 16, uiMuted)),
	)
}

func resultStatsV3() g.Widget {
	return g.Custom(func() {
		w, _ := g.GetAvailableRegion()
		gap := float32(14)
		cw := (w - gap*2) / 3
		g.Row(
			statCardV3("document", "解析卡片", parsedCardCount(), uiPrimary, cw),
			statCardV3("warning", "缺失 MOD", len(missingMods), uiDanger, cw),
			statCardV3("cube", "使用 MOD", len(usedMods), uiSuccess, cw),
		).Build()
	})
}

func statCardV3(kind, label string, value int, accent color.Color, width float32) g.Widget {
	return g.Style().SetColor(g.StyleColorChildBg, uiPanelAlt).SetColor(g.StyleColorBorder, uiBorder).
		SetStyle(g.StyleVarWindowPadding, 18, 14).SetStyleFloat(g.StyleVarChildRounding, 7).SetStyleFloat(g.StyleVarChildBorderSize, 1).
		To(g.Child().ID(g.ID("stat-v3-"+kind)).Size(width, 92).Border(true).Flags(g.WindowFlagsNoScrollbar).Layout(
			g.Row(
				iconV3("stat-icon-v3-"+kind, kind, 48, accent),
				g.Layout{textLabel(label, 17, uiMuted), boldTextV3(fmt.Sprintf("%d", value), 30, uiText)},
			),
		))
}

func resultFooterV3() string {
	if !checked {
		return "等待选择卡片并开始检测"
	}
	return "检测完成"
}
