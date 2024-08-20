package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
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

	err = copyExamples(fs, wd, "", cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Installed example workflows to: %s.\n", path.Join(cfg.WorkflowDirectory, "examples"))
}

func copyExamples(fs []fs.DirEntry, wd, prefix string, cfg *config.Config) error {
	for _, fs := range fs {
		if fs.IsDir() {
			ifs, err := os.ReadDir(path.Join(wd, "examples", prefix, fs.Name()))
			if err != nil {
				return err
			}
			err = copyExamples(ifs, wd, prefix+fs.Name(), cfg)
			if err != nil {
				return err
			}

			continue
		}

		pth := path.Join("examples", prefix, fs.Name())
		b, err := os.ReadFile(path.Join(wd, pth))
		if err != nil {
			panic(fmt.Sprintf("failed to read example workflow: %s! (%s)", pth, err))
		}

		if prefix != "" {
			err := os.MkdirAll(path.Join(cfg.WorkflowDirectory, "examples", prefix), 0o777)
			if err != nil {
				panic(fmt.Sprintf("failed to make directories for example workflow: %s! (%s)", pth, err))
			}
		}

		err = os.WriteFile(path.Join(cfg.WorkflowDirectory, pth), b, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to write example workflow: %s! (%s)", pth, err))
		}
	}

	return nil
}
