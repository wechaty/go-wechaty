package wechatypuppet

import (
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/wechaty/go-wechaty/wechaty-puppet/events"
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// iPuppet puppet concrete interface
type iPuppet interface {
	MessageImage(messageID string, imageType schemas.ImageType) (*filebox.FileBox, error)
	Start() (err error)
	Stop()
	Logout() error
	Ding(data string)
	SetContactAlias(contactID string, alias string) error
	ContactAlias(contactID string) (string, error)
	ContactList() ([]string, error)
	ContactQRCode(contactID string) (string, error)
	SetContactAvatar(contactID string, fileBox *filebox.FileBox) error
	ContactAvatar(contactID string) (*filebox.FileBox, error)
	ContactRawPayload(contactID string) (*schemas.ContactPayload, error)
	SetContactSelfName(name string) error
	ContactSelfQRCode() (string, error)
	SetContactSelfSignature(signature string) error
	MessageContact(messageID string) (string, error)
	MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error)
	MessageRecall(messageID string) (bool, error)
	MessageFile(id string) (*filebox.FileBox, error)
	MessageRawPayload(id string) (*schemas.MessagePayload, error)
	MessageSendText(conversationID string, text string, mentionIDList ...string) (string, error)
	MessageSendFile(conversationID string, fileBox *filebox.FileBox) (string, error)
	MessageSendContact(conversationID string, contactID string) (string, error)
	MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error)
	MessageURL(messageID string) (*schemas.UrlLinkPayload, error)
	RoomRawPayload(id string) (*schemas.RoomPayload, error)
	RoomList() ([]string, error)
	RoomDel(roomID, contactID string) error
	RoomAvatar(roomID string) (*filebox.FileBox, error)
	RoomAdd(roomID, contactID string) error
	SetRoomTopic(roomID string, topic string) error
	RoomTopic(roomID string) (string, error)
	RoomCreate(contactIDList []string, topic string) (string, error)
	RoomQuit(roomID string) error
	RoomQRCode(roomID string) (string, error)
	RoomMemberList(roomID string) ([]string, error)
	RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error)
	SetRoomAnnounce(roomID, text string) error
	RoomAnnounce(roomID string) (string, error)
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
	MessageRawMiniProgramPayload(messageID string) (*schemas.MiniProgramPayload, error)
}

// IPuppetAbstract puppet abstract class interface
type IPuppetAbstract interface {
	MessageSearch(query *schemas.MessageQueryFilter) ([]string, error)
	MessagePayload(messageID string) (payload *schemas.MessagePayload, err error)
	FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error)
	SetFriendshipPayload(friendshipID string, newPayload *schemas.FriendshipPayload)
	RoomPayload(roomID string) (payload *schemas.RoomPayload, err error)
	ContactPayload(contactID string) (*schemas.ContactPayload, error)
	ContactSearch(query interface{}, searchIDList []string) ([]string, error)
	FriendshipSearch(query *schemas.FriendshipSearchCondition) (string, error)
	SelfID() string
	iPuppet
	events.EventEmitter
	ContactValidate(contactID string) bool
	RoomValidate(roomID string) bool
	RoomMemberSearch(roomID string, query interface{}) ([]string, error)
	RoomMemberPayload(roomID, memberID string) (*schemas.RoomMemberPayload, error)
	MessageForward(conversationID string, messageID string) (string, error)
	RoomSearch(query *schemas.RoomQueryFilter) ([]string, error)
	RoomInvitationPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error)
	SetRoomInvitationPayload(payload *schemas.RoomInvitationPayload)
	DirtyPayload(payloadType schemas.PayloadType, id string) error
	MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error)
}

