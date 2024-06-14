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
	lib := lua.NewLib(LIB_NSFW, r.State, lg)

	/// @func skin()
	/// @arg image_id - the image to check for nudity using skin content.
	/// @returns boolean - if the skin content detector is over a threshold.
	/// @blocking
	/// @desc
	/// Not very accurate, but does not require an AI model.
	lib.CreateFunction("skin",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.State, args map[string]any) int {
			wait := make(chan bool, 1)
			result := false

			r.IC.Schedule(args["id"].(int), &img.ImageTask{
				Lib:  LIB_NSFW,
				Name: "skin",
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
}
