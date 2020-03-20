package user

type Contact struct {
	Id string
}

func (r *Contact) Load(id string) Contact {
	return Contact{}
}

func (r *Contact) Ready(forceSync bool) bool {
	return true
}
