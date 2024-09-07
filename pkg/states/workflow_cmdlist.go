package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

func WorkflowCMDList(sm *statemachine.StateMachine) error {
	cli.Clear()

	wf, _, err := workflow.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		fmt.Printf("failed to scan for workflows: %s\n", err)
		sm.SetState(STATE_EXIT)
		return err
	}

	for _, w := range *wf {
		fmt.Printf("> %s\n", w.Name)
		for v := range w.CliWorkflows {
			fmt.Printf("  - %s\n", v)
		}
	}

	sm.SetState(STATE_EXIT)
	return nil
}
