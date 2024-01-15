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
		workflows map[string]*workflow.Workflow
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

// Adds a state to the state machine.
func (sm *StateMachine) AddState(state State) *StateMachine {
	sm.states[state.id] = state
	return sm
}

// Transitions to the given state from the current state.
// Checks if the state exists, and if there is a connection from the current state to the given state.
// Calls the exit function of the current state, and the enter function of the given state, passing the previous state and the transition function.
// If the current state is "", the enter function of the given state is called, passing "" as the previous state and the transition function, and no exit function is called, and no connection is checked.
func (sm *StateMachine) Transition(to CliState) error {
	toState, exists := sm.states[to]
	if !exists {
		return errors.New(fmt.Sprintf("State %d does not exist.", to))
	}

	if sm.current == NONE {
		sm.current = to
		toState.enter(NONE, sm)
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

	if sm.states[sm.current].exit != nil {
		sm.states[sm.current].exit(to, sm)
	}
	prev := sm.current
	sm.current = to
	toState.enter(prev, sm)

	return nil
}

func (sm *StateMachine) SetWorkflowState(workflows map[string]*workflow.Workflow) {
	sm.programState.workflows = workflows
}

func (sm *StateMachine) GetWorkflowState() map[string]*workflow.Workflow {
	return sm.programState.workflows
}
