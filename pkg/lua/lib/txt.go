package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TXT = "txt"

/// @lib Text
/// @import txt
/// @desc
/// Library for reading and writing to '.txt' files.

func RegisterTXT(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TXT, r, r.State, lg)

	/// @func file_open(path, file, flag?) -> int<collection.FILE>
	/// @arg path {string} - The directory path to the file.
	/// @arg file {string} - The name of the file.
	/// @arg? flag {int<txt.FileFlags>} - Defaults to 'txt.CREATE'.
	/// @returns {int<collection.FILE>}
	/// @desc
	/// Will create the file if it does not exist,
	/// but will not create non-existant directories.
	/// Use 'bit.bitor' or 'bit.bitor_many' to combine flags.
	lib.CreateFunction(tab, "file_open",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "file"},
			{Type: lua.INT, Name: "flag", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fi, err := os.Stat(args["path"].(string))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("cannot find directory to txt file: %s", args["path"].(string)), log.LEVEL_ERROR)), 0)
			}

			if !fi.IsDir() {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("path provided is not a dir: %s", args["path"].(string)), log.LEVEL_ERROR)), 0)
			}

			chLog := log.NewLogger(fmt.Sprintf("file_%s", fi.Name()), lg)
			lg.Append(fmt.Sprintf("child log created: file_%s", fi.Name()), log.LEVEL_INFO)

			id := r.FC.AddItem(&chLog)

			r.FC.Schedule(id, &collection.Task[collection.ItemFile]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemFile]) {
					flag := args["flag"].(int)
					if flag == 0 {
						flag = os.O_CREATE
					}
					f, err := os.OpenFile(path.Join(args["path"].(string), args["file"].(string)), flag, 0o666)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("failed to open txt file: %s", args["file"].(string)), log.LEVEL_ERROR)), 0)
					}

					i.Self = &collection.ItemFile{
						Name: args["path"].(string),
						File: f,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func write(id, txt)
	/// @arg id {int<collection.FILE>}
	/// @arg txt {string}
	lib.CreateFunction(tab, "write",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "txt"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.FC.Schedule(args["id"].(int), &collection.Task[collection.ItemFile]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemFile]) {
					_, err := i.Self.File.WriteString(args["txt"].(string))
					if err != nil {
						state.Error(golua.LString(i.Lg.Append(fmt.Sprintf("failed to write to txt file: %d", args["id"]), log.LEVEL_ERROR)), 0)
					}
				},
			})
			return 0
		})

	/// @constants File Flags
	/// @const CREATE
	/// @const TRUNC
	/// @const EXCL
	/// @const APPEND
	/// @const RDWR
	/// @const RDONLY
	/// @const WRONLY
	tab.RawSetString("CREATE", golua.LNumber(os.O_CREATE))
	tab.RawSetString("TRUNC", golua.LNumber(os.O_TRUNC))
	tab.RawSetString("EXCL", golua.LNumber(os.O_EXCL))
	tab.RawSetString("APPEND", golua.LNumber(os.O_APPEND))
	tab.RawSetString("RDWR", golua.LNumber(os.O_RDWR))
	tab.RawSetString("RDONLY", golua.LNumber(os.O_RDONLY))
	tab.RawSetString("WRONLY", golua.LNumber(os.O_WRONLY))
}
