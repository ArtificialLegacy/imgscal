//go:generate goversioninfo -icon=assets/favicon.ico -manifest=imgscal.exe.manifest

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/config"
	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/states"
)

func main() {
	sm := statemachine.NewStateMachine(states.STATE_COUNT)

	sm.AddState(states.STATE_MAIN, states.MainMenu)
	sm.AddState(states.STATE_EXIT, states.Exit)
	sm.AddState(states.STATE_UTILITIES, states.Utilities)
	sm.AddState(states.STATE_WORKFLOW_LIST, states.WorkflowList)
	sm.AddState(states.STATE_WORKFLOW_CONFIRM, states.WorkflowConfirm)
	sm.AddState(states.STATE_WORKFLOW_RUN, states.WorkflowRun)
	sm.AddState(states.STATE_WORKFLOW_FAIL, states.WorkflowFail)
	sm.AddState(states.STATE_WORKFLOW_FINISH, states.WorkflowFinish)
	sm.AddState(states.STATE_WORKFLOW_HELP, states.WorkflowHelp)
	sm.AddState(states.STATE_WORKFLOW_CMD, states.WorkflowCMD)
	sm.AddState(states.STATE_WORKFLOW_CMDLIST, states.WorkflowCMDList)

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user home directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")
	if envCfg, exists := os.LookupEnv("IMGSCAL_CONFIG"); exists {
		cfgPath = envCfg
	}

	_, err = os.Stat(cfgPath)
	if err != nil {
		err = os.MkdirAll(path.Join(cfgDir, "imgscal"), 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make config directory! (%s)", err))
		}

		dfltConfig := config.NewConfigWithDefaults(homeDir)
		b, err := json.MarshalIndent(dfltConfig, "", "    ")
		if err != nil {
			panic(fmt.Sprintf("failed to marshal default config to json! (%s)", err))
		}

		err = os.WriteFile(cfgPath, b, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to write default config to file! (%s)", err))
		}

		sm.Config = dfltConfig
	} else {
		b, err := os.ReadFile(cfgPath)
		if err != nil {
			panic(fmt.Sprintf("failed to read config file! (%s)", err))
		}

		cfg := config.NewConfig()
		err = json.Unmarshal(b, cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal json config file! (%s)", err))
		}

		sm.Config = cfg
	}

	_, err = os.Stat(sm.Config.WorkflowDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.WorkflowDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make workflow directory! (%s)", err))
		}
	}

	_, err = os.Stat(sm.Config.ConfigDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.ConfigDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make config directory! (%s)", err))
		}
		f, err := os.OpenFile(path.Join(sm.Config.ConfigDirectory, ".gitignore"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to make config .gitignore! (%s)", err))
		}

		_, err = f.WriteString("**/*.secrets.json")
		if err != nil {
			panic(fmt.Sprintf("failed to write to config .gitignore! (%s)", err))
		}
	}

	_, err = os.Stat(sm.Config.OutputDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.OutputDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make output directory! (%s)", err))
		}
	}

	_, err = os.Stat(sm.Config.InputDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.InputDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make input directory! (%s)", err))
		}
	}

	_, err = os.Stat(sm.Config.PluginDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.PluginDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make plugin directory! (%s)", err))
		}
	}

	if len(os.Args) > 1 {
		pth := os.Args[1]
		sm.CliMode = true

		if pth == "list" {
			states.WorkflowCMDList(sm)
		} else if pth == "help" {
			if len(os.Args) > 2 {
				states.WorkflowHelpEnter(sm, os.Args[2])
			} else {
				fmt.Printf("%s'help' requires a workflow name to be specified!%s\n\n", cli.COLOR_RED, cli.COLOR_RESET)
				os.Exit(1)
			}
		} else {
			states.WorkflowCMDEnter(sm, os.Args[1])
		}
	}

	for {
		sm.Step()
	}
}
