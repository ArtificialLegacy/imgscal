package lib

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/fogleman/gg"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CONTEXT = "context"

/// @lib Context
/// @import context
/// @desc
/// Library for creating and drawing to canvases.

type Point map[string]float64

func RegisterContext(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CONTEXT, r, r.State, lg)

	/// @func degrees(radians) -> float
	/// @arg radians {float}
	/// @returns {float} - Degrees.
	lib.CreateFunction(tab, "degrees",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "rad"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			deg := gg.Degrees(args["rad"].(float64))
			state.Push(golua.LNumber(deg))
			return 1
		})

	/// @func radians(degrees) -> float
	/// @arg degrees {float}
	/// @returns {float} - Radians.
	lib.CreateFunction(tab, "radians",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "deg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			rad := gg.Radians(args["deg"].(float64))
			state.Push(golua.LNumber(rad))
			return 1
		})

	/// @func point(x, y) -> struct<context.Point>
	/// @arg x {float}
	/// @arg y {float}
	/// returns {struct<context.Point>}
	lib.CreateFunction(tab, "point",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Point
			/// @prop x {float}
			/// @prop y {float}

			t := state.NewTable()

			state.SetField(t, "x", golua.LNumber(args["x"].(float64)))
			state.SetField(t, "y", golua.LNumber(args["y"].(float64)))

			state.Push(t)
			return 1
		})

	/// @func distance(p1, p2) -> float
	/// @arg p1 {struct<context.Point>}
	/// @arg p2 {struct<context.Point>}
	/// @returns {float}
	lib.CreateFunction(tab, "distance",
		[]lua.Arg{
			{Type: lua.TABLE, Name: "p1", Table: &[]lua.Arg{
				{Type: lua.FLOAT, Name: "x"},
				{Type: lua.FLOAT, Name: "y"},
			}},
			{Type: lua.TABLE, Name: "p2", Table: &[]lua.Arg{
				{Type: lua.FLOAT, Name: "x"},
				{Type: lua.FLOAT, Name: "y"},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ap1 := args["p1"].(Point)
			ap2 := args["b2"].(Point)

			p1 := gg.Point{X: ap1["x"], Y: ap1["y"]}
			p2 := gg.Point{X: ap2["x"], Y: ap2["y"]}

			dist := p1.Distance(p2)

			state.Push(golua.LNumber(dist))
			return 1
		})

	/// @func interpolate(p1, p2, t) -> struct<context.Point>
	/// @arg p1 {struct<context.Point>}
	/// @arg p2 {struct<context.Point>}
	/// @arg t {float}
	/// @returns {struct<context.Point>}
	lib.CreateFunction(tab, "interpolate",
		[]lua.Arg{
			{Type: lua.TABLE, Name: "p1", Table: &[]lua.Arg{
				{Type: lua.FLOAT, Name: "x"},
				{Type: lua.FLOAT, Name: "y"},
			}},
			{Type: lua.TABLE, Name: "p2", Table: &[]lua.Arg{
				{Type: lua.FLOAT, Name: "x"},
				{Type: lua.FLOAT, Name: "y"},
			}},
			{Type: lua.FLOAT, Name: "t"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ap1 := args["p1"].(Point)
			ap2 := args["b2"].(Point)

			p1 := gg.Point{X: ap1["x"], Y: ap1["y"]}
			p2 := gg.Point{X: ap2["x"], Y: ap2["y"]}

			pi := p1.Interpolate(p2, args["t"].(float64))

			t := state.NewTable()

			state.SetField(t, "x", golua.LNumber(pi.X))
			state.SetField(t, "y", golua.LNumber(pi.Y))

			state.Push(t)
			return 1
		})

	/// @func new(width, height) -> int<collection.CONTEXT>
	/// @arg width {int}
	/// @arg height {int}
	/// returns {int<collection.CONTEXT>}
	lib.CreateFunction(tab, "new",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := fmt.Sprintf("context_%d", r.CC.Next())

			chLog := log.NewLogger(name, lg)
			lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

			id := r.CC.AddItem(&chLog)

			r.CC.Schedule(id, &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					c := gg.NewContext(args["width"].(int), args["height"].(int))
					i.Self.Context = c
					i.Lg.Append("new context created", log.LEVEL_INFO)
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func new_image(id) -> int<collection.CONTEXT>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {int<collection.CONTEXT>}
	lib.CreateFunction(tab, "new_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imageFinish := make(chan struct{}, 2)
			imageReady := make(chan struct{}, 2)
			var img image.Image

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imageReady <- struct{}{}
					<-imageFinish
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imageReady <- struct{}{}
				},
			})

			tempName := fmt.Sprintf("context_%d", r.CC.Next())

			chLog := log.NewLogger(tempName, lg)
			lg.Append(fmt.Sprintf("child log created: %s", tempName), log.LEVEL_INFO)

			id := r.CC.AddItem(&chLog)

			r.CC.Schedule(id, &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					<-imageReady

					i.Self = &collection.ItemContext{
						Context: gg.NewContextForImage(img),
					}

					imageFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					imageFinish <- struct{}{}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func new_direct(id) -> int<collection.CONTEXT>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {int<collection.CONTEXT>}
	/// @desc
	/// Creates a new context directly on the image,
	/// this requires the image to use the RGBA color model.
	/// If not it will be converted to RGBA.
	/// This is also not thread safe, as modifying either the image or the context will affect the other,
	/// with no guarantee of the order of operations.
	lib.CreateFunction(tab, "new_direct",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imageFinish := make(chan struct{})
			imageReady := make(chan struct{})
			var img *image.RGBA

			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					if rgba, ok := i.Self.Image.(*image.RGBA); ok {
						img = rgba
					} else {
						rgba := imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
						i.Self.Image = rgba
						img = rgba.(*image.RGBA)
					}
					imageReady <- struct{}{}
					<-imageFinish
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					close(imageReady)
				},
			})

			tempName := fmt.Sprintf("context_%d", r.CC.Next())

			chLog := log.NewLogger(tempName, lg)
			lg.Append(fmt.Sprintf("child log created: %s", tempName), log.LEVEL_INFO)

			id := r.CC.AddItem(&chLog)

			r.CC.Schedule(id, &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					<-imageReady

					i.Self = &collection.ItemContext{
						Context: gg.NewContextForRGBA(img),
					}

					imageFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					close(imageFinish)
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func to_image(id, name, encoding, model?, copy?) -> int<collection.IMAGE>
	/// @arg id {int<collection.CONTEXT>}
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg? model {int<image.Model>}
	/// @arg? copy {bool} - Set to true to copy the image, otherwise continuing to draw can affect the image.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
			{Type: lua.BOOL, Name: "copy", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			contextFinish := make(chan struct{}, 2)
			contextReady := make(chan struct{}, 2)

			var context *gg.Context

			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					context = i.Self.Context
					contextReady <- struct{}{}
					<-contextFinish
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					contextReady <- struct{}{}
				},
			})

			chLog := log.NewLogger(args["name"].(string), lg)
			lg.Append(fmt.Sprintf("child log created: %s", args["name"].(string)), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-contextReady

					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

					img := context.Image()
					if args["copy"].(bool) {
						img = imageutil.CopyImage(img, model)
					}

					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     args["name"].(string),
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Model:    model,
					}

					contextFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					contextFinish <- struct{}{}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func to_mask(id, name, encoding) -> int<collection.IMAGE>
	/// @arg id {int<collection.CONTEXT>}
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @returns {int<collection.IMAGE>} - The returned image will use the 'image.MODEL_ALPHA' color model.
	lib.CreateFunction(tab, "to_mask",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			contextFinish := make(chan struct{}, 2)
			contextReady := make(chan struct{}, 2)

			var context *gg.Context

			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					context = i.Self.Context
					contextReady <- struct{}{}
					<-contextFinish
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					contextReady <- struct{}{}
				},
			})

			chLog := log.NewLogger(args["name"].(string), lg)
			lg.Append(fmt.Sprintf("child log created: %s", args["name"].(string)), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-contextReady

					img := context.AsMask()

					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     args["name"].(string),
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Model:    imageutil.MODEL_ALPHA,
					}

					contextFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					contextFinish <- struct{}{}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func mask(id, img)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg img {int<collection.IMAGE>}
	/// @desc
	/// The image will be copied and converted to 'image.MODEL_ALPHA'.
	lib.CreateFunction(tab, "mask",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "img"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgReady := make(chan struct{}, 2)

			var img *image.Alpha

			r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = imageutil.CopyImage(img, imageutil.MODEL_ALPHA).(*image.Alpha)
					imgReady <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					<-imgReady
					err := i.Self.Context.SetMask(img)
					if err != nil {
						state.Error(golua.LString(lg.Append("failed to set image mask, image may be the wrong size.", log.LEVEL_ERROR)), 0)
					}
				},
			})
			return 0
		})

	/// @func size(id) -> int, int
	/// @arg id {int<collection.CONTEXT>}
	/// @returns {int} - The width of the context.
	/// @returns {int} - The height of the context.
	/// @blocking
	lib.CreateFunction(tab, "size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := 0
			height := 0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					width = i.Self.Context.Width()
					height = i.Self.Context.Height()
				},
			})

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func font_height(id) -> float
	/// @arg id {int<collection.CONTEXT>}
	/// @returns {float} - The height of text rendered with the current font.
	/// @blocking
	lib.CreateFunction(tab, "font_height",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					height = i.Self.Context.FontHeight()
				},
			})

			state.Push(golua.LNumber(height))
			return 1
		})

	/// @func string_measure(id, str) -> float, float
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @returns {float} - The width of text rendered with the current font.
	/// @returns {float} - The height of text rendered with the current font.
	/// @blocking
	lib.CreateFunction(tab, "string_measure",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := 0.0
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					width, height = i.Self.Context.MeasureString(args["str"].(string))
				},
			})

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func string_measure_multiline(id, str, spacing) -> float, float
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @arg spacing {float} - The space between each line.
	/// @returns {float} - The width of text rendered with the current font.
	/// @returns {float} - The height of text rendered with the current font.
	/// @blocking
	lib.CreateFunction(tab, "string_measure_multiline",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "spacing"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := 0.0
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					width, height = i.Self.Context.MeasureMultilineString(args["str"].(string), args["spacing"].(float64))
				},
			})

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func current_point(id) -> float, float, bool
	/// @arg id {int<collection.CONTEXT>}
	/// @returns {float} - The x location of the current point.
	/// @returns {float} - The y location of the current point.
	/// @returns {bool} - If the current point has been set.
	/// @blocking
	lib.CreateFunction(tab, "current_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x := 0.0
			y := 0.0
			exists := false

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					point, e := i.Self.Context.GetCurrentPoint()
					x = point.X
					y = point.Y
					exists = e
				},
			})

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			state.Push(golua.LBool(exists))
			return 3
		})

	/// @func transform_point(id, x, y) -> float, float
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @returns {float} -> New x position.
	/// @returns {float} -> New y position.
	/// @blocking
	/// @desc
	/// Multiplies a point by the current matrix.
	lib.CreateFunction(tab, "transform_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x := 0.0
			y := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					x, y = i.Self.Context.TransformPoint(args["x"].(float64), args["y"].(float64))
				},
			})

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func clear(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Fills the context with the current color.
	lib.CreateFunction(tab, "clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Clear()
				},
			})

			return 0
		})

	/// @func clip(id, preserve?)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg? preserve {bool} - Set in order to keep the current path.
	/// @desc
	/// Updates the clipping region by intersecting the current clipping region with the current path as it would be filled by fill().
	lib.CreateFunction(tab, "clip",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					if args["preserve"].(bool) {
						i.Self.Context.ClipPreserve()
					} else {
						i.Self.Context.Clip()
					}
				},
			})

			return 0
		})

	/// @func clip_reset(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Clears the clipping region.
	lib.CreateFunction(tab, "clip_reset",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.ResetClip()
				},
			})

			return 0
		})

	/// @func path_clear(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Removes all points from the current path.
	lib.CreateFunction(tab, "path_clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.ClearPath()
				},
			})

			return 0
		})

	/// @func path_close(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Adds a line segment from the current point to the initial point.
	lib.CreateFunction(tab, "path_close",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.ClosePath()
				},
			})

			return 0
		})

	/// @func path_to(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Starts a new subpath starting at the given point.
	lib.CreateFunction(tab, "path_to",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.MoveTo(
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func subpath(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Starts a new subpath starting at the current point.
	/// No current point will be set after.
	lib.CreateFunction(tab, "subpath",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.NewSubPath()
				},
			})

			return 0
		})

	/// @func draw_cubic(id, x1, y1, x2, y2, x3, y3)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg x2 {flaot}
	/// @arg y2 {float}
	/// @arg x3 {float}
	/// @arg y3 {float}
	/// @desc
	/// Draws a cubic bezier curve to the path starting at the current point,
	/// if there isn't a current point, it moves to (x1, y1).
	lib.CreateFunction(tab, "draw_cubic",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
			{Type: lua.FLOAT, Name: "x3"},
			{Type: lua.FLOAT, Name: "y3"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.CubicTo(
						args["x1"].(float64),
						args["y1"].(float64),
						args["x2"].(float64),
						args["y2"].(float64),
						args["x3"].(float64),
						args["y3"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_quadratic(id, x1, y1, x2, y2)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg x2 {float}
	/// @arg y2 {float}
	/// @desc
	/// Draws a quadratic bezier curve to the path starting at the current point,
	/// if there isn't a current point, it moves to (x1, y1).
	lib.CreateFunction(tab, "draw_quadratic",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.QuadraticTo(
						args["x1"].(float64),
						args["y1"].(float64),
						args["x2"].(float64),
						args["y2"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_arc(id, x, y, r, angle1, angle2)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg r {float}
	/// @arg angle1 {float}
	/// @arg angle2 {float}
	lib.CreateFunction(tab, "draw_arc",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "angle1"},
			{Type: lua.FLOAT, Name: "angle2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawArc(
						args["x"].(float64),
						args["y"].(float64),
						args["r"].(float64),
						args["angle1"].(float64),
						args["angle2"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_circle(id, x, y, r)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg r {float}
	lib.CreateFunction(tab, "draw_circle",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawCircle(
						args["x"].(float64),
						args["y"].(float64),
						args["r"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_ellipse(id, x, y, rx, xy)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg rx {float}
	/// @arg ry {float}
	lib.CreateFunction(tab, "draw_ellipse",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "rx"},
			{Type: lua.FLOAT, Name: "ry"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawEllipse(
						args["x"].(float64),
						args["y"].(float64),
						args["rx"].(float64),
						args["ry"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_elliptical_arc(id, x, y, rx, ry, angle1, angle2)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg rx {float}
	/// @arg ry {float}
	/// @arg angle1 {float}
	/// @arg angle2 {float}
	lib.CreateFunction(tab, "draw_elliptical_arc",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "rx"},
			{Type: lua.FLOAT, Name: "ry"},
			{Type: lua.FLOAT, Name: "angle1"},
			{Type: lua.FLOAT, Name: "angle2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawEllipticalArc(
						args["x"].(float64),
						args["y"].(float64),
						args["rx"].(float64),
						args["ry"].(float64),
						args["angle1"].(float64),
						args["angle2"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_image(id, img, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg img {int<collection.IMAGE>}
	/// @arg x {int}
	/// @arg y {int}
	lib.CreateFunction(tab, "draw_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgFinish := make(chan struct{}, 2)
			imgReady := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinish
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					<-imgReady
					i.Self.Context.DrawImage(img, args["x"].(int), args["y"].(int))
					imgFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					imgFinish <- struct{}{}
				},
			})
			return 0
		})

	/// @func draw_image_anchor(id, img, x, y, ax, ay)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg img {int<collection.IMAGE>}
	/// @arg x {int}
	/// @arg y {int}
	/// @arg ax {float}
	/// @arg ay {float}
	/// @desc
	/// The anchor is a point between (0,0) and (1,1), so (0.5,0.5) is centered.
	lib.CreateFunction(tab, "draw_image_anchor",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.FLOAT, Name: "ax"},
			{Type: lua.FLOAT, Name: "ay"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgFinish := make(chan struct{}, 2)
			imgReady := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["img"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinish
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					<-imgReady
					i.Self.Context.DrawImageAnchored(img, args["x"].(int), args["y"].(int), args["ax"].(float64), args["ay"].(float64))
					imgFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					imgFinish <- struct{}{}
				},
			})
			return 0
		})

	/// @func draw_line(id, x1, y1, x2, y2)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg x2 {float}
	/// @arg y2 {float}
	lib.CreateFunction(tab, "draw_line",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawLine(
						args["x1"].(float64),
						args["y1"].(float64),
						args["x2"].(float64),
						args["y2"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_line_to(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Draws a line to the point, starting from the current point.
	/// If there is no current point, no line will be drawn and the current point will be set.
	lib.CreateFunction(tab, "draw_line_to",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.LineTo(
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_point(id, x, y, r)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg r {float}
	/// @desc
	/// Similar to draw_circle but ensures that a circle of the specified size is drawn regardless of the current transformation matrix.
	/// The position is still transformed, but not the shape of the point.
	lib.CreateFunction(tab, "draw_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawPoint(
						args["x"].(float64),
						args["y"].(float64),
						args["r"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_rect(id, x, y, width, height)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "draw_rect",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawRectangle(
						args["x"].(float64),
						args["y"].(float64),
						args["width"].(float64),
						args["height"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_rect_round(id, x, y, width, height, r)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg width {float}
	/// @arg height {float}
	/// @arg r {float}
	lib.CreateFunction(tab, "draw_rect_round",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawRoundedRectangle(
						args["x"].(float64),
						args["y"].(float64),
						args["width"].(float64),
						args["height"].(float64),
						args["r"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_polygon(id, n, x, y, r, rotation)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg n {int} - The number of sides.
	/// @arg x {float}
	/// @arg y {float}
	/// @arg r {float}
	/// @arg rotation {float}
	lib.CreateFunction(tab, "draw_polygon",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "n"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "rotation"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawRegularPolygon(
						args["n"].(int),
						args["x"].(float64),
						args["y"].(float64),
						args["r"].(float64),
						args["rotation"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_string(id, str, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @arg x {float}
	/// @arg y {float}
	lib.CreateFunction(tab, "draw_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawString(
						args["str"].(string),
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_string_anchor(id, str, x, y, ax, ay)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg ax {float}
	/// @arg ay {float}
	/// @desc
	/// The anchor is a point between (0,0) and (1,1), so (0.5,0.5) is centered.
	lib.CreateFunction(tab, "draw_string_anchor",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "ax"},
			{Type: lua.FLOAT, Name: "ay"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawStringAnchored(
						args["str"].(string),
						args["x"].(float64),
						args["y"].(float64),
						args["ax"].(float64),
						args["ay"].(float64),
					)
				},
			})

			return 0
		})

	/// @func draw_string_wrap(id, str, x, y, ax, ay, width, spacing, align)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg ax {float}
	/// @arg ay {float}
	/// @arg width {float}
	/// @arg spacing {float} - The spacing between each line.
	/// @arg align {int<context.Align>}
	/// @desc
	/// The anchor is a point between (0,0) and (1,1), so (0.5,0.5) is centered.
	lib.CreateFunction(tab, "draw_string_wrap",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "ax"},
			{Type: lua.FLOAT, Name: "ay"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "spacing"},
			{Type: lua.INT, Name: "align"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.DrawStringWrapped(
						args["str"].(string),
						args["x"].(float64),
						args["y"].(float64),
						args["ax"].(float64),
						args["ay"].(float64),
						args["width"].(float64),
						args["spacing"].(float64),
						lua.ParseEnum(args["align"].(int), alignment, lib),
					)
				},
			})

			return 0
		})

	/// @func fill(id, preserve?)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg? preserve {bool} - Set in order to keep the current path.
	/// @desc
	/// Fills the current path with the current color.
	/// Closes open paths.
	lib.CreateFunction(tab, "fill",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					if args["preserve"].(bool) {
						i.Self.Context.FillPreserve()
					} else {
						i.Self.Context.Fill()
					}
				},
			})

			return 0
		})

	/// @func fill_rule(id, rule)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg rule {int<context.FillRule>}
	lib.CreateFunction(tab, "fill_rule",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "rule"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetFillRule(lua.ParseEnum(args["rule"].(int), fillRules, lib))
				},
			})

			return 0
		})

	/// @func stroke(id, preserve?)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg? preserve {bool} - Set in order to keep the current path.
	/// @desc
	/// Strokes the current path with the current color.
	lib.CreateFunction(tab, "stroke",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					if args["preserve"].(bool) {
						i.Self.Context.StrokePreserve()
					} else {
						i.Self.Context.Stroke()
					}
				},
			})

			return 0
		})

	/// @func identity(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Resets the current transformation matrix to the identity matrix.
	lib.CreateFunction(tab, "identity",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Identity()
				},
			})

			return 0
		})

	/// @func mask_invert(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Inverts the alpha values of the clipping mask.
	lib.CreateFunction(tab, "mask_invert",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.InvertMask()
				},
			})

			return 0
		})

	/// @func invert_y(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Flips the y axis so that Y=0 is at the bottom of the image.
	lib.CreateFunction(tab, "invert_y",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.InvertY()
				},
			})

			return 0
		})

	/// @func push(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Push the current context state to the stack.
	lib.CreateFunction(tab, "push",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Push()
				},
			})

			return 0
		})

	/// @func pop(id)
	/// @arg id {int<collection.CONTEXT>}
	/// @desc
	/// Pop the current context state to the stack.
	lib.CreateFunction(tab, "pop",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Pop()
				},
			})

			return 0
		})

	/// @func rotate(id, angle)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg angle {float}
	/// @desc
	/// Rotates the transformation matrix around the origin.
	lib.CreateFunction(tab, "rotate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "angle"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Rotate(args["angle"].(float64))
				},
			})

			return 0
		})

	/// @func rotate_about(id, angle, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg angle {float}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Rotates the transformation matrix around the point.
	lib.CreateFunction(tab, "rotate_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.RotateAbout(
						args["angle"].(float64),
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func scale(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Scales the transformation matrix by a factor.
	lib.CreateFunction(tab, "scale",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Scale(
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func scale_about(id, sx, sy, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg sx {float}
	/// @arg sy {float}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Scales the transformation matrix by a factor (sx,sy) starting at the point (x,y).
	lib.CreateFunction(tab, "scale_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "sx"},
			{Type: lua.FLOAT, Name: "sy"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.ScaleAbout(
						args["sx"].(float64),
						args["sy"].(float64),
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func color_hex(id, hex)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg hex {string}
	/// @desc
	/// Supports hex colors in the follow formats: #RGB #RRGGBB #RRGGBBAA.
	/// The leading # is optional.
	lib.CreateFunction(tab, "color_hex",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "hex"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetHexColor(args["hex"].(string))
				},
			})

			return 0
		})

	/// @func color(id, color)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "color",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
					i.Self.Context.SetColor(col)
				},
			})

			return 0
		})

	/// @func color_rgb(id, r, g, b)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg r {float}
	/// @arg g {float}
	/// @arg b {float}
	/// @desc
	/// Float values for r,g,b between 0 and 1.
	lib.CreateFunction(tab, "color_rgb",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "g"},
			{Type: lua.FLOAT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetRGB(
						args["r"].(float64),
						args["g"].(float64),
						args["b"].(float64),
					)
				},
			})

			return 0
		})

	/// @func color_rgb255(id, r, g, b)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg r {int}
	/// @arg g {int}
	/// @arg b {int}
	/// @desc
	/// Interger values for r,g,b between 0 and 255.
	lib.CreateFunction(tab, "color_rgb255",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetRGB255(
						args["r"].(int),
						args["g"].(int),
						args["b"].(int),
					)
				},
			})

			return 0
		})

	/// @func color_rgba(id, r, g, b, a)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg r {float}
	/// @arg g {float}
	/// @arg b {float}
	/// @arg a {float}
	/// @desc
	/// Float values for r,g,b,a between 0 and 1.
	lib.CreateFunction(tab, "color_rgba",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "g"},
			{Type: lua.FLOAT, Name: "b"},
			{Type: lua.FLOAT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetRGBA(
						args["r"].(float64),
						args["g"].(float64),
						args["b"].(float64),
						args["a"].(float64),
					)
				},
			})

			return 0
		})

	/// @func color_rgba255(id, r, g, b, a)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg r {int}
	/// @arg g {int}
	/// @arg b {int}
	/// @arg a {int}
	/// @desc
	/// Interger values for r,g,b,a between 0 and 255.
	lib.CreateFunction(tab, "color_rgba255",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetRGBA255(
						args["r"].(int),
						args["g"].(int),
						args["b"].(int),
						args["a"].(int),
					)
				},
			})

			return 0
		})

	/// @func dash_set(id, pattern)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg pattern {[]float} - Each float value is a dash with a length of the value.
	/// @desc
	/// Sets the dash pattern to use.
	/// Call with empty array to disable dashes.
	lib.CreateFunction(tab, "dash_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("pattern", lua.ArrayType{Type: lua.FLOAT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					pattern := []float64{}
					for _, v := range args["pattern"].(map[string]any) {
						pattern = append(pattern, v.(float64))
					}

					i.Self.Context.SetDash(pattern...)
				},
			})

			return 0
		})

	/// @func dash_set_offset(id, offset)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg offset {float}
	/// @desc
	/// The initial offset for the dash pattern.
	lib.CreateFunction(tab, "dash_set_offset",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "offset"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetDashOffset(args["offset"].(float64))
				},
			})

			return 0
		})

	/// @func line_cap(id, cap)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg cap {int<context.LineCap>}
	lib.CreateFunction(tab, "line_cap",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "cap"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetLineCap(lua.ParseEnum(args["cap"].(int), lineCaps, lib))
				},
			})

			return 0
		})

	/// @func line_join(id, join)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg join {int<context.LineJoin>}
	lib.CreateFunction(tab, "line_join",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "join"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetLineJoin(lua.ParseEnum(args["join"].(int), lineJoins, lib))
				},
			})

			return 0
		})

	/// @func line_width(id, width)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg width {float}
	lib.CreateFunction(tab, "line_width",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "width"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetLineWidth(args["width"].(float64))
				},
			})

			return 0
		})

	/// @func pixel_set(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {int}
	/// @arg y {int}
	/// @desc
	/// Sets a pixel to the current color.
	lib.CreateFunction(tab, "pixel_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetPixel(args["x"].(int), args["y"].(int))
				},
			})

			return 0
		})

	/// @func shear(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Updates the current matrix with a shearing angle, at the origin.
	lib.CreateFunction(tab, "shear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Shear(args["x"].(float64), args["y"].(float64))
				},
			})

			return 0
		})

	/// @func shear_about(id, sx, sy, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg sx {float}
	/// @arg sy {float}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Updates the current matrix with a shearing angle (sx,sy), at the given point (x,y).
	lib.CreateFunction(tab, "shear_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "sx"},
			{Type: lua.FLOAT, Name: "sy"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.ShearAbout(
						args["sx"].(float64),
						args["sy"].(float64),
						args["x"].(float64),
						args["y"].(float64),
					)
				},
			})

			return 0
		})

	/// @func translate(id, x, y)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg x {float}
	/// @arg y {float}
	/// @desc
	/// Updates the current matrix with a translation.
	lib.CreateFunction(tab, "translate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Translate(args["x"].(float64), args["y"].(float64))
				},
			})

			return 0
		})

	/// @func word_wrap(id, str, width) -> []string
	/// @arg id {int<collection.CONTEXT>}
	/// @arg str {string}
	/// @arg width {float}
	/// @returns {[]string} - The original string split into line breaks at the set width.
	/// @blocking
	lib.CreateFunction(tab, "word_wrap",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "width"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var wrappedStrings []string

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					wrappedStrings = i.Self.Context.WordWrap(args["str"].(string), args["width"].(float64))
				},
			})

			t := state.NewTable()
			for ind, str := range wrappedStrings {
				t.RawSetInt(ind+1, golua.LString(str))
			}

			state.Push(t)
			return 1
		})

	/// @func matrix_new(xx, yx, xy, yy, x0, y0) -> struct<context.Matrix>
	/// @arg xx {float}
	/// @arg yx {float}
	/// @arg xy {float}
	/// @arg yy {float}
	/// @arg x0 {float}
	/// @arg y0 {float}
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_new",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "xx"},
			{Type: lua.FLOAT, Name: "yx"},
			{Type: lua.FLOAT, Name: "xy"},
			{Type: lua.FLOAT, Name: "yy"},
			{Type: lua.FLOAT, Name: "x0"},
			{Type: lua.FLOAT, Name: "y0"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := matrixTable(state,
				args["xx"].(float64),
				args["yx"].(float64),
				args["xy"].(float64),
				args["yy"].(float64),
				args["x0"].(float64),
				args["y0"].(float64),
			)

			state.Push(t)
			return 1
		})

	/// @func matrix_identity() -> struct<context.Matrix>
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_identity",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			m := gg.Identity()
			t := matrixTable(state,
				m.XX,
				m.YX,
				m.XY,
				m.YY,
				m.X0,
				m.Y0,
			)

			state.Push(t)
			return 1
		})

	/// @func matrix_rotate(angle) -> struct<context.Matrix>
	/// @arg angle {float}
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_rotate",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "angle"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			m := gg.Rotate(args["angle"].(float64))
			t := matrixTable(state,
				m.XX,
				m.YX,
				m.XY,
				m.YY,
				m.X0,
				m.Y0,
			)

			state.Push(t)
			return 1
		})

	/// @func matrix_scale(x, y) -> struct<context.Matrix>
	/// @arg x {float}
	/// @arg y {float}
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_scale",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			m := gg.Scale(args["x"].(float64), args["y"].(float64))
			t := matrixTable(state,
				m.XX,
				m.YX,
				m.XY,
				m.YY,
				m.X0,
				m.Y0,
			)

			state.Push(t)
			return 1
		})

	/// @func matrix_shear(x, y) -> struct<context.Matrix>
	/// @arg x {float}
	/// @arg y {float}
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_shear",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			m := gg.Shear(args["x"].(float64), args["y"].(float64))
			t := matrixTable(state,
				m.XX,
				m.YX,
				m.XY,
				m.YY,
				m.X0,
				m.Y0,
			)

			state.Push(t)
			return 1
		})

	/// @func matrix_translate(x, y) -> struct<context.Matrix>
	/// @arg x {float}
	/// @arg y {float}
	/// @returns {struct<context.Matrix>}
	lib.CreateFunction(tab, "matrix_translate",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			m := gg.Translate(args["x"].(float64), args["y"].(float64))
			t := matrixTable(state,
				m.XX,
				m.YX,
				m.XY,
				m.YY,
				m.X0,
				m.Y0,
			)

			state.Push(t)
			return 1
		})

	/// @func point_cubic_bezier(x0, y0, x1, y1, x2, y2, x3, y3) -> []struct<context.Point>
	/// @arg x0 {float}
	/// @arg y0 {float}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg x2 {float}
	/// @arg y2 {float}
	/// @arg x3 {float}
	/// @arg y3 {float}
	/// @returns {[]struct<context.Point>}
	lib.CreateFunction(tab, "point_cubic_bezier",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x0"},
			{Type: lua.FLOAT, Name: "y0"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
			{Type: lua.FLOAT, Name: "x3"},
			{Type: lua.FLOAT, Name: "y3"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			points := gg.CubicBezier(
				args["x0"].(float64),
				args["y0"].(float64),
				args["x1"].(float64),
				args["y1"].(float64),
				args["x2"].(float64),
				args["y2"].(float64),
				args["x3"].(float64),
				args["y3"].(float64),
			)

			t := state.NewTable()
			for ind, p := range points {
				point := state.NewTable()
				point.RawSetString("x", golua.LNumber(p.X))
				point.RawSetString("y", golua.LNumber(p.Y))

				t.RawSetInt(ind+1, point)
			}

			state.Push(t)
			return 1
		})

	/// @func point_quadratic_bezier(x0, y0, x1, y1, x2, y2) -> []struct<context.Point>
	/// @arg x0 {float}
	/// @arg y0 {float}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg x2 {float}
	/// @arg y2 {float}
	/// @returns {[]struct<context.Point>}
	lib.CreateFunction(tab, "point_quadratic_bezier",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x0"},
			{Type: lua.FLOAT, Name: "y0"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
			{Type: lua.FLOAT, Name: "x3"},
			{Type: lua.FLOAT, Name: "y3"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			points := gg.QuadraticBezier(
				args["x0"].(float64),
				args["y0"].(float64),
				args["x1"].(float64),
				args["y1"].(float64),
				args["x2"].(float64),
				args["y2"].(float64),
			)

			t := state.NewTable()
			for ind, p := range points {
				point := state.NewTable()
				point.RawSetString("x", golua.LNumber(p.X))
				point.RawSetString("y", golua.LNumber(p.Y))

				t.RawSetInt(ind+1, point)
			}

			state.Push(t)
			return 1
		})

	/// @func pattern_solid(color) -> struct<context.PatternSolid>
	/// @arg color {struct<image.Color>}
	/// @returns {struct<context.PatternSolid>}
	lib.CreateFunction(tab, "pattern_solid",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternSolidTable(state, args["color"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func pattern_surface(id, repeat_op) -> struct<context.PatternSurface>
	/// @arg id {int<collection.IMAGE>}
	/// @arg repeat_op {int<context.RepeatOp>}
	/// @returns {struct<context.PatternSurface>}
	/// @desc
	/// Images are only grabbed when the pattern is set, not when this is called.
	lib.CreateFunction(tab, "pattern_surface",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "repeat_op"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternSurfaceTable(state, args["id"].(int), args["repeat_op"].(int))

			state.Push(t)
			return 1
		})

	/// @func pattern_surface_sync(id, repeat_op) -> struct<context.PatternSurfaceSync>
	/// @arg id {int<collection.IMAGE>}
	/// @arg repeat_op {int<context.RepeatOp>}
	/// @returns {struct<context.PatternSurfaceSync>}
	/// @desc
	/// This does not wait for the image to be ready or idle,
	/// if the image is not loaded it will grab an empy image.
	/// This also means it is not thread safe, it is unknown what state the image will be in when grabbed.
	/// Images are also only grabbed when the pattern is set, not when this is called.
	lib.CreateFunction(tab, "pattern_surface_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "repeat_op"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternSurfaceSyncTable(state, args["id"].(int), args["repeat_op"].(int))

			state.Push(t)
			return 1
		})

	/// @func pattern_custom(fn) -> struct<context.PatternCustom>
	/// @arg fn {function(x int, y int) -> struct<image.Color>}
	/// @returns {struct<context.PatternCustom>}
	lib.CreateFunction(tab, "pattern_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternCustomTable(state, args["fn"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func gradient_linear(x0, y0, x1, y1) -> struct<context.PatternGradientLinear>
	/// @arg x0 {float}
	/// @arg y0 {float}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @returns {struct<context.PatternGradientLinear>}
	lib.CreateFunction(tab, "gradient_linear",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x0"},
			{Type: lua.FLOAT, Name: "y0"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternGradientLinearTable(state, args["x0"].(float64), args["y0"].(float64), args["x1"].(float64), args["y1"].(float64))

			state.Push(t)
			return 1
		})

	/// @func gradient_radial(x0, y0, r0, x1, y1, r1) -> struct<context.PatternGradientRadial>
	/// @arg x0 {float}
	/// @arg y0 {float}
	/// @arg r0 {float}
	/// @arg x1 {float}
	/// @arg y1 {float}
	/// @arg r1 {float}
	/// @returns {struct<context.PatternGradientRadial>}
	lib.CreateFunction(tab, "gradient_radial",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x0"},
			{Type: lua.FLOAT, Name: "y0"},
			{Type: lua.FLOAT, Name: "r0"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "r1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternGradientRadialTable(state, args["x0"].(float64), args["y0"].(float64), args["r0"].(float64), args["x1"].(float64), args["y1"].(float64), args["r1"].(float64))

			state.Push(t)
			return 1
		})

	/// @func fill_style(id, pattern)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg pattern {struct<context.Pattern>}
	lib.CreateFunction(tab, "fill_style",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "pattern"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					pt := args["pattern"].(*golua.LTable)
					pattern := patternBuild(state, pt, r, lg)

					i.Self.Context.SetFillStyle(pattern)
				},
			})

			return 0
		})

	/// @func stroke_style(id, pattern)
	/// @arg id {int<collection.CONTEXT>}
	/// @arg pattern {struct<context.Pattern>}
	lib.CreateFunction(tab, "stroke_style",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "pattern"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					pt := args["pattern"].(*golua.LTable)
					pattern := patternBuild(state, pt, r, lg)

					i.Self.Context.SetStrokeStyle(pattern)
				},
			})

			return 0
		})

	/// @func font_load(path, points)
	/// @arg path {string}
	/// @arg points {float}
	lib.CreateFunction(tab, "font_load",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
			{Type: lua.FLOAT, Name: "points"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.LoadFontFace(args["path"].(string), args["points"].(float64))
				},
			})

			return 0
		})

	/// @constants Fill Rules
	/// @const FILLRULE_WINDING
	/// @const FILLRULE_EVENODD
	tab.RawSetString("FILLRULE_WINDING", golua.LNumber(gg.FillRuleWinding))
	tab.RawSetString("FILLRULE_EVENODD", golua.LNumber(gg.FillRuleEvenOdd))

	/// @constants Line Caps
	/// @const LINECAP_ROUND
	/// @const LINECAP_BUTT
	/// @const LINECAP_SQUARE
	tab.RawSetString("LINECAP_ROUND", golua.LNumber(gg.LineCapRound))
	tab.RawSetString("LINECAP_BUTT", golua.LNumber(gg.LineCapButt))
	tab.RawSetString("LINECAP_SQUARE", golua.LNumber(gg.LineCapSquare))

	/// @constants Line Joins
	/// @const LINEJOIN_ROUND
	/// @const LINEJOIN_BEVEL
	tab.RawSetString("LINEJOIN_ROUND", golua.LNumber(gg.LineJoinRound))
	tab.RawSetString("LINEJOIN_BEVEL", golua.LNumber(gg.LineJoinBevel))

	/// @constants Repeat Ops
	/// @const REPEAT_BOTH
	/// @const REPEAT_X
	/// @const REPEAT_Y
	/// @const REPEAT_NONE
	tab.RawSetString("REPEAT_BOTH", golua.LNumber(gg.RepeatBoth))
	tab.RawSetString("REPEAT_X", golua.LNumber(gg.RepeatX))
	tab.RawSetString("REPEAT_Y", golua.LNumber(gg.RepeatY))
	tab.RawSetString("REPEAT_NONE", golua.LNumber(gg.RepeatNone))

	/// @constants Alignment
	/// @const ALIGN_LEFT
	/// @const ALIGN_CENTER
	/// @const ALIGN_RIGHT
	tab.RawSetString("ALIGN_LEFT", golua.LNumber(gg.AlignLeft))
	tab.RawSetString("ALIGN_CENTER", golua.LNumber(gg.AlignCenter))
	tab.RawSetString("ALIGN_RIGHT", golua.LNumber(gg.AlignRight))

	/// @constants Patterns
	/// @const PATTERN_SOLID
	/// @const PATTERN_SURFACE
	/// @const PATTERN_SURFACE_SYNC
	/// @const PATTERN_GRADIENT_LINEAR
	/// @const PATTERN_GRADIENT_RADIAL
	/// @const PATTERN_CUSTOM
	tab.RawSetString("PATTERN_SOLID", golua.LString(PATTERN_SOLID))
	tab.RawSetString("PATTERN_SURFACE", golua.LString(PATTERN_SURFACE))
	tab.RawSetString("PATTERN_SURFACE_SYNC", golua.LString(PATTERN_SURFACE_SYNC))
	tab.RawSetString("PATTERN_GRADIENT_LINEAR", golua.LString(PATTERN_GRADIENT_LINEAR))
	tab.RawSetString("PATTERN_GRADIENT_RADIAL", golua.LString(PATTERN_GRADIENT_RADIAL))
	tab.RawSetString("PATTERN_CUSTOM", golua.LString(PATTERN_CUSTOM))
}

