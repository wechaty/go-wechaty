package wechaty_puppet

import (
	lru "github.com/hashicorp/golang-lru"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type PuppetInterface interface {
	MessageImage(messageId string, imageType schemas.ImageType) FileBox
}

type Puppet struct {
	PuppetInterface

	CacheMessagePayload *lru.Cache
}

// MessageList
func (p *Puppet) MessageList() (ks []string) {
	keys := p.CacheMessagePayload.Keys()
	for _, v := range keys {
		if k, ok := v.(string); ok {
			ks = append(ks, k)
		}
	}
	return
}

func (p *Puppet) MessageSearch(query schemas.MessageUserQueryFilter) []string {
	allMessageIdList := p.MessageList()
	if len(allMessageIdList) <= 0 {
		return allMessageIdList
	}

	// load
	var messagePayloadList []schemas.MessagePayload
	for _, v := range allMessageIdList {
		messagePayloadList = append(messagePayloadList, p.MessagePayload(v))
	}
	// Filter todo:: messageQueryFilterFactory
	var messageIdList []string
	for _, message := range messagePayloadList {
		if message.FromId == query.FromId || message.RoomId == query.RoomId || message.ToId == query.ToId {
			messageIdList = append(messageIdList, message.Id)
		}
	}

	return messageIdList
}

// messageQueryFilterFactory 实现正则和直接匹配
func (p *Puppet) messageQueryFilterFactory(query string) schemas.MessagePayloadFilterFunction {
	return nil
}

// todo:: no finish
func (p *Puppet) MessagePayload(messageId string) (payload schemas.MessagePayload) {
	return payload
}
