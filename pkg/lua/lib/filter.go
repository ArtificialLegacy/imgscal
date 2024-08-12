package lib

import (
	"image"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/disintegration/gift"
	golua "github.com/yuin/gopher-lua"
)

const LIB_FILTER = "filter"

func RegisterFilter(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_FILTER, r, r.State, lg)

	/// @func draw()
	/// @arg id1
	/// @arg id2
	/// @arg []filter
	/// desc
	/// applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.ANY, Name: "filters"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["id1"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.IC.Schedule(args["id2"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady

					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					g.Draw(imageutil.ImageGetDraw(i.Self.Image), img)

					state.Close()
					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					state.Close()
					imgFinished <- struct{}{}
				},
			})

			return -1
		})

	/// @func draw_at()
	/// @arg id1
	/// @arg id2
	/// @arg point
	/// @arg op
	/// @arg []filter
	/// desc
	/// applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw_at",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.ANY, Name: "point"},
			{Type: lua.INT, Name: "op"},
			{Type: lua.ANY, Name: "filters"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["id1"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.IC.Schedule(args["id2"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady

					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					pt := imageutil.TableToPoint(state, args["point"].(*golua.LTable))
					g.DrawAt(imageutil.ImageGetDraw(i.Self.Image), img, pt, gift.Operator(args["op"].(int)))

					state.Close()
					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					state.Close()
					imgFinished <- struct{}{}
				},
			})

			return -1
		})

	/// @func draw_at_xy()
	/// @arg id1
	/// @arg id2
	/// @arg x
	/// @arg y
	/// @arg op
	/// @arg []filter
	/// desc
	/// applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw_at_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.INT, Name: "op"},
			{Type: lua.ANY, Name: "filters"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			var img image.Image

			r.IC.Schedule(args["id1"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			r.IC.Schedule(args["id2"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady

					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					g.DrawAt(
						imageutil.ImageGetDraw(i.Self.Image), img,
						image.Point{X: args["x"].(int), Y: args["y"].(int)},
						gift.Operator(args["op"].(int)),
					)

					state.Close()
					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					state.Close()
					imgFinished <- struct{}{}
				},
			})

			return -1
		})

	/// @func bounds()
	/// @arg id
	/// @arg []filter
	/// @returns x1
	/// @returns y1
	/// @returns x2
	/// @returns y2
	/// @blocking
	/// desc
	/// Gets the resulting bounds of the image after the filters are applied.
	lib.CreateFunction(tab, "bounds",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "filters"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var dstBounds image.Rectangle
			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					dstBounds = g.Bounds(i.Self.Image.Bounds())
				},
			})

			state.Push(golua.LNumber(dstBounds.Min.X))
			state.Push(golua.LNumber(dstBounds.Min.Y))
			state.Push(golua.LNumber(dstBounds.Max.X))
			state.Push(golua.LNumber(dstBounds.Max.Y))
			return 4
		})

	/// @func bounds_size()
	/// @arg id
	/// @arg []filter
	/// @returns width
	/// @returns height
	/// @blocking
	/// desc
	/// Gets the resulting size of the image after the filters are applied.
	lib.CreateFunction(tab, "bounds_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "filters"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var dstBounds image.Rectangle
			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					dstBounds = g.Bounds(i.Self.Image.Bounds())
				},
			})

			state.Push(golua.LNumber(dstBounds.Dx()))
			state.Push(golua.LNumber(dstBounds.Dy()))
			return 2
		})

	/// @func brightness()
	/// @arg percent - between -100 and 100.
	/// @returns filter
	lib.CreateFunction(tab, "brightness",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := brightnessTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func color_balance()
	/// @arg percentRed - between -100 and 500.
	/// @arg percentGreen - between -100 and 500.
	/// @arg percentBlue - between -100 and 500.
	/// @returns filter
	lib.CreateFunction(tab, "color_balance",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percentRed"},
			{Type: lua.FLOAT, Name: "percentGreen"},
			{Type: lua.FLOAT, Name: "percentBlue"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorBalanceTable(state, args["percentRed"].(float64), args["percentGreen"].(float64), args["percentBlue"].(float64))

			state.Push(t)
			return 1
		})

	/// @func colorize()
	/// @arg hue - between 0 and 360.
	/// @arg saturation - between 0 and 100.
	/// @arg percent - between 0 and 100.
	/// @returns filter
	lib.CreateFunction(tab, "colorize",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "hue"},
			{Type: lua.FLOAT, Name: "saturation"},
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorizeTable(state, args["hue"].(float64), args["saturation"].(float64), args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func colorspace_linear_to_srgb()
	/// @returns filter
	lib.CreateFunction(tab, "colorspace_linear_to_srgb",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorspaceLinearSRGBTable(state)

			state.Push(t)
			return 1
		})

	/// @func colorspace_srgb_to_linear()
	/// @returns filter
	lib.CreateFunction(tab, "colorspace_srgb_to_linear",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorspaceSRGBLinearTable(state)

			state.Push(t)
			return 1
		})

	/// @func contrast()
	/// @arg percent - between -100 and 100.
	/// @returns filter
	lib.CreateFunction(tab, "contrast",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := contrastTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func convolution()
	/// @arg kernel - must be len of an odd square, eg 3x3=9 or 5x5=25
	/// @arg normalize
	/// @arg alpha
	/// @arg abs
	/// @arg delta
	/// @returns filter
	lib.CreateFunction(tab, "convolution",
		[]lua.Arg{
			{Type: lua.ANY, Name: "kernel"},
			{Type: lua.BOOL, Name: "normalize"},
			{Type: lua.BOOL, Name: "alpha"},
			{Type: lua.BOOL, Name: "abs"},
			{Type: lua.FLOAT, Name: "delta"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := convolutionTable(state,
				args["kernel"].(golua.LValue),
				args["normalize"].(bool),
				args["alpha"].(bool),
				args["abs"].(bool),
				args["delta"].(float64),
			)

			state.Push(t)
			return 1
		})

	/// @func crop()
	/// @arg min
	/// @arg max
	/// @returns filter
	lib.CreateFunction(tab, "crop",
		[]lua.Arg{
			{Type: lua.ANY, Name: "min"},
			{Type: lua.ANY, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			min := imageutil.TableToPoint(state, args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(state, args["max"].(*golua.LTable))

			t := cropTable(state, min.X, min.Y, max.X, max.Y)

			state.Push(t)
			return 1
		})

	/// @func crop_xy()
	/// @arg xmin
	/// @arg ymin
	/// @arg xmax
	/// @arg ymax
	/// @returns filter
	lib.CreateFunction(tab, "crop_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "xmin"},
			{Type: lua.INT, Name: "ymin"},
			{Type: lua.INT, Name: "xmax"},
			{Type: lua.INT, Name: "ymax"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cropTable(state,
				args["xmin"].(int),
				args["ymin"].(int),
				args["xmax"].(int),
				args["ymax"].(int),
			)

			state.Push(t)
			return 1
		})

	/// @func crop_to_size()
	/// @arg width
	/// @arg height
	/// @arg anchor
	/// @returns filter
	lib.CreateFunction(tab, "crop_to_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "anchor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cropToSizeTable(state,
				args["width"].(int),
				args["height"].(int),
				args["anchor"].(int),
			)

			state.Push(t)
			return 1
		})

	/// @func flip_horizontal()
	/// @returns filter
	lib.CreateFunction(tab, "flip_horizontal",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := flipHorizontalTable(state)

			state.Push(t)
			return 1
		})

	/// @func flip_vertical()
	/// @returns filter
	lib.CreateFunction(tab, "flip_vertical",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := flipVerticalTable(state)

			state.Push(t)
			return 1
		})

	/// @func gamma()
	/// @arg gamma
	/// @returns filter
	lib.CreateFunction(tab, "gamma",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "gamma"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := gammaTable(state, args["gamma"].(float64))

			state.Push(t)
			return 1
		})

	/// @func gaussian_blur()
	/// @arg sigma
	/// @returns filter
	lib.CreateFunction(tab, "gaussian_blur",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "sigma"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := gaussianBlurTable(state, args["sigma"].(float64))

			state.Push(t)
			return 1
		})

	/// @func grayscale()
	/// @returns filter
	lib.CreateFunction(tab, "grayscale",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := grayscaleTable(state)

			state.Push(t)
			return 1
		})

	/// @func invert()
	/// @returns filter
	lib.CreateFunction(tab, "invert",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := invertTable(state)

			state.Push(t)
			return 1
		})

	/// @func rotate()
	/// @arg angle
	/// @arg bgcolor
	/// @arg interpolation
	/// @returns filter
	lib.CreateFunction(tab, "rotate",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.ANY, Name: "bgcolor"},
			{Type: lua.INT, Name: "interpolation"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotateTable(state,
				args["angle"].(float64),
				args["bgcolor"].(golua.LValue),
				args["interpolation"].(int),
			)

			state.Push(t)
			return 1
		})

	/// @func rotate_180()
	/// @returns filter
	lib.CreateFunction(tab, "rotate_180",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate180Table(state)

			state.Push(t)
			return 1
		})

	/// @func rotate_270()
	/// @returns filter
	lib.CreateFunction(tab, "rotate_270",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate270Table(state)

			state.Push(t)
			return 1
		})

	/// @func rotate_90()
	/// @returns filter
	lib.CreateFunction(tab, "rotate_90",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate90Table(state)

			state.Push(t)
			return 1
		})

	/// @func hue()
	/// @arg shift -180 to 180
	/// @returns filter
	lib.CreateFunction(tab, "hue",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "shift"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := hueTable(state, args["shift"].(float64))

			state.Push(t)
			return 1
		})

	/// @func saturation()
	/// @arg percent -100 to 500
	/// @returns filter
	lib.CreateFunction(tab, "saturation",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := saturationTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func sepia()
	/// @arg percent 0 to 100
	/// @returns filter
	lib.CreateFunction(tab, "sepia",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := sepiaTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func pixelate()
	/// @arg size
	/// @returns filter
	lib.CreateFunction(tab, "pixelate",
		[]lua.Arg{
			{Type: lua.INT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := pixelateTable(state, args["size"].(int))

			state.Push(t)
			return 1
		})

	/// @func threshold()
	/// @arg percent 0 to 100
	/// @returns filter
	lib.CreateFunction(tab, "threshold",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := thresholdTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func transpose()
	/// @returns filter
	lib.CreateFunction(tab, "transpose",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := transposeTable(state)

			state.Push(t)
			return 1
		})

	/// @func transverse()
	/// @returns filter
	lib.CreateFunction(tab, "transverse",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := transverseTable(state)

			state.Push(t)
			return 1
		})

	/// @func sobel()
	/// @returns filter
	lib.CreateFunction(tab, "sobel",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := sobelTable(state)

			state.Push(t)
			return 1
		})

	/// @func maximum()
	/// @arg ksize - must be odd int, e.g. 3, 5, 7
	/// @arg disk
	/// @returns filter
	lib.CreateFunction(tab, "maximum",
		[]lua.Arg{
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.BOOL, Name: "disk"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := maximumTable(state, args["ksize"].(int), args["disk"].(bool))

			state.Push(t)
			return 1
		})

	/// @func mean()
	/// @arg ksize - must be odd int, e.g. 3, 5, 7
	/// @arg disk
	/// @returns filter
	lib.CreateFunction(tab, "mean",
		[]lua.Arg{
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.BOOL, Name: "disk"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := meanTable(state, args["ksize"].(int), args["disk"].(bool))

			state.Push(t)
			return 1
		})

	/// @func median()
	/// @arg ksize - must be odd int, e.g. 3, 5, 7
	/// @arg disk
	/// @returns filter
	lib.CreateFunction(tab, "median",
		[]lua.Arg{
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.BOOL, Name: "disk"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := medianTable(state, args["ksize"].(int), args["disk"].(bool))

			state.Push(t)
			return 1
		})

	/// @func minimum()
	/// @arg ksize - must be odd int, e.g. 3, 5, 7
	/// @arg disk
	/// @returns filter
	lib.CreateFunction(tab, "minimum",
		[]lua.Arg{
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.BOOL, Name: "disk"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := minimumTable(state, args["ksize"].(int), args["disk"].(bool))

			state.Push(t)
			return 1
		})

	/// @func sigmoid()
	/// @arg midpoint 0 to 1
	/// @arg factor
	/// @returns filter
	lib.CreateFunction(tab, "sigmoid",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "midpoint"},
			{Type: lua.FLOAT, Name: "factor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := sigmoidTable(state, args["midpoint"].(float64), args["factor"].(float64))

			state.Push(t)
			return 1
		})

	/// @func unsharp_mask()
	/// @arg sigma
	/// @arg amount
	/// @arg threshold
	/// @returns filter
	lib.CreateFunction(tab, "unsharp_mask",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "sigma"},
			{Type: lua.FLOAT, Name: "amount"},
			{Type: lua.FLOAT, Name: "threshold"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := unsharpMaskTable(state, args["sigma"].(float64), args["amount"].(float64), args["threshold"].(float64))

			state.Push(t)
			return 1
		})

	/// @func resize()
	/// @arg width
	/// @arg height
	/// @arg resampling
	/// @returns filter
	lib.CreateFunction(tab, "resize",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "resampling"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := resizeTable(state, args["width"].(int), args["height"].(int), args["resampling"].(int))

			state.Push(t)
			return 1
		})

	/// @func resize_to_fill()
	/// @arg width
	/// @arg height
	/// @arg resampling
	/// @arg anchor
	/// @returns filter
	lib.CreateFunction(tab, "resize_to_fill",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "resampling"},
			{Type: lua.INT, Name: "anchor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := resizeToFillTable(state, args["width"].(int), args["height"].(int), args["resampling"].(int), args["anchor"].(int))

			state.Push(t)
			return 1
		})

	/// @func resize_to_fit()
	/// @arg width
	/// @arg height
	/// @arg resampling
	/// @returns filter
	lib.CreateFunction(tab, "resize_to_fit",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "resampling"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := resizeToFitTable(state, args["width"].(int), args["height"].(int), args["resampling"].(int))

			state.Push(t)
			return 1
		})

	/// @func color_func()
	/// @arg fn - function(r,g,b,a) r,g,b,a
	/// @returns filter
	/// @desc
	/// Color values are floats between 0 and 1.
	lib.CreateFunction(tab, "color_func",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorFuncTable(state, args["fn"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @constants Anchor
	/// @const ANCHOR_CENTER
	/// @const ANCHOR_TOPLEFT
	/// @const ANCHOR_TOP
	/// @const ANCHOR_TOPRIGHT
	/// @const ANCHOR_LEFT
	/// @const ANCHOR_RIGHT
	/// @const ANCHOR_BOTTOMLEFT
	/// @const ANCHOR_BOTTOM
	/// @const ANCHOR_BOTTOMRIGHT
	r.State.SetTable(tab, golua.LString("ANCHOR_CENTER"), golua.LNumber(gift.CenterAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_TOPLEFT"), golua.LNumber(gift.TopLeftAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_TOP"), golua.LNumber(gift.TopAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_TOPRIGHT"), golua.LNumber(gift.TopRightAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_LEFT"), golua.LNumber(gift.LeftAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_RIGHT"), golua.LNumber(gift.RightAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_BOTTOMLEFT"), golua.LNumber(gift.BottomLeftAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_BOTTOM"), golua.LNumber(gift.BottomAnchor))
	r.State.SetTable(tab, golua.LString("ANCHOR_BOTTOMRIGHT"), golua.LNumber(gift.BottomRightAnchor))

	/// @constants Interpolation
	/// @const INTERPOLATION_NEARESTNEIGHBOR
	/// @const INTERPOLATION_LINEAR
	/// @const INTERPOLATION_CUBIC
	r.State.SetTable(tab, golua.LString("INTERPOLATION_NEARESTNEIGHBOR"), golua.LNumber(gift.NearestNeighborInterpolation))
	r.State.SetTable(tab, golua.LString("INTERPOLATION_LINEAR"), golua.LNumber(gift.LinearInterpolation))
	r.State.SetTable(tab, golua.LString("INTERPOLATION_CUBIC"), golua.LNumber(gift.CubicInterpolation))

	/// @constants Operators
	/// @const OPERATOR_COPY
	/// @const OPERATOR_OVER
	r.State.SetTable(tab, golua.LString("OPERATOR_COPY"), golua.LNumber(gift.CopyOperator))
	r.State.SetTable(tab, golua.LString("OPERATOR_OVER"), golua.LNumber(gift.OverOperator))

	/// @constants Resampling
	/// @const RESAMPLING_BOX
	/// @const RESAMPLING_CUBIC
	/// @const RESAMPLING_LANCZOS
	/// @const RESAMPLING_LINEAR
	/// @const RESAMPLING_NEARESTNEIGHBOR
	r.State.SetTable(tab, golua.LString("RESAMPLING_BOX"), golua.LNumber(RESAMPLING_BOX))
	r.State.SetTable(tab, golua.LString("RESAMPLING_CUBIC"), golua.LNumber(RESAMPLING_CUBIC))
	r.State.SetTable(tab, golua.LString("RESAMPLING_LANCZOS"), golua.LNumber(RESAMPLING_LANCZOS))
	r.State.SetTable(tab, golua.LString("RESAMPLING_LINEAR"), golua.LNumber(RESAMPLING_LINEAR))
	r.State.SetTable(tab, golua.LString("RESAMPLING_NEARESTNEIGHBOR"), golua.LNumber(RESAMPLING_NEARESTNEIGHBOR))
}

var samplers = []gift.Resampling{
	gift.BoxResampling,
	gift.CubicResampling,
	gift.LanczosResampling,
	gift.LinearResampling,
	gift.NearestNeighborResampling,
}

const (
	RESAMPLING_BOX int = iota
	RESAMPLING_CUBIC
	RESAMPLING_LANCZOS
	RESAMPLING_LINEAR
	RESAMPLING_NEARESTNEIGHBOR
)

const (
	FILTER_BRIGHTNESS                = "brightness"
	FILTER_COLOR_BALANCE             = "color_balance"
	FILTER_COLORIZE                  = "colorize"
	FILTER_COLORSPACE_LINEAR_TO_SRGB = "colorspace_linear_to_srgb"
	FILTER_COLORSPACE_SRGB_TO_LINEAR = "colorspace_srgb_to_linear"
	FILTER_CONTRAST                  = "contrast"
	FILTER_CONVOLUTION               = "convolution"
	FILTER_CROP                      = "crop"
	FILTER_CROP_TO_SIZE              = "crop_to_size"
	FILTER_FLIP_HORIZONTAL           = "flip_horizontal"
	FILTER_FLIP_VERTICAL             = "flip_vertical"
	FILTER_GAMMA                     = "gamma"
	FILTER_GAUSSIAN_BLUR             = "gaussian_blur"
	FILTER_GRAYSCALE                 = "grayscale"
	FILTER_INVERT                    = "invert"
	FILTER_ROTATE                    = "rotate"
	FILTER_ROTATE180                 = "rotate_180"
	FILTER_ROTATE270                 = "rotate_270"
	FILTER_ROTATE90                  = "rotate_90"
	FILTER_HUE                       = "hue"
	FILTER_SATURATION                = "saturation"
	FILTER_SEPIA                     = "sepia"
	FILTER_THRESHOLD                 = "threshold"
	FILTER_PIXELATE                  = "pixelate"
	FILTER_SOBEL                     = "sobel"
	FILTER_TRANSPOSE                 = "transpose"
	FILTER_TRANSVERSE                = "transverse"
	FILTER_MAXIMUM                   = "maximum"
	FILTER_MEAN                      = "mean"
	FILTER_MEDIAN                    = "median"
	FILTER_MINIMUM                   = "minimum"
	FILTER_SIGMOID                   = "sigmoid"
	FILTER_UNSHARP_MASK              = "unsharp_mask"
	FILTER_RESIZE                    = "resize"
	FILTER_RESIZE_TO_FILL            = "resize_to_fill"
	FILTER_RESIZE_TO_FIT             = "resize_to_fit"
	FILTER_COLOR_FUNC                = "color_func"
)

type filterList map[string]func(state *golua.LState, t *golua.LTable) gift.Filter

var filters = filterList{
	FILTER_BRIGHTNESS:                brightnessBuild,
	FILTER_COLOR_BALANCE:             colorBalanceBuild,
	FILTER_COLORIZE:                  colorizeBuild,
	FILTER_COLORSPACE_LINEAR_TO_SRGB: colorspaceLinearSRGBBuild,
	FILTER_COLORSPACE_SRGB_TO_LINEAR: colorspaceSRGBLinearBuild,
	FILTER_CONTRAST:                  contrastBuild,
	FILTER_CONVOLUTION:               convolutionBuild,
	FILTER_CROP:                      cropBuild,
	FILTER_CROP_TO_SIZE:              cropToSizeBuild,
	FILTER_FLIP_HORIZONTAL:           flipHorizontalBuild,
	FILTER_FLIP_VERTICAL:             flipVerticalBuild,
	FILTER_GAMMA:                     gammaBuild,
	FILTER_GAUSSIAN_BLUR:             gaussianBlurBuild,
	FILTER_GRAYSCALE:                 grayscaleBuild,
	FILTER_INVERT:                    invertBuild,
	FILTER_ROTATE180:                 rotate180Build,
	FILTER_ROTATE270:                 rotate270Build,
	FILTER_ROTATE90:                  rotate90Build,
	FILTER_ROTATE:                    rotateBuild,
	FILTER_HUE:                       hueBuild,
	FILTER_SATURATION:                saturationBuild,
	FILTER_SEPIA:                     sepiaBuild,
	FILTER_THRESHOLD:                 thresholdBuild,
	FILTER_PIXELATE:                  pixelateBuild,
	FILTER_SOBEL:                     sobelBuild,
	FILTER_TRANSPOSE:                 transposeBuild,
	FILTER_TRANSVERSE:                transverseBuild,
	FILTER_MAXIMUM:                   maximumBuild,
	FILTER_MEAN:                      meanBuild,
	FILTER_MEDIAN:                    medianBuild,
	FILTER_MINIMUM:                   minimumBuild,
	FILTER_SIGMOID:                   sigmoidBuild,
	FILTER_UNSHARP_MASK:              unsharpMaskBuild,
	FILTER_RESIZE:                    resizeBuild,
	FILTER_RESIZE_TO_FILL:            resizeToFillBuild,
	FILTER_RESIZE_TO_FIT:             resizeToFitBuild,
	FILTER_COLOR_FUNC:                colorFuncBuild,
}

func buildFilterList(state *golua.LState, filterList filterList, t *golua.LTable) *gift.GIFT {
	filters := []gift.Filter{}

	for i := range t.Len() {
		ft := state.GetTable(t, golua.LNumber(i+1)).(*golua.LTable)
		f := state.GetTable(ft, golua.LString("type")).(golua.LString)
		fltr := filterList[string(f)](state, ft)
		filters = append(filters, fltr)
	}

	g := gift.New(filters...)
	return g
}

func brightnessTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct flt_brightness
	/// @prop type
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_BRIGHTNESS))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func brightnessBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Brightness(float32(percent))
	return f
}

func colorBalanceTable(state *golua.LState, percentRed, percentGreen, percentBlue float64) *golua.LTable {
	/// @struct flt_color_balance
	/// @prop type
	/// @prop percentRed
	/// @prop percentGreen
	/// @prop percentBlue

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_COLOR_BALANCE))
	state.SetTable(t, golua.LString("percentRed"), golua.LNumber(percentRed))
	state.SetTable(t, golua.LString("percentGreen"), golua.LNumber(percentGreen))
	state.SetTable(t, golua.LString("percentBlue"), golua.LNumber(percentBlue))

	return t
}

func colorBalanceBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percentRed := state.GetTable(t, golua.LString("percentRed")).(golua.LNumber)
	percentGreen := state.GetTable(t, golua.LString("percentGreen")).(golua.LNumber)
	percentBlue := state.GetTable(t, golua.LString("percentBlue")).(golua.LNumber)

	f := gift.ColorBalance(float32(percentRed), float32(percentGreen), float32(percentBlue))
	return f
}

func colorizeTable(state *golua.LState, hue, saturation, percent float64) *golua.LTable {
	/// @struct flt_colorize
	/// @prop type
	/// @prop hue
	/// @prop saturation
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_COLORIZE))
	state.SetTable(t, golua.LString("hue"), golua.LNumber(hue))
	state.SetTable(t, golua.LString("saturation"), golua.LNumber(saturation))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func colorizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	hue := state.GetTable(t, golua.LString("hue")).(golua.LNumber)
	saturation := state.GetTable(t, golua.LString("saturation")).(golua.LNumber)
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Colorize(float32(hue), float32(saturation), float32(percent))
	return f
}

func colorspaceLinearSRGBTable(state *golua.LState) *golua.LTable {
	/// @struct flt_colorspace_linear_to_srgb
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_COLORSPACE_LINEAR_TO_SRGB))

	return t
}

func colorspaceLinearSRGBBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.ColorspaceLinearToSRGB()
	return f
}

func colorspaceSRGBLinearTable(state *golua.LState) *golua.LTable {
	/// @struct flt_colorspace_srgb_to_linear
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_COLORSPACE_SRGB_TO_LINEAR))

	return t
}

func colorspaceSRGBLinearBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.ColorspaceSRGBToLinear()
	return f
}

func contrastTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct flt_contrast
	/// @prop type
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_CONTRAST))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func contrastBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Contrast(float32(percent))
	return f
}

func convolutionTable(state *golua.LState, kernel golua.LValue, normalize, alpha, abs bool, delta float64) *golua.LTable {
	/// @struct flt_convolution
	/// @prop type
	/// @prop kernel
	/// @prop normalize
	/// @prop alpha
	/// @prop abs
	/// @prop delta

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_CONVOLUTION))
	state.SetTable(t, golua.LString("kernel"), kernel)
	state.SetTable(t, golua.LString("normalize"), golua.LBool(normalize))
	state.SetTable(t, golua.LString("alpha"), golua.LBool(alpha))
	state.SetTable(t, golua.LString("abs"), golua.LBool(abs))
	state.SetTable(t, golua.LString("delta"), golua.LNumber(delta))

	return t
}

func convolutionBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	kernel := state.GetTable(t, golua.LString("kernel")).(*golua.LTable)
	normalize := state.GetTable(t, golua.LString("normalize")).(golua.LBool)
	alpha := state.GetTable(t, golua.LString("alpha")).(golua.LBool)
	abs := state.GetTable(t, golua.LString("abs")).(golua.LBool)
	delta := state.GetTable(t, golua.LString("delta")).(golua.LNumber)

	kernalMatrix := []float32{}
	for i := range kernel.Len() {
		v := state.GetTable(kernel, golua.LNumber(i+1)).(golua.LNumber)
		kernalMatrix = append(kernalMatrix, float32(v))
	}

	f := gift.Convolution(kernalMatrix, bool(normalize), bool(alpha), bool(abs), float32(delta))
	return f
}

func cropTable(state *golua.LState, xmin, ymin, xmax, ymax int) *golua.LTable {
	/// @struct flt_crop
	/// @prop type
	/// @prop xmin
	/// @prop ymin
	/// @prop xmax
	/// @prop ymax

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_CROP))
	state.SetTable(t, golua.LString("xmin"), golua.LNumber(xmin))
	state.SetTable(t, golua.LString("ymin"), golua.LNumber(ymin))
	state.SetTable(t, golua.LString("xmax"), golua.LNumber(xmax))
	state.SetTable(t, golua.LString("ymax"), golua.LNumber(ymax))

	return t
}

func cropBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	xmin := state.GetTable(t, golua.LString("xmin")).(golua.LNumber)
	ymin := state.GetTable(t, golua.LString("ymin")).(golua.LNumber)
	xmax := state.GetTable(t, golua.LString("xmax")).(golua.LNumber)
	ymax := state.GetTable(t, golua.LString("ymax")).(golua.LNumber)

	f := gift.Crop(image.Rect(int(xmin), int(ymin), int(xmax), int(ymax)))
	return f
}

func cropToSizeTable(state *golua.LState, width, height, anchor int) *golua.LTable {
	/// @struct flt_crop_to_size
	/// @prop type
	/// @prop width
	/// @prop height
	/// @prop anchor

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_CROP_TO_SIZE))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))
	state.SetTable(t, golua.LString("anchor"), golua.LNumber(anchor))

	return t
}

func cropToSizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	height := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	anchor := state.GetTable(t, golua.LString("anchor")).(golua.LNumber)

	f := gift.CropToSize(int(width), int(height), gift.Anchor(anchor))
	return f
}

func flipHorizontalTable(state *golua.LState) *golua.LTable {
	/// @struct flt_flip_horizontal
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_FLIP_HORIZONTAL))

	return t
}

func flipHorizontalBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.FlipHorizontal()
	return f
}

