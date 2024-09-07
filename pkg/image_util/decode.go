package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	goico "github.com/ArtificialLegacy/go-ico"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func Decode(r io.ReadSeeker, encoding ImageEncoding) (image.Image, error) {
	switch encoding {
	case ENCODING_PNG:
		strip := PNGChunkStripper{
			Reader: r,
		}
		return png.Decode(&strip)

	case ENCODING_JPEG:
		return jpeg.Decode(r)

	case ENCODING_GIF:
		return gif.Decode(r)

	case ENCODING_TIFF:
		return tiff.Decode(r)

	case ENCODING_BMP:
		return bmp.Decode(r)

	// generic decoding of ICO only keeps the largest image
	case ENCODING_ICO:
		fallthrough
	case ENCODING_CUR:
		cfg, imgs, err := goico.Decode(r)
		if err != nil {
			return nil, err
		}

		return imgs[cfg.Largest], nil
	}

	return nil, fmt.Errorf("cannot decode unsupported encoding: %d", encoding)
}
