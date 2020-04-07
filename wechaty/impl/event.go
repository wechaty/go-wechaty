package impl

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/user"
  "time"
)

type (
  EventDong       func(data string)
  EventError      func(err error)
  EventFriendship func(friendship string) // TODO(dchaofei): Friendship struct
  EventHeartbeat  func(data string)
  EventLogin      func(user string)                // TODO(dchaofei): ContactSelf struct
  EventLogout     func(user string, reason string) // TODO(dchaofei): ContactSelf struct
  EventMessage    func(message user.Message)
  EventReady      func()
  EventRoomInvite func(roomInvitation string) // TODO(dchaofei): RoomInvitation struct
  EventRoomJoin   func(room user.Room, inviteeList []user.Contact, inviter user.Contact, date time.Time)
  EventRoomLeave  func(room user.Room, leaverList []user.Contact, remover user.Contact, date time.Time)
  EventRoomTopic  func(room user.Room, newTopic string, oldTopic string, changer user.Contact, date time.Time)
  EventScan       func(qrCode string, status schemas.ScanStatus, data string)
  EventStart      func()
  EventStop       func()
)
