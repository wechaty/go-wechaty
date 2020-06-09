package factory

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
)

type MessageFactory struct {
	_interface.IAccessory
}

func (m *MessageFactory) Load(id string) _interface.IMessage {
	return user.NewMessage(id, m.IAccessory)
}

// Find find message in cache
func (m *MessageFactory) Find(query interface{}) _interface.IMessage {
	var q *schemas.MessageQueryFilter
	switch v := query.(type) {
	case string:
		q = &schemas.MessageQueryFilter{Text: v}
	case *schemas.MessageQueryFilter:
		q = v
	default:
		log.Println("not support query type")
		return nil
	}
	messages := m.FindAll(q)
	if len(messages) < 1 {
		return nil
	}
	if len(messages) > 1 {
		log.Printf("Message FindAll() got more than one(%d) result\n", len(messages))
	}
	return messages[0]
}

// FindAll Find message in cache
func (m *MessageFactory) FindAll(query *schemas.MessageQueryFilter) []_interface.IMessage {
	var err error
	defer func() {
		if err != nil {
			log.Printf("MessageFactory FindAll rejected: %s\n", err)
		}
	}()
	messageIDs, err := m.GetPuppet().MessageSearch(query)
	if err != nil {
		return nil
	}

	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range messageIDs {
		id := id
		async.AddTask(func() (interface{}, error) {
			message := m.Load(id)
			return message, message.Ready()
		})
	}

	var messages []_interface.IMessage
	for _, v := range async.Result() {
		if v.Err != nil {
			continue
		}
		messages = append(messages, v.Value.(_interface.IMessage))
	}
	return messages
}
