package lib

import (
	"fmt"
	"image/color"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/AllenDang/giu"
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

	/// @func wg_column()
	/// @arg? widgets - []Widgets
	lib.CreateFunction(tab, "wg_column",
		[]lua.Arg{
			{Type: lua.ANY, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["widgets"]
			if v == nil {
				v = golua.LNil
			}
			t := columnTable(state, v.(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_row()
	/// @arg? widgets - []Widgets
	lib.CreateFunction(tab, "wg_row",
		[]lua.Arg{
			{Type: lua.ANY, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["widgets"]
			if v == nil {
				v = golua.LNil
			}
			t := rowTable(state, v.(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_combo_custom()
	/// @arg text
	/// @arg preview
	/// @returns widget
	lib.CreateFunction(tab, "wg_combo_custom",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.STRING, Name: "preview"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := comboCustomTable(state, args["text"].(string), args["preview"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_combo()
	/// @arg text
	/// @arg preview
	/// @arg items - []string
	/// @arg i32ref
	/// @returns widget
	lib.CreateFunction(tab, "wg_combo",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.STRING, Name: "preview"},
			{Type: lua.ANY, Name: "items"},
			{Type: lua.INT, Name: "i32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			text := args["text"].(string)
			preview := args["preview"].(string)
			items := args["items"].(golua.LValue)
			i32ref := args["i32ref"].(int)
			t := comboTable(state, text, preview, items, i32ref)

			state.Push(t)
			return 1
		})

	/// @func wg_condition()
	/// @arg condition - boolean
	/// @arg widgetIf
	/// @arg widgetElse
	lib.CreateFunction(tab, "wg_condition",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "condition"},
			{Type: lua.ANY, Name: "widgetIf"},
			{Type: lua.ANY, Name: "widgetElse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			widgetIf := args["widgetIf"].(golua.LValue)
			widgetElse := args["widgetElse"].(golua.LValue)
			t := conditionTable(state, args["condition"].(bool), widgetIf, widgetElse)

			state.Push(t)
			return 1
		})

	/// @func wg_context_menu()
	/// @returns widget
	lib.CreateFunction(tab, "wg_context_menu",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := contextMenuTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_date_picker()
	/// @arg id
	/// @arg timeref
	/// @returns widget
	lib.CreateFunction(tab, "wg_date_picker",
		[]lua.Arg{
			{Type: lua.STRING, Name: "id"},
			{Type: lua.INT, Name: "timeref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			id := args["id"].(string)
			timeref := args["timeref"].(int)
			t := datePickerTable(state, id, timeref)

			state.Push(t)
			return 1
		})

	/// @func wg_drag_int()
	/// @arg label
	/// @arg i32ref
	/// @arg minvalue
	/// @arg maxvalue
	/// @returns widget
	lib.CreateFunction(tab, "wg_drag_int",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
			{Type: lua.INT, Name: "i32ref"},
			{Type: lua.INT, Name: "minvalue"},
			{Type: lua.INT, Name: "maxvalue"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			label := args["label"].(string)
			i32ref := args["i32ref"].(int)
			minvalue := args["minvalue"].(int)
			maxvalue := args["maxvalue"].(int)
			t := dragIntTable(state, label, i32ref, minvalue, maxvalue)

			state.Push(t)
			return 1
		})

	/// @func wg_input_float()
	/// @arg f32ref
	/// @returns widget
	lib.CreateFunction(tab, "wg_input_float",
		[]lua.Arg{
			{Type: lua.INT, Name: "f32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			f32ref := args["f32ref"].(int)
			t := inputFloatTable(state, f32ref)

			state.Push(t)
			return 1
		})

	/// @func wg_input_int()
	/// @arg i32ref
	/// @returns widget
	lib.CreateFunction(tab, "wg_input_int",
		[]lua.Arg{
			{Type: lua.INT, Name: "i32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			i32ref := args["i32ref"].(int)
			t := inputIntTable(state, i32ref)

			state.Push(t)
			return 1
		})

	// @func wg_input_text()
	/// @arg strref
	/// @returns widget
	lib.CreateFunction(tab, "wg_input_text",
		[]lua.Arg{
			{Type: lua.INT, Name: "strref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			strref := args["strref"].(int)
			t := inputTextTable(state, strref)

			state.Push(t)
			return 1
		})

	// @func wg_input_text_multiline()
	/// @arg strref
	/// @returns widget
	lib.CreateFunction(tab, "wg_input_text_multiline",
		[]lua.Arg{
			{Type: lua.INT, Name: "strref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			strref := args["strref"].(int)
			t := inputMultilineTextTable(state, strref)

			state.Push(t)
			return 1
		})

	// @func wg_progress_bar()
	/// @arg fraction
	/// @returns widget
	lib.CreateFunction(tab, "wg_progress_bar",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "fraction"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fraction := args["fraction"].(float64)
			t := progressBarTable(state, fraction)

			state.Push(t)
			return 1
		})

	// @func wg_progress_indicator()
	/// @arg label
	/// @arg width
	/// @arg height
	/// @arg radius
	/// @returns widget
	lib.CreateFunction(tab, "wg_progress_indicator",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
			{Type: lua.FLOAT, Name: "radius"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			label := args["label"].(string)
			width := args["width"].(float64)
			height := args["height"].(float64)
			radius := args["radius"].(float64)

			t := progressIndicatorTable(state, label, width, height, radius)

			state.Push(t)
			return 1
		})

	/// @func wg_spacing()
	/// @returns widget
	lib.CreateFunction(tab, "wg_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spacingTable(state)

			state.Push(t)
			return 1
		})

	/// @constants Color Picker Flags
	/// @const FLAGCOLOREDIT_NONE
	/// @const FLAGCOLOREDIT_NOALPHA
	/// @const FLAGCOLOREDIT_NOPICKER
	/// @const FLAGCOLOREDIT_NOOPTIONS
	/// @const FLAGCOLOREDIT_NOSMALLPREVIEW
	/// @const FLAGCOLOREDIT_NOINPUTS
	/// @const FLAGCOLOREDIT_NOTOOLTIP
	/// @const FLAGCOLOREDIT_NOLABEL
	/// @const FLAGCOLOREDIT_NOSIDEPREVIEW
	/// @const FLAGCOLOREDIT_NODRAGDROP
	/// @const FLAGCOLOREDIT_NOBORDER
	/// @const FLAGCOLOREDIT_ALPHABAR
	/// @const FLAGCOLOREDIT_ALPHAPREVIEW
	/// @const FLAGCOLOREDIT_ALPHAPREVIEWHALF
	/// @const FLAGCOLOREDIT_HDR
	/// @const FLAGCOLOREDIT_DISPLAYRGB
	/// @const FLAGCOLOREDIT_DISPLAYHSV
	/// @const FLAGCOLOREDIT_DISPLAYHEX
	/// @const FLAGCOLOREDIT_UINT8
	/// @const FLAGCOLOREDIT_FLOAT
	/// @const FLAGCOLOREDIT_HUEBAR
	/// @const FLAGCOLOREDIT_HUEWHEEL
	/// @const FLAGCOLOREDIT_INPUTRGB
	/// @const FLAGCOLOREDIT_INPUTHSV
	/// @const FLAGCOLOREDIT_DEFAULTOPTIONS
	/// @const FLAGCOLOREDIT_DISPLAYMASK
	/// @const FLAGCOLOREDIT_DATATYPEMASK
	/// @const FLAGCOLOREDIT_PICKERMASK
	/// @const FLAGCOLOREDIT_INPUTMASK
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NONE"), golua.LNumber(FLAGCOLOREDIT_NONE))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOALPHA"), golua.LNumber(FLAGCOLOREDIT_NOALPHA))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOPICKER"), golua.LNumber(FLAGCOLOREDIT_NOPICKER))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOOPTIONS"), golua.LNumber(FLAGCOLOREDIT_NOOPTIONS))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOSMALLPREVIEW"), golua.LNumber(FLAGCOLOREDIT_NOSMALLPREVIEW))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOINPUTS"), golua.LNumber(FLAGCOLOREDIT_NOINPUTS))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOTOOLTIP"), golua.LNumber(FLAGCOLOREDIT_NOTOOLTIP))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOLABEL"), golua.LNumber(FLAGCOLOREDIT_NOLABEL))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOSIDEPREVIEW"), golua.LNumber(FLAGCOLOREDIT_NOSIDEPREVIEW))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NODRAGDROP"), golua.LNumber(FLAGCOLOREDIT_NODRAGDROP))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_NOBORDER"), golua.LNumber(FLAGCOLOREDIT_NOBORDER))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_ALPHABAR"), golua.LNumber(FLAGCOLOREDIT_ALPHABAR))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_ALPHAPREVIEW"), golua.LNumber(FLAGCOLOREDIT_ALPHAPREVIEW))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_ALPHAPREVIEWHALF"), golua.LNumber(FLAGCOLOREDIT_ALPHAPREVIEWHALF))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_HDR"), golua.LNumber(FLAGCOLOREDIT_HDR))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DISPLAYRGB"), golua.LNumber(FLAGCOLOREDIT_DISPLAYRGB))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DISPLAYHSV"), golua.LNumber(FLAGCOLOREDIT_DISPLAYHSV))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DISPLAYHEX"), golua.LNumber(FLAGCOLOREDIT_DISPLAYHEX))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_UINT8"), golua.LNumber(FLAGCOLOREDIT_UINT8))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_FLOAT"), golua.LNumber(FLAGCOLOREDIT_FLOAT))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_HUEBAR"), golua.LNumber(FLAGCOLOREDIT_HUEBAR))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_HUEWHEEL"), golua.LNumber(FLAGCOLOREDIT_HUEWHEEL))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_INPUTRGB"), golua.LNumber(FLAGCOLOREDIT_INPUTRGB))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_INPUTHSV"), golua.LNumber(FLAGCOLOREDIT_INPUTHSV))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DEFAULTOPTIONS"), golua.LNumber(FLAGCOLOREDIT_DEFAULTOPTIONS))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DISPLAYMASK"), golua.LNumber(FLAGCOLOREDIT_DISPLAYMASK))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_DATATYPEMASK"), golua.LNumber(FLAGCOLOREDIT_DATATYPEMASK))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_PICKERMASK"), golua.LNumber(FLAGCOLOREDIT_PICKERMASK))
	r.State.SetTable(tab, golua.LString("FLAGCOLOREDIT_INPUTMASK"), golua.LNumber(FLAGCOLOREDIT_INPUTMASK))

	/// @constants Combo Flags
	/// @const FLAGCOMBO_NONE
	/// @const FLAGCOMBO_POPUPALIGNLEFT
	/// @const FLAGCOMBO_HEIGHTSMALL
	/// @const FLAGCOMBO_HEIGHTREGULAR
	/// @const FLAGCOMBO_HEIGHTLARGE
	/// @const FLAGCOMBO_HEIGHTLARGEST
	/// @const FLAGCOMBO_NOARROWBUTTON
	/// @const FLAGCOMBO_NOPREVIEW
	/// @const FLAGCOMBO_WIDTHFITPREVIEW
	/// @const FLAGCOMBO_HEIGHTMASK
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_NONE"), golua.LNumber(FLAGCOMBO_NONE))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_POPUPALIGNLEFT"), golua.LNumber(FLAGCOMBO_POPUPALIGNLEFT))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_HEIGHTSMALL"), golua.LNumber(FLAGCOMBO_HEIGHTSMALL))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_HEIGHTREGULAR"), golua.LNumber(FLAGCOMBO_HEIGHTREGULAR))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_HEIGHTLARGEST"), golua.LNumber(FLAGCOMBO_HEIGHTLARGEST))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_NOARROWBUTTON"), golua.LNumber(FLAGCOMBO_NOARROWBUTTON))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_NOARROWBUTTON"), golua.LNumber(FLAGCOMBO_NOARROWBUTTON))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_NOPREVIEW"), golua.LNumber(FLAGCOMBO_NOPREVIEW))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_WIDTHFITPREVIEW"), golua.LNumber(FLAGCOMBO_WIDTHFITPREVIEW))
	r.State.SetTable(tab, golua.LString("FLAGCOMBO_HEIGHTMASK"), golua.LNumber(FLAGCOMBO_HEIGHTMASK))

	/// @constants Mouse Buttons
	/// @const MOUSEBUTTON_LEFT
	/// @const MOUSEBUTTON_RIGHT
	/// @const MOUSEBUTTON_MIDDLE
	r.State.SetTable(tab, golua.LString("MOUSEBUTTON_LEFT"), golua.LNumber(MOUSEBUTTON_LEFT))
	r.State.SetTable(tab, golua.LString("MOUSEBUTTON_RIGHT"), golua.LNumber(MOUSEBUTTON_RIGHT))
	r.State.SetTable(tab, golua.LString("MOUSEBUTTON_MIDDLE"), golua.LNumber(MOUSEBUTTON_MIDDLE))

	/// @constants Date Picker Labels
	/// @const DATEPICKERLABEL_MONTH
	/// @const DATEPICKERLABEL_YEAR
	r.State.SetTable(tab, golua.LString("DATEPICKERLABEL_MONTH"), golua.LString(DATEPICKERLABEL_MONTH))
	r.State.SetTable(tab, golua.LString("DATEPICKERLABEL_YEAR"), golua.LString(DATEPICKERLABEL_YEAR))

	/// @constants Input Text Flags
	/// @const FLAGINPUTTEXT_NONE
	/// @const FLAGINPUTTEXT_CHARSDECIMAL
	/// @const FLAGINPUTTEXT_CHARSHEXADECIMAL
	/// @const FLAGINPUTTEXT_CHARSUPPERCASE
	/// @const FLAGINPUTTEXT_CHARSNOBLANK
	/// @const FLAGINPUTTEXT_AUTOSELECTALL
	/// @const FLAGINPUTTEXT_ENTERRETURNSTRUE
	/// @const FLAGINPUTTEXT_CALLBACKCOMPLETION
	/// @const FLAGINPUTTEXT_CALLBACKHISTORY
	/// @const FLAGINPUTTEXT_CALLBACKALWAYS
	/// @const FLAGINPUTTEXT_CALLBACKCHARFILTER
	/// @const FLAGINPUTTEXT_ALLOWTABINPUT
	/// @const FLAGINPUTTEXT_CTRLENTERFORNEWLINE
	/// @const FLAGINPUTTEXT_NOHORIZONTALSCROLL
	/// @const FLAGINPUTTEXT_ALWAYSOVERWRITE
	/// @const FLAGINPUTTEXT_READONLY
	/// @const FLAGINPUTTEXT_PASSWORD
	/// @const FLAGINPUTTEXT_NOUNDOREDO
	/// @const FLAGINPUTTEXT_CHARSSCIENTIFIC
	/// @const FLAGINPUTTEXT_CALLBACKRESIZE
	/// @const FLAGINPUTTEXT_CALLBACKEDIT
	/// @const FLAGINPUTTEXT_ESCAPECLEARSALL
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_NONE"), golua.LNumber(FLAGINPUTTEXT_NONE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CHARSDECIMAL"), golua.LNumber(FLAGINPUTTEXT_CHARSDECIMAL))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CHARSHEXADECIMAL"), golua.LNumber(FLAGINPUTTEXT_CHARSHEXADECIMAL))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CHARSUPPERCASE"), golua.LNumber(FLAGINPUTTEXT_CHARSUPPERCASE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CHARSNOBLANK"), golua.LNumber(FLAGINPUTTEXT_CHARSNOBLANK))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_AUTOSELECTALL"), golua.LNumber(FLAGINPUTTEXT_AUTOSELECTALL))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_ENTERRETURNSTRUE"), golua.LNumber(FLAGINPUTTEXT_ENTERRETURNSTRUE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKCOMPLETION"), golua.LNumber(FLAGINPUTTEXT_CALLBACKCOMPLETION))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKHISTORY"), golua.LNumber(FLAGINPUTTEXT_CALLBACKHISTORY))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKALWAYS"), golua.LNumber(FLAGINPUTTEXT_CALLBACKALWAYS))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKCHARFILTER"), golua.LNumber(FLAGINPUTTEXT_CALLBACKCHARFILTER))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_ALLOWTABINPUT"), golua.LNumber(FLAGINPUTTEXT_ALLOWTABINPUT))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CTRLENTERFORNEWLINE"), golua.LNumber(FLAGINPUTTEXT_CTRLENTERFORNEWLINE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_NOHORIZONTALSCROLL"), golua.LNumber(FLAGINPUTTEXT_NOHORIZONTALSCROLL))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_ALWAYSOVERWRITE"), golua.LNumber(FLAGINPUTTEXT_ALWAYSOVERWRITE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_READONLY"), golua.LNumber(FLAGINPUTTEXT_READONLY))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_PASSWORD"), golua.LNumber(FLAGINPUTTEXT_PASSWORD))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_NOUNDOREDO"), golua.LNumber(FLAGINPUTTEXT_NOUNDOREDO))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CHARSSCIENTIFIC"), golua.LNumber(FLAGINPUTTEXT_CHARSSCIENTIFIC))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKRESIZE"), golua.LNumber(FLAGINPUTTEXT_CALLBACKRESIZE))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_CALLBACKEDIT"), golua.LNumber(FLAGINPUTTEXT_CALLBACKEDIT))
	r.State.SetTable(tab, golua.LString("FLAGINPUTTEXT_ESCAPECLEARSALL"), golua.LNumber(FLAGINPUTTEXT_ESCAPECLEARSALL))
}

