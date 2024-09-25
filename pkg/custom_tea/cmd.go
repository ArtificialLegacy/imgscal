package customtea

import (
	"time"

	teamodels "github.com/ArtificialLegacy/imgscal/pkg/custom_tea/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

func CMDBuild(state *golua.LState, item *teamodels.TeaItem, t *golua.LTable) tea.Cmd {
	cmdValue := t.RawGetString("cmd")
	cmdId := 0
	if cmdValue.Type() == golua.LTNumber {
		cmdId = int(cmdValue.(golua.LNumber))
	}
	cmd := cmdList[cmdId](item, state, t)

	return cmd
}

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
	CMD_STOPWATCHSTART
	CMD_STOPWATCHSTOP
	CMD_STOPWATCHTOGGLE
	CMD_STOPWATCHRESET
	CMD_TIMERINIT
	CMD_TIMERSTART
	CMD_TIMERSTOP
	CMD_TIMERTOGGLE
	CMD_VIEWPORTSYNC
	CMD_VIEWPORTUP
	CMD_VIEWPORTDOWN
	CMD_PRINTF
	CMD_PRINTLN
	CMD_WINDOWTITLE
	CMD_WINDOWSIZE
	CMD_SUSPEND
	CMD_QUIT
	CMD_SHOWCURSOR
	CMD_HIDECURSOR
	CMD_CLEARSCREEN
	CMD_CLEARSCROLLAREA
	CMD_SCROLLSYNC
	CMD_SCROLLUP
	CMD_SCROLLDOWN
	CMD_EVERY
	CMD_TICK
	CMD_TOGGLEREPORTFOCUS
	CMD_TOGGLEBRACKETEDPASTE
	CMD_DISABLEMOUSE
	CMD_ENABLEMOUSEALLMOTION
	CMD_ENABLEMOUSECELLMOTION
	CMD_ENTERALTSCREEN
	CMD_EXITALTSCREEN
)

type CMDBuilder func(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd

var cmdList = []CMDBuilder{}

func init() {
	cmdList = []CMDBuilder{
		func(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd { return nil },
		CMDStoredBuild,
		CMDBatchBuild,
		CMDSequenceBuild,
		CMDSpinnerTickBuild,
		CMDTextAreaFocusBuild,
		CMDTextInputFocusBuild,
		CMDBlinkBuild,
		CMDCursorFocusBuild,
		CMDFilePickerInitBuild,
		CMDListSetItemsBuild,
		CMDListInsertItemBuild,
		CMDListSetItemBuild,
		CMDListStatusMessageBuild,
		CMDListSpinnerStartBuild,
		CMDListSpinnerToggleBuild,
		CMDProgressSetBuild,
		CMDProgressDecBuild,
		CMDProgressIncBuild,
		CMDStopWatchStartBuild,
		CMDStopWatchStopBuild,
		CMDStopWatchToggleBuild,
		CMDStopWatchResetBuild,
		CMDTimerInitBuild,
		CMDTimerStartBuild,
		CMDTimerStopBuild,
		CMDTimerToggleBuild,
		CMDViewportSyncBuild,
		CMDViewportUpBuild,
		CMDViewportDownBuild,
		CMDPrintfBuild,
		CMDPrintlnBuild,
		CMDWindowTitleBuild,
		CMDWindowSizeBuild,
		CMDSuspendBuild,
		CMDQuitBuild,
		CMDShowCursorBuild,
		CMDHideCursorBuild,
		CMDClearScreenBuild,
		CMDClearScrollAreaBuild,
		CMDScrollSyncBuild,
		CMDScrollUpBuild,
		CMDScrollDownBuild,
		CMDEveryBuild,
		CMDTickBuild,
		CMDToggleReportFocusBuild,
		CMDToggleBracketedPasteBuild,
		CMDDisableMouseBuild,
		CMDEnableMouseAllMotionBuild,
		CMDEnableMouseCellMotionBuild,
		CMDEnterAltScreenBuild,
		CMDExitAltScreenBuild,
	}
}

func CMDNone(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_NONE))

	return t
}

