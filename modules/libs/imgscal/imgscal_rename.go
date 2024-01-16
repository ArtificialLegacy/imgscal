package imgscal

import (
	"fmt"
	"os"
	"strings"
)

func rename(file string, options map[string]interface{}) (string, error) {
	fileSplit := strings.Split(file, ".")
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

	err := os.Rename(fmt.Sprintf("%s\\temp\\%s", pwd, file), fmt.Sprintf("%s\\temp\\%s", pwd, filename))

	return filename, err
}
