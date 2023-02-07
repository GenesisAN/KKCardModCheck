// 角色卡插件数据结构
package Base

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

// PluginData 插件数据
type PluginData struct {
	Version int         //版本
	Data    interface{} //原始MsgPack数据
}

// PluginDataEx 插件数据扩展
type PluginDataEx struct {
	Version             int           //版本
	Name                string        //名称
	RequiredPluginGUIDs []string      //依赖插件
	RequiredZipmodGUIDs []ResolveInfo //依赖模组
}

// PrintMod 打印模组信息
func (pde *PluginDataEx) PrintMod() {

	if pde.RequiredZipmodGUIDs == nil || len(pde.RequiredZipmodGUIDs) == 0 {
		return
	}
	fmt.Println("插件内容依赖:")
	for i, i2 := range pde.RequiredZipmodGUIDs {
		fmt.Printf("  *[mod依赖%d]:%s (%s|LS:%d|CN:%d)\n", i, i2.GUID, i2.Property, i2.LocalSlot, i2.CategoryNo)
	}
}
func (data *PluginData) DeserializeObjects() PluginDataEx {
	var pluginDataEx PluginDataEx
	var resolveInfos []ResolveInfo
	ds := data.Data.(map[string]interface{})
	//提取 data中的info信息
	for s2, i := range ds {
		if s2 == "info" {
			bts := i.([]interface{})
			//将info内的[]byte数组，反序列化为ResolveInfo
			for _, bt := range bts {
				var bti MapSInterface
				bti.UnmarshalMsg(bt.([]byte))
				//	msgpack.Unmarshal(bt.([]byte), &bti)
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
