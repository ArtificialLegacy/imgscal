package lib

import (
	"errors"
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
	"github.com/charmbracelet/lipgloss"
	teaimage "github.com/mistakenelf/teacup/image"
	"github.com/mistakenelf/teacup/statusbar"
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
			id := r.CR_TEA.Add(&teamodels.TeaItem{})
			t := teaTable(r, lg, state, lib, id)

			state.Push(t)
			return 1
		},
	)

	/// @func program_options() -> struct<tui.ProgramOptions>
	/// @returns {struct<tui.ProgramOptions>}
	lib.CreateFunction(tab, "program_options",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := programOptions(lib, state)

			state.Push(t)
			return 1
		},
	)

	/// @func run(program, opts?)
	/// @arg program {struct<tui.Program>}
	/// @arg? opts {struct<tui.ProgramOptions>}
	lib.CreateFunction(tab, "run",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "program"},
			{Type: lua.RAW_TABLE, Name: "opts", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			program := args["program"].(*golua.LTable)
			id := int(program.RawGetString("id").(golua.LNumber))
			item, err := r.CR_TEA.Item(id)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			optArg := args["opts"].(*golua.LTable)
			opts := programOptionsBuild(state, optArg)

			pstate, _ := state.NewThread()
			p := tea.NewProgram(customtea.ProgramModel{Id: id, Item: item, State: pstate, R: r, Lg: lg}, opts...)
			_, err = p.Run()
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			spin := spinner.New(spinner.WithSpinner(spinnerList[args["type"].(int)]))
			id := spin.ID()
			item.Spinners[id] = &spin

			t := spinnerTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			t := spinnerTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			ta := textarea.New()
			id := len(item.TextAreas)
			item.TextAreas = append(item.TextAreas, &ta)

			t := textareaTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			ti := textinput.New()
			id := len(item.TextInputs)
			item.TextInputs = append(item.TextInputs, &ti)

			t := textinputTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			cu := cursor.New()
			id := len(item.Cursors)
			item.Cursors = append(item.Cursors, &cu)

			t := cursorTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			fp := filepicker.New()
			id := len(item.FilePickers)
			item.FilePickers = append(item.FilePickers, &fp)

			t := filePickerTable(r, lg, lib, state, prgrm, id)

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

	/// @func list(id, items, width, height, delegate) -> struct<tui.List>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the list to.
	/// @arg items {[]struct<tui.ListItem>} - Array of list items.
	/// @arg width {int}
	/// @arg height {int}
	/// @arg delegate {struct<tui.ListDelegate>}
	/// @returns {struct<tui.List>}
	lib.CreateFunction(tab, "list",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("items", lua.ArrayType{Type: lua.RAW_TABLE}, false),
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "delegate"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			itemList := args["items"].([]any)
			items := make([]list.Item, len(itemList))

			for i, v := range itemList {
				items[i] = customtea.ListItemBuild(v.(*golua.LTable))
			}

			did := args["delegate"].(*golua.LTable).RawGetString("id").(golua.LNumber)
			delegate := item.ListDelegates[int(did)]

			li := list.New(items, delegate, args["width"].(int), args["height"].(int))
			id := len(item.Lists)
			item.Lists = append(item.Lists, &li)

			t := listTable(r, lg, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func list_delegate(id) -> struct<tui.ListDelegate>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the list to.
	/// @returns {struct<tui.ListDelegate>}
	lib.CreateFunction(tab, "list_delegate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			di := list.NewDefaultDelegate()
			id := len(item.ListDelegates)
			item.ListDelegates = append(item.ListDelegates, &di)

			t := listDelegateTable(r, lg, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func list_filter_rank(index, matched) -> struct<tui.ListFilterRank>
	/// @arg index {int}
	/// @arg matched {[]int}
	/// @returns {struct<tui.ListFilterRank>}
	lib.CreateFunction(tab, "list_filter_rank",
		[]lua.Arg{
			{Type: lua.INT, Name: "index"},
			{Type: lua.RAW_TABLE, Name: "matched"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct ListFilterRank
			/// @prop index {int} - The index of the item.
			/// @prop matched {[]int} - The indexes of the matched words.

			t := state.NewTable()

			t.RawSetString("index", golua.LNumber(args["index"].(int)))
			t.RawSetString("matched", args["matched"].(*golua.LTable))

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			t := paginatorTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			opts := progressOptionsBuild(args["options"].(*golua.LTable))

			pr := progress.New(opts...)
			id := len(item.ProgressBars)
			item.ProgressBars = append(item.ProgressBars, &pr)

			t := progressTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			t := stopwatchTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			t := timerTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			opts := tableOptionsBuild(args["options"].(*golua.LTable), r)

			tb := table.New(opts...)
			id := len(item.Tables)
			item.Tables = append(item.Tables, &tb)

			t := tuitableTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			width := args["width"].(int)
			height := args["height"].(int)
			id := len(item.Viewports)
			vp := viewport.New(width, height)
			item.Viewports = append(item.Viewports, &vp)

			t := viewportTable(r, lg, lib, state, prgrm, id)

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
			/// @struct CMDViewportSync
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop id {int} - The viewport id.

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
			/// @struct CMDViewportUp
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop id {int} - The viewport id.
			// @prop lines {[]string} - The lines to display.

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
			/// @struct CMDViewportDown
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop id {int} - The viewport id.
			// @prop lines {[]string} - The lines to display.

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			init := args["init"].(*golua.LFunction)
			update := args["update"].(*golua.LFunction)
			view := args["view"].(*golua.LFunction)

			id := len(item.Customs)
			cm := teamodels.NewCustomModel(prgrm, init, update, view, state, item, customtea.CMDBuild, customtea.BuildMSG)
			item.Customs = append(item.Customs, &cm)

			t := tuicustomTable(r, lg, lib, state, prgrm, id)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			options := args["option"].(*golua.LTable)
			opts := keyOptionsBuild(options)

			id := len(item.KeyBindings)
			ky := key.NewBinding(opts...)
			item.KeyBindings = append(item.KeyBindings, &ky)

			t := tuikeyTable(r, lg, lib, state, prgrm, id)

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
					lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			hp := help.New()
			id := len(item.Helps)
			item.Helps = append(item.Helps, &hp)

			t := helpTable(r, lg, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func image(id, active, borderless, borderColor?) -> struct<tui.Image>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the image to.
	/// @arg active {bool}
	/// @arg borderless {bool}
	/// @arg? borderColor {struct<lipgloss.ColorAdaptive>} - Defaults to white and black.
	/// @returns {struct<tui.Image>}
	lib.CreateFunction(tab, "image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "active"},
			{Type: lua.BOOL, Name: "borderless"},
			{Type: lua.RAW_TABLE, Name: "borderColor", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			active := args["active"].(bool)
			borderless := args["borderless"].(bool)
			color := args["borderColor"].(*golua.LTable)

			lcolor := lgColorGenericBuild(color)
			if _, ok := lcolor.(lipgloss.AdaptiveColor); !ok {
				lcolor = lipgloss.AdaptiveColor{
					Light: "#000000",
					Dark:  "#FFFFFF",
				}
			}

			im := teaimage.New(active, borderless, lcolor.(lipgloss.AdaptiveColor))
			id := len(item.Images)
			item.Images = append(item.Images, &im)

			t := tuiimageTable(r, lg, lib, state, prgrm, id)

			state.Push(t)
			return 1
		})

	/// @func image_to_string(img, width?) -> string
	/// @arg img {int<collection.IMAGE>}
	/// @arg? width {int} - Defaults to the image's width.
	/// @returns {string}
	/// @blocking
	lib.CreateFunction(tab, "image_to_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "width", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var result string

			<-r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					width := args["width"].(int)
					if width == 0 {
						width = i.Self.Image.Bounds().Dx()
					}
					result = teaimage.ToString(width, i.Self.Image)
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					result = ""
				},
			})

			state.Push(golua.LString(result))
			return 1
		})

	/// @func statusbar(id, first_foreground, first_background, second_foreground, second_background, third_foreground, third_background, fourth_foreground, fourth_background) -> struct<tui.StatusBar>
	/// @arg id {int<collection.CRATE_TEA>} - The program id to add the statusbar to.
	/// @arg first_foreground {struct<lipgloss.AdaptiveColor>}
	/// @arg first_background {struct<lipgloss.AdaptiveColor>}
	/// @arg second_foreground {struct<lipgloss.AdaptiveColor>}
	/// @arg second_background {struct<lipgloss.AdaptiveColor>}
	/// @arg third_foreground {struct<lipgloss.AdaptiveColor>}
	/// @arg third_background {struct<lipgloss.AdaptiveColor>}
	/// @arg fourth_foreground {struct<lipgloss.AdaptiveColor>}
	/// @arg fourth_background {struct<lipgloss.AdaptiveColor>}
	/// @returns {struct<tui.StatusBar>}
	lib.CreateFunction(tab, "statusbar",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "first_foreground"},
			{Type: lua.RAW_TABLE, Name: "first_background"},
			{Type: lua.RAW_TABLE, Name: "second_foreground"},
			{Type: lua.RAW_TABLE, Name: "second_background"},
			{Type: lua.RAW_TABLE, Name: "third_foreground"},
			{Type: lua.RAW_TABLE, Name: "third_background"},
			{Type: lua.RAW_TABLE, Name: "fourth_foreground"},
			{Type: lua.RAW_TABLE, Name: "fourth_background"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			prgrm := args["id"].(int)
			item, err := r.CR_TEA.Item(prgrm)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			firstForeground := lgColorGenericBuild(args["first_foreground"].(*golua.LTable))
			if _, ok := firstForeground.(lipgloss.AdaptiveColor); !ok {
				firstForeground = defaultForeground
			}
			firstBackground := lgColorGenericBuild(args["first_background"].(*golua.LTable))
			if _, ok := firstBackground.(lipgloss.AdaptiveColor); !ok {
				firstBackground = defaultBackground
			}
			secondForeground := lgColorGenericBuild(args["second_foreground"].(*golua.LTable))
			if _, ok := secondForeground.(lipgloss.AdaptiveColor); !ok {
				secondForeground = defaultForeground
			}
			secondBackground := lgColorGenericBuild(args["second_background"].(*golua.LTable))
			if _, ok := secondBackground.(lipgloss.AdaptiveColor); !ok {
				secondBackground = defaultBackground
			}
			thirdForeground := lgColorGenericBuild(args["third_foreground"].(*golua.LTable))
			if _, ok := thirdForeground.(lipgloss.AdaptiveColor); !ok {
				thirdForeground = defaultForeground
			}
			thirdBackground := lgColorGenericBuild(args["third_background"].(*golua.LTable))
			if _, ok := thirdBackground.(lipgloss.AdaptiveColor); !ok {
				thirdBackground = defaultBackground
			}
			fourthForeground := lgColorGenericBuild(args["fourth_foreground"].(*golua.LTable))
			if _, ok := fourthForeground.(lipgloss.AdaptiveColor); !ok {
				fourthForeground = defaultForeground
			}
			fourthBackground := lgColorGenericBuild(args["fourth_background"].(*golua.LTable))
			if _, ok := fourthBackground.(lipgloss.AdaptiveColor); !ok {
				fourthBackground = defaultBackground
			}

			firstPairs := statusbar.ColorConfig{
				Foreground: firstForeground.(lipgloss.AdaptiveColor),
				Background: firstBackground.(lipgloss.AdaptiveColor),
			}
			secondPairs := statusbar.ColorConfig{
				Foreground: secondForeground.(lipgloss.AdaptiveColor),
				Background: secondBackground.(lipgloss.AdaptiveColor),
			}
			thirdPairs := statusbar.ColorConfig{
				Foreground: thirdForeground.(lipgloss.AdaptiveColor),
				Background: thirdBackground.(lipgloss.AdaptiveColor),
			}
			fourthPairs := statusbar.ColorConfig{
				Foreground: fourthForeground.(lipgloss.AdaptiveColor),
				Background: fourthBackground.(lipgloss.AdaptiveColor),
			}

			sb := statusbar.New(firstPairs, secondPairs, thirdPairs, fourthPairs)
			id := len(item.StatusBars)
			item.StatusBars = append(item.StatusBars, &sb)

			t := statusbarTable(r, lg, lib, state, prgrm, id)

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
			/// @interface CMD
			/// @prop cmd {int<tui.CMDID>} - The command type.

			/// @struct CMDNone
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDNone(state))
			return 1
		})

	/// @func cmd_suspend() -> struct<tui.CMDSuspend>
	/// @returns {struct<tui.CMDSuspend>}
	lib.CreateFunction(tab, "cmd_suspend",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDSuspend
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDSuspend(state))
			return 1
		})

	/// @func cmd_quit() -> struct<tui.CMDQuit>
	/// @returns {struct<tui.CMDQuit>}
	lib.CreateFunction(tab, "cmd_quit",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDQuit
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDQuit(state))
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
			/// @struct CMDBatch
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop cmds {[]struct<tui.CMD>} - The commands to execute.

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
			/// @struct CMDSequence
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop cmds {[]struct<tui.CMD>} - The commands to execute.

			t := customtea.CMDSequence(state, args["cmds"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func cmd_printf(format, args...) -> struct<tui.CMDPrintf>
	/// @arg format {string}
	/// @arg args {any...}
	/// @returns {struct<tui.CMDPrintf>}
	lib.CreateFunction(tab, "cmd_printf",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
			lua.ArgVariadic("args", lua.ArrayType{Type: lua.ANY}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDPrintf
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop format {string} - The format string.
			/// @prop args {[]any} - The arguments to format.

			a := args["args"].([]any)
			at := state.NewTable()
			for i, v := range a {
				at.RawSetInt(i+1, v.(golua.LValue))
			}
			state.Push(customtea.CMDPrintf(state, args["format"].(string), at))
			return 1
		})

	/// @func cmd_println(args...) -> struct<tui.CMDPrintln>
	/// @arg args {any...}
	/// @returns {struct<tui.CMDPrintln>}
	lib.CreateFunction(tab, "cmd_println",
		[]lua.Arg{
			lua.ArgVariadic("args", lua.ArrayType{Type: lua.ANY}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDPrintln
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop args {[]any} - The arguments to print.

			a := args["args"].([]any)
			at := state.NewTable()
			for i, v := range a {
				at.RawSetInt(i+1, v.(golua.LValue))
			}
			state.Push(customtea.CMDPrintln(state, at))
			return 1
		})

	/// @func cmd_window_title(title) -> struct<tui.CMDWindowTitle>
	/// @arg title {string}
	/// @returns {struct<tui.CMDWindowTitle>}
	lib.CreateFunction(tab, "cmd_window_title",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDWindowTitle
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop title {string} - The title to set.

			state.Push(customtea.CMDWindowTitle(state, args["title"].(string)))
			return 1
		})

	/// @func cmd_window_size() -> struct<tui.CMDWindowSize>
	/// @returns {struct<tui.CMDWindowSize>}
	lib.CreateFunction(tab, "cmd_window_size",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDWindowSize
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDWindowSize(state))
			return 1
		})

	/// @func cmd_show_cursor() -> struct<tui.CMDShowCursor>
	/// @returns {struct<tui.CMDShowCursor>}
	lib.CreateFunction(tab, "cmd_show_cursor",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDShowCursor
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDShowCursor(state))
			return 1
		})

	/// @func cmd_hide_cursor() -> struct<tui.CMDHideCursor>
	/// @returns {struct<tui.CMDHideCursor>}
	lib.CreateFunction(tab, "cmd_hide_cursor",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDHideCursor
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDHideCursor(state))
			return 1
		})

	/// @func cmd_clear_screen() -> struct<tui.CMDClearScreen>
	/// @returns {struct<tui.CMDClearScreen>}
	lib.CreateFunction(tab, "cmd_clear_screen",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDClearScreen
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDClearScreen(state))
			return 1
		})

	/// @func cmd_clear_scroll_area() -> struct<tui.CMDClearScrollArea>
	/// @returns {struct<tui.CMDClearScrollArea>}
	lib.CreateFunction(tab, "cmd_clear_scroll_area",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDClearScrollArea
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDClearScrollArea(state))
			return 1
		})

	/// @func cmd_scroll_sync(lines, top, bottom) -> struct<tui.CMDScrollSync>
	/// @arg lines {[]string}
	/// @arg top {int}
	/// @arg bottom {int}
	/// @returns {struct<tui.CMDScrollSync>}
	lib.CreateFunction(tab, "cmd_scroll_sync",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "lines"},
			{Type: lua.INT, Name: "top"},
			{Type: lua.INT, Name: "bottom"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDScrollSync
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop lines {[]string} - The lines to display.
			/// @prop top {int} - The top line.
			/// @prop bottom {int} - The bottom line.

			state.Push(customtea.CMDScrollSync(state, args["lines"].(*golua.LTable), args["top"].(int), args["bottom"].(int)))
			return 1
		})

	/// @func cmd_scroll_up(lines, top, bottom) -> struct<tui.CMDScrollUp>
	/// @arg lines {[]string}
	/// @arg top {int}
	/// @arg bottom {int}
	/// @returns {struct<tui.CMDScrollUp>}
	lib.CreateFunction(tab, "cmd_scroll_up",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "lines"},
			{Type: lua.INT, Name: "top"},
			{Type: lua.INT, Name: "bottom"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDScrollUp
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop lines {[]string} - The lines to display.
			/// @prop top {int} - The top line.
			/// @prop bottom {int} - The bottom line.

			state.Push(customtea.CMDScrollUp(state, args["lines"].(*golua.LTable), args["top"].(int), args["bottom"].(int)))
			return 1
		})

	/// @func cmd_scroll_down(lines, top, bottom) -> struct<tui.CMDScrollDown>
	/// @arg lines {[]string}
	/// @arg top {int}
	/// @arg bottom {int}
	/// @returns {struct<tui.CMDScrollDown>}
	lib.CreateFunction(tab, "cmd_scroll_down",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "lines"},
			{Type: lua.INT, Name: "top"},
			{Type: lua.INT, Name: "bottom"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDScrollDown
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop lines {[]string} - The lines to display.
			/// @prop top {int} - The top line.
			/// @prop bottom {int} - The bottom line.

			state.Push(customtea.CMDScrollDown(state, args["lines"].(*golua.LTable), args["top"].(int), args["bottom"].(int)))
			return 1
		})

	/// @func cmd_every(duration, fn) -> struct<tui.CMDEvery>
	/// @arg duration {int} - Time in ms.
	/// @arg fn {function(ms int) -> any}
	/// @returns {struct<tui.CMDEvery>}
	lib.CreateFunction(tab, "cmd_every",
		[]lua.Arg{
			{Type: lua.INT, Name: "duration"},
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDEvery
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop duration {int} - Time in ms.
			/// @prop fn {function(ms int) -> any} - The function to call.

			state.Push(customtea.CMDEvery(state, args["duration"].(int), args["fn"].(*golua.LFunction)))
			return 1
		})

	/// @func cmd_tick(duration, fn) -> struct<tui.CMDTick>
	/// @arg duration {int} - Time in ms.
	/// @arg fn {function(ms int) -> any}
	/// @returns {struct<tui.CMDTick>}
	lib.CreateFunction(tab, "cmd_tick",
		[]lua.Arg{
			{Type: lua.INT, Name: "duration"},
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDTick
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop duration {int} - Time in ms.
			/// @prop fn {function(ms int) -> any} - The function to call.

			state.Push(customtea.CMDTick(state, args["duration"].(int), args["fn"].(*golua.LFunction)))
			return 1
		})

	/// @func cmd_toggle_report_focus(enabled?) -> struct<tui.CMDToggleReportFocus>
	/// @arg? enabled {bool}
	/// @returns {struct<tui.CMDToggleReportFocus>}
	lib.CreateFunction(tab, "cmd_toggle_report_focus",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDToggleReportFocus
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop enabled {bool} - Whether to report focus.

			state.Push(customtea.CMDToggleReportFocus(state, args["enabled"].(bool)))
			return 1
		})

	/// @func cmd_toggle_bracketed_paste(enabled?) -> struct<tui.CMDToggleBracketedPaste>
	/// @arg? enabled {bool}
	/// @returns {struct<tui.CMDToggleBracketedPaste>}
	lib.CreateFunction(tab, "cmd_toggle_bracketed_paste",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDToggleBracketedPaste
			/// @prop cmd {int<tui.CMDID>} - The command type.
			/// @prop enabled {bool} - Whether to enable bracketed paste.

			state.Push(customtea.CMDToggleBracketedPaste(state, args["enabled"].(bool)))
			return 1
		})

	/// @func cmd_disable_mouse() -> struct<tui.CMDDisableMouse>
	/// @returns {struct<tui.CMDDisableMouse>}
	lib.CreateFunction(tab, "cmd_disable_mouse",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDDisableMouse
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDDisableMouse(state))
			return 1
		})

	/// @func cmd_enable_mouse_all_motion() -> struct<tui.CMDEnableMouseAllMotion>
	/// @returns {struct<tui.CMDEnableMouseAllMotion>}
	lib.CreateFunction(tab, "cmd_enable_mouse_all_motion",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDEnableMouseAllMotion
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDEnableMouseAllMotion(state))
			return 1
		})

	/// @func cmd_enable_mouse_cell_motion() -> struct<tui.CMDEnableMouseCellMotion>
	/// @returns {struct<tui.CMDEnableMouseCellMotion>}
	lib.CreateFunction(tab, "cmd_enable_mouse_cell_motion",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDEnableMouseCellMotion
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDEnableMouseCellMotion(state))
			return 1
		})

	/// @func cmd_enter_alt_screen() -> struct<tui.CMDEnterAltScreen>
	/// @returns {struct<tui.CMDEnterAltScreen>}
	lib.CreateFunction(tab, "cmd_enter_alt_screen",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDEnterAltScreen
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDEnterAltScreen(state))
			return 1
		})

	/// @func cmd_exit_alt_screen() -> struct<tui.CMDExitAltScreen>
	/// @returns {struct<tui.CMDExitAltScreen>}
	lib.CreateFunction(tab, "cmd_exit_alt_screen",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct CMDExitAltScreen
			/// @prop cmd {int<tui.CMDID>} - The command type.

			state.Push(customtea.CMDExitAltScreen(state))
			return 1
		})

	/// @struct CMDStored
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int}

	/// @struct CMDSpinnerTick
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The spinner id.

	/// @struct CMDTextAreaFocus
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The text area id.

	/// @struct CMDTextInputFocus
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The text input id.

	/// @struct CMDBlink
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The cursor id.

	/// @struct CMDCursorFocus
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The cursor id.

	/// @struct CMDFilePickerInit
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The file picker id.

	/// @struct CMDListSetItems
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.
	/// @prop items {[]struct<tui.ListItem>} - The items to set.

	/// @struct CMDListInsertItem
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.
	/// @prop index {int} - The index to insert at.
	/// @prop item {struct<tui.ListItem>} - The item to insert.

	/// @struct CMDListSetItem
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.
	/// @prop index {int} - The index to set.
	/// @prop item {struct<tui.ListItem>} - The item to set.

	/// @struct CMDListStatusMessage
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.
	/// @prop msg {string} - The message to display.

	/// @struct CMDListSpinnerStart
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.

	/// @struct CMDListSpinnerToggle
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The list id.

	/// @struct CMDProgressSet
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The progress id.
	/// @prop percent {float} - The percentage.

	/// @struct CMDProgressDec
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The progress id.
	/// @prop percent {float} - The percentage to decrease.

	/// @struct CMDProgressInc
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The progress id.
	/// @prop percent {float} - The percentage to increase.

	/// @struct CMDStopWatchStart
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The stopwatch id.

	/// @struct CMDStopWatchStop
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The stopwatch id.

	/// @struct CMDStopWatchReset
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The stopwatch id.

	/// @struct CMDStopWatchToggle
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The stopwatch id.

	/// @struct CMDTimerStart
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The timer id.

	/// @struct CMDTimerInit
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The timer id.

	/// @struct CMDTimerStop
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The timer id.

	/// @struct CMDTimerToggle
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The timer id.

	/// @struct CMDImageSize
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The image id.
	/// @prop width {int} - The width.
	/// @prop height {int} - The height.

	/// @struct CMDImageFile
	/// @prop cmd {int<tui.CMDID>} - The command type.
	/// @prop id {int} - The image id.
	/// @prop filename {string} - The filename.

	/// @constants CMDID {int}
	/// @const CMD_NONE
	/// @const CMD_STORED
	/// @const CMD_BATCH
	/// @const CMD_SEQUENCE
	/// @const CMD_SPINNERTICK
	/// @const CMD_TEXTAREAFOCUS
	/// @const CMD_TEXTINPUTFOCUS
	/// @const CMD_BLINK
	/// @const CMD_CURSORFOCUS
	/// @const CMD_FILEPICKERINIT
	/// @const CMD_LISTSETITEMS
	/// @const CMD_LISTINSERTITEM
	/// @const CMD_LISTSETITEM
	/// @const CMD_LISTSTATUSMESSAGE
	/// @const CMD_LISTSPINNERSTART
	/// @const CMD_LISTSPINNERTOGGLE
	/// @const CMD_PROGRESSSET
	/// @const CMD_PROGRESSDEC
	/// @const CMD_PROGRESSINC
	/// @const CMD_STOPWATCHSTART
	/// @const CMD_STOPWATCHSTOP
	/// @const CMD_STOPWATCHTOGGLE
	/// @const CMD_STOPWATCHRESET
	/// @const CMD_TIMERINIT
	/// @const CMD_TIMERSTART
	/// @const CMD_TIMERSTOP
	/// @const CMD_TIMERTOGGLE
	/// @const CMD_VIEWPORTSYNC
	/// @const CMD_VIEWPORTUP
	/// @const CMD_VIEWPORTDOWN
	/// @const CMD_PRINTF
	/// @const CMD_PRINTLN
	/// @const CMD_WINDOWTITLE
	/// @const CMD_WINDOWSIZE
	/// @const CMD_SUSPEND
	/// @const CMD_QUIT
	/// @const CMD_SHOWCURSOR
	/// @const CMD_HIDECURSOR
	/// @const CMD_CLEARSCREEN
	/// @const CMD_CLEARSCROLLAREA
	/// @const CMD_SCROLLSYNC
	/// @const CMD_SCROLLUP
	/// @const CMD_SCROLLDOWN
	/// @const CMD_EVERY
	/// @const CMD_TICK
	/// @const CMD_TOGGLEREPORTFOCUS
	/// @const CMD_TOGGLEBRACKETEDPASTE
	/// @const CMD_DISABLEMOUSE
	/// @const CMD_ENABLEMOUSEALLMOTION
	/// @const CMD_ENABLEMOUSECELLMOTION
	/// @const CMD_ENTERALTSCREEN
	/// @const CMD_EXITALTSCREEN
	/// @const CMD_IMAGESIZE
	/// @const CMD_IMAGEFILE
	tab.RawSetString("CMD_NONE", golua.LNumber(customtea.CMD_NONE))
	tab.RawSetString("CMD_STORED", golua.LNumber(customtea.CMD_STORED))
	tab.RawSetString("CMD_BATCH", golua.LNumber(customtea.CMD_BATCH))
	tab.RawSetString("CMD_SEQUENCE", golua.LNumber(customtea.CMD_SEQUENCE))
	tab.RawSetString("CMD_SPINNERTICK", golua.LNumber(customtea.CMD_SPINNERTICK))
	tab.RawSetString("CMD_TEXTAREAFOCUS", golua.LNumber(customtea.CMD_TEXTAREAFOCUS))
	tab.RawSetString("CMD_TEXTINPUTFOCUS", golua.LNumber(customtea.CMD_TEXTINPUTFOCUS))
	tab.RawSetString("CMD_BLINK", golua.LNumber(customtea.CMD_BLINK))
	tab.RawSetString("CMD_CURSORFOCUS", golua.LNumber(customtea.CMD_CURSORFOCUS))
	tab.RawSetString("CMD_FILEPICKERINIT", golua.LNumber(customtea.CMD_FILEPICKERINIT))
	tab.RawSetString("CMD_LISTSETITEMS", golua.LNumber(customtea.CMD_LISTSETITEMS))
	tab.RawSetString("CMD_LISTINSERTITEM", golua.LNumber(customtea.CMD_LISTINSERTITEM))
	tab.RawSetString("CMD_LISTSETITEM", golua.LNumber(customtea.CMD_LISTSETITEM))
	tab.RawSetString("CMD_LISTSTATUSMESSAGE", golua.LNumber(customtea.CMD_LISTSTATUSMESSAGE))
	tab.RawSetString("CMD_LISTSPINNERSTART", golua.LNumber(customtea.CMD_LISTSPINNERSTART))
	tab.RawSetString("CMD_LISTSPINNERTOGGLE", golua.LNumber(customtea.CMD_LISTSPINNERTOGGLE))
	tab.RawSetString("CMD_PROGRESSSET", golua.LNumber(customtea.CMD_PROGRESSSET))
	tab.RawSetString("CMD_PROGRESSDEC", golua.LNumber(customtea.CMD_PROGRESSDEC))
	tab.RawSetString("CMD_PROGRESSINC", golua.LNumber(customtea.CMD_PROGRESSINC))
	tab.RawSetString("CMD_STOPWATCHSTART", golua.LNumber(customtea.CMD_STOPWATCHSTART))
	tab.RawSetString("CMD_STOPWATCHSTOP", golua.LNumber(customtea.CMD_STOPWATCHSTOP))
	tab.RawSetString("CMD_STOPWATCHTOGGLE", golua.LNumber(customtea.CMD_STOPWATCHTOGGLE))
	tab.RawSetString("CMD_STOPWATCHRESET", golua.LNumber(customtea.CMD_STOPWATCHRESET))
	tab.RawSetString("CMD_TIMERINIT", golua.LNumber(customtea.CMD_TIMERINIT))
	tab.RawSetString("CMD_TIMERSTART", golua.LNumber(customtea.CMD_TIMERSTART))
	tab.RawSetString("CMD_TIMERSTOP", golua.LNumber(customtea.CMD_TIMERSTOP))
	tab.RawSetString("CMD_TIMERTOGGLE", golua.LNumber(customtea.CMD_TIMERTOGGLE))
	tab.RawSetString("CMD_VIEWPORTSYNC", golua.LNumber(customtea.CMD_VIEWPORTSYNC))
	tab.RawSetString("CMD_VIEWPORTUP", golua.LNumber(customtea.CMD_VIEWPORTUP))
	tab.RawSetString("CMD_VIEWPORTDOWN", golua.LNumber(customtea.CMD_VIEWPORTDOWN))
	tab.RawSetString("CMD_PRINTF", golua.LNumber(customtea.CMD_PRINTF))
	tab.RawSetString("CMD_PRINTLN", golua.LNumber(customtea.CMD_PRINTLN))
	tab.RawSetString("CMD_WINDOWTITLE", golua.LNumber(customtea.CMD_WINDOWTITLE))
	tab.RawSetString("CMD_WINDOWSIZE", golua.LNumber(customtea.CMD_WINDOWSIZE))
	tab.RawSetString("CMD_SUSPEND", golua.LNumber(customtea.CMD_SUSPEND))
	tab.RawSetString("CMD_QUIT", golua.LNumber(customtea.CMD_QUIT))
	tab.RawSetString("CMD_SHOWCURSOR", golua.LNumber(customtea.CMD_SHOWCURSOR))
	tab.RawSetString("CMD_HIDECURSOR", golua.LNumber(customtea.CMD_HIDECURSOR))
	tab.RawSetString("CMD_CLEARSCREEN", golua.LNumber(customtea.CMD_CLEARSCREEN))
	tab.RawSetString("CMD_CLEARSCROLLAREA", golua.LNumber(customtea.CMD_CLEARSCROLLAREA))
	tab.RawSetString("CMD_SCROLLSYNC", golua.LNumber(customtea.CMD_SCROLLSYNC))
	tab.RawSetString("CMD_SCROLLUP", golua.LNumber(customtea.CMD_SCROLLUP))
	tab.RawSetString("CMD_SCROLLDOWN", golua.LNumber(customtea.CMD_SCROLLDOWN))
	tab.RawSetString("CMD_EVERY", golua.LNumber(customtea.CMD_EVERY))
	tab.RawSetString("CMD_TICK", golua.LNumber(customtea.CMD_TICK))
	tab.RawSetString("CMD_TOGGLEREPORTFOCUS", golua.LNumber(customtea.CMD_TOGGLEREPORTFOCUS))
	tab.RawSetString("CMD_TOGGLEBRACKETEDPASTE", golua.LNumber(customtea.CMD_TOGGLEBRACKETEDPASTE))
	tab.RawSetString("CMD_DISABLEMOUSE", golua.LNumber(customtea.CMD_DISABLEMOUSE))
	tab.RawSetString("CMD_ENABLEMOUSEALLMOTION", golua.LNumber(customtea.CMD_ENABLEMOUSEALLMOTION))
	tab.RawSetString("CMD_ENABLEMOUSECELLMOTION", golua.LNumber(customtea.CMD_ENABLEMOUSECELLMOTION))
	tab.RawSetString("CMD_ENTERALTSCREEN", golua.LNumber(customtea.CMD_ENTERALTSCREEN))
	tab.RawSetString("CMD_EXITALTSCREEN", golua.LNumber(customtea.CMD_EXITALTSCREEN))
	tab.RawSetString("CMD_IMAGESIZE", golua.LNumber(customtea.CMD_IMAGESIZE))
	tab.RawSetString("CMD_IMAGEFILE", golua.LNumber(customtea.CMD_IMAGEFILE))

	/// @interface MSG
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGNone
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGBlur
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGFocus
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGQuit
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGResume
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGSuspend
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGCursorBlink
	/// @prop msg {int<tui.MSGID>} - The message type.

	/// @struct MSGKey
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop key {string} - The key pressed.
	/// @prop event {struct<tui.KeyEvent>} - The key event.

	/// @struct MSGMouse
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop key {string} - The key pressed.
	/// @prop event {struct<tui.MouseEvent>} - The mouse event.

	/// @struct MSGWindowSize
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop width {int} - The width.
	/// @prop height {int} - The height.

	/// @struct MSGSpinnerTick
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The spinner id.

	/// @struct MSGStopwatchReset
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The stopwatch id.

	/// @struct MSGStopwatchStartStop
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The stopwatch id.

	/// @struct MSGStopwatchTick
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The stopwatch id.

	/// @struct MSGTimerStartStop
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The timer id.

	/// @struct MSGTimerTimeout
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The timer id.

	/// @struct MSGTimerTick
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop id {int} - The timer id.
	/// @prop timeout {int} - If this tick is a timeout.

	/// @struct MSGLua
	/// @prop msg {int<tui.MSGID>} - The message type.
	/// @prop value {any} - The value.

	/// @struct KeyEvent
	/// @prop type {int<tui.Key} - The key event type.
	/// @prop alt {bool} - Whether the alt key was pressed.
	/// @prop paste {bool}
	/// @prop runes {[]int}

	/// @struct MouseEvent
	/// @prop x {int} - The x position.
	/// @prop y {int} - The y position.
	/// @prop shift {bool} - Whether the shift key was pressed.
	/// @prop alt {bool} - Whether the alt key was pressed.
	/// @prop ctrl {bool} - Whether the ctrl key was pressed.
	/// @prop action {int<tui.MouseAction>} - The mouse action.
	/// @prop button {int<tui.MouseButton>} - The mouse button.
	/// @prop is_wheel {bool} - Whether the event is a wheel event.

	/// @constants MSGID {int}
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
	/// @const MSG_LUA
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
	tab.RawSetString("MSG_LUA", golua.LNumber(customtea.MSG_LUA))

	/// @constants Spinner {int}
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

	/// @constants EchoMode {int}
	/// @const ECHO_NORMAL
	/// @const ECHO_PASSWORD
	/// @const ECHO_NONE
	tab.RawSetString("ECHO_NORMAL", golua.LNumber(textinput.EchoNormal))
	tab.RawSetString("ECHO_PASSWORD", golua.LNumber(textinput.EchoPassword))
	tab.RawSetString("ECHO_NONE", golua.LNumber(textinput.EchoNone))

	/// @constants CursorMode {int}
	/// @const CURSOR_BLINK
	/// @const CURSOR_STATIC
	/// @const CURSOR_HIDE
	tab.RawSetString("CURSOR_BLINK", golua.LNumber(cursor.CursorBlink))
	tab.RawSetString("CURSOR_STATIC", golua.LNumber(cursor.CursorStatic))
	tab.RawSetString("CURSOR_HIDE", golua.LNumber(cursor.CursorHide))

	/// @constants FilterState {int}
	/// @const FILTERSTATE_UNFILTERED
	/// @const FILTERSTATE_FILTERING
	/// @const FILTERSTATE_APPLIED
	tab.RawSetString("FILTERSTATE_UNFILTERED", golua.LNumber(list.Unfiltered))
	tab.RawSetString("FILTERSTATE_FILTERING", golua.LNumber(list.Filtering))
	tab.RawSetString("FILTERSTATE_APPLIED", golua.LNumber(list.FilterApplied))

	/// @constants PaginatorType {int}
	/// @const PAGINATOR_ARABIC
	/// @const PAGINATOR_DOT
	tab.RawSetString("PAGINATOR_ARABIC", golua.LNumber(paginator.Arabic))
	tab.RawSetString("PAGINATOR_DOT", golua.LNumber(paginator.Dots))

	/// @constants FilterFunc {int}
	/// @const FILTERFUNC_DEFAULT
	/// @const FILTERFUNC_UNSORTED
	tab.RawSetString("FILTERFUNC_DEFAULT", golua.LNumber(FILTERFUNC_DEFAULT))
	tab.RawSetString("FILTERFUNC_UNSORTED", golua.LNumber(FILTERFUNC_UNSORTED))

	/// @constants MouseAction {int}
	/// @const MOUSE_PRESS
	/// @const MOUSE_RELEASE
	/// @const MOUSE_MOTION
	tab.RawSetString("MOUSE_PRESS", golua.LNumber(tea.MouseActionPress))
	tab.RawSetString("MOUSE_RELEASE", golua.LNumber(tea.MouseActionRelease))
	tab.RawSetString("MOUSE_MOTION", golua.LNumber(tea.MouseActionMotion))

	/// @constants MouseButton {int}
	/// @const MOUSEBUTTON_NONE
	/// @const MOUSEBUTTON_LEFT
	/// @const MOUSEBUTTON_MIDDLE
	/// @const MOUSEBUTTON_RIGHT
	/// @const MOUSEBUTTON_WHEELUP
	/// @const MOUSEBUTTON_WHEELDOWN
	/// @const MOUSEBUTTON_WHEELLEFT
	/// @const MOUSEBUTTON_WHEELRIGHT
	/// @const MOUSEBUTTON_BACKWARD
	/// @const MOUSEBUTTON_FORWARD
	/// @const MOUSEBUTTON_10
	/// @const MOUSEBUTTON_11
	tab.RawSetString("MOUSEBUTTON_NONE", golua.LNumber(tea.MouseButtonNone))
	tab.RawSetString("MOUSEBUTTON_LEFT", golua.LNumber(tea.MouseButtonLeft))
	tab.RawSetString("MOUSEBUTTON_MIDDLE", golua.LNumber(tea.MouseButtonMiddle))
	tab.RawSetString("MOUSEBUTTON_RIGHT", golua.LNumber(tea.MouseButtonRight))
	tab.RawSetString("MOUSEBUTTON_WHEELUP", golua.LNumber(tea.MouseButtonWheelUp))
	tab.RawSetString("MOUSEBUTTON_WHEELDOWN", golua.LNumber(tea.MouseButtonWheelDown))
	tab.RawSetString("MOUSEBUTTON_WHEELLEFT", golua.LNumber(tea.MouseButtonWheelLeft))
	tab.RawSetString("MOUSEBUTTON_WHEELRIGHT", golua.LNumber(tea.MouseButtonWheelRight))
	tab.RawSetString("MOUSEBUTTON_BACKWARD", golua.LNumber(tea.MouseButtonBackward))
	tab.RawSetString("MOUSEBUTTON_FORWARD", golua.LNumber(tea.MouseButtonForward))
	tab.RawSetString("MOUSEBUTTON_10", golua.LNumber(tea.MouseButton10))
	tab.RawSetString("MOUSEBUTTON_11", golua.LNumber(tea.MouseButton11))

	/// @constants Key {int}
	/// @const KEY_NULL
	/// @const KEY_BREAK
	/// @const KEY_ENTER
	/// @const KEY_BACKSPACE
	/// @const KEY_TAB
	/// @const KEY_ESC
	/// @const KEY_ESCAPE
	/// @const KEY_CTRL_AT
	/// @const KEY_CTRL_A
	/// @const KEY_CTRL_B
	/// @const KEY_CTRL_C
	/// @const KEY_CTRL_D
	/// @const KEY_CTRL_E
	/// @const KEY_CTRL_F
	/// @const KEY_CTRL_G
	/// @const KEY_CTRL_H
	/// @const KEY_CTRL_I
	/// @const KEY_CTRL_J
	/// @const KEY_CTRL_K
	/// @const KEY_CTRL_L
	/// @const KEY_CTRL_M
	/// @const KEY_CTRL_N
	/// @const KEY_CTRL_O
	/// @const KEY_CTRL_P
	/// @const KEY_CTRL_Q
	/// @const KEY_CTRL_R
	/// @const KEY_CTRL_S
	/// @const KEY_CTRL_T
	/// @const KEY_CTRL_U
	/// @const KEY_CTRL_V
	/// @const KEY_CTRL_W
	/// @const KEY_CTRL_X
	/// @const KEY_CTRL_Y
	/// @const KEY_CTRL_Z
	/// @const KEY_CTRL_OPEN_BRACKET
	/// @const KEY_CTRL_BACKSLASH
	/// @const KEY_CTRL_CLOSE_BRACKET
	/// @const KEY_CTRL_CARET
	/// @const KEY_CTRL_UNDERSCORE
	/// @const KEY_CTRL_QUESTION_MARK
	/// @const KEY_RUNES
	/// @const KEY_UP
	/// @const KEY_DOWN
	/// @const KEY_RIGHT
	/// @const KEY_LEFT
	/// @const KEY_SHIFTTAB
	/// @const KEY_HOME
	/// @const KEY_END
	/// @const KEY_PGUP
	/// @const KEY_PGDOWN
	/// @const KEY_CTRL_PGUP
	/// @const KEY_CTRL_PGDOWN
	/// @const KEY_DELETE
	/// @const KEY_INSERT
	/// @const KEY_SPACE
	/// @const KEY_CTRL_UP
	/// @const KEY_CTRL_DOWN
	/// @const KEY_CTRL_RIGHT
	/// @const KEY_CTRL_LEFT
	/// @const KEY_CTRL_HOME
	/// @const KEY_CTRL_END
	/// @const KEY_SHIFT_UP
	/// @const KEY_SHIFT_DOWN
	/// @const KEY_SHIFT_RIGHT
	/// @const KEY_SHIFT_LEFT
	/// @const KEY_SHIFT_HOME
	/// @const KEY_SHIFT_END
	/// @const KEY_CTRL_SHIFT_UP
	/// @const KEY_CTRL_SHIFT_DOWN
	/// @const KEY_CTRL_SHIFT_LEFT
	/// @const KEY_CTRL_SHIFT_RIGHT
	/// @const KEY_CTRL_SHIFT_HOME
	/// @const KEY_CTRL_SHIFT_END
	/// @const KEY_F1
	/// @const KEY_F2
	/// @const KEY_F3
	/// @const KEY_F4
	/// @const KEY_F5
	/// @const KEY_F6
	/// @const KEY_F7
	/// @const KEY_F8
	/// @const KEY_F9
	/// @const KEY_F10
	/// @const KEY_F11
	/// @const KEY_F12
	/// @const KEY_F13
	/// @const KEY_F14
	/// @const KEY_F15
	/// @const KEY_F16
	/// @const KEY_F17
	/// @const KEY_F18
	/// @const KEY_F19
	/// @const KEY_F20
	tab.RawSetString("KEY_NULL", golua.LNumber(tea.KeyNull))
	tab.RawSetString("KEY_BREAK", golua.LNumber(tea.KeyBreak))
	tab.RawSetString("KEY_ENTER", golua.LNumber(tea.KeyEnter))
	tab.RawSetString("KEY_BACKSPACE", golua.LNumber(tea.KeyBackspace))
	tab.RawSetString("KEY_TAB", golua.LNumber(tea.KeyTab))
	tab.RawSetString("KEY_ESC", golua.LNumber(tea.KeyEsc))
	tab.RawSetString("KEY_ESCAPE", golua.LNumber(tea.KeyEscape))
	tab.RawSetString("KEY_CTRL_AT", golua.LNumber(tea.KeyCtrlAt))
	tab.RawSetString("KEY_CTRL_A", golua.LNumber(tea.KeyCtrlA))
	tab.RawSetString("KEY_CTRL_B", golua.LNumber(tea.KeyCtrlB))
	tab.RawSetString("KEY_CTRL_C", golua.LNumber(tea.KeyCtrlC))
	tab.RawSetString("KEY_CTRL_D", golua.LNumber(tea.KeyCtrlD))
	tab.RawSetString("KEY_CTRL_E", golua.LNumber(tea.KeyCtrlE))
	tab.RawSetString("KEY_CTRL_F", golua.LNumber(tea.KeyCtrlF))
	tab.RawSetString("KEY_CTRL_G", golua.LNumber(tea.KeyCtrlG))
	tab.RawSetString("KEY_CTRL_H", golua.LNumber(tea.KeyCtrlH))
	tab.RawSetString("KEY_CTRL_I", golua.LNumber(tea.KeyCtrlI))
	tab.RawSetString("KEY_CTRL_J", golua.LNumber(tea.KeyCtrlJ))
	tab.RawSetString("KEY_CTRL_K", golua.LNumber(tea.KeyCtrlK))
	tab.RawSetString("KEY_CTRL_L", golua.LNumber(tea.KeyCtrlL))
	tab.RawSetString("KEY_CTRL_M", golua.LNumber(tea.KeyCtrlM))
	tab.RawSetString("KEY_CTRL_N", golua.LNumber(tea.KeyCtrlN))
	tab.RawSetString("KEY_CTRL_O", golua.LNumber(tea.KeyCtrlO))
	tab.RawSetString("KEY_CTRL_P", golua.LNumber(tea.KeyCtrlP))
	tab.RawSetString("KEY_CTRL_Q", golua.LNumber(tea.KeyCtrlQ))
	tab.RawSetString("KEY_CTRL_R", golua.LNumber(tea.KeyCtrlR))
	tab.RawSetString("KEY_CTRL_S", golua.LNumber(tea.KeyCtrlS))
	tab.RawSetString("KEY_CTRL_T", golua.LNumber(tea.KeyCtrlT))
	tab.RawSetString("KEY_CTRL_U", golua.LNumber(tea.KeyCtrlU))
	tab.RawSetString("KEY_CTRL_V", golua.LNumber(tea.KeyCtrlV))
	tab.RawSetString("KEY_CTRL_W", golua.LNumber(tea.KeyCtrlW))
	tab.RawSetString("KEY_CTRL_X", golua.LNumber(tea.KeyCtrlX))
	tab.RawSetString("KEY_CTRL_Y", golua.LNumber(tea.KeyCtrlY))
	tab.RawSetString("KEY_CTRL_Z", golua.LNumber(tea.KeyCtrlZ))
	tab.RawSetString("KEY_CTRL_OPEN_BRACKET", golua.LNumber(tea.KeyCtrlOpenBracket))
	tab.RawSetString("KEY_CTRL_BACKSLASH", golua.LNumber(tea.KeyCtrlBackslash))
	tab.RawSetString("KEY_CTRL_CLOSE_BRACKET", golua.LNumber(tea.KeyCtrlCloseBracket))
	tab.RawSetString("KEY_CTRL_CARET", golua.LNumber(tea.KeyCtrlCaret))
	tab.RawSetString("KEY_CTRL_UNDERSCORE", golua.LNumber(tea.KeyCtrlUnderscore))
	tab.RawSetString("KEY_CTRL_QUESTION_MARK", golua.LNumber(tea.KeyCtrlQuestionMark))
	tab.RawSetString("KEY_RUNES", golua.LNumber(tea.KeyRunes))
	tab.RawSetString("KEY_UP", golua.LNumber(tea.KeyUp))
	tab.RawSetString("KEY_DOWN", golua.LNumber(tea.KeyDown))
	tab.RawSetString("KEY_RIGHT", golua.LNumber(tea.KeyRight))
	tab.RawSetString("KEY_LEFT", golua.LNumber(tea.KeyLeft))
	tab.RawSetString("KEY_SHIFTTAB", golua.LNumber(tea.KeyShiftTab))
	tab.RawSetString("KEY_HOME", golua.LNumber(tea.KeyHome))
	tab.RawSetString("KEY_END", golua.LNumber(tea.KeyEnd))
	tab.RawSetString("KEY_PGUP", golua.LNumber(tea.KeyPgUp))
	tab.RawSetString("KEY_PGDOWN", golua.LNumber(tea.KeyPgDown))
	tab.RawSetString("KEY_CTRL_PGUP", golua.LNumber(tea.KeyCtrlPgUp))
	tab.RawSetString("KEY_CTRL_PGDOWN", golua.LNumber(tea.KeyCtrlPgDown))
	tab.RawSetString("KEY_DELETE", golua.LNumber(tea.KeyDelete))
	tab.RawSetString("KEY_INSERT", golua.LNumber(tea.KeyInsert))
	tab.RawSetString("KEY_SPACE", golua.LNumber(tea.KeySpace))
	tab.RawSetString("KEY_CTRL_UP", golua.LNumber(tea.KeyCtrlUp))
	tab.RawSetString("KEY_CTRL_DOWN", golua.LNumber(tea.KeyCtrlDown))
	tab.RawSetString("KEY_CTRL_RIGHT", golua.LNumber(tea.KeyCtrlRight))
	tab.RawSetString("KEY_CTRL_LEFT", golua.LNumber(tea.KeyCtrlLeft))
	tab.RawSetString("KEY_CTRL_HOME", golua.LNumber(tea.KeyCtrlHome))
	tab.RawSetString("KEY_CTRL_END", golua.LNumber(tea.KeyCtrlEnd))
	tab.RawSetString("KEY_SHIFT_UP", golua.LNumber(tea.KeyShiftUp))
	tab.RawSetString("KEY_SHIFT_DOWN", golua.LNumber(tea.KeyShiftDown))
	tab.RawSetString("KEY_SHIFT_RIGHT", golua.LNumber(tea.KeyShiftRight))
	tab.RawSetString("KEY_SHIFT_LEFT", golua.LNumber(tea.KeyShiftLeft))
	tab.RawSetString("KEY_SHIFT_HOME", golua.LNumber(tea.KeyShiftHome))
	tab.RawSetString("KEY_SHIFT_END", golua.LNumber(tea.KeyShiftEnd))
	tab.RawSetString("KEY_CTRL_SHIFT_UP", golua.LNumber(tea.KeyCtrlShiftUp))
	tab.RawSetString("KEY_CTRL_SHIFT_DOWN", golua.LNumber(tea.KeyCtrlShiftDown))
	tab.RawSetString("KEY_CTRL_SHIFT_LEFT", golua.LNumber(tea.KeyCtrlShiftLeft))
	tab.RawSetString("KEY_CTRL_SHIFT_RIGHT", golua.LNumber(tea.KeyCtrlShiftRight))
	tab.RawSetString("KEY_CTRL_SHIFT_HOME", golua.LNumber(tea.KeyCtrlShiftHome))
	tab.RawSetString("KEY_CTRL_SHIFT_END", golua.LNumber(tea.KeyCtrlShiftEnd))
	tab.RawSetString("KEY_F1", golua.LNumber(tea.KeyF1))
	tab.RawSetString("KEY_F2", golua.LNumber(tea.KeyF2))
	tab.RawSetString("KEY_F3", golua.LNumber(tea.KeyF3))
	tab.RawSetString("KEY_F4", golua.LNumber(tea.KeyF4))
	tab.RawSetString("KEY_F5", golua.LNumber(tea.KeyF5))
	tab.RawSetString("KEY_F6", golua.LNumber(tea.KeyF6))
	tab.RawSetString("KEY_F7", golua.LNumber(tea.KeyF7))
	tab.RawSetString("KEY_F8", golua.LNumber(tea.KeyF8))
	tab.RawSetString("KEY_F9", golua.LNumber(tea.KeyF9))
	tab.RawSetString("KEY_F10", golua.LNumber(tea.KeyF10))
	tab.RawSetString("KEY_F11", golua.LNumber(tea.KeyF11))
	tab.RawSetString("KEY_F12", golua.LNumber(tea.KeyF12))
	tab.RawSetString("KEY_F13", golua.LNumber(tea.KeyF13))
	tab.RawSetString("KEY_F14", golua.LNumber(tea.KeyF14))
	tab.RawSetString("KEY_F15", golua.LNumber(tea.KeyF15))
	tab.RawSetString("KEY_F16", golua.LNumber(tea.KeyF16))
	tab.RawSetString("KEY_F17", golua.LNumber(tea.KeyF17))
	tab.RawSetString("KEY_F18", golua.LNumber(tea.KeyF18))
	tab.RawSetString("KEY_F19", golua.LNumber(tea.KeyF19))
	tab.RawSetString("KEY_F20", golua.LNumber(tea.KeyF20))
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

