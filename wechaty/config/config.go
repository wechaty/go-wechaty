package config

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"regexp"
)

// AtSepratorRegex mobile: \u2005, PC、mac: \u0020
// Deprecated: use AtSeparatorRegexStr
const AtSepratorRegex = "[\u2005\u0020]"

// AtSeparatorRegexStr mobile: \u2005, PC、mac: \u0020
const AtSeparatorRegexStr = "[\u2005\u0020]"

const FourPerEmSpace = string(rune(8197))

// AtSeparatorRegex regular expression split '@'
var AtSeparatorRegex = regexp.MustCompile(AtSeparatorRegexStr)

func QRCodeForChatie() *filebox.FileBox {
	const chatieOfficialAccountQrcode = "http://weixin.qq.com/r/qymXj7DEO_1ErfTs93y5"
	return filebox.FromQRCode(chatieOfficialAccountQrcode)
}
