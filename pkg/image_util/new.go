package imageutil

import (
	"image"
)

func NewImage(width, height int, model ColorModel) image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	rect := image.Rectangle{upLeft, lowRight}

	var img image.Image

	switch model {
	case MODEL_RGBA:
		img = image.NewRGBA(rect)
	case MODEL_RGBA64:
		img = image.NewRGBA64(rect)
	case MODEL_NRGBA:
		img = image.NewNRGBA(rect)
	case MODEL_NRGBA64:
		img = image.NewNRGBA64(rect)
	case MODEL_ALPHA:
		img = image.NewAlpha(rect)
	case MODEL_ALPHA16:
		img = image.NewAlpha16(rect)
	case MODEL_GRAY:
		img = image.NewGray(rect)
	case MODEL_GRAY16:
		img = image.NewGray16(rect)
	case MODEL_CMYK:
		img = image.NewCMYK(rect)
	}

	return img
}