// Puppet puppet abstract struct
type Puppet struct {
	Option

	id string
	// puppet implementation puppet_service or puppet_mock
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
func NewPuppet(option Option) (*Puppet, error) {
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

	p := &Puppet{
		Option:                     option,
		EventEmitter:               events.New(),
		cacheMessagePayload:        cacheMessage,
		cacheFriendshipPayload:     cacheFriendship,
		cacheRoomInvitationPayload: cacheRoomInvitation,
		cacheRoomPayload:           cacheRoomPayload,
		cacheRoomMemberPayload:     cacheRoomMemberPayload,
		cacheContactPayload:        cacheContactPayload,
	}

	p.On(schemas.PuppetEventNameDirty, func(i ...interface{}) {
		payload, ok := i[0].(*schemas.EventDirtyPayload)
		if !ok {
			return
		}
		_ = p.OnDirty(payload.PayloadType, payload.PayloadId)
	})
	return p, nil
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
func (p *Puppet) MessageSearch(query *schemas.MessageQueryFilter) ([]string, error) {
	allMessageIDList := p.MessageList()
	if query == nil {
		return allMessageIDList, nil
	}

	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range allMessageIDList {
		id := id
		async.AddTask(func() (interface{}, error) {
			return p.MessagePayload(id)
		})
	}

	var messagePayloadList []*schemas.MessagePayload
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		messagePayloadList = append(messagePayloadList, v.Value.(*schemas.MessagePayload))
	}

	filterFunction := p.messageQueryFilterFactory(query)
	var messageIDList []string
	for _, payload := range messagePayloadList {
		if !filterFunction(payload) {
			continue
		}
		messageIDList = append(messageIDList, payload.Id)
	}

	return messageIDList, nil
}

func (p *Puppet) messageQueryFilterFactory(query *schemas.MessageQueryFilter) schemas.MessagePayloadFilterFunction {
	var filters []schemas.MessagePayloadFilterFunction

	// Deprecated FromId compatible
	//nolint:staticcheck
	if query.FromId != "" && query.TalkerId == "" {
		query.TalkerId = query.FromId
	}
	if query.TalkerId != "" {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.TalkerId == payload.TalkerId
		})
	}
	if query.Id != "" {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.Id == payload.Id
		})
	}
	if query.RoomId != "" {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.RoomId == payload.RoomId
		})
	}
	if query.Text != "" {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.Text == payload.Text
		})
	}
	if query.TextRegExp != nil {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.TextRegExp.MatchString(payload.Text)
		})
	}
	// Deprecated ToId compatible
	//nolint:staticcheck
	if query.ToId != "" && query.ListenerId == "" {
		query.ListenerId = query.ToId
	}
	if query.ListenerId != "" {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.ListenerId == payload.ListenerId
		})
	}
	if query.Type != 0 {
		filters = append(filters, func(payload *schemas.MessagePayload) bool {
			return query.Type == payload.Type
		})
	}
	return func(payload *schemas.MessagePayload) bool {
		for _, v := range filters {
			if !v(payload) {
				return false
			}
		}
		return true
	}
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

	// 对 puppet 实现方返回的消息做统一处理
	NewMsgAdapter(payload.Type).Handle(payload)

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

// SetFriendshipPayload ...
func (p *Puppet) SetFriendshipPayload(friendshipID string, newPayload *schemas.FriendshipPayload) {
	p.cacheFriendshipPayload.Add(friendshipID, newPayload)
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

// ContactSearch query params is string or *schemas.ContactQueryFilter
func (p *Puppet) ContactSearch(query interface{}, searchIDList []string) ([]string, error) {
	if searchIDList == nil {
		var err error
		searchIDList, err = p.puppetImplementation.ContactList()
		if err != nil || len(searchIDList) == 0 {
			return nil, err
		}
	}

	if query == nil {
		return searchIDList, nil
	}

	switch v := query.(type) {
	case string:
		return p.contactSearchByQueryString(v, searchIDList)
	case *schemas.ContactQueryFilter:
		return p.contactSearchByQueryFilter(query.(*schemas.ContactQueryFilter), searchIDList)
	default:
		return nil, errors.New("unsupported query types")
	}
}

func (p *Puppet) contactSearchByQueryString(query string, searchIDList []string) ([]string, error) {
	nameIDList, err := p.contactSearchByQueryFilter(&schemas.ContactQueryFilter{Name: query}, searchIDList)
	if err != nil {
		return nil, err
	}
	aliasIDList, err := p.contactSearchByQueryFilter(&schemas.ContactQueryFilter{Alias: query}, searchIDList)
	if err != nil {
		return nil, err
	}
	return append(nameIDList, aliasIDList...), nil
}

func (p *Puppet) contactQueryFilterFactory(query *schemas.ContactQueryFilter) (schemas.ContactPayloadFilterFunction, error) {
	if query.Alias != "" {
		return func(payload *schemas.ContactPayload) bool {
			return payload.Alias == query.Alias
		}, nil
	}
	if query.AliasRegexp != nil {
		return func(payload *schemas.ContactPayload) bool {
			return query.AliasRegexp.MatchString(payload.Alias)
		}, nil
	}
	if query.Id != "" {
		return func(payload *schemas.ContactPayload) bool {
			return payload.Id == query.Id
		}, nil
	}
	if query.Name != "" {
		return func(payload *schemas.ContactPayload) bool {
			return payload.Name == query.Name
		}, nil
	}
	if query.NameRegexp != nil {
		return func(payload *schemas.ContactPayload) bool {
			return query.NameRegexp.MatchString(payload.Name)
		}, nil
	}
	if query.WeiXin != "" {
		return func(payload *schemas.ContactPayload) bool {
			return payload.WeiXin == query.WeiXin
		}, nil
	}
	return nil, errors.New("query must provide at least one key. current query is empty. ")
}

func (p *Puppet) contactSearchByQueryFilter(query *schemas.ContactQueryFilter, searchIDList []string) ([]string, error) {
	filterFun, err := p.contactQueryFilterFactory(query)
	if err != nil {
		return nil, err
	}
	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range searchIDList {
		id := id
		async.AddTask(func() (interface{}, error) {
			payload, err := p.ContactPayload(id)
			if err != nil {
				p.dirtyPayloadContact(id)
			}
			return payload, err
		})
	}
	var contactIDs []string
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		payload := v.Value.(*schemas.ContactPayload)
		if filterFun(payload) {
			contactIDs = append(contactIDs, payload.Id)
		}
	}
	return contactIDs, nil
}

