package imageutil

import (
	"image"
	"image/color"
	"image/draw"
)

func Draw(base image.Image, sub image.Image, x, y, width, height int) {
	r := image.Rect(x, y, x+width, y+height)
	DrawRect(base, sub, r)
}

func DrawRect(base image.Image, sub image.Image, r image.Rectangle) {
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

func AlphaSet(img image.Image, alpha uint8) {
	switch img.(type) {
	case *image.Gray:
		return
	case *image.Gray16:
		return
	case *image.CMYK:
		return
	}

	imgDraw := ImageGetDraw(img)

	minx := img.Bounds().Min.X
	miny := img.Bounds().Min.Y
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y

	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			c := imgDraw.At(x, y)
			switch col := c.(type) {
			case color.RGBA:
				col.A = alpha
				c = col
			case color.RGBA64:
				col.A = uint16(alpha)
				c = col
			case color.NRGBA:
				col.A = alpha
				c = col
			case color.NRGBA64:
				col.A = uint16(alpha)
				c = col
			case color.Alpha:
				col.A = alpha
				c = col
			case color.Alpha16:
				col.A = uint16(alpha)
				c = col
			}

			imgDraw.Set(x, y, c)
		}
	}
}
