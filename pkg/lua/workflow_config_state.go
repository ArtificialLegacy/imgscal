package lua

import (
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/Shopify/go-lua"
)

func WorkflowConfigState(wf *workflow.Workflow) *lua.State {
	state := lua.NewState()

	state.Register("config", func(state *lua.State) int {
		if state.IsTable(1) {
			state.RawGetValue(1, "name")
			wf.Name, _ = state.ToString(-1)

			state.RawGetValue(1, "version")
			wf.Version, _ = state.ToString(-1)

			state.RawGetValue(1, "author")
			wf.Author, _ = state.ToString(-1)

			state.RawGetValue(1, "desc")
			wf.Desc, _ = state.ToString(-1)

			state.RawGetValue(1, "requires")
			state.Length(-1)
			len, _ := state.ToInteger(-1)

			wf.Requires = make([]string, 0, len)

			for i := 1; i <= len; i++ {
				state.RawGetInt(-1-i, i)
				val, _ := state.ToString(-1)
				wf.Requires = append(wf.Requires, val)
			}
		}

		return 0
	})

	state.Register("main", func(state *lua.State) int {
		return 0
	})

	return state
}
