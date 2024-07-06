package lib

import (
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TIME = "time"

func RegisterTime(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TIME, r, r.State, lg)

	/// @func now_ms()
	/// @returns current time in ms
	lib.CreateFunction(tab, "now_ms", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano() / int64(time.Millisecond)

			state.Push(golua.LNumber(t))
			return 1
		})

	// @func now_date()
	/// @returns current date in MM-DD-YEAR format
	lib.CreateFunction(tab, "now_ms", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format("01-02-2006")

			state.Push(golua.LString(t))
			return 1
		})

	// @func now_timestamp()
	/// @returns current time in the default format
	lib.CreateFunction(tab, "now_timestamp", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().String()

			state.Push(golua.LString(t))
			return 1
		})

	// @func now_format()
	/// @arg format - go layout string
	/// @returns current time in given format
	lib.CreateFunction(tab, "now_ms",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format(args["format"].(string))

			state.Push(golua.LString(t))
			return 1
		})
}
