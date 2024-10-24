package lib

import (
	"fmt"
	"image"
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TEST = "test"

/// @lib Testing
/// @import test
/// @desc
/// A library for testing lua workflows.

func RegisterTest(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TEST, r, r.State, lg)

	/// @func assert(cond, msg?)
	/// @arg cond {bool}
	/// @arg? msg {string}
	lib.CreateFunction(tab, "assert",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "cond"},
			{Type: lua.STRING, Name: "msg", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cond := args["cond"].(bool)
			if cond {
				return 0
			}

			msg := args["msg"].(string)
			if msg != "" {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("assertion failed: %s", msg), log.LEVEL_ERROR)), 0)
				return 0
			}
			state.Error(golua.LString("assertion failed"), 0)

			return 0
		})

	/// @func assert_image(img1, img2, msg?)
	/// @arg img1 {int<collection.IMAGE>}
	/// @arg img2 {int<collection.IMAGE>}
	/// @arg? msg {string}
	lib.CreateFunction(tab, "assert_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "img1"},
			{Type: lua.INT, Name: "img2"},
			{Type: lua.STRING, Name: "msg", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			img1 := args["img1"].(int)
			img2 := args["img2"].(int)
			if img1 == img2 {
				return 0
			}

			msg := args["msg"].(string)

			var img image.Image
			imgReady := make(chan struct{}, 2)
			imgFinished := make(chan struct{}, 2)

			r.IC.Schedule(state, img1, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
					imgReady <- struct{}{}
					<-imgFinished
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgReady <- struct{}{}
					state.Error(golua.LString("compare failed while retrieving image1"), 0)
				},
			})

			r.IC.Schedule(state, img2, &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					<-imgReady
					equal := imageutil.ImageCompare(img, i.Self.Image)

					if !equal {
						if msg != "" {
							state.Error(golua.LString(lg.Append(fmt.Sprintf("assertion failed: %s", msg), log.LEVEL_ERROR)), 0)
							return
						}
						state.Error(golua.LString("assertion failed"), 0)
					}

					imgFinished <- struct{}{}
				},
				Fail: func(i *collection.Item[collection.ItemImage]) {
					imgFinished <- struct{}{}
					state.Error(golua.LString("compare failed while processing image2"), 0)
				},
			})

			return 0
		})

	/// @func benchmark_start() -> int
	/// @returns {int} - Start time.
	lib.CreateFunction(tab, "benchmark_start",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano()
			state.Push(golua.LNumber(t))
			return 1
		})

	/// @func benchmark_end(start) -> int
	/// @arg start {int}
	/// @returns {int} - Ellapsed time.
	lib.CreateFunction(tab, "benchmark_end",
		[]lua.Arg{
			{Type: lua.INT, Name: "start"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano()
			start := int64(args["start"].(int))
			ellapsed := t - start

			seconds := ellapsed / int64(time.Second)
			ms := (ellapsed - (seconds * int64(time.Second))) / int64(time.Millisecond)

			fmt.Printf("Benchmark finished in: %ds %dms.", seconds, ms)

			state.Push(golua.LNumber(ellapsed))
			return 1
		})
}
