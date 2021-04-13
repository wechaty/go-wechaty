package schemas

//go:generate stringer -type=ScanStatus
type ScanStatus uint8

const (
	ScanStatusUnknown   ScanStatus = 0
	ScanStatusCancel    ScanStatus = 1
	ScanStatusWaiting   ScanStatus = 2
	ScanStatusScanned   ScanStatus = 3
	ScanStatusConfirmed ScanStatus = 4
	ScanStatusTimeout   ScanStatus = 5
)

type EventFriendshipPayload struct {
	FriendshipID string
}

type EventLoginPayload struct {
	ContactId string
}

type EventLogoutPayload struct {
	ContactId string
	Data      string
}

type EventMessagePayload struct {
	MessageId string
}

type EventRoomInvitePayload struct {
	RoomInvitationId string
}

type EventRoomJoinPayload struct {
	InviteeIdList []string
	InviterId     string
	RoomId        string
	Timestamp     int64
}

type EventRoomLeavePayload struct {
	RemoveeIdList []string
	RemoverId     string
	RoomId        string
	Timestamp     int64
}

type EventRoomTopicPayload struct {
	ChangerId string
	NewTopic  string
	OldTopic  string
	RoomId    string
	Timestamp int64
}

type EventScanPayload struct {
	BaseEventPayload
	Status ScanStatus
	QrCode string
}

type EventDirtyPayload struct {
	PayloadType PayloadType
	PayloadId   string
}

type BaseEventPayload struct {
	Data string
}

type EventDongPayload = BaseEventPayload

type EventErrorPayload = BaseEventPayload

type EventReadyPayload = BaseEventPayload

type EventResetPayload = BaseEventPayload

type EventHeartbeatPayload = BaseEventPayload
