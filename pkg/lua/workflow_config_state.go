package lua

import (
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	lua "github.com/yuin/gopher-lua"
)

func WorkflowConfigState(wf *workflow.Workflow, lg *log.Logger) *lua.LState {
	state := lua.NewState()

	state.Register("config", func(state *lua.LState) int {
		lg.Append("config function called", log.LEVEL_INFO)

		t := state.Get(-1)

		if t.Type() != lua.LTTable {
			state.Error(lua.LString(lg.Append("value passed to config is not a table", log.LEVEL_ERROR)), 0)
		} else {
			name := state.GetField(t, "name")
			wf.Name = strings.Clone(string(name.(lua.LString)))

			version := state.GetField(t, "version")
			wf.Version = strings.Clone(string(version.(lua.LString)))

			desc := state.GetField(t, "desc")
			wf.Desc = strings.Clone(string(desc.(lua.LString)))

			requires := state.GetField(t, "requires").(*lua.LTable)
			wf.Requires = []string{}

			requires.ForEach(func(l1, l2 lua.LValue) {
				wf.Requires = append(wf.Requires, strings.Clone(string(l2.(lua.LString))))
			})
		}

		return 0
	})

	state.Register("main", func(state *lua.LState) int {
		lg.Append("main function called", log.LEVEL_INFO)
		return 0
	})

	return state
}
