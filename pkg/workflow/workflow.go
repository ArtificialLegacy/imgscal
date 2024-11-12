package workflow

import (
	"path"
	"strings"
)

const (
	API_VERSION = 1
	SCHEMA      = "https://gist.githubusercontent.com/ArtificialLegacy/9711f20511e76b519aedb729a6762b9f/raw/de77e999654060a38d7a4e7eea8aeb4f5ee1273e/imgscal_workflow.json"
)

type Workflow struct {
	Name         string
	Filepath     string
	Base         string
	Location     string
	Author       string
	Version      string
	APIVersion   int
	Desc         string
	Workflows    map[string]string
	CliWorkflows map[string]string
}

type WorkflowJSON struct {
	Schema       string            `json:"$schema"`
	Name         string            `json:"name"`
	Author       string            `json:"author"`
	Version      string            `json:"version"`
	APIVersion   int               `json:"api_version"`
	Desc         string            `json:"desc"`
	DescLong     []string          `json:"desc_long,omitempty"`
	Workflows    map[string]string `json:"workflows,omitempty"`
	CliWorkflows map[string]string `json:"cli_workflows,omitempty"`
}

func NewWorkflow(filepath, base string, input *WorkflowJSON) *Workflow {
	return &Workflow{
		Name:       input.Name,
		Filepath:   filepath,
		Base:       base,
		Location:   path.Dir(filepath),
		Author:     input.Author,
		Version:    input.Version,
		APIVersion: input.APIVersion,

		Desc: input.Desc + " " + strings.Join(input.DescLong, " "),

		Workflows:    input.Workflows,
		CliWorkflows: input.CliWorkflows,
	}
}
