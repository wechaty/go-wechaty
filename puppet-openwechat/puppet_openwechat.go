package puppet_openwechat

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/pkg/errors"
	wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	logger "github.com/wechaty/go-wechaty/wechaty-puppet/log"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

var log = logger.L.WithField("module", "puppet-openwechat")

type PuppetOpenWechat struct {
	*wechatyPuppet.Puppet

	bot  *openwechat.Bot
	self *openwechat.Self
}

func NewPuppetOpenWechat() (*PuppetOpenWechat, error) {
	puppet := &PuppetOpenWechat{}
	puppetBase, err := wechatyPuppet.NewPuppet(wechatyPuppet.Option{})
	puppetBase.SetPuppetImplementation(puppet)
	if err != nil {
		return nil, err
	}
	bot := openwechat.NewBot()
	bot.Caller.Client.SetMode(openwechat.Desktop)
	puppet.bot = bot
	puppet.Puppet = puppetBase

	puppet.initCallback()
	return puppet, nil
}

func (p PuppetOpenWechat) initCallback() {
	p.onScan()
	p.onMsg()
}

func (p PuppetOpenWechat) MessageSearch(query *schemas.MessageQueryFilter) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipPayload(friendshipID string) (*schemas.FriendshipPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetFriendshipPayload(friendshipID string, newPayload *schemas.FriendshipPayload) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactPayload(contactID string) (*schemas.ContactPayload, error) {
	log.Tracef("ContactPayload(%s)", contactID)
	self, err := p.bot.GetCurrentUser()
	if err != nil {
		return nil, errors.Wrap(err, "ContactPayload p.bot.GetCurrentUser")
	}
	if self.User.UserName == contactID {
		return openWechatUserToContact(self.User), nil
	}

	friends, err := self.Friends()
	if err != nil {
		return nil, errors.Wrap(err, "ContactPayload self.Friends")
	}

	friend := friends.GetByUsername(contactID)
	if friend == nil {
		return nil, fmt.Errorf("not found contactID=%s", contactID)
	}
	return openWechatUserToContact(friend.User), nil
}

func (p PuppetOpenWechat) ContactSearch(query interface{}, searchIDList []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipSearch(query *schemas.FriendshipSearchCondition) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageImage(messageID string, imageType schemas.ImageType) (*filebox.FileBox, error) {
	panic("implement me")
}

func (p *PuppetOpenWechat) Start() (err error) {
	log.Trace("PuppetOpenWechat.Start")

	//if err := p.bot.HotLogin(openwechat.NewJsonFileHotReloadStorage("puppet-open-wechat.memory-card.json"), true); err != nil {
	//	return errors.Wrap(err, "PuppetOpenWechat.Start HotLogin")
	//}
	if err := p.bot.HotLogin(openwechat.NewJsonFileHotReloadStorage("/Users/dingchaofei/work/github/go-wechaty/puppet-open-wechat.memory-card.json"), true); err != nil {
		return errors.Wrap(err, "PuppetOpenWechat.Start HotLogin")
	}
	self, err := p.bot.GetCurrentUser()
	if err != nil {
		return errors.Wrap(err, "PuppetOpenWechat.Start GetCurrentUser")
	}
	p.self = self
	go p.Emit(schemas.PuppetEventNameLogin, &schemas.EventLoginPayload{
		ContactId: self.UserName,
	})
	return err
}

func (p PuppetOpenWechat) Stop() {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) Logout() error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) Ding(data string) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetContactAlias(contactID string, alias string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactAlias(contactID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactList() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactQRCode(contactID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetContactAvatar(contactID string, fileBox *filebox.FileBox) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactAvatar(contactID string) (*filebox.FileBox, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactRawPayload(contactID string) (*schemas.ContactPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetContactSelfName(name string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactSelfQRCode() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetContactSelfSignature(signature string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageContact(messageID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageRecall(messageID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageFile(id string) (*filebox.FileBox, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageRawPayload(id string) (*schemas.MessagePayload, error) {
	return nil, fmt.Errorf("PuppetOpenWechat not implement MessageRawPayload method")
}

func (p PuppetOpenWechat) MessageSendText(conversationID string, text string, mentionIDList ...string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageSendFile(conversationID string, fileBox *filebox.FileBox) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageSendContact(conversationID string, contactID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageURL(messageID string) (*schemas.UrlLinkPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomRawPayload(id string) (*schemas.RoomPayload, error) {
	groups, err := p.self.Groups()
	if err != nil {
		return nil, err
	}
	rawRoom := groups.GetByUsername(id)
	if rawRoom == nil {
		return nil, fmt.Errorf("PuppetOpenWechat RoomRawPayload not found room id=[%s]", id)
	}

	payload := &schemas.RoomPayload{
		Id:           rawRoom.UserName,
		Topic:        rawRoom.NickName,
		Avatar:       "", // TODO 群头像？
		MemberIdList: nil,
		OwnerId:      "",
		AdminIdList:  nil,
	}

	for _, v := range rawRoom.MemberList {
		payload.MemberIdList = append(payload.MemberIdList, v.UserName)
	}
	return payload, nil
}

func (p PuppetOpenWechat) RoomList() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomDel(roomID, contactID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomAvatar(roomID string) (*filebox.FileBox, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomAdd(roomID, contactID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetRoomTopic(roomID string, topic string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomTopic(roomID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomCreate(contactIDList []string, topic string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomQuit(roomID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomQRCode(roomID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomMemberList(roomID string) ([]string, error) {
	log.Tracef("RoomMemberList(%s)", roomID)
	group, err := p.self.Groups()
	if err != nil {
		return nil, err
	}
	room := group.GetByUsername(roomID)
	if room == nil {
		return nil, fmt.Errorf("PuppetOpenWechat.RoomMemberList not found room id=[%s]", roomID)
	}
	memberIds := make([]string, 0, len(group))
	for _, v := range room.MemberList {
		memberIds = append(memberIds, v.UserName)
	}
	return memberIds, nil
}

func (p PuppetOpenWechat) RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetRoomAnnounce(roomID, text string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomAnnounce(roomID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomInvitationAccept(roomInvitationID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomInvitationRawPayload(id string) (*schemas.RoomInvitationPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipSearchPhone(phone string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipSearchWeixin(weixin string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipRawPayload(id string) (*schemas.FriendshipPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipAdd(contactID, hello string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) FriendshipAccept(friendshipID string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) TagContactAdd(id, contactID string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) TagContactRemove(id, contactID string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) TagContactDelete(id string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) TagContactList(contactID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageRawMiniProgramPayload(messageID string) (*schemas.MiniProgramPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) ContactValidate(contactID string) bool {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomValidate(roomID string) bool {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomMemberSearch(roomID string, query interface{}) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomMemberPayload(roomID, memberID string) (*schemas.RoomMemberPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageForward(conversationID string, messageID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomSearch(query *schemas.RoomQueryFilter) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) RoomInvitationPayload(roomInvitationID string) (*schemas.RoomInvitationPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) SetRoomInvitationPayload(payload *schemas.RoomInvitationPayload) {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) DirtyPayload(payloadType schemas.PayloadType, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p PuppetOpenWechat) MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error) {
	//TODO implement me
	panic("implement me")
}
