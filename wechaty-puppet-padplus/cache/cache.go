package cache

import (
	"errors"
	"sync"
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
func (p *MessagePayload) GetMessage(messageID string) (message interface{}, err error) {
	message, ok := p.payload.Load(messageID)
	if !ok {
		return nil, errors.New("not find message")
	}
	return
}
