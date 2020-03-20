package user

import (
	"github.com/wechaty/go-wechaty/src/wechaty"
	wechatyPuppet "github.com/wechaty/go-wechaty/src/wechaty-puppet"
	"github.com/wechaty/go-wechaty/src/wechaty-puppet/schemas"
)

type Images struct {
	wechaty.Accessory
	ImageId string
}

// NewImages create image struct
func NewImages(id string, accessory wechaty.Accessory) *Images {
	if accessory.Puppet == nil {
		panic("Image class can not be instanciated without a puppet!")
	}
	return &Images{accessory, id}
}

// Thumbnail message thumbnail images
func (img *Images) Thumbnail() wechatyPuppet.FileBox {
	return img.Accessory.Puppet.MessageImage(img.ImageId, schemas.ImageTypeThumbnail)
}

// HD message hd images
func (img *Images) HD() wechatyPuppet.FileBox {
	return img.Accessory.Puppet.MessageImage(img.ImageId, schemas.ImageTypeHD)
}

// Artwork message artwork images
func (img *Images) Artwork() wechatyPuppet.FileBox {
	return img.Accessory.Puppet.MessageImage(img.ImageId, schemas.ImageTypeArtwork)
}
