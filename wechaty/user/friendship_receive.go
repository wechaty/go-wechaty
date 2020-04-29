package user

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
)

type FriendshipReceive struct {
	_interface.Accessory
	payload *schemas.FriendshipPayloadReceive
	*Friendship
}

func NewFriendshipReceive(accessory _interface.Accessory, payload *schemas.FriendshipPayloadReceive) *FriendshipReceive {
	return &FriendshipReceive{
		Accessory:  accessory,
		payload:    payload,
		Friendship: NewFriendship(accessory, payload.FriendshipPayloadBase),
	}
}

func (f *FriendshipReceive) Accept() error {
	err := f.GetPuppet().FriendshipAccept(f.payload.Id)
	if err != nil {
		return err
	}
	return f.Contact().Sync()
}
