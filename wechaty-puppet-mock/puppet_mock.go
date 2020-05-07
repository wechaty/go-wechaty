package wechaty_puppet_mock

import (
  wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

var _ wechatyPuppet.IPuppetAbstract = &PuppetMock{}

type PuppetMock struct {
  *wechatyPuppet.Puppet
}

func NewPuppetMock(option *wechatyPuppet.Option) (*PuppetMock, error) {
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

func (p PuppetMock) RoomInvitationRawPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error) {
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

func (p *PuppetMock) Stop() {
  panic("implement me")
}

func (p *PuppetMock) Logout() error {
  panic("implement me")
}

func (p *PuppetMock) Ding(data string) {
  panic("implement me")
}

func (p *PuppetMock) SetContactAlias(contactID string, alias string) error {
  panic("implement me")
}

func (p *PuppetMock) ContactAlias(contactID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) ContactList() ([]string, error) {
  panic("implement me")
}

func (p *PuppetMock) ContactQRCode(contactID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) SetContactAvatar(contactID string, fileBox *file_box.FileBox) error {
  panic("implement me")
}

func (p *PuppetMock) GetContactAvatar(contactID string) (*file_box.FileBox, error) {
  panic("implement me")
}

func (p *PuppetMock) ContactRawPayload(contactID string) (*schemas.ContactPayload, error) {
  panic("implement me")
}

func (p *PuppetMock) SetContactSelfName(name string) error {
  panic("implement me")
}

func (p *PuppetMock) ContactSelfQRCode() (string, error) {
  panic("implement me")
}

func (p *PuppetMock) SetContactSelfSignature(signature string) error {
  panic("implement me")
}

func (p *PuppetMock) MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error) {
  panic("implement me")
}

func (p *PuppetMock) MessageContact(messageID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) MessageRecall(messageID string) (bool, error) {
  panic("implement me")
}

func (p *PuppetMock) MessageFile(id string) (*file_box.FileBox, error) {
  panic("implement me")
}

func (p *PuppetMock) MessageRawPayload(id string) (*schemas.MessagePayload, error) {
  panic("implement me")
}

func (p *PuppetMock) MessageURL(messageID string) (*schemas.UrlLinkPayload, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomRawPayload(id string) (*schemas.RoomPayload, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomList() ([]string, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomDel(roomID, contactID string) error {
  panic("implement me")
}

func (p *PuppetMock) RoomAvatar(roomID string) (*file_box.FileBox, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomAdd(roomID, contactID string) error {
  panic("implement me")
}

func (p *PuppetMock) SetRoomTopic(roomID string, topic string) error {
  panic("implement me")
}

func (p *PuppetMock) GetRoomTopic(roomID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomCreate(contactIDList []string, topic string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomQuit(roomID string) error {
  panic("implement me")
}

func (p *PuppetMock) RoomQRCode(roomID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomMemberList(roomID string) ([]string, error) {
  panic("implement me")
}

func (p *PuppetMock) RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error) {
  panic("implement me")
}

func (p *PuppetMock) SetRoomAnnounce(roomID, text string) error {
  panic("implement me")
}

func (p *PuppetMock) GetRoomAnnounce(roomID string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) FriendshipSearchPhone(phone string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) FriendshipSearchWeixin(weixin string) (string, error) {
  panic("implement me")
}

func (p *PuppetMock) FriendshipAdd(contactID, hello string) (err error) {
  panic("implement me")
}

func (p *PuppetMock) TagContactAdd(id, contactID string) (err error) {
  panic("implement me")
}

func (p *PuppetMock) TagContactRemove(id, contactID string) (err error) {
  panic("implement me")
}

func (p *PuppetMock) TagContactDelete(id string) (err error) {
  panic("implement me")
}

func (p *PuppetMock) TagContactList(contactID string) ([]string, error) {
  panic("implement me")
}
