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

func GetColor(img image.Image, state *golua.LState, x, y int) *golua.LTable {
	switch i := img.(type) {
	case *image.RGBA:
		r, g, b, a := i.RGBAAt(x, y).RGBA()
		return RGBAToColorTable(state, int(r), int(g), int(b), int(a))
	case *image.RGBA64:
		r, g, b, a := i.RGBA64At(x, y).RGBA()
		return RGBAToColorTable(state, int(r), int(g), int(b), int(a))
	case *image.NRGBA:
		c := i.NRGBAAt(x, y)
		return RGBAToColorTable(state, int(c.R), int(c.G), int(c.B), int(c.A))
	case *image.NRGBA64:
		c := i.NRGBA64At(x, y)
		return RGBAToColorTable(state, int(c.R), int(c.G), int(c.B), int(c.A))
	case *image.Alpha:
		_, _, _, a := i.AlphaAt(x, y).RGBA()
		return AlphaToColorTable(state, int(a))
	case *image.Alpha16:
		_, _, _, a := i.Alpha16At(x, y).RGBA()
		return Alpha16ToColorTable(state, int(a))
	case *image.Gray:
		r, _, _, _ := i.GrayAt(x, y).RGBA()
		return GrayToColorTable(state, int(r))
	case *image.Gray16:
		r, _, _, _ := i.Gray16At(x, y).RGBA()
		return Gray16ToColorTable(state, int(r))
	case *image.CMYK:
		r, g, b, a := i.CMYKAt(x, y).RGBA()
		return CMYKToColorTable(state, int(r), int(g), int(b), int(a))

	}

	return nil
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

func Color8BitTo16Bit(c uint8) uint16 {
	return uint16(c)<<8 | uint16(c)
}

func Color16BitTo8Bit(c uint16) uint8 {
	return uint8(c >> 8)
}

const (
	COLOR_TYPE_RGBA    string = "rgba"
	COLOR_TYPE_HSVA    string = "hsva"
	COLOR_TYPE_HSLA    string = "hsla"
	COLOR_TYPE_GRAYA   string = "graya"
	COLOR_TYPE_GRAYA16 string = "graya16"
	COLOR_TYPE_GRAY    string = "gray"
	COLOR_TYPE_GRAY16  string = "gray16"
	COLOR_TYPE_ALPHA   string = "alpha"
	COLOR_TYPE_ALPHA16 string = "alpha16"
	COLOR_TYPE_CMYKA   string = "cmyka"
	COLOR_TYPE_CMYK    string = "cmyk"
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

func GrayA16ToColorTable(state *golua.LState, gray, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_GRAYA16))
	t.RawSetString("gray", golua.LNumber(gray))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func GrayToColorTable(state *golua.LState, gray int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_GRAY))
	t.RawSetString("gray", golua.LNumber(gray))

	return t
}

func Gray16ToColorTable(state *golua.LState, gray int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_GRAY16))
	t.RawSetString("gray", golua.LNumber(gray))

	return t
}

func AlphaToColorTable(state *golua.LState, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_ALPHA))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func Alpha16ToColorTable(state *golua.LState, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_ALPHA16))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func CMYKAToColorTable(state *golua.LState, cyan, magenta, yellow, key, alpha int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_CMYKA))
	t.RawSetString("cyan", golua.LNumber(cyan))
	t.RawSetString("magenta", golua.LNumber(magenta))
	t.RawSetString("yellow", golua.LNumber(yellow))
	t.RawSetString("key", golua.LNumber(key))
	t.RawSetString("alpha", golua.LNumber(alpha))

	return t
}

