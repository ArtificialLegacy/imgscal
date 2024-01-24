package statemachine

type StateStepFunction func(sm *StateMachine)

type State struct {
	id          CliState
	step        StateStepFunction
	connections []CliState
}

func NewState(id CliState, step StateStepFunction, connections []CliState) State {
	return State{
		id:          id,
		step:        step,
		connections: connections,
	}
}
