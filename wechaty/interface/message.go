package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type IMessageFactory interface {
	Load(id string) IMessage
	// Find find message in cache
	Find(query interface{}) IMessage
	// FindAll Find message in cache
	FindAll(query *schemas.MessageQueryFilter) []IMessage
}

type IMessage interface {
	Ready() (err error)
	IsReady() bool
	String() string
	Room() IRoom
	Self() bool
}
