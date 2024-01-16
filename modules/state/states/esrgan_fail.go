package states

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/modules/utility/cli"
)

var esrganFailEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	println(fmt.Sprintf("\n%sCannot continue without ESRGAN, restart the program to attempt to install.%s\n", cli.RED, cli.RESET))

	os.Exit(1)
}

var ESRGANFail = statemachine.NewState(statemachine.ESRGAN_FAIL, esrganFailEnter, []statemachine.CliState{})
