package filebox

import "net/http"


type OptionsCommon struct {
	Name     string                 `json:"Name"`
	Metadata map[string]interface{} `json:"metadata"`
	BoxType  Type                   `json:"boxType"`
}

type OptionsBase64 struct {
	Base64 string `json:"base64"`
}

type OptionsUrl struct {
	RemoteUrl string      `json:"remoteUrl"`
	Headers   http.Header `json:"headers"`
}

type OptionsQRCode struct {
	QrCode string `json:"qrCode"`
}

type Options struct {
	OptionsCommon
	OptionsBase64
	OptionsQRCode
	OptionsUrl
}

//go:generate stringer -type=Type
type Type uint8

const (
	TypeUnknown Type = 0

	TypeBase64 = 1
	TypeUrl    = 2
	TypeQRCode = 3
	TypeBuffer = 4
	TypeFile   = 5
	TypeStream = 6
)
