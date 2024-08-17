package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/koyachi/go-nude"
	golua "github.com/yuin/gopher-lua"
)

const LIB_NSFW = "nsfw"

/// @lib NSFW
/// @import nsfw
/// @desc
/// Provides basic functionality for filtering image content, carry-over from when SD was supported.

func RegisterNSFW(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_NSFW, r, r.State, lg)

	/// @func skin()
	/// @arg image_id - the image to check for nudity using skin content.
	/// @returns boolean - if the skin content detector is over a threshold.
	/// @blocking
	/// @desc
	/// Not very accurate, but does not require an AI model.
	lib.CreateFunction(tab, "skin",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			result := false

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					res, err := nude.IsImageNude(i.Self.Image)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("nsfw skin check failed", log.LEVEL_ERROR)), 0)
					}

					result = res
				},
			})

			state.Push(golua.LBool(result))
			return 1
		})
}
