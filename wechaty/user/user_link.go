package user

import (
	"errors"
	"fmt"

	"github.com/otiai10/opengraph"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

var (
	ErrImageUrlOrDescNotFound = errors.New("imgUrl.or.desc.not.found")
)

type UrlLink struct {
	payload *schemas.UrlLinkPayload
}

func NewUrlLink(payload *schemas.UrlLinkPayload) *UrlLink {
	return &UrlLink{payload: payload}
}

func (ul *UrlLink) String() string {
	return fmt.Sprintf("UrlLink<%s>", ul.Url())
}

func (ul *UrlLink) Url() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Url
}

func (ul *UrlLink) Title() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Title
}

func (ul *UrlLink) ThumbnailUrl() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.ThumbnailUrl
}

func (ul *UrlLink) Description() string {
	if ul.payload == nil {
		return ""
	}
	return ul.payload.Description
}

func CreateUrlLink(url string) (*UrlLink, error) {
	var og, err = opengraph.Fetch(url)
	if err != nil {
		return nil, err
	}
	var payload = &schemas.UrlLinkPayload{
		Url:         url,
		Title:       og.Title,
		Description: og.Description,
	}

	if len(og.Image) != 0 {
		payload.ThumbnailUrl = og.Image[0].URL
	}

	if payload.ThumbnailUrl == "" || payload.Description == "" {
		return nil, ErrImageUrlOrDescNotFound
	}

	return NewUrlLink(payload), nil
}