var fillRules = []gg.FillRule{
	gg.FillRuleWinding,
	gg.FillRuleEvenOdd,
}

var lineCaps = []gg.LineCap{
	gg.LineCapRound,
	gg.LineCapButt,
	gg.LineCapSquare,
}

var lineJoins = []gg.LineJoin{
	gg.LineJoinRound,
	gg.LineJoinBevel,
}

var repeatOps = []gg.RepeatOp{
	gg.RepeatBoth,
	gg.RepeatX,
	gg.RepeatY,
	gg.RepeatNone,
}

var alignment = []gg.Align{
	gg.AlignLeft,
	gg.AlignCenter,
	gg.AlignRight,
}

const (
	PATTERN_SOLID           string = "solid"
	PATTERN_SURFACE         string = "surface"
	PATTERN_SURFACE_SYNC    string = "surface_sync"
	PATTERN_GRADIENT_LINEAR string = "gradient_linear"
	PATTERN_GRADIENT_RADIAL string = "gradient_radial"
	PATTERN_CUSTOM          string = "custom"
)

func matrixTable(state *golua.LState, xx, yx, xy, yy, x0, y0 float64) *golua.LTable {
	/// @struct Matrix
	/// @prop xx {float}
	/// @prop yx {float}
	/// @prop xy {float}
	/// @prop yy {float}
	/// @prop x0 {float}
	/// @prop y0 {float}
	/// @method multiply(struct<context.Matrix>) -> self
	/// @method rotate(angle float) -> self
	/// @method scale(x float, y float) -> self
	/// @method shear(x float, y float) -> self
	/// @method translate(x float, y float) -> self
	/// @method transform_point(x float, y float) -> float, float
	/// @method transform_vector(x float, y float) -> float, float

	t := state.NewTable()

	t.RawSetString("xx", golua.LNumber(xx))
	t.RawSetString("yx", golua.LNumber(yx))
	t.RawSetString("xy", golua.LNumber(xy))
	t.RawSetString("yy", golua.LNumber(yy))
	t.RawSetString("x0", golua.LNumber(x0))
	t.RawSetString("y0", golua.LNumber(y0))

	tableBuilderFunc(state, t, "multiply", func(state *golua.LState, t *golua.LTable) {
		mt := state.CheckTable(-1)
		m := matrixBuild(t)
		m2 := matrixBuild(mt)
		nm := m.Multiply(m2)
		matrixUpdate(t, nm)
	})

	tableBuilderFunc(state, t, "rotate", func(state *golua.LState, t *golua.LTable) {
		angle := state.CheckNumber(-1)
		m := matrixBuild(t)
		nm := m.Rotate(float64(angle))
		matrixUpdate(t, nm)
	})

	tableBuilderFunc(state, t, "scale", func(state *golua.LState, t *golua.LTable) {
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		m := matrixBuild(t)
		nm := m.Scale(float64(x), float64(y))
		matrixUpdate(t, nm)
	})

	tableBuilderFunc(state, t, "shear", func(state *golua.LState, t *golua.LTable) {
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		m := matrixBuild(t)
		nm := m.Shear(float64(x), float64(y))
		matrixUpdate(t, nm)
	})

	tableBuilderFunc(state, t, "translate", func(state *golua.LState, t *golua.LTable) {
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		m := matrixBuild(t)
		nm := m.Translate(float64(x), float64(y))
		matrixUpdate(t, nm)
	})

	t.RawSetString("transform_point", state.NewFunction(func(l *golua.LState) int {
		t := state.CheckTable(-3)
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		m := matrixBuild(t)
		tx, ty := m.TransformPoint(float64(x), float64(y))

		state.Push(golua.LNumber(tx))
		state.Push(golua.LNumber(ty))
		return 2
	}))

	t.RawSetString("transform_vector", state.NewFunction(func(l *golua.LState) int {
		t := state.CheckTable(-3)
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		m := matrixBuild(t)
		tx, ty := m.TransformVector(float64(x), float64(y))

		state.Push(golua.LNumber(tx))
		state.Push(golua.LNumber(ty))
		return 2
	}))

	return t
}

