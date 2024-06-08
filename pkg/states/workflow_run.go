package states

import (
	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/lua/lib"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowRun(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	req := []string{}
	for sm.Peek() {
		req = append(req, sm.PopString())
	}

	state := lua.WorkflowRunState()
	for _, plugin := range req {
		switch plugin {
		case lib.LIB_IMGSCAL:
			lib.RegisterImgscal(state)
		case lib.LIB_IMGSCALSHEET:
			lib.RegisterImgscalSheet(state)
		}
	}

	runner := lua.NewRunner(state)
	runner.Run(script)

	return nil
}
