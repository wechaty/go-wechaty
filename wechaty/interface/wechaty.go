package _interface

// Wechaty interface
type Wechaty interface {
	Room() IRoomFactory
	Contact() IContactFactory
	Message() IMessageFactory
	Tag() ITagFactory
}
