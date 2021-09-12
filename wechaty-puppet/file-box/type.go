package file_box

import "net/http"

// FileBoxOptionsCommon ...
// Deprecated: use filebox package
type FileBoxOptionsCommon struct {
	Name     string                 `json:"Name"`
	Metadata map[string]interface{} `json:"metadata"`
	BoxType  FileBoxType            `json:"boxType"`
}

// FileBoxOptionsBase64 ...
// Deprecated: use filebox package
type FileBoxOptionsBase64 struct {
	Base64 string `json:"base64"`
}

// FileBoxOptionsUrl ...
// Deprecated: use filebox package
type FileBoxOptionsUrl struct {
	RemoteUrl string      `json:"remoteUrl"`
	Headers   http.Header `json:"headers"`
}

// FileBoxOptionsQRCode ...
// Deprecated: use filebox package
type FileBoxOptionsQRCode struct {
	QrCode string `json:"qrCode"`
}

// FileBoxOptions ...
// Deprecated: use filebox package
type FileBoxOptions struct {
	FileBoxOptionsCommon
	FileBoxOptionsBase64
	FileBoxOptionsQRCode
	FileBoxOptionsUrl
}

//go:generate stringer -type=FileBoxType
// FileBoxType ...
// Deprecated: use filebox package
type FileBoxType uint8

const (
	FileBoxTypeUnknown FileBoxType = 0

	// Deprecated: use filebox package
	FileBoxTypeBase64 = 1
	// Deprecated: use filebox package
	FileBoxTypeUrl    = 2
	// Deprecated: use filebox package
	FileBoxTypeQRCode = 3
	// Deprecated: use filebox package
	FileBoxTypeBuffer = 4
	// Deprecated: use filebox package
	FileBoxTypeFile   = 5
	// Deprecated: use filebox package
	FileBoxTypeStream = 6
)
