package imageutil

import (
	"image"
	"image/draw"
)

func CopyImage(src image.Image, model ColorModel) image.Image {

	r := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy())
	var c draw.Image

	switch model {
	case MODEL_RGBA:
		c = image.NewRGBA(r)
	case MODEL_RGBA64:
		c = image.NewRGBA64(r)
	case MODEL_NRGBA:
		c = image.NewNRGBA(r)
	case MODEL_NRGBA64:
		c = image.NewNRGBA64(r)
	case MODEL_ALPHA:
		c = image.NewAlpha(r)
	case MODEL_ALPHA16:
		c = image.NewAlpha16(r)
	case MODEL_GRAY:
		c = image.NewGray(r)
	case MODEL_GRAY16:
		c = image.NewGray16(r)
	case MODEL_CMYK:
		c = image.NewCMYK(r)
	}

	copyImage(c, src)
	return c
}

func copyImage(dst draw.Image, src image.Image) {
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
}
