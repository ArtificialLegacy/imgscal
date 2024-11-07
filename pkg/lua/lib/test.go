package lib

import (
	"fmt"
	"image"
	"slices"
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
/// A library for testing lua workflows. This library is always available, and does not need to be imported.

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

	/// @func assert_schema(value, schema, msg?)
	/// @arg value {table<any>}
	/// @arg schema {table<any>}
	/// @arg? msg {string}
	lib.CreateFunction(tab, "assert_schema",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "value"},
			{Type: lua.RAW_TABLE, Name: "schema"},
			{Type: lua.STRING, Name: "msg", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			value := args["value"].(*golua.LTable)
			schema := args["schema"].(*golua.LTable)
			msg := args["msg"].(string)

			valid := validateSchema(value, schema)

			if !valid {
				if msg != "" {
					lua.Error(state, lg.Appendf("assertion failed: %s", log.LEVEL_ERROR, msg))
					return 0
				}
				lua.Error(state, "assertion failed")
			}

			return 0
		})

	/// @func assert_imported(name, msg?)
	/// @arg name {string}
	/// @arg? msg {string}
	lib.CreateFunction(tab, "assert_imported",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.STRING, Name: "msg", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)
			msg := args["msg"].(string)

			if !slices.Contains(r.Libraries, name) {
				if msg != "" {
					lua.Error(state, lg.Appendf("assertion failed: %s", log.LEVEL_ERROR, msg))
					return 0
				}
				lua.Error(state, "assertion failed")
			}

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

func validateSchema(value, schema *golua.LTable) bool {
	valid := true

	schema.ForEach(func(k, v1 golua.LValue) {
		v2 := value.RawGet(k)
		if v1.Type() != v2.Type() {
			valid = false
		} else if v1.Type() == golua.LTTable {
			validNested := validateSchema(v1.(*golua.LTable), v2.(*golua.LTable))
			if !validNested {
				valid = false
			}
		}
	})

	return valid
}
