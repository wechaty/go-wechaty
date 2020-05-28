package _interface

type IUrlLinkFactory interface {
	Create(url string) (IUrlLink, error)
}

type IUrlLink interface {
	String() string
	Url() string
	Title() string
	ThumbnailUrl() string
	Description() string
}
