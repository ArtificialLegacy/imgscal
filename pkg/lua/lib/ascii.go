package lib

import (
	"image"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
	"github.com/qeesung/image2ascii/convert"
)

const LIB_ASCII = "ascii"

func RegisterASCII(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_ASCII, r.State, lg)

	/// @func to_file()
	/// @arg image_id
	/// @arg filepath - directories to file must exist.
	/// @arg color - boolean, for terminal dislay
	/// @arg reverse - boolean
	lib.CreateFunction("to_file",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.State, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  LIB_ASCII,
				Name: "to_file",
				Fn: func(i *collection.Item[image.Image]) {
					converter := convert.NewImageConverter()
					str := converter.Image2ASCIIString(*i.Self, &convert.Options{
						Colored:  args["color"].(bool),
						Reversed: args["reverse"].(bool),
					})

					f, err := os.OpenFile(args["path"].(string), os.O_CREATE|os.O_TRUNC, 0o666)
					if err != nil {
						r.State.PushString(lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR))
						r.State.Error()
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
	lib.CreateFunction("to_file_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "path"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.State, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  LIB_ASCII,
				Name: "to_file_size",
				Fn: func(i *collection.Item[image.Image]) {
					converter := convert.NewImageConverter()
					str := converter.Image2ASCIIString(*i.Self, &convert.Options{
						FixedWidth:  args["width"].(int),
						FixedHeight: args["height"].(int),
						Colored:     args["color"].(bool),
						Reversed:    args["reverse"].(bool),
					})

					f, err := os.OpenFile(args["path"].(string), os.O_CREATE|os.O_TRUNC, 0o666)
					if err != nil {
						r.State.PushString(lg.Append("failed to open file for saving ascii string", log.LEVEL_ERROR))
						r.State.Error()
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
	lib.CreateFunction("to_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.State, args map[string]any) int {
			str := ""

			<-r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  LIB_ASCII,
				Name: "to_string",
				Fn: func(i *collection.Item[image.Image]) {
					converter := convert.NewImageConverter()
					str = converter.Image2ASCIIString(*i.Self, &convert.Options{
						Colored:  args["color"].(bool),
						Reversed: args["reverse"].(bool),
					})
				},
			})

			r.State.PushString(str)
			return 0
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
	lib.CreateFunction("to_string_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.BOOL, Name: "color"},
			{Type: lua.BOOL, Name: "reverse"},
		},
		func(state *golua.State, args map[string]any) int {
			str := ""

			<-r.IC.Schedule(args["id"].(int), &collection.Task[image.Image]{
				Lib:  LIB_ASCII,
				Name: "to_string_size",
				Fn: func(i *collection.Item[image.Image]) {
					converter := convert.NewImageConverter()
					str = converter.Image2ASCIIString(*i.Self, &convert.Options{
						FixedWidth:  args["width"].(int),
						FixedHeight: args["height"].(int),
						Colored:     args["color"].(bool),
						Reversed:    args["reverse"].(bool),
					})
				},
			})

			r.State.PushString(str)
			return 0
		},
	)

	r.State.SetGlobal(LIB_ASCII)
}
