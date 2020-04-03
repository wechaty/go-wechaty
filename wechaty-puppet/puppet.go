package wechatypuppet

import (
	lru "github.com/hashicorp/golang-lru"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// PuppetInterface puppet interface
type PuppetInterface interface {
	MessageImage(messageID string, imageType schemas.ImageType) FileBox
}

// Puppet puppet struce
type Puppet struct {
	PuppetInterface

	CacheMessagePayload *lru.Cache
}

// MessageList message list
func (p *Puppet) MessageList() (ks []string) {
	keys := p.CacheMessagePayload.Keys()
	for _, v := range keys {
		if k, ok := v.(string); ok {
			ks = append(ks, k)
		}
	}
	return
}

// MessageSearch search message
func (p *Puppet) MessageSearch(query schemas.MessageUserQueryFilter) []string {
	allMessageIDList := p.MessageList()
	if len(allMessageIDList) <= 0 {
		return allMessageIDList
	}

	// load
	var messagePayloadList []schemas.MessagePayload
	for _, v := range allMessageIDList {
		messagePayloadList = append(messagePayloadList, p.MessagePayload(v))
	}
	// Filter todo:: messageQueryFilterFactory
	var messageIDList []string
	for _, message := range messagePayloadList {
		if message.FromId == query.FromId || message.RoomId == query.RoomId || message.ToId == query.ToId {
			messageIDList = append(messageIDList, message.Id)
		}
	}

	return messageIDList
}

// messageQueryFilterFactory 实现正则和直接匹配
func (p *Puppet) messageQueryFilterFactory(query string) schemas.MessagePayloadFilterFunction {
	return nil
}

// MessagePayload message payload todo:: no finish
func (p *Puppet) MessagePayload(messageID string) (payload schemas.MessagePayload) {
	return payload
}