type FilterFunc int

const (
	FILTERFUNC_DEFAULT FilterFunc = iota
	FILTERFUNC_UNSORTED
)

func teaTable(r *lua.Runner, lg *log.Logger, state *golua.LState, lib *lua.Lib, id int) *golua.LTable {
	/// @struct Program
	/// @prop id {int}
	/// @method init(self, {function(id int<collection.CRATE_TEA>) -> any, struct<tui.CMD>}) -> self
	/// @method update(self, {function(data any, struct<tui.MSG>) -> struct<tui.CMD>}) -> self
	/// @method view(self, {function(data any) -> string}) -> self

	t := state.NewTable()
	t.RawSetString("id", golua.LNumber(id))

	lib.BuilderFunction(state, t, "init",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("id").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			item.FnView = args["fn"].(*golua.LFunction)
		},
	)

	return t
}

func programOptions(lib *lua.Lib, state *golua.LState) *golua.LTable {
	/// @struct ProgramOptions
	/// @method ansi_compressor(self) -> self
	/// @method alt_screen(self) -> self
	/// @method fps(self, fps int) -> self
	/// @method filter(self, filter {function(msg struct<tui.MSG>) -> bool}) -> self
	/// @method input_tty(self) -> self
	/// @method mouse_all_motion(self) -> self
	/// @method mouse_cell_motion(self) -> self
	/// @method report_focus(self) -> self
	/// @method no_bracketed_paste(self) -> self

	t := state.NewTable()

	t.RawSetString("__ansiCompressor", golua.LFalse)
	t.RawSetString("__altScreen", golua.LFalse)
	t.RawSetString("__fps", golua.LNil)
	t.RawSetString("__filter", golua.LNil)
	t.RawSetString("__inputTTY", golua.LFalse)
	t.RawSetString("__mouseAllMotion", golua.LFalse)
	t.RawSetString("__mouseCellMotion", golua.LFalse)
	t.RawSetString("__reportFocus", golua.LFalse)
	t.RawSetString("__noBracketedPaste", golua.LFalse)

	lib.BuilderFunction(state, t, "ansi_compressor",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__ansiCompressor", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "alt_screen",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__altScreen", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "fps",
		[]lua.Arg{
			{Type: lua.INT, Name: "fps"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__fps", golua.LNumber(args["fps"].(int)))
		})

	lib.BuilderFunction(state, t, "filter",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "filter"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__filter", args["filter"].(*golua.LFunction))
		})

	lib.BuilderFunction(state, t, "input_tty",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__inputTTY", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "mouse_all_motion",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__mouseAllMotion", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "mouse_cell_motion",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__mouseCellMotion", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "report_focus",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__reportFocus", golua.LTrue)
		})

	lib.BuilderFunction(state, t, "no_bracketed_paste",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__noBracketedPaste", golua.LTrue)
		})

	return t
}

