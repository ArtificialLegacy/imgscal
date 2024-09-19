package lib

import (
	"errors"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TUI = "tui"

/// @lib Terminal UI
/// @import tui
/// @desc
/// Library for creating BubbleTea TUIs.

func RegisterTUI(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TUI, r, r.State, lg)

	/// @func new() -> struct<tui.Program>
	/// @returns {struct<tui.Program>}
	lib.CreateFunction(tab, "new",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := r.CR_TEA.Add(&collection.TeaItem{
				Spinners: map[int]*spinner.Model{},
			})
			t := teaTable(r, state, lib, id)

			state.Push(t)
			return 1
		},
	)

	/// @func run(program)
	/// @arg program {struct<tui.Program>}
	lib.CreateFunction(tab, "run",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "program"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			program := args["program"].(*golua.LTable)
			id := int(program.RawGetString("id").(golua.LNumber))
			item, err := r.CR_TEA.Item(id)
			if err != nil {
				lua.Error(state, err.Error())
			}

			pstate, _ := state.NewThread()
			p := tea.NewProgram(teamodel{id: id, item: item, state: pstate, r: r, lg: lg})
			_, err = p.Run()
			if err != nil {
				lua.Error(state, err.Error())
			}

			pstate.Close()

			return 0
		},
	)

	/// @func spinner(id, from?) -> struct<tui.Spinner>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the spinner to.
	/// @arg? from {int<tui.Spinner>} - The built-in spinner to use.
	/// @returns {struct<tui.Spinner>}
	lib.CreateFunction(tab, "spinner",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "type", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			spin := spinner.New(spinner.WithSpinner(spinnerList[args["type"].(int)]))
			id := spin.ID()
			item.Spinners[id] = &spin

			t := spinnerTable(r, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func spinner_custom(id, frames, fps) -> struct<tui.Spinner>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the spinner to.
	/// @arg frames {[]string}
	/// @arg fps {int}
	/// @returns {struct<tui.Spinner>}
	lib.CreateFunction(tab, "spinner_custom",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("frames", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.INT, Name: "fps"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			frames := args["frames"].([]any)
			fps := args["fps"].(int)

			frameBuild := make([]string, len(frames))
			for i, f := range frames {
				frameBuild[i] = f.(string)
			}

			spin := spinner.New(spinner.WithSpinner(spinner.Spinner{
				Frames: frameBuild,
				FPS:    time.Second / time.Duration(fps),
			}))
			id := spin.ID()
			item.Spinners[id] = &spin

			t := spinnerTable(r, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func textarea(id) -> struct<tui.TextArea>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the text area to.
	/// @returns {struct<tui.TextArea>}
	lib.CreateFunction(tab, "textarea",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			ta := textarea.New()
			id := len(item.TextAreas)
			item.TextAreas = append(item.TextAreas, &ta)

			t := textareaTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func textinput(id) -> struct<tui.TextInput>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the text input to.
	/// @returns {struct<tui.TextInput>}
	lib.CreateFunction(tab, "textinput",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			ti := textinput.New()
			id := len(item.TextInputs)
			item.TextInputs = append(item.TextInputs, &ti)

			t := textinputTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func cursor(id) -> struct<tui.Cursor>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the cursor to.
	/// @returns {struct<tui.Cursor>}
	lib.CreateFunction(tab, "cursor",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			cu := cursor.New()
			id := len(item.Cursors)
			item.Cursors = append(item.Cursors, &cu)

			t := cursorTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func filepicker(id) -> struct<tui.FilePicker>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the file picker to.
	/// @returns {struct<tui.FilePicker>}
	lib.CreateFunction(tab, "filepicker",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			fp := filepicker.New()
			id := len(item.FilePickers)
			item.FilePickers = append(item.FilePickers, &fp)

			t := filePickerTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func cmd_batch(cmds) -> struct<tui.CMDBatch>
	/// @arg cmds {[]struct<tui.CMD>}
	/// @returns {struct<tui.CMDBatch>}
	lib.CreateFunction(tab, "cmd_batch",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "cmds"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cmdBatch(state, args["cmds"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func cmd_sequence(cmds) -> struct<tui.CMDSequence>
	/// @arg cmds {[]struct<tui.CMD>}
	/// @returns {struct<tui.CMDSequence>}
	lib.CreateFunction(tab, "cmd_sequence",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "cmds"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cmdSequence(state, args["cmds"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @constants Message
	/// @const MSG_KEY
	/// @const MSG_SPINNERTICK
	/// @const MSG_BLINK
	tab.RawSetString("MSG_KEY", golua.LNumber(MSG_KEY))
	tab.RawSetString("MSG_SPINNERTICK", golua.LNumber(MSG_SPINNERTICK))
	tab.RawSetString("MSG_BLINK", golua.LNumber(MSG_BLINK))

	/// @constants Spinners
	/// @const SPINNER_LINE
	/// @const SPINNER_DOT
	/// @const SPINNER_MINIDOT
	/// @const SPINNER_JUMP
	/// @const SPINNER_PULSE
	/// @const SPINNER_POINTS
	/// @const SPINNER_GLOBE
	/// @const SPINNER_MOON
	/// @const SPINNER_MONKEY
	/// @const SPINNER_METER
	/// @const SPINNER_HAMBURGER
	/// @const SPINNER_ELLIPSIS
	tab.RawSetString("SPINNER_LINE", golua.LNumber(SPINNER_LINE))
	tab.RawSetString("SPINNER_DOT", golua.LNumber(SPINNER_DOT))
	tab.RawSetString("SPINNER_MINIDOT", golua.LNumber(SPINNER_MINIDOT))
	tab.RawSetString("SPINNER_JUMP", golua.LNumber(SPINNER_JUMP))
	tab.RawSetString("SPINNER_PULSE", golua.LNumber(SPINNER_PULSE))
	tab.RawSetString("SPINNER_POINTS", golua.LNumber(SPINNER_POINTS))
	tab.RawSetString("SPINNER_GLOBE", golua.LNumber(SPINNER_GLOBE))
	tab.RawSetString("SPINNER_MOON", golua.LNumber(SPINNER_MOON))
	tab.RawSetString("SPINNER_MONKEY", golua.LNumber(SPINNER_MONKEY))
	tab.RawSetString("SPINNER_METER", golua.LNumber(SPINNER_METER))
	tab.RawSetString("SPINNER_HAMBURGER", golua.LNumber(SPINNER_HAMBURGER))
	tab.RawSetString("SPINNER_ELLIPSIS", golua.LNumber(SPINNER_ELLIPSIS))

	/// @constants Text Input Echo Mode
	/// @const ECHO_NORMAL
	/// @const ECHO_PASSWORD
	/// @const ECHO_NONE
	tab.RawSetString("ECHO_NORMAL", golua.LNumber(textinput.EchoNormal))
	tab.RawSetString("ECHO_PASSWORD", golua.LNumber(textinput.EchoPassword))
	tab.RawSetString("ECHO_NONE", golua.LNumber(textinput.EchoNone))

	/// @constants Cursor Mode
	/// @const CURSOR_BLINK
	/// @const CURSOR_STATIC
	/// @const CURSOR_HIDE
	tab.RawSetString("CURSOR_BLINK", golua.LNumber(cursor.CursorBlink))
	tab.RawSetString("CURSOR_STATIC", golua.LNumber(cursor.CursorStatic))
	tab.RawSetString("CURSOR_HIDE", golua.LNumber(cursor.CursorHide))
}

type TeaMSG int

const (
	MSG_NONE TeaMSG = iota
	MSG_KEY
	MSG_SPINNERTICK
	MSG_BLINK
)

type TeaCMD int

const (
	CMD_NONE TeaCMD = iota
	CMD_STORED
	CMD_BATCH
	CMD_SEQUENCE
	CMD_SPINNERTICK
	CMD_TEXTAREAFOCUS
	CMD_TEXTINPUTFOCUS
	CMD_BLINK
	CMD_CURSORFOCUS
	CMD_FILEPICKERINIT
)

var cmdList = []func(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd{}

func init() {
	cmdList = []func(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd{
		func(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd { return nil },
		cmdStoredBuild,
		cmdBatchBuild,
		cmdSequenceBuild,
		cmdSpinnerTickBuild,
		cmdTextAreaFocusBuild,
		cmdTextInputFocusBuild,
		cmdBlinkBuild,
		cmdCursorFocusBuild,
		cmdFilePickerInitBuild,
	}
}

type Spinners int

const (
	SPINNER_LINE Spinners = iota
	SPINNER_DOT
	SPINNER_MINIDOT
	SPINNER_JUMP
	SPINNER_PULSE
	SPINNER_POINTS
	SPINNER_GLOBE
	SPINNER_MOON
	SPINNER_MONKEY
	SPINNER_METER
	SPINNER_HAMBURGER
	SPINNER_ELLIPSIS
)

var spinnerList = []spinner.Spinner{
	spinner.Line,
	spinner.Dot,
	spinner.MiniDot,
	spinner.Jump,
	spinner.Pulse,
	spinner.Points,
	spinner.Globe,
	spinner.Moon,
	spinner.Monkey,
	spinner.Meter,
	spinner.Hamburger,
	spinner.Ellipsis,
}

type teamodel struct {
	item  *collection.TeaItem
	state *golua.LState
	id    int
	r     *lua.Runner
	lg    *log.Logger
}

func (m teamodel) Init() tea.Cmd {
	m.state.Push(m.item.FnInit)
	m.state.Push(golua.LNumber(m.id))
	m.state.Call(1, 2)
	model := m.state.CheckTable(-2)
	cmd := m.state.CheckTable(-1)
	m.state.Pop(2)

	bcmd := m.cmdBuild(m.state, cmd)

	m.item.LuaModel = model

	return bcmd
}

func (m teamodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.item.Msg = &msg
	defer func() {
		m.item.Msg = nil
	}()
	luaMsg := msgTableNone(m.state)
	var bcmd tea.Cmd

	m.item.Cmds = []tea.Cmd{}

	// Program should always be exittable
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		luaMsg = msgTableKey(m.state, msg.String())
	case spinner.TickMsg:
		luaMsg = msgTableSpinnerTick(m.state, msg.ID)
	case cursor.BlinkMsg:
		luaMsg = msgTableCursorBlink(m.state)
	}

	if luaMsg != nil {
		m.state.Push(m.item.FnUpdate)
		m.state.Push(m.item.LuaModel)
		m.state.Push(luaMsg)

		m.state.Call(2, 1)
		cmd := m.state.OptTable(-1, cmdNone(m.state))
		m.state.Pop(1)

		bcmd = m.cmdBuild(m.state, cmd)
	}

	return m, bcmd
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

func (m *teamodel) cmdBuild(state *golua.LState, t *golua.LTable) tea.Cmd {
	cmdValue := t.RawGetString("cmd")
	cmdId := 0
	if cmdValue.Type() == golua.LTNumber {
		cmdId = int(cmdValue.(golua.LNumber))
	}
	cmd := cmdList[cmdId](m.r, m, state, t)

	return cmd
}

func cmdNone(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_NONE))

	return t
}

func cmdStored(state *golua.LState, item *collection.TeaItem, cmd tea.Cmd) *golua.LTable {
	t := state.NewTable()

	id := len(item.Cmds)
	item.Cmds = append(item.Cmds, cmd)

	t.RawSetString("cmd", golua.LNumber(CMD_STORED))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdStoredBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmd := m.item.Cmds[int(t.RawGetString("id").(golua.LNumber))]

	return cmd
}

func cmdBatch(state *golua.LState, cmds *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_BATCH))
	t.RawSetString("cmds", cmds)

	return t
}

func cmdBatchBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmdList := t.RawGetString("cmds").(*golua.LTable)
	cmds := make([]tea.Cmd, cmdList.Len())

	for i := range cmdList.Len() {
		bcmd := m.cmdBuild(state, cmdList.RawGetInt(i+1).(*golua.LTable))
		cmds[i] = bcmd
	}

	return tea.Batch(cmds...)
}

func cmdSequence(state *golua.LState, cmds *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SEQUENCE))
	t.RawSetString("cmds", cmds)

	return t
}

func cmdSequenceBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmdList := t.RawGetString("cmds").(*golua.LTable)
	cmds := make([]tea.Cmd, cmdList.Len())

	for i := range cmdList.Len() {
		bcmd := m.cmdBuild(state, cmdList.RawGetInt(i+1).(*golua.LTable))
		cmds[i] = bcmd
	}

	return tea.Sequence(cmds...)
}

func cmdSpinnerTick(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SPINNERTICK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdTextAreaFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TEXTAREAFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdTextInputFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TEXTINPUTFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdTextAreaFocusBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	textArea := m.item.TextAreas[int(t.RawGetString("id").(golua.LNumber))]

	return textArea.Focus()
}

func cmdTextInputFocusBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	textInput := m.item.TextInputs[int(t.RawGetString("id").(golua.LNumber))]

	return textInput.Focus()
}

func cmdSpinnerTickBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	id := int(t.RawGetString("id").(golua.LNumber))
	return m.item.Spinners[id].Tick
}

func cmdBlink(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_BLINK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdBlinkBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	cursor := m.item.Cursors[int(t.RawGetString("id").(golua.LNumber))]

	return cursor.BlinkCmd()
}

func cmdCursorFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_CURSORFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdCursorFocusBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	cursor := m.item.Cursors[int(t.RawGetString("id").(golua.LNumber))]

	return cursor.Focus()
}

func cmdFilePickerInit(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_FILEPICKERINIT))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdFilePickerInitBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	fp := m.item.FilePickers[int(t.RawGetString("id").(golua.LNumber))]

	return fp.Init()
}

func (m teamodel) View() string {
	m.state.Push(m.item.FnView)
	m.state.Push(m.item.LuaModel)
	m.state.Call(1, 1)
	str := m.state.CheckString(-1)
	m.state.Pop(1)

	return str
}

func teaTable(r *lua.Runner, state *golua.LState, lib *lua.Lib, id int) *golua.LTable {
	t := state.NewTable()
	t.RawSetString("id", golua.LNumber(id))

	lib.BuilderFunction(state, t, "init",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("id").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}

			item.FnInit = args["fn"].(*golua.LFunction)
		},
	)

	lib.BuilderFunction(state, t, "update",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("id").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}

			item.FnUpdate = args["fn"].(*golua.LFunction)
		},
	)

	lib.BuilderFunction(state, t, "view",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("id").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}

			item.FnView = args["fn"].(*golua.LFunction)
		},
	)

	return t
}

