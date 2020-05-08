package user

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
)

type Friendship struct {
	_interface.Accessory
	id      string
	payload *schemas.FriendshipPayload
}

func NewFriendship(id string, accessory _interface.Accessory) *Friendship {
	return &Friendship{
		Accessory: accessory,
		id:        id,
	}
}

// Ready ...
func (f *Friendship) Ready() (err error) {
	if f.IsReady() {
		return nil
	}
	f.payload, err = f.GetPuppet().FriendshipPayload(f.id)
	if err != nil {
		return err
	}
	return f.Contact().Ready(false)
}

// IsReady ...
func (f *Friendship) IsReady() bool {
	return f.payload != nil
}

// Contact ...
func (f *Friendship) Contact() _interface.IContact {
	return f.GetWechaty().Contact().Load(f.payload.ContactId)
}

//func (f *Friendship) Hello() string {
//	return f.friendshipPayloadBase.Hello
//}
