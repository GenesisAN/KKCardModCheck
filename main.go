package main

import (
	card "KKCardModCheck/card"
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ModXml struct {
	GUID        string `xml:"guid"`
	Version     string `xml:"version"`
	Name        string `xml:"name"`
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Website     string `xml:"website"`
	Game        string `xml:"game"`
	BDMD5       string `gorm:"uniqueIndex"`
	Path        string
	Upload      bool
}

func main() {
	//if !RunningByDoubleClick() {
	if runtime.GOOS == "windows" && len(os.Args) == 1 {
		err := NoMoreDoubleClick()
		if err != nil {
			fmt.Printf("遇到错误: %v", err)
			time.Sleep(time.Second * 5)
		}
		os.Exit(0)
	}
	//}
	app := tview.NewApplication()
	pages := tview.NewPages()
	textView := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	var all int
	//numSelections := 0
	//textView.SetDoneFunc(func(key tcell.Key) {
	//	currentSelection := textView.GetHighlights()
	//	if key == tcell.KeyEnter {
	//		if len(currentSelection) > 0 {
	//			textView.Highlight()
	//		} else {
	//			textView.Highlight("0").ScrollToHighlight()
	//		}
	//	} else if len(currentSelection) > 0 {
	//		index, _ := strconv.Atoi(currentSelection[0])
	//		if key == tcell.KeyTab {
	//			index = (index + 1) % numSelections
	//		} else if key == tcell.KeyBacktab {
	//			index = (index - 1 + numSelections) % numSelections
	//		} else {
	//			return
	//		}
	//		textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
	//	}
	//})
	textView.SetBorder(true)
	mods := []ModXml{}
	lostmodname := make(map[string]card.ResolveInfo)
	list := tview.NewList().
		AddItem("生成游戏MOD数据文件", "生成的数据保存在本软件目录下的ModsInfo.json文件内", 'b', func() {
			pages.SwitchToPage("MOD读取页")
			textView.Clear()
			go func() {
				mods = []ModXml{}
				fs := GetAllFiles(`./`)
				all = len(fs)
				for i, v := range fs {
					mod, err := ReadZip("", v, i)
					if err != nil {
						log.Fatalln(err)
					}
					mods = append(mods, mod)
					fmt.Fprintf(textView, "[%d/%d]%s\t%s\t%s\n", len(mods), all, mod.Name, mod.GUID, mod.Version)

				}
				byt, err := json.Marshal(mods)
				if err != nil {
					return
				}
				os.WriteFile("./ModsInfo.json", byt, 0644)
				OKMsg(pages, "游戏mod信息统计完成，已写入ModesInfo.json文件!", "主页")
			}()
		}).
		AddItem("检查单个卡片缺失MOD", "根据卡片使用的MOD和MOD数据文件进行对比，找出缺少的MOD", 'c', func() {
			checkSingCardMods(pages, lostmodname)
		}).
		AddItem("批量检查卡片缺失MOD", "输入卡片文件夹，统计所有卡片缺失的Mod", 'f', func() {

		}).
		AddItem("[未完成]卡片MOD匹配", "查看下载的MOD是否是卡片所使用的Mod", 'p', nil).
		AddItem("[未完成]上传Mod到百度云，并将秒传地址提交到服务器", "该功能用于补充服务器没有的mod信息，为其他玩家造福", 'u', nil).
		AddItem("[未完成]从服务器获取缺失mod的秒传信息", "将缺失的mod名称，在服务器数据库中检索，尝试获取它的秒传地址", 'd', nil).
		AddItem("关于本软件", "本软件完全免费，具体社区信息请进入查看", 'a', func() {
			OKMsg(pages, "本软件由G_AN开发，欢迎入群讨论交流，入群方式可从doc.kkgkd.com获取", "主页")
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
		AddItem(newPrimitive("KK角色卡MOD缺失检测工具 v1.2.1"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("软件为作者免费发布，至于为什么是控制台UI？不觉得这很酷吗....很符合.....后面忘了"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(pages, 0, 0, 0, 0, 0, 0, false)
	//AddItem(main, 1, 0, 1, 3, 0, 0, false).
	//AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.c
	grid.AddItem(pages, 1, 0, 1, 3, 0, 0, false)
	//AddItem(main, 1, 1, 1, 1, 0, 100, false).
	//AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func checkSingCardMods(pages *tview.Pages, lostmodname map[string]card.ResolveInfo) {
	if IsNotExist("./ModsInfo.json") {
		OKMsg(pages, "未找到ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}
	cd := tview.NewInputField().
		SetLabel("输入卡片路径,或将卡片拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)
	cd.SetDoneFunc(func(key tcell.Key) {
		// 按下回车的处理
		if key == tcell.KeyEnter {
			// 初始化缺失mod map
			lostmodname = make(map[string]card.ResolveInfo)
			// 处理传入路径双引号处理
			cardpath := strings.Replace(cd.GetText(), "\"", "", -1)
			// ModsInfo文件读取
			data, err := os.ReadFile("./ModsInfo.json")
			if err != nil {
				OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			}
			// 反序列化Modsinfo.json的数据
			var modsinfo []ModXml
			json.Unmarshal(data, &modsinfo)
			ext := path.Ext(cardpath) // 提取路径后缀进行比对
			if ext == ".png" {
				// 读取KK卡片数据
				kkcard, err := card.ReadCardKK(cardpath)
				if err != nil {
					OKMsg(pages, fmt.Sprintf("%s", err.Error()), "路径输入")
					return
				}
				// 遍历卡片数据
				for _, cardmod := range kkcard.ExtendedList {
					// 找到mod相关数据
					if cardmod.Name == "com.bepis.sideloader.universalautoresolver" {
						// 遍历卡片mod数据里的zipmod数组
					NextZipmod:
						for _, zipmod := range cardmod.RequiredZipmodGUIDs {
							//是否找到了mod标志
							// 遍历本地modinfo的数据，进行对比
							for _, mod := range modsinfo {
								//如果找到了,直接break终止遍历，就开始下一个Zipmod匹配
								if mod.GUID == zipmod.GUID {
									break NextZipmod
								}
							}
							//没有break说明没有找到该mod,直接添加到map
							if _, OK := lostmodname[zipmod.GUID]; !OK {
								lostmodname[zipmod.GUID] = zipmod
							}
						}
					}
				}

				if len(lostmodname) != 0 { //当缺失mod数量不为0时，写入TXT，并弹出提示框
					var p []string
					//提取map中的缺失mod名称
					for s := range lostmodname {
						p = append(p, s)
					}
					//写入TXT文件
					os.WriteFile("mods.txt", []byte(strings.Join(p, "\n")), 0644)
					OKMsg(pages, fmt.Sprintf("检索完成，已生成mods.txt文件，共缺失%d个MOD!", len(lostmodname)), "主页")
				} else { // 缺失mod数量为0
					OKMsg(pages, "检索完成，没有缺失的MOD!", "主页")
				}
			} else { // 输入文件后缀错误
				OKMsg(pages, "请输入PNG文件!", "路径输入")
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
func checkAllCardMods(pages *tview.Pages, lostmodname map[string]card.ResolveInfo) {
	if IsNotExist("./ModsInfo.json") {
		OKMsg(pages, "未找到ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}
	data, err := os.ReadFile("./ModsInfo.json")
	if err != nil {
		OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
	}
	cd := tview.NewInputField().
		SetLabel("输入人物卡文件夹路径,或将文件夹拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)
	cd.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			lostmodname = make(map[string]card.ResolveInfo)
			cardpath := strings.Replace(cd.GetText(), "\"", "", -1)
			var modsinfo []ModXml
			json.Unmarshal(data, &modsinfo)
			ext := path.Ext(cardpath)
			if ext == ".png" {
				kkcard, err := card.ReadCardKK(cardpath)
				if err != nil {
					OKMsg(pages, fmt.Sprintf("卡片读取失败:%s", err.Error()), "路径输入")
					return
				}
				for _, cardmod := range kkcard.ExtendedList {
					if cardmod.Name == "com.bepis.sideloader.universalautoresolver" {
						for _, zipmod := range cardmod.RequiredZipmodGUIDs {
							find := false
							for _, mod := range modsinfo {
								if mod.GUID == zipmod.GUID {
									find = true
								}
							}
							if !find {
								if _, OK := lostmodname[zipmod.GUID]; !OK {
									lostmodname[zipmod.GUID] = zipmod
								}
							}
						}
					}
				}

				if len(lostmodname) != 0 {
					var p []string
					for s, _ := range lostmodname {
						p = append(p, s)
					}
					os.WriteFile("mods.txt", []byte(strings.Join(p, "\n")), 0644)
					OKMsg(pages, fmt.Sprintf("检索完成，已生成mods.txt文件，共缺失%d个MOD!", len(lostmodname)), "主页")
				} else {
					OKMsg(pages, "检索完成，没有缺失的MOD!", "主页")
				}

			} else {
				OKMsg(pages, "请输入PNG文件!", "路径输入")
			}
		} else if key == tcell.KeyESC {
			pages.SwitchToPage("主页")
		}

	})
	pages.AddPage("路径输入",
		cd,
		true,
		false)
	pages.SwitchToPage("路径输入")
}

func OKMsg(pages *tview.Pages, text string, page string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage(page)
				}
			}),
		true,
		false)
	pages.SwitchToPage("弹窗")
}

func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{}
	zr, err := zip.OpenReader(src) //打开modzip
	if err != nil {
		return mod, err
	}
	for _, v := range zr.File { //遍历里面的文件
		if v.FileInfo().Name() == "manifest.xml" { //找到manifest.xml
			fr, _ := v.Open()         //打开它
			data, _ := io.ReadAll(fr) //读取里面的所有内容
			//构造XML的结构体
			err = xml.Unmarshal(data, &mod) //内容反序列化，按指定格式拿到里面的数据
			if err != nil {
				return mod, err
			}
			//md5c, errc := GetModMD5Code("sf", src) //获取秒传链接
			//if err != nil {
			//	fmt.Errorf(errc.Error())
			//	return mod, err
			//}
			//mod.BDMD5 = md5c
			mod.Path = src
		}
	}
	return mod, err
}
func GetAllFiles(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".zipmod" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// IsExist 文件/路径存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsNotExist 文件/路径不存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}
