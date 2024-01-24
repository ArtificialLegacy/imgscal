package imgscal

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/utility/luautility"
	"github.com/Shopify/go-lua"
)

func Job(state *lua.State, file string, job string) (string, error) {
	filename := file

	switch job {
	case "rename":
		var options map[string]interface{}
		if state.IsTable(-1) {
			options = make(map[string]interface{})

			options["name"] = luautility.ParseOption(state, "name")
			options["prefix"] = luautility.ParseOption(state, "prefix")
			options["suffix"] = luautility.ParseOption(state, "suffix")

			name, err := rename(filename, options)
			if err != nil {
				return filename, err
			}
			filename = name

		} else {
			return filename, fmt.Errorf("options must be provided to run rename")
		}
	case "output":
		err := output(filename)
		if err != nil {
			return filename, err
		}
	case "copy":
		var options map[string]interface{}
		if state.IsTable(-1) {
			options = make(map[string]interface{})

			options["name"] = luautility.ParseOption(state, "name")
			options["prefix"] = luautility.ParseOption(state, "prefix")
			options["suffix"] = luautility.ParseOption(state, "suffix")

			name, err := copy(filename, options)
			if err != nil {
				return filename, err
			}
			filename = name
		} else {
			return filename, fmt.Errorf("options must be provided to run copy")
		}
	default:
		return filename, fmt.Errorf("unknown imgscal job: %s", job)
	}

	return filename, nil
}
