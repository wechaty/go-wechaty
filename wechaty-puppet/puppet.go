package wechatypuppet

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	"github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// iPuppet puppet concrete interface
type iPuppet interface {
	MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error)
	Start() (err error)
	Stop()
	Logout() error
	Ding(data string)
	SetContactAlias(contactID string, alias string) error
	GetContactAlias(contactID string) (string, error)
	ContactList() ([]string, error)
	ContactQRCode(contactID string) (string, error)
	SetContactAvatar(contactID string, fileBox *file_box.FileBox) error
	GetContactAvatar(contactID string) (*file_box.FileBox, error)
	ContactRawPayload(contactID string) (*schemas.ContactPayload, error)
	SetContactSelfName(name string) error
	ContactSelfQRCode() (string, error)
	SetContactSelfSignature(signature string) error
	MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error)
	MessageContact(messageID string) (string, error)
	MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error)
	MessageRecall(messageID string) (bool, error)
	MessageFile(id string) (*file_box.FileBox, error)
	MessageRawPayload(id string) (*schemas.MessagePayload, error)
	MessageSendText(conversationID string, text string) (string, error)
	MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error)
	MessageSendContact(conversationID string, contactID string) (string, error)
	MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error)
	MessageURL(messageID string) (*schemas.UrlLinkPayload, error)
	RoomRawPayload(id string) (*schemas.RoomPayload, error)
	RoomList() ([]string, error)
	RoomDel(roomID, contactID string) error
	RoomAvatar(roomID string) (*file_box.FileBox, error)
	RoomAdd(roomID, contactID string) error
	SetRoomTopic(roomID string, topic string) error
	GetRoomTopic(roomID string) (string, error)
	RoomCreate(contactIDList []string, topic string) (string, error)
	RoomQuit(roomID string) error
	RoomQRCode(roomID string) (string, error)
	RoomMemberList(roomID string) ([]string, error)
	RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error)
	SetRoomAnnounce(roomID, text string) error
	GetRoomAnnounce(roomID string) (string, error)
	RoomInvitationAccept(roomInvitationID string) error
	RoomInvitationRawPayload(id string) (*schemas.RoomInvitationPayload, error)
	FriendshipSearchPhone(phone string) (string, error)
	FriendshipSearchWeixin(weixin string) (string, error)
	FriendshipRawPayload(id string) (*schemas.FriendshipPayload, error)
	FriendshipAdd(contactID, hello string) (err error)
	FriendshipAccept(friendshipID string) (err error)
	TagContactAdd(id, contactID string) (err error)
	TagContactRemove(id, contactID string) (err error)
	TagContactDelete(id string) (err error)
	TagContactList(contactID string) ([]string, error)
}

// IPuppetAbstract puppet abstract class interface
type IPuppetAbstract interface {
	MessageSearch(query schemas.MessageUserQueryFilter) ([]string, error)
	MessagePayload(messageID string) (payload *schemas.MessagePayload, err error)
	FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error)
	RoomPayloadDirty(roomID string)
	RoomMemberPayloadDirty(roomID string) error
	RoomPayload(roomID string) (payload *schemas.RoomPayload, err error)
	ContactPayloadDirty(contactID string)
	ContactPayload(contactID string) (*schemas.ContactPayload, error)
	SelfID() string
	iPuppet
	events.EventEmitter
}

// Puppet puppet abstract struct
type Puppet struct {
	*Option

	id string
	// puppet implementation puppet_hostie or puppet_mock
	events.EventEmitter
	puppetImplementation       IPuppetAbstract
	cacheMessagePayload        *lru.Cache
	cacheFriendshipPayload     *lru.Cache
	cacheRoomInvitationPayload *lru.Cache
	cacheRoomPayload           *lru.Cache
	cacheRoomMemberPayload     *lru.Cache
	cacheContactPayload        *lru.Cache
}

