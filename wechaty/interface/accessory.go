package _interface

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
)

// Accessory accessory interface
type Accessory interface {
  GetPuppet() wechatyPuppet.IPuppetAbstract

  GetWechaty() Wechaty
}
