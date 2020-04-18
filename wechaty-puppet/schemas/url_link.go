package schemas

import "encoding/json"

type UrlLinkPayload struct {
  Description  string `json:"description"`
  ThumbnailUrl string `json:"thumbnailUrl"`
  Title        string `json:"title"`
  Url          string `json:"url"`
}

func (u *UrlLinkPayload) ToJson() string {
  b, _ := json.Marshal(u)
  return string(b)
}
