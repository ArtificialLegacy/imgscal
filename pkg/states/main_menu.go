package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

const (
	MAIN_MENU_OPTION_WORKFLOW int = iota
	MAIN_MENU_OPTION_EXIT
)

var options = []string{
	MAIN_MENU_OPTION_WORKFLOW: "Run Workflow",
	MAIN_MENU_OPTION_EXIT:     "Exit",
}

func MainMenu(setState statemachine.SetStateFunction) error {
	cli.Clear()

	result, err := cli.SelectMenu("Imgscal", options)
	if err != nil {
		return err
	}

	switch result {
	case MAIN_MENU_OPTION_WORKFLOW:
		setState(STATE_WORKFLOW_LIST)

	case MAIN_MENU_OPTION_EXIT:
		setState(STATE_EXIT)

	default:
		panic(fmt.Sprintf("MAIN_MENU_OPTION %d is not handled.", result))
	}

	return nil
}
