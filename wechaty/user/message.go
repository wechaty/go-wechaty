package user

import (
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type MessageUserQueryFilter struct {
	From Contact
	Text string // todo:: RegExp
	Room Room
	Type schemas.MessageType
	To   Contact
}

type Message struct {
	wechaty.Accessory

	Type schemas.MessageType

	InvalidDict map[string]bool

	Id      string
	Payload schemas.MessagePayload
}

// Find
//  userQuery: MessageUserQueryFilter | string
func (m *Message) Find(query interface{}) Message {
	var userQuery MessageUserQueryFilter
	if q, ok := query.(string); ok {
		userQuery.Text = q
	} else if q, ok := query.(MessageUserQueryFilter); ok {
		userQuery = q
	}

	messageList := m.FindAll(userQuery)
	if len(messageList) > 0 {
		return messageList[0]
	}
	return Message{}
}

// FindAll
func (m *Message) FindAll(userQuery MessageUserQueryFilter) (messages []Message) {
	puppetQuery := schemas.MessageUserQueryFilter{
		FromId: userQuery.From.Id,
		RoomId: userQuery.Room.Id,
		Text:   userQuery.Text,
		ToId:   userQuery.To.Id,
		Type:   userQuery.Type,
	}

	ids := m.GetPuppet().MessageSearch(puppetQuery)

	// check invalid message
	for _, v := range ids {
		message := m.Load(v)
		if message.Ready() {
			messages = append(messages, message)
		}
	}
	return
}

// Ready todo:: no finish
func (m *Message) Ready() bool {
	if m.IsReady() {
		return true
	}

	m.Payload = m.GetPuppet().MessagePayload(m.Id)

	if len(m.Payload.Id) == 0 {
		// todo:: should not panic, because not recover
		panic("no payload")
	}

	// todo::

	return false
}

func (m *Message) IsReady() bool {
	return len(m.Payload.Id) != 0
}

// Load load message
func (m *Message) Load(id string) Message {
	return Message{Id: id}
}

// Create load message alias
func (m *Message) Create(id string) Message {
	return m.Load(id)
}

// ToString message to print string
// todo:: no finish
func (m *Message) ToString() string {
	return fmt.Sprintf("%v", m.Payload)
}
