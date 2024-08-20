package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

func WorkflowList(workflowDir string) (*[]*Workflow, *[]error, error) {
	workflows := &[]*Workflow{}
	errorList := &[]error{}

	err := workflowScan(workflowDir, workflowDir, workflows, errorList)
	if err != nil {
		return nil, errorList, err
	}

	return workflows, errorList, nil
}

func workflowScan(dir, workflowDir string, workflows *[]*Workflow, errorList *[]error) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			pth := path.Join(dir, file.Name())
			err = workflowScan(pth, workflowDir, workflows, errorList)
			if err != nil {
				return err
			}
		}
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == "workflow.json" {
			pth := path.Join(dir, file.Name())
			w, err := WorkflowParse(pth, strings.TrimPrefix(pth, workflowDir))
			if err != nil {
				*errorList = append(*errorList, fmt.Errorf("failed to parse workflow: %s with error: %s", pth, err))
			}
			*workflows = append(*workflows, w)
		}
	}

	return nil
}

func WorkflowParse(name, base string) (*Workflow, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	w := &WorkflowJSON{}
	err = json.Unmarshal(b, w)
	if err != nil {
		return nil, err
	}

	workflow := NewWorkflow(name, base, w)
	return workflow, nil
}