func matrixUpdate(t *golua.LTable, m gg.Matrix) {
	t.RawSetString("xx", golua.LNumber(m.XX))
	t.RawSetString("yx", golua.LNumber(m.YX))
	t.RawSetString("xy", golua.LNumber(m.XY))
	t.RawSetString("yy", golua.LNumber(m.YY))
	t.RawSetString("x0", golua.LNumber(m.X0))
	t.RawSetString("y0", golua.LNumber(m.Y0))
}

func matrixBuild(t *golua.LTable) gg.Matrix {
	xx := float64(t.RawGetString("xx").(golua.LNumber))
	yx := float64(t.RawGetString("yx").(golua.LNumber))
	xy := float64(t.RawGetString("xy").(golua.LNumber))
	yy := float64(t.RawGetString("yy").(golua.LNumber))
	x0 := float64(t.RawGetString("x0").(golua.LNumber))
	y0 := float64(t.RawGetString("y0").(golua.LNumber))

	return gg.Matrix{
		XX: xx,
		YX: yx,
		XY: xy,
		YY: yy,
		X0: x0,
		Y0: y0,
	}
}

func patternBuild(state *golua.LState, t *golua.LTable, r *lua.Runner, lg *log.Logger) gg.Pattern {
	/// @struct Pattern
	/// @prop type {string<context.Pattern>}

	typ := t.RawGetString("type").(golua.LString)

	switch string(typ) {
	case PATTERN_SOLID:
		return patternSolidBuild(t)
	case PATTERN_SURFACE:
		return patternSurfaceBuild(t, r)
	case PATTERN_SURFACE_SYNC:
		return patternSurfaceSyncBuild(t, r)
	case PATTERN_GRADIENT_LINEAR:
		return patternGradientLinearBuild(t)
	case PATTERN_GRADIENT_RADIAL:
		return patternGradientRadialBuild(t)
	case PATTERN_CUSTOM:
		return patternCustomBuild(state, t)
	}

	state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown pattern type: %s", typ), log.LEVEL_ERROR)), 0)
	return gg.NewSolidPattern(color.RGBA{})
}

