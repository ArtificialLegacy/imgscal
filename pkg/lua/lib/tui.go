package lib

import (
	"errors"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
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

	/// @func list_item(title, desc, filter?) -> struct<tui.ListItem>
	/// @arg title {string}
	/// @arg desc {string}
	/// @arg? filter {string} - Defaults to the value of title.
	/// @returns {struct<tui.ListItem>}
	lib.CreateFunction(tab, "list_item",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.STRING, Name: "desc"},
			{Type: lua.STRING, Name: "filter", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			li := listItemTable(state, args["title"].(string), args["desc"].(string), args["filter"].(string))

			state.Push(li)
			return 1
		})

	/// @func list_filter_state_string(state) -> string
	/// @arg state {int<tui.FilterState>}
	/// @returns {string}
	lib.CreateFunction(tab, "list_filter_state_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "state"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := list.FilterState(args["state"].(int))
			state.Push(golua.LString(s.String()))
			return 1
		})

	/// @func list(id, items, width, height) -> struct<tui.List>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the list to.
	/// @arg items {[]struct<tui.ListItem>} - Array of list items.
	/// @arg width {int}
	/// @arg height {int}
	/// @returns {struct<tui.List>}
	lib.CreateFunction(tab, "list",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("items", lua.ArrayType{Type: lua.RAW_TABLE}, false),
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			itemList := args["items"].([]any)
			items := make([]list.Item, len(itemList))

			for i, v := range itemList {
				items[i] = listItemBuild(v.(*golua.LTable))
			}

			li := list.New(items, list.NewDefaultDelegate(), args["width"].(int), args["height"].(int))
			id := len(item.Lists)
			item.Lists = append(item.Lists, &li)

			t := listTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func paginator(id, per?, total?) -> struct<tui.Paginator>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the paginator to.
	/// @arg? per {int} - Not set if left at default value of 0.
	/// @arg? total {int} - Not set if left at default value of 0.
	/// @returns {struct<tui.Paginator>}
	lib.CreateFunction(tab, "paginator",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "per", Optional: true},
			{Type: lua.INT, Name: "total", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			per := args["per"].(int)
			total := args["total"].(int)

			opts := []paginator.Option{}
			if per > 0 {
				opts = append(opts, paginator.WithPerPage(per))
			}
			if total > 0 {
				opts = append(opts, paginator.WithTotalPages(total))
			}

			pg := paginator.New(opts...)
			id := len(item.Paginators)
			item.Paginators = append(item.Paginators, &pg)

			t := paginatorTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func progress_options() -> struct<tui.ProgressOptions>
	/// @returns {struct<tui.ProgressOptions>}
	lib.CreateFunction(tab, "progress_options",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(progressOptionsTable(lib, state))
			return 1
		})

	/// @func progress(id, options?) -> struct<tui.Progress>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the progress bar to.
	/// @arg? options {struct<tui.ProgressOptions>}
	/// @returns {struct<tui.Progress>}
	lib.CreateFunction(tab, "progress",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			opts := progressOptionsBuild(args["options"].(*golua.LTable))

			pr := progress.New(opts...)
			id := len(item.ProgressBars)
			item.ProgressBars = append(item.ProgressBars, &pr)

			t := progressTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func file_is_hidden(path) -> bool
	/// @arg path {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "file_is_hidden",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			is, _ := filepicker.IsHidden(args["path"].(string))

			state.Push(golua.LBool(is))
			return 1
		})

	/// @func cmd_none() -> struct<tui.CMDNone>
	/// @returns {struct<tui.CMDNone>}
	lib.CreateFunction(tab, "cmd_none",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(cmdNone(state))
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
	/// @const MSG_STOPWATCHRESET
	/// @const MSG_STOPWATCHSTARTSTOP
	/// @const MSG_STOPWATCHTICK
	tab.RawSetString("MSG_KEY", golua.LNumber(MSG_KEY))
	tab.RawSetString("MSG_SPINNERTICK", golua.LNumber(MSG_SPINNERTICK))
	tab.RawSetString("MSG_BLINK", golua.LNumber(MSG_BLINK))
	tab.RawSetString("MSG_STOPWATCHRESET", golua.LNumber(MSG_STOPWATCHRESET))
	tab.RawSetString("MSG_STOPWATCHSTARTSTOP", golua.LNumber(MSG_STOPWATCHSTARTSTOP))
	tab.RawSetString("MSG_STOPWATCHTICK", golua.LNumber(MSG_STOPWATCHTICK))

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

	/// @constants List Filter State
	/// @const FILTERSTATE_UNFILTERED
	/// @const FILTERSTATE_FILTERING
	/// @const FILTERSTATE_APPLIED
	tab.RawSetString("FILTERSTATE_UNFILTERED", golua.LNumber(list.Unfiltered))
	tab.RawSetString("FILTERSTATE_FILTERING", golua.LNumber(list.Filtering))
	tab.RawSetString("FILTERSTATE_APPLIED", golua.LNumber(list.FilterApplied))

	/// @constants Paginator Types
	/// @const PAGINATOR_ARABIC
	/// @const PAGINATOR_DOT
	tab.RawSetString("PAGINATOR_ARABIC", golua.LNumber(paginator.Arabic))
	tab.RawSetString("PAGINATOR_DOT", golua.LNumber(paginator.Dots))
}

