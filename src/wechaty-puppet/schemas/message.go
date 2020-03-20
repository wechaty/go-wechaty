package schemas

type MessageType uint8

const (
	MessageTypeUnknown MessageType = 0

	MessageTypeAttachment  // Attach(6)
	MessageTypeAudio       // Audio(1) Voice(34)
	MessageTypeContact     // ShareCard(42)
	MessageTypeChatHistory // ChatHistory(19)
	MessageTypeEmoticon    // Sticker Emoticon(15) Emoticon(47)
	MessageTypeImage       // Img(2) Image(3)
	MessageTypeText        // Text(1)
	MessageTypeLocation    // Location(48)
	MessageTypeMiniProgram // MiniProgram(33)
	MessageTypeTransfer    // Transfers(2000)
	MessageTypeRedEnvelope // RedEnvelopes(2001)
	MessageTypeRecalled    // Recalled(10002)
	MessageTypeUrl         // Url(5)
	MessageTypeVideo       // Video(4) Video(43)
)

type WechatAppMessageType uint8

const (
	WechatAppMessageTypeText                  WechatAppMessageType = 1
	WechatAppMessageTypeImg                                        = 2
	WechatAppMessageTypeAudio                                      = 3
	WechatAppMessageTypeVideo                                      = 4
	WechatAppMessageTypeUrl                                        = 5
	WechatAppMessageTypeAttach                                     = 6
	WechatAppMessageTypeOpen                                       = 7
	WechatAppMessageTypeEmoji                                      = 8
	WechatAppMessageTypeVoiceRemind                                = 9
	WechatAppMessageTypeScanGood                                   = 10
	WechatAppMessageTypeGood                                       = 13
	WechatAppMessageTypeEmotion                                    = 15
	WechatAppMessageTypeCardTicket                                 = 16
	WechatAppMessageTypeRealtimeShareLocation                      = 17
	WechatAppMessageTypeChatHistory                                = 19
	WechatAppMessageTypeMiniProgram                                = 33
	WechatAppMessageTypeTransfers                                  = 2000
	WechatAppMessageTypeRedEnvelopes                               = 2001
	WechatAppMessageTypeReaderType                                 = 100001
)

type WechatMessageType uint8

const (
	WechatMessageTypeText              = 1
	WechatMessageTypeImage             = 3
	WechatMessageTypeVoice             = 34
	WechatMessageTypeVerifyMsg         = 37
	WechatMessageTypePossibleFriendMsg = 40
	WechatMessageTypeShareCard         = 42
	WechatMessageTypeVideo             = 43
	WechatMessageTypeEmoticon          = 47
	WechatMessageTypeLocation          = 48
	WechatMessageTypeApp               = 49
	WechatMessageTypeVoipMsg           = 50
	WechatMessageTypeStatusNotify      = 51
	WechatMessageTypeVoipNotify        = 52
	WechatMessageTypeVoipInvite        = 53
	WechatMessageTypeMicroVideo        = 62
	WechatMessageTypeTransfer          = 2000 // 转账
	WechatMessageTypeRedEnvelope       = 2001 // 红包
	WechatMessageTypeMiniProgram       = 2002 // 小程序
	WechatMessageTypeGroupInvite       = 2003 // 群邀请
	WechatMessageTypeFile              = 2004 // 文件消息
	WechatMessageTypeSysNotice         = 9999
	WechatMessageTypeSys               = 10000
	WechatMessageTypeRecalled          = 10002 // NOTIFY 服务通知
)

type MessagePayloadBase struct {
	Id string

	// use message id to get rawPayload to get those informations when needed
	// contactId     string        // Contact ShareCard
	MentionIdList []string // Mentioned Contacts' Ids

	FileName  string
	Text      string
	Timestamp int
	Type      MessageType
}

type MessagePayloadRoom struct {
	MessagePayloadBase

	FromId string
	RoomId string
	ToId   string // if to is not set then room must be set
}

type MessagePayloadTo struct {
	MessagePayloadBase

	FromId string
	RoomId string
	ToId   string // if to is not set then room must be set
}

// todo:: MessagePayload
type MessagePayload interface{}

type MessageUserQueryFilter struct {
	FromId string
	Id     string
	RoomId string
	Text   string // todo:: RegExp
	ToId   string
	Type   MessageType
}

type MessagePayloadFilterFunction func(payload MessagePayload) bool
