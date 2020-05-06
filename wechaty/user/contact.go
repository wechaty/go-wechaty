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
	"github.com/wechaty/go-wechaty/wechaty/config"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"log"
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

func (c *Contact) Ready(forceSync bool) (err error) {
	if !forceSync && c.IsReady() {
		return nil
	}

	if forceSync {
		c.GetPuppet().ContactPayloadDirty(c.Id)
	}

	c.payload, err = c.GetPuppet().ContactPayload(c.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Contact) IsReady() bool {
	return c.payload != nil
}

func (c *Contact) Sync() error {
	return c.Ready(true)
}

func (c *Contact) String() string {
	return fmt.Sprintf("Contact<%s>", c.identity())
}

func (c *Contact) identity() string {
	identity := "loading..."
	if c.payload.Alias != "" {
		identity = c.payload.Alias
	} else if c.payload.Name != "" {
		identity = c.payload.Name
	} else if c.Id != "" {
		identity = c.Id
	}
	return identity
}

func (c *Contact) ID() string {
	return c.Id
}

func (c *Contact) Name() string {
	if c.payload == nil {
		return ""
	}
	return c.payload.Name
}

// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram)}
func (c *Contact) Say(something interface{}) (msg _interface.IMessage, err error) {
	var msgID string
	switch v := something.(type) {
	case string:
		msgID, err = c.GetPuppet().MessageSendText(c.Id, v)
	case *Contact:
		msgID, err = c.GetPuppet().MessageSendContact(c.Id, v.Id)
	case *file_box.FileBox:
		msgID, err = c.GetPuppet().MessageSendFile(c.Id, v)
	case *UrlLink:
		msgID, err = c.GetPuppet().MessageSendURL(c.Id, v.payload)
	case *MiniProgram:
		msgID, err = c.GetPuppet().MessageSendMiniProgram(c.Id, v.payload)
	default:
		return nil, fmt.Errorf("unsupported arg: %v", something)
	}
	if msgID == "" {
		return nil, nil
	}
	msg = c.GetWechaty().Message().Load(msgID)
	return msg, msg.Ready()
}

// TODO Alias()

// Friend true for friend of the bot, false for not friend of the bot
func (c *Contact) Friend() bool {
	return c.payload.Friend
}

// Type contact type
func (c *Contact) Type() schemas.ContactType {
	return c.payload.Type
}

// Star check if the contact is star contact
func (c *Contact) Star() bool {
	return c.payload.Star
}

// Gender gender
func (c *Contact) Gender() schemas.ContactGender {
	return c.payload.Gender
}

// Province Get the region 'province' from a contact
func (c *Contact) Province() string {
	return c.payload.Province
}

// City Get the region 'city' from a contact
func (c *Contact) City() string {
	return c.payload.City
}

// Avatar get avatar picture file stream
func (c *Contact) Avatar() *file_box.FileBox {
	avatar, err := c.GetPuppet().GetContactAvatar(c.Id)
	if err != nil {
		log.Printf("Contact Avatar() exception: %s\n", err)
		return config.QRCodeForChatie()
	}
	return avatar
}
