package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func Exit(setState *statemachine.StateMachine) error {
	os.Exit(0)
	return nil
}
