package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/config"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
)

type Room struct {
	id      string
	payLoad *schemas.RoomPayload
	_interface.IAccessory
}

// NewRoom ...
func NewRoom(id string, accessory _interface.IAccessory) *Room {
	return &Room{
		id:         id,
		IAccessory: accessory,
	}
}

// Ready is For FrameWork ONLY!
func (r *Room) Ready(forceSync bool) (err error) {
	if !forceSync && r.IsReady() {
		return nil
	}

	r.payLoad, err = r.GetPuppet().RoomPayload(r.id)
	if err != nil {
		return err
	}

	memberIDs, err := r.GetPuppet().RoomMemberList(r.id)
	if err != nil {
		return err
	}

	async := helper.NewAsync(helper.DefaultWorkerNum)
	for _, id := range memberIDs {
		id := id
		async.AddTask(func() (interface{}, error) {
			return nil, r.GetWechaty().Contact().Load(id).Ready(false)
		})
	}
	_ = async.Result()

	return nil
}

func (r *Room) IsReady() bool {
	return r.payLoad != nil
}

func (r *Room) String() string {
	str := "loading"
	if r.payLoad.Topic != "" {
		str = r.payLoad.Topic
	}
	return fmt.Sprintf("Room<%s>", str)
}

func (r *Room) ID() string {
	return r.id
}

// Find all contacts in a room
// params nil or string or RoomMemberQueryFilter
func (r *Room) MemberAll(query interface{}) ([]_interface.IContact, error) {
	if query == nil {
		return r.memberList()
	}
	idList, err := r.GetPuppet().RoomMemberSearch(r.id, query)
	if err != nil {
		return nil, err
	}
	var contactList []_interface.IContact
	for _, id := range idList {
		contact := r.GetWechaty().Contact().Load(id)
		if err := contact.Ready(false); err != nil {
			return nil, err
		}
		contactList = append(contactList, contact)
	}
	return contactList, nil
}

// Member Find all contacts in a room, if get many, return the first one.
// query params string or RoomMemberQueryFilter
func (r *Room) Member(query interface{}) (_interface.IContact, error) {
	memberList, err := r.MemberAll(query)
	if err != nil {
		return nil, err
	}
	if len(memberList) == 0 {
		return nil, nil
	}
	return memberList[0], nil
}

// get all room member from the room
func (r *Room) memberList() ([]_interface.IContact, error) {
	memberIDList, err := r.GetPuppet().RoomMemberList(r.id)
	if err != nil {
		return nil, err
	}
	if len(memberIDList) == 0 {
		return nil, nil
	}
	var contactList []_interface.IContact
	for _, id := range memberIDList {
		contactList = append(contactList, r.GetWechaty().Contact().Load(id))
	}
	return contactList, nil
}

// Alias return contact's roomAlias in the room
func (r *Room) Alias(contact _interface.IContact) (string, error) {
	memberPayload, err := r.GetPuppet().RoomMemberPayload(r.id, contact.ID())
	if err != nil {
		return "", err
	}
	return memberPayload.RoomAlias, nil
}

// Sync Force reload data for Room, Sync data from puppet API again.
func (r *Room) Sync() error {
	if err := r.GetPuppet().DirtyPayload(schemas.PayloadTypeRoom, r.id); err != nil {
		return err
	}
	if err := r.GetPuppet().DirtyPayload(schemas.PayloadTypeRoomMember, r.id); err != nil {
		return err
	}
	return r.Ready(true)
}

