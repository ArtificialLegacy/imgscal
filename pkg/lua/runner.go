package lua

import (
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/Shopify/go-lua"
)

type Runner[T any] struct {
	state *lua.State
	Data  *T
}

func NewRunner[T any](state *lua.State, data *T) Runner[T] {
	return Runner[T]{
		state: state,
		Data:  data,
	}
}

func (r *Runner[T]) Register(fn func(state *lua.State, data *T, lg *log.Logger), lg *log.Logger) {
	fn(r.state, r.Data, lg)
}

func (r *Runner[T]) Run(file string) error {
	pwd, _ := os.Getwd()

	err := lua.DoFile(r.state, path.Join(pwd, file))
	if err != nil {
		return err
	}

	return nil
}
