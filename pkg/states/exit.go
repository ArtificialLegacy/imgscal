package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func Exit(setState statemachine.SetStateFunction) error {
	os.Exit(0)
	return nil
}
