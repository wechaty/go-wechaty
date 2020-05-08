package file_box

import (
	"encoding/base64"
)

type fileBoxBase64 struct {
	base64Data string
}

func newFileBoxBase64(base64Data string) *fileBoxBase64 {
	return &fileBoxBase64{base64Data: base64Data}
}

func (fb *fileBoxBase64) toJSONMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"base64": fb.base64Data,
	}, nil
}

func (fb *fileBoxBase64) toBytes() ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(fb.base64Data)
	if err != nil {
		return nil, err
	}
	return dec, nil
}
