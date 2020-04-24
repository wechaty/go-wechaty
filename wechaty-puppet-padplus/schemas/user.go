package schemas

// 获取登录二维码 ResponseType_LOGIN_QRCODE
type PadPlusQrCode struct {
	QrCode   string `json:"qrcode"`
	QrCodeId string `json:"qrcodeId"`
}

type QrCodeStatus uint32

// 登录二维码状态
const (
	Waiting   QrCodeStatus = 0
	Scanned   QrCodeStatus = 1
	Confirmed QrCodeStatus = 2
	Canceled  QrCodeStatus = 4
	Expired   QrCodeStatus = 3
)

// 登录二维码扫描状态 ResponseType_QRCODE_SCAN
type PadPlusQrCodeStatus struct {
	HeadUrl  string       `json:"headUrl"`
	NickName string       `json:"nickName"`
	Status   QrCodeStatus `json:"status"`
	UserName string       `json:"userName"`
}

type LoginStatus uint32

const (
	Logined LoginStatus = 1
)

type PadPlusQrCodeLogin struct {
	HeadImgUrl string      `json:"headImgUrl"`
	NickName   string      `json:"nickName"`
	Status     LoginStatus `json:"status"`
	UserName   string      `json:"userName"`
	Uin        string      `json:"uin"`
	VerifyFlag string      `json:"verifyFlag"`
}

type GRPCLoginData struct {
	Event    string `json:"event"`
	HeadUrl  string `json:"head_url"`
	Loginer  string `json:"loginer"`
	Msg      string `json:"msg"`
	NickName string `json:"nick_name"`
	QrCodeId string `json:"qrcode_id"`
	ServerId string `json:"server_id"`
	Status   int    `json:"status"`
	UserName string `json:"user_name"`
}

// 扫描二维码事件
type ScanData struct {
	HeadUrl  string `json:"head_url"`
	Msg      string `json:"msg"`
	NickName string `json:"nick_name"`
	QrCodeId string `json:"qrcode_id"`
	Status   int    `json:"status"` // 1 不同步 2 同步消息
	UserName string `json:"user_name"`
}

type LogoutGRPCResponse struct {
	Code    int    `json:"code"`
	Uin     string `json:"uin"`
	Message string `json:"message"`
	MQType  int    `json:"mq_type"`
}

type GRPCQrCodeLogin struct {
	Alias      string `json:"alias"`
	HeadImgUrl string `json:"headImgUrl"`
	NickName   string `json:"nickName"`
	Status     int    `json:"status"`
	Uin        string `json:"uin"`
	UserName   string `json:"userName"`
	VerifyFlag string `json:"verifyFlag"`
}
