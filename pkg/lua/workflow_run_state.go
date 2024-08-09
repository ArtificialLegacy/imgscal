package lua

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	lua "github.com/yuin/gopher-lua"
)

func WorkflowRunState(lg *log.Logger) *lua.LState {
	state := lua.NewState()

	state.Register("config", func(state *lua.LState) int {
		lg.Append("config function called", log.LEVEL_SYSTEM)
		return 0
	})

	state.Register("main", func(state *lua.LState) int {
		lg.Append("main function called", log.LEVEL_SYSTEM)
		state.Call(0, 0)
		return 0
	})

	return state
}
