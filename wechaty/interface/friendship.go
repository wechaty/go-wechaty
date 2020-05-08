package _interface

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type IFriendshipFactory interface {
	Load(id string) IFriendship
	// Search search a Friend by phone or weixin.
	Search(query *schemas.FriendshipSearchCondition) (IContact, error)
}

type IFriendship interface {
	Ready() (err error)
	IsReady() bool
}
