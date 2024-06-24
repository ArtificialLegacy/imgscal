package imageutil

import (
	"image"
	"image/draw"
)

func CopyImage(src image.Image) image.Image {

	switch img := src.(type) {
	case *image.RGBA:
		c := image.NewRGBA(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.RGBA64:
		c := image.NewRGBA64(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.NRGBA:
		c := image.NewNRGBA(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.NRGBA64:
		c := image.NewNRGBA64(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Alpha:
		c := image.NewAlpha(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Alpha16:
		c := image.NewAlpha16(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Gray:
		c := image.NewGray(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.Gray16:
		c := image.NewGray16(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	case *image.CMYK:
		c := image.NewCMYK(img.Rect)
		draw.Draw(c, c.Rect, img, img.Rect.Min, draw.Src)
		return c
	}

	return nil
}
