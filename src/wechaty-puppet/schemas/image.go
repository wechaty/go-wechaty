package schemas

type ImageType uint8

const (
	ImageTypeUnknown   ImageType = 0
	ImageTypeThumbnail           = 1
	ImageTypeHD                  = 2
	ImageTypeArtwork             = 3
)
