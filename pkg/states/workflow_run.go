package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/lua/lib"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	golua "github.com/yuin/gopher-lua"
)

func WorkflowRun(sm *statemachine.StateMachine) error {
	script := sm.PopString()
	verbose := sm.PopBool()
	req := []string{}
	for sm.Peek() {
		req = append(req, sm.PopString())
	}

	if !sm.CliMode {
		cli.Clear()
	}

	var lg log.Logger

	if sm.Config.DisableLogs {
		lg = log.NewLoggerEmpty()
	} else {
		lg = log.NewLoggerBase("execute", sm.Config.LogDirectory, false)
	}

	if verbose {
		lg.EnableVerbose()
	}

	lg.Append("log started for workflow_run", log.LEVEL_SYSTEM)
	state := lua.WorkflowRunState(&lg)
	runner := lua.NewRunner(req, state, &lg, sm.CliMode)
	runner.Output = sm.Config.OutputDirectory

	defer func() {
		if r := recover(); r != nil {
			lg.Append(fmt.Sprintf("panic recovered: %+v", r), log.LEVEL_ERROR)
		}

		lg.Close()
	}()

	golua.OpenBase(state)
	golua.OpenMath(state)
	golua.OpenString(state)
	golua.OpenTable(state)
	lua.LoadPlugins("main", &runner, &lg, lib.Builtins, req, state)

	err := runner.Run(script)

	runner.Wg.Wait()

	lg.Append("All collections empty, exiting", log.LEVEL_SYSTEM)

	runner.CR_WIN.CleanAll()
	runner.CR_REF.CleanAll()

	runner.TC.CollectAll()
	runner.IC.CollectAll()
	runner.FC.CollectAll()
	runner.CC.CollectAll()
	runner.QR.CollectAll()

	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)

		if sm.CliMode {
			fmt.Printf("error occured while running script: %s\n", err)

			sm.SetState(STATE_EXIT)
			return nil
		}

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
		if sm.CliMode {
			sm.SetState(STATE_EXIT)
			return fmt.Errorf("error running script")
		}

		sm.PushString("error occurred within collection")
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_RUN)

		return fmt.Errorf("error running script")
	}

	lg.Append("workflow finished", log.LEVEL_INFO)

	if sm.CliMode {
		sm.SetState(STATE_EXIT)
		return nil
	}

	sm.PushString(script)
	sm.SetState(STATE_WORKFLOW_FINISH)

	return nil
}

func collErr(errs []error, name, script string, lg *log.Logger, sm *statemachine.StateMachine) bool {
	errExists := false

	for _, err := range errs {
		if err != nil {
			lg.Append(fmt.Sprintf("error occured within %s collection: %s", name, err), log.LEVEL_ERROR)

			if sm.CliMode {
				fmt.Printf("error occured within %s collection: %s\n", name, err)

				sm.SetState(STATE_EXIT)
			} else {
				sm.PushString(err.Error())
				sm.PushString(script)
				sm.SetState(STATE_WORKFLOW_FAIL_RUN)
			}

			errExists = true
		}
	}

	return errExists
}
