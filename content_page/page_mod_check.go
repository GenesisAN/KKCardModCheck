package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/util"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

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
var isChecking bool
var cancelRequested atomic.Bool
var progressParsed int
var progressErrors int
var checkMu sync.Mutex

func resetCheckStateLocked() {
	missingMods = nil
	missingModsText = ""
	usedMods = nil
	usedModsText = ""
	tips = ""
	checked = false
	isChecking = false
	progressParsed = 0
	progressErrors = 0
	cancelRequested.Store(false)
}

// 使用嵌套 Tab 展示“缺失的 Mod / 使用的 Mod”，默认先显示“缺失的 Mod”

func CardModVerification() *g.TabItemWidget {
	// 检查config.CardPath的后缀是否是.png，如果是则说明是单个卡片，否则是目录
	// 如果用户切换了卡片文件/目录（path变化），重置检查结果（隐藏子 Tab）
	if lastCardPath == "" {
		lastCardPath = config.Instance.CardPath
	} else if config.Instance.CardPath != lastCardPath {
		checkMu.Lock()
		if isChecking {
			cancelRequested.Store(true)
			tips = "路径变化，正在取消当前检查..."
		} else {
			resetCheckStateLocked()
		}
		checkMu.Unlock()
		lastCardPath = config.Instance.CardPath
	}
	checkMu.Lock()
	localTips := tips
	localChecked := checked
	localMissingText := missingModsText
	localUsedText := usedModsText
	localIsChecking := isChecking
	localProgressParsed := progressParsed
	localProgressErrors := progressErrors
	checkMu.Unlock()
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

	if localTips != "" {
		layout = append(layout,
			g.Align(g.AlignCenter).To(
				g.Label(localTips),
			),
		)
	}
	if localIsChecking {
		progressText := fmt.Sprintf("已解析 %d 张，错误 %d 个", localProgressParsed, localProgressErrors)
		layout = append(layout,
			g.Separator(),
			g.Align(g.AlignCenter).To(
				g.Label(progressText),
			),
			g.Align(g.AlignCenter).To(
				g.Button("取消").Size(120, 28).OnClick(func() {
					cancelRequested.Store(true)
					checkMu.Lock()
					tips = "正在取消，请稍候..."
					checkMu.Unlock()
					g.Update()
				}),
			),
		)
	}
	// 使用嵌套 TabBar 展示子页面：先显示“缺失的 Mod”，每个子页内包含复制/清除按钮
	if localChecked {
		layout = append(layout,
			g.TabBar().TabItems(
				g.TabItem(" 缺失的 Mod ").Layout(
					g.Row(
						g.Button("清除结果").Size(120, 28).OnClick(func() {
							checkMu.Lock()
							missingMods = nil
							missingModsText = ""
							usedMods = nil
							usedModsText = ""
							tips = ""
							checked = false
							checkMu.Unlock()
							g.Update()
						}),
						g.Button("复制到剪贴板").Size(140, 28).OnClick(func() {
							checkMu.Lock()
							text := missingModsText
							checkMu.Unlock()
							if text == "" {
								return
							}
							checkMu.Lock()
							tips = "复制中..."
							checkMu.Unlock()
							g.Update()
							go func(payload string) {
								if err := util.CopyToClipboard(payload); err != nil {
									fmt.Printf("复制到剪切板失败: %v\n", err)
									checkMu.Lock()
									tips = "复制到剪贴板失败"
									checkMu.Unlock()
								} else {
									checkMu.Lock()
									tips = "已复制到剪贴板"
									checkMu.Unlock()
								}
								g.Update()
							}(text)
						}),
					),
					g.InputTextMultiline(&localMissingText).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
				),
				g.TabItem(" 使用的 Mod ").Layout(
					g.Row(
						g.Button("清除结果").Size(120, 28).OnClick(func() {
							checkMu.Lock()
							missingMods = nil
							missingModsText = ""
							usedMods = nil
							usedModsText = ""
							tips = ""
							checked = false
							checkMu.Unlock()
							g.Update()
						}),
						g.Button("复制到剪贴板").Size(140, 28).OnClick(func() {
							checkMu.Lock()
							text := usedModsText
							checkMu.Unlock()
							if text == "" {
								return
							}
							checkMu.Lock()
							tips = "复制中..."
							checkMu.Unlock()
							g.Update()
							go func(payload string) {
								if err := util.CopyToClipboard(payload); err != nil {
									fmt.Printf("复制到剪切板失败: %v\n", err)
									checkMu.Lock()
									tips = "复制到剪贴板失败"
									checkMu.Unlock()
								} else {
									checkMu.Lock()
									tips = "已复制到剪贴板"
									checkMu.Unlock()
								}
								g.Update()
							}(text)
						}),
					),
					g.InputTextMultiline(&localUsedText).Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
				),
			),
		)
	}
	return g.TabItem(" 卡片Mod缺失检测 ").Layout(layout...)
}

