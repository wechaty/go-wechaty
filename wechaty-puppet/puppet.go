package wechatypuppet

import (
  lru "github.com/hashicorp/golang-lru"
  "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"

  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// puppetInterface puppet interface
type puppetInterface interface {
  MessageImage(messageID string, imageType schemas.ImageType) file_box.FileBox
  FriendshipPayloadReceive(friendshipID string) schemas.FriendshipPayloadReceive
  FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm
  FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify
  FriendshipAccept(friendshipID string)
}

// Puppet puppet struct
type Puppet struct {
  puppetInterface

  CacheMessagePayload    *lru.Cache
  CacheFriendshipPayload *lru.Cache
}

// MessageList message list
func (p *Puppet) MessageList() (ks []string) {
  keys := p.CacheMessagePayload.Keys()
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
  cachePayload, ok := p.CacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadReceive)
  }
  payload := p.puppetInterface.FriendshipPayloadReceive(friendshipID)
  p.CacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}

// FriendshipPayloadConfirm ...
func (p *Puppet) FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm {
  cachePayload, ok := p.CacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadConfirm)
  }
  payload := p.puppetInterface.FriendshipPayloadConfirm(friendshipID)
  p.CacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}

// FriendshipPayloadVerify ...
func (p *Puppet) FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify {
  cachePayload, ok := p.CacheFriendshipPayload.Get(friendshipID)
  if ok {
    return cachePayload.(schemas.FriendshipPayloadVerify)
  }
  payload := p.puppetInterface.FriendshipPayloadVerify(friendshipID)
  p.CacheFriendshipPayload.Add(friendshipID, payload)
  return payload
}
