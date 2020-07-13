package cache

import (
	"errors"
	"sync"

	"github.com/wechaty/go-wechaty/wechaty-puppet-padplus/payload"
)

// ContactPayload contact cache payload
type ContactPayload struct {
	payload sync.Map
}

// ContactPayload new contact payload cache
func NewContactPayload() *ContactPayload {
	return &ContactPayload{payload: sync.Map{}}
}

// GetMessage get message
func (p *ContactPayload) Load(contactId string) (contact payload.ContactPayload, err error) {
	pay, ok := p.payload.Load(contactId)
	if !ok {
		return contact, errors.New("not find contact")
	}
	if contact, ok := pay.(payload.ContactPayload); ok {
		return contact, nil
	}
	return contact, errors.New("contact asset error")
}

func (p *ContactPayload) Store(contactId string, contact payload.ContactPayload) {
	p.payload.Store(contactId, contact)
}
