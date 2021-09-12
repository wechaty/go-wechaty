package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/filebox"

type IContactSelfFactory interface {
	IContactFactory
}

type IContactSelf interface {
	IContact
	SetAvatar(box *filebox.FileBox) error
	// QRCode get bot qrcode
	QRCode() (string, error)
	SetName(name string) error
	// Signature change bot signature
	Signature(signature string) error
}
