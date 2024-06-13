package lib

import (
	"fmt"
	"strings"

	img "github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
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
	r.State.NewTable()

	/// @func name()
	/// @arg image_id - the id of the image to rename.
	/// @arg new_name - the new name to use for the image, including the file extension.
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("image.name called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-2)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to image.name", log.LEVEL_ERROR))
			r.State.Error()
		}

		name, ok := r.State.ToString(-1)
		if !ok {
			r.State.PushString(lg.Append("invalid image name provided to image.name", log.LEVEL_ERROR))
			r.State.Error()
		}

		r.IC.Schedule(id, &img.ImageTask{
			Fn: func(i *img.Image) {
				lg.Append("image.name task ran", log.LEVEL_INFO)
				i.Name = name
				lg.Append("image.name task finished", log.LEVEL_INFO)
			},
		})
		return 0
	})
	r.State.SetField(-2, "name")

	/// @func name_ext()
	/// @arg image_id - the id of the image to rename.
	/// @arg options - a table containing each rename step. [name, prefix, suffix, ext]
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("image.name_ext called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-2)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to image.name", log.LEVEL_ERROR))
			r.State.Error()
		}

		state.Field(-1, "prefix")
		state.Field(-2, "suffix")
		state.Field(-3, "name")
		state.Field(-4, "ext")

		prefix, prefixOk := state.ToString(-4)
		suffix, suffixOk := state.ToString(-3)
		name, nameOk := state.ToString(-2)
		ext, extOk := state.ToString(-1)

		r.IC.Schedule(id, &img.ImageTask{
			Fn: func(i *img.Image) {
				lg.Append("image.name_ext task ran", log.LEVEL_INFO)

				fileSplit := strings.Split(i.Name, ".")
				fileName := strings.Join(fileSplit[:len(fileSplit)-1], ".")
				fileExt := fileSplit[len(fileSplit)-1]

				if nameOk {
					fileName = name
				}
				if prefixOk {
					fileName = prefix + fileName
				}
				if suffixOk {
					fileName += suffix
				}
				if extOk {
					fileExt = ext
				}

				i.Name = fileName + fileExt

				lg.Append(fmt.Sprintf("new image name: %s", i.Name), log.LEVEL_INFO)
				lg.Append("image.name_ext task finished", log.LEVEL_INFO)
			},
		})

		return 0
	})
	r.State.SetField(-2, "name_ext")

	/// @func collect()
	/// @arg image_id - the id of the image to collect.
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("image.collect called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-1)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to image.collect", log.LEVEL_ERROR))
			r.State.Error()
		}

		r.IC.CollectImage(id)

		return 0
	})
	r.State.SetField(-2, "collect")

	/// @func size()
	/// @arg image_id - the id of the image to get the size from.
	/// @returns width
	/// @returns height
	/// @blocking
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("image.size called", log.LEVEL_INFO)

		id, ok := r.State.ToInteger(-1)
		if !ok {
			r.State.PushString(lg.Append("invalid image id provided to image.size", log.LEVEL_ERROR))
			r.State.Error()
		}

		wait := make(chan bool, 1)
		width := 0
		height := 0

		r.IC.Schedule(id, &img.ImageTask{
			Fn: func(i *img.Image) {
				b := i.Img.Bounds()
				width = b.Dx()
				height = b.Dy()
				wait <- true
			},
		})

		<-wait
		r.State.PushInteger(width)
		r.State.PushInteger(height)
		return 2
	})
	r.State.SetField(-2, "size")

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

	r.State.SetGlobal(LIB_IMAGE)
}
