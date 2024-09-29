package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_STRINGS = "strings"

/// @lib Strings
/// @import strings

func RegisterStrings(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_STRINGS, r, r.State, lg)

	/// @func to_rune(str) -> int
	/// @arg str {string}
	/// @returns {int}
	/// @desc
	/// Only returns the first character as a rune, or 0.
	lib.CreateFunction(tab, "to_rune",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := args["str"].(string)
			runes := []rune(str)
			v := rune(0)
			if len(runes) > 0 {
				v = runes[0]
			}

			state.Push(golua.LNumber(v))
			return 1
		})

	/// @func to_runes(str) -> []int
	/// @arg str {string}
	/// @returns {[]int}
	lib.CreateFunction(tab, "to_runes",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := args["str"].(string)
			runes := []rune(str)

			runeList := state.NewTable()
			for i, v := range runes {
				runeList.RawSetInt(i, golua.LNumber(v))
			}

			state.Push(runeList)
			return 1
		})

	/// @func from_rune(rune) -> string
	/// @arg rune {int}
	/// @returns {string}
	lib.CreateFunction(tab, "from_rune",
		[]lua.Arg{
			{Type: lua.INT, Name: "rune"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			rn := args["rune"].(int)
			str := string([]rune{rune(rn)})

			state.Push(golua.LString(str))
			return 1
		})

	/// @func from_runes(runes) -> string
	/// @arg runes {[]int}
	/// @returns {string}
	lib.CreateFunction(tab, "from_runes",
		[]lua.Arg{
			lua.ArgArray("runes", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			runes := args["runes"].([]any)
			runeList := make([]rune, len(runes))

			for i, v := range runes {
				runeList[i] = rune(v.(int))
			}

			state.Push(golua.LString(runeList))
			return 1
		})
}
