package lib

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_PIPE = "pipe"

/// @lib Pipe
/// @import pipe
/// @desc
/// Library for piping data in and out of a cli workflow.

func RegisterPipe(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_PIPE, r, r.State, lg)

	/// @func in_string() -> string
	/// @returns {string}
	lib.CreateFunction(tab, "in_string",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to read from stdin: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(b))
			return 1
		})

	/// @func in_bytes() -> []int
	/// @returns {[]int}
	lib.CreateFunction(tab, "in_bytes",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to read from stdin: %s", log.LEVEL_ERROR, err))
			}

			t := state.NewTable()
			for i, v := range b {
				t.RawSetInt(i+1, golua.LNumber(v))
			}

			state.Push(t)
			return 1
		})

	/// @func in_image(name, encoding, model?) -> int<collection.IMAGE>
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "in_image",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to read from stdin: %s", log.LEVEL_ERROR, err))
			}

			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			id := r.IC.ScheduleAdd(state, name, lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemImage]) {
				img, err := imageutil.Decode(strings.NewReader(string(b)), encoding)
				if err != nil {
					lua.Error(state, i.Lg.Appendf("piped data is an invalid image: %s", log.LEVEL_ERROR, err))
				}

				img, model = imageutil.Limit(img, model)

				i.Self = &collection.ItemImage{
					Name:     name,
					Image:    img,
					Encoding: encoding,
					Model:    model,
				}
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func in_image_cached(encoding, model?) -> int<collection.CRATE_CACHEDIMAGE>
	/// @arg encoding {int<image.Encoding>}
	/// @arg? model {int<image.ColorModel>} - Used only to specify default when there is an unsupported color model.
	/// @returns {int<collection.CRATE_CACHEDIMAGE>}
	lib.CreateFunction(tab, "in_image_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to read from stdin: %s", log.LEVEL_ERROR, err))
			}

			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

			img, err := imageutil.Decode(strings.NewReader(string(b)), encoding)
			if err != nil {
				lua.Error(state, lg.Appendf("piped data is an invalid image: %s", log.LEVEL_ERROR, err))
			}

			img, model = imageutil.Limit(img, model)

			id := r.CR_CIM.Add(&collection.CachedImageItem{
				Image: img,
				Model: model,
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func out_string(str)
	/// @arg str {string}
	lib.CreateFunction(tab, "out_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			str := args["str"].(string)

			_, err := fmt.Fprint(os.Stdout, str)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to write to stdout: %s", log.LEVEL_ERROR, err))
			}

			return 0
		})

	/// @func out_bytes(b)
	/// @arg b {[]int}
	lib.CreateFunction(tab, "out_bytes",
		[]lua.Arg{
			lua.ArgArray("b", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			bl := args["b"].([]any)

			for _, v := range bl {
				_, err := fmt.Fprint(os.Stdout, byte(v.(int)))
				if err != nil {
					lua.Error(state, lg.Appendf("failed to write to stdout: %s", log.LEVEL_ERROR, err))
				}
			}

			return 0
		})

	/// @func out_image(id)
	/// @arg id {int<collection.IMAGE>}
	/// @blocking
	lib.CreateFunction(tab, "out_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			id := args["id"].(int)

			<-r.IC.Schedule(state, id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					err := imageutil.Encode(os.Stdout, i.Self.Image, i.Self.Encoding)
					if err != nil {
						lua.Error(state, i.Lg.Appendf("failed to encode image: %s", log.LEVEL_ERROR, err))
					}
				},
			})

			return 0
		})

	/// @func out_image_cached(id, encoding)
	/// @arg id {int<collection.CRATE_CACHEDIMAGE>}
	/// @blocking
	lib.CreateFunction(tab, "out_image_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				lua.Error(state, lg.Append("can only use the pipe library in cli mode", log.LEVEL_ERROR))
			}

			id := args["id"].(int)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

			item, err := r.CR_CIM.Item(id)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to retrieve image: %s", log.LEVEL_ERROR, err))
			}

			err = imageutil.Encode(os.Stdout, item.Image, encoding)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to encode image: %s", log.LEVEL_ERROR, err))
			}

			return 0
		})
}
