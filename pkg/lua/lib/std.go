package lib

import (
	"fmt"
	"time"

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

	/// @func sleep(ms)
	/// @arg ms {int} - The number of milliseconds to sleep.
	lib.CreateFunction(tab, "sleep",
		[]lua.Arg{
			{Type: lua.INT, Name: "ms"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ms := args["ms"].(int)
			time.Sleep(time.Duration(time.UnixMilli(int64(ms)).UnixNano()))
			return 0
		})

	/// @func fmt(str, values...) -> string
	/// @arg str {string}
	/// @arg values {any...}
	/// @returns {string}
	lib.CreateFunction(tab, "fmt",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			lua.ArgVariadic("values", lua.ArrayType{Type: lua.ANY}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			values := args["values"].([]any)
			valueList := make([]any, len(values))
			for i, v := range values {
				valueList[i] = lua.GetValue(v.(golua.LValue))
			}

			format := fmt.Sprintf(args["str"].(string), valueList...)
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

	/// @func call_thread(func)
	/// @arg func {function} - The function to call in a new thread.
	/// @desc
	/// This function will call the provided function in a new thread, unlike tasks there is no control over the thread.
	/// There is also no guarantee that the thread will finish before the script ends.
	lib.CreateFunction(tab, "call_thread",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			scheduledState, _ := state.NewThread()
			go func() {
				scheduledState.Push(args["func"].(golua.LValue))
				scheduledState.Call(0, 0)
				scheduledState.Close()
			}()
			return 0
		})
}
