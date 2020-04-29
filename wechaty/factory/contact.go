package factory

import (
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"sync"
)

type ContactFactory struct {
	_interface.Accessory
	pool *sync.Map
}

func NewContactFactory(accessory _interface.Accessory) *ContactFactory {
	return &ContactFactory{
		Accessory: accessory,
		pool:      &sync.Map{},
	}
}

func (c *ContactFactory) Load(id string) _interface.IContact {
	load, ok := c.pool.Load(id)
	if !ok {
		contact := user.NewContact(id, c.Accessory)
		c.pool.Store(id, contact)
		return contact
	}
	switch load.(type) {
	case *user.ContactSelf:
		return load.(*user.ContactSelf).Contact
	default:
		return load.(*user.Contact)
	}
}

func (c *ContactFactory) LoadSelf(id string) _interface.IContact {
	load, ok := c.pool.Load(id)
	if ok {
		return load.(*user.ContactSelf)
	}
	contact := user.NewContactSelf(id, c.Accessory)
	c.pool.Store(id, contact)
	return contact
}
