package user

import (
	"encoding/json"
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

var ErrFriendshipNil = fmt.Errorf("ErrFriendshipNil")

type Friendship struct {
	_interface.IAccessory
	id      string
	payload *schemas.FriendshipPayload
}

// NewFriendship ...
func NewFriendship(id string, accessory _interface.IAccessory) *Friendship {
	return &Friendship{
		IAccessory: accessory,
		id:         id,
	}
}

// Ready ...
func (f *Friendship) Ready() (err error) {
	if f == nil {
		return ErrFriendshipNil
	}
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
	if f == nil {
		return false
	}
	return f.payload != nil
}

// Contact ...
func (f *Friendship) Contact() _interface.IContact {
	if f == nil || f.payload == nil {
		return nil
	}
	return f.GetWechaty().Contact().Load(f.payload.ContactId)
}

func (f *Friendship) String() string {
	if f == nil || f.payload == nil {
		return "Friendship not payload"
	}
	return fmt.Sprintf("Friendship#%s<%s>", f.payload.Type, f.payload.ContactId)
}

// Accept friend request
func (f *Friendship) Accept() error {
	if f == nil || f.payload == nil {
		return ErrFriendshipNil
	}
	if f.payload.Type != schemas.FriendshipTypeReceive {
		return fmt.Errorf("accept() need type to be FriendshipType.Receive, but it got a %s", f.payload.Type)
	}
	err := f.GetPuppet().FriendshipAccept(f.id)
	if err != nil {
		return err
	}
	return f.Contact().Sync()
}

func (f *Friendship) Type() schemas.FriendshipType {
	if f == nil || f.payload == nil {
		return schemas.FriendshipTypeUnknown
	}
	return f.payload.Type
}

// Hello get verify message from
func (f *Friendship) Hello() string {
	if f == nil || f.payload == nil {
		return ""
	}
	return f.payload.Hello
}

// toJSON get friendShipPayload Json
func (f *Friendship) ToJSON() (string, error) {
	if f == nil || f.payload == nil {
		return "", ErrFriendshipNil
	}
	marshal, err := json.Marshal(f.payload)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
