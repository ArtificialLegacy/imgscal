package lib

import (
	"fmt"
	"image"
	"path"
	"strconv"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/crazy3lf/colorconv"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IMAGE = "image"

/// @lib Image
/// @import image
/// @desc
/// Library including the basic tools for handling images.
/// Also handles colors.

func RegisterImage(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_IMAGE, r, r.State, lg)

	/// @func new(name, encoding, width, height, model?) -> int<collection.IMAGE>
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg? model {int<image.ColorModel>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "new",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

					i.Self = &collection.ItemImage{
						Image:    imageutil.NewImage(args["width"].(int), args["height"].(int), model),
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Name:     args["name"].(string),
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func name(id, name)
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string} - The new name to use for the image, not including the file extension.
	lib.CreateFunction(tab, "name",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self.Name = args["name"].(string)
				},
			})
			return 0
		})

	/// @func name_ext(id, options)
	/// @arg id {int<collection.IMAGE>}
	/// @arg options {struct<image.ImageNameOptions>}
	lib.CreateFunction(tab, "name_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "options", Table: &[]lua.Arg{
				{Type: lua.STRING, Name: "name", Optional: true},
				{Type: lua.STRING, Name: "prefix", Optional: true},
				{Type: lua.STRING, Name: "suffix", Optional: true},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ImageNameOptions
			/// @prop name {string}
			/// @prop prefix {string}
			/// @prop suffix {string}

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					opt := args["options"].(map[string]any)
					newName := ""

					if opt["name"] != "" {
						newName = opt["name"].(string)
					}
					if opt["prefix"] != "" {
						newName = opt["prefix"].(string) + i.Self.Name
					}
					if opt["suffix"] != "" {
						newName += opt["suffix"].(string)
					}

					i.Self.Name = newName
					i.Lg.Append(fmt.Sprintf("new image name: %s", i.Self.Name), log.LEVEL_INFO)
				},
			})
			return 0
		})

	/// @func encoding(id, encoding)
	/// @arg id {int<collection.IMAGE>}
	/// @arg encoding {int<image.Encoding>}
	lib.CreateFunction(tab, "encoding",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self.Encoding = lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
				},
			})
			return 0
		})

	/// @func model(id) -> int<image.ColorModel>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {int<image.ColorModel>}
	/// @blocking
	lib.CreateFunction(tab, "model",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			model := 0
			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model = int(i.Self.Model)
				},
			})

			state.Push(golua.LNumber(model))
			return 1
		})

	/// @func size(id) -> int, int
	/// @arg id {int<collection.IMAGE>}
	/// @returns {int} - Image width.
	/// @returns {int} - Image height.
	/// @blocking
	lib.CreateFunction(tab, "size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := 0
			height := 0

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					b := i.Self.Image.Bounds()
					width = b.Dx()
					height = b.Dy()
				},
			})

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func crop(id, x1, y1, x2, y2)
	/// @arg id {int<collection.IMAGE>}
	/// @arg x1 {int}
	/// @arg y1 {int}
	/// @arg x2 {int}
	/// @arg y2 {int}
	/// @desc
	/// Overwrites the image.
	lib.CreateFunction(tab, "crop",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x1"},
			{Type: lua.INT, Name: "y1"},
			{Type: lua.INT, Name: "x2"},
			{Type: lua.INT, Name: "y2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x1 := args["x1"].(int) + i.Self.Image.Bounds().Min.X
					y1 := args["y1"].(int) + i.Self.Image.Bounds().Min.Y
					x2 := args["x2"].(int) + i.Self.Image.Bounds().Min.X
					y2 := args["y2"].(int) + i.Self.Image.Bounds().Min.Y

					i.Self.Image = imageutil.SubImage(
						i.Self.Image,
						x1,
						y1,
						x2,
						y2,
						true,
					)
				},
			})

			return 0
		})

	/// @func subimg(id, name, x1, y1, x2, y2, copy?) -> int<collection.IMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string} - Name for the new image.
	/// @arg x1 {int}
	/// @arg y1 {int}
	/// @arg x2 {int}
	/// @arg y2 {int}
	/// @arg? copy {bool}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "subimg",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "x1"},
			{Type: lua.INT, Name: "y1"},
			{Type: lua.INT, Name: "x2"},
			{Type: lua.INT, Name: "y2"},
			{Type: lua.BOOL, Name: "copy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var simg image.Image
			var encoding imageutil.ImageEncoding
			simgReady := make(chan struct{}, 1)

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x1 := args["x1"].(int) + i.Self.Image.Bounds().Min.X
					y1 := args["y1"].(int) + i.Self.Image.Bounds().Min.Y
					x2 := args["x2"].(int) + i.Self.Image.Bounds().Min.X
					y2 := args["y2"].(int) + i.Self.Image.Bounds().Min.Y

					simg = imageutil.SubImage(
						i.Self.Image,
						x1,
						y1,
						x2,
						y2,
						args["copy"].(bool),
					)

					encoding = i.Self.Encoding
					simgReady <- struct{}{}
				},
			})

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-simgReady
					i.Self = &collection.ItemImage{
						Image:    simg,
						Name:     name,
						Encoding: encoding,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func copy(id, name, model) -> int<collection.IMAGE>
	/// @arg id {int<collection.IMAGE>}
	/// @arg name {string}
	/// @arg model {int<image.ColorModel>} - Use -1 to maintain the color model.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "copy",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "model"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var cimg image.Image
			var encoding imageutil.ImageEncoding
			cimgReady := make(chan struct{}, 1)

			var model imageutil.ColorModel

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model = i.Self.Model
					if args["model"].(int) != -1 {
						model = lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					}

					cimg = imageutil.CopyImage(i.Self.Image, model)
					encoding = i.Self.Encoding
					cimgReady <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					i.Lg.Append("failed to copy image", log.LEVEL_ERROR)
					cimgReady <- struct{}{}
				},
			})

			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-cimgReady
					i.Self = &collection.ItemImage{
						Image:    cimg,
						Name:     name,
						Encoding: encoding,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func convert(id, model)
	/// @arg id {int<collection.IMAGE>}
	/// @arg model (int<image.ColorModel>)
	/// @desc
	/// Replaces the image inplace with a new image with the new model.
	lib.CreateFunction(tab, "convert",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "model"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					i.Self.Image = imageutil.CopyImage(i.Self.Image, model)
					i.Self.Model = model
				},
			})

			return 0
		})

	/// @func refresh(id)
	/// @arg id {int<collection.IMAGE>}
	/// @desc
	/// Shortcut for redrawing the image to guarantee the bounds of the image start at (0,0).
	/// This is sometimes needed as thirdparty libraries don't always account for non-zero min bounds.
	lib.CreateFunction(tab, "refresh",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self.Image = imageutil.CopyImage(i.Self.Image, i.Self.Model)
				},
			})

			return 0
		})

	/// @func clear(id)
	/// @arg id {int<collection.IMAGE>}
	/// @desc
	/// Resets all pixels to 0,0,0,0.
	lib.CreateFunction(tab, "clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					b := i.Self.Image.Bounds()
					iNew := image.NewRGBA(b)
					imageutil.DrawRect(i.Self.Image, iNew, b)
				},
			})

			return 0
		})

	/// @func pixel(id, x, y) -> struct<image.ColorRGBA>
	/// @arg id {int<collection.IMAGE>}
	/// @arg x {int}
	/// @arg y {int}
	/// @returns {struct<image.ColorRGBA>}
	/// @blocking
	lib.CreateFunction(tab, "pixel",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var red, green, blue, alpha uint32

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x := args["x"].(int) + i.Self.Image.Bounds().Min.X
					y := args["y"].(int) + i.Self.Image.Bounds().Min.Y

					col := i.Self.Image.At(x, y)
					red, green, blue, alpha = col.RGBA()
				},
			})

			t := imageutil.RGBAToColorTable(state, int(red), int(green), int(blue), int(alpha))
			state.Push(t)
			return 1
		})

	/// @func pixel_set(id, x, y, color)
	/// @arg id {int<collection.IMAGE>}
	/// @arg x {int}
	/// @arg y {int}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "pixel_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x := args["x"].(int) + i.Self.Image.Bounds().Min.X
					y := args["y"].(int) + i.Self.Image.Bounds().Min.Y

					red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

					imageutil.Set(
						i.Self.Image,
						x,
						y,
						int(red),
						int(green),
						int(blue),
						int(alpha),
					)
				},
			})
			return 0
		})

	/// @func point(x?, y?) -> struct<image.Point>
	/// @arg? x {int}
	/// @arg? y {int}
	/// @returns {struct<image.Point>}
	lib.CreateFunction(tab, "point",
		[]lua.Arg{
			{Type: lua.INT, Name: "x", Optional: true},
			{Type: lua.INT, Name: "y", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Point
			/// @prop x {int}
			/// @prop y {int}

			t := state.NewTable()
			state.SetTable(t, golua.LString("x"), golua.LNumber(args["x"].(int)))
			state.SetTable(t, golua.LString("y"), golua.LNumber(args["y"].(int)))

			state.Push(t)
			return 1
		})

	/// @func color_hex_to_rgba(hex) -> struct<image.ColorRGBA>
	/// @arg hex {string} - Accepts the formats RGBA, RGB, RRGGBBAA, RRGGBB, and an optional prefix of either # or 0x.
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "color_hex_to_rgba",
		[]lua.Arg{
			{Type: lua.STRING, Name: "hex"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			hex := args["hex"].(string)
			hex = strings.TrimPrefix(hex, "0x")
			hex = strings.TrimPrefix(hex, "#")

			red := 0
			green := 0
			blue := 0
			alpha := 255

			switch len(hex) {
			case 4:
				c, err := strconv.ParseInt(string(hex[3])+string(hex[3]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR)), 0)
				}
				alpha = int(c)
				fallthrough
			case 3:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[0]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s", hex), log.LEVEL_ERROR)), 0)
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[1])+string(hex[1]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR)), 0)
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[2])+string(hex[2]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on blue): %s", hex), log.LEVEL_ERROR)), 0)
				}
				blue = int(c)

			case 8:
				c, err := strconv.ParseInt(string(hex[6])+string(hex[7]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR)), 0)
				}
				alpha = int(c)
				fallthrough
			case 6:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[1]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s", hex), log.LEVEL_ERROR)), 0)
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[2])+string(hex[3]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR)), 0)
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[4])+string(hex[5]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on blue): %s", hex), log.LEVEL_ERROR)), 0)
				}
				blue = int(c)
			default:
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string: %s", hex), log.LEVEL_ERROR)), 0)
			}

			t := imageutil.RGBAToColorTable(state, red, green, blue, alpha)
			state.Push(t)
			return 1
		})

	/// @func color_to_hex(color, prefix?, lowercase?) -> string
	/// @arg color {struct<image.Color>}
	/// @arg? prefix {string} - Should be "", "#" or "0x".
	/// @arg? lowercase {bool} - Set to true to use lowercase letters in the hex string.
	/// @returns {string}
	lib.CreateFunction(tab, "color_to_hex",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
			{Type: lua.STRING, Name: "prefix", Optional: true},
			{Type: lua.BOOL, Name: "lowercase", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var hex string
			red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			if args["lowercase"].(bool) {
				hex = fmt.Sprintf("%s%02x%02x%02x%02x", args["prefix"], red, green, blue, alpha)
			} else {
				hex = fmt.Sprintf("%s%02X%02X%02X%02X", args["prefix"], red, green, blue, alpha)
			}

			state.Push(golua.LString(hex))
			return 1
		})

	/// @func color_rgb(r, g, b) -> struct<image.ColorRGBA>
	/// @arg r {int}
	/// @arg g {int}
	/// @arg b {int}
	/// @returns {struct<image.ColorRGBA>}
	/// @desc
	/// Alpha channel is set to 255.
	lib.CreateFunction(tab, "color_rgb",
		[]lua.Arg{
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Color
			/// @prop type {string<image.ColorType>}
			/// @desc
			/// Color structs are automatically converted into the needed type.
			/// Do mind that some functions may return a different type than what was passed into it.

			/// @struct ColorRGBA
			/// @prop type {string<image.ColorType>}
			/// @prop red {int}
			/// @prop green {int}
			/// @prop blue {int}
			/// @prop alpha {int}

			t := imageutil.RGBAToColorTable(state, args["r"].(int), args["g"].(int), args["b"].(int), 255)
			state.Push(t)
			return 1
		})

	/// @func color_rgba(r, g, b, a) -> struct<image.ColorRGBA>
	/// @arg r {int}
	/// @arg g {int}
	/// @arg b {int}
	/// @arg a {int}
	/// @returns struct<image.ColorRGBA>
	lib.CreateFunction(tab, "color_rgba",
		[]lua.Arg{
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageutil.RGBAToColorTable(state, args["r"].(int), args["g"].(int), args["b"].(int), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_gray(v) -> struct<image.ColorGRAYA>
	/// @arg v {int}
	/// @returns struct<image.ColorGRAYA>
	/// @desc
	/// Alpha channel is set to 255.
	lib.CreateFunction(tab, "color_gray",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorGRAYA
			/// @prop type {string<image.ColorType>}
			/// @prop gray {int}
			/// @prop alpha {int}

			t := imageutil.GrayAToColorTable(state, args["v"].(int), 255)
			state.Push(t)
			return 1
		})

	/// @func color_graya(v, a) -> struct<image.ColorGRAYA>
	/// @arg v {int}
	/// @arg a {int}
	/// @returns {struct<image.ColorGRAYA>}
	lib.CreateFunction(tab, "color_graya",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageutil.GrayAToColorTable(state, args["v"].(int), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_hsv(h, s, v) -> struct<image.ColorHSVA>
	/// @arg h {int}
	/// @arg s {int}
	/// @arg v {int}
	/// @returns struct<image.ColorHSVA>
	/// @desc
	/// Alpha channel is set to 255.
	lib.CreateFunction(tab, "color_hsv",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorHSVA
			/// @prop type {string<image.ColorType>}
			/// @prop hue {int}
			/// @prop sat {int}
			/// @prop value {int}
			/// @prop alpha {int}

			t := imageutil.HSVAToColorTable(state, args["h"].(float64), args["s"].(float64), args["v"].(float64), 255)
			state.Push(t)
			return 1
		})

	/// @func color_hsva(h, s, v, a) -> struct<image.ColorHSVA>
	/// @arg h {int}
	/// @arg s {int}
	/// @arg v {int}
	/// @arg a {int}
	/// @returns struct<image.ColorHSVA>
	lib.CreateFunction(tab, "color_hsva",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "v"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageutil.HSVAToColorTable(state, args["h"].(float64), args["s"].(float64), args["v"].(float64), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_hsl(h, s, l) -> struct<image.ColorHSLA>
	/// @arg h {int}
	/// @arg s {int}
	/// @arg l {int}
	/// @returns struct<image.ColorHSLA>
	/// @desc
	/// Alpha channel is set to 255.
	lib.CreateFunction(tab, "color_hsl",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "l"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorHSLA
			/// @prop type {string<image.ColorType>}
			/// @prop hue {int}
			/// @prop light {int}
			/// @prop value {int}
			/// @prop alpha {int}

			t := imageutil.HSLAToColorTable(state, args["h"].(float64), args["s"].(float64), args["l"].(float64), 255)
			state.Push(t)
			return 1
		})

	/// @func color_hsla(h, s, l, a) -> struct<image.ColorHSLA>
	/// @arg h {int}
	/// @arg s {int}
	/// @arg l {int}
	/// @arg a {int}
	/// @returns struct<image.ColorHSLA>
	lib.CreateFunction(tab, "color_hsla",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "l"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageutil.HSLAToColorTable(state, args["h"].(float64), args["s"].(float64), args["l"].(float64), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_to_rgb(color) -> struct<image.ColorRGBA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorRGBA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_rgb",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			t := imageutil.RGBAToColorTable(state, int(cr), int(cg), int(cb), int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_hsv(color) -> struct<image.ColorHSVA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorHSVA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_hsv",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ch, cs, cv, ca := imageutil.ColorTableToHSVA(args["color"].(*golua.LTable))
			t := imageutil.HSVAToColorTable(state, ch, cs, cv, int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_hsl(color) -> struct<image.ColorHSLA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorHSLA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_hsl",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ch, cs, cl, ca := imageutil.ColorTableToHSLA(args["color"].(*golua.LTable))
			t := imageutil.HSLAToColorTable(state, ch, cs, cl, int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_gray(color) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_gray",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cy, ca := imageutil.ColorTableToGrayA(args["color"].(*golua.LTable))
			t := imageutil.GrayAToColorTable(state, int(cy), int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_gray_average(color) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_gray_average",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayAverage(cr, cg, cb)
			t := imageutil.GrayAToColorTable(state, int(g.Y), int(ca))

			state.Push(t)
			return 1
		})

	/// @func color_to_gray_weight(color, rWeight, gWeight, bWeight) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @arg rWeight {int}
	/// @arg gWeight {int}
	/// @arg bWeight {int}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_gray_weight",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
			{Type: lua.INT, Name: "rWeight"},
			{Type: lua.INT, Name: "gWeight"},
			{Type: lua.INT, Name: "bWeight"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			rWeight := args["rWeight"].(int)
			gWeight := args["gWeight"].(int)
			bWeight := args["bWeight"].(int)

			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayWithWeight(cr, cg, cb, uint(rWeight), uint(gWeight), uint(bWeight))
			t := imageutil.GrayAToColorTable(state, int(g.Y), int(ca))

			state.Push(t)
			return 1
		})

	/// @func convert_color(model, color) -> struct<image.ColorRGBA>
	/// @arg model {int<image.ColorModel>}
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "convert_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			red, green, blue, alpha := imageutil.ConvertColor(
				lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib),
				int(cr),
				int(cg),
				int(cb),
				int(ca),
			)

			t := imageutil.RGBAToColorTable(state, red, green, blue, alpha)
			state.Push(t)
			return 1
		})

	/// @func draw(id, src, x, y, width?, height?)
	/// @arg id {int<collection.IMAGE>}
	/// @arg src {int<collection.IMAGE>} - To draw onto the base image.
	/// @arg x {int}
	/// @arg y {int}
	/// @arg? width {int}
	/// @arg? height {int}
	lib.CreateFunction(tab, "draw",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "src"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.INT, Name: "width", Optional: true},
			{Type: lua.INT, Name: "height", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["src"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady
					x := args["x"].(int) + i.Self.Image.Bounds().Min.X
					y := args["y"].(int) + i.Self.Image.Bounds().Min.Y
					width := args["width"].(int)
					height := args["height"].(int)

					if width == 0 {
						width = i.Self.Image.Bounds().Dx() - args["x"].(int)
					}
					if height == 0 {
						height = i.Self.Image.Bounds().Dy() - args["y"].(int)
					}

					imageutil.Draw(i.Self.Image, img, x, y, width, height)
					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgFinished <- struct{}{}
				},
			})

			return 0
		})

	/// @func map(id, fn, invert?)
	/// @arg id {int<collection.IMAGE>}
	/// @arg fn {function(x int, y int, color struct<image.ColorRGBA>) -> struct<image.ColorRGBA>}
	/// @arg? invert {bool} - Reverses the looping order from columns to rows.
	lib.CreateFunction(tab, "map",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "func"},
			{Type: lua.BOOL, Name: "invert", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			scheduledState, _ := state.NewThread()

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x := i.Self.Image.Bounds().Min.X
					y := i.Self.Image.Bounds().Min.Y
					width := i.Self.Image.Bounds().Dx()
					height := i.Self.Image.Bounds().Dy()

					d1Start := x
					d1End := x + width
					d2Start := y
					d2End := y + height
					invert := args["invert"].(bool)
					if invert {
						d1Start = y
						d1End = y + height
						d2Start = x
						d2End = x + width
					}

					for d1 := d1Start; d1 < d1End; d1++ {
						for d2 := d2Start; d2 < d2End; d2++ {
							ix := d1
							iy := d2
							if invert {
								ix = d2
								iy = d1
							}

							px := i.Self.Image.At(ix, iy)
							cr, cg, cb, ca := px.RGBA()

							scheduledState.Push(args["func"].(*golua.LFunction))

							scheduledState.Push(golua.LNumber(ix - x))
							scheduledState.Push(golua.LNumber(iy - y))

							t := imageutil.RGBAToColorTable(scheduledState, int(cr), int(cg), int(cb), int(ca))
							scheduledState.Push(t)
							scheduledState.Call(3, 1)
							c := scheduledState.ToTable(-1)
							scheduledState.Pop(1)

							nr, ng, nb, na := imageutil.ColorTableToRGBA(c)

							imageutil.Set(i.Self.Image, ix, iy, int(nr), int(ng), int(nb), int(na))
						}
					}

					scheduledState.Close()
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					scheduledState.Close()
				},
			})

			return 0
		})

	/// @func ext_to_encoding(ext) -> int<image.Encoding>
	/// @arg ext {string}
	/// @returns {int<image.Encoding>}
	lib.CreateFunction(tab, "ext_to_encoding",
		[]lua.Arg{
			{Type: lua.STRING, Name: "ext"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			encoding := imageutil.ExtensionEncoding(args["ext"].(string))

			state.Push(golua.LNumber(encoding))
			return 1
		})

	/// @func path_to_encoding(pth) -> int<image.Encoding>
	/// @arg pth {string}
	/// @returns {int<image.Encoding>}
	/// @desc
	/// First gets the ext from the path.
	lib.CreateFunction(tab, "path_to_encoding",
		[]lua.Arg{
			{Type: lua.STRING, Name: "pth"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ext := path.Ext(args["pth"].(string))
			encoding := imageutil.ExtensionEncoding(ext)

			state.Push(golua.LNumber(encoding))
			return 1
		})

	/// @func encoding_to_ext(encoding) -> string
	/// @arg encoding {int<image.Encoding>}
	/// @returns {string}
	lib.CreateFunction(tab, "encoding_to_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ext := imageutil.EncodingExtension(imageutil.ImageEncoding(args["encoding"].(int)))

			state.Push(golua.LString(ext))
			return 1
		})

	/// @constants Color Models
	/// @const RGBA
	/// @const RGBA64
	/// @const NRGBA
	/// @const NRGBA64
	/// @const ALPHA
	/// @const ALPHA16
	/// @const GRAY
	/// @const GRAY16
	/// @const CMYK
	tab.RawSetString("MODEL_RGBA", golua.LNumber(imageutil.MODEL_RGBA))
	tab.RawSetString("MODEL_RGBA64", golua.LNumber(imageutil.MODEL_RGBA64))
	tab.RawSetString("MODEL_NRGBA", golua.LNumber(imageutil.MODEL_NRGBA))
	tab.RawSetString("MODEL_NRGBA64", golua.LNumber(imageutil.MODEL_NRGBA64))
	tab.RawSetString("MODEL_ALPHA", golua.LNumber(imageutil.MODEL_ALPHA))
	tab.RawSetString("MODEL_ALPHA16", golua.LNumber(imageutil.MODEL_ALPHA16))
	tab.RawSetString("MODEL_GRAY", golua.LNumber(imageutil.MODEL_GRAY))
	tab.RawSetString("MODEL_GRAY16", golua.LNumber(imageutil.MODEL_GRAY16))
	tab.RawSetString("MODEL_CMYK", golua.LNumber(imageutil.MODEL_CMYK))

	/// @constants Encodings
	/// @const ENCODING_PNG
	/// @const ENCODING_JPEG
	/// @const ENCODING_GIF
	tab.RawSetString("ENCODING_PNG", golua.LNumber(imageutil.ENCODING_PNG))
	tab.RawSetString("ENCODING_JPEG", golua.LNumber(imageutil.ENCODING_JPEG))
	tab.RawSetString("ENCODING_GIF", golua.LNumber(imageutil.ENCODING_GIF))

	/// @constants Color Types
	/// @const COLOR_TYPE_RGBA
	/// @const COLOR_TYPE_HSVA
	/// @const COLOR_TYPE_HSLA
	/// @const COLOR_TYPE_GRAYA
	tab.RawSetString("COLOR_TYPE_RGBA", golua.LString(imageutil.COLOR_TYPE_RGBA))
	tab.RawSetString("COLOR_TYPE_HSVA", golua.LString(imageutil.COLOR_TYPE_HSVA))
	tab.RawSetString("COLOR_TYPE_HSLA", golua.LString(imageutil.COLOR_TYPE_HSLA))
	tab.RawSetString("COLOR_TYPE_GRAYA", golua.LString(imageutil.COLOR_TYPE_GRAYA))
}
