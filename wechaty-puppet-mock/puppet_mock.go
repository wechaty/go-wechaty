package wechaty_puppet_mock

import (
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  option2 "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type PuppetMock struct {
  option *option2.Option
}

func NewPuppetMock(option *option2.Option) *PuppetMock {
  pm := &PuppetMock{option: option}
  return pm
}

func (p *PuppetMock) RoomInvitationPayload(roomInvitationID string) schemas.RoomInvitationPayload {
  panic("implement me")
}

func (p *PuppetMock) RoomInvitationAccept(roomInvitationID string) {
  panic("implement me")
}

func (p *PuppetMock) MessageSendText(conversationID string, text string) string {
  panic("implement me")
}

func (p *PuppetMock) MessageSendContact(conversationID string, contactID string) string {
  panic("implement me")
}

func (p *PuppetMock) MessageSendFile(conversationID string, fileBox file_box.FileBox) string {
  panic("implement me")
}

func (p *PuppetMock) MessageSendUrl(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) string {
  panic("implement me")
}

func (p *PuppetMock) MessageSendMiniProgram(conversationID string, urlLinkPayload *schemas.MiniProgramPayload) string {
  panic("implement me")
}

func (p *PuppetMock) MessageImage(messageID string, imageType schemas.ImageType) file_box.FileBox {
  panic("implement me")
}

func (p *PuppetMock) FriendshipPayload(friendshipID string) schemas.FriendshipPayload {
  panic("implement me")
}

func (p *PuppetMock) FriendshipAccept(friendshipID string) {
  panic("implement me")
}

func (p *PuppetMock) Start() error {
  go func() {
    // emit scan
    p.option.Emit(schemas.PuppetEventNameScan, &schemas.EventScanPayload{
      BaseEventPayload: schemas.BaseEventPayload{},
      Status:           schemas.ScanStatusWaiting,
      QrCode:           "https://not-exist.com",
    })
  }()
  return nil
}
