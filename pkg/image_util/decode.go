package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func Decode(r io.Reader, encoding ImageEncoding) (image.Image, error) {
	switch encoding {
	case ENCODING_PNG:
		strip := PNGChunkStripper{
			Reader: r,
		}
		return png.Decode(io.Reader(&strip))

	case ENCODING_JPEG:
		return jpeg.Decode(r)

	case ENCODING_GIF:
		return gif.Decode(r)
	}

	return nil, fmt.Errorf("cannot decode unsupported encoding: %d", encoding)
}
