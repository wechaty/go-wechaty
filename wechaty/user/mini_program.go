package user

import "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

type MiniProgram struct {
	payload *schemas.MiniProgramPayload
}

func NewMiniProgram(payload *schemas.MiniProgramPayload) *MiniProgram {
	return &MiniProgram{payload: payload}
}

func (mp *MiniProgram) Appid() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Appid
}

func (mp *MiniProgram) Description() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Description
}

func (mp *MiniProgram) PagePath() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.PagePath
}

func (mp *MiniProgram) ThumbUrl() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.ThumbUrl
}

func (mp *MiniProgram) Title() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Title
}

func (mp *MiniProgram) Username() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.Username
}

func (mp *MiniProgram) ThumbKey() string {
	if mp.payloadIsNil() {
		return ""
	}
	return mp.payload.ThumbKey
}

func (mp *MiniProgram) payloadIsNil() bool {
	return mp.payload == nil
}
