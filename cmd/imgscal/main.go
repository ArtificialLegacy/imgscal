//go:generate goversioninfo -icon=assets/favicon.ico -manifest=imgscal.exe.manifest

package main

import (
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/states"
)

func main() {
	sm := statemachine.NewStateMachine(states.STATE_COUNT)

	sm.AddState(states.STATE_MAIN, states.MainMenu)
	sm.AddState(states.STATE_EXIT, states.Exit)
	sm.AddState(states.STATE_WORKFLOW_LIST, states.WorkflowList)
	sm.AddState(states.STATE_WORKFLOW_CONFIRM, states.WorkflowConfirm)
	sm.AddState(states.STATE_WORKFLOW_FAIL_LOAD, states.WorkflowFailLoad)
	sm.AddState(states.STATE_WORKFLOW_RUN, states.WorkflowRun)

	for {
		sm.Step()
	}
}
