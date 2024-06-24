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
)

const LIB_IMAGE = "image"

func RegisterImage(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_IMAGE, r.State, lg)

	/// @func new()
	/// @arg name
	/// @arg encoding
	/// @arg width
	/// @arg height
	/// @arg model
	/// @returns id
	lib.CreateFunction("new",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name))
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					upLeft := image.Point{0, 0}
					lowRight := image.Point{args["width"].(int), args["height"].(int)}
					rect := image.Rectangle{upLeft, lowRight}

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

					var img image.Image

					switch model {
					case imageutil.MODEL_RGBA:
						img = image.NewRGBA(rect)
					case imageutil.MODEL_RGBA64:
						img = image.NewRGBA64(rect)
					case imageutil.MODEL_NRGBA:
						img = image.NewNRGBA(rect)
					case imageutil.MODEL_NRGBA64:
						img = image.NewNRGBA64(rect)
					case imageutil.MODEL_ALPHA:
						img = image.NewAlpha(rect)
					case imageutil.MODEL_ALPHA16:
						img = image.NewAlpha16(rect)
					case imageutil.MODEL_GRAY:
						img = image.NewGray(rect)
					case imageutil.MODEL_GRAY16:
						img = image.NewGray16(rect)
					case imageutil.MODEL_CMYK:
						img = image.NewCMYK(rect)
					default:
						i.Lg.Append(fmt.Sprintf("invalid model: %d", args["model"].(int)), log.LEVEL_ERROR)
					}

					i.Self = &collection.ItemImage{
						Image:    img,
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Name:     args["name"].(string),
						Model:    model,
					}
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func name()
	/// @arg image_id - the id of the image to rename.
	/// @arg new_name - the new name to use for the image, not including the file extension.
	lib.CreateFunction("name",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("name_ext",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.TABLE, Name: "options", Table: &[]lua.Arg{
				{Type: lua.STRING, Name: "name", Optional: true},
				{Type: lua.STRING, Name: "prefix", Optional: true},
				{Type: lua.STRING, Name: "suffix", Optional: true},
				{Type: lua.STRING, Name: "ext", Optional: true},
			}},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("encoding",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("model",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			model := 0
			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model = int(i.Self.Model)
				},
			})

			r.State.PushInteger(model)
			return 1
		})

	/// @func size()
	/// @arg image_id - the id of the image to get the size from.
	/// @returns width
	/// @returns height
	/// @blocking
	lib.CreateFunction("size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

			r.State.PushInteger(width)
			r.State.PushInteger(height)
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
	lib.CreateFunction("crop",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x1"},
			{Type: lua.INT, Name: "y1"},
			{Type: lua.INT, Name: "x2"},
			{Type: lua.INT, Name: "y2"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self.Image = imageutil.SubImage(
						i.Self.Image,
						args["x1"].(int),
						args["y1"].(int),
						args["x2"].(int),
						args["y2"].(int),
						false,
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
	/// @returns new_id
	lib.CreateFunction("subimg",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "x1"},
			{Type: lua.INT, Name: "y1"},
			{Type: lua.INT, Name: "x2"},
			{Type: lua.INT, Name: "y2"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var simg image.Image
			var encoding imageutil.ImageEncoding
			simgReady := make(chan struct{}, 1)

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					simg = imageutil.SubImage(
						i.Self.Image,
						args["x1"].(int),
						args["y1"].(int),
						args["x2"].(int),
						args["y2"].(int),
						true,
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

			r.State.PushInteger(id)

			return 1
		})

	/// @func copy()
	/// @arg id
	/// @arg name
	/// @returns new_id
	lib.CreateFunction("copy",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var cimg image.Image
			var encoding imageutil.ImageEncoding
			cimgReady := make(chan struct{}, 1)

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					cimg = imageutil.CopyImage(i.Self.Image)
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
					}
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func pixel()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @returns {red, green, blue, alpha}
	/// @blocking
	lib.CreateFunction("pixel",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
			var red, green, blue, alpha uint32

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					col := i.Self.Image.At(args["x"].(int), args["y"].(int))
					red, green, blue, alpha = col.RGBA()
				},
			})

			r.State.NewTable()
			r.State.PushInteger(int(red))
			r.State.SetField(-2, "red")
			r.State.PushInteger(int(green))
			r.State.SetField(-2, "green")
			r.State.PushInteger(int(blue))
			r.State.SetField(-2, "blue")
			r.State.PushInteger(int(alpha))
			r.State.SetField(-2, "alpha")
			return 4
		})

	/// @func color_hex()
	/// @arg hex
	/// @returns {red, green, blue, alpha}
	lib.CreateFunction("color_hex",
		[]lua.Arg{
			{Type: lua.STRING, Name: "hex"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				alpha = int(c)
				fallthrough
			case 3:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[0]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s %s", hex, err), log.LEVEL_ERROR))
					r.State.Error()
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[1])+string(hex[1]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[2])+string(hex[2]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on blue): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				blue = int(c)

			case 8:
				c, err := strconv.ParseInt(string(hex[6])+string(hex[7]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on alpha): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				alpha = int(c)
				fallthrough
			case 6:
				c, err := strconv.ParseInt(string(hex[0])+string(hex[1]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on red): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				red = int(c)

				c, err = strconv.ParseInt(string(hex[2])+string(hex[3]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on green): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				green = int(c)

				c, err = strconv.ParseInt(string(hex[4])+string(hex[5]), 16, 64)
				if err != nil {
					r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string (failed on blue): %s", hex), log.LEVEL_ERROR))
					r.State.Error()
				}
				blue = int(c)
			default:
				r.State.PushString(lg.Append(fmt.Sprintf("invalid hex string: %s", hex), log.LEVEL_ERROR))
				r.State.Error()
			}

			r.State.NewTable()
			r.State.PushInteger(red)
			r.State.SetField(-2, "red")
			r.State.PushInteger(green)
			r.State.SetField(-2, "green")
			r.State.PushInteger(blue)
			r.State.SetField(-2, "blue")
			r.State.PushInteger(alpha)
			r.State.SetField(-2, "alpha")
			return 1
		})

	/// @func pixel_set()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg {red, green, blue, alpha}
	lib.CreateFunction("pixel_set",
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
		func(d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					color := args["color"].(map[string]any)
					imageutil.Set(
						i.Self.Image,
						args["x"].(int),
						args["y"].(int),
						color["red"].(int),
						color["green"].(int),
						color["blue"].(int),
						color["alpha"].(int),
					)
				},
			})
			return 0
		})

	/// @func convert_color()
	/// @arg model
	/// @arg color {red, green, blue, alpha}
	/// @returns new color {red, green, blue, alpha}
	lib.CreateFunction("convert_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
			{Type: lua.TABLE, Name: "color", Table: &[]lua.Arg{
				{Type: lua.INT, Name: "red"},
				{Type: lua.INT, Name: "green"},
				{Type: lua.INT, Name: "blue"},
				{Type: lua.INT, Name: "alpha"},
			}},
		},
		func(d lua.TaskData, args map[string]any) int {
			color := args["color"].(map[string]any)
			red, green, blue, alpha := imageutil.ConvertColor(
				lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib),
				color["red"].(int),
				color["green"].(int),
				color["blue"].(int),
				color["alpha"].(int),
			)

			r.State.NewTable()
			r.State.PushInteger(red)
			r.State.SetField(-2, "red")
			r.State.PushInteger(green)
			r.State.SetField(-2, "green")
			r.State.PushInteger(blue)
			r.State.SetField(-2, "blue")
			r.State.PushInteger(alpha)
			r.State.SetField(-2, "alpha")
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
	r.State.PushInteger(int(imageutil.MODEL_RGBA))
	r.State.SetField(-2, "MODEL_RGBA")
	r.State.PushInteger(int(imageutil.MODEL_RGBA64))
	r.State.SetField(-2, "MODEL_RGBA64")
	r.State.PushInteger(int(imageutil.MODEL_NRGBA))
	r.State.SetField(-2, "MODEL_NRGBA")
	r.State.PushInteger(int(imageutil.MODEL_NRGBA64))
	r.State.SetField(-2, "MODEL_NRGBA64")
	r.State.PushInteger(int(imageutil.MODEL_ALPHA))
	r.State.SetField(-2, "MODEL_ALPHA")
	r.State.PushInteger(int(imageutil.MODEL_ALPHA16))
	r.State.SetField(-2, "MODEL_ALPHA16")
	r.State.PushInteger(int(imageutil.MODEL_GRAY))
	r.State.SetField(-2, "MODEL_GRAY")
	r.State.PushInteger(int(imageutil.MODEL_GRAY16))
	r.State.SetField(-2, "MODEL_GRAY16")
	r.State.PushInteger(int(imageutil.MODEL_CMYK))
	r.State.SetField(-2, "MODEL_CMYK")

	/// @constants Encodings
	r.State.PushInteger(int(imageutil.ENCODING_PNG))
	r.State.SetField(-2, "ENCODING_PNG")
	r.State.PushInteger(int(imageutil.ENCODING_JPEG))
	r.State.SetField(-2, "ENCODING_JPEG")
	r.State.PushInteger(int(imageutil.ENCODING_GIF))
	r.State.SetField(-2, "ENCODING_GIF")
}