func CMDStored(state *golua.LState, item *teamodels.TeaItem, cmd tea.Cmd) *golua.LTable {
	t := state.NewTable()

	id := len(item.Cmds)
	item.Cmds = append(item.Cmds, cmd)

	t.RawSetString("cmd", golua.LNumber(CMD_STORED))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDStoredBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmd := item.Cmds[int(t.RawGetString("id").(golua.LNumber))]

	return cmd
}

func CMDBatch(state *golua.LState, cmds *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_BATCH))
	t.RawSetString("cmds", cmds)

	return t
}

func CMDBatchBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmdList := t.RawGetString("cmds").(*golua.LTable)
	cmds := make([]tea.Cmd, cmdList.Len())

	for i := range cmdList.Len() {
		bcmd := CMDBuild(state, item, cmdList.RawGetInt(i+1).(*golua.LTable))
		cmds[i] = bcmd
	}

	return tea.Batch(cmds...)
}

func CMDSequence(state *golua.LState, cmds *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SEQUENCE))
	t.RawSetString("cmds", cmds)

	return t
}

func CMDSequenceBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	cmdList := t.RawGetString("cmds").(*golua.LTable)
	cmds := make([]tea.Cmd, cmdList.Len())

	for i := range cmdList.Len() {
		bcmd := CMDBuild(state, item, cmdList.RawGetInt(i+1).(*golua.LTable))
		cmds[i] = bcmd
	}

	return tea.Sequence(cmds...)
}

func CMDSpinnerTick(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SPINNERTICK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDSpinnerTickBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	id := int(t.RawGetString("id").(golua.LNumber))
	return item.Spinners[id].Tick
}

func CMDTextAreaFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TEXTAREAFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTextInputFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TEXTINPUTFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTextAreaFocusBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	textArea := item.TextAreas[int(t.RawGetString("id").(golua.LNumber))]

	return textArea.Focus()
}

func CMDTextInputFocusBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	textInput := item.TextInputs[int(t.RawGetString("id").(golua.LNumber))]

	return textInput.Focus()
}

func CMDBlink(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_BLINK))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDBlinkBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	cursor := item.Cursors[int(t.RawGetString("id").(golua.LNumber))]

	return cursor.BlinkCmd()
}

func CMDCursorFocus(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_CURSORFOCUS))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDCursorFocusBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	cursor := item.Cursors[int(t.RawGetString("id").(golua.LNumber))]

	return cursor.Focus()
}

func CMDFilePickerInit(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_FILEPICKERINIT))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDFilePickerInitBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	fp := item.FilePickers[int(t.RawGetString("id").(golua.LNumber))]

	return fp.Init()
}

func CMDListSetItems(state *golua.LState, id int, items *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSETITEMS))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("items", items)

	return t
}

func CMDListSetItemsBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	itemList := t.RawGetString("items").(*golua.LTable)
	items := make([]list.Item, itemList.Len())

	for i := range itemList.Len() {
		items[i] = ListItemBuild(itemList.RawGetInt(i + 1).(*golua.LTable))
	}

	return li.SetItems(items)
}

func CMDListInsertItem(state *golua.LState, id, index int, item *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTINSERTITEM))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("index", golua.LNumber(index))
	t.RawSetString("item", item)

	return t
}

func CMDListInsertItemBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	it := t.RawGetString("item").(*golua.LTable)
	index := int(t.RawGetString("index").(golua.LNumber))

	return li.InsertItem(index, ListItemBuild(it))
}

func CMDListSetItem(state *golua.LState, id, index int, item *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSETITEM))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("index", golua.LNumber(index))
	t.RawSetString("item", item)

	return t
}

func CMDListSetItemBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	it := t.RawGetString("item").(*golua.LTable)
	index := int(t.RawGetString("index").(golua.LNumber))

	return li.SetItem(index, ListItemBuild(it))
}

func CMDListStatusMessage(state *golua.LState, id int, msg string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSTATUSMESSAGE))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("msg", golua.LString(msg))

	return t
}

func CMDListStatusMessageBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	msg := string(t.RawGetString("msg").(golua.LString))
	return li.NewStatusMessage(msg)
}

func CMDListSpinnerStart(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSPINNERSTART))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDListSpinnerStartBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	return li.StartSpinner()
}

func CMDListSpinnerToggle(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_LISTSPINNERTOGGLE))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDListSpinnerToggleBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.Lists[int(t.RawGetString("id").(golua.LNumber))]

	return li.ToggleSpinner()
}

func CMDProgressSet(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSSET))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func CMDProgressSetBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.SetPercent(percent)
}

func CMDProgressDec(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSDEC))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func CMDProgressDecBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.DecrPercent(percent)
}

func CMDProgressInc(state *golua.LState, id int, percent float64) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PROGRESSINC))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func CMDProgressIncBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	li := item.ProgressBars[int(t.RawGetString("id").(golua.LNumber))]

	percent := float64(t.RawGetString("percent").(golua.LNumber))
	return li.IncrPercent(percent)
}

func CMDStopWatchStart(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_STOPWATCHSTART))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDStopWatchStartBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	sw := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))]

	return sw.Start()
}

func CMDStopWatchStop(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_STOPWATCHSTOP))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDStopWatchStopBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	sw := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))]

	return sw.Stop()
}

func CMDStopWatchReset(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_STOPWATCHRESET))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDStopWatchResetBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	sw := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))]

	return sw.Reset()
}

func CMDStopWatchToggle(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_STOPWATCHTOGGLE))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDStopWatchToggleBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	sw := item.StopWatches[int(t.RawGetString("id").(golua.LNumber))]

	return sw.Toggle()
}

func CMDTimerStart(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TIMERSTART))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTimerStartBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	ti := item.Timers[int(t.RawGetString("id").(golua.LNumber))]

	return ti.Start()
}

func CMDTimerInit(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TIMERINIT))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTimerInitBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	ti := item.Timers[int(t.RawGetString("id").(golua.LNumber))]

	return ti.Init()
}

func CMDTimerStop(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TIMERSTOP))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTimerStopBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	ti := item.Timers[int(t.RawGetString("id").(golua.LNumber))]

	return ti.Stop()
}

func CMDTimerToggle(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TIMERTOGGLE))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDTimerToggleBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	ti := item.Timers[int(t.RawGetString("id").(golua.LNumber))]

	return ti.Toggle()
}

func CMDViewportSync(state *golua.LState, id int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_VIEWPORTSYNC))
	t.RawSetString("id", golua.LNumber(id))

	return t
}

func CMDViewportSyncBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	vp := item.Viewports[int(t.RawGetString("id").(golua.LNumber))]

	return viewport.Sync(*vp)
}

func CMDViewportUp(state *golua.LState, id int, lines *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_VIEWPORTUP))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("lines", lines)

	return t
}

func CMDViewportUpBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	vp := item.Viewports[int(t.RawGetString("id").(golua.LNumber))]
	lineList := t.RawGetString("lines").(*golua.LTable)
	lines := make([]string, lineList.Len())

	for i := range lineList.Len() {
		l := lineList.RawGetInt(i + 1).(golua.LString)
		lines[i] = string(l)
	}

	return viewport.ViewUp(*vp, lines)
}

func CMDViewportDown(state *golua.LState, id int, lines *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_VIEWPORTDOWN))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("lines", lines)

	return t
}

func CMDViewportDownBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	vp := item.Viewports[int(t.RawGetString("id").(golua.LNumber))]
	lineList := t.RawGetString("lines").(*golua.LTable)
	lines := make([]string, lineList.Len())

	for i := range lineList.Len() {
		l := lineList.RawGetInt(i + 1).(golua.LString)
		lines[i] = string(l)
	}

	return viewport.ViewDown(*vp, lines)
}

func CMDPrintf(state *golua.LState, format string, args *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PRINTF))
	t.RawSetString("format", golua.LString(format))
	t.RawSetString("args", args)

	return t
}

func CMDPrintfBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	format := t.RawGetString("format").(golua.LString)
	args := t.RawGetString("args").(*golua.LTable)

	values := make([]any, args.Len())
	for i := range args.Len() {
		v := args.RawGetInt(i + 1)

		switch v := v.(type) {
		case golua.LString:
			values[i] = string(v)
		case golua.LNumber:
			if float64(v) == float64(int(v)) {
				values[i] = int(v)
			} else {
				values[i] = float64(v)
			}
		case golua.LBool:
			values[i] = bool(v)
		default:
			values[i] = v.String()
		}
	}

	return tea.Printf(string(format), values...)
}

func CMDPrintln(state *golua.LState, args *golua.LTable) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_PRINTLN))
	t.RawSetString("args", args)

	return t
}

func CMDPrintlnBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	args := t.RawGetString("args").(*golua.LTable)

	values := make([]any, args.Len())
	for i := range args.Len() {
		v := args.RawGetInt(i + 1)

		switch v := v.(type) {
		case golua.LString:
			values[i] = string(v)
		case golua.LNumber:
			if float64(v) == float64(int(v)) {
				values[i] = int(v)
			} else {
				values[i] = float64(v)
			}
		case golua.LBool:
			values[i] = bool(v)
		default:
			values[i] = v.String()
		}
	}

	return tea.Println(values...)
}

func CMDWindowTitle(state *golua.LState, title string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_WINDOWTITLE))
	t.RawSetString("title", golua.LString(title))

	return t
}

func CMDWindowTitleBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	title := t.RawGetString("title").(golua.LString)

	return tea.SetWindowTitle(string(title))
}

func CMDWindowSize(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_WINDOWSIZE))

	return t
}

func CMDWindowSizeBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.WindowSize()
}

func CMDSuspend(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SUSPEND))

	return t
}

func CMDSuspendBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.Suspend
}

func CMDQuit(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_QUIT))

	return t
}

func CMDQuitBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.Quit
}

func CMDShowCursor(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SHOWCURSOR))

	return t
}

func CMDShowCursorBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.ShowCursor
}

func CMDHideCursor(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_HIDECURSOR))

	return t
}

func CMDHideCursorBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.HideCursor
}

func CMDClearScreen(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_CLEARSCREEN))

	return t
}

func CMDClearScreenBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.ClearScreen
}

func CMDClearScrollArea(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_CLEARSCROLLAREA))

	return t
}

func CMDClearScrollAreaBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.ClearScrollArea
}

func CMDScrollSync(state *golua.LState, lines *golua.LTable, top, bottom int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SCROLLSYNC))
	t.RawSetString("lines", lines)
	t.RawSetString("top", golua.LNumber(top))
	t.RawSetString("bottom", golua.LNumber(bottom))

	return t
}

func CMDScrollSyncBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	lineList := t.RawGetString("lines").(*golua.LTable)
	lines := make([]string, lineList.Len())

	for i := range lineList.Len() {
		l := lineList.RawGetInt(i + 1).(golua.LString)
		lines[i] = string(l)
	}

	top := int(t.RawGetString("top").(golua.LNumber))
	bottom := int(t.RawGetString("bottom").(golua.LNumber))

	return tea.SyncScrollArea(lines, top, bottom)
}

func CMDScrollUp(state *golua.LState, lines *golua.LTable, top, bottom int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SCROLLUP))
	t.RawSetString("lines", lines)
	t.RawSetString("top", golua.LNumber(top))
	t.RawSetString("bottom", golua.LNumber(bottom))

	return t
}

func CMDScrollUpBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	lineList := t.RawGetString("lines").(*golua.LTable)
	lines := make([]string, lineList.Len())

	for i := range lineList.Len() {
		l := lineList.RawGetInt(i + 1).(golua.LString)
		lines[i] = string(l)
	}

	top := int(t.RawGetString("top").(golua.LNumber))
	bottom := int(t.RawGetString("bottom").(golua.LNumber))

	return tea.ScrollUp(lines, top, bottom)
}

