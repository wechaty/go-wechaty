package filebox

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/tuotoo/qrcode"
	helper_functions "github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"path/filepath"
)

// ErrToJSON err ToJSON
var ErrToJSON = errors.New("FileBox.toJSON() only support TypeUrl,TypeQRCode,TypeBase64")

type fileImplInterface interface {
	toJSONMap() map[string]interface{}
	toReader() (io.Reader, error)
}

// FileBox struct
type FileBox struct {
	fileImpl fileImplInterface
	Name     string
	metadata map[string]interface{}
	boxType  Type
	mimeType string
}

func newFileBox(common *OptionsCommon, fileImpl fileImplInterface) *FileBox {
	return &FileBox{
		fileImpl: fileImpl,
		Name:     common.Name,
		metadata: common.Metadata,
		boxType:  common.BoxType,
		mimeType: mime.TypeByExtension(filepath.Ext(common.Name)),
	}
}

// FromJSON create FileBox from JSON
func FromJSON(s string) (*FileBox, error) {
	newJson, err := simplejson.NewJson([]byte(s))
	if err != nil {
		return nil, err
	}
	boxType, err := newJson.Get("boxType").Int64()
	if err != nil {
		return nil, err
	}
	options := new(Options)
	if err := json.Unmarshal([]byte(s), options); err != nil {
		return nil, err
	}
	switch boxType {
	case TypeBase64:
		return FromBase64(options.Base64, options.Name), nil
	case TypeQRCode:
		return FromQRCode(options.QrCode), nil
	case TypeUrl:
		return FromUrl(options.RemoteUrl, options.Name, nil)
	default:
		return nil, errors.New("invalid value boxType")
	}
}

// FromBase64 create FileBox from Base64
func FromBase64(base64 string, name string) *FileBox {
	return newFileBox(&OptionsCommon{
		Name:    name,
		BoxType: TypeBase64,
	}, newFileBoxBase64(base64))
}

// FromUrl create FileBox from url
func FromUrl(urlString string, name string, headers http.Header) (*FileBox, error) {
	if name == "" {
		u, err := url.Parse(urlString)
		if err != nil {
			return nil, err
		}
		name = u.Path
	}
	return newFileBox(&OptionsCommon{
		Name:    name,
		BoxType: TypeUrl,
	}, newFileBoxUrl(urlString, headers)), nil
}

// FromFile create FileBox from file
func FromFile(path, name string) *FileBox {
	if name == "" {
		name = path2.Base(path)
	}
	return newFileBox(&OptionsCommon{
		Name:    name,
		BoxType: TypeFile,
	}, newFileBoxFile(path))
}

// FromQRCode create FileBox from QRCode
func FromQRCode(qrCode string) *FileBox {
	return newFileBox(&OptionsCommon{
		Name:    "qrcode.png",
		BoxType: TypeQRCode,
	}, newFileBoxQRCode(qrCode))
}

// FromStream from io.Reader
func FromStream(reader io.Reader, name string) *FileBox {
	return newFileBox(&OptionsCommon{
		Name:    name,
		BoxType: TypeStream,
	}, newFileBoxStream(reader))
}

// ToJSON to json string
func (fb *FileBox) ToJSON() (string, error) {
	jsonMap := map[string]interface{}{
		"name":     fb.Name,
		"metadata": fb.metadata,
		"type":     fb.boxType,
	}
	switch fb.boxType {
	case TypeUrl, TypeQRCode, TypeBase64:
	default:
		return "", ErrToJSON
	}
	implJsonMap := fb.fileImpl.toJSONMap()
	for k, v := range implJsonMap {
		jsonMap[k] = v
	}
	marshal, err := json.Marshal(jsonMap)
	return string(marshal), err
}

// ToFile save to file
func (fb *FileBox) ToFile(filePath string, overwrite bool) error {
	if filePath == "" {
		filePath = fb.Name
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(path, filePath)
	if !overwrite && helper_functions.FileExists(fullPath) {
		return os.ErrExist
	}

	reader, err := fb.ToReader()
	if err != nil {
		return err
	}

	writer := bufio.NewReader(reader)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := writer.WriteTo(file); err != nil {
		return err
	}
	return nil
}

// ToBytes to bytes
func (fb *FileBox) ToBytes() ([]byte, error) {
	reader, err := fb.ToReader()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

// ToBase64 to base64 string
func (fb *FileBox) ToBase64() (string, error) {
	if fb.boxType == TypeBase64 {
		return fb.fileImpl.(*fileBoxBase64).base64Data, nil
	}

	fileBytes, err := fb.ToBytes()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

// ToDataURL to dataURL
func (fb *FileBox) ToDataURL() (string, error) {
	toBase64, err := fb.ToBase64()
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("data:%s;base64,%s", fb.mimeType, toBase64), nil
}

// ToQRCode to QRCode
func (fb *FileBox) ToQRCode() (string, error) {
	reader, err := fb.ToReader()
	if err != nil {
		return "", err
	}
	decode, err := qrcode.Decode(reader)
	if err != nil {
		return "", nil
	}
	return decode.Content, nil
}

// String ...
func (fb *FileBox) String() string {
	return fmt.Sprintf("FileBox#%s<%s>", fb.boxType, fb.Name)
}

// ToReader to io.Reader
func (fb *FileBox) ToReader() (io.Reader, error) {
	return fb.fileImpl.toReader()
}

func (fb *FileBox) Type() Type {
	return fb.boxType
}
