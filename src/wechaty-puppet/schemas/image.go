package schemas

type ImageType uint8

const (
  Unknown   ImageType = 0
  Thumbnail ImageType = 1
  HD        ImageType = 2
  Artwork   ImageType = 3
)
