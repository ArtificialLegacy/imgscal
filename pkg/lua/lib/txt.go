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

func RegisterTXT(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TXT, r, r.State, lg)

	/// @func file_open()
	/// @arg path - the directory path to the file
	/// @arg file - the name of the file
	/// @arg? flag - int, defaults to O_CREATE.
	/// @returns id - id of the opened file
	/// @desc
	/// Will create the file if it does not exist,
	/// but will not create non-existant directories.
	/// Use bitwise OR to combine flags.
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

			chLog := log.NewLogger(fmt.Sprintf("file_%s", fi.Name()))
			chLog.Parent = lg
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

	/// @func write()
	/// @arg id - id of the file to write to
	/// @arg txt - string of text to write
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

	/// @constants File open flags
	/// @const O_CREATE
	/// @const O_TRUNC
	/// @const O_EXCL
	/// @const O_APPEND
	/// @const O_RDWR
	/// @const O_RDONLY
	/// @const O_WRONLY
	r.State.SetField(tab, "O_CREATE", golua.LNumber(os.O_CREATE))
	r.State.SetField(tab, "O_TRUNC", golua.LNumber(os.O_TRUNC))
	r.State.SetField(tab, "O_EXCL", golua.LNumber(os.O_EXCL))
	r.State.SetField(tab, "O_APPEND", golua.LNumber(os.O_APPEND))
	r.State.SetField(tab, "O_RDWR", golua.LNumber(os.O_RDWR))
	r.State.SetField(tab, "O_RDONLY", golua.LNumber(os.O_RDONLY))
	r.State.SetField(tab, "O_WRONLY", golua.LNumber(os.O_WRONLY))
}
