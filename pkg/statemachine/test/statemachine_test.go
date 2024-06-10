package statemachine_test

import (
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

func createStateMachine() *statemachine.StateMachine {
	sm := statemachine.NewStateMachine(2)

	sm.AddState(0, func(sm *statemachine.StateMachine) error { return nil })
	sm.AddState(1, func(sm *statemachine.StateMachine) error { return nil })

	return sm
}

func TestStackInt(t *testing.T) {
	sm := createStateMachine()

	sm.PushInt(10)
	i := sm.PopInt()
	if i != 10 {
		t.Errorf("Incorrect int returned from stack; got=%d, expected=%d", i, 10)
	}
}

func TestStackString(t *testing.T) {
	sm := createStateMachine()

	sm.PushString("test")
	s := sm.PopString()
	if s != "test" {
		t.Errorf("Incorrect string returned from stack; got=%s, expected=%s", s, "test")
	}
}

func TestStackPeek(t *testing.T) {
	sm := createStateMachine()

	sm.PushString("test")
	notEmpty := sm.Peek()
	if !notEmpty {
		t.Error("Stack peek returned empty when there should be an entry.")
	}

	sm.PopString()
	empty := sm.Peek()
	if empty {
		t.Error("Stack peek returned not empty when there shouldn't be an entry.")
	}
}
