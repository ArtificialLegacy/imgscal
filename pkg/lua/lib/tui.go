package lib

import (
	"errors"
	"math"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	customtea "github.com/ArtificialLegacy/imgscal/pkg/custom_tea"
	teamodels "github.com/ArtificialLegacy/imgscal/pkg/custom_tea/models"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
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
			id := r.CR_TEA.Add(&teamodels.TeaItem{
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
			p := tea.NewProgram(customtea.ProgramModel{Id: id, Item: item, State: pstate, R: r, Lg: lg})
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
			li := customtea.ListItemTable(state, args["title"].(string), args["desc"].(string), args["filter"].(string))

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
				items[i] = customtea.ListItemBuild(v.(*golua.LTable))
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

	/// @func stopwatch(id, interval?) -> struct<tui.StopWatch>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the stopwatch to.
	/// @arg? interval {int}
	/// @returns {struct<tui.StopWatch>}
	lib.CreateFunction(tab, "stopwatch",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "interval", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			interval := args["interval"].(int)

			var sw stopwatch.Model
			if interval >= 0 {
				sw = stopwatch.NewWithInterval(time.Duration(interval * 1e6))
			} else {
				sw = stopwatch.New()
			}

			id := sw.ID()
			item.StopWatches[id] = &sw

			t := stopwatchTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func timer(id, timeout, interval?) -> struct<tui.Timer>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the timer to.
	/// @arg timeout {int}
	/// @arg? interval {int}
	/// @returns {struct<tui.Timer>}
	lib.CreateFunction(tab, "timer",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "timeout"},
			{Type: lua.INT, Name: "interval", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			timeout := time.Duration(args["timeout"].(int) * 1e6)
			interval := args["interval"].(int)

			var ti timer.Model
			if interval >= 0 {
				ti = timer.NewWithInterval(timeout, time.Duration(interval*1e6))
			} else {
				ti = timer.New(timeout)
			}

			id := ti.ID()
			item.Timers[id] = &ti

			t := timerTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func table_options() -> struct<tui.TableOptions>
	/// @returns {struct<tui.TableOptions>}
	lib.CreateFunction(tab, "table_options",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(tableOptionsTable(lib, state))
			return 1
		})

	/// @func table_column(title, width) -> struct<tui.TableColumn>
	/// @arg title {string}
	/// @arg width {int}
	/// @returns {struct<tui.TableColumn>}
	lib.CreateFunction(tab, "table_column",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tuitableColTable(state, args["title"].(string), args["width"].(int))
			state.Push(t)
			return 1
		})

	/// @func table(id, options?) -> struct<tui.Table>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the table to.
	/// @arg? options {struct<tui.TableOptions>}
	/// @returns {struct<tui.Table>}
	lib.CreateFunction(tab, "table",
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

			opts := tableOptionsBuild(args["options"].(*golua.LTable))

			tb := table.New(opts...)
			id := len(item.Tables)
			item.Tables = append(item.Tables, &tb)

			t := tuitableTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func viewport(id, width, height) -> struct<tui.Viewport>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the viewport to.
	/// @arg width {int}
	/// @arg height {int}
	/// @returns {struct<tui.Viewport>}
	lib.CreateFunction(tab, "viewport",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			width := args["width"].(int)
			height := args["height"].(int)
			id := len(item.Viewports)
			vp := viewport.New(width, height)
			item.Viewports = append(item.Viewports, &vp)

			t := viewportTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func viewport_sync(model) -> struct<tui.CMDViewportSync>
	/// @arg model {int} - ID of the viewport model.
	/// @returns {struct<tui.CMDViewportSync>}
	lib.CreateFunction(tab, "viewport_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := customtea.CMDViewportSync(state, args["model"].(int))
			state.Push(t)
			return 1
		})

	/// @func viewport_view_up(model, lines) -> struct<tui.CMDViewportUp>
	/// @arg model {int} - ID of the viewport model.
	/// @arg lines {[]string}
	/// @returns {struct<tui.CMDViewportUp>}
	lib.CreateFunction(tab, "viewport_view_up",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
			{Type: lua.RAW_TABLE, Name: "lines"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := customtea.CMDViewportUp(state, args["model"].(int), args["lines"].(*golua.LTable))
			state.Push(t)
			return 1
		})

	/// @func viewport_view_down(model, lines) -> struct<tui.CMDViewportDown>
	/// @arg model {int} - ID of the viewport model.
	/// @arg lines {[]string}
	/// @returns {struct<tui.CMDViewportDown>}
	lib.CreateFunction(tab, "viewport_view_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "model"},
			{Type: lua.RAW_TABLE, Name: "lines"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := customtea.CMDViewportDown(state, args["model"].(int), args["lines"].(*golua.LTable))
			state.Push(t)
			return 1
		})

	/// @func custom(id, init, update, view) -> struct<tui.Custom>
	/// @arg id {int<collection.CRATE_TEA>}
	/// @arg init {function(id int) -> table<any>, struct<tui.CMD>}
	/// @arg update {function(data table<any>, msg struct<tui.MSG>) -> struct<tui.CMD>}
	/// @arg view {function(data table<any>) -> string}
	/// @returns {struct<tui.Custom>}
	lib.CreateFunction(tab, "custom",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "init"},
			{Type: lua.FUNC, Name: "update"},
			{Type: lua.FUNC, Name: "view"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			init := args["init"].(*golua.LFunction)
			update := args["update"].(*golua.LFunction)
			view := args["view"].(*golua.LFunction)

			id := len(item.Customs)
			cm := teamodels.NewCustomModel(prgrm, init, update, view, state, item, customtea.CMDBuild, customtea.BuildMSG)
			item.Customs = append(item.Customs, &cm)

			t := tuicustomTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func keybinding_option() -> struct<tui.KeyOption>
	/// @returns {struct<tui.KeyOption>}
	lib.CreateFunction(tab, "keybinding_option",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := keyOptionsTable(state, lib)
			state.Push(t)
			return 1
		})

	/// @func keybinding(id, option) -> struct<tui.Keybinding>
	/// @arg id {int<collection.CRATE_TEA>}
	/// @arg? option {struct<tui.KeyOption>}
	/// @returns {struct<tui.Keybinding>}
	lib.CreateFunction(tab, "keybinding",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "option", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			options := args["option"].(*golua.LTable)
			opts := keyOptionsBuild(options)

			id := len(item.KeyBindings)
			ky := key.NewBinding(opts...)
			item.KeyBindings = append(item.KeyBindings, &ky)

			t := tuikeyTable(r, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func key_match(msg, keybindings...) -> bool
	/// @arg msg {struct<tui.MSGKey>}
	/// @arg keybinding {struct<tui.Keybinding>...}
	/// @returns {bool}
	lib.CreateFunction(tab, "key_match",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "msg"},
			lua.ArgVariadic("bindings", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			msg := args["msg"].(*golua.LTable)
			bindings := args["bindings"].([]any)

			mtype := msg.RawGetString("msg")
			if mtype.Type() != golua.LTNumber {
				state.Push(golua.LFalse)
				return 1
			}
			if customtea.TeaMSG(mtype.(golua.LNumber)) != customtea.MSG_KEY {
				state.Push(golua.LFalse)
				return 1
			}

			mk := msg.RawGetString("key").(golua.LString)

			blist := make([]key.Binding, len(bindings))
			for i, v := range bindings {
				vt := v.(*golua.LTable)
				prgrm := int(vt.RawGetString("program").(golua.LNumber))
				item, err := r.CR_TEA.Item(prgrm)
				if err != nil {
					lua.Error(state, err.Error())
				}
				id := int(vt.RawGetString("id").(golua.LNumber))
				blist[i] = *item.KeyBindings[id]
			}

			matches := key.Matches(mk, blist...)

			state.Push(golua.LBool(matches))
			return 1
		})

	/// @func help(id) -> struct<tui.Help>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the help to.
	/// @returns {struct<tui.Help>}
	lib.CreateFunction(tab, "help",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, err.Error())
			}

			hp := help.New()
			id := len(item.Helps)
			item.Helps = append(item.Helps, &hp)

			t := helpTable(r, lib, state, prgrm, id)

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
			state.Push(customtea.CMDNone(state))
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
			t := customtea.CMDBatch(state, args["cmds"].(*golua.LTable))

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
			t := customtea.CMDSequence(state, args["cmds"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @constants Message
	/// @const MSG_NONE
	/// @const MSG_BLUR
	/// @const MSG_FOCUS
	/// @const MSG_QUIT
	/// @const MSG_RESUME
	/// @const MSG_SUSPEND
	/// @const MSG_WINDOWSIZE
	/// @const MSG_KEY
	/// @const MSG_MOUSE
	/// @const MSG_SPINNERTICK
	/// @const MSG_BLINK
	/// @const MSG_STOPWATCHRESET
	/// @const MSG_STOPWATCHSTARTSTOP
	/// @const MSG_STOPWATCHTICK
	/// @const MSG_TIMERSTARTSTOP
	/// @const MSG_TIMERTICK
	/// @const MSG_TIMERTIMEOUT
	tab.RawSetString("MSG_NONE", golua.LNumber(customtea.MSG_NONE))
	tab.RawSetString("MSG_BLUR", golua.LNumber(customtea.MSG_BLUR))
	tab.RawSetString("MSG_FOCUS", golua.LNumber(customtea.MSG_FOCUS))
	tab.RawSetString("MSG_QUIT", golua.LNumber(customtea.MSG_QUIT))
	tab.RawSetString("MSG_RESUME", golua.LNumber(customtea.MSG_RESUME))
	tab.RawSetString("MSG_SUSPEND", golua.LNumber(customtea.MSG_SUSPEND))
	tab.RawSetString("MSG_WINDOWSIZE", golua.LNumber(customtea.MSG_WINDOWSIZE))
	tab.RawSetString("MSG_KEY", golua.LNumber(customtea.MSG_KEY))
	tab.RawSetString("MSG_MOUSE", golua.LNumber(customtea.MSG_MOUSE))
	tab.RawSetString("MSG_SPINNERTICK", golua.LNumber(customtea.MSG_SPINNERTICK))
	tab.RawSetString("MSG_BLINK", golua.LNumber(customtea.MSG_BLINK))
	tab.RawSetString("MSG_STOPWATCHRESET", golua.LNumber(customtea.MSG_STOPWATCHRESET))
	tab.RawSetString("MSG_STOPWATCHSTARTSTOP", golua.LNumber(customtea.MSG_STOPWATCHSTARTSTOP))
	tab.RawSetString("MSG_STOPWATCHTICK", golua.LNumber(customtea.MSG_STOPWATCHTICK))
	tab.RawSetString("MSG_TIMERSTARTSTOP", golua.LNumber(customtea.MSG_TIMERSTARTSTOP))
	tab.RawSetString("MSG_TIMERTICK", golua.LNumber(customtea.MSG_TIMERTICK))
	tab.RawSetString("MSG_TIMERTIMEOUT", golua.LNumber(customtea.MSG_TIMERTIMEOUT))

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

func teaTable(r *lua.Runner, state *golua.LState, lib *lua.Lib, id int) *golua.LTable {
	/// @struct Program
	/// @prop id {int}
	/// @method init(fn: function(id) -> any, struct<tui.CMD>)
	/// @method update(fn: function(data: any, msg: struct<tui.MSG>) -> struct<tui.CMD>)
	/// @method view(fn: function(data: any) -> string)

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

func spinnerTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Spinner
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method tick() -> struct<tui.CMD>

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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("tick", state.NewFunction(func(state *golua.LState) int {
		cmd := customtea.CMDSpinnerTick(state, int(t.RawGetString("id").(golua.LNumber)))

		state.Push(cmd)
		return 1
	}))

	lib.TableFunction(state, t, "spinner",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			spinner := item.Spinners[id].Spinner

			fps := math.Floor(float64(time.Second / spinner.FPS))

			frames := state.NewTable()
			for i, v := range spinner.Frames {
				frames.RawSetInt(i+1, golua.LString(v))
			}

			state.Push(frames)
			state.Push(golua.LNumber(fps))
			return 2
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
			item.Spinners[id].Spinner = spinnerList[from]
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

			item.Spinners[id].Spinner = spin
		})

	t.RawSetString("__style", golua.LNil)
	lib.TableFunction(state, t, "style",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__style")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Spinners[id].Style
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__style", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Spinners[id].Style = *style.Style
			t.RawSetString("__style", st)
		})

	return t
}

func textareaTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct TextArea
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method reset()
	/// @method focus() -> struct<tui.CMD>
	/// @method blur()
	/// @method cursor_down()
	/// @method cursor_end()
	/// @method cursor_up()
	/// @method cursor_down()
	/// @method focused() -> bool
	/// @method size() -> int, int
	/// @method width() -> int
	/// @method height() -> int
	/// @method size_set(width int, height int)
	/// @method width_set(width int)
	/// @method height_set(height int)
	/// @method insert_rune(rune int)
	/// @method insert_string(str string)
	/// @method length() -> int
	/// @method line() -> int
	/// @method line_count() -> int
	/// @method cursor_set(col int)
	/// @method value() -> string
	/// @method value_set(str string)
	/// @method line_info() -> struct<tui.LineInfo>

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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
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
		t := customtea.CMDTextAreaFocus(state, int(t.RawGetString("id").(golua.LNumber)))

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

	t.RawSetString("__cursor", golua.LNil)
	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		oc := t.RawGetString("__cursor")
		if oc.Type() == golua.LTTable {
			state.Push(oc)
			return 1
		}

		program := int(t.RawGetString("program").(golua.LNumber))
		item, err := r.CR_TEA.Item(program)
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextAreas[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)
		cu := cursorTable(r, lib, state, program, cid)

		state.Push(cu)
		t.RawSetString("__cursor", cu)
		return 1
	}))

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.TextAreas[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.CharacterForward,
				&value.CharacterBackward,
				&value.DeleteAfterCursor,
				&value.DeleteBeforeCursor,
				&value.DeleteCharacterBackward,
				&value.DeleteCharacterForward,
				&value.DeleteWordBackward,
				&value.DeleteWordForward,
				&value.InsertNewline,
				&value.LineEnd,
				&value.LineNext,
				&value.LinePrevious,
				&value.LineStart,
				&value.Paste,
				&value.WordBackward,
				&value.WordForward,
				&value.InputBegin,
				&value.InputEnd,
				&value.UppercaseWordForward,
				&value.LowercaseWordForward,
				&value.CapitalizeWordForward,
				&value.TransposeCharacterBackward,
			)

			ids := [22]int{}
			for i := range 22 {
				ids[i] = start + i
			}

			kmt := textareaKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	return t
}

func textareaKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [22]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("character_backward", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("character_forward", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("delete_after_cursor", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("delete_before_cursor", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("delete_character_backward", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("delete_character_forward", tuikeyTable(r, lib, state, program, ids[5]))
	t.RawSetString("delete_word_backward", tuikeyTable(r, lib, state, program, ids[6]))
	t.RawSetString("delete_word_forward", tuikeyTable(r, lib, state, program, ids[7]))
	t.RawSetString("insert_newline", tuikeyTable(r, lib, state, program, ids[8]))
	t.RawSetString("line_end", tuikeyTable(r, lib, state, program, ids[9]))
	t.RawSetString("line_next", tuikeyTable(r, lib, state, program, ids[10]))
	t.RawSetString("line_previous", tuikeyTable(r, lib, state, program, ids[11]))
	t.RawSetString("line_start", tuikeyTable(r, lib, state, program, ids[12]))
	t.RawSetString("paste", tuikeyTable(r, lib, state, program, ids[13]))
	t.RawSetString("word_backward", tuikeyTable(r, lib, state, program, ids[14]))
	t.RawSetString("word_forward", tuikeyTable(r, lib, state, program, ids[15]))
	t.RawSetString("input_begin", tuikeyTable(r, lib, state, program, ids[16]))
	t.RawSetString("input_end", tuikeyTable(r, lib, state, program, ids[17]))
	t.RawSetString("uppercase_word", tuikeyTable(r, lib, state, program, ids[18]))
	t.RawSetString("lowercase_word", tuikeyTable(r, lib, state, program, ids[19]))
	t.RawSetString("capitalize_word", tuikeyTable(r, lib, state, program, ids[20]))
	t.RawSetString("transpose_character_backward", tuikeyTable(r, lib, state, program, ids[21]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			ta := item.TextAreas[id]
			ta.KeyMap = textarea.DefaultKeyMap
			item.KeyBindings[ids[0]] = &ta.KeyMap.CharacterBackward
			item.KeyBindings[ids[1]] = &ta.KeyMap.CharacterForward
			item.KeyBindings[ids[2]] = &ta.KeyMap.DeleteAfterCursor
			item.KeyBindings[ids[3]] = &ta.KeyMap.DeleteBeforeCursor
			item.KeyBindings[ids[4]] = &ta.KeyMap.DeleteCharacterBackward
			item.KeyBindings[ids[5]] = &ta.KeyMap.DeleteCharacterForward
			item.KeyBindings[ids[6]] = &ta.KeyMap.DeleteWordBackward
			item.KeyBindings[ids[7]] = &ta.KeyMap.DeleteWordForward
			item.KeyBindings[ids[8]] = &ta.KeyMap.InsertNewline
			item.KeyBindings[ids[9]] = &ta.KeyMap.LineEnd
			item.KeyBindings[ids[10]] = &ta.KeyMap.LineNext
			item.KeyBindings[ids[11]] = &ta.KeyMap.LinePrevious
			item.KeyBindings[ids[12]] = &ta.KeyMap.LineStart
			item.KeyBindings[ids[13]] = &ta.KeyMap.Paste
			item.KeyBindings[ids[14]] = &ta.KeyMap.WordBackward
			item.KeyBindings[ids[15]] = &ta.KeyMap.WordForward
			item.KeyBindings[ids[16]] = &ta.KeyMap.InputBegin
			item.KeyBindings[ids[17]] = &ta.KeyMap.InputEnd
			item.KeyBindings[ids[18]] = &ta.KeyMap.UppercaseWordForward
			item.KeyBindings[ids[19]] = &ta.KeyMap.LowercaseWordForward
			item.KeyBindings[ids[20]] = &ta.KeyMap.CapitalizeWordForward
			item.KeyBindings[ids[21]] = &ta.KeyMap.TransposeCharacterBackward
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()
			kt.RawSetInt(1, t.RawGetString("paste"))
			kt.RawSetInt(2, t.RawGetString("uppercase_word"))
			kt.RawSetInt(3, t.RawGetString("lowercase_word"))
			kt.RawSetInt(4, t.RawGetString("capitalize_word"))
			kt.RawSetInt(5, t.RawGetString("transpose_character_backward"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("character_forward"))
			kt1.RawSetInt(2, t.RawGetString("character_backward"))
			kt1.RawSetInt(3, t.RawGetString("word_forward"))
			kt1.RawSetInt(4, t.RawGetString("word_backward"))
			kt1.RawSetInt(5, t.RawGetString("line_start"))
			kt1.RawSetInt(6, t.RawGetString("line_end"))
			kt1.RawSetInt(7, t.RawGetString("line_next"))
			kt1.RawSetInt(8, t.RawGetString("line_previous"))
			kt1.RawSetInt(9, t.RawGetString("input_begin"))
			kt1.RawSetInt(10, t.RawGetString("input_end"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("delete_character_backward"))
			kt2.RawSetInt(2, t.RawGetString("delete_character_forward"))
			kt2.RawSetInt(3, t.RawGetString("delete_word_forward"))
			kt2.RawSetInt(4, t.RawGetString("delete_word_backward"))
			kt2.RawSetInt(5, t.RawGetString("delete_before_cursor"))
			kt2.RawSetInt(6, t.RawGetString("delete_after_cursor"))

			kt3 := state.NewTable()
			kt3.RawSetInt(1, t.RawGetString("insert_newline"))
			kt3.RawSetInt(2, t.RawGetString("paste"))
			kt3.RawSetInt(3, t.RawGetString("uppercase_word"))
			kt3.RawSetInt(4, t.RawGetString("lowercase_word"))
			kt3.RawSetInt(6, t.RawGetString("capitalize_word"))
			kt3.RawSetInt(7, t.RawGetString("transpose_character_backward"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)
			kt.RawSetInt(3, kt3)

			state.Push(kt)
			return 1
		})

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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("focus", state.NewFunction(func(state *golua.LState) int {
		t := customtea.CMDTextInputFocus(state, int(t.RawGetString("id").(golua.LNumber)))

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

	t.RawSetString("__cursor", golua.LNil)
	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		oc := t.RawGetString("__cursor")
		if oc.Type() == golua.LTTable {
			state.Push(oc)
			return 1
		}

		program := int(t.RawGetString("program").(golua.LNumber))
		item, err := r.CR_TEA.Item(program)
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextInputs[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)
		cu := cursorTable(r, lib, state, program, cid)

		state.Push(cu)
		t.RawSetString("__cursor", cu)
		return 1
	}))

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.TextInputs[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.CharacterForward,
				&value.CharacterBackward,
				&value.WordForward,
				&value.WordBackward,
				&value.DeleteWordBackward,
				&value.DeleteWordForward,
				&value.DeleteAfterCursor,
				&value.DeleteBeforeCursor,
				&value.DeleteCharacterBackward,
				&value.DeleteCharacterForward,
				&value.LineStart,
				&value.LineEnd,
				&value.Paste,
				&value.AcceptSuggestion,
				&value.NextSuggestion,
				&value.PrevSuggestion,
			)

			ids := [16]int{}
			for i := range 16 {
				ids[i] = start + i
			}

			kmt := textinputKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	t.RawSetString("__stylePrompt", golua.LNil)
	lib.TableFunction(state, t, "style_prompt",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePrompt")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextInputs[id].PromptStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePrompt", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_prompt_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextInputs[id].PromptStyle = *style.Style
			t.RawSetString("__stylePrompt", st)
		})

	t.RawSetString("__styleText", golua.LNil)
	lib.TableFunction(state, t, "style_text",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleText")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextInputs[id].TextStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleText", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_text_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextInputs[id].TextStyle = *style.Style
			t.RawSetString("__styleText", st)
		})

	t.RawSetString("__stylePlaceholder", golua.LNil)
	lib.TableFunction(state, t, "style_placeholder",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePlaceholder")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextInputs[id].PlaceholderStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePlaceholder", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_placeholder_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextInputs[id].PlaceholderStyle = *style.Style
			t.RawSetString("__stylePlaceholder", st)
		})

	t.RawSetString("__styleCompletion", golua.LNil)
	lib.TableFunction(state, t, "style_completion",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleCompletion")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextInputs[id].CompletionStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleCompletion", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_completion_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextInputs[id].CompletionStyle = *style.Style
			t.RawSetString("__styleCompletion", st)
		})

	return t
}

func textinputKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [16]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("character_forward", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("character_backward", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("word_forward", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("word_backward", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("delete_word_backward", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("delete_word_forward", tuikeyTable(r, lib, state, program, ids[5]))
	t.RawSetString("delete_after_cursor", tuikeyTable(r, lib, state, program, ids[6]))
	t.RawSetString("delete_before_cursor", tuikeyTable(r, lib, state, program, ids[7]))
	t.RawSetString("delete_character_backward", tuikeyTable(r, lib, state, program, ids[8]))
	t.RawSetString("delete_character_forward", tuikeyTable(r, lib, state, program, ids[9]))
	t.RawSetString("line_start", tuikeyTable(r, lib, state, program, ids[10]))
	t.RawSetString("line_end", tuikeyTable(r, lib, state, program, ids[11]))
	t.RawSetString("paste", tuikeyTable(r, lib, state, program, ids[12]))
	t.RawSetString("suggestion_accept", tuikeyTable(r, lib, state, program, ids[13]))
	t.RawSetString("suggestion_next", tuikeyTable(r, lib, state, program, ids[14]))
	t.RawSetString("suggestion_prev", tuikeyTable(r, lib, state, program, ids[15]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			ti := item.TextInputs[id]
			ti.KeyMap = textinput.DefaultKeyMap
			item.KeyBindings[ids[0]] = &ti.KeyMap.CharacterForward
			item.KeyBindings[ids[1]] = &ti.KeyMap.CharacterBackward
			item.KeyBindings[ids[2]] = &ti.KeyMap.WordForward
			item.KeyBindings[ids[3]] = &ti.KeyMap.WordBackward
			item.KeyBindings[ids[4]] = &ti.KeyMap.DeleteWordBackward
			item.KeyBindings[ids[5]] = &ti.KeyMap.DeleteWordForward
			item.KeyBindings[ids[6]] = &ti.KeyMap.DeleteAfterCursor
			item.KeyBindings[ids[7]] = &ti.KeyMap.DeleteBeforeCursor
			item.KeyBindings[ids[8]] = &ti.KeyMap.DeleteCharacterBackward
			item.KeyBindings[ids[9]] = &ti.KeyMap.DeleteCharacterForward
			item.KeyBindings[ids[10]] = &ti.KeyMap.LineStart
			item.KeyBindings[ids[11]] = &ti.KeyMap.LineEnd
			item.KeyBindings[ids[12]] = &ti.KeyMap.Paste
			item.KeyBindings[ids[13]] = &ti.KeyMap.AcceptSuggestion
			item.KeyBindings[ids[14]] = &ti.KeyMap.NextSuggestion
			item.KeyBindings[ids[15]] = &ti.KeyMap.PrevSuggestion
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()
			kt.RawSetInt(1, t.RawGetString("paste"))
			kt.RawSetInt(2, t.RawGetString("suggestion_accept"))
			kt.RawSetInt(3, t.RawGetString("suggestion_next"))
			kt.RawSetInt(4, t.RawGetString("suggestion_prev"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("character_forward"))
			kt1.RawSetInt(2, t.RawGetString("character_backward"))
			kt1.RawSetInt(3, t.RawGetString("word_forward"))
			kt1.RawSetInt(4, t.RawGetString("word_backward"))
			kt1.RawSetInt(5, t.RawGetString("line_start"))
			kt1.RawSetInt(6, t.RawGetString("line_end"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("delete_character_backward"))
			kt2.RawSetInt(2, t.RawGetString("delete_character_forward"))
			kt2.RawSetInt(3, t.RawGetString("delete_word_forward"))
			kt2.RawSetInt(4, t.RawGetString("delete_word_backward"))
			kt2.RawSetInt(5, t.RawGetString("delete_before_cursor"))
			kt2.RawSetInt(6, t.RawGetString("delete_after_cursor"))

			kt3 := state.NewTable()
			kt3.RawSetInt(1, t.RawGetString("paste"))
			kt3.RawSetInt(2, t.RawGetString("suggestion_accept"))
			kt3.RawSetInt(3, t.RawGetString("suggestion_next"))
			kt3.RawSetInt(4, t.RawGetString("suggestion_prev"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)
			kt.RawSetInt(3, kt3)

			state.Push(kt)
			return 1
		})

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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("blink", state.NewFunction(func(state *golua.LState) int {
		state.Push(customtea.CMDBlink(state, int(t.RawGetString("id").(golua.LNumber))))
		return 1
	}))

	t.RawSetString("focus", state.NewFunction(func(state *golua.LState) int {
		t := customtea.CMDCursorFocus(state, int(t.RawGetString("id").(golua.LNumber)))

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

	t.RawSetString("__style", golua.LNil)
	lib.TableFunction(state, t, "style",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__style")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Cursors[id].Style
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__style", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Cursors[id].Style = *style.Style
			t.RawSetString("__style", st)
		})

	t.RawSetString("__styleText", golua.LNil)
	lib.TableFunction(state, t, "style_text",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleText")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Cursors[id].TextStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleText", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_text_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Cursors[id].TextStyle = *style.Style
			t.RawSetString("__styleText", st)
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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
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
		state.Push(customtea.CMDFilePickerInit(state, int(t.RawGetString("id").(golua.LNumber))))
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

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.FilePickers[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.GoToTop,
				&value.GoToLast,
				&value.Down,
				&value.Up,
				&value.PageUp,
				&value.PageDown,
				&value.Back,
				&value.Open,
				&value.Select,
			)

			ids := [9]int{}
			for i := range 9 {
				ids[i] = start + i
			}

			kmt := filepickerKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	t.RawSetString("__styleCursor", golua.LNil)
	lib.TableFunction(state, t, "style_cursor",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleCursor")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.Cursor
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleCursor", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_cursor_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.Cursor = *style.Style
			t.RawSetString("__styleCursor", st)
		})

	t.RawSetString("__styleCursorDisabled", golua.LNil)
	lib.TableFunction(state, t, "style_cursor_disabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleCursorDisabled")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.DisabledCursor
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleCursorDisabled", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_cursor_disabled_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.DisabledCursor = *style.Style
			t.RawSetString("__styleCursorDisabled", st)
		})

	t.RawSetString("__styleSymlink", golua.LNil)
	lib.TableFunction(state, t, "style_symlink",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleSymlink")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.Symlink
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleSymlink", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_symlink_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.Symlink = *style.Style
			t.RawSetString("__styleSymlink", st)
		})

	t.RawSetString("__styleDirectory", golua.LNil)
	lib.TableFunction(state, t, "style_directory",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDirectory")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.Directory
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDirectory", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_directory_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.Directory = *style.Style
			t.RawSetString("__styleDirectory", st)
		})

	t.RawSetString("__styleDirectoryEmpty", golua.LNil)
	lib.TableFunction(state, t, "style_directory_empty",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDirectoryEmpty")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.EmptyDirectory
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDirectoryEmpty", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_directory_empty_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.EmptyDirectory = *style.Style
			t.RawSetString("__styleDirectoryEmpty", st)
		})

	t.RawSetString("__styleFile", golua.LNil)
	lib.TableFunction(state, t, "style_file",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFile")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.File
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFile", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_file_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.File = *style.Style
			t.RawSetString("__styleFile", st)
		})

	t.RawSetString("__styleFileSize", golua.LNil)
	lib.TableFunction(state, t, "style_file_size",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFileSize")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.FileSize
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFileSize", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_file_size_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.FileSize = *style.Style
			t.RawSetString("__styleFileSize", st)
		})

	t.RawSetString("__styleFileDisabled", golua.LNil)
	lib.TableFunction(state, t, "style_file_disabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFileDisabled")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.DisabledFile
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFileDisabled", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_file_disabled_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.DisabledFile = *style.Style
			t.RawSetString("__styleFileDisabled", st)
		})

	t.RawSetString("__stylePermission", golua.LNil)
	lib.TableFunction(state, t, "style_permission",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePermission")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.Permission
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePermission", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_permission_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.Permission = *style.Style
			t.RawSetString("__stylePermission", st)
		})

	t.RawSetString("__styleSelected", golua.LNil)
	lib.TableFunction(state, t, "style_selected",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleSelected")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.Selected
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleSelected", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_selected_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.Selected = *style.Style
			t.RawSetString("__styleSelected", st)
		})

	t.RawSetString("__styleSelectedDisabled", golua.LNil)
	lib.TableFunction(state, t, "style_selected_disabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleSelectedDisabled")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.FilePickers[id].Styles.DisabledSelected
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleSelectedDisabled", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_selected_disabled_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.FilePickers[id].Styles.DisabledSelected = *style.Style
			t.RawSetString("__styleSelectedDisabled", st)
		})

	return t
}

func filepickerKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [9]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("goto_top", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("goto_last", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("down", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("up", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("page_up", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("page_down", tuikeyTable(r, lib, state, program, ids[5]))
	t.RawSetString("back", tuikeyTable(r, lib, state, program, ids[6]))
	t.RawSetString("open", tuikeyTable(r, lib, state, program, ids[7]))
	t.RawSetString("select", tuikeyTable(r, lib, state, program, ids[8]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fp := item.FilePickers[id]
			fp.KeyMap = filepicker.DefaultKeyMap()
			item.KeyBindings[ids[0]] = &fp.KeyMap.GoToTop
			item.KeyBindings[ids[1]] = &fp.KeyMap.GoToLast
			item.KeyBindings[ids[2]] = &fp.KeyMap.Down
			item.KeyBindings[ids[3]] = &fp.KeyMap.Up
			item.KeyBindings[ids[4]] = &fp.KeyMap.PageUp
			item.KeyBindings[ids[5]] = &fp.KeyMap.PageDown
			item.KeyBindings[ids[6]] = &fp.KeyMap.Back
			item.KeyBindings[ids[7]] = &fp.KeyMap.Open
			item.KeyBindings[ids[8]] = &fp.KeyMap.Select
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()
			kt.RawSetInt(1, t.RawGetString("up"))
			kt.RawSetInt(2, t.RawGetString("down"))
			kt.RawSetInt(3, t.RawGetString("goto_top"))
			kt.RawSetInt(4, t.RawGetString("goto_last"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("up"))
			kt1.RawSetInt(2, t.RawGetString("down"))
			kt1.RawSetInt(3, t.RawGetString("goto_top"))
			kt1.RawSetInt(4, t.RawGetString("goto_last"))
			kt1.RawSetInt(5, t.RawGetString("page_up"))
			kt1.RawSetInt(6, t.RawGetString("page_down"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("back"))
			kt2.RawSetInt(2, t.RawGetString("open"))
			kt2.RawSetInt(3, t.RawGetString("select"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)

			state.Push(kt)
			return 1
		})

	return t
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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
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
			li := v.(customtea.ListItem)
			items.RawSetInt(i+1, customtea.ListItemTableFrom(state, li))
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
			li := v.(customtea.ListItem)
			items.RawSetInt(i+1, customtea.ListItemTableFrom(state, li))
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

			state.Push(customtea.CMDListSetItems(state, id, args["items"].(*golua.LTable)))
			return 1
		})

	lib.TableFunction(state, t, "item_insert",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
			{Type: lua.RAW_TABLE, Name: "item"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDListInsertItem(state, id, args["index"].(int), args["item"].(*golua.LTable)))
			return 1
		})

	lib.TableFunction(state, t, "item_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
			{Type: lua.RAW_TABLE, Name: "item"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDListSetItem(state, id, args["index"].(int), args["item"].(*golua.LTable)))
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

		li := item.Lists[id].SelectedItem().(customtea.ListItem)

		state.Push(customtea.ListItemTableFrom(state, li))
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

			state.Push(customtea.CMDListStatusMessage(state, id, args["msg"].(string)))
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

			state.Push(customtea.CMDListSpinnerStart(state, id))
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

			state.Push(customtea.CMDListSpinnerToggle(state, id))
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

	t.RawSetString("__filterInput", golua.LNil)
	lib.TableFunction(state, t, "filter_input",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			ofi := t.RawGetString("__filterInput")
			if ofi.Type() == golua.LTTable {
				state.Push(ofi)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].FilterInput
			mid := len(item.TextInputs)
			item.TextInputs = append(item.TextInputs, model)

			fi := textinputTable(r, lib, state, program, mid)
			state.Push(fi)
			t.RawSetString("__filterInput", fi)
			return 1
		})

	t.RawSetString("__paginator", golua.LNil)
	lib.TableFunction(state, t, "paginator",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			opg := t.RawGetString("__paginator")
			if opg.Type() == golua.LTTable {
				state.Push(opg)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].Paginator
			mid := len(item.Paginators)
			item.Paginators = append(item.Paginators, model)

			pg := paginatorTable(r, lib, state, program, mid)
			state.Push(pg)
			t.RawSetString("__paginator", pg)
			return 1
		})

	t.RawSetString("__help", golua.LNil)
	lib.TableFunction(state, t, "help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			oh := t.RawGetString("__help")
			if oh.Type() == golua.LTTable {
				state.Push(oh)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].Help
			mid := len(item.Helps)
			item.Helps = append(item.Helps, model)

			hp := helpTable(r, lib, state, program, mid)
			state.Push(hp)
			t.RawSetString("__help", hp)
			return 1
		})

	lib.TableFunction(state, t, "help_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Lists[id].ShowHelp()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "help_show_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowHelp(args["enabled"].(bool))
		})

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.Lists[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.CursorUp,
				&value.CursorDown,
				&value.NextPage,
				&value.PrevPage,
				&value.GoToStart,
				&value.GoToEnd,
				&value.Filter,
				&value.ClearFilter,
				&value.CancelWhileFiltering,
				&value.AcceptWhileFiltering,
				&value.ShowFullHelp,
				&value.CloseFullHelp,
				&value.Quit,
				&value.ForceQuit,
			)

			ids := [14]int{}
			for i := range 14 {
				ids[i] = start + i
			}

			kmt := listKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	lib.TableFunction(state, t, "view_help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			li := item.Lists[id]
			help := li.Help
			str := help.View(li)

			state.Push(golua.LString(str))
			return 1
		})

	lib.TableFunction(state, t, "view_help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			li := item.Lists[id]
			help := li.Help
			str := help.ShortHelpView(li.ShortHelp())

			state.Push(golua.LString(str))
			return 1
		})

	lib.TableFunction(state, t, "view_help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			li := item.Lists[id]
			help := li.Help
			str := help.FullHelpView(li.FullHelp())

			state.Push(golua.LString(str))
			return 1
		})

	lib.BuilderFunction(state, t, "help_short_additional",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.Lists[id].AdditionalShortHelpKeys = func() []key.Binding {
				state.Push(fn)
				state.Call(0, 1)
				bindings := state.CheckTable(-1)
				state.Pop(1)

				bindingList := make([]key.Binding, bindings.Len())
				for i := range bindings.Len() {
					b := bindings.RawGetInt(i + 1).(*golua.LTable)
					bid := b.RawGetString("id").(golua.LNumber)
					bindingList[i] = *item.KeyBindings[int(bid)]
				}

				return bindingList
			}
		})

	lib.BuilderFunction(state, t, "help_full_additional",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.Lists[id].AdditionalFullHelpKeys = func() []key.Binding {
				state.Push(fn)
				state.Call(0, 1)
				bindings := state.CheckTable(-1)
				state.Pop(1)

				bindingList := make([]key.Binding, bindings.Len())
				for i := range bindings.Len() {
					b := bindings.RawGetInt(i + 1).(*golua.LTable)
					bid := b.RawGetString("id").(golua.LNumber)
					bindingList[i] = *item.KeyBindings[int(bid)]
				}

				return bindingList
			}
		})

	t.RawSetString("__styleTitlebar", golua.LNil)
	lib.TableFunction(state, t, "style_titlebar",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleTitlebar")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.TitleBar
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleTitlebar", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_titlebar_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.TitleBar = *style.Style
			t.RawSetString("__styleTitlebar", st)
		})

	t.RawSetString("__styleTitle", golua.LNil)
	lib.TableFunction(state, t, "style_title",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleTitle")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.Title
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleTitle", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_title_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.Title = *style.Style
			t.RawSetString("__styleTitle", st)
		})

	t.RawSetString("__styleSpinner", golua.LNil)
	lib.TableFunction(state, t, "style_spinner",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleSpinner")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.Spinner
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleSpinner", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_spinner_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.Spinner = *style.Style
			t.RawSetString("__styleSpinner", st)
		})

	t.RawSetString("__styleFilterPrompt", golua.LNil)
	lib.TableFunction(state, t, "style_filter_prompt",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFilterPrompt")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.FilterPrompt
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFilterPrompt", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_filter_prompt_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.FilterPrompt = *style.Style
			t.RawSetString("__styleFilterPrompt", st)
		})

	t.RawSetString("__styleFilterCursor", golua.LNil)
	lib.TableFunction(state, t, "style_filter_cursor",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFilterCursor")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.FilterCursor
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFilterCursor", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_filter_cursor_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.FilterCursor = *style.Style
			t.RawSetString("__styleFilterCursor", st)
		})

	t.RawSetString("__styleFilterCharMatch", golua.LNil)
	lib.TableFunction(state, t, "style_filter_char_match",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFilterCursor")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.DefaultFilterCharacterMatch
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFilterCharMatch", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_filter_char_match_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.DefaultFilterCharacterMatch = *style.Style
			t.RawSetString("__styleFilterCharMatch", st)
		})

	t.RawSetString("__styleStatusbar", golua.LNil)
	lib.TableFunction(state, t, "style_statusbar",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleStatusbar")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.StatusBar
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleStatusbar", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_statusbar_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.StatusBar = *style.Style
			t.RawSetString("__styleStatusbar", st)
		})

	t.RawSetString("__styleStatusEmpty", golua.LNil)
	lib.TableFunction(state, t, "style_status_empty",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleStatusEmpty")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.StatusEmpty
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleStatusEmpty", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_status_empty_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.StatusEmpty = *style.Style
			t.RawSetString("__styleStatusEmpty", st)
		})

	t.RawSetString("__styleStatusbarFilterActive", golua.LNil)
	lib.TableFunction(state, t, "style_statusbar_filter_active",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleStatusbarFilterActive")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.StatusBarActiveFilter
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleStatusbarFilterActive", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_statusbar_filter_active_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.StatusBarActiveFilter = *style.Style
			t.RawSetString("__styleStatusbarFilterActive", st)
		})

	t.RawSetString("__styleStatusbarFilterCount", golua.LNil)
	lib.TableFunction(state, t, "style_statusbar_filter_count",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleStatusbarFilterCount")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.StatusBarFilterCount
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleStatusbarFilterCount", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_statusbar_filter_count_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.StatusBarFilterCount = *style.Style
			t.RawSetString("__styleStatusbarFilterCount", st)
		})

	t.RawSetString("__styleNoItems", golua.LNil)
	lib.TableFunction(state, t, "style_no_items",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleNoItems")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.NoItems
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleNoItems", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_no_items_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.NoItems = *style.Style
			t.RawSetString("__styleNoItems", st)
		})

	t.RawSetString("__styleHelp", golua.LNil)
	lib.TableFunction(state, t, "style_help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleHelp")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.HelpStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleHelp", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_help_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.HelpStyle = *style.Style
			t.RawSetString("__styleHelp", st)
		})

	t.RawSetString("__stylePagination", golua.LNil)
	lib.TableFunction(state, t, "style_pagination",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePagination")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.PaginationStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePagination", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_pagination_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.PaginationStyle = *style.Style
			t.RawSetString("__stylePagination", st)
		})

	t.RawSetString("__stylePaginationDotActive", golua.LNil)
	lib.TableFunction(state, t, "style_pagination_dot_active",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePaginationDotActive")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.ActivePaginationDot
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePaginationDotActive", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_pagination_dot_active_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.ActivePaginationDot = *style.Style
			t.RawSetString("__stylePaginationDotActive", st)
		})

	t.RawSetString("__stylePaginationDotInactive", golua.LNil)
	lib.TableFunction(state, t, "style_pagination_dot_inactive",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePaginationDotInactive")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.InactivePaginationDot
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePaginationDotInactive", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_pagination_dot_inactive_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.InactivePaginationDot = *style.Style
			t.RawSetString("__stylePaginationDotInactive", st)
		})

	t.RawSetString("__styleDividerDot", golua.LNil)
	lib.TableFunction(state, t, "style_divider_dot",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDividerDot")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Lists[id].Styles.DividerDot
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDividerDot", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_divider_dot_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.DividerDot = *style.Style
			t.RawSetString("__styleDividerDot", st)
		})

	return t
}

func listKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [14]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("cursor_up", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("cursor_down", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("page_next", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("page_prev", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("goto_start", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("goto_end", tuikeyTable(r, lib, state, program, ids[5]))
	t.RawSetString("filter", tuikeyTable(r, lib, state, program, ids[6]))
	t.RawSetString("filter_clear", tuikeyTable(r, lib, state, program, ids[7]))
	t.RawSetString("filter_cancel", tuikeyTable(r, lib, state, program, ids[8]))
	t.RawSetString("filter_accept", tuikeyTable(r, lib, state, program, ids[9]))
	t.RawSetString("show_full_help", tuikeyTable(r, lib, state, program, ids[10]))
	t.RawSetString("close_full_help", tuikeyTable(r, lib, state, program, ids[11]))
	t.RawSetString("quit", tuikeyTable(r, lib, state, program, ids[12]))
	t.RawSetString("force_quit", tuikeyTable(r, lib, state, program, ids[13]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			li := item.Lists[id]
			li.KeyMap = list.DefaultKeyMap()
			item.KeyBindings[ids[0]] = &li.KeyMap.CursorUp
			item.KeyBindings[ids[1]] = &li.KeyMap.CursorDown
			item.KeyBindings[ids[2]] = &li.KeyMap.NextPage
			item.KeyBindings[ids[3]] = &li.KeyMap.PrevPage
			item.KeyBindings[ids[4]] = &li.KeyMap.GoToStart
			item.KeyBindings[ids[5]] = &li.KeyMap.GoToEnd
			item.KeyBindings[ids[6]] = &li.KeyMap.Filter
			item.KeyBindings[ids[7]] = &li.KeyMap.ClearFilter
			item.KeyBindings[ids[8]] = &li.KeyMap.CancelWhileFiltering
			item.KeyBindings[ids[9]] = &li.KeyMap.AcceptWhileFiltering
			item.KeyBindings[ids[10]] = &li.KeyMap.ShowFullHelp
			item.KeyBindings[ids[11]] = &li.KeyMap.CloseFullHelp
			item.KeyBindings[ids[12]] = &li.KeyMap.Quit
			item.KeyBindings[ids[13]] = &li.KeyMap.ForceQuit
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()
			kt.RawSetInt(1, t.RawGetString("cursor_up"))
			kt.RawSetInt(2, t.RawGetString("cursor_down"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("cursor_up"))
			kt1.RawSetInt(2, t.RawGetString("cursor_down"))
			kt1.RawSetInt(3, t.RawGetString("page_next"))
			kt1.RawSetInt(4, t.RawGetString("page_prev"))
			kt1.RawSetInt(5, t.RawGetString("goto_start"))
			kt1.RawSetInt(6, t.RawGetString("goto_end"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("filter"))
			kt2.RawSetInt(2, t.RawGetString("filter_clear"))
			kt2.RawSetInt(3, t.RawGetString("filter_accept"))
			kt2.RawSetInt(4, t.RawGetString("filter_cancel"))

			kt3 := state.NewTable()
			kt3.RawSetInt(1, t.RawGetString("quit"))
			kt3.RawSetInt(2, t.RawGetString("force_quit"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)
			kt.RawSetInt(3, kt3)

			state.Push(kt)
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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
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

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.Paginators[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.PrevPage,
				&value.NextPage,
			)

			ids := [2]int{}
			for i := range 2 {
				ids[i] = start + i
			}

			kmt := paginatorKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	return t
}

func paginatorKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [2]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("page_prev", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("page_next", tuikeyTable(r, lib, state, program, ids[1]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			pg := item.Paginators[id]
			pg.KeyMap = paginator.DefaultKeyMap
			item.KeyBindings[ids[0]] = &pg.KeyMap.PrevPage
			item.KeyBindings[ids[1]] = &pg.KeyMap.NextPage
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()
			kt.RawSetInt(1, t.RawGetString("page_prev"))
			kt.RawSetInt(2, t.RawGetString("page_next"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("page_prev"))
			kt1.RawSetInt(2, t.RawGetString("page_next"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)

			state.Push(kt)
			return 1
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
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
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

			state.Push(customtea.CMDProgressSet(state, id, percent))
			return 1
		})

	lib.TableFunction(state, t, "percent_dec",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			percent := args["percent"].(float64)

			state.Push(customtea.CMDProgressDec(state, id, percent))
			return 1
		})

	lib.TableFunction(state, t, "percent_inc",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			percent := args["percent"].(float64)

			state.Push(customtea.CMDProgressInc(state, id, percent))
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

	t.RawSetString("__stylePercent", golua.LNil)
	lib.TableFunction(state, t, "style_percentage",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__stylePercent")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ProgressBars[id].PercentageStyle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__stylePercent", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_percentage_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ProgressBars[id].PercentageStyle = *style.Style
			t.RawSetString("__stylePercent", st)
		})

	return t
}

func stopwatchTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		sw, cmd := item.StopWatches[id].Update(*item.Msg)
		item.StopWatches[id] = &sw

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "start",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDStopWatchStart(state, id))
			return 1
		})

	lib.TableFunction(state, t, "stop",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDStopWatchStop(state, id))
			return 1
		})

	lib.TableFunction(state, t, "toggle",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDStopWatchToggle(state, id))
			return 1
		})

	lib.TableFunction(state, t, "reset",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDStopWatchReset(state, id))
			return 1
		})

	lib.TableFunction(state, t, "elapsed",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StopWatches[id].Elapsed()

			state.Push(golua.LNumber(value.Milliseconds()))
			return 1
		})

	lib.TableFunction(state, t, "running",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StopWatches[id].Running()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "interval",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StopWatches[id].Interval

			state.Push(golua.LNumber(value.Milliseconds()))
			return 1
		})

	lib.BuilderFunction(state, t, "interval_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "interval"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StopWatches[id].Interval = time.Duration(args["interval"].(int) * 1e6)
		})

	return t
}

func timerTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Timers[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		ti, cmd := item.Timers[id].Update(*item.Msg)
		item.Timers[id] = &ti

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "init",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDTimerInit(state, id))
			return 1
		})

	lib.TableFunction(state, t, "start",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDTimerStart(state, id))
			return 1
		})

	lib.TableFunction(state, t, "stop",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDTimerStop(state, id))
			return 1
		})

	lib.TableFunction(state, t, "toggle",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))

			state.Push(customtea.CMDTimerToggle(state, id))
			return 1
		})

	lib.TableFunction(state, t, "running",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Timers[id].Running()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "timed_out",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Timers[id].Timedout()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "timeout",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Timers[id].Timeout

			state.Push(golua.LNumber(value.Milliseconds()))
			return 1
		})

	lib.BuilderFunction(state, t, "timeout_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "timeout"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Timers[id].Timeout = time.Duration(args["timeout"].(int) * 1e6)
		})

	lib.TableFunction(state, t, "interval",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Timers[id].Interval

			state.Push(golua.LNumber(value.Milliseconds()))
			return 1
		})

	lib.BuilderFunction(state, t, "interval_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "interval"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Timers[id].Interval = time.Duration(args["interval"].(int) * 1e6)
		})

	return t
}

func tableOptionsTable(lib *lua.Lib, state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("__columns", golua.LNil)
	t.RawSetString("__rows", golua.LNil)
	t.RawSetString("__focused", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__styleHeader", golua.LNil)
	t.RawSetString("__styleCell", golua.LNil)
	t.RawSetString("__styleSelected", golua.LNil)

	lib.BuilderFunction(state, t, "focused",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "focused"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__focused", golua.LBool(args["focused"].(bool)))
		})

	lib.BuilderFunction(state, t, "width",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__width", golua.LNumber(args["width"].(int)))
		})

	lib.BuilderFunction(state, t, "height",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__height", golua.LNumber(args["height"].(int)))
		})

	lib.BuilderFunction(state, t, "columns",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "cols"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__columns", args["cols"].(*golua.LTable))
		})

	lib.BuilderFunction(state, t, "rows",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "rows"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__rows", args["rows"].(*golua.LTable))
		})

	lib.BuilderFunction(state, t, "styles",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "header"},
			{Type: lua.RAW_TABLE, Name: "cell"},
			{Type: lua.RAW_TABLE, Name: "selected"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__styleHeader", args["header"].(*golua.LTable))
			t.RawSetString("__styleCell", args["cell"].(*golua.LTable))
			t.RawSetString("__styleSelected", args["selected"].(*golua.LTable))
		})

	return t
}

