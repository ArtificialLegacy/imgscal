package imageutil

import (
	"image"
	"image/draw"
	"strings"

	golua "github.com/yuin/gopher-lua"
)

type ImageEncoding int

const (
	ENCODING_PNG ImageEncoding = iota
	ENCODING_JPEG
	ENCODING_GIF
	ENCODING_UNKNOWN
)

func EncodingExtension(encoding ImageEncoding) string {
	switch encoding {
	case ENCODING_PNG:
		return ".png"
	case ENCODING_JPEG:
		return ".jpg"
	case ENCODING_GIF:
		return ".gif"
	default:
		return ".unknown"
	}
}

func ExtensionEncoding(ext string) ImageEncoding {
	ext = strings.ToLower(ext)

	switch ext {
	case ".png":
		return ENCODING_PNG
	case ".jpg":
		return ENCODING_JPEG
	case ".gif":
		return ENCODING_GIF
	}

	return ENCODING_UNKNOWN
}

func ImageGetDraw(img image.Image) draw.Image {
	switch i := img.(type) {
	case *image.RGBA:
		return i
	case *image.RGBA64:
		return i
	case *image.NRGBA:
		return i
	case *image.NRGBA64:
		return i
	case *image.Alpha:
		return i
	case *image.Alpha16:
		return i
	case *image.Gray:
		return i
	case *image.Gray16:
		return i
	case *image.CMYK:
		return i
	default:
		return nil
	}
}

func ImageCompare(img1, img2 image.Image) bool {
	draw1 := ImageGetDraw(img1)
	draw2 := ImageGetDraw(img2)
	if draw1 == nil || draw2 == nil {
		return false
	}

	bounds1 := draw1.Bounds()
	bounds2 := draw2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return false
	}

	for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
		for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
			zx1 := x - bounds1.Min.X
			zy1 := y - bounds1.Min.Y
			x2 := bounds2.Min.X + zx1
			y2 := bounds2.Min.Y + zy1

			c1 := draw1.At(x, y)
			c2 := draw2.At(x2, y2)

			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}

	return true
}

func TableToPoint(t *golua.LTable) image.Point {
	x := t.RawGetString("x").(golua.LNumber)
	y := t.RawGetString("y").(golua.LNumber)

	return image.Point{
		X: int(x),
		Y: int(y),
	}
}

func PointToTable(state *golua.LState, p image.Point) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("x", golua.LNumber(p.X))
	t.RawSetString("y", golua.LNumber(p.Y))

	return t
}

var EncodingList = []ImageEncoding{
	ENCODING_PNG,
	ENCODING_JPEG,
	ENCODING_GIF,
	ENCODING_UNKNOWN,
}
