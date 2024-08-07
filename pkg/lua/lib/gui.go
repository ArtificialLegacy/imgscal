package lib

import (
	"fmt"
	"image"
	"image/color"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
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

	/// @func window_master()
	/// @arg name
	/// @arg width
	/// @arg height
	/// @arg? flags
	/// @returns id of the window.
	lib.CreateFunction(tab, "window_master",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.INT, Name: "flags", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w := g.NewMasterWindow(args["name"].(string), args["width"].(int), args["height"].(int), g.MasterWindowFlags(args["flags"].(int)))
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

	/// @func window_single()
	/// @returns window widget
	lib.CreateFunction(tab, "window_single",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, true, false, "")

			state.Push(win)
			return 1
		})

	/// @func window_single_with_menu_bar()
	/// @returns window widget
	lib.CreateFunction(tab, "window_single_with_menu_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, true, true, "")

			state.Push(win)
			return 1
		})

	/// @func window()
	/// @arg title
	/// @returns window widget
	lib.CreateFunction(tab, "window",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, false, false, args["title"].(string))

			state.Push(win)
			return 1
		})

	/// @func layout()
	/// @arg widgets - []Widgets
	/// @desc
	/// Builds a list of widgets in place.
	lib.CreateFunction(tab, "layout",
		[]lua.Arg{
			{Type: lua.ANY, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			widgets := args["widgets"].(*golua.LTable)
			layout := g.Layout(layoutBuild(r, state, parseWidgets(parseTable(widgets, state), state, lg), lg))
			layout.Build()

			return 0
		})

	/// @func popup_open()
	/// @arg name
	lib.CreateFunction(tab, "popup_open",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			name := args["name"].(string)

			g.OpenPopup(name)
			return 0
		})

	/// @func popup_close()
	lib.CreateFunction(tab, "popup_close",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.CloseCurrentPopup()
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

	/// @func wg_number()
	/// @arg number
	/// @returns widget
	/// @desc
	/// A float->string wrapper around wg_label
	lib.CreateFunction(tab, "wg_number",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "number"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := labelTable(state, fmt.Sprintf("%v", args["number"].(float64)))

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
	/// @returns widget
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
	/// @returns widget
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

	/// @func wg_button_small()
	/// @arg text
	/// @returns widget
	lib.CreateFunction(tab, "wg_button_small",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonSmallTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_button_radio()
	/// @arg text
	/// @arg active
	/// @returns widget
	lib.CreateFunction(tab, "wg_button_radio",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.BOOL, Name: "active"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonRadioTable(state, args["text"].(string), args["active"].(bool))

			state.Push(t)
			return 1
		})

	/// @func wg_image_url()
	/// @arg url
	/// @returns widget
	lib.CreateFunction(tab, "wg_image_url",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageUrlTable(state, args["url"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_image()
	/// @arg id
	/// @returns widget
	/// @blocking
	lib.CreateFunction(tab, "wg_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageTable(state, args["id"].(int), false)

			state.Push(t)
			return 1
		})

	/// @func wg_image_sync()
	/// @arg id
	/// @returns widget
	/// @desc
	/// Note: this does not wait for the image to be ready or idle,
	/// if the image is not loaded it will dislay an empy image
	/// May look weird if the image is also being processed while displayed here.
	lib.CreateFunction(tab, "wg_image_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageTable(state, args["id"].(int), true)

			state.Push(t)
			return 1
		})

	/// @func wg_list_box()
	/// @arg items
	/// @returns widget
	lib.CreateFunction(tab, "wg_list_box",
		[]lua.Arg{
			{Type: lua.ANY, Name: "items"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := listBoxTable(state, args["items"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_list_clipper()
	/// @returns widget
	lib.CreateFunction(tab, "wg_list_clipper",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := listClipperTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_bar_main()
	/// @returns widget
	lib.CreateFunction(tab, "wg_menu_bar_main",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := mainMenuBarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_bar()
	/// @returns widget
	lib.CreateFunction(tab, "wg_menu_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuBarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_item()
	/// @arg label
	/// @returns widget
	lib.CreateFunction(tab, "wg_menu_item",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuItemTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_menu()
	/// @arg label
	/// @returns widget
	lib.CreateFunction(tab, "wg_menu",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_selectable()
	/// @arg label
	/// @returns widget
	lib.CreateFunction(tab, "wg_selectable",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := selectableTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_slider_float()
	/// @arg f32ref
	/// @arg min
	/// @arg max
	/// @returns widget
	lib.CreateFunction(tab, "wg_slider_float",
		[]lua.Arg{
			{Type: lua.INT, Name: "f32ref"},
			{Type: lua.FLOAT, Name: "min"},
			{Type: lua.FLOAT, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			floatref := args["f32ref"].(int)
			min := args["min"].(float64)
			max := args["max"].(float64)
			t := sliderFloatTable(state, floatref, min, max)

			state.Push(t)
			return 1
		})

	/// @func wg_slider_int()
	/// @arg i32ref
	/// @arg min
	/// @arg max
	/// @returns widget
	lib.CreateFunction(tab, "wg_slider_int",
		[]lua.Arg{
			{Type: lua.INT, Name: "i32ref"},
			{Type: lua.INT, Name: "min"},
			{Type: lua.INT, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			intref := args["i32ref"].(int)
			min := args["min"].(int)
			max := args["max"].(int)
			t := sliderIntTable(state, intref, min, max)

			state.Push(t)
			return 1
		})

	/// @func wg_vslider_int()
	/// @arg i32ref
	/// @arg min
	/// @arg max
	/// @returns widget
	lib.CreateFunction(tab, "wg_vslider_int",
		[]lua.Arg{
			{Type: lua.INT, Name: "i32ref"},
			{Type: lua.INT, Name: "min"},
			{Type: lua.INT, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			intref := args["i32ref"].(int)
			min := args["min"].(int)
			max := args["max"].(int)
			t := vsliderIntTable(state, intref, min, max)

			state.Push(t)
			return 1
		})

	/// @func wg_tab_bar()
	/// @returns widget
	lib.CreateFunction(tab, "wg_tab_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tabbarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_tab_item()
	/// @arg label
	/// @returns tab item
	lib.CreateFunction(tab, "wg_tab_item",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tabitemTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_tooltip()
	/// @arg tip
	/// @returns widget
	lib.CreateFunction(tab, "wg_tooltip",
		[]lua.Arg{
			{Type: lua.STRING, Name: "tip"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tooltipTable(state, args["tip"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_table_column()
	/// @arg label
	/// @returns table column
	lib.CreateFunction(tab, "wg_table_column",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tableColumnTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_table_row()
	/// @arg? widgets - []Widgets
	/// @returns table row
	lib.CreateFunction(tab, "wg_table_row",
		[]lua.Arg{
			{Type: lua.ANY, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["widgets"]
			if v == nil {
				v = golua.LNil
			}
			t := tableRowTable(state, v.(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_table()
	/// @returns widget
	lib.CreateFunction(tab, "wg_table",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tableTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_button_arrow()
	/// @arg dir
	/// @returns widget
	lib.CreateFunction(tab, "wg_button_arrow",
		[]lua.Arg{
			{Type: lua.INT, Name: "dir"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonArrowTable(state, args["dir"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_tree_table_node()
	/// @arg label
	/// @returns widget
	lib.CreateFunction(tab, "wg_tree_table_node",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := treeTableNodeTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	// @func wg_tree_table_row()
	/// @arg label
	/// @arg? widgets - []Widgets
	/// @returns tree table row
	lib.CreateFunction(tab, "wg_tree_table_row",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
			{Type: lua.ANY, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["widgets"]
			if v == nil {
				v = golua.LNil
			}
			t := treeTableRowTable(state, args["label"].(string), v.(golua.LValue))

			state.Push(t)
			return 1
		})

	// @func wg_tree_table()
	/// @returns widget
	lib.CreateFunction(tab, "wg_tree_table",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := treeTableTable(state)

			state.Push(t)
			return 1
		})

	// @func wg_popup_modal()
	/// @arg name
	/// @returns widget
	lib.CreateFunction(tab, "wg_popup_modal",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := popupModalTable(state, args["name"].(string))

			state.Push(t)
			return 1
		})

	// @func wg_popup()
	/// @arg name
	/// @returns widget
	lib.CreateFunction(tab, "wg_popup",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := popupTable(state, args["name"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_layout_split()
	/// @arg direction
	/// @arg f32ref
	/// @arg layout1 - []Widgets
	/// @arg layout2 - []Widgets
	/// @returns widget
	lib.CreateFunction(tab, "wg_layout_split",
		[]lua.Arg{
			{Type: lua.INT, Name: "direction"},
			{Type: lua.INT, Name: "f32ref"},
			{Type: lua.ANY, Name: "layout1"},
			{Type: lua.ANY, Name: "layout2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			layout1 := args["layout1"].(golua.LValue)
			layout2 := args["layout2"].(golua.LValue)
			t := splitLayoutTable(state, args["direction"].(int), args["f32ref"].(int), layout1, layout2)

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

	/// @constants Selectable Flags
	/// @const FLAGSELECTABLE_NONE
	/// @const FLAGSELECTABLE_DONTCLOSEPOPUPS
	/// @const FLAGSELECTABLE_SPANALLCOLUMNS
	/// @const FLAGSELECTABLE_ALLOWDOUBLECLICK
	/// @const FLAGSELECTABLE_DISABLED
	/// @const FLAGSELECTABLE_ALLOWOVERLAP
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_NONE"), golua.LNumber(FLAGSELECTABLE_NONE))
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_DONTCLOSEPOPUPS"), golua.LNumber(FLAGSELECTABLE_DONTCLOSEPOPUPS))
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_SPANALLCOLUMNS"), golua.LNumber(FLAGSELECTABLE_SPANALLCOLUMNS))
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_ALLOWDOUBLECLICK"), golua.LNumber(FLAGSELECTABLE_ALLOWDOUBLECLICK))
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_DISABLED"), golua.LNumber(FLAGSELECTABLE_DISABLED))
	r.State.SetTable(tab, golua.LString("FLAGSELECTABLE_ALLOWOVERLAP"), golua.LNumber(FLAGSELECTABLE_ALLOWOVERLAP))

	/// @constants Slider Flags
	/// @const FLAGSLIDER_NONE
	/// @const FLAGSLIDER_ALWAYSCLAMP
	/// @const FLAGSLIDER_LOGARITHMIC
	/// @const FLAGSLIDER_NOROUNDTOFORMAT
	/// @const FLAGSLIDER_NOINPUT
	/// @const FLAGSLIDER_INVALIDMASK
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_NONE"), golua.LNumber(FLAGSLIDER_NONE))
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_ALWAYSCLAMP"), golua.LNumber(FLAGSLIDER_ALWAYSCLAMP))
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_LOGARITHMIC"), golua.LNumber(FLAGSLIDER_LOGARITHMIC))
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_NOROUNDTOFORMAT"), golua.LNumber(FLAGSLIDER_NOROUNDTOFORMAT))
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_NOINPUT"), golua.LNumber(FLAGSLIDER_NOINPUT))
	r.State.SetTable(tab, golua.LString("FLAGSLIDER_INVALIDMASK"), golua.LNumber(FLAGSLIDER_INVALIDMASK))

	/// @constants Tab Bar Flags
	/// @const FLAGTABBAR_NONE
	/// @const FLAGTABBAR_REORDERABLE
	/// @const FLAGTABBAR_AUTOSELECTNEWTABS
	/// @const FLAGTABBAR_TABLLISTPOPUPBUTTON
	/// @const FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON
	/// @const FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS
	/// @const FLAGTABBAR_NOTOOLTIP
	/// @const FLAGTABBAR_FITTINGPOLICYRESIZEDOWN
	/// @const FLAGTABBAR_FITTINGPOLICYSCROLL
	/// @const FLAGTABBAR_FITTINGPOLICYMASK
	/// @const FLAGTABBAR_FITTINGPOLICYDEFAULT
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_NONE"), golua.LNumber(FLAGTABBAR_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_REORDERABLE"), golua.LNumber(FLAGTABBAR_REORDERABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_AUTOSELECTNEWTABS"), golua.LNumber(FLAGTABBAR_AUTOSELECTNEWTABS))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_TABLLISTPOPUPBUTTON"), golua.LNumber(FLAGTABBAR_TABLLISTPOPUPBUTTON))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON"), golua.LNumber(FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS"), golua.LNumber(FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_NOTOOLTIP"), golua.LNumber(FLAGTABBAR_NOTOOLTIP))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_FITTINGPOLICYRESIZEDOWN"), golua.LNumber(FLAGTABBAR_FITTINGPOLICYRESIZEDOWN))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_FITTINGPOLICYSCROLL"), golua.LNumber(FLAGTABBAR_FITTINGPOLICYSCROLL))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_FITTINGPOLICYMASK"), golua.LNumber(FLAGTABBAR_FITTINGPOLICYMASK))
	r.State.SetTable(tab, golua.LString("FLAGTABBAR_FITTINGPOLICYDEFAULT"), golua.LNumber(FLAGTABBAR_FITTINGPOLICYDEFAULT))

	/// @constants Tab Item Flags
	/// @const FLAGTABITEM_NONE
	/// @const FLAGTABITEM_UNSAVEDOCUMENT
	/// @const FLAGTABITEM_SETSELECTED
	/// @const FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON
	/// @const FLAGTABITEM_NOPUSHID
	/// @const FLAGTABITEM_NOTOOLTIP
	/// @const FLAGTABITEM_NOREORDER
	/// @const FLAGTABITEM_LEADING
	/// @const FLAGTABITEM_TRAILING
	/// @const FLAGTABITEM_NOASSUMEDCLOSURE
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NONE"), golua.LNumber(FLAGTABITEM_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_UNSAVEDOCUMENT"), golua.LNumber(FLAGTABITEM_UNSAVEDOCUMENT))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_SETSELECTED"), golua.LNumber(FLAGTABITEM_SETSELECTED))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON"), golua.LNumber(FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NOPUSHID"), golua.LNumber(FLAGTABITEM_NOPUSHID))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NOTOOLTIP"), golua.LNumber(FLAGTABITEM_NOTOOLTIP))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NOREORDER"), golua.LNumber(FLAGTABITEM_NOREORDER))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_LEADING"), golua.LNumber(FLAGTABITEM_LEADING))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_TRAILING"), golua.LNumber(FLAGTABITEM_TRAILING))
	r.State.SetTable(tab, golua.LString("FLAGTABITEM_NOASSUMEDCLOSURE"), golua.LNumber(FLAGTABITEM_NOASSUMEDCLOSURE))

	/// @constants Table Column Flags
	/// @const FLAGTABLECOLUMN_NONE
	/// @const FLAGTABLECOLUMN_DEFAULTHIDE
	/// @const FLAGTABLECOLUMN_DEFAULTSORT
	/// @const FLAGTABLECOLUMN_WIDTHSTRETCH
	/// @const FLAGTABLECOLUMN_WIDTHFIXED
	/// @const FLAGTABLECOLUMN_NORESIZE
	/// @const FLAGTABLECOLUMN_NOREORDER
	/// @const FLAGTABLECOLUMN_NOHIDE
	/// @const FLAGTABLECOLUMN_NOCLIP
	/// @const FLAGTABLECOLUMN_NOSORT
	/// @const FLAGTABLECOLUMN_NOSORTASCENDING
	/// @const FLAGTABLECOLUMN_NOSORTDESCENDING
	/// @const FLAGTABLECOLUMN_NOHEADERWIDTH
	/// @const FLAGTABLECOLUMN_PREFERSORTASCENDING
	/// @const FLAGTABLECOLUMN_PREFERSORTDESCENDING
	/// @const FLAGTABLECOLUMN_INDENTENABLE
	/// @const FLAGTABLECOLUMN_INDENTDISABLE
	/// @const FLAGTABLECOLUMN_ISENABLED
	/// @const FLAGTABLECOLUMN_ISVISIBLE
	/// @const FLAGTABLECOLUMN_ISSORTED
	/// @const FLAGTABLECOLUMN_ISHOVERED
	/// @const FLAGTABLECOLUMN_WIDTHMASK
	/// @const FLAGTABLECOLUMN_INDENTMASK
	/// @const FLAGTABLECOLUMN_STATUSMASK
	/// @const FLAGTABLECOLUMN_NODIRECTRESIZE
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NONE"), golua.LNumber(FLAGTABLECOLUMN_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_DEFAULTHIDE"), golua.LNumber(FLAGTABLECOLUMN_DEFAULTHIDE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_DEFAULTSORT"), golua.LNumber(FLAGTABLECOLUMN_DEFAULTSORT))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_WIDTHSTRETCH"), golua.LNumber(FLAGTABLECOLUMN_WIDTHSTRETCH))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_WIDTHFIXED"), golua.LNumber(FLAGTABLECOLUMN_WIDTHFIXED))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NORESIZE"), golua.LNumber(FLAGTABLECOLUMN_NORESIZE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOREORDER"), golua.LNumber(FLAGTABLECOLUMN_NOREORDER))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOHIDE"), golua.LNumber(FLAGTABLECOLUMN_NOHIDE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOCLIP"), golua.LNumber(FLAGTABLECOLUMN_NOCLIP))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOSORT"), golua.LNumber(FLAGTABLECOLUMN_NOSORT))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOSORTASCENDING"), golua.LNumber(FLAGTABLECOLUMN_NOSORTASCENDING))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOSORTDESCENDING"), golua.LNumber(FLAGTABLECOLUMN_NOSORTDESCENDING))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NOHEADERWIDTH"), golua.LNumber(FLAGTABLECOLUMN_NOHEADERWIDTH))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_PREFERSORTASCENDING"), golua.LNumber(FLAGTABLECOLUMN_PREFERSORTASCENDING))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_PREFERSORTDESCENDING"), golua.LNumber(FLAGTABLECOLUMN_PREFERSORTDESCENDING))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_INDENTENABLE"), golua.LNumber(FLAGTABLECOLUMN_INDENTENABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_INDENTDISABLE"), golua.LNumber(FLAGTABLECOLUMN_INDENTDISABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_ISENABLED"), golua.LNumber(FLAGTABLECOLUMN_ISENABLED))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_ISVISIBLE"), golua.LNumber(FLAGTABLECOLUMN_ISVISIBLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_ISSORTED"), golua.LNumber(FLAGTABLECOLUMN_ISSORTED))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_ISHOVERED"), golua.LNumber(FLAGTABLECOLUMN_ISHOVERED))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_WIDTHMASK"), golua.LNumber(FLAGTABLECOLUMN_WIDTHMASK))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_INDENTMASK"), golua.LNumber(FLAGTABLECOLUMN_INDENTMASK))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_STATUSMASK"), golua.LNumber(FLAGTABLECOLUMN_STATUSMASK))
	r.State.SetTable(tab, golua.LString("FLAGTABLECOLUMN_NODIRECTRESIZE"), golua.LNumber(FLAGTABLECOLUMN_NODIRECTRESIZE))

	/// @constants Table Row Flags
	/// @const FLAGTABLEROW_NONE
	/// @const FLAGTABLEROW_HEADERS
	r.State.SetTable(tab, golua.LString("FLAGTABLEROW_NONE"), golua.LNumber(FLAGTABLEROW_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTABLEROW_HEADERS"), golua.LNumber(FLAGTABLEROW_HEADERS))

	/// @constants Table Flags
	/// @const FLAGTABLE_NONE
	/// @const FLAGTABLE_RESIZEABLE
	/// @const FLAGTABLE_REORDERABLE
	/// @const FLAGTABLE_HIDEABLE
	/// @const FLAGTABLE_SORTABLE
	/// @const FLAGTABLE_NOSAVEDSETTINGS
	/// @const FLAGTABLE_CONTEXTMENUINBODY
	/// @const FLAGTABLE_ROWBG
	/// @const FLAGTABLE_BORDERSINNERH
	/// @const FLAGTABLE_BORDERSOUTERH
	/// @const FLAGTABLE_BORDERSINNERV
	/// @const FLAGTABLE_BORDERSOUTERV
	/// @const FLAGTABLE_BORDERSH
	/// @const FLAGTABLE_BORDERSV
	/// @const FLAGTABLE_BORDERSINNER
	/// @const FLAGTABLE_BORDERSOUTER
	/// @const FLAGTABLE_BORDERS
	/// @const FLAGTABLE_NOBORDERSINBODY
	/// @const FLAGTABLE_NOBORDERSINBODYUNTILRESIZE
	/// @const FLAGTABLE_SIZINGFIXEDFIT
	/// @const FLAGTABLE_SIZINGFIXEDSAME
	/// @const FLAGTABLE_SIZINGSTRETCHPROP
	/// @const FLAGTABLE_SIZINGSTRETCHSAME
	/// @const FLAGTABLE_NOHOSTEXTENDX
	/// @const FLAGTABLE_NOHOSTEXTENDY
	/// @const FLAGTABLE_NOKEEPCOLUMNSVISIBLE
	/// @const FLAGTABLE_PRECISEWIDTHS
	/// @const FLAGTABLE_NOCLIP
	/// @const FLAGTABLE_PADOUTERX
	/// @const FLAGTABLE_NOPADOUTERX
	/// @const FLAGTABLE_NOPADINNERX
	/// @const FLAGTABLE_SCROLLX
	/// @const FLAGTABLE_SCROLLY
	/// @const FLAGTABLE_SORTMULTI
	/// @const FLAGTABLE_SORTTRISTATE
	/// @const FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN
	/// @const FLAGTABLE_SIZINGMASK
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NONE"), golua.LNumber(FLAGTABLE_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_RESIZEABLE"), golua.LNumber(FLAGTABLE_RESIZEABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_REORDERABLE"), golua.LNumber(FLAGTABLE_REORDERABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_HIDEABLE"), golua.LNumber(FLAGTABLE_HIDEABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SORTABLE"), golua.LNumber(FLAGTABLE_SORTABLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOSAVEDSETTINGS"), golua.LNumber(FLAGTABLE_NOSAVEDSETTINGS))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_CONTEXTMENUINBODY"), golua.LNumber(FLAGTABLE_CONTEXTMENUINBODY))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_ROWBG"), golua.LNumber(FLAGTABLE_ROWBG))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSINNERH"), golua.LNumber(FLAGTABLE_BORDERSINNERH))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSOUTERH"), golua.LNumber(FLAGTABLE_BORDERSOUTERH))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSINNERV"), golua.LNumber(FLAGTABLE_BORDERSINNERV))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSOUTERV"), golua.LNumber(FLAGTABLE_BORDERSOUTERV))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSH"), golua.LNumber(FLAGTABLE_BORDERSH))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSV"), golua.LNumber(FLAGTABLE_BORDERSV))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSINNER"), golua.LNumber(FLAGTABLE_BORDERSINNER))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERSOUTER"), golua.LNumber(FLAGTABLE_BORDERSOUTER))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_BORDERS"), golua.LNumber(FLAGTABLE_BORDERS))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOBORDERSINBODY"), golua.LNumber(FLAGTABLE_NOBORDERSINBODY))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOBORDERSINBODYUNTILRESIZE"), golua.LNumber(FLAGTABLE_NOBORDERSINBODYUNTILRESIZE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SIZINGFIXEDFIT"), golua.LNumber(FLAGTABLE_SIZINGFIXEDFIT))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SIZINGFIXEDSAME"), golua.LNumber(FLAGTABLE_SIZINGFIXEDSAME))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SIZINGSTRETCHPROP"), golua.LNumber(FLAGTABLE_SIZINGSTRETCHPROP))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SIZINGSTRETCHSAME"), golua.LNumber(FLAGTABLE_SIZINGSTRETCHSAME))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOHOSTEXTENDX"), golua.LNumber(FLAGTABLE_NOHOSTEXTENDX))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOHOSTEXTENDY"), golua.LNumber(FLAGTABLE_NOHOSTEXTENDY))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOKEEPCOLUMNSVISIBLE"), golua.LNumber(FLAGTABLE_NOKEEPCOLUMNSVISIBLE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_PRECISEWIDTHS"), golua.LNumber(FLAGTABLE_PRECISEWIDTHS))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOCLIP"), golua.LNumber(FLAGTABLE_NOCLIP))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_PADOUTERX"), golua.LNumber(FLAGTABLE_PADOUTERX))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOPADOUTERX"), golua.LNumber(FLAGTABLE_NOPADOUTERX))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_NOPADINNERX"), golua.LNumber(FLAGTABLE_NOPADINNERX))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SCROLLX"), golua.LNumber(FLAGTABLE_SCROLLX))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SCROLLY"), golua.LNumber(FLAGTABLE_SCROLLY))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SORTMULTI"), golua.LNumber(FLAGTABLE_SORTMULTI))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SORTTRISTATE"), golua.LNumber(FLAGTABLE_SORTTRISTATE))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN"), golua.LNumber(FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN))
	r.State.SetTable(tab, golua.LString("FLAGTABLE_SIZINGMASK"), golua.LNumber(FLAGTABLE_SIZINGMASK))

	/// @constants Directions
	/// @const DIR_NONE
	/// @const DIR_LEFT
	/// @const DIR_RIGHT
	/// @const DIR_UP
	/// @const DIR_DOWN
	/// @const DIR_COUNT
	r.State.SetTable(tab, golua.LString("DIR_NONE"), golua.LNumber(DIR_NONE))
	r.State.SetTable(tab, golua.LString("DIR_LEFT"), golua.LNumber(DIR_LEFT))
	r.State.SetTable(tab, golua.LString("DIR_RIGHT"), golua.LNumber(DIR_RIGHT))
	r.State.SetTable(tab, golua.LString("DIR_UP"), golua.LNumber(DIR_UP))
	r.State.SetTable(tab, golua.LString("DIR_DOWN"), golua.LNumber(DIR_DOWN))
	r.State.SetTable(tab, golua.LString("DIR_COUNT"), golua.LNumber(DIR_COUNT))

	/// @constants Tree Node Flags
	/// @const FLAGTREENODE_NONE
	/// @const FLAGTREENODE_SELECTED
	/// @const FLAGTREENODE_FRAMED
	/// @const FLAGTREENODE_ALLOWOVERLAP
	/// @const FLAGTREENODE_NOTREEPUSHONOPEN
	/// @const FLAGTREENODE_NOAUTOOPENONLOG
	/// @const FLAGTREENODE_DEFAULTOPEN
	/// @const FLAGTREENODE_OPENONDOUBLECLICK
	/// @const FLAGTREENODE_OPENONARROW
	/// @const FLAGTREENODE_LEAF
	/// @const FLAGTREENODE_BULLET
	/// @const FLAGTREENODE_FRAMEPADDING
	/// @const FLAGTREENODE_SPANAVAILWIDTH
	/// @const FLAGTREENODE_SPANFULLWIDTH
	/// @const FLAGTREENODE_SPANALLCOLUMNS
	/// @const FLAGTREENODE_NAVLEFTJUMPSBACKHERE
	/// @const FLAGTREENODE_COLLAPSINGHEADER
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_NONE"), golua.LNumber(FLAGTREENODE_NONE))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_SELECTED"), golua.LNumber(FLAGTREENODE_SELECTED))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_FRAMED"), golua.LNumber(FLAGTREENODE_FRAMED))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_ALLOWOVERLAP"), golua.LNumber(FLAGTREENODE_ALLOWOVERLAP))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_NOTREEPUSHONOPEN"), golua.LNumber(FLAGTREENODE_NOTREEPUSHONOPEN))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_NOAUTOOPENONLOG"), golua.LNumber(FLAGTREENODE_NOAUTOOPENONLOG))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_DEFAULTOPEN"), golua.LNumber(FLAGTREENODE_DEFAULTOPEN))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_OPENONDOUBLECLICK"), golua.LNumber(FLAGTREENODE_OPENONDOUBLECLICK))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_OPENONARROW"), golua.LNumber(FLAGTREENODE_OPENONARROW))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_LEAF"), golua.LNumber(FLAGTREENODE_LEAF))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_BULLET"), golua.LNumber(FLAGTREENODE_BULLET))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_FRAMEPADDING"), golua.LNumber(FLAGTREENODE_FRAMEPADDING))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_SPANAVAILWIDTH"), golua.LNumber(FLAGTREENODE_SPANAVAILWIDTH))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_SPANFULLWIDTH"), golua.LNumber(FLAGTREENODE_SPANFULLWIDTH))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_SPANALLCOLUMNS"), golua.LNumber(FLAGTREENODE_SPANALLCOLUMNS))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_NAVLEFTJUMPSBACKHERE"), golua.LNumber(FLAGTREENODE_NAVLEFTJUMPSBACKHERE))
	r.State.SetTable(tab, golua.LString("FLAGTREENODE_COLLAPSINGHEADER"), golua.LNumber(FLAGTREENODE_COLLAPSINGHEADER))

	/// @constants Master Window Flags
	/// @const FLAGMASTERWINDOW_NOTRESIZABLE
	/// @const FLAGMASTERWINDOW_MAXIMIZED
	/// @const FLAGMASTERWINDOW_FLOATING
	/// @const FLAGMASTERWINDOW_FRAMELESS
	/// @const FLAGMASTERWINDOW_TRANSPARENT
	r.State.SetTable(tab, golua.LString("FLAGMASTERWINDOW_NOTRESIZABLE"), golua.LNumber(FLAGMASTERWINDOW_NOTRESIZABLE))
	r.State.SetTable(tab, golua.LString("FLAGMASTERWINDOW_MAXIMIZED"), golua.LNumber(FLAGMASTERWINDOW_MAXIMIZED))
	r.State.SetTable(tab, golua.LString("FLAGMASTERWINDOW_FLOATING"), golua.LNumber(FLAGMASTERWINDOW_FLOATING))
	r.State.SetTable(tab, golua.LString("FLAGMASTERWINDOW_FRAMELESS"), golua.LNumber(FLAGMASTERWINDOW_FRAMELESS))
	r.State.SetTable(tab, golua.LString("FLAGMASTERWINDOW_TRANSPARENT"), golua.LNumber(FLAGMASTERWINDOW_TRANSPARENT))

	/// @constants Window Flags
	/// @const FLAGWINDOW_NONE
	/// @const FLAGWINDOW_NOTITLEBAR
	/// @const FLAGWINDOW_NORESIZE
	/// @const FLAGWINDOW_NOMOVE
	/// @const FLAGWINDOW_NOSCROLLBAR
	/// @const FLAGWINDOW_NOSCROLLWITHMOUSE
	/// @const FLAGWINDOW_NOCOLLAPSE
	/// @const FLAGWINDOW_ALWAYSAUTORESIZE
	/// @const FLAGWINDOW_NOBACKGROUND
	/// @const FLAGWINDOW_NOSAVEDSETTINGS
	/// @const FLAGWINDOW_NOMOUSEINPUTS
	/// @const FLAGWINDOW_MENUBAR
	/// @const FLAGWINDOW_HORIZONTALSCROLLBAR
	/// @const FLAGWINDOW_NOFOCUSONAPPEARING
	/// @const FLAGWINDOW_NOBRINGTOFRONTONFOCUS
	/// @const FLAGWINDOW_ALWAYSVERTICALSCROLLBAR
	/// @const FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR
	/// @const FLAGWINDOW_NONAVINPUTS
	/// @const FLAGWINDOW_NONAVFOCUS
	/// @const FLAGWINDOW_UNSAVEDDOCUMENT
	/// @const FLAGWINDOW_NONAV
	/// @const FLAGWINDOW_NODECORATION
	/// @const FLAGWINDOW_NOINPUTS
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NONE"), golua.LNumber(FLAGWINDOW_NONE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOTITLEBAR"), golua.LNumber(FLAGWINDOW_NOTITLEBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NORESIZE"), golua.LNumber(FLAGWINDOW_NORESIZE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOMOVE"), golua.LNumber(FLAGWINDOW_NOMOVE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOSCROLLBAR"), golua.LNumber(FLAGWINDOW_NOSCROLLBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOSCROLLWITHMOUSE"), golua.LNumber(FLAGWINDOW_NOSCROLLWITHMOUSE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOCOLLAPSE"), golua.LNumber(FLAGWINDOW_NOCOLLAPSE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_ALWAYSAUTORESIZE"), golua.LNumber(FLAGWINDOW_ALWAYSAUTORESIZE))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOBACKGROUND"), golua.LNumber(FLAGWINDOW_NOBACKGROUND))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOSAVEDSETTINGS"), golua.LNumber(FLAGWINDOW_NOSAVEDSETTINGS))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOMOUSEINPUTS"), golua.LNumber(FLAGWINDOW_NOMOUSEINPUTS))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_MENUBAR"), golua.LNumber(FLAGWINDOW_MENUBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_HORIZONTALSCROLLBAR"), golua.LNumber(FLAGWINDOW_HORIZONTALSCROLLBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOFOCUSONAPPEARING"), golua.LNumber(FLAGWINDOW_NOFOCUSONAPPEARING))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOBRINGTOFRONTONFOCUS"), golua.LNumber(FLAGWINDOW_NOBRINGTOFRONTONFOCUS))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_ALWAYSVERTICALSCROLLBAR"), golua.LNumber(FLAGWINDOW_ALWAYSVERTICALSCROLLBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR"), golua.LNumber(FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NONAVINPUTS"), golua.LNumber(FLAGWINDOW_NONAVINPUTS))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NONAVFOCUS"), golua.LNumber(FLAGWINDOW_NONAVFOCUS))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_UNSAVEDDOCUMENT"), golua.LNumber(FLAGWINDOW_UNSAVEDDOCUMENT))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NONAV"), golua.LNumber(FLAGWINDOW_NONAV))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NODECORATION"), golua.LNumber(FLAGWINDOW_NODECORATION))
	r.State.SetTable(tab, golua.LString("FLAGWINDOW_NOINPUTS"), golua.LNumber(FLAGWINDOW_NOINPUTS))

	/// @constants Split Direction
	/// @const SPLITDIRECTION_HORIZONTAL
	/// @const SPLITDIRECTION_VERTICAL
	r.State.SetTable(tab, golua.LString("SPLITDIRECTION_HORIZONTAL"), golua.LNumber(SPLITDIRECTION_HORIZONTAL))
	r.State.SetTable(tab, golua.LString("SPLITDIRECTION_VERTICAL"), golua.LNumber(SPLITDIRECTION_VERTICAL))
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
	DATEPICKERLABEL_MONTH = g.DatePickerLabelMonth
	DATEPICKERLABEL_YEAR  = g.DatePickerLabelYear
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
	FLAGSELECTABLE_NONE             int = 0b0000_0000
	FLAGSELECTABLE_DONTCLOSEPOPUPS  int = 0b0000_0001
	FLAGSELECTABLE_SPANALLCOLUMNS   int = 0b0000_0010
	FLAGSELECTABLE_ALLOWDOUBLECLICK int = 0b0000_0100
	FLAGSELECTABLE_DISABLED         int = 0b0000_1000
	FLAGSELECTABLE_ALLOWOVERLAP     int = 0b0001_0000
)