func tuitableColTable(state *golua.LState, title string, width int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("title", golua.LString(title))
	t.RawSetString("width", golua.LNumber(width))

	return t
}

func tableOptionsBuild(t *golua.LTable, r *lua.Runner) []table.Option {
	opts := []table.Option{}

	focused := t.RawGetString("__focused")
	if focused.Type() == golua.LTBool {
		opts = append(opts, table.WithFocused(bool(focused.(golua.LBool))))
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		opts = append(opts, table.WithWidth(int(width.(golua.LNumber))))
	}

	height := t.RawGetString("__height")
	if height.Type() == golua.LTNumber {
		opts = append(opts, table.WithHeight(int(height.(golua.LNumber))))
	}

	cols := t.RawGetString("__columns")
	if cols.Type() == golua.LTTable {
		colt := cols.(*golua.LTable)
		colList := make([]table.Column, colt.Len())

		for i := range colt.Len() {
			c := colt.RawGetInt(i + 1).(*golua.LTable)
			colList[i] = table.Column{
				Title: string(c.RawGetString("title").(golua.LString)),
				Width: int(c.RawGetString("width").(golua.LNumber)),
			}
		}

		opts = append(opts, table.WithColumns(colList))
	}

	rows := t.RawGetString("__rows")
	if rows.Type() == golua.LTTable {
		rowt := rows.(*golua.LTable)
		rowList := make([]table.Row, rowt.Len())

		for i := range rowt.Len() {
			r := rowt.RawGetInt(i + 1).(*golua.LTable)
			rowData := make(table.Row, r.Len())
			for z := range r.Len() {
				rowData[z] = string(r.RawGetInt(z + 1).(golua.LString))
			}
			rowList[i] = rowData
		}

		opts = append(opts, table.WithRows(rowList))
	}

	header := t.RawGetString("__styleHeader")
	cell := t.RawGetString("__styleCell")
	selected := t.RawGetString("__styleSelected")
	if header.Type() == golua.LTTable && cell.Type() == golua.LTTable && cell.Type() == golua.LTTable {
		hid := header.(*golua.LTable).RawGetString("id").(golua.LNumber)
		cid := cell.(*golua.LTable).RawGetString("id").(golua.LNumber)
		sid := selected.(*golua.LTable).RawGetString("id").(golua.LNumber)

		header, _ := r.CR_LIP.Item(int(hid))
		cell, _ := r.CR_LIP.Item(int(cid))
		selected, _ := r.CR_LIP.Item(int(sid))

		opts = append(opts, table.WithStyles(table.Styles{
			Header:   *header.Style,
			Cell:     *cell.Style,
			Selected: *selected.Style,
		}))
	}

	return opts
}

func tuitableTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Tables[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		ti, cmd := item.Tables[id].Update(*item.Msg)
		item.Tables[id] = &ti

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.BuilderFunction(state, t, "update_viewport",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].UpdateViewport()
		})

	lib.TableFunction(state, t, "focused",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Tables[id].Focused()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "focus",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].Focus()
		})

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].Blur()
		})

	lib.BuilderFunction(state, t, "goto_top",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].GotoTop()
		})

	lib.BuilderFunction(state, t, "goto_bottom",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].GotoBottom()
		})

	lib.BuilderFunction(state, t, "move_up",
		[]lua.Arg{
			{Type: lua.INT, Name: "n"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].MoveUp(args["n"].(int))
		})

	lib.BuilderFunction(state, t, "move_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "n"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].MoveDown(args["n"].(int))
		})

	lib.TableFunction(state, t, "cursor",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Tables[id].Cursor()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "cursor_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "n"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].SetCursor(args["n"].(int))
		})

	lib.TableFunction(state, t, "columns",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			cols := item.Tables[id].Columns()
			colt := state.NewTable()
			for i, v := range cols {
				colt.RawSetInt(i+1, tuitableColTable(state, v.Title, v.Width))
			}

			state.Push(colt)
			return 1
		})

	lib.TableFunction(state, t, "rows",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			rows := item.Tables[id].Rows()
			rowt := state.NewTable()
			for i, v := range rows {
				r := state.NewTable()
				for z, s := range v {
					r.RawSetInt(z+1, golua.LString(s))
				}
				rowt.RawSetInt(i+1, r)
			}

			state.Push(rowt)
			return 1
		})

	lib.BuilderFunction(state, t, "columns_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "cols"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			cols := args["cols"].(*golua.LTable)
			colList := make([]table.Column, cols.Len())
			for i := range cols.Len() {
				c := cols.RawGetInt(i + 1).(*golua.LTable)
				colList[i] = table.Column{
					Title: string(c.RawGetString("title").(golua.LString)),
					Width: int(c.RawGetString("width").(golua.LNumber)),
				}
			}

			item.Tables[id].SetColumns(colList)
		})

	lib.BuilderFunction(state, t, "rows_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "rows"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			rows := args["rows"].(*golua.LTable)
			rowsList := make([]table.Row, rows.Len())

			for i := range rows.Len() {
				r := rows.RawGetInt(i + 1).(*golua.LTable)
				row := make([]string, r.Len())

				for z := range r.Len() {
					row[z] = string(r.RawGetInt(z + 1).(golua.LString))
				}
				rowsList[i] = row
			}

			item.Tables[id].SetRows(rowsList)
		})

	lib.BuilderFunction(state, t, "from_values",
		[]lua.Arg{
			{Type: lua.STRING, Name: "value"},
			{Type: lua.STRING, Name: "separator"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].FromValues(args["value"].(string), args["separator"].(string))
		})

	lib.TableFunction(state, t, "row_selected",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Tables[id].SelectedRow()
			rows := state.NewTable()

			for i, s := range value {
				rows.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(rows)
			return 1
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

			value := item.Tables[id].Width()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.TableFunction(state, t, "height",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Tables[id].Height()

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

			item.Tables[id].SetWidth(args["width"].(int))
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

			item.Tables[id].SetHeight(args["height"].(int))
		})

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.Tables[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.LineUp,
				&value.LineDown,
				&value.PageUp,
				&value.PageDown,
				&value.HalfPageUp,
				&value.HalfPageDown,
				&value.GotoTop,
				&value.GotoBottom,
			)

			ids := [8]int{}
			for i := range 8 {
				ids[i] = start + i
			}

			kmt := tableKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	t.RawSetString("__help", golua.LNil)
	lib.TableFunction(state, t, "help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			oh := t.RawGetString("__help")
			if oh.Type() == golua.LTTable {
				state.Push(oh)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Tables[id].Help
			mid := len(item.Helps)
			item.Helps = append(item.Helps, model)

			hp := helpTable(r, lib, state, program, mid)
			state.Push(hp)
			t.RawSetString("__help", hp)
			return 1
		})

	lib.TableFunction(state, t, "help_view",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			str := item.Tables[id].HelpView()

			state.Push(golua.LString(str))
			return 1
		})

	lib.BuilderFunction(state, t, "styles",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "header"},
			{Type: lua.RAW_TABLE, Name: "cell"},
			{Type: lua.RAW_TABLE, Name: "selected"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			hid := args["header"].(*golua.LTable).RawGetString("id").(golua.LNumber)
			cid := args["cell"].(*golua.LTable).RawGetString("id").(golua.LNumber)
			sid := args["selected"].(*golua.LTable).RawGetString("id").(golua.LNumber)

			header, _ := r.CR_LIP.Item(int(hid))
			cell, _ := r.CR_LIP.Item(int(cid))
			selected, _ := r.CR_LIP.Item(int(sid))

			item.Tables[id].SetStyles(table.Styles{
				Header:   *header.Style,
				Cell:     *cell.Style,
				Selected: *selected.Style,
			})
		})

	lib.BuilderFunction(state, t, "styles_default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].SetStyles(table.DefaultStyles())
		})

	return t
}

func tableKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [8]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("line_up", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("line_down", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("page_up", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("page_down", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("half_page_up", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("half_page_down", tuikeyTable(r, lib, state, program, ids[5]))
	t.RawSetString("goto_top", tuikeyTable(r, lib, state, program, ids[6]))
	t.RawSetString("goto_bottom", tuikeyTable(r, lib, state, program, ids[7]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			tb := item.Tables[id]
			tb.KeyMap = table.DefaultKeyMap()
			item.KeyBindings[ids[0]] = &tb.KeyMap.LineUp
			item.KeyBindings[ids[1]] = &tb.KeyMap.LineDown
			item.KeyBindings[ids[2]] = &tb.KeyMap.PageUp
			item.KeyBindings[ids[3]] = &tb.KeyMap.PageDown
			item.KeyBindings[ids[4]] = &tb.KeyMap.HalfPageUp
			item.KeyBindings[ids[5]] = &tb.KeyMap.HalfPageDown
			item.KeyBindings[ids[6]] = &tb.KeyMap.GotoTop
			item.KeyBindings[ids[7]] = &tb.KeyMap.GotoBottom
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()

			kt.RawSetInt(1, t.RawGetString("line_up"))
			kt.RawSetInt(2, t.RawGetString("line_up"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("line_up"))
			kt1.RawSetInt(2, t.RawGetString("line_down"))
			kt1.RawSetInt(3, t.RawGetString("goto_top"))
			kt1.RawSetInt(4, t.RawGetString("goto_bottom"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("page_up"))
			kt2.RawSetInt(2, t.RawGetString("page_down"))
			kt2.RawSetInt(3, t.RawGetString("page_up_half"))
			kt2.RawSetInt(4, t.RawGetString("page_down_half"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)
			state.Push(kt)
			return 1
		})

	return t
}

func viewportTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Viewports[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		vp, cmd := item.Viewports[id].Update(*item.Msg)
		item.Viewports[id] = &vp

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "view_up",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].ViewUp()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "view_down",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].ViewDown()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "view_up_half",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].HalfViewUp()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "view_down_half",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].HalfViewDown()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "at_top",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].AtTop()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "at_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].AtBottom()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "goto_top",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].GotoTop()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "goto_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].GotoBottom()
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "line_up",
		[]lua.Arg{
			{Type: lua.INT, Name: "n"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].LineUp(args["n"].(int))
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "line_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "n"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].LineDown(args["n"].(int))
			lines := state.NewTable()

			for i, s := range value {
				lines.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(lines)
			return 1
		})

	lib.TableFunction(state, t, "past_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].PastBottom()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.TableFunction(state, t, "scroll_percent",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].ScrollPercent()

			state.Push(golua.LNumber(value))
			return 1
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

			value := item.Viewports[id].Width

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.TableFunction(state, t, "height",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].Height

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

			item.Viewports[id].Width = args["width"].(int)
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

			item.Viewports[id].Height = args["height"].(int)
		})

	lib.BuilderFunction(state, t, "content_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].SetContent(args["content"].(string))
		})

	lib.TableFunction(state, t, "line_count_total",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].TotalLineCount()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.TableFunction(state, t, "line_count_visible",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].VisibleLineCount()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.TableFunction(state, t, "mouse_wheel_enabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].MouseWheelEnabled

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "mouse_wheel_enabled_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].MouseWheelEnabled = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "mouse_wheel_delta",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].MouseWheelDelta

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "mouse_wheel_delta_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "delta"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].MouseWheelDelta = args["delta"].(int)
		})

	lib.TableFunction(state, t, "offset_y",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].YOffset

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "offset_y_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "offset"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].SetYOffset(args["offset"].(int))
		})

	lib.BuilderFunction(state, t, "offset_y_set_direct",
		[]lua.Arg{
			{Type: lua.INT, Name: "offset"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].YOffset = args["offset"].(int)
		})

	lib.TableFunction(state, t, "position_y",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].YPosition

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "position_y_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "position"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].YPosition = args["position"].(int)
		})

	lib.TableFunction(state, t, "high_performance",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Viewports[id].HighPerformanceRendering

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "high_performance_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Viewports[id].HighPerformanceRendering = args["enabled"].(bool)
		})

	t.RawSetString("__keymap", golua.LNil)
	lib.TableFunction(state, t, "keymap",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kmto := t.RawGetString("__keymap")
			if kmto.Type() == golua.LTTable {
				state.Push(kmto)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := &item.Viewports[id].KeyMap
			start := len(item.KeyBindings)
			item.KeyBindings = append(item.KeyBindings,
				&value.PageDown,
				&value.PageUp,
				&value.HalfPageUp,
				&value.HalfPageDown,
				&value.Down,
				&value.Up,
			)

			ids := [6]int{}
			for i := range 6 {
				ids[i] = start + i
			}

			kmt := viewportKeymapTable(r, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	t.RawSetString("__style", golua.LNil)
	lib.TableFunction(state, t, "style",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__style")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Viewports[id].Style
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__style", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Viewports[id].Style = *style.Style
			t.RawSetString("__style", st)
		})

	return t
}

func viewportKeymapTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int, ids [6]int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("page_down", tuikeyTable(r, lib, state, program, ids[0]))
	t.RawSetString("page_up", tuikeyTable(r, lib, state, program, ids[1]))
	t.RawSetString("page_up_half", tuikeyTable(r, lib, state, program, ids[2]))
	t.RawSetString("page_down_half", tuikeyTable(r, lib, state, program, ids[3]))
	t.RawSetString("down", tuikeyTable(r, lib, state, program, ids[4]))
	t.RawSetString("up", tuikeyTable(r, lib, state, program, ids[5]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			vp := item.Viewports[id]
			vp.KeyMap = viewport.DefaultKeyMap()
			item.KeyBindings[ids[0]] = &vp.KeyMap.PageDown
			item.KeyBindings[ids[1]] = &vp.KeyMap.PageUp
			item.KeyBindings[ids[2]] = &vp.KeyMap.HalfPageUp
			item.KeyBindings[ids[3]] = &vp.KeyMap.HalfPageDown
			item.KeyBindings[ids[4]] = &vp.KeyMap.Down
			item.KeyBindings[ids[5]] = &vp.KeyMap.Up
		})

	lib.TableFunction(state, t, "help_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt := state.NewTable()

			kt.RawSetInt(1, t.RawGetString("up"))
			kt.RawSetInt(2, t.RawGetString("down"))

			state.Push(kt)
			return 1
		})

	lib.TableFunction(state, t, "help_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			kt1 := state.NewTable()
			kt1.RawSetInt(1, t.RawGetString("up"))
			kt1.RawSetInt(2, t.RawGetString("down"))

			kt2 := state.NewTable()
			kt2.RawSetInt(1, t.RawGetString("page_up"))
			kt2.RawSetInt(2, t.RawGetString("page_down"))

			kt3 := state.NewTable()
			kt3.RawSetInt(1, t.RawGetString("page_up_half"))
			kt3.RawSetInt(2, t.RawGetString("page_down_half"))

			kt := state.NewTable()
			kt.RawSetInt(1, kt1)
			kt.RawSetInt(2, kt2)
			kt.RawSetInt(3, kt3)
			state.Push(kt)
			return 1
		})

	return t
}

func tuicustomTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("init", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		cmd := item.Customs[id].Init()

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, err.Error())
		}

		str := item.Customs[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	lib.TableFunction(state, t, "update",
		[]lua.Arg{
			lua.ArgVariadic("values", lua.ArrayType{Type: lua.ANY}, true),
		},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			values := args["values"].([]any)
			msg := *item.Msg

			if len(values) > 0 {
				msg = teamodels.CustomMSG{
					Original: msg,
					Values:   values,
				}
			}

			cuv, cmd := item.Customs[id].Update(msg)
			cu := cuv.(teamodels.CustomModel)
			item.Customs[id] = &cu

			var bcmd *golua.LTable

			if cmd == nil {
				bcmd = customtea.CMDNone(state)
			} else {
				bcmd = customtea.CMDStored(state, item, cmd)
			}

			state.Push(bcmd)
			return 1
		})

	return t
}

func keyOptionsTable(state *golua.LState, lib *lua.Lib) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("__disabled", golua.LNil)
	t.RawSetString("__helpKey", golua.LNil)
	t.RawSetString("__helpDesc", golua.LNil)
	t.RawSetString("__keys", golua.LNil)

	lib.BuilderFunction(state, t, "disabled",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__disabled", golua.LBool(args["enabled"].(bool)))
		})

	lib.BuilderFunction(state, t, "help",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			{Type: lua.STRING, Name: "desc"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__helpKey", golua.LString(args["key"].(string)))
			t.RawSetString("__helpDesc", golua.LString(args["desc"].(string)))
		})

	lib.BuilderFunction(state, t, "keys",
		[]lua.Arg{
			lua.ArgVariadic("keys", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			keyList := args["keys"].([]any)

			keys := state.NewTable()
			for i, s := range keyList {
				keys.RawSetInt(i+1, golua.LString(s.(string)))
			}

			t.RawSetString("__keys", keys)
		})

	return t
}

func keyOptionsBuild(t *golua.LTable) []key.BindingOpt {
	opts := []key.BindingOpt{}

	disabled := t.RawGetString("__disabled")
	if disabled.Type() == golua.LTBool {
		d := bool(disabled.(golua.LBool))
		if d {
			opts = append(opts, key.WithDisabled())
		}
	}

	helpKey := t.RawGetString("__helpKey")
	helpDesc := t.RawGetString("__helpDesc")
	if helpKey.Type() == golua.LTString && helpDesc.Type() == golua.LTString {
		opts = append(opts, key.WithHelp(
			string(helpKey.(golua.LString)),
			string(helpDesc.(golua.LString)),
		))
	}

	keys := t.RawGetString("__keys")
	if keys.Type() == golua.LTTable {
		kt := keys.(*golua.LTable)
		keyList := make([]string, kt.Len())

		for i := range kt.Len() {
			ki := string(kt.RawGetInt(i + 1).(golua.LString))
			keyList[i] = ki
		}

		opts = append(opts, key.WithKeys(keyList...))
	}

	return opts
}

func tuikeyTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "enabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.KeyBindings[id].Enabled()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "enabled_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.KeyBindings[id].SetEnabled(args["enabled"].(bool))
		})

	lib.TableFunction(state, t, "help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.KeyBindings[id].Help()

			state.Push(golua.LString(value.Key))
			state.Push(golua.LString(value.Desc))
			return 2
		})

	lib.BuilderFunction(state, t, "help_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			{Type: lua.STRING, Name: "desc"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.KeyBindings[id].SetHelp(args["key"].(string), args["desc"].(string))
		})

	lib.TableFunction(state, t, "keys",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.KeyBindings[id].Keys()
			keys := state.NewTable()

			for i, s := range value {
				keys.RawSetInt(i+1, golua.LString(s))
			}

			state.Push(keys)
			return 1
		})

	lib.BuilderFunction(state, t, "keys_set",
		[]lua.Arg{
			lua.ArgVariadic("keys", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			keys := args["keys"].([]any)
			keyList := make([]string, len(keys))

			for i, s := range keys {
				keyList[i] = s.(string)
			}

			item.KeyBindings[id].SetKeys(keyList...)
		})

	lib.BuilderFunction(state, t, "unbind",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.KeyBindings[id].Unbind()
		})

	return t
}

func helpTable(r *lua.Runner, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "view",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "keymap"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			showAll := item.Helps[id].ShowAll
			value := ""
			keymap := args["keymap"].(*golua.LTable)

			if showAll {
				full := keymap.RawGetString("help_full")
				if full.Type() == golua.LTFunction {
					state.Push(full)
					state.Call(0, 1)
					groups := state.CheckTable(-1)
					state.Pop(1)

					state.Push(t.RawGetString("view_help_full"))
					state.Push(groups)
					state.Call(1, 1)
					value = state.CheckString(-1)
					state.Pop(1)
				}
			} else {
				short := keymap.RawGetString("help_short")
				if short.Type() == golua.LTFunction {
					state.Push(short)
					state.Call(0, 1)
					bindings := state.CheckTable(-1)
					state.Pop(1)

					state.Push(t.RawGetString("view_help_short"))
					state.Push(bindings)
					state.Call(1, 1)
					value = state.CheckString(-1)
					state.Pop(1)
				}
			}

			state.Push(golua.LString(value))
			return 1
		})

	lib.TableFunction(state, t, "view_help_short",
		[]lua.Arg{
			lua.ArgArray("bindings", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			bindings := args["bindings"].([]any)
			bindList := make([]key.Binding, len(bindings))
			for i, v := range bindings {
				b := v.(*golua.LTable)
				bid := b.RawGetString("id").(golua.LNumber)
				bindList[i] = *item.KeyBindings[int(bid)]
			}

			value := item.Helps[id].ShortHelpView(bindList)
			state.Push(golua.LString(value))
			return 1
		})

	lib.TableFunction(state, t, "view_help_full",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "groups"},
		},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			groups := args["groups"].(*golua.LTable)
			groupList := make([][]key.Binding, groups.Len())
			for i := range groups.Len() {
				group := groups.RawGetInt(i + 1).(*golua.LTable)
				groupList[i] = make([]key.Binding, group.Len())

				for z := range group.Len() {
					b := group.RawGetInt(z + 1).(*golua.LTable)
					bid := b.RawGetString("id").(golua.LNumber)
					groupList[i][z] = *item.KeyBindings[int(bid)]
				}
			}

			value := item.Helps[id].FullHelpView(groupList)
			state.Push(golua.LString(value))
			return 1
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

			value := item.Helps[id].Width

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

			item.Helps[id].Width = args["width"].(int)
		})

	lib.TableFunction(state, t, "show_all",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Helps[id].ShowAll

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "show_all_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Helps[id].ShowAll = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "separator_short",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Helps[id].ShortSeparator

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "separator_short_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "sep"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Helps[id].ShortSeparator = args["sep"].(string)
		})

	lib.TableFunction(state, t, "separator_full",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Helps[id].FullSeparator

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "separator_full_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "sep"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Helps[id].FullSeparator = args["sep"].(string)
		})

	lib.TableFunction(state, t, "ellipsis",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Helps[id].Ellipsis

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "ellipsis_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "ellipsis"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Helps[id].Ellipsis = args["ellipsis"].(string)
		})

	t.RawSetString("__styleEllipsis", golua.LNil)
	lib.TableFunction(state, t, "style_ellipsis",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleEllipsis")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.Ellipsis
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleEllipsis", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_ellipsis_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.Ellipsis = *style.Style
			t.RawSetString("__styleEllipsis", st)
		})

	t.RawSetString("__styleShortKey", golua.LNil)
	lib.TableFunction(state, t, "style_short_key",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleShortKey")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.ShortKey
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleShortKey", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_short_key_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.ShortKey = *style.Style
			t.RawSetString("__styleShortKey", st)
		})

	t.RawSetString("__styleShortDesc", golua.LNil)
	lib.TableFunction(state, t, "style_short_desc",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleShortDesc")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.ShortDesc
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleShortDesc", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_short_desc_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.ShortDesc = *style.Style
			t.RawSetString("__styleShortDesc", st)
		})

	t.RawSetString("__styleShortSep", golua.LNil)
	lib.TableFunction(state, t, "style_short_separator",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleShortSep")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.ShortSeparator
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleShortSep", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_short_separator_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.ShortSeparator = *style.Style
			t.RawSetString("__styleShortSep", st)
		})

	t.RawSetString("__styleFullKey", golua.LNil)
	lib.TableFunction(state, t, "style_full_key",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFullKey")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.FullKey
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFullKey", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_full_key_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.FullKey = *style.Style
			t.RawSetString("__styleFullKey", st)
		})

	t.RawSetString("__styleFullDesc", golua.LNil)
	lib.TableFunction(state, t, "style_full_desc",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFullDesc")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.FullDesc
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFullDesc", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_full_desc_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.FullDesc = *style.Style
			t.RawSetString("__styleFullDesc", st)
		})

	t.RawSetString("__styleFullSep", golua.LNil)
	lib.TableFunction(state, t, "style_full_separator",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFullSep")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.Helps[id].Styles.FullSeparator
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFullSep", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_full_separator_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, err.Error())
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Helps[id].Styles.FullSeparator = *style.Style
			t.RawSetString("__styleFullSep", st)
		})

	return t
}
