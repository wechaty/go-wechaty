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

func (f *Friendship) Contact() *Contact {
  return NewContact(f.Accessory, f.friendshipPayloadBase.ContactId)
}

func (f *Friendship) Hello() string {
  return f.friendshipPayloadBase.Hello
}
