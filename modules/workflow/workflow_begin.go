package workflow

import (
	"errors"
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/modules/cli"
)

// returns the path to the image(s) the user wants to upscale and an error if the path does not exist
func WorkflowBegin() (string, error) {
	answer, _ := cli.Question("Enter the path to the image(s) you want to upscale: ", cli.QuestionOptions{
		Normalize: false,
		Accepts:   nil,
		Fallback:  "",
	})

	_, err := os.Stat(answer)

	if err != nil || os.IsNotExist(err) {
		println(fmt.Sprintf("%s! The path you entered does not exist. Please try again.%s", cli.RED, cli.RESET))
		return "", errors.New("Path does not exist.")
	}

	return answer, nil
}
