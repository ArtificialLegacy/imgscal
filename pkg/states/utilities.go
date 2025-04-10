package states

import (
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/config"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	UTILITIES_OPTION_CONFIG int = iota
	UTILITIES_OPTION_RMLogs
	UTILITIES_OPTION_BACK
)

var utilities_style = lipgloss.NewStyle().MarginTop(2).MarginBottom(2)

func Utilities(sm *statemachine.StateMachine) error {
	cli.Clear()

	result := -1

	options := []list.Item{
		utilities_item{index: UTILITIES_OPTION_CONFIG, title: "View Config"},
		utilities_item{index: UTILITIES_OPTION_RMLogs, title: "Delete Logs"},
		utilities_item{index: UTILITIES_OPTION_BACK, title: "Back"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	m := utilities_model{list: list.New(options, delegate, 0, 0), selected: &result}
	m.list.Title = "ImgScal Utilities"
	m.list.SetShowStatusBar(false)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run the utilities menu: %s\n", err)
		sm.SetState(STATE_EXIT)
		return nil
	}

	switch result {
	case UTILITIES_OPTION_CONFIG:
		viewConfig(sm.Config)

	case UTILITIES_OPTION_RMLogs:
		removeLogs(sm.Config)

	case UTILITIES_OPTION_BACK:
		fallthrough
	default:
		sm.SetState(STATE_MAIN)
	}

	return nil
}

const (
	configPathColor  = "\u001b[38;5;240m\u001b[4m"
	configTrueColor  = "\u001b[38;5;195m"
	configFalseColor = "\u001b[38;5;209m"
)

func viewConfig(cfg *config.Config) {
	cli.Clear()

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")

	dlColor := configFalseColor
	if cfg.DisableLogs {
		dlColor = configTrueColor
	}
	acColor := configFalseColor
	if cfg.AlwaysConfirm {
		acColor = configTrueColor
	}

	fmt.Printf("\n\n  %sImgScal Config%s v%s\n", cli.COLOR_BOLD, cli.COLOR_RESET, cfg.ConfigVersion)
	fmt.Printf("  %s%s%s\n\n", configPathColor, cfgPath, cli.COLOR_RESET)

	cfgFields := []string{
		"config_directory",
		"workflow_directory",
		"log_directory",
		"output_directory",
		"input_directory",
		"plugin_directory",
		"disable_logs",
		"always_confirm",
		"disable_bell",
		"default_author",
	}

	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Underline(true)

	cfgValues := []string{
		pathStyle.Render(cfg.ConfigDirectory),
		pathStyle.Render(cfg.WorkflowDirectory),
		pathStyle.Render(cfg.LogDirectory),
		pathStyle.Render(cfg.OutputDirectory),
		pathStyle.Render(cfg.InputDirectory),
		pathStyle.Render(cfg.PluginDirectory),
		fmt.Sprintf("%s\"%s\"%s", cli.COLOR_YELLOW, cfg.DefaultAuthor, cli.COLOR_RESET),
		fmt.Sprintf("%s%t%s", dlColor, cfg.DisableLogs, cli.COLOR_RESET),
		fmt.Sprintf("%s%t%s", acColor, cfg.AlwaysConfirm, cli.COLOR_RESET),
		fmt.Sprintf("%s%t%s", acColor, cfg.DisableBell, cli.COLOR_RESET),
	}

	strFields := ""
	for i, s := range cfgFields {
		if strFields == "" {
			strFields = lipgloss.JoinVertical(lipgloss.Left, s)
		} else {
			strFields = lipgloss.JoinVertical(lipgloss.Left, strFields, s)
		}

		if i < len(cfgFields)-1 {
			strFields = lipgloss.JoinVertical(lipgloss.Left, strFields, "")
		}
	}

	strValues := ""
	for i, s := range cfgValues {
		if strValues == "" {
			strValues = lipgloss.JoinVertical(lipgloss.Left, s)
		} else {
			strValues = lipgloss.JoinVertical(lipgloss.Left, strValues, s)
		}

		if i < len(cfgValues)-1 {
			strValues = lipgloss.JoinVertical(lipgloss.Left, strValues, "")
		}
	}

	str := lipgloss.JoinHorizontal(lipgloss.Top, strFields, "  ", strValues)

	fmt.Print(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).MarginLeft(2).Padding(1, 2).Render(str))
	fmt.Print("\n\n    ")

	cli.Question("Press any key to continue...", cli.QuestionOptions{})
}

