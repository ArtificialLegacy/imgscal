package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
)

const LIB_STD = "std"

func RegisterStd(r *lua.Runner, lg *log.Logger) {
	r.State.NewTable()

	/// @func panic()
	/// @arg msg - the message to display in the error
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("std.panic called", log.LEVEL_INFO)

		msg, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid msg provided to panic", log.LEVEL_ERROR))
			state.Error()
		}

		state.PushString(lg.Append(fmt.Sprintf("lua panic: %s", msg), log.LEVEL_ERROR))
		state.Error()

		return 0
	})
	r.State.SetField(-2, "panic")

	/// @func log()
	/// @arg msg - the message to display in the log
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("std.log called", log.LEVEL_INFO)

		msg := r.State.ToValue(-1)
		lg.Append(fmt.Sprintf("lua log: %s", msg), log.LEVEL_INFO)

		return 0
	})
	r.State.SetField(-2, "log")

	r.State.SetGlobal(LIB_STD)
}
