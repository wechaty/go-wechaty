package file_box

import (
	"encoding/base64"
	"io/ioutil"
)

type fileBoxFile struct {
	path string
}

func newFileBoxFile(path string) *fileBoxFile {
	return &fileBoxFile{path: path}
}

func (fb *fileBoxFile) toJSONMap() (map[string]interface{}, error) {
	b, err := fb.toBase64()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"base64": b,
	}, nil
}

func (fb *fileBoxFile) toBytes() ([]byte, error) {
	file, err := ioutil.ReadFile(fb.path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fb *fileBoxFile) toBase64() (string, error) {
	file, err := fb.toBytes()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
}
