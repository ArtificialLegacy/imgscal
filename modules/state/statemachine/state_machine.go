package statemachine

import (
	"errors"
	"fmt"

	"github.com/ArtificialLegacy/imgscal/modules/workflow"
)

// A finite state machine used to control the flow of the program.
type StateMachine struct {
	states       map[CliState]State
	current      CliState
	programState struct {
		workflows       map[string]*workflow.Workflow
		currentWorkflow string
	}
}

// Initializes a new state machine, with no states.
// The initial state is set to "".
func NewStateMachine() *StateMachine {
	return &StateMachine{
		states:  make(map[CliState]State),
		current: NONE,
	}
}

func (sm *StateMachine) AddStates(states []State) {
	for _, state := range states {
		sm.states[state.id] = state
	}
}

func (sm *StateMachine) Step() {
	sm.states[sm.current].step(sm)
}

func (sm *StateMachine) Transition(to CliState) error {
	_, exists := sm.states[to]
	if !exists {
		return errors.New(fmt.Sprintf("State %d does not exist.", to))
	}

	if sm.current == NONE {
		sm.current = to
		return nil
	}

	connections := sm.states[sm.current].connections
	var found bool
	for _, connection := range connections {
		if connection == to {
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("No connection from %d to %d.", sm.current, to))
	}

	sm.current = to
	return nil
}

func (sm *StateMachine) SetWorkflowsState(workflows map[string]*workflow.Workflow) {
	sm.programState.workflows = workflows
}

func (sm *StateMachine) GetWorkflowsState() map[string]*workflow.Workflow {
	return sm.programState.workflows
}

func (sm *StateMachine) SetCurrentWorkflowState(workflow string) {
	sm.programState.currentWorkflow = workflow
}

func (sm *StateMachine) GetCurrentWorkflowState() string {
	return sm.programState.currentWorkflow
}
