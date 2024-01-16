//go:generate goversioninfo -icon=assets/favicon.ico -manifest=imgscal.exe.manifest

package main

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/state/statemachine"
	"github.com/ArtificialLegacy/imgscal/modules/state/states"
	"github.com/ArtificialLegacy/imgscal/modules/workflow"
)

func main() {
	os.Setenv("ESRGAN_DOWNLOAD_URL", "https://github.com/xinntao/Real-ESRGAN/releases/download/v0.2.5.0/")
	os.Setenv("ESRGAN_FOLDER_NAME", "realesrgan-ncnn-vulkan-20220424-windows")

	pwd, _ := os.Getwd()
	if _, err := os.Stat(fmt.Sprintf("%s\\outputs", pwd)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s\\outputs", pwd), 0777)
	}

	wfs := workflow.WorkflowsLoad()

	sm := statemachine.NewStateMachine()
	sm.SetWorkflowsState(wfs)

	sm.AddStates([]statemachine.State{
		states.ESRGANVerify,
		states.ESRGANDownload,
		states.ESRGANFail,
		states.LandingMenu,
		states.ESRGANManage,
		states.WorkflowMenu,
		states.WorkflowRun,
		states.WorkflowFinish,
	})

	sm.Transition(statemachine.ESRGAN_VERIFY)

	for true {
		sm.Step()
	}
}
