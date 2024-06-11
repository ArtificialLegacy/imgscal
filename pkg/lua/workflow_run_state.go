package lua

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/Shopify/go-lua"
)

func WorkflowRunState(lg *log.Logger) *lua.State {
	state := lua.NewState()

	state.Register("config", func(state *lua.State) int {
		lg.Append("config function called", log.LEVEL_INFO)
		return 0
	})

	state.Register("main", func(state *lua.State) int {
		lg.Append("main function called", log.LEVEL_INFO)
		state.Call(0, 0)
		return 0
	})

	return state
}
