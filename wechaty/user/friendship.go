package user

import (
  wechaty "github.com/wechaty/go-wechaty/wechaty"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type friendship interface {
  Contact() *Contact
}

type Friendship struct {
  wechaty.Accessory
  friendshipPayloadBase schemas.FriendshipPayloadBase
}

func NewFriendship(accessory wechaty.Accessory, friendshipPayloadBase schemas.FriendshipPayloadBase) *Friendship {
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
