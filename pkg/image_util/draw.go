package imageutil

import (
	"image"
	"image/draw"
)

func Draw(base image.Image, sub image.Image, x, y, width, height int) {
	r := image.Rect(x, y, x+width, y+height)

	switch img := base.(type) {
	case *image.RGBA:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.RGBA64:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.NRGBA:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.NRGBA64:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.Alpha:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.Alpha16:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.Gray:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.Gray16:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	case *image.CMYK:
		draw.Draw(img, r, sub, sub.Bounds().Min, draw.Src)
	}
}
