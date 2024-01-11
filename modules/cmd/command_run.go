package cmd

import (
	"os"
	"os/exec"
)

// CommandRun runs a command and prints the output to the console.
func CommandRun(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	println(cmd.String())
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
