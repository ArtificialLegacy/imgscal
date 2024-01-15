//go:generate goversioninfo -icon=assets/favicon.ico -manifest=imgscal.exe.manifest

package main

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/statemachine"
	"github.com/ArtificialLegacy/imgscal/modules/states"
	"github.com/ArtificialLegacy/imgscal/modules/workflow"
)

func main() {
	os.Setenv("ESRGAN_DOWNLOAD_URL", "https://github.com/xinntao/Real-ESRGAN/releases/download/v0.2.5.0/")
	os.Setenv("ESRGAN_FOLDER_NAME", "realesrgan-ncnn-vulkan-20220424-windows")

	workflows := workflow.WorkflowsLoad()

	stateMachine := statemachine.NewStateMachine()
	stateMachine.SetWorkflowState(workflows)

	stateMachine.AddState(states.ESRGANVerify)
	stateMachine.AddState(states.ESRGANDownload)
	stateMachine.AddState(states.ESRGANFail)
	stateMachine.AddState(states.LandingMenu)
	stateMachine.AddState(states.ESRGANManage)
	stateMachine.AddState(states.WorkflowMenu)
	stateMachine.AddState(states.ESRGANX4)
	stateMachine.AddState(states.ESRGANAnimeX4)
	stateMachine.AddState(states.WorkflowFinish)

	stateMachine.Transition(statemachine.ESRGAN_VERIFY)
}
