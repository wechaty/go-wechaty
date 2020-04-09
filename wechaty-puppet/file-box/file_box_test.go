package file_box

import (
  "encoding/base64"
  "io/ioutil"
  "log"
  "os"
  "reflect"
  "testing"
)

func TestFileBox_ToFile(t *testing.T) {
  expect := "test content"
  fileBox := NewFileBoxFromJSONObjectBase64(&FileBoxJsonObjectBase64{
    FileBoxJsonObjectCommon: &FileBoxJsonObjectCommon{
      Name:     "test.text",
      Metadata: nil,
      BoxType:  FileBoxTypeBase64,
    },
    Base64: base64.StdEncoding.EncodeToString([]byte(expect)),
  })
  const filename = "testdata/test.text"
  t.Run("toFile success", func(t *testing.T) {
    err := fileBox.ToFile(filename, true)
    if err != nil {
      log.Fatal(err)
    }
    file, err := os.Open(filename)
    if err != nil {
      log.Fatal(err)
    }
    got, err := ioutil.ReadAll(file)
    if err != nil {
      log.Fatal(err)
    }
    if expect != string(got) {
      log.Fatalf("got %s expect %s", got, expect)
    }
  })
  t.Run("file exists", func(t *testing.T) {
    err := fileBox.ToFile(filename, false)
    if err != os.ErrExist {
      log.Fatalf("got %s expect %s", err, os.ErrExist)
    }
  })
}

func TestNewFileBoxFromJSONString(t *testing.T) {
  tests := []struct {
    jsonString     string
    expectFileImpl reflect.Type
  }{
    {
      jsonString: `{
"name":"test.png",
"metadata": null,
"boxType":1,
"base64":"dGVzdCBjb250ZW50"
}`,
      expectFileImpl: reflect.TypeOf(new(fileBoxBase64)),
    },
    {
      jsonString: `{
"name":"test.png",
"metadata": null,
"boxType":2,
"remoteUrl":"http://www.example.com",
"header":null
}`,
      expectFileImpl: reflect.TypeOf(new(fileBoxUrl)),
    },
    {
      jsonString: `{
"name":"test.png",
"metadata": null,
"boxType":3,
"qrCode":"test content"
}`,
      expectFileImpl: reflect.TypeOf(new(fileBoxQRCode)),
    },
  }
  for _, t := range tests {
    fileBox, err := NewFileBoxFromJSONString(t.jsonString)
    if err != nil {
      log.Fatal(err)
    }
    gotReflectType := reflect.TypeOf(fileBox.fileImpl)
    if gotReflectType != t.expectFileImpl {
      log.Fatalf("got %v expect %v", gotReflectType, t.expectFileImpl)
    }
  }
}
