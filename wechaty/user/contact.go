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
	"log"

	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/config"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Contact struct {
	_interface.IAccessory

	Id      string
	payload *schemas.ContactPayload
}

// NewContact ...
func NewContact(id string, accessory _interface.IAccessory) *Contact {
	return &Contact{
		IAccessory: accessory,
		Id:         id,
	}
}

// Ready is For FrameWork ONLY!
func (c *Contact) Ready(forceSync bool) (err error) {
	if !forceSync && c.IsReady() {
		return nil
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

// Sync force reload data for Contact, sync data from lowlevel API again.
func (c *Contact) Sync() error {
	err := c.GetPuppet().DirtyPayload(schemas.PayloadTypeContact, c.Id)
	if err != nil {
		return err
	}
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
	case *filebox.FileBox:
		msgID, err = c.GetPuppet().MessageSendFile(c.Id, v)
	case *UrlLink:
		msgID, err = c.GetPuppet().MessageSendURL(c.Id, v.payload)
	case *MiniProgram:
		msgID, err = c.GetPuppet().MessageSendMiniProgram(c.Id, v.payload)
	default:
		return nil, fmt.Errorf("unsupported arg: %v", something)
	}
	if err != nil {
		return nil, err
	}
	if msgID == "" {
		return nil, nil
	}
	msg = c.GetWechaty().Message().Load(msgID)
	return msg, msg.Ready()
}

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
func (c *Contact) Avatar() *filebox.FileBox {
	avatar, err := c.GetPuppet().ContactAvatar(c.Id)
	if err != nil {
		log.Printf("Contact Avatar() exception: %s\n", err)
		return config.QRCodeForChatie()
	}
	return avatar
}

// Self Check if contact is self
func (c *Contact) Self() bool {
	return c.GetPuppet().SelfID() == c.Id
}

// Weixin get the weixin number from a contact
// Sometimes cannot get weixin number due to weixin security mechanism, not recommend.
func (c *Contact) Weixin() string {
	return c.payload.WeiXin
}

// Alias get alias
func (c *Contact) Alias() string {
	return c.payload.Alias
}

// SetAlias set alias
func (c *Contact) SetAlias(newAlias string) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Contact SetAlias(%s) rejected: %s\n", newAlias, err)
		}
	}()
	err = c.GetPuppet().SetContactAlias(c.Id, newAlias)
	if err != nil {
		return
	}
	err = c.GetPuppet().DirtyPayload(schemas.PayloadTypeContact, c.Id)
	if err != nil {
		log.Println("SetAlias DirtyPayload err:", err)
	}
	c.payload, err = c.GetPuppet().ContactPayload(c.Id)
	if err != nil {
		log.Println("SetAlias ContactPayload err:", err)
		return
	}
	if c.payload.Alias != newAlias {
		log.Printf("Contact SetAlias(%s) sync with server fail: set(%s) is not equal to get(%s)\n", newAlias, newAlias, c.payload.Alias)
	}
}
