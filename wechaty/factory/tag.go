package factory

import (
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"sync"
)

type TagFactory struct {
	_interface.IAccessory
	pool *sync.Map
}

// NewTagFactory ...
func NewTagFactory(accessory _interface.IAccessory) *TagFactory {
	return &TagFactory{
		IAccessory: accessory,
		pool:       &sync.Map{},
	}
}

func (r *TagFactory) Load(id string) _interface.ITag {
	load, ok := r.pool.Load(id)
	if ok {
		return load.(*user.Tag)
	}
	tag := user.NewTag(id, r.IAccessory)
	r.pool.Store(id, tag)
	return tag
}

func (r *TagFactory) Get(tag string) _interface.ITag {
	return r.Load(tag)
}

func (r *TagFactory) Delete(tag _interface.ITag) error {
	return r.GetPuppet().TagContactDelete(tag.ID())
}
