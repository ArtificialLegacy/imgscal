package lib

import (
	"fmt"
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/fogleman/gg"
)

const LIB_CONTEXT = "context"

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

			id := r.CC.AddItem(name, &chLog)

			r.CC.Schedule(id, &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					c := gg.NewContext(args["width"].(int), args["height"].(int))
					i.Self = c
					i.Lg.Append("new context created", log.LEVEL_INFO)
				},
			})

			r.State.PushInteger(id)
			return 1
		})

	/// @func to_image()
	/// @arg id
	/// @arg ext - defaults to png
	/// @returns id - new image id
	lib.CreateFunction("to_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "ext", Optional: true},
		},
		func(d lua.TaskData, args map[string]any) int {
			ext := "png"
			if args["ext"] != "" {
				ext = args["ext"].(string)
			}
			name := fmt.Sprintf("image_context_%d.%s", args["id"], ext)

			chLog := log.NewLogger(name)
			chLog.Parent = lg
			lg.Append(fmt.Sprintf("child log created: %s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(name, &chLog)
			contextFinish := make(chan struct{}, 1)
			contextReady := make(chan struct{}, 1)

			var context *gg.Context

			r.CC.Schedule(id, &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					context = i.Self
					contextReady <- struct{}{}
					<-contextFinish
				},
			})

			r.IC.Schedule(id, &collection.Task[image.Image]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[image.Image]) {
					<-contextReady

					img := context.Image()
					i.Self = &img

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

			<-r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					width = i.Self.Width()
					height = i.Self.Height()
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

			<-r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					height = i.Self.FontHeight()
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

			<-r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					width, height = i.Self.MeasureString(args["str"].(string))
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

			<-r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					width, height = i.Self.MeasureMultilineString(args["str"].(string), args["spacing"].(float64))
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

			<-r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					point, e := i.Self.GetCurrentPoint()
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

	/// @func clear()
	/// @arg id
	/// @desc
	/// fills the context with the current color
	lib.CreateFunction("clear",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Clear()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					if args["preserve"].(bool) {
						i.Self.ClipPreserve()
					} else {
						i.Self.Clip()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.ResetClip()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.ClearPath()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.ClosePath()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.MoveTo(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.NewSubPath()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.CubicTo(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.QuadraticTo(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawArc(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawCircle(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawEllipse(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawEllipticalArc(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawLine(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.LineTo(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawPoint(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawRectangle(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawRoundedRectangle(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawRegularPolygon(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawString(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawStringAnchored(
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
	/// @arg  align
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.DrawStringWrapped(
						args["str"].(string),
						args["x"].(float64),
						args["y"].(float64),
						args["ax"].(float64),
						args["ay"].(float64),
						args["width"].(float64),
						args["spacing"].(float64),
						args["align"].(gg.Align),
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					if args["preserve"].(bool) {
						i.Self.FillPreserve()
					} else {
						i.Self.Fill()
					}
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					if args["preserve"].(bool) {
						i.Self.StrokePreserve()
					} else {
						i.Self.Stroke()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Identity()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.InvertMask()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.InvertY()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Push()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Pop()
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Rotate(args["angle"].(float64))
				},
			})

			return 0
		})

	/// @func rotate()
	/// @arg id
	/// @arg angle
	/// @arg x
	/// @arg y
	/// @desc
	/// rotates the transformation matrix around the point
	lib.CreateFunction("rotate",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.RotateAbout(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.Scale(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.ScaleAbout(
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.SetHexColor(args["hex"].(string))
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
			{Type: lua.ARRAY, Name: "pattern", Table: &[]lua.Arg{{Type: lua.FLOAT}}},
		},
		func(d lua.TaskData, args map[string]any) int {
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.SetDash(args["pattern"].([]float64)...)
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
			r.CC.Schedule(args["id"].(int), &collection.Task[gg.Context]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[gg.Context]) {
					i.Self.SetDashOffset(args["offset"].(float64))
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
