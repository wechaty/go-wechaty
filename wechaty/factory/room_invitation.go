package factory

import (
	"encoding/json"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type RoomInvitationFactory struct {
	_interface.IAccessory
}

func (r *RoomInvitationFactory) Load(id string) _interface.IRoomInvitation {
	return user.NewRoomInvitation(id, r.IAccessory)
}

func (r *RoomInvitationFactory) FromJSON(s string) (_interface.IRoomInvitation, error) {
	payload := new(schemas.RoomInvitationPayload)
	err := json.Unmarshal([]byte(s), payload)
	if err != nil {
		return nil, err
	}
	return r.FromPayload(payload), nil
}

func (r *RoomInvitationFactory) FromPayload(payload *schemas.RoomInvitationPayload) _interface.IRoomInvitation {
	r.GetPuppet().SetRoomInvitationPayload(payload)
	return r.Load(payload.Id)
}
