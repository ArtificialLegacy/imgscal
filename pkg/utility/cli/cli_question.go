package cli

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

type QuestionOptions struct {
	Normalize bool
	Accepts   []string
	Fallback  string
}

func Question(question string, options QuestionOptions) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if options.Normalize {
		result = strings.ToLower(result)
	}

	if (options.Accepts != nil) && (len(options.Accepts) > 0) {
		var found bool
		for _, accept := range options.Accepts {
			if result == accept {
				found = true
				break
			}
		}

		if !found {
			if options.Fallback != "" {
				return options.Fallback, nil
			} else {
				return "", errors.New("No fallback provided and input does not match any of the accepted values.")
			}
		}
	}

	return result, nil
}
