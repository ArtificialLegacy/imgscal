package lib

import (
	"fmt"
	"image"
	"math"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ernyoke/imger/blend"
	"github.com/ernyoke/imger/blur"
	"github.com/ernyoke/imger/convolution"
	"github.com/ernyoke/imger/edgedetection"
	"github.com/ernyoke/imger/effects"
	"github.com/ernyoke/imger/generate"
	"github.com/ernyoke/imger/grayscale"
	"github.com/ernyoke/imger/histogram"
	"github.com/ernyoke/imger/padding"
	"github.com/ernyoke/imger/resize"
	"github.com/ernyoke/imger/threshold"
	"github.com/ernyoke/imger/transform"
	golua "github.com/yuin/gopher-lua"
)

const LIB_IMGER = "imger"

/// @lib Imger
/// @import imger
/// @desc
/// Library including imger image filters.

func RegisterImger(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_IMGER, r, r.State, lg)

	/// @func gray_add(img1, img2)
	/// @arg img1 {int<collection.IMAGE>} - First image.
	/// @arg img2 {int<collection.IMAGE>} - Second image.
	/// @desc
	/// Result is stored in img1.
	lib.CreateFunction(tab, "gray_add",
		[]lua.Arg{
			{Type: lua.INT, Name: "img1"},
			{Type: lua.INT, Name: "img2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var img image.Image

			r.IC.SchedulePipe(args["img2"].(int), args["img1"].(int),
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						if i.Self.Model == imageutil.MODEL_GRAY {
							img = i.Self.Image
						} else {
							img = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}
					},
				},
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_GRAY {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}

						iOut, err := blend.AddGray(imgCopy.(*image.Gray), img.(*image.Gray))
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blend images: %s", err), log.LEVEL_ERROR)), 0)
						}
						i.Self.Image = iOut
						i.Self.Model = imageutil.MODEL_GRAY
					},
				})

			return 0
		})

	/// @func gray_add_weighted(img1, img2, weight1, weight2)
	/// @arg img1 {int<collection.IMAGE>} - First image.
	/// @arg img2 {int<collection.IMAGE>} - Second image.
	/// @arg weight1 {float} - Weight of the first image.
	/// @arg weight2 {float} - Weight of the second image.
	/// @desc
	/// Result is stored in img1.
	lib.CreateFunction(tab, "gray_add_weighted",
		[]lua.Arg{
			{Type: lua.INT, Name: "img1"},
			{Type: lua.INT, Name: "img2"},
			{Type: lua.FLOAT, Name: "weight1"},
			{Type: lua.FLOAT, Name: "weight2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			weight1 := args["weight1"].(float64)
			weight2 := args["weight2"].(float64)

			var img image.Image

			r.IC.SchedulePipe(args["img2"].(int), args["img1"].(int),
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						if i.Self.Model == imageutil.MODEL_GRAY {
							img = i.Self.Image
						} else {
							img = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}
					},
				},
				&collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_GRAY {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}

						iOut, err := blend.AddGrayWeighted(imgCopy.(*image.Gray), weight1, img.(*image.Gray), weight2)
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blend images: %s", err), log.LEVEL_ERROR)), 0)
						}
						i.Self.Image = iOut
						i.Self.Model = imageutil.MODEL_GRAY
					},
				})

			return 0
		})

	/// @func gray_add_scalar(img, value)
	/// @arg img {int<collection.IMAGE>}
	/// @arg value {int}
	lib.CreateFunction(tab, "gray_add_scalar",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "value"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			value := args["value"].(int)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_GRAY {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
					}

					i.Self.Image = blend.AddScalarToGray(imgCopy.(*image.Gray), value)
					i.Self.Model = imageutil.MODEL_GRAY
				},
			})

			return 0
		})

	/// @func blur_box(img, ksize, anchor, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg ksize {struct<image.Point>}
	/// @arg anchor {struct<image.Point>} - Point inside of kernel.
	/// @arg border {int<imger.Border>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "blur_box",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "ksize"},
			{Type: lua.RAW_TABLE, Name: "anchor"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			ksize := imageutil.TableToPoint(args["ksize"].(*golua.LTable))
			anchor := imageutil.TableToPoint(args["anchor"].(*golua.LTable))
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.BoxGray(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.BoxRGBA(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func blur_box_xy(img, kwidth, kheight, ax, ay, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg kwidth {int}
	/// @arg kheight {int}
	/// @arg ax {int} - Anchor x inside of kernel.
	/// @arg ay {int} - Anchor y inside of kernel.
	/// @arg border {int<imger.Border>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "blur_box_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "kwidth"},
			{Type: lua.INT, Name: "kheight"},
			{Type: lua.INT, Name: "ax"},
			{Type: lua.INT, Name: "ay"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			ksize := image.Point{X: args["kwidth"].(int), Y: args["kheight"].(int)}
			anchor := image.Point{X: args["ax"].(int), Y: args["ay"].(int)}
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.BoxGray(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.BoxRGBA(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func blur_gaussian(img, radius, sigma, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg radius {float}
	/// @arg sigma {float}
	/// @arg border {int<imger.Border>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "blur_gaussian",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.FLOAT, Name: "sigma"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			radius := args["radius"].(float64)
			sigma := args["sigma"].(float64)
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.GaussianBlurGray(img, radius, sigma, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := blur.GaussianBlurRGBA(img, radius, sigma, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to blur image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func kernel(width, height, content) -> struct<imger.Kernel>
	/// @arg width {int} - Width of the kernel.
	/// @arg height {int} - Height of the kernel.
	/// @arg content {[][]float} - Kernel content.
	/// @returns {struct<imger.Kernel>}
	lib.CreateFunction(tab, "kernel",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "content"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := args["width"].(int)
			height := args["height"].(int)
			content := args["content"].(*golua.LTable)

			k := kernelTable(lib, state, lg, width, height, content)
			state.Push(k)
			return 1
		})

	/// @func kernel_zero(width, height) -> struct<imger.Kernel>
	/// @arg width {int} - Width of the kernel.
	/// @arg height {int} - Height of the kernel.
	/// @returns {struct<imger.Kernel>}
	lib.CreateFunction(tab, "kernel_zero",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := args["width"].(int)
			height := args["height"].(int)

			content := state.NewTable()
			for i := 0; i < width; i++ {
				col := state.NewTable()
				for j := 0; j < height; j++ {
					col.RawSetInt(j+1, golua.LNumber(0))
				}
				content.RawSetInt(i+1, col)
			}

			k := kernelTable(lib, state, lg, width, height, content)
			state.Push(k)
			return 1
		})

	/// @func kernel_flat(width, height, content_flat) -> struct<imger.Kernel>
	/// @arg width {int} - Width of the kernel.
	/// @arg height {int} - Height of the kernel.
	/// @arg content_flat {[]float} - Kernel content as a flat array. The length of the array must be width * height. Should be row-major order.
	/// @returns {struct<imger.Kernel>}
	lib.CreateFunction(tab, "kernel_flat",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "content_flat"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := args["width"].(int)
			height := args["height"].(int)
			contentFlat := args["content_flat"].(*golua.LTable)

			if contentFlat.Len() != width*height {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Invalid kernel content size: %d, expected %d", contentFlat.Len(), width*height), log.LEVEL_ERROR)), 0)
				return 0
			}

			content := state.NewTable()

			for i := 0; i < width; i++ {
				col := state.NewTable()
				for j := 0; j < height; j++ {
					col.RawSetInt(j+1, contentFlat.RawGetInt(i*height+j+1))
				}
				content.RawSetInt(i+1, col)
			}

			k := kernelTable(lib, state, lg, width, height, content)
			state.Push(k)
			return 1
		})

	/// @func convolve(img, kernel, anchor, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg kernel {struct<imger.Kernel>}
	/// @arg anchor {struct<image.Point>} - Point inside of kernel.
	/// @arg border {int<imger.Border>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "convolve",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "kernel"},
			{Type: lua.RAW_TABLE, Name: "anchor"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			kernel := kernelBuild(args["kernel"].(*golua.LTable))
			anchor := imageutil.TableToPoint(args["anchor"].(*golua.LTable))
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := convolution.ConvolveGray(img, kernel, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to convolve image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := convolution.ConvolveRGBA(img, kernel, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to convolve image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func convolve_xy(img, kernel, ax, ay, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg kernel {struct<imger.Kernel>}
	/// @arg ax {int} - Anchor x inside of kernel.
	/// @arg ay {int} - Anchor y inside of kernel.
	/// @arg border {int<imger.Border>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "convolve",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "kernel"},
			{Type: lua.INT, Name: "ax"},
			{Type: lua.INT, Name: "ay"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			kernel := kernelBuild(args["kernel"].(*golua.LTable))
			anchor := image.Point{X: args["ax"].(int), Y: args["ay"].(int)}
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := convolution.ConvolveGray(img, kernel, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to convolve image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := convolution.ConvolveRGBA(img, kernel, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to convolve image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func edge_canny(img, lower, upper, ksize, name, encoding, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg lower {float} - Lower threshold.
	/// @arg upper {float} - Upper threshold.
	/// @arg ksize {int} - Size of the kernel.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.GRAY color model.
	lib.CreateFunction(tab, "edge_canny",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "lower"},
			{Type: lua.FLOAT, Name: "upper"},
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			lower := args["lower"].(float64)
			upper := args["upper"].(float64)
			ksize := uint(args["ksize"].(int))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut, err := edgedetection.CannyGray(img, lower, upper, ksize)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut, err := edgedetection.CannyRGBA(img, lower, upper, ksize)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func edge_canny_inplace(img, lower, upper, ksize, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg lower {float} - Lower threshold.
	/// @arg upper {float} - Upper threshold.
	/// @arg ksize {int} - Size of the kernel.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "edge_canny_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "lower"},
			{Type: lua.FLOAT, Name: "upper"},
			{Type: lua.INT, Name: "ksize"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			lower := args["lower"].(float64)
			upper := args["upper"].(float64)
			ksize := uint(args["ksize"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.CannyGray(img, lower, upper, ksize)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.CannyRGBA(img, lower, upper, ksize)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			})

			return 0
		})

	/// @func edge_sobel(img, name, encoding, border, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.GRAY color model.
	lib.CreateFunction(tab, "edge_sobel",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut, err := edgedetection.SobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut, err := edgedetection.SobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func edge_sobel_inplace(img, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "edge_sobel_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.SobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.SobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			})

			return 0
		})

	/// @func edge_sobel_horizontal(img, name, encoding, border, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.GRAY color model.
	lib.CreateFunction(tab, "edge_sobel_horizontal",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut, err := edgedetection.HorizontalSobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut, err := edgedetection.HorizontalSobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func edge_sobel_horizontal_inplace(img, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "edge_sobel_horizontal_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.HorizontalSobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.HorizontalSobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			})

			return 0
		})

	/// @func edge_sobel_vertical(img, name, encoding, border, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.GRAY color model.
	lib.CreateFunction(tab, "edge_sobel_vertical",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut, err := edgedetection.VerticalSobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut, err := edgedetection.VerticalSobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func edge_sobel_vertical_inplace(img, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "edge_sobel_vertical_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.VerticalSobelGray(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.VerticalSobelRGBA(img, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			})

			return 0
		})

	/// @func edge_laplacian(img, kernel, name, encoding, border, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg kernel {int<imger.LaplacianKernel>} - Laplacian kernel.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.GRAY color model.
	lib.CreateFunction(tab, "edge_laplacian",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "kernel"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)
			border := padding.Border(args["border"].(int))
			kernel := edgedetection.LaplacianKernel(args["kernel"].(int))

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut, err := edgedetection.LaplacianGray(img, border, kernel)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut, err := edgedetection.LaplacianRGBA(img, border, kernel)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func edge_laplacian_inplace(img, kernel, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg kernel {int<imger.LaplacianKernel>} - Laplacian kernel.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "edge_laplacian_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "kernel"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			kernel := edgedetection.LaplacianKernel(args["kernel"].(int))
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.LaplacianGray(img, border, kernel)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := edgedetection.LaplacianRGBA(img, border, kernel)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to detect edges: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			})

			return 0
		})

	/// @func emboss(img, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "emboss",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			gray := args["gray"].(bool)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					var imgOut *image.Gray
					var err error

					if gray {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_GRAY {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}

						imgOut, err = effects.EmbossGray(imgCopy.(*image.Gray))
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to emboss image: %s", err), log.LEVEL_ERROR)), 0)
						}
					} else {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_RGBA {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
						}

						imgOut, err = effects.EmbossRGBA(imgCopy.(*image.RGBA))
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to emboss image: %s", err), log.LEVEL_ERROR)), 0)
						}
					}

					i.Self.Image = imgOut
					i.Self.Model = imageutil.MODEL_GRAY
				},
			})

			return 0
		})

	/// @func invert(img, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "invert",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				return effects.InvertGray(img), imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				return effects.InvertRGBA(img), imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func pixelate(img, factor, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg factor {float} - Pixelation factor.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "pixelate",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "factor"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			factor := args["factor"].(float64)
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := effects.PixelateGray(img, factor)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to pixelate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := effects.PixelateRGBA(img, factor)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to pixelate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func sharpen(img, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "sharpen",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := effects.SharpenGray(img)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to sharpen image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := effects.SharpenRGBA(img)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to sharpen image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func sepia(img)
	/// @arg img {int<collection.IMAGE>}
	lib.CreateFunction(tab, "sepia",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_RGBA {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
					}

					i.Self.Image = effects.Sepia(imgCopy.(*image.RGBA))
					i.Self.Model = imageutil.MODEL_RGBA
				},
			})

			return 0
		})

	/// @func gradient_linear(size, startColor, endColor, direction, name, encoding) -> int<collection.IMAGE>
	/// @arg size {struct<image.Point>} - Size of the gradient.
	/// @arg startColor {struct<image.Color>} - Start color of the gradient.
	/// @arg endColor {struct<image.Color>} - End color of the gradient.
	/// @arg direction {int<imger.Direction>} - Direction of the gradient.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.RGBA color model.
	lib.CreateFunction(tab, "gradient_linear",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "size"},
			{Type: lua.RAW_TABLE, Name: "startColor"},
			{Type: lua.RAW_TABLE, Name: "endColor"},
			{Type: lua.INT, Name: "direction"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			size := imageutil.TableToPoint(args["size"].(*golua.LTable))
			startColor := imageutil.ColorTableToRGBAColor(args["startColor"].(*golua.LTable))
			endColor := imageutil.ColorTableToRGBAColor(args["endColor"].(*golua.LTable))
			direction := generate.Direction(args["direction"].(int))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self = &collection.ItemImage{
						Image:    generate.LinearGradient(size, *startColor, *endColor, direction),
						Encoding: encoding,
						Name:     name,
						Model:    imageutil.MODEL_RGBA,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func gradient_linear_xy(width, height, startColor, endColor, direction, name, encoding) -> int<collection.IMAGE>
	/// @arg width {int} - Width of the gradient.
	/// @arg height {int} - Height of the gradient.
	/// @arg startColor {struct<image.Color>} - Start color of the gradient.
	/// @arg endColor {struct<image.Color>} - End color of the gradient.
	/// @arg direction {int<imger.Direction>} - Direction of the gradient.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.RGBA color model.
	lib.CreateFunction(tab, "gradient_linear_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "startColor"},
			{Type: lua.RAW_TABLE, Name: "endColor"},
			{Type: lua.INT, Name: "direction"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			size := image.Point{X: args["width"].(int), Y: args["height"].(int)}
			startColor := imageutil.ColorTableToRGBAColor(args["startColor"].(*golua.LTable))
			endColor := imageutil.ColorTableToRGBAColor(args["endColor"].(*golua.LTable))
			direction := generate.Direction(args["direction"].(int))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self = &collection.ItemImage{
						Image:    generate.LinearGradient(size, *startColor, *endColor, direction),
						Encoding: encoding,
						Name:     name,
						Model:    imageutil.MODEL_RGBA,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func gradient_sigmoidal(size, startColor, endColor, direction, name, encoding) -> int<collection.IMAGE>
	/// @arg size {struct<image.Point>} - Size of the gradient.
	/// @arg startColor {struct<image.Color>} - Start color of the gradient.
	/// @arg endColor {struct<image.Color>} - End color of the gradient.
	/// @arg direction {int<imger.Direction>} - Direction of the gradient.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.RGBA color model.
	lib.CreateFunction(tab, "gradient_sigmoidal",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "size"},
			{Type: lua.RAW_TABLE, Name: "startColor"},
			{Type: lua.RAW_TABLE, Name: "endColor"},
			{Type: lua.INT, Name: "direction"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			size := imageutil.TableToPoint(args["size"].(*golua.LTable))
			startColor := imageutil.ColorTableToRGBAColor(args["startColor"].(*golua.LTable))
			endColor := imageutil.ColorTableToRGBAColor(args["endColor"].(*golua.LTable))
			direction := generate.Direction(args["direction"].(int))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self = &collection.ItemImage{
						Image:    generate.SigmoidalGradient(size, *startColor, *endColor, direction),
						Encoding: encoding,
						Name:     name,
						Model:    imageutil.MODEL_RGBA,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func gradient_sigmoidal_xy(width, height, startColor, endColor, direction, name, encoding) -> int<collection.IMAGE>
	/// @arg width {int} - Width of the gradient.
	/// @arg height {int} - Height of the gradient.
	/// @arg startColor {struct<image.Color>} - Start color of the gradient.
	/// @arg endColor {struct<image.Color>} - End color of the gradient.
	/// @arg direction {int<imger.Direction>} - Direction of the gradient.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @returns {int<collection.IMAGE>} - The resulting image will be in the image.RGBA color model.
	lib.CreateFunction(tab, "gradient_sigmoidal_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.RAW_TABLE, Name: "startColor"},
			{Type: lua.RAW_TABLE, Name: "endColor"},
			{Type: lua.INT, Name: "direction"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			size := image.Point{X: args["width"].(int), Y: args["height"].(int)}
			startColor := imageutil.ColorTableToRGBAColor(args["startColor"].(*golua.LTable))
			endColor := imageutil.ColorTableToRGBAColor(args["endColor"].(*golua.LTable))
			direction := generate.Direction(args["direction"].(int))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					i.Self = &collection.ItemImage{
						Image:    generate.SigmoidalGradient(size, *startColor, *endColor, direction),
						Encoding: encoding,
						Name:     name,
						Model:    imageutil.MODEL_RGBA,
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func grayscale(img, name, encoding, to16?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg? to16 {bool} - If true, the image is converted to 16-bit grayscale, otherwise it is converted to 8-bit grayscale.
	/// @returns {int<collection.IMAGE>}
	lib.CreateFunction(tab, "grayscale",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.BOOL, Name: "to16", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			to16 := args["to16"].(bool)

			var imgOut image.Image
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgOut = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
				},
			})

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady

					var imgSave image.Image
					var model imageutil.ColorModel

					if to16 {
						imgSave = grayscale.Grayscale16(imgOut)
						model = imageutil.MODEL_GRAY16
					} else {
						imgSave = grayscale.Grayscale(imgOut)
						model = imageutil.MODEL_GRAY
					}

					i.Self = &collection.ItemImage{
						Image:    imgSave,
						Encoding: encoding,
						Name:     name,
						Model:    model,
					}

					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgFinished <- struct{}{}
				},
			})

			return 0
		})

	/// @func grayscale_inplace(img, to16?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg? to16 {bool} - If true, the image is converted to 16-bit grayscale, otherwise it is converted to 8-bit grayscale.
	lib.CreateFunction(tab, "grayscale_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.BOOL, Name: "to16", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			to16 := args["to16"].(bool)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					var imgSave image.Image
					var model imageutil.ColorModel

					if to16 {
						imgSave = grayscale.Grayscale16(i.Self.Image)
						model = imageutil.MODEL_GRAY16
					} else {
						imgSave = grayscale.Grayscale(i.Self.Image)
						model = imageutil.MODEL_GRAY
					}

					i.Self.Image = imgSave
					i.Self.Model = model
				},
			})

			return 0
		})

	/// @func padding(img, ksize, anchor, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg ksize {struct<image.Point>} - Size of the kernel.
	/// @arg anchor {struct<image.Point>} - Anchor point of the kernel.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "padding",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "ksize"},
			{Type: lua.RAW_TABLE, Name: "anchor"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			ksize := imageutil.TableToPoint(args["ksize"].(*golua.LTable))
			anchor := imageutil.TableToPoint(args["anchor"].(*golua.LTable))
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingGray(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingRGBA(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func padding_xy(img, kwidth, kheight, ax, ay, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg kwidth {int} - Width of the kernel.
	/// @arg kheight {int} - Height of the kernel.
	/// @arg ax {int} - X coordinate of the anchor point.
	/// @arg ay {int} - Y coordinate of the anchor point.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "padding_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "kwidth"},
			{Type: lua.INT, Name: "kheight"},
			{Type: lua.INT, Name: "ax"},
			{Type: lua.INT, Name: "ay"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			ksize := image.Point{X: args["kwidth"].(int), Y: args["kheight"].(int)}
			anchor := image.Point{X: args["ax"].(int), Y: args["ay"].(int)}
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingGray(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingRGBA(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func padding_size(img, top, bottom, left, right, border, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg top {int} - Top padding.
	/// @arg bottom {int} - Bottom padding.
	/// @arg left {int} - Left padding.
	/// @arg right {int} - Right padding.
	/// @arg border {int<imger.Border>} - Border type.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "padding_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "top"},
			{Type: lua.INT, Name: "bottom"},
			{Type: lua.INT, Name: "left"},
			{Type: lua.INT, Name: "right"},
			{Type: lua.INT, Name: "border"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			top := args["top"].(int)
			bottom := args["bottom"].(int)
			left := args["left"].(int)
			right := args["right"].(int)
			ksize := image.Point{X: left + right + 1, Y: top + bottom + 1}
			anchor := image.Point{X: left, Y: top}
			border := padding.Border(args["border"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingGray(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := padding.PaddingRGBA(img, ksize, anchor, border)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to add padding to image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func threshold(img, t, method, to16?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg t {int} - Threshold value, either an uint8 or uint16 value.
	/// @arg method {int<imger.Threshold>} - Threshold method.
	/// @arg? to16 {bool} - If true, the image is converted to 16-bit grayscale, otherwise it is converted to 8-bit grayscale.
	lib.CreateFunction(tab, "threshold",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "t"},
			{Type: lua.INT, Name: "method"},
			{Type: lua.BOOL, Name: "to16", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			t := args["t"].(int)
			method := threshold.Method(args["method"].(int))
			to16 := args["to16"].(bool)

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					var imgSave image.Image
					var model imageutil.ColorModel
					var err error

					if !to16 {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_GRAY {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
						}

						imgSave, err = threshold.Threshold(imgCopy.(*image.Gray), uint8(t), method)
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to apply threshold to image: %s", err), log.LEVEL_ERROR)), 0)
						}
						model = imageutil.MODEL_GRAY
					} else {
						imgCopy := i.Self.Image
						if i.Self.Model != imageutil.MODEL_GRAY16 {
							imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY16)
						}

						imgSave, err = threshold.Threshold16(imgCopy.(*image.Gray16), uint16(t), method)
						if err != nil {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to apply threshold to image: %s", err), log.LEVEL_ERROR)), 0)
						}
						model = imageutil.MODEL_GRAY16
					}

					i.Self.Image = imgSave
					i.Self.Model = model
				},
			})

			return 0
		})

	/// @func threshold_otsu(img, method)
	/// @arg img {int<collection.IMAGE>}
	/// @arg method {int<imger.Threshold>} - Threshold method.
	lib.CreateFunction(tab, "threshold_otsu",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "method"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			method := threshold.Method(args["method"].(int))

			r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_GRAY {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
					}

					imgSave, err := threshold.OtsuThreshold(imgCopy.(*image.Gray), method)
					if err != nil {
						state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to apply otsu threshold to image: %s", err), log.LEVEL_ERROR)), 0)
					}

					i.Self.Image = imgSave
					i.Self.Model = imageutil.MODEL_GRAY
				},
			})

			return 0
		})

	/// @func rotate(img, angle, anchor, resizeToFit, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg angle {float} - Angle of rotation in degrees, counter-clockwise.
	/// @arg anchor {struct<image.Point>} - Anchor point of rotation.
	/// @arg resizeToFit {bool} - If true, the image is resized to fit the rotated image.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "rotate",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.RAW_TABLE, Name: "anchor"},
			{Type: lua.BOOL, Name: "resizeToFit"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			angle := args["angle"].(float64)
			anchor := imageutil.TableToPoint(args["anchor"].(*golua.LTable))
			resizeToFit := args["resizeToFit"].(bool)
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := transform.RotateGray(img, angle, anchor, resizeToFit)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to rotate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := transform.RotateRGBA(img, angle, anchor, resizeToFit)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to rotate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func rotate_xy(img, angle, ax, ay, resizeToFit, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg angle {float} - Angle of rotation in degrees, counter-clockwise.
	/// @arg ax {int} - X coordinate of the anchor point.
	/// @arg ay {int} - Y coordinate of the anchor point.
	/// @arg resizeToFit {bool} - If true, the image is resized to fit the rotated image.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "rotate_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "angle"},
			{Type: lua.INT, Name: "ax"},
			{Type: lua.INT, Name: "ay"},
			{Type: lua.BOOL, Name: "resizeToFit"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			angle := args["angle"].(float64)
			anchor := image.Point{X: args["ax"].(int), Y: args["ay"].(int)}
			resizeToFit := args["resizeToFit"].(bool)
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := transform.RotateGray(img, angle, anchor, resizeToFit)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to rotate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := transform.RotateRGBA(img, angle, anchor, resizeToFit)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to rotate image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func resize(img, scalex, scaley, inter, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg scalex {float} - Scale factor along the horizontal axis.
	/// @arg scaley {float} - Scale factor along the vertical axis.
	/// @arg inter {int<imger.Interpolation>} - Interpolation method.
	/// @arg? gray {bool} - If true, the image is converted to gray before applying the filter, otherwise it is converted to RGBA.
	lib.CreateFunction(tab, "resize",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.FLOAT, Name: "scalex"},
			{Type: lua.FLOAT, Name: "scaley"},
			{Type: lua.INT, Name: "inter"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			scalex := args["scalex"].(float64)
			scaley := args["scaley"].(float64)
			inter := resize.Interpolation(args["inter"].(int))
			gray := args["gray"].(bool)

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut, err := resize.ResizeGray(img, scalex, scaley, inter)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to resize image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut, err := resize.ResizeRGBA(img, scalex, scaley, inter)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("Failed to resize image: %s", err), log.LEVEL_ERROR)), 0)
				}
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func histogram_draw(img, scale, name, encoding, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg scale {struct<image.Point>} - Scale of the histogram.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg? gray {bool}
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Size is (256*scale.x, 256*scale.y).
	lib.CreateFunction(tab, "histogram_draw",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "scale"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			scale := imageutil.TableToPoint(args["scale"].(*golua.LTable))
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)

			scale.X *= 256
			scale.Y *= 256

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut := histogram.DrawHistogramGray(img, scale)
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut := histogram.DrawHistogramRGBA(img, scale)
				imageutil.AlphaSet(iOut, 255)
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func histogram_draw_xy(img, scalex, scaley, name, encoding, gray?) -> int<collection.IMAGE>
	/// @arg img {int<collection.IMAGE>}
	/// @arg scalex {int} - Scale of the histogram along the horizontal axis.
	/// @arg scaley {int} - Scale of the histogram along the vertical axis.
	/// @arg name {string} - Name of the image.
	/// @arg encoding {int<image.Encoding>} - Encoding of the image.
	/// @arg? gray {bool}
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Size is (256*scalex, 256*scaley).
	lib.CreateFunction(tab, "histogram_draw_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "scalex"},
			{Type: lua.INT, Name: "scaley"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			scale := image.Point{X: args["scalex"].(int), Y: args["scaley"].(int)}
			name := args["name"].(string)
			encoding := lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib)
			gray := args["gray"].(bool)

			scale.X *= 256
			scale.Y *= 256

			id := imgerFilterNew(r, lg, d, img, name, encoding, gray, func(img *image.Gray) image.Image {
				iOut := histogram.DrawHistogramGray(img, scale)
				return iOut
			}, func(img *image.RGBA) image.Image {
				iOut := histogram.DrawHistogramRGBA(img, scale)
				imageutil.AlphaSet(iOut, 255)
				return iOut
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func histogram_draw_inplace(img, scale, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg scale {struct<image.Point>} - Scale of the histogram.
	/// @arg? gray {bool}
	/// @desc
	/// Size is (256*scale.x, 256*scale.y).
	lib.CreateFunction(tab, "histogram_draw_inplace",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.RAW_TABLE, Name: "scale"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			scale := imageutil.TableToPoint(args["scale"].(*golua.LTable))
			gray := args["gray"].(bool)

			scale.X *= 256
			scale.Y *= 256

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut := histogram.DrawHistogramGray(img, scale)
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut := histogram.DrawHistogramRGBA(img, scale)
				imageutil.AlphaSet(iOut, 255)
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func histogram_draw_inplace_xy(img, scalex, scaley, gray?)
	/// @arg img {int<collection.IMAGE>}
	/// @arg scalex {int} - Scale of the histogram along the horizontal axis.
	/// @arg scaley {int} - Scale of the histogram along the vertical axis.
	/// @arg? gray {bool}
	/// @desc
	/// Size is (256*scalex, 256*scaley).
	lib.CreateFunction(tab, "histogram_draw_inplace_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "scalex"},
			{Type: lua.INT, Name: "scaley"},
			{Type: lua.BOOL, Name: "gray", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)
			scale := image.Point{X: args["scalex"].(int), Y: args["scaley"].(int)}
			gray := args["gray"].(bool)

			scale.X *= 256
			scale.Y *= 256

			imgerFilter(r, d, img, gray, func(img *image.Gray) (image.Image, imageutil.ColorModel) {
				iOut := histogram.DrawHistogramGray(img, scale)
				return iOut, imageutil.MODEL_GRAY
			}, func(img *image.RGBA) (image.Image, imageutil.ColorModel) {
				iOut := histogram.DrawHistogramRGBA(img, scale)
				imageutil.AlphaSet(iOut, 255)
				return iOut, imageutil.MODEL_RGBA
			})

			return 0
		})

	/// @func histogram_gray(img, max, removeZero?) -> []int
	/// @arg img {int<collection.IMAGE>}
	/// @arg max {int} - Maximum value of the histogram.
	/// @arg? removeZero {bool} - If true, the zero value is removed from the histogram.
	/// @returns {[]int} - Length is 256.
	/// @blocking
	lib.CreateFunction(tab, "histogram_gray",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "max"},
			{Type: lua.BOOL, Name: "removeZero", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			var data [256]uint64

			<-r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_GRAY {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
					}

					data = histogram.HistogramGray(imgCopy.(*image.Gray))
				},
			})

			lst := state.NewTable()
			max := uint64(0)
			removeZero := args["removeZero"].(bool)
			maxValue := uint64(args["max"].(int))

			for i, v := range data {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}

			for i, v := range data {
				if removeZero && i == 0 {
					lst.Append(golua.LNumber(0))
					continue
				}
				lst.Append(golua.LNumber(v * maxValue / max))
			}

			state.Push(lst)
			return 1
		})

	/// @func histogram_rgb(img, max, removeZero?) -> []int, []int, []int
	/// @arg img {int<collection.IMAGE>}
	/// @arg max {int} - Maximum value of the histogram.
	/// @arg? removeZero {bool} - If true, the zero value is removed from the histogram.
	/// @returns {[]int} - Histogram of the red channel, length is 256.
	/// @returns {[]int} - Histogram of the green channel, length is 256.
	/// @returns {[]int} - Histogram of the blue channel, length is 256.
	/// @blocking
	lib.CreateFunction(tab, "histogram_rgb",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "max"},
			{Type: lua.BOOL, Name: "removeZero", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			var dataRed [256]uint64
			var dataGreen [256]uint64
			var dataBlue [256]uint64

			<-r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_RGBA {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
					}

					data := histogram.HistogramRGBA(imgCopy.(*image.RGBA))
					dataRed = data[0]
					dataGreen = data[1]
					dataBlue = data[2]
				},
			})

			lstRed := state.NewTable()
			lstGreen := state.NewTable()
			lstBlue := state.NewTable()
			max := uint64(0)
			removeZero := args["removeZero"].(bool)
			maxValue := uint64(args["max"].(int))

			for i, v := range dataRed {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}
			for i, v := range dataGreen {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}
			for i, v := range dataBlue {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}

			for i, v := range dataRed {
				if removeZero && i == 0 {
					lstRed.Append(golua.LNumber(0))
					continue
				}
				lstRed.Append(golua.LNumber(v * maxValue / max))
			}
			for i, v := range dataGreen {
				if removeZero && i == 0 {
					lstGreen.Append(golua.LNumber(0))
					continue
				}
				lstGreen.Append(golua.LNumber(v * maxValue / max))
			}
			for i, v := range dataBlue {
				if removeZero && i == 0 {
					lstBlue.Append(golua.LNumber(0))
					continue
				}
				lstBlue.Append(golua.LNumber(v * maxValue / max))
			}

			state.Push(lstRed)
			state.Push(lstGreen)
			state.Push(lstBlue)
			return 3
		})

	/// @func histogram_red(img, max, removeZero?) -> []int
	/// @arg img {int<collection.IMAGE>}
	/// @arg max {int} - Maximum value of the histogram.
	/// @arg? removeZero {bool} - If true, the zero value is removed from the histogram.
	/// @returns {[]int} - Length is 256.
	/// @blocking
	lib.CreateFunction(tab, "histogram_red",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "max"},
			{Type: lua.BOOL, Name: "removeZero", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			var data [256]uint64

			<-r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_RGBA {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
					}

					data = histogram.HistogramRGBARed(imgCopy.(*image.RGBA))
				},
			})

			lst := state.NewTable()
			max := uint64(0)
			removeZero := args["removeZero"].(bool)
			maxValue := uint64(args["max"].(int))

			for i, v := range data {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}

			for i, v := range data {
				if removeZero && i == 0 {
					lst.Append(golua.LNumber(0))
					continue
				}
				lst.Append(golua.LNumber(v * maxValue / max))
			}

			state.Push(lst)
			return 1
		})

	/// @func histogram_green(img, max, removeZero?) -> []int
	/// @arg img {int<collection.IMAGE>}
	/// @arg max {int} - Maximum value of the histogram.
	/// @arg? removeZero {bool} - If true, the zero value is removed from the histogram.
	/// @returns {[]int} - Length is 256.
	/// @blocking
	lib.CreateFunction(tab, "histogram_green",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "max"},
			{Type: lua.BOOL, Name: "removeZero", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			var data [256]uint64

			<-r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_RGBA {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
					}

					data = histogram.HistogramRGBAGreen(imgCopy.(*image.RGBA))
				},
			})

			lst := state.NewTable()
			max := uint64(0)
			removeZero := args["removeZero"].(bool)
			maxValue := uint64(args["max"].(int))

			for i, v := range data {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}

			for i, v := range data {
				if removeZero && i == 0 {
					lst.Append(golua.LNumber(0))
					continue
				}
				lst.Append(golua.LNumber(v * maxValue / max))
			}

			state.Push(lst)
			return 1
		})

	/// @func histogram_blue(img, max, removeZero?) -> []int
	/// @arg img {int<collection.IMAGE>}
	/// @arg max {int} - Maximum value of the histogram.
	/// @arg? removeZero {bool} - If true, the zero value is removed from the histogram.
	/// @returns {[]int} - Length is 256.
	/// @blocking
	lib.CreateFunction(tab, "histogram_blue",
		[]lua.Arg{
			{Type: lua.INT, Name: "img"},
			{Type: lua.INT, Name: "max"},
			{Type: lua.BOOL, Name: "removeZero", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img := args["img"].(int)

			var data [256]uint64

			<-r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					imgCopy := i.Self.Image
					if i.Self.Model != imageutil.MODEL_RGBA {
						imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
					}

					data = histogram.HistogramRGBABlue(imgCopy.(*image.RGBA))
				},
			})

			lst := state.NewTable()
			max := uint64(0)
			removeZero := args["removeZero"].(bool)
			maxValue := uint64(args["max"].(int))

			for i, v := range data {
				if removeZero && i == 0 {
					continue
				}
				if v > max {
					max = v
				}
			}

			for i, v := range data {
				if removeZero && i == 0 {
					lst.Append(golua.LNumber(0))
					continue
				}
				lst.Append(golua.LNumber(v * maxValue / max))
			}

			state.Push(lst)
			return 1
		})

	/// @constants Border {int}
	/// @const BORDER_CONSTANT
	/// @const BORDER_REPLICATE
	/// @const BORDER_REFLECT
	tab.RawSetString("BORDER_CONSTANT", golua.LNumber(padding.BorderConstant))
	tab.RawSetString("BORDER_REPLICATE", golua.LNumber(padding.BorderReplicate))
	tab.RawSetString("BORDER_REFLECT", golua.LNumber(padding.BorderReflect))

	/// @constants LaplacianKernel {int}
	/// @const LAPLACIAN_K4
	/// @const LAPLACIAN_K8
	tab.RawSetString("LAPLACIAN_K4", golua.LNumber(edgedetection.K4))
	tab.RawSetString("LAPLACIAN_K8", golua.LNumber(edgedetection.K8))

	/// @constants Direction {int}
	/// @const DIR_HORIZONTAL
	/// @const DIR_VERTICAL
	tab.RawSetString("DIR_HORIZONTAL", golua.LNumber(generate.H))
	tab.RawSetString("DIR_VERTICAL", golua.LNumber(generate.V))

	/// @constants Threshold {int}
	/// @const THRESHOLD_BINARY
	/// @const THRESHOLD_BINARYINV
	/// @const THRESHOLD_TRUNC
	/// @const THRESHOLD_TOZERO
	/// @const THRESHOLD_TOZEROINV
	tab.RawSetString("THRESHOLD_BINARY", golua.LNumber(threshold.ThreshBinary))
	tab.RawSetString("THRESHOLD_BINARYINV", golua.LNumber(threshold.ThreshBinaryInv))
	tab.RawSetString("THRESHOLD_TRUNC", golua.LNumber(threshold.ThreshTrunc))
	tab.RawSetString("THRESHOLD_TOZERO", golua.LNumber(threshold.ThreshToZero))
	tab.RawSetString("THRESHOLD_TOZEROINV", golua.LNumber(threshold.ThreshToZeroInv))

	/// @constants Interpolation {int}
	/// @const INTER_NEAREST
	/// @const INTER_LINEAR
	/// @const INTER_CATMULLROM
	/// @const INTER_LANCZOS
	tab.RawSetString("INTER_NEAREST", golua.LNumber(resize.InterNearest))
	tab.RawSetString("INTER_LINEAR", golua.LNumber(resize.InterLinear))
	tab.RawSetString("INTER_CATMULLROM", golua.LNumber(resize.InterCatmullRom))
	tab.RawSetString("INTER_LANCZOS", golua.LNumber(resize.InterLanczos))
}

