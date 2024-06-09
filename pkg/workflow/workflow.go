package workflow

import "github.com/ArtificialLegacy/imgscal/pkg/image"

type Workflow struct {
	Name     string
	Version  string
	Author   string
	Desc     string
	Requires []string
}

func NewWorkflow() Workflow {
	return Workflow{
		Name:     "",
		Version:  "",
		Author:   "",
		Desc:     "",
		Requires: []string{},
	}
}

type WorkflowData struct {
	IC image.ImageCollection
}
