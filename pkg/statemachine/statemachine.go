package statemachine

type SetStateFunction func(state int)

type StateFunction func(setState SetStateFunction) error

type StateMachine struct {
	states       []StateFunction
	currentState int
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
	return sm.states[sm.currentState](sm.SetState)
}
