package schemas

import "regexp"

//go:generate stringer -type=ContactGender
type ContactGender uint8

const (
	ContactGenderUnknown ContactGender = 0
	ContactGenderMale    ContactGender = 1
	ContactGenderFemale  ContactGender = 2
)

//go:generate stringer -type=ContactType
type ContactType uint8

const (
	ContactTypeUnknown  ContactType = 0
	ContactTypePersonal ContactType = 1
	ContactTypeOfficial ContactType = 2
)

// ContactQueryFilter use the first non-empty parameter of all parameters to search
type ContactQueryFilter struct {
	// 别名过滤
	Alias string
	// 别名正则表达式过滤
	AliasRegexp *regexp.Regexp
	// id 过滤
	Id string
	// 昵称过滤
	Name string
	// 昵称正则表达式过滤
	NameRegexp *regexp.Regexp
	WeiXin     string
}

type ContactPayload struct {
	Id        string
	Gender    ContactGender
	Type      ContactType
	Name      string
	Avatar    string
	Address   string
	Alias     string
	City      string
	Friend    bool
	Province  string
	Signature string
	Star      bool
	WeiXin    string
}

type ContactPayloadFilterFunction func(payload *ContactPayload) bool