type TeaMSG int

const (
	MSG_NONE TeaMSG = iota
	MSG_KEY
	MSG_SPINNERTICK
	MSG_BLINK
	MSG_STOPWATCHRESET
	MSG_STOPWATCHSTARTSTOP
	MSG_STOPWATCHTICK
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
	CMD_LISTSETITEMS
	CMD_LISTINSERTITEM
	CMD_LISTSETITEM
	CMD_LISTSTATUSMESSAGE
	CMD_LISTSPINNERSTART
	CMD_LISTSPINNERTOGGLE
	CMD_PROGRESSSET
	CMD_PROGRESSDEC
	CMD_PROGRESSINC
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
		cmdListSetItemsBuild,
		cmdListInsertItemBuild,
		cmdListSetItemBuild,
		cmdListStatusMessageBuild,
		cmdListSpinnerStartBuild,
		cmdListSpinnerToggleBuild,
		cmdProgressSetBuild,
		cmdProgressDecBuild,
		cmdProgressIncBuild,
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

func cmdListSetItems(state *golua.LState, id int, items *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSETITEMS))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("items", items)

	return t
}

func cmdListSetItemsBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	itemList := t.RawGetString("items").(*golua.LTable)
	items := make([]list.Item, itemList.Len())

	for i := range itemList.Len() {
		items[i] = listItemBuild(itemList.RawGetInt(i + 1).(*golua.LTable))
	}

	return li.SetItems(items)
}

func cmdListInsertItem(state *golua.LState, id, index int, item *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTINSERTITEM))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("index", golua.LNumber(index))
	t.RawSetString("item", item)

	return t
}

func cmdListInsertItemBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	it := t.RawGetString("item").(*golua.LTable)
	index := int(t.RawGetString("index").(golua.LNumber))

	return li.InsertItem(index, listItemBuild(it))
}

func cmdListSetItem(state *golua.LState, id, index int, item *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSETITEM))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("index", golua.LNumber(index))
	t.RawSetString("item", item)

	return t
}

func cmdListSetItemBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	it := t.RawGetString("item").(*golua.LTable)
	index := int(t.RawGetString("index").(golua.LNumber))

	return li.SetItem(index, listItemBuild(it))
}

func cmdListStatusMessage(state *golua.LState, id int, msg string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSTATUSMESSAGE))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("msg", golua.LString(msg))

	return t
}

func cmdListStatusMessageBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	msg := string(t.RawGetString("msg").(golua.LString))
	return li.NewStatusMessage(msg)
}

func cmdListSpinnerStart(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSPINNERSTART))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdListSpinnerStartBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	return li.StartSpinner()
}

func cmdListSpinnerToggle(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSPINNERTOGGLE))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdListSpinnerToggleBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	return li.ToggleSpinner()
}

func cmdProgressSet(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSSET))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func cmdProgressSetBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.SetPercent(percent)
}

func cmdProgressDec(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSDEC))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdProgressDecBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.DecrPercent(percent)
}

