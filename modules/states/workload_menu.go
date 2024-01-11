package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var workloadMenuEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	cli.Clear()

	response, _ := cli.Menu("Select workload to run", []string{
		"Real-ESRGAN-x4plus",
		"Real-ESRGAN-x4plus-anime",
		"Back",
	})

	switch response {
	case 0:
		transition(statemachine.ESRGAN_X4)
		return
	case 1:
		transition(statemachine.ESRGAN_ANIMEX4)
		return
	case 2:
		transition(statemachine.LANDING_MENU)
		return
	}
}

var WorkloadMenu = statemachine.NewState(statemachine.WORKLOAD_MENU, workloadMenuEnter, nil, []statemachine.CliState{statemachine.ESRGAN_X4, statemachine.ESRGAN_ANIMEX4, statemachine.LANDING_MENU})
