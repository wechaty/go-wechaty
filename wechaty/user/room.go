package user

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
)

type Room struct {
	id      string
	payLoad *schemas.RoomPayload
	_interface.Accessory
}

func NewRoom(id string, accessory _interface.Accessory) *Room {
	return &Room{
		id:        id,
		Accessory: accessory,
	}
}

func (r *Room) Ready(forceSync bool) (err error) {
	if !forceSync && r.IsReady() {
		return nil
	}

	if forceSync {
		r.GetPuppet().RoomPayloadDirty(r.id)
		if err := r.GetPuppet().RoomMemberPayloadDirty(r.id); err != nil {
			return err
		}
	}

	r.payLoad, err = r.GetPuppet().RoomPayload(r.id)
	if err != nil {
		return err
	}

	memberIDs, err := r.GetPuppet().RoomMemberList(r.id)
	if err != nil {
		return err
	}

	// TODO change to concurrent
	for _, id := range memberIDs {
		if err := r.GetWechaty().Contact().Load(id).Ready(false); err != nil {
			log.Printf("Room Ready() member.Ready() rejection: %s\n", err)
		}
	}

	return nil
}

func (r *Room) IsReady() bool {
	return r.payLoad != nil
}

func (r *Room) String() string {
	str := "loading"
	if r.payLoad.Topic != "" {
		str = r.payLoad.Topic
	}
	return fmt.Sprintf("Room<%s>", str)
}

func (r *Room) ID() string {
	return r.id
}
