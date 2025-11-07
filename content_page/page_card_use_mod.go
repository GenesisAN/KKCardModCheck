package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
	ic "github.com/GenesisAN/illusionsCard"
)

var (
	cuCardPath   string
	cuStatus     string
	cuLog        string
	cuInProgress bool
	cuOutputPath string
)

// ModeCardUseTab 生成卡片借物表的 GUI 实现
func ModeCardUseTab() *g.TabItemWidget {
	if cuCardPath == "" {
		cuCardPath = config.Instance.CardPath
	}

	layout := []g.Widget{
		g.Label("根据卡片生成借物表（card_use_mod.txt）"),
		g.Separator(),
		util.FilePicker("卡片文件", &cuCardPath).ForFiles().Filters("*.png"),
		g.Separator(),
		g.Row(
			g.Align(g.AlignCenter).To(g.Button(func() string {
				if cuInProgress {
					return "正在生成..."
				}
				return "生成借物表"
			}()).Size(180, 60).Disabled(cuInProgress).OnClick(func() {
				if cuInProgress {
					return
				}
				startGenerateCardMods()
			})),
			g.Spacing(),
			g.Align(g.AlignCenter).To(g.Button("打开借物表文件").Size(140, 60).OnClick(func() {
				if cuOutputPath == "" {
					return
				}
				util.OpenFile(cuOutputPath)
			})),
			g.Spacing(),
			g.Align(g.AlignCenter).To(g.Button("清除").Size(100, 60).OnClick(func() {
				cuCardPath = ""
				cuStatus = ""
				cuLog = ""
				cuOutputPath = ""
				g.Update()
			})),
		),
		g.Separator(),
		g.Label(func() string {
			if cuStatus == "" {
				return "状态: 等待操作"
			}
			return "状态: " + cuStatus
		}()),
		g.InputTextMultiline(&cuLog).Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsAllowTabInput | g.InputTextFlagsNoHorizontalScroll),
	}

	return g.TabItem(" 角色卡借物表生成 ").Layout(layout...)
}

func startGenerateCardMods() {
	cuInProgress = true
	cuStatus = "准备生成..."
	cuLog = ""
	g.Update()

	go func() {
		defer func() {
			cuInProgress = false
			g.Update()
		}()

		// 确定 ModsInfo.json 路径：优先使用配置中的路径
		modsInfoPath := config.Instance.ModInfoPath
		if modsInfoPath == "" {
			// fallback to appDir/ModsInfo.json
			modsInfoPath = filepath.Join(appDir, "ModsInfo.json")
		}

		if util.IsNotExist(modsInfoPath) {
			cuStatus = "错误: 未找到 ModsInfo.json，请先生成 Mod 信息文件"
			cuLog += fmt.Sprintf("未找到: %s\n", modsInfoPath)
			g.Update()
			return
		}

		data, err := os.ReadFile(modsInfoPath)
		if err != nil {
			cuStatus = "错误: 读取 ModsInfo.json 失败"
			cuLog += err.Error() + "\n"
			g.Update()
			return
		}

		var modsinfo []util.ModXml
		if err := json.Unmarshal(data, &modsinfo); err != nil {
			cuStatus = "错误: 解析 ModsInfo.json 失败"
			cuLog += err.Error() + "\n"
			g.Update()
			return
		}

		// 清理输入路径
		cardpath := strings.ReplaceAll(cuCardPath, "\"", "")
		ext := strings.ToLower(filepath.Ext(cardpath))
		var sb strings.Builder
		sb.WriteString(cardpath + "\n\n")

		if ext != ".png" {
			cuStatus = "错误: 请输入 PNG 卡片文件"
			cuLog += "请输入 PNG 卡片文件\n"
			g.Update()
			return
		}

		cuStatus = "正在读取卡片..."
		g.Update()
		kkcard, err := ic.ReadCardFromPath(cardpath)
		if err != nil {
			cuStatus = "错误: 读取卡片失败"
			cuLog += err.Error() + "\n"
			g.Update()
			return
		}

		sb.WriteString("依赖DLL:\n")
		for _, dllname := range kkcard.GetDLLDependencies() {
			sb.WriteString(fmt.Sprintf(" - %s\n", dllname))
		}
		sb.WriteString("\n依赖ZipMod:\n")

		cardmods := make(map[string]struct{})
		localmods := make(map[string]util.ModXml)

		v := kkcard.GetZipmodsDependencies()
		matched := 0
		if len(v) > 0 {
			for _, zipmod := range v {
				cardmods[zipmod] = struct{}{}
			}
			for _, cm := range modsinfo {
				if cm.GUID != "" {
					localmods[cm.GUID] = cm
				}
			}
			for zipmod_guid := range cardmods {
				if lm, ok := localmods[zipmod_guid]; ok {
					matched++
					sb.WriteString(fmt.Sprintf("%s:\n - 名称:%s\n - 作者:%s\n - 相关web:%s\n\n", lm.GUID, lm.Name, lm.Author, lm.Website))
				} else {
					sb.WriteString(fmt.Sprintf("%s:\n - 本地无此MOD\n\n", zipmod_guid))
				}
			}
		}

		// 生成输出内容并显示在多行文本框
		outPath := filepath.Join(appDir, "card_use_mod.txt")
		content := sb.String()
		// 首先在 UI 中显示完整内容
		cuLog = content
		g.Update()

		// 尝试写入输出文件到 appDir（不成功也不影响已显示的内容）
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			cuStatus = "警告: 写入输出文件失败"
			cuLog += "\n" + err.Error() + "\n"
			g.Update()
			return
		}

		cuOutputPath = outPath

		cuStatus = fmt.Sprintf("完成: 共从本地匹配到 (%d/%d) 个 mod 的相关信息", matched, len(cardmods))
		// 在已有内容下追加状态信息
		//cuLog = content + "\n\n" + cuStatus + "\n"
		g.Update()
	}()
}
