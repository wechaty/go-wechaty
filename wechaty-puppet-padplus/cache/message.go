package cache

import (
	"errors"
	"sync"

	"github.com/wechaty/go-wechaty/wechaty-puppet-padplus/payload"
)

// MessagePayload message cache payload
type MessagePayload struct {
	payload sync.Map
}

// NewMessagePayload new message payload cache
func NewMessagePayload() *MessagePayload {
	return &MessagePayload{payload: sync.Map{}}
}

// GetMessage get message
func (p *MessagePayload) GetMessage(messageID string) (pay payload.MessagePayload, err error) {
	message, ok := p.payload.Load(messageID)
	if !ok {
		return pay, errors.New("not find message")
	}
	if pay, ok := message.(payload.MessagePayload); ok {
		return pay, nil
	}
	return pay, errors.New("contact asset error")
}

func (p *MessagePayload) Store(messageID string, message payload.MessagePayload) {
	p.payload.Store(messageID, message)
}
