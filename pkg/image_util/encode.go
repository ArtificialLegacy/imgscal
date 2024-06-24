package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func Encode(w io.Writer, img image.Image, encoding ImageEncoding) error {
	switch encoding {
	case ENCODING_PNG:
		return png.Encode(w, img)
	case ENCODING_JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{})
	case ENCODING_GIF:
		return gif.Encode(w, img, &gif.Options{})
	}

	return fmt.Errorf("cannot encode unsupported encoding: %d", encoding)
}
