package lib

import (
	img "github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
	"github.com/koyachi/go-nude"
)

const LIB_NSFW = "nsfw"

func RegisterNSFW(r *lua.Runner, lg *log.Logger) {
	r.State.NewTable()

	/// @func skin()
	/// @arg image_id - the image to check for nudity using skin content.
	/// @returns boolean - if the skin content detector is over a threshold.
	/// @blocking
	/// @desc
	/// Not very accurate, but does not require an AI model.
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("nsfw.skin called", log.LEVEL_INFO)

		id, ok := state.ToInteger(-1)
		if !ok {
			state.PushString(lg.Append("invalid image id provided to nsfw.skin", log.LEVEL_ERROR))
			state.Error()
		}

		wait := make(chan bool, 1)
		result := false

		r.IC.Schedule(id, &img.ImageTask{
			Fn: func(i *img.Image) {
				r, err := nude.IsImageNude(i.Img)
				if err != nil {
					state.PushString(lg.Append("nsfw skin check failed", log.LEVEL_ERROR))
					state.Error()
				}
				result = r
				wait <- true
			},
		})

		<-wait
		r.State.PushBoolean(result)
		return 1
	})
	r.State.SetField(-2, "skin")

	r.State.SetGlobal(LIB_NSFW)
}
