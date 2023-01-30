package card

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/vmihailenco/msgpack/v5"
)

type KoiCard struct {
	*Card
	CharParmeter *KoiChaFileParameter
}

func ParseKoiChara(buff *bytes.Buffer) (KoiCard, error) {
	kc := KoiCard{&Card{}, &KoiChaFileParameter{}}

	//版本号	//buff.Next(int(versionlen[0]))
	versionlen, err := BufRead(buff, 1, "读取失败：卡片版本无法识别")
	if err != nil {
		return kc, err
	}
	_, err = BufRead(buff, int(versionlen[0]), "读取失败：卡片版本无法读取")
	if err != nil {
		return kc, err
	}

	//kc.LoadVersion = "123"
	//fmt.Println("版本:", version)
	FLBuf, err := BufRead(buff, 4, "读取失败：卡片头长度无法识别")
	if err != nil {
		return kc, err
	}
	faceLength := binary.LittleEndian.Uint32(FLBuf)

	_, err = BufRead(buff, int(faceLength), "读取失败：卡片头长度无法识别")
	if err != nil {
		return kc, err
	}

	countBuf, err := BufRead(buff, 4, "读取失败：卡片头长度无法识别")
	if err != nil {
		return kc, err
	}
	var count = binary.LittleEndian.Uint32(countBuf)

	bytes, err := BufRead(buff, int(count), "读取失败：卡片头无法读取")
	if err != nil {
		return kc, err
	}
	//BlockHeader 获取卡片结构头
	var bh map[string][]interface{}
	err = msgpack.Unmarshal(bytes, &bh)
	if err != nil {
		fmt.Print(err)
	}
	var bhinfos []BlockHeaderInfo
	for s, infov := range bh {
		if s == "lstInfo" {
			for _, v := range infov {
				var info BlockHeaderInfo
				mapstructure.Decode(v, &info)
				bhinfos = append(bhinfos, info)
			}
		}
	}
	buff.Next(8)
	//num2 := binary.LittleEndian.Uint64(buff.Next(8)) //int64
	//fmt.Print(num2)
	//cap := buff.Cap()
	//len := buff.Len()
	//pos := cap - len - 1

	var P_Pos, P_Size, K_Pos int
	var parameterBytes, extDataByte []byte
	//fmt.Println("卡片头部信息:")
	//遍历 头部信息，获取Parameter位置
	for _, v := range bhinfos {
		if v.Name == "Parameter" {
			if v.Version == "0.0.5" {
				P_Pos = int(v.Pos)
				P_Size = int(v.Size)
				buff.Next(P_Pos)
				parameterBytes = buff.Next(int(v.Size))
			}
		}
		//fmt.Println(v.Name, "    [Pos:", int(v.Pos), "Size:", int(v.Size), "]")
	}

	//根据位置信息，反序列化 MsgPack 得到 KoiChaFileParameter
	var Cfp KoiChaFileParameter
	var cfp interface{}
	err = msgpack.Unmarshal(parameterBytes, &cfp)
	mapstructure.Decode(cfp, &Cfp)

	kc.CharParmeter = &Cfp
	kc.ChaFileParameterEx(kc)

	//遍历头部信息，获取KKEx位置
	for _, v := range bhinfos {
		if v.Name == "KKEx" {
			K_Pos = int(v.Pos)
			buff.Next(K_Pos - (P_Size + P_Pos))
			extDataByte = buff.Next(int(v.Size))
			//mt.Print(extDataByte)
		}
	}

	//根据KKEx位置信息，反序列化 得到 extDataO
	var extDataO map[string]interface{}
	err = msgpack.Unmarshal(extDataByte, &extDataO)

	//exDataO处理后的数据 exData
	exData := make(map[string]*PluginData)
	kc.Extended = exData
	//exData原始数据
	for S, v := range extDataO {
		//取出PluginData
		var pd PluginData
		if v != nil {
			vlist := v.([]interface{})
			pd.Version = int(vlist[0].(int8))
			pd.Data = vlist[1].(map[string]interface{})
			pd.RequiredPluginGUIDs = []string{}
			pd.RequiredZipmodGUIDs = []string{}
			exData[S] = &pd
		}
	}

	// 遍历exData提取RequiredZipmodGUIDs
	exDataEx := make(map[string]*PluginDataEx)
	for s, data := range exData {
		// 根据GUID找出插件对应的Data
		dex := DeserializeObjects(data)
		dex.Name = s
		dex.Version = data.Version
		exDataEx[dex.Name] = &dex
		//exDataEx = append(exDataEx, &dex)
	}
	kc.ExtendedList = exDataEx
	return kc, nil
}

func DeserializeObjects(data *PluginData) PluginDataEx {
	var pluginDataEx PluginDataEx
	var resolveInfos []ResolveInfo
	ds := data.Data.(map[string]interface{})
	//提取 data中的info信息
	for s2, i := range ds {
		if s2 == "info" {
			bts := i.([]interface{})
			//将info内的[]byte数组，反序列化为ResolveInfo
			for _, bt := range bts {
				var bti map[string]interface{}
				msgpack.Unmarshal(bt.([]byte), &bti)
				var ri ResolveInfo
				//从中提取ResolveInfoEx
				mapstructure.Decode(bti, &ri)
				//将提取的ResolveInfoEx放入pluginDataEx
				resolveInfos = append(resolveInfos, ri)
				pluginDataEx.RequiredZipmodGUIDs = append(pluginDataEx.RequiredZipmodGUIDs, ri)
			}
		}
	}
	return pluginDataEx
}

func (kc *KoiCard) PrintCardInfo() {
	fmt.Println("插件依赖:")
	for _, ex := range kc.ExtendedList {
		fmt.Printf("[插件]%s(版本:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}
