package esrgan

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cmd"
)

func X4(infile string, index int, total int) {
	pwd, _ := os.Getwd()
	cmd.CommandRun(
		fmt.Sprintf("%s\\esrgan-tool\\realesrgan-ncnn-vulkan.exe %s\\temp\\%s -o %s\\temp\\%s -n realesrgan-x4plus", pwd, pwd, infile, pwd, infile),
	)
}
