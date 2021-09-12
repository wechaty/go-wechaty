package _interface

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type IRoomFactory interface {
	// Create a new room.
	Create(contactList []IContact, topic string) (IRoom, error)
	FindAll(query *schemas.RoomQueryFilter) []IRoom
	// Find query params is string or *schemas.RoomQueryFilter
	Find(query interface{}) IRoom
	Load(id string) IRoom
}

type IRoom interface {
	// Ready is For FrameWork ONLY!
	Ready(forceSync bool) (err error)
	IsReady() bool
	String() string
	ID() string
	// Find all contacts in a room
	// params nil or string or RoomMemberQueryFilter
	MemberAll(query interface{}) ([]IContact, error)
	// Member Find all contacts in a room, if get many, return the first one.
	// query params string or RoomMemberQueryFilter
	Member(query interface{}) (IContact, error)
	// Alias return contact's roomAlias in the room
	Alias(contact IContact) (string, error)
	// Sync Force reload data for Room, Sync data from puppet API again.
	Sync() error
	// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram )}
	// mentionList @ contact list
	Say(something interface{}, mentionList ...IContact) (msg IMessage, err error)
	// Add contact in a room
	Add(contact IContact) error
	// Del delete a contact from the room
	// it works only when the bot is the owner of the room
	Del(contact IContact) error
	// Quit the room itself
	Quit() error
	// Topic get topic from the room
	Topic() string
	// Topic set topic from the room
	SetTopic(topic string) error
	// Announce get announce from the room
	Announce() (string, error)
	// Announce set announce from the room
	// It only works when bot is the owner of the room.
	SetAnnounce(text string) error
	// QrCode Get QR Code Value of the Room from the room, which can be used as scan and join the room.
	QrCode() (string, error)
	// Has check if the room has member `contact`
	Has(contact IContact) (bool, error)
	// Owner get room's owner from the room.
	Owner() IContact
	// Avatar get avatar from the room.
	Avatar() (*filebox.FileBox, error)
}
