package lib

import (
	"fmt"
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_GUI = "gui"

func RegisterGUI(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_GUI, r, r.State, lg)

	/// @func window()
	/// @arg name
	/// @arg width
	/// @arg height
	/// @returns id of the window.
	lib.CreateFunction(tab, "window",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w := g.NewMasterWindow(args["name"].(string), args["width"].(int), args["height"].(int), 0)
			ind := r.CR_WIN.Add(w)

			state.Push(golua.LNumber(ind))
			return 1
		})

	/// @func window_pos()
	/// @arg id
	/// @returns x, y
	lib.CreateFunction(tab, "window_pos",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			x, y := w.GetPos()
			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func window_set_pos()
	/// @arg id
	/// @arg x
	/// @arg y
	lib.CreateFunction(tab, "window_set_pos",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetPos(args["x"].(int), args["y"].(int))
			return 0
		})

	/// @func window_size()
	/// @arg id
	/// @returns width, height
	lib.CreateFunction(tab, "window_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			width, height := w.GetSize()
			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func window_set_size()
	/// @arg id
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "window_set_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetSize(args["width"].(int), args["height"].(int))
			return 0
		})

	/// @func window_set_size_limits()
	/// @arg id
	/// @arg minw
	/// @arg minh
	/// @arg maxw
	/// @arg maxh
	lib.CreateFunction(tab, "window_set_size_limits",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "minw"},
			{Type: lua.INT, Name: "minh"},
			{Type: lua.INT, Name: "maxw"},
			{Type: lua.INT, Name: "maxh"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetSizeLimits(args["minw"].(int), args["minh"].(int), args["maxw"].(int), args["maxh"].(int))
			return 0
		})

	/// @func window_set_bg_color()
	/// @arg id
	/// @arg r
	/// @arg g
	/// @arg b
	/// @arg a
	lib.CreateFunction(tab, "window_set_bg_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "r"},
			{Type: lua.INT, Name: "g"},
			{Type: lua.INT, Name: "b"},
			{Type: lua.INT, Name: "a"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			c := color.NRGBA{
				R: uint8(args["r"].(int)),
				G: uint8(args["g"].(int)),
				B: uint8(args["b"].(int)),
				A: uint8(args["a"].(int)),
			}

			w.SetBgColor(c)
			return 0
		})

	/// @func window_should_close()
	/// @arg id
	/// @arg v - bool
	lib.CreateFunction(tab, "window_should_close",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "v"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetShouldClose(args["v"].(bool))
			return 0
		})

	/// @func window_set_icon()
	/// @arg id
	/// @arg icon_id
	lib.CreateFunction(tab, "window_set_icon",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "icon_id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			r.IC.Schedule(args["icon_id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					w.SetIcon(imageutil.CopyImage(i.Self.Image, imageutil.MODEL_NRGBA))
				},
			})

			return 0
		})

	/// @func window_set_fps()
	/// @arg id
	/// @arg fps
	lib.CreateFunction(tab, "window_set_fps",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "fps"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetTargetFPS(uint(args["fps"].(int)))
			return 0
		})

	/// @func window_set_title()
	/// @arg id
	/// @arg title
	lib.CreateFunction(tab, "window_set_size_limits",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.INT, Name: "minw"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetSizeLimits(args["minw"].(int), args["minh"].(int), args["maxw"].(int), args["maxh"].(int))
			return 0
		})

	/// @func window_run()
	/// @arg id
	/// @arg fn
	/// @blocking
	lib.CreateFunction(tab, "window_run",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.Run(func() {
				defer func() {
					if p := recover(); p != nil {
						w.Close()
						g.Update()
						panic(p)
					}
				}()

				state.Push(args["fn"].(*golua.LFunction))
				state.Call(0, 0)
			})

			return 0
		})

	/// @func wg_single_window()
	/// @arg? widgets - []Widgets
	lib.CreateFunction(tab, "wg_single_window",
		[]lua.Arg{
			lua.ArgArray("widgets", lua.ArrayType{Type: lua.ANY}, true),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			wts := []*golua.LTable{}

			wa := args["widgets"].(map[string]any)

			for i := range len(wa) {
				wt := wa[string(i+1)]

				if t, ok := wt.(*golua.LTable); ok {
					wts = append(wts, t)
				} else {
					state.Error(golua.LString(lg.Append("invalid table provided as widget to wg_single_window", log.LEVEL_ERROR)), 0)
				}
			}

			w := layoutBuild(state, wts)
			g.SingleWindow().Layout(w...)

			return 0
		})

	/// @func wg_label()
	/// @arg? text
	/// @returns widget
	lib.CreateFunction(tab, "wg_label",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := labelTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_button()
	/// @arg text
	/// @returns widget
	lib.CreateFunction(tab, "wg_button",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text", Optional: false},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})
}

const (
	WIDGET_LABEL  = "label"
	WIDGET_BUTTON = "button"
)

func layoutBuild(state *golua.LState, widgets []*golua.LTable) []g.Widget {
	w := []g.Widget{}

	for _, wt := range widgets {
		t := state.GetTable(wt, golua.LString("type")).String()

		switch t {
		case WIDGET_LABEL:
			w = append(w, labelBuild(state, wt))
		case WIDGET_BUTTON:
			w = append(w, buttonBuild(state, wt))
		}
	}

	return w
}

func labelTable(state *golua.LState, text string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LABEL))
	state.SetTable(t, golua.LString("label"), golua.LString(text))

	return t
}

func labelBuild(state *golua.LState, t *golua.LTable) *g.LabelWidget {
	return g.Label(state.GetTable(t, golua.LString("label")).String())
}

func buttonTable(state *golua.LState, text string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("__disabled"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__click"), golua.LNil)

	tableBuilderFunc(state, t, "disabled", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__disabled"), golua.LBool(v))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), golua.LNumber(width))
		state.SetTable(t, golua.LString("__height"), golua.LNumber(height))
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	return t
}

func buttonBuild(state *golua.LState, t *golua.LTable) *g.ButtonWidget {
	b := g.Button(state.GetTable(t, golua.LString("label")).String())

	disabled := state.GetTable(t, golua.LString("__disabled"))
	if disabled.Type() == golua.LTBool {
		b.Disabled(bool(disabled.(golua.LBool)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func tableBuilderFunc(state *golua.LState, t *golua.LTable, name string, fn func(state *golua.LState, t *golua.LTable)) {
	state.SetTable(t, golua.LString(name), state.NewFunction(func(state *golua.LState) int {
		self := state.CheckTable(1)

		fn(state, self)

		state.Push(self)
		return 1
	}))
}
