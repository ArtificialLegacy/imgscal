//go:generate goversioninfo

package main

import (
	"os"

	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
	"github.com/ArtificialLegacy/imgscal/modules/states"
)

func main() {
	os.Setenv("ESRGAN_DOWNLOAD_URL", "https://github.com/xinntao/Real-ESRGAN/releases/download/v0.2.5.0/")
	os.Setenv("ESRGAN_FOLDER_NAME", "realesrgan-ncnn-vulkan-20220424-windows")

	stateMachine := statemachine.NewStateMachine()

	stateMachine.AddState(states.ESRGANVerify)
	stateMachine.AddState(states.ESRGANDownload)
	stateMachine.AddState(states.ESRGANFail)
	stateMachine.AddState(states.LandingMenu)
	stateMachine.AddState(states.ESRGANManage)
	stateMachine.AddState(states.WorkloadMenu)
	stateMachine.AddState(states.ESRGANX4)
	stateMachine.AddState(states.ESRGANAnimeX4)
	stateMachine.AddState(states.WorkloadFinish)

	stateMachine.Transition(statemachine.ESRGAN_VERIFY)
}
