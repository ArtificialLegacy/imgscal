package customtea

import (
	teamodels "github.com/ArtificialLegacy/imgscal/pkg/custom_tea/models"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	tea "github.com/charmbracelet/bubbletea"
	golua "github.com/yuin/gopher-lua"
)

type ProgramModel struct {
	Item  *teamodels.TeaItem
	State *golua.LState
	Id    int
	R     *lua.Runner
	Lg    *log.Logger
}

func (m ProgramModel) Init() tea.Cmd {
	m.State.Push(m.Item.FnInit)
	m.State.Push(golua.LNumber(m.Id))
	m.State.Call(1, 2)
	model := m.State.CheckTable(-2)
	cmd := m.State.CheckTable(-1)
	m.State.Pop(2)

	bcmd := CMDBuild(m.State, m.Item, cmd)

	m.Item.LuaModel = model

	return bcmd
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Item.Msg = &msg
	defer func() {
		m.Item.Msg = nil
	}()
	var bcmd tea.Cmd

	m.Item.Cmds = []tea.Cmd{}

	// Program should always be exittable
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	luaMsg := BuildMSG(msg, m.State)

	m.State.Push(m.Item.FnUpdate)
	m.State.Push(m.Item.LuaModel)
	m.State.Push(luaMsg)

	m.State.Call(2, 1)
	cmd := m.State.OptTable(-1, CMDNone(m.State))
	m.State.Pop(1)

	bcmd = CMDBuild(m.State, m.Item, cmd)

	return m, bcmd
}

func (m ProgramModel) View() string {
	m.State.Push(m.Item.FnView)
	m.State.Push(m.Item.LuaModel)
	m.State.Call(1, 1)
	str := m.State.CheckString(-1)
	m.State.Pop(1)

	return str
}
