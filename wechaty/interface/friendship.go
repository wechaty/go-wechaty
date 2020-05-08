package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type IFriendshipFactory interface {
	Load(id string) IFriendship
	// Search search a Friend by phone or weixin.
	Search(query *schemas.FriendshipSearchCondition) (IContact, error)
	// Add send a Friend Request to a `contact` with message `hello`.
	// The best practice is to send friend request once per minute.
	// Remember not to do this too frequently, or your account may be blocked.
	Add(contact IContact, hello string) error
	// FromJSON create friendShip by friendshipJson
	FromJSON(payload string) (IFriendship, error)
	// FromPayload create friendShip by friendshipPayload
	FromPayload(payload *schemas.FriendshipPayload) (IFriendship, error)
}

type IFriendship interface {
	Ready() (err error)
	IsReady() bool
	Contact() IContact
	String() string
	// Accept accept friend request
	Accept() error
	Type() schemas.FriendshipType
	// Hello get verify message from
	Hello() string
	// ToJSON get friendShipPayload Json
	ToJSON() (string, error)
}
