package filebox

import "io"

var _ fileImplInterface = &fileBoxUnknown{}

type fileBoxUnknown struct {
}

func (f fileBoxUnknown) toJSONMap() (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (f fileBoxUnknown) toReader() (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}
