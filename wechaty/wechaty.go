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
	"errors"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	puppethostie "github.com/wechaty/go-wechaty/wechaty-puppet-hostie"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	mc "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/factory"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
	"reflect"
	"time"
)

// Wechaty ...
type Wechaty struct {
	*Option

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
	var wy = &Wechaty{events: events.New(), Option: &Option{}}

	for _, fn := range optFns {
		fn(wy.Option)
	}

	return wy
}

// register event
func (w *Wechaty) registerEvent(name schemas.PuppetEventName, f interface{}) {
	w.events.On(name, func(data ...interface{}) {
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
		puppet, err := puppethostie.NewPuppetHostie(w.puppetOption)
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
	w.message = &factory.MessageFactory{Accessory: accessory}
	w.contact = factory.NewContactFactory(accessory)
	w.room = factory.NewRoomFactory(accessory)
	w.tag = factory.NewTagFactory(accessory)
	w.friendship = &factory.FriendshipFactory{Accessory: accessory}
	w.image = &factory.ImageFactory{Accessory: accessory}
	w.urlLink = &factory.UrlLinkFactory{}
	w.roomInvitation = &factory.RoomInvitationFactory{Accessory: accessory}
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

	return nil
}

func (w *Wechaty) initPuppetEventBridge() {
	// TODO temporary
	for _, name := range schemas.GetEventNames() {
		name := name
		switch name {
		case schemas.PuppetEventNameDong:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, i[0].(*schemas.EventDongPayload).Data)
			})
		case schemas.PuppetEventNameError:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, i[0].(*schemas.EventErrorPayload).Data)
			})
		case schemas.PuppetEventNameHeartbeat:
			w.puppet.On(name, func(i ...interface{}) {
				w.emit(name, i[0].(*schemas.EventHeartbeatPayload).Data)
			})
		case schemas.PuppetEventNameLogin:
			w.puppet.On(name, func(i ...interface{}) {
				contact := w.contact.LoadSelf(i[0].(*schemas.EventLoginPayload).ContactId)
				if err := contact.Ready(false); err != nil {
					log.Printf("emit login contact.Ready err: %s\n", err.Error())
					return
				}
				w.emit(name, contact)
			})
		case schemas.PuppetEventNameLogout:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventLogoutPayload)
				contact := w.contact.LoadSelf(payload.ContactId)
				if err := contact.Ready(false); err != nil {
					log.Printf("emit logout contact.Ready err: %s\n", err.Error())
				}
				w.emit(name, contact, payload.Data)
			})
		case schemas.PuppetEventNameScan:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventScanPayload)
				w.emit(name, payload.QrCode, payload.Status, payload.Data)
			})
		case schemas.PuppetEventNameMessage:
			w.puppet.On(name, func(i ...interface{}) {
				messageID := i[0].(*schemas.EventMessagePayload).MessageId
				message := w.message.Load(messageID)
				if err := message.Ready(); err != nil {
					log.Printf("emit message message.Ready() err: %s\n", err.Error())
					return
				}
				w.emit(name, message)
			})
		case schemas.PuppetEventNameFriendship:
			w.puppet.On(name, func(i ...interface{}) {
				friendship := w.friendship.Load(i[0].(*schemas.EventFriendshipPayload).FriendshipID)
				if err := friendship.Ready(); err != nil {
					log.Printf("emit friendship friendship.Ready() err: %s\n", err.Error())
					return
				}
				w.emit(name, friendship)
			})
		case schemas.PuppetEventNameRoomInvite:
			w.puppet.On(name, func(i ...interface{}) {
				roomInvitation := w.roomInvitation.Load(i[0].(*schemas.EventRoomInvitePayload).RoomInvitationId)
				w.emit(name, roomInvitation)
			})
		case schemas.PuppetEventNameRoomJoin:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomJoinPayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomjoin room.Sync() err: %s\n", err.Error())
					return
				}
				var inviteeList []_interface.IContact
				for _, id := range payload.InviteeIdList {
					c := w.contact.Load(id)
					if err := c.Ready(false); err != nil {
						log.Printf("emit roomjoin contact.Ready() err: %s\n", err.Error())
						return
					}
					inviteeList = append(inviteeList, c)
				}
				inviter := w.contact.Load(payload.InviterId)
				if err := inviter.Ready(false); err != nil {
					log.Printf("emit roomjoin inviter.Ready() err: %s\n", err.Error())
					return
				}
				w.emit(name, room, inviteeList, inviter, time.Unix(payload.Timestamp, 0))
			})
		case schemas.PuppetEventNameRoomLeave:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomLeavePayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomleave room.Sync() err: %s\n", err.Error())
					return
				}
				var leaverList []_interface.IContact
				for _, id := range payload.RemoveeIdList {
					c := w.contact.Load(id)
					if err := c.Ready(false); err != nil {
						log.Printf("emit roomleave contact.Ready() err: %s\n", err.Error())
						return
					}
					leaverList = append(leaverList, c)
				}
				remover := w.contact.Load(payload.RemoverId)
				if err := remover.Ready(false); err != nil {
					log.Printf("emit roomleave inviter.Ready() err: %s\n", err.Error())
					return
				}
				w.emit(name, room, leaverList, remover, time.Unix(payload.Timestamp, 0))
				selfID := w.puppet.SelfID()
				for _, id := range payload.RemoveeIdList {
					if id != selfID {
						continue
					}
					w.puppet.RoomPayloadDirty(payload.RoomId)
					_ = w.puppet.RoomMemberPayloadDirty(payload.RoomId)
				}
			})
		case schemas.PuppetEventNameRoomTopic:
			w.puppet.On(name, func(i ...interface{}) {
				payload := i[0].(*schemas.EventRoomTopicPayload)
				room := w.room.Load(payload.RoomId)
				if err := room.Sync(); err != nil {
					log.Printf("emit roomtopic room.Sync() err: %s\n", err.Error())
					return
				}
				changer := w.contact.Load(payload.ChangerId)
				if err := changer.Ready(false); err != nil {
					log.Printf("emit roomtopic changer.Ready() err: %s\n", err.Error())
					return
				}
				w.emit(name, room, payload.NewTopic, payload.OldTopic, changer, time.Unix(payload.Timestamp, 0))
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

// UrlLink ...
func (w *Wechaty) UrlLink() _interface.IUrlLinkFactory {
	return w.urlLink
}

// RoomInvitation ...
func (w *Wechaty) RoomInvitation() _interface.IRoomInvitationFactory {
	return w.roomInvitation
}
