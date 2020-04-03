package schemas

type MiniProgramPayload struct {
  Appid       string // optional, Appid, get from wechat (mp.weixin.qq.com)
  Description string // optional, mini program title
  PagePath    string // optional, mini program page path
  ThumbUrl    string // optional, default picture, convert to thumbnail
  Title       string // optional, mini program title
  Username    string // original ID, get from wechat (mp.weixin.qq.com)
  ThumbKey    string // original, thumbnailurl and thumbkey will make the headphoto of mini-program better
}
