package card

import "fmt"

type KoiChaFileParameter struct {
	CurrentVersion string
	Version        string
	Sex            int8
	Lastname       string
	Firstname      string
	Nickname       string
	CallType       int32
	Personality    int32
	BloodType      int8
	BirthMonth     int8
	BirthDay       int8
	ClubActivities int8
	VoiceRate      float32
	WeakPoint      int32
	Awnser         Awnser
	Denial         Denial
	Attribute      Attribute
	Aggressive     int32
	Diligence      int32
	Kindness       int32
}

func (kfp *KoiChaFileParameter) GetSex() string {
	if kfp.Sex == 0 {
		return "男性"
	}
	return "女性"
}

func (kfp *KoiChaFileParameter) Print() {
	cc := *kfp
	if cc.Nickname == "" {
		cc.Nickname = "无"
	}
	if cc.Firstname == "" {
		cc.Firstname = "无"
	}
	if cc.Lastname == "" {
		cc.Lastname = "无"
	}
	fmt.Printf(`昵称:%s
姓:%s
名:%s
性别:%s
				`, cc.Nickname, cc.Firstname, cc.Lastname, kfp.GetSex())
}

type Awnser struct {
	Animal      bool //喜欢动物
	Eat         bool
	Cook        bool
	Exercise    bool
	Study       bool
	Fashionable bool
	BlackCoffee bool
	Spicy       bool
	Sweet       bool //喜欢甜食
}
type Denial struct {
	Kiss      bool
	Aibu      bool
	Anal      bool
	Massage   bool
	NotCondom bool
}
type Attribute struct {
	Hinnyo    bool
	Harapeko  bool
	Donkan    bool
	Choroi    bool
	Bitch     bool
	Mutturi   bool
	Dokusyo   bool
	Ongaku    bool
	Kappatu   bool
	Ukemi     bool
	Friendly  bool
	Kireizuki bool
	Taida     bool
	Sinsyutu  bool
	Hitori    bool
	Undo      bool
	Majime    bool
	LikeGirls bool
}
