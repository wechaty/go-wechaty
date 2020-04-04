package storage

import (
  "log"
  "reflect"
  "testing"
)

var data = map[string]interface{}{
  "key1": "key1",
  "key2": "key2",
}

func TestFileStorage_Save(t *testing.T) {
  storage := newFileStorage(t)
  err := storage.Save(data)
  if err != nil {
    t.Fatalf(err.Error())
  }
}

func TestFileStorage_Load(t *testing.T) {
  storage := newFileStorage(t)
  got, err := storage.Load()
  if err != nil {
    log.Fatalf(err.Error())
  }
  if !reflect.DeepEqual(got, data) {
    log.Fatalf("got %v expect %v", got, data)
  }
}

func TestNopStorage_Destroy(t *testing.T) {
  storage := newFileStorage(t)
  err := storage.Destroy()
  if err != nil {
    t.Fatalf(err.Error())
  }
}

func newFileStorage(t *testing.T) *FileStorage {
  storage, err := NewFileStorage("testdata/file")
  if err != nil {
    t.Fatalf(err.Error())
  }
  return storage
}
