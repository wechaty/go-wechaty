package wechaty

import (
	wechatyPuppet "github.com/wechaty/go-wechaty/src/wechaty-puppet"
)

type Accessory interface {
	SetPuppet(puppet wechatyPuppet.Puppet)
	GetPuppet() *wechatyPuppet.Puppet

	SetWechaty(wechaty Wechaty)
	GetWechaty() *Wechaty
}
