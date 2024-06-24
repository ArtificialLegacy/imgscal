package imageutil

import "image"

func SubImage(img image.Image, x1, y1, x2, y2 int, copy bool) image.Image {
	upLeft := image.Point{x1, y1}
	lowRight := image.Point{x2, y2}
	rect := image.Rectangle{upLeft, lowRight}

	n := img
	if copy {
		n = CopyImage(img)
	}

	switch nimg := n.(type) {
	case *image.RGBA:
		return nimg.SubImage(rect)
	case *image.RGBA64:
		return nimg.SubImage(rect)
	case *image.NRGBA:
		return nimg.SubImage(rect)
	case *image.NRGBA64:
		return nimg.SubImage(rect)
	case *image.Alpha:
		return nimg.SubImage(rect)
	case *image.Alpha16:
		return nimg.SubImage(rect)
	case *image.Gray:
		return nimg.SubImage(rect)
	case *image.Gray16:
		return nimg.SubImage(rect)
	case *image.CMYK:
		return nimg.SubImage(rect)
	}

	return nil
}
