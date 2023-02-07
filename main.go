package main

import (
	card "KKCardModCheck/IllusionsCard"
	"KKCardModCheck/IllusionsCard/Base"
	"KKCardModCheck/IllusionsCard/KK"
	"archive/zip"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"github.com/rivo/tview"
	"github.com/tcnksm/go-latest"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
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

const version = "0.1.4"

var app *tview.Application

func main() {
	if isWin() && len(os.Args) == 1 {
		err := NoMoreDoubleClick()
		if err != nil {
			fmt.Printf("遇到错误: %v", err)
			time.Sleep(time.Second * 5)
		}
		os.Exit(0)
	}

	if IsNotExist("./Koikatu.exe") {
		fmt.Printf("请将软件放入Koikatu游戏根目录后再运行（暂时不支持在KKS下运行）")
		os.Exit(0)
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
	var all int
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
	mods := []ModXml{}
	lostmodname := make(map[string]Base.ResolveInfo)
	list := tview.NewList().
		AddItem("生成游戏MOD数据文件", "生成的数据保存在本软件目录下的ModsInfo.json文件内", 'b', func() {
			pages.SwitchToPage("MOD读取页")
			textView.Clear()
			go func() {
				mods = []ModXml{}
				fs := GetAllFiles(`./`, ".zipmod")
				all = len(fs)

				if all == 0 {
					OKMsg(pages, "未检测到任何mod!", "主页")
					return
				}

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
		AddItem("[未完成]从本机分享mod信息到服务器", "该功能用于获取本地存在，但服务器没有的mod，并调用BadiduPCS.exe上传到百度云，造福其他用户", 'u', nil).
		AddItem("<试验性>从服务器获取缺失mod的秒传信息", "尝试获取缺失mod的秒传地址，自己通过百度云下载", 'd', func() {
			tryGetBDMD5(pages)
		}).
		AddItem("[未完成]用BadiduPCS直接下载缺失mod", "尝试获取服务器中缺失mod的秒传地址，并调用BadiduPCS.exe下载", 'd', nil)

	if res.Outdated {
		list.AddItem(fmt.Sprintf("下载新版本:v%s", res.Current), "喂!三点几，饮茶先啦！(哼啊啊啊啊...)", 'a', func() {
			if isWin() {
				MsgWeb(pages, fmt.Sprintf("当前软件版本：%s,Github最新Release版本:%s", version, res.Current), "主页", "访问下载页", "https://github.com/GenesisAN/KKCardModCheck/releases")
				return
			}
			OKMsg(pages, "Github最新Release版本下载地址:https://github.com/GenesisAN/KKCardModCheck/releases", "主页")
		})
	}
	list.AddItem("关于本软件", "本软件完全免费，具体社区信息请进入查看", 'a', func() {
		if isWin() {
			MsgWeb(pages, "本软件由G_AN开发，欢迎入群讨论交流，入群方式可从doc.kkgkd.com获取", "主页", "访问网页", "https://doc.kkgkd.com")
			return
		}
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
		AddItem(newPrimitive(fmt.Sprintf("KK角色卡MOD缺失检测工具 %s", newVersionText)), 0, 0, 1, 3, 0, 0, false).
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
func checkSingCardMods(pages *tview.Pages, lostmodname map[string]Base.ResolveInfo) {
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
			lostmodname = make(map[string]Base.ResolveInfo)

			// ModsInfo文件读取
			data, err := os.ReadFile("./ModsInfo.json")
			if err != nil {
				OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			}
			// 反序列化Modsinfo.json的数据
			var modsinfo []ModXml
			json.Unmarshal(data, &modsinfo)
			// 处理传入路径双引号处理
			cardpath := strings.Replace(cd.GetText(), "\"", "", -1)
			ext := path.Ext(cardpath) // 提取路径后缀进行比对
			if ext == ".png" {
				// 读取KK卡片数据
				kkcard, err := card.ReadCardKK(cardpath)
				if err != nil {
					OKMsg(pages, fmt.Sprintf("%s", err.Error()), "路径输入")
					return
				}
				v, Ok := kkcard.ExtendedList["com.bepis.sideloader.universalautoresolver"]
				if Ok {
				NextZipmod:
					for _, zipmod := range v.RequiredZipmodGUIDs {
						// 是否找到了mod标志
						// 遍历本地modinfo的数据，进行对比
						for _, mod := range modsinfo {
							// 如果找到了,直接break终止遍历，就开始下一个Zipmod匹配
							if mod.GUID == zipmod.GUID {
								break NextZipmod
							}
						}
						// 没有break说明没有找到该mod,直接添加到map
						if _, OK := lostmodname[zipmod.GUID]; !OK {
							lostmodname[zipmod.GUID] = zipmod
						}
					}
				}
				if len(lostmodname) != 0 { // 当缺失mod数量不为0时，写入TXT，并弹出提示框
					var p []string
					// 提取map中的缺失mod名称
					for s := range lostmodname {
						p = append(p, s)
					}
					//写入TXT文件
					os.WriteFile("mods.txt", []byte(strings.Join(p, "\n")), 0644)
					MsgWeb(pages, fmt.Sprintf("检索完成，已生成mods.txt文件，共缺失%d个MOD!", len(lostmodname)), "主页", "查看缺失mod", "mods.txt")
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

// 检查所有卡片缺失mod
func checkAllCardMods(pages *tview.Pages, textView *tview.TextView, lostmodname map[string]Base.ResolveInfo) {
	if IsNotExist("./ModsInfo.json") {
		OKMsg(pages, "未找到ModsInfo.json\n请先生成游戏MOD数据文件", "主页")
		return
	}
	textView.Clear()
	cd := tview.NewInputField().
		SetLabel("输入卡片文件夹路径,或将文件夹拖入(回车确定，ESC返回): ").
		SetFieldWidth(100)
	cd.SetDoneFunc(func(key tcell.Key) {
		// 按下回车的处理
		if key == tcell.KeyEnter {
			// 初始化缺失mod map
			lostmodname = make(map[string]Base.ResolveInfo)
			// 处理传入路径双引号处理
			cardpath := strings.Replace(cd.GetText(), "\"", "", -1)
			// ModsInfo文件读取
			data, err := os.ReadFile("./ModsInfo.json")
			if err != nil {
				OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			}
			// 反序列化Modsinfo.json的数据
			var cards []*KK.KoiCard
			var modsinfo []ModXml
			var frc [][]string
			json.Unmarshal(data, &modsinfo)
			ext := path.Ext(cardpath) // 提取路径后缀进行比对
			if ext == "" && cardpath != "" {
				pages.SwitchToPage("MOD读取页")
				textView.Clear()
				go func() {
					fs := GetAllFiles(cardpath, ".png")
					for i, v := range fs {
						card, err := card.ReadCardKK(v)
						if err != nil {
							frc = append(frc, []string{cardpath, err.Error()})
							continue
						}
						cards = append(cards, card)
						fmt.Fprintf(textView, "[%d/%d]%s\n", i+1, len(fs), v)
					}
					for _, kkcard := range cards {
						v, Ok := kkcard.ExtendedList["com.bepis.sideloader.universalautoresolver"]
						if Ok {
						NextZipmod:
							for _, zipmod := range v.RequiredZipmodGUIDs {
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
						err := os.WriteFile("mods.txt", []byte(strings.Join(p, "\n")), 0777)
						if err != nil {
							OKMsg(pages, fmt.Sprintf("检索完成，但生成mods.txt文件失败，原因:%s!", err.Error()), "主页")
							app.Draw()
						} else {
							if isWin() {
								MsgWeb(pages, fmt.Sprintf("检索完成，已生成mods.txt文件，共缺失%d个MOD!", len(lostmodname)), "主页", "查看缺失mod", "mods.txt")
								app.Draw()
							} else {
								OKMsg(pages, fmt.Sprintf("检索完成，已生成mods.txt文件，共缺失%d个MOD!", len(lostmodname)), "主页")
								app.Draw()
							}
						}
					} else { // 缺失mod数量为0
						OKMsg(pages, "检索完成，没有缺失的MOD!", "主页")
						app.Draw()
					}
				}()

			} else { // 输入文件后缀错误
				OKMsg(pages, "请输入路径!", "路径输入")
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
				OKMsg(pages, "卡片路径请输入PNG文件", "卡MOD比对")
				return
			}
			if mpext != ".zipmod" {
				OKMsg(pages, "mod路径请输入文件zipmod文件", "卡MOD比对")
				return
			}
			// 读取KK卡片数据
			kkcard, err := card.ReadCardKK(cardpath)
			if err != nil {
				OKMsg(pages, fmt.Sprintf("%s", err.Error()), "卡MOD比对")
				return
			}
			mod, err := ReadZip("", modpath, 0)
			if err != nil {
				OKMsg(pages, fmt.Sprintf("%s", err.Error()), "卡MOD比对")
				app.Draw()
				return
			}
			v, Ok := kkcard.ExtendedList["com.bepis.sideloader.universalautoresolver"]
			if Ok {
				for _, zipmod := range v.RequiredZipmodGUIDs {
					//如果找到了,直接break终止遍历，就开始下一个Zipmod匹配
					if mod.GUID == zipmod.GUID {
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

func tryGetBDMD5(pages *tview.Pages) {
	if IsNotExist("./mods.txt") {
		OKMsg(pages, "未找到mods.txt\n请先检查卡片缺失mod，生成mods.txt信息后再尝试", "主页")
		return
	}
	go func() {
		data, err := os.ReadFile("./mods.txt")
		if err != nil {
			OKMsg(pages, fmt.Sprintf("文件读取失败:%s", err.Error()), "路径输入")
			app.Draw()
			return
		}

		mods := strings.Split(string(data), "\n")
		var bdmd5 []string
		var nobdmd5 []string
		count := len(mods)
		for i, i2 := range mods {
			MsgTips(pages, fmt.Sprintf("正在向服务器请求mod信息[%d/%d]", i, count))
			app.Draw()
			data, err := getData(fmt.Sprintf("https://anweb.asuscomm.com:3000/api/v1/mods/search?s=%s&p=0&t=0", i2), map[string]string{"Accept": "application/json"})
			if err != nil {
				//fmt.Println(err.Error())
				OKMsg(pages, fmt.Sprintf("请求失败:%s", err.Error()), "主页")
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
	//app.Draw()
}

func MsgWeb(pages *tview.Pages, text, page, button, url string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text).
			AddButtons([]string{"OK", button}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage(page)
				} else {
					exec.Command("cmd", "/C", "start", url).Start()
				}
			}),
		true,
		false)
	pages.SwitchToPage("弹窗")
}

func MsgTips(pages *tview.Pages, text string) {
	pages.AddPage("弹窗",
		tview.NewModal().
			SetBackgroundColor(tcell.ColorBlack).
			SetText(text),
		true,
		false)
	pages.SwitchToPage("弹窗")
}

func ReadZip(dst, src string, i int) (ModXml, error) {
	mod := ModXml{}
	//===============BDMD5计算===================
	lp, err := checksum.GetFileSum(src, checksum.CHECKSUM_MD5|checksum.CHECKSUM_SLICE_MD5|checksum.CHECKSUM_CRC32)
	if err != nil {
		return mod, err
	}
	strLength, strMd5, strSliceMd5 := strconv.FormatInt(lp.Length, 10), hex.EncodeToString(lp.MD5), hex.EncodeToString(lp.SliceMD5)
	fileName := filepath.Base(src)
	regFileName := strings.Replace(fileName, " ", "_", -1)
	regFileName = strings.Replace(regFileName, "#", "_", -1)
	bdmd5 := strMd5 + "#" + strSliceMd5 + "#" + strLength + "#" + regFileName
	//==========================================

	zr, err := zip.OpenReader(src) //打开modzip
	if err != nil {
		return mod, errors.New(fmt.Sprintf("打开zip失败:%s", err.Error()))
	}

	for _, v := range zr.File { //遍历里面的文件
		if v.FileInfo().Name() == "manifest.xml" { //找到manifest.xml
			fr, _ := v.Open()         //打开它
			data, _ := io.ReadAll(fr) //读取里面的所有内容
			//构造XML的结构体
			err = xml.Unmarshal(data, &mod) //内容反序列化，按指定格式拿到里面的数据
			if err != nil {
				return mod, errors.New(fmt.Sprintf("读取zipmod信息失败:%s", err.Error()))
			}

			if err != nil {
				fmt.Errorf(err.Error())
				return mod, err
			}
			mod.BDMD5 = bdmd5
			mod.Path = strings.Replace(src, "mods\\", "", -1)
		}
	}
	return mod, err
}
func GetAllFiles(root, ext string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ext {
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
func isWin() bool {
	return runtime.GOOS == "windows"
}

// NoMoreDoubleClick 提示用户不要双击运行，并生成安全启动脚本
func NoMoreDoubleClick() error {
	r := boxW(0, "请勿通过双击直接运行本程序\n这会因CMD字符格式问题导致界面渲染出现问题\n点击确认将释出启动脚本，点击取消则关闭程序", "警告", 0x00000030|0x00000001)
	if r == 2 {
		return nil
	}
	r = boxW(0, "点击确认将覆盖start.bat，点击取消则关闭程序", "警告", 0x00000030|0x00000001)
	if r == 2 {
		return nil
	}
	f, err := os.OpenFile("start.bat", os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Printf("打开start.bat失败: %v", err)
		return err
	}
	_ = f.Truncate(0)

	ex, _ := os.Executable()
	exPath := filepath.Base(ex)
	_, err = f.WriteString("%Created by KKCardModCheck. DO NOT EDIT ME!%\nchcp 65001 \n\"" + exPath + "\" -s\npause")
	if err != nil {
		fmt.Printf("写入启动.bat失败: %v", err)
		return err
	}
	f.Close()
	boxW(0, "启动脚本已生成，请双击start.bat启动", "提示", 0x00000040|0x00000000)
	return nil
}

// BoxW of Win32 API. Check https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw for more detail.
func boxW(hwnd uintptr, caption, title string, flags uint) int {
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(captionPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags))

	return int(ret)
}
