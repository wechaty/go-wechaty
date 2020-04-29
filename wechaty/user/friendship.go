package user

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
)

type friendship interface {
	Contact() *Contact
}

type Friendship struct {
	_interface.Accessory
	friendshipPayloadBase schemas.FriendshipPayloadBase
}

func NewFriendship(accessory _interface.Accessory, friendshipPayloadBase schemas.FriendshipPayloadBase) *Friendship {
	return &Friendship{
		Accessory:             accessory,
		friendshipPayloadBase: friendshipPayloadBase,
	}
}

func (f *Friendship) Contact() _interface.IContact {
	return f.GetWechaty().Contact().Load(f.friendshipPayloadBase.Id)
}

func (f *Friendship) Hello() string {
	return f.friendshipPayloadBase.Hello
}
