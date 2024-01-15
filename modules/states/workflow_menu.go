package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var workflowMenuEnter statemachine.StateEnterFunction = func(from statemachine.CliState, sm *statemachine.StateMachine) {
	cli.Clear()

	response, _ := cli.Menu("Select workflow to run", []string{
		"Real-ESRGAN-x4plus",
		"Real-ESRGAN-x4plus-anime",
		"Back",
	})

	switch response {
	case 0:
		sm.Transition(statemachine.ESRGAN_X4)
		return
	case 1:
		sm.Transition(statemachine.ESRGAN_ANIMEX4)
		return
	case 2:
		sm.Transition(statemachine.LANDING_MENU)
		return
	}
}

var WorkflowMenu = statemachine.NewState(statemachine.WORKFLOW_MENU, workflowMenuEnter, nil, []statemachine.CliState{statemachine.ESRGAN_X4, statemachine.ESRGAN_ANIMEX4, statemachine.LANDING_MENU})
