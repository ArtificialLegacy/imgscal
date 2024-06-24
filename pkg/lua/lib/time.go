package lib

import (
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

const LIB_TIME = "time"

func RegisterTime(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_TIME, r.State, lg)

	/// @func now_ms()
	/// @returns current time in ms
	lib.CreateFunction("now_ms", []lua.Arg{},
		func(d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano() / int64(time.Millisecond)

			r.State.PushInteger(int(t))
			return 1
		})

	// @func now_date()
	/// @returns current date in MM-DD-YEAR format
	lib.CreateFunction("now_ms", []lua.Arg{},
		func(d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format("01-02-2006")

			r.State.PushString(t)
			return 1
		})

	// @func now_timestamp()
	/// @returns current time in the default format
	lib.CreateFunction("now_timestamp", []lua.Arg{},
		func(d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().String()

			r.State.PushString(t)
			return 1
		})

	// @func now_format()
	/// @arg format - go layout string
	/// @returns current time in given format
	lib.CreateFunction("now_ms",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format(args["format"].(string))

			r.State.PushString(t)
			return 1
		})
}
