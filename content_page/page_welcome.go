package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"os"
	"path/filepath"

	g "github.com/AllenDang/giu"
)

// Welcome 显示统一视觉风格的首次使用配置页。
func Welcome() g.Layout {
	return g.Layout{
		g.Custom(func() {
			availableWidth, availableHeight := g.GetAvailableRegion()
			panelWidth := float32(760)
			panelHeight := float32(570)
			if availableWidth < panelWidth+48 {
				panelWidth = availableWidth - 48
			}
			if availableHeight > panelHeight {
				g.Dummy(0, (availableHeight-panelHeight)/2).Build()
			}
			g.Align(g.AlignCenter).To(
				panelStyle("welcome-panel", panelWidth, panelHeight,
					g.Align(g.AlignCenter).To(textLabel("◆", 36, uiPrimary)),
					g.Align(g.AlignCenter).To(titleLabel("欢迎使用 KKCardModCheck")),
					g.Align(g.AlignCenter).To(mutedLabel("角色卡 MOD 依赖检查与管理工具")),
					g.Dummy(0, 14),
					insetStyle("welcome-steps", -1, 88,
						g.Row(
							textLabel("1", 24, uiPrimary),
							g.Layout{g.Label("选择游戏目录"), mutedLabel("定位 Koikatu_Data")},
							g.Spacing(),
							textLabel("2", 24, uiPrimary),
							g.Layout{g.Label("构建 MOD 信息"), mutedLabel("扫描本地 Zipmod")},
							g.Spacing(),
							textLabel("3", 24, uiPrimary),
							g.Layout{g.Label("开始检测"), mutedLabel("核对卡片依赖")},
						),
					),
					g.Dummy(0, 10),
					sectionLabel("设置游戏目录"),
					mutedLabel("请选择 Koikatsu / Koikatsu Sunshine 安装目录，或直接选择 Koikatu_Data"),
					g.Dummy(0, 4),
					util.FilePicker("游戏目录", &config.Instance.GamePath).ForDirectories().LabelWidth(110).ButtonWidth(140),
					welcomeWarning(),
					g.Dummy(0, 12),
					g.Align(g.AlignCenter).To(
						primaryButton("保存并继续", 220, 52, func() {
							config.Instance.InitConfig = true
							config.Save()
							g.Update()
						}),
					),
					g.Dummy(0, 10),
					g.Align(g.AlignCenter).To(g.Row(
						secondaryButton("BUG 反馈", 120, 34, func() {
							util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/issues")
						}),
						secondaryButton("讨论区", 100, 34, func() {
							util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/discussions")
						}),
						secondaryButton("QQ群", 90, 34, func() {
							util.OpenURL("https://qm.qq.com/q/rza05u1550")
						}),
					)),
				),
			).Build()
		}),
	}
}

func welcomeWarning() g.Widget {
	if config.Instance.GamePath == "" {
		return mutedLabel("选择目录后将进入 MOD 信息构建步骤")
	}
	if isKoikatsuGameDirectory(config.Instance.GamePath) {
		return textLabel("目录识别成功", 15, uiSuccess)
	}
	return textLabel("未发现 Koikatu_Data；仍可把当前路径作为 MOD 存放目录使用", 15, uiWarning)
}

func isKoikatsuGameDirectory(path string) bool {
	if filepath.Base(path) == "Koikatu_Data" {
		return true
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == "Koikatu_Data" {
			return true
		}
	}
	return false
}