func programOptionsBuild(state *golua.LState, t *golua.LTable) []tea.ProgramOption {
	opts := []tea.ProgramOption{}

	ansiCompressor := t.RawGetString("__ansiCompressor")
	if ansiCompressor.Type() == golua.LTBool && bool(ansiCompressor.(golua.LBool)) {
		opts = append(opts, tea.WithANSICompressor())
	}

	altScreen := t.RawGetString("__altScreen")
	if altScreen.Type() == golua.LTBool && bool(altScreen.(golua.LBool)) {
		opts = append(opts, tea.WithAltScreen())
	}

	fps := t.RawGetString("__fps")
	if altScreen.Type() == golua.LTNumber {
		opts = append(opts, tea.WithFPS(int(fps.(golua.LNumber))))
	}

	filter := t.RawGetString("__filter")
	if filter.Type() == golua.LTFunction {
		fn := filter.(*golua.LFunction)
		opts = append(opts, tea.WithFilter(func(m1 tea.Model, m2 tea.Msg) tea.Msg {
			if _, ok := m2.(tea.QuitMsg); ok {
				return m2
			}

			mt := customtea.BuildMSG(m2, state)
			state.Push(fn)
			state.Push(mt)
			state.Call(1, 1)
			pass := state.CheckBool(-1)
			state.Pop(1)

			if !pass {
				return nil
			}

			return m2
		}))
	}

	inputTTY := t.RawGetString("__inputTTY")
	if inputTTY.Type() == golua.LTBool && bool(inputTTY.(golua.LBool)) {
		opts = append(opts, tea.WithInputTTY())
	}

	mouseAllMotion := t.RawGetString("__mouseAllMotion")
	if mouseAllMotion.Type() == golua.LTBool && bool(mouseAllMotion.(golua.LBool)) {
		opts = append(opts, tea.WithMouseAllMotion())
	}

	mouseCellMotion := t.RawGetString("__mouseCellMotion")
	if mouseCellMotion.Type() == golua.LTBool && bool(mouseCellMotion.(golua.LBool)) {
		opts = append(opts, tea.WithMouseCellMotion())
	}

	reportFocus := t.RawGetString("__reportFocus")
	if reportFocus.Type() == golua.LTBool && bool(reportFocus.(golua.LBool)) {
		opts = append(opts, tea.WithReportFocus())
	}

	noBracketedPaste := t.RawGetString("__noBracketedPaste")
	if noBracketedPaste.Type() == golua.LTBool && bool(noBracketedPaste.(golua.LBool)) {
		opts = append(opts, tea.WithoutBracketedPaste())
	}

	return opts
}

func spinnerTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Spinner
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method tick() -> struct<tui.CMDSpinnerTick>
	/// @method spinner() -> []string, int
	/// @method spinner_set(self, from int) -> self
	/// @method spinner_set_custom(self, frames []string, fps int) -> self
	/// @method style() -> struct<lipgloss.Style>
	/// @method style_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Spinners[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			spinner := item.Spinners[id].Spinner

			fps := time.Second / spinner.FPS

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func textareaTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct TextArea
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method reset(self) -> self
	/// @method focus() -> struct<tui.CMDTextAreaFocus>
	/// @method blur(self) -> self
	/// @method cursor_down(self) -> self
	/// @method cursor_end(self) -> self
	/// @method cursor_up(self) -> self
	/// @method cursor_down(self) -> self
	/// @method focused() -> bool
	/// @method size() -> int, int
	/// @method width() -> int
	/// @method height() -> int
	/// @method size_set(self, width int, height int) -> self
	/// @method width_set(self, width int) -> self
	/// @method height_set(self, height int) -> self
	/// @method insert_rune(self, rune int) -> self
	/// @method insert_string(self, str string) -> self
	/// @method length() -> int
	/// @method line() -> int
	/// @method line_count() -> int
	/// @method cursor_set(self, col int) -> self
	/// @method value() -> string
	/// @method value_set(str string)
	/// @method line_info() -> struct<tui.LineInfo>
	/// @method prompt() -> string
	/// @method prompt_set(self, str string) -> self
	/// @method line_numbers() -> bool
	/// @method line_numbers_set(self, enabled bool) -> self
	/// @method char_end() -> int
	/// @method char_end_set(self, rune int) -> self
	/// @method char_limit() -> int
	/// @method char_limit_set(self, limit int) -> self
	/// @method width_max() -> int
	/// @method width_max_set(self, width int) -> self
	/// @method height_max() -> int
	/// @method height_max_set(self, height int) -> self
	/// @method prompt_func(self, width int, fn {function(lineIndex int) -> string}) -> self
	/// @method cursor() -> struct<tui.Cursor>
	/// @method keymap() -> struct<tui.TextAreaKeymap>
	/// @method style_focus_base() -> struct<lipgloss.Style>
	/// @method style_focus_base_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_base() -> struct<lipgloss.Style>
	/// @method style_blur_base_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_cursor_line() -> struct<lipgloss.Style>
	/// @method style_focus_cursor_line_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_cursor_line() -> struct<lipgloss.Style>
	/// @method style_blur_cursor_line_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_cursor_line_number() -> struct<lipgloss.Style>
	/// @method style_focus_cursor_line_number_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_cursor_line_number() -> struct<lipgloss.Style>
	/// @method style_blur_cursor_line_number_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_buffer_end() -> struct<lipgloss.Style>
	/// @method style_focus_buffer_end_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_buffer_end() -> struct<lipgloss.Style>
	/// @method style_blur_buffer_end_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_line_number() -> struct<lipgloss.Style>
	/// @method style_focus_line_number_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_line_number() -> struct<lipgloss.Style>
	/// @method style_blur_line_number_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_placeholder() -> struct<lipgloss.Style>
	/// @method style_focus_placeholder_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_placeholder() -> struct<lipgloss.Style>
	/// @method style_blur_placeholder_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_prompt() -> struct<lipgloss.Style>
	/// @method style_focus_prompt_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_prompt() -> struct<lipgloss.Style>
	/// @method style_blur_prompt_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_focus_text() -> struct<lipgloss.Style>
	/// @method style_focus_text_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_blur_text() -> struct<lipgloss.Style>
	/// @method style_blur_text_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.TextAreas[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].Blur()
		})

	lib.BuilderFunction(state, t, "cursor_down",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorDown()
		})

	lib.BuilderFunction(state, t, "cursor_end",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorEnd()
		})

	lib.BuilderFunction(state, t, "cursor_start",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorStart()
		})

	lib.BuilderFunction(state, t, "cursor_up",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].CursorUp()
		})

	t.RawSetString("focused", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		focused := item.TextAreas[id].Focused()

		state.Push(golua.LBool(focused))
		return 1
	}))

	t.RawSetString("size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.TextAreas[id].Width()

		state.Push(golua.LNumber(width))
		return 1
	}))

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].InsertString(args["str"].(string))
		})

	t.RawSetString("length", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		length := item.TextAreas[id].Length()

		state.Push(golua.LNumber(length))
		return 1
	}))

	t.RawSetString("line", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		line := item.TextAreas[id].Line()

		state.Push(golua.LNumber(line))
		return 1
	}))

	t.RawSetString("line_count", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetCursor(args["col"].(int))
		})

	t.RawSetString("value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextAreas[id].SetValue(args["str"].(string))
		})

	t.RawSetString("line_info", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		info := item.TextAreas[id].LineInfo()

		state.Push(lineInfoTable(state, &info))
		return 1
	}))

	t.RawSetString("prompt", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.Prompt = args["str"].(string)
			ta.SetWidth(ta.Width())
		})

	t.RawSetString("line_numbers", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.ShowLineNumbers = args["enable"].(bool)
			ta.SetWidth(ta.Width())
		})

	t.RawSetString("char_end", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.EndOfBufferCharacter = rune(args["rune"].(int))
		})

	t.RawSetString("char_limit", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.CharLimit = args["limit"].(int)
		})

	t.RawSetString("width_max", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))
			ta := item.TextAreas[id]

			ta.MaxWidth = args["width"].(int)
		})

	t.RawSetString("height_max", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextAreas[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)
		cu := cursorTable(r, lg, lib, state, program, cid)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := textareaKeymapTable(r, lg, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	t.RawSetString("__styleFocusBase", golua.LNil)
	lib.TableFunction(state, t, "style_focus_base",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusBase")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.Base
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusBase", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_base_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.Base = *style.Style
			t.RawSetString("__styleFocusBase", st)
		})

	t.RawSetString("__styleBlurBase", golua.LNil)
	lib.TableFunction(state, t, "style_blur_base",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurBase")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.Base
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurBase", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_base_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.Base = *style.Style
			t.RawSetString("__styleBlurBase", st)
		})

	t.RawSetString("__styleFocusCursorLine", golua.LNil)
	lib.TableFunction(state, t, "style_focus_cursor_line",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusCursorLine")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.CursorLine
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusCursorLine", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_cursor_line_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.CursorLine = *style.Style
			t.RawSetString("__styleFocusCursorLine", st)
		})

	t.RawSetString("__styleBlurCursorLine", golua.LNil)
	lib.TableFunction(state, t, "style_blur_cursor_line",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurCursorLine")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.CursorLine
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurCursorLine", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_cursor_line_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.CursorLine = *style.Style
			t.RawSetString("__styleBlurCursorLine", st)
		})

	t.RawSetString("__styleFocusCursorLineNumber", golua.LNil)
	lib.TableFunction(state, t, "style_focus_cursor_line_number",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusCursorLineNumber")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.CursorLineNumber
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusCursorLineNumber", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_cursor_line_number_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.CursorLineNumber = *style.Style
			t.RawSetString("__styleFocusCursorLineNumber", st)
		})

	t.RawSetString("__styleBlurCursorLineNumber", golua.LNil)
	lib.TableFunction(state, t, "style_blur_cursor_line_number",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurCursorLineNumber")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.CursorLineNumber
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurCursorLineNumber", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_cursor_line_number_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.CursorLineNumber = *style.Style
			t.RawSetString("__styleBlurCursorLineNumber", st)
		})

	t.RawSetString("__styleFocusBufferEnd", golua.LNil)
	lib.TableFunction(state, t, "style_focus_buffer_end",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusBufferEnd")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.EndOfBuffer
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusBufferEnd", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_buffer_end_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.EndOfBuffer = *style.Style
			t.RawSetString("__styleFocusEndOfBuffer", st)
		})

	t.RawSetString("__styleBlurBufferEnd", golua.LNil)
	lib.TableFunction(state, t, "style_blur_buffer_end",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurBufferEnd")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.EndOfBuffer
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurBufferEnd", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_buffer_end_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.EndOfBuffer = *style.Style
			t.RawSetString("__styleBlurEndOfBuffer", st)
		})

	t.RawSetString("__styleFocusLineNumber", golua.LNil)
	lib.TableFunction(state, t, "style_focus_line_number",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusLineNumber")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.LineNumber
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusLineNumber", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_line_number_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.LineNumber = *style.Style
			t.RawSetString("__styleFocusLineNumber", st)
		})

	t.RawSetString("__styleBlurLineNumber", golua.LNil)
	lib.TableFunction(state, t, "style_blur_line_number",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurLineNumber")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.LineNumber
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurLineNumber", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_line_number_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.LineNumber = *style.Style
			t.RawSetString("__styleBlurLineNumber", st)
		})

	t.RawSetString("__styleFocusPlaceholder", golua.LNil)
	lib.TableFunction(state, t, "style_focus_placeholder",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusPlaceholder")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.Placeholder
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusPlaceholder", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_placeholder_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.Placeholder = *style.Style
			t.RawSetString("__styleFocusPlaceholder", st)
		})

	t.RawSetString("__styleBlurPlaceholder", golua.LNil)
	lib.TableFunction(state, t, "style_blur_placeholder",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurPlaceholder")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.Placeholder
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurPlaceholder", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_placeholder_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.Placeholder = *style.Style
			t.RawSetString("__styleBlurPlaceholder", st)
		})

	t.RawSetString("__styleFocusPrompt", golua.LNil)
	lib.TableFunction(state, t, "style_focus_prompt",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusPrompt")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.Prompt
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusPrompt", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_prompt_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.Prompt = *style.Style
			t.RawSetString("__styleFocusPrompt", st)
		})

	t.RawSetString("__styleBlurPrompt", golua.LNil)
	lib.TableFunction(state, t, "style_blur_prompt",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurPrompt")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.Prompt
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurPrompt", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_prompt_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.Prompt = *style.Style
			t.RawSetString("__styleBlurPrompt", st)
		})

	t.RawSetString("__styleFocusText", golua.LNil)
	lib.TableFunction(state, t, "style_focus_text",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFocusText")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].FocusedStyle.Text
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFocusText", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_focus_text_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].FocusedStyle.Text = *style.Style
			t.RawSetString("__styleFocusText", st)
		})

	t.RawSetString("__styleBlurText", golua.LNil)
	lib.TableFunction(state, t, "style_blur_text",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleBlurText")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.TextAreas[id].BlurredStyle.Text
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleBlurText", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_blur_text_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.TextAreas[id].BlurredStyle.Text = *style.Style
			t.RawSetString("__styleBlurText", st)
		})

	return t
}

func textareaKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [22]int) *golua.LTable {
	/// @struct TextAreaKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop character_backward {struct<tui.KeyBinding>}
	/// @prop character_forward {struct<tui.KeyBinding>}
	/// @prop delete_after_cursor {struct<tui.KeyBinding>}
	/// @prop delete_before_cursor {struct<tui.KeyBinding>}
	/// @prop delete_character_backward {struct<tui.KeyBinding>}
	/// @prop delete_character_forward {struct<tui.KeyBinding>}
	/// @prop delete_word_backward {struct<tui.KeyBinding>}
	/// @prop delete_word_forward {struct<tui.KeyBinding>}
	/// @prop insert_newline {struct<tui.KeyBinding>}
	/// @prop line_end {struct<tui.KeyBinding>}
	/// @prop line_next {struct<tui.KeyBinding>}
	/// @prop line_previous {struct<tui.KeyBinding>}
	/// @prop line_start {struct<tui.KeyBinding>}
	/// @prop paste {struct<tui.KeyBinding>}
	/// @prop word_backward {struct<tui.KeyBinding>}
	/// @prop word_forward {struct<tui.KeyBinding>}
	/// @prop input_begin {struct<tui.KeyBinding>}
	/// @prop input_end {struct<tui.KeyBinding>}
	/// @prop uppercase_word {struct<tui.KeyBinding>}
	/// @prop lowercase_word {struct<tui.KeyBinding>}
	/// @prop capitalize_word {struct<tui.KeyBinding>}
	/// @prop transpose_character_backward {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("character_backward", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("character_forward", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("delete_after_cursor", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("delete_before_cursor", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("delete_character_backward", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("delete_character_forward", tuikeyTable(r, lg, lib, state, program, ids[5]))
	t.RawSetString("delete_word_backward", tuikeyTable(r, lg, lib, state, program, ids[6]))
	t.RawSetString("delete_word_forward", tuikeyTable(r, lg, lib, state, program, ids[7]))
	t.RawSetString("insert_newline", tuikeyTable(r, lg, lib, state, program, ids[8]))
	t.RawSetString("line_end", tuikeyTable(r, lg, lib, state, program, ids[9]))
	t.RawSetString("line_next", tuikeyTable(r, lg, lib, state, program, ids[10]))
	t.RawSetString("line_previous", tuikeyTable(r, lg, lib, state, program, ids[11]))
	t.RawSetString("line_start", tuikeyTable(r, lg, lib, state, program, ids[12]))
	t.RawSetString("paste", tuikeyTable(r, lg, lib, state, program, ids[13]))
	t.RawSetString("word_backward", tuikeyTable(r, lg, lib, state, program, ids[14]))
	t.RawSetString("word_forward", tuikeyTable(r, lg, lib, state, program, ids[15]))
	t.RawSetString("input_begin", tuikeyTable(r, lg, lib, state, program, ids[16]))
	t.RawSetString("input_end", tuikeyTable(r, lg, lib, state, program, ids[17]))
	t.RawSetString("uppercase_word", tuikeyTable(r, lg, lib, state, program, ids[18]))
	t.RawSetString("lowercase_word", tuikeyTable(r, lg, lib, state, program, ids[19]))
	t.RawSetString("capitalize_word", tuikeyTable(r, lg, lib, state, program, ids[20]))
	t.RawSetString("transpose_character_backward", tuikeyTable(r, lg, lib, state, program, ids[21]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
	/// @struct LineInfo
	/// @prop width {int}
	/// @prop width_char {int}
	/// @prop height {int}
	/// @prop column_start {int}
	/// @prop column_offset {int}
	/// @prop row_offset {int}
	/// @prop char_offset {int}

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

func textinputTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct TextInput
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tea.Cmd>
	/// @method focus() -> struct<tea.CmdTextInputFocus>
	/// @method reset(self) -> self
	/// @method blur(self) -> self
	/// @method cursor_start(self) -> self
	/// @method cursor_end(self) -> self
	/// @method current_suggestion() -> string
	/// @method available_suggestions() -> []string
	/// @method suggestions_set(self, suggestions []string) -> self
	/// @method focused() -> bool
	/// @method position() -> int
	/// @method position_set(self, pos int) -> self
	/// @method value() -> string
	/// @method value_set(self, val string) -> self
	/// @method validate(self, fn func(string) -> bool, string) -> self
	/// @method prompt() -> string
	/// @method prompt_set(self, string) -> self
	/// @method placeholder() -> string
	/// @method placeholder_set(self, string) -> self
	/// @method echomode() -> int<tui.EchoMode>
	/// @method echomode_set(self, int<tui.EchoMode>) -> self
	/// @method echo_char() -> int
	/// @method echo_char_set(self, rune int) -> self
	/// @method char_limit() -> int
	/// @method char_limit_set(self, int) -> self
	/// @method width() -> int
	/// @method width_set(self, int) -> self
	/// @method suggestions_show() -> bool
	/// @method suggestions_show_set(self, bool) -> self
	/// @method cursor() -> struct<tui.Cursor>
	/// @method keymap() -> struct<TextInputKeymap>
	/// @method style_prompt() -> struct<lipgloss.Style>
	/// @method style_prompt_set(self, struct<lipgloss.Style>) -> self
	/// @method style_text() -> struct<lipgloss.Style>
	/// @method style_text_set(self, struct<lipgloss.Style>) -> self
	/// @method style_placeholder() -> struct<lipgloss.Style>
	/// @method style_placeholder_set(self, struct<lipgloss.Style>) -> self
	/// @method style_completion() -> struct<lipgloss.Style>
	/// @method style_completion_set(self, struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.TextInputs[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Reset()
		})

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Blur()
		})

	lib.BuilderFunction(state, t, "cursor_start",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CursorStart()
		})

	lib.BuilderFunction(state, t, "cursor_end",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CursorEnd()
		})

	t.RawSetString("current_suggestion", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		suggestion := item.TextInputs[id].CurrentSuggestion()

		state.Push(golua.LString(suggestion))
		return 1
	}))

	t.RawSetString("available_suggestions", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		focused := item.TextInputs[id].Focused()

		state.Push(golua.LBool(focused))
		return 1
	}))

	t.RawSetString("position", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].SetCursor(args["pos"].(int))
		})

	t.RawSetString("value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Prompt = args["value"].(string)
		})

	t.RawSetString("placeholder", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Placeholder = args["value"].(string)
		})

	t.RawSetString("echomode", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].EchoMode = textinput.EchoMode(args["echomode"].(int))
		})

	t.RawSetString("echo_char", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].EchoCharacter = rune(args["rune"].(int))
		})

	t.RawSetString("char_limit", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].CharLimit = args["limit"].(int)
		})

	t.RawSetString("width", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.TextInputs[id].Width = args["width"].(int)
		})

	t.RawSetString("suggestions_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		ta := item.TextInputs[id]
		cid := len(item.Cursors)
		item.Cursors = append(item.Cursors, &ta.Cursor)
		cu := cursorTable(r, lg, lib, state, program, cid)

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := textinputKeymapTable(r, lg, lib, state, program, id, ids)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func textinputKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [16]int) *golua.LTable {
	/// @struct TextInputKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop character_forward {struct<tui.KeyBinding>}
	/// @prop character_backward {struct<tui.KeyBinding>}
	/// @prop word_forward {struct<tui.KeyBinding>}
	/// @prop word_backward {struct<tui.KeyBinding>}
	/// @prop delete_word_backward {struct<tui.KeyBinding>}
	/// @prop delete_word_forward {struct<tui.KeyBinding>}
	/// @prop delete_after_cursor {struct<tui.KeyBinding>}
	/// @prop delete_before_cursor {struct<tui.KeyBinding>}
	/// @prop delete_character_backward {struct<tui.KeyBinding>}
	/// @prop delete_character_forward {struct<tui.KeyBinding>}
	/// @prop line_start {struct<tui.KeyBinding>}
	/// @prop line_end {struct<tui.KeyBinding>}
	/// @prop paste {struct<tui.KeyBinding>}
	/// @prop suggestion_accept {struct<tui.KeyBinding>}
	/// @prop suggestion_next {struct<tui.KeyBinding>}
	/// @prop suggestion_prev {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("character_forward", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("character_backward", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("word_forward", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("word_backward", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("delete_word_backward", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("delete_word_forward", tuikeyTable(r, lg, lib, state, program, ids[5]))
	t.RawSetString("delete_after_cursor", tuikeyTable(r, lg, lib, state, program, ids[6]))
	t.RawSetString("delete_before_cursor", tuikeyTable(r, lg, lib, state, program, ids[7]))
	t.RawSetString("delete_character_backward", tuikeyTable(r, lg, lib, state, program, ids[8]))
	t.RawSetString("delete_character_forward", tuikeyTable(r, lg, lib, state, program, ids[9]))
	t.RawSetString("line_start", tuikeyTable(r, lg, lib, state, program, ids[10]))
	t.RawSetString("line_end", tuikeyTable(r, lg, lib, state, program, ids[11]))
	t.RawSetString("paste", tuikeyTable(r, lg, lib, state, program, ids[12]))
	t.RawSetString("suggestion_accept", tuikeyTable(r, lg, lib, state, program, ids[13]))
	t.RawSetString("suggestion_next", tuikeyTable(r, lg, lib, state, program, ids[14]))
	t.RawSetString("suggestion_prev", tuikeyTable(r, lg, lib, state, program, ids[15]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func cursorTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Cursor
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method blink() -> struct<tui.CMDBlink>
	/// @method focus() -> struct<tui.CMDCursorFocus>
	/// @method blur(self) -> self
	/// @method mode() -> int<tui.CursorMode>
	/// @method mode_set(self, mode int<tui.CursorMode>) -> self
	/// @method char_set(self, str string) -> self
	/// @method style() -> struct<lipgloss.Style>
	/// @method style_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_text() -> struct<lipgloss.Style>
	/// @method style_text_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Cursors[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Cursors[id].Blur()
		})

	t.RawSetString("mode", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func filePickerTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct FilePicker
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method did_select_file() -> bool, string
	/// @method did_select_disabled() -> bool, string
	/// @method init() -> struct<tui.CMDFilePickerInit>
	/// @method path() -> string
	/// @method path_set(self, path string) -> self
	/// @method current_directory() -> string
	/// @method current_directory_set(self, dir string) -> self
	/// @method allowed_types() -> []string
	/// @method allowed_types_set(self, types []string) -> self
	/// @method show_perm() -> bool
	/// @method show_perm_set(self, show bool) -> self
	/// @method show_size() -> bool
	/// @method show_size_set(self, show bool) -> self
	/// @method show_hidden() -> bool
	/// @method show_hidden_set(self, show bool) -> self
	/// @method dir_allowed() -> bool
	/// @method dir_allowed_set(self, allowed bool) -> self
	/// @method file_allowed() -> bool
	/// @method file_allowed_set(self, allowed bool) -> self
	/// @method file_selected() -> string
	/// @method file_selected_set(self, file string) -> self
	/// @method height() -> int
	/// @method height_set(self, height int) -> self
	/// @method height_auto() -> bool
	/// @method height_auto_set(self, auto bool) -> self
	/// @method cursor() -> string
	/// @method cursor_set(self, cursor string) -> self
	/// @method keymap() -> struct<tui.FilePickerKeymap>
	/// @method style_cursor() -> struct<lipgloss.Style>
	/// @method style_cursor_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_cursor_disabled() -> struct<lipgloss.Style>
	/// @method style_cursor_disabled_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_symlink() -> struct<lipgloss.Style>
	/// @method style_symlink_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_directory() -> struct<lipgloss.Style>
	/// @method style_directory_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_directory_empty() -> struct<lipgloss.Style>
	/// @method style_directory_empty_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_file() -> struct<lipgloss.Style>
	/// @method style_file_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_file_size() -> struct<lipgloss.Style>
	/// @method style_file_size_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_file_disabled() -> struct<lipgloss.Style>
	/// @method style_file_disabled_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_permission() -> struct<lipgloss.Style>
	/// @method style_permission_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_selected() -> struct<lipgloss.Style>
	/// @method style_selected_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_selected_disabled() -> struct<lipgloss.Style>
	/// @method style_selected_disabled_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.FilePickers[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		did, str := item.FilePickers[id].DidSelectDisabledFile(*item.Msg)

		state.Push(golua.LBool(did))
		state.Push(golua.LString(str))
		return 2
	}))

	t.RawSetString("init", state.NewFunction(func(state *golua.LState) int {
		state.Push(customtea.CMDFilePickerInit(state, int(t.RawGetString("id").(golua.LNumber))))
		return 1
	}))

	t.RawSetString("path", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].Path = args["path"].(string)
		})

	t.RawSetString("current_directory", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].CurrentDirectory = args["dir"].(string)
		})

	t.RawSetString("allowed_types", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowPermissions = args["enabled"].(bool)
		})

	t.RawSetString("show_size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowSize = args["enabled"].(bool)
		})

	t.RawSetString("show_hidden", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].ShowHidden = args["enabled"].(bool)
		})

	t.RawSetString("dir_allowed", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].DirAllowed = args["enabled"].(bool)
		})

	t.RawSetString("file_allowed", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].FileAllowed = args["enabled"].(bool)
		})

	t.RawSetString("file_selected", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].FileSelected = args["file"].(string)
		})

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].Height = args["height"].(int)
		})

	t.RawSetString("height_auto", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.FilePickers[id].AutoHeight = args["enabled"].(bool)
		})

	t.RawSetString("cursor", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := filepickerKeymapTable(r, lg, lib, state, program, id, ids)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func filepickerKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [9]int) *golua.LTable {
	/// @struct FilePickerKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop goto_top {struct<tui.KeyBinding>}
	/// @prop goto_last {struct<tui.KeyBinding>}
	/// @prop down {struct<tui.KeyBinding>}
	/// @prop up {struct<tui.KeyBinding>}
	/// @prop page_up {struct<tui.KeyBinding>}
	/// @prop page_down {struct<tui.KeyBinding>}
	/// @prop back {struct<tui.KeyBinding>}
	/// @prop open {struct<tui.KeyBinding>}
	/// @prop select {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("goto_top", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("goto_last", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("down", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("up", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("page_up", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("page_down", tuikeyTable(r, lg, lib, state, program, ids[5]))
	t.RawSetString("back", tuikeyTable(r, lg, lib, state, program, ids[6]))
	t.RawSetString("open", tuikeyTable(r, lg, lib, state, program, ids[7]))
	t.RawSetString("select", tuikeyTable(r, lg, lib, state, program, ids[8]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func listTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct List
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method cursor() -> int
	/// @method cursor_up(self) -> self
	/// @method cursor_down(self) -> self
	/// @method page_next(self) -> self
	/// @method page_prev(self) -> self
	/// @method pagination_show() -> bool
	/// @method pagination_show_set(self, enabled bool) -> self
	/// @method disable_quit(self) -> self
	/// @method size() -> int, int
	/// @method width() -> int
	/// @method height() -> int
	/// @method size_set(self, width int, height int) -> self
	/// @method width_set(self, width int) -> self
	/// @method height_set(self, height int) -> self
	/// @method filter_state() -> int<tui.FilterState>
	/// @method filter_value() -> string
	/// @method filter_enabled() -> bool
	/// @method filter_enabled_set(self, enabled bool) -> self
	/// @method filter_show() -> bool
	/// @method filter_show_set(self, enabled bool) -> self
	/// @method filter_reset(self) -> self
	/// @method is_filtered() -> bool
	/// @method filter_setting() -> bool
	/// @method filter_func(self, fn int<tui.FilterFunc>) -> self
	/// @method filter_func_custom(self, fn {function(string, []string) -> []struct<tui.ListFilterRank>}) -> self
	/// @method index() -> int
	/// @method items() -> []struct<tui.ListItem>
	/// @method items_visible() -> []struct<tui.ListItem>
	/// @method items_set(self, items []struct<tui.ListItem>) -> self
	/// @method item_insert(self, index int, item struct<tui.ListItem>) -> self
	/// @method item_set(self, index int, item struct<tui.ListItem>) -> self
	/// @method item_remove(self, index int) -> self
	/// @method selected() -> struct<tui.ListItem>
	/// @method select(self, index int) -> self
	/// @method matches(index int) -> []int
	/// @method status_message() -> struct<tui.CMDListStatusMessage>
	/// @method status_message_lifetime() -> int
	/// @method status_message_lifetime_set(self, ms int) -> self
	/// @method statusbar_show() -> bool
	/// @method statusbar_show_set(self, enabled bool) -> self
	/// @method statusbar_item_name() -> string, string
	/// @method statusbar_item_name_set(self, singular string, plural string) -> self
	/// @method title_show() -> bool
	/// @method title_show_set(self, enabled bool) -> self
	/// @method spinner_set(self, from int<tui.SpinnerType>) -> self
	/// @method spinner_set_custom(self, frames []string, fps int) -> self
	/// @method spinner_start() -> struct<tui.CMDListSpinnerStart>
	/// @method spinner_stop(self) -> self
	/// @method spinner_toggle() -> struct<tui.CMDListSpinnerToggle>
	/// @method infinite_scroll() -> bool
	/// @method infinite_scroll_set(self, enabled bool) -> self
	/// @method filter_input() -> struct<tui.TextInput>
	/// @method paginator() -> struct<tui.Paginator>
	/// @method help() -> struct<tui.Help>
	/// @method help_show() -> bool
	/// @method help_show_set(self, enabled bool) -> self
	/// @method keymap() -> struct<tui.ListKeymap>
	/// @method view_help() -> string
	/// @method view_help_short() -> string
	/// @method view_help_full() -> string
	/// @method help_short_additional(self, {function() -> []struct<tui.KeyBinding>}) -> self
	/// @method help_full_additional(self, {function() -> [][]struct<tui.KeyBinding>}) -> self
	/// @method style_titlebar() -> struct<lipgloss.Style>
	/// @method style_titlebar_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_title() -> struct<lipgloss.Style>
	/// @method style_title_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_spinner() -> struct<lipgloss.Style>
	/// @method style_spinner_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_filter_prompt() -> struct<lipgloss.Style>
	/// @method style_filter_prompt_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_filter_cursor() -> struct<lipgloss.Style>
	/// @method style_filter_cursor_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_filter_char_match() -> struct<lipgloss.Style>
	/// @method style_filter_char_match_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_statusbar() -> struct<lipgloss.Style>
	/// @method style_statusbar_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_status_empty() -> struct<lipgloss.Style>
	/// @method style_status_empty_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_statusbar_filter_active() -> struct<lipgloss.Style>
	/// @method style_statusbar_filter_active_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_statusbar_filter_count() -> struct<lipgloss.Style>
	/// @method style_statusbar_filter_count_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_no_items() -> struct<lipgloss.Style>
	/// @method style_no_items_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_help() -> struct<lipgloss.Style>
	/// @method style_help_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_pagination() -> struct<lipgloss.Style>
	/// @method style_pagination_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_pagination_dot_active() -> struct<lipgloss.Style>
	/// @method style_pagination_dot_active_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_pagination_dot_inactive() -> struct<lipgloss.Style>
	/// @method style_pagination_dot_inactive_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_divider_dot() -> struct<lipgloss.Style>
	/// @method style_divider_dot_set(self, style struct<lipgloss.Style>) -> self
	/// @method delegate_set(self, delegate struct<tui.ListDelegate>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Lists[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].CursorUp()
		})

	lib.BuilderFunction(state, t, "cursor_down",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].CursorDown()
		})

	lib.BuilderFunction(state, t, "page_next",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].NextPage()
		})

	lib.BuilderFunction(state, t, "page_prev",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].PrevPage()
		})

	t.RawSetString("pagination_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowPagination(args["enabled"].(bool))
		})

	lib.BuilderFunction(state, t, "disable_quit",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].DisableQuitKeybindings()
		})

	t.RawSetString("size", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		width := item.Lists[id].Width()

		state.Push(golua.LNumber(width))
		return 1
	}))

	t.RawSetString("height", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetHeight(args["height"].(int))
		})

	t.RawSetString("filter_state", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].FilterState()

		state.Push(golua.LNumber(value))
		return 1
	}))

	t.RawSetString("filter_value", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].FilterValue()

		state.Push(golua.LString(value))
		return 1
	}))

	t.RawSetString("filter_enabled", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetFilteringEnabled(args["enabled"].(bool))
		})

	t.RawSetString("filter_show", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowFilter(args["enabled"].(bool))
		})

	lib.BuilderFunction(state, t, "filter_reset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].ResetFilter()
		})

	t.RawSetString("is_filtered", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].IsFiltered()

		state.Push(golua.LBool(value))
		return 1
	}))

	t.RawSetString("filter_setting", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].SettingFilter()

		state.Push(golua.LBool(value))
		return 1
	}))

	lib.BuilderFunction(state, t, "filter_func",
		[]lua.Arg{
			{Type: lua.INT, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			filter := list.DefaultFilter
			if args["fn"].(int) == int(FILTERFUNC_UNSORTED) {
				filter = list.UnsortedFilter
			}

			item.Lists[id].Filter = filter
		})

	lib.BuilderFunction(state, t, "filter_func_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.Lists[id].Filter = func(s1 string, s2 []string) []list.Rank {
				s2t := state.NewTable()
				for i, v := range s2 {
					s2t.RawSetInt(i+1, golua.LString(v))
				}

				state.Push(fn)
				state.Push(golua.LString(s1))
				state.Push(s2t)
				state.Call(2, 1)
				ranks := state.CheckTable(-1)
				state.Pop(1)

				rankList := make([]list.Rank, ranks.Len())

				for i := range ranks.Len() {
					rank := ranks.RawGetInt(i + 1).(*golua.LTable)
					index := rank.RawGetString("index").(golua.LNumber)

					matched := rank.RawGetString("matched").(*golua.LTable)
					matchList := make([]int, matched.Len())
					for z := range matched.Len() {
						m := matched.RawGetInt(z + 1).(golua.LNumber)
						matchList[z] = int(m)
					}

					rankList[i] = list.Rank{
						Index:          int(index),
						MatchedIndexes: matchList,
					}
				}

				return rankList
			}
		})

	t.RawSetString("index", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))

		value := item.Lists[id].Index()

		state.Push(golua.LNumber(value))
		return 1
	}))

	t.RawSetString("items", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].RemoveItem(args["index"].(int))
		})

	t.RawSetString("selected", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].StatusMessageLifetime = time.Duration(args["duration"].(int) * 1e6)
		})

	lib.TableFunction(state, t, "statusbar_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetShowStatusBar(args["enabled"].(bool))
		})

	lib.TableFunction(state, t, "statusbar_item_name",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Lists[id].SetStatusBarItemName(args["singular"].(string), args["plural"].(string))
		})

	lib.TableFunction(state, t, "title_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].FilterInput
			mid := len(item.TextInputs)
			item.TextInputs = append(item.TextInputs, model)

			fi := textinputTable(r, lg, lib, state, program, mid)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].Paginator
			mid := len(item.Paginators)
			item.Paginators = append(item.Paginators, model)

			pg := paginatorTable(r, lg, lib, state, program, mid)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Lists[id].Help
			mid := len(item.Helps)
			item.Helps = append(item.Helps, model)

			hp := helpTable(r, lg, lib, state, program, mid)
			state.Push(hp)
			t.RawSetString("__help", hp)
			return 1
		})

	lib.TableFunction(state, t, "help_show",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := listKeymapTable(r, lg, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	lib.TableFunction(state, t, "view_help",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.Lists[id].Styles.DividerDot = *style.Style
			t.RawSetString("__styleDividerDot", st)
		})

	lib.BuilderFunction(state, t, "delegate_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "delegate"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			did := args["delegate"].(*golua.LTable).RawGetString("id").(golua.LNumber)
			delegate := item.ListDelegates[int(did)]

			item.Lists[id].SetDelegate(delegate)
		})

	return t
}

func listDelegateTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int) *golua.LTable {
	/// @struct ListDelegate
	/// @prop program {int}
	/// @prop id {int}
	/// @method show_description() -> bool
	/// @method show_description_set(self, enabled bool) -> self
	/// @method update_func(self, fn {function(msg struct<tui.MSG>) -> struct<tui.CMD>}) -> self
	/// @method short_help_func(self, fn {function() -> []struct<tui.KeyBinding>}) -> self
	/// @method full_help_func(self, fn {function() -> [][]struct<tui.KeyBinding>}) -> self
	/// @method height() -> int
	/// @method height_set(self, height int) -> self
	/// @method spacing() -> int
	/// @method spacing_set(self, spacing int) -> self
	/// @method style_title_normal() -> struct<tui.Style>
	/// @method style_title_normal_set(self, style struct<tui.Style>) -> self
	/// @method style_title_selected() -> struct<tui.Style>
	/// @method style_title_selected_set(self, style struct<tui.Style>) -> self
	/// @method style_title_dimmed() -> struct<tui.Style>
	/// @method style_title_dimmed_set(self, style struct<tui.Style>) -> self
	/// @method style_desc_normal() -> struct<tui.Style>
	/// @method style_desc_normal_set(self, style struct<tui.Style>) -> self
	/// @method style_desc_selected() -> struct<tui.Style>
	/// @method style_desc_selected_set(self, style struct<tui.Style>) -> self
	/// @method style_desc_dimmed() -> struct<tui.Style>
	/// @method style_desc_dimmed_set(self, style struct<tui.Style>) -> self
	/// @method style_filter_match() -> struct<tui.Style>
	/// @method style_filter_match_set(self, style struct<tui.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "show_description",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ListDelegates[id].ShowDescription

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "show_description_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ListDelegates[id].ShowDescription = args["enabled"].(bool)
		})

	lib.BuilderFunction(state, t, "update_func",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.ListDelegates[id].UpdateFunc = func(msg tea.Msg, _ *list.Model) tea.Cmd {
				luaMsg := customtea.BuildMSG(msg, state)

				state.Push(fn)
				state.Push(luaMsg)
				state.Call(1, 1)
				cmd := state.CheckTable(-1)
				state.Pop(1)

				return customtea.CMDBuild(state, item, cmd)
			}
		})

	lib.BuilderFunction(state, t, "short_help_func",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.ListDelegates[id].ShortHelpFunc = func() []key.Binding {
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

	lib.BuilderFunction(state, t, "full_help_func",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			fn := args["fn"].(*golua.LFunction)

			item.ListDelegates[id].FullHelpFunc = func() [][]key.Binding {
				state.Push(fn)
				state.Call(0, 1)
				groups := state.CheckTable(-1)
				state.Pop(1)

				groupsList := make([][]key.Binding, groups.Len())

				for i := range groups.Len() {
					bindings := groups.RawGetInt(i + 1).(*golua.LTable)
					bindingList := make([]key.Binding, bindings.Len())

					for z := range bindings.Len() {
						b := bindings.RawGetInt(z + 1).(*golua.LTable)
						bid := b.RawGetString("id").(golua.LNumber)

						bindingList[z] = *item.KeyBindings[int(bid)]
					}

					groupsList[i] = bindingList
				}

				return groupsList
			}
		})

	lib.TableFunction(state, t, "height",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ListDelegates[id].Height()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "height_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ListDelegates[id].SetHeight(args["height"].(int))
		})

	lib.TableFunction(state, t, "spacing",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ListDelegates[id].Spacing()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "spacing_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "spacing"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.ListDelegates[id].SetSpacing(args["spacing"].(int))
		})

	t.RawSetString("__styleTitleNormal", golua.LNil)
	lib.TableFunction(state, t, "style_title_normal",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleTitleNormal")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.NormalTitle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleTitleNormal", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_title_normal_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.NormalTitle = *style.Style
			t.RawSetString("__styleTitleNormal", st)
		})

	t.RawSetString("__styleTitleSelected", golua.LNil)
	lib.TableFunction(state, t, "style_title_selected",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleTitleSelected")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.SelectedTitle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleTitleSelected", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_title_selected_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.SelectedTitle = *style.Style
			t.RawSetString("__styleTitleSelected", st)
		})

	t.RawSetString("__styleTitleDimmed", golua.LNil)
	lib.TableFunction(state, t, "style_title_dimmed",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleTitleDimmed")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.DimmedTitle
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleTitleDimmed", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_title_dimmed_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.DimmedTitle = *style.Style
			t.RawSetString("__styleTitleDimmed", st)
		})

	t.RawSetString("__styleDescNormal", golua.LNil)
	lib.TableFunction(state, t, "style_desc_normal",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDescNormal")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.NormalDesc
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDescNormal", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_desc_normal_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.NormalDesc = *style.Style
			t.RawSetString("__styleDescNormal", st)
		})

	t.RawSetString("__styleDescSelected", golua.LNil)
	lib.TableFunction(state, t, "style_desc_selected",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDescSelected")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.SelectedDesc
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDescSelected", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_desc_selected_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.SelectedDesc = *style.Style
			t.RawSetString("__styleDescSelected", st)
		})

	t.RawSetString("__styleDescDimmed", golua.LNil)
	lib.TableFunction(state, t, "style_desc_dimmed",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleDescDimmed")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.DimmedDesc
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleDescDimmed", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_desc_dimmed_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.DimmedDesc = *style.Style
			t.RawSetString("__styleDescDimmed", st)
		})

	t.RawSetString("__styleFilterMatch", golua.LNil)
	lib.TableFunction(state, t, "style_filter_match",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			so := t.RawGetString("__styleFilterMatch")
			if so.Type() == golua.LTTable {
				state.Push(so)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			style := &item.ListDelegates[id].Styles.FilterMatch
			mid := r.CR_LIP.Add(&collection.StyleItem{
				Style: style,
			})

			st := lipglossStyleTable(state, lib, r, mid)
			state.Push(st)
			t.RawSetString("__styleFilterMatch", st)
			return 1
		})

	lib.BuilderFunction(state, t, "style_filter_match_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			st := args["style"].(*golua.LTable)
			styleid := st.RawGetString("id").(golua.LNumber)
			style, _ := r.CR_LIP.Item(int(styleid))

			item.ListDelegates[id].Styles.FilterMatch = *style.Style
			t.RawSetString("__styleFilterMatch", st)
		})

	return t
}

func listKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [14]int) *golua.LTable {
	/// @struct ListKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop cursor_up {struct<tui.KeyBinding>}
	/// @prop cursor_down {struct<tui.KeyBinding>}
	/// @prop page_next {struct<tui.KeyBinding>}
	/// @prop page_prev {struct<tui.KeyBinding>}
	/// @prop goto_start {struct<tui.KeyBinding>}
	/// @prop goto_end {struct<tui.KeyBinding>}
	/// @prop filter {struct<tui.KeyBinding>}
	/// @prop filter_clear {struct<tui.KeyBinding>}
	/// @prop filter_cancel {struct<tui.KeyBinding>}
	/// @prop filter_accept {struct<tui.KeyBinding>}
	/// @prop show_full_help {struct<tui.KeyBinding>}
	/// @prop close_full_help {struct<tui.KeyBinding>}
	/// @prop quit {struct<tui.KeyBinding>}
	/// @prop force_quit {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("cursor_up", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("cursor_down", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("page_next", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("page_prev", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("goto_start", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("goto_end", tuikeyTable(r, lg, lib, state, program, ids[5]))
	t.RawSetString("filter", tuikeyTable(r, lg, lib, state, program, ids[6]))
	t.RawSetString("filter_clear", tuikeyTable(r, lg, lib, state, program, ids[7]))
	t.RawSetString("filter_cancel", tuikeyTable(r, lg, lib, state, program, ids[8]))
	t.RawSetString("filter_accept", tuikeyTable(r, lg, lib, state, program, ids[9]))
	t.RawSetString("show_full_help", tuikeyTable(r, lg, lib, state, program, ids[10]))
	t.RawSetString("close_full_help", tuikeyTable(r, lg, lib, state, program, ids[11]))
	t.RawSetString("quit", tuikeyTable(r, lg, lib, state, program, ids[12]))
	t.RawSetString("force_quit", tuikeyTable(r, lg, lib, state, program, ids[13]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func paginatorTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Paginator
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method slice_bounds(length int) -> int, int
	/// @method page_next(self) -> self
	/// @method page_prev(self) -> self
	/// @method page_items(total int) -> int
	/// @method page_on_first() -> bool
	/// @method page_on_last() -> bool
	/// @method total_pages_set(self, items int) -> self
	/// @method type() -> int<tui.PaginatorType>
	/// @method type_set(self, t int<tui.PaginatorType>) -> self
	/// @method page() -> int
	/// @method page_set(self, p int) -> self
	/// @method page_per() -> int
	/// @method page_per_set(self, p int) -> self
	/// @method page_total() -> int
	/// @method page_total_set(self, p int) -> self
	/// @method format_dot() -> string, string
	/// @method format_dot_set(self, active string, inactive string) -> self
	/// @method format_arabic() -> string
	/// @method format_arabic_set(self, f string) -> self
	/// @method keymap() -> struct<tui.PaginatorKeymap>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Paginators[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Paginators[id].NextPage()
		})

	lib.BuilderFunction(state, t, "page_prev",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := paginatorKeymapTable(r, lg, lib, state, program, id, ids)
			t.RawSetString("__keymap", kmt)
			state.Push(kmt)
			return 1
		})

	return t
}

func paginatorKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [2]int) *golua.LTable {
	/// @struct PaginatorKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop page_prev {struct<tui.KeyBinding>}
	/// @prop page_next {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("page_prev", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("page_next", tuikeyTable(r, lg, lib, state, program, ids[1]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
	/// @struct ProgressOptions
	/// @method width(self, width int) -> self
	/// @method gradient_default(self) -> self
	/// @method gradient_default_scaled(self) -> self
	/// @method gradient(self, colorA string, colorB string) -> self
	/// @method gradient_scaled(self, colorA string, colorB string) -> self
	/// @method solid(self, colorA string) -> self
	/// @method fill_char(self, full int, empty int) -> self
	/// @method spring_options(self, freq float, damp float) -> self

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

func progressTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Progress
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method view_as(percent float) -> string
	/// @method update() -> struct<tui.CMD>
	/// @method percent() -> float
	/// @method percent_set(percent float) -> struct<tui.CMDProgressSet>
	/// @method percent_dec(percent float) -> struct<tui.CMDProgressDec>
	/// @method percent_inc(percent float) -> struct<tui.CMDProgressInc>
	/// @method percent_show() -> bool
	/// @method percent_show_set(self, enabled bool) -> self
	/// @method percent_format() -> string
	/// @method percent_format_set(self, format string) -> self
	/// @method is_animating() -> bool
	/// @method spring_options(self, freq float, damp float) -> self
	/// @method width() -> int
	/// @method width_set(self, width int) -> self
	/// @method full() -> int
	/// @method full_set(self, rune int) -> self
	/// @method full_color() -> string
	/// @method full_color_set(self, color string) -> self
	/// @method empty() -> int
	/// @method empty_set(self, rune int) -> self
	/// @method empty_color() -> string
	/// @method empty_color_set(self, color string) -> self
	/// @method style_percentage() -> struct<lipgloss.Style>
	/// @method style_percentage_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}

			str := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))].ViewAs(args["percent"].(float64))

			state.Push(golua.LString(str))
			return 1
		})

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.ProgressBars[id].PercentFormat

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "percent_format_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func stopwatchTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct StopWatch
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method start() -> struct<tui.CMDStopWatchStart>
	/// @method stop() -> struct<tui.CMDStopWatchStop>
	/// @method toggle() -> struct<tui.CMDStopWatchToggle>
	/// @method reset() -> struct<tui.CMDStopWatchReset>
	/// @method elapsed() -> int
	/// @method running() -> bool
	/// @method interval() -> int
	/// @method interval_set(self, ms int) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StopWatches[id].Interval = time.Duration(args["interval"].(int) * 1e6)
		})

	return t
}

func timerTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Timer
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method init() -> struct<tui.CMDTimerInit>
	/// @method start() -> struct<tui.CMDTimerStart>
	/// @method stop() -> struct<tui.CMDTimerStop>
	/// @method toggle() -> struct<tui.CMDTimerToggle>
	/// @method running() -> bool
	/// @method timed_out() -> bool
	/// @method timeout() -> int
	/// @method timeout_set(self, ms int) -> self
	/// @method interval() -> int
	/// @method interval_set(self, ms int) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Timers[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Timers[id].Interval = time.Duration(args["interval"].(int) * 1e6)
		})

	return t
}

func tableOptionsTable(lib *lua.Lib, state *golua.LState) *golua.LTable {
	/// @struct TableOptions
	/// @method focused(self, focused bool) -> self
	/// @method width(self, width int) -> self
	/// @method height(self, height int) -> self
	/// @method columns(self, cols struct<tui.TableColumn>) -> self
	/// @method rows(self, rows [][]string) -> self
	/// @method styles(self, header struct<lipgloss.Style>, cell struct<lipgloss.Style>, selected struct<lipgloss.Style>) -> self

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
	/// @struct TableColumn
	/// @prop title {string}
	/// @prop width {int}

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

func tuitableTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Table
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method update_viewport(self) -> self
	/// @method focused() -> bool
	/// @method focus(self) -> self
	/// @method blur(self) -> self
	/// @method goto_top(self) -> self
	/// @method goto_bottom(self) -> self
	/// @method move_up(self, n int) -> self
	/// @method move_down(self, n int) -> self
	/// @method cursor() -> int
	/// @method cursor_set(self, n int) -> self
	/// @method columns() -> []struct<tui.TableColumn>
	/// @method rows() -> [][]string
	/// @method columns_set(self, cols []struct<tui.TableColumn>) -> self
	/// @method rows_set(self, rows [][]string) -> self
	/// @method from_values(self, value string, separator string) -> self
	/// @method row_selected() -> []string
	/// @method width() -> int
	/// @method height() -> int
	/// @method width_set(self, width int) -> self
	/// @method height_set(self, height int) -> self
	/// @method keymap() -> struct<tui.TableKeymap>
	/// @method help() -> struct<tui.Help>
	/// @method help_view() -> string
	/// @method styles(self, header struct<lipgloss.Style>, cell struct<lipgloss.Style>, selected struct<lipgloss.Style>) -> self
	/// @method styles_default(self) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Tables[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].Focus()
		})

	lib.BuilderFunction(state, t, "blur",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].Blur()
		})

	lib.BuilderFunction(state, t, "goto_top",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].GotoTop()
		})

	lib.BuilderFunction(state, t, "goto_bottom",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := tableKeymapTable(r, lg, lib, state, program, id, ids)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Tables[id].Help
			mid := len(item.Helps)
			item.Helps = append(item.Helps, model)

			hp := helpTable(r, lg, lib, state, program, mid)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Tables[id].SetStyles(table.DefaultStyles())
		})

	return t
}

func tableKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [8]int) *golua.LTable {
	/// @struct TableKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop line_up {int}
	/// @prop line_down {int}
	/// @prop page_up {int}
	/// @prop page_down {int}
	/// @prop half_page_up {int}
	/// @prop half_page_down {int}
	/// @prop goto_top {int}
	/// @prop goto_bottom {int}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("line_up", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("line_down", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("page_up", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("page_down", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("half_page_up", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("half_page_down", tuikeyTable(r, lg, lib, state, program, ids[5]))
	t.RawSetString("goto_top", tuikeyTable(r, lg, lib, state, program, ids[6]))
	t.RawSetString("goto_bottom", tuikeyTable(r, lg, lib, state, program, ids[7]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func viewportTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Viewport
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tui.CMD>
	/// @method view_up() -> []string
	/// @method view_down() -> []string
	/// @method view_up_half() -> []string
	/// @method view_down_half() -> []string
	/// @method at_top() -> bool
	/// @method at_bottom() -> bool
	/// @method goto_top() -> []string
	/// @method goto_bottom() -> []string
	/// @method line_up() -> []string
	/// @method line_down() -> []string
	/// @method past_bottom() -> bool
	/// @method scroll_percent() -> float
	/// @method width() -> int
	/// @method height() -> int
	/// @method width_set(self, width int) -> self
	/// @method height_set(self, height int) -> self
	/// @method content_set(self, content string) -> self
	/// @method line_count_total() -> int
	/// @method line_count_visible() -> int
	/// @method mouse_wheel_enabled() -> bool
	/// @method mouse_wheel_enabled_set(self, enabled bool) -> self
	/// @method mouse_wheel_delta() -> int
	/// @method mouse_wheel_delta_set(self, delta int) -> self
	/// @method offset_y() -> int
	/// @method offset_y_set(self, offset int) -> self
	/// @method offset_y_set_direct(self, offset int) -> self
	/// @method position_y() -> int
	/// @method position_y_set(self, position int) -> self
	/// @method high_performance() -> bool
	/// @method high_performance_set(self, enabled bool) -> self
	/// @method keymap() -> struct<tui.ViewportKeymap>
	/// @method style() -> struct<lipgloss.Style>
	/// @method style_set(self, style struct<lipgloss.Style>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Viewports[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

			kmt := viewportKeymapTable(r, lg, lib, state, program, id, ids)
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func viewportKeymapTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int, ids [6]int) *golua.LTable {
	/// @struct ViewportKeymap
	/// @prop program {int}
	/// @prop id {int}
	/// @prop page_down {struct<tui.KeyBinding>}
	/// @prop page_up {struct<tui.KeyBinding>}
	/// @prop page_up_half {struct<tui.KeyBinding>}
	/// @prop page_down_half {struct<tui.KeyBinding>}
	/// @prop down {struct<tui.KeyBinding>}
	/// @prop up {struct<tui.KeyBinding>}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("page_down", tuikeyTable(r, lg, lib, state, program, ids[0]))
	t.RawSetString("page_up", tuikeyTable(r, lg, lib, state, program, ids[1]))
	t.RawSetString("page_up_half", tuikeyTable(r, lg, lib, state, program, ids[2]))
	t.RawSetString("page_down_half", tuikeyTable(r, lg, lib, state, program, ids[3]))
	t.RawSetString("down", tuikeyTable(r, lg, lib, state, program, ids[4]))
	t.RawSetString("up", tuikeyTable(r, lg, lib, state, program, ids[5]))

	lib.BuilderFunction(state, t, "default",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func tuicustomTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program, id int) *golua.LTable {
	/// @struct Custom
	/// @prop program {int}
	/// @prop id {int}
	/// @method init() -> struct<tui.CMD>
	/// @method view() -> string
	/// @method update(values []any?) -> struct<tui.CMD>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("init", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
	/// @struct KeyOptions
	/// @method disabled(self, enabled bool) -> self
	/// @method help(self, key string, desc string) -> self
	/// @method keys(self, keys []string) -> self

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

func tuikeyTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct KeyBinding
	/// @prop program {int}
	/// @prop id {int}
	/// @method enabled() -> bool
	/// @method enabled_set(enabled bool)
	/// @method help() -> string, string
	/// @method help_set(key string, desc string)
	/// @method keys() -> []string
	/// @method keys_set(keys []string)
	/// @method unbind(self) -> self

	/// @struct Keymap
	/// @prop program {int}
	/// @prop id {int}
	/// @method default(self) -> self
	/// @method help_short() -> []struct<tui.KeyBinding>
	/// @method help_full() -> [][]struct<tui.KeyBinding>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "enabled",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.KeyBindings[id].Unbind()
		})

	return t
}

func helpTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Help
	/// @prop program {int}
	/// @prop id {int}
	/// @method view(keymap struct<tui.Keymap>) -> string
	/// @method view_help_short(bindings []struct<tui.KeyBinding>) -> string
	/// @method view_help_full(groups [][]struct<tui.KeyBinding>) -> string
	/// @method width() -> int
	/// @method width_set(self, width int) -> self
	/// @method show_all() -> bool
	/// @method show_all_set(self, show_all bool) -> self
	/// @method separator_short() -> string
	/// @method separator_short_set(self, separator string) -> self
	/// @method separator_full() -> string
	/// @method separator_full_set(self, separator string) -> self
	/// @method ellipsis() -> string
	/// @method ellipsis_set(self, ellipsis string) -> self
	/// @method style_ellipsis() -> struct<lipgloss.Style>
	/// @method style_ellipsis_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_short_key() -> struct<lipgloss.Style>
	/// @method style_short_key_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_short_desc() -> struct<lipgloss.Style>
	/// @method style_short_desc_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_short_separator() -> struct<lipgloss.Style>
	/// @method style_short_separator_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_full_key() -> struct<lipgloss.Style>
	/// @method style_full_key_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_full_desc() -> struct<lipgloss.Style>
	/// @method style_full_desc_set(self, style struct<lipgloss.Style>) -> self
	/// @method style_full_separator() -> struct<lipgloss.Style>
	/// @method style_full_separator_set(self, style struct<lipgloss.Style>) -> self

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
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

func tuiimageTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct Image
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tea.CMD>
	/// @method image_string() -> string
	/// @method image_string_set(self, img string) -> self
	/// @method image_file() -> string
	/// @method image_file_set(filename string) -> struct<tea.CMDImageFile>
	/// @method size_set(width int, height int) -> struct<tea.CMDImageSize>
	/// @method is_active() -> bool
	/// @method is_active_set(self, enabled bool) -> self
	/// @method borderless() -> bool
	/// @method borderless_set(self, enabled bool) -> self
	/// @method border_color() -> struct<lipgloss.ColorAdaptive>
	/// @method border_color_set(self, color struct<lipgloss.ColorAny>?) -> self
	/// @method goto_top(self) -> self
	/// @method viewport() -> struct<tui.Viewport>

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.Images[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		im, cmd := item.Images[id].Update(*item.Msg)
		item.Images[id] = &im

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "image_string",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Images[id].ImageString

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "image_string_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "img"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			imgstr := args["img"].(string)
			img := item.Images[id]

			img.ImageString = imgstr
			img.Viewport.SetContent(imgstr)
		})

	lib.TableFunction(state, t, "image_file",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Images[id].FileName

			state.Push(golua.LString(value))
			return 1
		})

	lib.TableFunction(state, t, "image_file_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "filename"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := int(t.RawGetString("id").(golua.LNumber))
			cmd := customtea.CMDImageFile(state, id, args["filename"].(string))

			state.Push(cmd)
			return 1
		})

	lib.TableFunction(state, t, "size_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, args map[string]any) int {
			id := t.RawGetString("id").(golua.LNumber)
			cmd := customtea.CMDImageSize(state, int(id), args["width"].(int), args["height"].(int))

			state.Push(cmd)
			return 1
		})

	lib.TableFunction(state, t, "is_active",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Images[id].Active

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "is_active_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Images[id].Active = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "borderless",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.Images[id].Borderless

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "borderless_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Images[id].Borderless = args["enabled"].(bool)
		})

	lib.TableFunction(state, t, "border_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			color := item.Images[id].BorderColor
			value := lgColorGenericTable(state, color)

			state.Push(value)
			return 1
		})

	lib.BuilderFunction(state, t, "border_color_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color", Optional: true},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			color := args["color"].(*golua.LTable)
			lcolor := lgColorGenericBuild(color)

			if _, ok := lcolor.(lipgloss.AdaptiveColor); !ok {
				lcolor = lipgloss.AdaptiveColor{
					Light: "#000000",
					Dark:  "#FFFFFF",
				}
			}

			item.Images[id].BorderColor = lcolor.(lipgloss.AdaptiveColor)
		})

	lib.BuilderFunction(state, t, "goto_top",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.Images[id].GotoTop()
		})

	t.RawSetString("__viewport", golua.LNil)
	lib.TableFunction(state, t, "viewport",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			ovp := t.RawGetString("__viewport")
			if ovp.Type() == golua.LTTable {
				state.Push(ovp)
				return 1
			}

			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			model := &item.Images[id].Viewport
			mid := len(item.Viewports)
			item.Viewports = append(item.Viewports, model)

			vp := viewportTable(r, lg, lib, state, program, mid)
			state.Push(vp)
			t.RawSetString("__viewport", vp)
			return 1
		})

	return t
}

