package telegram

const (
	TextPost PostType = iota
	GalleryPost
	GifPost
	VideoPost
	AudioPost
)

type (
	PostType int8
)
