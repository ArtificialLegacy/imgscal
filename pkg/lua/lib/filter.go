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

/// @lib Filter
/// @import filter
/// @desc
/// Library for applying lists of filters onto images.

func RegisterFilter(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_FILTER, r, r.State, lg)

	/// @func draw(id1, id2, filters, disableParallelization?)
	/// @arg id1 {int<collection.IMAGE>}
	/// @arg id2 {int<collection.IMAGE>}
	/// @arg filters {[]struct<filter.Filter>}
	/// @arg? disableParallelization {bool}
	/// @desc
	/// Applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.RAW_TABLE, Name: "filters"},
			{Type: lua.BOOL, Name: "disableParallelization", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image
			scheduledState, _ := state.NewThread()

			r.IC.SchedulePipe(args["id1"].(int), args["id2"].(int),
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
						g := buildFilterList(scheduledState, filters, args["filters"].(*golua.LTable))
						if args["disableParallelization"].(bool) {
							g.SetParallelization(false)
						}
						g.Draw(imageutil.ImageGetDraw(i.Self.Image), img)

						scheduledState.Close()
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						scheduledState.Close()
					},
				})

			return 0
		})

	/// @func draw_at(id1, id2, point, op, filters, disableParallelization?)
	/// @arg id1 {int<collection.IMAGE>}
	/// @arg id2 {int<collection.IMAGE>}
	/// @arg point {struct<image.Point>}
	/// @arg op {int<filter.Operator>}
	/// @arg filters {[]struct<filter.Filter>}
	/// @arg? disableParallelization {bool}
	/// @desc
	/// Applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw_at",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.RAW_TABLE, Name: "point"},
			{Type: lua.INT, Name: "op"},
			{Type: lua.RAW_TABLE, Name: "filters"},
			{Type: lua.BOOL, Name: "disableParallelization", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image
			scheduledState, _ := state.NewThread()

			r.IC.SchedulePipe(args["id1"].(int), args["id2"].(int),
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
						g := buildFilterList(scheduledState, filters, args["filters"].(*golua.LTable))
						pt := imageutil.TableToPoint(args["point"].(*golua.LTable))
						if args["disableParallelization"].(bool) {
							g.SetParallelization(false)
						}
						g.DrawAt(imageutil.ImageGetDraw(i.Self.Image), img, pt, gift.Operator(args["op"].(int)))

						scheduledState.Close()
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						scheduledState.Close()
					},
				})

			return 0
		})

	/// @func draw_at_xy(id1, id2, x, y, op, filters, disableParallelization?)
	/// @arg id1 {int<collection.IMAGE>}
	/// @arg id2 {int<collection.IMAGE>}
	/// @arg x {int}
	/// @arg y {int}
	/// @arg op {int<filter.Operator>}
	/// @arg filters {[]struct<filter.Filter>}
	/// @arg? disableParallelization {bool}
	/// @desc
	/// Applies the filters to image1 with the output going into image2.
	lib.CreateFunction(tab, "draw_at_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "id1"},
			{Type: lua.INT, Name: "id2"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.INT, Name: "op"},
			{Type: lua.RAW_TABLE, Name: "filters"},
			{Type: lua.BOOL, Name: "disableParallelization", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image
			scheduledState, _ := state.NewThread()

			r.IC.SchedulePipe(args["id1"].(int), args["id2"].(int),
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
						g := buildFilterList(scheduledState, filters, args["filters"].(*golua.LTable))
						if args["disableParallelization"].(bool) {
							g.SetParallelization(false)
						}
						g.DrawAt(
							imageutil.ImageGetDraw(i.Self.Image), img,
							image.Point{X: args["x"].(int), Y: args["y"].(int)},
							gift.Operator(args["op"].(int)),
						)

						scheduledState.Close()
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						scheduledState.Close()
					},
				})

			return 0
		})

	/// @func bounds(id, filters, disableParallelization?) -> int, int, int, int
	/// @arg id {int<collection.IMAGE>}
	/// @arg filters {[]struct<filter.Filter>}
	/// @arg? disableParallelization {bool}
	/// @returns {int} - X position of the top left corner of the images bounds after the filters are applied.
	/// @returns {int} - Y position of the top left corner of the images bounds after the filters are applied.
	/// @returns {int} - X position of the bottom right corner of the images bounds after the filters are applied.
	/// @returns {int} - Y position of the bottom right corner of the images bounds after the filters are applied.
	/// @blocking
	/// @desc
	/// Gets the resulting bounds of the image after the filters are applied.
	lib.CreateFunction(tab, "bounds",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "filters"},
			{Type: lua.BOOL, Name: "disableParallelization", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var dstBounds image.Rectangle

			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					if args["disableParallelization"].(bool) {
						g.SetParallelization(false)
					}
					dstBounds = g.Bounds(i.Self.Image.Bounds())
				},
			})

			state.Push(golua.LNumber(dstBounds.Min.X))
			state.Push(golua.LNumber(dstBounds.Min.Y))
			state.Push(golua.LNumber(dstBounds.Max.X))
			state.Push(golua.LNumber(dstBounds.Max.Y))
			return 4
		})

	/// @func bounds_size(id, filters, disableParallelization?) -> int, int
	/// @arg id {int<collection.IMAGE>}
	/// @arg filters {[]struct<filter.Filter>}
	/// @arg? disableParallelization {bool}
	/// @returns {int} - Width of the image after the filters are applied.
	/// @returns {int} - Height of the image after the filters are applied.
	/// @blocking
	/// @desc
	/// Gets the resulting size of the image after the filters are applied.
	lib.CreateFunction(tab, "bounds_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "filters"},
			{Type: lua.BOOL, Name: "disableParallelization", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var dstBounds image.Rectangle
			<-r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					g := buildFilterList(state, filters, args["filters"].(*golua.LTable))
					if args["disableParallelization"].(bool) {
						g.SetParallelization(false)
					}
					dstBounds = g.Bounds(i.Self.Image.Bounds())
				},
			})

			state.Push(golua.LNumber(dstBounds.Dx()))
			state.Push(golua.LNumber(dstBounds.Dy()))
			return 2
		})

	/// @func brightness(percent) -> struct<filter.FilterBrightness>
	/// @arg percent {float} - Percent value between -100 and 100.
	/// @returns {struct<filter.FilterBrightness>}
	lib.CreateFunction(tab, "brightness",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := brightnessTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func color_balance(percentRed, percentGreen, percentBlue) -> struct<filter.FilterColorBalance>
	/// @arg percentRed {float} - Percent values between -100 and 500.
	/// @arg percentGreen {float} - Percent values between -100 and 500.
	/// @arg percentBlue {float} - Percent values between -100 and 500.
	/// @returns {struct<filter.FilterColorBalance>}
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

	/// @func colorize(hue, saturation, percent) -> struct<filter.FilterColorize>
	/// @arg hue {float} - Hue value between 0 and 360.
	/// @arg saturation {float} - Saturation value between 0 and 100.
	/// @arg percent {float} - Percent value between 0 and 100.
	/// @returns {struct<filter.FilterColorize>}
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

	/// @func colorspace_linear_to_srgb() -> struct<filter.FilterColorspaceLinearToSRGB>
	/// @returns {struct<filter.FilterColorspaceLinearToSRGB>}
	lib.CreateFunction(tab, "colorspace_linear_to_srgb",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorspaceLinearSRGBTable(state)

			state.Push(t)
			return 1
		})

	/// @func colorspace_srgb_to_linear() -> struct<filter.FilterColorspaceSRGBToLinear>
	/// @returns {struct<filter.FilterColorspaceSRGBToLinear>}
	lib.CreateFunction(tab, "colorspace_srgb_to_linear",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorspaceSRGBLinearTable(state)

			state.Push(t)
			return 1
		})

	/// @func contrast(percent) -> struct<filter.FilterContrast>
	/// @arg percent {float} - Percent values between -100 and 100.
	/// @returns {struct<filter.FilterContrast>}
	lib.CreateFunction(tab, "contrast",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := contrastTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func convolution(kernel, normalize, alpha, abs, delta) -> struct<filter.FilterConvolution>
	/// @arg kernel {[]float} - Must have a length of an odd square, e.g. 3x3=9 or 5x5=25.
	/// @arg normalize {bool}
	/// @arg alpha {bool}
	/// @arg abs {bool}
	/// @arg delta {float}
	/// @returns {struct<filter.FilterConvolution>}
	lib.CreateFunction(tab, "convolution",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "kernel"},
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

	/// @func crop(min, max) -> struct<filter.FilterCrop>
	/// @arg min {struct<image.Point>} - Top left corner of new image.
	/// @arg max {struct<image.Point>} - Bottom right corner of new image.
	/// @returns {struct<filter.FilterCrop>}
	lib.CreateFunction(tab, "crop",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "min"},
			{Type: lua.RAW_TABLE, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			min := imageutil.TableToPoint(args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(args["max"].(*golua.LTable))

			t := cropTable(state, min.X, min.Y, max.X, max.Y)

			state.Push(t)
			return 1
		})

	/// @func crop_xy(xmin, ymin, xmax, ymax) -> struct<filter.FilterCrop>
	/// @arg xmin {int} - X position of the top left corner.
	/// @arg ymin {int} - Y position of the top left corner.
	/// @arg xmax {int} - X position of the bottom right corner.
	/// @arg ymax {int} - Y position of the bottom right corner.
	/// @returns {struct<filter.FilterCrop>}
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

	/// @func crop_to_size(width, height, anchor) -> struct<filter.FilterCropToSize>
	/// @arg width {int} - Width of the new image.
	/// @arg height {int} - Height of the new image.
	/// @arg anchor {int<filter.Anchor>}
	/// @returns {struct<filter.FilterCropToSize>}
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

	/// @func flip_horizontal() -> struct<filter.FilterFlipHorizontal>
	/// @returns {struct<filter.FilterFlipHorizontal>}
	lib.CreateFunction(tab, "flip_horizontal",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := flipHorizontalTable(state)

			state.Push(t)
			return 1
		})

	/// @func flip_vertical() -> struct<filter.FilterFlipVertical>
	/// @returns {struct<filter.FilterFlipVertical>}
	lib.CreateFunction(tab, "flip_vertical",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := flipVerticalTable(state)

			state.Push(t)
			return 1
		})

	/// @func gamma(gamma) -> struct<filter.FilterGamma>
	/// @arg gamma {float} - Must be positive, a value of 1 maintains the image.
	/// @returns {struct<filter.FilterGamma>}
	lib.CreateFunction(tab, "gamma",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "gamma"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := gammaTable(state, args["gamma"].(float64))

			state.Push(t)
			return 1
		})

	/// @func gaussian_blur(sigma) -> struct<filter.FilterGaussianBlur>
	/// @arg sigma {float} - Radius of blur is ~3*sigma.
	/// @returns {struct<filter.FilterGaussianBlur>}
	lib.CreateFunction(tab, "gaussian_blur",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "sigma"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := gaussianBlurTable(state, args["sigma"].(float64))

			state.Push(t)
			return 1
		})

	/// @func grayscale() -> struct<filter.FilterGreyscale>
	/// @returns {struct<filter.FilterGreyscale>}
	lib.CreateFunction(tab, "grayscale",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := grayscaleTable(state)

			state.Push(t)
			return 1
		})

	/// @func invert() -> struct<filter.FilterInvert>
	/// @returns {struct<filter.FilterInvert>}
	lib.CreateFunction(tab, "invert",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := invertTable(state)

			state.Push(t)
			return 1
		})

	/// @func rotate(angle, bgcolor, interpolation) -> struct<filter.FilterRotate>
	/// @arg angle {float} - Angle value in degrees.
	/// @arg bgcolor {struct<image.Color>}
	/// @arg interpolation {int<filter.Interpolation>}
	/// @returns {struct<filter.FilterRotate>}
	/// @desc
	/// When doing rotations on multiples of 90 degrees, it is best to use the rotate filters for 90, 180, and 270.
	lib.CreateFunction(tab, "rotate",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.RAW_TABLE, Name: "bgcolor"},
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

	/// @func rotate_90() -> struct<filter.FilterRotate90>
	/// @returns {struct<filter.FilterRotate90>}
	lib.CreateFunction(tab, "rotate_90",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate90Table(state)

			state.Push(t)
			return 1
		})

	/// @func rotate_180() -> struct<filter.FilterRotate180>
	/// @returns {struct<filter.FilterRotate180>}
	lib.CreateFunction(tab, "rotate_180",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate180Table(state)

			state.Push(t)
			return 1
		})

	/// @func rotate_270() -> struct<filter.FilterRotate270>
	/// @returns {struct<filter.FilterRotate270>}
	lib.CreateFunction(tab, "rotate_270",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := rotate270Table(state)

			state.Push(t)
			return 1
		})

	/// @func hue(shift) -> struct<filter.FilterHue>
	/// @arg shift {float} - Shift value between -180 and 180.
	/// @returns {struct<filter.FilterHue>}
	lib.CreateFunction(tab, "hue",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "shift"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := hueTable(state, args["shift"].(float64))

			state.Push(t)
			return 1
		})

	/// @func saturation(percent) -> struct<filter.FilterSaturation>
	/// @arg percent {float} - Percent value between -100 and 500.
	/// @returns {struct<filter.FilterSaturation>}
	lib.CreateFunction(tab, "saturation",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := saturationTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func sepia(percent) -> struct<filter.FilterSepia>
	/// @arg percent {float} - Percent value between 0 and 100.
	/// @returns {struct<filter.FilterSepia>}
	lib.CreateFunction(tab, "sepia",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := sepiaTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func pixelate(size) -> struct<filter.FilterPixelate>
	/// @arg size {int}
	/// @returns {struct<filter.FilterPixelate>}
	lib.CreateFunction(tab, "pixelate",
		[]lua.Arg{
			{Type: lua.INT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := pixelateTable(state, args["size"].(int))

			state.Push(t)
			return 1
		})

	/// @func threshold(percent) -> struct<filter.FilterThreshold>
	/// @arg percent {float} - Percent value between 0 and 100.
	/// @returns {struct<filter.FilterThreshold>}
	lib.CreateFunction(tab, "threshold",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "percent"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := thresholdTable(state, args["percent"].(float64))

			state.Push(t)
			return 1
		})

	/// @func transpose() -> struct<filter.FilterTranspose>
	/// @returns {struct<filter.FilterTranspose>}
	lib.CreateFunction(tab, "transpose",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := transposeTable(state)

			state.Push(t)
			return 1
		})

	/// @func transverse() -> struct<filter.FilterTransverse>
	/// @returns {struct<filter.FilterTranverse>}
	lib.CreateFunction(tab, "transverse",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := transverseTable(state)

			state.Push(t)
			return 1
		})

	/// @func sobel() -> struct<filter.FilterSobel>
	/// @returns {struct<filter.FilterSobel>}
	lib.CreateFunction(tab, "sobel",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := sobelTable(state)

			state.Push(t)
			return 1
		})

	/// @func maximum(ksize, disk) -> struct<filter.FilterMaximum>
	/// @arg ksize {int} - Must be odd int, e.g. 3, 5, 7.
	/// @arg disk {bool} - If the kernel used should be disk shaped instead of a square.
	/// @returns {struct<filter.FilterMaximum>}
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

	/// @func mean(ksize, disk) -> struct<filter.FilterMean>
	/// @arg ksize {int} - Must be odd int, e.g. 3, 5, 7.
	/// @arg disk {bool} - If the kernel used should be disk shaped instead of a square.
	/// @returns {struct<filter.FilterMean>}
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

	/// @func median(ksize, disk) -> struct<filter.FilterMedian>
	/// @arg ksize {int} - Must be odd int, e.g. 3, 5, 7.
	/// @arg disk {bool} - If the kernel used should be disk shaped instead of a square.
	/// @returns {struct<filter.FilterMedian>}
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

	/// @func minimum(ksize, disk) -> struct<filter.FilterMinimum>
	/// @arg ksize {int} - Must be odd int, e.g. 3, 5, 7.
	/// @arg disk {bool} - If the kernel used should be disk shaped instead of a square.
	/// @returns {struct<filter.FilterMinimum>}
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

	/// @func sigmoid(midpoint, factor) -> struct<filter.FilterSigmoid>
	/// @arg midpoint {float} - Value between 0 and 1.
	/// @arg factor {float}
	/// @returns {struct<filter.FilterSigmoid>}
	/// @desc
	/// Adjusts the contrast of the image, but on a curve instead of linearly.
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

	/// @func unsharp_mask(sigma, amount, threshold) -> struct<filter.FilterUnsharpMask>
	/// @arg sigma {float} - Radius is ~3*sigma.
	/// @arg amount {float} - Typically between 0.5 and 1.5.
	/// @arg threshold {float} - Minimum brightness change to sharpen, typically between 0 and 0.05.
	/// @returns struct<filter.FilterUnsharpMask>
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

	/// @func resize(width, height, resampling) -> struct<filter.FilterResize>
	/// @arg width {int}
	/// @arg height {int}
	/// @arg resampling {int<filter.Resampling>}
	/// @returns {struct<filter.FilterResize>}
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

	/// @func resize_to_fill(width, height, resampling, anchor) -> struct<filter.FilterResizeToFill>
	/// @arg width {int}
	/// @arg height {int}
	/// @arg resampling {int<filter.Resampling>}
	/// @arg anchor {int<filter.Anchor>}
	/// @returns {struct<filter.FilterResizeToFill>}
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

	/// @func resize_to_fit(width, height, resampling) -> struct<filter.FilterResizeToFit>
	/// @arg width {int}
	/// @arg height {int}
	/// @arg resampling {int<filter.Resampling>}
	/// @returns {struct<filter.FilterResizeToFit>}
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

	/// @func color_func(fn) -> struct<filter.FilterColorFunc>
	/// @arg fn {function(r float, g float, b float, a float) -> float, float, float, float}
	/// @returns struct<filter.FilterColorFunc>
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

	/// @func color_func_unsafe(fn) -> struct<filter.FilterColorFuncUnsafe>
	/// @arg fn {function(r float, g float, b float, a float) -> float, float, float, float}
	/// @returns {struct<filter.FilterColorFuncUnsafe>}
	/// @desc
	/// Color values are floats between 0 and 1.
	/// Note parallelization must be disabled when drawing for this to work.
	/// This has the benefit of not requiring new lua threads on each function call,
	/// but this means it is not thread safe.
	lib.CreateFunction(tab, "color_func_unsafe",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorFuncUnsafeTable(state, args["fn"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @constants Anchor {int}
	/// @const ANCHOR_CENTER
	/// @const ANCHOR_TOPLEFT
	/// @const ANCHOR_TOP
	/// @const ANCHOR_TOPRIGHT
	/// @const ANCHOR_LEFT
	/// @const ANCHOR_RIGHT
	/// @const ANCHOR_BOTTOMLEFT
	/// @const ANCHOR_BOTTOM
	/// @const ANCHOR_BOTTOMRIGHT
	tab.RawSetString("ANCHOR_CENTER", golua.LNumber(gift.CenterAnchor))
	tab.RawSetString("ANCHOR_TOPLEFT", golua.LNumber(gift.TopLeftAnchor))
	tab.RawSetString("ANCHOR_TOP", golua.LNumber(gift.TopAnchor))
	tab.RawSetString("ANCHOR_TOPRIGHT", golua.LNumber(gift.TopRightAnchor))
	tab.RawSetString("ANCHOR_LEFT", golua.LNumber(gift.LeftAnchor))
	tab.RawSetString("ANCHOR_RIGHT", golua.LNumber(gift.RightAnchor))
	tab.RawSetString("ANCHOR_BOTTOMLEFT", golua.LNumber(gift.BottomLeftAnchor))
	tab.RawSetString("ANCHOR_BOTTOM", golua.LNumber(gift.BottomAnchor))
	tab.RawSetString("ANCHOR_BOTTOMRIGHT", golua.LNumber(gift.BottomRightAnchor))

	/// @constants Interpolation {int}
	/// @const INTERPOLATION_NEARESTNEIGHBOR
	/// @const INTERPOLATION_LINEAR
	/// @const INTERPOLATION_CUBIC
	tab.RawSetString("INTERPOLATION_NEARESTNEIGHBOR", golua.LNumber(gift.NearestNeighborInterpolation))
	tab.RawSetString("INTERPOLATION_LINEAR", golua.LNumber(gift.LinearInterpolation))
	tab.RawSetString("INTERPOLATION_CUBIC", golua.LNumber(gift.CubicInterpolation))

	/// @constants Operator {int}
	/// @const OPERATOR_COPY
	/// @const OPERATOR_OVER
	tab.RawSetString("OPERATOR_COPY", golua.LNumber(gift.CopyOperator))
	tab.RawSetString("OPERATOR_OVER", golua.LNumber(gift.OverOperator))

	/// @constants Resampling {int}
	/// @const RESAMPLING_BOX
	/// @const RESAMPLING_CUBIC
	/// @const RESAMPLING_LANCZOS
	/// @const RESAMPLING_LINEAR
	/// @const RESAMPLING_NEARESTNEIGHBOR
	tab.RawSetString("RESAMPLING_BOX", golua.LNumber(RESAMPLING_BOX))
	tab.RawSetString("RESAMPLING_CUBIC", golua.LNumber(RESAMPLING_CUBIC))
	tab.RawSetString("RESAMPLING_LANCZOS", golua.LNumber(RESAMPLING_LANCZOS))
	tab.RawSetString("RESAMPLING_LINEAR", golua.LNumber(RESAMPLING_LINEAR))
	tab.RawSetString("RESAMPLING_NEARESTNEIGHBOR", golua.LNumber(RESAMPLING_NEARESTNEIGHBOR))

	/// @constants FilterType {string}
	/// @const FILTER_BRIGHTNESS
	/// @const FILTER_COLOR_BALANCE
	/// @const FILTER_COLORIZE
	/// @const FILTER_COLORSPACE_LINEAR_TO_SRGB
	/// @const FILTER_COLORSPACE_SRGB_TO_LINEAR
	/// @const FILTER_CONTRAST
	/// @const FILTER_CONVOLUTION
	/// @const FILTER_CROP
	/// @const FILTER_CROP_TO_SIZE
	/// @const FILTER_FLIP_HORIZONTAL
	/// @const FILTER_FLIP_VERTICAL
	/// @const FILTER_GAMMA
	/// @const FILTER_GAUSSIAN_BLUR
	/// @const FILTER_GRAYSCALE
	/// @const FILTER_INVERT
	/// @const FILTER_ROTATE
	/// @const FILTER_ROTATE90
	/// @const FILTER_ROTATE180
	/// @const FILTER_ROTATE270
	/// @const FILTER_HUE
	/// @const FILTER_SATURATION
	/// @const FILTER_SEPIA
	/// @const FILTER_THRESHOLD
	/// @const FILTER_PIXELATE
	/// @const FILTER_SOBEL
	/// @const FILTER_TRANSPOSE
	/// @const FILTER_TRANSVERSE
	/// @const FILTER_MAXIMUM
	/// @const FILTER_MEAN
	/// @const FILTER_MEDIAN
	/// @const FILTER_MINIMUM
	/// @const FILTER_SIGMOID
	/// @const FILTER_UNSHARP_MASK
	/// @const FILTER_RESIZE
	/// @const FILTER_RESIZE_TO_FILL
	/// @const FILTER_RESIZE_TO_FIT
	/// @const FILTER_COLOR_FUNC
	/// @const FILTER_COLOR_FUNC_UNSAFE
	tab.RawSetString("FILTER_BRIGHTNESS", golua.LString(FILTER_BRIGHTNESS))
	tab.RawSetString("FILTER_COLOR_BALANCE", golua.LString(FILTER_COLOR_BALANCE))
	tab.RawSetString("FILTER_COLORIZE", golua.LString(FILTER_COLORIZE))
	tab.RawSetString("FILTER_COLORSPACE_LINEAR_TO_SRGB", golua.LString(FILTER_COLORSPACE_LINEAR_TO_SRGB))
	tab.RawSetString("FILTER_COLORSPACE_SRGB_TO_LINEAR", golua.LString(FILTER_COLORSPACE_SRGB_TO_LINEAR))
	tab.RawSetString("FILTER_CONTRAST", golua.LString(FILTER_CONTRAST))
	tab.RawSetString("FILTER_CONVOLUTION", golua.LString(FILTER_CONVOLUTION))
	tab.RawSetString("FILTER_CROP", golua.LString(FILTER_CROP))
	tab.RawSetString("FILTER_CROP_TO_SIZE", golua.LString(FILTER_CROP_TO_SIZE))
	tab.RawSetString("FILTER_FLIP_HORIZONTAL", golua.LString(FILTER_FLIP_HORIZONTAL))
	tab.RawSetString("FILTER_FLIP_VERTICAL", golua.LString(FILTER_FLIP_VERTICAL))
	tab.RawSetString("FILTER_GAMMA", golua.LString(FILTER_GAMMA))
	tab.RawSetString("FILTER_GAUSSIAN_BLUR", golua.LString(FILTER_GAUSSIAN_BLUR))
	tab.RawSetString("FILTER_GRAYSCALE", golua.LString(FILTER_GRAYSCALE))
	tab.RawSetString("FILTER_INVERT", golua.LString(FILTER_INVERT))
	tab.RawSetString("FILTER_ROTATE", golua.LString(FILTER_ROTATE))
	tab.RawSetString("FILTER_ROTATE90 ", golua.LString(FILTER_ROTATE90))
	tab.RawSetString("FILTER_ROTATE180 ", golua.LString(FILTER_ROTATE180))
	tab.RawSetString("FILTER_ROTATE270 ", golua.LString(FILTER_ROTATE270))
	tab.RawSetString("FILTER_HUE", golua.LString(FILTER_HUE))
	tab.RawSetString("FILTER_SATURATION", golua.LString(FILTER_SATURATION))
	tab.RawSetString("FILTER_SEPIA", golua.LString(FILTER_SEPIA))
	tab.RawSetString("FILTER_THRESHOLD", golua.LString(FILTER_THRESHOLD))
	tab.RawSetString("FILTER_PIXELATE", golua.LString(FILTER_PIXELATE))
	tab.RawSetString("FILTER_SOBEL", golua.LString(FILTER_SOBEL))
	tab.RawSetString("FILTER_TRANSPOSE", golua.LString(FILTER_TRANSPOSE))
	tab.RawSetString("FILTER_TRANSVERSE", golua.LString(FILTER_TRANSVERSE))
	tab.RawSetString("FILTER_MAXIMUM", golua.LString(FILTER_MAXIMUM))
	tab.RawSetString("FILTER_MEAN", golua.LString(FILTER_MEAN))
	tab.RawSetString("FILTER_MEDIAN", golua.LString(FILTER_MEDIAN))
	tab.RawSetString("FILTER_MINIMUM", golua.LString(FILTER_MINIMUM))
	tab.RawSetString("FILTER_SIGMOID", golua.LString(FILTER_SIGMOID))
	tab.RawSetString("FILTER_UNSHARP_MASK", golua.LString(FILTER_UNSHARP_MASK))
	tab.RawSetString("FILTER_RESIZE", golua.LString(FILTER_RESIZE))
	tab.RawSetString("FILTER_RESIZE_TO_FILL", golua.LString(FILTER_RESIZE_TO_FILL))
	tab.RawSetString("FILTER_RESIZE_TO_FIT", golua.LString(FILTER_RESIZE_TO_FIT))
	tab.RawSetString("FILTER_COLOR_FUNC", golua.LString(FILTER_COLOR_FUNC))
	tab.RawSetString("FILTER_COLOR_FUNC_UNSAFE", golua.LString(FILTER_COLOR_FUNC_UNSAFE))
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
	FILTER_ROTATE90                  = "rotate_90"
	FILTER_ROTATE180                 = "rotate_180"
	FILTER_ROTATE270                 = "rotate_270"
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
	FILTER_COLOR_FUNC_UNSAFE         = "color_func_unsafe"
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
	FILTER_COLOR_FUNC_UNSAFE:         colorFuncUnsafeBuild,
}

func buildFilterList(state *golua.LState, filterList filterList, t *golua.LTable) *gift.GIFT {
	/// @interface Filter
	/// @prop type {string<filter.FilterType>}

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
	/// @struct FilterBrightness
	/// @prop type {string<filter.FilterType>}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_BRIGHTNESS))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func brightnessBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Brightness(float32(percent))
	return f
}

func colorBalanceTable(state *golua.LState, percentRed, percentGreen, percentBlue float64) *golua.LTable {
	/// @struct FilterColorBalance
	/// @prop type {string<filter.FilterType>}
	/// @prop percentRed {float}
	/// @prop percentGreen {float}
	/// @prop percentBlue {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLOR_BALANCE))
	t.RawSetString("percentRed", golua.LNumber(percentRed))
	t.RawSetString("percentGreen", golua.LNumber(percentGreen))
	t.RawSetString("percentBlue", golua.LNumber(percentBlue))

	return t
}

func colorBalanceBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percentRed := t.RawGetString("percentRed").(golua.LNumber)
	percentGreen := t.RawGetString("percentGreen").(golua.LNumber)
	percentBlue := t.RawGetString("percentBlue").(golua.LNumber)

	f := gift.ColorBalance(float32(percentRed), float32(percentGreen), float32(percentBlue))
	return f
}

func colorizeTable(state *golua.LState, hue, saturation, percent float64) *golua.LTable {
	/// @struct FilterColorize
	/// @prop type {string<filter.FilterType>}
	/// @prop hue {float}
	/// @prop saturation {float}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLORIZE))
	t.RawSetString("hue", golua.LNumber(hue))
	t.RawSetString("saturation", golua.LNumber(saturation))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func colorizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	hue := t.RawGetString("hue").(golua.LNumber)
	saturation := t.RawGetString("saturation").(golua.LNumber)
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Colorize(float32(hue), float32(saturation), float32(percent))
	return f
}

func colorspaceLinearSRGBTable(state *golua.LState) *golua.LTable {
	/// @struct FilterColorspaceLinearToSRGB
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLORSPACE_LINEAR_TO_SRGB))

	return t
}

func colorspaceLinearSRGBBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.ColorspaceLinearToSRGB()
	return f
}

func colorspaceSRGBLinearTable(state *golua.LState) *golua.LTable {
	/// @struct FilterColorspaceSRGBToLinear
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLORSPACE_SRGB_TO_LINEAR))

	return t
}

func colorspaceSRGBLinearBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.ColorspaceSRGBToLinear()
	return f
}

func contrastTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct FilterContrast
	/// @prop type {string<filter.FilterType>}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_CONTRAST))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func contrastBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Contrast(float32(percent))
	return f
}

func convolutionTable(state *golua.LState, kernel golua.LValue, normalize, alpha, abs bool, delta float64) *golua.LTable {
	/// @struct FilterConvolution
	/// @prop type {string<filter.FilterType>}
	/// @prop kernel {[]float}
	/// @prop normalize {bool}
	/// @prop alpha {bool}
	/// @prop abs {bool}
	/// @prop delta {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_CONVOLUTION))
	t.RawSetString("kernel", kernel)
	t.RawSetString("normalize", golua.LBool(normalize))
	t.RawSetString("alpha", golua.LBool(alpha))
	t.RawSetString("abs", golua.LBool(abs))
	t.RawSetString("delta", golua.LNumber(delta))

	return t
}

func convolutionBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	kernel := t.RawGetString("kernel").(*golua.LTable)
	normalize := t.RawGetString("normalize").(golua.LBool)
	alpha := t.RawGetString("alpha").(golua.LBool)
	abs := t.RawGetString("abs").(golua.LBool)
	delta := t.RawGetString("delta").(golua.LNumber)

	kernalMatrix := []float32{}
	for i := range kernel.Len() {
		v := state.GetTable(kernel, golua.LNumber(i+1)).(golua.LNumber)
		kernalMatrix = append(kernalMatrix, float32(v))
	}

	f := gift.Convolution(kernalMatrix, bool(normalize), bool(alpha), bool(abs), float32(delta))
	return f
}

func cropTable(state *golua.LState, xmin, ymin, xmax, ymax int) *golua.LTable {
	/// @struct FilterCrop
	/// @prop type {string<filter.FilterType>}
	/// @prop xmin {int}
	/// @prop ymin {int}
	/// @prop xmax {int}
	/// @prop ymax {int}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_CROP))
	t.RawSetString("xmin", golua.LNumber(xmin))
	t.RawSetString("ymin", golua.LNumber(ymin))
	t.RawSetString("xmax", golua.LNumber(xmax))
	t.RawSetString("ymax", golua.LNumber(ymax))

	return t
}

func cropBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	xmin := t.RawGetString("xmin").(golua.LNumber)
	ymin := t.RawGetString("ymin").(golua.LNumber)
	xmax := t.RawGetString("xmax").(golua.LNumber)
	ymax := t.RawGetString("ymax").(golua.LNumber)

	f := gift.Crop(image.Rect(int(xmin), int(ymin), int(xmax), int(ymax)))
	return f
}

func cropToSizeTable(state *golua.LState, width, height, anchor int) *golua.LTable {
	/// @struct FilterCropToSize
	/// @prop type {string<filter.FilterType>}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop anchor {int<filter.Anchor>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_CROP_TO_SIZE))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("anchor", golua.LNumber(anchor))

	return t
}

func cropToSizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := t.RawGetString("width").(golua.LNumber)
	height := t.RawGetString("height").(golua.LNumber)
	anchor := t.RawGetString("anchor").(golua.LNumber)

	f := gift.CropToSize(int(width), int(height), gift.Anchor(anchor))
	return f
}

func flipHorizontalTable(state *golua.LState) *golua.LTable {
	/// @struct FilterFlipHorizontal
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_FLIP_HORIZONTAL))

	return t
}

func flipHorizontalBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.FlipHorizontal()
	return f
}

func flipVerticalTable(state *golua.LState) *golua.LTable {
	/// @struct FilterFlipVertical
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_FLIP_VERTICAL))

	return t
}

func flipVerticalBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.FlipVertical()
	return f
}

func gammaTable(state *golua.LState, gamma float64) *golua.LTable {
	/// @struct FilterGamma
	/// @prop type {string<filter.FilterType>}
	/// @prop gamma {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_GAMMA))
	t.RawSetString("gamma", golua.LNumber(gamma))

	return t
}

func gammaBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	gamma := t.RawGetString("gamma").(golua.LNumber)

	f := gift.Gamma(float32(gamma))
	return f
}

func gaussianBlurTable(state *golua.LState, sigma float64) *golua.LTable {
	/// @struct FilterGaussianBlur
	/// @prop type {string<filter.FilterType>}
	/// @prop sigma {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_GAUSSIAN_BLUR))
	t.RawSetString("sigma", golua.LNumber(sigma))

	return t
}

func gaussianBlurBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	sigma := t.RawGetString("sigma").(golua.LNumber)

	f := gift.GaussianBlur(float32(sigma))
	return f
}

func grayscaleTable(state *golua.LState) *golua.LTable {
	/// @struct FilterGrayscale
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_GRAYSCALE))

	return t
}

func grayscaleBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Grayscale()
	return f
}

func invertTable(state *golua.LState) *golua.LTable {
	/// @struct FilterInvert
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_INVERT))

	return t
}

func invertBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Invert()
	return f
}

func rotateTable(state *golua.LState, angle float64, bgcolor golua.LValue, interpolation int) *golua.LTable {
	/// @struct FilterRotate
	/// @prop type {string<filter.FilterType>}
	/// @prop angle {float}
	/// @prop bgcolor {struct<image.Color>}
	/// @prop interpolation {int<filter.Interpolation>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_ROTATE))
	t.RawSetString("angle", golua.LNumber(angle))
	t.RawSetString("bgcolor", bgcolor)
	t.RawSetString("interpolation", golua.LNumber(interpolation))

	return t
}

func rotateBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	angle := t.RawGetString("angle").(golua.LNumber)
	bgcolor := t.RawGetString("bgcolor").(*golua.LTable)
	interpolation := t.RawGetString("interpolation").(golua.LNumber)

	c := imageutil.ColorTableToRGBAColor(bgcolor)

	f := gift.Rotate(float32(angle), c, gift.Interpolation(interpolation))
	return f
}

func rotate90Table(state *golua.LState) *golua.LTable {
	/// @struct FilterRotate90
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_ROTATE90))

	return t
}

func rotate90Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate90()
	return f
}

func rotate180Table(state *golua.LState) *golua.LTable {
	/// @struct FilterRotate180
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_ROTATE180))

	return t
}

