package _interface

import (
	"time"

	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type IRoomInvitationFactory interface {
	Load(id string) IRoomInvitation
	FromJSON(s string) (IRoomInvitation, error)
	FromPayload(payload *schemas.RoomInvitationPayload) IRoomInvitation
}

type IRoomInvitation interface {
	String() string
	ToStringAsync() (string, error)
	Accept() error
	Inviter() (IContact, error)
	Room() (IRoom, error)
	Topic() (string, error)
	MemberCount() (int, error)
	MemberList() ([]IContact, error)
	Date() (time.Time, error)
	Age() (time.Duration, error)
	ToJson() (string, error)
	RawPayload() (schemas.RoomInvitationPayload, error)
}
