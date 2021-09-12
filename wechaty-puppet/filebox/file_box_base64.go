package filebox

import (
	"encoding/base64"
	"io"
	"strings"
)

type fileBoxBase64 struct {
	base64Data string
}

func newFileBoxBase64(base64Data string) *fileBoxBase64 {
	return &fileBoxBase64{base64Data: base64Data}
}

func (fb *fileBoxBase64) toJSONMap() map[string]interface{} {
	return map[string]interface{}{
		"base64": fb.base64Data,
	}
}

func (fb *fileBoxBase64) toBytes() ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(fb.base64Data)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func (fb *fileBoxBase64) toReader() (io.Reader, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fb.base64Data))
	return reader, nil
}
