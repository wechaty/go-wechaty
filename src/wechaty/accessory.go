package wechaty

import (
	wechatyPuppet "github.com/wechaty/go-wechaty/src/wechaty-puppet"
)

type Accessory struct {
	Puppet wechatyPuppet.Puppet // wechat-puppet 的 Puppet 接口
}
