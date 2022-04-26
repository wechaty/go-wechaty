package filebox

import (
	"errors"
	"fmt"
	"io"
)

var _ fileImplInterface = fileBoxUuid{}

// UuidLoader uuid loader
type UuidLoader func(uuid string) (io.Reader, error)

// UuidSaver uuid saver
type UuidSaver func(reader io.Reader) (uuid string, err error)

var uuidToStream UuidLoader
var uuidFromStream UuidSaver

// SetUuidLoader set uuid loader
func SetUuidLoader(loader UuidLoader) {
	uuidToStream = loader
}

// SetUuidSaver set uuid saver
func SetUuidSaver(saver UuidSaver) {
	uuidFromStream = saver
}

type fileBoxUuid struct {
	uuid string
}

func (f fileBoxUuid) toJSONMap() (map[string]interface{}, error) {
	if f.uuid == "" {
		return nil, fmt.Errorf("fileBoxUuid.toJSONMap %w", ErrNoUuid)
	}
	return map[string]interface{}{
		"uuid": f.uuid,
	}, nil
}

func (f fileBoxUuid) toReader() (io.Reader, error) {
	if uuidToStream == nil {
		return nil, errors.New("need to call filebox.setUuidLoader() to set UUID loader first")
	}
	return uuidToStream(f.uuid)
}

func newFileBoxUuid(uuid string) *fileBoxUuid {
	return &fileBoxUuid{uuid: uuid}
}
