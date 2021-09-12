package _interface

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type IContactFactory interface {
	Load(id string) IContact
	LoadSelf(id string) IContactSelf
	// Find query params is string or *schemas.ContactQueryFilter
	Find(query interface{}) IContact
	// FindAll query params is string or *schemas.ContactQueryFilter
	FindAll(query interface{}) []IContact
	// Tags get tags for all contact
	Tags() []ITag
}

type IContact interface {
	// Ready is For FrameWork ONLY!
	Ready(forceSync bool) (err error)
	IsReady() bool
	// Sync force reload data for Contact, sync data from lowlevel API again.
	Sync() error
	String() string
	ID() string
	Name() string
	// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram)}
	Say(something interface{}) (msg IMessage, err error)
	// Friend true for friend of the bot, false for not friend of the bot
	Friend() bool
	Type() schemas.ContactType
	// Star check if the contact is star contact
	Star() bool
	Gender() schemas.ContactGender
	Province() string
	City() string
	// Avatar get avatar picture file stream
	Avatar() *filebox.FileBox
	// Self Check if contact is self
	Self() bool
	// Weixin get the weixin number from a contact
	// Sometimes cannot get weixin number due to weixin security mechanism, not recommend.
	Weixin() string
	// Alias get alias
	Alias() string
	// SetAlias set alias
	SetAlias(newAlias string)
}
