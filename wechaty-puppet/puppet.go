package wechatypuppet

import (
  lru "github.com/hashicorp/golang-lru"
  wph "github.com/wechaty/go-wechaty/wechaty-puppet-hostie"
  "github.com/wechaty/go-wechaty/wechaty-puppet/events"
  "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// Option Puppet option
type Option struct {
  token string
}

// PuppetInterface puppet interface
type PuppetInterface interface {
  MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error)
  FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error)
  FriendshipAccept(friendshipID string) error
  Start() error
  RoomInvitationPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error)
  RoomInvitationAccept(roomInvitationID string) error
  MessageSendText(conversationID string, text string) (string, error)
  MessageSendContact(conversationID string, contactID string) (string, error)
  MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error)
  MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error)
  MessageSendMiniProgram(conversationID string, urlLinkPayload *schemas.MiniProgramPayload) (string, error)
}

// Puppet puppet struct
type Puppet struct {
  *option.Option
  PuppetInterface

  cacheMessagePayload        *lru.Cache
  cacheFriendshipPayload     *lru.Cache
  cacheRoomInvitationPayload *lru.Cache
}

// NewPuppet instance
func NewPuppet(option *option.Option) (*Puppet, error) {
  cacheMessage, err := lru.New(1024)
  if err != nil {
    return nil, err
  }
  cacheFriendship, err := lru.New(1024)
  if err != nil {
    return nil, err
  }
  cacheRoomInvitation, err := lru.New(1024)
  if err != nil {
    return nil, err
  }
  if option.EventEmitter == nil {
    option.EventEmitter = events.New()
  }
  return &Puppet{
    Option:                     option,
    PuppetInterface:            wph.NewPuppetHostie(option),
    cacheMessagePayload:        cacheMessage,
    cacheFriendshipPayload:     cacheFriendship,
    cacheRoomInvitationPayload: cacheRoomInvitation,
  }, nil
}

// MessageList message list
func (p *Puppet) MessageList() (ks []string) {
  keys := p.cacheMessagePayload.Keys()
  for _, v := range keys {
    if k, ok := v.(string); ok {
      ks = append(ks, k)
    }
  }
  return
}

// MessageSearch search message
func (p *Puppet) MessageSearch(query schemas.MessageUserQueryFilter) []string {
  allMessageIDList := p.MessageList()
  if len(allMessageIDList) <= 0 {
    return allMessageIDList
  }

  // load
  var messagePayloadList []schemas.MessagePayload
  for _, v := range allMessageIDList {
    messagePayloadList = append(messagePayloadList, p.MessagePayload(v))
  }
  // Filter todo:: messageQueryFilterFactory
  var messageIDList []string
  for _, message := range messagePayloadList {
    if message.FromId == query.FromId || message.RoomId == query.RoomId || message.ToId == query.ToId {
      messageIDList = append(messageIDList, message.Id)
    }
  }

  return messageIDList
}

// messageQueryFilterFactory 实现正则和直接匹配
func (p *Puppet) messageQueryFilterFactory(query string) schemas.MessagePayloadFilterFunction {
  return nil
}

// MessagePayload message payload todo:: no finish
func (p *Puppet) MessagePayload(messageID string) (payload schemas.MessagePayload) {
  return payload
}

// FriendshipPayload ...
func (p *Puppet) FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error) {
  cachePayload, ok := p.cacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(*schemas.FriendshipPayload), nil
  }
  payload, err := p.PuppetInterface.FriendshipPayload(friendshipID)
  if err != nil {
    return nil, err
  }
  p.cacheFriendshipPayload.Add(friendshipID, payload)
  return payload, nil
}