func patternSolidTable(state *golua.LState, color *golua.LTable) *golua.LTable {
	/// @struct PatternSolid
	/// @prop type {string<context.Pattern>}
	/// @prop color {struct<image.Color>}

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_SOLID))
	t.RawSetString("color", color)

	return t
}

func patternSolidBuild(t *golua.LTable) gg.Pattern {
	ct := t.RawGetString("color").(*golua.LTable)
	col := imageutil.ColorTableToRGBAColor(ct)

	p := gg.NewSolidPattern(col)
	return p
}

func patternSurfaceTable(state *golua.LState, id, repeatOp int) *golua.LTable {
	/// @struct PatternSurface
	/// @prop type {string<context.Pattern>}
	/// @prop id {int<collection.IMAGE>}
	/// @prop repeatOp {int<context.RepeatOp>}

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_SURFACE))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("repeatOp", golua.LNumber(repeatOp))

	return t
}

func patternSurfaceBuild(t *golua.LTable, r *lua.Runner) gg.Pattern {
	id := t.RawGetString("id").(golua.LNumber)
	repeatOp := t.RawGetString("repeatOp").(golua.LNumber)

	var img image.Image

	<-r.IC.Schedule(int(id), &collection.Task[collection.ItemImage]{
		Lib:  LIB_CONTEXT,
		Name: "pattern_surface",
		Fn: func(i *collection.Item[collection.ItemImage]) {
			img = i.Self.Image
		},
	})

	p := gg.NewSurfacePattern(img, repeatOps[int(repeatOp)])
	return p
}