func cmdProgressInc(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSINC))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func cmdProgressIncBuild(r *lua.Runner, m *teamodel, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := m.item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.IncrPercent(percent)
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

type teaListItem struct {
	title       string
	description string
	filter      string
}

func (i teaListItem) Title() string       { return i.title }
func (i teaListItem) Description() string { return i.description }
func (i teaListItem) FilterValue() string { return i.filter }

func listItemTable(state *golua.LState, title, description, filter string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("title", golua.LString(title))
	t.RawSetString("description", golua.LString(description))
	t.RawSetString("filter", golua.LString(filter))

	return t
}

func listItemBuild(t *golua.LTable) teaListItem {
	return teaListItem{
		title:       string(t.RawGetString("title").(golua.LString)),
		description: string(t.RawGetString("description").(golua.LString)),
		filter:      string(t.RawGetString("filter").(golua.LString)),
	}
}

func listTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Lists[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		li, cmd := item.Lists[id].Update(*item.Msg)
		item.Lists[id] = &li

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].Cursor()

		state.Push(golua.LNumber(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "cursor_up",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].CursorUp()
		})

	lib.BuilderFunction(state, t, "cursor_down",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].CursorDown()
		})

	lib.BuilderFunction(state, t, "page_next",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].NextPage()
		})

	lib.BuilderFunction(state, t, "page_prev",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].PrevPage()
		})

	t.RawSetString("pagination_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].ShowPagination()

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "pagination_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowPagination(args["enabled"].(bool))
		})

	lib.BuilderFunction(state, t, "disable_quit",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].DisableQuitKeybindings()
		})

	t.RawSetString("size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.Lists[id].Width()
		height := item.Lists[id].Height()

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

		width := item.Lists[id].Width()

		state.Push(golua.LNumber(width))
		return 1
	}))

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		height := item.Lists[id].Height()

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

			item.Lists[id].SetWidth(args["width"].(int))
			item.Lists[id].SetHeight(args["height"].(int))
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

			item.Lists[id].SetWidth(args["width"].(int))
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

			item.Lists[id].SetHeight(args["height"].(int))
		})

	t.RawSetString("filter_state", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].FilterState()

		state.Push(golua.LNumber(value))
		return 1
	}))

	t.RawSetString("filter_value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].FilterValue()

		state.Push(golua.LString(value))
		return 1
	}))

	t.RawSetString("filter_enabled", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].FilteringEnabled()

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "filter_enabled_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetFilteringEnabled(args["enabled"].(bool))
		})

	t.RawSetString("filter_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].ShowFilter()

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "filter_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowFilter(args["enabled"].(bool))
		})

	lib.BuilderFunction(state, t, "filter_reset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].ResetFilter()
		})

	t.RawSetString("is_filtered", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].IsFiltered()

		state.Push(golua.LBool(value))
		return 1
	}))

	t.RawSetString("filter_setting", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].SettingFilter()

		state.Push(golua.LBool(value))
		return 1
	}))

	t.RawSetString("index", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].Index()

		state.Push(golua.LNumber(value))
		return 1
	}))

	t.RawSetString("items", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		itemList := item.Lists[id].Items()
		items := state.NewTable()

		for i, v := range itemList {
			li := v.(teaListItem)
			items.RawSetInt(i+1, listItemTable(state, li.title, li.description, li.filter))
		}

		state.Push(items)
		return 1
	}))

	t.RawSetString("items_visible", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		itemList := item.Lists[id].VisibleItems()
		items := state.NewTable()

		for i, v := range itemList {
			li := v.(teaListItem)
			items.RawSetInt(i+1, listItemTable(state, li.title, li.description, li.filter))
		}

		state.Push(items)
		return 1
	}))

	lib.TableFunction(state, t, "items_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "items"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListSetItems(state, id, args["items"].(*golua.LTable)))
			return 1
		})

	lib.TableFunction(state, t, "item_insert",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
			{Type: lua.RAW_TABLE, Name: "item"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListInsertItem(state, id, args["index"].(int), args["item"].(*golua.LTable)))
			return 1
		})

	lib.TableFunction(state, t, "item_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
			{Type: lua.RAW_TABLE, Name: "item"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListSetItem(state, id, args["index"].(int), args["item"].(*golua.LTable)))
			return 1
		})

	lib.BuilderFunction(state, t, "item_remove",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].RemoveItem(args["index"].(int))
		})

	t.RawSetString("selected", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].SelectedItem().(teaListItem)

		state.Push(listItemTable(state, value.title, value.description, value.filter))
		return 1
	}))

	lib.BuilderFunction(state, t, "select",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].Select(args["index"].(int))
		})

	lib.TableFunction(state, t, "matches",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			m := item.Lists[id].MatchesForItem(args["index"].(int))
			matches := state.NewTable()

			for i, v := range m {
				matches.RawSetInt(i+1, golua.LNumber(v))
			}

			state.Push(matches)
			return 1
		})

	lib.TableFunction(state, t, "status_message",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListStatusMessage(state, id, args["msg"].(string)))
			return 1
		})

	lib.TableFunction(state, t, "status_message_lifetime",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Lists[id].StatusMessageLifetime

			state.Push(golua.LNumber(value.Milliseconds()))
			return 1
		})

	lib.BuilderFunction(state, t, "status_message_lifetime_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "duration"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].StatusMessageLifetime = time.Duration(args["duration"].(int) * 1e6)
		})

	lib.TableFunction(state, t, "statusbar_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Lists[id].ShowStatusBar()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "statusbar_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowStatusBar(args["enabled"].(bool))
		})

	lib.TableFunction(state, t, "statusbar_item_name",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			sing, plur := item.Lists[id].StatusBarItemName()

			state.Push(golua.LString(sing))
			state.Push(golua.LString(plur))
			return 2
		})

	lib.BuilderFunction(state, t, "statusbar_item_name_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "singular"},
			{Type: lua.STRING, Name: "plural"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetStatusBarItemName(args["singular"].(string), args["plural"].(string))
		})

	lib.TableFunction(state, t, "title_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Lists[id].ShowTitle()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "title_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowTitle(args["enabled"].(bool))
		})

	lib.BuilderFunction(state, t, "spinner_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "from"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			from := args["from"].(int)
			item.Lists[id].SetSpinner(spinnerList[from])
		})

	lib.BuilderFunction(state, t, "spinner_set_custom",
		[]lua.Arg{
			lua.ArgArray("frames", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.INT, Name: "fps"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			frames := args["frames"].([]any)
			frameBuild := make([]string, len(frames))

			for i, v := range frames {
				frameBuild[i] = v.(string)
			}

			spin := spinner.Spinner{
				Frames: frameBuild,
				FPS:    time.Second / time.Duration(args["fps"].(int)),
			}

			item.Lists[id].SetSpinner(spin)
		})

	lib.TableFunction(state, t, "spinner_start",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListSpinnerStart(state, id))
			return 1
		})

	lib.BuilderFunction(state, t, "spinner_stop",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].StopSpinner()
		})

	lib.TableFunction(state, t, "spinner_toggle",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(cmdListSpinnerToggle(state, id))
			return 1
		})

	lib.TableFunction(state, t, "infinite_scroll",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Lists[id].InfiniteScrolling

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "infinite_scroll_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].InfiniteScrolling = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "filter_input",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].FilterInput
			mid := len(item.TextInputs)
			item.TextInputs = append(item.TextInputs, model)

			state.Push(textinputTable(r, lib, state, program, mid))
			return 1
		})

	lib.TableFunction(state, t, "paginator",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].Paginator
			mid := len(item.Paginators)
			item.Paginators = append(item.Paginators, model)

			state.Push(paginatorTable(r, lib, state, program, mid))
			return 1
		})

	return t
}

func paginatorTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Paginators[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		pg, cmd := item.Paginators[id].Update(*item.Msg)
		item.Paginators[id] = &pg

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "slice_bounds",
		[]lua.Arg{
			{Type: lua.INT, Name: "length"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			start, end := item.Paginators[id].GetSliceBounds(args["length"].(int))

			state.Push(golua.LNumber(start))
			state.Push(golua.LNumber(end))
			return 2
		})

	lib.BuilderFunction(state, t, "page_next",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].NextPage()
		})

	lib.BuilderFunction(state, t, "page_prev",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].PrevPage()
		})

	lib.TableFunction(state, t, "page_items",
		[]lua.Arg{
			{Type: lua.INT, Name: "total"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			per := item.Paginators[id].ItemsOnPage(args["total"].(int))

			state.Push(golua.LNumber(per))
			return 1
		})

	lib.TableFunction(state, t, "page_on_first",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].OnFirstPage()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "page_on_last",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].OnLastPage()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "total_pages_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "items"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			total := item.Paginators[id].SetTotalPages(args["items"].(int))

			state.Push(golua.LNumber(total))
			return 1
		})

	lib.TableFunction(state, t, "type",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			ptype := item.Paginators[id].Type

			state.Push(golua.LNumber(ptype))
			return 1
		})

	lib.BuilderFunction(state, t, "type_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "type"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].Type = paginator.Type(args["type"].(int))
		})

	lib.TableFunction(state, t, "page",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].Page

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "page_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "page"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].Page = args["page"].(int)
		})

	lib.TableFunction(state, t, "page_per",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].PerPage

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "page_per_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "per"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].PerPage = args["per"].(int)
		})

	lib.TableFunction(state, t, "page_total",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].TotalPages

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "page_total_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "total"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].TotalPages = args["total"].(int)
		})

	lib.TableFunction(state, t, "format_dot",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			active := item.Paginators[id].ActiveDot
			inactive := item.Paginators[id].InactiveDot

			state.Push(golua.LString(active))
			state.Push(golua.LString(inactive))
			return 2
		})

	lib.BuilderFunction(state, t, "format_dot_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "active"},
			{Type: lua.STRING, Name: "inactive"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].ActiveDot = args["active"].(string)
			item.Paginators[id].InactiveDot = args["inactive"].(string)
		})

	lib.TableFunction(state, t, "format_arabic",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Paginators[id].ArabicFormat

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "format_arabic_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].ArabicFormat = args["format"].(string)
		})

	return t
}

