package imgscal

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/utility/file"
)

func copy(filestart string, options map[string]interface{}) (string, error) {
	fileSplit := strings.Split(filestart, ".")
	filename := strings.Join(fileSplit[:len(fileSplit)-1], ".")

	if options["name"] != "" {
		filename = options["name"].(string)
	}
	if options["prefix"] != "" {
		filename = options["prefix"].(string) + filename
	}
	if options["suffix"] != "" {
		filename += options["suffix"].(string)
	}

	filename = filename + "." + fileSplit[len(fileSplit)-1]

	pwd, _ := os.Getwd()

	_, err := file.Copy(fmt.Sprintf("%s\\temp\\%s", pwd, filestart), fmt.Sprintf("%s\\temp\\%s", pwd, filename))

	return filename, err
}
