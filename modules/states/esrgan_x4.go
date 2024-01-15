package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var esrganx4Enter statemachine.StateEnterFunction = func(from statemachine.CliState, sm *statemachine.StateMachine) {
	cli.Clear()

	answer, err := esrgan.WorkflowBegin()
	if err != nil {
		sm.Transition(statemachine.WORKFLOW_FINISH)
		return
	}

	file, _ := os.Stat(answer)

	if file.IsDir() {
		files, err := os.ReadDir(answer)
		if err != nil {
			sm.Transition(statemachine.WORKFLOW_FINISH)
			return
		}
		for index, file := range files {
			esrgan.X4(answer+"\\"+file.Name(), index+1, len(files))
		}
	} else {
		esrgan.X4(answer, 1, 1)
	}

	sm.Transition(statemachine.WORKFLOW_FINISH)
}

var ESRGANX4 = statemachine.NewState(statemachine.ESRGAN_X4, esrganx4Enter, nil, []statemachine.CliState{statemachine.WORKFLOW_FINISH})
