package schemas

import "encoding/json"

type MiniProgramPayload struct {
	Appid       string `json:"appid"`       // optional, Appid, get from wechat (mp.weixin.qq.com)
	Description string `json:"description"` // optional, mini program title
	PagePath    string `json:"pagePath"`    // optional, mini program page path
	ThumbUrl    string `json:"thumbUrl"`    // optional, default picture, convert to thumbnail
	Title       string `json:"title"`       // optional, mini program title
	Username    string `json:"username"`    // original ID, get from wechat (mp.weixin.qq.com)
	ThumbKey    string `json:"thumbKey"`    // original, thumbnailurl and thumbkey will make the headphoto of mini-program better
	ShareId     string `json:"shareId"`     // optional, the unique userId for who share this mini program
	IconUrl     string `json:"iconUrl"`     // optional, mini program icon url
}

func (m *MiniProgramPayload) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}
