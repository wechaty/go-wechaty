package schemas

type PuppetOptions struct {
  endpoint string
  timeout  int64
  token    string
}

type PuppetEventName int

const (
  PuppetEventNameUnknown PuppetEventName = iota
  PuppetEventNameFriendShip
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