const (
	logWarningColor    = "\u001b[38;5;209m"
	logNameColor       = "\u001b[38;5;240m"
	logDisplayLimit    = 20
	logAdditionalColor = "\u001b[38;5;248m"
)

func removeLogs(cfg *config.Config) {
	cli.Clear()

	fs, err := os.ReadDir(cfg.LogDirectory)
	if err != nil {
		panic(fmt.Sprintf("cannot read log directory! (%s)", err))
	}

	files := []string{}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".txt" {
			continue
		}
		files = append(files, f.Name())
	}

	if len(files) == 0 {
		fmt.Printf("\n%s%s!%s No log files to delete.\n\n", cli.COLOR_BOLD, logWarningColor, cli.COLOR_RESET)
		cli.Question("Press any key to continue...", cli.QuestionOptions{})
		return
	}

	displayFiles := files
	if len(files) > logDisplayLimit {
		displayFiles = files[:logDisplayLimit]
	}

	fmt.Printf("%sDeleting %s%s%s%sall%s%s log files:%s\n\n", logWarningColor, cli.COLOR_RESET, cli.COLOR_RED, cli.COLOR_BOLD, cli.COLOR_UNDERLINE, cli.COLOR_RESET, logWarningColor, cli.COLOR_RESET)
	fmt.Printf("%s%s%s%s\n", configPathColor, cli.COLOR_UNDERLINE, cfg.LogDirectory, cli.COLOR_RESET)

	for _, s := range displayFiles {
		fmt.Printf("  > %s%s%s\n", logNameColor, s, cli.COLOR_RESET)
	}

	if len(files) > logDisplayLimit {
		fmt.Printf("  > %s%s+%d%s more", logAdditionalColor, cli.COLOR_BOLD, len(files)-logDisplayLimit, cli.COLOR_RESET)
	}

	fmt.Printf("\n\n")

	answer, err := cli.Question(
		fmt.Sprintf("Are you sure? %sY%s/%s%s(N)%s", cli.COLOR_GREEN, cli.COLOR_RESET, cli.COLOR_BOLD, cli.COLOR_RED, cli.COLOR_RESET),
		cli.QuestionOptions{
			Normalize: true,
			Accepts:   []string{"y", "n"},
			Fallback:  "n",
		},
	)
	if err != nil {
		panic(fmt.Sprintf("error with result received from question! (%s)", err))
	}

	if answer == "y" {
		for _, f := range files {
			err := os.Remove(path.Join(cfg.LogDirectory, f))
			if err != nil {
				panic(fmt.Sprintf("failed to remove log file: %s! (%s)", f, err))
			}
		}

		fmt.Printf("\n%s%s!%s All log files deleted.\n\n", cli.COLOR_BOLD, logWarningColor, cli.COLOR_RESET)
		cli.Question("Press any key to continue...", cli.QuestionOptions{})
	}
}

type utilities_item struct {
	title, desc string
	index       int
}

func (i utilities_item) Title() string       { return i.title }
func (i utilities_item) Description() string { return i.desc }
func (i utilities_item) FilterValue() string { return i.title }

type utilities_model struct {
	list     list.Model
	selected *int
}

func (m utilities_model) Init() tea.Cmd {
	return nil
}

func (m utilities_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			result := m.list.SelectedItem()
			*m.selected = result.(utilities_item).index
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
		h, v := utilities_style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m utilities_model) View() string {
	return utilities_style.Render(m.list.View())
}
