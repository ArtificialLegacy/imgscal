package imageutil

import "image"

func Limit(img image.Image, model ColorModel) (image.Image, ColorModel) {
	switch i := img.(type) {
	case *image.RGBA:
		return i, MODEL_RGBA
	case *image.RGBA64:
		return i, MODEL_RGBA64
	case *image.NRGBA:
		return i, MODEL_NRGBA
	case *image.NRGBA64:
		return i, MODEL_NRGBA64
	case *image.Alpha:
		return i, MODEL_ALPHA
	case *image.Alpha16:
		return i, MODEL_ALPHA16
	case *image.Gray:
		return i, MODEL_GRAY
	case *image.Gray16:
		return i, MODEL_GRAY16
	case *image.CMYK:
		return i, MODEL_CMYK
	default:
		return CopyImage(i, model), model
	}
}
