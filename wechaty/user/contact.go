/**
 * Go Wechaty - https://github.com/wechaty/go-wechaty
 *
 * Authors: Huan LI (李卓桓) <https://github.com/huan>
 *          Bojie LI (李博杰) <https://github.com/SilkageNet>
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

package user

import (
	"fmt"
	file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
)

type Contact struct {
	_interface.Accessory

	Id      string
	payload *schemas.ContactPayload
}

func NewContact(id string, accessory _interface.Accessory) *Contact {
	return &Contact{
		Accessory: accessory,
		Id:        id,
	}
}

func (r *Contact) Ready(forceSync bool) (err error) {
	if !forceSync && r.IsReady() {
		return nil
	}

	if forceSync {
		r.GetPuppet().ContactPayloadDirty(r.Id)
	}

	r.payload, err = r.GetPuppet().ContactPayload(r.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Contact) IsReady() bool {
	return r.payload != nil
}

func (r *Contact) Sync() error {
	return r.Ready(true)
}

func (r *Contact) String() string {
	return fmt.Sprintf("Contact<%s>", r.identity())
}

func (r *Contact) identity() string {
	identity := "loading..."
	if r.payload.Alias != "" {
		identity = r.payload.Alias
	} else if r.payload.Name != "" {
		identity = r.payload.Name
	} else if r.Id != "" {
		identity = r.Id
	}
	return identity
}

func (r *Contact) ID() string {
	return r.Id
}

func (r *Contact) Name() string {
	if r.payload == nil {
		return ""
	}
	return r.payload.Name
}

// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram)}
func (r *Contact) Say(something interface{}) (msg _interface.IMessage, err error) {
	var msgID string
	switch v := something.(type) {
	case string:
		msgID, err = r.GetPuppet().MessageSendText(r.Id, v)
	case *Contact:
		msgID, err = r.GetPuppet().MessageSendContact(r.Id, v.Id)
	case *file_box.FileBox:
		msgID, err = r.GetPuppet().MessageSendFile(r.Id, v)
	case *UrlLink:
		msgID, err = r.GetPuppet().MessageSendURL(r.Id, v.payload)
	case *MiniProgram:
		msgID, err = r.GetPuppet().MessageSendMiniProgram(r.Id, v.payload)
	default:
		return nil, fmt.Errorf("unsupported arg: %v", something)
	}
	if msgID == "" {
		return nil, nil
	}
	msg = r.GetWechaty().Message().Load(msgID)
	return msg, msg.Ready()
}

// TODO Alias()

// Friend true for friend of the bot, false for not friend of the bot
func (r *Contact) Friend() bool {
	return r.payload.Friend
}

// Type contact type
func (r *Contact) Type() schemas.ContactType {
	return r.payload.Type
}
