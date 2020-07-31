package wechaty

import (
	"errors"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"reflect"
	"time"
)

type PluginManager struct {
	priorityChanged bool
	nextRound       bool
	plugins         []*Plugin
}

func NewPluginManager() PluginManager{
	return PluginManager{
		priorityChanged: false,
		nextRound:       false,
		plugins:         nil,
	}
}

func (m *PluginManager) SetPriority(p *Plugin, priority int) {
	m.priorityChanged = true
	p.Priority = priority
}

func (m *PluginManager) sort() {
	// TODO: sort
	// TODO: sort when priority changed(immediately, and reset priorityCahnged field)
}

func (m *PluginManager) AddPlugin(p *Plugin, w *Wechaty) {
	findResult := m.GetPlugin(p.Name)
	if findResult != nil {
		// TODO: 如何处理
		log.Fatal("a plugin name was used, please use a different name.")
	}
	p.Wechaty = w
	p.Manager = m
	m.plugins = append(m.plugins, p)
	m.sort()
}

func (m *PluginManager) GetPlugin(name string) *Plugin {
	for _, p := range m.plugins {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (m *PluginManager) NextRound() {
	m.nextRound = true
}

// In the order of priority, emit every callback functions in every plugins.
func (m *PluginManager) Emit(name schemas.PuppetEventName, i ...interface{}) {
	if m.priorityChanged {
		m.sort()
		m.priorityChanged = false
	}
	for _, plugin := range m.plugins {
		if plugin.Enable {
			plugin.Emit(name, i...)
		}
		if m.nextRound {
			m.nextRound = false
			break
		}
	}
}

type Plugin struct {
	Name            string // TODO: name是必要的吗？在创建变量的时候是可以获取plugin指针的
	Enable          bool
	Priority        int // TODO: a better type to describe Priority

	Wechaty         *Wechaty
	Manager			*PluginManager

	priorityChanged bool
	data            map[string]interface{}
	events          events.EventEmitter
}

func NewPlugin(name string) *Plugin {
	return &Plugin{
		Name:            name,
		Enable:          true,
		Priority:        0,
		priorityChanged: false,
		data:            nil,
		events:          nil,
		Wechaty:         nil,
		Manager:		 nil,
	}
}

// TODO: test usage, type convert
// TODO: mux
func (p *Plugin) GetData(name string) interface{} {
	return p.data[name]
}

func (p *Plugin) SetData(name string, newData interface{}) {
	p.data[name] = newData
}

func (p *Plugin) Emit(name schemas.PuppetEventName, i ...interface{}) {
	// reference: wechaty.initPuppetEventBridge()
	// TODO: when error occur, log messages will be printed more than twice.
	// TODO: this part of code is heavily duplicated with wechaty.initPuppetEventBridge()
	// TODO: some code will execute more than once.
	switch name {
	case schemas.PuppetEventNameDong:
		p.events.Emit(name, i[0].(*schemas.EventDongPayload).Data)
	case schemas.PuppetEventNameError:
		p.events.Emit(name, errors.New(i[0].(*schemas.EventErrorPayload).Data))
	case schemas.PuppetEventNameHeartbeat:
		p.events.Emit(name, i[0].(*schemas.EventHeartbeatPayload).Data)
	case schemas.PuppetEventNameLogin:
		contact := p.Wechaty.contact.LoadSelf(i[0].(*schemas.EventLoginPayload).ContactId)
		if err := contact.Ready(false); err != nil {
			log.Printf("emit login contact.Ready err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, contact)
	case schemas.PuppetEventNameLogout:
		payload := i[0].(*schemas.EventLogoutPayload)
		contact := p.Wechaty.contact.LoadSelf(payload.ContactId)
		if err := contact.Ready(false); err != nil {
			log.Printf("emit logout contact.Ready err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, contact, payload.Data)
	case schemas.PuppetEventNameScan:
		payload := i[0].(*schemas.EventScanPayload)
		p.events.Emit(name, payload.QrCode, payload.Status, payload.Data)
	case schemas.PuppetEventNameMessage:
		messageID := i[0].(*schemas.EventMessagePayload).MessageId
		message := p.Wechaty.message.Load(messageID)
		if err := message.Ready(); err != nil {
			log.Printf("emit message message.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, message)
	case schemas.PuppetEventNameFriendship:
		friendship := p.Wechaty.friendship.Load(i[0].(*schemas.EventFriendshipPayload).FriendshipID)
		if err := friendship.Ready(); err != nil {
			log.Printf("emit friendship friendship.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, friendship)
	case schemas.PuppetEventNameRoomInvite:
		roomInvitation := p.Wechaty.roomInvitation.Load(i[0].(*schemas.EventRoomInvitePayload).RoomInvitationId)
		p.events.Emit(name, roomInvitation)
	case schemas.PuppetEventNameRoomJoin:
		payload := i[0].(*schemas.EventRoomJoinPayload)
		room := p.Wechaty.room.Load(payload.RoomId)
		if err := room.Sync(); err != nil {
			log.Printf("emit roomjoin room.Sync() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		var inviteeList []_interface.IContact
		for _, id := range payload.InviteeIdList {
			c := p.Wechaty.contact.Load(id)
			if err := c.Ready(false); err != nil {
				log.Printf("emit roomjoin contact.Ready() err: %s\n", err.Error())
				p.Wechaty.emit(schemas.PuppetEventNameError, err)
				return
			}
			inviteeList = append(inviteeList, c)
		}
		inviter := p.Wechaty.contact.Load(payload.InviterId)
		if err := inviter.Ready(false); err != nil {
			log.Printf("emit roomjoin inviter.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, room, inviteeList, inviter, time.Unix(payload.Timestamp, 0))
	case schemas.PuppetEventNameRoomLeave:
		payload := i[0].(*schemas.EventRoomLeavePayload)
		room := p.Wechaty.room.Load(payload.RoomId)
		if err := room.Sync(); err != nil {
			log.Printf("emit roomleave room.Sync() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		var leaverList []_interface.IContact
		for _, id := range payload.RemoveeIdList {
			c := p.Wechaty.contact.Load(id)
			if err := c.Ready(false); err != nil {
				log.Printf("emit roomleave contact.Ready() err: %s\n", err.Error())
				p.events.Emit(schemas.PuppetEventNameError, err)
				return
			}
			leaverList = append(leaverList, c)
		}
		remover := p.Wechaty.contact.Load(payload.RemoverId)
		if err := remover.Ready(false); err != nil {
			log.Printf("emit roomleave inviter.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, room, leaverList, remover, time.Unix(payload.Timestamp, 0))
		selfID := p.Wechaty.puppet.SelfID()
		for _, id := range payload.RemoveeIdList {
			if id != selfID {
				continue
			}
			p.Wechaty.puppet.RoomPayloadDirty(payload.RoomId)
			_ = p.Wechaty.puppet.RoomMemberPayloadDirty(payload.RoomId)
		}
	case schemas.PuppetEventNameRoomTopic:
		payload := i[0].(*schemas.EventRoomTopicPayload)
		room := p.Wechaty.room.Load(payload.RoomId)
		if err := room.Sync(); err != nil {
			log.Printf("emit roomtopic room.Sync() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		changer := p.Wechaty.contact.Load(payload.ChangerId)
		if err := changer.Ready(false); err != nil {
			log.Printf("emit roomtopic changer.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, err)
			return
		}
		p.events.Emit(name, room, payload.NewTopic, payload.OldTopic, changer, time.Unix(payload.Timestamp, 0))
	default:

	}
}

// TODO: understand reflect
func (p *Plugin) registerEvent(name schemas.PuppetEventName, f interface{}) {
	p.events.On(name, func(data ...interface{}) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: ???
				p.events.Emit(schemas.PuppetEventNameError, fmt.Errorf("panic: event %s %v", name, err))
			}
		}()
		values := make([]reflect.Value, 0, len(data))
		for _, v := range data {
			values = append(values, reflect.ValueOf(v))
		}
		_ = reflect.ValueOf(f).Call(values)
	})
}

// TODO: this part of code is duplicated with wechtay.go, may be hard to maitain.
// OnScan ...
func (p *Plugin) OnScan(f EventScan) *Plugin {
	p.registerEvent(schemas.PuppetEventNameScan, f)
	return p
}

// OnLogin ...
func (p *Plugin) OnLogin(f EventLogin) *Plugin {
	p.registerEvent(schemas.PuppetEventNameLogin, f)
	return p
}

// OnMessage ...
func (p *Plugin) OnMessage(f EventMessage) *Plugin {
	p.registerEvent(schemas.PuppetEventNameMessage, f)
	return p
}

// OnDong ...
func (p *Plugin) OnDong(f EventDong) *Plugin {
	p.registerEvent(schemas.PuppetEventNameDong, f)
	return p
}

// OnError ...
func (p *Plugin) OnError(f EventError) *Plugin {
	p.registerEvent(schemas.PuppetEventNameError, f)
	return p
}

// OnFriendship ...
func (p *Plugin) OnFriendship(f EventFriendship) *Plugin {
	p.registerEvent(schemas.PuppetEventNameFriendship, f)
	return p
}

// OnHeartbeat ...
func (p *Plugin) OnHeartbeat(f EventHeartbeat) *Plugin {
	p.registerEvent(schemas.PuppetEventNameHeartbeat, f)
	return p
}

// OnLogout ...
func (p *Plugin) OnLogout(f EventLogout) *Plugin {
	p.registerEvent(schemas.PuppetEventNameLogout, f)
	return p
}

// OnReady ...
func (p *Plugin) OnReady(f EventReady) *Plugin {
	p.registerEvent(schemas.PuppetEventNameReady, f)
	return p
}

// OnRoomInvite ...
func (p *Plugin) OnRoomInvite(f EventRoomInvite) *Plugin {
	p.registerEvent(schemas.PuppetEventNameRoomInvite, f)
	return p
}

// OnRoomJoin ...
func (p *Plugin) OnRoomJoin(f EventRoomJoin) *Plugin {
	p.registerEvent(schemas.PuppetEventNameRoomJoin, f)
	return p
}

// OnRoomLeave ...
func (p *Plugin) OnRoomLeave(f EventRoomLeave) *Plugin {
	p.registerEvent(schemas.PuppetEventNameRoomLeave, f)
	return p
}

// OnRoomTopic ...
func (p *Plugin) OnRoomTopic(f EventRoomTopic) *Plugin {
	p.registerEvent(schemas.PuppetEventNameRoomTopic, f)
	return p
}

// OnStart ...
func (p *Plugin) OnStart(f EventStart) *Plugin {
	p.registerEvent(schemas.PuppetEventNameStart, f)
	return p
}

// OnStop ...
func (p *Plugin) OnStop(f EventStop) *Plugin {
	p.registerEvent(schemas.PuppetEventNameStop, f)
	return p
}