const (
	FLAGSLIDER_NONE            int = 0b0000_0000_0000_0000_0000_0000_0000_0000
	FLAGSLIDER_ALWAYSCLAMP     int = 0b0000_0000_0000_0000_0000_0000_0001_0000
	FLAGSLIDER_LOGARITHMIC     int = 0b0000_0000_0000_0000_0000_0000_0010_0000
	FLAGSLIDER_NOROUNDTOFORMAT int = 0b0000_0000_0000_0000_0000_0000_0100_0000
	FLAGSLIDER_NOINPUT         int = 0b0000_0000_0000_0000_0000_0000_1000_0000
	FLAGSLIDER_INVALIDMASK     int = 0b0111_0000_0000_0000_0000_0000_0000_1111
)

const (
	FLAGTABBAR_NONE                         int = 0b0000_0000
	FLAGTABBAR_REORDERABLE                  int = 0b0000_0001
	FLAGTABBAR_AUTOSELECTNEWTABS            int = 0b0000_0010
	FLAGTABBAR_TABLLISTPOPUPBUTTON          int = 0b0000_0100
	FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON int = 0b0000_1000
	FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS    int = 0b0001_0000
	FLAGTABBAR_NOTOOLTIP                    int = 0b0010_0000
	FLAGTABBAR_FITTINGPOLICYRESIZEDOWN      int = 0b0100_0000
	FLAGTABBAR_FITTINGPOLICYSCROLL          int = 0b1000_0000
	FLAGTABBAR_FITTINGPOLICYMASK            int = 0b1100_0000
	FLAGTABBAR_FITTINGPOLICYDEFAULT         int = 0b0100_0000
)

