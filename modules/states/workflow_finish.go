package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var workflowFinishEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	print("\n")

	cli.Question("Workflow finished. Press enter to continue...", cli.QuestionOptions{
		Normalize: false,
		Accepts:   []string{},
		Fallback:  "",
	})

	sm.Transition(statemachine.LANDING_MENU)
}

var WorkflowFinish = statemachine.NewState(
	statemachine.WORKFLOW_FINISH,
	workflowFinishEnter,
	[]statemachine.CliState{statemachine.LANDING_MENU},
)
