package wechaty_puppet

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type Puppet interface {
	MessageImage(messageId string, imageType schemas.ImageType) FileBox
}
