package states

import (
	"fmt"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/statemachine"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
)

func WorkflowCMDEnter(sm *statemachine.StateMachine, name string) {
	sm.SetState(STATE_WORKFLOW_CMD)
	sm.Data = name
}

func WorkflowCMD(sm *statemachine.StateMachine) error {
	name := sm.Data.(string)

	wf, errlist, err := workflow.WorkflowList(sm.Config.WorkflowDirectory)
	if err != nil {
		fmt.Printf("failed to scan for workflows: %s\n", err)
		sm.SetState(STATE_EXIT)
		return err
	}
	if len(*errlist) > 0 {
		fmt.Printf("failed to scan for workflows: %+v\n", *errlist)
		sm.SetState(STATE_EXIT)
		return err
	}

	foundPath := ""
	var foundWf *workflow.Workflow
	for _, w := range *wf {
		if w.Name == name {
			found, ok := w.CliWorkflows["*"]
			if !ok {
				fmt.Printf("cannot use workflow base name when there is no star workflow: %s\n", err)
				sm.SetState(STATE_EXIT)
				return err
			}

			foundWf = w
			foundPath = path.Join(path.Dir(w.Base), found)
			break
		}

		if !strings.HasPrefix(name, path.Base(path.Dir(w.Base))) {
			continue
		}

		found, ok := w.CliWorkflows[strings.TrimPrefix(name, path.Base(path.Dir(w.Base))+"/")]
		if !ok {
			continue
		}

		foundWf = w
		foundPath = path.Join(path.Dir(w.Base), found)
		break
	}

	if foundWf.APIVersion > workflow.API_VERSION {
		fmt.Printf("Workflow %s has a newer API version than supported: %d. Current API: %d\n", foundWf.Name, foundWf.APIVersion, workflow.API_VERSION)
		sm.SetState(STATE_EXIT)
		return err
	}

	WorkflowRunEnter(sm, WorkflowRunData{Script: foundPath, Name: name})
	return nil
}