// ContactValidate ...
func (p *Puppet) ContactValidate(contactID string) bool {
	return true
}

// FriendshipSearch ...
func (p *Puppet) FriendshipSearch(query *schemas.FriendshipSearchCondition) (string, error) {
	if query.Phone != "" {
		return p.puppetImplementation.FriendshipSearchPhone(query.Phone)
	} else if query.WeiXin != "" {
		return p.puppetImplementation.FriendshipSearchWeixin(query.WeiXin)
	} else {
		return "", errors.New("query must provide at least one key. current query is empty. ")
	}
}

// RoomMemberSearch ...
func (p *Puppet) RoomMemberSearch(roomID string, query interface{}) ([]string, error) {
	switch v := query.(type) {
	case string:
		return p.roomMemberSearchByString(roomID, v)
	case *schemas.RoomMemberQueryFilter:
		return p.roomMemberSearchByQueryFilter(roomID, v)
	default:
		return nil, errors.New("unsupported query types")
	}
}

func (p *Puppet) roomMemberSearchByString(roomID string, query string) ([]string, error) {
	roomAliasIDList, err := p.puppetImplementation.RoomMemberSearch(roomID, &schemas.RoomMemberQueryFilter{
		RoomAlias: query,
	})
	if err != nil {
		return nil, err
	}
	nameIDList, err := p.puppetImplementation.RoomMemberSearch(roomID, &schemas.RoomMemberQueryFilter{
		Name: query,
	})
	if err != nil {
		return nil, err
	}
	contactAliasIDList, err := p.puppetImplementation.RoomMemberSearch(roomID, &schemas.RoomMemberQueryFilter{
		ContactAlias: query,
	})
	if err != nil {
		return nil, err
	}
	return append(roomAliasIDList, append(nameIDList, contactAliasIDList...)...), nil
}

func (p *Puppet) roomMemberSearchByQueryFilter(roomID string, query *schemas.RoomMemberQueryFilter) ([]string, error) {
	memberIDList, err := p.puppetImplementation.RoomMemberList(roomID)
	if err != nil {
		return nil, err
	}
	var idList []string
	if query.ContactAlias != "" || query.Name != "" {
		contactQueryFilter := &schemas.ContactQueryFilter{
			Alias: query.ContactAlias,
			Name:  query.Name,
		}
		contactIDList, err := p.puppetImplementation.ContactSearch(contactQueryFilter, memberIDList)
		if err != nil {
			return nil, err
		}
		idList = append(idList, contactIDList...)
	}
	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range memberIDList {
		id := id
		async.AddTask(func() (interface{}, error) {
			return p.RoomMemberPayload(roomID, id)
		})
	}
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		if query.RoomAlias == "" {
			continue
		}
		payload := v.Value.(*schemas.RoomMemberPayload)
		if payload.RoomAlias == query.RoomAlias {
			idList = append(idList, payload.Id)
		}
	}
	return idList, nil
}