const (
	FLAGTABITEM_NONE                         int = 0b0000_0000_0000
	FLAGTABITEM_UNSAVEDOCUMENT               int = 0b0000_0000_0001
	FLAGTABITEM_SETSELECTED                  int = 0b0000_0000_0010
	FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON int = 0b0000_0000_0100
	FLAGTABITEM_NOPUSHID                     int = 0b0000_0000_1000
	FLAGTABITEM_NOTOOLTIP                    int = 0b0000_0001_0000
	FLAGTABITEM_NOREORDER                    int = 0b0000_0010_0000
	FLAGTABITEM_LEADING                      int = 0b0000_0100_0000
	FLAGTABITEM_TRAILING                     int = 0b0000_1000_0000
	FLAGTABITEM_NOASSUMEDCLOSURE             int = 0b0001_0000_0000
)

const (
	FLAGTABLECOLUMN_NONE                 int = 0b0000_0000_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_DEFAULTHIDE          int = 0b0000_0000_0000_0000_0000_0000_0000_0010
	FLAGTABLECOLUMN_DEFAULTSORT          int = 0b0000_0000_0000_0000_0000_0000_0000_0100
	FLAGTABLECOLUMN_WIDTHSTRETCH         int = 0b0000_0000_0000_0000_0000_0000_0000_1000
	FLAGTABLECOLUMN_WIDTHFIXED           int = 0b0000_0000_0000_0000_0000_0000_0001_0000
	FLAGTABLECOLUMN_NORESIZE             int = 0b0000_0000_0000_0000_0000_0000_0010_0000
	FLAGTABLECOLUMN_NOREORDER            int = 0b0000_0000_0000_0000_0000_0000_0100_0000
	FLAGTABLECOLUMN_NOHIDE               int = 0b0000_0000_0000_0000_0000_0000_1000_0000
	FLAGTABLECOLUMN_NOCLIP               int = 0b0000_0000_0000_0000_0000_0001_0000_0000
	FLAGTABLECOLUMN_NOSORT               int = 0b0000_0000_0000_0000_0000_0010_0000_0000
	FLAGTABLECOLUMN_NOSORTASCENDING      int = 0b0000_0000_0000_0000_0000_0100_0000_0000
	FLAGTABLECOLUMN_NOSORTDESCENDING     int = 0b0000_0000_0000_0000_0000_1000_0000_0000
	FLAGTABLECOLUMN_NOHEADERWIDTH        int = 0b0000_0000_0000_0000_0010_0000_0000_0000
	FLAGTABLECOLUMN_PREFERSORTASCENDING  int = 0b0000_0000_0000_0000_0100_0000_0000_0000
	FLAGTABLECOLUMN_PREFERSORTDESCENDING int = 0b0000_0000_0000_0000_1000_0000_0000_0000
	FLAGTABLECOLUMN_INDENTENABLE         int = 0b0000_0000_0000_0001_0000_0000_0000_0000
	FLAGTABLECOLUMN_INDENTDISABLE        int = 0b0000_0000_0000_0010_0000_0000_0000_0000
	FLAGTABLECOLUMN_ISENABLED            int = 0b0000_0001_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_ISVISIBLE            int = 0b0000_0010_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_ISSORTED             int = 0b0000_0100_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_ISHOVERED            int = 0b0000_1000_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_WIDTHMASK            int = 0b0000_0000_0000_0000_0000_0000_0001_1000
	FLAGTABLECOLUMN_INDENTMASK           int = 0b0000_0000_0000_0011_0000_0000_0000_0000
	FLAGTABLECOLUMN_STATUSMASK           int = 0b0000_1111_0000_0000_0000_0000_0000_0000
	FLAGTABLECOLUMN_NODIRECTRESIZE       int = 0b0100_0000_0000_0000_0000_0000_0000_0000
)

