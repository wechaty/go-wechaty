package config

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
)

// AtSepratorRegex mobile: \u2005, PC、mac: \u0020
const AtSepratorRegex = "[\u2005\u0020]"

const FourPerEmSpace = string(rune(8197))

func QRCodeForChatie() *filebox.FileBox {
	const chatieOfficialAccountQrcode = "http://weixin.qq.com/r/qymXj7DEO_1ErfTs93y5"
	return filebox.FromQRCode(chatieOfficialAccountQrcode)
}
