package states

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
	"github.com/ArtificialLegacy/imgscal/modules/workflow"
)

var workflowRunEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	currentWorkflow := sm.GetCurrentWorkflowState()
	wf := sm.GetWorkflowsState()[currentWorkflow]

	filepath, err := workflow.WorkflowBegin()
	if err != nil {
		fmt.Printf("%s! The path you entered does not exist. Please try again.%s\n", cli.RED, cli.RESET)
		sm.Transition(statemachine.WORKFLOW_FINISH)
	}

	file, _ := os.Stat(filepath)
	if file.IsDir() {
		files, err := os.ReadDir(filepath)
		if err != nil {
			fmt.Printf("%s! The path you entered is not valid. Please try again.%s\n", cli.RED, cli.RESET)
			sm.Transition(statemachine.WORKFLOW_FINISH)
			return
		}

		for index, file := range files {
			println(fmt.Sprintf("%s!%s Running %s on %s. (Image %d of %d)", cli.CYAN, cli.RESET, wf.Config.Name, file.Name(), index, len(files)))
			wf.Run(filepath+"\\"+file.Name(), file.Name(), wf.Config.Requires)
		}
	} else {
		fileSplit := strings.Split(filepath, "\\")
		filename := fileSplit[len(fileSplit)-1]

		fmt.Printf("%s!%s Running %s on %s.", cli.CYAN, cli.RESET, wf.Config.Name, file.Name())
		wf.Run(filepath, filename, wf.Config.Requires)
	}

	sm.Transition(statemachine.WORKFLOW_FINISH)
}

var WorkflowRun = statemachine.NewState(
	statemachine.WORKFLOW_RUN,
	workflowRunEnter,
	[]statemachine.CliState{statemachine.WORKFLOW_FINISH},
)
