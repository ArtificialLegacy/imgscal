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
			wts := parseWidgets(args["widgets"].(map[string]any), state, lg)
			w := layoutBuild(r, state, wts, lg)
			g.SingleWindow().Layout(w...)

			return 0
		})

	/// @func wg_label()
	/// @arg text
	/// @returns widget
	lib.CreateFunction(tab, "wg_label",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
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
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_dummy()
	/// @arg width
	/// @arg height
	/// @returns widget
	lib.CreateFunction(tab, "wg_dummy",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := dummyTable(state, args["width"].(float64), args["height"].(float64))

			state.Push(t)
			return 1
		})

	/// @func wg_separator()
	/// @returns widget
	lib.CreateFunction(tab, "wg_separator",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := separatorTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_bullet_text()
	/// @arg text
	/// @returns widget
	lib.CreateFunction(tab, "wg_bullet_text",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := bulletTextTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_bullet()
	/// @returns widget
	lib.CreateFunction(tab, "wg_bullet",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := bulletTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_checkbox()
	/// @arg text
	/// @arg boolref
	/// @returns widget
	lib.CreateFunction(tab, "wg_checkbox",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.INT, Name: "boolref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := checkboxTable(state, args["text"].(string), args["boolref"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_child()
	/// @returns widget
	lib.CreateFunction(tab, "wg_child",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := childTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_color_edit()
	/// @arg text
	/// @arg colorref
	/// @returns widget
	lib.CreateFunction(tab, "wg_color_edit",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.INT, Name: "colorref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := colorEditTable(state, args["text"].(string), args["colorref"].(int))

			state.Push(t)
			return 1
		})
}

func tableBuilderFunc(state *golua.LState, t *golua.LTable, name string, fn func(state *golua.LState, t *golua.LTable)) {
	state.SetTable(t, golua.LString(name), state.NewFunction(func(state *golua.LState) int {
		self := state.CheckTable(1)

		fn(state, self)

		state.Push(self)
		return 1
	}))
}

const (
	WIDGET_LABEL       = "label"
	WIDGET_BUTTON      = "button"
	WIDGET_DUMMY       = "dummy"
	WIDGET_SEPARATOR   = "separator"
	WIDGET_BULLET_TEXT = "bullet_text"
	WIDGET_BULLET      = "bullet"
	WIDGET_CHECKBOX    = "checkbox"
	WIDGET_CHILD       = "child"
	WIDGET_COLOR_EDIT  = "color_edit"
)

var buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{}

func init() {
	buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{
		WIDGET_LABEL:       labelBuild,
		WIDGET_BUTTON:      buttonBuild,
		WIDGET_DUMMY:       dummyBuild,
		WIDGET_SEPARATOR:   separatorBuild,
		WIDGET_BULLET_TEXT: bulletTextBuild,
		WIDGET_BULLET:      bulletBuild,
		WIDGET_CHECKBOX:    checkboxBuild,
		WIDGET_CHILD:       childBuild,
		WIDGET_COLOR_EDIT:  colorEditBuild,
	}
}

func parseWidgets(widgetTable map[string]any, state *golua.LState, lg *log.Logger) []*golua.LTable {
	wts := []*golua.LTable{}

	for i := range len(widgetTable) {
		wt := widgetTable[string(i+1)]

		if t, ok := wt.(*golua.LTable); ok {
			wts = append(wts, t)
		} else {
			state.Error(golua.LString(lg.Append("invalid table provided as widget to wg_single_window", log.LEVEL_ERROR)), 0)
		}
	}

	return wts
}

func parseTable(t *golua.LTable, state *golua.LState) map[string]any {
	v := map[string]any{}

	ln := t.Len()
	for i := range ln {
		v[string(i+1)] = state.GetTable(t, golua.LNumber(i+1))
	}

	return v
}

func layoutBuild(r *lua.Runner, state *golua.LState, widgets []*golua.LTable, lg *log.Logger) []g.Widget {
	w := []g.Widget{}

	for _, wt := range widgets {
		t := state.GetTable(wt, golua.LString("type")).String()

		build, ok := buildList[t]
		if !ok {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown widget: %s", t), log.LEVEL_ERROR)), 0)
		}

		w = append(w, build(r, lg, state, wt))
	}

	return w
}

func labelTable(state *golua.LState, text string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LABEL))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("__wrapped"), golua.LNil)

	tableBuilderFunc(state, t, "wrapped", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__wrapped"), golua.LBool(v))
	})

	return t
}

func labelBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	l := g.Label(state.GetTable(t, golua.LString("label")).String())

	wrapped := state.GetTable(t, golua.LString("__wrapped"))
	if wrapped.Type() == golua.LTBool {
		l.Wrapped(bool(wrapped.(golua.LBool)))
	}

	return l
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

func buttonBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
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

func dummyTable(state *golua.LState, width, height float64) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_DUMMY))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))

	return t
}

func dummyBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	w := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	h := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	d := g.Dummy(float32(w), float32(h))

	return d
}

func separatorTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SEPARATOR))

	return t
}

func separatorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	s := g.Separator()

	return s
}

func bulletTextTable(state *golua.LState, text string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BULLET_TEXT))
	state.SetTable(t, golua.LString("text"), golua.LString(text))

	return t
}

func bulletTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.BulletText(state.GetTable(t, golua.LString("text")).String())

	return b
}

func bulletTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BULLET))

	return t
}

func bulletBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Bullet()

	return b
}

func checkboxTable(state *golua.LState, text string, boolref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CHECKBOX))
	state.SetTable(t, golua.LString("text"), golua.LString(text))
	state.SetTable(t, golua.LString("boolref"), golua.LNumber(boolref))
	state.SetTable(t, golua.LString("__change"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	return t
}

func checkboxBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("boolref")).(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	selected := sref.Value.(*bool)
	c := g.Checkbox(state.GetTable(t, golua.LString("text")).String(), selected)

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LBool(*selected))
			state.Call(1, 0)
		})
	}

	return c
}

func childTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CHILD))
	state.SetTable(t, golua.LString("__border"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__border"), golua.LBool(b))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), golua.LNumber(width))
		state.SetTable(t, golua.LString("__height"), golua.LNumber(height))
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func childBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.Child()

	border := state.GetTable(t, golua.LString("__border"))
	if border.Type() == golua.LTBool {
		c.Border(bool(border.(golua.LBool)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func colorEditTable(state *golua.LState, text string, colorref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_COLOR_EDIT))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("colorref"), golua.LNumber(colorref))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), golua.LNumber(width))
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	return t
}

func colorEditBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("colorref")).(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	selected := sref.Value.(*color.RGBA)
	c := g.ColorEdit(state.GetTable(t, golua.LString("label")).String(), selected)

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			ct := imageutil.RGBAToTable(state, selected)

			state.Push(change)
			state.Push(ct)
			state.Call(1, 0)
		})
	}

	return c
}
