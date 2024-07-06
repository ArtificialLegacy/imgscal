package lib

import (
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/qeesung/image2ascii/convert"

	golua "github.com/yuin/gopher-lua"
)

const LIB_ASCII = "ascii"

func RegisterASCII(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_ASCII, r, r.State, lg)

	/// @func to_file()
	/// @arg image_id
	/// @arg filepath - directories to file must exist.
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	lib.CreateFunction(tab, "to_file",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					converter := convert.NewImageConverter()
					str := converter.Image2ASCIIString(i.Self.Image, &convert.Options{
						Colored:  args["color"].(bool),
						Reversed: args["reverse"].(bool),
					})

					f, err := os.OpenFile(args["path"].(string), os.O_CREATE|os.O_TRUNC, 0o666)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR)), 0)
					}
					defer f.Close()

					f.WriteString(str)
				},
			})

			return 0
		})

	/// @func to_file_size()
	/// @arg image_id
	/// @arg filepath - directories to file must exist.
	/// @arg width
	/// @arg height
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	lib.CreateFunction(tab, "to_file_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					converter := convert.NewImageConverter()
					str := converter.Image2ASCIIString(i.Self.Image, &convert.Options{
						FixedWidth:  args["width"].(int),
						FixedHeight: args["height"].(int),
						Colored:     args["color"].(bool),
						Reversed:    args["reverse"].(bool),
					})

					f, err := os.OpenFile(args["path"].(string), os.O_CREATE|os.O_TRUNC, 0o666)
					if err != nil {
						state.Error(golua.LString(i.Lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR)), 0)
					}
					defer f.Close()

					f.WriteString(str)
				},
			})

			return 0
		},
	)

	/// @func to_string()
	/// @arg image_id
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	/// @returns the ascii art as a string
	/// @blocking
	lib.CreateFunction(tab, "to_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := ""

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					converter := convert.NewImageConverter()
					str = converter.Image2ASCIIString(i.Self.Image, &convert.Options{
						Colored:  args["color"].(bool),
						Reversed: args["reverse"].(bool),
					})
				},
			})

			state.Push(golua.LString(str))
			return 1
		},
	)

	/// @func to_string_size()
	/// @arg image_id
	/// @arg width
	/// @arg height
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	/// @returns the ascii art as a string
	/// @blocking
	lib.CreateFunction(tab, "to_string_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := ""

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					converter := convert.NewImageConverter()
					str = converter.Image2ASCIIString(i.Self.Image, &convert.Options{
						FixedWidth:  args["width"].(int),
						FixedHeight: args["height"].(int),
						Colored:     args["color"].(bool),
						Reversed:    args["reverse"].(bool),
					})
				},
			})

			state.Push(golua.LString(str))
			return 1
		},
	)
}
