package lib

import (
	"strings"

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

	/// @func compare(a, b) -> string
	/// @arg a {string}
	/// @arg b {string}
	/// @returns {int} - 0 if a == b, -1 if a < b, and 1 if a > b.
	lib.CreateFunction(tab, "compare",
		[]lua.Arg{
			{Type: lua.STRING, Name: "a"},
			{Type: lua.STRING, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			a := args["a"].(string)
			b := args["b"].(string)

			value := strings.Compare(a, b)

			state.Push(golua.LNumber(value))
			return 1
		})

	/// @func has_prefix(s, prefix) -> bool
	/// @arg s {string}
	/// @arg prefix {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "has_prefix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "prefix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			prefix := args["prefix"].(string)

			found := strings.HasPrefix(s, prefix)

			state.Push(golua.LBool(found))
			return 1
		})

	/// @func trim_prefix(s, prefix) -> string
	/// @arg s {string}
	/// @arg prefix {string}
	/// @returns {string}
	lib.CreateFunction(tab, "trim_prefix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "prefix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			prefix := args["prefix"].(string)

			after := strings.TrimPrefix(s, prefix)

			state.Push(golua.LString(after))
			return 1
		})

	/// @func cut_prefix(s, prefix) -> string, bool
	/// @arg s {string}
	/// @arg prefix {string}
	/// @returns {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "cut_prefix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "prefix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			prefix := args["prefix"].(string)

			after, found := strings.CutPrefix(s, prefix)

			state.Push(golua.LString(after))
			state.Push(golua.LBool(found))
			return 2
		})

	/// @func has_suffix(s, suffix) -> bool
	/// @arg s {string}
	/// @arg suffix {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "has_suffix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "suffix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			suffix := args["suffix"].(string)

			found := strings.HasSuffix(s, suffix)

			state.Push(golua.LBool(found))
			return 1
		})

	/// @func trim_suffix(s, suffix) -> string
	/// @arg s {string}
	/// @arg suffix {string}
	/// @returns {string}
	lib.CreateFunction(tab, "trim_suffix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "suffix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			suffix := args["suffix"].(string)

			after := strings.TrimSuffix(s, suffix)

			state.Push(golua.LString(after))
			return 1
		})

	/// @func cut_suffix(s, suffix) -> string, bool
	/// @arg s {string}
	/// @arg suffix {string}
	/// @returns {string}
	/// @returns {bool}
	lib.CreateFunction(tab, "cut_suffix",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
			{Type: lua.STRING, Name: "suffix"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := args["s"].(string)
			suffix := args["suffix"].(string)

			after, found := strings.CutSuffix(s, suffix)

			state.Push(golua.LString(after))
			state.Push(golua.LBool(found))
			return 2
		})
}