func flipVerticalTable(state *golua.LState) *golua.LTable {
	/// @struct flt_flip_vertical
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_FLIP_VERTICAL))

	return t
}

func flipVerticalBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.FlipVertical()
	return f
}

func gammaTable(state *golua.LState, gamma float64) *golua.LTable {
	/// @struct flt_gamma
	/// @prop type
	/// @prop gamma

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_GAMMA))
	state.SetTable(t, golua.LString("gamma"), golua.LNumber(gamma))

	return t
}

func gammaBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	gamma := state.GetTable(t, golua.LString("gamma")).(golua.LNumber)

	f := gift.Gamma(float32(gamma))
	return f
}

func gaussianBlurTable(state *golua.LState, sigma float64) *golua.LTable {
	/// @struct flt_gaussian_blur
	/// @prop type
	/// @prop sigma

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_GAUSSIAN_BLUR))
	state.SetTable(t, golua.LString("sigma"), golua.LNumber(sigma))

	return t
}

func gaussianBlurBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	sigma := state.GetTable(t, golua.LString("sigma")).(golua.LNumber)

	f := gift.GaussianBlur(float32(sigma))
	return f
}

func grayscaleTable(state *golua.LState) *golua.LTable {
	/// @struct flt_grayscale
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_GRAYSCALE))

	return t
}

func grayscaleBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Grayscale()
	return f
}

func invertTable(state *golua.LState) *golua.LTable {
	/// @struct flt_invert
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_INVERT))

	return t
}

func invertBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Invert()
	return f
}

func rotateTable(state *golua.LState, angle float64, bgcolor golua.LValue, interpolation int) *golua.LTable {
	/// @struct flt_rotate
	/// @prop type
	/// @prop angle
	/// @prop bgcolor
	/// @prop interpolation

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_ROTATE))
	state.SetTable(t, golua.LString("angle"), golua.LNumber(angle))
	state.SetTable(t, golua.LString("bgcolor"), bgcolor)
	state.SetTable(t, golua.LString("interpolation"), golua.LNumber(interpolation))

	return t
}

func rotateBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	angle := state.GetTable(t, golua.LString("angle")).(golua.LNumber)
	bgcolor := state.GetTable(t, golua.LString("bgcolor")).(*golua.LTable)
	interpolation := state.GetTable(t, golua.LString("interpolation")).(golua.LNumber)

	c := imageutil.TableToRGBA(state, bgcolor)

	f := gift.Rotate(float32(angle), c, gift.Interpolation(interpolation))
	return f
}

func rotate180Table(state *golua.LState) *golua.LTable {
	/// @struct flt_rotate_180
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_ROTATE180))

	return t
}

func rotate180Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate180()
	return f
}

func rotate270Table(state *golua.LState) *golua.LTable {
	/// @struct flt_rotate_270
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_ROTATE270))

	return t
}

func rotate270Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate270()
	return f
}

func rotate90Table(state *golua.LState) *golua.LTable {
	/// @struct flt_rotate_90
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_ROTATE90))

	return t
}

func rotate90Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate90()
	return f
}

func hueTable(state *golua.LState, shift float64) *golua.LTable {
	/// @struct flt_hue
	/// @prop type
	/// @prop shift

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_HUE))
	state.SetTable(t, golua.LString("shift"), golua.LNumber(shift))

	return t
}

func hueBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	shift := state.GetTable(t, golua.LString("shift")).(golua.LNumber)

	f := gift.Hue(float32(shift))
	return f
}

func saturationTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct flt_saturation
	/// @prop type
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_SATURATION))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func saturationBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Saturation(float32(percent))
	return f
}

func sepiaTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct flt_sepia
	/// @prop type
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_SEPIA))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func sepiaBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Sepia(float32(percent))
	return f
}

func thresholdTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct flt_threshold
	/// @prop type
	/// @prop percent

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_THRESHOLD))
	state.SetTable(t, golua.LString("percent"), golua.LNumber(percent))

	return t
}

func thresholdBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := state.GetTable(t, golua.LString("percent")).(golua.LNumber)

	f := gift.Threshold(float32(percent))
	return f
}

func pixelateTable(state *golua.LState, size int) *golua.LTable {
	/// @struct flt_pixelate
	/// @prop type
	/// @prop size

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_PIXELATE))
	state.SetTable(t, golua.LString("size"), golua.LNumber(size))

	return t
}

func pixelateBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	size := state.GetTable(t, golua.LString("size")).(golua.LNumber)

	f := gift.Pixelate(int(size))
	return f
}

func transposeTable(state *golua.LState) *golua.LTable {
	/// @struct flt_transpose
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_TRANSPOSE))

	return t
}

func transposeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Transpose()
	return f
}

func transverseTable(state *golua.LState) *golua.LTable {
	/// @struct flt_transverse
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_TRANSVERSE))

	return t
}

func transverseBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Transverse()
	return f
}

func sobelTable(state *golua.LState) *golua.LTable {
	/// @struct flt_sobel
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_SOBEL))

	return t
}

func sobelBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Sobel()
	return f
}

func maximumTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct flt_maximum
	/// @prop type
	/// @prop ksize
	/// @prop disk

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_MAXIMUM))
	state.SetTable(t, golua.LString("ksize"), golua.LNumber(ksize))
	state.SetTable(t, golua.LString("disk"), golua.LBool(disk))

	return t
}

func maximumBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := state.GetTable(t, golua.LString("ksize")).(golua.LNumber)
	disk := state.GetTable(t, golua.LString("disk")).(golua.LBool)

	f := gift.Maximum(int(ksize), bool(disk))
	return f
}

func meanTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct flt_mean
	/// @prop type
	/// @prop ksize
	/// @prop disk

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_MEAN))
	state.SetTable(t, golua.LString("ksize"), golua.LNumber(ksize))
	state.SetTable(t, golua.LString("disk"), golua.LBool(disk))

	return t
}

func meanBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := state.GetTable(t, golua.LString("ksize")).(golua.LNumber)
	disk := state.GetTable(t, golua.LString("disk")).(golua.LBool)

	f := gift.Mean(int(ksize), bool(disk))
	return f
}

func medianTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct flt_median
	/// @prop type
	/// @prop ksize
	/// @prop disk

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_MEDIAN))
	state.SetTable(t, golua.LString("ksize"), golua.LNumber(ksize))
	state.SetTable(t, golua.LString("disk"), golua.LBool(disk))

	return t
}

func medianBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := state.GetTable(t, golua.LString("ksize")).(golua.LNumber)
	disk := state.GetTable(t, golua.LString("disk")).(golua.LBool)

	f := gift.Median(int(ksize), bool(disk))
	return f
}

func minimumTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct flt_minimum
	/// @prop type
	/// @prop ksize
	/// @prop disk

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_MINIMUM))
	state.SetTable(t, golua.LString("ksize"), golua.LNumber(ksize))
	state.SetTable(t, golua.LString("disk"), golua.LBool(disk))

	return t
}

func minimumBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := state.GetTable(t, golua.LString("ksize")).(golua.LNumber)
	disk := state.GetTable(t, golua.LString("disk")).(golua.LBool)

	f := gift.Minimum(int(ksize), bool(disk))
	return f
}

func sigmoidTable(state *golua.LState, midpoint, factor float64) *golua.LTable {
	/// @struct flt_sigmoid
	/// @prop type
	/// @prop midpoint
	/// @prop factor

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_SIGMOID))
	state.SetTable(t, golua.LString("midpoint"), golua.LNumber(midpoint))
	state.SetTable(t, golua.LString("factor"), golua.LNumber(factor))

	return t
}

func sigmoidBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	midpoint := state.GetTable(t, golua.LString("midpoint")).(golua.LNumber)
	factor := state.GetTable(t, golua.LString("factor")).(golua.LNumber)

	f := gift.Sigmoid(float32(midpoint), float32(factor))
	return f
}

func unsharpMaskTable(state *golua.LState, sigma, amount, threshold float64) *golua.LTable {
	/// @struct flt_unsharp_mask
	/// @prop type
	/// @prop sigma
	/// @prop amount
	/// @prop threshold

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_UNSHARP_MASK))
	state.SetTable(t, golua.LString("sigma"), golua.LNumber(sigma))
	state.SetTable(t, golua.LString("amount"), golua.LNumber(amount))
	state.SetTable(t, golua.LString("threshold"), golua.LNumber(threshold))

	return t
}

func unsharpMaskBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	sigma := state.GetTable(t, golua.LString("sigma")).(golua.LNumber)
	amount := state.GetTable(t, golua.LString("amount")).(golua.LNumber)
	threshold := state.GetTable(t, golua.LString("threshold")).(golua.LNumber)

	f := gift.UnsharpMask(float32(sigma), float32(amount), float32(threshold))
	return f
}

func resizeTable(state *golua.LState, width, height, resampling int) *golua.LTable {
	/// @struct flt_resize
	/// @prop type
	/// @prop width
	/// @prop height
	/// @prop resampling

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_RESIZE))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))
	state.SetTable(t, golua.LString("resampling"), golua.LNumber(resampling))

	return t
}

func resizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	height := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	resampling := state.GetTable(t, golua.LString("resampling")).(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.Resize(int(width), int(height), s)
	return f
}

func resizeToFillTable(state *golua.LState, width, height, resampling, anchor int) *golua.LTable {
	/// @struct flt_resize_to_fill
	/// @prop type
	/// @prop width
	/// @prop height
	/// @prop resampling
	/// @prop anchor

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_RESIZE_TO_FILL))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))
	state.SetTable(t, golua.LString("resampling"), golua.LNumber(resampling))
	state.SetTable(t, golua.LString("anchor"), golua.LNumber(anchor))

	return t
}

func resizeToFillBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	height := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	resampling := state.GetTable(t, golua.LString("resampling")).(golua.LNumber)
	anchor := state.GetTable(t, golua.LString("anchor")).(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.ResizeToFill(int(width), int(height), s, gift.Anchor(anchor))
	return f
}

func resizeToFitTable(state *golua.LState, width, height, resampling int) *golua.LTable {
	/// @struct flt_resize_to_fit
	/// @prop type
	/// @prop width
	/// @prop height
	/// @prop resampling

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_RESIZE_TO_FIT))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))
	state.SetTable(t, golua.LString("resampling"), golua.LNumber(resampling))

	return t
}

func resizeToFitBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	height := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	resampling := state.GetTable(t, golua.LString("resampling")).(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.ResizeToFit(int(width), int(height), s)
	return f
}

func colorFuncTable(state *golua.LState, fn *golua.LFunction) *golua.LTable {
	/// @struct flt_color_func
	/// @prop type
	/// @prop fn

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(FILTER_COLOR_FUNC))
	state.SetTable(t, golua.LString("fn"), fn)

	return t
}

func colorFuncBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	fn := state.GetTable(t, golua.LString("fn")).(*golua.LFunction)

	f := gift.ColorFunc(func(r0, g0, b0, a0 float32) (r float32, g float32, b float32, a float32) {
		cfInner, _ := state.NewThread()
		cfInner.Push(fn)
		cfInner.Push(golua.LNumber(r0))
		cfInner.Push(golua.LNumber(g0))
		cfInner.Push(golua.LNumber(b0))
		cfInner.Push(golua.LNumber(a0))
		cfInner.Call(4, 4)

		r1 := cfInner.CheckNumber(-4)
		g1 := cfInner.CheckNumber(-3)
		b1 := cfInner.CheckNumber(-2)
		a1 := cfInner.CheckNumber(-1)
		return float32(r1), float32(g1), float32(b1), float32(a1)
	})
	return f
}