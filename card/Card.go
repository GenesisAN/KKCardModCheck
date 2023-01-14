package card

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

func (cfpEx *Card) ChaFileParameterEx(card KoiCard) {
	cfpEx.CharInfo = &ChaFileParameterEx{}
	cfpEx.CharInfo.Lastname = card.CharParmeter.Lastname
	cfpEx.CharInfo.Firstname = card.CharParmeter.Firstname
	cfpEx.CharInfo.Version = card.CharParmeter.Version
	cfpEx.CharInfo.Nickname = card.CharParmeter.Nickname
}
