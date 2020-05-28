package factory

import (
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"sync"
)

type TagFactory struct {
	_interface.Accessory
	pool *sync.Map
}

func (r *TagFactory) Load(id string) _interface.ITag {
	load, ok := r.pool.Load(id)
	if ok {
		return load.(*user.Tag)
	}
	tag := user.NewTag(id, r.Accessory)
	r.pool.Store(id, tag)
	return tag
}

func (r *TagFactory) Get(tag string) _interface.ITag {
	return r.Load(tag)
}

func (r *TagFactory) Delete(tag _interface.ITag) error {
	return r.GetPuppet().TagContactDelete(tag.ID())
}
