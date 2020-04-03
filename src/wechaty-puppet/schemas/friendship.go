package schemas

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
  Id        string
  ContactId string
  Hello     string
  Timestamp int64
}

type FriendshipPayloadConfirm struct {
  FriendshipPayloadBase
  Type FriendshipType // FriendshipTypeConfirm
}

type FriendshipPayloadReceive struct {
  FriendshipPayloadBase
  Type     FriendshipType // FriendshipTypeReceive
  Scene    FriendshipSceneType
  Stranger string
  Ticket   string
}

type FriendshipPayloadVerify struct {
  FriendshipPayloadBase
  Type FriendshipType // FriendshipTypeVerify
}

type FriendshipSearchCondition struct {
  Phone  string
  WeiXin string
}
