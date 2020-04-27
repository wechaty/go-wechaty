package factory

import (
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"sync"
)

type RoomFactory struct {
	_interface.Accessory
	pool *sync.Map
}

func NewRoomFactory(accessory _interface.Accessory) *RoomFactory {
	return &RoomFactory{
		Accessory: accessory,
		pool:      &sync.Map{},
	}
}

func (r *RoomFactory) Load(id string) _interface.IRoom {
	load, ok := r.pool.Load(id)
	if ok {
		return load.(*user.Room)
	}
	room := user.NewRoom(id, r.Accessory)
	r.pool.Store(id, room)
	return room
}
