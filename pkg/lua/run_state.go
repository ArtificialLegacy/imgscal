package lua

import (
	"os"
	"path"

	"github.com/Shopify/go-lua"
)

type Runner struct {
	state *lua.State
}

func NewRunner(state *lua.State) Runner {
	return Runner{
		state: state,
	}
}

func (r *Runner) Run(file string) error {
	pwd, _ := os.Getwd()

	err := lua.DoFile(r.state, path.Join(pwd, file))
	if err != nil {
		return err
	}

	return nil
}
