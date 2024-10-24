package lib

import (
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/anthonynsimon/bild/blend"
	golua "github.com/yuin/gopher-lua"
)

const LIB_BLEND = "blend"

/// @lib Blend
/// @import blend
/// @desc
/// Combines images using different blend modes.

func RegisterBlend(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BLEND, r, r.State, lg)

	/// @func add(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "add",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Add)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func color_burn(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "color_burn",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.ColorBurn)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func color_dodge(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "color_dodge",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.ColorDodge)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func darken(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "darken",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Darken)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func difference(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "difference",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Difference)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func divide(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "divide",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Divide)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func exclusion(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "exclusion",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Exclusion)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func lighten(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "lighten",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Lighten)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func linear_burn(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "linear_burn",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.LinearBurn)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func linear_light(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "linear_light",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.LinearLight)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func multiply(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "multiply",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Multiply)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func normal(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "normal",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Normal)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func overlay(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "overlay",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Overlay)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func screen(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "screen",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Screen)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func soft_light(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "soft_light",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.SoftLight)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func subtract(bg, fg, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "subtract",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := blendImages(r, lib, state, lg, args["bg"].(int), args["fg"].(int), args["name"].(string), d.Lib, d.Name, args["encoding"].(int), blend.Subtract)

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func opacity(bg, fg, percent, name, encoding) -> int<collection.IMAGE>
	/// @arg bg {int<collection.IMAGE>}
	/// @arg fg {int<collection.IMAGE>}
	/// @arg percent {float} - Between 0 and 1.
	/// @arg name {string} - The name of the new image.
	/// @arg encoding {int<image.ENCODING>}
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "opacity",
		[]lua.Arg{
			{Type: lua.INT, Name: "bg"},
			{Type: lua.INT, Name: "fg"},
			{Type: lua.FLOAT, Name: "percent"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			blendReady := make(chan struct{}, 2)

			var img image.Image
			var blended image.Image

			r.IC.SchedulePipe(state, args["bg"].(int), args["fg"].(int),
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						img = i.Self.Image
					},
				},
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						blended = blend.Opacity(img, i.Self.Image, args["percent"].(float64))
						blendReady <- struct{}{}
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						blendReady <- struct{}{}
					},
				})

			name := args["name"].(string)
			id := r.IC.ScheduleAdd(state, name, lg, d.Lib, d.Name, func(i *collection.Item[collection.ItemImage]) {
				<-blendReady
				i.Self = &collection.ItemImage{
					Image:    blended,
					Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
					Name:     name,
					Model:    imageutil.MODEL_RGBA,
				}
			})

			state.Push(golua.LNumber(id))
			return 1
		})
}

func blendImages(r *lua.Runner, lib *lua.Lib, state *golua.LState, lg *log.Logger, id1, id2 int, name, dl, dn string, encoding int, fn func(image.Image, image.Image) *image.RGBA) int {
	blendReady := make(chan struct{}, 2)

	var img image.Image
	var blended image.Image

	r.IC.SchedulePipe(state, id1, id2,
		&collection.Task[collection.ItemImage]{
			Lib:  dl,
			Name: dn,
			Fn: func(i *collection.Item[collection.ItemImage]) {
				img = i.Self.Image
			},
		},
		&collection.Task[collection.ItemImage]{
			Lib:  dl,
			Name: dn,
			Fn: func(i *collection.Item[collection.ItemImage]) {
				blended = fn(img, i.Self.Image)
				blendReady <- struct{}{}
			},
			Fail: func(i *collection.Item[collection.ItemImage]) {
				blendReady <- struct{}{}
			},
		})

	id := r.IC.ScheduleAdd(state, name, lg, dl, dn, func(i *collection.Item[collection.ItemImage]) {
		<-blendReady
		i.Self = &collection.ItemImage{
			Image:    blended,
			Encoding: lua.ParseEnum(encoding, imageutil.EncodingList, lib),
			Name:     name,
			Model:    imageutil.MODEL_RGBA,
		}
	})

	return id
}
