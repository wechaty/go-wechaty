package payload

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// ImageType image type
type ImageType uint8

const (
	ImageTypeUnknown   ImageType = 0
	ImageTypeThumbnail           = 1
	ImageTypeHD                  = 2
	ImageTypeArtwork             = 3
)

type WechatMessageType uint8

const (
	WechatMessageTypeText              WechatMessageType = 1
	WechatMessageTypeImage                               = 3
	WechatMessageTypeVoice                               = 34
	WechatMessageTypeVerifyMsg                           = 37
	WechatMessageTypePossibleFriendMsg                   = 40
	WechatMessageTypeShareCard                           = 42
	WechatMessageTypeVideo                               = 43
	WechatMessageTypeEmoticon                            = 47
	WechatMessageTypeLocation                            = 48
	WechatMessageTypeApp                                 = 49
	WechatMessageTypeVoipMsg                             = 50
	WechatMessageTypeStatusNotify                        = 51
	WechatMessageTypeVoipNotify                          = 52
	WechatMessageTypeVoipInvite                          = 53
	WechatMessageTypeMicroVideo                          = 62
	WechatMessageTypeTransfer                            = 2000 // 转账
	WechatMessageTypeRedEnvelope                         = 2001 // 红包
	WechatMessageTypeMiniProgram                         = 2002 // 小程序
	WechatMessageTypeGroupInvite                         = 2003 // 群邀请
	WechatMessageTypeFile                                = 2004 // 文件消息
	WechatMessageTypeSysNotice                           = 9999
	WechatMessageTypeSys                                 = 10000
	WechatMessageTypeRecalled                            = 10002 // NOTIFY 服务通知
)

// MessageType message type
type MessageType int

const (
	MessageTypeText              MessageType = 1
	MessageTypeContact                       = 2
	MessageTypeImage                         = 3
	MessageTypeDeleted                       = 4
	MessageTypeVoice                         = 34
	MessageTypeSelfAvatar                    = 35
	MessageTypeVerifyMsg                     = 37
	MessageTypePossibleFriendMsg             = 40
	MessageTypeShareCard                     = 42
	MessageTypeVideo                         = 43
	MessageTypeEmoticon                      = 47
	MessageTypeLocation                      = 48
	MessageTypeApp                           = 49
	MessageTypeVoipMsg                       = 50
	MessageTypeStatusNotify                  = 51
	MessageTypeVoipNotify                    = 52
	MessageTypeVoipInvite                    = 53
	MessageTypeMicroVideo                    = 62
	MessageTypeSelfInfo                      = 101
	MessageTypeSysNotice                     = 9999
	MessageTypeSys                           = 10000
	MessageTypeRecalled                      = 10002
	MessageTypeN112048                       = 2048  // 2048 = 1 << 11
	MessageTypeN1532768                      = 32768 // 32768  = 1 << 15
)

// MessagePayload message payload
type MessagePayload struct {
	AppMsgType          int         `json:"AppMsgType"`
	Content             string      `json:"Content"`
	CreateTime          int         `json:"CreateTime"`
	FileName2           string      `json:"fileName"`
	FileName            string      `json:"FileName"`
	FromMemberNickName2 string      `json:"fromMemberNickName"`
	FromMemberNickName  string      `json:"FromMemberNickName"`
	FromMemberUserName  string      `json:"FromMemberUserName"`
	FromMemberUserName2 string      `json:"fromMemberUserName"`
	FromUserName        string      `json:"FromUserName"`
	ImgBuf              string      `json:"ImgBuf"`
	ImgStatus           int         `json:"ImgStatus"`
	L1MsgType           int         `json:"L1MsgType"`
	MsgId               string      `json:"MsgId"`
	MsgSource           string      `json:"MsgSource"`
	MsgSourceCd         int         `json:"msgSourceCd"`
	MsgType             MessageType `json:"MsgType"`
	NewMsgId            int         `json:"NewMsgId"`
	PushContent         string      `json:"PushContent"`
	Status              int         `json:"Status"`
	ToUserName          string      `json:"ToUserName"`
	Uin                 string      `json:"Uin"`
	Url                 string      `json:"Url"`
	Url2                string      `json:"url"`
	WechatUserName      string      `json:"wechatUserName"`
}

// MessageTypeUnknown     MessageType = 0
// MessageTypeAttachment  MessageType = 1
// MessageTypeAudio       MessageType = 2
// MessageTypeContact     MessageType = 3
// MessageTypeChatHistory MessageType = 4
// MessageTypeEmoticon    MessageType = 5
// MessageTypeImage       MessageType = 6
// MessageTypeText        MessageType = 7
// MessageTypeLocation    MessageType = 8
// MessageTypeMiniProgram MessageType = 9
// MessageTypeGroupNote   MessageType = 10
// MessageTypeTransfer    MessageType = 11
// MessageTypeRedEnvelope MessageType = 12
// MessageTypeRecalled    MessageType = 13
// MessageTypeURL         MessageType = 14
// MessageTypeVideo       MessageType = 15

func (t MessageType) ToSchemasMessageType() schemas.MessageType {
	switch t {
	case MessageTypeText:
		return schemas.MessageTypeText
	case MessageTypeImage:
		return schemas.MessageTypeImage
	case MessageTypeLocation:
		return schemas.MessageTypeLocation
	case MessageTypeVoice:
		return schemas.MessageTypeAudio
	case MessageTypeVideo:
		return schemas.MessageTypeVideo
	case MessageTypeContact:
		return schemas.MessageTypeContact
	case MessageTypeDeleted:
		return schemas.MessageTypeRecalled
	}
	return 0
}

// ToSchemasContactPayload to puppet contact payload
func (pay MessagePayload) ToSchemasContactPayload() (scp *schemas.MessagePayload) {
	payload := schemas.MessagePayload{
		MessagePayloadBase: schemas.MessagePayloadBase{
			Id:            pay.MsgId,
			MentionIdList: nil,
			FileName:      pay.FileName,
			Text:          pay.Content,
			Timestamp:     uint64(pay.CreateTime),
			Type:          pay.MsgType.ToSchemasMessageType(),
		},
		MessagePayloadRoom: schemas.MessagePayloadRoom{
			FromId: pay.FromUserName,
			RoomId: "",
			ToId:   pay.ToUserName,
		},
	}
	return &payload
}

// PadPlusMediaData media payload
type PadPlusMediaData struct {
	Content string `json:"content"`
	MsgId   string `json:"msg_id"`
	Src     string `json:"src"`
	Status  string `json:"status"`
	Thumb   string `json:"thumb"`
}

// SendMessageResponse send message client request response
type SendMessageResponse struct {
	MsgId     string `json:"msgId"`
	Timestamp int    `json:"timestamp"`
	Success   bool   `json:"success"`
}
