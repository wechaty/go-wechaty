package puppet_openwechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
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