func rotate180Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate180()
	return f
}

func rotate270Table(state *golua.LState) *golua.LTable {
	/// @struct FilterRotate270
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_ROTATE270))

	return t
}

func rotate270Build(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Rotate270()
	return f
}

func hueTable(state *golua.LState, shift float64) *golua.LTable {
	/// @struct FilterHue
	/// @prop type {string<filter.FilterType>}
	/// @prop shift {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_HUE))
	t.RawSetString("shift", golua.LNumber(shift))

	return t
}

func hueBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	shift := t.RawGetString("shift").(golua.LNumber)

	f := gift.Hue(float32(shift))
	return f
}

func saturationTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct FilterSaturation
	/// @prop type {string<filter.FilterType>}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_SATURATION))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func saturationBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Saturation(float32(percent))
	return f
}

func sepiaTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct FilterSepia
	/// @prop type {string<filter.FilterType>}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_SEPIA))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func sepiaBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Sepia(float32(percent))
	return f
}

func thresholdTable(state *golua.LState, percent float64) *golua.LTable {
	/// @struct FilterThreshold
	/// @prop type {string<filter.FilterType>}
	/// @prop percent {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_THRESHOLD))
	t.RawSetString("percent", golua.LNumber(percent))

	return t
}

func thresholdBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	percent := t.RawGetString("percent").(golua.LNumber)

	f := gift.Threshold(float32(percent))
	return f
}

func pixelateTable(state *golua.LState, size int) *golua.LTable {
	/// @struct FilterPixelate
	/// @prop type {string<filter.FilterType>}
	/// @prop size {int}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_PIXELATE))
	t.RawSetString("size", golua.LNumber(size))

	return t
}

func pixelateBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	size := t.RawGetString("size").(golua.LNumber)

	f := gift.Pixelate(int(size))
	return f
}

func transposeTable(state *golua.LState) *golua.LTable {
	/// @struct FilterTranspose
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_TRANSPOSE))

	return t
}

func transposeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Transpose()
	return f
}

func transverseTable(state *golua.LState) *golua.LTable {
	/// @struct FilterTransverse
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_TRANSVERSE))

	return t
}

func transverseBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Transverse()
	return f
}

func sobelTable(state *golua.LState) *golua.LTable {
	/// @struct FilterSobel
	/// @prop type {string<filter.FilterType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_SOBEL))

	return t
}

func sobelBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	f := gift.Sobel()
	return f
}

func maximumTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct FilterMaximum
	/// @prop type {string<filter.FilterType>}
	/// @prop ksize {int}
	/// @prop disk {bool}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_MAXIMUM))
	t.RawSetString("ksize", golua.LNumber(ksize))
	t.RawSetString("disk", golua.LBool(disk))

	return t
}

func maximumBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := t.RawGetString("ksize").(golua.LNumber)
	disk := t.RawGetString("disk").(golua.LBool)

	f := gift.Maximum(int(ksize), bool(disk))
	return f
}

func meanTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct FilterMean
	/// @prop type {string<filter.FilterType>}
	/// @prop ksize {int}
	/// @prop disk {bool}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_MEAN))
	t.RawSetString("ksize", golua.LNumber(ksize))
	t.RawSetString("disk", golua.LBool(disk))

	return t
}

func meanBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := t.RawGetString("ksize").(golua.LNumber)
	disk := t.RawGetString("disk").(golua.LBool)

	f := gift.Mean(int(ksize), bool(disk))
	return f
}

func medianTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct FilterMedian
	/// @prop type {string<filter.FilterType>}
	/// @prop ksize {int}
	/// @prop disk {bool}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_MEDIAN))
	t.RawSetString("ksize", golua.LNumber(ksize))
	t.RawSetString("disk", golua.LBool(disk))

	return t
}

func medianBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := t.RawGetString("ksize").(golua.LNumber)
	disk := t.RawGetString("disk").(golua.LBool)

	f := gift.Median(int(ksize), bool(disk))
	return f
}

