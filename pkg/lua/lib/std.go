package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_STD = "std"

/// @lib Standard
/// @import std
/// @desc
/// A library of miscellaneous functions.

func RegisterStd(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_STD, r, r.State, lg)

	/// @func log()
	/// @arg msg - the message to display in the log
	lib.CreateFunction(tab, "log",
		[]lua.Arg{
			{Type: lua.ANY, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua log: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func warn()
	/// @arg msg - the message to display as a warning in the log
	lib.CreateFunction(tab, "warn",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua warn: %s", args["msg"]), log.LEVEL_WARN)
			return 0
		})

	/// @func panic()
	/// @arg msg - the message to display in the error
	lib.CreateFunction(tab, "panic",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("lua panic: %s", args["msg"]), log.LEVEL_ERROR)), 0)
			return 0
		})

	/// @func fmt()
	/// @arg str
	/// @arg []values
	/// @returns string
	lib.CreateFunction(tab, "fmt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			{Type: lua.ANY, Name: "values"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := []any{}
			values := args["values"].(*golua.LTable)
			for i := range values.Len() {
				v = append(v, state.GetTable(values, golua.LNumber(i+1)))
			}

			format := fmt.Sprintf(args["str"].(string), v...)
			state.Push(golua.LString(format))
			return 1
		})
}
