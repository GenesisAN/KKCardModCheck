package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/giu/custom"
	"KKCardModCheck/util"
	"os"
	"path/filepath"

	g "github.com/AllenDang/giu"
)

// Welcome shows the initial welcome/configuration UI
func Welcome() g.Layout {
	// 先构建 child 内部的 widget 列表，以便在需要时插入提示信息
	childWidgets := g.Layout{
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label("★"),
				g.Label("欢迎使用 KKCardModCheck！"),
				g.Label("★"),
			),
		),

		g.Align(g.AlignCenter).To(
			g.Label("遇到 BUG 可通过顶部菜单 \"关于 > BUG反馈\" 反馈"),
		),
		g.Align(g.AlignCenter).To(
			g.Label("功能建议/提问请到 \"关于 > 讨论区\" 发帖"),
		),
		g.Align(g.AlignCenter).To(
			g.Row(
				g.Label("也可以来我们的"),
				custom.ButtonStyleBlue.To(
					g.Button("QQ群").OnClick(func() {
						util.OpenURL("https://qm.qq.com/q/rza05u1550")
					}),
				),
				g.Label("交流"),
			),
		),
		g.Separator(),

		g.Label("请先配置游戏目录以继续使用"),

		g.Label("游戏目录是指 Koikatu_Data 文件夹，选中后点击保存并继续"),

		g.Separator(),
		util.FilePicker("游戏目录:", &config.Instance.GamePath).ForDirectories(),
		g.Align(g.AlignCenter).To(
			g.Button("保存并继续").Size(200, 80).OnClick(func() {
				// 保存配置，若所选目录不是 Koikatu_Data 则提示但仍保存为 Mod 存放目录
				config.Instance.InitConfig = true
				config.Save()
				// 检查最后一级目录名是否为 Koikatu_Data
			}),
		),
	}
	if config.Instance.GamePath != "" {
		// 检查所选目录是否为 Koikatu_Data 本身，或者该目录下是否包含名为 Koikatu_Data 的子目录（仅查一层，不递归）
		valid := false
		if filepath.Base(config.Instance.GamePath) == "Koikatu_Data" {
			valid = true
		} else {
			// 尝试读取所选目录的直接子项，查找名为 Koikatu_Data 的目录
			if entries, err := os.ReadDir(config.Instance.GamePath); err == nil {
				for _, e := range entries {
					if e.IsDir() && e.Name() == "Koikatu_Data" {
						valid = true
						break
					}
				}
			}
		}
		if !valid {
			// 标记在界面显示提示（仍然保存路径并继续）
			childWidgets = append(childWidgets,
				g.Separator(),
				g.Align(g.AlignCenter).To(
					g.Label("警告: 所选文件夹不是 Koikatu游戏目录，因为其下不包含 Koikatu_Data 子文件夹。但程序可以把当前路径作为 Mod 存放目录来进行使用。"),
				),
			)
		}
	}

	return g.Layout{
		g.Child().Border(true).Layout(childWidgets...),
	}
}
