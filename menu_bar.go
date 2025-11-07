package main

import (
	"KKCardModCheck/util"
	"fmt"

	g "github.com/AllenDang/giu"
)

// App的菜单栏
func AppMenuBar() g.Widget {
	return g.MenuBar().Layout(
		g.Menu("文件").Layout(
			// g.MenuItem("更新Mod信息").OnClick(func() {
			// 	util.OpenFile(config.GetConfigPath())
			// }),
			g.MenuItem("打开软件所在位置").OnClick(func() { util.OpenURL(".") }),
			g.MenuItem("打开游戏目录").OnClick(func() { util.OpenURL("../") }),
			g.Separator(),
			g.MenuItem("退出").OnClick(func() { wnd.Close() }),
		),
		g.Menu("关于").Layout(
			g.MenuItem("BUG反馈").OnClick(func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/issues")
			}),
			g.MenuItem("讨论区").OnClick(func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/discussions")
			}),
			g.MenuItem("GitHub仓库").OnClick(func() {
				util.OpenURL("https://github.com/GenesisAN/KKCardModCheck")
			}),
			g.MenuItem("QQ群").OnClick(func() {
				util.OpenURL("https://qm.qq.com/q/rza05u1550")
			}),
			g.Separator(),
			g.Menu("版本信息").Layout(
				g.Label(fmt.Sprintf("当前:%s", version)),
				g.Label(fmt.Sprintf("构建日期:%s", buildDate)),
				g.Label(fmt.Sprintf("Git提交:%s", gitHash)),
				g.Separator(),
				func() g.Widget {
					if version_res != nil {
						return g.Label(fmt.Sprintf("最新:%s", version_res.Current))
					} else if version_res == nil && version_err != nil {
						return g.Label("最新:检查失败")
					}
					return g.Label("最新:检查中...")
				}(),
			),
		),
		func() g.Widget {
			if version_res != nil && version_res.Outdated {
				return g.Menu("有新版本！").Layout(
					g.Label(fmt.Sprintf("%s → %s", version, version_res.Current)),
					g.MenuItem("打开发布页").OnClick(func() {
						util.OpenURL("https://github.com/GenesisAN/KKCardModCheck/releases")
					}),
				)
			}
			return g.Custom(func() {})
		}(),
	)
}
