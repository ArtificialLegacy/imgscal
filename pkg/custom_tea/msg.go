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
	MSG_BLUR
	MSG_FOCUS
	MSG_QUIT
	MSG_RESUME
	MSG_SUSPEND
	MSG_WINDOWSIZE
	MSG_KEY
	MSG_MOUSE
	MSG_SPINNERTICK
	MSG_BLINK
	MSG_STOPWATCHRESET
	MSG_STOPWATCHSTARTSTOP
	MSG_STOPWATCHTICK
	MSG_TIMERSTARTSTOP
	MSG_TIMERTICK
	MSG_TIMERTIMEOUT
	MSG_LUA
)

func BuildMSG(msg tea.Msg, state *golua.LState) *golua.LTable {
	var luaMsg *golua.LTable

	switch msg := msg.(type) {
	case tea.KeyMsg:
		luaMsg = msgTableKey(state, tea.Key(msg))
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
	case tea.BlurMsg:
		luaMsg = msgTableSimple(state, MSG_BLUR)
	case tea.FocusMsg:
		luaMsg = msgTableSimple(state, MSG_FOCUS)
	case tea.MouseMsg:
		luaMsg = msgTableMouse(state, msg.String(), tea.MouseEvent(msg))
	case tea.QuitMsg:
		luaMsg = msgTableSimple(state, MSG_QUIT)
	case tea.ResumeMsg:
		luaMsg = msgTableSimple(state, MSG_RESUME)
	case tea.SuspendMsg:
		luaMsg = msgTableSimple(state, MSG_SUSPEND)
	case tea.WindowSizeMsg:
		luaMsg = msgTableWindowSize(state, msg.Width, msg.Height)
	case *golua.LTable:
		luaMsg = msg
	case golua.LValue:
		luaMsg = msgTableLua(state, msg)
	default:
		luaMsg = msgTableSimple(state, MSG_NONE)
	}

	return luaMsg
}

func msgTableSimple(state *golua.LState, msg TeaMSG) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(msg))

	return t
}

func msgTableKey(state *golua.LState, key tea.Key) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_KEY))
	t.RawSetString("key", golua.LString(key.String()))
	t.RawSetString("event", KeyEventTable(state, key))

	return t
}

func msgTableMouse(state *golua.LState, mouse string, event tea.MouseEvent) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_MOUSE))
	t.RawSetString("key", golua.LString(mouse))
	t.RawSetString("event", MouseEventTable(state, event))

	return t
}

func msgTableWindowSize(state *golua.LState, width, height int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_WINDOWSIZE))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))

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

func msgTableLua(state *golua.LState, value golua.LValue) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("msg", golua.LNumber(MSG_LUA))
	t.RawSetString("value", value)

	return t
}