func statusbarTable(r *lua.Runner, lg *log.Logger, lib *lua.Lib, state *golua.LState, program int, id int) *golua.LTable {
	/// @struct StatusBar
	/// @prop program {int}
	/// @prop id {int}
	/// @method view() -> string
	/// @method update() -> struct<tea.CMD>
	/// @method content() -> string, string, string, string
	/// @method content_set(self, first string, second string, third string, fourth string) -> self
	/// @method colors() -> struct<lipgloss.ColorGeneric>, struct<lipgloss.ColorGeneric>, struct<lipgloss.ColorGeneric>, struct<lipgloss.ColorGeneric>
	/// @method colors_set(self, first_foreground struct<lipgloss.ColorAny>, first_background struct<lipgloss.ColorAny>, second_foreground struct<lipgloss.ColorAny>, second_background struct<lipgloss.ColorAny>, third_foreground struct<lipgloss.ColorAny>, third_background struct<lipgloss.ColorAny>, fourth_foreground struct<lipgloss.ColorAny>, fourth_background struct<lipgloss.ColorAny>) -> self
	/// @method width() -> int
	/// @method width_set(self, width int) -> self
	/// @method height() -> int
	/// @method height_set(self, height int) -> self
	/// @method column_first() -> string
	/// @method column_first_set(self, first string) -> self
	/// @method column_second() -> string
	/// @method column_second_set(self, second string) -> self
	/// @method column_third() -> string
	/// @method column_third_set(self, third string) -> self
	/// @method column_fourth() -> string
	/// @method column_fourth_set(self, fourth string) -> self
	/// @method column_first_colors() -> struct<lipgloss.ColorAdaptive>, struct<lipgloss.ColorAdaptive>
	/// @method column_first_colors_set(self, foreground struct<lipgloss.ColorAny>, background struct<lipgloss.ColorAny>) -> self
	/// @method column_second_colors() -> struct<lipgloss.ColorAdaptive>, struct<lipgloss.ColorAdaptive>
	/// @method column_second_colors_set(self, foreground struct<lipgloss.ColorAny>, background struct<lipgloss.ColorAny>) -> self
	/// @method column_third_colors() -> struct<lipgloss.ColorAdaptive>, struct<lipgloss.ColorAdaptive>
	/// @method column_third_colors_set(self, foreground struct<lipgloss.ColorAny>, background struct<lipgloss.ColorAny>) -> self
	/// @method column_fourth_colors() -> struct<lipgloss.ColorAdaptive>, struct<lipgloss.ColorAdaptive>
	/// @method column_fourth_colors_set(self, foreground struct<lipgloss.ColorAny>, background struct<lipgloss.ColorAny>) -> self

	t := state.NewTable()

	t.RawSetString("program", golua.LNumber(program))
	t.RawSetString("id", golua.LNumber(id))

	t.RawSetString("view", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}

		str := item.StatusBars[int(t.RawGetString("id").(golua.LNumber))].View()

		state.Push(golua.LString(str))
		return 1
	}))

	t.RawSetString("update", state.NewFunction(func(state *golua.LState) int {
		item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
		if err != nil {
			lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
		}
		id := int(t.RawGetString("id").(golua.LNumber))
		sb, cmd := item.StatusBars[id].Update(*item.Msg)
		item.StatusBars[id] = &sb

		var bcmd *golua.LTable

		if cmd == nil {
			bcmd = customtea.CMDNone(state)
		} else {
			bcmd = customtea.CMDStored(state, item, cmd)
		}

		state.Push(bcmd)
		return 1
	}))

	lib.TableFunction(state, t, "content",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			first := item.StatusBars[id].FirstColumn
			second := item.StatusBars[id].SecondColumn
			third := item.StatusBars[id].ThirdColumn
			fourth := item.StatusBars[id].FourthColumn

			state.Push(golua.LString(first))
			state.Push(golua.LString(second))
			state.Push(golua.LString(third))
			state.Push(golua.LString(fourth))
			return 4
		})

	lib.BuilderFunction(state, t, "content_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "first"},
			{Type: lua.STRING, Name: "second"},
			{Type: lua.STRING, Name: "third"},
			{Type: lua.STRING, Name: "fourth"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			sb := item.StatusBars[id]

			sb.FirstColumn = args["first"].(string)
			sb.SecondColumn = args["second"].(string)
			sb.ThirdColumn = args["third"].(string)
			sb.FourthColumn = args["fourth"].(string)
		})

	lib.TableFunction(state, t, "colors",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			first := item.StatusBars[id].FirstColumnColors
			second := item.StatusBars[id].SecondColumnColors
			third := item.StatusBars[id].ThirdColumnColors
			fourth := item.StatusBars[id].FourthColumnColors

			state.Push(lgColorGenericTable(state, first.Foreground))
			state.Push(lgColorGenericTable(state, first.Background))
			state.Push(lgColorGenericTable(state, second.Foreground))
			state.Push(lgColorGenericTable(state, second.Background))
			state.Push(lgColorGenericTable(state, third.Foreground))
			state.Push(lgColorGenericTable(state, third.Background))
			state.Push(lgColorGenericTable(state, fourth.Foreground))
			state.Push(lgColorGenericTable(state, fourth.Background))
			return 8
		})

	lib.BuilderFunction(state, t, "colors_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "first_foreground"},
			{Type: lua.RAW_TABLE, Name: "first_background"},
			{Type: lua.RAW_TABLE, Name: "second_foreground"},
			{Type: lua.RAW_TABLE, Name: "second_background"},
			{Type: lua.RAW_TABLE, Name: "third_foreground"},
			{Type: lua.RAW_TABLE, Name: "third_background"},
			{Type: lua.RAW_TABLE, Name: "fourth_foreground"},
			{Type: lua.RAW_TABLE, Name: "fourth_background"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			firstForeground := lgColorGenericBuild(args["first_foreground"].(*golua.LTable))
			if _, ok := firstForeground.(lipgloss.AdaptiveColor); !ok {
				firstForeground = defaultForeground
			}
			firstBackground := lgColorGenericBuild(args["first_background"].(*golua.LTable))
			if _, ok := firstBackground.(lipgloss.AdaptiveColor); !ok {
				firstBackground = defaultBackground
			}
			secondForeground := lgColorGenericBuild(args["second_foreground"].(*golua.LTable))
			if _, ok := secondForeground.(lipgloss.AdaptiveColor); !ok {
				secondForeground = defaultForeground
			}
			secondBackground := lgColorGenericBuild(args["second_background"].(*golua.LTable))
			if _, ok := secondBackground.(lipgloss.AdaptiveColor); !ok {
				secondBackground = defaultBackground
			}
			thirdForeground := lgColorGenericBuild(args["third_foreground"].(*golua.LTable))
			if _, ok := thirdForeground.(lipgloss.AdaptiveColor); !ok {
				thirdForeground = defaultForeground
			}
			thirdBackground := lgColorGenericBuild(args["third_background"].(*golua.LTable))
			if _, ok := thirdBackground.(lipgloss.AdaptiveColor); !ok {
				thirdBackground = defaultBackground
			}
			fourthForeground := lgColorGenericBuild(args["fourth_foreground"].(*golua.LTable))
			if _, ok := fourthForeground.(lipgloss.AdaptiveColor); !ok {
				fourthForeground = defaultForeground
			}
			fourthBackground := lgColorGenericBuild(args["fourth_background"].(*golua.LTable))
			if _, ok := fourthBackground.(lipgloss.AdaptiveColor); !ok {
				fourthBackground = defaultBackground
			}

			firstPairs := statusbar.ColorConfig{
				Foreground: firstForeground.(lipgloss.AdaptiveColor),
				Background: firstBackground.(lipgloss.AdaptiveColor),
			}
			secondPairs := statusbar.ColorConfig{
				Foreground: secondForeground.(lipgloss.AdaptiveColor),
				Background: secondBackground.(lipgloss.AdaptiveColor),
			}
			thirdPairs := statusbar.ColorConfig{
				Foreground: thirdForeground.(lipgloss.AdaptiveColor),
				Background: thirdBackground.(lipgloss.AdaptiveColor),
			}
			fourthPairs := statusbar.ColorConfig{
				Foreground: fourthForeground.(lipgloss.AdaptiveColor),
				Background: fourthBackground.(lipgloss.AdaptiveColor),
			}

			sb := item.StatusBars[id]

			sb.FirstColumnColors = firstPairs
			sb.SecondColumnColors = secondPairs
			sb.ThirdColumnColors = thirdPairs
			sb.FourthColumnColors = fourthPairs
		})

	lib.TableFunction(state, t, "width",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].Width

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
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].Width = args["width"].(int)
		})

	lib.TableFunction(state, t, "height",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].Height

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "height_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].Height = args["height"].(int)
		})

	lib.TableFunction(state, t, "column_first",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].FirstColumn

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "column_first_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].FirstColumn = args["content"].(string)
		})

	lib.TableFunction(state, t, "column_second",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].SecondColumn

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "column_second_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].SecondColumn = args["content"].(string)
		})

	lib.TableFunction(state, t, "column_third",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].ThirdColumn

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "column_third_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].ThirdColumn = args["content"].(string)
		})

	lib.TableFunction(state, t, "column_fourth",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			value := item.StatusBars[id].FourthColumn

			state.Push(golua.LString(value))
			return 1
		})

	lib.BuilderFunction(state, t, "column_fourth_set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			item.StatusBars[id].FourthColumn = args["content"].(string)
		})

	lib.TableFunction(state, t, "column_first_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			colors := item.StatusBars[id].FirstColumnColors

			state.Push(lgColorGenericTable(state, colors.Foreground))
			state.Push(lgColorGenericTable(state, colors.Background))
			return 2
		})

	lib.BuilderFunction(state, t, "column_first_color_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "foreground"},
			{Type: lua.RAW_TABLE, Name: "background"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			foreground := lgColorGenericBuild(args["foreground"].(*golua.LTable))
			if _, ok := foreground.(lipgloss.AdaptiveColor); !ok {
				foreground = defaultForeground
			}
			background := lgColorGenericBuild(args["background"].(*golua.LTable))
			if _, ok := background.(lipgloss.AdaptiveColor); !ok {
				background = defaultBackground
			}

			item.StatusBars[id].FirstColumnColors = statusbar.ColorConfig{
				Foreground: foreground.(lipgloss.AdaptiveColor),
				Background: background.(lipgloss.AdaptiveColor),
			}
		})

	lib.TableFunction(state, t, "column_second_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			colors := item.StatusBars[id].SecondColumnColors

			state.Push(lgColorGenericTable(state, colors.Foreground))
			state.Push(lgColorGenericTable(state, colors.Background))
			return 2
		})

	lib.BuilderFunction(state, t, "column_second_color_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "foreground"},
			{Type: lua.RAW_TABLE, Name: "background"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			foreground := lgColorGenericBuild(args["foreground"].(*golua.LTable))
			if _, ok := foreground.(lipgloss.AdaptiveColor); !ok {
				foreground = defaultForeground
			}
			background := lgColorGenericBuild(args["background"].(*golua.LTable))
			if _, ok := background.(lipgloss.AdaptiveColor); !ok {
				background = defaultBackground
			}

			item.StatusBars[id].SecondColumnColors = statusbar.ColorConfig{
				Foreground: foreground.(lipgloss.AdaptiveColor),
				Background: background.(lipgloss.AdaptiveColor),
			}
		})

	lib.TableFunction(state, t, "column_third_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			colors := item.StatusBars[id].ThirdColumnColors

			state.Push(lgColorGenericTable(state, colors.Foreground))
			state.Push(lgColorGenericTable(state, colors.Background))
			return 2
		})

	lib.BuilderFunction(state, t, "column_third_color_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "foreground"},
			{Type: lua.RAW_TABLE, Name: "background"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			foreground := lgColorGenericBuild(args["foreground"].(*golua.LTable))
			if _, ok := foreground.(lipgloss.AdaptiveColor); !ok {
				foreground = defaultForeground
			}
			background := lgColorGenericBuild(args["background"].(*golua.LTable))
			if _, ok := background.(lipgloss.AdaptiveColor); !ok {
				background = defaultBackground
			}

			item.StatusBars[id].ThirdColumnColors = statusbar.ColorConfig{
				Foreground: foreground.(lipgloss.AdaptiveColor),
				Background: background.(lipgloss.AdaptiveColor),
			}
		})

	lib.TableFunction(state, t, "column_fourth_color",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			program := int(t.RawGetString("program").(golua.LNumber))
			item, err := r.CR_TEA.Item(program)
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			colors := item.StatusBars[id].FourthColumnColors

			state.Push(lgColorGenericTable(state, colors.Foreground))
			state.Push(lgColorGenericTable(state, colors.Background))
			return 2
		})

	lib.BuilderFunction(state, t, "column_fourth_color_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "foreground"},
			{Type: lua.RAW_TABLE, Name: "background"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, err := r.CR_TEA.Item(int(t.RawGetString("program").(golua.LNumber)))
			if err != nil {
				lua.Error(state, lg.Append(err.Error(), log.LEVEL_ERROR))
			}
			id := int(t.RawGetString("id").(golua.LNumber))

			defaultForeground := lipgloss.AdaptiveColor{
				Light: "#000000",
				Dark:  "#FFFFFF",
			}
			defaultBackground := lipgloss.AdaptiveColor{
				Light: "#FFFFFF",
				Dark:  "#000000",
			}

			foreground := lgColorGenericBuild(args["foreground"].(*golua.LTable))
			if _, ok := foreground.(lipgloss.AdaptiveColor); !ok {
				foreground = defaultForeground
			}
			background := lgColorGenericBuild(args["background"].(*golua.LTable))
			if _, ok := background.(lipgloss.AdaptiveColor); !ok {
				background = defaultBackground
			}

			item.StatusBars[id].FourthColumnColors = statusbar.ColorConfig{
				Foreground: foreground.(lipgloss.AdaptiveColor),
				Background: background.(lipgloss.AdaptiveColor),
			}
		})

	return t
}
