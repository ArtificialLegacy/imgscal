package cli

import "fmt"

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func ClearLine() {
	fmt.Print("\033[1000D\033[K")
}