func patternSurfaceSyncTable(state *golua.LState, id, repeatOp int) *golua.LTable {
	/// @struct PatternSurfaceSync
	/// @prop type {string<context.Pattern>}
	/// @prop id {int<collection.IMAGE>}
	/// @prop repeatOp {int<context.RepeatOp>}

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_SURFACE_SYNC))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("repeatOp", golua.LNumber(repeatOp))

	return t
}

func patternSurfaceSyncBuild(t *golua.LTable, r *lua.Runner) gg.Pattern {
	id := t.RawGetString("id").(golua.LNumber)
	repeatOp := t.RawGetString("repeatOp").(golua.LNumber)

	self := r.IC.Item(int(id)).Self

	var img image.Image
	if self != nil {
		img = self.Image
	} else {
		img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	p := gg.NewSurfacePattern(img, repeatOps[int(repeatOp)])
	return p
}

func patternGradientLinearTable(state *golua.LState, x0, y0, x1, y1 float64) *golua.LTable {
	/// @struct PatternGradientLinear
	/// @prop type {string<context.Pattern>}
	/// @prop x0 {float}
	/// @prop y0 {float}
	/// @prop x1 {float}
	/// @prop y1 {float}
	/// @method color_stop(offset float, struct<image.Color>)

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_GRADIENT_LINEAR))
	t.RawSetString("x0", golua.LNumber(x0))
	t.RawSetString("y0", golua.LNumber(y0))
	t.RawSetString("x1", golua.LNumber(x1))
	t.RawSetString("y1", golua.LNumber(y1))
	t.RawSetString("__colorStops", state.NewTable())

	tableBuilderFunc(state, t, "color_stop", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-2)
		col := state.CheckTable(-1)

		stops := t.RawGetString("__colorStops").(*golua.LTable)
		cs := state.NewTable()
		cs.RawSetString("offset", golua.LNumber(offset))
		cs.RawSetString("color", col)

		stops.Append(cs)
	})

	return t
}

