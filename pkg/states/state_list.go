package states

const (
	STATE_MAIN int = iota
	STATE_EXIT
	STATE_WORKFLOW_LIST
	STATE_WORKFLOW_CONFIRM
	STATE_WORKFLOW_FAIL_LOAD
	STATE_WORKFLOW_RUN
	STATE_WORKFLOW_FINISH
	STATE_WORKFLOW_FAIL_RUN

	STATE_COUNT
)
