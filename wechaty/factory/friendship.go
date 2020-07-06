package factory

import (
	"encoding/json"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type FriendshipFactory struct {
	_interface.IAccessory
}

func (m *FriendshipFactory) Load(id string) _interface.IFriendship {
	return user.NewFriendship(id, m.IAccessory)
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

// Add send a Friend Request to a `contact` with message `hello`.
// The best practice is to send friend request once per minute.
// Remember not to do this too frequently, or your account may be blocked.
func (m *FriendshipFactory) Add(contact _interface.IContact, hello string) error {
	return m.GetPuppet().FriendshipAdd(contact.ID(), hello)
}

// FromJSON create friendShip by friendshipJson
func (m *FriendshipFactory) FromJSON(payload string) (_interface.IFriendship, error) {
	f := new(schemas.FriendshipPayload)
	err := json.Unmarshal([]byte(payload), f)
	if err != nil {
		return nil, err
	}
	return m.FromPayload(f)
}

// FromPayload create friendShip by friendshipPayload
func (m *FriendshipFactory) FromPayload(payload *schemas.FriendshipPayload) (_interface.IFriendship, error) {
	m.GetPuppet().SetFriendshipPayload(payload.Id, payload)
	friendship := m.Load(payload.Id)
	err := friendship.Ready()
	if err != nil {
		return nil, err
	}
	return friendship, nil
}
