package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	goico "github.com/ArtificialLegacy/go-ico"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func Encode(w io.WriteSeeker, img image.Image, encoding ImageEncoding) error {
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
	case ENCODING_WEBP:
		options, err := encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 100)
		if err != nil {
			return err
		}
		return webp.Encode(w, img, options)
	case ENCODING_ICO:
		imgs := []image.Image{img}
		ico, err := goico.NewICOConfig(imgs)
		if err != nil {
			return err
		}
		return goico.Encode(w, ico, imgs)
	case ENCODING_CUR:
		imgs := []image.Image{img}
		ico, err := goico.NewCURConfig(imgs, []int{img.Bounds().Dx() / 2, img.Bounds().Dy() / 2})
		if err != nil {
			return err
		}
		return goico.Encode(w, ico, imgs)
	}

	return fmt.Errorf("cannot encode unsupported encoding: %d", encoding)
}
