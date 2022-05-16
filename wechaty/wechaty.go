/**
 * Go Wechaty - https://github.com/wechaty/go-wechaty
 *
 * Authors: Huan LI (李卓桓) <https://github.com/huan>
 *          Xiaoyu DING (丁小雨） <https://github.com/dingdayu>
 *
 * 2020-now @ Copyright Wechaty <https://github.com/wechaty>
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package wechaty ...
package wechaty

import (
	"context"
	"errors"
	"fmt"
	"github.com/lucsky/cuid"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	puppetservice "github.com/wechaty/go-wechaty/wechaty-puppet-service"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	mc "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/factory"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"time"
)

// Wechaty ...
type Wechaty struct {
	*Option

	// cuid
	id string

	puppet wp.IPuppetAbstract
	events events.EventEmitter

	message        _interface.IMessageFactory
	room           _interface.IRoomFactory
	contact        _interface.IContactFactory
	tag            _interface.ITagFactory
	friendship     _interface.IFriendshipFactory
	image          _interface.IImageFactory
	urlLink        _interface.IUrlLinkFactory
	roomInvitation _interface.IRoomInvitationFactory
}

// NewWechaty ...
// instance by golang.
func NewWechaty(optFns ...OptionFn) *Wechaty {
	var w = &Wechaty{events: events.New(), Option: &Option{}}

	for _, fn := range optFns {
		fn(w.Option)
	}

	w.id = cuid.New()

	return w
}

func (w *Wechaty) String() string {
	return fmt.Sprintf("Wechaty#%s", w.id)
}

// Name Wechaty bot name set by `options.name`
func (w *Wechaty) Name() string {
	if len(w.Option.name) != 0 {
		return w.Option.name
	}
	return "wechaty"
}

// register event
func (w *Wechaty) registerEvent(name schemas.PuppetEventName, f interface{}) {
	w.events.On(name, func(data ...interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic: ", err)
				log.Println(string(debug.Stack()))
				w.emit(schemas.PuppetEventNameError, NewContext(), fmt.Errorf("panic: event %s %v", name, err))
			}
		}()
		values := make([]reflect.Value, 0, len(data))
		for _, v := range data {
			values = append(values, reflect.ValueOf(v))
		}
		_ = reflect.ValueOf(f).Call(values)
	})
}

// OnScan ...
func (w *Wechaty) OnScan(f EventScan) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameScan, f)
	return w
}

// OnLogin ...
func (w *Wechaty) OnLogin(f EventLogin) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameLogin, f)
	return w
}

// OnMessage ...
func (w *Wechaty) OnMessage(f EventMessage) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameMessage, f)
	return w
}

// OnDong ...
func (w *Wechaty) OnDong(f EventDong) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameDong, f)
	return w
}

// OnError ...
func (w *Wechaty) OnError(f EventError) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameError, f)
	return w
}

// OnFriendship ...
func (w *Wechaty) OnFriendship(f EventFriendship) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameFriendship, f)
	return w
}

// OnHeartbeat ...
func (w *Wechaty) OnHeartbeat(f EventHeartbeat) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameHeartbeat, f)
	return w
}

// OnLogout ...
func (w *Wechaty) OnLogout(f EventLogout) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameLogout, f)
	return w
}

// OnReady ...
func (w *Wechaty) OnReady(f EventReady) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameReady, f)
	return w
}

// OnRoomInvite ...
func (w *Wechaty) OnRoomInvite(f EventRoomInvite) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameRoomInvite, f)
	return w
}

// OnRoomJoin ...
func (w *Wechaty) OnRoomJoin(f EventRoomJoin) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameRoomJoin, f)
	return w
}

// OnRoomLeave ...
func (w *Wechaty) OnRoomLeave(f EventRoomLeave) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameRoomLeave, f)
	return w
}

// OnRoomTopic ...
func (w *Wechaty) OnRoomTopic(f EventRoomTopic) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameRoomTopic, f)
	return w
}

// OnStart ...
func (w *Wechaty) OnStart(f EventStart) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameStart, f)
	return w
}

// OnStop ...
func (w *Wechaty) OnStop(f EventStop) *Wechaty {
	w.registerEvent(schemas.PuppetEventNameStop, f)
	return w
}

// Use loads a plugin.
func (w *Wechaty) Use(plugin *Plugin) *Wechaty {
	plugin.registerPluginEvent(w)
	return w
}

func (w *Wechaty) emit(name schemas.PuppetEventName, data ...interface{}) {
	w.events.Emit(name, data...)
}

// init puppet
func (w *Wechaty) initPuppet() error {
	if w.puppet != nil {
		log.Fatalln("Puppet already inited.")
		return nil
	}
	if w.memoryCard == nil {
		return errors.New("memory card not init")
	}

	// TODO: set puppet memory

	if w.Option.puppet == nil {
		puppet, err := puppetservice.NewPuppetService(w.puppetOption)
		if err != nil {
			return err
		}
		w.puppet = puppet
	} else {
		w.puppet = w.Option.puppet
	}

	w.initPuppetEventBridge()
	w.initPuppetAccessory()

	return nil
}

func (w *Wechaty) initPuppetAccessory() {
	accessory := &Accessory{
		puppet:  w.puppet,
		wechaty: w,
	}
	w.message = &factory.MessageFactory{IAccessory: accessory}
	w.contact = factory.NewContactFactory(accessory)
	w.room = factory.NewRoomFactory(accessory)
	w.tag = factory.NewTagFactory(accessory)
	w.friendship = &factory.FriendshipFactory{IAccessory: accessory}
	w.image = &factory.ImageFactory{IAccessory: accessory}
	w.urlLink = &factory.UrlLinkFactory{}
	w.roomInvitation = &factory.RoomInvitationFactory{IAccessory: accessory}
}

// Start ...
func (w *Wechaty) Start() error {

	var err error

	// TODO: check wechaty on, impl state events

	if w.memoryCard == nil {
		w.memoryCard, err = mc.NewMemoryCard(w.name)
		if err != nil {
			log.Println("memory card new err: ", err)
			return err
		}
	}

	err = w.memoryCard.Load()
	if err != nil {
		log.Println("memory card load err: ", err)
		return err
	}

	err = w.initPuppet()
	if err != nil {
		log.Println("memory card load err: ", err)
		return err
	}

	err = w.puppet.Start()
	if err != nil {
		log.Println("puppet start err: ", err)
		return err
	}

	// TODO: io start

	go w.emit(schemas.PuppetEventNameStart, NewContext())

	return nil
}

// DaemonStart 守护进程运行
func (w *Wechaty) DaemonStart() {
	if err := w.Start(); err != nil {
		panic(err)
	}
	var quitSig = make(chan os.Signal, 1)
	signal.Notify(quitSig, os.Interrupt)

	<-quitSig
	log.Fatal("exit.by.signal")
}

func (w *Wechaty) initPuppetEventBridge() {
	// TODO temporary
	for _, name := range schemas.GetEventNames() {
		name := name
		switch name {
		case schemas.PuppetEventNameDong:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, NewContext(), i[0].(*schemas.EventDongPayload).Data)
			})
		case schemas.PuppetEventNameError:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, NewContext(), errors.New(i[0].(*schemas.EventErrorPayload).Data))
			})
		case schemas.PuppetEventNameHeartbeat:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, NewContext(), i[0].(*schemas.EventHeartbeatPayload).Data)
			})
		case schemas.PuppetEventNameLogin:
			w.puppet.On(name, func(i ...interface{}) {
				contact := w.contact.LoadSelf(i[0].(*schemas.EventLoginPayload).ContactId)
				if err := contact.Ready(false); err != nil {
					log.Printf("emit login contact.Ready err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), contact)
			})
		case schemas.PuppetEventNameLogout:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventLogoutPayload)
				contact := w.contact.LoadSelf(payload.ContactId)
				if err := contact.Ready(false); err != nil {
					log.Printf("emit logout contact.Ready err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), contact, payload.Data)
			})
		case schemas.PuppetEventNameScan:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventScanPayload)
				w.emit(name, NewContext(), payload.QrCode, payload.Status, payload.Data)
			})
		case schemas.PuppetEventNameMessage:
			w.puppet.On(name, func(i ...interface{}) {
				messageID := i[0].(*schemas.EventMessagePayload).MessageId
				message := w.message.Load(messageID)
				if err := message.Ready(); err != nil {
					log.Printf("emit message message.Ready() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), message)
			})
		case schemas.PuppetEventNameFriendship:
			w.puppet.On(name, func(i ...interface{}) {
				friendship := w.friendship.Load(i[0].(*schemas.EventFriendshipPayload).FriendshipID)
				if err := friendship.Ready(); err != nil {
					log.Printf("emit friendship friendship.Ready() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), friendship)
			})
		case schemas.PuppetEventNameRoomInvite:
			w.puppet.On(name, func(i ...interface{}) {
				roomInvitation := w.roomInvitation.Load(i[0].(*schemas.EventRoomInvitePayload).RoomInvitationId)
				w.emit(name, NewContext(), roomInvitation)
			})
		case schemas.PuppetEventNameRoomJoin:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomJoinPayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomjoin room.Sync() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				var inviteeList []_interface.IContact
				for _, id := range payload.InviteeIdList {
					c := w.contact.Load(id)
					if err := c.Ready(false); err != nil {
						log.Printf("emit roomjoin contact.Ready() err: %s\n", err.Error())
						w.emit(schemas.PuppetEventNameError, NewContext(), err)
						return
					}
					inviteeList = append(inviteeList, c)
				}
				inviter := w.contact.Load(payload.InviterId)
				if err := inviter.Ready(false); err != nil {
					log.Printf("emit roomjoin inviter.Ready() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), room, inviteeList, inviter, time.Unix(payload.Timestamp, 0))
			})
		case schemas.PuppetEventNameRoomLeave:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomLeavePayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomleave room.Sync() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				var leaverList []_interface.IContact
				for _, id := range payload.RemoveeIdList {
					c := w.contact.Load(id)
					if err := c.Ready(false); err != nil {
						log.Printf("emit roomleave contact.Ready() err: %s\n", err.Error())
						w.emit(schemas.PuppetEventNameError, NewContext(), err)
						return
					}
					leaverList = append(leaverList, c)
				}
				remover := w.contact.Load(payload.RemoverId)
				if err := remover.Ready(false); err != nil {
					log.Printf("emit roomleave inviter.Ready() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), room, leaverList, remover, time.Unix(payload.Timestamp, 0))
				selfID := w.puppet.SelfID()
				for _, id := range payload.RemoveeIdList {
					if id != selfID {
						continue
					}
					_ = w.puppet.DirtyPayload(schemas.PayloadTypeRoom, payload.RoomId)
					_ = w.puppet.DirtyPayload(schemas.PayloadTypeRoomMember, payload.RoomId)
				}
			})
		case schemas.PuppetEventNameRoomTopic:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomTopicPayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomtopic room.Sync() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				changer := w.contact.Load(payload.ChangerId)
				if err := changer.Ready(false); err != nil {
					log.Printf("emit roomtopic changer.Ready() err: %s\n", err.Error())
					w.emit(schemas.PuppetEventNameError, NewContext(), err)
					return
				}
				w.emit(name, NewContext(), room, payload.NewTopic, payload.OldTopic, changer, time.Unix(payload.Timestamp, 0))
			})
		case schemas.PuppetEventNameDirty:
			/**
			 * https://github.com/wechaty/go-wechaty/issues/72
			 */
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventDirtyPayload)
				switch payload.PayloadType {
				case schemas.PayloadTypeRoomMember,
					schemas.PayloadTypeContact:
					if err := w.contact.Load(payload.PayloadId).Ready(true); err != nil {
						log.Printf("emit dirty contact.Ready() err: %s\n", err.Error())
						w.emit(schemas.PuppetEventNameError, NewContext(), err)
						return
					}
				case schemas.PayloadTypeRoom:
					if err := w.room.Load(payload.PayloadId).Ready(true); err != nil {
						log.Printf("emit dirty room.Ready() err: %s\n", err.Error())
						w.emit(schemas.PuppetEventNameError, NewContext(), err)
						return
					}

				case schemas.PayloadTypeFriendship:
					// Friendship has no payload
					return
				case schemas.PayloadTypeMessage:
					// Message does not need to dirty (?)
					return
				case schemas.PayloadTypeUnknown:
					fallthrough
				default:
					log.Printf("unknown payload type:  %s\n", payload.PayloadType)
				}
			})
		default:

		}
	}
}

