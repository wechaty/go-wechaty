package wechatypuppet

import (
  "fmt"
  lru "github.com/hashicorp/golang-lru"
  wpm "github.com/wechaty/go-wechaty/wechaty-puppet-mock"
  "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "log"
)

// Option Puppet option
type Option struct {
  token string
}

// PuppetInterface puppet interface
type PuppetInterface interface {
  MessageImage(messageID string, imageType schemas.ImageType) file_box.FileBox
  FriendshipPayloadReceive(friendshipID string) schemas.FriendshipPayloadReceive
  FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm
  FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify
  FriendshipAccept(friendshipID string)
  Start(emitChan chan<- schemas.EmitStruct) error
  RoomInvitationPayload(roomInvitationID string) schemas.RoomInvitationPayload
  RoomInvitationAccept(roomInvitationID string)
  MessageSendText(conversationID string, text string) string
  MessageSendContact(conversationID string, contactID string) string
  MessageSendFile(conversationID string, fileBox file_box.FileBox) string
  MessageSendUrl(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) string
  MessageSendMiniProgram(conversationID string, urlLinkPayload *schemas.MiniProgramPayload) string
}

// Puppet puppet struct
type Puppet struct {
  PuppetInterface

  cacheMessagePayload        *lru.Cache
  cacheFriendshipPayload     *lru.Cache
  cacheRoomInvitationPayload *lru.Cache
  eventParamsChan            chan<- schemas.EventParams
}

// NewPuppet instance
func NewPuppet(eventParamsChan chan<- schemas.EventParams, token string) (*Puppet, error) {
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
  return &Puppet{
    PuppetInterface:        wpm.NewPuppetMock(token),
    cacheMessagePayload:    cacheMessage,
    cacheFriendshipPayload: cacheFriendship,
    cacheRoomInvitationPayload: cacheRoomInvitation,
    eventParamsChan:        eventParamsChan,
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

// FriendshipPayloadReceive ...
func (p *Puppet) FriendshipPayloadReceive(friendshipID string) schemas.FriendshipPayloadReceive {
  cachePayload, ok := p.cacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadReceive)
  }
  payload := p.PuppetInterface.FriendshipPayloadReceive(friendshipID)
  p.cacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}

// FriendshipPayloadConfirm ...
func (p *Puppet) FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm {
  cachePayload, ok := p.cacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadConfirm)
  }
  payload := p.PuppetInterface.FriendshipPayloadConfirm(friendshipID)
  p.cacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}

// FriendshipPayloadVerify ...
func (p *Puppet) FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify {
  cachePayload, ok := p.cacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadVerify)
  }
  payload := p.PuppetInterface.FriendshipPayloadVerify(friendshipID)
  p.cacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}

// RoomInvitationPayload ...
func (p *Puppet) RoomInvitationPayload(roomInvitationID string) schemas.RoomInvitationPayload {
  cachePayload, ok := p.cacheRoomInvitationPayload.Get(roomInvitationID)
  if ok {
    return cachePayload.(schemas.RoomInvitationPayload)
  }
  payload := p.PuppetInterface.RoomInvitationPayload(roomInvitationID)
  p.cacheRoomInvitationPayload.Add(roomInvitationID, payload)
  return payload
}

// Start puppet
func (p *Puppet) Start() error {
  emitChan := make(chan schemas.EmitStruct)
  errChan := make(chan error)
  go func() {
    errChan <- p.PuppetInterface.Start(emitChan)
  }()
  go func() {
    for v := range emitChan {
      if err := p.emit(v); err != nil {
        // TODO log
        log.Printf("Puppet.Start emit err: %s\n", err)
      }
    }
  }()
  return <-errChan
}

func (p *Puppet) emitScan(payload schemas.EventScanPayload) {
  p.eventParamsChan <- schemas.EventParams{
    EventName: schemas.PuppetEventNameScan,
    Params: []interface{}{
      payload.QrCode,
      payload.Status,
      payload.Data,
    },
  }
}

func (p *Puppet) emit(emitStruct schemas.EmitStruct) error {
  switch emitStruct.EventName {
  case schemas.PuppetEventNameScan:
    p.emitScan(emitStruct.Payload.(schemas.EventScanPayload))
  default:
    return fmt.Errorf("not support envent, %v", emitStruct.EventName)
  }
  return nil
}
