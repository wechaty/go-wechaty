package wechaty_puppet_padplus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"path/filepath"

	uuid "github.com/satori/go.uuid"

	"github.com/wechaty/go-wechaty/wechaty-puppet-padplus/payload"
	pd "github.com/wechaty/go-wechaty/wechaty-puppet-padplus/proto"
	file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

// Request 发送请求
// 有些事件应该在一定时间后得到响应，否则则超市
// SEND_MESSAGE,SEND_FILE 3 * 60 * 1000
// GET_MESSAGE_MEDIA 5 * 60 * 1000
// SEARCH_CONTACT,ADD_CONTACT,CREATE_ROOM,GET_ROOM_QRCODE 1 * 60 * 1000
func (p *PuppetPadPlus) Request(apiType pd.ApiType, data interface{}) (response string, err error) {
	us := uuid.NewV4().String()
	trace := uuid.NewV4().String()
	paramsByte, err := json.Marshal(data)
	params := string(paramsByte)

	resp, err := p.grpcClient.Request(context.Background(), &pd.RequestObject{
		RequestId: &us, // uuid
		Token:     &p.Token,
		ApiType:   &apiType,
		TraceId:   &trace, // uuid
		Params:    &params,
		Uin:       &p.Uin,
	})
	if err != nil {
		log.Printf("Type: %s, err: %v", apiType, err)
		return "", err
	}
	log.Printf("Type: %s, Result: %s", apiType, *resp.Result)
	return *resp.Result, err
}

// Login 登录
func (p *PuppetPadPlus) Login() (err error) {
	_, err = p.Request(pd.ApiType_GET_QRCODE, nil)
	return
}

// AutoLogin 自动登录
func (p *PuppetPadPlus) AutoLogin() (err error) {
	_, err = p.Request(pd.ApiType_INIT, nil)
	return
}

// MessageImage load image message
func (p *PuppetPadPlus) MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error) {
	log.Printf("PuppetPadPlus MessageImage(%s, %s)\n", messageID, imageType.String())

	rawPayLoad, err := p.MessageRawPayload(messageID)
	if err != nil {
		return nil, fmt.Errorf("load cache message: %w", err)
	}

	switch imageType {
	case schemas.ImageTypeThumbnail:
		return file_box.FromUrl("uri", "", nil)
	case schemas.ImageTypeArtwork:
		data := payload.PadPlusRichMediaData{
			Content:      "",
			MsgType:      0, // "rawPayLoad.Type",
			ContentType:  "img",
			Src:          "uri",
			AppMsgType:   0,
			FileName:     "",
			MsgId:        messageID,
			CreateTime:   0,
			FromUserName: rawPayLoad.FromId,
			ToUserName:   rawPayLoad.ToId,
		}

		pay, _ := p.loadRichMediaData(data)
		if len(pay.Src) > 0 {
			return file_box.FromUrl(pay.Src, pay.Content, nil)
		}
	}
	return nil, errors.New("message image type error")
}

func (p *PuppetPadPlus) MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
	panic("implement me")
}

// MessageSendFile send file message
func (p *PuppetPadPlus) MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error) {
	var data map[string]interface{}

	mimeType := mime.TypeByExtension(filepath.Ext(fileBox.Name))
	switch mimeType {
	case "image/jpeg", "image/png", ".jpg", ".jpeg", ".png":
		// PadplusMessageType.Image == 3
		j, _ := fileBox.ToJSON()
		var udata map[string]string
		_ = json.Unmarshal([]byte(j), &udata)
		data = map[string]interface{}{"fileName": fileBox.Name, "fromUserName": p.SelfID(), "messageType": 3, "subType": "pic", "toUserName": conversationID, "url": udata["remoteUrl"]}
	}

	res, err := p.Request(pd.ApiType_SEND_FILE, data)
	if err != nil {
		return "", fmt.Errorf("PuppetPadPlus MessageSendText err: %w", err)
	}
	var pay payload.SendMessageResponse
	p.unMarshal(res, &pay)

	return pay.MsgId, nil
}

