package main

import (
	"fmt"
	"os"
	"path"
)

const (
	workspace_file = ".luarc.json"
	workspace_data = "{}"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get working directory: %s", err))
	}

	pth := path.Join(wd, workspace_file)

	err = os.WriteFile(pth, []byte(workspace_data), 0o666)
	if err != nil {
		panic(fmt.Sprintf("failed to write %s: %s", workspace_file, err))
	}

	fmt.Printf("Created lua workspace file in the current directory: %s\n", pth)
}
