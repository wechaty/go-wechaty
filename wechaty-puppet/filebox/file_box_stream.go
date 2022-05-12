package filebox

import (
	"io"
)

var _ fileImplInterface = &fileBoxStream{}

type fileBoxStream struct {
	Reader io.Reader
}

func newFileBoxStream(reader io.Reader) *fileBoxStream {
	return &fileBoxStream{Reader: reader}
}

func (fb *fileBoxStream) toJSONMap() (map[string]interface{}, error) {
	return nil, nil
}

func (fb *fileBoxStream) toBytes() ([]byte, error) { // nolint:unused
	panic("im")
}

func (fb *fileBoxStream) toReader() (io.Reader, error) {
	return fb.Reader, nil
}
