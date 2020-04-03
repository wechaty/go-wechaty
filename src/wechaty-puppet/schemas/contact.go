package schemas

type ContactGender uint8

const (
  ContactGenderUnknown ContactGender = 0
  ContactGenderMale    ContactGender = 1
  ContactGenderFemale  ContactGender = 2
)

type ContactType uint8

const (
  ContactTypeUnknown  ContactType = 0
  ContactTypePersonal ContactType = 1
  ContactTypeOfficial ContactType = 2
)

type ContactQueryFilter struct {
  Alias  string
  Id     string
  Name   string
  WeiXin string
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
  Start     bool
  WeiXin    string
}