func tableBuilderFunc(state *golua.LState, t *golua.LTable, name string, fn func(state *golua.LState, t *golua.LTable)) {
	state.SetTable(t, golua.LString(name), state.NewFunction(func(state *golua.LState) int {
		self := state.CheckTable(1)

		fn(state, self)

		state.Push(self)
		return 1
	}))
}

// -- flags
const (
	FLAGCOMBO_NONE            int = 0b0000_0000
	FLAGCOMBO_POPUPALIGNLEFT  int = 0b0000_0001
	FLAGCOMBO_HEIGHTSMALL     int = 0b0000_0010
	FLAGCOMBO_HEIGHTREGULAR   int = 0b0000_0100
	FLAGCOMBO_HEIGHTLARGE     int = 0b0000_1000
	FLAGCOMBO_HEIGHTLARGEST   int = 0b0001_0000
	FLAGCOMBO_NOARROWBUTTON   int = 0b0010_0000
	FLAGCOMBO_NOPREVIEW       int = 0b0100_0000
	FLAGCOMBO_WIDTHFITPREVIEW int = 0b1000_0000

	FLAGCOMBO_HEIGHTMASK int = 0b0001_1110
)

const (
	FLAGCOLOREDIT_NONE             int = 0b0000_0000_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_NOALPHA          int = 0b0000_0000_0000_0000_0000_0000_0000_0010
	FLAGCOLOREDIT_NOPICKER         int = 0b0000_0000_0000_0000_0000_0000_0000_0100
	FLAGCOLOREDIT_NOOPTIONS        int = 0b0000_0000_0000_0000_0000_0000_0000_1000
	FLAGCOLOREDIT_NOSMALLPREVIEW   int = 0b0000_0000_0000_0000_0000_0000_0001_0000
	FLAGCOLOREDIT_NOINPUTS         int = 0b0000_0000_0000_0000_0000_0000_0010_0000
	FLAGCOLOREDIT_NOTOOLTIP        int = 0b0000_0000_0000_0000_0000_0000_0100_0000
	FLAGCOLOREDIT_NOLABEL          int = 0b0000_0000_0000_0000_0000_0000_1000_0000
	FLAGCOLOREDIT_NOSIDEPREVIEW    int = 0b0000_0000_0000_0000_0000_0001_0000_0000
	FLAGCOLOREDIT_NODRAGDROP       int = 0b0000_0000_0000_0000_0000_0010_0000_0000
	FLAGCOLOREDIT_NOBORDER         int = 0b0000_0000_0000_0000_0000_0100_0000_0000
	FLAGCOLOREDIT_ALPHABAR         int = 0b0000_0000_0000_0001_0000_0000_0000_0000
	FLAGCOLOREDIT_ALPHAPREVIEW     int = 0b0000_0000_0000_0010_0000_0000_0000_0000
	FLAGCOLOREDIT_ALPHAPREVIEWHALF int = 0b0000_0000_0000_0100_0000_0000_0000_0000
	FLAGCOLOREDIT_HDR              int = 0b0000_0000_0000_1000_0000_0000_0000_0000
	FLAGCOLOREDIT_DISPLAYRGB       int = 0b0000_0000_0001_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_DISPLAYHSV       int = 0b0000_0000_0010_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_DISPLAYHEX       int = 0b0000_0000_0100_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_UINT8            int = 0b0000_0000_1000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_FLOAT            int = 0b0000_0001_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_HUEBAR           int = 0b0000_0010_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_HUEWHEEL         int = 0b0000_0100_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_INPUTRGB         int = 0b0000_1000_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_INPUTHSV         int = 0b0001_0000_0000_0000_0000_0000_0000_0000

	FLAGCOLOREDIT_DEFAULTOPTIONS int = 0b0000_1010_1001_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_DISPLAYMASK    int = 0b0000_0000_0111_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_DATATYPEMASK   int = 0b0000_0001_1000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_PICKERMASK     int = 0b0000_0110_0000_0000_0000_0000_0000_0000
	FLAGCOLOREDIT_INPUTMASK      int = 0b0001_1000_0000_0000_0000_0000_0000_0000
)

