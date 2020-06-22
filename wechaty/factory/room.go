package factory

import (
	"errors"
	"log"
	"sync"

	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

type RoomFactory struct {
	_interface.IAccessory
	pool *sync.Map
}

// NewRoomFactory ...
func NewRoomFactory(accessory _interface.IAccessory) *RoomFactory {
	return &RoomFactory{
		IAccessory: accessory,
		pool:       &sync.Map{},
	}
}

// Create a new room.
func (r *RoomFactory) Create(contactList []_interface.IContact, topic string) (_interface.IRoom, error) {
	if len(contactList) < 2 {
		return nil, errors.New("contactList need at least 2 contact to create a new room")
	}
	contactIDList := make([]string, len(contactList))
	for index, contact := range contactList {
		contactIDList[index] = contact.ID()
	}
	roomID, err := r.GetPuppet().RoomCreate(contactIDList, topic)
	if err != nil {
		return nil, err
	}
	return r.Load(roomID), nil
}

// FindAll query param is string or *schemas.RoomQueryFilter
func (r *RoomFactory) FindAll(query *schemas.RoomQueryFilter) []_interface.IRoom {
	roomIDList, err := r.GetPuppet().RoomSearch(query)
	if err != nil {
		log.Println("RoomFactory err: ", err)
		return nil
	}
	if len(roomIDList) == 0 {
		return nil
	}
	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range roomIDList {
		id := id
		async.AddTask(func() (interface{}, error) {
			room := r.Load(id)
			return room, room.Ready(false)
		})
	}
	var roomList []_interface.IRoom
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		roomList = append(roomList, v.Value.(_interface.IRoom))
	}
	return roomList
}

// Find query params is string or *schemas.RoomQueryFilter
func (r *RoomFactory) Find(query interface{}) _interface.IRoom {
	var q *schemas.RoomQueryFilter
	switch v := query.(type) {
	case string:
		q = &schemas.RoomQueryFilter{Topic: v}
	case *schemas.RoomQueryFilter:
		q = v
	default:
		log.Printf("not support query type %T\n", query)
		return nil
	}
	roomList := r.FindAll(q)
	if len(roomList) == 0 {
		return nil
	}
	for _, room := range roomList {
		if r.GetPuppet().RoomValidate(room.ID()) {
			return room
		}
	}
	return nil
}

// Load query param is string
func (r *RoomFactory) Load(id string) _interface.IRoom {
	load, ok := r.pool.Load(id)
	if ok {
		return load.(*user.Room)
	}
	room := user.NewRoom(id, r.IAccessory)
	r.pool.Store(id, room)
	return room
}
