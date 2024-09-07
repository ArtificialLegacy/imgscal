package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_BIT = "bit"

/// @lib Bit
/// @import bit
/// @desc
/// Utility library for performing bitwise operations.

func RegisterBit(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BIT, r, r.State, lg)

	/// @func bitor(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a | b).
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

	/// @func bitor_many(operands...) -> int
	/// @arg operands {int...}
	/// @returns {int} - The result of all operands on (0 | operand[1] | operand[2]...).
	lib.CreateFunction(tab, "bitor_many",
		[]lua.Arg{
			lua.ArgVariadic("operands", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			list := args["operands"].([]any)
			acc := 0

			for i := range len(list) {
				v := list[i].(int)
				acc |= v
			}

			state.Push(golua.LNumber(acc))
			return 1
		})

	/// @func bitand(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a & b).
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

	/// @func bitxor(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a ^ b).
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

	/// @func bitclear(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a &^ b).
	/// @desc
	/// Equivalent to (a & (~b)).
	/// Keeps the bits in 'a' where the bit in 'b' is 0, otherwise the bit is 0.
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

	/// @func bitnot(a) -> int
	/// @arg a {int}
	/// @returns {int} - The result of (^a).
	lib.CreateFunction(tab, "bitnot",
		[]lua.Arg{
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(int)

			state.Push(golua.LNumber(^a))
			return 1
		})

	/// @func bit_rshift(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a >> b).
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

	/// @func bit_lshift(a, b) -> int
	/// @arg a {int}
	/// @arg b {int}
	/// @returns {int} - The result of (a << b).
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

	/// @func byte_string(nums...) -> string
	/// @arg nums {int...} - A list of numbers to convert to bytes. These will be treated as uint8.
	/// @returns {string} - The byte string representation of the numbers.
	lib.CreateFunction(tab, "byte_string",
		[]lua.Arg{
			lua.ArgVariadic("nums", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			list := args["nums"].([]any)
			bytes := make([]byte, len(list))

			for i, v := range list {
				bytes[i] = byte(v.(int))
			}

			state.Push(golua.LString(bytes))
			return 1
		})
}
