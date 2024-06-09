package states

import (
	"fmt"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

func WorkflowConfirm(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	wf := workflow.NewWorkflow()

	state := lua.WorkflowConfigState(&wf)
	runner := lua.NewRunner(state, &struct{}{})
	err := runner.Run(script)

	if err != nil || len(wf.Requires) == 0 || wf.Version == "" || wf.Name == "" {
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_LOAD)
		return nil
	}

	fmt.Printf("\n%s%s%s [%s] by %s.\n\n", cli.COLOR_BOLD, wf.Name, cli.COLOR_RESET, wf.Version, wf.Author)
	fmt.Printf("%s\n\n", wf.Desc)
	fmt.Printf("Required plugins: \n - %s\n\n;", strings.Join(wf.Requires, "\n - "))

	answer, err := cli.Question(
		fmt.Sprintf("Do you wish to run the above workflow? %s(Y)%s/%sN%s", cli.COLOR_GREEN, cli.COLOR_RESET, cli.COLOR_RED, cli.COLOR_RESET),
		cli.QuestionOptions{
			Normalize: true,
			Accepts:   []string{"y", "n"},
			Fallback:  "y",
		},
	)

	if err != nil {
		return err
	}

	switch answer {
	case "y":
		for _, s := range wf.Requires {
			sm.PushString(s)
		}
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_RUN)
	case "n":
		sm.SetState(STATE_WORKFLOW_LIST)
	default:
		panic("Impossible answer provided.")
	}

	return nil
}
