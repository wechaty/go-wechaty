package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"

type IImageFactory interface {
	Create(id string) IImage
}

type IImage interface {
	Thumbnail() (*file_box.FileBox, error)
	HD() (*file_box.FileBox, error)
	Artwork() (*file_box.FileBox, error)
}
