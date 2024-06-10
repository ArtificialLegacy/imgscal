package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
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

	lg := log.NewLogger()
	lg.Append("log started for workflow_run", log.LEVEL_INFO)
	state := lua.WorkflowRunState(&lg)
	runner := lua.NewRunner(state, &data)

	for _, plugin := range req {
		switch plugin {
		case lib.LIB_IMGSCAL:
			lg.Append("registering standard plugin IMGSCAL", log.LEVEL_INFO)
			runner.Register(lib.RegisterImgscal, &lg)
		case lib.LIB_IMGSCALSHEET:
			lg.Append("registering standard plugin IMGSCAL_SHEET", log.LEVEL_INFO)
			runner.Register(lib.RegisterImgscalSheet, &lg)
		}
	}

	err := runner.Run(script)
	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)
		lg.Dump("./log")
		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)
		return nil
	}

	lg.Append("Collecting images", log.LEVEL_INFO)
	data.IC.Collect()

	lg.Dump("./log")

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}
