package content_page

import (
	"KKCardModCheck/config"
	"KKCardModCheck/giu/custom"
	"KKCardModCheck/util"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	g "github.com/AllenDang/giu"
)

var (
	modBuildInProgress bool
	modBuildProgress   float32
	modBuildStatus     string
	modBuildLog        string
	// 当没有检测到任何 .zipmod 时置为 true，以便在 UI 中提供返回 Welcome 的选项
	modBuildNoZipmod bool
	// 解析失败的文件路径列表（构建后用于一键移动）
	modBuildFailedFiles []string
	// 正在移动损坏 mod 的标志
	movingBadMods bool
)

// startModBuild 启动一个后台协程来执行构建流程（此处为示例性的模拟进度）。
func startModBuild() {
	if modBuildInProgress {
		return
	}
	modBuildInProgress = true
	// 重置无zipmod标志以及失败列表
	modBuildNoZipmod = false
	modBuildFailedFiles = nil
	modBuildProgress = 0
	modBuildStatus = "准备构建..."
	modBuildLog = ""

	go func() {
		defer func() {
			modBuildInProgress = false
			g.Update()
		}()

		gamePath := config.Instance.GamePath
		if gamePath == "" {
			modBuildStatus = "未设置游戏路径，请先在设置中选择游戏目录。"
			modBuildLog += modBuildStatus + "\n"
			g.Update()
			return
		}

		files := util.GetAllFiles(gamePath, ".zipmod")
		total := len(files)
		if total == 0 {
			modBuildStatus = "未检测到任何 .zipmod 文件。"
			modBuildLog += modBuildStatus + "\n"
			// 标记没有检测到 zipmod，供界面显示返回 Welcome 的按钮
			modBuildNoZipmod = true
			g.Update()
			return
		}

		var mods []util.ModXml
		var failmods []util.ModXml

		for i, p := range files {
			// process each file
			modBuildStatus = fmt.Sprintf("解析 %d/%d: %s", i+1, total, p)
			g.Update()

			m, err := util.ReadZip("", p, i)
			if err != nil {
				failmods = append(failmods, m)
				// 记录失败文件路径到全局列表，供一键移动使用
				modBuildFailedFiles = append(modBuildFailedFiles, p)
				// 使用当前文件在总列表中的位置(i+1)作为编号，而不是 failmods 的长度，避免从 1 单独计数
				modBuildLog += fmt.Sprintf("[%d/%d] Read Fail: %s\n Error: %s\n", i+1, total, m.Path, err.Error())
			} else {
				// 如果命令行参数中带 -p，则计算哈希
				if len(os.Args) > 1 && os.Args[1] == "-p" {
					if hash, err := util.GetFileHash64(p); err == nil {
						m.MD5 = fmt.Sprintf("%x", hash)
					}
				}
				mods = append(mods, m)

				// 剔除换行符，避免日志格式混乱
				name := util.RemoveNewlines(m.Name)
				guid := util.RemoveNewlines(m.GUID)
				version := util.RemoveNewlines(m.Version)
				// 使用当前文件在总列表中的位置(i+1)作为编号，便于与失败日志一致
				modBuildLog += fmt.Sprintf("[%d/%d] %s\t%s\t%s\n", i+1, total, name, guid, version)
			}

			modBuildProgress = float32(i+1) / float32(total)
			g.Update()
		}

		// 写入 json 文件
		eyt, _ := json.MarshalIndent(failmods, "", "  ")
		byt, err := json.MarshalIndent(mods, "", "  ")
		if err != nil {
			modBuildStatus = "序列化 Mod 数据失败。"
			modBuildLog += modBuildStatus + "\n"
			g.Update()
			return
		}

		modsPath := filepath.Join(config.GetDataDir(), "ModsInfo.json")
		if err := os.WriteFile(modsPath, byt, 0644); err != nil {
			modBuildLog += fmt.Sprintf("写入 %s 失败: %v\n", modsPath, err)
		}

		if len(failmods) > 0 {
			failPath := filepath.Join(config.GetDataDir(), "ModsInfoFail.json")
			if err := os.WriteFile(failPath, eyt, 0644); err != nil {
				modBuildLog += fmt.Sprintf("写入 %s 失败: %v\n", failPath, err)
			} else {
				modBuildLog += fmt.Sprintf("相关解析失败Mod已写入 %s\n", failPath)
			}
		}

		// 在日志末尾统一列出所有解析失败的文件，便于查看和复制
		if len(failmods) > 0 {
			modBuildLog += fmt.Sprintf("失败文件汇总（共%d个）：\n", len(failmods))
			for i, fm := range failmods {
				modBuildLog += fmt.Sprintf("[%d/%d] %s\n", i+1, len(failmods), fm.Path)
			}
			modBuildLog += "\n"
		}
		config.Instance.ModInfoPath = modsPath
		// 更新配置
		config.Instance.ModInfoBuild = true
		config.Instance.ModInfoBuildTime = time.Now().Format(time.RFC3339)
		config.Save()

		modBuildStatus = fmt.Sprintf("检测完成，共检测到%d个MOD，成功%d个，失败%d个", total, len(mods), len(failmods))
		modBuildLog += modBuildStatus + "\n"
		g.Update()
	}()
}

