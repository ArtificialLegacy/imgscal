package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TXT = "txt"

/// @lib Text
/// @import txt
/// @desc
/// Library for reading and writing to text files.

func RegisterTXT(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TXT, r, r.State, lg)

	/// @func write(path, text)
	/// @arg path {string}
	/// @arg text {string}
	lib.CreateFunction(tab, "write",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			text := args["text"].(string)

			err := os.WriteFile(path, []byte(text), 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to write to file: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func write_lines(path, lines)
	/// @arg path {string}
	/// @arg lines {[]string}
	lib.CreateFunction(tab, "write_lines",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			lua.ArgArray("lines", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			lines := args["lines"].([]any)

			text := ""
			for _, line := range lines {
				text += line.(string) + "\n"
			}

			err := os.WriteFile(path, []byte(text), 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to write to file: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func truncate(path)
	/// @arg path {string}
	lib.CreateFunction(tab, "truncate",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)

			err := os.WriteFile(path, []byte{}, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to truncate file: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func append(path, text)
	/// @arg path {string}
	/// @arg text {string}
	lib.CreateFunction(tab, "append",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			text := args["text"].(string)

			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to append to file: %s", err), log.LEVEL_ERROR)), 0)
			}
			defer file.Close()

			_, err = file.WriteString(text)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to append to file: %s", err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func append_lines(path, lines)
	/// @arg path {string}
	/// @arg lines {[]string}
	lib.CreateFunction(tab, "append_lines",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			lua.ArgArray("lines", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			lines := args["lines"].([]any)

			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to append to file: %s", err), log.LEVEL_ERROR)), 0)
			}
			defer file.Close()

			for _, line := range lines {
				_, err = file.WriteString(line.(string) + "\n")
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to append to file: %s", err), log.LEVEL_ERROR)), 0)
				}
			}

			return 0
		})

	/// @func read(path) -> string
	/// @arg path {string}
	/// @returns {string}
	lib.CreateFunction(tab, "read",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(string(text)))
			return 1
		})

	/// @func read_lines(path) -> []string
	/// @arg path {string}
	/// @returns {[]string}
	lib.CreateFunction(tab, "read_lines",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			lines := strings.Split(string(text), "\n")
			arr := state.NewTable()
			for i, line := range lines {
				arr.RawSetInt(i+1, golua.LString(line))
			}

			state.Push(arr)
			return 1
		})

	/// @func line_count(path) -> int
	/// @arg path {string}
	/// @returns {int}
	lib.CreateFunction(tab, "line_count",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			b := strings.Split(string(text), "\n")
			state.Push(golua.LNumber(len(b)))

			return 1
		})

	/// @func line(path, index) -> string
	/// @arg path {string}
	/// @arg index {int}
	/// @returns {string} - Returns an empty string if the index is out of bounds.
	lib.CreateFunction(tab, "line",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "index"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			index := args["index"].(int)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			lines := strings.Split(string(text), "\n")
			if index < 0 || index >= len(lines) {
				state.Push(golua.LString(""))
			} else {
				state.Push(golua.LString(lines[index]))
			}

			return 1
		})

	/// @func iter(path, func)
	/// @arg path {string}
	/// @arg func {function(line string, index int)}
	lib.CreateFunction(tab, "iter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			fn := args["func"].(*golua.LFunction)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			lines := strings.Split(string(text), "\n")

			for i, line := range lines {
				state.Push(fn)
				state.Push(golua.LString(line))
				state.Push(golua.LNumber(i + 1))
				state.Call(2, 0)
			}

			return 0
		})

	/// @func map(path, func) -> string
	/// @arg path {string}
	/// @arg func {function(line string, index int) -> string}
	/// @returns {string}
	lib.CreateFunction(tab, "map",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			fn := args["func"].(*golua.LFunction)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			lines := strings.Split(string(text), "\n")
			var out string

			for i, line := range lines {
				state.Push(fn)
				state.Push(golua.LString(line))
				state.Push(golua.LNumber(i + 1))
				state.Call(2, 1)
				out += state.ToString(-1)
				state.Pop(1)
			}

			state.Push(golua.LString(out))
			return 1
		})

	/// @func map_lines(path, func) -> []string
	/// @arg path {string}
	/// @arg func {function(line string, index int) -> string}
	/// @returns {[]string}
	lib.CreateFunction(tab, "map_lines",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.FUNC, Name: "func"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			path := args["path"].(string)
			fn := args["func"].(*golua.LFunction)

			text, err := os.ReadFile(path)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to read file: %s", err), log.LEVEL_ERROR)), 0)
			}

			lines := strings.Split(string(text), "\n")
			arr := state.NewTable()

			for i, line := range lines {
				state.Push(fn)
				state.Push(golua.LString(line))
				state.Push(golua.LNumber(i + 1))
				state.Call(2, 1)
				arr.RawSetInt(i+1, golua.LString(state.ToString(-1)))
				state.Pop(1)
			}

			state.Push(arr)
			return 1
		})
}
