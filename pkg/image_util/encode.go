package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func Encode(w io.Writer, img image.Image, encoding ImageEncoding) error {
	switch encoding {
	case ENCODING_PNG:
		return png.Encode(w, img)
	case ENCODING_JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{})
	case ENCODING_GIF:
		return gif.Encode(w, img, &gif.Options{})
	case ENCODING_TIFF:
		return tiff.Encode(w, img, &tiff.Options{})
	case ENCODING_BMP:
		return bmp.Encode(w, img)
	}

	return fmt.Errorf("cannot encode unsupported encoding: %d", encoding)
}
