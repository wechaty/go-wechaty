package schemas

//go:generate stringer -type=PayloadType

// PayloadType ...
type PayloadType int32

const (
	// PayloadTypeUnknown unknown
	PayloadTypeUnknown    PayloadType = 0
	PayloadTypeMessage                = 1
	PayloadTypeContact                = 2
	PayloadTypeRoom                   = 3
	PayloadTypeRoomMember             = 4
	PayloadTypeFriendship             = 5
)
