package user

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

type RoomInvitation struct {
	_interface.IAccessory
	id string
}

// NewRoomInvitation ...
func NewRoomInvitation(id string, accessory _interface.IAccessory) *RoomInvitation {
	return &RoomInvitation{
		IAccessory: accessory,
		id:         id,
	}
}

func (ri *RoomInvitation) String() string {
	id := "loading"
	if ri.id != "" {
		id = ri.id
	}
	return fmt.Sprintf("RoomInvitation#%s", id)
}

func (ri *RoomInvitation) ToStringAsync() (string, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("RoomInvitation#%s<%s,%s>", ri.id, payload.Topic, payload.InviterId), nil
}

// Accept Room Invitation
func (ri *RoomInvitation) Accept() error {
	err := ri.GetPuppet().RoomInvitationAccept(ri.id)
	if err != nil {
		return err
	}
	inviter, err := ri.Inviter()
	if err != nil {
		return err
	}
	topic, err := ri.Topic()
	if err != nil {
		return err
	}
	log.Printf("RoomInvitation accept() with room(%s) & inviter(%s) ready()", topic, inviter)
	return inviter.Ready(false)
}

// Inviter get the inviter from room invitation
func (ri *RoomInvitation) Inviter() (_interface.IContact, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return nil, err
	}
	return ri.GetWechaty().Contact().Load(payload.InviterId), nil
}

// Topic get the room topic from room invitation
func (ri *RoomInvitation) Topic() (string, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return "", err
	}
	return payload.Topic, nil
}

func (ri *RoomInvitation) MemberCount() (int, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return 0, err
	}
	return payload.MemberCount, nil
}

// MemberList list of Room Members that you known(is friend)
func (ri *RoomInvitation) MemberList() ([]_interface.IContact, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return nil, err
	}
	contactList := make([]_interface.IContact, 0, len(payload.MemberIdList))
	for _, id := range payload.MemberIdList {
		c := ri.GetWechaty().Contact().Load(id)
		if err := c.Ready(false); err != nil {
			return nil, err
		}
		contactList = append(contactList, c)
	}
	return contactList, nil
}

// Date get the invitation time
func (ri *RoomInvitation) Date() (time.Time, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return time.Time{}, err
	}
	return payload.Timestamp, nil
}

// Age returns the room invitation age in seconds
func (ri *RoomInvitation) Age() (time.Duration, error) {
	date, err := ri.Date()
	if err != nil {
		return 0, err
	}
	return time.Since(date), nil
}

func (ri *RoomInvitation) ToJson() (string, error) {
	payload, err := ri.GetPuppet().RoomInvitationPayload(ri.id)
	if err != nil {
		return "", err
	}
	marshal, err := json.Marshal(payload)
	return string(marshal), err
}
