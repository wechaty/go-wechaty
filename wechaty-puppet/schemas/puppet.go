package schemas

import pbwechaty "github.com/wechaty/go-grpc/wechaty/puppet"

//go:generate stringer -type=PuppetEventName
type PuppetEventName uint8

const (
	PuppetEventNameUnknown PuppetEventName = iota
	PuppetEventNameFriendship
	PuppetEventNameLogin
	PuppetEventNameLogout
	PuppetEventNameMessage
	PuppetEventNameRoomInvite
	PuppetEventNameRoomJoin
	PuppetEventNameRoomLeave
	PuppetEventNameRoomTopic
	PuppetEventNameScan

	PuppetEventNameDong
	PuppetEventNameError
	PuppetEventNameHeartbeat
	PuppetEventNameReady
	PuppetEventNameReset
	PuppetEventNameDirty

	PuppetEventNameStop
	PuppetEventNameStart
)

var eventNames = []PuppetEventName{
	PuppetEventNameFriendship,
	PuppetEventNameLogin,
	PuppetEventNameLogout,
	PuppetEventNameMessage,
	PuppetEventNameRoomInvite,
	PuppetEventNameRoomJoin,
	PuppetEventNameRoomLeave,
	PuppetEventNameRoomTopic,
	PuppetEventNameScan,

	PuppetEventNameDong,
	PuppetEventNameError,
	PuppetEventNameHeartbeat,
	PuppetEventNameReady,
	PuppetEventNameReset,
	PuppetEventNameDirty,

	PuppetEventNameStop,
	PuppetEventNameStart,
}

func GetEventNames() []PuppetEventName {
	return eventNames
}

var pbEventType2PuppetEventName = map[pbwechaty.EventType]PuppetEventName{
	pbwechaty.EventType_EVENT_TYPE_DONG:        PuppetEventNameDong,
	pbwechaty.EventType_EVENT_TYPE_ERROR:       PuppetEventNameError,
	pbwechaty.EventType_EVENT_TYPE_HEARTBEAT:   PuppetEventNameHeartbeat,
	pbwechaty.EventType_EVENT_TYPE_FRIENDSHIP:  PuppetEventNameFriendship,
	pbwechaty.EventType_EVENT_TYPE_LOGIN:       PuppetEventNameLogin,
	pbwechaty.EventType_EVENT_TYPE_LOGOUT:      PuppetEventNameLogout,
	pbwechaty.EventType_EVENT_TYPE_MESSAGE:     PuppetEventNameMessage,
	pbwechaty.EventType_EVENT_TYPE_READY:       PuppetEventNameReady,
	pbwechaty.EventType_EVENT_TYPE_ROOM_INVITE: PuppetEventNameRoomInvite,
	pbwechaty.EventType_EVENT_TYPE_ROOM_JOIN:   PuppetEventNameRoomJoin,
	pbwechaty.EventType_EVENT_TYPE_ROOM_LEAVE:  PuppetEventNameRoomLeave,
	pbwechaty.EventType_EVENT_TYPE_ROOM_TOPIC:  PuppetEventNameRoomTopic,
	pbwechaty.EventType_EVENT_TYPE_SCAN:        PuppetEventNameScan,
	pbwechaty.EventType_EVENT_TYPE_RESET:       PuppetEventNameReset,
	pbwechaty.EventType_EVENT_TYPE_UNSPECIFIED: PuppetEventNameUnknown,
	pbwechaty.EventType_EVENT_TYPE_DIRTY:       PuppetEventNameDirty,
}

// PbEventType2PuppetEventName grpc event map wechaty-puppet event name
func PbEventType2PuppetEventName() map[pbwechaty.EventType]PuppetEventName {
	return pbEventType2PuppetEventName
}

var pbEventType2GeneratePayloadFunc = map[pbwechaty.EventType]func() interface{}{
	pbwechaty.EventType_EVENT_TYPE_DONG:        func() interface{} { return &EventDongPayload{} },
	pbwechaty.EventType_EVENT_TYPE_ERROR:       func() interface{} { return &EventErrorPayload{} },
	pbwechaty.EventType_EVENT_TYPE_HEARTBEAT:   func() interface{} { return &EventHeartbeatPayload{} },
	pbwechaty.EventType_EVENT_TYPE_FRIENDSHIP:  func() interface{} { return &EventFriendshipPayload{} },
	pbwechaty.EventType_EVENT_TYPE_LOGIN:       func() interface{} { return &EventLoginPayload{} },
	pbwechaty.EventType_EVENT_TYPE_LOGOUT:      func() interface{} { return &EventLogoutPayload{} },
	pbwechaty.EventType_EVENT_TYPE_MESSAGE:     func() interface{} { return &EventMessagePayload{} },
	pbwechaty.EventType_EVENT_TYPE_READY:       func() interface{} { return &EventReadyPayload{} },
	pbwechaty.EventType_EVENT_TYPE_ROOM_INVITE: func() interface{} { return &EventRoomInvitePayload{} },
	pbwechaty.EventType_EVENT_TYPE_ROOM_JOIN:   func() interface{} { return &EventRoomJoinPayload{} },
	pbwechaty.EventType_EVENT_TYPE_ROOM_LEAVE:  func() interface{} { return &EventRoomLeavePayload{} },
	pbwechaty.EventType_EVENT_TYPE_ROOM_TOPIC:  func() interface{} { return &EventRoomTopicPayload{} },
	pbwechaty.EventType_EVENT_TYPE_SCAN:        func() interface{} { return &EventScanPayload{} },
	pbwechaty.EventType_EVENT_TYPE_RESET:       func() interface{} { return &EventResetPayload{} },
	pbwechaty.EventType_EVENT_TYPE_DIRTY:       func() interface{} { return &EventDirtyPayload{} },
}

// PbEventType2GeneratePayloadFunc grpc event map wechaty-puppet event payload
func PbEventType2GeneratePayloadFunc() map[pbwechaty.EventType]func() interface{} {
	return pbEventType2GeneratePayloadFunc
}
