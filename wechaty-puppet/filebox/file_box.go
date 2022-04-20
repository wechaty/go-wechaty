package filebox

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tuotoo/qrcode"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"io"
	"io/ioutil"
	"mime"
	"net/url"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
)

var (
	// ErrToJSON err to json
	ErrToJSON = errors.New("FileBox.toJSON() only support TypeUrl,TypeQRCode,TypeBase64, TypeUuid")
	// ErrNoBase64Data no base64 data
	ErrNoBase64Data = errors.New("no Base64 data")
	// ErrNoUrl no url
	ErrNoUrl = errors.New("no url")
	// ErrNoPath no path
	ErrNoPath = errors.New("no path")
	// ErrNoQRCode no QR Code
	ErrNoQRCode = errors.New("no QR Code")
	// ErrNoUuid no uuid
	ErrNoUuid = errors.New("no uuid")
)

type fileImplInterface interface {
	toJSONMap() (map[string]interface{}, error)
	toReader() (io.Reader, error)
}

// FileBox struct
type FileBox struct {
	fileImpl fileImplInterface
	Name     string
	metadata map[string]interface{}
	boxType  Type
	mimeType string
	size     int64
	md5      string

	err error
}

func newFileBox(boxType Type, fileImpl fileImplInterface, options Options) *FileBox {
	return &FileBox{
		fileImpl: fileImpl,
		Name:     options.Name,
		metadata: options.Metadata,
		boxType:  boxType,
		size:     options.Size,
		md5:      options.Md5,
		mimeType: mime.TypeByExtension(filepath.Ext(options.Name)),
	}
}

// FromJSON create FileBox from JSON
func FromJSON(s string) *FileBox {
	options := new(Options)
	if err := json.Unmarshal([]byte(s), options); err != nil {
		err = fmt.Errorf("FromJSON json.Unmarshal: %w", err)
		return newFileBox(TypeUnknown, &fileBoxUnknown{}, newOptions()).setErr(err)
	}

	// 对未来要弃用的 json.BoxTypeDeprecated 做兼容处理
	if options.BoxTypeDeprecated != 0 && options.BoxType == 0 {
		options.BoxType = options.BoxTypeDeprecated
	}

	switch options.BoxType {
	case TypeBase64:
		return FromBase64(options.Base64, WithOptions(*options))
	case TypeQRCode:
		return FromQRCode(options.QrCode, WithOptions(*options))
	case TypeUrl:
		return FromUrl(options.RemoteUrl, WithOptions(*options))
	case TypeUuid:
		return FromUuid(options.Uuid, WithOptions(*options))
	default:
		err := fmt.Errorf("FromJSON invalid value boxType: %v", options.BoxType)
		return newFileBox(TypeUnknown, &fileBoxUnknown{}, newOptions()).setErr(err)
	}
}

// FromBase64 create FileBox from Base64
func FromBase64(encode string, options ...Option) *FileBox {
	var err error
	if encode == "" {
		err = fmt.Errorf("FromBase64 %w", ErrNoBase64Data)
	}

	o := newOptions(options...)
	if o.Name == "" {
		o.Name = "base64.dat"
	}
	o.Size = helper.Base64OrigLength(encode)
	return newFileBox(TypeBase64,
		newFileBoxBase64(encode), o).setErr(err)
}

// FromUrl create FileBox from url
func FromUrl(urlString string, options ...Option) *FileBox {
	var err error
	if urlString == "" {
		err = fmt.Errorf("FromUrl %w", ErrNoUrl)
	}

	o := newOptions(options...)
	if o.Name == "" && err == nil {
		if u, e := url.Parse(urlString); e != nil {
			err = e
		} else {
			o.Name = strings.TrimLeft(u.Path, "/")
		}
	}
	return newFileBox(TypeUrl,
		newFileBoxUrl(urlString, o.Headers), o).setErr(err)
}

// FromFile create FileBox from file
func FromFile(path string, options ...Option) *FileBox {
	var err error
	if path == "" {
		err = fmt.Errorf("FromFile %w", ErrNoPath)
	}

	o := newOptions(options...)
	if o.Name == "" {
		o.Name = path2.Base(path)
	}

	if err == nil {
		if file, e := os.Stat(path); e != nil {
			err = e
		} else {
			o.Size = file.Size()
		}
	}

	return newFileBox(TypeFile,
		newFileBoxFile(path), o).setErr(err)
}

