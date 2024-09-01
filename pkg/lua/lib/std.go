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

	/// @func log(msg)
	/// @arg msg {string} - The message to display in the log.
	lib.CreateFunction(tab, "log",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua log: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func log_value(value)
	/// @arg value {any} - The value to display in the log.
	lib.CreateFunction(tab, "log_value",
		[]lua.Arg{
			{Type: lua.ANY, Name: "value"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua log: %+v", args["value"]), log.LEVEL_INFO)
			return 0
		})

	/// @func warn(msg)
	/// @arg msg {string} - The message to display as a warning in the log.
	lib.CreateFunction(tab, "warn",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua warn: %s", args["msg"]), log.LEVEL_WARN)
			return 0
		})

	/// @func panic(msg)
	/// @arg msg {string} - The message to display in the error.
	/// @desc
	/// This results in a lua panic.
	lib.CreateFunction(tab, "panic",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("lua panic: %s", args["msg"]), log.LEVEL_ERROR)), 0)
			return 0
		})

	/// @func fmt(str, values) -> string
	/// @arg str {string}
	/// @arg values {[]any} - The value in each index should be compatible with the Go fmt string provided.
	/// @returns {string}
	lib.CreateFunction(tab, "fmt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "values"},
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

	/// @func config() -> table<any>
	/// @returns {table<any>}
	lib.CreateFunction(tab, "config",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := lua.CreateValue(r.ConfigData, state)

			state.Push(data)
			return 1
		})

	/// @func secrets() -> table<any>
	/// @returns {table<any>}
	lib.CreateFunction(tab, "secrets",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := lua.CreateValue(r.SecretData, state)

			state.Push(data)
			return 1
		})
}
