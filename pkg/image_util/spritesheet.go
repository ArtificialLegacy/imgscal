package imageutil

import (
	"image"
	"math"
)

func SpritesheetToFrames(img image.Image, copy bool, count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep int) []image.Image {
	imgs := make([]image.Image, count)

	offsetx := hpixel + (hcell * width) + (hsep * hcell)
	offsety := vpixel + (vcell * height) + (vsep * vcell)

	col := index % perRow
	row := index / perRow

	topx := offsetx + img.Bounds().Min.X + (col*width + col*hsep)
	topy := offsety + img.Bounds().Min.Y + (row*height + row*vsep)
	bottomx := topx + width
	bottomy := topy + height

	for ind := range count {
		simg := SubImage(img, topx, topy, bottomx, bottomy, copy)
		imgs[ind] = simg

		if (ind+1+index)%perRow == 0 {
			topx = offsetx
			bottomx = topx + width

			topy += height + vsep
			bottomy = topy + width
		} else {
			topx += width + hsep
			bottomx = topx + width
		}
	}

	return imgs
}

func FramesToSpritesheet(imgs []image.Image, model ColorModel, count, width, height, perRow, hpixel, vpixel, hcell, vcell, index, hsep, vsep int) image.Image {
	imgs = imgs[index:]
	if count > len(imgs) {
		count = len(imgs)
	}

	if perRow == 0 {
		perRow = count
	}

	rows := int(math.Ceil(float64(count) / float64(perRow)))

	offsetx := hpixel + (hcell * width) + (hsep * hcell)
	offsety := vpixel + (vcell * height) + (vsep * vcell)

	ssWidth := offsetx*2 + perRow*width + hsep*(perRow-1)
	ssHeight := offsety*2 + rows*height + vsep*(rows-1)

	img := NewImage(ssWidth, ssHeight, model)

	for i, frame := range imgs {
		col := i % perRow
		row := i / perRow
		x := offsetx + col*width + hsep*col
		y := offsety + row*height + vsep*row

		Draw(img, frame, x, y, width, height)
	}

	return img
}
