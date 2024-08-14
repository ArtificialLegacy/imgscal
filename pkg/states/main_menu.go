package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

const (
	MAIN_MENU_OPTION_WORKFLOW int = iota
	MAIN_MENU_OPTION_UTILITIES
	MAIN_MENU_OPTION_EXIT
)

var mainMenuOptions = []string{
	MAIN_MENU_OPTION_WORKFLOW:  "Run Workflow",
	MAIN_MENU_OPTION_UTILITIES: "Utilities",
	MAIN_MENU_OPTION_EXIT:      fmt.Sprintf("%sExit%s", cli.COLOR_RED, cli.COLOR_RESET),
}

func MainMenu(sm *statemachine.StateMachine) error {
	cli.Clear()

	result, err := cli.SelectMenu("Imgscal", mainMenuOptions)
	if err != nil {
		return err
	}

	switch result {
	case MAIN_MENU_OPTION_WORKFLOW:
		sm.SetState(STATE_WORKFLOW_LIST)

	case MAIN_MENU_OPTION_UTILITIES:
		sm.SetState(STATE_UTILITIES)

	case MAIN_MENU_OPTION_EXIT:
		sm.SetState(STATE_EXIT)

	default:
		panic(fmt.Sprintf("MAIN_MENU_OPTION %d is not handled.", result))
	}

	return nil
}
