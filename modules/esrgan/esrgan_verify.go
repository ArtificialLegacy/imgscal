package esrgan

import (
	"fmt"
	"os"
)

func Verify() bool {
	pwd, _ := os.Getwd()
	exists, err := os.Stat(fmt.Sprintf("%s\\esrgan-tool\\", pwd))

	if err != nil {
		return false
	}

	return exists != nil
}
