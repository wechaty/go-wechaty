package _interface

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
)

// Accessory accessory interface
type Accessory interface {
  SetPuppet(puppet wechatyPuppet.Puppet)
  GetPuppet() *wechatyPuppet.Puppet

  SetWechaty(wechaty Wechaty)
  GetWechaty() Wechaty
}
