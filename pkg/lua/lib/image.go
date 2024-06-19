package lib

import (
	"fmt"
	"image"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

const LIB_IMAGE = "image"

const (
	model_RGBA int = iota
	model_RGBA64
	model_NRGBA
	model_NRGBA64
	model_ALPHA
	model_ALPHA16
	model_GRAY
	model_GRAY16
	model_CMYK
	model_NYCBCRA
	model_YCBCR
)

func RegisterImage(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_IMAGE, r.State, lg)

	/// @func name()
	/// @arg image_id - the id of the image to rename.
	/// @arg new_name - the new name to use for the image, including the file extension.
	lib.CreateFunction("name",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					i.Name = args["name"].(string)
				},
			})
			return 0
		})

	/// @func name_ext()
	/// @arg image_id - the id of the image to rename.
	/// @arg options - a table containing each rename step. [name, prefix, suffix, ext]
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
			r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					fileSplit := strings.Split(i.Name, ".")
					fileName := strings.Join(fileSplit[:len(fileSplit)-1], ".")
					fileExt := fileSplit[len(fileSplit)-1]

					opt := args["options"].(map[string]any)

					if opt["name"] != "" {
						fileName = opt["name"].(string)
					}
					if opt["prefix"] != "" {
						fileName = opt["prefix"].(string) + fileName
					}
					if opt["suffix"] != "" {
						fileName += opt["suffix"].(string)
					}
					if opt["ext"] != "" {
						fileExt = opt["ext"].(string)
					}

					i.Name = fileName + fileExt
					i.Lg.Append(fmt.Sprintf("new image name: %s", i.Name), log.LEVEL_INFO)
				},
			})
			return 0
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

			<-r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					b := (*i.Self).Bounds()
					width = b.Dx()
					height = b.Dy()
				},
			})

			r.State.PushInteger(width)
			r.State.PushInteger(height)
			return 2
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
	/// @const NYCBCRA
	/// @const YCBCR
	r.State.PushInteger(model_RGBA)
	r.State.SetField(-2, "RGBA")
	r.State.PushInteger(model_RGBA64)
	r.State.SetField(-2, "RGBA64")
	r.State.PushInteger(model_NRGBA)
	r.State.SetField(-2, "NRGBA")
	r.State.PushInteger(model_NRGBA64)
	r.State.SetField(-2, "NRGBA64")
	r.State.PushInteger(model_ALPHA)
	r.State.SetField(-2, "ALPHA")
	r.State.PushInteger(model_ALPHA16)
	r.State.SetField(-2, "ALPHA16")
	r.State.PushInteger(model_GRAY)
	r.State.SetField(-2, "GRAY")
	r.State.PushInteger(model_GRAY16)
	r.State.SetField(-2, "GRAY16")
	r.State.PushInteger(model_CMYK)
	r.State.SetField(-2, "CMYK")
	r.State.PushInteger(model_NYCBCRA)
	r.State.SetField(-2, "NYCBCRA")
	r.State.PushInteger(model_YCBCR)
	r.State.SetField(-2, "YCBCR")
}
