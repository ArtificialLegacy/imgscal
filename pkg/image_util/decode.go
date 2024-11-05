package imageutil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	goico "github.com/ArtificialLegacy/go-ico"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
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

	case ENCODING_WEBP:
		return webp.Decode(r, &decoder.Options{})

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

func DecodeConfig(r io.Reader, encoding ImageEncoding) (int, int, error) {
	var cfg image.Config
	var err error

	switch encoding {
	case ENCODING_PNG:
		strip := PNGChunkStripper{
			Reader: r,
		}
		cfg, err = png.DecodeConfig(&strip)

	case ENCODING_JPEG:
		cfg, err = jpeg.DecodeConfig(r)

	case ENCODING_GIF:
		cfg, err = gif.DecodeConfig(r)

	case ENCODING_TIFF:
		cfg, err = tiff.DecodeConfig(r)

	case ENCODING_BMP:
		cfg, err = bmp.DecodeConfig(r)

	case ENCODING_WEBP:
		cfg, err = webp.DecodeConfig(r, &decoder.Options{})

	case ENCODING_ICO:
		fallthrough
	case ENCODING_CUR:
		fcfg, ferr := goico.DecodeConfig(r)
		cfg, err = image.Config{
			Width:  fcfg.Entries[fcfg.Largest].Width,
			Height: fcfg.Entries[fcfg.Largest].Height,
		}, ferr

	default:
		return 0, 0, fmt.Errorf("unsupported encoding: %d", encoding)
	}

	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}
