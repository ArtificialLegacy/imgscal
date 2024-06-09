package lua

import (
	"github.com/Shopify/go-lua"
)

func WorkflowRunState() *lua.State {
	state := lua.NewState()

	state.Register("config", func(state *lua.State) int {
		return 0
	})

	state.Register("main", func(state *lua.State) int {
		state.Call(0, 0)
		return 0
	})

	return state
}