const (
	FLAGTABLE_NONE                       int = 0b0000_0000_0000_0000_0000_0000_0000_0000
	FLAGTABLE_RESIZEABLE                 int = 0b0000_0000_0000_0000_0000_0000_0000_0001
	FLAGTABLE_REORDERABLE                int = 0b0000_0000_0000_0000_0000_0000_0000_0010
	FLAGTABLE_HIDEABLE                   int = 0b0000_0000_0000_0000_0000_0000_0000_0100
	FLAGTABLE_SORTABLE                   int = 0b0000_0000_0000_0000_0000_0000_0000_1000
	FLAGTABLE_NOSAVEDSETTINGS            int = 0b0000_0000_0000_0000_0000_0000_0001_0000
	FLAGTABLE_CONTEXTMENUINBODY          int = 0b0000_0000_0000_0000_0000_0000_0010_0000
	FLAGTABLE_ROWBG                      int = 0b0000_0000_0000_0000_0000_0000_0100_0000
	FLAGTABLE_BORDERSINNERH              int = 0b0000_0000_0000_0000_0000_0000_1000_0000
	FLAGTABLE_BORDERSOUTERH              int = 0b0000_0000_0000_0000_0000_0001_0000_0000
	FLAGTABLE_BORDERSINNERV              int = 0b0000_0000_0000_0000_0000_0010_0000_0000
	FLAGTABLE_BORDERSOUTERV              int = 0b0000_0000_0000_0000_0000_0100_0000_0000
	FLAGTABLE_BORDERSH                   int = 0b0000_0000_0000_0000_0000_0001_1000_0000
	FLAGTABLE_BORDERSV                   int = 0b0000_0000_0000_0000_0000_0110_0000_0000
	FLAGTABLE_BORDERSINNER               int = 0b0000_0000_0000_0000_0000_0010_1000_0000
	FLAGTABLE_BORDERSOUTER               int = 0b0000_0000_0000_0000_0000_0101_0000_0000
	FLAGTABLE_BORDERS                    int = 0b0000_0000_0000_0000_0000_0111_1000_0000
	FLAGTABLE_NOBORDERSINBODY            int = 0b0000_0000_0000_0000_0000_1000_0000_0000
	FLAGTABLE_NOBORDERSINBODYUNTILRESIZE int = 0b0000_0000_0000_0000_0001_0000_0000_0000
	FLAGTABLE_SIZINGFIXEDFIT             int = 0b0000_0000_0000_0000_0010_0000_0000_0000
	FLAGTABLE_SIZINGFIXEDSAME            int = 0b0000_0000_0000_0000_0100_0000_0000_0000
	FLAGTABLE_SIZINGSTRETCHPROP          int = 0b0000_0000_0000_0000_0110_0000_0000_0000
	FLAGTABLE_SIZINGSTRETCHSAME          int = 0b0000_0000_0000_0000_1000_0000_0000_0000
	FLAGTABLE_NOHOSTEXTENDX              int = 0b0000_0000_0000_0001_0000_0000_0000_0000
	FLAGTABLE_NOHOSTEXTENDY              int = 0b0000_0000_0000_0010_0000_0000_0000_0000
	FLAGTABLE_NOKEEPCOLUMNSVISIBLE       int = 0b0000_0000_0000_0100_0000_0000_0000_0000
	FLAGTABLE_PRECISEWIDTHS              int = 0b0000_0000_0000_1000_0000_0000_0000_0000
	FLAGTABLE_NOCLIP                     int = 0b0000_0000_0001_0000_0000_0000_0000_0000
	FLAGTABLE_PADOUTERX                  int = 0b0000_0000_0010_0000_0000_0000_0000_0000
	FLAGTABLE_NOPADOUTERX                int = 0b0000_0000_0100_0000_0000_0000_0000_0000
	FLAGTABLE_NOPADINNERX                int = 0b0000_0000_1000_0000_0000_0000_0000_0000
	FLAGTABLE_SCROLLX                    int = 0b0000_0001_0000_0000_0000_0000_0000_0000
	FLAGTABLE_SCROLLY                    int = 0b0000_0010_0000_0000_0000_0000_0000_0000
	FLAGTABLE_SORTMULTI                  int = 0b0000_0100_0000_0000_0000_0000_0000_0000
	FLAGTABLE_SORTTRISTATE               int = 0b0000_1000_0000_0000_0000_0000_0000_0000
	FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN     int = 0b0001_0000_0000_0000_0000_0000_0000_0000
	FLAGTABLE_SIZINGMASK                 int = 0b0000_0000_0000_0000_1110_0000_0000_0000
)

