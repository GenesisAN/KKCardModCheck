package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
	ic "github.com/GenesisAN/illusionsCard"
)

var (
	mmModPath    string
	mmCardPath   string
	mmStatus     string
	mmLog        string
	mmInProgress bool
)

func ModeMatchTab() *g.TabItemWidget {
	// keep fields in sync with global config defaults
	if mmCardPath == "" {
		mmCardPath = config.Instance.CardPath
	}
	layout := []g.Widget{
		g.Label("卡片与 Mod 配对检查"),
		g.Separator(),
		g.Label("提示: 可拖入文件或使用浏览按钮选择文件。"),
		g.Separator(),
		util.FilePicker("MOD 文件", &mmModPath).ForFiles().Filters("*.zipmod"),
		util.FilePicker("卡片文件", &mmCardPath).ForFiles().Filters("*.png"),
		g.Separator(),
		g.Row(
			g.Align(g.AlignCenter).To(g.Button(func() string {
				if mmInProgress {
					return "正在比对..."
				}
				return "比对"
			}()).Size(180, 60).Disabled(mmInProgress).OnClick(func() {
				if mmInProgress {
					return
				}
				startModMatch()
			})),
			g.Spacing(),
			g.Align(g.AlignCenter).To(g.Button("清除").Size(120, 60).OnClick(func() {
				mmModPath = ""
				mmCardPath = ""
				mmStatus = ""
				mmLog = ""
				g.Update()
			})),
		),
		g.Separator(),
		g.Label(func() string {
			if mmStatus == "" {
				return "状态: 等待操作"
			}
			return "状态: " + mmStatus
		}()),
		g.InputTextMultiline(&mmLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsAllowTabInput | g.InputTextFlagsNoHorizontalScroll),
	}

	return g.TabItem(" 卡片与Mod配对 ").Layout(layout...)
}

func startModMatch() {
	mmInProgress = true
	mmStatus = "准备比对..."
	mmLog = ""
	g.Update()

	go func() {
		defer func() {
			mmInProgress = false
			g.Update()
		}()

		// 清理双引号
		cardpath := strings.ReplaceAll(mmCardPath, "\"", "")
		modpath := strings.ReplaceAll(mmModPath, "\"", "")

		cpext := strings.ToLower(filepath.Ext(cardpath))
		mpext := strings.ToLower(filepath.Ext(modpath))

		if cpext != ".png" {
			mmStatus = "错误: 卡片路径请输入 PNG 文件"
			mmLog += mmStatus + "\n"
			g.Update()
			return
		}
		if mpext != ".zipmod" {
			mmStatus = "错误: MOD 路径请输入 zipmod 文件"
			mmLog += mmStatus + "\n"
			g.Update()
			return
		}

		mmStatus = "正在读取卡片..."
		g.Update()
		card, err := ic.ReadCardFromPath(cardpath)
		if err != nil {
			mmStatus = "读取卡片失败"
			mmLog += err.Error() + "\n"
			g.Update()
			return
		}

		mmStatus = "正在解析 MOD..."
		g.Update()
		mod, err := util.ReadZip("", modpath, 0)
		if err != nil {
			mmStatus = "读取 MOD 失败"
			mmLog += err.Error() + "\n"
			g.Update()
			return
		}

		deps := card.GetZipmodsDependencies()
		found := false
		for _, guid := range deps {
			if guid == mod.GUID {
				found = true
				break
			}
		}

		if found {
			mmStatus = "该卡使用了此 MOD"
			mmLog += "比对结果: 该卡使用了此 MOD\n"
		} else {
			mmStatus = "该卡没有使用此 MOD"
			mmLog += "比对结果: mod 没有被这张卡片使用\n"
		}

		// 添加一些 MOD / 卡片 信息用于诊断
		mmLog += "--- 详细信息 ---\n"
		mmLog += "卡片路径: " + card.GetPath() + "\n"
		mmLog += "MOD 路径: " + mod.Path + "\n"
		mmLog += "MOD GUID: " + mod.GUID + "\n"
		mmLog += "卡片引用的 GUID 列表:\n"
		for _, gguid := range deps {
			mmLog += "  " + gguid + "\n"
		}
		g.Update()
	}()
}
