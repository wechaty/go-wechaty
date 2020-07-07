package payload

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// EventPadPlusQrCode ResponseType_LOGIN_QRCODE
type EventPadPlusQrCode struct {
	QrCode   string `json:"qrcode"`
	QrCodeId string `json:"qrcodeId"`
}

// QrCodeStatus scan qr code status
type QrCodeStatus uint8

const (
	ScanStatusWaiting   QrCodeStatus = 0
	ScanStatusScanned   QrCodeStatus = 1
	ScanStatusConfirmed QrCodeStatus = 2
	ScanStatusExpired   QrCodeStatus = 3
	ScanStatusCanceled  QrCodeStatus = 4
)

// ToPuppetStatus padplus scan status to puppet status
func (q QrCodeStatus) ToPuppetStatus() schemas.ScanStatus {
	switch q {
	case ScanStatusWaiting:
		return schemas.ScanStatusWaiting
	case ScanStatusScanned:
		return schemas.ScanStatusScanned
	case ScanStatusConfirmed:
		return schemas.ScanStatusConfirmed
	case ScanStatusExpired:
		return schemas.ScanStatusTimeout
	case ScanStatusCanceled:
		return schemas.ScanStatusCancel
	}
	return schemas.ScanStatusUnknown
}

// PadPlusQrCodeStatus 登录二维码扫描状态 ResponseType_QRCODE_SCAN
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

// EventScanData 扫描二维码事件
type EventScanData struct {
	HeadUrl  string       `json:"head_url"`
	Msg      string       `json:"msg"`
	NickName string       `json:"nick_name"`
	QrCodeId string       `json:"qrcode_id"`
	Status   QrCodeStatus `json:"status"` // 1 不同步 2 同步消息
	UserName string       `json:"user_name"`
}

type LogoutGRPCResponse struct {
	Code    int    `json:"code"`
	Uin     string `json:"uin"`
	Message string `json:"message"`
	MQType  int    `json:"mq_type"`
}

// QrCodeLogin 登录成功的事件
type QrCodeLogin struct {
	Alias      string `json:"alias"`
	HeadImgUrl string `json:"headImgUrl"`
	NickName   string `json:"nickName"`
	Status     int    `json:"status"`
	Uin        string `json:"uin"`
	UserName   string `json:"userName"`
	VerifyFlag string `json:"verifyFlag"`
}
