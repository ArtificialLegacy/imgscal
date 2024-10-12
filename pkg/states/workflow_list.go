package states

import (
	"fmt"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var wflist_style = lipgloss.NewStyle().Margin(1, 2)

func WorkflowList(sm *statemachine.StateMachine) error {
	cli.Clear()

	workflows, errList, err := workflow.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		return err
	}

	if errList != nil && len(*errList) > 0 {
		lg := log.NewLoggerBase("scan", sm.Config.LogDirectory, false)
		defer lg.Close()
		lg.Append("Encountered errors while scanning for workflows: ", log.LEVEL_SYSTEM)
		for _, e := range *errList {
			lg.Append(e.Error(), log.LEVEL_ERROR)
		}
	}

	if len(*workflows) == 0 {
		fmt.Printf("\nWorkflow directory empty, nothing to run.\n")
		fmt.Printf("%s%s%s\n\n", configPathColor, sm.Config.WorkflowDirectory, cli.COLOR_RESET)

		fmt.Printf(" > Try \u001b[48;5;234mmake install-examples%s\n\n", cli.COLOR_RESET)

		cli.Question("Press any key to continue...", cli.QuestionOptions{})
		sm.SetState(STATE_MAIN)
		return nil
	}

	options := []string{}
	optionsWorkflows := []*workflow.Workflow{}
	optionsPaths := []string{}
	listOptions := []list.Item{}

	i := 0
	for _, w := range *workflows {
		starUsed := false
		for s, ws := range w.Workflows {
			optName := ""
			if s == "*" {
				if starUsed {
					continue
				}
				starUsed = true
				optName = w.Name
			} else {
				optName = w.Name + "/" + s
			}
			options = append(options, optName)
			optionsWorkflows = append(optionsWorkflows, w)
			optionsPaths = append(optionsPaths, path.Join(path.Dir(w.Base), ws))
			listOptions = append(listOptions, wflist_item{index: i, title: fmt.Sprintf("%s (%s)", optName, w.Version), desc: fmt.Sprintf("Author: %s", w.Author)})
			i++
		}
	}

	result := 0

	m := wflist_model{list: list.New(listOptions, list.NewDefaultDelegate(), 0, 0), selected: &result}
	m.list.Title = "Select Workflow:"

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run workflow selection list: %s\n", err)
		sm.SetState(STATE_EXIT)
		return nil
	}

	if result == -1 {
		sm.SetState(STATE_MAIN)
	} else {
		WorkflowConfirmEnter(sm, WorkflowConfirmData{
			Workflow: optionsWorkflows[result],
			Entry:    optionsPaths[result],
			Name:     options[result],
		})
	}

	return nil
}

type wflist_item struct {
	title, desc string
	index       int
}

func (i wflist_item) Title() string       { return i.title }
func (i wflist_item) Description() string { return i.desc }
func (i wflist_item) FilterValue() string { return i.title }

type wflist_model struct {
	list     list.Model
	selected *int
}

func (m wflist_model) Init() tea.Cmd {
	return nil
}

func (m wflist_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			result := m.list.SelectedItem()
			*m.selected = result.(wflist_item).index
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
		h, v := wflist_style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m wflist_model) View() string {
	return wflist_style.Render(m.list.View())
}
