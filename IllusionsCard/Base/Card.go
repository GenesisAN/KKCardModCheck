package Base

type Card struct {
	Extended     map[string]*PluginData
	ExtendedList map[string]*PluginDataEx
	CharInfo     *ChaFileParameterEx
	Image        []byte
	CardType     string
	LoadVersion  string
	Path         string
}

type ChaFileParameterEx struct {
	Version   string
	Lastname  string
	Firstname string
	Nickname  string
}
