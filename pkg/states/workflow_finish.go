package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowFinishEnter(sm *statemachine.StateMachine, script string) {
	sm.SetState(STATE_WORKFLOW_FINISH)
	sm.Data = script
}

func WorkflowFinish(sm *statemachine.StateMachine) error {
	script := sm.Data.(string)
	sm.Data = nil

	fmt.Printf("\n\n")

	cli.Question(fmt.Sprintf("Script %s%s%s%s ran successfully...", cli.COLOR_GREEN, script, cli.COLOR_RESET, cli.COLOR_BOLD), cli.QuestionOptions{})

	sm.SetState(STATE_WORKFLOW_LIST)
	return nil
}
