package file_box

import (
  "github.com/skip2/go-qrcode"
)

type fileBoxQRCode struct {
  qrCode string
}

func NewFileBoxQRCode(qrCode string) *fileBoxQRCode {
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