type ProgressGradient int

const (
	PROGRESSGRADIENT_DEFAULT ProgressGradient = iota
	PROGRESSGRADIENT_DEFAULTSCALED
	PROGRESSGRADIENT_NORMAL
	PROGRESSGRADIENT_NORMALSCALED
	PROGRESSGRADIENT_SOLID
)

func progressOptionsTable(lib *lua.Lib, state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__gradient", golua.LNil)
	t.RawSetString("__colorA", golua.LNil)
	t.RawSetString("__colorB", golua.LNil)
	t.RawSetString("__fullchar", golua.LNil)
	t.RawSetString("__emptychar", golua.LNil)
	t.RawSetString("__springFreq", golua.LNil)
	t.RawSetString("__springDamp", golua.LNil)
	t.RawSetString("__withoutPercent", golua.LNil)

	lib.BuilderFunction(state, t, "width",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__width", golua.LNumber(args["width"].(int)))
		})

	lib.BuilderFunction(state, t, "gradient_default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__gradient", golua.LNumber(PROGRESSGRADIENT_DEFAULT))
		})

	lib.BuilderFunction(state, t, "gradient_default_scaled",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__gradient", golua.LNumber(PROGRESSGRADIENT_DEFAULTSCALED))
		})

	lib.BuilderFunction(state, t, "gradient",
		[]lua.Arg{
			{Type: lua.STRING, Name: "colorA"},
			{Type: lua.STRING, Name: "colorB"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__gradient", golua.LNumber(PROGRESSGRADIENT_NORMAL))
			t.RawSetString("__colorA", golua.LString(args["colorA"].(string)))
			t.RawSetString("__colorB", golua.LString(args["colorB"].(string)))
		})

	lib.BuilderFunction(state, t, "gradient_scaled",
		[]lua.Arg{
			{Type: lua.STRING, Name: "colorA"},
			{Type: lua.STRING, Name: "colorB"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__gradient", golua.LNumber(PROGRESSGRADIENT_NORMALSCALED))
			t.RawSetString("__colorA", golua.LString(args["colorA"].(string)))
			t.RawSetString("__colorB", golua.LString(args["colorB"].(string)))
		})

	lib.BuilderFunction(state, t, "solid",
		[]lua.Arg{
			{Type: lua.STRING, Name: "colorA"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__gradient", golua.LNumber(PROGRESSGRADIENT_SOLID))
			t.RawSetString("__colorA", golua.LString(args["colorA"].(string)))
		})

	lib.BuilderFunction(state, t, "fill_char",
		[]lua.Arg{
			{Type: lua.INT, Name: "full"},
			{Type: lua.INT, Name: "empty"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__fullchar", golua.LNumber(args["full"].(int)))
			t.RawSetString("__emptychar", golua.LNumber(args["empty"].(int)))
		})

	lib.BuilderFunction(state, t, "spring_options",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "freq"},
			{Type: lua.FLOAT, Name: "damp"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__springFreq", golua.LNumber(args["freq"].(float64)))
			t.RawSetString("__springDamp", golua.LNumber(args["damp"].(float64)))
		})

	return t
}

func progressOptionsBuild(t *golua.LTable) []progress.Option {
	opts := []progress.Option{}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		opts = append(opts, progress.WithWidth(int(width.(golua.LNumber))))
	}

	gradient := t.RawGetString("__gradient")
	if gradient.Type() == golua.LTNumber {
		switch ProgressGradient(gradient.(golua.LNumber)) {
		case PROGRESSGRADIENT_DEFAULT:
			opts = append(opts, progress.WithDefaultGradient())
		case PROGRESSGRADIENT_DEFAULTSCALED:
			opts = append(opts, progress.WithDefaultScaledGradient())
		case PROGRESSGRADIENT_NORMAL:
			opts = append(opts, progress.WithGradient(
				string(t.RawGetString("__colorA").(golua.LString)),
				string(t.RawGetString("__colorB").(golua.LString)),
			))
		case PROGRESSGRADIENT_NORMALSCALED:
			opts = append(opts, progress.WithScaledGradient(
				string(t.RawGetString("__colorA").(golua.LString)),
				string(t.RawGetString("__colorB").(golua.LString)),
			))
		case PROGRESSGRADIENT_SOLID:
			opts = append(opts, progress.WithSolidFill(
				string(t.RawGetString("__colorA").(golua.LString)),
			))

		}
	}

	fullchar := t.RawGetString("__fullchar")
	emptychar := t.RawGetString("__emptychar")
	if fullchar.Type() == golua.LTNumber && emptychar.Type() == golua.LTNumber {
		opts = append(opts, progress.WithFillCharacters(rune(fullchar.(golua.LNumber)), rune(emptychar.(golua.LNumber))))
	}

	springFreq := t.RawGetString("__springFreq")
	springDamp := t.RawGetString("__springDamp")
	if springFreq.Type() == golua.LTNumber && springDamp.Type() == golua.LTNumber {
		opts = append(opts, progress.WithSpringOptions(float64(springFreq.(golua.LNumber)), float64(springDamp.(golua.LNumber))))
	}

	return opts
}

func progressTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	lib.TableFunction(state, t, "view_as",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}

			str := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))].ViewAs(args["percent"].(float64))

			state.Push(golua.LString(str))
			return 1
		})

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		pb, cmd := item.ProgressBars[id].Update(*item.Msg)
		pbp := pb.(progress.Model)
		item.ProgressBars[id] = &pbp

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = cmdNone(state)
		} else {
			bcmd = cmdStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "percent",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].Percent()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.TableFunction(state, t, "percent_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			percent := args["percent"].(float64)

			state.Push(cmdProgressSet(state, id, percent))
			return 1
		})

	lib.TableFunction(state, t, "percent_dec",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			percent := args["percent"].(float64)

			state.Push(cmdProgressDec(state, id, percent))
			return 1
		})

	lib.TableFunction(state, t, "percent_inc",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			percent := args["percent"].(float64)

			state.Push(cmdProgressInc(state, id, percent))
			return 1
		})

	lib.TableFunction(state, t, "percent_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].ShowPercentage

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "percent_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].ShowPercentage = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "percent_format",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].PercentFormat

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "percent_format",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].PercentFormat = args["format"].(string)
		})

	lib.TableFunction(state, t, "is_animating",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].IsAnimating()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "spring_options",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "freq"},
			{Type: lua.FLOAT, Name: "damp"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].SetSpringOptions(args["freq"].(float64), args["damp"].(float64))
		})

	lib.TableFunction(state, t, "width",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].Width

			state.Push(golua.LNumber(value))
			return 1
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

			item.ProgressBars[id].Width = args["width"].(int)
		})

	lib.TableFunction(state, t, "full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].Full

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "full_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].Full = rune(args["rune"].(int))
		})

	lib.TableFunction(state, t, "full_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].FullColor

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "full_color_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "color"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].FullColor = args["color"].(string)
		})

	lib.TableFunction(state, t, "empty",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].Empty

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "empty_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].Empty = rune(args["rune"].(int))
		})

	lib.TableFunction(state, t, "empty_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].EmptyColor

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "empty_color_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "color"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ProgressBars[id].EmptyColor = args["color"].(string)
		})

	return t
}
