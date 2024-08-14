package states

import (
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/config"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
)

const (
	UTILITIES_OPTION_CONFIG int = iota
	UTILITIES_OPTION_RMLogs
	UTILITIES_OPTION_BACK
)

var utilitiesOptions = []string{
	UTILITIES_OPTION_CONFIG: "View Config",
	UTILITIES_OPTION_RMLogs: "Delete Logs",
	UTILITIES_OPTION_BACK:   fmt.Sprintf("%sReturn%s", cli.COLOR_RED, cli.COLOR_RESET),
}

func Utilities(sm *statemachine.StateMachine) error {
	cli.Clear()

	result, err := cli.SelectMenu("Utilities", utilitiesOptions)
	if err != nil {
		return err
	}

	switch result {
	case UTILITIES_OPTION_CONFIG:
		viewConfig(sm.Config)

	case UTILITIES_OPTION_RMLogs:
		removeLogs(sm.Config)

	case UTILITIES_OPTION_BACK:
		sm.SetState(STATE_MAIN)

	default:
		panic(fmt.Sprintf("UTILITIES_OPTION %d is not handled.", result))
	}

	return nil
}

const configPathColor = "\u001b[38;5;240m\u001b[4m"
const configTrueColor = "\u001b[38;5;195m"
const configFalseColor = "\u001b[38;5;209m"

func viewConfig(cfg *config.Config) {
	cli.Clear()

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")

	var dlColor = configFalseColor
	if cfg.DisableLogs {
		dlColor = configTrueColor
	}
	var acColor = configFalseColor
	if cfg.AlwaysConfirm {
		acColor = configTrueColor
	}

	fmt.Printf("%sImgScal Config%s v%s\n", cli.COLOR_BOLD, cli.COLOR_RESET, cfg.ConfigVersion)
	fmt.Printf("%s%s%s\n\n", configPathColor, cfgPath, cli.COLOR_RESET)

	fmt.Printf("workflow_directory: %s%s%s\n", configPathColor, cfg.WorkflowDirectory, cli.COLOR_RESET)
	fmt.Printf("log_directory:      %s%s%s\n", configPathColor, cfg.LogDirectory, cli.COLOR_RESET)
	fmt.Printf("output_directory:   %s%s%s\n", configPathColor, cfg.OutputDirectory, cli.COLOR_RESET)
	fmt.Printf("disable_logs:       %s%t%s\n", dlColor, cfg.DisableLogs, cli.COLOR_RESET)
	fmt.Printf("always_confirm:     %s%t%s\n\n", acColor, cfg.AlwaysConfirm, cli.COLOR_RESET)

	cli.Question("Press any key to continue...", cli.QuestionOptions{})
}

const logWarningColor = "\u001b[38;5;209m"
const logNameColor = "\u001b[38;5;240m"
const logDisplayLimit = 20
const logAdditionalColor = "\u001b[38;5;248m"

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
