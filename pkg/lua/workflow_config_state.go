package lua

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/Shopify/go-lua"
)

func WorkflowConfigState(wf *workflow.Workflow, lg *log.Logger) *lua.State {
	state := lua.NewState()

	state.Register("config", func(state *lua.State) int {
		lg.Append("config function called", log.LEVEL_INFO)

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
		} else {
			lg.Append("value passed to config is not a table", log.LEVEL_ERROR)
			state.PushString("value passed to config is not a table")
			state.Error()
		}

		return 0
	})

	state.Register("main", func(state *lua.State) int {
		lg.Append("main function called", log.LEVEL_INFO)
		return 0
	})

	return state
}
