package states

import (
	"fmt"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/lua/lib"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	golua "github.com/yuin/gopher-lua"
)

type WorkflowRunData struct {
	Script string
	Name   string
}

func WorkflowRunEnter(sm *statemachine.StateMachine, data WorkflowRunData) {
	sm.SetState(STATE_WORKFLOW_RUN)
	sm.Data = data
}

func WorkflowRun(sm *statemachine.StateMachine) error {
	data := sm.Data.(WorkflowRunData)
	pth := data.Script
	name := data.Name
	sm.Data = nil

	if !sm.CliMode {
		cli.Clear()
	}

	var lg log.Logger

	if sm.Config.DisableLogs {
		lg = log.NewLoggerEmpty()
	} else {
		lg = log.NewLoggerBase("execute", sm.Config.LogDirectory, false)
	}

	lg.Append("log started for workflow_run", log.LEVEL_SYSTEM)
	state := golua.NewState(golua.Options{
		SkipOpenLibs: false,
	})
	collection.CreateContext(state)

	runner := lua.NewRunner(state, &lg, sm.CliMode)
	runner.Config = sm.Config
	runner.Entry = name

	defer func() {
		if r := recover(); r != nil {
			lg.Append(fmt.Sprintf("panic recovered: %+v", r), log.LEVEL_ERROR)

			if sm.CliMode {
				sm.SetState(STATE_EXIT)
			} else {
				WorkflowFailEnter(sm, WorkflowFailData{
					Name:  pth,
					Error: fmt.Errorf(runner.Failed),
				})
			}
		}

		lg.Close()
	}()

	luaPth := path.Join(sm.Config.WorkflowDirectory, pth)

	err := runner.Run(luaPth, lib.Builtins)
	runner.Wg.Wait()

	lg.Append("All collections empty, exiting", log.LEVEL_SYSTEM)

	runner.CR_WIN.CleanAll()
	runner.CR_REF.CleanAll()
	runner.CR_GMP.CleanAll()
	runner.CR_LIP.CleanAll()
	runner.CR_TEA.CleanAll()
	runner.CR_CIM.CleanAll()

	runner.TC.CollectAll(state)
	runner.IC.CollectAll(state)
	runner.CC.CollectAll(state)
	runner.QR.CollectAll(state)

	if runner.Failed != "" {
		lg.Append(fmt.Sprintf("error occured while running script: %s", runner.Failed), log.LEVEL_ERROR)

		if sm.CliMode {
			fmt.Printf("error occured while running script: %s\n", runner.Failed)

			sm.SetState(STATE_EXIT)
			return nil
		}

		WorkflowFailEnter(sm, WorkflowFailData{
			Name:  pth,
			Error: fmt.Errorf(runner.Failed),
		})

		return nil
	}

	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)

		if sm.CliMode {
			fmt.Printf("error occured while running script: %s\n", err)

			sm.SetState(STATE_EXIT)
			return nil
		}

		WorkflowFailEnter(sm, WorkflowFailData{
			Name:  pth,
			Error: err,
		})

		return nil
	}

	ert := collErr(runner.TC.Errs, "TC", pth, &lg, sm)
	eri := collErr(runner.IC.Errs, "IC", pth, &lg, sm)
	erc := collErr(runner.CC.Errs, "CC", pth, &lg, sm)
	erq := collErr(runner.QR.Errs, "QR", pth, &lg, sm)
	if ert || eri || erc || erq {
		if sm.CliMode {
			sm.SetState(STATE_EXIT)
			return fmt.Errorf("error running script")
		}

		WorkflowFailEnter(sm, WorkflowFailData{
			Name:  pth,
			Error: fmt.Errorf("error occurred within collection"),
		})

		return nil
	}

	lg.Append("workflow finished", log.LEVEL_INFO)

	if sm.CliMode {
		sm.SetState(STATE_EXIT)
		return nil
	}

	WorkflowFinishEnter(sm, pth)
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
				WorkflowFailEnter(sm, WorkflowFailData{
					Name:  script,
					Error: err,
				})
			}

			errExists = true
		}
	}

	return errExists
}
