package factory

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type FriendshipFactory struct {
	_interface.Accessory
}

func (m *FriendshipFactory) Load(id string) _interface.IFriendship {
	return user.NewFriendship(id, m.Accessory)
}

// Search search a Friend by phone or weixin.
func (m *FriendshipFactory) Search(query *schemas.FriendshipSearchCondition) (_interface.IContact, error) {
	contactID, err := m.GetPuppet().FriendshipSearch(query)
	if err != nil {
		return nil, err
	}
	if contactID == "" {
		return nil, nil
	}
	contact := m.GetWechaty().Contact().Load(contactID)
	err = contact.Ready(false)
	if err != nil {
		return nil, err
	}
	return contact, nil
}
