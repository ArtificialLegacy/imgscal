package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var workloadFinishEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	print("\n")

	cli.Question("Workload finished. Press enter to continue...", cli.QuestionOptions{
		Normalize: false,
		Accepts:   []string{},
		Fallback:  "",
	})

	transition(statemachine.LANDING_MENU)
}

var WorkloadFinish = statemachine.NewState(statemachine.WORKLOAD_FINISH, workloadFinishEnter, nil, []statemachine.CliState{statemachine.LANDING_MENU})
