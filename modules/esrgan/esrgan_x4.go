package esrgan

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cmd"
)

func X4(infile string, index int, total int) {
	filename := workflowInit(infile, "RealESRGAN-x4plus", index, total)
	pwd, _ := os.Getwd()

	cmd.CommandRun(fmt.Sprintf("%s\\esrgan-tool\\realesrgan-ncnn-vulkan.exe -i %s -o %s\\outputs\\up_%s -n realesrgan-x4plus-anime", pwd, infile, pwd, filename))
}
