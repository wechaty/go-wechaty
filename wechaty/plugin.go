package wechaty

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
)

// PluginEvent stores the event name and the callback function.
type PluginEvent struct {
	name schemas.PuppetEventName
	f    interface{} // callback function
}

// Plugin ...
type Plugin struct {
	Wechaty *Wechaty
	mu      sync.RWMutex
	enable  bool
	events  []PluginEvent
}

// NewPlugin ...
func NewPlugin() *Plugin {
	p := &Plugin{
		enable: true,
	}
	return p
}

// SetEnable enable or disable a plugin.
func (p *Plugin) SetEnable(value bool) {
	p.mu.Lock()
	p.enable = value
	p.mu.Unlock()
}

// IsEnable returns whether the plugin is activated.
func (p *Plugin) IsEnable() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enable
}

func (p *Plugin) registerPluginEvent(wechaty *Wechaty) {
	for _, pluginEvent := range p.events {
		pluginEvent := pluginEvent
		f := func(data ...interface{}) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("panic: ", err)
					log.Println(string(debug.Stack()))
					wechaty.events.Emit(schemas.PuppetEventNameError, NewContext(), fmt.Errorf("panic: event %s %v", pluginEvent.name, err))
				}
			}()
			values := make([]reflect.Value, 0, len(data))
			for _, v := range data {
				values = append(values, reflect.ValueOf(v))
			}
			// check whether the plugin can be used.
			if values[0].Interface().(*Context).IsActive(p) &&
				!values[0].Interface().(*Context).abort {
				_ = reflect.ValueOf(pluginEvent.f).Call(values)
			}
		}
		wechaty.registerEvent(pluginEvent.name, f)
	}
}

// OnScan ...
func (p *Plugin) OnScan(f EventScan) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameScan,
		f:    f,
	})
	return p
}

// OnLogin ...
func (p *Plugin) OnLogin(f EventLogin) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameLogin,
		f:    f,
	})
	return p
}

// OnMessage ...
func (p *Plugin) OnMessage(f EventMessage) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameMessage,
		f:    f,
	})
	return p
}

// OnDong ...
func (p *Plugin) OnDong(f EventDong) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameDong,
		f:    f,
	})
	return p
}

// OnError ...
func (p *Plugin) OnError(f EventError) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameError,
		f:    f,
	})
	return p
}

// OnFriendship ...
func (p *Plugin) OnFriendship(f EventFriendship) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameFriendship,
		f:    f,
	})
	return p
}

// OnHeartbeat ...
func (p *Plugin) OnHeartbeat(f EventHeartbeat) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameHeartbeat,
		f:    f,
	})
	return p
}

// OnLogout ...
func (p *Plugin) OnLogout(f EventLogout) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameLogout,
		f:    f,
	})
	return p
}

// OnReady ...
func (p *Plugin) OnReady(f EventReady) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameReady,
		f:    f,
	})
	return p
}

// OnRoomInvite ...
func (p *Plugin) OnRoomInvite(f EventRoomInvite) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameRoomInvite,
		f:    f,
	})
	return p
}

// OnRoomJoin ...
func (p *Plugin) OnRoomJoin(f EventRoomJoin) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameRoomJoin,
		f:    f,
	})
	return p
}

// OnRoomLeave ...
func (p *Plugin) OnRoomLeave(f EventRoomLeave) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameRoomLeave,
		f:    f,
	})
	return p
}

// OnRoomTopic ...
func (p *Plugin) OnRoomTopic(f EventRoomTopic) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameRoomTopic,
		f:    f,
	})
	return p
}

// OnStart ...
func (p *Plugin) OnStart(f EventStart) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameStart,
		f:    f,
	})
	return p
}

// OnStop ...
func (p *Plugin) OnStop(f EventStop) *Plugin {
	p.events = append(p.events, PluginEvent{
		name: schemas.PuppetEventNameStop,
		f:    f,
	})
	return p
}
