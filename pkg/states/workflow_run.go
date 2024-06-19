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

	lg := log.NewLogger("main")
	defer lg.Close()

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
	errIC := runner.IC.CollectAll()
	errFC := runner.FC.CollectAll()
	errCC := runner.CC.CollectAll()
	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)

		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)

		return nil
	}

	eri := collErr(errIC, "IC", script, &lg, sm)
	erf := collErr(errFC, "FC", script, &lg, sm)
	erc := collErr(errCC, "CC", script, &lg, sm)
	if eri || erf || erc {
		return nil
	}

	lg.Append("workflow finished", log.LEVEL_INFO)

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}

func collErr(err error, name, script string, lg *log.Logger, sm *statemachine.StateMachine) bool {
	if err != nil {
		lg.Append(fmt.Sprintf("error occured within %s collection: %s", name, err), log.LEVEL_ERROR)

		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)

		return true
	}

	return false
}
