package _interface

type IMiniProgram interface {
	AppID() string
	Description() string
	PagePath() string
	ThumbUrl() string
	Title() string
	Username() string
	ThumbKey() string
}
