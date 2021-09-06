package file_box

import (
	"encoding/base64"
)

type fileBoxStream struct {
	bytes []byte
}

func newFileBoxStream(bytes []byte) *fileBoxStream {
	return &fileBoxStream{bytes: bytes}
}

func (fb *fileBoxStream) toJSONMap() (map[string]interface{}, error) {
	bytes, err := fb.toBytes()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"base64": base64.StdEncoding.EncodeToString(bytes),
	}, nil
}

func (fb *fileBoxStream) toBytes() ([]byte, error) {
	return fb.bytes, nil
}
