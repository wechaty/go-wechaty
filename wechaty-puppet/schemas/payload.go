package schemas

type PayloadType int32

const (
	PayloadTypeUnknown    = 0
	PayloadTypeMessage    = 1
	PayloadTypeContact    = 2
	PayloadTypeRoom       = 3
	PayloadTypeRoomMember = 4
	PayloadTypeFriendship = 5
)
