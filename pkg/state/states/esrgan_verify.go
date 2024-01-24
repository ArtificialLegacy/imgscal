package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/libs/esrgan"
	"github.com/ArtificialLegacy/imgscal/pkg/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/utility/cli"
)

var esrganVerifyEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
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

var ESRGANVerify = statemachine.NewState(
	statemachine.ESRGAN_VERIFY,
	esrganVerifyEnter,
	[]statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU},
)