// Room ...
func (w *Wechaty) Room() _interface.IRoomFactory {
	return w.room
}

// Message ...
func (w *Wechaty) Message() _interface.IMessageFactory {
	return w.message
}

// Contact ...
func (w *Wechaty) Contact() _interface.IContactFactory {
	return w.contact
}

// Tag ...
func (w *Wechaty) Tag() _interface.ITagFactory {
	return w.tag
}

// Friendship ...
func (w *Wechaty) Friendship() _interface.IFriendshipFactory {
	return w.friendship
}

// Image ...
func (w *Wechaty) Image() _interface.IImageFactory {
	return w.image
}

// URLLink ...
func (w *Wechaty) URLLink() _interface.IUrlLinkFactory {
	return w.urlLink
}

// RoomInvitation ...
func (w *Wechaty) RoomInvitation() _interface.IRoomInvitationFactory {
	return w.roomInvitation
}

// Puppet return puppet impl
func (w *Wechaty) Puppet() wp.IPuppetAbstract {
	return w.puppet
}

// UserSelf return contact self
func (w *Wechaty) UserSelf() _interface.IContactSelf {
	userID := w.puppet.SelfID()
	return w.Contact().LoadSelf(userID)
}

// Context ...
type Context struct {
	context.Context

	cancel             func()
	abort              bool
	disableOncePlugins []*Plugin
	data               map[string]interface{}
}

// NewContext ...
func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		abort:   false,
		Context: ctx,
		cancel:  cancel,
		data:    map[string]interface{}{},
	}
}

// IsActive returns whether the plugin is active now.
func (c *Context) IsActive(plugin *Plugin) bool {
	if !plugin.IsEnable() {
		return false
	}
	for _, p := range c.disableOncePlugins {
		if p == plugin {
			return false
		}
	}
	return true
}

// DisableOnce disables a plugin temperarily.
// The plugin will be active again(if it is enable).
func (c *Context) DisableOnce(plugin *Plugin) {
	c.disableOncePlugins = append(c.disableOncePlugins, plugin)
}

// Abort stops executing all follow-up plugins
// and terminates goroutuines which listen to Context.Done() (See go programming language context.Context. https://golang.org/pkg/context/)
func (c *Context) Abort() {
	c.abort = true
	c.cancel()
}

// GetData returns temperary data
// which only exists in the current context.
func (c *Context) GetData(name string) interface{} {
	return c.data[name]
}

// SetData sets temperary data
// which only exists in the current context.
func (c *Context) SetData(name string, value interface{}) {
	c.data[name] = value
}
