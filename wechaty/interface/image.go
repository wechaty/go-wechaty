package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/filebox"

type IImageFactory interface {
	Create(id string) IImage
}

type IImage interface {
	Thumbnail() (*filebox.FileBox, error)
	HD() (*filebox.FileBox, error)
	Artwork() (*filebox.FileBox, error)
}
