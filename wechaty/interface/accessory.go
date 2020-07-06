package _interface

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
)

// IAccessory accessory interface
type IAccessory interface {
  GetPuppet() wechatyPuppet.IPuppetAbstract

  GetWechaty() IWechaty
}
