package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/akamensky/argparse"
)

const helpstr = `---@param info imgscal_WorkflowInfo
function help(info)
    return [[
Usage:
    > 
    ]]
end
`

const initstr = `---@param workflow imgscal_WorkflowInit
function init(workflow)
    workflow.import({

    })
end
`

const mainstr = `function main()

end
`

func main() {
	parser := argparse.NewParser("imgscal-entrypoint", "Creates entrypoints in a workflow.")

	name := parser.StringPositional(&argparse.Options{Required: true})
	pth := parser.StringPositional(&argparse.Options{Required: true})

	iscli := parser.Flag("c", "cli", &argparse.Options{})

	err := parser.Parse(os.Args)
	if err != nil {
		panic(fmt.Sprintf("failed to parse cmd line args: %s", err))
	}

	if *name == "" {
		fmt.Print("! Entry point name is required.\n")
		os.Exit(1)
	}
	if !validateEntrypointName(*name) {
		fmt.Print("! Entry point name has invalid characters.\n")
		os.Exit(1)
	}

	if *pth == "" {
		fmt.Print("! Entry point path is required.\n")
		os.Exit(1)
	}
	if !strings.HasSuffix(*pth, ".lua") {
		fmt.Print("! Entry point path must have the .lua extension.\n")
		os.Exit(1)
	}
	if !validateEntrypointPath(*pth) {
		fmt.Print("! Entry point path has invalid characters.\n")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get working directory: %s", err))
	}

	jsonpth := path.Join(wd, "workflow.json")
	b, err := os.ReadFile(jsonpth)
	if err != nil {
		fmt.Print("Failed to find workflow.json in the current directory.\n")
		os.Exit(1)
	}

	wf := &workflow.WorkflowJSON{}
	err = json.Unmarshal(b, wf)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal workflow.json: %s", err))
	}

	wfmap := &wf.Workflows
	if *iscli {
		wfmap = &wf.CliWorkflows
	}
	if *wfmap == nil {
		*wfmap = map[string]string{}
	}

	if _, ok := (*wfmap)[*name]; ok {
		fmt.Printf("Workflow with the name %s already exists.\n", *name)
		os.Exit(1)
	}

	entrypath := path.Join(wd, *pth)
	dir := path.Dir(entrypath)
	if dir != "." {
		err = os.MkdirAll(dir, 0o777)
		if err != nil {
			panic(fmt.Sprintf("failed to create path to directories: %s", err))
		}
	}

	fe, err := os.Stat(entrypath)
	if err != nil && fe != nil {
		fmt.Printf("Lua file at path %s already exists.\n", *pth)
		os.Exit(1)
	}

	fs, err := os.OpenFile(entrypath, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to open lua file for writing: %s", err))
	}
	defer fs.Close()

	if *iscli {
		fs.WriteString(helpstr)
		fs.WriteString("\n")
	}

	fs.WriteString(initstr)
	fs.WriteString("\n")
	fs.WriteString(mainstr)

	(*wfmap)[*name] = *pth

	b, err = json.MarshalIndent(wf, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json: %s", err))
	}

	err = os.WriteFile(jsonpth, b, 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to write updated workflow.json: %s", err))
	}

	fmt.Printf("Created new entry point named %s%s%s at path %s%s%s.\n", cli.COLOR_BLUE, *name, cli.COLOR_RESET, "\u001b[38;5;240m\u001b[4m", *pth, cli.COLOR_RESET)
}

var nameHasValidChars = regexp.MustCompile(`^[a-zA-Z0-9_*\-]+$`)

func validateEntrypointName(value string) bool {
	return nameHasValidChars.MatchString(value)
}

var pathHasValidChars = regexp.MustCompile(`^[a-zA-Z0-9_\-/]+.lua$`)

func validateEntrypointPath(value string) bool {
	return pathHasValidChars.MatchString(value)
}
