package lib

import (
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

	r.State.SetGlobal(LIB_IMAGE)
}
