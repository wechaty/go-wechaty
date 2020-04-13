package wechaty_puppet_mock

import (
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type PuppetMock struct {
  token string
}

func NewPuppetMock(token string) *PuppetMock {
  return &PuppetMock{token: token}
}

func (p *PuppetMock) MessageImage(messageID string, imageType schemas.ImageType) file_box.FileBox {
  panic("implement me")
}

func (p *PuppetMock) FriendshipPayloadReceive(friendshipID string) schemas.FriendshipPayloadReceive {
  panic("implement me")
}

func (p *PuppetMock) FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm {
  panic("implement me")
}

func (p *PuppetMock) FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify {
  panic("implement me")
}

func (p *PuppetMock) FriendshipAccept(friendshipID string) {
  panic("implement me")
}

func (p *PuppetMock) Start(emitChan chan<- schemas.EmitStruct) error {
  go func() {
    // emit scan
    emitChan <- schemas.EmitStruct{
      EventName: schemas.PuppetEventNameScan,
      Payload: schemas.EventScanPayload{
        BaseEventPayload: schemas.BaseEventPayload{},
        Status:           schemas.ScanStatusWaiting,
        QrCode:           "https://not-exist.com",
      },
    }
  }()
  select {}
}
