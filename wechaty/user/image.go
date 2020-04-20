package user

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Images struct {
  _interface.Accessory
	ImageId string
}

// NewImages create image struct
func NewImages(id string, accessory _interface.Accessory) *Images {
	if accessory.GetPuppet() == nil {
		panic("Image class can not be instantiated without a puppet!")
	}
	return &Images{accessory, id}
}

// Thumbnail message thumbnail images
func (img *Images) Thumbnail() (*file_box.FileBox, error) {
	return img.Accessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeThumbnail)
}

// HD message hd images
func (img *Images) HD() (*file_box.FileBox, error) {
	return img.Accessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeHD)
}

// Artwork message artwork images
func (img *Images) Artwork() (*file_box.FileBox, error) {
	return img.Accessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeArtwork)
}
