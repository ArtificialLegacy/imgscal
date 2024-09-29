package lib

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"path"
	"strconv"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/crazy3lf/colorconv"
	color_extractor "github.com/marekm4/color-extractor"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IMAGE = "image"

/// @lib Image
/// @import image
/// @desc
/// Library including the basic tools for handling images and colors.

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

	/// @func new_filled(name, encoding, width, height, color, model?) -> int<collection.IMAGE>
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg color {struct<image.Color>}
	/// @arg? model {int<image.ColorModel>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "new_filled",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "color"},
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

					width := args["width"].(int)
					height := args["height"].(int)

					img := imageutil.NewImage(width, height, model)
					red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

					for ix := 0; ix < width; ix++ {
						for iy := 0; iy < height; iy++ {
							imageutil.Set(img, ix, iy, int(red), int(green), int(blue), int(alpha))
						}
					}

					i.Self = &collection.ItemImage{
						Image:    img,
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Name:     name,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func new_random(name, encoding, width, height, enableAlpha?, model?) -> int<collection.IMAGE>
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg? enableAlpha {bool}
	/// @arg? model {int<image.ColorModel>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "new_random",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.BOOL, Name: "enableAlpha", Optional: true},
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

					width := args["width"].(int)
					height := args["height"].(int)

					img := imageutil.NewImage(width, height, model)

					for ix := 0; ix < width; ix++ {
						for iy := 0; iy < height; iy++ {
							red := rand.Intn(256)
							green := rand.Intn(256)
							blue := rand.Intn(256)
							alpha := 255
							if args["enableAlpha"].(bool) {
								alpha = rand.Intn(256)
							}
							imageutil.Set(img, ix, iy, int(red), int(green), int(blue), int(alpha))
						}
					}

					i.Self = &collection.ItemImage{
						Image:    img,
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Name:     name,
						Model:    model,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func remove(id)
	/// @arg id {int<collection.IMAGE>}
	/// @desc
	/// Removes the image from the collection.
	/// This is a shortcut for collection.collect(collection.IMAGE, id).
	lib.CreateFunction(tab, "remove",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Collect(args["id"].(int))
			return 0
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

	/// @func clone(id, src, model)
	/// @arg id {int<collection.IMAGE>}
	/// @arg src {int<collection.IMAGE>}
	/// @arg model {int<image.ColorModel>} - Use -1 to maintain the color model.
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Clones image src to id with the new model.
	lib.CreateFunction(tab, "clone",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "src"},
			{Type: lua.INT, Name: "model"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var cimg image.Image
			cimgReady := make(chan struct{}, 1)
			var model imageutil.ColorModel

			r.IC.Schedule(args["src"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model = i.Self.Model
					if args["model"].(int) != -1 {
						model = lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)
					}

					cimg = imageutil.CopyImage(i.Self.Image, model)
					cimgReady <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					i.Lg.Append("failed to copy image", log.LEVEL_ERROR)
					cimgReady <- struct{}{}
				},
			})

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-cimgReady
					i.Self.Image = cimg
					i.Self.Model = model
				},
			})

			return 0
		})

	/// @func convert(id, model)
	/// @arg id {int<collection.IMAGE>}
	/// @arg model {int<image.ColorModel>}
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

	/// @func alpha_set(id, alpha)
	/// @arg id {int<collection.IMAGE>}
	/// @arg alpha {int}
	/// @desc
	/// This has no effect on gray, gray16, or cmyk images.
	lib.CreateFunction(tab, "alpha_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "alpha"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					alpha := uint8(args["alpha"].(int))
					imageutil.AlphaSet(i.Self.Image, alpha)
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
			{Type: lua.RAW_TABLE, Name: "color"},
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
	/// @arg hex {string} - Accepts the formats RGBA, RGB, RRGGBBAA, RRGGBB, and an optional prefix of "#", "$", or "0x".
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "color_hex_to_rgba",
		[]lua.Arg{
			{Type: lua.STRING, Name: "hex"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			hex := args["hex"].(string)
			hex = strings.TrimPrefix(hex, "0x")
			hex = strings.TrimPrefix(hex, "#")
			hex = strings.TrimPrefix(hex, "$")

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

	/// @func color_hex_to_rgba_bgr(hex) -> struct<image.ColorRGBA>
	/// @arg hex {string} - Accepts the formats ABGR, BGR, AABBGGRR, BBGGRR, and an optional prefix of "#", "$", or "0x".
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "color_hex_to_rgba_bgr",
		[]lua.Arg{
			{Type: lua.STRING, Name: "hex"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			hex := args["hex"].(string)
			hex = strings.TrimPrefix(hex, "0x")
			hex = strings.TrimPrefix(hex, "#")
			hex = strings.TrimPrefix(hex, "$")

			red := 0
			green := 0
			blue := 0
			alpha := 255

			offset := 0

			switch len(hex) {
			case 4:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[0]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR)), 0)
				}
				alpha = int(c)
				offset = 1
				fallthrough
			case 3:
				c, err := strconv.ParseInt(string(hex[2+offset])+string(hex[2+offset]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s", hex), log.LEVEL_ERROR)), 0)
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[1+offset])+string(hex[1+offset]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR)), 0)
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[offset])+string(hex[offset]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on blue): %s", hex), log.LEVEL_ERROR)), 0)
				}
				blue = int(c)

			case 8:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[1]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR)), 0)
				}
				alpha = int(c)
				offset = 2
				fallthrough
			case 6:
				c, err := strconv.ParseInt(string(hex[4+offset])+string(hex[5+offset]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s", hex), log.LEVEL_ERROR)), 0)
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[2+offset])+string(hex[3+offset]), 16, 64)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR)), 0)
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[offset])+string(hex[offset]), 16, 64)
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

	/// @func color_to_hex(color, noalpha?, prefix?, lowercase?) -> string
	/// @arg color {struct<image.Color>}
	/// @arg? noalpha {bool} - Set to true to exclude the alpha channel.
	/// @arg? prefix {string} - Should be "", "#", '$', or "0x".
	/// @arg? lowercase {bool} - Set to true to use lowercase letters in the hex string.
	/// @returns {string}
	/// @desc
	/// In the format RRGGBBAA or RRGGBB.
	lib.CreateFunction(tab, "color_to_hex",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.BOOL, Name: "noalpha", Optional: true},
			{Type: lua.STRING, Name: "prefix", Optional: true},
			{Type: lua.BOOL, Name: "lowercase", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var hex string
			red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			var redString string
			var greenString string
			var blueString string
			var alphaString string

			prefix := args["prefix"].(string)

			if args["lowercase"].(bool) {
				redString = fmt.Sprintf("%02x", red)
				greenString = fmt.Sprintf("%02x", green)
				blueString = fmt.Sprintf("%02x", blue)
				alphaString = fmt.Sprintf("%02x", alpha)
			} else {
				redString = fmt.Sprintf("%02X", red)
				greenString = fmt.Sprintf("%02X", green)
				blueString = fmt.Sprintf("%02X", blue)
				alphaString = fmt.Sprintf("%02X", alpha)
			}

			if args["noalpha"].(bool) {
				hex = fmt.Sprintf("%s%s%s%s", prefix, redString, greenString, blueString)
			} else {
				hex = fmt.Sprintf("%s%s%s%s%s", prefix, redString, greenString, blueString, alphaString)
			}

			state.Push(golua.LString(hex))
			return 1
		})

	/// @func color_to_hex_bgr(color, noalpha?, prefix?, lowercase?) -> string
	/// @arg color {struct<image.Color>}
	/// @arg? noalpha {bool} - Set to true to exclude the alpha channel.
	/// @arg? prefix {string} - Should be "", "#", '$', or "0x".
	/// @arg? lowercase {bool} - Set to true to use lowercase letters in the hex string.
	/// @returns {string}
	/// @desc
	/// In the format AABBGGRR or BBGGRR.
	lib.CreateFunction(tab, "color_to_hex_bgr",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.BOOL, Name: "noalpha", Optional: true},
			{Type: lua.STRING, Name: "prefix", Optional: true},
			{Type: lua.BOOL, Name: "lowercase", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var hex string
			red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			var redString string
			var greenString string
			var blueString string
			var alphaString string

			prefix := args["prefix"].(string)

			if args["lowercase"].(bool) {
				redString = fmt.Sprintf("%02x", red)
				greenString = fmt.Sprintf("%02x", green)
				blueString = fmt.Sprintf("%02x", blue)
				alphaString = fmt.Sprintf("%02x", alpha)
			} else {
				redString = fmt.Sprintf("%02X", red)
				greenString = fmt.Sprintf("%02X", green)
				blueString = fmt.Sprintf("%02X", blue)
				alphaString = fmt.Sprintf("%02X", alpha)
			}

			if args["noalpha"].(bool) {
				hex = fmt.Sprintf("%s%s%s%s", prefix, blueString, greenString, redString)
			} else {
				hex = fmt.Sprintf("%s%s%s%s%s", prefix, alphaString, blueString, greenString, redString)
			}

			state.Push(golua.LString(hex))
			return 1
		})

	/// @func color_8bit_to_16bit(c) -> int
	/// @arg c {int}
	/// @returns {int}
	lib.CreateFunction(tab, "color_8bit_to_16bit",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := args["c"].(int)
			state.Push(golua.LNumber(imageutil.Color8BitTo16Bit(uint8(c))))
			return 1
		})

	/// @func color_16bit_to_8bit(c) -> int
	/// @arg c {int}
	/// @returns {int}
	lib.CreateFunction(tab, "color_16bit_to_8bit",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := args["c"].(int)
			state.Push(golua.LNumber(imageutil.Color16BitTo8Bit(uint16(c))))
			return 1
		})

	/// @func color_24bit_to_rgba(c) -> struct<image.ColorRGBA>
	/// @arg c {int}
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "color_24bit_to_rgba",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := args["c"].(int)

			red := c & 0xFF
			green := (c >> 8) & 0xFF
			blue := (c >> 16) & 0xFF
			alpha := 255

			t := imageutil.RGBAToColorTable(state, red, green, blue, alpha)
			state.Push(t)
			return 1
		})

	/// @func color_32bit_to_rgba(c) -> struct<image.ColorRGBA>
	/// @arg c {int}
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "color_32bit_to_rgba",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := args["c"].(int)

			red := c & 0xFF
			green := (c >> 8) & 0xFF
			blue := (c >> 16) & 0xFF
			alpha := (c >> 24) & 0xFF

			t := imageutil.RGBAToColorTable(state, red, green, blue, alpha)
			state.Push(t)
			return 1
		})

	/// @func color_to_24bit(color) -> int
	/// @arg color {struct<image.Color>}
	/// @returns {int}
	lib.CreateFunction(tab, "color_to_24bit",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			red, green, blue, _ := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			c := int(red) | (int(green) << 8) | (int(blue) << 16)
			state.Push(golua.LNumber(c))
			return 1
		})

	/// @func color_to_32bit(color) -> int
	/// @arg color {struct<image.Color>}
	/// @returns {int}
	lib.CreateFunction(tab, "color_to_32bit",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			red, green, blue, alpha := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			c := int(red) | (int(green) << 8) | (int(blue) << 16) | (int(alpha) << 24)
			state.Push(golua.LNumber(c))
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

	/// @func color_alpha(a) -> struct<image.ColorALPHA>
	/// @arg a {int}
	/// @returns {struct<image.ColorALPHA>}
	lib.CreateFunction(tab, "color_alpha",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorALPHA
			/// @prop type {string<image.ColorType>}
			/// @prop alpha {int}

			t := imageutil.AlphaToColorTable(state, args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_alpha16(a) -> struct<image.ColorALPHA16>
	/// @arg a {int}
	/// @returns {struct<image.ColorALPHA16>}
	lib.CreateFunction(tab, "color_alpha16",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorALPHA16
			/// @prop type {string<image.ColorType>}
			/// @prop alpha {int}

			t := imageutil.Alpha16ToColorTable(state, args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_gray(v) -> struct<image.ColorGRAY>
	/// @arg v {int}
	/// @returns struct<image.ColorGRAY>
	lib.CreateFunction(tab, "color_gray",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorGRAY
			/// @prop type {string<image.ColorType>}
			/// @prop gray {int}

			t := imageutil.GrayToColorTable(state, args["v"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_gray16(v) -> struct<image.ColorGRAY16>
	/// @arg v {int}
	/// @returns {struct<image.ColorGRAY16>}
	lib.CreateFunction(tab, "color_gray16",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorGRAY16
			/// @prop type {string<image.ColorType>}
			/// @prop gray {int}

			t := imageutil.Gray16ToColorTable(state, args["v"].(int))
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
			/// @struct ColorGRAYA
			/// @prop type {string<image.ColorType>}
			/// @prop gray {int}
			/// @prop alpha {int}

			t := imageutil.GrayAToColorTable(state, args["v"].(int), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_graya16(v, a) -> struct<image.ColorGRAYA16>
	/// @arg v {int}
	/// @arg a {int}
	/// @returns {struct<image.ColorGRAYA16>}
	lib.CreateFunction(tab, "color_graya16",
		[]lua.Arg{
			{Type: lua.INT, Name: "v"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorGRAYA16
			/// @prop type {string<image.ColorType>}
			/// @prop gray {int}
			/// @prop alpha {int}

			t := imageutil.GrayA16ToColorTable(state, args["v"].(int), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_cmyk(c, m, y, k) -> struct<image.ColorCMYK>
	/// @arg c {int}
	/// @arg m {int}
	/// @arg y {int}
	/// @arg k {int}
	/// @returns {struct<image.ColorCMYK>}
	lib.CreateFunction(tab, "color_cmyk",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
			{Type: lua.INT, Name: "m"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.INT, Name: "k"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorCMYK
			/// @prop type {string<image.ColorType>}
			/// @prop cyan {int}
			/// @prop magenta {int}
			/// @prop yellow {int}
			/// @prop key {int}

			t := imageutil.CMYKToColorTable(state, args["c"].(int), args["m"].(int), args["y"].(int), args["k"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_cmyka(c, m, y, k, a) -> struct<image.ColorCMYKA>
	/// @arg c {int}
	/// @arg m {int}
	/// @arg y {int}
	/// @arg k {int}
	/// @arg a {int}
	/// @returns {struct<image.ColorCMYKA>}
	lib.CreateFunction(tab, "color_cmyka",
		[]lua.Arg{
			{Type: lua.INT, Name: "c"},
			{Type: lua.INT, Name: "m"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.INT, Name: "k"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorCMYKA
			/// @prop type {string<image.ColorType>}
			/// @prop cyan {int}
			/// @prop magenta {int}
			/// @prop yellow {int}
			/// @prop key {int}
			/// @prop alpha {int}

			t := imageutil.CMYKAToColorTable(state, args["c"].(int), args["m"].(int), args["y"].(int), args["k"].(int), args["a"].(int))
			state.Push(t)
			return 1
		})

	/// @func color_zero() -> struct<image.ColorZERO>
	/// @returns {struct<image.ColorZERO>}
	lib.CreateFunction(tab, "color_zero",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ColorZERO
			/// @prop type {string<image.ColorType>}

			t := imageutil.ZeroColorTable(state)
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
			{Type: lua.RAW_TABLE, Name: "color"},
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
			{Type: lua.RAW_TABLE, Name: "color"},
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
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ch, cs, cl, ca := imageutil.ColorTableToHSLA(args["color"].(*golua.LTable))
			t := imageutil.HSLAToColorTable(state, ch, cs, cl, int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_gray(color) -> struct<image.ColorGRAY>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAY>}
	lib.CreateFunction(tab, "color_to_gray",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cy := imageutil.ColorTableToGray(args["color"].(*golua.LTable))
			t := imageutil.GrayToColorTable(state, int(cy))
			state.Push(t)
			return 1
		})

	/// @func color_to_gray_average(color) -> struct<image.ColorGRAY>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAY>}
	lib.CreateFunction(tab, "color_to_gray_average",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, _ := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayAverage(cr, cg, cb)
			t := imageutil.GrayToColorTable(state, int(g.Y))

			state.Push(t)
			return 1
		})

	/// @func color_to_gray_weight(color, rWeight, gWeight, bWeight) -> struct<image.ColorGRAY>
	/// @arg color {struct<image.Color>}
	/// @arg rWeight {int}
	/// @arg gWeight {int}
	/// @arg bWeight {int}
	/// @returns {struct<image.ColorGRAY>}
	lib.CreateFunction(tab, "color_to_gray_weight",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.INT, Name: "rWeight"},
			{Type: lua.INT, Name: "gWeight"},
			{Type: lua.INT, Name: "bWeight"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			rWeight := args["rWeight"].(int)
			gWeight := args["gWeight"].(int)
			bWeight := args["bWeight"].(int)

			cr, cg, cb, _ := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayWithWeight(cr, cg, cb, uint(rWeight), uint(gWeight), uint(bWeight))
			t := imageutil.GrayToColorTable(state, int(g.Y))

			state.Push(t)
			return 1
		})

	/// @func color_to_graya(color) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cy, ca := imageutil.ColorTableToGrayA(args["color"].(*golua.LTable))
			t := imageutil.GrayAToColorTable(state, int(cy), int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_graya_average(color) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya_average",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayAverage(cr, cg, cb)
			t := imageutil.GrayAToColorTable(state, int(g.Y), int(ca))

			state.Push(t)
			return 1
		})

	/// @func color_to_graya_weight(color, rWeight, gWeight, bWeight) -> struct<image.ColorGRAYA>
	/// @arg color {struct<image.Color>}
	/// @arg rWeight {int}
	/// @arg gWeight {int}
	/// @arg bWeight {int}
	/// @returns {struct<image.ColorGRAYA>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya_weight",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
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

	/// @func color_to_graya16(color) -> struct<image.ColorGRAYA16>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA16>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya16",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cy, ca := imageutil.ColorTableToGrayA16(args["color"].(*golua.LTable))
			t := imageutil.GrayA16ToColorTable(state, int(cy), int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_graya16_average(color) -> struct<image.ColorGRAYA16>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorGRAYA16>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya16_average",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cr, cg, cb, ca := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))
			g := colorconv.RGBToGrayAverage(cr, cg, cb)
			t := imageutil.GrayA16ToColorTable(state, int(imageutil.Color8BitTo16Bit(g.Y)), int(imageutil.Color8BitTo16Bit(ca)))

			state.Push(t)
			return 1
		})

	/// @func color_to_graya16_weight(color, rWeight, gWeight, bWeight) -> struct<image.ColorGRAYA16>
	/// @arg color {struct<image.Color>}
	/// @arg rWeight {int}
	/// @arg gWeight {int}
	/// @arg bWeight {int}
	/// @returns {struct<image.ColorGRAYA16>}
	/// @desc
	/// Alpha is maintained.
	lib.CreateFunction(tab, "color_to_graya16_weight",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
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
			t := imageutil.GrayA16ToColorTable(state, int(imageutil.Color8BitTo16Bit(g.Y)), int(imageutil.Color8BitTo16Bit(ca)))

			state.Push(t)
			return 1
		})

	/// @func color_to_alpha(color) -> struct<image.ColorALPHA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorALPHA>}
	lib.CreateFunction(tab, "color_to_alpha",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ca := imageutil.ColorTableToAlpha(args["color"].(*golua.LTable))
			t := imageutil.AlphaToColorTable(state, int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_alpha16(color) -> struct<image.ColorALPHA16>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorALPHA16>}
	lib.CreateFunction(tab, "color_to_alpha16",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ca := imageutil.ColorTableToAlpha16(args["color"].(*golua.LTable))
			t := imageutil.Alpha16ToColorTable(state, int(ca))
			state.Push(t)
			return 1
		})

	/// @func color_to_cmyk(color) -> struct<image.ColorCMYK>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorCMYK>}
	lib.CreateFunction(tab, "color_to_cmyk",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cc, cm, cy, ck := imageutil.ColorTableToCMYK(args["color"].(*golua.LTable))
			t := imageutil.CMYKToColorTable(state, int(cc), int(cm), int(cy), int(ck))
			state.Push(t)
			return 1
		})

	/// @func color_to_cmyka(color) -> struct<image.ColorCMYKA>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<image.ColorCMYKA>}
	lib.CreateFunction(tab, "color_to_cmyka",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cc, cm, cy, ck, ca := imageutil.ColorTableToCMYKA(args["color"].(*golua.LTable))
			t := imageutil.CMYKAToColorTable(state, int(cc), int(cm), int(cy), int(ck), int(ca))
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
			{Type: lua.RAW_TABLE, Name: "color"},
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

	/// @func compare(id1, id2) -> bool
	/// @arg id1 {int<collection.IMAGE>}
	/// @arg id2 {int<collection.IMAGE>}
	/// @returns {bool}
	/// @blocking
	/// @desc
	/// Compares two images pixel by pixel.
	/// Early returns if the image ids are the same, without scheduling any tasks.
	lib.CreateFunction(tab, "compare",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id1 := args["id1"].(int)
			id2 := args["id2"].(int)
			if id1 == id2 {
				state.Push(golua.LTrue)
				return 1
			}

			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			var img image.Image
			var equal bool

			r.IC.Schedule(id1, &collection.Task[collection.ItemImage]{
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

			<-r.IC.Schedule(id2, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady
					equal = imageutil.ImageCompare(img, i.Self.Image)
					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgFinished <- struct{}{}
				},
			})

			state.Push(golua.LBool(equal))
			return 1
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

	/// @func extract_colors(img) -> []struct<image.ColorRGBA>
	/// @arg img {int<collection.IMAGE>}
	/// @returns {[]struct<image.ColorRGBA>}
	/// @blocking
	lib.CreateFunction(tab, "extract_colors",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image

			<-r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})

			colors := color_extractor.ExtractColors(img)

			t := state.NewTable()

			for _, c := range colors {
				cr := c.(color.RGBA)
				t.Append(imageutil.RGBAColorToColorTable(state, &cr))
			}

			state.Push(t)
			return 1
		})

	/// @func extract_colors_config(img, downSizeTo, smallBucket) -> []struct<image.ColorRGBA>
	/// @arg img {int<collection.IMAGE>}
	/// @arg downSizeTo {float}
	/// @arg smallBucket {float}
	/// @returns {[]struct<image.ColorRGBA>}
	/// @blocking
	lib.CreateFunction(tab, "extract_colors_config",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image

			<-r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})

			downSizeTo := args["downSizeTo"].(float64)
			smallBucket := args["smallBucket"].(float64)

			colors := color_extractor.ExtractColorsWithConfig(img, color_extractor.Config{
				DownSizeTo:  downSizeTo,
				SmallBucket: smallBucket,
			})

			t := state.NewTable()

			for _, c := range colors {
				cr := c.(color.RGBA)
				t.Append(imageutil.RGBAColorToColorTable(state, &cr))
			}

			state.Push(t)
			return 1
		})

	/// @func png_data_chunk(key, data) -> struct<image.PNGDataChunk>
	/// @arg key {string}
	/// @arg data {string}
	/// @returns {struct<image.PNGDataChunk>}
	lib.CreateFunction(tab, "png_data_chunk",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			{Type: lua.STRING, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct PNGDataChunk
			/// @prop key {string}
			/// @prop data {string}

			t := state.NewTable()
			t.RawSetString("key", golua.LString(args["key"].(string)))
			t.RawSetString("data", golua.LString(args["data"].(string)))

			state.Push(t)
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
	/// @const ENCODING_TIFF
	/// @const ENCODING_BMP
	/// @const ENCODING_ICO
	/// @const ENCODING_CUR
	tab.RawSetString("ENCODING_PNG", golua.LNumber(imageutil.ENCODING_PNG))
	tab.RawSetString("ENCODING_JPEG", golua.LNumber(imageutil.ENCODING_JPEG))
	tab.RawSetString("ENCODING_GIF", golua.LNumber(imageutil.ENCODING_GIF))
	tab.RawSetString("ENCODING_TIFF", golua.LNumber(imageutil.ENCODING_TIFF))
	tab.RawSetString("ENCODING_BMP", golua.LNumber(imageutil.ENCODING_BMP))
	tab.RawSetString("ENCODING_ICO", golua.LNumber(imageutil.ENCODING_ICO))
	tab.RawSetString("ENCODING_CUR", golua.LNumber(imageutil.ENCODING_CUR))

	/// @constants Color Types
	/// @const COLOR_TYPE_RGBA
	/// @const COLOR_TYPE_HSVA
	/// @const COLOR_TYPE_HSLA
	/// @const COLOR_TYPE_GRAY
	/// @const COLOR_TYPE_GRAY16
	/// @const COLOR_TYPE_GRAYA
	/// @const COLOR_TYPE_GRAYA16
	/// @const COLOR_TYPE_ALPHA
	/// @const COLOR_TYPE_ALPHA16
	/// @const COLOR_TYPE_CMYK
	/// @const COLOR_TYPE_CMYKA
	/// @const COLOR_TYPE_ZERO
	tab.RawSetString("COLOR_TYPE_RGBA", golua.LString(imageutil.COLOR_TYPE_RGBA))
	tab.RawSetString("COLOR_TYPE_HSVA", golua.LString(imageutil.COLOR_TYPE_HSVA))
	tab.RawSetString("COLOR_TYPE_HSLA", golua.LString(imageutil.COLOR_TYPE_HSLA))
	tab.RawSetString("COLOR_TYPE_GRAY", golua.LString(imageutil.COLOR_TYPE_GRAY))
	tab.RawSetString("COLOR_TYPE_GRAY16", golua.LString(imageutil.COLOR_TYPE_GRAY16))
	tab.RawSetString("COLOR_TYPE_GRAYA", golua.LString(imageutil.COLOR_TYPE_GRAYA))
	tab.RawSetString("COLOR_TYPE_GRAYA16", golua.LString(imageutil.COLOR_TYPE_GRAYA16))
	tab.RawSetString("COLOR_TYPE_ALPHA", golua.LString(imageutil.COLOR_TYPE_ALPHA))
	tab.RawSetString("COLOR_TYPE_ALPHA16", golua.LString(imageutil.COLOR_TYPE_ALPHA16))
	tab.RawSetString("COLOR_TYPE_CMYK", golua.LString(imageutil.COLOR_TYPE_CMYK))
	tab.RawSetString("COLOR_TYPE_CMYKA", golua.LString(imageutil.COLOR_TYPE_CMYKA))
	tab.RawSetString("COLOR_TYPE_ZERO", golua.LString(imageutil.COLOR_TYPE_ZERO))
}
