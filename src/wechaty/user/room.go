package user

type Room struct {
	Id string
}

func (r *Room) Load(id string) Room {
	return Room{}
}

func (r *Room) Ready(forceSync bool) bool {
	return true
}
