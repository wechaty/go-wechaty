package wechaty

import (
	"errors"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"reflect"
	"sync"
	"time"
)

// Manage all plugins.
type PluginManager struct {
	plugins []*Plugin
}

func newPluginManager() PluginManager {
	return PluginManager{
		plugins: nil,
	}
}

func (m *PluginManager) addPlugin(p *Plugin, w *Wechaty) {
	p.Wechaty = w
	m.plugins = append(m.plugins, p)
}

// Control whether to run the plugin by its property.
// In the order of priority, emit every callback functions in every plugins.
func (m *PluginManager) emit(name schemas.PuppetEventName, i ...interface{}) {
	context := PluginContext{
		abort:              false,
		disableOncePlugins: nil,
		data:               nil,
	}
	for _, plugin := range m.plugins {
		if context.IsActive(plugin) {
			plugin.emit(name, &context, i...)
		}
		if context.abort {
			break
		}
	}
}

type PluginContext struct {
	abort              bool
	disableOncePlugins []*Plugin
	data               map[string]interface{}
}

func (c *PluginContext) IsActive(plugin *Plugin) bool {
	if plugin.enable == false {
		return false
	}
	for _, p := range c.disableOncePlugins {
		if p == plugin {
			return false
		}
	}
	return true
}

func (c *PluginContext) DisableOnce(plugin *Plugin) {
	c.disableOncePlugins = append(c.disableOncePlugins, plugin)
}

func (c *PluginContext) Abort() {
	c.abort = true
}

func (c *PluginContext) GetData(name string) interface{}{
	return c.data[name]
}


func (c *PluginContext) SetData(name string, value interface{}){
	c.data[name] = value
}

type Plugin struct {
	Wechaty *Wechaty
	mu sync.Mutex
	enable  bool
	events  events.EventEmitter
}

func NewPlugin() *Plugin {
		p := &Plugin{
			enable:  true,
			events:  events.New(),
		}
		return p
}

func (p *Plugin) SetEnable(value bool) {
	p.mu.Lock()
	p.enable = value
	p.mu.Unlock()
}

func (p *Plugin) emit(name schemas.PuppetEventName, context *PluginContext, i ...interface{}) {
	// reference: wechaty.initPuppetEventBridge()
	// TODO: when error occur, log messages may be printed more than twice.
	// TODO: some code will execute more than once.
	switch name {
	case schemas.PuppetEventNameDong:
		p.events.Emit(name, context, i[0].(*schemas.EventDongPayload).Data)
	case schemas.PuppetEventNameError:
		p.events.Emit(name, context, errors.New(i[0].(*schemas.EventErrorPayload).Data))
	case schemas.PuppetEventNameHeartbeat:
		p.events.Emit(name, context, i[0].(*schemas.EventHeartbeatPayload).Data)
	case schemas.PuppetEventNameLogin:
		contact := p.Wechaty.contact.LoadSelf(i[0].(*schemas.EventLoginPayload).ContactId)
		if err := contact.Ready(false); err != nil {
			log.Printf("emit login contact.Ready err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, contact)
	case schemas.PuppetEventNameLogout:
		payload := i[0].(*schemas.EventLogoutPayload)
		contact := p.Wechaty.contact.LoadSelf(payload.ContactId)
		if err := contact.Ready(false); err != nil {
			log.Printf("emit logout contact.Ready err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, contact, payload.Data)
	case schemas.PuppetEventNameScan:
		payload := i[0].(*schemas.EventScanPayload)
		p.events.Emit(name, context, payload.QrCode, payload.Status, payload.Data)
	case schemas.PuppetEventNameMessage:
		messageID := i[0].(*schemas.EventMessagePayload).MessageId
		message := p.Wechaty.message.Load(messageID)
		if err := message.Ready(); err != nil {
			log.Printf("emit message message.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, message)
	case schemas.PuppetEventNameFriendship:
		friendship := p.Wechaty.friendship.Load(i[0].(*schemas.EventFriendshipPayload).FriendshipID)
		if err := friendship.Ready(); err != nil {
			log.Printf("emit friendship friendship.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, friendship)
	case schemas.PuppetEventNameRoomInvite:
		roomInvitation := p.Wechaty.roomInvitation.Load(i[0].(*schemas.EventRoomInvitePayload).RoomInvitationId)
		p.events.Emit(name, context, roomInvitation)
	case schemas.PuppetEventNameRoomJoin:
		payload := i[0].(*schemas.EventRoomJoinPayload)
		room := p.Wechaty.room.Load(payload.RoomId)
		if err := room.Sync(); err != nil {
			log.Printf("emit roomjoin room.Sync() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		var inviteeList []_interface.IContact
		for _, id := range payload.InviteeIdList {
			c := p.Wechaty.contact.Load(id)
			if err := c.Ready(false); err != nil {
				log.Printf("emit roomjoin contact.Ready() err: %s\n", err.Error())
				p.Wechaty.emit(schemas.PuppetEventNameError, context, err)
				return
			}
			inviteeList = append(inviteeList, c)
		}
		inviter := p.Wechaty.contact.Load(payload.InviterId)
		if err := inviter.Ready(false); err != nil {
			log.Printf("emit roomjoin inviter.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, room, inviteeList, inviter, time.Unix(payload.Timestamp, 0))
	case schemas.PuppetEventNameRoomLeave:
		payload := i[0].(*schemas.EventRoomLeavePayload)
		room := p.Wechaty.room.Load(payload.RoomId)
		if err := room.Sync(); err != nil {
			log.Printf("emit roomleave room.Sync() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		var leaverList []_interface.IContact
		for _, id := range payload.RemoveeIdList {
			c := p.Wechaty.contact.Load(id)
			if err := c.Ready(false); err != nil {
				log.Printf("emit roomleave contact.Ready() err: %s\n", err.Error())
				p.events.Emit(schemas.PuppetEventNameError, context, err)
				return
			}
			leaverList = append(leaverList, c)
		}
		remover := p.Wechaty.contact.Load(payload.RemoverId)
		if err := remover.Ready(false); err != nil {
			log.Printf("emit roomleave inviter.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, room, leaverList, remover, time.Unix(payload.Timestamp, 0))
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
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		changer := p.Wechaty.contact.Load(payload.ChangerId)
		if err := changer.Ready(false); err != nil {
			log.Printf("emit roomtopic changer.Ready() err: %s\n", err.Error())
			p.events.Emit(schemas.PuppetEventNameError, context, err)
			return
		}
		p.events.Emit(name, context, room, payload.NewTopic, payload.OldTopic, changer, time.Unix(payload.Timestamp, 0))
	default:

	}
}

func (p *Plugin) registerEvent(name schemas.PuppetEventName, f interface{}) {
	// TODO: use wechaty events
	p.events.On(name, func(data ...interface{}) {
		defer func() {
			if err := recover(); err != nil {
				p.Wechaty.events.Emit(schemas.PuppetEventNameError, fmt.Errorf("panic: event %s %v", name, err))
			}
		}()
		values := make([]reflect.Value, 0, len(data))
		for _, v := range data {
			values = append(values, reflect.ValueOf(v))
		}

		// TODO: if active, call
		// TODO: learn reflex
		_ = reflect.ValueOf(f).Call(values)
	})
}

// reference: wechaty.go
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
