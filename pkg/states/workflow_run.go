package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/lua/lib"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	golua "github.com/Shopify/go-lua"
)

func WorkflowRun(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	req := []string{}
	for sm.Peek() {
		req = append(req, sm.PopString())
	}

	lg := log.NewLogger("./log")

	lg.Append("log started for workflow_run", log.LEVEL_INFO)
	state := lua.WorkflowRunState(&lg)
	runner := lua.NewRunner(state, &lg)

	golua.Require(state, "basic", golua.BaseOpen, true)

	for _, plugin := range req {
		builtin, ok := lib.Builtins[plugin]
		if !ok {
			lg.Append(fmt.Sprintf("plugin %s does not exist", plugin), log.LEVEL_WARN)
		} else {
			builtin(&runner, &lg)
			state.Pop(1)
			lg.Append(fmt.Sprintf("registered plugin %s", plugin), log.LEVEL_INFO)
		}
	}

	err := runner.Run(script)
	runner.IC.Collect()
	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)
		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)
		return nil
	}

	lg.Append("workflow finished", log.LEVEL_INFO)

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}
