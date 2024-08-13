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

func TableToPoint(state *golua.LState, t *golua.LTable) image.Point {
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
