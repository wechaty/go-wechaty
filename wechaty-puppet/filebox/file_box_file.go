package filebox

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
)

type fileBoxFile struct {
	path string
}

func newFileBoxFile(path string) *fileBoxFile {
	return &fileBoxFile{path: path}
}

func (fb *fileBoxFile) toJSONMap() map[string]interface{} {
	return nil
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

func (fb *fileBoxFile) toReader() (io.Reader, error) {
	return os.Open(fb.path)
}
