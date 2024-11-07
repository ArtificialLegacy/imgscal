package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/config"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/charmbracelet/huh"
)

const (
	workspace_file = ".luarc.json"
	workspace_data = "{}"
)

func main() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}
	_, err = os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user home directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")

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
		err := os.MkdirAll(cfg.WorkflowDirectory, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to make workflow directory! (%s)", err))
		}
	}

	workflows, errs, err := workflow.WorkflowList(cfg.WorkflowDirectory)
	if err != nil {
		panic(fmt.Sprintf("failed to scan for pre-existing workflows: (%s)", err))
	}

	if errs != nil && len(*errs) > 0 {
		errList := ""
		for _, e := range *errs {
			errList += e.Error() + " | "
		}

		panic(fmt.Sprintf("failed to scan for pre-existing workflows: (%s)", errList))
	}

	cli.Clear()

	name := ""
	author := cfg.DefaultAuthor
	version := "1.0.0"
	desc := ""

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Description("Only allows the characters: A-Z, a-z (case sensitive), 0-9 and _.").Value(&name).Validate(validateName(workflows)),
			huh.NewInput().Title("Author").Value(&author).Validate(validateAuthor),
			huh.NewInput().Title("Version").Value(&version).Placeholder("1.0.0"),
			huh.NewText().Title("Description").Value(&desc),
		),
	).Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			fmt.Print("Tool exitted early, no workflow created.\n")
			os.Exit(1)
		}
		panic(fmt.Sprintf("failed to run form: %s", err))
	}

	wfPath := path.Join(cfg.WorkflowDirectory, name)

	f, err := os.Stat(wfPath)
	if err == nil || f != nil {
		panic(fmt.Sprintf("path in workflow directory already exists: %s", name))
	}

	descLong := []string{}

	descList := strings.Split(strings.TrimSpace(desc), "\n")
	if len(descList) == 0 {
		desc = ""
	} else {
		desc = descList[0]
		if len(descList) > 1 {
			descLong = descList[1:]
		}
	}

	wf := workflow.WorkflowJSON{
		Name:       name,
		Author:     author,
		Version:    version,
		APIVersion: workflow.API_VERSION,
		Desc:       desc,
		DescLong:   descLong,
	}

	err = os.Mkdir(wfPath, 0o777)
	if err != nil {
		panic(fmt.Sprintf("failed to make workflow directory: %s with error (%s)", wfPath, err))
	}

	wfb, err := json.MarshalIndent(wf, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal workflow json: %s", err))
	}

	err = os.WriteFile(path.Join(wfPath, "workflow.json"), wfb, 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to write workflow.json: %s", err))
	}

	err = os.WriteFile(path.Join(wfPath, workspace_file), []byte(workspace_data), 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to write %s: %s", workspace_file, err))
	}

	fmt.Printf("%sCreated Workflow: %s.%s\n", cli.COLOR_GREEN, name, cli.COLOR_RESET)
	fmt.Printf("Full Path: %s`%s`%s.\n", "\u001b[38;5;240m\u001b[4m", wfPath, cli.COLOR_RESET)
}

var nameHasValidChars = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func validateName(workflows *[]*workflow.Workflow) func(string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("workflow name is required, cannot be empty")
		}

		if !nameHasValidChars.MatchString(value) {
			return fmt.Errorf("workflow name provided contains invalid characters")
		}

		for _, wf := range *workflows {
			if value == wf.Name {
				return fmt.Errorf("workflow name already exists")
			}
		}

		return nil
	}
}

func validateAuthor(value string) error {
	if value == "" {
		return fmt.Errorf("workflow author is required, cannot be empty")
	}

	return nil
}
