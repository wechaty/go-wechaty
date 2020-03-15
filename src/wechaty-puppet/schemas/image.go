package schemas

type ImageType uint8

const (
	Unknown   ImageType = 0
	Thumbnail           = 1
	HD                  = 2
	Artwork             = 3
)
