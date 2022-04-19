package filebox

import (
	"io"
	"strings"
	"testing"
)

func TestSetUuidLoader(t *testing.T) {
	const data = "hello"
	const uuid = "xxxx-xxxx"
	loader := func(uuid string) (io.Reader, error) {
		return strings.NewReader(data), nil
	}

	SetUuidLoader(loader)

	bytes, err := FromUuid(uuid).ToBytes()
	if err != nil {
		t.Log(err)
	}
	if string(bytes) != data {
		t.Errorf("got %s want %s", string(bytes), data)
	}
}

func TestSetUuidSaver(t *testing.T) {
	const data = "https://github.com/dchaofei"
	const uuid = "xxxx-xxxx"
	saver := func(reader io.Reader) (string, error) {
		return uuid, nil
	}
	SetUuidSaver(saver)

	bytes, err := FromFile("testdata/dchaofei.txt").ToUuid()
	if err != nil {
		t.Log(err)
	}
	if string(bytes) != uuid {
		t.Errorf("got %s want %s", string(bytes), data)
	}
}
