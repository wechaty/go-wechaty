package factory

import (
	"errors"
	"github.com/otiai10/opengraph"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

var (
	ErrImageUrlOrDescNotFound = errors.New("imgUrl.or.desc.not.found")
)

type UrlLinkFactory struct{}

func (u *UrlLinkFactory) Create(url string) (_interface.IUrlLink, error) {
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

	return user.NewUrlLink(payload), nil
}
