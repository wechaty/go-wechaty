package wechaty_puppet_padplus

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/tuotoo/qrcode"
)

// ShowQrCode console print qrcode
func ShowQrCode(qrCode string) (err error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(qrCode))
	jpegI, err := jpeg.Decode(reader)
	if err != nil {
		return err
	}
	qrMatrix, err := qrcode.DecodeImg(jpegI, filepath.Join(os.TempDir(), "tuotoo", "qrcode", uuid.NewV4().String()))
	if err != nil {
		return err
	}

	for _, point := range qrMatrix.Points {
		fmt.Print("\n  ")
		for _, p := range point {
			if p {
				fmt.Print("\033[40;40m  \033[0m")
			} else {
				fmt.Print("\033[47;30m  \033[0m")
			}
		}
	}
	fmt.Print("\n\n")
	return nil
}
