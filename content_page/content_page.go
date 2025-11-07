package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"os"
	"path/filepath"

	g "github.com/AllenDang/giu"
)

// 软件所在目录
var appDir string = "."

func init() {
	appDir = filepath.Dir(os.Args[0])
}

// 显示初始配置窗口
func InitialConfigWindow() g.Layout {
	// 如果尚未完成初始配置，则显示欢迎页面
	if !config.Instance.InitConfig {
		return Welcome()
	}
	if !config.Instance.HideModeInfoBuild {
		return ModInfoBuild()
	}
	// 否则显示主内容页面
	return MainContent()
}

// 内容页面显示初始配置完成后的主界面
func MainContent() g.Layout {
	return g.Layout{
		g.Child().Border(true).Layout(
			g.TabBar().TabItems(
				CardModVerification(),
				ModeMatchTab(),
				ModeCardUseTab(),
				g.TabItem(" Mod信息构建 ").Layout(
					ModInfoBuild(),
				),
				g.TabItem(" 设置 ").Layout(
					g.Label("软件安装目录："),
					g.InputText(&appDir).Flags(g.InputTextFlagsReadOnly),
					g.Button("打开软件所在目录").OnClick(func() {
						util.OpenURL(".")
					}),
					g.Separator(),
					g.Row(
						g.Label("软件配置文件: "),
						func() g.Widget {
							cfgPath := config.GetConfigPath()
							return g.InputText(&cfgPath).Flags(g.InputTextFlagsReadOnly)
						}(),
					),
					g.Row(
						g.Button("重置").OnClick(func() {
							config.ResetToDefault()
						}),
						g.Button("打开配置文件").OnClick(func() {
							util.OpenFile(config.GetConfigPath())
						}),
					),
					g.Separator(),
					g.Label("检测目录下的所有卡片或者是某一卡片的Mod是否缺失"),
					util.FilePicker("游戏目录", &config.Instance.GamePath).ForDirectories(),
					g.Separator(),
				),
			),
		),
	}

}
