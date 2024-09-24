package customtea

import (
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

func MouseEventTable(state *golua.LState, event tea.MouseEvent) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("x", golua.LNumber(event.X))
	t.RawSetString("y", golua.LNumber(event.Y))
	t.RawSetString("shift", golua.LBool(event.Shift))
	t.RawSetString("alt", golua.LBool(event.Alt))
	t.RawSetString("ctrl", golua.LBool(event.Ctrl))
	t.RawSetString("action", golua.LNumber(event.Action))
	t.RawSetString("button", golua.LNumber(event.Button))
	t.RawSetString("is_wheel", golua.LBool(event.IsWheel()))

	return t
}

func KeyEventTable(state *golua.LState, event tea.Key) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(event.Type))
	t.RawSetString("alt", golua.LBool(event.Alt))
	t.RawSetString("paste", golua.LBool(event.Paste))

	rt := state.NewTable()
	for i, v := range event.Runes {
		rt.RawSetInt(i+1, golua.LNumber(v))
	}
	t.RawSetString("runes", rt)

	return t
}
