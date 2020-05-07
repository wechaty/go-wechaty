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
	Alias       string
	AliasRegexp *regexp.Regexp
	Id          string
	Name        string
	NameRegexp  *regexp.Regexp
	WeiXin      string
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
