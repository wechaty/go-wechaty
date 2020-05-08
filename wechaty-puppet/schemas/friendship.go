package schemas

//go:generate stringer -type=FriendshipType
type FriendshipType uint8

const (
	FriendshipTypeUnknown FriendshipType = 0
	FriendshipTypeConfirm FriendshipType = 1
	FriendshipTypeReceive FriendshipType = 2
	FriendshipTypeVerify  FriendshipType = 3
)

type FriendshipSceneType uint8

const (
	FriendshipSceneTypeUnknown  FriendshipSceneType = 0
	FriendshipSceneTypeQQ       FriendshipSceneType = 1
	FriendshipSceneTypeEmail    FriendshipSceneType = 2
	FriendshipSceneTypeWeiXin   FriendshipSceneType = 3
	FriendshipSceneTypeQQTBD    FriendshipSceneType = 12 // QQ号搜索
	FriendshipSceneTypeRoom     FriendshipSceneType = 14
	FriendshipSceneTypePhone    FriendshipSceneType = 15
	FriendshipSceneTypeCard     FriendshipSceneType = 17 // 名片分享
	FriendshipSceneTypeLocation FriendshipSceneType = 18
	FriendshipSceneTypeBottle   FriendshipSceneType = 25
	FriendshipSceneTypeShaking  FriendshipSceneType = 29
	FriendshipSceneTypeQRCode   FriendshipSceneType = 30
)

type FriendshipPayloadBase struct {
	Id        string `json:"id"`
	ContactId string `json:"contactId"`
	Hello     string `json:"hello"`
	Timestamp int64  `json:"timestamp"`
}

type FriendshipPayloadConfirm struct {
	FriendshipPayloadBase
	Type FriendshipType // FriendshipTypeConfirm
}

type FriendshipPayloadReceive struct {
	FriendshipPayloadBase
	Type     FriendshipType      `json:"type"` // FriendshipTypeReceive
	Scene    FriendshipSceneType `json:"scene"`
	Stranger string              `json:"stranger"`
	Ticket   string              `json:"ticket"`
}

type FriendshipPayloadVerify struct {
	FriendshipPayloadBase
	Type FriendshipType // FriendshipTypeVerify
}

type FriendshipPayload struct {
	FriendshipPayloadReceive
}

// FriendshipSearchCondition use the first non-empty parameter of all parameters to search
type FriendshipSearchCondition struct {
	Phone  string
	WeiXin string
}
