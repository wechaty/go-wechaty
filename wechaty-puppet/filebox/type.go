package filebox

import "net/http"

type Option func(options *Options)

// WithName set name
func WithName(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func WithOptions(o Options) Option {
	return func(options *Options) {
		*options = o
	}
}

// WithMetadata set metadata
func WithMetadata(metadata map[string]interface{}) Option {
	return func(options *Options) {
		options.Metadata = metadata
	}
}

// WithMd5 set md5
func WithMd5(md5 string) Option {
	return func(options *Options) {
		options.Md5 = md5
	}
}

// WithSize set size
func WithSize(size int64) Option {
	return func(options *Options) {
		options.Size = size
	}
}

// OptionsCommon common options
type OptionsCommon struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
	BoxType  Type                   `json:"type"`
	// Deprecated
	BoxTypeDeprecated Type   `json:"boxType"`
	Size              int64  `json:"size"`
	Md5               string `json:"md5"`
}

// OptionsBase64 base64
type OptionsBase64 struct {
	Base64 string `json:"base64"`
}

// OptionsUrl url
type OptionsUrl struct {
	RemoteUrl string      `json:"url"`
	Headers   http.Header `json:"headers"`
}

// OptionsQRCode QRCode
type OptionsQRCode struct {
	QrCode string `json:"qrCode"`
}

// OptionsUuid uuid
type OptionsUuid struct {
	Uuid string `json:"uuid"`
}

// Options ...
type Options struct {
	OptionsCommon
	OptionsBase64
	OptionsQRCode
	OptionsUrl
	OptionsUuid
}

func newOptions(options ...Option) Options {
	option := Options{}
	for _, f := range options {
		f(&option)
	}
	return option
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
	TypeUuid   = 7
)
