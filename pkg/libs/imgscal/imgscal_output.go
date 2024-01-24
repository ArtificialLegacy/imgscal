package imgscal

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/utility/file"
)

func output(filename string) error {
	pwd, _ := os.Getwd()

	file1 := fmt.Sprintf("%s\\temp\\%s", pwd, filename)
	file2 := fmt.Sprintf("%s\\outputs\\%s", pwd, filename)

	if _, err := file.Copy(file1, file2); err != nil {
		println(err.Error())
		return err
	}

	return nil
}
