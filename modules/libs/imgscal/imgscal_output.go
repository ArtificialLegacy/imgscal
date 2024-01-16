package imgscal

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/utility/file"
)

func output(filename string) error {
	pwd, _ := os.Getwd()
	_, err := file.Copy(fmt.Sprintf("%s\\temp\\%s", pwd, filename), fmt.Sprintf("%s\\outputs\\%s", pwd, filename))
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}
