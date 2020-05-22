package config

import (
	file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
)

// AtSepratorRegex mobile: \u2005, PC、mac: \u0020
const AtSepratorRegex = "[\u2005\u0020]"

const FourPerEmSpace = string(rune(8197))

func QRCodeForChatie() *file_box.FileBox {
	const chatieOfficialAccountQrcode = "http://weixin.qq.com/r/qymXj7DEO_1ErfTs93y5"
	return file_box.FromQRCode(chatieOfficialAccountQrcode)
}
