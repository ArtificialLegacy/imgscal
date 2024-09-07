package states

import (
	"fmt"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	golua "github.com/yuin/gopher-lua"
)

func WorkflowHelpEnter(sm *statemachine.StateMachine, name string) {
	sm.SetState(STATE_WORKFLOW_HELP)
	sm.Data = name
}

func WorkflowHelp(sm *statemachine.StateMachine) error {
	name := sm.Data.(string)

	wf, _, err := workflow.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		fmt.Printf("failed to scan for workflows: %s\n", err)
		sm.SetState(STATE_EXIT)
		return err
	}

	foundPath := ""
	var foundWf *workflow.Workflow
	for _, w := range *wf {
		if w.Name == name {
			found, ok := w.CliWorkflows["*"]
			if !ok {
				fmt.Printf("cannot use workflow base name when there is no star workflow: %s\n", err)
				sm.SetState(STATE_EXIT)
				return err
			}

			foundPath = path.Join(w.Location, found)
			foundWf = w
			break
		}

		if !strings.HasPrefix(name, path.Base(path.Dir(w.Base))) {
			continue
		}

		found, ok := w.CliWorkflows[strings.TrimPrefix(name, path.Base(path.Dir(w.Base))+"/")]
		if !ok {
			continue
		}
		foundPath = path.Join(w.Location, found)
		foundWf = w
		break
	}

	var lg log.Logger
	if sm.Config.DisableLogs {
		lg = log.NewLoggerEmpty()
	} else {
		lg = log.NewLoggerBase("help", sm.Config.LogDirectory, false)
	}
	defer lg.Close()

	lg.Append("log started for workflow_help", log.LEVEL_SYSTEM)
	state := golua.NewState()
	runner := lua.NewRunner(state, &lg, sm.CliMode)
	str, err := runner.Help(foundPath, foundWf)
	if err != nil {
		lg.Append(fmt.Sprintf("%s", err), log.LEVEL_ERROR)

		fmt.Printf("error: %s\n", err)
		fmt.Printf("for: %s\n", foundPath)

		sm.SetState(STATE_EXIT)
		return nil
	}

	fmt.Printf("Result of help for %s:\n\n", name)
	fmt.Printf("%s\n\n", str)

	sm.SetState(STATE_EXIT)
	return nil
}
