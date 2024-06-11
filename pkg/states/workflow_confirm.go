package states

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	golua "github.com/Shopify/go-lua"
)

func WorkflowConfirm(sm *statemachine.StateMachine) error {
	cli.Clear()

	script := sm.PopString()
	wf := workflow.NewWorkflow()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	lg := log.NewLogger()
	defer lg.Dump("./log")

	lg.Append("log started for workflow_confirm", log.LEVEL_INFO)
	state := lua.WorkflowConfigState(&wf, &lg)
	err = golua.DoFile(state, path.Join(pwd, script))

	if err != nil {
		lg.Append(fmt.Sprintf("error occured while running script: %s", err), log.LEVEL_ERROR)
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_FAIL_LOAD)
		return nil
	}

	if len(wf.Requires) == 0 {
		lg.Append("requires array was empty", log.LEVEL_WARN)
	}
	if wf.Version == "" {
		lg.Append("version field was empty", log.LEVEL_WARN)
	}
	if wf.Name == "" {
		lg.Append("name field was empty", log.LEVEL_WARN)
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
		lg.Append(fmt.Sprintf("confirmation aborted from err during prompt: %s", err), log.LEVEL_ERROR)
		return err
	}

	switch answer {
	case "y":
		for _, s := range wf.Requires {
			sm.PushString(s)
		}
		sm.PushString(script)
		sm.SetState(STATE_WORKFLOW_RUN)
		lg.Append("confirmation answer y", log.LEVEL_INFO)
	case "n":
		sm.SetState(STATE_WORKFLOW_LIST)
		lg.Append("confirmation answer n", log.LEVEL_INFO)
	default:
		lg.Append("impossible answer provided", log.LEVEL_ERROR)
		panic("Impossible answer provided.")
	}

	return nil
}
