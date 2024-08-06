package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_BIT = "bit"

func RegisterBit(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BIT, r, r.State, lg)

	/// @func bitor()
	/// @arg a
	/// @arg b
	/// @returns a | b
	lib.CreateFunction(tab, "bitor",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a | b))
			return 1
		})

	/// @func bitor_many()
	/// @arg list - []int
	/// @returns list[1] | list[2] | list[3]...
	lib.CreateFunction(tab, "bitor_many",
		[]lua.Arg{
			lua.ArgArray("list", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			list := args["list"].(map[string]any)
			acc := 0

			for i := range len(list) {
				v := list[string(i+1)].(int)
				acc |= v
			}

			state.Push(golua.LNumber(acc))
			return 1
		})

	/// @func bitand()
	/// @arg a
	/// @arg b
	/// @returns a & b
	lib.CreateFunction(tab, "bitand",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a & b))
			return 1
		})

	/// @func bitxor()
	/// @arg a
	/// @arg b
	/// @returns a ^ b
	lib.CreateFunction(tab, "bitxor",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a ^ b))
			return 1
		})

	/// @func bitclear()
	/// @arg a
	/// @arg b
	/// @returns a &^ b
	/// @desc
	/// equivalent to a & (~b)
	/// keeps the bits in a where the bit in b is 0, otherwise the bit is 0.
	lib.CreateFunction(tab, "bitclear",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a &^ b))
			return 1
		})

	/// @func bitnot()
	/// @arg a
	/// @returns ^a
	lib.CreateFunction(tab, "bitnot",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)

			state.Push(golua.LNumber(^a))
			return 1
		})

	/// @func bit_rshift()
	/// @arg a
	/// @arg b
	/// @returns a >> b
	lib.CreateFunction(tab, "bit_rshift",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a >> b))
			return 1
		})

	/// @func bit_lshift()
	/// @arg a
	/// @arg b
	/// @returns a << b
	lib.CreateFunction(tab, "bit_lshift",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)
			b := args["b"].(int)

			state.Push(golua.LNumber(a << b))
			return 1
		})
}