func CMYKToColorTable(state *golua.LState, cyan, magenta, yellow, key int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("type", golua.LString(COLOR_TYPE_CMYK))
	t.RawSetString("cyan", golua.LNumber(cyan))
	t.RawSetString("magenta", golua.LNumber(magenta))
	t.RawSetString("yellow", golua.LNumber(yellow))
	t.RawSetString("key", golua.LNumber(key))

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

func ParseGrayA16Table(t *golua.LTable) (uint16, uint16) {
	cy := t.RawGetString("gray").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint16(cy), uint16(ca)
}

func ParseGrayTable(t *golua.LTable) uint8 {
	cy := t.RawGetString("gray").(golua.LNumber)

	return uint8(cy)
}

func ParseGray16Table(t *golua.LTable) uint16 {
	cy := t.RawGetString("gray").(golua.LNumber)

	return uint16(cy)
}

func ParseAlphaTable(t *golua.LTable) uint8 {
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint8(ca)
}

func ParseAlpha16Table(t *golua.LTable) uint16 {
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint16(ca)
}

func ParseCMYKATable(t *golua.LTable) (uint8, uint8, uint8, uint8, uint8) {
	cc := t.RawGetString("cyan").(golua.LNumber)
	cm := t.RawGetString("magenta").(golua.LNumber)
	cy := t.RawGetString("yellow").(golua.LNumber)
	ck := t.RawGetString("key").(golua.LNumber)
	ca := t.RawGetString("alpha").(golua.LNumber)

	return uint8(cc), uint8(cm), uint8(cy), uint8(ck), uint8(ca)
}

func ParseCMYKTable(t *golua.LTable) (uint8, uint8, uint8, uint8) {
	cc := t.RawGetString("cyan").(golua.LNumber)
	cm := t.RawGetString("magenta").(golua.LNumber)
	cy := t.RawGetString("yellow").(golua.LNumber)
	ck := t.RawGetString("key").(golua.LNumber)

	return uint8(cc), uint8(cm), uint8(cy), uint8(ck)
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

	case COLOR_TYPE_GRAYA16:
		cy, ca := ParseGrayA16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		return cy8, cy8, cy8, Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		return cy, cy, cy, 255

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		return cy8, cy8, cy8, 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, 0, 0, ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		ca8 := Color16BitTo8Bit(ca)
		return 0, 0, 0, ca8

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		return cr, cg, cb, ca

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		return cr, cg, cb, 255
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

	case COLOR_TYPE_GRAYA16:
		cy, ca := ParseGrayA16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		ch, cs, cv := colorconv.RGBToHSV(cy8, cy8, cy8)
		return ch, cs, cv, Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		ch, cs, cv := colorconv.RGBToHSV(cy, cy, cy)
		return ch, cs, cv, 255

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cy)

		ch, cs, cv := colorconv.RGBToHSV(cy8, cy8, cy8)
		return ch, cs, cv, 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, 0, 0, ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return 0, 0, 0, Color16BitTo8Bit(ca)

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		ch, cs, cv := colorconv.RGBToHSV(cr, cg, cb)
		return ch, cs, cv, ca

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		ch, cs, cv := colorconv.RGBToHSV(cr, cg, cb)
		return ch, cs, cv, 255
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

	case COLOR_TYPE_GRAYA16:
		cy, ca := ParseGrayA16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		ch, cs, cl := colorconv.RGBToHSL(cy8, cy8, cy8)
		return ch, cs, cl, Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		ch, cs, cl := colorconv.RGBToHSL(cy, cy, cy)
		return ch, cs, cl, 255

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		ch, cs, cl := colorconv.RGBToHSL(cy8, cy8, cy8)
		return ch, cs, cl, 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, 0, 0, ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return 0, 0, 0, Color16BitTo8Bit(ca)

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		ch, cs, cl := colorconv.RGBToHSL(cr, cg, cb)
		return ch, cs, cl, ca

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		ch, cs, cl := colorconv.RGBToHSL(cr, cg, cb)
		return ch, cs, cl, 255
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

	case COLOR_TYPE_GRAYA16:
		cy, ca := ParseGrayA16Table(t)
		return Color16BitTo8Bit(cy), Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		return cy, 255

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cy)
		return cy8, 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return 0, Color16BitTo8Bit(ca)

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y, ca

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y, 255
	}

	return 0, 0
}

