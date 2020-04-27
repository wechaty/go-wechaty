package _interface

type IContactFactory interface {
	Load(id string) IContact
	LoadSelf(id string) IContact
}

type IContact interface {
	Ready(forceSync bool) (err error)
	IsReady() bool
	Sync() error
	String() string
	ID() string
	Name() string
}
