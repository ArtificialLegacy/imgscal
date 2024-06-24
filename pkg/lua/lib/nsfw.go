package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
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
		func(d lua.TaskData, args map[string]any) int {
			result := false

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					res, err := nude.IsImageNude(i.Self.Image)
					if err != nil {
						r.State.PushString(i.Lg.Append("nsfw skin check failed", log.LEVEL_ERROR))
						r.State.Error()
					}

					result = res
				},
			})

			r.State.PushBoolean(result)
			return 1
		})
}
