package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type IContactFactory interface {
	Load(id string) IContact
	LoadSelf(id string) IContact
	// Find query params is string or *schemas.ContactQueryFilter
	Find(query interface{}) IContact
	// FindAll query params is string or *schemas.ContactQueryFilter
	FindAll(query interface{}) []IContact
	// Tags get tags for all contact
	Tags() []ITag
}

type IContact interface {
	Ready(forceSync bool) (err error)
	IsReady() bool
	Sync() error
	String() string
	ID() string
	Name() string
	// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram)}
	Say(something interface{}) (msg IMessage, err error)
	// Friend true for friend of the bot, false for not friend of the bot
	Friend() bool
	Type() schemas.ContactType
}
