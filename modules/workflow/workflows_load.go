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
			workflow.config.name, _ = state.ToString(-1)

			state.RawGetValue(1, "version")
			workflow.config.version, _ = state.ToString(-1)

			state.RawGetValue(1, "requires")
			state.Length(-1)
			len, _ := state.ToInteger(-1)

			workflow.config.requires = make([]string, 0, len)

			for i := 1; i <= len; i++ {
				state.RawGetInt(-1-i, i)
				val, _ := state.ToString(-1)
				workflow.config.requires = append(workflow.config.requires, val)
			}
		}
		return 1
	})
}

func WorkflowsLoad() map[string]*Workflow {
	workflows := make(map[string]*Workflow)

	scripts, err := os.ReadDir("workflows")
	if err != nil {
		panic(err)
	}

	for _, script := range scripts {
		if script.IsDir() {
			continue
		}

		state := lua.NewState()
		workflow := &Workflow{
			file:   script.Name(),
			config: WorkflowConfig{},
		}

		workflows[workflow.file] = workflow

		registerConfig(state, workflow)

		pwd, _ := os.Getwd()

		err := lua.DoFile(state, fmt.Sprintf("%s\\workflows\\%s", pwd, script.Name()))
		workflow.succeed = err == nil
	}

	return workflows
}
