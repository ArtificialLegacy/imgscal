package lib

import (
	"fmt"
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/fogleman/gg"
)

const LIB_CONTEXT = "context"

type Point map[string]float64

func RegisterContext(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_CONTEXT, r.State, lg)

	/// @func degrees()
	/// @arg radians - float
	/// @returns degrees - float
	lib.CreateFunction("degrees",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "rad"},
		},
		func(d lua.TaskData, args map[string]any) int {
			deg := gg.Degrees(args["rad"].(float64))
			r.State.PushNumber(deg)
			return 1
		})

	/// @func radians()
	/// @arg degrees - float
	/// @returns radians - float
	lib.CreateFunction("radians",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "deg"},
		},
		func(d lua.TaskData, args map[string]any) int {
			rad := gg.Radians(args["deg"].(float64))
			r.State.PushNumber(rad)
			return 1
		})

	/// @func point()
	/// @arg x
	/// @arg y
	/// returns point{x, y}
	lib.CreateFunction("point",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.State.NewTable()

			r.State.PushNumber(args["x"].(float64))
			r.State.SetField(-2, "x")
			r.State.PushNumber(args["y"].(float64))
			r.State.SetField(-2, "y")

			return 1
		})

	/// @func distance()
	/// @arg p1 - point{x, y}
	/// @arg p2 - point{x, y}
	/// @returns float
	lib.CreateFunction("distance",
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
		func(d lua.TaskData, args map[string]any) int {
			ap1 := args["p1"].(Point)
			ap2 := args["b2"].(Point)

			p1 := gg.Point{X: ap1["x"], Y: ap1["y"]}
			p2 := gg.Point{X: ap2["x"], Y: ap2["y"]}

			dist := p1.Distance(p2)

			r.State.PushNumber(dist)
			return 1
		})

	/// @func interpolate()
	/// @arg p1 - point{x, y}
	/// @arg p2 - point{x, y}
	/// @arg t - float
	/// @returns point{x, y}
	lib.CreateFunction("interpolate",
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
		func(d lua.TaskData, args map[string]any) int {
			ap1 := args["p1"].(Point)
			ap2 := args["b2"].(Point)

			p1 := gg.Point{X: ap1["x"], Y: ap1["y"]}
			p2 := gg.Point{X: ap2["x"], Y: ap2["y"]}

			pi := p1.Interpolate(p2, args["t"].(float64))

			r.State.NewTable()

			r.State.PushNumber(pi.X)
			r.State.SetField(-2, "x")
			r.State.PushNumber(pi.Y)
			r.State.SetField(-2, "y")
			return 1
		})

	/// @func new()
	/// @arg width - int
	/// @arg height - int
	/// returns id
	lib.CreateFunction("new",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(d lua.TaskData, args map[string]any) int {
			name := fmt.Sprintf("context_%d", r.CC.Next())

			chLog := log.NewLogger(name)
			chLog.Parent = lg
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

			r.State.PushInteger(id)
			return 1
		})

	/// @func new_image()
	/// @arg id - image id to create a context for
	/// @returns new context id
	lib.CreateFunction("new_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

			chLog := log.NewLogger(tempName)
			chLog.Parent = lg
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

			r.State.PushInteger(id)
			return 1
		})

	/// @func to_image()
	/// @arg id
	/// @arg name
	/// @arg encoding
	/// @arg model
	/// @returns id - new image id
	lib.CreateFunction("to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "model", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
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

			chLog := log.NewLogger(args["name"].(string))
			chLog.Parent = lg
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

			r.State.PushInteger(id)
			return 1
		})

	/// @func size()
	/// @arg id
	/// @returns width
	/// @returns height
	/// @blocking
	lib.CreateFunction("size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

			r.State.PushInteger(width)
			r.State.PushInteger(height)
			return 2
		})

	/// @func font_height()
	/// @arg id
	/// @returns height
	/// @blocking
	lib.CreateFunction("font_height",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					height = i.Self.Context.FontHeight()
				},
			})

			r.State.PushNumber(height)
			return 1
		})

	/// @func string_measure()
	/// @arg id
	/// @arg str
	/// @returns width
	/// @returns height
	/// @blocking
	lib.CreateFunction("string_measure",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
		},
		func(d lua.TaskData, args map[string]any) int {
			width := 0.0
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					width, height = i.Self.Context.MeasureString(args["str"].(string))
				},
			})

			r.State.PushNumber(width)
			r.State.PushNumber(height)
			return 2
		})

	/// @func string_measure_multiline()
	/// @arg id
	/// @arg str
	/// @arg spacing
	/// @returns width
	/// @returns height
	/// @blocking
	lib.CreateFunction("string_measure_multiline",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "spacing"},
		},
		func(d lua.TaskData, args map[string]any) int {
			width := 0.0
			height := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					width, height = i.Self.Context.MeasureMultilineString(args["str"].(string), args["spacing"].(float64))
				},
			})

			r.State.PushNumber(width)
			r.State.PushNumber(height)
			return 2
		})

	/// @func current_point()
	/// @arg id
	/// @returns x
	/// @returns y
	/// @returns exists
	/// @blocking
	lib.CreateFunction("current_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

			r.State.PushNumber(x)
			r.State.PushNumber(y)
			r.State.PushBoolean(exists)
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
	lib.CreateFunction("transform_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
			x := 0.0
			y := 0.0

			<-r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					x, y = i.Self.Context.TransformPoint(args["x"].(float64), args["y"].(float64))
				},
			})

			r.State.PushNumber(x)
			r.State.PushNumber(y)
			return 2
		})

	/// @func clear()
	/// @arg id
	/// @desc
	/// fills the context with the current color
	lib.CreateFunction("clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	/// @arg preserve - keep the path or not
	/// @desc
	/// updates the clipping region by intersecting the current clipping region with the current path as it would be filled by fill().
	lib.CreateFunction("clip",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("clip_reset",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("path_clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("path_close",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("path_to",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("subpath",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_cubic",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
			{Type: lua.FLOAT, Name: "x3"},
			{Type: lua.FLOAT, Name: "y3"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_quadratic",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_arc",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "angle1"},
			{Type: lua.FLOAT, Name: "angle2"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_circle",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_ellipse",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "rx"},
			{Type: lua.FLOAT, Name: "ry"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_elliptical_arc",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "rx"},
			{Type: lua.FLOAT, Name: "ry"},
			{Type: lua.FLOAT, Name: "angle1"},
			{Type: lua.FLOAT, Name: "angle2"},
		},
		func(d lua.TaskData, args map[string]any) int {
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

	/// @func draw_line()
	/// @arg id
	/// @arg x1
	/// @arg y1
	/// @arg x2
	/// @arg y2
	lib.CreateFunction("draw_line",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x1"},
			{Type: lua.FLOAT, Name: "y1"},
			{Type: lua.FLOAT, Name: "x2"},
			{Type: lua.FLOAT, Name: "y2"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_line_to",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_point",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_rect",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_rect_round",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
			{Type: lua.FLOAT, Name: "r"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_polygon",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "n"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "rotation"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_string_anchor",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "ax"},
			{Type: lua.FLOAT, Name: "ay"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("draw_string",
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
		func(d lua.TaskData, args map[string]any) int {
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
	/// @arg preserve
	/// @desc
	/// fills the current path with the current color.
	/// closes open paths.
	lib.CreateFunction("fill",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("fill_rule",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "rule"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	/// @arg preserve
	/// @desc
	/// strokes the current path with the current color.
	lib.CreateFunction("stroke",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "preserve", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("identity",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("mask_invert",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("invert_y",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("push",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("pop",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("rotate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "angle"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("rotate_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("scale",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("scale_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "sx"},
			{Type: lua.FLOAT, Name: "sy"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("color_hex",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "hex"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.SetHexColor(args["hex"].(string))
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
	lib.CreateFunction("color_rgb",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "g"},
			{Type: lua.FLOAT, Name: "b"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("color_rgb255",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("color_rgba",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "r"},
			{Type: lua.FLOAT, Name: "g"},
			{Type: lua.FLOAT, Name: "b"},
			{Type: lua.FLOAT, Name: "a"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("color_rgba255",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
			{Type: lua.INT, Name: "a"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("dash_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("pattern", lua.ArrayType{Type: lua.FLOAT}, false),
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("dash_set_offset",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "offset"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("line_cap",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "cap"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("line_join",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "join"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("line_width",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "width"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("pixel_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("shear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("shear_about",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "sx"},
			{Type: lua.FLOAT, Name: "sy"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
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
	lib.CreateFunction("translate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[collection.ItemContext]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemContext]) {
					i.Self.Context.Translate(args["x"].(float64), args["y"].(float64))
				},
			})

			return 0
		})

	/// @constants Fill Rules
	/// @const FILLRULE_WINDING
	/// @const FILLRULE_EVENODD
	lib.State.PushInteger(int(gg.FillRuleWinding))
	lib.State.SetField(-2, "FILLRULE_WINDING")
	lib.State.PushInteger(int(gg.FillRuleEvenOdd))
	lib.State.SetField(-2, "FILLRULE_EVENODD")

	/// @constants Line Caps
	/// @const LINECAP_ROUND
	/// @const LINECAP_BUTT
	/// @const LINCAP_SQUARE
	lib.State.PushInteger(int(gg.LineCapRound))
	lib.State.SetField(-2, "LINECAP_ROUND")
	lib.State.PushInteger(int(gg.LineCapButt))
	lib.State.SetField(-2, "LINECAP_BUTT")
	lib.State.PushInteger(int(gg.LineCapSquare))
	lib.State.SetField(-2, "LINECAP_SQUARE")

	/// @constants Line Joins
	/// @const LINEJOIN_ROUND
	/// @const LINEJOIN_BEVEL
	lib.State.PushInteger(int(gg.LineJoinRound))
	lib.State.SetField(-2, "LINEJOIN_ROUND")
	lib.State.PushInteger(int(gg.LineJoinBevel))
	lib.State.SetField(-2, "LINEJOIN_BEVEL")

	/// @constants Repeat Ops
	/// @const REPEAT_BOTH
	/// @const REPEAT_X
	/// @const REPEAT_Y
	/// @const REPEAT_NONE
	lib.State.PushInteger(int(gg.RepeatBoth))
	lib.State.SetField(-2, "REPEAT_BOTH")
	lib.State.PushInteger(int(gg.RepeatX))
	lib.State.SetField(-2, "REPEAT_X")
	lib.State.PushInteger(int(gg.RepeatY))
	lib.State.SetField(-2, "REPEAT_Y")
	lib.State.PushInteger(int(gg.RepeatNone))
	lib.State.SetField(-2, "REPEAT_NONE")

	/// @constants Alignment
	/// @const ALIGN_LEFT
	/// @const ALIGN_CENTER
	/// @const ALIGN_RIGHT
	lib.State.PushInteger(int(gg.AlignLeft))
	lib.State.SetField(-2, "ALIGN_LEFT")
	lib.State.PushInteger(int(gg.AlignCenter))
	lib.State.SetField(-2, "ALIGN_CENTER")
	lib.State.PushInteger(int(gg.AlignRight))
	lib.State.SetField(-2, "ALIGN_RIGHT")
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
