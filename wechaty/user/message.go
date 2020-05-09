package user

import (
	"errors"
	"fmt"
	file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/interface"
	"time"
)

//type MessageQueryFilter struct {
//	From Contact
//	Text string // todo:: RegExp
//	Room Room
//	Type schemas.MessageType
//	To   Contact
//}

type Message struct {
	_interface.Accessory

	id      string
	payload *schemas.MessagePayload
}

func NewMessage(id string, accessory _interface.Accessory) _interface.IMessage {
	return &Message{
		Accessory: accessory,
		id:        id,
	}
}

//// Find
////  userQuery: MessageQueryFilter | string
//func (m *Message) Find(query interface{}) Message {
//	var userQuery MessageQueryFilter
//	if q, ok := query.(string); ok {
//		userQuery.Text = q
//	} else if q, ok := query.(MessageQueryFilter); ok {
//		userQuery = q
//	}
//
//	messageList := m.FindAll(userQuery)
//	if len(messageList) > 0 {
//		return messageList[0]
//	}
//	return Message{}
//}
//
//// FindAll
//func (m *Message) FindAll(userQuery MessageQueryFilter) (messages []Message) {
//	puppetQuery := schemas.MessageQueryFilter{
//		FromId: userQuery.From.id,
//		RoomId: userQuery.Room.id,
//		Text:   userQuery.Text,
//		ToId:   userQuery.To.id,
//		Type:   userQuery.Type,
//	}
//
//	ids := m.GetPuppet().MessageSearch(puppetQuery)
//
//	// check invalid message
//	for _, v := range ids {
//		message := m.Load(v)
//		if message.Ready() {
//			messages = append(messages, message)
//		}
//	}
//	return
//}

// Ready todo:: no finish
func (m *Message) Ready() (err error) {
	if m.IsReady() {
		return nil
	}

	m.payload, err = m.GetPuppet().MessagePayload(m.id)

	if err != nil {
		return err
	}

	fromId := m.payload.FromId
	roomId := m.payload.RoomId
	toId := m.payload.ToId

	if roomId != "" {
		if err := m.GetWechaty().Room().Load(roomId).Ready(false); err != nil {
			return err
		}
	}

	if fromId != "" {
		if err := m.GetWechaty().Contact().Load(fromId).Ready(false); err != nil {
			return err
		}
	}

	if toId != "" {
		if err := m.GetWechaty().Contact().Load(toId).Ready(false); err != nil {
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
	fromStr := ""
	roomStr := ""
	from := m.From()
	if from != nil {
		fromStr = "🗣" + from.String()
	}
	room := m.Room()
	if room != nil {
		roomStr = "@👥'" + room.String()
	}
	str := fmt.Sprintf("Message#%s[%s%s]", m.Type(), fromStr, roomStr)
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

func (m *Message) Room() _interface.IRoom {
	roomId := m.payload.RoomId
	if roomId == "" {
		return nil
	}
	return m.GetWechaty().Room().Load(roomId)
}

func (m *Message) Type() schemas.MessageType {
	return m.payload.Type
}

func (m *Message) From() _interface.IContact {
	if m.payload.FromId == "" {
		return nil
	}
	return m.GetWechaty().Contact().Load(m.payload.FromId)
}

func (m *Message) Text() string {
	return m.payload.Text
}

func (m *Message) Self() bool {
	userID := m.GetPuppet().SelfID()
	from := m.From()
	return userID == from.ID()
}

func (m *Message) Age() time.Duration {
	return time.Now().Sub(m.Date())
}

func (m *Message) Date() time.Time {
	return time.Unix(int64(m.payload.Timestamp), 0)
}

// SayText ... TODO unified say()
func (m *Message) SayText(text string) (_interface.IMessage, error) {
	conversationId, err := m.getSayId()
	if err != nil {
		return nil, err
	}
	messageID, err := m.GetPuppet().MessageSendText(conversationId, text)
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

// SayFile ...
func (m *Message) SayFile(fileBox *file_box.FileBox) (_interface.IMessage, error) {
	conversationId, err := m.getSayId()
	if err != nil {
		return nil, err
	}
	messageID, err := m.GetPuppet().MessageSendFile(conversationId, fileBox)
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

func (m *Message) getSayId() (string, error) {
	room := m.Room()
	if room != nil {
		return room.ID(), nil
	}
	from := m.From()
	if from != nil {
		return from.ID(), nil
	}
	return "", errors.New("neither room nor from? ")
}
