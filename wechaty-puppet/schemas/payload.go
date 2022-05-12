package schemas

//go:generate stringer -type=PayloadType

// PayloadType ...
type PayloadType int32

const (
	// PayloadTypeUnknown unknown
	PayloadTypeUnknown PayloadType = iota
	PayloadTypeMessage
	PayloadTypeContact
	PayloadTypeRoom
	PayloadTypeRoomMember
	PayloadTypeFriendship
)
