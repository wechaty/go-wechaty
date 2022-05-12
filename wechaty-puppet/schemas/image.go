package schemas

//go:generate stringer -type=ImageType
type ImageType uint8

const (
	ImageTypeUnknown ImageType = iota
	ImageTypeThumbnail
	ImageTypeHD
	ImageTypeArtwork
)