// NewPuppet instance
func NewPuppet(option *Option) (*Puppet, error) {
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
	cacheRoomPayload, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	cacheRoomMemberPayload, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	cacheContactPayload, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	return &Puppet{
		Option:                     option,
		EventEmitter:               events.New(),
		cacheMessagePayload:        cacheMessage,
		cacheFriendshipPayload:     cacheFriendship,
		cacheRoomInvitationPayload: cacheRoomInvitation,
		cacheRoomPayload:           cacheRoomPayload,
		cacheRoomMemberPayload:     cacheRoomMemberPayload,
		cacheContactPayload:        cacheContactPayload,
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
func (p *Puppet) MessageSearch(query schemas.MessageUserQueryFilter) ([]string, error) {
	allMessageIDList := p.MessageList()
	if len(allMessageIDList) <= 0 {
		return allMessageIDList, nil
	}

	// load
	var messagePayloadList []*schemas.MessagePayload
	for _, v := range allMessageIDList {
		payload, err := p.MessagePayload(v)
		if err != nil {
			return nil, err
		}
		messagePayloadList = append(messagePayloadList, payload)
	}
	// Filter todo:: messageQueryFilterFactory
	var messageIDList []string
	for _, message := range messagePayloadList {
		if message.FromId == query.FromId || message.RoomId == query.RoomId || message.ToId == query.ToId {
			messageIDList = append(messageIDList, message.Id)
		}
	}

	return messageIDList, nil
}

// messageQueryFilterFactory 实现正则和直接匹配
func (p *Puppet) messageQueryFilterFactory(query string) schemas.MessagePayloadFilterFunction {
	return nil
}

// MessagePayload message payload todo:: no finish
func (p *Puppet) MessagePayload(messageID string) (*schemas.MessagePayload, error) {
	cachePayload, ok := p.cacheMessagePayload.Get(messageID)
	if ok {
		return cachePayload.(*schemas.MessagePayload), nil
	}
	payload, err := p.puppetImplementation.MessageRawPayload(messageID)
	if err != nil {
		return nil, err
	}
	p.cacheMessagePayload.Add(messageID, payload)
	return payload, nil
}

// FriendshipPayload ...
func (p *Puppet) FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error) {
	cachePayload, ok := p.cacheFriendshipPayload.Get(friendshipID)
	if ok {
		return cachePayload.(*schemas.FriendshipPayload), nil
	}
	payload, err := p.puppetImplementation.FriendshipRawPayload(friendshipID)
	if err != nil {
		return nil, err
	}
	p.cacheFriendshipPayload.Add(friendshipID, payload)
	return payload, nil
}

// SetPuppetImplementation set puppet implementation
func (p *Puppet) SetPuppetImplementation(i IPuppetAbstract) {
	p.puppetImplementation = i
}

// SetID set login id
func (p *Puppet) SetID(id string) {
	p.id = id
}

// SelfID self id
func (p *Puppet) SelfID() string {
	return p.id
}

// RoomPayloadDirty ...
func (p *Puppet) RoomPayloadDirty(roomID string) {
	p.cacheRoomPayload.Remove(roomID)
}

// RoomMemberPayloadDirty ...
func (p *Puppet) RoomMemberPayloadDirty(roomID string) error {
	contactIds, err := p.puppetImplementation.RoomMemberList(roomID)
	if err != nil {
		return err
	}
	for _, id := range contactIds {
		p.cacheRoomMemberPayload.Remove(p.cacheKeyRoomMember(roomID, id))
	}
	return nil
}

func (p *Puppet) cacheKeyRoomMember(roomID string, contactID string) string {
	return contactID + "@@@" + roomID
}

// RoomPayload ...
func (p *Puppet) RoomPayload(roomID string) (payload *schemas.RoomPayload, err error) {
	cachePayload, ok := p.cacheRoomPayload.Get(roomID)
	if ok {
		return cachePayload.(*schemas.RoomPayload), nil
	}
	payload, err = p.puppetImplementation.RoomRawPayload(roomID)
	if err != nil {
		return nil, err
	}
	p.cacheRoomPayload.Add(roomID, payload)
	return payload, nil
}

// ContactPayloadDirty ...
func (p *Puppet) ContactPayloadDirty(contactID string) {
	p.cacheContactPayload.Remove(contactID)
}

// ContactPayload ...
func (p *Puppet) ContactPayload(contactID string) (*schemas.ContactPayload, error) {
	cachePayload, ok := p.cacheContactPayload.Get(contactID)
	if ok {
		return cachePayload.(*schemas.ContactPayload), nil
	}
	payload, err := p.puppetImplementation.ContactRawPayload(contactID)
	if err != nil {
		return nil, err
	}
	p.cacheContactPayload.Add(contactID, payload)
	return payload, nil
}
