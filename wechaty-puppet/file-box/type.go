package file_box

import "net/http"

type FileBoxJsonObjectCommon struct {
  Name     string                 `json:"name"`
  Metadata map[string]interface{} `json:"metadata"`
  BoxType  FileBoxType            `json:"boxType"`
}

type FileBoxJsonObjectBase64 struct {
  *FileBoxJsonObjectCommon
  Base64 string `json:"base64"`
}

type FileBoxJsonObjectUrl struct {
  *FileBoxJsonObjectCommon
  RemoteUrl string      `json:"remoteUrl"`
  Headers   http.Header `json:"headers"`
}

type FileBoxJsonObjectQRCode struct {
  *FileBoxJsonObjectCommon
  QrCode string `json:"qrCode"`
}

type FileBoxType uint8

const (
  FileBoxTypeUnknown FileBoxType = 0

  FileBoxTypeBase64 = 1
  FileBoxTypeUrl    = 2
  FileBoxTypeQRCode = 3
  FileBoxTypeBuffer = 4
  FileBoxTypeFile   = 5
  FileBoxTypeStream = 6
)
