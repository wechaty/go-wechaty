package payload

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// ContactGender contact gender
type ContactGender int8

const (
	ContactGenderUnknown ContactGender = 0
	ContactGenderMale                  = 1
	ContactGenderFemale                = 2
)

// ToSchemasContactGender to puppet contact gender
func (g ContactGender) ToSchemasContactGender() schemas.ContactGender {
	switch g {
	case ContactGenderUnknown:
		return schemas.ContactGenderUnknown
	case ContactGenderMale:
		return schemas.ContactGenderMale
	case ContactGenderFemale:
		return schemas.ContactGenderFemale
	}
	return schemas.ContactGenderUnknown
}

// ContactPayload ipadplus contact payload
type ContactPayload struct {
	Alias        string        `json:"alias"`
	ContactType  int           `json:"contactType"`
	TagList      string        `json:"tagList"`
	BigHeadUrl   string        `json:"big_head_url"` // "http://wx.qlogo.cn/mmhead/ver_1/xfCMmibHH74xGLoyeDFJadrZXX3eOEznPefiaCa3iczxZGMwPtDuSbRQKx3Xdm18un303mf0NFia3USY2nO2VEYILw/0",
	City         string        `json:"city"`         // 'Haidian'
	Country      string        `json:"country"`      // "CN"
	NickName     string        `json:"nickName"`     // "梦君君", RPCContact: 用户昵称， Room: 群昵称
	Province     string        `json:"province"`     // "Beijing",
	Remark       string        `json:"remark"`       // "女儿",
	Sex          ContactGender `json:"sex"`
	Signature    string        `json:"signature"` // "且行且珍惜",
	SmallHeadUrl string        `json:"smallHeadUrl"`
	Stranger     string        `json:"stranger"` // 用户v1码，从未加过好友则为空 "v1_0468f2cd3f0efe7ca2589d57c3f9ba952a3789e41b6e78ee00ed53d1e6096b88@stranger"
	Ticket       string        `json:"ticket"`   // 用户v2码，如果非空则为单向好友(非对方好友) 'v2_xxx@stranger'
	UserName     string        `json:"userName"` // "mengjunjun001" | "qq512436430" Unique name
	VerifyFlag   int           `json:"verifyFlag"`
	ContactFlag  int           `json:"contactFlag"`
}

// ToSchemasContactPayload to puppet contact
func (pay ContactPayload) ToSchemasContactPayload() (scp *schemas.ContactPayload) {
	payload := schemas.ContactPayload{
		Id:        pay.UserName,
		Gender:    pay.Sex.ToSchemasContactGender(),
		Type:      0,
		Name:      pay.NickName,
		Avatar:    pay.BigHeadUrl,
		Address:   "",
		Alias:     pay.Alias,
		City:      pay.City,
		Friend:    false,
		Province:  pay.Province,
		Signature: pay.Signature,
		Star:      false,
		WeiXin:    pay.Stranger,
	}
	return &payload
}

// RPCContact contact
type RPCContactPayload struct {
	Alias           string
	BigHeadImgUrl   string
	ChatRoomOwner   string
	ChatroomVersion int
	City            string
	ContactFlag     int
	ContactType     int
	EncryptUsername string
	ExtInfo         string
	ExtInfoExt      string
	HeadImgUrl      string
	LabelLists      string
	MsgType         int
	NickName        string
	Province        string
	PYInitial       string
	PYQuanPin       string
	Remark          string
	RemarkName      string
	RemarkPYInitial string
	RemarkPYQuanPin string
	Seq             string
	Sex             int
	Signature       string
	SmallHeadImgUrl string
	Type7           string
	Uin             int
	UserName        string
	VerifyFlag      int
	WechatUserName  string `json:"wechatUserName"`
}

// ToSchemasContactGender to puppet contact gender
func (pay RPCContactPayload) ToSchemasContactGender() schemas.ContactGender {
	switch pay.Sex {
	case 0:
		return schemas.ContactGenderUnknown
	case 1:
		return schemas.ContactGenderMale
	case 2:
		return schemas.ContactGenderFemale
	}
	return schemas.ContactGenderUnknown
}

func (pay RPCContactPayload) ToSchemasContactPayload() (scp *schemas.ContactPayload) {
	payload := schemas.ContactPayload{
		Id:        pay.UserName,
		Gender:    pay.ToSchemasContactGender(),
		Type:      0,
		Name:      pay.NickName,
		Avatar:    pay.BigHeadImgUrl,
		Address:   "",
		Alias:     pay.Alias,
		City:      pay.City,
		Friend:    false,
		Province:  pay.Province,
		Signature: pay.Signature,
		Star:      false,
		WeiXin:    pay.WechatUserName,
	}
	return &payload
}

// ToContactPayload RPCContactPayload to ContactPayload
func (pay RPCContactPayload) ToContactPayload() (scp *ContactPayload) {
	payload := ContactPayload{
		Alias:        pay.Alias,
		ContactType:  pay.ContactType,
		TagList:      pay.Type7,
		BigHeadUrl:   pay.BigHeadImgUrl,
		City:         pay.City,
		Country:      "",
		NickName:     pay.NickName,
		Province:     pay.Province,
		Remark:       pay.Remark,
		Sex:          ContactGender(pay.Sex),
		Signature:    pay.Signature,
		SmallHeadUrl: pay.SmallHeadImgUrl,
		Stranger:     "",
		Ticket:       "",
		UserName:     pay.UserName,
		VerifyFlag:   pay.VerifyFlag,
		ContactFlag:  pay.ContactFlag,
	}
	return &payload
}

// ContactRoom contact room
type ContactRoom struct {
	ContactType     int
	ExtInfoExt      string
	Sex             int
	EncryptUsername string
	wechatUserName  string
	PYQuanPin       string
	Remark          string
	LabelLists      string
	ChatroomVersion int
	ExtInfo         string
	ChatRoomOwner   string
	VerifyFlag      int
	ContactFlag     int
	Ticket          string
	UserName        string
	src             int
	HeadImgUrl      string
	RemarkPYInitial string
	MsgType         int
	City            string
	NickName        string
	Province        string
	Alias           string
	Signature       string
	RemarkName      string
	RemarkPYQuanPin string
	Uin             int
	SmallHeadImgUrl string
	PYInitial       string
	Seq             string
	BigHeadImgUrl   string
}