// RoomMemberPayload ...
func (p *Puppet) RoomMemberPayload(roomID, memberID string) (*schemas.RoomMemberPayload, error) {
	cacheKey := p.cacheKeyRoomMember(roomID, memberID)
	cachePayload, ok := p.cacheRoomMemberPayload.Get(cacheKey)
	if ok {
		return cachePayload.(*schemas.RoomMemberPayload), nil
	}
	payload, err := p.puppetImplementation.RoomMemberRawPayload(roomID, memberID)
	if err != nil {
		return nil, err
	}
	p.cacheRoomMemberPayload.Add(cacheKey, payload)
	return payload, nil
}

// RoomInvitationPayload ...
func (p *Puppet) RoomInvitationPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error) {
	cachePayload, ok := p.cacheRoomInvitationPayload.Get(roomInvitationID)
	if ok {
		return cachePayload.(*schemas.RoomInvitationPayload), nil
	}
	payload, err := p.puppetImplementation.RoomInvitationRawPayload(roomInvitationID)
	if err != nil {
		return nil, err
	}
	p.cacheRoomInvitationPayload.Add(roomInvitationID, payload)
	return payload, nil
}

// SetRoomInvitationPayload ...
func (p *Puppet) SetRoomInvitationPayload(payload *schemas.RoomInvitationPayload) {
	p.cacheRoomInvitationPayload.Add(payload.Id, payload)
}

// MessageForward ...
func (p *Puppet) MessageForward(conversationID string, messageID string) (string, error) {
	return p.puppetImplementation.MessageForward(conversationID, messageID)
	//payload, err := p.MessagePayload(messageID)
	//if err != nil {
	//	return "", err
	//}
	//var newMsgID string
	//switch payload.Type {
	//case schemas.MessageTypeAttachment,
	//	schemas.MessageTypeVideo,
	//	schemas.MessageTypeAudio,
	//	schemas.MessageTypeImage:
	//	newMsgID, err = p.messageForwardFile(conversationID, messageID)
	//case schemas.MessageTypeText:
	//	newMsgID, err = p.puppetImplementation.MessageSendText(conversationID, payload.Text)
	//case schemas.MessageTypeMiniProgram:
	//	newMsgID, err = p.messageForwardMiniProgram(conversationID, messageID)
	//case schemas.MessageTypeURL:
	//	newMsgID, err = p.messageForwardURL(conversationID, messageID)
	//case schemas.MessageTypeContact:
	//	newMsgID, err = p.messageForwardContact(conversationID, messageID)
	//default:
	//	return "", fmt.Errorf("unsupported forward message type: %s", payload.Type)
	//}
	//if err != nil {
	//	return "", err
	//}
	//return newMsgID, nil
}

func (p *Puppet) messageForwardFile(conversationID string, messageID string) (string, error) { //nolint:unused
	file, err := p.puppetImplementation.MessageFile(messageID)
	if err != nil {
		return "", err
	}
	newMsgID, err := p.puppetImplementation.MessageSendFile(conversationID, file)
	if err != nil {
		return "", err
	}
	return newMsgID, nil
}

func (p *Puppet) messageForwardMiniProgram(conversationID string, messageID string) (string, error) { //nolint:unused
	payload, err := p.puppetImplementation.MessageMiniProgram(messageID)
	if err != nil {
		return "", err
	}
	newMsgID, err := p.puppetImplementation.MessageSendMiniProgram(conversationID, payload)
	if err != nil {
		return "", err
	}
	return newMsgID, nil
}

func (p *Puppet) messageForwardURL(conversationID string, messageID string) (string, error) { //nolint:unused
	payload, err := p.puppetImplementation.MessageURL(messageID)
	if err != nil {
		return "", err
	}
	newMsgID, err := p.puppetImplementation.MessageSendURL(conversationID, payload)
	if err != nil {
		return "", err
	}
	return newMsgID, nil
}

func (p *Puppet) messageForwardContact(conversationID string, messageID string) (string, error) { //nolint:unused
	payload, err := p.puppetImplementation.MessageContact(messageID)
	if err != nil {
		return "", err
	}
	newMsgID, err := p.puppetImplementation.MessageSendContact(conversationID, payload)
	if err != nil {
		return "", err
	}
	return newMsgID, nil
}

