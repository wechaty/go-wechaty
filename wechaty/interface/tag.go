package _interface

type ITagFactory interface {
	Load(id string) ITag
	Get(tag string) ITag
	Delete(tag ITag) error
}

type ITag interface {
	ID() string
	Add(to IContact) error
	Remove(from IContact) error
}