const (
	FLAGTABLEROW_NONE    int = 0b0
	FLAGTABLEROW_HEADERS int = 0b1
)

const (
	DIR_NONE int = iota - 1
	DIR_LEFT
	DIR_RIGHT
	DIR_UP
	DIR_DOWN
	DIR_COUNT
)

const (
	FLAGTREENODE_NONE                 int = 0b0000_0000_0000_0000
	FLAGTREENODE_SELECTED             int = 0b0000_0000_0000_0001
	FLAGTREENODE_FRAMED               int = 0b0000_0000_0000_0010
	FLAGTREENODE_ALLOWOVERLAP         int = 0b0000_0000_0000_0100
	FLAGTREENODE_NOTREEPUSHONOPEN     int = 0b0000_0000_0000_1000
	FLAGTREENODE_NOAUTOOPENONLOG      int = 0b0000_0000_0001_0000
	FLAGTREENODE_DEFAULTOPEN          int = 0b0000_0000_0010_0000
	FLAGTREENODE_OPENONDOUBLECLICK    int = 0b0000_0000_0100_0000
	FLAGTREENODE_OPENONARROW          int = 0b0000_0000_1000_0000
	FLAGTREENODE_LEAF                 int = 0b0000_0001_0000_0000
	FLAGTREENODE_BULLET               int = 0b0000_0010_0000_0000
	FLAGTREENODE_FRAMEPADDING         int = 0b0000_0100_0000_0000
	FLAGTREENODE_SPANAVAILWIDTH       int = 0b0000_1000_0000_0000
	FLAGTREENODE_SPANFULLWIDTH        int = 0b0001_0000_0000_0000
	FLAGTREENODE_SPANALLCOLUMNS       int = 0b0010_0000_0000_0000
	FLAGTREENODE_NAVLEFTJUMPSBACKHERE int = 0b0100_0000_0000_0000
	FLAGTREENODE_COLLAPSINGHEADER     int = 0b0000_0000_0001_1010
)

const (
	FLAGMASTERWINDOW_NOTRESIZABLE int = 1 << iota
	FLAGMASTERWINDOW_MAXIMIZED
	FLAGMASTERWINDOW_FLOATING
	FLAGMASTERWINDOW_FRAMELESS
	FLAGMASTERWINDOW_TRANSPARENT
)

const (
	FLAGWINDOW_NONE                      int = 0b0000_0000_0000_0000_0000
	FLAGWINDOW_NOTITLEBAR                int = 0b0000_0000_0000_0000_0001
	FLAGWINDOW_NORESIZE                  int = 0b0000_0000_0000_0000_0010
	FLAGWINDOW_NOMOVE                    int = 0b0000_0000_0000_0000_0100
	FLAGWINDOW_NOSCROLLBAR               int = 0b0000_0000_0000_0000_1000
	FLAGWINDOW_NOSCROLLWITHMOUSE         int = 0b0000_0000_0000_0001_0000
	FLAGWINDOW_NOCOLLAPSE                int = 0b0000_0000_0000_0010_0000
	FLAGWINDOW_ALWAYSAUTORESIZE          int = 0b0000_0000_0000_0100_0000
	FLAGWINDOW_NOBACKGROUND              int = 0b0000_0000_0000_1000_0000
	FLAGWINDOW_NOSAVEDSETTINGS           int = 0b0000_0000_0001_0000_0000
	FLAGWINDOW_NOMOUSEINPUTS             int = 0b0000_0000_0010_0000_0000
	FLAGWINDOW_MENUBAR                   int = 0b0000_0000_0100_0000_0000
	FLAGWINDOW_HORIZONTALSCROLLBAR       int = 0b0000_0000_1000_0000_0000
	FLAGWINDOW_NOFOCUSONAPPEARING        int = 0b0000_0001_0000_0000_0000
	FLAGWINDOW_NOBRINGTOFRONTONFOCUS     int = 0b0000_0010_0000_0000_0000
	FLAGWINDOW_ALWAYSVERTICALSCROLLBAR   int = 0b0000_0100_0000_0000_0000
	FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR int = 0b0000_1000_0000_0000_0000
	FLAGWINDOW_NONAVINPUTS               int = 0b0001_0000_0000_0000_0000
	FLAGWINDOW_NONAVFOCUS                int = 0b0010_0000_0000_0000_0000
	FLAGWINDOW_UNSAVEDDOCUMENT           int = 0b0100_0000_0000_0000_0000
	FLAGWINDOW_NONAV                     int = 0b0011_0000_0000_0000_0000
	FLAGWINDOW_NODECORATION              int = 0b0000_0000_0000_0010_1011
	FLAGWINDOW_NOINPUTS                  int = 0b0011_0000_0010_0000_0000
)

