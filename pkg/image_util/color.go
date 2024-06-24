package imageutil

import (
	"image"
	"image/color"
)

func Set(img image.Image, x, y, red, green, blue, alpha int) {
	col := color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha)}

	switch i := img.(type) {
	case *image.RGBA:
		c := color.RGBAModel.Convert(col)
		i.Set(x, y, c)
	case *image.RGBA64:
		c := color.RGBA64Model.Convert(col)
		i.Set(x, y, c)
	case *image.NRGBA:
		c := color.NRGBAModel.Convert(col)
		i.Set(x, y, c)
	case *image.NRGBA64:
		c := color.NRGBA64Model.Convert(col)
		i.Set(x, y, c)
	case *image.Alpha:
		c := color.AlphaModel.Convert(col)
		i.Set(x, y, c)
	case *image.Alpha16:
		c := color.Alpha16Model.Convert(col)
		i.Set(x, y, c)
	case *image.Gray:
		c := color.GrayModel.Convert(col)
		i.Set(x, y, c)
	case *image.Gray16:
		c := color.Gray16Model.Convert(col)
		i.Set(x, y, c)
	case *image.CMYK:
		c := color.CMYKModel.Convert(col)
		i.Set(x, y, c)
	}
}

func ConvertColor(model ColorModel, red, green, blue, alpha int) (int, int, int, int) {
	col := color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha)}

	var re, gr, bl, al uint32

	switch model {
	case MODEL_RGBA:
		c := color.RGBAModel.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_RGBA64:
		c := color.RGBA64Model.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_NRGBA:
		c := color.NRGBAModel.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_NRGBA64:
		c := color.NRGBA64Model.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_ALPHA:
		c := color.AlphaModel.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_ALPHA16:
		c := color.Alpha16Model.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_GRAY:
		c := color.GrayModel.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_GRAY16:
		c := color.Gray16Model.Convert(col)
		re, gr, bl, al = c.RGBA()
	case MODEL_CMYK:
		c := color.CMYKModel.Convert(col)
		re, gr, bl, al = c.RGBA()
	}

	return int(re), int(gr), int(bl), int(al)
}
