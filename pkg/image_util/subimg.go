package imageutil

import "image"

func SubImage(img image.Image, x1, y1, x2, y2 int, copy bool) image.Image {
	upLeft := image.Point{x1, y1}
	lowRight := image.Point{x2, y2}
	rect := image.Rectangle{upLeft, lowRight}

	switch nimg := img.(type) {
	case *image.RGBA:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.RGBA64:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.NRGBA:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.NRGBA64:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.Alpha:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.Alpha16:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.Gray:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.Gray16:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	case *image.CMYK:
		i := nimg.SubImage(rect)
		if copy {
			return CopyImage(i)
		} else {
			return i
		}
	}

	return nil
}
