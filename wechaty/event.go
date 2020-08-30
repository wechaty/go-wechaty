package wechaty

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"time"
)

type (
	// EventDong ...
	EventDong       func(data string)
	PluginEventDong func(context PluginContext, data string)
	// EventError ...
	EventError       func(err error)
	PluginEventError func(context PluginContext, err error)
	// EventFriendship ...
	EventFriendship       func(friendship *user.Friendship)
	PluginEventFriendship func(context PluginContext, friendship *user.Friendship)
	// EventHeartbeat ...
	EventHeartbeat       func(data string)
	PluginEventHeartbeat func(context PluginContext, data string)
	// EventLogin ...
	EventLogin       func(user *user.ContactSelf)
	PluginEventLogin func(context PluginContext, user *user.ContactSelf)
	// EventLogout ...
	EventLogout       func(user *user.ContactSelf, reason string)
	PluginEventLogout func(context PluginContext, user *user.ContactSelf, reason string)
	// EventMessage ...
	EventMessage       func(message *user.Message)
	PluginEventMessage func(context PluginContext, message *user.Message)
	// EventReady ...
	EventReady       func()
	PluginEventReady func(context PluginContext)
	// EventRoomInvite ...
	EventRoomInvite       func(roomInvitation *user.RoomInvitation)
	PluginEventRoomInvite func(context PluginContext, roomInvitation *user.RoomInvitation)
	// EventRoomJoin ...
	EventRoomJoin       func(room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time)
	PluginEventRoomJoin func(context PluginContext, room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time)
	// EventRoomLeave ...
	EventRoomLeave       func(room *user.Room, leaverList []_interface.IContact, remover _interface.IContact, date time.Time)
	PluginEventRoomLeave func(context PluginContext, room *user.Room, leaverList []_interface.IContact, remover _interface.IContact, date time.Time)
	// EventRoomTopic ...
	EventRoomTopic       func(room *user.Room, newTopic string, oldTopic string, changer _interface.IContact, date time.Time)
	PluginEventRoomTopic func(context PluginContext, room *user.Room, newTopic string, oldTopic string, changer _interface.IContact, date time.Time)
	// EventScan ...
	EventScan       func(qrCode string, status schemas.ScanStatus, data string)
	PluginEventScan func(context PluginContext, qrCode string, status schemas.ScanStatus, data string)
	// EventStart ...
	EventStart       func()
	PluginEventStart func(context PluginContext)
	// EventStop ...
	EventStop       func()
	PluginEventStop func(context PluginContext)
)
