package config

import file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"

func QRCodeForChatie() *file_box.FileBox {
	const chatieOfficialAccountQrcode = "http://weixin.qq.com/r/qymXj7DEO_1ErfTs93y5"
	return file_box.FromQRCode(chatieOfficialAccountQrcode)
}
