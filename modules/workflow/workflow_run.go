package workflow

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/file"
	"github.com/Shopify/go-lua"
)

func registerMain(state *lua.State, file string) {
	state.Register("main", func(state *lua.State) int {
		if state.IsFunction(-1) {
			state.PushString(file)
			state.Call(1, 0)
		}
		return 1
	})
}

func registerJob(state *lua.State) {
	state.Register("job", func(state *lua.State) int {
		println("job run")
		return 1
	})
}

func emptyConfig(state *lua.State) {
	state.Register("config", func(state *lua.State) int {
		return 1
	})
}

func (workflow *Workflow) Run(filepath string, filename string, requires []string) error {
	pwd, _ := os.Getwd()

	_, err := file.Copy(filepath, fmt.Sprintf("%s\\temp\\%s", pwd, filename))
	if err != nil {
		return err
	}
	defer os.Remove(fmt.Sprintf("%s\\temp\\%s", pwd, filename))

	state := lua.NewState()
	registerMain(state, filename)
	registerJob(state)
	emptyConfig(state)

	lua.DoFile(state, fmt.Sprintf("%s\\workflows\\%s", pwd, workflow.File))

	return nil
}
