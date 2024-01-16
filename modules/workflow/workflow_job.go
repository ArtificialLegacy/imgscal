package workflow

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ArtificialLegacy/imgscal/modules/libs/esrgan"
	"github.com/ArtificialLegacy/imgscal/modules/libs/imgscal"
	"github.com/Shopify/go-lua"
)

func workflowJob(state *lua.State, wf *Workflow) error {
	var job string
	var file string

	if state.IsString(1) {
		job, _ = state.ToString(1)
	} else {
		return fmt.Errorf("job must be a string")
	}

	if state.IsString(2) {
		file, _ = state.ToString(2)
	} else {
		return fmt.Errorf("file must be a string")
	}

	jobSplit := strings.Split(job, ".")
	lib := jobSplit[0]
	action := jobSplit[1]

	if !slices.Contains[[]string](wf.Config.Requires, lib) {
		return fmt.Errorf("job lib %s not required", lib)
	}

	returnValue := ""

	switch lib {
	case "imgscal":
		returnValue, _ = imgscal.Job(state, file, action)
	case "esrgan":
		returnValue, _ = esrgan.Job(state, file, action)
	default:
		return fmt.Errorf("job lib %s not found", lib)
	}

	state.PushString(returnValue)

	return nil
}
