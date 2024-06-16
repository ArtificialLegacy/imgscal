package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
)

const LIB_TXT = "txt"

func RegisterTXT(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_TXT, r.State, lg)

	/// @func file_open()
	/// @arg path - the directory path to the file
	/// @arg file - the name of the file
	/// @arg flag - int, defaults to O_CREATE.
	/// @returns id - id of the opened file
	/// @desc
	/// Will create the file if it does not exist,
	/// but will not create non-existant directories.
	/// Use bitwise OR to combine flags.
	lib.CreateFunction("file_open",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.STRING, Name: "file"},
			{Type: lua.INT, Name: "flag", Optional: true},
		},
		func(state *golua.State, args map[string]any) int {
			fi, err := os.Stat(args["path"].(string))
			if err != nil {
				state.PushString(lg.Append(fmt.Sprintf("cannot find directory to txt file: %s", args["path"].(string)), log.LEVEL_ERROR))
				state.Error()
			}

			if !fi.IsDir() {
				state.PushString(lg.Append(fmt.Sprintf("path provided is not a dir: %s", args["path"].(string)), log.LEVEL_ERROR))
				state.Error()
			}

			id := r.FC.AddItem(args["path"].(string))

			r.FC.Schedule(id, &collection.Task[os.File]{
				Lib:  LIB_TXT,
				Name: "file_open",
				Fn: func(i *collection.Item[os.File]) {
					flag := args["flag"].(int)
					if flag == 0 {
						flag = os.O_CREATE
					}
					f, err := os.OpenFile(path.Join(args["path"].(string), args["file"].(string)), flag, 0o666)
					if err != nil {
						state.PushString(lg.Append(fmt.Sprintf("failed to open txt file: %s", args["file"].(string)), log.LEVEL_ERROR))
						state.Error()
					}

					i.Self = f
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func write()
	/// @arg id - id of the file to write to
	/// @arg txt - string of text to write
	lib.CreateFunction("write",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "txt"},
		},
		func(state *golua.State, args map[string]any) int {
			r.FC.Schedule(args["id"].(int), &collection.Task[os.File]{
				Lib:  LIB_TXT,
				Name: "write",
				Fn: func(i *collection.Item[os.File]) {
					_, err := i.Self.WriteString(args["txt"].(string))
					if err != nil {
						state.PushString(lg.Append(fmt.Sprintf("failed to write to txt file: %d", args["id"]), log.LEVEL_ERROR))
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
	r.State.PushInteger(os.O_CREATE)
	r.State.SetField(-1, "O_CREATE")
	r.State.PushInteger(os.O_TRUNC)
	r.State.SetField(-1, "O_TRUNC")
	r.State.PushInteger(os.O_EXCL)
	r.State.SetField(-1, "O_EXCL")
	r.State.PushInteger(os.O_APPEND)
	r.State.SetField(-1, "O_APPEND")
	r.State.PushInteger(os.O_RDWR)
	r.State.SetField(-1, "O_RDWR")
	r.State.PushInteger(os.O_RDONLY)
	r.State.SetField(-1, "O_RDONLY")
	r.State.PushInteger(os.O_WRONLY)
	r.State.SetField(-1, "O_WRONLY")
}
