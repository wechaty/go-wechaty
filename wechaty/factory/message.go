package factory

import (
  _interface "github.com/wechaty/go-wechaty/wechaty/interface"
  "github.com/wechaty/go-wechaty/wechaty/user"
)

type MessageFactory struct {
  _interface.Accessory
}

func (m *MessageFactory) Load(id string) _interface.IMessage {
  return user.NewMessage(id, m.Accessory)
}
