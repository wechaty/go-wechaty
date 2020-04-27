package _interface

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
)

// Accessory accessory interface
type Accessory interface {
  SetPuppet(puppet wechatyPuppet.IPuppetAbstract)
  GetPuppet() wechatyPuppet.IPuppetAbstract

  SetWechaty(wechaty Wechaty)
  GetWechaty() Wechaty
}
