package workflow

type WorkflowConfig struct {
	name     string
	version  string
	requires []string
}

type Workflow struct {
	file    string
	succeed bool
	config  WorkflowConfig
}
