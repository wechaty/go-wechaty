package user

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"time"
)

type RoomInvitation struct {
	_interface.Accessory
	id string
}

func NewRoomInvitation(id string, accessory _interface.Accessory) *RoomInvitation {
	return &RoomInvitation{
		Accessory: accessory,
		id:        id,
	}
}

func (ri *RoomInvitation) String() string {
	id := "loading"
	if ri.id != "" {
		id = ri.id
	}
	return fmt.Sprintf("RoomInvitation#%s", id)
}

func (ri *RoomInvitation) ToStringAsync() string {
	payload := ri.getPayload()
	return fmt.Sprintf("RoomInvitation#%s<%s,%s>", ri.id, payload.Topic, payload.InviterId)
}

func (ri *RoomInvitation) Accept() {
	ri.GetPuppet().RoomInvitationAccept(ri.id)
	inviter := ri.Inviter()
	topic := ri.ToPic()
	inviter.Ready(false)
	log.Printf("RoomInvitation accept() with room(%s) & inviter(%s) ready()", topic, inviter)
}

func (ri *RoomInvitation) Inviter() *Contact {
	return NewContact(ri.getPayload().InviterId, ri.Accessory)
}

func (ri *RoomInvitation) ToPic() string {
	return ri.getPayload().Topic
}

func (ri *RoomInvitation) MemberCount() int {
	return ri.getPayload().MemberCount
}

func (ri *RoomInvitation) MemberList() ([]*Contact, error) {
	payload := ri.getPayload()
	contactList := make([]*Contact, 0, len(payload.MemberIdList))
	for _, id := range payload.MemberIdList {
		c := NewContact(id, ri.Accessory)
		if err := c.Ready(false); err != nil {
			return nil, err
		}
		contactList = append(contactList, c)
	}
	return contactList, nil
}

func (ri *RoomInvitation) Date() time.Time {
	return time.Unix(ri.getPayload().Timestamp, 0)
}

func (ri *RoomInvitation) getPayload() *schemas.RoomInvitationPayload {
	payload, _ := ri.GetPuppet().RoomInvitationRawPayload(ri.id)
	return payload
}