const (
	SPLITDIRECTION_HORIZONTAL int = 1 << iota
	SPLITDIRECTION_VERTICAL
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
	WIDGET_BUTTON_SMALL         = "button_small"
	WIDGET_BUTTON_RADIO         = "button_radio"
	WIDGET_IMAGE_URL            = "image_url"
	WIDGET_IMAGE                = "image"
	WIDGET_LIST_BOX             = "list_box"
	WIDGET_LIST_CLIPPER         = "list_clipper"
	WIDGET_MENU_BAR_MAIN        = "menu_bar_main"
	WIDGET_MENU_BAR             = "menu_bar"
	WIDGET_MENU_ITEM            = "menu_item"
	WIDGET_MENU                 = "menu"
	WIDGET_SELECTABLE           = "selectable"
	WIDGET_SLIDER_FLOAT         = "slider_float"
	WIDGET_SLIDER_INT           = "slider_int"
	WIDGET_VSLIDER_INT          = "vslider_int"
	WIDGET_TAB_BAR              = "tab_bar"
	WIDGET_TAB_ITEM             = "tab_item"
	WIDGET_TOOLTIP              = "tooltip"
	WIDGET_TABLE_COLUMN         = "table_column"
	WIDGET_TABLE_ROW            = "table_row"
	WIDGET_TABLE                = "table"
	WIDGET_BUTTON_ARROW         = "button_arrow"
	WIDGET_TREE_TABLE_NODE      = "tree_table_node"
	WIDGET_TREE_TABLE_ROW       = "tree_table_row"
	WIDGET_TREE_TABLE           = "tree_table"
	WIDGET_WINDOW_SINGLE        = "window_single"
	WIDGET_POPUP_MODAL          = "popup_modal"
	WIDGET_POPUP                = "popup"
	WIDGET_LAYOUT_SPLIT         = "layout_split"
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
		WIDGET_BUTTON_SMALL:         buttonSmallBuild,
		WIDGET_BUTTON_RADIO:         buttonRadioBuild,
		WIDGET_IMAGE_URL:            imageUrlBuild,
		WIDGET_IMAGE:                imageBuild,
		WIDGET_LIST_BOX:             listBoxBuild,
		WIDGET_LIST_CLIPPER:         listClipperBuild,
		WIDGET_MENU_BAR_MAIN:        mainMenuBarBuild,
		WIDGET_MENU_BAR:             menuBarBuild,
		WIDGET_MENU_ITEM:            menuItemBuild,
		WIDGET_MENU:                 menuBuild,
		WIDGET_SELECTABLE:           selectableBuild,
		WIDGET_SLIDER_FLOAT:         sliderFloatBuild,
		WIDGET_SLIDER_INT:           sliderIntBuild,
		WIDGET_VSLIDER_INT:          vsliderIntBuild,
		WIDGET_TAB_BAR:              tabbarBuild,
		WIDGET_TOOLTIP:              tooltipBuild,
		WIDGET_TABLE:                tableBuild,
		WIDGET_BUTTON_ARROW:         buttonArrowBuild,
		WIDGET_TREE_TABLE_NODE:      treeTableNodeBuild,
		WIDGET_TREE_TABLE:           treeTableBuild,
		WIDGET_POPUP_MODAL:          popupModalBuild,
		WIDGET_POPUP:                popupBuild,
		WIDGET_LAYOUT_SPLIT:         splitLayoutBuild,
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

func buttonSmallTable(state *golua.LState, text string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON_SMALL))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("__click"), golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	return t
}

func buttonSmallBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := state.GetTable(t, golua.LString("label")).(golua.LString)
	b := g.SmallButton(string(text))

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func buttonRadioTable(state *golua.LState, text string, active bool) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON_RADIO))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("active"), golua.LBool(active))
	state.SetTable(t, golua.LString("__change"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	return t
}

func buttonRadioBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := state.GetTable(t, golua.LString("label")).(golua.LString)
	active := state.GetTable(t, golua.LString("active")).(golua.LBool)
	b := g.RadioButton(string(text), bool(active))

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Call(0, 0)
		})
	}

	return b
}

func imageUrlTable(state *golua.LState, url string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_IMAGE_URL))
	state.SetTable(t, golua.LString("url"), golua.LString(url))
	state.SetTable(t, golua.LString("__click"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__timeout"), golua.LNil)
	state.SetTable(t, golua.LString("__failwidgets"), golua.LNil)
	state.SetTable(t, golua.LString("__loadwidgets"), golua.LNil)
	state.SetTable(t, golua.LString("__failure"), golua.LNil)
	state.SetTable(t, golua.LString("__ready"), golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "timeout", func(state *golua.LState, t *golua.LTable) {
		timeout := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__timeout"), timeout)
	})

	tableBuilderFunc(state, t, "layout_for_failure", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__failwidgets"), lt)
	})

	tableBuilderFunc(state, t, "layout_for_loading", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__loadwidgets"), lt)
	})

	tableBuilderFunc(state, t, "on_failure", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__failure"), fn)
	})

	tableBuilderFunc(state, t, "on_ready", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__ready"), fn)
	})

	return t
}

func imageUrlBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	url := state.GetTable(t, golua.LString("url")).(golua.LString)
	i := g.ImageWithURL(string(url))

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		i.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	timeout := state.GetTable(t, golua.LString("__timeout"))
	if timeout.Type() == golua.LTNumber {
		i.Timeout(time.Duration(timeout.(golua.LNumber)))
	}

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		i.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	lfail := state.GetTable(t, golua.LString("__failwidgets"))
	if lfail.Type() == golua.LTTable {
		i.LayoutForFailure(layoutBuild(r, state, parseWidgets(parseTable(lfail.(*golua.LTable), state), state, lg), lg)...)
	}

	lload := state.GetTable(t, golua.LString("__loadwidgets"))
	if lload.Type() == golua.LTTable {
		i.LayoutForFailure(layoutBuild(r, state, parseWidgets(parseTable(lload.(*golua.LTable), state), state, lg), lg)...)
	}

	failure := state.GetTable(t, golua.LString("__failure"))
	if failure.Type() == golua.LTFunction {
		i.OnFailure(func(err error) {
			lg.Append(fmt.Sprintf("error occured while loading image url: %s", err), log.LEVEL_WARN)
			state.Push(failure)
			state.Call(0, 0)
		})
	}

	ready := state.GetTable(t, golua.LString("__ready"))
	if ready.Type() == golua.LTFunction {
		i.OnReady(func() {
			state.Push(ready)
			state.Call(0, 0)
		})
	}

	return i
}

func imageTable(state *golua.LState, image int, sync bool) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_IMAGE))
	state.SetTable(t, golua.LString("image"), golua.LNumber(image))
	state.SetTable(t, golua.LString("sync"), golua.LBool(sync))
	state.SetTable(t, golua.LString("__click"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	return t
}

func imageBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ig := state.GetTable(t, golua.LString("image")).(golua.LNumber)
	var img image.Image

	sync := state.GetTable(t, golua.LString("sync")).(golua.LBool)

	if !sync {
		<-r.IC.Schedule(int(ig), &collection.Task[collection.ItemImage]{
			Lib:  LIB_GUI,
			Name: "wg_image",
			Fn: func(i *collection.Item[collection.ItemImage]) {
				img = i.Self.Image
			},
		})
	} else {
		item := r.IC.Item(int(ig))
		if item.Self.Image == nil {
			img = image.NewRGBA(image.Rectangle{
				Min: image.Pt(0, 0),
				Max: image.Pt(1, 1), // image must have at least 1 pixel for imgui.
			})
		} else {
			img = item.Self.Image
		}
	}

	i := g.ImageWithRgba(img)

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		i.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		i.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return i
}

func listBoxTable(state *golua.LState, items golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LIST_BOX))
	state.SetTable(t, golua.LString("items"), items)
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__border"), golua.LNil)
	state.SetTable(t, golua.LString("__context"), golua.LNil)
	state.SetTable(t, golua.LString("__dclick"), golua.LNil)
	state.SetTable(t, golua.LString("__menu"), golua.LNil)
	state.SetTable(t, golua.LString("__sel"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__border"), golua.LBool(b))
	})

	tableBuilderFunc(state, t, "context_menu", func(state *golua.LState, t *golua.LTable) {
		cmt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__context"), cmt)
	})

	tableBuilderFunc(state, t, "on_double_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__dclick"), fn)
	})

	tableBuilderFunc(state, t, "on_menu", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__menu"), fn)
	})

	tableBuilderFunc(state, t, "selected_index", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__sel"), sel)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	return t
}

func listBoxBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	it := state.GetTable(t, golua.LString("items")).(*golua.LTable)
	items := []string{}
	for i := range it.Len() {
		is := state.GetTable(it, golua.LNumber(i+1)).(golua.LString)
		items = append(items, string(is))
	}

	b := g.ListBox(items)

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		b.OnChange(func(index int) {
			state.Push(change)
			state.Push(golua.LNumber(index))
			state.Call(1, 0)
		})
	}

	dclick := state.GetTable(t, golua.LString("__dclick"))
	if dclick.Type() == golua.LTFunction {
		b.OnDClick(func(index int) {
			state.Push(dclick)
			state.Push(golua.LNumber(index))
			state.Call(1, 0)
		})
	}

	selmenu := state.GetTable(t, golua.LString("__menu"))
	if selmenu.Type() == golua.LTFunction {
		b.OnMenu(func(index int, menu string) {
			state.Push(selmenu)
			state.Push(golua.LNumber(index))
			state.Push(golua.LString(menu))
			state.Call(2, 0)
		})
	}

	sel := state.GetTable(t, golua.LString("__sel"))
	if sel.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(sel.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		b.SelectedIndex(ref.Value.(*int32))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return b
}

func listClipperTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LIST_CLIPPER))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func listClipperBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.ListClipper()

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func mainMenuBarTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MENU_BAR_MAIN))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func mainMenuBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.MainMenuBar()

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func menuBarTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MENU_BAR))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func menuBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.MenuBar()

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func menuItemTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MENU_ITEM))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__enabled"), golua.LNil)
	state.SetTable(t, golua.LString("__click"), golua.LNil)
	state.SetTable(t, golua.LString("__sel"), golua.LNil)
	state.SetTable(t, golua.LString("__shortcut"), golua.LNil)

	tableBuilderFunc(state, t, "enabled", func(state *golua.LState, t *golua.LTable) {
		en := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__enabled"), golua.LBool(en))
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	tableBuilderFunc(state, t, "selected", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__sel"), golua.LBool(sel))
	})

	tableBuilderFunc(state, t, "shortcut", func(state *golua.LState, t *golua.LTable) {
		sc := state.CheckString(-1)
		state.SetTable(t, golua.LString("__shortcut"), golua.LString(sc))
	})

	return t
}

func menuItemBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	m := g.MenuItem(string(label))

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		m.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	enabled := state.GetTable(t, golua.LString("__enabled"))
	if enabled.Type() == golua.LTBool {
		m.Enabled(bool(enabled.(golua.LBool)))
	}

	sel := state.GetTable(t, golua.LString("__sel"))
	if sel.Type() == golua.LTBool {
		m.Selected(bool(sel.(golua.LBool)))
	}

	shortcut := state.GetTable(t, golua.LString("__shortcut"))
	if shortcut.Type() == golua.LTString {
		m.Shortcut(string(shortcut.(golua.LString)))
	}

	return m
}

func menuTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MENU))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__enabled"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "enabled", func(state *golua.LState, t *golua.LTable) {
		en := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__enabled"), golua.LBool(en))
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func menuBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	m := g.Menu(string(label))

	enabled := state.GetTable(t, golua.LString("__enabled"))
	if enabled.Type() == golua.LTBool {
		m.Enabled(bool(enabled.(golua.LBool)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return m
}

func selectableTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SELECTABLE))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__click"), golua.LNil)
	state.SetTable(t, golua.LString("__dclick"), golua.LNil)
	state.SetTable(t, golua.LString("__sel"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	tableBuilderFunc(state, t, "on_double_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__dclick"), fn)
	})

	tableBuilderFunc(state, t, "selected", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__sel"), golua.LBool(sel))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	return t
}

func selectableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	b := g.Selectable(string(label))

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	dclick := state.GetTable(t, golua.LString("__dclick"))
	if dclick.Type() == golua.LTFunction {
		b.OnDClick(func() {
			state.Push(dclick)
			state.Call(0, 0)
		})
	}

	sel := state.GetTable(t, golua.LString("__sel"))
	if sel.Type() == golua.LTBool {
		b.Selected(bool(sel.(golua.LBool)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		b.Flags(g.SelectableFlags(flags.(golua.LNumber)))
	}

	return b
}

func sliderFloatTable(state *golua.LState, f32ref int, min, max float64) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SLIDER_FLOAT))
	state.SetTable(t, golua.LString("f32ref"), golua.LNumber(f32ref))
	state.SetTable(t, golua.LString("min"), golua.LNumber(min))
	state.SetTable(t, golua.LString("max"), golua.LNumber(max))
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	return t
}

func sliderFloatBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	floatref := state.GetTable(t, golua.LString("f32ref")).(golua.LNumber)
	ref, err := r.CR_REF.Item(int(floatref))
	value := ref.Value.(*float32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := state.GetTable(t, golua.LString("min")).(golua.LNumber)
	max := state.GetTable(t, golua.LString("max")).(golua.LNumber)
	b := g.SliderFloat(value, float32(min), float32(max))

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(floatref)
			state.Call(2, 0)
		})
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)))
	}

	return b
}

func sliderIntTable(state *golua.LState, i32ref int, min, max int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SLIDER_INT))
	state.SetTable(t, golua.LString("i32ref"), golua.LNumber(i32ref))
	state.SetTable(t, golua.LString("min"), golua.LNumber(min))
	state.SetTable(t, golua.LString("max"), golua.LNumber(max))
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	return t
}

func sliderIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	intref := state.GetTable(t, golua.LString("i32ref")).(golua.LNumber)
	ref, err := r.CR_REF.Item(int(intref))
	value := ref.Value.(*int32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := state.GetTable(t, golua.LString("min")).(golua.LNumber)
	max := state.GetTable(t, golua.LString("max")).(golua.LNumber)
	b := g.SliderInt(value, int32(min), int32(max))

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(intref)
			state.Call(2, 0)
		})
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)))
	}

	return b
}

func vsliderIntTable(state *golua.LState, i32ref int, min, max int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_VSLIDER_INT))
	state.SetTable(t, golua.LString("i32ref"), golua.LNumber(i32ref))
	state.SetTable(t, golua.LString("min"), golua.LNumber(min))
	state.SetTable(t, golua.LString("max"), golua.LNumber(max))
	state.SetTable(t, golua.LString("__change"), golua.LNil)
	state.SetTable(t, golua.LString("__label"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__change"), fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		state.SetTable(t, golua.LString("__label"), golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	return t
}

func vsliderIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	intref := state.GetTable(t, golua.LString("i32ref")).(golua.LNumber)
	ref, err := r.CR_REF.Item(int(intref))
	value := ref.Value.(*int32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := state.GetTable(t, golua.LString("min")).(golua.LNumber)
	max := state.GetTable(t, golua.LString("max")).(golua.LNumber)
	b := g.VSliderInt(value, int32(min), int32(max))

	change := state.GetTable(t, golua.LString("__change"))
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(intref)
			state.Call(2, 0)
		})
	}

	label := state.GetTable(t, golua.LString("__label"))
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		b.Flags(g.SliderFlags(flags.(golua.LNumber)))
	}

	return b
}

func tabbarTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TAB_BAR))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "tab_items", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func tabbarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.TabBar()

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TabBarFlags(flags.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg)
		wdi := []*g.TabItemWidget{}
		for _, w := range wd {
			i := tabitemBuild(r, lg, state, w)
			wdi = append(wdi, i)
		}
		tb.TabItems(wdi...)
	}

	return tb
}

func tabitemTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TAB_ITEM))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__open"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__open"), flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func tabitemBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TabItemWidget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	i := g.TabItem(string(label))

	open := state.GetTable(t, golua.LString("__open"))
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		i.IsOpen(ref.Value.(*bool))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		i.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		i.Flags(g.TabItemFlags(flags.(golua.LNumber)))
	}

	return i
}

func tooltipTable(state *golua.LState, tip string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TOOLTIP))
	state.SetTable(t, golua.LString("tip"), golua.LString(tip))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func tooltipBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tip := state.GetTable(t, golua.LString("tip")).(golua.LString)
	i := g.Tooltip(string(tip))

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		i.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return i
}

func tableColumnTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TABLE_COLUMN))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "inner_width_or_weight", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), flags)
	})

	return t
}

func tableColumnBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TableColumnWidget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	c := g.TableColumn(string(label))

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.TableColumnFlags(flags.(golua.LNumber)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		c.InnerWidthOrWeight(float32(width.(golua.LNumber)))
	}

	return c
}

func tableRowTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TABLE_ROW))
	state.SetTable(t, golua.LString("widgets"), widgets)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__color"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "bg_color", func(state *golua.LState, t *golua.LTable) {
		clr := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__color"), clr)
	})

	tableBuilderFunc(state, t, "min_height", func(state *golua.LState, t *golua.LTable) {
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__height"), height)
	})

	return t
}

func tableRowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TableRowWidget {
	var widgets []g.Widget

	wid := state.GetTable(t, golua.LString("widgets"))
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable), state), state, lg), lg)
	}

	s := g.TableRow(widgets...)

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		s.Flags(g.TableRowFlags(flags.(golua.LNumber)))
	}

	height := state.GetTable(t, golua.LString("__height"))
	if height.Type() == golua.LTNumber {
		s.MinHeight(float64(height.(golua.LNumber)))
	}

	clr := state.GetTable(t, golua.LString("__color"))
	if clr.Type() == golua.LTTable {
		rgba := imageutil.TableToRGBA(state, clr.(*golua.LTable))
		s.BgColor(rgba)
	}

	return s
}

func tableTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TABLE))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__columns"), golua.LNil)
	state.SetTable(t, golua.LString("__rows"), golua.LNil)
	state.SetTable(t, golua.LString("__fast"), golua.LNil)
	state.SetTable(t, golua.LString("__freeze"), state.NewTable())
	state.SetTable(t, golua.LString("__innerwidth"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "fast_mode", func(state *golua.LState, t *golua.LTable) {
		fast := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__fast"), golua.LBool(fast))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "columns", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__columns"), lt)
	})

	tableBuilderFunc(state, t, "rows", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__rows"), lt)
	})

	tableBuilderFunc(state, t, "inner_width", func(state *golua.LState, t *golua.LTable) {
		innerwidth := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__innerwidth"), innerwidth)
	})

	tableBuilderFunc(state, t, "freeze", func(state *golua.LState, t *golua.LTable) {
		col := state.CheckNumber(-2)
		row := state.CheckNumber(-1)
		pt := state.NewTable()
		state.SetTable(pt, golua.LString("col"), col)
		state.SetTable(pt, golua.LString("row"), row)

		ft := state.GetTable(t, golua.LString("__freeze")).(*golua.LTable)
		ft.Append(pt)
	})

	return t
}

func tableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.Table()

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TableFlags(flags.(golua.LNumber)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		tb.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	innerwidth := state.GetTable(t, golua.LString("__innerwidth"))
	if innerwidth.Type() == golua.LTNumber {
		tb.InnerWidth(float64(innerwidth.(golua.LNumber)))
	}

	fast := state.GetTable(t, golua.LString("__fast"))
	if fast.Type() == golua.LTBool {
		tb.FastMode(bool(fast.(golua.LBool)))
	}

	freeze := state.GetTable(t, golua.LString("__freeze")).(*golua.LTable)
	for i := range freeze.Len() {
		pt := state.GetTable(freeze, golua.LNumber(i+1)).(*golua.LTable)
		col := state.GetTable(pt, golua.LString("col")).(golua.LNumber)
		row := state.GetTable(pt, golua.LString("row")).(golua.LNumber)

		tb.Freeze(int(col), int(row))
	}

	columns := state.GetTable(t, golua.LString("__columns"))
	if columns.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(columns.(*golua.LTable), state), state, lg)
		wdi := []*g.TableColumnWidget{}
		for _, w := range wd {
			i := tableColumnBuild(r, lg, state, w)
			wdi = append(wdi, i)
		}
		tb.Columns(wdi...)
	}

	rows := state.GetTable(t, golua.LString("__rows"))
	if rows.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(rows.(*golua.LTable), state), state, lg)
		wdi := []*g.TableRowWidget{}
		for _, w := range wd {
			i := tableRowBuild(r, lg, state, w)
			wdi = append(wdi, i)
		}
		tb.Rows(wdi...)
	}

	return tb
}

func buttonArrowTable(state *golua.LState, dir int) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON_ARROW))
	state.SetTable(t, golua.LString("dir"), golua.LNumber(dir))
	state.SetTable(t, golua.LString("__click"), golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__click"), fn)
	})

	return t
}

func buttonArrowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	dir := state.GetTable(t, golua.LString("dir")).(golua.LNumber)
	b := g.ArrowButton(g.Direction(dir))

	click := state.GetTable(t, golua.LString("__click"))
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func treeTableNodeTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TREE_TABLE_NODE))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func treeTableNodeBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	n := g.TreeNode(string(label))

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		n.Flags(g.TreeNodeFlags(flags.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		n.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return n
}

func treeTableRowTable(state *golua.LState, label string, widgets golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TREE_TABLE_ROW))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("widgets"), widgets)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__children"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "children", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__children"), lt)
	})

	return t
}

func treeTableRowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TreeTableRowWidget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	var widgets []g.Widget

	wid := state.GetTable(t, golua.LString("widgets"))
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable), state), state, lg), lg)
	}

	n := g.TreeTableRow(string(label), widgets...)

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		n.Flags(g.TreeNodeFlags(flags.(golua.LNumber)))
	}

	children := state.GetTable(t, golua.LString("__children"))
	if children.Type() == golua.LTTable {
		rwid := parseWidgets(parseTable(children.(*golua.LTable), state), state, lg)
		childs := []*g.TreeTableRowWidget{}
		for _, c := range rwid {
			childs = append(childs, treeTableRowBuild(r, lg, state, c))
		}
		n.Children(childs...)
	}

	return n
}

func treeTableTable(state *golua.LState) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TREE_TABLE))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__columns"), golua.LNil)
	state.SetTable(t, golua.LString("__rows"), golua.LNil)
	state.SetTable(t, golua.LString("__freeze"), state.NewTable())
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "columns", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__columns"), lt)
	})

	tableBuilderFunc(state, t, "rows", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__rows"), lt)
	})

	tableBuilderFunc(state, t, "freeze", func(state *golua.LState, t *golua.LTable) {
		col := state.CheckNumber(-2)
		row := state.CheckNumber(-1)
		pt := state.NewTable()
		state.SetTable(pt, golua.LString("col"), col)
		state.SetTable(pt, golua.LString("row"), row)

		ft := state.GetTable(t, golua.LString("__freeze")).(*golua.LTable)
		ft.Append(pt)
	})

	return t
}

func treeTableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.TreeTable()

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TableFlags(flags.(golua.LNumber)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		tb.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	freeze := state.GetTable(t, golua.LString("__freeze")).(*golua.LTable)
	for i := range freeze.Len() {
		pt := state.GetTable(freeze, golua.LNumber(i+1)).(*golua.LTable)
		col := state.GetTable(pt, golua.LString("col")).(golua.LNumber)
		row := state.GetTable(pt, golua.LString("row")).(golua.LNumber)

		tb.Freeze(int(col), int(row))
	}

	columns := state.GetTable(t, golua.LString("__columns"))
	if columns.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(columns.(*golua.LTable), state), state, lg)
		wdi := []*g.TableColumnWidget{}
		for _, w := range wd {
			i := tableColumnBuild(r, lg, state, w)
			wdi = append(wdi, i)
		}
		tb.Columns(wdi...)
	}

	rows := state.GetTable(t, golua.LString("__rows"))
	if rows.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(rows.(*golua.LTable), state), state, lg)
		wdi := []*g.TreeTableRowWidget{}
		for _, w := range wd {
			i := treeTableRowBuild(r, lg, state, w)
			wdi = append(wdi, i)
		}
		tb.Rows(wdi...)
	}

	return tb
}

func windowTable(r *lua.Runner, lg *log.Logger, state *golua.LState, single bool, menubar bool, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_WINDOW_SINGLE))
	state.SetTable(t, golua.LString("single"), golua.LBool(single))
	state.SetTable(t, golua.LString("menubar"), golua.LBool(menubar))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__front"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__open"), golua.LNil)
	state.SetTable(t, golua.LString("__posx"), golua.LNil)
	state.SetTable(t, golua.LString("__posy"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__ready"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "pos", func(state *golua.LState, t *golua.LTable) {
		posx := state.CheckNumber(-2)
		posy := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__posx"), posx)
		state.SetTable(t, golua.LString("__posy"), posy)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		open := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__open"), open)
	})

	tableBuilderFunc(state, t, "bring_to_front", func(state *golua.LState, t *golua.LTable) {
		state.SetTable(t, golua.LString("__front"), golua.LTrue)
	})

	tableBuilderFunc(state, t, "ready", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__ready"), fn)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
		windowBuild(r, lg, state, t)
	})

	return t
}

func windowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.WindowWidget {
	var w *g.WindowWidget

	single := state.GetTable(t, golua.LString("single")).(golua.LBool)
	if single {
		menubar := state.GetTable(t, golua.LString("menubar")).(golua.LBool)
		if menubar {
			w = g.SingleWindowWithMenuBar()
		} else {
			w = g.SingleWindow()
		}
	} else {
		label := state.GetTable(t, golua.LString("label")).(golua.LString)
		w = g.Window(string(label))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		w.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		w.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	posx := state.GetTable(t, golua.LString("__posx"))
	posy := state.GetTable(t, golua.LString("__posy"))
	if posx.Type() == golua.LTNumber && posy.Type() == golua.LTNumber {
		w.Pos(float32(posx.(golua.LNumber)), float32(posy.(golua.LNumber)))
	}

	front := state.GetTable(t, golua.LString("__front"))
	if front.Type() == golua.LTBool {
		if front.(golua.LBool) {
			w.BringToFront()
		}
	}

	open := state.GetTable(t, golua.LString("__open"))
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		w.IsOpen(ref.Value.(*bool))
	}

	ready := state.GetTable(t, golua.LString("__ready"))
	if ready.Type() == golua.LTFunction {
		fnt := state.NewTable()

		state.SetTable(fnt, golua.LString("current_position"), state.NewFunction(func(state *golua.LState) int {
			x, y := w.CurrentPosition()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		}))

		state.SetTable(fnt, golua.LString("current_size"), state.NewFunction(func(state *golua.LState) int {
			w, h := w.CurrentSize()

			state.Push(golua.LNumber(w))
			state.Push(golua.LNumber(h))
			return 2
		}))

		state.SetTable(fnt, golua.LString("has_focus"), state.NewFunction(func(state *golua.LState) int {
			f := w.HasFocus()

			state.Push(golua.LBool(f))
			return 1
		}))

		state.Push(ready)
		state.Push(fnt)
		state.Call(1, 0)
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		w.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return w
}

func popupModalTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_POPUP_MODAL))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__open"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__open"), flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func popupModalBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	m := g.PopupModal(string(label))

	open := state.GetTable(t, golua.LString("__open"))
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		m.IsOpen(ref.Value.(*bool))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		m.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	return m
}

func popupTable(state *golua.LState, label string) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_POPUP))
	state.SetTable(t, golua.LString("label"), golua.LString(label))
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func popupBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := state.GetTable(t, golua.LString("label")).(golua.LString)
	m := g.Popup(string(label))

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		m.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	return m
}

func splitLayoutTable(state *golua.LState, direction, floatref int, layout1 golua.LValue, layout2 golua.LValue) *golua.LTable {
	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LAYOUT_SPLIT))
	state.SetTable(t, golua.LString("direction"), golua.LNumber(direction))
	state.SetTable(t, golua.LString("floatref"), golua.LNumber(floatref))
	state.SetTable(t, golua.LString("layout1"), layout1)
	state.SetTable(t, golua.LString("layout2"), layout2)
	state.SetTable(t, golua.LString("__border"), golua.LNil)

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		border := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__border"), golua.LBool(border))
	})

	return t
}

func splitLayoutBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	direction := state.GetTable(t, golua.LString("direction")).(golua.LNumber)

	floatref := state.GetTable(t, golua.LString("floatref"))
	ref, err := r.CR_REF.Item(int(floatref.(golua.LNumber)))
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	pos := ref.Value.(*float32)

	var widgets1 []g.Widget
	wid1 := state.GetTable(t, golua.LString("layout1"))
	if wid1.Type() == golua.LTTable {
		widgets1 = layoutBuild(r, state, parseWidgets(parseTable(wid1.(*golua.LTable), state), state, lg), lg)
	}

	var widgets2 []g.Widget
	wid2 := state.GetTable(t, golua.LString("layout2"))
	if wid2.Type() == golua.LTTable {
		widgets2 = layoutBuild(r, state, parseWidgets(parseTable(wid2.(*golua.LTable), state), state, lg), lg)
	}

	s := g.SplitLayout(g.SplitDirection(direction), pos, g.Layout(widgets1), g.Layout(widgets2))

	border := state.GetTable(t, golua.LString("__border"))
	if border.Type() == golua.LTBool {
		s.Border(bool(border.(golua.LBool)))
	}

	return s
}
