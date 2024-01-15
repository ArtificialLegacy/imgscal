package statemachine

type StateEnterFunction func(from CliState, sm *StateMachine)

type StateExitFunction func(to CliState, sm *StateMachine)

type State struct {
	id          CliState
	enter       StateEnterFunction
	exit        StateExitFunction
	connections []CliState
}

func NewState(id CliState, enter StateEnterFunction, exit StateExitFunction, connections []CliState) State {
	return State{
		id:          id,
		enter:       enter,
		exit:        exit,
		connections: connections,
	}
}
