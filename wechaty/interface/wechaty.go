package _interface

// IWechaty interface
type IWechaty interface {
	Room() IRoomFactory
	Contact() IContactFactory
	Message() IMessageFactory
	Tag() ITagFactory
	Friendship() IFriendshipFactory
	Image() IImageFactory
	UserSelf() IContactSelf
}
