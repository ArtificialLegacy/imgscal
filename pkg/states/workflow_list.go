package states

import (
	"fmt"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

func WorkflowList(sm *statemachine.StateMachine) error {
	cli.Clear()

	workflows, errList, err := workflow.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		return err
	}

	if errList != nil && len(*errList) > 0 {
		lg := log.NewLoggerBase("scan", sm.Config.LogDirectory, false)
		defer lg.Close()
		lg.Append("Encountered errors while scanning for workflows: ", log.LEVEL_SYSTEM)
		for _, e := range *errList {
			lg.Append(e.Error(), log.LEVEL_ERROR)
		}
	}

	if len(*workflows) == 0 {
		fmt.Printf("\nWorkflow directory empty, nothing to run.\n")
		fmt.Printf("%s%s%s\n\n", configPathColor, sm.Config.WorkflowDirectory, cli.COLOR_RESET)

		fmt.Printf(" > Try \u001b[48;5;234mmake install-examples%s\n\n", cli.COLOR_RESET)

		cli.Question("Press any key to continue...", cli.QuestionOptions{})
		sm.SetState(STATE_MAIN)
		return nil
	}

	options := []string{}
	optionsWorkflows := []*workflow.Workflow{}
	optionsPaths := []string{}

	for _, w := range *workflows {
		starUsed := false
		for s, ws := range w.Workflows {
			if s == "*" {
				if starUsed {
					continue
				}
				starUsed = true
				options = append(options, w.Name)
			} else {
				options = append(options, w.Name+"/"+s)
			}
			optionsWorkflows = append(optionsWorkflows, w)
			optionsPaths = append(optionsPaths, path.Join(path.Dir(w.Base), ws))
		}
	}

	options = append(options, fmt.Sprintf("%sReturn%s", cli.COLOR_RED, cli.COLOR_RESET))

	result, err := cli.SelectMenu(fmt.Sprintf("Select %sworkflow%s to run.", cli.COLOR_BOLD, cli.COLOR_RESET), options)
	if err != nil {
		return err
	}

	if result == len(options)-1 {
		sm.SetState(STATE_MAIN)
	} else {
		WorkflowConfirmEnter(sm, WorkflowConfirmData{
			Workflow: optionsWorkflows[result],
			Entry:    optionsPaths[result],
			Name:     options[result],
		})
	}

	return nil
}