func SingleCardCheck() {
	checkMu.Lock()
	if isChecking {
		tips = "正在检查，请稍候..."
		checkMu.Unlock()
		g.Update()
		return
	}
	isChecking = true
	checked = false
	progressParsed = 0
	progressErrors = 0
	cancelRequested.Store(false)
	tips = "开始解析..."
	checkMu.Unlock()
	g.Update()

	go func() {
		localGUIDs, err := util.LoadModGUIDsFromJSON(config.Instance.ModInfoPath)
		if err != nil {
			checkMu.Lock()
			tips = fmt.Sprintf("读取 MOD 文件失败：%s", err.Error())
			isChecking = false
			checkMu.Unlock()
			g.Update()
			return
		}

		ext := filepath.Ext(config.Instance.CardPath)
		missing := make(map[string]Base.ResolveInfo)
		usedSet := make(map[string]struct{})
		usedList := make([]string, 0, 128)
		parsedCount := 0
		errorCount := 0
		foundCount := 0
		updateEvery := 20

		updateProgress := func(force bool) {
			if !force && parsedCount%updateEvery != 0 {
				return
			}
			checkMu.Lock()
			progressParsed = parsedCount
			progressErrors = errorCount
			tips = fmt.Sprintf("正在解析...已处理 %d 张", parsedCount)
			checkMu.Unlock()
			g.Update()
		}

		processCard := func(card Base.CardInterface) {
			if card == nil {
				return
			}
			parsedCount++
			for _, guid := range card.GetZipmodsDependencies() {
				if guid == "" {
					continue
				}
				if _, ok := usedSet[guid]; !ok {
					usedSet[guid] = struct{}{}
					usedList = append(usedList, guid)
				}
			}
			if comparer, ok := card.(interface {
				CompareMissingMods([]string) map[string]Base.ResolveInfo
			}); ok {
				for guid, info := range comparer.CompareMissingMods(localGUIDs) {
					missing[guid] = info
				}
			}
			updateProgress(false)
		}

		if ext == ".png" {
			card, err := ic.ReadCardFromPath(config.Instance.CardPath)
			if err != nil {
				checkMu.Lock()
				tips = fmt.Sprintf("卡片读取失败: %s", err.Error())
				isChecking = false
				checkMu.Unlock()
				g.Update()
				return
			}
			fmt.Println("→ 解析卡片成功:", card.GetPath())
			processCard(card)
		} else if ext == "" {
			if !util.IsExist(config.Instance.CardPath) {
				checkMu.Lock()
				tips = "所选目录不存在，无法解析"
				isChecking = false
				checkMu.Unlock()
				g.Update()
				return
			}

			errCanceled := errors.New("canceled")
			walkErr := filepath.Walk(config.Instance.CardPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if cancelRequested.Load() {
					return errCanceled
				}
				if info == nil || info.IsDir() {
					return nil
				}
				if filepath.Ext(path) != ".png" {
					return nil
				}
				foundCount++
				card, err := ic.ReadCardFromPath(path)
				if err != nil {
					fmt.Printf("解析卡片失败 (%s)：%v\n", path, err)
					errorCount++
					updateProgress(false)
					return nil
				}
				fmt.Println("→ 解析卡片成功:", card.GetPath())
				processCard(card)
				if cancelRequested.Load() {
					return errCanceled
				}
				return nil
			})

			if walkErr != nil {
				if errors.Is(walkErr, errCanceled) {
					checkMu.Lock()
					tips = "已取消"
					isChecking = false
					checkMu.Unlock()
					g.Update()
					return
				}
				checkMu.Lock()
				tips = fmt.Sprintf("目录遍历失败：%v", walkErr)
				isChecking = false
				checkMu.Unlock()
				g.Update()
				return
			}

			if foundCount == 0 {
				checkMu.Lock()
				tips = "目录中未找到任何 .png 卡片"
				isChecking = false
				checkMu.Unlock()
				g.Update()
				return
			}
		} else {
			checkMu.Lock()
			tips = "所选路径不是 .png 文件或目录，无法解析"
			isChecking = false
			checkMu.Unlock()
			g.Update()
			return
		}

		if cancelRequested.Load() {
			checkMu.Lock()
			tips = "已取消"
			isChecking = false
			checkMu.Unlock()
			g.Update()
			return
		}

		var missingLines []string
		for guid, mod := range missing {
			if guid == "" {
				missingLines = append(missingLines, fmt.Sprintf("%s (GUID未知)", mod.Property))
			} else {
				missingLines = append(missingLines, guid)
			}
		}
		sort.Strings(missingLines)
		sort.Strings(usedList)

		checkMu.Lock()
		missingMods = missingLines
		missingModsText = strings.Join(missingLines, "\n")
		usedMods = usedList
		usedModsText = strings.Join(usedList, "\n")
		if len(missingLines) > 0 {
			tips = fmt.Sprintf("检索完成，共解析 %d 张卡片，缺失 %d 个 MOD", parsedCount, len(missingLines))
		} else {
			tips = fmt.Sprintf("检索完成，共解析 %d 张卡片，没有缺失的MOD！", parsedCount)
		}
		checked = true
		isChecking = false
		progressParsed = parsedCount
		progressErrors = errorCount
		checkMu.Unlock()
		g.Update()
	}()
}