func CMDScrollDown(state *golua.LState, lines *golua.LTable, top, bottom int) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_SCROLLDOWN))
	t.RawSetString("lines", lines)
	t.RawSetString("top", golua.LNumber(top))
	t.RawSetString("bottom", golua.LNumber(bottom))

	return t
}

func CMDScrollDownBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	lineList := t.RawGetString("lines").(*golua.LTable)
	lines := make([]string, lineList.Len())

	for i := range lineList.Len() {
		l := lineList.RawGetInt(i + 1).(golua.LString)
		lines[i] = string(l)
	}

	top := int(t.RawGetString("top").(golua.LNumber))
	bottom := int(t.RawGetString("bottom").(golua.LNumber))

	return tea.ScrollDown(lines, top, bottom)
}

func CMDEvery(state *golua.LState, duration int, fn *golua.LFunction) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_EVERY))
	t.RawSetString("duration", golua.LNumber(duration))
	t.RawSetString("fn", fn)

	return t
}

func CMDEveryBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	duration := int(t.RawGetString("duration").(golua.LNumber))
	fn := t.RawGetString("fn").(*golua.LFunction)

	return tea.Every(time.Duration(duration)*time.Millisecond, func(tm time.Time) tea.Msg {
		ms := tm.UnixMilli()

		state.Push(fn)
		state.Push(golua.LNumber(ms))
		state.Call(1, 1)
		v := state.CheckAny(-1)
		state.Pop(1)

		return tea.Msg(v)
	})
}

func CMDTick(state *golua.LState, duration int, fn *golua.LFunction) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TICK))
	t.RawSetString("duration", golua.LNumber(duration))
	t.RawSetString("fn", fn)

	return t
}

func CMDTickBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	duration := int(t.RawGetString("duration").(golua.LNumber))
	fn := t.RawGetString("fn").(*golua.LFunction)

	return tea.Tick(time.Duration(duration)*time.Millisecond, func(tm time.Time) tea.Msg {
		ms := tm.UnixMilli()

		state.Push(fn)
		state.Push(golua.LNumber(ms))
		state.Call(1, 1)
		v := state.CheckAny(-1)
		state.Pop(1)

		return tea.Msg(v)
	})
}

func CMDToggleReportFocus(state *golua.LState, enable bool) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TOGGLEREPORTFOCUS))
	t.RawSetString("enable", golua.LBool(enable))

	return t
}

func CMDToggleReportFocusBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	if bool(t.RawGetString("enable").(golua.LBool)) {
		return tea.EnableReportFocus
	}
	return tea.DisableReportFocus
}

func CMDToggleBracketedPaste(state *golua.LState, enable bool) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_TOGGLEBRACKETEDPASTE))
	t.RawSetString("enable", golua.LBool(enable))

	return t
}

func CMDToggleBracketedPasteBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	if bool(t.RawGetString("enable").(golua.LBool)) {
		return tea.EnableBracketedPaste
	}
	return tea.DisableBracketedPaste
}

func CMDDisableMouse(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_DISABLEMOUSE))

	return t
}

func CMDDisableMouseBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.DisableMouse
}

func CMDEnableMouseAllMotion(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_ENABLEMOUSEALLMOTION))

	return t
}

func CMDEnableMouseAllMotionBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.EnableMouseAllMotion
}

func CMDEnableMouseCellMotion(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_ENABLEMOUSECELLMOTION))

	return t
}

func CMDEnableMouseCellMotionBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.EnableMouseCellMotion
}

func CMDEnterAltScreen(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_ENTERALTSCREEN))

	return t
}

func CMDEnterAltScreenBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.EnterAltScreen
}

func CMDExitAltScreen(state *golua.LState) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("cmd", golua.LNumber(CMD_EXITALTSCREEN))

	return t
}

func CMDExitAltScreenBuild(item *teamodels.TeaItem, state *golua.LState, t *golua.LTable) tea.Cmd {
	return tea.EnterAltScreen
}
