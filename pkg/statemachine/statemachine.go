package statemachine

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/config"
)

type StateFunction func(sm *StateMachine) error

const STACK_SIZE = 64

type StateMachine struct {
	states       []StateFunction
	currentState int

	stack        []any
	stackPointer int

	Config  *config.Config
	CliMode bool
}

func NewStateMachine(stateCount int) *StateMachine {
	return &StateMachine{
		states:       make([]StateFunction, stateCount),
		currentState: 0,
		stack:        make([]any, STACK_SIZE),
		stackPointer: -1,
	}
}

func (sm *StateMachine) AddState(id int, fn StateFunction) {
	sm.states[id] = fn
}

func (sm *StateMachine) SetState(state int) {
	sm.currentState = state
}

func (sm *StateMachine) push(value any) {
	if sm.stackPointer == STACK_SIZE-1 {
		panic("Attempting to push value to stack when stack is full.")
	}

	sm.stackPointer++
	sm.stack[sm.stackPointer] = value
}

func (sm *StateMachine) pop() any {
	if sm.stackPointer == -1 {
		panic("Attempting to pop from stack when stack is empty.")
	}

	val := sm.stack[sm.stackPointer]
	sm.stackPointer--

	return val
}

func (sm *StateMachine) Peek() bool {
	return sm.stackPointer >= 0
}

func (sm *StateMachine) PushInt(value int) {
	sm.push(value)
}

func (sm *StateMachine) PushString(value string) {
	sm.push(value)
}

func (sm *StateMachine) PushBool(value bool) {
	sm.push(value)
}

func (sm *StateMachine) PopInt() int {
	val := sm.pop()

	switch v := val.(type) {
	case int:
		return v
	default:
		panic(fmt.Sprintf("Attemping to pop a non-int off the stack as an int. got=%T", val))
	}
}

func (sm *StateMachine) PopString() string {
	val := sm.pop()

	switch v := val.(type) {
	case string:
		return v
	default:
		panic(fmt.Sprintf("Attemping to pop a non-string off the stack as an string. got=%T", val))
	}
}

func (sm *StateMachine) PopBool() bool {
	val := sm.pop()

	switch v := val.(type) {
	case bool:
		return v
	default:
		panic(fmt.Sprintf("Attemping to pop a non-bool off the stack as a bool. got=%T", val))
	}
}

func (sm *StateMachine) Step() error {
	return sm.states[sm.currentState](sm)
}
