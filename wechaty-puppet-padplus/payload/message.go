package payload

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

type MessagePayload struct {
	AppMsgType          int    `json:"AppMsgType"`
	Content             string `json:"Content"`
	CreateTime          int    `json:"CreateTime"`
	FileName2           string `json:"fileName"`
	FileName            string `json:"FileName"`
	FromMemberNickName2 string `json:"fromMemberNickName"`
	FromMemberNickName  string `json:"FromMemberNickName"`
	FromMemberUserName  string `json:"FromMemberUserName"`
	FromMemberUserName2 string `json:"fromMemberUserName"`
	FromUserName        string `json:"FromUserName"`
	ImgBuf              string `json:"ImgBuf"`
	ImgStatus           int    `json:"ImgStatus"`
	L1MsgType           int    `json:"L1MsgType"`
	MsgId               string `json:"MsgId"`
	MsgSource           string `json:"MsgSource"`
	MsgSourceCd         int    `json:"msgSourceCd"`
	MsgType             int    `json:"MsgType"`
	NewMsgId            int    `json:"NewMsgId"`
	PushContent         string `json:"PushContent"`
	Status              int    `json:"Status"`
	ToUserName          string `json:"ToUserName"`
	Uin                 string `json:"Uin"`
	Url                 string `json:"Url"`
	Url2                string `json:"url"`
	WechatUserName      string `json:"wechatUserName"`
}

// 媒体信息响应结构体
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
