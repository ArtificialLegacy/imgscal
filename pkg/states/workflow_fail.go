package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

type WorkflowFailData struct {
	Name  string
	Error error
}

func WorkflowFailEnter(sm *statemachine.StateMachine, data WorkflowFailData) {
	sm.SetState(STATE_WORKFLOW_FAIL)
	sm.Data = data
}

func WorkflowFail(sm *statemachine.StateMachine) error {
	data := sm.Data.(WorkflowFailData)
	sm.Data = nil
	script := data.Name
	err := data.Error

	fmt.Printf("\n%s\n\n", err)

	cli.Question(fmt.Sprintf("Script %s%s%s%s failed to run...", cli.COLOR_RED, script, cli.COLOR_RESET, cli.COLOR_BOLD), cli.QuestionOptions{})

	sm.SetState(STATE_WORKFLOW_LIST)
	return nil
}
