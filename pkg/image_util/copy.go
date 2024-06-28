package imageutil

import (
	"image"
	"image/draw"
)

func CopyImage(src image.Image) image.Image {

	r := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy())

	switch img := src.(type) {
	case *image.RGBA:
		c := image.NewRGBA(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.RGBA64:
		c := image.NewRGBA64(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.NRGBA:
		c := image.NewNRGBA(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.NRGBA64:
		c := image.NewNRGBA64(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Alpha:
		c := image.NewAlpha(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Alpha16:
		c := image.NewAlpha16(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Gray:
		c := image.NewGray(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Gray16:
		c := image.NewGray16(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.CMYK:
		c := image.NewCMYK(r)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	}

	return nil
}
