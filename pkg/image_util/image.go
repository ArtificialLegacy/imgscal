package imageutil

import (
	"image"
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

func TableToPoint(state *golua.LState, t *golua.LTable) image.Point {
	x := state.GetTable(t, golua.LString("x")).(golua.LNumber)
	y := state.GetTable(t, golua.LString("y")).(golua.LNumber)

	return image.Point{
		X: int(x),
		Y: int(y),
	}
}

func PointToTable(state *golua.LState, p image.Point) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("x"), golua.LNumber(p.X))
	state.SetTable(t, golua.LString("y"), golua.LNumber(p.Y))

	return t
}

var EncodingList = []ImageEncoding{
	ENCODING_PNG,
	ENCODING_JPEG,
	ENCODING_GIF,
	ENCODING_UNKNOWN,
}