const (
	MOUSEBUTTON_LEFT   int = 0
	MOUSEBUTTON_RIGHT  int = 1
	MOUSEBUTTON_MIDDLE int = 2
)

const (
	DATEPICKERLABEL_MONTH = giu.DatePickerLabelMonth
	DATEPICKERLABEL_YEAR  = giu.DatePickerLabelYear
)

const (
	FLAGINPUTTEXT_NONE                int = 0b0000_0000_0000_0000_0000_0000
	FLAGINPUTTEXT_CHARSDECIMAL        int = 0b0000_0000_0000_0000_0000_0001
	FLAGINPUTTEXT_CHARSHEXADECIMAL    int = 0b0000_0000_0000_0000_0000_0010
	FLAGINPUTTEXT_CHARSUPPERCASE      int = 0b0000_0000_0000_0000_0000_0100
	FLAGINPUTTEXT_CHARSNOBLANK        int = 0b0000_0000_0000_0000_0000_1000
	FLAGINPUTTEXT_AUTOSELECTALL       int = 0b0000_0000_0000_0000_0001_0000
	FLAGINPUTTEXT_ENTERRETURNSTRUE    int = 0b0000_0000_0000_0000_0010_0000
	FLAGINPUTTEXT_CALLBACKCOMPLETION  int = 0b0000_0000_0000_0000_0100_0000
	FLAGINPUTTEXT_CALLBACKHISTORY     int = 0b0000_0000_0000_0000_1000_0000
	FLAGINPUTTEXT_CALLBACKALWAYS      int = 0b0000_0000_0000_0001_0000_0000
	FLAGINPUTTEXT_CALLBACKCHARFILTER  int = 0b0000_0000_0000_0010_0000_0000
	FLAGINPUTTEXT_ALLOWTABINPUT       int = 0b0000_0000_0000_0100_0000_0000
	FLAGINPUTTEXT_CTRLENTERFORNEWLINE int = 0b0000_0000_0000_1000_0000_0000
	FLAGINPUTTEXT_NOHORIZONTALSCROLL  int = 0b0000_0000_0001_0000_0000_0000
	FLAGINPUTTEXT_ALWAYSOVERWRITE     int = 0b0000_0000_0010_0000_0000_0000
	FLAGINPUTTEXT_READONLY            int = 0b0000_0000_0100_0000_0000_0000
	FLAGINPUTTEXT_PASSWORD            int = 0b0000_0000_1000_0000_0000_0000
	FLAGINPUTTEXT_NOUNDOREDO          int = 0b0000_0001_0000_0000_0000_0000
	FLAGINPUTTEXT_CHARSSCIENTIFIC     int = 0b0000_0010_0000_0000_0000_0000
	FLAGINPUTTEXT_CALLBACKRESIZE      int = 0b0000_0100_0000_0000_0000_0000
	FLAGINPUTTEXT_CALLBACKEDIT        int = 0b0000_1000_0000_0000_0000_0000
	FLAGINPUTTEXT_ESCAPECLEARSALL     int = 0b0001_0000_0000_0000_0000_0000
)

