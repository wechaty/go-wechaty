package _interface

type IMessageFactory interface {
	Load(id string) IMessage
}

type IMessage interface {
	Ready() (err error)
	IsReady() bool
	String() string
	Room() IRoom
	Self() bool
}
