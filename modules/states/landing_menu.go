package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var landingMenuEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	cli.Clear()

	response, _ := cli.Menu("Select task to perform", []string{
		"Run Workload",
		"Manage Real-ESRGAN",
		"Exit",
	})

	switch response {
	case 0:
		transition(statemachine.WORKLOAD_MENU)
		return
	case 1:
		transition(statemachine.ESRGAN_MANAGE)
		return
	case 2:
		os.Exit(0)
		return
	}
}

var LandingMenu = statemachine.NewState(statemachine.LANDING_MENU, landingMenuEnter, nil, []statemachine.CliState{statemachine.WORKLOAD_MENU, statemachine.ESRGAN_MANAGE})
