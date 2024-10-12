package states

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MAIN_MENU_OPTION_WORKFLOW int = iota
	MAIN_MENU_OPTION_UTILITIES
	MAIN_MENU_OPTION_EXIT
)

var main_style = lipgloss.NewStyle().Margin(1, 2)

func MainMenu(sm *statemachine.StateMachine) error {
	cli.Clear()

	options := []list.Item{
		main_item{index: MAIN_MENU_OPTION_WORKFLOW, title: "Run Workflow"},
		main_item{index: MAIN_MENU_OPTION_UTILITIES, title: "Utilities"},
		main_item{index: MAIN_MENU_OPTION_EXIT, title: "Exit"},
	}

	result := -1

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	m := main_model{list: list.New(options, delegate, 0, 0), selected: &result}
	m.list.Title = "ImgScal"
	m.list.SetShowStatusBar(false)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run the main menu: %s\n", err)
		sm.SetState(STATE_EXIT)
		return nil
	}

	switch result {
	case MAIN_MENU_OPTION_WORKFLOW:
		sm.SetState(STATE_WORKFLOW_LIST)

	case MAIN_MENU_OPTION_UTILITIES:
		sm.SetState(STATE_UTILITIES)

	case MAIN_MENU_OPTION_EXIT:
		fallthrough
	default:
		sm.SetState(STATE_EXIT)
		cli.Clear()
	}

	return nil
}

type main_item struct {
	title, desc string
	index       int
}

func (i main_item) Title() string       { return i.title }
func (i main_item) Description() string { return i.desc }
func (i main_item) FilterValue() string { return i.title }

type main_model struct {
	list     list.Model
	selected *int
}

func (m main_model) Init() tea.Cmd {
	return nil
}

func (m main_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			result := m.list.SelectedItem()
			*m.selected = result.(main_item).index
			return m, tea.Quit
		case "q":
			fallthrough
		case "ctrl+c":
			*m.selected = -1
			return m, tea.Quit
		}

	case tea.QuitMsg:
		*m.selected = -1

	case tea.WindowSizeMsg:
		h, v := main_style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m main_model) View() string {
	return main_style.Render(m.list.View())
}
