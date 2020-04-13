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
  wechatypuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "reflect"
)

// Wechaty ...
type Wechaty struct {
  eventMap map[schemas.PuppetEventName]interface{}
  puppet   *wechatypuppet.Puppet
  token    string
}

// NewWechaty ...
// instance by golang.
func NewWechaty(token string) *Wechaty {
  return &Wechaty{
    eventMap: map[schemas.PuppetEventName]interface{}{},
    token:    token,
  }
}

func (w *Wechaty) registerEvent(name schemas.PuppetEventName, event interface{}) {
  w.eventMap[name] = event
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

// OnFriendshipConfirm ...
func (w *Wechaty) OnFriendshipConfirm(f EventFriendshipConfirm) *Wechaty {
  w.registerEvent(schemas.PuppetEventNameFriendShipConfirm, f)
  return w
}

// OnFriendshipVerify ...
func (w *Wechaty) OnFriendshipVerify(f EventFriendshipVerify) *Wechaty {
  w.registerEvent(schemas.PuppetEventNameFriendShipVerify, f)
  return w
}

// OnFriendshipReceive ...
func (w *Wechaty) OnFriendshipReceive(f EventFriendshipReceive) *Wechaty {
  w.registerEvent(schemas.PuppetEventNameFriendShipReceive, f)
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

// Emit ...
func (w *Wechaty) emit(name schemas.PuppetEventName, data ...interface{}) {
  f, ok := w.eventMap[name]
  if ok {
    values := make([]reflect.Value, 0, len(data))
    for _, v := range data {
      values = append(values, reflect.ValueOf(v))
    }
    _ = reflect.ValueOf(f).Call(values)
  }
}

// Start ...
func (w *Wechaty) Start() error {
  eventParamsChan := make(chan schemas.EventParams)
  puppet, err := wechatypuppet.NewPuppet(eventParamsChan, w.token)
  if err != nil {
    return err
  }
  w.puppet = puppet
  go func() {
    for v := range eventParamsChan {
      w.emit(v.EventName, v.Params...)
    }
  }()
  return w.puppet.Start()
}
