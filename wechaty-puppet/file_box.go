package wechatypuppet

// FileBox file struct
type FileBox struct {
}

// ToJSON struct to map
func (f *FileBox) ToJSON() map[string]interface{} {
	return nil
}

// ToFile save to file
func (f *FileBox) ToFile(path string) {
	return
}

// FromQrCode from qr code
func (f *FileBox) FromQrCode(path string) {
	return
}
