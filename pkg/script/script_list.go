package script

import (
	"os"
	"path"
	"path/filepath"
)

type Script struct {
	Filepath string
	Name     string
}

func WorkflowList(workflowDir string) ([]Script, error) {
	scripts := []Script{}

	err := scriptScan(workflowDir, "", &scripts)
	if err != nil {
		return nil, err
	}

	return scripts, nil
}

func scriptScan(dir, prefix string, scripts *[]Script) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			pth := path.Join(dir, file.Name())
			err := scriptScan(pth, path.Join(prefix, file.Name()), scripts)
			if err != nil {
				return err
			}
		}
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".lua" {
			pth := path.Join(dir, file.Name())
			script := Script{Filepath: pth, Name: path.Join(prefix, file.Name())}
			*scripts = append(*scripts, script)
		}
	}

	return nil
}