// Say something params {(string | Contact | FileBox | UrlLink | MiniProgram )}
// mentionList @ contact list
func (r *Room) Say(something interface{}, mentionList ..._interface.IContact) (msg _interface.IMessage, err error) {
	var msgID string
	switch v := something.(type) {
	case string:
		msgID, err = r.sayText(v, mentionList...)
	case *Contact:
		msgID, err = r.GetPuppet().MessageSendContact(r.id, v.Id)
	case *filebox.FileBox:
		msgID, err = r.GetPuppet().MessageSendFile(r.id, v)
	case *UrlLink:
		msgID, err = r.GetPuppet().MessageSendURL(r.id, v.payload)
	case *MiniProgram:
		msgID, err = r.GetPuppet().MessageSendMiniProgram(r.id, v.payload)
	default:
		return nil, fmt.Errorf("unsupported arg: %v", something)
	}
	if err != nil {
		return nil, err
	}
	if msgID == "" {
		return nil, nil
	}
	msg = r.GetWechaty().Message().Load(msgID)
	return msg, msg.Ready()
}

func (r *Room) sayText(text string, mentionList ..._interface.IContact) (string, error) {
	var mentionIDList []string
	if len(mentionList) > 0 {
		mentionAlias := make([]string, 0, len(mentionList))
		const atSeparator = config.FourPerEmSpace
		for _, contact := range mentionList {
			mentionIDList = append(mentionIDList, contact.ID())
			alias, _ := r.Alias(contact)
			if alias == "" {
				alias = contact.Name()
			}
			alias = strings.ReplaceAll(alias, " ", atSeparator)
			mentionAlias = append(mentionAlias, "@"+alias)
		}
		text = strings.Join(mentionAlias, atSeparator) + " " + text
	}
	return r.GetPuppet().MessageSendText(r.id, text, mentionIDList...)
}

// Add contact in a room
func (r *Room) Add(contact _interface.IContact) error {
	return r.GetPuppet().RoomAdd(r.id, contact.ID())
}

// Del delete a contact from the room
// it works only when the bot is the owner of the room
func (r *Room) Del(contact _interface.IContact) error {
	return r.GetPuppet().RoomDel(r.id, contact.ID())
}

// Quit the room itself
func (r *Room) Quit() error {
	return r.GetPuppet().RoomQuit(r.id)
}

// Topic get topic from the room
func (r *Room) Topic() string {
	if r.payLoad.Topic != "" {
		return r.payLoad.Topic
	}
	memberList, err := r.memberList()
	if err != nil {
		log.Println("Room Topic err: ", err)
		return ""
	}
	i := 1
	defaultTopic := ""
	for _, member := range memberList {
		if i >= 3 {
			break
		}
		if member.ID() == r.GetPuppet().SelfID() {
			continue
		}
		defaultTopic += member.Name() + ","
		i++
	}
	return strings.TrimRight(defaultTopic, ",")
}

// SetTopic set topic from the room
func (r *Room) SetTopic(topic string) error {
	return r.GetPuppet().SetRoomTopic(r.id, topic)
}

// Announce get announce from the room
func (r *Room) Announce() (string, error) {
	return r.GetPuppet().RoomAnnounce(r.id)
}

// SetAnnounce set announce from the room
// It only works when bot is the owner of the room.
func (r *Room) SetAnnounce(text string) error {
	return r.GetPuppet().SetRoomAnnounce(r.id, text)
}

// QrCode Get QR Code Value of the Room from the room, which can be used as scan and join the room.
func (r *Room) QrCode() (string, error) {
	return r.GetPuppet().RoomQRCode(r.id)
}

// Has check if the room has member `contact`
func (r *Room) Has(contact _interface.IContact) (bool, error) {
	memberIDList, err := r.GetPuppet().RoomMemberList(r.id)
	if err != nil {
		return false, err
	}
	for _, id := range memberIDList {
		if id == contact.ID() {
			return true, nil
		}
	}
	return false, nil
}

// Owner get room's owner from the room.
func (r *Room) Owner() _interface.IContact {
	if r.payLoad.OwnerId == "" {
		return nil
	}
	return r.GetWechaty().Contact().Load(r.payLoad.OwnerId)
}

// Avatar get avatar from the room.
func (r *Room) Avatar() (*filebox.FileBox, error) {
	return r.GetPuppet().RoomAvatar(r.id)
}
