package user

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Images struct {
  _interface.IAccessory
	ImageId string
}

// NewImages create image struct
func NewImages(id string, accessory _interface.IAccessory) *Images {
	if accessory.GetPuppet() == nil {
		panic("Image class can not be instantiated without a puppet!")
	}
	return &Images{accessory, id}
}

// Thumbnail message thumbnail images
func (img *Images) Thumbnail() (*filebox.FileBox, error) {
	return img.IAccessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeThumbnail)
}

// HD message hd images
func (img *Images) HD() (*filebox.FileBox, error) {
	return img.IAccessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeHD)
}

// Artwork message artwork images
func (img *Images) Artwork() (*filebox.FileBox, error) {
	return img.IAccessory.GetPuppet().MessageImage(img.ImageId, schemas.ImageTypeArtwork)
}
