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

type Point map[string]float64

func RegisterContext(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CONTEXT, r, r.State, lg)

	/// @func degrees()
	/// @arg radians - float
	/// @returns degrees - float
	lib.CreateFunction(tab, "degrees",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "rad"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			deg := gg.Degrees(args["rad"].(float64))
			state.Push(golua.LNumber(deg))
			return 1
		})

	/// @func radians()
	/// @arg degrees - float
	/// @returns radians - float
	lib.CreateFunction(tab, "radians",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "deg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			rad := gg.Radians(args["deg"].(float64))
			state.Push(golua.LNumber(rad))
			return 1
		})

	/// @func point()
	/// @arg x
	/// @arg y
	/// returns point{x, y}
	lib.CreateFunction(tab, "point",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := state.NewTable()

			state.SetField(t, "x", golua.LNumber(args["x"].(float64)))
			state.SetField(t, "y", golua.LNumber(args["y"].(float64)))

			state.Push(t)
			return 1
		})

	/// @func distance()
	/// @arg p1 - point{x, y}
	/// @arg p2 - point{x, y}
	/// @returns float
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

	/// @func interpolate()
	/// @arg p1 - point{x, y}
	/// @arg p2 - point{x, y}
	/// @arg t - float
	/// @returns point{x, y}
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

	/// @func new()
	/// @arg width - int
	/// @arg height - int
	/// returns id
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

	/// @func new_image()
	/// @arg id - image id to create a context for
	/// @returns new context id
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

	/// @func new_direct()
	/// @arg id - image id to create a context for
	/// @returns new context id
	/// @desc
	/// creates a new context directly on the image,
	/// this requires the image to use the RGBA color model
	/// or it will be converted.
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

	/// @func to_image()
	/// @arg id
	/// @arg name
	/// @arg encoding
	/// @arg? model
	/// @returns id - new image id
	lib.CreateFunction(tab, "to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
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

					img := context.Image()

					i.Self = &collection.ItemImage{
						Image:    img,
						Name:     args["name"].(string),
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Model:    lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib),
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

	/// @func to_mask()
	/// @arg id
	/// @arg name
	/// @arg encoding
	/// @returns id - new alpha image
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

	/// @func mask()
	/// @arg id
	/// @arg img_id
	lib.CreateFunction(tab, "mask",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "img"},
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
					aimg, ok := img.(*image.Alpha)
					if !ok {
						state.Error(golua.LString(lg.Append("invalid image provided to context.mask", log.LEVEL_ERROR)), 0)
					}
					err := i.Self.Context.SetMask(aimg)
					if err != nil {
						state.Error(golua.LString(lg.Append("failed to set image mask, image may be the wrong size.", log.LEVEL_ERROR)), 0)
					}
					imgFinish <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemContext]) {
					imgFinish <- struct{}{}
				},
			})
			return 0
		})

	/// @func size()
	/// @arg id
	/// @returns width
	/// @returns height
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

	/// @func font_height()
	/// @arg id
	/// @returns height
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

	/// @func string_measure()
	/// @arg id
	/// @arg str
	/// @returns width
	/// @returns height
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

	/// @func string_measure_multiline()
	/// @arg id
	/// @arg str
	/// @arg spacing
	/// @returns width
	/// @returns height
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

	/// @func current_point()
	/// @arg id
	/// @returns x
	/// @returns y
	/// @returns exists
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

	/// @arg transform_point()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @returns x
	/// @returns y
	/// @blocking
	/// @desc
	/// multiplies a point by the current matrix.
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

	/// @func clear()
	/// @arg id
	/// @desc
	/// fills the context with the current color
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

	/// @func clip()
	/// @arg id
	/// @arg? preserve - keep the path or not
	/// @desc
	/// updates the clipping region by intersecting the current clipping region with the current path as it would be filled by fill().
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

	/// @func clip_reset()
	/// @arg id
	/// @desc
	/// clears the clipping region
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

	/// @func path_clear()
	/// @arg id
	/// @desc
	/// removes all points from the current path
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

	/// @func path_close()
	/// @arg id
	/// @desc
	/// adds a line segment from the current point to the first point
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

	/// @func path_to()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @desc
	/// starts a new subpath starting at the given point.
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

	/// @func subpath()
	/// @arg id
	/// @desc
	/// starts a new subpath starting at the current point.
	/// no current point will be set after.
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

	/// @func draw_cubic()
	/// @arg id
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @arg x3
	/// @arg y3
	/// @desc
	/// draws a bezier curve to the path starting at the current point
	/// if this isn't a current point, it moves to (x1, y1)
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

	/// @func draw_quadratic()
	/// @arg id
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @desc
	/// draws a quadratic bezier curve to the path starting at the current point
	/// if this isn't a current point, it moves to (x1, y1)
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

	/// @func draw_arc()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg r
	/// @arg angle1
	/// @arg angle2
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

	/// @func draw_circle()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg r
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

	/// @func draw_ellipse()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg rx
	/// @arg ry
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

	/// @func draw_elliptical_arc()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg rx
	/// @arg ry
	/// @arg angle1
	/// @arg angle2
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

	/// @func draw_image()
	/// @arg id
	/// @arg img_id
	/// @arg x
	/// @arg y
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

	/// @func draw_image_anchor()
	/// @arg id
	/// @arg img_id
	/// @arg x
	/// @arg y
	/// @arg ax - float
	/// @arg ay - float
	/// @desc
	/// anchor is between 0 and 1, so 0.5 is centered.
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

	/// @func draw_line()
	/// @arg id
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
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

	/// @func draw_line_to()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @desc
	/// draws a line relative to the current point.
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

	/// @func draw_point()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg r
	/// @desc
	/// similar to draw_circle but ensures that a circle of the specified size is drawn regardless of the current transformation matrix.
	/// the position is still transformed, but not the shape of the point.
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

	/// @func draw_rect()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg width
	/// @arg height
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

	/// @func draw_rect_round()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @arg width
	/// @arg height
	/// @arg r
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

	/// @func draw_polygon()
	/// @arg id
	/// @arg n
	/// @arg x
	/// @arg y
	/// @arg r
	/// @arg rotation
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

	/// @func draw_string()
	/// @arg id
	/// @arg str
	/// @arg x
	/// @arg y
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

	/// @func draw_string_anchor()
	/// @arg id
	/// @arg str
	/// @arg x
	/// @arg y
	/// @arg ax
	/// @arg ay
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

	/// @func draw_string_wrap()
	/// @arg id
	/// @arg str
	/// @arg x
	/// @arg y
	/// @arg ax
	/// @arg ay
	/// @arg width
	/// @arg spacing
	/// @arg align
	lib.CreateFunction(tab, "draw_string",
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

	/// @func fill()
	/// @arg id
	/// @arg? preserve
	/// @desc
	/// fills the current path with the current color.
	/// closes open paths.
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

	/// @func fill_rule()
	/// @arg id
	/// @arg rule
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

	/// @func stroke()
	/// @arg id
	/// @arg? preserve
	/// @desc
	/// strokes the current path with the current color.
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

	/// @func identity()
	/// @arg id
	/// @desc
	/// resets the current transformation matrix to the identity matrix.
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

	/// @func mask_invert()
	/// @arg id
	/// @desc
	/// inverts the alpha values of the clipping mask.
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

	/// @func invert_y()
	/// @arg id
	/// @desc
	/// flips the y axis so that Y=0 is at the bottom of the image.
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

	/// @func push()
	/// @arg id
	/// @desc
	/// push the current context state to the stack.
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

	/// @func pop()
	/// @arg id
	/// @desc
	/// pop the current context state to the stack.
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

	/// @func rotate()
	/// @arg id
	/// @arg angle
	/// @desc
	/// rotates the transformation matrix around the origin
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

	/// @func rotate_about()
	/// @arg id
	/// @arg angle
	/// @arg x
	/// @arg y
	/// @desc
	/// rotates the transformation matrix around the point
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

	/// @func scale()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @desc
	/// scales the transformation matrix by a factor
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

	/// @func scale_about()
	/// @arg id
	/// @arg sx
	/// @arg sy
	/// @arg x
	/// @arg y
	/// @desc
	/// scales the transformation matrix by a factor starting at the point.
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

	/// @func color_hex()
	/// @arg id
	/// @arg hex
	/// @desc
	/// supports hex colors in the follow formats: #RGB #RRGGBB #RRGGBBAA
	/// the leading # is optional
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

	/// @func color()
	/// @arg id
	/// @arg color struct
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

	/// @func color_rgb()
	/// @arg id
	/// @arg r
	/// @arg g
	/// @arg b
	/// @desc
	/// values between 0 and 1.
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

	/// @func color_rgb255()
	/// @arg id
	/// @arg r
	/// @arg g
	/// @arg b
	/// @desc
	/// interger values between 0 and 255.
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

	/// @func color_rgba()
	/// @arg id
	/// @arg r
	/// @arg g
	/// @arg b
	/// @arg a
	/// @desc
	/// values between 0 and 1.
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

	/// @func color_rgba255()
	/// @arg id
	/// @arg r
	/// @arg g
	/// @arg b
	/// @arg a
	/// @desc
	/// interger values between 0 and 255.
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

	/// @func dash_set()
	/// @arg id
	/// @arg pattern - [length...]
	/// @desc
	/// sets the dash length pattern to use.
	/// call with empty array to disable dashes.
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

	/// @func dash_set_offset()
	/// @arg id
	/// @arg offset
	/// @desc
	/// the initial offset for the dash pattern
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

	/// @func line_cap()
	/// @arg id
	/// @arg cap
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

	/// @func line_join()
	/// @arg id
	/// @arg join
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

	/// @func line_width()
	/// @arg id
	/// @arg width
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

	/// @func pixel_set()
	/// @arg id
	/// @arg x
	/// @arg y
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

	/// @func shear()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @desc
	/// updates the current matrix with a shearing angle, at the origin.
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

	/// @func shear_about()
	/// @arg id
	/// @arg sx
	/// @arg sy
	/// @arg x
	/// @arg y
	/// @desc
	/// updates the current matrix with a shearing angle, at the given point.
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

	/// @func translate()
	/// @arg id
	/// @arg x
	/// @arg y
	/// @desc
	/// updates the current matrix with a translation.
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

	/// @func word_wrap()
	/// @arg id
	/// @arg str
	/// @arg width
	/// @returns []string
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

	/// @func matrix_new()
	/// @arg xx
	/// @arg yx
	/// @arg xy
	/// @arg yy
	/// @arg x0
	/// @arg y0
	/// @returns matrix struct
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

	/// @func matrix_identity()
	/// @returns matrix struct
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

	/// @func matrix_rotate()
	/// @arg angle
	/// @returns matrix struct
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

	/// @func matrix_scale()
	/// @arg x
	/// @arg y
	/// @returns matrix struct
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

	/// @func matrix_shear()
	/// @arg x
	/// @arg y
	/// @returns matrix struct
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

	/// @func matrix_translate()
	/// @arg x
	/// @arg y
	/// @returns matrix struct
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

	/// @func point_cubic_bezier()
	/// @arg x0
	/// @arg y0
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @arg x3
	/// @arg y3
	/// @returns []point
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

	/// @func point_quadratic_bezier()
	/// @arg x0
	/// @arg y0
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	/// @returns []point
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

	/// @func pattern_solid()
	/// @arg color struct
	/// @returns pattern struct
	lib.CreateFunction(tab, "pattern_solid",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternSolidTable(state, args["color"].(*golua.LTable))

			state.Push(t)
			return 1
		})

	/// @func pattern_surface()
	/// @arg id
	/// @arg repeat_op
	/// @returns pattern struct
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

	/// @func pattern_surface_sync()
	/// @arg id
	/// @arg repeat_op
	/// @returns pattern struct
	/// @desc
	/// Note: this does not wait for the image to be ready or idle,
	/// if the image is not loaded it will grab an empy image
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

	/// @func pattern_custom()
	/// @arg fn (x, y) color struct
	/// @returns pattern struct
	lib.CreateFunction(tab, "pattern_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := patternCustomTable(state, args["fn"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func gradient_linear()
	/// @arg x0
	/// @arg y0
	/// @arg x1
	/// @arg y1
	/// @returns pattern struct
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

	/// @func gradient_radial()
	/// @arg x0
	/// @arg y0
	/// @arg r0
	/// @arg x1
	/// @arg y1
	/// @arg r1
	/// @returns pattern struct
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

	/// @func fill_style()
	/// @arg id
	/// @arg pattern struct
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

	/// @func stroke_style()
	/// @arg id
	/// @arg pattern struct
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

	/// @func font_load()
	/// @arg path
	/// @arg points
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
	/// @const LINCAP_SQUARE
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
	/// @prop xx
	/// @prop yx
	/// @prop xy
	/// @prop yy
	/// @prop x0
	/// @prop y0
	/// @method multiply(matrix)
	/// @method rotate(angle)
	/// @method scale(x, y)
	/// @method shear(x, y)
	/// @method translate(x, y)
	/// @method transform_point(x, y) x, y
	/// @method transform_vector(x, y) x, y

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
	/// @prop type

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
	/// @prop type
	/// @prop color

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
	/// @prop type
	/// @prop id
	/// @prop repeatOp

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
	/// @prop type
	/// @prop id
	/// @prop repeatOp

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
	/// @prop type
	/// @prop x0
	/// @prop y0
	/// @prop x1
	/// @prop y1
	/// @method color_stop(offset, color)

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
	/// @prop type
	/// @prop x0
	/// @prop y0
	/// @prop r0
	/// @prop x1
	/// @prop y1
	/// @prop r1
	/// @method color_stop(offset, color)

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
	/// @prop type
	/// @prop fn(x, y) color struct

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
