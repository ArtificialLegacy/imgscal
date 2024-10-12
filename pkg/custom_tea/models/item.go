package teamodels

import (
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
	teaimage "github.com/mistakenelf/teacup/image"
	"github.com/mistakenelf/teacup/statusbar"
	golua "github.com/yuin/gopher-lua"
)

type TeaItem struct {
	LuaModel *golua.LTable

	FnInit   *golua.LFunction
	FnUpdate *golua.LFunction
	FnView   *golua.LFunction

	Msg  *tea.Msg
	Cmds []tea.Cmd

	KeyBindings []*key.Binding

	Spinners      map[int]*spinner.Model
	TextAreas     []*textarea.Model
	TextInputs    []*textinput.Model
	Cursors       []*cursor.Model
	FilePickers   []*filepicker.Model
	Lists         []*list.Model
	ListDelegates []*list.DefaultDelegate
	Paginators    []*paginator.Model
	ProgressBars  []*progress.Model
	StopWatches   map[int]*stopwatch.Model
	Timers        map[int]*timer.Model
	Tables        []*table.Model
	Viewports     []*viewport.Model
	Customs       []*CustomModel
	Helps         []*help.Model
	Images        []*teaimage.Model
	StatusBars    []*statusbar.Model
}
