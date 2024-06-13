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

	r.State.SetGlobal(LIB_IMAGE)
}
