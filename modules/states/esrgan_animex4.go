package states

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
	"github.com/ArtificialLegacy/imgscal/modules/esrgan"
	statemachine "github.com/ArtificialLegacy/imgscal/modules/state_machine"
)

var esrganAnimeX4Enter statemachine.StateEnterFunction = func(from statemachine.CliState, transition func(to statemachine.CliState) error) {
	cli.Clear()

	answer, err := esrgan.WorkloadBegin()
	if err != nil {
		transition(statemachine.WORKLOAD_FINISH)
		return
	}

	file, _ := os.Stat(answer)

	if file.IsDir() {
		files, err := os.ReadDir(answer)
		if err != nil {
			transition(statemachine.WORKLOAD_FINISH)
			return
		}
		for index, file := range files {
			esrgan.AnimeX4(answer+"\\"+file.Name(), index+1, len(files))
		}
	} else {
		esrgan.AnimeX4(answer, 1, 1)
	}

	transition(statemachine.WORKLOAD_FINISH)
}

var ESRGANAnimeX4 = statemachine.NewState(statemachine.ESRGAN_ANIMEX4, esrganAnimeX4Enter, nil, []statemachine.CliState{statemachine.WORKLOAD_FINISH})
