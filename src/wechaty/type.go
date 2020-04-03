package wechaty

import (
	"github.com/wechaty/go-wechaty/src/wechaty/user"
)

type Sayable interface {
	Say(text string, replyTo ...user.Contact)
}