// FromQRCode create FileBox from QRCode
func FromQRCode(qrCode string, options ...Option) *FileBox {
	var err error
	if qrCode == "" {
		err = fmt.Errorf("FromQRCode %w", ErrNoQRCode)
	}

	return newFileBox(TypeQRCode,
		newFileBoxQRCode(qrCode),
		newOptions(append(options, WithName("qrcode.png"))...)).setErr(err)
}

// FromStream from io.Reader
func FromStream(reader io.Reader, options ...Option) *FileBox {
	o := newOptions(options...)
	if o.Name == "" {
		o.Name = "stream.dat"
	}
	return newFileBox(TypeStream,
		newFileBoxStream(reader), o)
}

func FromUuid(uuid string, options ...Option) *FileBox {
	var err error
	if uuid == "" {
		err = fmt.Errorf("FromUuid %w", ErrNoUuid)
	}

	o := newOptions(options...)
	if o.Name == "" {
		o.Name = uuid + ".dat"
	}
	return newFileBox(TypeUuid, newFileBoxUuid(uuid), o).setErr(err)
}

// ToJSON to json string
func (fb *FileBox) ToJSON() (string, error) {
	if fb.err != nil {
		return "", fb.err
	}

	jsonMap := map[string]interface{}{
		"name":     fb.Name,
		"metadata": fb.metadata,
		"type":     fb.boxType,
		"boxType":  fb.boxType, //Deprecated
		"size":     fb.size,
		"md5":      fb.md5,
	}

	switch fb.boxType {
	case TypeUrl, TypeQRCode, TypeBase64, TypeUuid:
		break
	default:
		return "", ErrToJSON
	}
	implJsonMap, err := fb.fileImpl.toJSONMap()
	if err != nil {
		return "", err
	}
	for k, v := range implJsonMap {
		jsonMap[k] = v
	}
	marshal, err := json.Marshal(jsonMap)
	return string(marshal), err
}

// ToFile save to file
func (fb *FileBox) ToFile(filePath string, overwrite bool) error {
	if fb.err != nil {
		return fb.err
	}

	if filePath == "" {
		filePath = fb.Name
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(path, filePath)
	if !overwrite && helper.FileExists(fullPath) {
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
	if fb.err != nil {
		return nil, fb.err
	}

	reader, err := fb.ToReader()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

// ToBase64 to base64 string
func (fb *FileBox) ToBase64() (string, error) {
	if fb.err != nil {
		return "", fb.err
	}

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
	if fb.err != nil {
		return "", fb.err
	}

	toBase64, err := fb.ToBase64()
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("data:%s;base64,%s", fb.mimeType, toBase64), nil
}

// ToQRCode to QRCode
func (fb *FileBox) ToQRCode() (string, error) {
	if fb.err != nil {
		return "", fb.err
	}

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

// ToUuid to uuid
func (fb *FileBox) ToUuid() (string, error) {
	if fb.err != nil {
		return "", fb.err
	}

	if fb.boxType == TypeUuid {
		return fb.fileImpl.(*fileBoxUuid).uuid, nil
	}

	reader, err := fb.ToReader()
	if err != nil {
		return "", err
	}

	if uuidFromStream == nil {
		return "", errors.New("need to use filebox.SetUuidSaver() before dealing with UUID")
	}

	return uuidFromStream(reader)
}

// String ...
func (fb *FileBox) String() string {
	return fmt.Sprintf("FileBox#%s<%s>", fb.boxType, fb.Name)
}

// ToReader to io.Reader
func (fb *FileBox) ToReader() (io.Reader, error) {
	if fb.err != nil {
		return nil, fb.err
	}

	return fb.fileImpl.toReader()
}

// Type get type
func (fb *FileBox) Type() Type {
	return fb.boxType
}

// Error ret err
func (fb *FileBox) Error() error {
	return fb.err
}

func (fb *FileBox) setErr(err error) *FileBox {
	fb.err = err
	return fb
}

func (fb *FileBox) Size() {

}
