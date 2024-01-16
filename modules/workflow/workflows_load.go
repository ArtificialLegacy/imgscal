package workflow

import (
	"fmt"
	"os"

	"github.com/Shopify/go-lua"
)

func registerConfig(state *lua.State, workflow *Workflow) {
	state.Register("config", func(state *lua.State) int {
		if state.IsTable(1) {
			state.RawGetValue(1, "name")
			workflow.Config.Name, _ = state.ToString(-1)

			state.RawGetValue(1, "version")
			workflow.Config.Version, _ = state.ToString(-1)

			state.RawGetValue(1, "requires")
			state.Length(-1)
			len, _ := state.ToInteger(-1)

			workflow.Config.Requires = make([]string, 0, len)

			for i := 1; i <= len; i++ {
				state.RawGetInt(-1-i, i)
				val, _ := state.ToString(-1)
				workflow.Config.Requires = append(workflow.Config.Requires, val)
			}
		}
		return 1
	})
}

func emptyMain(state *lua.State) {
	state.Register("main", func(state *lua.State) int {
		return 1
	})
}

func WorkflowsLoad() map[string]*Workflow {
	workflows := make(map[string]*Workflow)

	scripts, err := os.ReadDir("workflows")
	if err != nil {
		panic(err)
	}

	var count int8 = 0

	for _, script := range scripts {
		if script.IsDir() {
			continue
		}

		count++
		if count >= 127 {
			break
		}

		workflow := &Workflow{
			File:   script.Name(),
			Config: WorkflowConfig{},
		}
		workflows[workflow.File] = workflow

		state := lua.NewState()
		registerConfig(state, workflow)
		emptyMain(state)

		pwd, _ := os.Getwd()

		err := lua.DoFile(state, fmt.Sprintf("%s\\workflows\\%s", pwd, script.Name()))
		workflow.Succeed = err == nil
	}

	return workflows
}
