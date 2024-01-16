package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var landingMenuEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	cli.Clear()

	response, _ := cli.Menu("Select task to perform", []string{
		"Run Workflow",
		"Manage Real-ESRGAN",
		"Exit",
	})

	switch response {
	case 0:
		sm.Transition(statemachine.WORKFLOW_MENU)
		return
	case 1:
		sm.Transition(statemachine.ESRGAN_MANAGE)
		return
	case 2:
		os.Exit(0)
		return
	}
}

var LandingMenu = statemachine.NewState(
	statemachine.LANDING_MENU,
	landingMenuEnter,
	[]statemachine.CliState{statemachine.WORKFLOW_MENU, statemachine.ESRGAN_MANAGE},
)
