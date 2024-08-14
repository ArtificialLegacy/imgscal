//go:generate goversioninfo -icon=assets/favicon.ico -manifest=imgscal.exe.manifest

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

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
	sm.AddState(states.STATE_WORKFLOW_FAIL_LOAD, states.WorkflowFailLoad)
	sm.AddState(states.STATE_WORKFLOW_RUN, states.WorkflowRun)
	sm.AddState(states.STATE_WORKFLOW_FAIL_RUN, states.WorkflowFailRun)
	sm.AddState(states.STATE_WORKFLOW_FINISH, states.WorkflowFinish)

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user home directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")

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

	_, err = os.Stat(sm.Config.OutputDirectory)
	if err != nil {
		err := os.MkdirAll(sm.Config.OutputDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make output directory! (%s)", err))
		}
	}

	if len(os.Args) > 1 {
		pth := os.Args[1]
		if path.Ext(pth) != ".lua" {
			pth += ".lua"
		}

		sm.CliMode = true
		sm.PushString(path.Join("workflows/", pth))
		sm.SetState(states.STATE_WORKFLOW_CONFIRM)
	}

	for {
		sm.Step()
	}
}
