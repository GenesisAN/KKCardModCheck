package KK

//go:generate msgp
type BlockHeaderListInfo struct {
	LstInfo []*BlockHeader `msg:"lstInfo"`
}

type BlockHeader struct {
	Name    string `msg:"name"`
	Version string `msg:"version"`
	Pos     int    `msg:"pos"`
	Size    int    `msg:"size"`
	Data    []byte
}