func imgerFilter(r *lua.Runner, d lua.TaskData, img int, gray bool, fnGRAYA func(*image.Gray) (image.Image, imageutil.ColorModel), fnRGBA func(*image.RGBA) (image.Image, imageutil.ColorModel)) {
	r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
		Lib:  d.Lib,
		Name: d.Name,
		Fn: func(i *collection.Item[collection.ItemImage]) {
			if gray {
				imgCopy := i.Self.Image
				if i.Self.Model != imageutil.MODEL_GRAY {
					imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
				}
				i.Self.Image, i.Self.Model = fnGRAYA(imgCopy.(*image.Gray))
			} else {
				imgCopy := i.Self.Image
				if i.Self.Model != imageutil.MODEL_RGBA {
					imgCopy = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
				}
				i.Self.Image, i.Self.Model = fnRGBA(imgCopy.(*image.RGBA))
			}
		},
	})
}

func imgerFilterNew(r *lua.Runner, lg *log.Logger, d lua.TaskData, img int, name string, encoding imageutil.ImageEncoding, gray bool, fnGRAY func(*image.Gray) image.Image, fnRGBA func(*image.RGBA) image.Image) int {
	var imgOut image.Image
	imgReady := make(chan struct{}, 2)

	r.IC.Schedule(img, &collection.Task[collection.ItemImage]{
		Lib:  d.Lib,
		Name: d.Name,
		Fn: func(i *collection.Item[collection.ItemImage]) {
			if gray {
				imgOut = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_GRAY)
			} else {
				imgOut = imageutil.CopyImage(i.Self.Image, imageutil.MODEL_RGBA)
			}
			imgReady <- struct{}{}
		},
		Fail: func(i *collection.Item[collection.ItemImage]) {
			imgReady <- struct{}{}
		},
	})

	chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
	lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

	id := r.IC.AddItem(&chLog)

	r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
		Lib:  d.Lib,
		Name: d.Name,
		Fn: func(i *collection.Item[collection.ItemImage]) {
			var model imageutil.ColorModel

			<-imgReady

			var imgSave image.Image

			if gray {
				model = imageutil.MODEL_GRAY
				imgSave = fnGRAY(imgOut.(*image.Gray))
			} else {
				model = imageutil.MODEL_RGBA
				imgSave = fnRGBA(imgOut.(*image.RGBA))
			}

			i.Self = &collection.ItemImage{
				Image:    imgSave,
				Encoding: encoding,
				Name:     name,
				Model:    model,
			}
		},
	})

	return id
}

