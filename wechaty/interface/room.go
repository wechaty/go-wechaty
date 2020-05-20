package _interface

type IRoomFactory interface {
	Load(id string) IRoom
}

type IRoom interface {
	Ready(forceSync bool) (err error)
	IsReady() bool
	String() string
	ID() string
	// Find all contacts in a room
	// params nil or string or RoomMemberQueryFilter
	MemberAll(query interface{}) ([]IContact, error)
	// Alias return contact's roomAlias in the room
	Alias(contact IContact) (string, error)
}
