package file_box

import (
	"github.com/skip2/go-qrcode"
)

type fileBoxQRCode struct {
	qrCode string
}

func newFileBoxQRCode(qrCode string) *fileBoxQRCode {
	return &fileBoxQRCode{qrCode: qrCode}
}

func (fb *fileBoxQRCode) toJSONMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"qrCode": fb.qrCode,
	}, nil
}

func (fb *fileBoxQRCode) toBytes() ([]byte, error) {
	qr, err := qrcode.New(fb.qrCode, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	return qr.PNG(256)
}
