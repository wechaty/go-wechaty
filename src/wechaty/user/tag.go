package user

import (
  "github.com/wechaty/go-wechaty/src/wechaty"
)

type Tag struct {
  wechaty.Accessory
  TagId string
}

func NewTag(id string, accessory wechaty.Accessory) *Tag {
  if accessory.Puppet == nil {
    panic("Tag class can not be instantiated without a puppet!")
  }
  return &Tag{accessory, id}
}
