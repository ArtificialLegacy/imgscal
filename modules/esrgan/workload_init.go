package esrgan

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
)

func workloadInit(infile string, workload string, index int, total int) string {
	fileSplit := strings.Split(infile, "\\")
	filename := fileSplit[len(fileSplit)-1]

	println(fmt.Sprintf("%s!%s Running %s on %s (Image %d of %d)", cli.CYAN, cli.RESET, workload, filename, index, total))

	pwd, _ := os.Getwd()

	if _, err := os.Stat(fmt.Sprintf("%s\\outputs", pwd)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s\\outputs", pwd), 0777)
	}

	return filename
}
