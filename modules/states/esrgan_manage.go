package states

import (
	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var esrganManageEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	cli.Clear()

	response, _ := cli.Menu("Select task to perform", []string{
		"Repair/Update Real-ESRGAN",
		"Uninstall Real-ESRGAN",
		"Back",
	})

	switch response {
	case 0:
		esrgan.Remove()
		transition(statemachine.ESRGAN_DOWNLOAD)
		return
	case 1:
		esrgan.Remove()
		transition(statemachine.ESRGAN_FAIL)
		return
	case 2:
		transition(statemachine.LANDING_MENU)
		return
	}
}

var ESRGANManage = statemachine.NewState(statemachine.ESRGAN_MANAGE, esrganManageEnter, nil, []statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU})
