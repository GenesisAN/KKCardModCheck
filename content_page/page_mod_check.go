package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"fmt"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
	ic "github.com/GenesisAN/illusionsCard"
	"github.com/GenesisAN/illusionsCard/Base"
)

var tips = ""
var missingMods []string
var missingModsText string
var usedMods []string
var usedModsText string
var checked bool        // 表示是否已执行过检查，未检查时不显示子 Tab
var lastCardPath string // 记录上一次选择的卡片路径，路径变化则重置检查结果

// 使用嵌套 Tab 展示“缺失的 Mod / 使用的 Mod”，默认先显示“缺失的 Mod”

func CardModVerification() *g.TabItemWidget {
	// 检查config.CardPath的后缀是否是.png，如果是则说明是单个卡片，否则是目录
	// 如果用户切换了卡片文件/目录（path变化），重置检查结果（隐藏子 Tab）
	if lastCardPath == "" {
		lastCardPath = config.Instance.CardPath
	} else if config.Instance.CardPath != lastCardPath {
		// 路径发生变化 -> 重置结果
		missingMods = nil
		missingModsText = ""
		usedMods = nil
		usedModsText = ""
		tips = ""
		checked = false
		lastCardPath = config.Instance.CardPath
	}
	layout := []g.Widget{
		g.Label(" 选择你要检查的卡片/卡片目录 "),
		g.Separator(),
		util.FilePicker("卡片文件/目录", &config.Instance.CardPath).WithDualMode().Filters("*.png"),
		g.Separator(),
		g.Align(g.AlignCenter).To(
			g.Button("检查").Size(200, 60).OnClick(SingleCardCheck),
		),
		g.Separator(),
	}
	// 仅在渲染阶段显示图片预览，不在渲染中执行检查（检查应当在点击 "检查" 时触发）
	if filepath.Ext(config.Instance.CardPath) == ".png" {
		layout = append(layout,
			g.Align(g.AlignCenter).To(g.ImageWithFile(config.Instance.CardPath).Size(256, 360)))
		layout = append(layout, g.Separator())
	}

	if tips != "" {
		layout = append(layout,
			g.Align(g.AlignCenter).To(
				g.Label(tips),
			),
		)
	}
	// 使用嵌套 TabBar 展示子页面：先显示“缺失的 Mod”，每个子页内包含复制/清除按钮
	if checked {
		layout = append(layout,
			g.TabBar().TabItems(
				g.TabItem(" 缺失的 Mod ").Layout(
					g.Row(
						g.Button("清除结果").Size(120, 28).OnClick(func() {
							missingMods = nil
							missingModsText = ""
							usedMods = nil
							usedModsText = ""
							tips = ""
							checked = false
							g.Update()
						}),
						g.Button("复制到剪贴板").Size(140, 28).OnClick(func() {
							if missingModsText == "" {
								return
							}
							if err := util.CopyToClipboard(missingModsText); err != nil {
								fmt.Printf("复制到剪切板失败: %v\n", err)
								tips = "复制到剪贴板失败"
							} else {
								tips = "已复制到剪贴板"
							}
							g.Update()
						}),
					),
					g.InputTextMultiline(&missingModsText).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
				),
				g.TabItem(" 使用的 Mod ").Layout(
					g.Row(
						g.Button("清除结果").Size(120, 28).OnClick(func() {
							missingMods = nil
							missingModsText = ""
							usedMods = nil
							usedModsText = ""
							tips = ""
							checked = false
							g.Update()
						}),
						g.Button("复制到剪贴板").Size(140, 28).OnClick(func() {
							if usedModsText == "" {
								return
							}
							if err := util.CopyToClipboard(usedModsText); err != nil {
								fmt.Printf("复制到剪切板失败: %v\n", err)
								tips = "复制到剪贴板失败"
							} else {
								tips = "已复制到剪贴板"
							}
							g.Update()
						}),
					),
					g.InputTextMultiline(&usedModsText).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
				),
			),
		)
	}
	return g.TabItem(" 卡片Mod缺失检测 ").Layout(layout...)
}

func SingleCardCheck() {

	localGUIDs, err := util.LoadModGUIDsFromJSON(config.Instance.ModInfoPath)
	if err != nil {
		tips = fmt.Sprintf("读取 MOD 文件失败：%s", err.Error())
		return
	}

	// 根据路径后缀决定解析模式：
	//  - 以 .png 结尾 -> 单卡解析
	//  - 无后缀 -> 视为目录，批量解析目录下所有 .png 文件
	//  - 其他有后缀 -> 不解析
	ext := filepath.Ext(config.Instance.CardPath)
	var cards []Base.CardInterface
	if ext == ".png" {
		// 单张卡片
		card, err := ic.ReadCardFromPath((config.Instance.CardPath))
		if err != nil {
			tips = fmt.Sprintf("卡片读取失败: %s", err.Error())
			return
		}
		fmt.Println("→ 解析卡片成功:", card.GetPath())
		cards = []Base.CardInterface{card}
	} else if ext == "" {
		// 目录模式 - 递归查找所有 png 文件并尝试解析
		if !util.IsExist(config.Instance.CardPath) {
			tips = "所选目录不存在，无法解析"
			return
		}
		files := util.GetAllFiles(config.Instance.CardPath, ".png")
		if len(files) == 0 {
			tips = "目录中未找到任何 .png 卡片"
			return
		}
		for _, f := range files {
			card, err := ic.ReadCardFromPath(f)
			if err != nil {
				// 解析单个卡片失败则记录到控制台并跳过
				fmt.Printf("解析卡片失败 (%s)：%v\n", f, err)
				continue
			}
			fmt.Println("→ 解析卡片成功:", card.GetPath())
			cards = append(cards, card)
		}
		if len(cards) == 0 {
			tips = "未能解析任何卡片，检查目录或卡片格式是否正确"
			return
		}
	} else {
		tips = "所选路径不是 .png 文件或目录，无法解析"
		return
	}

	missing := util.CollectMissingMods(cards, localGUIDs)
	// 聚合所有卡片引用的 Mod 并去重
	var allUsed []string
	for _, c := range cards {
		if c == nil {
			continue
		}
		allUsed = append(allUsed, c.GetZipmodsDependencies()...)
	}
	usedMods = util.DedupeStrings(allUsed)
	usedModsText = strings.Join(usedMods, "\n")
	//ch

	parsedCount := len(cards)
	if len(missing) > 0 {
		var lines []string
		for guid, mod := range missing {
			//如果guid为空则使用mod.Property
			if guid == "" {
				lines = append(lines, fmt.Sprintf("%s (GUID未知)", mod.Property))
			} else {
				lines = append(lines, guid)
			}
		}
		// 在 GUI 中显示缺失的 Mod，而不写文件或弹窗
		missingMods = lines
		missingModsText = strings.Join(lines, "\n")
		tips = fmt.Sprintf("检索完成，共解析 %d 张卡片，缺失 %d 个 MOD", parsedCount, len(lines))
		checked = true
		g.Update()
	} else {
		missingMods = nil
		missingModsText = ""
		tips = fmt.Sprintf("检索完成，共解析 %d 张卡片，没有缺失的MOD！", parsedCount)
		checked = true
		g.Update()
	}

}
