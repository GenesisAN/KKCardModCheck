package util

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

// FilePickerWidget 封装的文件选择器组件，遵循giu的Widget接口
type FilePickerWidget struct {
	readOnly       bool
	label          string
	path           *string
	buttonLabel    string
	dialogTitle    string
	filters        []string
	labelWidth     float32
	buttonWidth    float32
	minInputWidth  float32
	dualMode       bool   // 是否启用双按钮模式
	fileButtonText string // 文件按钮文本
	dirButtonText  string // 文件夹按钮文本
	Click          func()
}

// FilePicker 创建一个统一的文件/文件夹选择器组件
func FilePicker(label string, path *string) *FilePickerWidget {
	return &FilePickerWidget{
		label:          label,
		path:           path,
		buttonLabel:    "浏览",
		dialogTitle:    "选择",
		labelWidth:     120,
		buttonWidth:    144,
		minInputWidth:  200,
		dualMode:       false,
		fileButtonText: "文件",
		dirButtonText:  "文件夹",
	}
}

// ReadOnly keeps informational paths selectable without allowing edits.
func (fp *FilePickerWidget) ReadOnly() *FilePickerWidget {
	fp.readOnly = true
	return fp
}

// ButtonLabel 设置按钮文本
func (fp *FilePickerWidget) ButtonLabel(label string) *FilePickerWidget {
	fp.buttonLabel = label
	return fp
}

// DialogTitle 设置对话框标题
func (fp *FilePickerWidget) DialogTitle(title string) *FilePickerWidget {
	fp.dialogTitle = title
	return fp
}

// Filters 设置文件过滤器 (仅对文件选择有效)
func (fp *FilePickerWidget) Filters(filters ...string) *FilePickerWidget {
	fp.filters = filters
	return fp
}

// ForFiles 配置为文件选择模式
func (fp *FilePickerWidget) ForFiles() *FilePickerWidget {
	fp.buttonLabel = "选择文件"
	fp.dialogTitle = "选择文件"
	fp.dualMode = false
	return fp
}

// ForDirectories 配置为文件夹选择模式
func (fp *FilePickerWidget) ForDirectories() *FilePickerWidget {
	fp.buttonLabel = "选择文件夹"
	fp.dialogTitle = "选择文件夹"
	fp.dualMode = false
	return fp
}

// WithDualMode 配置为双模式（文件+文件夹两个按钮）
func (fp *FilePickerWidget) WithDualMode() *FilePickerWidget {
	fp.dualMode = true
	// 设置更合理的默认按钮文本
	if fp.fileButtonText == "文件" {
		fp.fileButtonText = "选文件"
	}
	if fp.dirButtonText == "文件夹" {
		fp.dirButtonText = "选文件夹"
	}
	return fp
}

// FileButtonText 设置文件按钮的文本（仅在双模式下有效）
func (fp *FilePickerWidget) FileButtonText(text string) *FilePickerWidget {
	fp.fileButtonText = text
	return fp
}

// DirButtonText 设置文件夹按钮的文本（仅在双模式下有效）
func (fp *FilePickerWidget) DirButtonText(text string) *FilePickerWidget {
	fp.dirButtonText = text
	return fp
}

// LabelWidth 设置标签宽度
func (fp *FilePickerWidget) LabelWidth(width float32) *FilePickerWidget {
	fp.labelWidth = width
	return fp
}

// ButtonWidth 设置按钮宽度
func (fp *FilePickerWidget) ButtonWidth(width float32) *FilePickerWidget {
	fp.buttonWidth = width
	return fp
}

// MinInputWidth 设置输入框最小宽度
func (fp *FilePickerWidget) MinInputWidth(width float32) *FilePickerWidget {
	fp.minInputWidth = width
	return fp
}

func (fp *FilePickerWidget) OnClick(f func()) *FilePickerWidget {
	// 占位函数，实际点击事件在按钮点击处理函数中调用
	fp.Click = f
	return fp
}

