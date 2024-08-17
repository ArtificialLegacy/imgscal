package workflow

type Workflow struct {
	Name         string
	Version      string
	Author       string
	Desc         string
	Requires     []string
	CliExclusive bool
	Verbose      bool
}

func NewWorkflow() Workflow {
	return Workflow{
		Name:         "",
		Version:      "",
		Author:       "",
		Desc:         "",
		Requires:     []string{},
		CliExclusive: false,
		Verbose:      false,
	}
}
