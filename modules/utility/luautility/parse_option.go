package luautility

import (
	"github.com/Shopify/go-lua"
)

func ParseOption(state *lua.State, option string) string {
	value := ""

	state.RawGetValue(-1, option)

	if state.IsString(-1) {
		value, _ = state.ToString(-1)
	}

	state.Remove(-1)

	return value
}
