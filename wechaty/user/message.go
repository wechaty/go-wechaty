package user

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/config"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Message struct {
	_interface.IAccessory

	id      string
	payload *schemas.MessagePayload
}

// NewMessage ...
func NewMessage(id string, accessory _interface.IAccessory) _interface.IMessage {
	return &Message{
		IAccessory: accessory,
		id:         id,
	}
}

func (m *Message) Ready() (err error) {
	if m.IsReady() {
		return nil
	}

	m.payload, err = m.GetPuppet().MessagePayload(m.id)

	if err != nil {
		return err
	}

	talkerId := m.payload.TalkerId
	roomId := m.payload.RoomId
	listenerId := m.payload.ListenerId

	if roomId != "" {
		if err := m.GetWechaty().Room().Load(roomId).Ready(false); err != nil {
			return err
		}
	}

	if talkerId != "" {
		if err := m.GetWechaty().Contact().Load(talkerId).Ready(false); err != nil {
			return err
		}
	}

	if listenerId != "" {
		if err := m.GetWechaty().Contact().Load(listenerId).Ready(false); err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) IsReady() bool {
	return m.payload != nil
}

// String() message to print string
func (m *Message) String() string {
	talkerStr := ""
	roomStr := ""
	talker := m.Talker()
	if talker != nil {
		talkerStr = "🗣" + talker.String()
	}
	room := m.Room()
	if room != nil {
		roomStr = "@👥'" + room.String()
	}
	str := fmt.Sprintf("Message#%s[%s%s]", m.Type(), talkerStr, roomStr)
	switch m.Type() {
	case schemas.MessageTypeText, schemas.MessageTypeUnknown:
		r := []rune(m.Text())
		if len(r) > 70 {
			r = r[0:70]
		}
		str += "\t" + string(r)
	}
	return str
}

// Room get the room from the message.
func (m *Message) Room() _interface.IRoom {
	roomId := m.payload.RoomId
	if roomId == "" {
		return nil
	}
	return m.GetWechaty().Room().Load(roomId)
}

// Type get the type from the message.
func (m *Message) Type() schemas.MessageType {
	return m.payload.Type
}

// From get the sender from a message
// Deprecated: please use Talker()
func (m *Message) From() _interface.IContact {
	return m.Talker()
}

// Talker Get the talker of a message.
func (m *Message) Talker() _interface.IContact {
	if m.payload.TalkerId == "" {
		return nil
	}
	return m.GetWechaty().Contact().Load(m.payload.TalkerId)
}

// Text get the text content of the message
func (m *Message) Text() string {
	return m.payload.Text
}

// Self check if a message is sent by self
func (m *Message) Self() bool {
	userID := m.GetPuppet().SelfID()
	talker := m.Talker()
	return userID == talker.ID()
}

func (m *Message) Age() time.Duration {
	return time.Since(m.Date())
}

// Date sent date
func (m *Message) Date() time.Time {
	return m.payload.Timestamp
}

// Say reply a Text or Media File message to the sender.
func (m *Message) Say(textOrContactOrFileOrUrlOrMini interface{}) (_interface.IMessage, error) {
	conversationId, err := m.sayId()
	if err != nil {
		return nil, err
	}
	var messageID string
	switch v := textOrContactOrFileOrUrlOrMini.(type) {
	case string:
		messageID, err = m.GetPuppet().MessageSendText(conversationId, v)
	case *Contact:
		messageID, err = m.GetPuppet().MessageSendContact(conversationId, v.Id)
	case *filebox.FileBox:
		messageID, err = m.GetPuppet().MessageSendFile(conversationId, v)
	case *UrlLink:
		messageID, err = m.GetPuppet().MessageSendURL(conversationId, v.payload)
	case *MiniProgram:
		messageID, err = m.GetPuppet().MessageSendMiniProgram(conversationId, v.payload)
	default:
		return nil, fmt.Errorf("unknown msg: %v", textOrContactOrFileOrUrlOrMini)
	}
	if err != nil {
		return nil, err
	}
	if messageID == "" {
		return nil, nil
	}
	message := m.GetWechaty().Message().Load(messageID)
	if err := message.Ready(); err != nil {
		return nil, err
	}
	return message, nil
}

func (m *Message) sayId() (string, error) {
	room := m.Room()
	if room != nil {
		return room.ID(), nil
	}
	talker := m.Talker()
	if talker != nil {
		return talker.ID(), nil
	}
	return "", errors.New("neither room nor talker? ")
}

// To get the destination of the message
// Deprecated: please use Listener()
func (m *Message) To() _interface.IContact {
	return m.Listener()
}

/*
Listener Get the destination of the message.
listener() will return nil if a message is in a room
use Room() to get the room.
*/
func (m *Message) Listener() _interface.IContact {
	if m.payload.ListenerId == "" {
		return nil
	}
	return m.GetWechaty().Contact().Load(m.payload.ListenerId)
}

// ToRecalled get the recalled message
func (m *Message) ToRecalled() (_interface.IMessage, error) {
	if m.Type() != schemas.MessageTypeRecalled {
		return nil, errors.New("can not call toRecalled() on message which is not recalled type")
	}
	originalMessageId := m.Text()
	if originalMessageId == "" {
		return nil, errors.New("can not find recalled message")
	}
	message := m.GetWechaty().Message().Load(originalMessageId)
	err := message.Ready()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Recall recall a message
func (m *Message) Recall() (bool, error) {
	recall, err := m.GetPuppet().MessageRecall(m.id)
	if err != nil {
		return false, err
	}
	return recall, nil
}

// MentionList get message mentioned contactList.
func (m *Message) MentionList() []_interface.IContact {
	room := m.Room()
	if m.Type() != schemas.MessageTypeText || room == nil {
		return nil
	}
	var contactList []_interface.IContact
	if len(m.payload.MentionIdList) > 0 {
		async := helper.NewAsync(helper.DefaultWorkerNum)
		for _, id := range m.payload.MentionIdList {
			id := id
			async.AddTask(func() (interface{}, error) {
				contact := m.GetWechaty().Contact().Load(id)
				return contact, contact.Ready(false)
			})
		}
		for _, v := range async.Result() {
			if v.Err != nil {
				continue
			}
			contactList = append(contactList, v.Value.(_interface.IContact))
		}
		return contactList
	}

	atList := strings.Split(m.Text(), config.AtSepratorRegex)
	if len(atList) == 0 {
		return nil
	}
	var mentionNameList []string
	for _, v := range atList {
		if !strings.Contains(v, "@") {
			continue
		}
		for _, v := range m.multipleAt(v) {
			if v == "" {
				continue
			}
			mentionNameList = append(mentionNameList, v)
		}
	}
	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, name := range mentionNameList {
		name := name
		async.AddTask(func() (interface{}, error) {
			return room.MemberAll(name)
		})
	}
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		contactList = append(contactList, v.Value.([]_interface.IContact)...)
	}
	return contactList
}

// convert 'hello@a@b@c' to [ 'c', 'b@c', 'a@b@c' ]
func (m *Message) multipleAt(str string) []string {
	r, _ := regexp.Compile("^.*?@")
	strs := strings.Split(r.ReplaceAllString(str, "@"), "@")
	var name string
	var nameList []string
	var filterStrs []string
	for _, mentionName := range strs {
		if mentionName == "" {
			continue
		}
		filterStrs = append(filterStrs, mentionName)
	}
	//reverse
	sort.Slice(filterStrs, func(i, j int) bool {
		return filterStrs[i] > filterStrs[j]
	})
	for _, mentionName := range filterStrs {
		name = mentionName + "@" + name
		r := []rune(name)
		nameList = append(nameList, string(r[0:len(r)-1]))
	}
	return nameList
}

func (m *Message) MentionText() string {
	text := m.Text()
	room := m.Room()
	mentionList := m.MentionList()

	if room == nil || len(mentionList) == 0 {
		return text
	}

	toAliasName := func(member _interface.IContact) string {
		alias, _ := room.Alias(member)
		if alias != "" {
			return alias
		}
		return member.Name()
	}

	var mentionNameList []string
	for _, v := range mentionList {
		mentionNameList = append(mentionNameList, toAliasName(v))
	}

	for _, v := range mentionNameList {
		reg := regexp.MustCompile("@" + regexp.QuoteMeta(v) + "(\u2005|\u0020|$)")
		text = reg.ReplaceAllString(text, "")
	}
	return strings.TrimSpace(text)
}

func (m *Message) MentionSelf() bool {
	selfID := m.GetPuppet().SelfID()
	mentionList := m.MentionList()
	for _, v := range mentionList {
		if v.ID() == selfID {
			return true
		}
	}
	return false
}

// Forward Forward the received message.
func (m *Message) Forward(contactOrRoomId string) error {
	_, err := m.GetPuppet().MessageForward(contactOrRoomId, m.id)
	return err
}

// ToFileBox extract the Media File from the Message, and put it into the FileBox.
func (m *Message) ToFileBox() (*filebox.FileBox, error) {
	if m.Type() == schemas.MessageTypeText {
		return nil, errors.New("text message no file")
	}
	return m.GetPuppet().MessageFile(m.id)
}

// ToImage extract the Image File from the Message, so that we can use different image sizes.
func (m *Message) ToImage() (_interface.IImage, error) {
	if m.Type() != schemas.MessageTypeImage {
		return nil, errors.New("not a image type message")
	}
	return m.GetWechaty().Image().Create(m.id), nil
}

// ToContact Get Share Card of the Message
// Extract the Contact Card from the Message, and encapsulate it into Contact class
func (m *Message) ToContact() (_interface.IContact, error) {
	if m.Type() != schemas.MessageTypeContact {
		return nil, errors.New("message not a ShareCard")
	}
	contactID, err := m.GetPuppet().MessageContact(m.id)
	if err != nil {
		return nil, err
	}
	contact := m.GetWechaty().Contact().Load(contactID)
	err = contact.Ready(false)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (m *Message) ToUrlLink() (*UrlLink, error) {
	if m.Type() != schemas.MessageTypeURL {
		return nil, errors.New("message not a Url Link")
	}
	urlPayload, err := m.GetPuppet().MessageURL(m.id)
	if err != nil {
		return nil, err
	}
	return NewUrlLink(urlPayload), nil
}

func (m *Message) ToMiniProgram() (*MiniProgram, error) {
	if m.Type() != schemas.MessageTypeMiniProgram {
		return nil, errors.New("message not a MiniProgram")
	}
	miniProgramPayload, err := m.GetPuppet().MessageMiniProgram(m.id)
	if err != nil {
		return nil, err
	}
	return NewMiniProgram(miniProgramPayload), nil
}

// ID message id
func (m *Message) ID() string {
	return m.id
}
