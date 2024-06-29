package imageutil

import "image"

func SubImage(img image.Image, x1, y1, x2, y2 int, copy bool) image.Image {
	upLeft := image.Point{x1, y1}
	lowRight := image.Point{x2, y2}
	rect := image.Rectangle{upLeft, lowRight}

	switch nimg := img.(type) {
	case *image.RGBA:
		return subimg(nimg, rect, copy, MODEL_RGBA)
	case *image.RGBA64:
		return subimg(nimg, rect, copy, MODEL_RGBA64)
	case *image.NRGBA:
		return subimg(nimg, rect, copy, MODEL_NRGBA)
	case *image.NRGBA64:
		return subimg(nimg, rect, copy, MODEL_NRGBA64)
	case *image.Alpha:
		return subimg(nimg, rect, copy, MODEL_ALPHA)
	case *image.Alpha16:
		return subimg(nimg, rect, copy, MODEL_ALPHA16)
	case *image.Gray:
		return subimg(nimg, rect, copy, MODEL_GRAY)
	case *image.Gray16:
		return subimg(nimg, rect, copy, MODEL_GRAY16)
	case *image.CMYK:
		return subimg(nimg, rect, copy, MODEL_CMYK)
	}

	return nil
}

type imgsubimg interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func subimg(img imgsubimg, rect image.Rectangle, copy bool, model ColorModel) image.Image {
	i := img.SubImage(rect)
	if copy {
		return CopyImage(i, model)
	} else {
		return i
	}
}