func (p *PuppetPadPlus) MessageSendContact(conversationID string, contactID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomInvitationAccept(roomInvitationID string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) FriendshipAccept(friendshipID string) (err error) {
	panic("implement me")
}

// Logout logout
func (p *PuppetPadPlus) Logout() (err error) {
	log.Println("PuppetHostie Logout()")
	if !p.isLogin() {
		return errors.New("logout before login? ")
	}
	_, err = p.Request(pd.ApiType_LOGOUT, p.SelfID())
	if err != nil {
		return fmt.Errorf("PuppetHostie Logout() err: %w", err)
	}
	go p.Emit(schemas.PuppetEventNameLogout, &schemas.EventLogoutPayload{
		ContactId: p.SelfID(),
	})
	p.SetID("")
	return
}

func (p *PuppetPadPlus) Ding(data string) {
	panic("implement me")
}

// SetContactAlias set contact alias
func (p *PuppetPadPlus) SetContactAlias(contactID string, alias string) error {
	_, err := p.Request(pd.ApiType_CONTACT_ALIAS, map[string]string{"newRemarkName": alias, "userName": contactID, "wechatUser": p.SelfID()})
	if err != nil {
		return fmt.Errorf("PuppetHostie SetContactAlias err: %w", err)
	}
	return nil
}

func (p *PuppetPadPlus) ContactAlias(contactID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) ContactList() ([]string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) ContactQRCode(contactID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) SetContactAvatar(contactID string, fileBox *file_box.FileBox) error {
	panic("implement me")
}

func (p *PuppetPadPlus) ContactAvatar(contactID string) (*file_box.FileBox, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) ContactRawPayload(contactID string) (*schemas.ContactPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) SetContactSelfName(name string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) ContactSelfQRCode() (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) SetContactSelfSignature(signature string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) MessageContact(messageID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) MessageRecall(messageID string) (bool, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) MessageFile(id string) (*file_box.FileBox, error) {
	panic("implement me")
}

// MessageRawPayload load message at cache
func (p *PuppetPadPlus) MessageRawPayload(id string) (*schemas.MessagePayload, error) {
	pay, err := p.messagePayload.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if p, ok := pay.(schemas.MessagePayload); ok {
		return &p, nil
	}
	return nil, errors.New("message to payload error")
}

// MessageSendText send text message
func (p *PuppetPadPlus) MessageSendText(conversationID string, text string, mentionIDList ...string) (string, error) {
	log.Printf("PuppetPadPlus messageSendText(%s, %s)\n", conversationID, text)

	res, err := p.Request(pd.ApiType_SEND_MESSAGE, map[string]interface{}{"content": text, "fromUserName": p.SelfID(), "messageType": "text", "toUserName": conversationID, "mentionListStr": mentionIDList})
	if err != nil {
		return "", fmt.Errorf("PuppetPadPlus MessageSendText err: %w", err)
	}
	var pay payload.SendMessageResponse
	p.unMarshal(res, &pay)

	return pay.MsgId, nil
}

func (p *PuppetPadPlus) MessageURL(messageID string) (*schemas.UrlLinkPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomRawPayload(id string) (*schemas.RoomPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomList() ([]string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomDel(roomID, contactID string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomAvatar(roomID string) (*file_box.FileBox, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomAdd(roomID, contactID string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) SetRoomTopic(roomID string, topic string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomTopic(roomID string) (string, error) {
	panic("implement me")
}

// RoomCreate create room
func (p *PuppetPadPlus) RoomCreate(contactIDList []string, topic string) (string, error) {
	log.Printf("PuppetHostie roomCreate(%s, %s)\n", contactIDList, topic)

	res, err := p.Request(pd.ApiType_CREATE_ROOM, map[string]interface{}{"memberList": contactIDList, "topic": topic})
	if err != nil {
		return "", fmt.Errorf("PuppetHostie SetContactAlias err: %w", err)
	}
	var pay payload.CreateRoomResponse
	p.unMarshal(res, &pay)
	return pay.RoomId, nil
}

func (p *PuppetPadPlus) RoomQuit(roomID string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomQRCode(roomID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomMemberList(roomID string) ([]string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) SetRoomAnnounce(roomID, text string) error {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomAnnounce(roomID string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) RoomInvitationRawPayload(id string) (*schemas.RoomInvitationPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) FriendshipSearchPhone(phone string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) FriendshipSearchWeixin(weixin string) (string, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) FriendshipRawPayload(id string) (*schemas.FriendshipPayload, error) {
	panic("implement me")
}

func (p *PuppetPadPlus) FriendshipAdd(contactID, hello string) (err error) {
	panic("implement me")
}

func (p *PuppetPadPlus) TagContactAdd(id, contactID string) (err error) {
	panic("implement me")
}

func (p *PuppetPadPlus) TagContactRemove(id, contactID string) (err error) {
	panic("implement me")
}

func (p *PuppetPadPlus) TagContactDelete(id string) (err error) {
	panic("implement me")
}

func (p *PuppetPadPlus) TagContactList(contactID string) ([]string, error) {
	panic("implement me")
}