// moveFailedMods 将解析失败的 mod 文件移动到数据目录下的 BadMods 子目录，异步执行并写入日志
func moveFailedMods() {
	if movingBadMods || len(modBuildFailedFiles) == 0 {
		return
	}
	movingBadMods = true
	g.Update()

	go func() {
		destDir := filepath.Join(config.GetDataDir(), "BadMods")
		if err := os.MkdirAll(destDir, 0755); err != nil {
			modBuildLog += fmt.Sprintf("创建目录 %s 失败: %v\n", destDir, err)
			movingBadMods = false
			g.Update()
			return
		}

		for i, p := range modBuildFailedFiles {
			if p == "" {
				continue
			}
			base := filepath.Base(p)
			dst := filepath.Join(destDir, base)

			// 先尝试重命名（效率高），若失败则回退到拷贝+删除的方式
			if err := os.Rename(p, dst); err != nil {
				in, errOpen := os.Open(p)
				if errOpen != nil {
					modBuildLog += fmt.Sprintf("[%d/%d] 移动失败: 无法打开源文件 %s: %v\n", i+1, len(modBuildFailedFiles), p, errOpen)
					continue
				}
				out, errCreate := os.Create(dst)
				if errCreate != nil {
					modBuildLog += fmt.Sprintf("[%d/%d] 移动失败: 无法创建目标文件 %s: %v\n", i+1, len(modBuildFailedFiles), dst, errCreate)
					in.Close()
					continue
				}
				if _, errCopy := io.Copy(out, in); errCopy != nil {
					modBuildLog += fmt.Sprintf("[%d/%d] 移动失败: 拷贝 %s -> %s 失败: %v\n", i+1, len(modBuildFailedFiles), p, dst, errCopy)
					out.Close()
					in.Close()
					continue
				}
				out.Close()
				in.Close()
				if errRem := os.Remove(p); errRem != nil {
					modBuildLog += fmt.Sprintf("[%d/%d] 移动警告: 删除源文件 %s 失败: %v\n", i+1, len(modBuildFailedFiles), p, errRem)
				} else {
					modBuildLog += fmt.Sprintf("[%d/%d] 已移动: %s -> %s\n", i+1, len(modBuildFailedFiles), p, dst)
				}
			} else {
				modBuildLog += fmt.Sprintf("[%d/%d] 已移动: %s -> %s\n", i+1, len(modBuildFailedFiles), p, dst)
			}
			g.Update()
		}

		// 移动完成后清空列表
		modBuildFailedFiles = nil
		movingBadMods = false
		g.Update()
	}()
}

