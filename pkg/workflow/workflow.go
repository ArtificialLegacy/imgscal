package workflow

type WorkflowConfig struct {
	Name     string
	Version  string
	Requires []string
}

type Workflow struct {
	File    string
	Succeed bool
	Config  WorkflowConfig
}
