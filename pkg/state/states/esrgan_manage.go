package states

import (
	"github.com/ArtificialLegacy/imgscal/pkg/libs/esrgan"
	"github.com/ArtificialLegacy/imgscal/pkg/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/utility/cli"
)

var esrganManageEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	cli.Clear()

	response, _ := cli.Menu("Select task to perform", []string{
		"Repair/Update Real-ESRGAN",
		"Uninstall Real-ESRGAN",
		"Back",
	})

	switch response {
	case 0:
		esrgan.Remove()
		sm.Transition(statemachine.ESRGAN_DOWNLOAD)
		return
	case 1:
		esrgan.Remove()
		sm.Transition(statemachine.ESRGAN_FAIL)
		return
	case 2:
		sm.Transition(statemachine.LANDING_MENU)
		return
	}
}

var ESRGANManage = statemachine.NewState(
	statemachine.ESRGAN_MANAGE,
	esrganManageEnter,
	[]statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU},
)