const (
	WIDGET_LABEL                = "label"
	WIDGET_BUTTON               = "button"
	WIDGET_DUMMY                = "dummy"
	WIDGET_SEPARATOR            = "separator"
	WIDGET_BULLET_TEXT          = "bullet_text"
	WIDGET_BULLET               = "bullet"
	WIDGET_CHECKBOX             = "checkbox"
	WIDGET_CHILD                = "child"
	WIDGET_COLOR_EDIT           = "color_edit"
	WIDGET_COLUMN               = "column"
	WIDGET_ROW                  = "row"
	WIDGET_COMBO_CUSTOM         = "combo_custom"
	WIDGET_COMBO                = "combo"
	WIDGET_CONDITION            = "condition"
	WIDGET_CONTEXT_MENU         = "context_menu"
	WIDGET_DATE_PICKER          = "date_picker"
	WIDGET_DRAG_INT             = "drag_int"
	WIDGET_INPUT_FLOAT          = "input_float"
	WIDGET_INPUT_INT            = "input_int"
	WIDGET_INPUT_TEXT           = "input_text"
	WIDGET_INPUT_MULTILINE_TEXT = "input_multiline_text"
	WIDGET_PROGRESS_BAR         = "progress_bar"
	WIDGET_PROGRESS_INDICATOR   = "progress_indicator"
	WIDGET_SPACING              = "spacing"
)

var buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{}

func init() {
	buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{
		WIDGET_LABEL:                labelBuild,
		WIDGET_BUTTON:               buttonBuild,
		WIDGET_DUMMY:                dummyBuild,
		WIDGET_SEPARATOR:            separatorBuild,
		WIDGET_BULLET_TEXT:          bulletTextBuild,
		WIDGET_BULLET:               bulletBuild,
		WIDGET_CHECKBOX:             checkboxBuild,
		WIDGET_CHILD:                childBuild,
		WIDGET_COLOR_EDIT:           colorEditBuild,
		WIDGET_COLUMN:               columnBuild,
		WIDGET_ROW:                  rowBuild,
		WIDGET_COMBO_CUSTOM:         comboCustomBuild,
		WIDGET_COMBO:                comboBuild,
		WIDGET_CONDITION:            conditionBuild,
		WIDGET_CONTEXT_MENU:         contextMenuBuild,
		WIDGET_DATE_PICKER:          datePickerBuild,
		WIDGET_DRAG_INT:             dragIntBuild,
		WIDGET_INPUT_FLOAT:          inputFloatBuild,
		WIDGET_INPUT_INT:            inputIntBuild,
		WIDGET_INPUT_TEXT:           inputTextBuild,
		WIDGET_INPUT_MULTILINE_TEXT: inputMultilineTextBuild,
		WIDGET_PROGRESS_BAR:         progressBarBuild,
		WIDGET_PROGRESS_INDICATOR:   progressIndicatorBuild,
		WIDGET_SPACING:              spacingBuild,
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
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
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
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
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
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
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
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
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
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ColorEditFlags(flags.(golua.LNumber)))
	}

	return c
}

func columnTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_COLUMN))
	state.SetTable(t, golua.LString("widgets"), widgets)

	return t
}

func columnBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	var widgets []g.Widget

	wid := state.GetTable(t, golua.LString("widgets"))
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable), state), state, lg), lg)
	}

	s := g.Column(widgets...)

	return s
}

func rowTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_ROW))
	state.SetTable(t, golua.LString("widgets"), widgets)

	return t
}

func rowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	var widgets []g.Widget

	wid := state.GetTable(t, golua.LString("widgets"))
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable), state), state, lg), lg)
	}

	s := g.Row(widgets...)

	return s
}

func comboCustomTable(state *golua.LState, text, preview string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_COMBO_CUSTOM))
	state.SetTable(t, golua.LString("text"), golua.LString(text))
	state.SetTable(t, golua.LString("preview"), golua.LString(preview))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	return t
}

func comboCustomBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := state.GetTable(t, golua.LString("text")).String()
	preview := state.GetTable(t, golua.LString("preview")).String()
	c := g.ComboCustom(text, preview)

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ComboFlags(flags.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func comboTable(state *golua.LState, text, preview string, items golua.LValue, i32Ref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_COMBO))
	state.SetTable(t, golua.LString("text"), golua.LString(text))
	state.SetTable(t, golua.LString("preview"), golua.LString(preview))
	state.SetTable(t, golua.LString("items"), items)
	state.SetTable(t, golua.LString("i32ref"), golua.LNumber(i32Ref))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	return t
}

func comboBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := state.GetTable(t, golua.LString("text")).String()
	preview := state.GetTable(t, golua.LString("preview")).String()

	items := []string{}
	it := state.GetTable(t, golua.LString("items")).(*golua.LTable)
	for i := range it.Len() {
		v := state.GetTable(it, golua.LNumber(i+1)).(golua.LString)
		items = append(items, string(v))
	}

	ref := int(state.GetTable(t, golua.LString("i32ref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.Combo(text, preview, items, selected)

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ComboFlags(flags.(golua.LNumber)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	return c
}

func conditionTable(state *golua.LState, condition bool, layoutIf, layoutElse golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CONDITION))
	state.SetTable(t, golua.LString("condition"), golua.LBool(condition))
	state.SetTable(t, golua.LString("layoutIf"), layoutIf)
	state.SetTable(t, golua.LString("layoutElse"), layoutElse)

	return t
}

func conditionBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	condition := state.GetTable(t, golua.LString("condition")).(golua.LBool)

	widIf := state.GetTable(t, golua.LString("layoutIf")).(*golua.LTable)
	widgetsIf := layoutBuild(r, state, []*golua.LTable{widIf}, lg)
	widElse := state.GetTable(t, golua.LString("layoutElse")).(*golua.LTable)
	widgetsElse := layoutBuild(r, state, []*golua.LTable{widElse}, lg)

	s := g.Condition(bool(condition), widgetsIf[0], widgetsElse[0])

	return s
}

func contextMenuTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CONTEXT_MENU))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__button"), golua.LNil)

	tableBuilderFunc(state, t, "mouse_button", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__button"), lt)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func contextMenuBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.ContextMenu()

	button := state.GetTable(t, golua.LString("__button"))
	if button.Type() == golua.LTNumber {
		c.MouseButton(g.MouseButton(button.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func datePickerTable(state *golua.LState, id string, timeref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_DATE_PICKER))
	state.SetTable(t, golua.LString("id"), golua.LString(id))
	state.SetTable(t, golua.LString("timeref"), golua.LNumber(timeref))
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__startofweek"), golua.LNil)
	state.SetTable(t, golua.LString("__translationMonth"), golua.LNil)
	state.SetTable(t, golua.LString("__translationYear"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "start_of_week", func(state *golua.LState, t *golua.LTable) {
		sow := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__startofweek"), sow)
	})

	tableBuilderFunc(state, t, "translation", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-2)
		value := state.CheckString(-1)

		if label == string(DATEPICKERLABEL_MONTH) {
			state.SetTable(t, golua.LString("__translationMonth"), golua.LString(value))
		} else if label == string(DATEPICKERLABEL_YEAR) {
			state.SetTable(t, golua.LString("__translationYear"), golua.LString(value))
		}
	})

	return t
}

func datePickerBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("timeref")).(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	date := sref.Value.(*time.Time)
	c := g.DatePicker(state.GetTable(t, golua.LString("id")).String(), date)

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(date.UnixMilli()))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	sow := state.GetTable(t, golua.LString("__startofweek"))
	if sow.Type() == golua.LTNumber {
		c.StartOfWeek(time.Weekday((sow.(golua.LNumber))))
	}

	translationMonth := state.GetTable(t, golua.LString("__translationMonth"))
	if translationMonth.Type() == golua.LTString {
		c.Translation(DATEPICKERLABEL_MONTH, translationMonth.String())
	}
	translationYear := state.GetTable(t, golua.LString("__translationYear"))
	if translationYear.Type() == golua.LTString {
		c.Translation(DATEPICKERLABEL_YEAR, translationYear.String())
	}

	return c
}

func dragIntTable(state *golua.LState, text string, i32Ref, minValue, maxValue int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_DRAG_INT))
	state.SetTable(t, golua.LString("text"), golua.LString(text))
	state.SetTable(t, golua.LString("i32ref"), golua.LNumber(i32Ref))
	state.SetTable(t, golua.LString("minvalue"), golua.LNumber(minValue))
	state.SetTable(t, golua.LString("maxvalue"), golua.LNumber(maxValue))
	state.SetTable(t, golua.LString("__speed"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)

	tableBuilderFunc(state, t, "speed", func(state *golua.LState, t *golua.LTable) {
		speed := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__speed"), speed)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	return t
}

func dragIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := state.GetTable(t, golua.LString("text")).String()
	min := state.GetTable(t, golua.LString("minvalue")).(golua.LNumber)
	max := state.GetTable(t, golua.LString("maxvalue")).(golua.LNumber)

	ref := int(state.GetTable(t, golua.LString("i32ref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.DragInt(text, selected, int32(min), int32(max))

	speed := state.GetTable(t, golua.LString("__speed"))
	if speed.Type() == golua.LTNumber {
		c.Speed(float32(speed.(golua.LNumber)))
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	return c
}

func inputFloatTable(state *golua.LState, floatref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_INPUT_FLOAT))
	state.SetTable(t, golua.LString("f32ref"), golua.LNumber(floatref))
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__stepsize"), golua.LNil)
	state.SetTable(t, golua.LString("__stepsizefast"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "step_size", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__stepsize"), stepsize)
	})

	tableBuilderFunc(state, t, "step_size_fast", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__stepsizefast"), stepsize)
	})

	return t
}

func inputFloatBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("f32ref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*float32)

	c := g.InputFloat(selected)

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	stepsize := state.GetTable(t, golua.LString("__stepsize"))
	if stepsize.Type() == golua.LTNumber {
		c.StepSize(float32(stepsize.(golua.LNumber)))
	}

	stepsizefast := state.GetTable(t, golua.LString("__stepsizefast"))
	if stepsizefast.Type() == golua.LTNumber {
		c.StepSizeFast(float32(stepsizefast.(golua.LNumber)))
	}

	return c
}

func inputIntTable(state *golua.LState, intref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_INPUT_INT))
	state.SetTable(t, golua.LString("i32ref"), golua.LNumber(intref))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__stepsize"), golua.LNil)
	state.SetTable(t, golua.LString("__stepsizefast"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "step_size", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__stepsize"), stepsize)
	})

	tableBuilderFunc(state, t, "step_size_fast", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__stepsizefast"), stepsize)
	})

	return t
}

func inputIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("i32ref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.InputInt(selected)

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	stepsize := state.GetTable(t, golua.LString("__stepsize"))
	if stepsize.Type() == golua.LTNumber {
		c.StepSize(int(stepsize.(golua.LNumber)))
	}

	stepsizefast := state.GetTable(t, golua.LString("__stepsizefast"))
	if stepsizefast.Type() == golua.LTNumber {
		c.StepSizeFast(int(stepsizefast.(golua.LNumber)))
	}

	return c
}

func inputTextTable(state *golua.LState, strref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_INPUT_TEXT))
	state.SetTable(t, golua.LString("strref"), golua.LNumber(strref))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__autocomplete"), golua.LNil)
	state.SetTable(t, golua.LString("__callback"), golua.LNil)
	state.SetTable(t, golua.LString("__hint"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "autocomplete", func(state *golua.LState, t *golua.LTable) {
		ac := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__autocomplete"), ac)
	})

	tableBuilderFunc(state, t, "callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__callback"), fn)
	})

	tableBuilderFunc(state, t, "hint", func(state *golua.LState, t *golua.LTable) {
		hint := state.CheckString(-1)
		state.SetTable(t, golua.LString("__hint"), golua.LString(hint))
	})

	return t
}

func inputTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("strref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*string)

	c := g.InputText(selected)

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	hint := state.GetTable(t, golua.LString("__hint"))
	if hint.Type() == golua.LTString {
		c.Hint(string(hint.(golua.LString)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	callback := state.GetTable(t, golua.LString("__callback"))
	if callback.Type() == golua.LTFunction {
		c.Callback(func(data imgui.InputTextCallbackData) int {
			state.Push(callback)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
			return 0
		})
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	ac := state.GetTable(t, golua.LString("__autocomplete"))
	if ac.Type() == golua.LTTable {
		acList := []string{}
		at := ac.(*golua.LTable)
		for i := range at.Len() {
			ai := state.GetTable(at, golua.LNumber(i+1)).(golua.LString)
			acList = append(acList, string(ai))
		}

		c.AutoComplete(acList)
	}

	return c
}

func inputMultilineTextTable(state *golua.LState, strref int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_INPUT_MULTILINE_TEXT))
	state.SetTable(t, golua.LString("strref"), golua.LNumber(strref))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__callback"), golua.LNil)
	state.SetTable(t, golua.LString("__autoscroll"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__callback"), fn)
	})

	tableBuilderFunc(state, t, "autoscroll_to_bottom", func(state *golua.LState, t *golua.LTable) {
		as := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__autoscroll"), golua.LBool(as))
	})

	return t
}

func inputMultilineTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(state.GetTable(t, golua.LString("strref")).(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*string)

	c := g.InputTextMultiline(selected)

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	callback := state.GetTable(t, golua.LString("__callback"))
	if callback.Type() == golua.LTFunction {
		c.Callback(func(data imgui.InputTextCallbackData) int {
			state.Push(callback)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
			return 0
		})
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	autoscroll := state.GetTable(t, golua.LString("__autoscroll"))
	if autoscroll.Type() == golua.LTBool {
		c.AutoScrollToBottom(bool(autoscroll.(golua.LBool)))
	}

	return c
}

func progressBarTable(state *golua.LState, fraction float64) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_PROGRESS_BAR))
	state.SetTable(t, golua.LString("fraction"), golua.LNumber(fraction))
	state.SetTable(t, golua.LString("__overlay"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "overlay", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__overlay"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	return t
}

func progressBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	fraction := state.GetTable(t, golua.LString("fraction")).(golua.LNumber)
	p := g.ProgressBar(float32(fraction))

	overlay := state.GetTable(t, golua.LString("__overlay"))
	if overlay.Type() == golua.LTString {
		p.Overlay(string(overlay.(golua.LString)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		p.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return p
}

func progressIndicatorTable(state *golua.LState, label string, width, height, radius float64) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_PROGRESS_INDICATOR))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("width"), golua.LNumber(width))
	state.SetTable(t, golua.LString("height"), golua.LNumber(height))
	state.SetTable(t, golua.LString("radius"), golua.LNumber(radius))

	return t
}

func progressIndicatorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	width := state.GetTable(t, golua.LString("width")).(golua.LNumber)
	height := state.GetTable(t, golua.LString("height")).(golua.LNumber)
	radius := state.GetTable(t, golua.LString("radius")).(golua.LNumber)
	p := g.ProgressIndicator(string(label), float32(width), float32(height), float32(radius))

	return p
}

func spacingTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SPACING))

	return t
}

func spacingBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Spacing()

	return b
}
