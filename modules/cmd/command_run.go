package cmd

import (
	"os"
	"os/exec"
	"strings"
)

// CommandRun runs a command and prints the output to the console.
func CommandRun(command string) error {
	commandSplit := strings.Split(command, " ")

	cmd := exec.Command(commandSplit[0], commandSplit[1:]...)
	println(cmd.String())
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