// RoomSearch ...
func (p *Puppet) RoomSearch(query *schemas.RoomQueryFilter) ([]string, error) {
	allRoomIDList, err := p.puppetImplementation.RoomList()
	if err != nil {
		return nil, err
	}
	if query == nil || query.Empty() {
		return allRoomIDList, nil
	}
	filterFunc, err := p.roomQueryFilterFactory(query)
	if err != nil {
		return nil, err
	}
	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range allRoomIDList {
		id := id
		async.AddTask(func() (interface{}, error) {
			payload, err := p.RoomPayload(id)
			if err != nil {
				p.dirtyPayloadRoom(id)
				p.dirtyPayloadRoomMember(id)
				return nil, err
			}
			return payload, nil
		})
	}
	var roomIDList []string
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		payload := v.Value.(*schemas.RoomPayload)
		if !filterFunc(payload) {
			continue
		}
		roomIDList = append(roomIDList, payload.Id)
	}
	return roomIDList, nil
}

func (p *Puppet) roomQueryFilterFactory(query *schemas.RoomQueryFilter) (schemas.RoomPayloadFilterFunction, error) {
	if query.Empty() {
		return nil, errors.New("query must provide at least one key. current query is empty")
	} else if query.All() {
		return nil, errors.New("query only support one key. multi key support is not availble now")
	}
	if query.TopicRegexp != nil {
		return func(payload *schemas.RoomPayload) bool {
			return query.TopicRegexp.MatchString(payload.Topic)
		}, nil
	}
	if query.Id != "" {
		return func(payload *schemas.RoomPayload) bool {
			return query.Id == payload.Id
		}, nil
	}
	if query.Topic != "" {
		return func(payload *schemas.RoomPayload) bool {
			return query.Topic == payload.Topic
		}, nil
	}
	return nil, nil
}

// RoomValidate ...
func (p *Puppet) RoomValidate(roomID string) bool {
	return true
}

func (p *Puppet) dirtyPayloadMessage(messageID string) {
	p.cacheMessagePayload.Remove(messageID)
}

func (p *Puppet) dirtyPayloadContact(contactID string) {
	p.cacheContactPayload.Remove(contactID)
}

func (p *Puppet) dirtyPayloadRoom(roomID string) {
	p.cacheRoomPayload.Remove(roomID)
}

func (p *Puppet) dirtyPayloadRoomMember(roomID string) {
	contactIds, _ := p.puppetImplementation.RoomMemberList(roomID)
	for _, id := range contactIds {
		p.cacheRoomMemberPayload.Remove(p.cacheKeyRoomMember(roomID, id))
	}
}

func (p *Puppet) dirtyPayloadFriendship(friendshipID string) {
	p.cacheFriendshipPayload.Remove(friendshipID)
}

// OnDirty clean cache
func (p *Puppet) OnDirty(payloadType schemas.PayloadType, id string) error {
	switch payloadType {
	case schemas.PayloadTypeMessage:
		p.dirtyPayloadMessage(id)
	case schemas.PayloadTypeContact:
		p.dirtyPayloadContact(id)
	case schemas.PayloadTypeRoom:
		p.dirtyPayloadRoom(id)
	case schemas.PayloadTypeRoomMember:
		p.dirtyPayloadRoomMember(id)
	case schemas.PayloadTypeFriendship:
		p.dirtyPayloadFriendship(id)
	default:
		return fmt.Errorf("unknown payload type: %v", payloadType)
	}
	return nil
}

// MessageMiniProgram ...
func (p *Puppet) MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error) {
	get, ok := p.cacheMessagePayload.Get(messageID)
	if !ok {
		return p.puppetImplementation.MessageRawMiniProgramPayload(messageID)
	}
	payload := get.(*schemas.MessagePayload)
	if !payload.FixMiniApp {
		return p.puppetImplementation.MessageRawMiniProgramPayload(messageID)
	}
	miniapp, err := helper.ParseMiniApp(payload)
	if err != nil {
		return nil, err
	}
	return miniapp, nil
}

// DirtyPayload base clean cache
func (p *Puppet) DirtyPayload(payloadType schemas.PayloadType, id string) error {
	return p.OnDirty(payloadType, id)
}
