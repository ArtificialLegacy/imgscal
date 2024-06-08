package states

import (
	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/script"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func WorkflowList(setState statemachine.SetStateFunction) error {
	cli.Clear()

	scripts, err := script.WorkflowList()
	if err != nil {
		return err
	}

	options := []string{}

	for _, s := range scripts {
		options = append(options, s.Name)
	}

	options = append(options, "Return")

	result, err := cli.SelectMenu("Select workflow to run.", options)
	if err != nil {
		return err
	}

	if result == len(options)-1 {
		setState(STATE_MAIN)
	}

	return nil
}
