package payload

// PadPlusRichMediaData load media message
type PadPlusRichMediaData struct {
	Content      string `json:"content"`
	MsgType      int    `json:"msgType"`
	ContentType  string `json:"contentType"`
	Src          string `json:"src"`
	AppMsgType   int    `json:"appMsgType"`
	FileName     string `json:"fileName"`
	MsgId        string `json:"msgId"`
	CreateTime   int    `json:"createTime"`
	FromUserName string `json:"fromUserName"`
	ToUserName   string `json:"toUserName"`
}