func ColorTableToGrayA16(t *golua.LTable) (uint16, uint16) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y), Color8BitTo16Bit(ca)

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y), Color8BitTo16Bit(ca)

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y), Color8BitTo16Bit(ca)

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		return Color8BitTo16Bit(cy), Color8BitTo16Bit(ca)

	case COLOR_TYPE_GRAYA16:
		cy, ca := ParseGrayA16Table(t)
		return cy, ca

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		return Color8BitTo16Bit(cy), 65535

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		return cy, 65535

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, Color8BitTo16Bit(ca)

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return 0, ca

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y), Color8BitTo16Bit(ca)

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y), 65535
	}

	return 0, 0
}

func ColorTableToGray(t *golua.LTable) uint8 {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, _ := ParseRGBATable(t)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y

	case COLOR_TYPE_HSVA:
		ch, cs, cv, _ := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y

	case COLOR_TYPE_HSLA:
		ch, cs, cl, _ := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y

	case COLOR_TYPE_GRAYA:
		cy, _ := ParseGrayATable(t)
		return cy

	case COLOR_TYPE_GRAYA16:
		cy, _ := ParseGrayA16Table(t)
		return Color16BitTo8Bit(cy)

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		return cy

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		return Color16BitTo8Bit(cy)

	case COLOR_TYPE_ALPHA:
		fallthrough
	case COLOR_TYPE_ALPHA16:
		return 0

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, _ := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return g.Y
	}

	return 0
}

func ColorTableToGray16(t *golua.LTable) uint16 {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, _ := ParseRGBATable(t)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y)

	case COLOR_TYPE_HSVA:
		ch, cs, cv, _ := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y)

	case COLOR_TYPE_HSLA:
		ch, cs, cl, _ := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0
		}

		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y)

	case COLOR_TYPE_GRAYA:
		cy, _ := ParseGrayATable(t)
		return Color8BitTo16Bit(cy)

	case COLOR_TYPE_GRAYA16:
		cy, _ := ParseGrayA16Table(t)
		return cy

	case COLOR_TYPE_GRAY:
		cy := ParseGrayTable(t)
		return Color8BitTo16Bit(cy)

	case COLOR_TYPE_GRAY16:
		cy := ParseGray16Table(t)
		return cy

	case COLOR_TYPE_ALPHA:
		fallthrough
	case COLOR_TYPE_ALPHA16:
		return 0

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, _ := ParseCMYKATable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y)

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		cr, cg, cb := color.CMYKToRGB(cc, cm, cy, ck)
		g := colorconv.RGBToGrayAverage(cr, cg, cb)
		return Color8BitTo16Bit(g.Y)
	}

	return 0
}

func ColorTableToAlpha(t *golua.LTable) uint8 {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		_, _, _, ca := ParseRGBATable(t)
		return ca

	case COLOR_TYPE_HSVA:
		_, _, _, ca := ParseHSVATable(t)
		return ca

	case COLOR_TYPE_HSLA:
		_, _, _, ca := ParseHSLATable(t)
		return ca

	case COLOR_TYPE_GRAYA:
		_, ca := ParseGrayATable(t)
		return ca

	case COLOR_TYPE_GRAYA16:
		_, ca := ParseGrayA16Table(t)
		return Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		fallthrough
	case COLOR_TYPE_GRAY16:
		return 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return Color16BitTo8Bit(ca)

	case COLOR_TYPE_CMYKA:
		_, _, _, _, ca := ParseCMYKATable(t)
		return ca

	case COLOR_TYPE_CMYK:
		return 255
	}

	return 0
}

func ColorTableToAlpha16(t *golua.LTable) uint16 {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		_, _, _, ca := ParseRGBATable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_HSVA:
		_, _, _, ca := ParseHSVATable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_HSLA:
		_, _, _, ca := ParseHSLATable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_GRAYA:
		_, ca := ParseGrayATable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_GRAYA16:
		_, ca := ParseGrayA16Table(t)
		return ca

	case COLOR_TYPE_GRAY:
		fallthrough
	case COLOR_TYPE_GRAY16:
		return 65535

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return ca

	case COLOR_TYPE_CMYKA:
		_, _, _, _, ca := ParseCMYKATable(t)
		return Color8BitTo16Bit(ca)

	case COLOR_TYPE_CMYK:
		return 65535
	}

	return 0
}