func kernelTable(lib *lua.Lib, state *golua.LState, lg *log.Logger, width, height int, content *golua.LTable) *golua.LTable {
	/// @struct Kernel
	/// @prop width {int} - Width of the kernel.
	/// @prop height {int} - Height of the kernel.
	/// @prop content {[][]float} - Kernel content.
	/// @method center() -> struct<image.Point> - Returns the center of the kernel, rounded down.
	/// @method center_xy() -> int, int - Returns the center of the kernel, rounded down.
	/// @method sum() -> float - Returns the sum of all elements in the kernel.
	/// @method sum_abs() -> float - Returns the sum of the absolute values of all elements in the kernel.
	/// @method size() -> struct<image.Point> - Returns the size of the kernel as a point.
	/// @method at(x int, y int) -> float - Returns the value at the given position, or nil if out of bounds. Uses 0-based indexing, set directly to the kernel content for 1-based indexing.
	/// @method set(self, x int, y int, value float) -> self - Sets the value at the given position. Uses 0-based indexing, set directly to the kernel content for 1-based indexing.
	/// @method normalize(self) -> self - Normalizes the kernel.
	/// @method normalize_set(self, normalize bool) -> self - Sets the normalize flag of the kernel, this normalization is only run when the kernel is built.

	// verify content size
	if content.Len() != width {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("Invalid kernel content size: %d, expected %d", content.Len(), width), log.LEVEL_ERROR)), 0)
		return nil
	}

	for i := 0; i < width; i++ {
		col := content.RawGetInt(i + 1).(*golua.LTable)
		if col.Len() != height {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("Invalid kernel row size: %d, expected %d", col.Len(), height), log.LEVEL_ERROR)), 0)
			return nil
		}
	}

	t := state.NewTable()

	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("content", content)

	t.RawSetString("__normalize", golua.LFalse)

	t.RawSetString("center", state.NewFunction(func(state *golua.LState) int {
		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)

		centerx := int(width / 2)
		centery := int(height / 2)

		p := state.NewTable()
		p.RawSetString("x", golua.LNumber(centerx))
		p.RawSetString("y", golua.LNumber(centery))

		state.Push(p)
		return 1
	}))

	t.RawSetString("center_xy", state.NewFunction(func(state *golua.LState) int {
		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)

		centerx := int(width / 2)
		centery := int(height / 2)

		state.Push(golua.LNumber(centerx))
		state.Push(golua.LNumber(centery))
		return 2
	}))

	t.RawSetString("sum", state.NewFunction(func(state *golua.LState) int {
		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)
		content := t.RawGetString("content").(*golua.LTable)

		sum := 0.0

		for i := 0; i < int(width); i++ {
			col := content.RawGetInt(i + 1).(*golua.LTable)
			for j := 0; j < int(height); j++ {
				sum += float64(col.RawGetInt(j + 1).(golua.LNumber))
			}
		}

		state.Push(golua.LNumber(sum))
		return 1
	}))

	t.RawSetString("sum_abs", state.NewFunction(func(state *golua.LState) int {
		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)
		content := t.RawGetString("content").(*golua.LTable)

		sum := 0.0

		for i := 0; i < int(width); i++ {
			col := content.RawGetInt(i + 1).(*golua.LTable)
			for j := 0; j < int(height); j++ {
				sum += math.Abs(float64(col.RawGetInt(j + 1).(golua.LNumber)))
			}
		}

		state.Push(golua.LNumber(sum))
		return 1
	}))

	t.RawSetString("size", state.NewFunction(func(state *golua.LState) int {
		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)

		p := state.NewTable()
		p.RawSetString("x", width)
		p.RawSetString("y", height)

		state.Push(p)
		return 1
	}))

	t.RawSetString("at", state.NewFunction(func(state *golua.LState) int {
		x := state.CheckInt(2)
		y := state.CheckInt(3)

		width := t.RawGetString("width").(golua.LNumber)
		height := t.RawGetString("height").(golua.LNumber)
		content := t.RawGetString("content").(*golua.LTable)

		if x < 0 || x >= int(width) || y < 0 || y >= int(height) {
			state.Push(golua.LNil)
		} else {
			col := content.RawGetInt(x + 1).(*golua.LTable)
			state.Push(col.RawGetInt(y + 1))
		}

		return 1
	}))

	lib.BuilderFunction(state, t, "set",
		[]lua.Arg{
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
			{Type: lua.FLOAT, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			width := t.RawGetString("width").(golua.LNumber)
			height := t.RawGetString("height").(golua.LNumber)
			content := t.RawGetString("content").(*golua.LTable)

			x := args["x"].(int)
			y := args["y"].(int)
			value := args["value"].(float64)

			if x < 0 || x >= int(width) || y < 0 || y >= int(height) {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("Index out of bounds: (%d,%d) of (%d,%d)", x, y, width, height), log.LEVEL_ERROR)), 0)
				return
			}

			col := content.RawGetInt(x + 1).(*golua.LTable)
			col.RawSetInt(y+1, golua.LNumber(value))
		})

	lib.BuilderFunction(state, t, "normalize",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			width := t.RawGetString("width").(golua.LNumber)
			height := t.RawGetString("height").(golua.LNumber)
			content := t.RawGetString("content").(*golua.LTable)

			sum := 0.0

			for i := 0; i < int(width); i++ {
				col := content.RawGetInt(i + 1).(*golua.LTable)
				for j := 0; j < int(height); j++ {
					sum += math.Abs(float64(col.RawGetInt(j + 1).(golua.LNumber)))
				}
			}

			if sum == 0 {
				return
			}

			for i := 0; i < int(width); i++ {
				col := content.RawGetInt(i + 1).(*golua.LTable)
				for j := 0; j < int(height); j++ {
					value := float64(col.RawGetInt(j + 1).(golua.LNumber))
					col.RawSetInt(j+1, golua.LNumber(value/sum))
				}
			}
		})

	lib.BuilderFunction(state, t, "normalize_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "normalize"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__normalize", golua.LBool(args["normalize"].(bool)))
		})

	return t
}

func kernelBuild(t *golua.LTable) *convolution.Kernel {
	width := int(t.RawGetString("width").(golua.LNumber))
	height := int(t.RawGetString("height").(golua.LNumber))

	content := t.RawGetString("content").(*golua.LTable)
	k := make([][]float64, width)
	for i := 0; i < width; i++ {
		k[i] = make([]float64, height)
		row := content.RawGetInt(i + 1).(*golua.LTable)
		for j := 0; j < height; j++ {
			k[i][j] = float64(row.RawGetInt(j + 1).(golua.LNumber))
		}
	}

	kernel := convolution.Kernel{
		Width:   width,
		Height:  height,
		Content: k,
	}

	if t.RawGetString("__normalize").(golua.LBool) {
		kernel.Normalize()
	}

	return &kernel
}
