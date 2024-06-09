package states

import (
	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/lua/lib"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

func WorkflowRun(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	req := []string{}
	for sm.Peek() {
		req = append(req, sm.PopString())
	}

	data := workflow.WorkflowData{
		IC: *image.NewImageCollection(),
	}

	state := lua.WorkflowRunState()
	runner := lua.NewRunner(state, &data)

	for _, plugin := range req {
		switch plugin {
		case lib.LIB_IMGSCAL:
			runner.Register(lib.RegisterImgscal)
		case lib.LIB_IMGSCALSHEET:
			runner.Register(lib.RegisterImgscalSheet)
		}
	}

	err := runner.Run(script)
	if err != nil {
		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)
		return nil
	}

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}
