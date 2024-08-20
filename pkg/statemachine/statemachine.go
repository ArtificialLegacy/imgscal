package statemachine

import (
	"github.com/ArtificialLegacy/imgscal/pkg/config"
)

type StateFunction func(sm *StateMachine) error

type StateMachine struct {
	states       []StateFunction
	currentState int

	Data any

	Config  *config.Config
	CliMode bool
}

func NewStateMachine(stateCount int) *StateMachine {
	return &StateMachine{
		states:       make([]StateFunction, stateCount),
		currentState: 0,
	}
}

func (sm *StateMachine) AddState(id int, fn StateFunction) {
	sm.states[id] = fn
}

func (sm *StateMachine) SetState(state int) {
	sm.currentState = state
}

func (sm *StateMachine) Step() error {
	return sm.states[sm.currentState](sm)
}