// Build 构建组件（实现g.Widget接口）
func (fp *FilePickerWidget) Build() {
	const gap float32 = 8
	const height float32 = 32
	available, _ := g.GetAvailableRegion()
	origin := g.GetCursorPos()
	buttonWidth := fp.buttonWidth
	if buttonWidth <= 0 {
		buttonWidth = 144
	}
	labelWidth := fp.labelWidth
	if labelWidth <= 0 {
		labelWidth = 120
	}
	buttons := 1
	if fp.dualMode {
		buttons = 2
	}
	inputWidth := available - labelWidth - gap - float32(buttons)*(buttonWidth+gap)
	// Stack the label above the controls only when the form becomes too narrow.
	controlY := origin.Y
	inputX := origin.X + int(labelWidth+gap)
	g.SetCursorPos(image.Pt(origin.X, origin.Y+3))
	g.AlignTextToFramePadding()
	g.Label(fp.label).Build()
	if inputWidth < fp.minInputWidth {
		controlY += int(height + gap)
		inputX = origin.X
		inputWidth = available - float32(buttons)*(buttonWidth+gap)
	}
	if inputWidth < 1 {
		inputWidth = 1
	}
	g.SetCursorPos(image.Pt(inputX, controlY+3))
	input := g.InputText(fp.path).Size(inputWidth)
	if fp.readOnly {
		input.Flags(g.InputTextFlagsReadOnly)
	}
	input.Build()
	buttonX := inputX + int(inputWidth+gap)
	g.SetCursorPos(image.Pt(buttonX, controlY))
	if fp.dualMode {
		g.Button(fp.fileButtonText).Size(buttonWidth, height).OnClick(fp.onFileButtonClick).Build()
		g.SetCursorPos(image.Pt(buttonX+int(buttonWidth+gap), controlY))
		g.Button(fp.dirButtonText).Size(buttonWidth, height).OnClick(fp.onDirButtonClick).Build()
	} else {
		click := fp.onBrowseClick
		if fp.Click != nil {
			click = fp.Click
		}
		g.Button(fp.buttonLabel).Size(buttonWidth, height).OnClick(click).Build()
	}
	// Register the complete row so the following field starts at the same column.
	g.SetCursorPos(origin)
	g.Dummy(available, float32(controlY-origin.Y)+height).Build()
}

// calculateButtonWidth 计算按钮宽度的辅助方法
func (fp *FilePickerWidget) calculateButtonWidth(text string) float32 {
	// 更精确的文本宽度计算
	runes := []rune(text)
	charCount := float32(len(runes))

	// 如果字符很少，使用较大的字符宽度；如果字符较多，使用较小的字符宽度
	var avgCharWidth float32
	if charCount <= 2 {
		avgCharWidth = 15 // 单个或双字符按钮需要更多空间
	} else {
		avgCharWidth = 12 // 多字符按钮可以紧凑一些
	}

	padding := float32(25)  // 按钮内边距
	minWidth := float32(40) // 最小按钮宽度

	calculatedWidth := charCount*avgCharWidth + padding
	if calculatedWidth < minWidth {
		calculatedWidth = minWidth
	}

	return calculatedWidth
}

// onBrowseClick 浏览按钮点击处理
func (fp *FilePickerWidget) onBrowseClick() {
	var selectedPath string
	var err error

	// 根据按钮标签或对话框标题判断选择类型
	if strings.Contains(fp.buttonLabel, "文件夹") || strings.Contains(fp.dialogTitle, "文件夹") {
		selectedPath, err = showDirectoryDialog(fp.dialogTitle, *fp.path)
	} else {
		selectedPath, err = showFileDialog(fp.dialogTitle, *fp.path, fp.filters)
	}

	if err == nil && selectedPath != "" {
		*fp.path = selectedPath
	}
}

// onFileButtonClick 文件按钮点击处理（双模式下使用）
func (fp *FilePickerWidget) onFileButtonClick() {
	selectedPath, err := showFileDialog("选择文件", *fp.path, fp.filters)
	if err == nil && selectedPath != "" {
		*fp.path = selectedPath
	}
}

// onDirButtonClick 文件夹按钮点击处理（双模式下使用）
func (fp *FilePickerWidget) onDirButtonClick() {
	selectedPath, err := showDirectoryDialog("选择文件夹", *fp.path)
	if err == nil && selectedPath != "" {
		*fp.path = selectedPath
	}
}

