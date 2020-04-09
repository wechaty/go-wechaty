package file_box

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/bitly/go-simplejson"
  "github.com/tuotoo/qrcode"
  helper_functions "github.com/wechaty/go-wechaty/wechaty-puppet/helper-functions"
  "io/ioutil"
  "mime"
  "os"
  "path/filepath"
)

type fileImplInterface interface {
  toJSONMap() map[string]interface{}
  toBytes() ([]byte, error)
}

// FileBox struct
type FileBox struct {
  fileImpl  fileImplInterface
  name      string
  metadata  map[string]interface{}
  boxType   FileBoxType
  fileBytes []byte
  mimeType  string
}

func newFileBox(common *FileBoxJsonObjectCommon, fileImpl fileImplInterface) *FileBox {
  return &FileBox{
    fileImpl: fileImpl,
    name:     common.Name,
    metadata: common.Metadata,
    boxType:  common.BoxType,
    mimeType: mime.TypeByExtension(filepath.Ext(common.Name)),
  }
}

func NewFileBoxFromJSONString(s string) (*FileBox, error) {
  newJson, err := simplejson.NewJson([]byte(s))
  if err != nil {
    return nil, err
  }
  boxType, err := newJson.Get("boxType").Int64()
  if err != nil {
    return nil, err
  }
  switch boxType {
  case FileBoxTypeBase64:
    fileBoxStruct := new(FileBoxJsonObjectBase64)
    if err := json.Unmarshal([]byte(s), fileBoxStruct); err != nil {
      return nil, err
    }
    return NewFileBoxFromJSONObjectBase64(fileBoxStruct), nil
  case FileBoxTypeQRCode:
    fileBoxStruct := new(FileBoxJsonObjectQRCode)
    if err := json.Unmarshal([]byte(s), fileBoxStruct); err != nil {
      return nil, err
    }
    return NewFileBoxFromJSONObjectQRCode(fileBoxStruct), nil
  case FileBoxTypeUrl:
    fileBoxStruct := new(FileBoxJsonObjectUrl)
    if err := json.Unmarshal([]byte(s), fileBoxStruct); err != nil {
      return nil, err
    }
    return NewFileBoxFromJSONObjectUrl(fileBoxStruct), nil
  default:
    return nil, errors.New("invalid value boxType")
  }
}

func NewFileBoxFromJSONObjectBase64(data *FileBoxJsonObjectBase64) *FileBox {
  return newFileBox(data.FileBoxJsonObjectCommon, newFileBoxBase64(data.Base64))
}

func NewFileBoxFromJSONObjectUrl(data *FileBoxJsonObjectUrl) *FileBox {
  return newFileBox(data.FileBoxJsonObjectCommon, NewFileBoxUrl(data.RemoteUrl, data.Headers))
}

func NewFileBoxFromJSONObjectQRCode(data *FileBoxJsonObjectQRCode) *FileBox {
  return newFileBox(data.FileBoxJsonObjectCommon, NewFileBoxQRCode(data.QrCode))
}

func (fb *FileBox) ToJSONString() (string, error) {
  jsonMap := map[string]interface{}{
    "name":     fb.name,
    "metadata": fb.metadata,
    "boxType":  fb.boxType,
  }
  implJsonMap := fb.fileImpl.toJSONMap()
  for k, v := range implJsonMap {
    jsonMap[k] = v
  }
  marshal, err := json.Marshal(jsonMap)
  return string(marshal), err
}

func (fb *FileBox) ToFile(filePath string, overwrite bool) error {
  if filePath == "" {
    filePath = fb.name
  }
  path, err := os.Getwd()
  if err != nil {
    return err
  }
  fullPath := filepath.Join(path, filePath)
  if !overwrite && helper_functions.FileExists(fullPath) {
    return os.ErrExist
  }
  fileBytes, err := fb.ToBytes()
  if err != nil {
    return err
  }
  return ioutil.WriteFile(filePath, fileBytes, os.ModePerm)
}

func (fb *FileBox) ToBytes() ([]byte, error) {
  if fb.fileBytes != nil {
    return fb.fileBytes, nil
  }
  toBytes, err := fb.fileImpl.toBytes()
  if err != nil {
    return nil, err
  }
  fb.fileBytes = toBytes
  return fb.fileBytes, nil
}

func (fb *FileBox) ToBase64() (string, error) {
  fileBytes, err := fb.ToBytes()
  if err != nil {
    return "", err
  }
  return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func (fb *FileBox) ToDataURL() (string, error) {
  toBase64, err := fb.ToBase64()
  if err != nil {
    return "", nil
  }
  return fmt.Sprintf("data:%s;base64,%s", fb.mimeType, toBase64), nil
}

func (fb *FileBox) ToQrCode() (string, error) {
  fileBytes, err := fb.ToBytes()
  if err != nil {
    return "", err
  }
  decode, err := qrcode.Decode(bytes.NewReader(fileBytes))
  if err != nil {
    return "", nil
  }
  return decode.Content, nil
}
