package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var esrganVerifyEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	exists := esrgan.Verify()
	if exists {
		transition(statemachine.LANDING_MENU)
		return
	}

	cli.Clear()

	response, _ := cli.Question(fmt.Sprintf("%s!%s Real-ESRGAN was not found. Would you like to download it? (%sY%s/%sN%s) ", cli.RED, cli.RESET, cli.GREEN, cli.RESET, cli.RED, cli.RESET), cli.QuestionOptions{
		Normalize: true,
		Accepts:   []string{"y", "n"},
		Fallback:  "n",
	})

	if response == "y" {
		transition(statemachine.ESRGAN_DOWNLOAD)
		return
	} else if response == "n" {
		transition(statemachine.ESRGAN_FAIL)
		return
	}
}

var ESRGANVerify = statemachine.NewState(statemachine.ESRGAN_VERIFY, esrganVerifyEnter, nil, []statemachine.CliState{statemachine.ESRGAN_DOWNLOAD, statemachine.ESRGAN_FAIL, statemachine.LANDING_MENU})
