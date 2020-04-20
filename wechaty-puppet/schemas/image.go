package schemas

//go:generate stringer -type=ImageType
type ImageType uint8

const (
	ImageTypeUnknown   ImageType = 0
	ImageTypeThumbnail           = 1
	ImageTypeHD                  = 2
	ImageTypeArtwork             = 3
)
