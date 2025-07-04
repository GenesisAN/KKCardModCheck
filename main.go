package main

import (
	"KKCardModCheck/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	ic "github.com/GenesisAN/illusionsCard"
	icb "github.com/GenesisAN/illusionsCard/Base"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tcnksm/go-latest"
)

const version = "0.1.9"

var app *tview.Application

var gamePath = "./"

func main() {
	if util.IsWin() && len(os.Args) == 1 {
		err := util.NoMoreDoubleClick()
		if err != nil {
			fmt.Printf("遇到错误: %v", err)
			time.Sleep(time.Second * 5)
		}
		os.Exit(0)
	}
	// 检查是否在Koikatu根目录下运行
	if len(os.Args) > 1 && os.Args[1] != "-p" {
		if util.IsNotExist("./Koikatu.exe") {
			fmt.Printf("请将软件放入Koikatu游戏根目录后再运行（暂时不支持在KKS下运行）")
			os.Exit(0)
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "-p" {
		gamePath = os.Args[2]
	}

	// 主线程发生恐慌时清空控制台，让控制台能更好的显示报错信息，只对主线程生效
	defer RecoverClearCMD()

	app = tview.NewApplication()
	pages := tview.NewPages()
	textView := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	var 所有统计 int
	// 更新检测

	githubTag := &latest.GithubTag{
		Owner:      "GenesisAN",
		Repository: "KKCardModCheck",
	}

	res, targer := latest.Check(githubTag, version)
	var newVersionText string
	if targer != nil {
		newVersionText = fmt.Sprintf("%s (*无法检测版本,请检查https://github.com/GenesisAN/KKCardModCheck访问是否顺畅)", version)
	} else {
		if res.Outdated {
			newVersionText = fmt.Sprintf("%s (*有新版本:%s)", version, res.Current)
		} else {
			newVersionText = fmt.Sprintf("v%s", version)
		}
	}

	textView.SetBorder(true)
	mods := []util.ModXml{}
	failmods := []util.ModXml{}
	lostmodname := make(map[string]icb.ResolveInfo)
	list := tview.NewList().
		AddItem("生成游戏MOD数据文件", "生成的数据保存在本软件目录下的ModsInfo.json文件内", 'b', func() {
			pages.SwitchToPage("MOD读取页")
			textView.Clear()
			go func() {
				mods = []util.ModXml{}
				fs := util.GetAllFiles(gamePath, ".zipmod")
				所有统计 = len(fs)

				if 所有统计 == 0 {
					util.OKMsg(pages, "未检测到任何mod!", "主页")
					return
				}

				for i, v := range fs {
					mod, err := util.ReadZip("", v, i)
					if err != nil {
						failmods = append(failmods, mod)
						fmt.Fprintf(textView, "[%d/%d]Read Fail:%s\nError:%s\n", len(failmods), 所有统计, mod.Path, err.Error())
					} else {
						//计算md5
						if os.Args[1] == "-p" {
							hash, err := util.GetFileHash64(v)
							if err != nil {
								panic(err)
							}
							mod.MD5 = fmt.Sprintf("%x", hash)
						}
						mods = append(mods, mod)
						fmt.Fprintf(textView, "[%d/%d]%s\t%s\t%s\n", len(mods), 所有统计, mod.Name, mod.GUID, mod.Version)
					}
				}

				eyt, err := json.Marshal(failmods)
				byt, err := json.Marshal(mods)
				if err != nil {
					return
				}
				os.WriteFile("./ModsInfo.json", byt, 0644)
				os.WriteFile("./ModsInfoFail.json", eyt, 0644)
				tips := fmt.Sprintf("检测完成，共检测到%d个MOD，成功%d个，失败%d个\n成功的数据已写入ModsInfo.json文件\n失败的数据已写入ModsInfoFail.json文件\n", 所有统计, len(mods), len(failmods))
				util.OKMsg(pages, tips, "主页")
				app.Draw()
			}()
		}).
		AddItem("批量检查卡片缺失MOD", "输入卡片文件夹和MOD数据文件进行对比，统计所有卡片缺失的Mod", 'f', func() {
			checkAllCardMods(pages, textView, lostmodname)
		}).
		AddItem("单个检查卡片缺失MOD", "根据卡片使用的MOD和MOD数据文件进行对比，统计单个卡片缺失的Mod", 'c', func() {
			checkSingCardMods(pages, lostmodname)
		}).
		AddItem("卡片MOD匹配", "查看下载的MOD是否是卡片所使用的Mod", 'p', func() {
			checkCardUseMod(pages)
		}).
		//AddItem("[未完成]从本机分享mod信息到服务器", "该功能用于获取本地存在，但服务器没有的mod，并调用BadiduPCS.exe上传到百度云，造福其他用户", 'u', nil).
		// AddItem("<试验性>从服务器获取缺失mod的秒传信息", "尝试获取缺失mod的秒传地址，自己通过百度云下载", 'g', func() {
		// 	tryGetBDMD5(pages)
		// }).
		AddItem("生成角色卡借物表", "尝通过本地MOD信息，生成自制角色的MOD借物表", 'm', func() {
			GetCardModsInfo(pages)
		})
	//AddItem("[未完成]用BadiduPCS直接下载缺失mod", "尝试获取服务器中缺失mod的秒传地址，并调用BadiduPCS.exe下载", 'd', nil)

	if res.Outdated {
		list.AddItem(fmt.Sprintf("下载新版本:v%s", res.Current), "喂!三点几，饮茶先啦！(哼啊啊啊啊...)", 'd', func() {
			if util.IsWin() {
				util.MsgWeb(pages, fmt.Sprintf("当前软件版本：%s,Github最新Release版本:%s", version, res.Current), "主页", "访问下载页", "https://github.com/GenesisAN/KKCardModCheck/releases")
				return
			}
			util.OKMsg(pages, "Github最新Release版本下载地址:https://github.com/GenesisAN/KKCardModCheck/releases", "主页")
		})
	}
	list.AddItem("关于本软件", "本软件完全免费，具体社区信息请进入查看", 'a', func() {
		if util.IsWin() {
			util.MsgWeb(pages, "本软件由G_AN开发，欢迎入群讨论交流", "主页", "入群", "https://qm.qq.com/q/uzGnetojdY")
			return
		}
		util.OKMsg(pages, "本软件由G_AN开发，欢迎入群讨论交流:625083414", "主页")
	}).
		AddItem("退出", "按下退出程序", 'q', func() {
			app.Stop()
		})

	pages.AddPage("主页", list, true, true)
	pages.AddPage("MOD读取页", textView, true, false)

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	//menu := newPrimitive("主菜单")
	//main := newPrimitive("Main content")
	//sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive(fmt.Sprintf("KK角色卡MOD缺失检测工具 %s", newVersionText)), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("软件为免费发布，至于为什么是控制台UI？不觉得这很酷吗....很符合.....后面忘了"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(pages, 0, 0, 0, 0, 0, 0, false)
	//AddItem(main, 1, 0, 1, 3, 0, 0, false).
	//AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.c
	grid.AddItem(pages, 1, 0, 1, 3, 0, 0, false)
	//AddItem(main, 1, 1, 1, 1, 0, 100, false).
	//AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).SetFocus(pages).Run(); err != nil {
		fmt.Print("\033[H\033[2J")
		panic(err)
	}
}

// 恐慌时清空CMD界面
func RecoverClearCMD() {
	if x := recover(); x != nil {
		fmt.Print("\033[H\033[2J")
	}
}

// 检查单个卡片缺失mod
func checkSingCardMods(pages *tview.Pages, _ map[string]icb.ResolveInfo) {
	if util.IsNotExist("./ModsInfo.json") {
		util.OKMsg(pages, "未找到 ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}

	cd := tview.NewInputField().
		SetLabel("输入卡片路径,或将卡片拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)

	cd.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cardPath := strings.Replace(cd.GetText(), "\"", "", -1)
			if path.Ext(cardPath) != ".png" {
				util.OKMsg(pages, "请输入PNG格式的卡片文件路径！", "路径输入")
				return
			}

			// 读取本地 ModsInfo.json
			localGUIDs, err := util.LoadModGUIDsFromJSON("./ModsInfo.json")
			if err != nil {
				util.OKMsg(pages, fmt.Sprintf("读取 MOD 文件失败：%s", err.Error()), "路径输入")
				return
			}

			// 解析卡片
			card, err := ic.ReadCardFromPath(cardPath)
			if err != nil {
				util.OKMsg(pages, fmt.Sprintf("卡片读取失败: %s", err.Error()), "路径输入")
				return
			}
			fmt.Println("→ 解析卡片成功:", card.GetPath())
			// 输出卡片信息（可选）
			//card.PrintCardInfo()
			// 缺失 MOD 比对
			var missing map[string]icb.ResolveInfo
			if checker, ok := card.(interface {
				CompareMissingMods([]string) map[string]icb.ResolveInfo
			}); ok {
				missing = checker.CompareMissingMods(localGUIDs)
			} else {
				util.OKMsg(pages, "当前卡片类型不支持MOD依赖解析", "路径输入")
				return
			}

			// 写入或提示
			if len(missing) > 0 {
				var lines []string
				for guid := range missing {
					lines = append(lines, guid)
				}
				if err := os.WriteFile("mods.txt", []byte(strings.Join(lines, "\n")), 0644); err != nil {
					util.OKMsg(pages, fmt.Sprintf("生成 mods.txt 文件失败: %s", err.Error()), "主页")
				} else {
					util.MsgWeb(pages,
						fmt.Sprintf("检索完成，共缺失 %d 个 MOD，已写入 mods.txt", len(missing)),
						"主页", "查看缺失mod", "mods.txt")
				}
			} else {
				util.OKMsg(pages, "检索完成，没有缺失的MOD！", "主页")
			}
		} else if key == tcell.KeyESC {
			pages.SwitchToPage("主页")
		}
	})

	pages.AddPage("路径输入", cd, true, false)
	pages.SwitchToPage("路径输入")
}

// // 检查所有卡片缺失mod
func checkAllCardMods(pages *tview.Pages, textView *tview.TextView, lostmodname map[string]icb.ResolveInfo) {
	if util.IsNotExist("./ModsInfo.json") {
		util.OKMsg(pages, "未找到 ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}

	textView.Clear()
	cd := tview.NewInputField().
		SetLabel("输入卡片文件夹路径,或将文件夹拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)

	cd.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			// 清空旧数据
			lostmodname = make(map[string]icb.ResolveInfo)

			cardPath := strings.Replace(cd.GetText(), "\"", "", -1)
			if path.Ext(cardPath) != "" || cardPath == "" {
				util.OKMsg(pages, "请输入合法的文件夹路径！", "路径输入")
				return
			}

			// 加载本地 MOD GUID 列表
			localGUIDs, err := util.LoadModGUIDsFromJSON("./ModsInfo.json")
			if err != nil {
				util.OKMsg(pages, fmt.Sprintf("读取 ModsInfo.json 失败: %s", err.Error()), "路径输入")
				return
			}

			pages.SwitchToPage("MOD读取页")
			textView.Clear()

			go func() {
				files := util.GetAllFiles(cardPath, ".png")
				var cards []icb.CardInterface

				for i, f := range files {
					card, err := ic.ReadCardFromPath(f)
					if err != nil {
						fmt.Fprintf(textView, "[%d/%d] 解析失败: %s\n", i+1, len(files), err.Error())
						continue
					}
					cards = append(cards, card)
					fmt.Fprintf(textView, "[%d/%d] 成功读取: %s\n", i+1, len(files), card.GetPath())
				}

				missing := util.CollectMissingMods(cards, localGUIDs)

				if len(missing) > 0 {
					var lines []string
					for guid := range missing {
						lines = append(lines, guid)
					}

					if err := os.WriteFile("mods.txt", []byte(strings.Join(lines, "\n")), 0644); err != nil {
						util.OKMsg(pages, fmt.Sprintf("缺失MOD已识别，但写入 mods.txt 失败: %s", err.Error()), "主页")
						app.Draw()
						return
					}

					util.MsgWeb(
						pages,
						fmt.Sprintf("检索完成，共缺失 %d 个 MOD，已生成 mods.txt", len(missing)),
						"主页",
						"查看缺失mod",
						"mods.txt",
					)
				} else {
					util.OKMsg(pages, "检索完成，没有缺失的MOD!", "主页")
				}
				app.Draw()
			}()
		} else if key == tcell.KeyESC {
			pages.SwitchToPage("主页")
		}
	})

	pages.AddPage("路径输入", cd, true, false)
	pages.SwitchToPage("路径输入")
}

// 卡片路径，mod路径
var cp, mp string

// 检查卡片使用的mod
func checkCardUseMod(pages *tview.Pages) {
	f := tview.NewForm().
		AddTextView("Tips", "按Tab可切换选项，输入框可以拖入文件夹,或者png来快速填写,", 100, 1, true, false).
		AddInputField("MOD路径", mp, 100, nil, func(text string) { mp = text }).
		SetCancelFunc(func() { pages.SwitchToPage("主页") }).
		AddInputField("卡片路径", cp, 100, nil, func(text string) { cp = text }).
		AddButton("比对", func() {
			// 处理传入路径双引号处理
			cardpath := strings.Replace(cp, "\"", "", -1)
			modpath := strings.Replace(mp, "\"", "", -1)
			cpext := path.Ext(cardpath) // 提取路径后缀进行比对
			mpext := path.Ext(modpath)
			if cpext != ".png" {
				util.OKMsg(pages, "卡片路径请输入PNG文件", "卡MOD比对")
				return
			}
			if mpext != ".zipmod" {
				util.OKMsg(pages, "mod路径请输入文件zipmod文件", "卡MOD比对")
				return
			}
			// 读取KK卡片数据
			kkcard, err := ic.ReadCardFromPath(cardpath)
			if err != nil {
				util.OKMsg(pages, err.Error(), "卡MOD比对")
				return
			}
			mod, err := util.ReadZip("", modpath, 0)
			if err != nil {
				util.OKMsg(pages, err.Error(), "卡MOD比对")
				app.Draw()
				return
			}
			v := kkcard.GetZipmodsDependencies()

			for _, zipmod_guid := range v {
				//如果找到了,直接break终止遍历，就开始下一个Zipmod匹配
				if mod.GUID == zipmod_guid {
					cp, mp = "", ""
					pages.AddPage("弹窗",
						tview.NewModal().
							SetBackgroundColor(tcell.ColorBlack).
							SetText("这张卡使用了该Mod").
							AddButtons([]string{"继续比对", "返回主界面"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								if buttonIndex == 0 {
									mp = ""
									checkCardUseMod(pages)
								} else {
									cp, mp = "", ""
									pages.SwitchToPage("主页")
								}

							}),
						true,
						false)
					pages.SwitchToPage("弹窗")
					return
				}
			}

			mp = ""
			pages.AddPage("弹窗",
				tview.NewModal().
					SetBackgroundColor(tcell.ColorBlack).
					SetText("mod没有被这张卡片使用").
					AddButtons([]string{"继续比对", "返回主界面"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonIndex == 0 {
							mp = ""
							checkCardUseMod(pages)
						} else {
							cp, mp = "", ""
							pages.SwitchToPage("主页")
						}
					}),
				true,
				false)
			pages.SwitchToPage("弹窗")
		}).
		AddButton("返回", func() { pages.SwitchToPage("主页") })
	pages.AddPage("卡MOD比对",
		f,
		true,
		false)
	pages.SwitchToPage("卡MOD比对")
}

// 生成卡片借物表
func GetCardModsInfo(pages *tview.Pages) {
	if util.IsNotExist("./ModsInfo.json") {
		util.OKMsg(pages, "未找到ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}
	cd := tview.NewInputField().
		SetLabel("输入卡片路径,或将卡片拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)
	cd.SetDoneFunc(func(key tcell.Key) {
		// 按下回车的处理
		if key == tcell.KeyEnter {
			// ModsInfo文件读取
			data, err := os.ReadFile("./ModsInfo.json")
			if err != nil {
				util.OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			}
			// 反序列化Modsinfo.json的数据
			var modsinfo []util.ModXml
			json.Unmarshal(data, &modsinfo)

			var sb strings.Builder
			// 处理传入路径双引号处理
			cardpath := strings.Replace(cd.GetText(), "\"", "", -1)
			ext := path.Ext(cardpath) // 提取路径后缀进行比对
			sb.WriteString(cardpath + "\n\n")
			if ext == ".png" {
				// 读取KK卡片数据
				kkcard, err := ic.ReadCardFromPath(cardpath)
				if err != nil {
					util.OKMsg(pages, err.Error(), "路径输入")
					return
				}
				sb.WriteString("依赖DLL:\n")
				for _, dllname := range kkcard.GetDLLDependencies() {
					sb.WriteString(fmt.Sprintf(" - %s\n", dllname))
				}
				sb.WriteString("\n依赖ZipMod:\n")
				cardmods := make(map[string]string)
				localmods := make(map[string]util.ModXml)
				v := kkcard.GetZipmodsDependencies()
				bdmd5 := 0
				if len(v) > 0 {
					// 提取出卡片mod map
					for _, zipmod := range v {
						if _, OK := cardmods[zipmod]; !OK {
							cardmods[zipmod] = zipmod
						}
					}
					for _, cm := range modsinfo {
						if _, OK := localmods[cm.GUID]; !OK {
							localmods[cm.GUID] = cm
						}
					}
					for _, zipmod_guid := range cardmods {
						if lm, OK := localmods[zipmod_guid]; OK {
							bdmd5++
							//	sb.WriteString(fmt.Sprintf("%s:\n - 名称:%s\n - 作者:%s\n - 相关web:%s\n - 秒传地址:%s\n\n", lm.GUID, lm.Name, lm.Author, lm.Website /*lm.BDMD5*/))
							sb.WriteString(fmt.Sprintf("%s:\n - 名称:%s\n - 作者:%s\n - 相关web:%s\n\n\n", lm.GUID, lm.Name, lm.Author, lm.Website /*lm.BDMD5*/))
						} else {
							sb.WriteString(fmt.Sprintf("%s:\n - 本地无此MOD\n\n", zipmod_guid))
						}
					}
					//data, err := getData(fmt.Sprintf("https://anweb.asuscomm.com:3000/api/v1/mods/search?s=%s&p=0&t=0", zipmod.GUID), map[string]string{"Accept": "application/json"})
					//if err != nil {
					//	//fmt.Println(err.Error())
					//	util.OKMsg(pages, fmt.Sprintf("请求失败:%s", err.Error()), "主页")
					//	app.Draw()
					//	return
					//}
					//var res Response
					//json.Unmarshal(data, &res)
					//if res.Code == 200 {
					//	for _, datum := range res.Data {
					//		if datum.GUID == zipmod.GUID {
					//			bdmd5++
					//			sb.WriteString(fmt.Sprintf("%s:\n - 作者:%s\n - 相关web:%s\n - 秒传地址:%s\n", zipmod.GUID, datum.Author, datum.Website, datum.BDMD5))
					//		}
					//	}
					//} else {
					//	sb.WriteString(fmt.Sprintf("%s:暂未找到相关信息\n", zipmod.GUID))
					//	continue
					//}
				}
				os.WriteFile("card_use_mod.txt", []byte(sb.String()), 0777)
				pages.AddPage("弹窗",
					tview.NewModal().
						SetBackgroundColor(tcell.ColorBlack).
						SetText(fmt.Sprintf("检索完成，card_use_mod.txt文件，共从本地匹配到（%d/%d）个mod的相关信息!", bdmd5, len(cardmods))).
						AddButtons([]string{"返回主页", "查看卡借物表文件"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							if buttonIndex == 0 {
								pages.SwitchToPage("主页")
							} else if buttonIndex == 1 {
								exec.Command("cmd", "/C", "start", "card_use_mod.txt").Start()
							}
						}),
					true,
					false)
				pages.SwitchToPage("弹窗")
			} else { // 输入文件后缀错误
				util.OKMsg(pages, "请输入PNG文件!", "路径输入")
			}
		} else if key == tcell.KeyESC { // 如果按下ESC，返回主页
			pages.SwitchToPage("主页")
		}
	})
	pages.AddPage("路径输入",
		cd,
		true,
		false)
	pages.SwitchToPage("路径输入")
}

func tryGetBDMD5(pages *tview.Pages) {
	if util.IsNotExist("./mods.txt") {
		util.OKMsg(pages, "未找到mods.txt\n请先检查卡片缺失mod，生成mods.txt信息后再尝试", "主页")
		return
	}
	go func() {
		data, err := os.ReadFile("./mods.txt")
		if err != nil {
			util.OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			app.Draw()
			return
		}

		mods := strings.Split(string(data), "\n")
		var bdmd5 []string
		var nobdmd5 []string
		count := len(mods)
		for i, i2 := range mods {
			util.MsgTips(pages, fmt.Sprintf("正在向服务器请求mod信息[%d/%d]", i, count))
			app.Draw()
			data, err := getData(fmt.Sprintf("https://zjpfj.asuscomm.com:4040/api/v1/mods/search?s=%s&p=0&t=0", i2), map[string]string{"Accept": "application/json"})
			if err != nil {
				//fmt.Println(err.Error())
				util.OKMsg(pages, fmt.Sprintf("请求失败:%s", err.Error()), "主页")
				app.Draw()
				return
			}
			var res Response
			json.Unmarshal(data, &res)
			if res.Code == 200 {
				if res.Data[0].GUID == i2 /*&& res.Data[0].Upload*/ {
					bdmd5 = append(bdmd5, res.Data[0].BDMD5)
				}
			} else {
				nobdmd5 = append(nobdmd5, i2)
				continue
			}
		}
		//写入TXT文件
		os.WriteFile("bdmd5s.txt", []byte(strings.Join(bdmd5, "\n")), 0644)
		os.WriteFile("nobdmd5s.txt", []byte(strings.Join(nobdmd5, "\n")), 0644)
		pages.AddPage("弹窗",
			tview.NewModal().
				SetBackgroundColor(tcell.ColorBlack).
				SetText(fmt.Sprintf("检索完成，已生成bdmd5s.txt文件，共从服务器匹配到（%d/%d）个mod的秒传下载地址!", len(bdmd5), len(mods))).
				AddButtons([]string{"OK", "查看秒传文件", "查看无下载地址的mod"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonIndex == 0 {
						pages.SwitchToPage("主页")
					} else if buttonIndex == 1 {
						exec.Command("cmd", "/C", "start", "bdmd5s.txt").Start()
					} else {
						exec.Command("cmd", "/C", "start", "nobdmd5s.txt").Start()
					}
				}),
			true,
			false)
		pages.SwitchToPage("弹窗")
		app.Draw()
	}()
}

type ModInfo struct {
	GUID        string
	Version     string
	Name        string
	Author      string
	Description string
	Website     string
	Game        string
	BDMD5       string `gorm:"uniqueIndex"`
	Upload      bool
}

type Response struct {
	//错误代码
	Code int `json:"code"`
	//数据
	Data []ModInfo `json:"data,omitempty"`
	//返回消息
	Msg string `json:"msg"`
	//错误信息
	Error string `json:"error,omitempty"`
}

func getData(url string, head map[string]string) (data []byte, err error) {
	//提交请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	reqest, err := http.NewRequest("GET", url, nil)
	//增加header选项
	for i, v := range head {
		reqest.Header.Add(i, v)
	}
	if err != nil {
		return nil, err
	}
	//处理返回结果
	//发起http请求的client实例
	client := &http.Client{Transport: tr}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	data, err = io.ReadAll(response.Body)
	response.Body.Close()
	return data, err
	//func jsonSave(v interface{}, path string) {
	//	f, _ := os.Create(path)
	//	defer f.Close()
	//	json.NewEncoder(f).Encode(v)
	//}
}
