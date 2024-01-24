package states

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/utility/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

var workflowRunEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	currentWorkflow := sm.GetCurrentWorkflowState()
	wf := sm.GetWorkflowsState()[currentWorkflow]
	pwd, _ := os.Getwd()

	imgFile := fmt.Sprintf("%s\\temp", pwd)

	if _, err := os.Stat(fmt.Sprintf("%s\\temp", pwd)); os.IsNotExist(err) {
		os.Mkdir(imgFile, 0777)
	}

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
			fmt.Printf("%s!%s Running %s on %s. (Image %d of %d)\n", cli.CYAN, cli.RESET, wf.Config.Name, file.Name(), index+1, len(files))
			wf.Run(filepath+"\\"+file.Name(), file.Name(), wf.Config.Requires)
		}
	} else {
		fileSplit := strings.Split(filepath, "\\")
		filename := fileSplit[len(fileSplit)-1]

		fmt.Printf("%s!%s Running %s on %s.\n", cli.CYAN, cli.RESET, wf.Config.Name, file.Name())
		wf.Run(filepath, filename, wf.Config.Requires)
	}

	tempFiles, _ := os.ReadDir(imgFile)
	for _, file := range tempFiles {
		os.Remove(fmt.Sprintf("%s\\temp\\%s\n", pwd, file.Name()))
	}
	os.Remove(imgFile)

	sm.Transition(statemachine.WORKFLOW_FINISH)
}

var WorkflowRun = statemachine.NewState(
	statemachine.WORKFLOW_RUN,
	workflowRunEnter,
	[]statemachine.CliState{statemachine.WORKFLOW_FINISH},
)
