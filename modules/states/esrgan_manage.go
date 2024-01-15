package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var esrganManageEnter statemachine.StateEnterFunction = func(from statemachine.CliState, sm *statemachine.StateMachine) {
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

var ESRGANManage = statemachine.NewState(statemachine.ESRGAN_MANAGE, esrganManageEnter, nil, []statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU})
