package cli

import "github.com/manifoldco/promptui"

func SelectMenu(question string, options []string) (int, error) {
	prompt := promptui.Select{
		Label: question,
		Items: options,
	}

	result, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return result, nil
}
