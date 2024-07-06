package lib

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IMAGE = "image"

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

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
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

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
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

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
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

	/// @func pixel()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @returns {red, green, blue, alpha}
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

			t := state.NewTable()
			state.SetField(t, "red", golua.LNumber(red))
			state.SetField(t, "green", golua.LNumber(green))
			state.SetField(t, "blue", golua.LNumber(blue))
			state.SetField(t, "alpha", golua.LNumber(alpha))
			state.Push(t)
			return 4
		})

	/// @func pixel_set()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg {red, green, blue, alpha}
	lib.CreateFunction(tab, "pixel_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					color := args["color"].(map[string]any)
					x := args["x"].(int) + i.Self.Image.Bounds().Min.X
					y := args["y"].(int) + i.Self.Image.Bounds().Min.Y

					imageutil.Set(
						i.Self.Image,
						x,
						y,
						color["red"].(int),
						color["green"].(int),
						color["blue"].(int),
						color["alpha"].(int),
					)
				},
			})
			return 0
		})

	/// @func color_hex()
	/// @arg hex
	/// @returns {red, green, blue, alpha}
	lib.CreateFunction(tab, "color_hex",
		[]lua.Arg{
			{Type: lua.STRING, Name: "hex"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			hex := args["hex"].(string)
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

			t := state.NewTable()
			state.SetField(t, "red", golua.LNumber(red))
			state.SetField(t, "green", golua.LNumber(green))
			state.SetField(t, "blue", golua.LNumber(blue))
			state.SetField(t, "alpha", golua.LNumber(alpha))
			state.Push(t)
			return 1
		})

	/// @func convert_color()
	/// @arg model
	/// @arg color {red, green, blue, alpha}
	/// @returns new color {red, green, blue, alpha}
	lib.CreateFunction(tab, "convert_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			color := args["color"].(map[string]any)
			red, green, blue, alpha := imageutil.ConvertColor(
				lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib),
				color["red"].(int),
				color["green"].(int),
				color["blue"].(int),
				color["alpha"].(int),
			)

			t := state.NewTable()
			state.SetField(t, "red", golua.LNumber(red))
			state.SetField(t, "green", golua.LNumber(green))
			state.SetField(t, "blue", golua.LNumber(blue))
			state.SetField(t, "alpha", golua.LNumber(alpha))
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
	/// @blocking
	lib.CreateFunction(tab, "map",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					x := i.Self.Image.Bounds().Min.X
					y := i.Self.Image.Bounds().Min.Y
					width := i.Self.Image.Bounds().Dx()
					height := i.Self.Image.Bounds().Dy()

					for ix := x; ix < x+width; ix++ {
						for iy := y; iy < y+height; iy++ {
							px := i.Self.Image.At(ix, iy)
							cr, cg, cb, ca := px.RGBA()

							state.Push(args["func"].(*golua.LFunction))

							state.Push(golua.LNumber(ix - x))
							state.Push(golua.LNumber(iy - y))

							t := state.NewTable()
							state.SetField(t, "red", golua.LNumber(cr))
							state.SetField(t, "green", golua.LNumber(cg))
							state.SetField(t, "blue", golua.LNumber(cb))
							state.SetField(t, "alpha", golua.LNumber(ca))
							state.Push(t)

							state.Call(3, 1)

							c := state.ToTable(-1)

							nr := state.GetField(c, "red")
							if nr.Type() != golua.LTNumber {
								state.Error(golua.LString(lg.Append("invalid red field returned into image.map", log.LEVEL_ERROR)), 0)
							}
							ng := state.GetField(c, "green")
							if ng.Type() != golua.LTNumber {
								state.Error(golua.LString(lg.Append("invalid green field returned into image.map", log.LEVEL_ERROR)), 0)
							}
							nb := state.GetField(c, "blue")
							if nb.Type() != golua.LTNumber {
								state.Error(golua.LString(lg.Append("invalid blue field returned into image.map", log.LEVEL_ERROR)), 0)
							}
							na := state.GetField(c, "alpha")
							if na.Type() != golua.LTNumber {
								state.Error(golua.LString(lg.Append("invalid alpha field returned into image.map", log.LEVEL_ERROR)), 0)
							}

							imageutil.Set(i.Self.Image, ix, iy, int(nr.(golua.LNumber)), int(ng.(golua.LNumber)), int(nb.(golua.LNumber)), int(na.(golua.LNumber)))
						}
					}
				},
			})

			return 0
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
	r.State.SetField(tab, "MODEL_RGBA", golua.LNumber(imageutil.MODEL_RGBA))
	r.State.SetField(tab, "MODEL_RGBA64", golua.LNumber(imageutil.MODEL_RGBA64))
	r.State.SetField(tab, "MODEL_NRGBA", golua.LNumber(imageutil.MODEL_NRGBA))
	r.State.SetField(tab, "MODEL_NRGBA64", golua.LNumber(imageutil.MODEL_NRGBA64))
	r.State.SetField(tab, "MODEL_ALPHA", golua.LNumber(imageutil.MODEL_ALPHA))
	r.State.SetField(tab, "MODEL_ALPHA16", golua.LNumber(imageutil.MODEL_ALPHA16))
	r.State.SetField(tab, "MODEL_GRAY", golua.LNumber(imageutil.MODEL_GRAY))
	r.State.SetField(tab, "MODEL_GRAY16", golua.LNumber(imageutil.MODEL_GRAY16))
	r.State.SetField(tab, "MODEL_CMYK", golua.LNumber(imageutil.MODEL_CMYK))

	/// @constants Encodings
	r.State.SetField(tab, "ENCODING_PNG", golua.LNumber(imageutil.ENCODING_PNG))
	r.State.SetField(tab, "ENCODING_JPEG", golua.LNumber(imageutil.ENCODING_JPEG))
	r.State.SetField(tab, "ENCODING_GIF", golua.LNumber(imageutil.ENCODING_GIF))
}