func patternGradientLinearBuild(t *golua.LTable) gg.Pattern {
	x0 := t.RawGetString("x0").(golua.LNumber)
	y0 := t.RawGetString("y0").(golua.LNumber)
	x1 := t.RawGetString("x1").(golua.LNumber)
	y1 := t.RawGetString("y1").(golua.LNumber)

	p := gg.NewLinearGradient(float64(x0), float64(y0), float64(x1), float64(y1))

	colorStops := t.RawGetString("__colorStops").(*golua.LTable)
	for i := range colorStops.Len() {
		cs := colorStops.RawGetInt(i + 1).(*golua.LTable)

		offset := cs.RawGetString("offset").(golua.LNumber)
		ct := cs.RawGetString("color").(*golua.LTable)

		col := imageutil.ColorTableToRGBAColor(ct)
		p.AddColorStop(float64(offset), col)
	}

	return p
}

func patternGradientRadialTable(state *golua.LState, x0, y0, r0, x1, y1, r1 float64) *golua.LTable {
	/// @struct PatternGradientRadial
	/// @prop type {string<context.Pattern>}
	/// @prop x0 {float}
	/// @prop y0 {float}
	/// @prop r0 {float}
	/// @prop x1 {float}
	/// @prop y1 {float}
	/// @prop r1 {float}
	/// @method color_stop(offset float, struct<image.Color>)

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_GRADIENT_RADIAL))
	t.RawSetString("x0", golua.LNumber(x0))
	t.RawSetString("y0", golua.LNumber(y0))
	t.RawSetString("r0", golua.LNumber(r0))
	t.RawSetString("x1", golua.LNumber(x1))
	t.RawSetString("y1", golua.LNumber(y1))
	t.RawSetString("r1", golua.LNumber(r1))
	t.RawSetString("__colorStops", state.NewTable())

	tableBuilderFunc(state, t, "color_stop", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-2)
		col := state.CheckTable(-1)

		stops := t.RawGetString("__colorStops").(*golua.LTable)
		cs := state.NewTable()
		cs.RawSetString("offset", golua.LNumber(offset))
		cs.RawSetString("color", col)

		stops.Append(cs)
	})

	return t
}

