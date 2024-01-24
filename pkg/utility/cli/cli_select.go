package cli

import "github.com/manifoldco/promptui"

// Menu displays a menu and returns the selected option.
func Menu(question string, options []string) (int8, error) {
	prompt := promptui.Select{
		Label: question,
		Items: options,
	}

	result, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return int8(result), nil
}
