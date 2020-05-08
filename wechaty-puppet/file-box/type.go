package file_box

import "net/http"

type FileBoxOptionsCommon struct {
	Name     string                 `json:"Name"`
	Metadata map[string]interface{} `json:"metadata"`
	BoxType  FileBoxType            `json:"boxType"`
}

type FileBoxOptionsBase64 struct {
	Base64 string `json:"base64"`
}

type FileBoxOptionsUrl struct {
	RemoteUrl string      `json:"remoteUrl"`
	Headers   http.Header `json:"headers"`
}

type FileBoxOptionsQRCode struct {
	QrCode string `json:"qrCode"`
}

type FileBoxOptions struct {
	FileBoxOptionsCommon
	FileBoxOptionsBase64
	FileBoxOptionsQRCode
	FileBoxOptionsUrl
}

//go:generate stringer -type=FileBoxType
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
