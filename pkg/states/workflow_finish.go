package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowFinish(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()

	cli.Question(fmt.Sprintf("Script %s%s%s%s ran successfully...", cli.COLOR_GREEN, script, cli.COLOR_RESET, cli.COLOR_BOLD), cli.QuestionOptions{})

	sm.SetState(STATE_WORKFLOW_LIST)
	return nil
}
