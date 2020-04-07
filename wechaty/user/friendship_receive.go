package user

import (
  "errors"
  "github.com/wechaty/go-wechaty/wechaty"
  helper_functions "github.com/wechaty/go-wechaty/wechaty-puppet/helper-functions"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type FriendshipReceive struct {
  wechaty.Accessory
  payload *schemas.FriendshipPayloadReceive
  *Friendship
}

func NewFriendshipReceive(accessory wechaty.Accessory, payload *schemas.FriendshipPayloadReceive) *FriendshipReceive {
  return &FriendshipReceive{
    Accessory:  accessory,
    payload:    payload,
    Friendship: NewFriendship(accessory, payload.FriendshipPayloadBase),
  }
}

func (f *FriendshipReceive) Accept() {
  f.GetPuppet().FriendshipAccept(f.payload.Id)
  contact := f.Contact()
  _ = helper_functions.TryWait(func() error {
    contact.Ready(false)
    if contact.isReady() {
      return nil
    }
    return errors.New("friendshipReceive.accept() contact.ready() not ready")
  })
  contact.Sync()
}
