package states

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var esrganFailEnter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	println(fmt.Sprintf("\n%sCannot continue without ESRGAN, restart the program to attempt to install.%s\n", cli.RED, cli.RESET))

	os.Exit(1)
}

var ESRGANFail = statemachine.NewState(statemachine.ESRGAN_FAIL, esrganFailEnter, nil, []statemachine.CliState{})
