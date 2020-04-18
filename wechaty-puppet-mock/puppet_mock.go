package wechaty_puppet_mock

import (
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  option2 "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type PuppetMock struct {
  option *option2.Option
}

func NewPuppetMock(optFns ...option2.OptionFn) *PuppetMock {
  pm := &PuppetMock{option: &option2.Option{}}

  for _, fn := range optFns {
    fn(pm.option)
  }
  return pm
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

func (p *PuppetMock) Start() error {
  go func() {
    // emit scan
    p.option.EmitChan <- schemas.EmitStruct{
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
