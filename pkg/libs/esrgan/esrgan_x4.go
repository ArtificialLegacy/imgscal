package esrgan

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/utility/cmd"
)

func X4(infile string, options map[string]interface{}) error {
	pwd, _ := os.Getwd()
	err := cmd.CommandRun(
		fmt.Sprintf(
			"%s\\esrgan-tool\\realesrgan-ncnn-vulkan.exe -i %s\\temp\\%s -o %s\\temp\\%s -n realesrgan-x4plus -s %s",
			pwd, pwd, infile, pwd, infile, options["scale"],
		),
	)

	return err
}