func minimumTable(state *golua.LState, ksize int, disk bool) *golua.LTable {
	/// @struct FilterMinimum
	/// @prop type {string<filter.FilterType>}
	/// @prop ksize {int}
	/// @prop disk {bool}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_MINIMUM))
	t.RawSetString("ksize", golua.LNumber(ksize))
	t.RawSetString("disk", golua.LBool(disk))

	return t
}

func minimumBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	ksize := t.RawGetString("ksize").(golua.LNumber)
	disk := t.RawGetString("disk").(golua.LBool)

	f := gift.Minimum(int(ksize), bool(disk))
	return f
}

func sigmoidTable(state *golua.LState, midpoint, factor float64) *golua.LTable {
	/// @struct FilterSigmoid
	/// @prop type {string<filter.FilterType>}
	/// @prop midpoint {float}
	/// @prop factor {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_SIGMOID))
	t.RawSetString("midpoint", golua.LNumber(midpoint))
	t.RawSetString("factor", golua.LNumber(factor))

	return t
}

func sigmoidBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	midpoint := t.RawGetString("midpoint").(golua.LNumber)
	factor := t.RawGetString("factor").(golua.LNumber)

	f := gift.Sigmoid(float32(midpoint), float32(factor))
	return f
}

func unsharpMaskTable(state *golua.LState, sigma, amount, threshold float64) *golua.LTable {
	/// @struct FilterUnsharpMask
	/// @prop type {string<filter.FilterType>}
	/// @prop sigma {float}
	/// @prop amount {float}
	/// @prop threshold {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_UNSHARP_MASK))
	t.RawSetString("sigma", golua.LNumber(sigma))
	t.RawSetString("amount", golua.LNumber(amount))
	t.RawSetString("threshold", golua.LNumber(threshold))

	return t
}

func unsharpMaskBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	sigma := t.RawGetString("sigma").(golua.LNumber)
	amount := t.RawGetString("amount").(golua.LNumber)
	threshold := t.RawGetString("threshold").(golua.LNumber)

	f := gift.UnsharpMask(float32(sigma), float32(amount), float32(threshold))
	return f
}

func resizeTable(state *golua.LState, width, height, resampling int) *golua.LTable {
	/// @struct FilterResize
	/// @prop type {string<filter.FilterType>}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop resampling {int<filter.Resampling>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_RESIZE))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("resampling", golua.LNumber(resampling))

	return t
}

func resizeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := t.RawGetString("width").(golua.LNumber)
	height := t.RawGetString("height").(golua.LNumber)
	resampling := t.RawGetString("resampling").(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.Resize(int(width), int(height), s)
	return f
}

func resizeToFillTable(state *golua.LState, width, height, resampling, anchor int) *golua.LTable {
	/// @struct FilterResizeToFill
	/// @prop type {string<filter.FilterType>}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop resampling {int<filter.Resampling>}
	/// @prop anchor {int<filter.Anchor>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_RESIZE_TO_FILL))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("resampling", golua.LNumber(resampling))
	t.RawSetString("anchor", golua.LNumber(anchor))

	return t
}

func resizeToFillBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := t.RawGetString("width").(golua.LNumber)
	height := t.RawGetString("height").(golua.LNumber)
	resampling := t.RawGetString("resampling").(golua.LNumber)
	anchor := t.RawGetString("anchor").(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.ResizeToFill(int(width), int(height), s, gift.Anchor(anchor))
	return f
}

func resizeToFitTable(state *golua.LState, width, height, resampling int) *golua.LTable {
	/// @struct FilterResizeToFit
	/// @prop type {string<filter.FilterType>}
	/// @prop width {int}
	/// @prop height {int}
	/// @prop resampling {int<filter.Resampling>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_RESIZE_TO_FIT))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("resampling", golua.LNumber(resampling))

	return t
}

func resizeToFitBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	width := t.RawGetString("width").(golua.LNumber)
	height := t.RawGetString("height").(golua.LNumber)
	resampling := t.RawGetString("resampling").(golua.LNumber)

	s := samplers[int(resampling)]
	f := gift.ResizeToFit(int(width), int(height), s)
	return f
}

func colorFuncTable(state *golua.LState, fn *golua.LFunction) *golua.LTable {
	/// @struct FilterColorFunc
	/// @prop type {string<filter.FilterType>}
	/// @prop fn {function(r float, g float, b float, a float) -> float, float, float, float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLOR_FUNC))
	t.RawSetString("fn", fn)

	return t
}

func colorFuncBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	fn := t.RawGetString("fn").(*golua.LFunction)

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
		cfInner.Close()
		return float32(r1), float32(g1), float32(b1), float32(a1)
	})
	return f
}

func colorFuncUnsafeTable(state *golua.LState, fn *golua.LFunction) *golua.LTable {
	/// @struct FilterColorFuncUnsafe
	/// @prop type {string<filter.FilterType>}
	/// @prop fn {function(r float, g float, b float, a float) -> float, float, float, float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(FILTER_COLOR_FUNC_UNSAFE))
	t.RawSetString("fn", fn)

	return t
}

func colorFuncUnsafeBuild(state *golua.LState, t *golua.LTable) gift.Filter {
	fn := t.RawGetString("fn").(*golua.LFunction)

	f := gift.ColorFunc(func(r0, g0, b0, a0 float32) (r float32, g float32, b float32, a float32) {
		state.Push(fn)
		state.Push(golua.LNumber(r0))
		state.Push(golua.LNumber(g0))
		state.Push(golua.LNumber(b0))
		state.Push(golua.LNumber(a0))
		state.Call(4, 4)

		r1 := state.CheckNumber(-4)
		g1 := state.CheckNumber(-3)
		b1 := state.CheckNumber(-2)
		a1 := state.CheckNumber(-1)
		state.Pop(4)
		return float32(r1), float32(g1), float32(b1), float32(a1)
	})
	return f
}
