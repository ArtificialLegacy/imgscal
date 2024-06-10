package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowFailRun(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	err := sm.PopString()

	fmt.Printf("\n%s\n\n", err)

	cli.Question(fmt.Sprintf("Script %s%s%s%s failed to run...", cli.COLOR_RED, script, cli.COLOR_RESET, cli.COLOR_BOLD), cli.QuestionOptions{})

	sm.SetState(STATE_WORKFLOW_LIST)
	return nil
}