func patternGradientRadialBuild(t *golua.LTable) gg.Pattern {
	x0 := t.RawGetString("x0").(golua.LNumber)
	y0 := t.RawGetString("y0").(golua.LNumber)
	r0 := t.RawGetString("r0").(golua.LNumber)
	x1 := t.RawGetString("x1").(golua.LNumber)
	y1 := t.RawGetString("y1").(golua.LNumber)
	r1 := t.RawGetString("r1").(golua.LNumber)

	p := gg.NewRadialGradient(float64(x0), float64(y0), float64(r0), float64(x1), float64(y1), float64(r1))

	colorStops := t.RawGetString("__colorStops").(*golua.LTable)
	for i := range colorStops.Len() {
		cs := colorStops.RawGetInt(i + 1).(*golua.LTable)

		offset := cs.RawGetString("offset").(golua.LNumber)
		ct := cs.RawGetString("color").(*golua.LTable)

		col := imageutil.ColorTableToRGBAColor(ct)
		p.AddColorStop(float64(offset), col)
	}

	return p
}

type PatternCustom struct {
	fn    *golua.LFunction
	state *golua.LState
}

func (p PatternCustom) ColorAt(x, y int) color.Color {
	p.state.Push(p.fn)
	p.state.Push(golua.LNumber(x))
	p.state.Push(golua.LNumber(y))
	p.state.Call(2, 1)
	ct := p.state.CheckTable(-1)
	p.state.Pop(1)

	col := imageutil.ColorTableToRGBAColor(ct)
	return col
}

func patternCustomTable(state *golua.LState, fn *golua.LFunction) *golua.LTable {
	/// @struct PatternCustom
	/// @prop type {string<context.Pattern>}
	/// @prop fn {function(x int, y int) -> struct<image.ColorRGBA>}

	t := state.NewTable()

	t.RawSetString("type", golua.LString(PATTERN_CUSTOM))
	t.RawSetString("fn", fn)

	return t
}

func patternCustomBuild(state *golua.LState, t *golua.LTable) gg.Pattern {
	fn := t.RawGetString("fn").(*golua.LFunction)

	p := PatternCustom{
		fn:    fn,
		state: state,
	}

	return p
}
