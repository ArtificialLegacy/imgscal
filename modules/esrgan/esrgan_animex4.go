package esrgan

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cmd"
)

func AnimeX4(infile string, index int, total int) {
	filename := workloadInit(infile, "RealESRGAN-x4plus Anime", index, total)
	pwd, _ := os.Getwd()

	cmd.CommandRun(fmt.Sprintf("%s\\esrgan-tool\\realesrgan-ncnn-vulkan.exe", pwd), "-i", infile, "-o", fmt.Sprintf("%s\\outputs\\up_%s", pwd, filename), "-n", "realesrgan-x4plus-anime")
}