func ColorTableToCMYKA(t *golua.LTable) (uint8, uint8, uint8, uint8, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, ca := ParseRGBATable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck, ca

	case COLOR_TYPE_HSVA:
		ch, cs, cv, ca := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0, 0, 0, 0
		}

		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck, ca

	case COLOR_TYPE_HSLA:
		ch, cs, cl, ca := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0, 0, 0, 0
		}

		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck, ca

	case COLOR_TYPE_GRAYA:
		cy, ca := ParseGrayATable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cy, cy, cy)
		return cc, cm, cy, ck, ca

	case COLOR_TYPE_GRAYA16:
		cyg, ca := ParseGrayA16Table(t)
		cy8 := Color16BitTo8Bit(cyg)
		cc, cm, cy, ck := color.RGBToCMYK(cy8, cy8, cy8)
		return cc, cm, cy, ck, Color16BitTo8Bit(ca)

	case COLOR_TYPE_GRAY:
		cyg := ParseGrayTable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cyg, cyg, cyg)
		return cc, cm, cy, ck, 255

	case COLOR_TYPE_GRAY16:
		cyg := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cyg)
		cc, cm, cy, ck := color.RGBToCMYK(cy8, cy8, cy8)
		return cc, cm, cy, ck, 255

	case COLOR_TYPE_ALPHA:
		ca := ParseAlphaTable(t)
		return 0, 0, 0, 0, ca

	case COLOR_TYPE_ALPHA16:
		ca := ParseAlpha16Table(t)
		return 0, 0, 0, 0, Color16BitTo8Bit(ca)

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, ca := ParseCMYKATable(t)
		return cc, cm, cy, ck, ca

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		return cc, cm, cy, ck, 255
	}

	return 0, 0, 0, 0, 0
}

func ColorTableToCMYK(t *golua.LTable) (uint8, uint8, uint8, uint8) {
	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case COLOR_TYPE_RGBA:
		cr, cg, cb, _ := ParseRGBATable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck

	case COLOR_TYPE_HSVA:
		ch, cs, cv, _ := ParseHSVATable(t)
		cr, cg, cb, err := colorconv.HSVToRGB(ch, cs, cv)
		if err != nil {
			return 0, 0, 0, 0
		}

		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck

	case COLOR_TYPE_HSLA:
		ch, cs, cl, _ := ParseHSLATable(t)
		cr, cg, cb, err := colorconv.HSLToRGB(ch, cs, cl)
		if err != nil {
			return 0, 0, 0, 0
		}

		cc, cm, cy, ck := color.RGBToCMYK(cr, cg, cb)
		return cc, cm, cy, ck

	case COLOR_TYPE_GRAYA:
		cyg, _ := ParseGrayATable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cyg, cyg, cyg)
		return cc, cm, cy, ck

	case COLOR_TYPE_GRAYA16:
		cyg, _ := ParseGrayA16Table(t)
		cy8 := Color16BitTo8Bit(cyg)
		cc, cm, cy, ck := color.RGBToCMYK(cy8, cy8, cy8)
		return cc, cm, cy, ck

	case COLOR_TYPE_GRAY:
		cyg := ParseGrayTable(t)
		cc, cm, cy, ck := color.RGBToCMYK(cyg, cyg, cyg)
		return cc, cm, cy, ck

	case COLOR_TYPE_GRAY16:
		cyg := ParseGray16Table(t)
		cy8 := Color16BitTo8Bit(cyg)
		cc, cm, cy, ck := color.RGBToCMYK(cy8, cy8, cy8)
		return cc, cm, cy, ck

	case COLOR_TYPE_ALPHA:
		fallthrough
	case COLOR_TYPE_ALPHA16:
		return 0, 0, 0, 0

	case COLOR_TYPE_CMYKA:
		cc, cm, cy, ck, _ := ParseCMYKATable(t)
		return cc, cm, cy, ck

	case COLOR_TYPE_CMYK:
		cc, cm, cy, ck := ParseCMYKTable(t)
		return cc, cm, cy, ck
	}

	return 0, 0, 0, 0
}
