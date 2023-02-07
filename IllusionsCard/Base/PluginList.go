package Base

//go:generate msgp
type MapSArrayInterface map[string][]interface{}

type MapSInterface map[string]interface{}

type ResolveInfo struct {
	GUID       string `mapstructure:"ModID"` //GUID of the mod as defined in the manifest.xml
	Slot       int    //在mod的列表文件中定义的项目ID
	LocalSlot  int    //Resolved item ID.  IDs greater than 100000000 are resolved IDs belonging to Sideloader.  Use the resolved ID (local slot) to look up the original ID (slot)
	Property   string // Property of the object as defined in Sideloader's StructReference. If ever you need to know what to use for this, enable debug resolve info logging and see what Sideloader generates at the start of the game.
	CategoryNo int    //ChaListDefine.CategoryNo. Typically only used for hard mod resolving in cases where the GUID is not known.
}
