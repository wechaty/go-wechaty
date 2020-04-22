package wechaty_puppet_mock

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  option2 "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

var _ wechatyPuppet.IPuppetAbstract = &PuppetMock{}

type PuppetMock struct {
  *wechatyPuppet.Puppet
}

func NewPuppetMock(option *option2.Option) (*PuppetMock, error) {
  puppetAbstract, err := wechatyPuppet.NewPuppet(option)
  if err != nil {
    return nil, err
  }
  puppetMock := &PuppetMock{
    Puppet: puppetAbstract,
  }
  puppetAbstract.SetPuppetImplementation(puppetMock)
  return puppetMock, nil
}

func (p *PuppetMock) Start() error {
  go func() {
    // emit scan
    p.Emit(schemas.PuppetEventNameScan, &schemas.EventScanPayload{
      BaseEventPayload: schemas.BaseEventPayload{},
      Status:           schemas.ScanStatusWaiting,
      QrCode:           "https://not-exist.com",
    })
  }()
  return nil
}

func (p PuppetMock) MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error) {
  panic("implement me")
}

func (p PuppetMock) FriendshipRawPayload(friendshipID string) (*schemas.FriendshipPayload, error) {
  panic("implement me")
}

func (p PuppetMock) FriendshipAccept(friendshipID string) error {
  panic("implement me")
}

func (p PuppetMock) RoomInvitationPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error) {
  panic("implement me")
}

func (p PuppetMock) RoomInvitationAccept(roomInvitationID string) error {
  panic("implement me")
}

func (p PuppetMock) MessageSendText(conversationID string, text string) (string, error) {
  panic("implement me")
}

func (p PuppetMock) MessageSendContact(conversationID string, contactID string) (string, error) {
  panic("implement me")
}

func (p PuppetMock) MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error) {
  panic("implement me")
}

func (p PuppetMock) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
  panic("implement me")
}

func (p PuppetMock) MessageSendMiniProgram(conversationID string, urlLinkPayload *schemas.MiniProgramPayload) (string, error) {
  panic("implement me")
}
