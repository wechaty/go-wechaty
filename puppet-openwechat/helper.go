package puppet_openwechat

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"os"
	"strings"
	"time"
)

func openWechatUserToContact(user *openwechat.User) *schemas.ContactPayload {
	return &schemas.ContactPayload{
		Id:        user.UserName,
		Gender:    schemas.ContactGender(user.Sex),
		Type:      contactType(user),
		Name:      user.NickName,
		Avatar:    "",
		Address:   "",
		Alias:     user.Alias,
		City:      user.City,
		Friend:    user.IsFriend(),
		Province:  user.Province,
		Signature: user.Signature,
		Star:      user.StarFriend != 0,
		WeiXin:    "",
	}
}

func contactType(user *openwechat.User) schemas.ContactType {
	if user.IsFriend() {
		return schemas.ContactTypePersonal
	}
	if user.IsMP() {
		return schemas.ContactTypeOfficial
	}
	return schemas.ContactTypeUnknown
}

func rawMsgToPayload(rawMsg *openwechat.Message) (*schemas.MessagePayload, error) {
	payload := &schemas.MessagePayload{
		MessagePayloadBase: schemas.MessagePayloadBase{
			Id:            rawMsg.MsgId,
			MentionIdList: nil, // TODO 从消息中解析被 @ 列表
			FileName:      "",  // TODO 暂时不确定有什么用
			Text:          rawMsg.Content,
			Timestamp:     time.Unix(rawMsg.CreateTime, 0),
			Type:          webMsgToPuppetMsgType(rawMsg),
		},
	}

	if rawMsg.IsSendByGroup() {
		fmt.Println(rawMsg.FromUserName)
		s, _ := rawMsg.Bot.GetCurrentUser()
		g, _ := s.Groups(true)

		room := g.GetByUsername(rawMsg.FromUserName)
		if room == nil {
			f, err := os.OpenFile("temp.out", os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, os.ModePerm)
			if err != nil {
				panic(err)
			}
			_, err = f.Write(rawMsg.Raw)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(rawMsg.Raw) + "\n\n")
		}
		fmt.Println(g.GetByUsername(rawMsg.FromUserName).NickName)
		if isRoomId(rawMsg.FromUserName) {
			payload.RoomId = rawMsg.FromUserName
		} else if isRoomId(rawMsg.ToUserName) {
			payload.RoomId = rawMsg.ToUserName
		} else {
			return nil, fmt.Errorf("解析群消息，但是没有找到 roomId, raw: %s", string(rawMsg.Raw))
		}

		payload.TalkerId = rawMsg.SenderInGroupUserName
	} else {
		payload.TalkerId = rawMsg.FromUserName
	}

	if !isRoomId(rawMsg.ToUserName) {
		payload.ListenerId = rawMsg.ToUserName
	}
	if payload.ListenerId == "" && payload.RoomId == "" {
		return nil, fmt.Errorf("ListenerId 和 RoomId 都为空， raw:%s", string(rawMsg.Raw))
	}

	return payload, nil
}

func isRoomId(s string) bool {
	return strings.HasPrefix(s, "@@")
}

func webMsgToPuppetMsgType(rawMsg *openwechat.Message) schemas.MessageType {
	switch rawMsg.MsgType {
	case openwechat.MsgTypeText:
		switch rawMsg.SubMsgType {
		case int(openwechat.MsgTypeLocation):
			return schemas.MessageTypeAttachment
		default:
			return schemas.MessageTypeText
		}

	case openwechat.MsgTypeEmoticon, openwechat.MsgTypeImage:
		return schemas.MessageTypeImage

	case openwechat.MsgTypeVoice:
		return schemas.MessageTypeAudio

	case openwechat.MsgTypeMicroVideo, openwechat.MsgTypeVideo:
		return schemas.MessageTypeVideo

	case openwechat.MsgTypeApp:
		switch rawMsg.AppMsgType {
		case openwechat.AppMsgTypeAttach, openwechat.AppMsgTypeUrl, openwechat.AppMsgTypeReaderType:
			return schemas.MessageTypeAttachment
		default:
			return schemas.MessageTypeText
		}

	case openwechat.MsgTypeSys:
		return schemas.MessageTypeText

	case openwechat.MsgTypeRecalled:
		return schemas.MessageTypeRecalled

	default:
		return schemas.MessageTypeText
	}
}

//
//func rawUserToRoom(rawUser *openwechat.User) (*schemas.RoomPayload, error) {
//	rawUser.Get
//}
