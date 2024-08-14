package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/config"
)

func main() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")

	_, err = os.Stat(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("no config file, cannot access log directory setting! (%s)", err))
	}

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file! (%s)", err))
	}

	cfg := config.NewConfig()
	err = json.Unmarshal(b, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal json config file! (%s)", err))
	}

	_, err = os.Stat(cfg.WorkflowDirectory)
	if err != nil {
		panic(fmt.Sprintf("failed to find workflow directory! (%s)", err))
	}

	err = os.MkdirAll(path.Join(cfg.WorkflowDirectory, "examples"), 0o777)
	if err != nil {
		panic(fmt.Sprintf("failed to create examples directory! (%s)", err))
	}

	wd, _ := os.Getwd()
	fs, err := os.ReadDir(path.Join(wd, "examples"))
	if err != nil {
		panic(fmt.Sprintf("failed to read examples directory! (%s)", err))
	}

	for _, fs := range fs {
		b, err := os.ReadFile(path.Join(wd, "examples", fs.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed to read example workflow: %s! (%s)", fs.Name(), err))
		}

		err = os.WriteFile(path.Join(cfg.WorkflowDirectory, "examples", fs.Name()), b, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to write example workflow: %s! (%s)", fs.Name(), err))
		}
	}

	fmt.Printf("Installed example workflows to: %s.\n", path.Join(cfg.WorkflowDirectory, "examples"))
}
