package factory

import (
	"fmt"
	"log"
	"sync"

	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type ContactFactory struct {
	_interface.IAccessory
	pool *sync.Map
}

// NewContactFactory ...
func NewContactFactory(accessory _interface.IAccessory) *ContactFactory {
	return &ContactFactory{
		IAccessory: accessory,
		pool:       &sync.Map{},
	}
}

// Load query param is string
func (c *ContactFactory) Load(id string) _interface.IContact {
	load, ok := c.pool.Load(id)
	if !ok {
		contact := user.NewContact(id, c.IAccessory)
		c.pool.Store(id, contact)
		return contact
	}
	switch v := load.(type) {
	case *user.ContactSelf:
		return v.Contact
	case *user.Contact:
		return v
	default:
		panic(fmt.Sprintf("ContactFactory Load unknow type: %#v", v))
	}
}

// LoadSelf query param is string
func (c *ContactFactory) LoadSelf(id string) _interface.IContactSelf {
	load, ok := c.pool.Load(id)
	if !ok {
		contact := user.NewContactSelf(id, c.IAccessory)
		c.pool.Store(id, contact)
		return contact
	}
	switch v := load.(type) {
	case *user.ContactSelf:
		return v
	case *user.Contact:
		return &user.ContactSelf{Contact: v}
	default:
		panic(fmt.Sprintf("ContactFactory LoadSelf unknow type: %#v", v))
	}
}

// Find query params is string or *schemas.ContactQueryFilter
func (c *ContactFactory) Find(query interface{}) _interface.IContact {
	contacts := c.FindAll(query)
	if len(contacts) == 0 {
		return nil
	}
	if len(contacts) > 1 {
		log.Printf("Contact Find() got more than one(%d) result\n", len(contacts))
	}
	for _, v := range contacts {
		if c.GetPuppet().ContactValidate(v.ID()) {
			return v
		}
	}
	return nil
}

// FindAll query params is string or *schemas.ContactQueryFilter
func (c *ContactFactory) FindAll(query interface{}) []_interface.IContact {
	contactIds, err := c.GetPuppet().ContactSearch(query, nil)
	if err != nil {
		log.Printf("Contact c.GetPuppet().ContactSearch() rejected: %s\n", err)
		return nil
	}

	if len(contactIds) == 0 {
		return nil
	}

	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range contactIds {
		id := id
		async.AddTask(func() (interface{}, error) {
			contact := c.Load(id)
			return contact, contact.Ready(false)
		})
	}

	var contacts []_interface.IContact
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		contacts = append(contacts, v.Value.(_interface.IContact))
	}
	return contacts
}

// Tags get tags for all contact
func (c *ContactFactory) Tags() []_interface.ITag {
	tagIDList, err := c.GetPuppet().TagContactList("")
	if err != nil {
		log.Printf("ContactFactory Tags() exception: %s\n", err)
		return nil
	}
	tagList := make([]_interface.ITag, 0, len(tagIDList))
	for _, id := range tagIDList {
		tagList = append(tagList, c.GetWechaty().Tag().Load(id))
	}
	return tagList
}
