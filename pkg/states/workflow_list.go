package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/script"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowList(sm *statemachine.StateMachine) error {
	cli.Clear()

	scripts, err := script.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		return err
	}

	options := []string{}

	for _, s := range scripts {
		options = append(options, s.Name)
	}

	options = append(options, fmt.Sprintf("%sReturn%s", cli.COLOR_RED, cli.COLOR_RESET))

	result, err := cli.SelectMenu(fmt.Sprintf("Select %sworkflow%s to run.", cli.COLOR_BOLD, cli.COLOR_RESET), options)
	if err != nil {
		return err
	}

	if result == len(options)-1 {
		sm.SetState(STATE_MAIN)
	} else {
		sm.PushString(scripts[result].Filepath)
		sm.SetState(STATE_WORKFLOW_CONFIRM)
	}

	return nil
}
