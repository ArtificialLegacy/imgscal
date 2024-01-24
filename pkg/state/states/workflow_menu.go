package states

import (
	"github.com/ArtificialLegacy/imgscal/pkg/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/utility/cli"
)

var workflowMenuEnter statemachine.StateStepFunction = func(sm *statemachine.StateMachine) {
	cli.Clear()

	workflowList := []string{}
	workflowResponse := []string{}
	workflows := sm.GetWorkflowsState()

	for _, workflow := range workflows {
		workflowList = append(workflowList, workflow.Config.Name)
		workflowResponse = append(workflowResponse, workflow.File)
	}

	workflowList = append(workflowList, "Back")

	response, _ := cli.Menu("Select workflow to run", workflowList)

	if response == int8(len(workflowList)-1) {
		sm.Transition(statemachine.LANDING_MENU)
		return
	}

	sm.SetCurrentWorkflowState(workflowResponse[response])
	sm.Transition(statemachine.WORKFLOW_RUN)
}

var WorkflowMenu = statemachine.NewState(
	statemachine.WORKFLOW_MENU,
	workflowMenuEnter,
	[]statemachine.CliState{statemachine.WORKFLOW_RUN, statemachine.LANDING_MENU},
)
