package workflow

import (
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
)

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
	IC collection.Collection[image.Image]
}
