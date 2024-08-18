package lib

import (
	"fmt"
	"math"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/ojrac/opensimplex-go"
	golua "github.com/yuin/gopher-lua"
)

const LIB_NOISE = "noise"

/// @lib Noise
/// @import noise
/// @desc
/// Library for generating and interacting with noise maps.

func RegisterNoise(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_NOISE, r, r.State, lg)

	/// @func simplex_image_new(seed, coef, normalize, name, encoding, width, model?, disableColor?, disableAlpha?) -> int<collection.IMAGE>
	/// @arg seed {int}
	/// @arg coef {float}
	/// @arg normalize {bool} - Use noise values between (0,1) instead of (-1,1).
	/// @arg name {string}
	/// @arg encoding {int<image.Encoding>}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg? model {int<image.ColorModel>}
	/// @arg? disableColor {bool} - Don't set r,g,b values using noise.
	/// @arg? disableAlpha {bool} - Don't set alpha values using noise.
	/// @returns {int<collection.IMAGE>}
	/// @desc
	/// Creates a new image, setting each pixel of the image to
	/// the simplex_2d result of the x,y pos multiplied by coef.
	/// If color is disabled, r,g,b values are set to 255.
	/// If alpha is disabled, alpha is set to 255.
	lib.CreateFunction(tab, "simplex_image_new",
		[]lua.Arg{
			{Type: lua.INT, Name: "seed"},
			{Type: lua.FLOAT, Name: "coef"},
			{Type: lua.BOOL, Name: "normalize"},
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "encoding"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "model", Optional: true},
			{Type: lua.BOOL, Name: "disableColor", Optional: true},
			{Type: lua.BOOL, Name: "disableAlpha", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			chLog := log.NewLogger(fmt.Sprintf("image_%s", name), lg)
			lg.Append(fmt.Sprintf("child log created: image_%s", name), log.LEVEL_INFO)

			id := r.IC.AddItem(&chLog)

			r.IC.Schedule(id, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					model := lua.ParseEnum(args["model"].(int), imageutil.ModelList, lib)

					i.Self = &collection.ItemImage{
						Image:    imageutil.NewImage(args["width"].(int), args["height"].(int), model),
						Encoding: lua.ParseEnum(args["encoding"].(int), imageutil.EncodingList, lib),
						Name:     name,
						Model:    model,
					}

					coef := args["coef"].(float64)

					x := i.Self.Image.Bounds().Min.X
					y := i.Self.Image.Bounds().Min.Y
					width := i.Self.Image.Bounds().Dx()
					height := i.Self.Image.Bounds().Dy()

					var noise opensimplex.Noise
					if args["normalize"].(bool) {
						noise = opensimplex.NewNormalized(int64(args["seed"].(int)))
					} else {
						noise = opensimplex.New(int64(args["seed"].(int)))
					}

					dc := args["disableColor"].(bool)
					da := args["disableAlpha"].(bool)

					for ix := x; ix < x+width; ix++ {
						for iy := y; iy < y+height; iy++ {
							px := ix - x
							py := iy - y

							v := noise.Eval2(float64(px)*coef, float64(py)*coef)
							cv := int(math.Round(255 * v))
							c := 255
							a := 255

							if !dc {
								c = cv
							}
							if !da {
								a = cv
							}

							imageutil.Set(i.Self.Image, ix, iy, c, c, c, a)
						}
					}
				},
			})

			state.Push(golua.LNumber(id))
			return 1
		})

	/// @func simplex_image_map(seed, coef, id, normalize, disableColor?, disableAlpha?, keep?)
	/// @arg seed {int}
	/// @arg coef {float}
	/// @arg id {int<collection.IMAGE>}
	/// @arg normalize {bool} - Use noise values between (0,1) instead of (-1,1).
	/// @arg? disableColor {bool} - Don't set r,g,b values using noise.
	/// @arg? disableAlpha {bool} - Don't set alpha values using noise.
	/// @arg? keep {bool} - Maintain color or alpha channels if they're disabled.
	/// @desc
	/// Loops over the pixels of an image, setting each pixel of the image to
	/// the simplex_2d result of the x,y pos multiplied by coef.
	/// If color is disabled, r,g,b values are set to 255.
	/// If alpha is disabled, alpha is set to 255.
	/// If keep is set, the disabled values will be kept in the image, useful for only changing the alpha channel of an image.
	lib.CreateFunction(tab, "simplex_image_map",
		[]lua.Arg{
			{Type: lua.INT, Name: "seed"},
			{Type: lua.FLOAT, Name: "coef"},
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "normalize"},
			{Type: lua.BOOL, Name: "disableColor", Optional: true},
			{Type: lua.BOOL, Name: "disableAlpha", Optional: true},
			{Type: lua.BOOL, Name: "keep", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			r.IC.Schedule(args["id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					coef := args["coef"].(float64)

					x := i.Self.Image.Bounds().Min.X
					y := i.Self.Image.Bounds().Min.Y
					width := i.Self.Image.Bounds().Dx()
					height := i.Self.Image.Bounds().Dy()

					var noise opensimplex.Noise
					if args["normalize"].(bool) {
						noise = opensimplex.NewNormalized(int64(args["seed"].(int)))
					} else {
						noise = opensimplex.New(int64(args["seed"].(int)))
					}

					dc := args["disableColor"].(bool)
					da := args["disableAlpha"].(bool)
					keep := args["keep"].(bool)

					for ix := x; ix < x+width; ix++ {
						for iy := y; iy < y+height; iy++ {
							px := ix - x
							py := iy - y

							v := noise.Eval2(float64(px)*coef, float64(py)*coef)
							cv := int(math.Round(255 * v))
							cr, cg, cb, ca := imageutil.Get(i.Self.Image, ix, iy)

							if dc {
								if !keep {
									cr = 255
									cg = 255
									cb = 255
								}
							} else {
								cr = cv
								cg = cv
								cb = cv
							}

							if da {
								if !keep {
									ca = 255
								}
							} else {
								ca = cv
							}

							imageutil.Set(i.Self.Image, ix, iy, cr, cg, cb, ca)
						}
					}
				},
			})

			return 0
		})

	/// @func simplex_2d(seed, x, y, normalize?) - float
	/// @arg seed {int}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg? normalize {bool} - Use noise values between (0,1) instead of (-1,1).
	/// @returns {float}
	lib.CreateFunction(tab, "simplex_2d",
		[]lua.Arg{
			{Type: lua.INT, Name: "seed"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.BOOL, Name: "normalize", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var noise opensimplex.Noise
			if args["normalize"].(bool) {
				noise = opensimplex.NewNormalized(int64(args["seed"].(int)))
			} else {
				noise = opensimplex.New(int64(args["seed"].(int)))
			}

			v := noise.Eval2(args["x"].(float64), args["y"].(float64))

			state.Push(golua.LNumber(v))
			return 1
		})

	/// @func simplex_3d(seed, x, y, z, normalize?) -> float
	/// @arg seed {int}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg z {float}
	/// @arg? normalize {bool} - Use noise values between (0,1) instead of (-1,1).
	/// @returns {float}
	lib.CreateFunction(tab, "simplex_3d",
		[]lua.Arg{
			{Type: lua.INT, Name: "seed"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "z"},
			{Type: lua.BOOL, Name: "normalize", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var noise opensimplex.Noise
			if args["normalize"].(bool) {
				noise = opensimplex.NewNormalized(int64(args["seed"].(int)))
			} else {
				noise = opensimplex.New(int64(args["seed"].(int)))
			}

			v := noise.Eval3(args["x"].(float64), args["y"].(float64), args["z"].(float64))

			state.Push(golua.LNumber(v))
			return 1
		})

	/// @func simplex_4d(seed, x, y, z, w, normalize?) -> float
	/// @arg seed {int}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg z {float}
	/// @arg w {float}
	/// @arg? normalize {bool} - Use noise values between (0,1) instead of (-1,1).
	/// @returns {float}
	lib.CreateFunction(tab, "simplex_4d",
		[]lua.Arg{
			{Type: lua.INT, Name: "seed"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "z"},
			{Type: lua.FLOAT, Name: "w"},
			{Type: lua.BOOL, Name: "normalize", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var noise opensimplex.Noise
			if args["normalize"].(bool) {
				noise = opensimplex.NewNormalized(int64(args["seed"].(int)))
			} else {
				noise = opensimplex.New(int64(args["seed"].(int)))
			}

			v := noise.Eval4(args["x"].(float64), args["y"].(float64), args["z"].(float64), args["w"].(float64))

			state.Push(golua.LNumber(v))
			return 1
		})
}