func ModInfoBuild() g.Layout {
	// 顶部标题/说明部分
	header := g.Layout{
		g.Label("Mod信息构建"),
		g.Label("需要通过扫描 Mod 目录生成一个 Mod 信息文件，该文件用于辅助判断卡片所需的 Mod 是否存在"),
		g.Label("如果你添加了新的 Mod，记得重新生成 Mod 信息文件。"),
		g.Separator(),
		g.Label("可选配置:"),
		g.Checkbox("启动时自动检测 Mod 目录变化，如果有新 Mod 自动进行增量构建", &config.Instance.AutoDetectModChanges),
		g.Separator(),
	}

	// 按钮和进度显示
	action := g.Layout{
		g.Row(
			g.Align(g.AlignCenter).To(

				// 按钮组：第一个为构建按钮，第二个（开始使用）在配置允许时才显示
				func() g.Widget {
					// 构建按钮始终显示
					buildBtn := g.Button(func() string {
						if modBuildInProgress {
							return "正在构建..."
						}
						if config.Instance.ModInfoBuild {
							return "重新构建"
						}
						return "开始构建"
					}()).Size(200, 80).Disabled(modBuildInProgress).OnClick(func() {
						if !modBuildInProgress {
							startModBuild()
						}
					})

					// 如果没有检测到 zipmod，则在构建按钮旁显示“返回配置”按钮（样式与“开始使用”保持一致）
					var backBtn g.Widget
					if modBuildNoZipmod {
						backBtn = g.Button("返回配置").Size(200, 80).OnClick(func() {
							// 把流程重置为未完成初始配置，以便显示 Welcome 界面
							config.Instance.InitConfig = false
							// 隐藏模式构建界面的标记（可选，保证回到初始配置流程）
							config.Instance.HideModeInfoBuild = false
							config.Save()
							// 复位标志
							modBuildNoZipmod = false
						})
					}

					// 仅当已构建（ModInfoBuild 为真）且未被隐藏（HideModeInfoBuild 为假）时显示“开始使用”按钮
					if config.Instance.ModInfoBuild && !config.Instance.HideModeInfoBuild {
						startUseBtn := custom.ButtonStyleBlue.To(
							g.Button("开始使用").Size(200, 80).OnClick(func() {
								config.Instance.HideModeInfoBuild = true
							}),
						)
						// 可能存在的移动损坏 Mod 按钮
						var moveBtn g.Widget
						if len(modBuildFailedFiles) > 0 {
							moveBtn =
								g.Button(func() string {
									if movingBadMods {
										return "移动中..."
									}
									return "移动损坏Mod"
								}()).Size(200, 80).Disabled(movingBadMods).OnClick(func() {
									if !movingBadMods {
										moveFailedMods()
									}
								})
						}
						// 组合按钮显示顺序：构建, 返回（如果有）, 移动失败Mod（如果有）, 开始使用
						if backBtn != nil && moveBtn != nil {
							return g.Row(buildBtn, backBtn, moveBtn, startUseBtn)
						}
						if backBtn != nil {
							return g.Row(buildBtn, backBtn, startUseBtn)
						}
						if moveBtn != nil {
							return g.Row(buildBtn, moveBtn, startUseBtn)
						}
						return g.Row(buildBtn, startUseBtn)
					}

					// 其它情况：如果有 backBtn 或 moveBtn 一起显示，否则只显示构建按钮
					var moveBtn g.Widget
					if len(modBuildFailedFiles) > 0 {
						moveBtn = g.Button(func() string {
							if movingBadMods {
								return "移动中..."
							}
							return "移动损坏Mod"
						}()).Size(200, 80).Disabled(movingBadMods).OnClick(func() {
							if !movingBadMods {
								moveFailedMods()
							}
						})
					}
					if backBtn != nil && moveBtn != nil {
						return g.Row(buildBtn, backBtn, moveBtn)
					}
					if backBtn != nil {
						return g.Row(buildBtn, backBtn)
					}
					if moveBtn != nil {
						return g.Row(buildBtn, moveBtn)
					}
					return g.Row(buildBtn)
				}(),
			),
			g.Spacing(),
			// g.Button("打开 Mod 目录").OnClick(func() {
			// 	// 可选：在此处添加打开目录的实现
			// }),

		),
		g.Spacing(),
		g.ProgressBar(modBuildProgress).Size(-1, 0),
		g.Label(func() string {
			if modBuildStatus == "" {
				return "构建过程可能需要一些时间，请耐心等待..."
			}
			return modBuildStatus
		}()),
	}

	// 注意：当没有检测到 .zipmod 时，提示信息通过 modBuildStatus 展示，返回按钮已并入顶部按钮组，避免重复输出。

	return g.Layout{
		g.Child().Border(true).Layout(
			header,
			g.Separator(),
			action,
			g.Separator(),
			// 日志区域：只读可复制文本容器

			g.InputTextMultiline(&modBuildLog).
				Size(-1, -1).
				Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsAllowTabInput|g.InputTextFlagsNoHorizontalScroll),
		),
	}
}
