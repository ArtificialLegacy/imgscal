package imageutil

import (
	"image"
	"image/color"

	"github.com/crazy3lf/colorconv"
	golua "github.com/yuin/gopher-lua"
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

func Get(img image.Image, x, y int) (int, int, int, int) {
	cr, cg, cb, ca := img.At(x, y).RGBA()
	return int(cr), int(cg), int(cb), int(ca)
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

const (
	COLOR_TYPE_RGBA  string = "rgba"
	COLOR_TYPE_HSVA  string = "hsva"
	COLOR_TYPE_HSLA  string = "hsla"
	COLOR_TYPE_GRAYA string = "graya"
)

func RGBAColorToColorTable(state *golua.LState, rgba *color.RGBA) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_RGBA))
	t.RawSetString("red", golua.LNumber(rgba.R))
	t.RawSetString("green", golua.LNumber(rgba.G))
	t.RawSetString("blue", golua.LNumber(rgba.B))
	t.RawSetString("alpha", golua.LNumber(rgba.A))

	return t
}

func RGBAToColorTable(state *golua.LState, r, g, b, a int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_RGBA))
	t.RawSetString("red", golua.LNumber(r))
	t.RawSetString("green", golua.LNumber(g))
	t.RawSetString("blue", golua.LNumber(b))
	t.RawSetString("alpha", golua.LNumber(a))

	return t
}

func HSVAToColorTable(state *golua.LState, hue, sat, val float64, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_HSVA))
	t.RawSetString("hue", golua.LNumber(hue))
	t.RawSetString("saturation", golua.LNumber(sat))
	t.RawSetString("value", golua.LNumber(val))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func HSLAToColorTable(state *golua.LState, hue, sat, light float64, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_HSLA))
	t.RawSetString("hue", golua.LNumber(hue))
	t.RawSetString("saturation", golua.LNumber(sat))
	t.RawSetString("light", golua.LNumber(light))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func GrayAToColorTable(state *golua.LState, gray, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_GRAYA))
	t.RawSetString("gray", golua.LNumber(gray))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func ParseRGBATable(t *golua.LTable) (uint8, uint8, uint8, uint8) {
	cr := t.RawGetString("red").(golua.LNumber)
	cg := t.RawGetString("green").(golua.LNumber)
	cb := t.RawGetString("blue").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint8(cr), uint8(cg), uint8(cb), uint8(ca)
}

func ParseHSVATable(t *golua.LTable) (float64, float64, float64, uint8) {
	ch := t.RawGetString("hue").(golua.LNumber)
	cs := t.RawGetString("sat").(golua.LNumber)
	cv := t.RawGetString("value").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return float64(ch), float64(cs), float64(cv), uint8(ca)
}

func ParseHSLATable(t *golua.LTable) (float64, float64, float64, uint8) {
	ch := t.RawGetString("hue").(golua.LNumber)
	cs := t.RawGetString("sat").(golua.LNumber)
	cl := t.RawGetString("light").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return float64(ch), float64(cs), float64(cl), uint8(ca)
}

func ParseGrayATable(t *golua.LTable) (uint8, uint8) {
	cy := t.RawGetString("gray").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint8(cy), uint8(ca)
}

func ColorTableToRGBAColor(t *golua.LTable) *color.RGBA {
	cr, cg, cb, ca := ColorTableToRGBA(t)

	return &color.RGBA{
		R: cr,
		G: cg,
		B: cb,
		A: ca,
	}
}

func ColorTableToRGBA(t *golua.LTable) (uint8, uint8, uint8, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		return cr, cg, cb, ca

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0, 0, 0
		}

		return cr, cg, cb, ca

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0, 0, 0
		}

		return cr, cg, cb, ca

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		return cy, cy, cy, ca
	}

	return 0, 0, 0, 0
}

func ColorTableToHSVA(t *golua.LTable) (float64, float64, float64, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		ch, cs, cv := colorconv.RGBToHSV(cr, cg, cb)
		return ch, cs, cv, ca

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		return ch, cs, cv, ca

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0, 0, 0
		}
		ch, cs, cv := colorconv.RGBToHSV(cr, cg, cb)
		return ch, cs, cv, ca

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		ch, cs, cv := colorconv.RGBToHSV(cy, cy, cy)
		return ch, cs, cv, ca
	}

	return 0, 0, 0, 0
}

func ColorTableToHSLA(t *golua.LTable) (float64, float64, float64, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		ch, cs, cl := colorconv.RGBToHSL(cr, cg, cb)
		return ch, cs, cl, ca

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0, 0, 0
		}
		ch, cs, cl := colorconv.RGBToHSL(cr, cg, cb)
		return ch, cs, cl, ca

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		return ch, cs, cl, ca

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		ch, cs, cl := colorconv.RGBToHSL(cy, cy, cy)
		return ch, cs, cl, ca
	}

	return 0, 0, 0, 0
}

func ColorTableToGrayA(t *golua.LTable) (uint8, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y, ca

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y, ca

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y, ca

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		return cy, ca
	}

	return 0, 0
}
