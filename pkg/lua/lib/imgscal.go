package lib

import "github.com/Shopify/go-lua"

const LIB_IMGSCAL = "imgscal"

func RegisterImgscal(state *lua.State) {
	state.NewTable()
	state.SetGlobal("imgscal")
	state.Global("imgscal")

	state.PushGoFunction(func(state *lua.State) int {
		return 0
	})
	state.Field(-1, "name")

	state.PushGoFunction(func(state *lua.State) int {
		return 0
	})
	state.Field(-1, "prompt_file")

	state.PushGoFunction(func(state *lua.State) int {
		return 0
	})
	state.Field(-1, "out")
}
