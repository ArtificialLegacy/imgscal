package lua

import (
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/image"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/Shopify/go-lua"
)

type Runner struct {
	State *lua.State
	IC    *image.ImageCollection
	lg    *log.Logger
}

func NewRunner(state *lua.State, lg *log.Logger) Runner {
	return Runner{
		State: state,
		IC:    image.NewImageCollection(lg),
		lg:    lg,
	}
}

func (r *Runner) Run(file string) error {
	pwd, _ := os.Getwd()

	err := lua.DoFile(r.State, path.Join(pwd, file))
	if err != nil {
		return err
	}

	return nil
}