func spinnerTable(r *lua.Runner, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Spinners[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		nm, cmd := item.Spinners[id].Update(*item.Msg)
		item.Spinners[id] = &nm

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("tick", state.NewFunction(func(state *golua.LState) int {
		cmd := cmdSpinnerTick(state, int(t.RawGetString("id").(golua.LNumber)))

		state.Push(cmd)
		return 1
	}))

	return t
}

func textareaTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.TextAreas[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		nm, cmd := item.TextAreas[id].Update(*item.Msg)
		item.TextAreas[id] = &nm

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.BuilderFunction(state, t, "reset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].Reset()
		})

	t.RawSetString("focus", state.NewFunction(func(state *golua.LState) int {
		t := cmdTextAreaFocus(state, int(t.RawGetString("id").(golua.LNumber)))

		state.Push(t)
		return 1
	}))

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].Blur()
		})

	lib.BuilderFunction(state, t, "cursor_down",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorDown()
		})

	lib.BuilderFunction(state, t, "cursor_end",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorEnd()
		})

	lib.BuilderFunction(state, t, "cursor_start",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorStart()
		})

	lib.BuilderFunction(state, t, "cursor_up",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorUp()
		})

	t.RawSetString("focused", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		focused := item.TextAreas[id].Focused()

		state.Push(golua.LBool(focused))
		return 1
	}))

	t.RawSetString("size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.TextAreas[id].Width()
		height := item.TextAreas[id].Height()

		state.Push(golua.LNumber(width))
		state.Push(golua.LNumber(height))
		return 2
	}))

	t.RawSetString("width", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.TextAreas[id].Width()

		state.Push(golua.LNumber(width))
		return 1
	}))

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		height := item.TextAreas[id].Height()

		state.Push(golua.LNumber(height))
		return 1
	}))

	lib.BuilderFunction(state, t, "size_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetWidth(args["width"].(int))
			item.TextAreas[id].SetHeight(args["height"].(int))
		})

	lib.BuilderFunction(state, t, "width_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetWidth(args["width"].(int))
		})

	lib.BuilderFunction(state, t, "height_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetHeight(args["height"].(int))
		})

	lib.BuilderFunction(state, t, "insert_rune",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].InsertRune(rune(args["rune"].(int)))
		})

	lib.BuilderFunction(state, t, "insert_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].InsertString(args["str"].(string))
		})

	t.RawSetString("length", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		length := item.TextAreas[id].Length()

		state.Push(golua.LNumber(length))
		return 1
	}))

	t.RawSetString("line", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		line := item.TextAreas[id].Line()

		state.Push(golua.LNumber(line))
		return 1
	}))

	t.RawSetString("line_count", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		count := item.TextAreas[id].LineCount()

		state.Push(golua.LNumber(count))
		return 1
	}))

	lib.BuilderFunction(state, t, "cursor_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetCursor(args["col"].(int))
		})

	t.RawSetString("value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextAreas[id].Value()

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "value_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetValue(args["str"].(string))
		})

	t.RawSetString("line_info", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		info := item.TextAreas[id].LineInfo()

		state.Push(lineInfoTable(state, &info))
		return 1
	}))

	t.RawSetString("prompt", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		prompt := item.TextAreas[id].Prompt

		state.Push(golua.LString(prompt))
		return 1
	}))

	lib.BuilderFunction(state, t, "prompt_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.Prompt = args["str"].(string)
			ta.SetWidth(ta.Width())
		})

	t.RawSetString("line_numbers", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		linenum := item.TextAreas[id].ShowLineNumbers

		state.Push(golua.LBool(linenum))
		return 1
	}))

	lib.BuilderFunction(state, t, "line_numbers_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enable"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.ShowLineNumbers = args["enable"].(bool)
			ta.SetWidth(ta.Width())
		})

	t.RawSetString("char_end", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		char := item.TextAreas[id].EndOfBufferCharacter

		state.Push(golua.LNumber(char))
		return 1
	}))

	lib.BuilderFunction(state, t, "char_end_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.EndOfBufferCharacter = rune(args["rune"].(int))
		})

	t.RawSetString("char_limit", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		limit := item.TextAreas[id].CharLimit

		state.Push(golua.LNumber(limit))
		return 1
	}))

	lib.BuilderFunction(state, t, "char_limit_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "limit"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.CharLimit = args["limit"].(int)
		})

	t.RawSetString("width_max", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.TextAreas[id].MaxWidth

		state.Push(golua.LNumber(width))
		return 1
	}))

	lib.BuilderFunction(state, t, "width_max_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.MaxWidth = args["width"].(int)
		})

	t.RawSetString("height_max", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		height := item.TextAreas[id].MaxHeight

		state.Push(golua.LNumber(height))
		return 1
	}))

	lib.BuilderFunction(state, t, "height_max_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.MaxHeight = args["height"].(int)
		})

	lib.BuilderFunction(state, t, "prompt_func",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			promptState, _ := state.NewThread()
			ta.SetPromptFunc(args["width"].(int), func(lineIdx int) string {
				promptState.Push(args["fn"].(*golua.LFunction))
				promptState.Push(golua.LNumber(lineIdx))
				promptState.Call(1, 1)
				str := promptState.CheckString(-1)
				promptState.Pop(1)

				return str
			})
		})

	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextAreas[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)

		state.Push(golua.LNumber(cid))
		return 1
	}))

	return t
}

