package customtea

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

type TeaMSG int

const (
	MSG_NONE TeaMSG = iota
	MSG_KEY
	MSG_SPINNERTICK
	MSG_BLINK
	MSG_STOPWATCHRESET
	MSG_STOPWATCHSTARTSTOP
	MSG_STOPWATCHTICK
	MSG_TIMERSTARTSTOP
	MSG_TIMERTICK
	MSG_TIMERTIMEOUT
)

func BuildMSG(msg tea.Msg, state *golua.LState) *golua.LTable {
	var luaMsg *golua.LTable

	switch msg := msg.(type) {
	case tea.KeyMsg:
		luaMsg = msgTableKey(state, msg.String())
	case spinner.TickMsg:
		luaMsg = msgTableSpinnerTick(state, msg.ID)
	case cursor.BlinkMsg:
		luaMsg = msgTableCursorBlink(state)
	case stopwatch.ResetMsg:
		luaMsg = msgTableStopWatchReset(state, msg.ID)
	case stopwatch.StartStopMsg:
		luaMsg = msgTableStopWatchStartStop(state, msg.ID)
	case stopwatch.TickMsg:
		luaMsg = msgTableStopWatchTick(state, msg.ID)
	case timer.StartStopMsg:
		luaMsg = msgTableTimerStartStop(state, msg.ID)
	case timer.TickMsg:
		luaMsg = msgTableTimerTick(state, msg.ID, msg.Timeout)
	case timer.TimeoutMsg:
		luaMsg = msgTableTimerTimeout(state, msg.ID)
	default:
		luaMsg = msgTableNone(state)
	}

	return luaMsg
}

func msgTableNone(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_NONE))

	return t
}

func msgTableKey(state *golua.LState, key string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_KEY))
	t.RawSetString("key", golua.LString(key))

	return t
}

func msgTableSpinnerTick(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_SPINNERTICK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableCursorBlink(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_BLINK))

	return t
}

func msgTableStopWatchReset(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_STOPWATCHRESET))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableStopWatchStartStop(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_STOPWATCHSTARTSTOP))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableStopWatchTick(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_STOPWATCHTICK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableTimerStartStop(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_TIMERSTARTSTOP))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableTimerTimeout(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_TIMERTIMEOUT))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func msgTableTimerTick(state *golua.LState, id int, timeout bool) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_TIMERTICK))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("timeout", golua.LBool(timeout))

	return t
}
