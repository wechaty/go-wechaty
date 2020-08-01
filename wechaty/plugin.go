package wechaty

import (
	"errors"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"reflect"
	"sort"
	"time"
)

// Manage all plugins.
type PluginManager struct {
	priorityChanged bool
	nextRound       bool
	plugins         []*Plugin
}

func NewPluginManager() PluginManager {
	return PluginManager{
		priorityChanged: false,
		nextRound:       false,
		plugins:         nil,
	}
}

// Skip all remain plugins.
func (m *PluginManager) NextRound() {
	m.nextRound = true
}

// Sort by priority.
type PluginSlice []*Plugin
func (s PluginSlice) Len() int { return len(s)}
func (s PluginSlice) Less(i, j int) bool { return s[i].priority > s[j].priority}
func (s PluginSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (m *PluginManager) sortPlugins() {
	sort.Sort(PluginSlice(m.plugins))
}

func (m *PluginManager) addPlugin(p *Plugin, w *Wechaty) {
	p.Wechaty = w
	p.Manager = m
	m.plugins = append(m.plugins, p)
	m.sortPlugins()
}

// Control whether to run the plugin by its property.
// In the order of priority, emit every callback functions in every plugins.
func (m *PluginManager) emit(name schemas.PuppetEventName, i ...interface{}) {
	if m.priorityChanged {
		m.sortPlugins()
		m.priorityChanged = false
	}

	for _, plugin := range m.plugins {
		if plugin.IsActive() {
			plugin.emit(name, i...)
		}
		if m.nextRound {
			m.nextRound = false
			break
		}
	}

	for _, plugin := range m.plugins {
		plugin.disableOnce = false
	}

}

type Plugin struct {
	enable      bool

	// Disable the plugin, until next event starts
	disableOnce bool

	// A plugin with a bigger priority value will be called earlier.
	// If two plugins have the same priority value, run the one which registered earlier first.
	priority    int

	Wechaty *Wechaty
	Manager *PluginManager

	data   map[string]interface{}
	events events.EventEmitter
}

func NewPlugin() *Plugin {
	p := &Plugin{
		enable:      true,
		disableOnce: false,
		priority:    0,
		data:        make(map[string]interface{}),
		events:      events.New(),
		Wechaty:     nil,
		Manager:     nil,
	}
	return p
}

func (p *Plugin) SetEnable(enable bool) {
	p.enable = enable
}

func (p *Plugin) IsActive() bool {
	return p.enable && !p.disableOnce
}

func (p *Plugin) DisableOnce() {
	p.disableOnce = true
}

// The default priority value is 0.
// A plugin with a bigger priority value will be called earlier.
func (p *Plugin) SetPriority(priority int) {
	if p.Manager != nil {
		p.Manager.priorityChanged = true
	}
	p.priority = priority
}

func (p *Plugin) GetData(name string) interface{} {
	return p.data[name]
}

func (p *Plugin) SetData(name string, newData interface{}) {
	p.data[name] = newData
}

func (p *Plugin) emit(name schemas.PuppetEventName, i ...interface{}) {
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
				// TODO: 这个事件的错误处理是否应该放到wechaty实例处理?
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
