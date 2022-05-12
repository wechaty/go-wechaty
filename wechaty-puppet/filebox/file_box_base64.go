package filebox

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

var _ fileImplInterface = &fileBoxBase64{}

type fileBoxBase64 struct {
	base64Data string
}

func newFileBoxBase64(base64Data string) *fileBoxBase64 {
	return &fileBoxBase64{base64Data: base64Data}
}

func (fb *fileBoxBase64) toJSONMap() (map[string]interface{}, error) {
	if fb.base64Data == "" {
		return nil, fmt.Errorf("fileBoxBase64.toJSONMap %w", ErrNoBase64Data)
	}

	return map[string]interface{}{
		"base64": fb.base64Data,
	}, nil
}

func (fb *fileBoxBase64) toBytes() ([]byte, error) { //nolint:unused
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
