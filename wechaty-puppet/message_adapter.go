package wechatypuppet

import (
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"regexp"
)

var numRegex = regexp.MustCompile(`^\d+$`)

var rawMsgAdapter = RawMsgAdapter{}
var unknownMsgAdapter = UnknownMsgAdapter{}
var recalledMsgAdapter = RecalledMsgAdapter{}

// MsgAdapter 消息适配器
type MsgAdapter interface {
	Handle(payload *schemas.MessagePayload)
}

// NewMsgAdapter 各种 puppet 返回的消息有出入，这里做统一
func NewMsgAdapter(msgType schemas.MessageType) MsgAdapter {
	switch msgType {
	case schemas.MessageTypeUnknown:
		return unknownMsgAdapter
	case schemas.MessageTypeRecalled:
		return recalledMsgAdapter
	}
	return rawMsgAdapter
}

// RawMsgAdapter 不需要处理的消息
type RawMsgAdapter struct{}

// Handle ~
func (r RawMsgAdapter) Handle(msg *schemas.MessagePayload) {}

// UnknownMsgAdapter Unknown 类型的消息适配器
type UnknownMsgAdapter struct{}

// Handle 对 Unknown 类型的消息做适配
func (u UnknownMsgAdapter) Handle(payload *schemas.MessagePayload) {
	// 有些消息，puppet 服务端没有解析出来，这里尝试解析
	helper.FixUnknownMessage(payload)
}

// RecalledMsgAdapter 撤回类型的消息适配器
type RecalledMsgAdapter struct{}

// Handle padlocal 返回的是 xml，需要解析出 msgId
func (r RecalledMsgAdapter) Handle(payload *schemas.MessagePayload) {
	if numRegex.MatchString(payload.Text) {
		return
	}
	// padlocal 返回的是 xml，需要解析出 msgId
	// https://github.com/wechaty/go-wechaty/issues/87
	payload.Text = helper.ParseRecalledID(payload.Text)
}
