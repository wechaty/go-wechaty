package schemas

import (
	"regexp"
	"time"
)

//go:generate stringer -type=MessageType
type MessageType uint8

const (
	MessageTypeUnknown     MessageType = 0
	MessageTypeAttachment  MessageType = 1
	MessageTypeAudio       MessageType = 2
	MessageTypeContact     MessageType = 3
	MessageTypeChatHistory MessageType = 4
	MessageTypeEmoticon    MessageType = 5
	MessageTypeImage       MessageType = 6
	MessageTypeText        MessageType = 7
	MessageTypeLocation    MessageType = 8
	MessageTypeMiniProgram MessageType = 9
	MessageTypeGroupNote   MessageType = 10
	MessageTypeTransfer    MessageType = 11
	MessageTypeRedEnvelope MessageType = 12
	MessageTypeRecalled    MessageType = 13
	MessageTypeURL         MessageType = 14
	MessageTypeVideo       MessageType = 15
)

type WeChatAppMessageType int

const (
	WeChatAppMessageTypeText                  WeChatAppMessageType = 1
	WeChatAppMessageTypeImg                   WeChatAppMessageType = 2
	WeChatAppMessageTypeAudio                 WeChatAppMessageType = 3
	WeChatAppMessageTypeVideo                 WeChatAppMessageType = 4
	WeChatAppMessageTypeUrl                   WeChatAppMessageType = 5
	WeChatAppMessageTypeAttach                WeChatAppMessageType = 6
	WeChatAppMessageTypeOpen                  WeChatAppMessageType = 7
	WeChatAppMessageTypeEmoji                 WeChatAppMessageType = 8
	WeChatAppMessageTypeVoiceRemind           WeChatAppMessageType = 9
	WeChatAppMessageTypeScanGood              WeChatAppMessageType = 10
	WeChatAppMessageTypeGood                  WeChatAppMessageType = 13
	WeChatAppMessageTypeEmotion               WeChatAppMessageType = 15
	WeChatAppMessageTypeCardTicket            WeChatAppMessageType = 16
	WeChatAppMessageTypeRealtimeShareLocation WeChatAppMessageType = 17
	WeChatAppMessageTypeChatHistory           WeChatAppMessageType = 19
	WeChatAppMessageTypeMiniProgram           WeChatAppMessageType = 33
	WeChatAppMessageTypeTransfers             WeChatAppMessageType = 2000
	WeChatAppMessageTypeRedEnvelopes          WeChatAppMessageType = 2001
	WeChatAppMessageTypeReaderType            WeChatAppMessageType = 100001
)

type WeChatMessageType int

const (
	WeChatMessageTypeText              WeChatMessageType = 1
	WeChatMessageTypeImage             WeChatMessageType = 3
	WeChatMessageTypeVoice             WeChatMessageType = 34
	WeChatMessageTypeVerifyMsg         WeChatMessageType = 37
	WeChatMessageTypePossibleFriendMsg WeChatMessageType = 40
	WeChatMessageTypeShareCard         WeChatMessageType = 42
	WeChatMessageTypeVideo             WeChatMessageType = 43
	WeChatMessageTypeEmoticon          WeChatMessageType = 47
	WeChatMessageTypeLocation          WeChatMessageType = 48
	WeChatMessageTypeApp               WeChatMessageType = 49
	WeChatMessageTypeVOIPMsg           WeChatMessageType = 50
	WeChatMessageTypeStatusNotify      WeChatMessageType = 51
	WeChatMessageTypeVOIPNotify        WeChatMessageType = 52
	WeChatMessageTypeVOIPInvite        WeChatMessageType = 53
	WeChatMessageTypeMicroVideo        WeChatMessageType = 62
	WeChatMessageTypeTransfer          WeChatMessageType = 2000 // 转账
	WeChatMessageTypeRedEnvelope       WeChatMessageType = 2001 // 红包
	WeChatMessageTypeMiniProgram       WeChatMessageType = 2002 // 小程序
	WeChatMessageTypeGroupInvite       WeChatMessageType = 2003 // 群邀请
	WeChatMessageTypeFile              WeChatMessageType = 2004 // 文件消息
	WeChatMessageTypeSysNotice         WeChatMessageType = 9999
	WeChatMessageTypeSys               WeChatMessageType = 10000
	WeChatMessageTypeRecalled          WeChatMessageType = 10002
)

type MessagePayloadBase struct {
	Id string

	// use message id to get rawPayload to get those informations when needed
	// contactId     string        // Contact ShareCard
	MentionIdList []string // Mentioned Contacts' Ids

	FileName  string
	Text      string
	Timestamp time.Time
	Type      MessageType

	// 小程序有些消息类型，wechaty服务端解析不处理，框架端解析。 xml type 36 是小程序
	FixMiniApp bool
}

type MessagePayloadRoom struct {
	TalkerId   string
	RoomId     string
	ListenerId string
}

type MessagePayloadTo = MessagePayloadRoom

type MessagePayload struct {
	MessagePayloadBase
	MessagePayloadRoom
}

type MessageQueryFilter struct {
	TalkerId string
	// Deprecated: please use TalkerId
	FromId     string
	Id         string
	RoomId     string
	Text       string
	TextRegExp *regexp.Regexp
	// Deprecated: please use ListenerId
	ToId       string
	ListenerId string
	Type       MessageType
}

type MessagePayloadFilterFunction func(payload *MessagePayload) bool
