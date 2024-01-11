package esrgan

import (
	"fmt"
	"os"
)

func Remove() error {
	pwd, _ := os.Getwd()
	err := os.RemoveAll(fmt.Sprintf("%s\\esrgan-tool\\", pwd))

	return err
}
