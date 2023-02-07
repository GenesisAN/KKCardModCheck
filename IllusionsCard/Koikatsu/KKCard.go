package Koikatsu

import (
	"KKCardModCheck/IllusionsCard/Base"
	"KKCardModCheck/IllusionsCard/Tools"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

type KoiCard struct {
	*Base.Card
	CharParmeter *KKChaFileParameter
}

func (card *KoiCard) KKChaFileParameterEx(cfp *KKChaFileParameter) {
	card.CharParmeter = cfp
	card.CharInfo = &Base.ChaFileParameterEx{}
	card.CharInfo.Lastname = cfp.Lastname
	card.CharInfo.Firstname = cfp.Firstname
	card.CharInfo.Version = cfp.Version
	card.CharInfo.Nickname = cfp.Nickname
}

func ParseKoiChara(buff *bytes.Buffer) (KoiCard, error) {
	kc := KoiCard{&Base.Card{}, &KKChaFileParameter{}}

	versionlen, err := Tools.BufRead(buff, 1, "BufRead Fail:Version len")
	if err != nil {
		return kc, err
	}

	_, err = Tools.BufRead(buff, int(versionlen[0]), "BufRead Fail:Version")
	if err != nil {
		return kc, err
	}

	FLBuf, err := Tools.BufRead(buff, 4, "BufRead Fail:Face img len")
	if err != nil {
		return kc, err
	}

	faceLength := binary.LittleEndian.Uint32(FLBuf)
	_, err = Tools.BufRead(buff, int(faceLength), "BufRead Fail:Face img")
	if err != nil {
		return kc, err
	}

	countBuf, err := Tools.BufRead(buff, 4, "BufRead Fail:Card BlockHeader len")
	if err != nil {
		return kc, err
	}

	var count = binary.LittleEndian.Uint32(countBuf)
	bytes, err := Tools.BufRead(buff, int(count), "BufRead Fail:Card BlockHeader")
	if err != nil {
		return kc, err
	}

	var bhls BlockHeaderListInfo
	_, err = bhls.UnmarshalMsg(bytes)
	if err != nil {
		return kc, errors.New("BlockHeaderListInfo Unmarshal Fail")
	}

	sort.SliceStable(bhls.LstInfo, func(i, j int) bool {
		return bhls.LstInfo[i].Pos < bhls.LstInfo[j].Pos
	})

	buff.Next(8)
	bhmap := make(map[string]*BlockHeader)
	for _, bh := range bhls.LstInfo {
		cBuf, err := Tools.BufRead(buff, int(bh.Size), "BufRead Fail:"+bh.Name)
		if err != nil {
			return kc, err
		}
		bh.Data = cBuf
		//encodedString := hex.EncodeToString(cBuf)
		//os.WriteFile(fmt.Sprintf("%s-hex.txt", bh.Name), []byte(encodedString), 0777)
		bhmap[bh.Name] = bh
	}
	var parameterBytes, extDataByte []byte
	//遍历 头部信息，获取Parameter位置
	if par, ok := bhmap["Parameter"]; ok {
		if par.Version != "0.0.5" {
			return kc, errors.New("BlockHeaderListInfo Unmarshal Fail")
		}
		parameterBytes = par.Data
	} else {
		return kc, errors.New("parameter not found")
	}

	//根据位置信息，反序列化 MsgPack 得到 KKChaFileParameter
	var Cfp KKChaFileParameter
	_, err = Cfp.UnmarshalMsg(parameterBytes)
	if err != nil {
		return kc, errors.New("KKChaFileParameter Unmarshal Fail")
	}
	kc.KKChaFileParameterEx(&Cfp)
	if par, ok := bhmap["KKEx"]; ok {
		extDataByte = par.Data
	}

	//根据KKEx位置信息，反序列化 得到 extDataO
	var extDataO Base.MapSArrayInterface

	extDataO.UnmarshalMsg(extDataByte)
	//exDataO处理后的数据 exData
	exData := make(map[string]*Base.PluginData)
	kc.Extended = exData
	//exData原始数据
	for S, v := range extDataO {
		if v != nil {
			//取出PluginData
			var pd Base.PluginData
			pd.Version = int(v[0].(int64))
			pd.Data = v[1]
			exData[S] = &pd
		}
	}
	// 遍历exData提取RequiredZipmodGUIDs
	exDataEx := make(map[string]*Base.PluginDataEx)
	for s, data := range exData {
		// 根据GUID找出插件对应的Data
		dex := data.DeserializeObjects()
		dex.Name = s
		dex.Version = data.Version
		exDataEx[dex.Name] = &dex
	}
	kc.ExtendedList = exDataEx
	return kc, nil
}

func (kc *KoiCard) PrintCardInfo() {
	fmt.Println("插件依赖:")
	for _, ex := range kc.ExtendedList {
		fmt.Printf("[插件]%s(版本:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}
