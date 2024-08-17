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

	/// @func new()
	/// @arg name
	/// @arg encoding
	/// @arg width
	/// @arg height
	/// @arg? model
	/// @returns id
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

	/// @func name()
	/// @arg image_id - the id of the image to rename.
	/// @arg new_name - the new name to use for the image, not including the file extension.
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

	/// @func name_ext()
	/// @arg image_id - the id of the image to rename.
	/// @arg options - a table containing each rename step. [name, prefix, suffix]
	lib.CreateFunction(tab, "name_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "options", Table: &[]lua.Arg{
				{Type: lua.STRING, Name: "name", Optional: true},
				{Type: lua.STRING, Name: "prefix", Optional: true},
				{Type: lua.STRING, Name: "suffix", Optional: true},
				{Type: lua.STRING, Name: "ext", Optional: true},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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

	/// @func encoding()
	/// @arg id
	/// @arg encoding
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

	/// @func model()
	/// @arg id
	/// @returns model
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

	/// @func size()
	/// @arg image_id - the id of the image to get the size from.
	/// @returns width
	/// @returns height
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

	/// @func crop()
	/// @arg id
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @desc
	/// overwrites the img
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

	/// @func subimg()
	/// @arg id
	/// @arg name
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @arg? copy
	/// @returns new_id
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

	/// @func copy()
	/// @arg id
	/// @arg name
	/// @arg model - use -1 to maintain color model
	/// @returns new_id
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

	/// @func convert()
	/// @arg id
	/// @arg model
	/// @desc
	/// replaces the image inplace with a new image with the new model
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

	/// @func refresh()
	/// @arg id
	/// @desc
	/// shortcut for redrawing the image to guarantee the bounds of the image start at (0,0)
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

	/// @func clear()
	/// @arg id
	/// @desc
	/// Resets all pixels to 0,0,0,0
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

	/// @func pixel()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @returns color struct [RGBA]
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

	/// @func pixel_set()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg color struct
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

	/// @func point()
	/// @arg? x
	/// @arg? y
	/// @returns {x, y}
	lib.CreateFunction(tab, "point",
		[]lua.Arg{
			{Type: lua.INT, Name: "x", Optional: true},
			{Type: lua.INT, Name: "y", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct point
			/// @prop x
			/// @prop y

			t := state.NewTable()
			state.SetTable(t, golua.LString("x"), golua.LNumber(args["x"].(int)))
			state.SetTable(t, golua.LString("y"), golua.LNumber(args["y"].(int)))

			state.Push(t)
			return 1
		})

	/// @func color_hex_to_rgba()
	/// @arg hex
	/// @returns color struct [RGBA]
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

	/// @func color_to_hex()
	/// @arg color struct
	/// @arg? prefix - should be "", "#" or "0x"
	/// @arg? lowercase - set to true to use lowercase letters in the hex string
	/// @returns hex string
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

	/// @func color_rgb()
	/// @arg r
	/// @arg g
	/// @arg b
	/// @returns color struct [RGBA]
	/// @desc
	/// alpha channel is set to 255.
	lib.CreateFunction(tab, "color_rgb",
		[]lua.Arg{
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Color
			/// @prop type
			/// @desc
			/// Color structs are automatically converted into the needed type.
			/// Do mind that some functions may return a different type than what was passed into it.

			/// @struct Color [RGBA]
			/// @prop type
			/// @prop red
			/// @prop green
			/// @prop blue
			/// @prop alpha

			t := imageutil.RGBAToColorTable(state, args["r"].(int), args["g"].(int), args["b"].(int), 255)
			state.Push(t)
			return 1
		})

	/// @func color_rgba()
	/// @arg r
	/// @arg g
	/// @arg b
	/// @arg a
	/// @returns color struct [RGBA]
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

	/// @func color_gray()
	/// @arg v
	/// @returns color struct [GRAYA]
	/// @desc
	/// alpha channel is set to 255.
	lib.CreateFunction(tab, "color_gray",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Color [GRAYA]
			/// @prop type
			/// @prop gray
			/// @prop alpha

			t := imageutil.GrayAToColorTable(state, args["v"].(int), 255)
			state.Push(t)
			return 1
		})

	/// @func color_graya()
	/// @arg v
	/// @arg a
	/// @returns color struct [GRAYA]
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

	/// @func color_hsv()
	/// @arg h
	/// @arg s
	/// @arg v
	/// @returns color struct [HSVA]
	/// @desc
	/// alpha channel is set to 255.
	lib.CreateFunction(tab, "color_hsv",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Color [HSVA]
			/// @prop type
			/// @prop hue
			/// @prop sat
			/// @prop value
			/// @prop alpha

			t := imageutil.HSVAToColorTable(state, args["h"].(float64), args["s"].(float64), args["v"].(float64), 255)
			state.Push(t)
			return 1
		})

	/// @func color_hsva()
	/// @arg h
	/// @arg s
	/// @arg v
	/// @arg a
	/// @returns color struct [HSVA]
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

	/// @func color_hsl()
	/// @arg h
	/// @arg s
	/// @arg l
	/// @returns color struct [HSLA]
	/// @desc
	/// alpha channel is set to 255.
	lib.CreateFunction(tab, "color_hsl",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "h"},
			{Type: lua.FLOAT, Name: "s"},
			{Type: lua.FLOAT, Name: "l"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Color [HSLA]
			/// @prop type
			/// @prop hue
			/// @prop light
			/// @prop value
			/// @prop alpha

			t := imageutil.HSLAToColorTable(state, args["h"].(float64), args["s"].(float64), args["l"].(float64), 255)
			state.Push(t)
			return 1
		})

	/// @func color_hsla()
	/// @arg h
	/// @arg s
	/// @arg l
	/// @arg a
	/// @returns color struct [HSLA]
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

	/// @func color_to_rgb()
	/// @arg color struct
	/// @returns color struct [RGBA]
	/// @desc
	/// alpha is maintained
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

	/// @func color_to_hsv()
	/// @arg color struct
	/// @returns color struct [HSVA]
	/// @desc
	/// alpha is maintained
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

	/// @func color_to_hsl()
	/// @arg color struct
	/// @returns color struct [HSLA]
	/// @desc
	/// alpha is maintained
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

	/// @func color_to_gray()
	/// @arg color struct
	/// @returns color struct [GRAYA]
	/// @desc
	/// alpha is maintained
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

	/// @func color_to_gray_average()
	/// @arg color struct
	/// @returns color struct [GRAYA]
	/// @desc
	/// alpha is maintained
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

	/// @func color_to_gray_weight()
	/// @arg color struct
	/// @arg rWeight
	/// @arg gWeight
	/// @arg bWeight
	/// @returns color struct [GRAYA]
	/// @desc
	/// alpha is maintained
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

	/// @func convert_color()
	/// @arg model
	/// @arg color struct
	/// @returns color struct [RGBA]
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

	/// @func draw()
	/// @arg id
	/// @arg id - to draw onto the base image
	/// @arg x
	/// @arg y
	/// @arg? width
	/// @arg? height
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

	/// @func map()
	/// @arg id
	/// @arg fn - takes in x, y , {red, green, blue, alpha} and returns a new {red, green, blue, alpha}
	/// @arg? invert - reverses the looping order from columns to rows
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

	/// @func ext_to_encoding()
	/// @arg ext
	/// @returns encoding
	lib.CreateFunction(tab, "ext_to_encoding",
		[]lua.Arg{
			{Type: lua.STRING, Name: "ext"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			encoding := imageutil.ExtensionEncoding(args["ext"].(string))

			state.Push(golua.LNumber(encoding))
			return 1
		})

	/// @func path_to_encoding()
	/// @arg pth
	/// @returns encoding
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

	/// @func encoding_to_ext()
	/// @arg encoding
	/// @returns ext
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