// showFileDialog 显示文件选择对话框
func showFileDialog(title, defaultPath string, filters []string) (string, error) {
	dlg := dialog.File().Title(title)

	// 设置默认路径 - 统一使用 SetStartDir
	if defaultPath != "" {
		if filepath.IsAbs(defaultPath) {
			// 如果是文件路径，使用其目录；如果是目录，直接使用
			if !isDirectory(defaultPath) {
				dlg = dlg.SetStartDir(filepath.Dir(defaultPath))
			} else {
				dlg = dlg.SetStartDir(defaultPath)
			}
		}
	}

	// 添加文件过滤器（同时添加大小写变体以兼容不同平台/对话框行为）
	if len(filters) > 0 {
		for _, filter := range filters {
			if filter == "*.*" || filter == "*" {
				dlg = dlg.Filter("所有文件", "*.*")
				continue
			}

			// 处理如 "*.txt" 格式的过滤器
			ext := strings.TrimPrefix(filter, "*")
			if ext != "" {
				displayName := strings.ToUpper(strings.TrimPrefix(ext, ".")) + " 文件"

				// 构建大小写变体，例如 "*.png;*.PNG"
				alt := ""
				if strings.HasPrefix(filter, "*.") {
					upper := "*." + strings.ToUpper(strings.TrimPrefix(filter, "*."))
					if upper != filter {
						alt = filter + ";" + upper
					}
				}

				pattern := filter
				if alt != "" {
					pattern = alt
				}

				dlg = dlg.Filter(displayName, pattern)
			}
		}
	}

	return dlg.Load()
}

// showDirectoryDialog 显示文件夹选择对话框
func showDirectoryDialog(title, defaultPath string) (string, error) {
	dlg := dialog.Directory().Title(title)

	// 统一使用 SetStartDir
	if defaultPath != "" {
		dlg = dlg.SetStartDir(defaultPath)
	}

	return dlg.Browse()
}

// isDirectory 简单检查路径是否为目录
func isDirectory(path string) bool {
	// 简单的启发式检查：如果没有扩展名，可能是目录
	return filepath.Ext(path) == ""
}

// 便捷函数

// FilePickerRow 创建文件选择行 (兼容旧版本)
func FilePickerRow(label string, path *string, filters ...string) *FilePickerWidget {
	fp := FilePicker(label, path).ForFiles()
	if len(filters) > 0 {
		fp = fp.Filters(filters...)
	}
	return fp
}

// DirectoryPickerRow 创建文件夹选择行 (兼容旧版本)
func DirectoryPickerRow(label string, path *string) *FilePickerWidget {
	return FilePicker(label, path).ForDirectories()
}

// DualPickerRow 创建双模式选择行（文件+文件夹两个按钮）
func DualPickerRow(label string, path *string, filters ...string) *FilePickerWidget {
	fp := FilePicker(label, path).WithDualMode()
	if len(filters) > 0 {
		fp = fp.Filters(filters...)
	}
	return fp
}

// 为了保持向后兼容，添加这些别名函数
// DirectoryPicker 文件夹选择器（别名）
func DirectoryPicker(label string, path *string) *FilePickerWidget {
	return FilePicker(label, path).ForDirectories()
}

// DualPicker 双模式选择器（别名）
func DualPicker(label string, path *string) *FilePickerWidget {
	return FilePicker(label, path).WithDualMode()
}

// CheckCardMods 检查卡片所需的 Mod（占位实现）
// 参数：cardPath - 卡片文件路径，modPath - Mod目录路径
// 返回：缺失的Mod列表，错误信息
func CheckCardMods(cardPath, modPath string) ([]string, error) {
	// TODO: 实现实际的卡片Mod检查逻辑
	// 这里提供一个占位实现，您需要根据实际需求来实现检查逻辑

	if cardPath == "" {
		return nil, fmt.Errorf("卡片路径不能为空")
	}

	if modPath == "" {
		return nil, fmt.Errorf("Mod路径不能为空")
	}

	// 占位返回，表示没有缺失的Mod
	// 实际实现中，您需要：
	// 1. 解析卡片文件，提取所需的Mod信息
	// 2. 扫描Mod目录，检查哪些Mod存在
	// 3. 返回缺失的Mod列表
	return []string{}, nil
}
