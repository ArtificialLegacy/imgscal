package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
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

	lg := log.NewLogger("execute")

	lg.Append("log started for workflow_run", log.LEVEL_INFO)
	state := lua.WorkflowRunState(&lg)
	runner := lua.NewRunner(state, &lg)

	defer func() {
		if r := recover(); r != nil {
			tcCount, tcBusy := runner.TC.TaskCount()
			lg.Append(fmt.Sprintf("collection [%T] left: %d, (busy: %t)", runner.TC, tcCount, tcBusy), log.LEVEL_WARN)
			icCount, icBusy := runner.IC.TaskCount()
			lg.Append(fmt.Sprintf("collection [%T] left: %d, (busy: %t)", runner.IC, icCount, icBusy), log.LEVEL_WARN)
			fcCount, fcBusy := runner.FC.TaskCount()
			lg.Append(fmt.Sprintf("collection [%T] left: %d, (busy: %t)", runner.FC, fcCount, fcBusy), log.LEVEL_WARN)
			ccCount, ccBusy := runner.CC.TaskCount()
			lg.Append(fmt.Sprintf("collection [%T] left: %d, (busy: %t)", runner.CC, ccCount, ccBusy), log.LEVEL_WARN)
			qrCount, qrBusy := runner.QR.TaskCount()
			lg.Append(fmt.Sprintf("collection [%T] left: %d, (busy: %t)", runner.QR, qrCount, qrBusy), log.LEVEL_WARN)
		} else {
			lg.Close()
		}
	}()

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

	for checkState(runner.IC) || checkState(runner.FC) || checkState(runner.CC) || checkState(runner.QR) {
	}

	runner.TC.CollectAll()
	runner.IC.CollectAll()
	runner.FC.CollectAll()
	runner.CC.CollectAll()
	runner.QR.CollectAll()

	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)

		sm.PushString(err.Error())
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)

		return nil
	}

	ert := collErr(runner.TC.Errs, "TC", script, &lg, sm)
	eri := collErr(runner.IC.Errs, "IC", script, &lg, sm)
	erf := collErr(runner.FC.Errs, "FC", script, &lg, sm)
	erc := collErr(runner.CC.Errs, "CC", script, &lg, sm)
	erq := collErr(runner.QR.Errs, "QR", script, &lg, sm)
	if ert || eri || erf || erc || erq {
		return nil
	}

	lg.Append("workflow finished", log.LEVEL_INFO)

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}

func collErr(errs []error, name, script string, lg *log.Logger, sm *statemachine.StateMachine) bool {
	errExists := false

	for _, err := range errs {
		if err != nil {
			lg.Append(fmt.Sprintf("error occured within %s collection: %s", name, err), log.LEVEL_ERROR)

			sm.PushString(err.Error())
			sm.PushString(script)
			sm.SetState(STATE_WORKFLOW_FAIL_RUN)

			errExists = true
		}
	}

	return errExists
}

func checkState[T collection.ItemSelf](c *collection.Collection[T]) bool {
	co, b := c.TaskCount()

	return co > 0 || b
}
