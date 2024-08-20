package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

type WorkflowConfirmData struct {
	Workflow *workflow.Workflow
	Entry    string
	Base     string
}

func WorkflowConfirmEnter(sm *statemachine.StateMachine, data WorkflowConfirmData) {
	sm.SetState(STATE_WORKFLOW_CONFIRM)
	sm.Data = data
}

func WorkflowConfirm(sm *statemachine.StateMachine) error {
	data := sm.Data.(WorkflowConfirmData)
	sm.Data = nil
	workflow := data.Workflow
	pth := data.Entry

	autoConfirm := sm.CliMode || sm.Config.AlwaysConfirm

	var answer string

	if !autoConfirm {
		fmt.Printf("\n%s%s%s [%s] by %s.\n", cli.COLOR_BOLD, workflow.Name, cli.COLOR_RESET, workflow.Version, workflow.Author)
		fmt.Printf("%s%s%s\n\n", configPathColor, pth, cli.COLOR_RESET)
		fmt.Printf("%s\n\n", workflow.Desc)

		var err error
		answer, err = cli.Question(
			fmt.Sprintf("Do you wish to run the above workflow? %s(Y)%s/%s%sN%s", cli.COLOR_GREEN, cli.COLOR_RESET, cli.COLOR_BOLD, cli.COLOR_RED, cli.COLOR_RESET),
			cli.QuestionOptions{
				Normalize: true,
				Accepts:   []string{"y", "n"},
				Fallback:  "y",
			},
		)

		if err != nil {
			return fmt.Errorf("confirmation aborted from err during prompt: %s", err)
		}
	} else {
		answer = "y"
	}

	switch answer {
	case "y":
		WorkflowRunEnter(sm, pth)
	case "n":
		sm.SetState(STATE_WORKFLOW_LIST)
	default:
		panic("Impossible answer provided.")
	}

	return nil
}
