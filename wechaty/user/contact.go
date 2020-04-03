package user

import "github.com/wechaty/go-wechaty/wechaty"

type Contact struct {
	wechaty.Accessory

	Id string
}

func (r *Contact) Load(id string) Contact {
	return Contact{}
}

func (r *Contact) Ready(forceSync bool) bool {
	return true
}
