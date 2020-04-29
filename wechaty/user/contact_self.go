package user

import _interface "github.com/wechaty/go-wechaty/wechaty/interface"

type ContactSelf struct {
	*Contact
}

func NewContactSelf(id string, accessory _interface.Accessory) *ContactSelf {
	return &ContactSelf{&Contact{
		Accessory: accessory,
		Id:        id,
	}}
}
