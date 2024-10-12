package teamodels

import (
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

type (
	cmdBuilder func(state *golua.LState, item *TeaItem, t *golua.LTable) tea.Cmd
	msgBuilder func(msg tea.Msg, state *golua.LState) *golua.LTable
)

type CustomData struct {
	Data *golua.LTable
}

func newCustomData(data *golua.LTable) *CustomData {
	return &CustomData{
		Data: data,
	}
}

// Data is double nested so the pointer can be modified even though the recievers are not pointers
type CustomModel struct {
	Program  int
	Data     *CustomData
	InitFn   *golua.LFunction
	UpdateFn *golua.LFunction
	ViewFn   *golua.LFunction
	State    *golua.LState
	Item     *TeaItem
	CMDBuild cmdBuilder
	MSGBuild msgBuilder
}

func NewCustomModel(program int, init, update, view *golua.LFunction, state *golua.LState, item *TeaItem, cmd cmdBuilder, msg msgBuilder) CustomModel {
	return CustomModel{
		Program:  program,
		Data:     newCustomData(state.NewTable()),
		InitFn:   init,
		UpdateFn: update,
		ViewFn:   view,
		State:    state,
		Item:     item,
		CMDBuild: cmd,
		MSGBuild: msg,
	}
}

func (m CustomModel) Init() tea.Cmd {
	m.State.Push(m.InitFn)
	m.State.Push(golua.LNumber(m.Program))
	m.State.Call(1, 2)
	data := m.State.CheckTable(-2)
	cmd := m.State.CheckTable(-1)
	m.State.Pop(2)

	m.Data.Data = data

	bcmd := m.CMDBuild(m.State, m.Item, cmd)
	return bcmd
}

type CustomMSG struct {
	Original tea.Msg
	Values   []any
}

func (m CustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmsg CustomMSG
	if c, ok := msg.(CustomMSG); ok {
		cmsg = c
	} else {
		cmsg = CustomMSG{
			Original: msg,
		}
	}
	luaMsg := m.MSGBuild(cmsg.Original, m.State)

	m.State.Push(m.UpdateFn)
	m.State.Push(m.Data.Data)
	m.State.Push(luaMsg)

	additional := len(cmsg.Values)
	for _, v := range cmsg.Values {
		m.State.Push(v.(golua.LValue))
	}

	m.State.Call(2+additional, 1)
	cmd := m.State.CheckTable(-1)
	m.State.Pop(1)

	bcmd := m.CMDBuild(m.State, m.Item, cmd)

	return m, bcmd
}

func (m CustomModel) View() string {
	m.State.Push(m.ViewFn)
	m.State.Push(m.Data.Data)
	m.State.Call(1, 1)
	str := m.State.CheckString(-1)
	m.State.Pop(1)

	return str
}
