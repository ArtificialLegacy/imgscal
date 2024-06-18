package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

const LIB_STD = "std"

func RegisterStd(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_STD, r.State, lg)

	/// @func log()
	/// @arg msg - the message to display in the log
	lib.CreateFunction("log",
		[]lua.Arg{
			{Type: lua.ANY, Name: "msg"},
		},
		func(d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua log: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func warn()
	/// @arg msg - the message to display as a warning in the log
	lib.CreateFunction("warn",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(d lua.TaskData, args map[string]any) int {
			lg.Append(fmt.Sprintf("lua warn: %s", args["msg"]), log.LEVEL_WARN)
			return 0
		})

	/// @func panic()
	/// @arg msg - the message to display in the error
	lib.CreateFunction("panic",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.State.PushString(lg.Append(fmt.Sprintf("lua panic: %s", args["msg"]), log.LEVEL_ERROR))
			r.State.Error()

			return 0
		})
}
