package KK

import (
	"fmt"
)

func (kfp *KKChaFileParameter) GetSex() string {
	if kfp.Sex == 0 {
		return "男性"
	}
	return "女性"
}

func (kfp *KKChaFileParameter) Print() {
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

//go:generate msgp
type KKChaFileParameter struct {
	Version          string      `msg:"version"`
	Sex              int         `msg:"sex"`
	ExType           int         `msg:"exType"`
	Lastname         string      `msg:"lastname"`
	Firstname        string      `msg:"firstname"`
	Nickname         string      `msg:"nickname"`
	CallType         int         `msg:"callType"`
	Personality      int         `msg:"personality"`
	BloodType        int         `msg:"bloodType"`
	BirthMonth       int         `msg:"birthMonth"`
	BirthDay         int         `msg:"birthDay"`
	ClubActivities   int         `msg:"clubActivities"`
	VoiceRate        float64     `msg:"voiceRate"`
	WeakPoint        int         `msg:"weakPoint"`
	Awnser           Awnser      `msg:"awnser"`
	Denial           Denial      `msg:"denial"`
	Attribute        Attribute   `msg:"attribute"`
	Aggressive       int         `msg:"aggressive"`
	Diligence        int         `msg:"diligence"`
	Kindness         int         `msg:"kindness"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
type Awnser struct {
	Animal           bool        `msg:"animal"`
	Eat              bool        `msg:"eat"`
	Cook             bool        `msg:"cook"`
	Exercise         bool        `msg:"exercise"`
	Study            bool        `msg:"study"`
	Fashionable      bool        `msg:"fashionable"`
	BlackCoffee      bool        `msg:"blackCoffee"`
	Spicy            bool        `msg:"spicy"`
	Sweet            bool        `msg:"sweet"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
type Denial struct {
	Kiss             bool        `msg:"kiss"`
	Aibu             bool        `msg:"aibu"`
	Anal             bool        `msg:"anal"`
	Massage          bool        `msg:"massage"`
	NotCondom        bool        `msg:"notCondom"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
type Attribute struct {
	Hinnyo           bool        `msg:"hinnyo"`
	Harapeko         bool        `msg:"harapeko"`
	Donkan           bool        `msg:"donkan"`
	Choroi           bool        `msg:"choroi"`
	Bitch            bool        `msg:"bitch"`
	Mutturi          bool        `msg:"mutturi"`
	Dokusyo          bool        `msg:"dokusyo"`
	Ongaku           bool        `msg:"ongaku"`
	Kappatu          bool        `msg:"kappatu"`
	Ukemi            bool        `msg:"ukemi"`
	Friendly         bool        `msg:"friendly"`
	Kireizuki        bool        `msg:"kireizuki"`
	Taida            bool        `msg:"taida"`
	Sinsyutu         bool        `msg:"sinsyutu"`
	Hitori           bool        `msg:"hitori"`
	Undo             bool        `msg:"undo"`
	Majime           bool        `msg:"majime"`
	LikeGirls        bool        `msg:"likeGirls"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
