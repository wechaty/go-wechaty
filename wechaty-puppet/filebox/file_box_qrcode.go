package filebox

import (
	"bytes"
	"github.com/skip2/go-qrcode"
	"io"
)

type fileBoxQRCode struct {
	qrCode string
}

func newFileBoxQRCode(qrCode string) *fileBoxQRCode {
	return &fileBoxQRCode{qrCode: qrCode}
}

func (fb *fileBoxQRCode) toJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"qrCode": fb.qrCode,
	}
}

func (fb *fileBoxQRCode) toBytes() ([]byte, error) {
	qr, err := qrcode.New(fb.qrCode, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	return qr.PNG(256)
}

func (fb *fileBoxQRCode) toReader() (io.Reader, error) {
	byteData, err := fb.toBytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byteData), nil
}
