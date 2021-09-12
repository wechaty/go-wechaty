package user

import (
	"errors"

	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

type ContactSelf struct {
	*Contact
}

// NewContactSelf ...
func NewContactSelf(id string, accessory _interface.IAccessory) *ContactSelf {
	return &ContactSelf{&Contact{
		IAccessory: accessory,
		Id:         id,
	}}
}

// SetAvatar SET the avatar for a bot
func (c *ContactSelf) SetAvatar(box *filebox.FileBox) error {
	if c.Id != c.GetPuppet().SelfID() {
		return errors.New("set avatar only available for user self")
	}
	return c.GetPuppet().SetContactAvatar(c.Id, box)
}

// QRCode get bot qrcode
func (c *ContactSelf) QRCode() (string, error) {
	puppetId := c.GetPuppet().SelfID()
	if puppetId == "" {
		return "", errors.New("can not get qrcode, user might be either not logged in or already logged out")
	}
	if c.Id != puppetId {
		return "", errors.New("only can get qrcode for the login userself")
	}
	code, err := c.GetPuppet().ContactSelfQRCode()
	if err != nil {
		return "", err
	}
	return code, nil
}

// SetName change bot name
func (c *ContactSelf) SetName(name string) error {
	puppetId := c.GetPuppet().SelfID()
	if puppetId == "" {
		return errors.New("can not set name for user self, user might be either not logged in or already logged out")
	}
	if c.Id != puppetId {
		return errors.New("only can set name for user self")
	}
	err := c.GetPuppet().SetContactSelfName(name)
	if err != nil {
		return err
	}
	_ = c.Sync()
	return nil
}

// Signature change bot signature
func (c *ContactSelf) Signature(signature string) error {
	puppetId := c.GetPuppet().SelfID()
	if puppetId == "" {
		return errors.New("can not set signature for user self, user might be either not logged in or already logged out")
	}
	if c.Id != puppetId {
		return errors.New("only can change signature for user self")
	}
	return c.GetPuppet().SetContactSelfSignature(signature)
}