func lineInfoTable(state *golua.LState, info *textarea.LineInfo) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("width", golua.LNumber(info.Width))
	t.RawSetString("width_char", golua.LNumber(info.CharWidth))
	t.RawSetString("height", golua.LNumber(info.Height))
	t.RawSetString("column_start", golua.LNumber(info.StartColumn))
	t.RawSetString("column_offset", golua.LNumber(info.ColumnOffset))
	t.RawSetString("row_offset", golua.LNumber(info.RowOffset))
	t.RawSetString("char_offset", golua.LNumber(info.CharOffset))

	return t
}

func textinputTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.TextInputs[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		nm, cmd := item.TextInputs[id].Update(*item.Msg)
		item.TextInputs[id] = &nm

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("focus", state.NewFunction(func(state *golua.LState) int {
		t := cmdTextInputFocus(state, int(t.RawGetString("id").(golua.LNumber)))

		state.Push(t)
		return 1
	}))

	lib.BuilderFunction(state, t, "reset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Reset()
		})

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Blur()
		})

	lib.BuilderFunction(state, t, "cursor_start",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CursorStart()
		})

	lib.BuilderFunction(state, t, "cursor_end",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CursorEnd()
		})

	t.RawSetString("current_suggestion", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		suggestion := item.TextInputs[id].CurrentSuggestion()

		state.Push(golua.LString(suggestion))
		return 1
	}))

	t.RawSetString("available_suggestions", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		suggestions := item.TextInputs[id].AvailableSuggestions()

		slist := state.NewTable()
		for i, s := range suggestions {
			slist.RawSetInt(i+1, golua.LString(s))
		}

		state.Push(slist)
		return 1
	}))

	lib.BuilderFunction(state, t, "suggestions_set",
		[]lua.Arg{
			lua.ArgArray("suggestions", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			suggestions := args["suggestions"].([]any)
			slist := make([]string, len(suggestions))
			for i, s := range suggestions {
				slist[i] = s.(string)
			}

			item.TextInputs[id].SetSuggestions(slist)
		})

	t.RawSetString("focused", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		focused := item.TextInputs[id].Focused()

		state.Push(golua.LBool(focused))
		return 1
	}))

	t.RawSetString("position", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		pos := item.TextInputs[id].Position()

		state.Push(golua.LNumber(pos))
		return 1
	}))

	lib.BuilderFunction(state, t, "position_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "pos"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].SetCursor(args["pos"].(int))
		})

	t.RawSetString("value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].Value()

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "value_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].SetValue(args["value"].(string))
		})

	lib.BuilderFunction(state, t, "validate",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			newState, _ := state.NewThread()
			item.TextInputs[id].Validate = func(str string) error {
				newState.Push(args["fn"].(*golua.LFunction))
				newState.Push(golua.LString(str))
				newState.Call(1, 2)

				ok := newState.CheckBool(-2)
				err := newState.CheckString(-1)
				newState.Pop(2)

				if !ok {
					return errors.New(err)
				}
				return nil
			}
		})

	t.RawSetString("prompt", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].Prompt

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "prompt_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Prompt = args["value"].(string)
		})

	t.RawSetString("placeholder", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].Placeholder

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "placeholder_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Placeholder = args["value"].(string)
		})

	t.RawSetString("echomode", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].EchoMode

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "echomode_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "echomode"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].EchoMode = textinput.EchoMode(args["echomode"].(int))
		})

	t.RawSetString("echo_char", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].EchoCharacter

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "echo_char_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].EchoCharacter = rune(args["rune"].(int))
		})

	t.RawSetString("char_limit", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].CharLimit

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "char_limit_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "limit"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CharLimit = args["limit"].(int)
		})

	t.RawSetString("width", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].Width

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "width_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "limit"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Width = args["width"].(int)
		})

	t.RawSetString("suggestions_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.TextInputs[id].ShowSuggestions

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "suggestions_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "show"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].ShowSuggestions = args["show"].(bool)
		})

	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextInputs[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)

		state.Push(golua.LNumber(cid))
		return 1
	}))

	return t
}

func cursorTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Cursors[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		nm, cmd := item.Cursors[id].Update(*item.Msg)
		item.Cursors[id] = &nm

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("blink", state.NewFunction(func(state *golua.LState) int {
		state.Push(cmdBlink(state, int(t.RawGetString("id").(golua.LNumber))))
		return 1
	}))

	t.RawSetString("focus", state.NewFunction(func(state *golua.LState) int {
		t := cmdCursorFocus(state, int(t.RawGetString("id").(golua.LNumber)))

		state.Push(t)
		return 1
	}))

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Cursors[id].Blur()
		})

	t.RawSetString("mode", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Cursors[id].Mode()

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "mode_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "mode"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Cursors[id].SetMode(cursor.Mode(args["mode"].(int)))
		})

	lib.BuilderFunction(state, t, "char_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Cursors[id].SetChar(args["str"].(string))
		})

	return t
}

func filePickerTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.FilePickers[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		nm, cmd := item.FilePickers[id].Update(*item.Msg)
		item.FilePickers[id] = &nm

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("did_select_file", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		did, str := item.FilePickers[id].DidSelectFile(*item.Msg)

		state.Push(golua.LBool(did))
		state.Push(golua.LString(str))
		return 2
	}))

	t.RawSetString("did_select_disabled", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		did, str := item.FilePickers[id].DidSelectDisabledFile(*item.Msg)

		state.Push(golua.LBool(did))
		state.Push(golua.LString(str))
		return 2
	}))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.FilePickers[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("init", state.NewFunction(func(state *golua.LState) int {
		state.Push(cmdFilePickerInit(state, int(t.RawGetString("id").(golua.LNumber))))
		return 1
	}))

	t.RawSetString("path", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].Path

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "path_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].Path = args["path"].(string)
		})

	t.RawSetString("current_directory", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].CurrentDirectory

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "current_directory_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "dir"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].CurrentDirectory = args["dir"].(string)
		})

	t.RawSetString("allowed_types", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		allowed := item.FilePickers[id].AllowedTypes
		list := state.NewTable()

		for i, s := range allowed {
			list.RawSetInt(i+1, golua.LString(s))
		}

		state.Push(list)
		return 1
	}))

	lib.BuilderFunction(state, t, "allowed_types_set",
		[]lua.Arg{
			lua.ArgArray("allowed", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			allowed := args["allowed"].([]any)
			list := make([]string, len(allowed))
			for i, v := range allowed {
				list[i] = v.(string)
			}
			item.FilePickers[id].AllowedTypes = list
		})

	t.RawSetString("show_perm", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].ShowPermissions

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "show_perm_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowPermissions = args["enabled"].(bool)
		})

	t.RawSetString("show_size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].ShowSize

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "show_size_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowSize = args["enabled"].(bool)
		})

	t.RawSetString("show_hidden", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].ShowHidden

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "show_hidden_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowHidden = args["enabled"].(bool)
		})

	t.RawSetString("dir_allowed", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].DirAllowed

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "dir_allowed_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].DirAllowed = args["enabled"].(bool)
		})

	t.RawSetString("file_allowed", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].FileAllowed

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "file_allowed_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].FileAllowed = args["enabled"].(bool)
		})

	t.RawSetString("file_selected", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].FileSelected

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "file_selected_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "file"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].FileSelected = args["file"].(string)
		})

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].Height

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "height_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].Height = args["height"].(int)
		})

	t.RawSetString("height_auto", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].AutoHeight

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "height_auto_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].AutoHeight = args["enabled"].(bool)
		})

	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.FilePickers[id].Cursor

		state.Push(golua.LString(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "cursor_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "cursor"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].Cursor = args["cursor"].(string)
		})

	return t
}
