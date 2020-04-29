package _interface

type IRoomFactory interface {
	Load(id string) IRoom
}

type IRoom interface {
	Ready(forceSync bool) (err error)
	IsReady() bool
	String() string
	ID() string
}
