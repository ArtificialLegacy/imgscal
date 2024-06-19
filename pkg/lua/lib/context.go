package lib

import (
	"fmt"
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/fogleman/gg"
)

const LIB_CONTEXT = "context"

func RegisterContext(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_CONTEXT, r.State, lg)

	/// @func degrees()
	/// @arg radians - float
	/// @returns degrees - float
	lib.CreateFunction("degrees",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "rad"},
		},
		func(d lua.TaskData, args map[string]any) int {
			deg := gg.Degrees(args["rad"].(float64))
			r.State.PushNumber(deg)
			return 1
		})

	/// @func radians()
	/// @arg degrees - float
	/// @returns radians - float
	lib.CreateFunction("radians",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "deg"},
		},
		func(d lua.TaskData, args map[string]any) int {
			rad := gg.Radians(args["deg"].(float64))
			r.State.PushNumber(rad)
			return 1
		})

	/// @func new()
	/// @arg width - int
	/// @arg height - int
	/// returns id
	lib.CreateFunction("new",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(d lua.TaskData, args map[string]any) int {
			name := fmt.Sprintf("context_%d", r.CC.Next())

			chLog := log.NewLogger(name)
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

			id := r.CC.AddItem(name, &chLog)

			r.CC.Schedule(id, &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					c := gg.NewContext(args["width"].(int), args["height"].(int))
					i.Self = c
					i.Lg.Append("new context created", log.LEVEL_INFO)
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func to_image()
	/// @arg id
	/// @arg ext - defaults to png
	/// @returns id - new image id
	lib.CreateFunction("to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "ext", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			ext := "png"
			if args["ext"] != "" {
				ext = args["ext"].(string)
			}
			name := fmt.Sprintf("image_context_%d.%s", args["id"], ext)

			chLog := log.NewLogger(name)
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(name, &chLog)
			contextFinish := make(chan struct{}, 1)
			contextReady := make(chan struct{}, 1)

			var context *gg.Context

			r.CC.Schedule(id, &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					context = i.Self
					contextReady <- struct{}{}
					<-contextFinish
				},
			})

			r.IC.Schedule(id, &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					<-contextReady

					img := context.Image()
					i.Self = &img

					contextFinish <- struct{}{}
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @constants Fill Rules
	/// @const FILLRULE_WINDING
	/// @const FILLRULE_EVENODD
	lib.State.PushInteger(int(gg.FillRuleWinding))
	lib.State.SetField(-2, "FILLRULE_WINDING")
	lib.State.PushInteger(int(gg.FillRuleEvenOdd))
	lib.State.SetField(-2, "FILLRULE_EVENODD")

	/// @constants Line Caps
	/// @const LINECAP_ROUND
	/// @const LINECAP_BUTT
	/// @const LINCAP_SQUARE
	lib.State.PushInteger(int(gg.LineCapRound))
	lib.State.SetField(-2, "LINECAP_ROUND")
	lib.State.PushInteger(int(gg.LineCapButt))
	lib.State.SetField(-2, "LINECAP_BUTT")
	lib.State.PushInteger(int(gg.LineCapSquare))
	lib.State.SetField(-2, "LINECAP_SQUARE")

	/// @constants Line Joins
	/// @const LINEJOIN_ROUND
	/// @const LINEJOIN_BEVEL
	lib.State.PushInteger(int(gg.LineJoinRound))
	lib.State.SetField(-2, "LINEJOIN_ROUND")
	lib.State.PushInteger(int(gg.LineJoinBevel))
	lib.State.SetField(-2, "LINEJOIN_BEVEL")

	/// @constants Repeat Ops
	/// @const REPEAT_BOTH
	/// @const REPEAT_X
	/// @const REPEAT_Y
	/// @const REPEAT_NONE
	lib.State.PushInteger(int(gg.RepeatBoth))
	lib.State.SetField(-2, "REPEAT_BOTH")
	lib.State.PushInteger(int(gg.RepeatX))
	lib.State.SetField(-2, "REPEAT_X")
	lib.State.PushInteger(int(gg.RepeatY))
	lib.State.SetField(-2, "REPEAT_Y")
	lib.State.PushInteger(int(gg.RepeatNone))
	lib.State.SetField(-2, "REPEAT_NONE")

	/// @constants Alignment
	/// @const ALIGN_LEFT
	/// @const ALIGN_CENTER
	/// @const ALIGN_RIGHT
	lib.State.PushInteger(int(gg.AlignLeft))
	lib.State.SetField(-2, "ALIGN_LEFT")
	lib.State.PushInteger(int(gg.AlignCenter))
	lib.State.SetField(-2, "ALIGN_CENTER")
	lib.State.PushInteger(int(gg.AlignRight))
	lib.State.SetField(-2, "ALIGN_RIGHT")
}
