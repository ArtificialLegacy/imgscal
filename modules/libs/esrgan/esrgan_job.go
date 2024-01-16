package esrgan

import (
	"github.com/ArtificialLegacy/imgscal/modules/utility/luautility"
	"github.com/Shopify/go-lua"
)

func Job(state *lua.State, file string, job string) (string, error) {
	filename := file

	options := make(map[string]interface{})

	if state.IsTable(-1) {
		options["scale"] = luautility.ParseOption(state, "scale")
	} else {
		options["scale"] = "4"
	}

	switch job {
	case "x4":
		err := X4(filename, options)
		if err != nil {
			return filename, err
		}
	case "animex4":
		err := AnimeX4(filename, options)
		if err != nil {
			return filename, err
		}
	}

	return filename, nil
}
