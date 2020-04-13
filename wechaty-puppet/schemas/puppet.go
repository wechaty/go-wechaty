package schemas

type PuppetOptions struct {
  endpoint string
  timeout  int64
  token    string
}

//go:generate stringer -type=PuppetEventName
type PuppetEventName uint8

const (
  PuppetEventNameUnknown PuppetEventName = iota
  PuppetEventNameFriendShipConfirm
  PuppetEventNameFriendShipReceive
  PuppetEventNameFriendShipVerify
  PuppetEventNameLogin
  PuppetEventNameLogout
  PuppetEventNameMessage
  PuppetEventNameRoomInvite
  PuppetEventNameRoomJoin
  PuppetEventNameRoomLeave
  PuppetEventNameRoomTopic
  PuppetEventNameScan

  PuppetEventNameDong
  PuppetEventNameError
  PuppetEventNameHeartbeat
  PuppetEventNameReady
  PuppetEventNameReset

  PuppetEventNameStop
  PuppetEventNameStart
)
