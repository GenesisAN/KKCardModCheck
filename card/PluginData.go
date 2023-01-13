// 角色卡插件数据结构
package card

import "fmt"

// PluginData 插件数据
type PluginData struct {
	RequiredPluginGUIDs []string    //插件GUID
	RequiredZipmodGUIDs []string    //ZIPMOD GUID
	Version             int         //版本
	Data                interface{} //原始MsgPack数据
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

type ResolveInfo struct {
	GUID       string `mapstructure:"ModID"` //GUID of the mod as defined in the manifest.xml
	Slot       int    //在mod的列表文件中定义的项目ID
	LocalSlot  int    //Resolved item ID.  IDs greater than 100000000 are resolved IDs belonging to Sideloader.  Use the resolved ID (local slot) to look up the original ID (slot)
	Property   string // Property of the object as defined in Sideloader's StructReference. If ever you need to know what to use for this, enable debug resolve info logging and see what Sideloader generates at the start of the game.
	CategoryNo int    //ChaListDefine.CategoryNo. Typically only used for hard mod resolving in cases where the GUID is not known.
}
