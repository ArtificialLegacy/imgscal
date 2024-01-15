package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
)

var esrganVerifyEnter statemachine.StateEnterFunction = func(from statemachine.CliState, sm *statemachine.StateMachine) {
	exists := esrgan.Verify()
	if exists {
		sm.Transition(statemachine.LANDING_MENU)
		return
	}

	cli.Clear()

	response, _ := cli.Question(fmt.Sprintf("%s!%s Real-ESRGAN was not found. Would you like to download it? (%sY%s/%sN%s) ", cli.RED, cli.RESET, cli.GREEN, cli.RESET, cli.RED, cli.RESET), cli.QuestionOptions{
		Normalize: true,
		Accepts:   []string{"y", "n"},
		Fallback:  "n",
	})

	if response == "y" {
		sm.Transition(statemachine.ESRGAN_DOWNLOAD)
		return
	} else if response == "n" {
		sm.Transition(statemachine.ESRGAN_FAIL)
		return
	}
}

var ESRGANVerify = statemachine.NewState(statemachine.ESRGAN_VERIFY, esrganVerifyEnter, nil, []statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU})
