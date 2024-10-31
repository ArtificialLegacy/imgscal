package lib

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	g "github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/assets"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_GUI = "gui"

/// @lib GUI
/// @import gui
/// @desc
/// Library for creating custom interfaces.
/// @section
/// Currently does not work well on windows, so WSL is recommended.

func RegisterGUI(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_GUI, r, r.State, lg)

	/// @func window_master(name, width, height, flags?) -> int<collection.CRATE_WINDOW>
	/// @arg name {string}
	/// @arg width {int}
	/// @arg height {int}
	/// @arg? flags {int<gui.MasterWindowFlags>} - Use 'bit.bitor' or 'bit.bitor_many' to combine flags.
	/// @returns {int<collection.CRATE_WINDOW>} - The id of the new window.
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

	/// @func window_pos(id) -> int, int
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @returns {int} - The x position of the window.
	/// @returns {int} - The y position of the window.
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

	/// @func window_set_pos(id, x, y)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg x {int}
	/// @arg y {int}
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

	/// @func window_size(id) -> int, int
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @returns {int} - The width of the window.
	/// @returns {int} - The height of the window.
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

	/// @func window_set_size(id, width, height)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg width {int}
	/// @arg height {int}
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

	/// @func window_set_size_limits(id, minw, minh, maxw, maxh)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg minw {int}
	/// @arg minh {int}
	/// @arg maxw {int}
	/// @arg maxh {int}
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

	/// @func window_set_bg_color_rgba(id, r, g, b, a)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg r {int}
	/// @arg g {int}
	/// @arg b {int}
	/// @arg a {int}
	lib.CreateFunction(tab, "window_set_bg_color_rgba",
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

	/// @func window_set_bg_color(id, color)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "window_set_bg_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			r, g, b, a := imageutil.ColorTableToRGBA(args["color"].(*golua.LTable))

			c := color.NRGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			}

			w.SetBgColor(c)
			return 0
		})

	/// @func window_should_close(id, v)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg v {bool}
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

	/// @func window_set_icon_imgscal(id, circled?)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg? circled {bool}
	/// @desc
	/// Uses the 16x16 and 32x32 imscal application icons for the window icon.
	/// Note the imgscal icon is light green with a transparent background,
	/// against a white background it will be hard to see.
	/// The circled versions have a dark background and are more readable.
	lib.CreateFunction(tab, "window_set_icon_imgscal",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.BOOL, Name: "circled", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			circled := args["circled"].(bool)
			var iconBytes [][]byte

			if !circled {
				iconBytes = [][]byte{
					assets.FAVICON_16x16,
					assets.FAVICON_32x32,
				}
			} else {
				iconBytes = [][]byte{
					assets.FAVICON_16x16_circle,
					assets.FAVICON_32x32_circle,
				}
			}

			icons := []image.Image{}

			for _, f := range iconBytes {
				ic, err := imageutil.Decode(bytes.NewReader(f), imageutil.ENCODING_PNG)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid image: %s", err), log.LEVEL_ERROR)), 0)
				}

				icons = append(icons, ic)
			}

			w.SetIcon(icons...)
			return 0
		})

	/// @func window_set_icon(id, img)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg img {int<collection.IMAGE>}
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

			r.IC.Schedule(state, args["icon_id"].(int), &collection.Task[collection.ItemImage]{
				Lib:  d.Lib,
				Name: d.Name,
				Fn: func(i *collection.Item[collection.ItemImage]) {
					w.SetIcon(imageutil.CopyImage(i.Self.Image, imageutil.MODEL_NRGBA))
				},
			})

			return 0
		})

	/// @func window_set_icon_many(id, imgs)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg imgs {[]int<collection.IMAGE>}
	/// @blocking
	/// @desc
	/// Setting multiple icons allows it select the closest to the system's desired size.
	lib.CreateFunction(tab, "window_set_icon_many",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			lua.ArgArray("icon_ids", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			imgids := args["icon_ids"].([]any)
			imgList := []image.Image{}
			wg := sync.WaitGroup{}

			for _, id := range imgids {
				wg.Add(1)
				r.IC.Schedule(state, id.(int), &collection.Task[collection.ItemImage]{
					Lib:  d.Lib,
					Name: d.Name,
					Fn: func(i *collection.Item[collection.ItemImage]) {
						imgList = append(imgList, imageutil.CopyImage(i.Self.Image, imageutil.MODEL_NRGBA))
						wg.Done()
					},
					Fail: func(i *collection.Item[collection.ItemImage]) {
						wg.Done()
					},
				})
			}

			wg.Wait()
			w.SetIcon(imgList...)

			return 0
		})

	/// @func window_clear_icon(id)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @desc
	/// Resets window icon to the default, same as 'window_set_icon_many(id, {})'.
	lib.CreateFunction(tab, "window_clear_icon",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetIcon()
			return 0
		})

	/// @func window_set_fps(id, fps)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg fps {int}
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

	/// @func window_set_title(id, title)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg title {string}
	lib.CreateFunction(tab, "window_set_title",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetTitle(args["title"].(string))
			return 0
		})

	/// @func window_register_keyboard_shortcuts(id, shortcuts)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg shortcuts {[]struct<gui.Shortcut>}
	lib.CreateFunction(tab, "window_register_keyboard_shortcuts",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "shortcuts"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			st := args["shortcuts"].(*golua.LTable)
			stList := []g.WindowShortcut{}
			for i := range st.Len() {
				s := st.RawGetInt(i + 1).(*golua.LTable)

				key := s.RawGetString("key").(golua.LNumber)
				mod := s.RawGetString("mod").(golua.LNumber)
				callback := s.RawGetString("callback")

				shortcut := g.WindowShortcut{
					Key:      g.Key(key),
					Modifier: g.Modifier(mod),
					Callback: func() {
						state.Push(callback)
						state.Call(0, 0)
					},
				}

				stList = append(stList, shortcut)
			}

			w.RegisterKeyboardShortcuts(stList...)

			return 0
		})

	/// @func window_set_close_callback(id, callback)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg callback {function() -> bool}
	lib.CreateFunction(tab, "window_set_close_callback",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "callback"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetCloseCallback(func() bool {
				state.Push(args["callback"].(*golua.LFunction))
				state.Call(0, 1)
				res := bool(state.ToBool(-1))
				state.Pop(1)
				return res
			})
			return 0
		})

	/// @func window_set_drop_callback(id, callback)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg callback {function([]string)}
	lib.CreateFunction(tab, "window_set_drop_callback",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "callback"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetDropCallback(func(items []string) {
				state.Push(args["callback"].(*golua.LFunction))
				t := state.NewTable()
				for _, s := range items {
					t.Append(golua.LString(s))
				}
				state.Push(t)
				state.Call(1, 0)
			})
			return 0
		})

	/// @func window_additional_input_handler_callback(id, callback)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg callback {function(key int<gui.Key>, mod int<gui.Key>, action int<gui.Action>)}
	lib.CreateFunction(tab, "window_additional_input_handler_callback",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.FUNC, Name: "callback"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.SetAdditionalInputHandlerCallback(func(k g.Key, m g.Modifier, a g.Action) {
				state.Push(args["callback"].(*golua.LFunction))
				state.Push(golua.LNumber(k))
				state.Push(golua.LNumber(m))
				state.Push(golua.LNumber(a))
				state.Call(3, 0)
			})

			return 0
		})

	/// @func window_close(id)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @desc
	/// Same as 'window_should_close(id, true)'.
	lib.CreateFunction(tab, "window_close",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			w.Close()
			return 0
		})

	/// @func window_run(id, fn)
	/// @arg id {int<collection.CRATE_WINDOW>}
	/// @arg fn {function()}
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
				state.Push(args["fn"].(*golua.LFunction))
				state.Call(0, 0)
			})

			return 0
		})

	/// @func window_single() -> struct<gui.WidgetWindow>
	/// @returns {struct<gui.WidgetWindow>}
	lib.CreateFunction(tab, "window_single",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, true, false, "")

			state.Push(win)
			return 1
		})

	/// @func window_single_with_menu_bar() -> struct<gui.WidgetWindow>
	/// @returns {struct<gui.WidgetWindow>}
	lib.CreateFunction(tab, "window_single_with_menu_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, true, true, "")

			state.Push(win)
			return 1
		})

	/// @func window() -> struct<gui.WidgetWindow>
	/// @arg title {string}
	/// @returns {struct<gui.WidgetWindow>}
	lib.CreateFunction(tab, "window",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			win := windowTable(r, lg, state, false, false, args["title"].(string))

			state.Push(win)
			return 1
		})

	/// @func layout(widgets)
	/// @arg widgets {[]struct<gui.Widget>}
	/// @desc
	/// Builds a list of widgets when called.
	lib.CreateFunction(tab, "layout",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "widgets", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			widgets := args["widgets"].(*golua.LTable)
			layout := g.Layout(layoutBuild(r, state, parseWidgets(parseTable(widgets), state, lg), lg))
			layout.Build()

			return 0
		})

	/// @func popup_open(name)
	/// @arg name {string}
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

	/// @func prepare_msg_box() -> struct<gui.WidgetMSGBoxPrepare>
	/// @returns {struct<gui.WidgetMSGBoxPrepare>}
	lib.CreateFunction(tab, "prepare_msg_box",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := msgBoxPrepareTable(state)

			state.Push(t)
			return 1
		})

	/// @func style_var_is_vec2(var) -> bool
	/// @arg var {int<gui.StyleVarID>}
	/// @returns {bool}
	lib.CreateFunction(tab, "style_var_is_vec2",
		[]lua.Arg{
			{Type: lua.INT, Name: "var"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["var"].(int)
			sv := g.StyleVarID(v)
			b := sv.IsVec2()

			state.Push(golua.LBool(b))
			return 1
		})

	/// @func style_var_string(var) -> string
	/// @arg var {int<gui.StyleVarID>}
	/// @returns {string}
	lib.CreateFunction(tab, "style_var_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "var"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["var"].(int)
			sv := g.StyleVarID(v)
			s := sv.String()

			state.Push(golua.LString(s))
			return 1
		})

	/// @func style_var_from_string(str) -> int<gui.StyleVarID>
	/// @arg str {string}
	/// @returns {int<gui.StyleVarID>}
	lib.CreateFunction(tab, "style_var_from_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "s"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			v := args["s"].(string)
			sv := g.StyleVarIDFromString(v)

			state.Push(golua.LNumber(sv))
			return 1
		})

	/// @func shortcut(key, mod, callback) -> struct<gui.Shortcut>
	/// @arg key {int<gui.Key>}
	/// @arg mod {int<gui.Key>}
	/// @arg callback {function()}
	/// @returns {struct<gui.Shortcut>}
	lib.CreateFunction(tab, "shortcut",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
			{Type: lua.INT, Name: "mod"},
			{Type: lua.FUNC, Name: "callback"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Shortcut
			/// @prop key {int<gui.Key>}
			/// @prop mod {int<gui.Key>}
			/// @prop callback {function()}

			key := args["key"].(int)
			mod := args["mod"].(int)
			callback := args["callback"].(*golua.LFunction)

			t := state.NewTable()
			t.RawSetString("key", golua.LNumber(key))
			t.RawSetString("mod", golua.LNumber(mod))
			t.RawSetString("callback", callback)

			state.Push(t)
			return 1
		})

	/// @func plot_ticker(position, label) -> struct<gui.PlotTicker>
	/// @arg position {float}
	/// @arg label {string}
	/// @returns {struct<gui.PlotTicker>}
	lib.CreateFunction(tab, "plot_ticker",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "position"},
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pos := args["position"].(float64)
			label := args["label"].(string)

			t := state.NewTable()
			t.RawSetString("position", golua.LNumber(pos))
			t.RawSetString("label", golua.LString(label))

			state.Push(t)
			return 1
		})

	/// @func css_parse(path)
	/// @arg path {string}
	lib.CreateFunction(tab, "css_parse",
		[]lua.Arg{
			{Type: lua.STRING, Name: "path"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pth := args["path"].(string)
			b, err := os.ReadFile(pth)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to read css file: %s with error: %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			err = g.ParseCSSStyleSheet(b)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("failed to parse css file: %s with error: %s", pth, err), log.LEVEL_ERROR)), 0)
			}

			return 0
		})

	/// @func align_text_to_frame_padding()
	lib.CreateFunction(tab, "align_text_to_frame_padding",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.AlignTextToFramePadding()
			return 0
		})

	/// @func calc_text_size(text) -> float, float
	/// @arg text {string}
	/// @returns {float} - Width of the text string.
	/// @returns {float} - Height of the text string.
	lib.CreateFunction(tab, "calc_text_size",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.CalcTextSize(args["text"].(string))

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func calc_text_size_width(text) -> float
	/// @arg text {string}
	/// @returns {float} - Width of the text string.
	lib.CreateFunction(tab, "calc_text_size_width",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, _ := g.CalcTextSize(args["text"].(string))

			state.Push(golua.LNumber(width))
			return 1
		})

	/// @func calc_text_size_height(text) -> float
	/// @arg text {string}
	/// @returns {float} - Height of the text string.
	lib.CreateFunction(tab, "calc_text_size_height",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, height := g.CalcTextSize(args["text"].(string))

			state.Push(golua.LNumber(height))
			return 1
		})

	/// @func calc_text_size_v(text, hideAfterDoubleHash, wrapWidth) -> float, float
	/// @arg text {string}
	/// @arg hideAfterDoubleHash {bool}
	/// @arg wrapWidth {float}
	/// @returns {float} - Width of the text string.
	/// @returns {float} - Height of the text string.
	lib.CreateFunction(tab, "calc_text_size_v",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.BOOL, Name: "hideAfterDoubleHash"},
			{Type: lua.FLOAT, Name: "wrapWidth"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.CalcTextSizeV(args["text"].(string), args["hideAfterDoubleHash"].(bool), float32(args["wrapWidth"].(float64)))

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func available_region() -> float, float
	/// @returns {float} - Width.
	/// @returns {float} - Height.
	lib.CreateFunction(tab, "available_region",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetAvailableRegion()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func frame_padding() -> float, float
	/// @returns {float} - X padding.
	/// @returns {float} - Y Padding.
	lib.CreateFunction(tab, "frame_padding",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x, y := g.GetFramePadding()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func item_inner_spacing() -> float, float
	/// @returns {float} - Width.
	/// @returns {float} - Height.
	lib.CreateFunction(tab, "item_inner_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetItemInnerSpacing()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func item_spacing() -> float, float
	/// @returns {float} - Width.
	/// @returns {float} - Height.
	lib.CreateFunction(tab, "item_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetItemSpacing()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func mouse_pos_xy() -> int, int
	/// @returns {int} - Mouse x position.
	/// @returns {int} - Mouse y position.
	lib.CreateFunction(tab, "mouse_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetMousePos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func mouse_pos() -> struct<image.Point>
	/// @returns {struct<image.Point>}
	lib.CreateFunction(tab, "mouse_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetMousePos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func window_padding() -> float, float
	/// @returns {float} X padding.
	/// @returns {float} Y padding.
	lib.CreateFunction(tab, "window_padding",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x, y := g.GetWindowPadding()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func is_item_active() -> bool
	/// @returns {bool}
	lib.CreateFunction(tab, "is_item_active",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemActive()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_item_clicked(button) -> bool
	/// @arg button {int<gui.MouseButton>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_item_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_item_hovered() -> bool
	/// @returns {bool}
	lib.CreateFunction(tab, "is_item_hovered",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemHovered()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_down(key) -> bool
	/// @arg key {int<gui.Key>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_key_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyDown(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_pressed(key) -> bool
	/// @arg key {int<gui.Key>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_key_pressed",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyPressed(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_released(key) -> bool
	/// @arg key {int<gui.Key>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_key_released",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyReleased(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_clicked(button) -> bool
	/// @arg button {int<gui.MouseButton>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_mouse_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_double_clicked(button) -> bool
	/// @arg button {int<gui.MouseButton>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_mouse_double_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseDoubleClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_down(button) -> bool
	/// @arg button {int<gui.MouseButton>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_mouse_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseDown(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_released(button) -> bool
	/// @arg button {int<gui.MouseButton>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_mouse_released",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseReleased(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_appearing() -> bool
	/// @returns {bool}
	lib.CreateFunction(tab, "is_window_appearing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowAppearing()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_collapsed() -> bool
	/// @returns {bool}
	lib.CreateFunction(tab, "is_window_collapsed",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowCollapsed()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_focused(flags) -> bool
	/// @arg flags {int<gui.FocusedFlags>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_window_focused",
		[]lua.Arg{
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowFocused(g.FocusedFlags(args["flags"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_hovered(flags) -> bool
	/// @arg flags {int<gui.HoveredFlags>}
	/// @returns {bool}
	lib.CreateFunction(tab, "is_window_hovered",
		[]lua.Arg{
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowHovered(g.HoveredFlags(args["flags"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func open_url(url)
	/// @arg url {string}
	lib.CreateFunction(tab, "open_url",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.OpenURL(args["url"].(string))
			return 0
		})

	/// @func pop_clip_rect()
	lib.CreateFunction(tab, "pop_clip_rect",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopClipRect()
			return 0
		})

	/// @func pop_font()
	lib.CreateFunction(tab, "pop_font",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopFont()
			return 0
		})

	/// @func pop_item_width()
	lib.CreateFunction(tab, "pop_item_width",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopItemWidth()
			return 0
		})

	/// @func pop_style()
	lib.CreateFunction(tab, "pop_style",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopStyle()
			return 0
		})

	/// @func pop_style_color()
	lib.CreateFunction(tab, "pop_style_color",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopStyleColor()
			return 0
		})

	/// @func pop_style_color_v(count)
	/// @arg count {int}
	lib.CreateFunction(tab, "pop_style_color_v",
		[]lua.Arg{
			{Type: lua.INT, Name: "count"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopStyleColorV(args["count"].(int))
			return 0
		})

	/// @func pop_style_v(count)
	/// @arg count {int}
	lib.CreateFunction(tab, "pop_style_v",
		[]lua.Arg{
			{Type: lua.INT, Name: "count"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopStyleV(args["count"].(int))
			return 0
		})

	/// @func pop_text_wrap_pos()
	lib.CreateFunction(tab, "pop_text_wrap_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopTextWrapPos()
			return 0
		})

	/// @func push_button_text_align(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "push_button_text_align",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushButtonTextAlign(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_clip_rect(min, max, intersect)
	/// @arg min {struct<image.Point>}
	/// @arg max {struct<image.Point>}
	/// @arg intersect {bool}
	lib.CreateFunction(tab, "push_clip_rect",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "min"},
			{Type: lua.RAW_TABLE, Name: "max"},
			{Type: lua.BOOL, Name: "intersect"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			min := imageutil.TableToPoint(args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(args["max"].(*golua.LTable))
			g.PushClipRect(min, max, args["intersect"].(bool))
			return 0
		})

	/// @func push_color_button(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_button",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButton(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_button_active(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_button_active",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButtonActive(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_button_hovered(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_button_hovered",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButtonHovered(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_frame_bg(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_frame_bg",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorFrameBg(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_text(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_text",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorText(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_text_disabled(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_text_disabled",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorTextDisabled(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_window_bg(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_color_window_bg",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorWindowBg(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_font(fontref) -> bool
	/// @arg fontref {int<ref.FONT>}
	/// @returns {bool}
	lib.CreateFunction(tab, "push_font",
		[]lua.Arg{
			{Type: lua.INT, Name: "fontref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ref := args["fontref"].(int)
			sref, err := r.CR_REF.Item(ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			font := sref.Value.(*g.FontInfo)

			ok := g.PushFont(font)

			state.Push(golua.LBool(ok))
			return 1
		})

	/// @func push_frame_padding(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "push_frame_padding",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushFramePadding(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_item_spacing(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "push_item_spacing",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushItemSpacing(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_item_width(width)
	/// @arg width {float}
	lib.CreateFunction(tab, "push_item_width",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushItemWidth(float32(args["width"].(float64)))
			return 0
		})

	/// @func push_selectable_text_align(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "push_selectable_text_align",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushSelectableTextAlign(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_style_color(id, color)
	/// @arg id {int<gui.StyleColorID>}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "push_style_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushStyleColor(g.StyleColorID(args["id"].(int)), imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_text_wrap_pos()
	lib.CreateFunction(tab, "push_text_wrap_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushTextWrapPos()
			return 0
		})

	/// @func push_window_padding(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "push_window_padding",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushWindowPadding(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func same_line()
	lib.CreateFunction(tab, "same_line",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SameLine()
			return 0
		})

	/// @func cursor_pos_set_xy(x, y)
	/// @arg x {int}
	/// @arg y {int}
	lib.CreateFunction(tab, "cursor_pos_set_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorPos(image.Point{
				X: args["x"].(int),
				Y: args["y"].(int),
			})
			return 0
		})

	/// @func cursor_pos_set(point)
	/// @arg point {struct<image.Point>}
	lib.CreateFunction(tab, "cursor_pos_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "point"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorPos(imageutil.TableToPoint(args["point"].(*golua.LTable)))
			return 0
		})

	/// @func cursor_screeen_pos_set_xy(x, y)
	/// @arg x {int}
	/// @arg y {int}
	lib.CreateFunction(tab, "cursor_screen_pos_set_xy",
		[]lua.Arg{
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorScreenPos(image.Point{
				X: args["x"].(int),
				Y: args["y"].(int),
			})
			return 0
		})

	/// @func cursor_screen_pos_set(point)
	/// @arg point {struct<image.Point>}
	lib.CreateFunction(tab, "cursor_screen_pos_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "point"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorScreenPos(imageutil.TableToPoint(args["point"].(*golua.LTable)))
			return 0
		})

	/// @func item_default_focus_set()
	lib.CreateFunction(tab, "item_default_focus_set",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetItemDefaultFocus()
			return 0
		})

	/// @func keyboard_focus_here()
	lib.CreateFunction(tab, "keyboard_focus_here",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetKeyboardFocusHere()
			return 0
		})

	/// @func keyboard_focus_here_v(i)
	/// @arg i {int} - Widget offset, e.g. -1 is the previous widget.
	lib.CreateFunction(tab, "keyboard_focus_here_v",
		[]lua.Arg{
			{Type: lua.INT, Name: "i"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetKeyboardFocusHereV(args["i"].(int))
			return 0
		})

	/// @func mouse_cursor_set(cursor)
	/// @arg cursor {int<gui.Cursor>}
	lib.CreateFunction(tab, "mouse_cursor_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "cursor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetMouseCursor(g.MouseCursorType(args["cursor"].(int)))
			return 0
		})

	/// @func next_window_pos_set(x, y)
	/// @arg x {int}
	/// @arg y {int}
	lib.CreateFunction(tab, "next_window_pos_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetNextWindowPos(float32(args["x"].(float64)), float32(args["y"].(float64)))
			return 0
		})

	/// @func next_window_size_set(width, height)
	/// @arg width {float}
	/// @arg height {float}
	lib.CreateFunction(tab, "next_window_size_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetNextWindowSize(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func next_window_size_v_set(width, height, cond)
	/// @arg width {float}
	/// @arg height {float}
	/// @arg cond {int<gui.Condition>}
	lib.CreateFunction(tab, "next_window_size_v_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
			{Type: lua.INT, Name: "cond"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetNextWindowSizeV(float32(args["width"].(float64)), float32(args["height"].(float64)), g.ExecCondition(args["cond"].(int)))
			return 0
		})

	/// @func update()
	lib.CreateFunction(tab, "update",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.Update()
			return 0
		})

	/// @func color_to_uint32(color) -> int
	/// @arg color {struct<image.Color>}
	/// @returns {int} - Number representation of a color.
	/// @desc
	/// Returns the uint32 to lua as a float64 (The type lua uses for numbers).
	lib.CreateFunction(tab, "color_to_uint32",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cuint := g.ColorToUint(imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable)))

			state.Push(golua.LNumber(cuint))
			return 1
		})

	/// @func uint32_to_color(ucolor) -> struct<image.ColorRGBA>
	/// @arg ucolor {int}
	/// @returns {struct<image.ColorRGBA>}
	lib.CreateFunction(tab, "uint32_to_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "ucolor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.UintToColor(uint32(args["ucolor"].(int)))
			state.Push(imageutil.RGBAColorToColorTable(state, c))
			return 1
		})

	/// @func wg_label(text) -> struct<gui.WidgetLabel>
	/// @arg text {string}
	/// @returns {struct<gui.WidgetLabel>}
	lib.CreateFunction(tab, "wg_label",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := labelTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_number(number) -> struct<gui.WidgetLabel>
	/// @arg number {float}
	/// @returns {struct<gui.WidgetLabel>}
	/// @desc
	/// Converts the number into a string before creating the label widget.
	lib.CreateFunction(tab, "wg_number",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "number"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := labelTable(state, fmt.Sprintf("%v", args["number"].(float64)))

			state.Push(t)
			return 1
		})

	/// @func wg_button(text) -> struct<gui.WidgetButton>
	/// @arg text {string}
	/// @returns {struct<gui.WidgetButton>}
	lib.CreateFunction(tab, "wg_button",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_dummy(width, height) -> struct<gui.WidgetDummy>
	/// @arg width {float}
	/// @arg height {float}
	/// @returns {struct<gui.WidgetDummy>}
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

	/// @func wg_separator() -> struct<gui.WidgetSeparator>
	/// @returns {struct<gui.WidgetSeparator>}
	lib.CreateFunction(tab, "wg_separator",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := separatorTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_bullet_text(text) -> struct<gui.WidgetBulletText>
	/// @arg text {string}
	/// @returns {struct<gui.WidgetBulletText>}
	lib.CreateFunction(tab, "wg_bullet_text",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := bulletTextTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_bullet() -> struct<gui.WidgetBullet>
	/// @returns {struct<gui.WidgetBullet>}
	lib.CreateFunction(tab, "wg_bullet",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := bulletTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_checkbox(text, boolref) -> struct<gui.WidgetCheckbox>
	/// @arg text {string}
	/// @arg boolref {int<ref.BOOL>}
	/// @returns {struct<gui.WidgetCheckbox>}
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

	/// @func wg_child() -> struct<gui.WidgetChild>
	/// @returns {struct<gui.WidgetChild>}
	lib.CreateFunction(tab, "wg_child",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := childTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_color_edit(text, colorref) -> struct<gui.WidgetColorEdit>
	/// @arg text {string}
	/// @arg colorref {int<ref.RGBA>}
	/// @returns {struct<gui.WidgetColorEdit>}
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

	/// @func wg_column(widgets?) -> struct<gui.WidgetColumn>
	/// @arg? widgets {[]struct<gui.Widget>}
	/// @returns {struct<gui.WidgetColumn>}
	lib.CreateFunction(tab, "wg_column",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "widgets", Optional: true},
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

	/// @func wg_row(widgets?) -> struct<gui.WidgetRow>
	/// @arg? widgets {[]struct<gui.Widget>}
	/// @returns {struct<gui.WidgetRow>}
	lib.CreateFunction(tab, "wg_row",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "widgets", Optional: true},
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

	/// @func wg_combo_custom(text, preview) -> struct<gui.WidgetComboCustom>
	/// @arg text {string}
	/// @arg preview {string}
	/// @returns {struct<gui.WidgetComboCustom>}
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

	/// @func wg_combo(text, preview, items, i32ref) -> struct<gui.WidgetCombo>
	/// @arg text {string}
	/// @arg preview {string}
	/// @arg items {[]string}
	/// @arg i32ref {int<ref.INT32>}
	/// @returns {struct<gui.WidgetCombo>}
	lib.CreateFunction(tab, "wg_combo",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.STRING, Name: "preview"},
			{Type: lua.RAW_TABLE, Name: "items"},
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

	/// @func wg_combo_preview(text, items, i32ref) -> struct<gui.WidgetCombo>
	/// @arg text {string}
	/// @arg items {[]string}
	/// @arg i32ref {int<ref.INT32>}
	/// @returns {struct<gui.WidgetCombo>}
	/// @desc
	/// Same as wg_combo but sets preview to the selected value in items.
	lib.CreateFunction(tab, "wg_combo_preview",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.RAW_TABLE, Name: "items"},
			{Type: lua.INT, Name: "i32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			text := args["text"].(string)
			items := args["items"].(*golua.LTable)
			i32ref := args["i32ref"].(int)

			sref, err := r.CR_REF.Item(i32ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			selected := sref.Value.(*int32)
			sel := int(*selected)
			preview := items.RawGetInt(sel + 1).(golua.LString)

			t := comboTable(state, text, string(preview), items, i32ref)

			state.Push(t)
			return 1
		})

	/// @func wg_condition(condition, widgetIf, widgetElse) -> struct<gui.WidgetCondition>
	/// @arg condition {bool}
	/// @arg widgetIf {struct<gui.Widget>}
	/// @arg widgetElse {struct<gui.Widget>}
	/// @returns {struct<gui.WidgetCondition>}
	lib.CreateFunction(tab, "wg_condition",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "condition"},
			{Type: lua.RAW_TABLE, Name: "widgetIf"},
			{Type: lua.RAW_TABLE, Name: "widgetElse"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			widgetIf := args["widgetIf"].(golua.LValue)
			widgetElse := args["widgetElse"].(golua.LValue)
			t := conditionTable(state, args["condition"].(bool), widgetIf, widgetElse)

			state.Push(t)
			return 1
		})

	/// @func wg_context_menu() -> struct<gui.WidgetContextMenu>
	/// @returns {struct<gui.WidgetContextMenu>}
	lib.CreateFunction(tab, "wg_context_menu",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := contextMenuTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_date_picker(id, timeref) -> struct<gui.WidgetDatePicker>
	/// @arg id {string}
	/// @arg timeref {int<ref.TIME>}
	/// @returns {struct<gui.WidgetDatePicker>}
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

	/// @func wg_drag_int(label, i32ref, minvalue, maxvalue) -> struct<gui.WidgetDragInt>
	/// @arg label {string}
	/// @arg i32ref {int<ref.INT32>}
	/// @arg minvalue {int}
	/// @arg maxvalue {int}
	/// @returns {struct<gui.WidgetDragInt>}
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

	/// @func wg_input_float(f32ref) -> struct<gui.WidgetInputFloat>
	/// @arg f32ref {int<ref.FLOAT32>}
	/// @returns {struct<gui.WidgetInputFloat>}
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

	/// @func wg_input_int(i32ref) -> struct<gui.WidgetInputInt>
	/// @arg i32ref {int<ref.INT32>}
	/// @returns {struct<gui.WidgetInputInt>}
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

	/// @func wg_input_text(strref) -> struct<gui.WidgetInputText>
	/// @arg strref {int<ref.STRING>}
	/// @returns {struct<gui.WidgetInputText>}
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

	/// @func wg_input_text_multiline(strref) -> struct<gui.WidgetInputTextMultiline>
	/// @arg strref {int<ref.STRING>}
	/// @returns {struct<gui.WidgetInputTextMultiline>}
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

	/// @func wg_progress_bar(fraction) -> struct<gui.WidgetProgressBar>
	/// @arg fraction {float}
	/// @returns {struct<gui.WidgetProgressBar>}
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

	/// @func wg_progress_indicator(label, width, height, radius) -> struct<gui.WidgetProgressIndicator>
	/// @arg label {string}
	/// @arg width {float}
	/// @arg height {float}
	/// @arg radius {float}
	/// @returns {struct<gui.WidgetProgressIndicator>}
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

	/// @func wg_spacing() -> struct<gui.WidgetSpacing>
	/// @returns {struct<gui.WidgetSpacing>}
	lib.CreateFunction(tab, "wg_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := spacingTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_button_small(text) -> struct<gui.WidgetButtonSmall>
	/// @arg text {string}
	/// @returns {struct<gui.WidgetButtonSmall>}
	lib.CreateFunction(tab, "wg_button_small",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonSmallTable(state, args["text"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_button_radio(text, active) -> struct<gui.WidgetButtonRadio>
	/// @arg text {string}
	/// @arg active {bool}
	/// @returns {struct<gui.WidgetButtonRadio>}
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

	/// @func wg_image_url(url) -> struct<gui.WidgetImageURL>
	/// @arg url {string}
	/// @returns struct<gui.WidgetImageURL>
	lib.CreateFunction(tab, "wg_image_url",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageUrlTable(state, args["url"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_image(id) -> struct<gui.WidgetImage>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<gui.WidgetImage>}
	/// @blocking
	lib.CreateFunction(tab, "wg_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageTable(state, args["id"].(int), false, false)

			state.Push(t)
			return 1
		})

	/// @func wg_image_sync(id) -> struct<gui.WidgetImage>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<gui.WidgetImage>}
	/// @desc
	/// Note: this does not wait for the image to be ready or idle,
	/// if the image is not loaded it will dislay an empy image.
	/// May look weird if the image is also being processed while displayed here.
	lib.CreateFunction(tab, "wg_image_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageTable(state, args["id"].(int), true, false)

			state.Push(t)
			return 1
		})

	/// @func wg_image_cached(id) -> struct<gui.WidgetImage>
	/// @arg id {int<collection.CRATE_CACHEDIMAGE>}
	/// @returns {struct<gui.WidgetImage>}
	lib.CreateFunction(tab, "wg_image_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := imageTable(state, args["id"].(int), false, true)

			state.Push(t)
			return 1
		})

	/// @func wg_list_box(items) -> struct<gui.WidgetListbox>
	/// @arg items {[]string}
	/// @returns {struct<gui.WidgetListbox>}
	lib.CreateFunction(tab, "wg_list_box",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "items"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := listBoxTable(state, args["items"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_list_clipper() -> struct<gui.WidgetListClipper>
	/// @returns {struct<gui.WidgetListClipper>}
	lib.CreateFunction(tab, "wg_list_clipper",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := listClipperTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_bar_main() -> struct<gui.WidgetMenuBarMain>
	/// @returns {struct<gui.WidgetMenuBarMain>}
	lib.CreateFunction(tab, "wg_menu_bar_main",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := mainMenuBarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_bar() -> struct<gui.WidgetMenuBar>
	/// @returns {struct<gui.WidgetMenuBar>}
	lib.CreateFunction(tab, "wg_menu_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuBarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_menu_item(label) -> struct<gui.WidgetMenuItem>
	/// @arg label {string}
	/// @returns {struct<gui.WidgetMenuItem>}
	lib.CreateFunction(tab, "wg_menu_item",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuItemTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_menu(label) -> struct<gui.WidgetMenu>
	/// @arg label {string}
	/// @returns {struct<gui.WidgetMenu>}
	lib.CreateFunction(tab, "wg_menu",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := menuTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_selectable(label) -> struct<gui.WidgetSelectable>
	/// @arg label {string}
	/// @returns {struct<gui.WidgetSelectable>}
	lib.CreateFunction(tab, "wg_selectable",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := selectableTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_slider_float(f32ref, min, max) -> struct<gui.WidgetSliderFloat>
	/// @arg f32ref {int<ref.FLOAT32>}
	/// @arg min {float}
	/// @arg max {float}
	/// @returns {struct<gui.WidgetSliderFloat>}
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

	/// @func wg_slider_int(i32ref, min, max) -> struct<gui.WidgetSliderInt>
	/// @arg i32ref {int<ref.INT32>}
	/// @arg min {int}
	/// @arg max {int}
	/// @returns {struct<gui.WidgetSliderInt>}
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

	/// @func wg_vslider_int(i32ref, min, max) -> struct<gui.WidgetVSliderInt>
	/// @arg i32ref {int<ref.INT32>}
	/// @arg min {int}
	/// @arg max {int}
	/// @returns {struct<gui.WidgetVSliderInt>}
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

	/// @func wg_tab_bar() -> struct<gui.WidgetTabBar>
	/// @returns {struct<gui.WidgetTabBar>}
	lib.CreateFunction(tab, "wg_tab_bar",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tabbarTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_tab_item(label) -> struct<gui.TabItem>
	/// @arg label {string}
	/// @returns {struct<gui.TabItem>}
	lib.CreateFunction(tab, "wg_tab_item",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tabitemTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_tooltip(tip) -> struct<gui.WidgetTooltip>
	/// @arg tip {string}
	/// @returns {struct<gui.WidgetTooltip>}
	lib.CreateFunction(tab, "wg_tooltip",
		[]lua.Arg{
			{Type: lua.STRING, Name: "tip"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tooltipTable(state, args["tip"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_table_column(label) -> struct<gui.TableColumn>
	/// @arg label {string}
	/// @returns {struct<gui.TableColumn>}
	lib.CreateFunction(tab, "wg_table_column",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tableColumnTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_table_row(widgets?) -> struct<gui.TableRow>
	/// @arg? widgets {[]struct<gui.Widget>}
	/// @returns {struct<gui.TableRow>}
	lib.CreateFunction(tab, "wg_table_row",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "widgets", Optional: true},
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

	/// @func wg_table() -> struct<gui.WidgetTable>
	/// @returns {struct<gui.WidgetTable>}
	lib.CreateFunction(tab, "wg_table",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tableTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_button_arrow(dir) -> {struct<gui.WidgetButtonArrow>}
	/// @arg dir {int<gui.Direction>}
	/// @returns {struct<gui.WidgetButtonArrow>}
	lib.CreateFunction(tab, "wg_button_arrow",
		[]lua.Arg{
			{Type: lua.INT, Name: "dir"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonArrowTable(state, args["dir"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_tree_node(label) -> struct<gui.WidgetTreeNode>
	/// @arg label {string}
	/// @returns {struct<gui.WidgetTreeNode>}
	lib.CreateFunction(tab, "wg_tree_node",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := treeNodeTable(state, args["label"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_tree_table_row(label, widgets?) -> struct<gui.TreeTableRow>
	/// @arg label {string}
	/// @arg? widgets {[]struct<gui.Widget>}
	/// @returns {struct<gui.TreeTableRow>}
	lib.CreateFunction(tab, "wg_tree_table_row",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
			{Type: lua.RAW_TABLE, Name: "widgets", Optional: true},
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

	/// @func wg_tree_table() -> struct<gui.WidgetTreeTable>
	/// @returns {struct<gui.WidgetTreeTable>}
	lib.CreateFunction(tab, "wg_tree_table",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := treeTableTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_popup_modal(name) -> struct<gui.WidgetPopupModal>
	/// @arg name {string}
	/// @returns {struct<gui.WidgetPopupModal>}
	lib.CreateFunction(tab, "wg_popup_modal",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := popupModalTable(state, args["name"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_popup(name) -> struct<gui.WidgetPopup>
	/// @arg name {string}
	/// @returns {struct<gui.WidgetPopup>}
	lib.CreateFunction(tab, "wg_popup",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := popupTable(state, args["name"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_layout_split(direction, f32ref, layout1, layout2) -> struct<gui.WidgetLayoutSplit>
	/// @arg direction {int<gui.Direction>}
	/// @arg f32ref {int<ref.FLOAT32>}
	/// @arg layout1 {[]struct<gui.Widget>}
	/// @arg layout2 {[]struct<gui.Widget>}
	/// @returns {struct<gui.WidgetLayoutSplit>}
	lib.CreateFunction(tab, "wg_layout_split",
		[]lua.Arg{
			{Type: lua.INT, Name: "direction"},
			{Type: lua.INT, Name: "f32ref"},
			{Type: lua.RAW_TABLE, Name: "layout1"},
			{Type: lua.RAW_TABLE, Name: "layout2"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			layout1 := args["layout1"].(golua.LValue)
			layout2 := args["layout2"].(golua.LValue)
			t := splitLayoutTable(state, args["direction"].(int), args["f32ref"].(int), layout1, layout2)

			state.Push(t)
			return 1
		})

	/// @func wg_splitter(direction, f32ref) -> struct<gui.WidgetSplitter>
	/// @arg direction {int<gui.Direction>}
	/// @arg f32ref {int<ref.FLOAT32>}
	/// @returns {struct<gui.WidgetSplitter>}
	lib.CreateFunction(tab, "wg_splitter",
		[]lua.Arg{
			{Type: lua.INT, Name: "direction"},
			{Type: lua.INT, Name: "f32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := splitterTable(state, args["direction"].(int), args["f32ref"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_stack(visible, widgets) -> struct<gui.WidgetStack>
	/// @arg visible {int} - The index in widgets that is visible.
	/// @arg widgets {[]struct<gui.Widget>}
	/// @returns {struct<gui.WidgetStack>}
	lib.CreateFunction(tab, "wg_stack",
		[]lua.Arg{
			{Type: lua.INT, Name: "visible"},
			{Type: lua.RAW_TABLE, Name: "widgets"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := stackTable(state, args["visible"].(int), args["widgets"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_align(at) -> struct<gui.WidgetAlign>
	/// @arg at {int<gui.Alignment>}
	/// @returns {struct<gui.WidgetAlign>}
	lib.CreateFunction(tab, "wg_align",
		[]lua.Arg{
			{Type: lua.INT, Name: "at"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := alignTable(state, args["at"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_msg_box(title, content) -> struct<gui.WidgetMSGBox>
	/// @arg title {string}
	/// @arg content {string}
	/// @returns {struct<gui.WidgetMSGBox>}
	/// @desc
	/// There must be a call to 'prepare_msg_box()' once a loop when using a msg box.
	lib.CreateFunction(tab, "wg_msg_box",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := msgBoxTable(state, args["title"].(string), args["content"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_button_invisible() -> struct<gui.WidgetButtonInvisible>
	/// @returns {struct<gui.WidgetButtonInvisible>}
	lib.CreateFunction(tab, "wg_button_invisible",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonInvisibleTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_button_image(id) -> struct<gui.WidgetButtonImage>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<gui.WidgetButtonImage>}
	lib.CreateFunction(tab, "wg_button_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonImageTable(state, args["id"].(int), false, false)

			state.Push(t)
			return 1
		})

	/// @func wg_button_image_sync(id) -> struct<gui.WidgetButtonImage>
	/// @arg id {int<collection.IMAGE>}
	/// @returns {struct<gui.WidgetButtonImage>}
	/// @desc
	/// Note: this does not wait for the image to be ready or idle,
	/// if the image is not loaded it will dislay an empy image.
	/// May look weird if the image is also being processed while displayed here.
	lib.CreateFunction(tab, "wg_button_image_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonImageTable(state, args["id"].(int), true, false)

			state.Push(t)
			return 1
		})

	/// @func wg_button_image_cached(id) -> struct<gui.WidgetButtonImage>
	/// @arg id {int<collection.CRATE_CACHEDIMAGE>}
	/// @returns {struct<gui.WidgetButtonImage>}
	lib.CreateFunction(tab, "wg_button_image_cached",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonImageTable(state, args["id"].(int), false, true)

			state.Push(t)
			return 1
		})

	/// @func wg_style() -> struct<gui.WidgetStyle>
	/// @returns {struct<gui.WidgetStyle>}
	lib.CreateFunction(tab, "wg_style",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := styleTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_custom(builder) -> struct<gui.WidgetCustom>
	/// @arg builder {function()}
	/// @returns {struct<gui.WidgetCustom>}
	lib.CreateFunction(tab, "wg_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "builder"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := customTable(state, args["builder"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func wg_event() -> struct<gui.WidgetEvent>
	/// @returns {struct<gui.WidgetEvent>}
	lib.CreateFunction(tab, "wg_event",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := eventHandlerTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_css_tag(tag) -> struct<gui.WidgetCSSTag>
	/// @arg tag {string}
	/// @returns {struct<gui.WidgetCSSTag>}
	lib.CreateFunction(tab, "wg_css_tag",
		[]lua.Arg{
			{Type: lua.STRING, Name: "tag"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cssTagTable(state, args["tag"].(string))

			state.Push(t)
			return 1
		})

	/// @func wg_plot(title) -> struct<gui.WidgetPlot>
	/// @arg title {string}
	/// @returns {struct<gui.WidgetPlot>}
	lib.CreateFunction(tab, "wg_plot",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotTable(state, args["title"].(string))

			state.Push(t)
			return 1
		})

	/// @func pt_bar_h(title, data) -> struct<gui.PlotBarH>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<gui.PlotBarH>}
	lib.CreateFunction(tab, "pt_bar_h",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarHTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_bar(title, data) -> struct<gui.PlotBar>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<gui.PlotBar>}
	lib.CreateFunction(tab, "pt_bar",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line(title, data) -> struct<gui.PlotLine>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<gui.PlotLine>}
	lib.CreateFunction(tab, "pt_line",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line_xy(title, xdata, ydata) -> struct<gui.PlotLineXY>
	/// @arg title {string}
	/// @arg xdata {[]float}
	/// @arg ydata {[]float}
	/// @returns {struct<gui.PlotLineXY>}
	lib.CreateFunction(tab, "pt_line_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "xdata"},
			{Type: lua.RAW_TABLE, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_pie_chart(labels, data, x, y, radius) -> struct<gui.PlotPieChart>
	/// @arg labels {[]string}
	/// @arg data {[]float}
	/// @arg x {float}
	/// @arg y {float}
	/// @arg radius {float}
	/// @returns {struct<gui.PlotPieChart>}
	lib.CreateFunction(tab, "pt_pie_chart",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "labels"},
			{Type: lua.RAW_TABLE, Name: "data"},
			{Type: lua.FLOAT, Name: "x"},
			{Type: lua.FLOAT, Name: "y"},
			{Type: lua.FLOAT, Name: "radius"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			labels := args["labels"].(golua.LValue)
			data := args["data"].(golua.LValue)
			x := args["x"].(float64)
			y := args["y"].(float64)
			radius := args["radius"].(float64)
			t := plotPieTable(state, labels, data, x, y, radius)

			state.Push(t)
			return 1
		})

	/// @func pt_scatter(title, data) -> struct<gui.PlotScatter>
	/// @arg title {string}
	/// @arg data {[]float}
	/// @returns {struct<gui.PlotScatter>}
	lib.CreateFunction(tab, "pt_scatter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_scatter_xy(title, xdata, ydata) -> struct<gui.PlotScatterXY>
	/// @arg title {string}
	/// @arg xdata {[]float}
	/// @arg ydata {[]float}
	/// @returns {struct<gui.PlotScatterXY>}
	lib.CreateFunction(tab, "pt_scatter_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.RAW_TABLE, Name: "xdata"},
			{Type: lua.RAW_TABLE, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_custom(builder) -> struct<gui.PlotCustom>
	/// @arg builder {function()}
	/// @returns {struct<gui.PlotCustom>}
	lib.CreateFunction(tab, "pt_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "builder"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotCustomTable(state, args["builder"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func cursor_screen_pos_xy() -> int, int
	/// @returns {int} - Cursor x position.
	/// @returns {int} - Cursor y position.
	lib.CreateFunction(tab, "cursor_screen_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorScreenPos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func cursor_screen_pos() -> struct<image.Point>
	/// @returns {struct<image.Point>}
	lib.CreateFunction(tab, "cursor_screen_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorScreenPos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func cursor_pos_xy() -> int, int
	/// @returns {int} - Cursor x position.
	/// @returns {int} - Cursor y position.
	lib.CreateFunction(tab, "cursor_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorPos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func cursor_pos() -> struct<image.Point>
	/// @returns {struct<image.Point>}
	lib.CreateFunction(tab, "cursor_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorPos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func canvas_bezier_cubic(pos0, cp0, cp1, pos1, color, thickness, segments)
	/// @arg pos0 {struct<image.Point>}
	/// @arg cp0 {struct<image.Point>}
	/// @arg cp1 {struct<image.Point>}
	/// @arg pos1 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg thickness {float}
	/// @arg segments {int}
	lib.CreateFunction(tab, "canvas_bezier_cubic",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "pos0"},
			{Type: lua.RAW_TABLE, Name: "cp0"},
			{Type: lua.RAW_TABLE, Name: "cp1"},
			{Type: lua.RAW_TABLE, Name: "pos1"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			pos0 := imageutil.TableToPoint(args["pos0"].(*golua.LTable))
			cp0 := imageutil.TableToPoint(args["cp0"].(*golua.LTable))
			cp1 := imageutil.TableToPoint(args["cp1"].(*golua.LTable))
			pos1 := imageutil.TableToPoint(args["pos1"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)
			segments := args["segments"].(int)

			c.AddBezierCubic(pos0, cp0, cp1, pos1, col, float32(thickness), int32(segments))
			return 0
		})

	/// @func canvas_circle(center, radius, color, segments, thickness)
	/// @arg center {struct<image.Point>}
	/// @arg radius {float}
	/// @arg color {struct<image.Color>}
	/// @arg segments {int}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_circle",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.INT, Name: "segments"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			segments := args["segments"].(int)
			thickness := args["thickness"].(float64)

			c.AddCircle(center, float32(radius), col, int32(segments), float32(thickness))
			return 0
		})

	/// @func canvas_circle_filled(center, radius, color)
	/// @arg center {struct<image.Point>}
	/// @arg radius {float}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "canvas_circle_filled",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))

			c.AddCircleFilled(center, float32(radius), col)
			return 0
		})

	/// @func canvas_line(p1, b2, color, thickness)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_line",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddLine(p1, p2, col, float32(thickness))
			return 0
		})

	/// @func canvas_quad(p1, p2, p3, p4, color, thickness)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg p3 {struct<image.Point>}
	/// @arg p4 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_quad",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "p3"},
			{Type: lua.RAW_TABLE, Name: "p4"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(args["p3"].(*golua.LTable))
			p4 := imageutil.TableToPoint(args["p4"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddQuad(p1, p2, p3, p4, col, float32(thickness))
			return 0
		})

	/// @func canvas_quad_filled(p1, p2, p3, p4, color)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg p3 {struct<image.Point>}
	/// @arg p4 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "canvas_quad_filled",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "p3"},
			{Type: lua.RAW_TABLE, Name: "p4"},
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(args["p3"].(*golua.LTable))
			p4 := imageutil.TableToPoint(args["p4"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))

			c.AddQuadFilled(p1, p2, p3, p4, col)
			return 0
		})

	/// @func canvas_rect(min, max, color, rounding, flags, thickness)
	/// @arg min {struct<image.Point>}
	/// @arg max {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg rounding {float}
	/// @arg flags {int<gui.DrawFlags>}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_rect",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "min"},
			{Type: lua.RAW_TABLE, Name: "max"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "rounding"},
			{Type: lua.INT, Name: "flags"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			min := imageutil.TableToPoint(args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(args["max"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			rounding := args["rounding"].(float64)
			flags := args["flags"].(int)
			thickness := args["thickness"].(float64)

			c.AddRect(min, max, col, float32(rounding), g.DrawFlags(flags), float32(thickness))
			return 0
		})

	/// @func canvas_rect_filled(min, max, color, rounding, flags)
	/// @arg min {struct<image.Point>}
	/// @arg max {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg rounding {float}
	/// @arg flags {int<gui.DrawFlags>}
	lib.CreateFunction(tab, "canvas_rect_filled",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "min"},
			{Type: lua.RAW_TABLE, Name: "max"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "rounding"},
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			min := imageutil.TableToPoint(args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(args["max"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			rounding := args["rounding"].(float64)
			flags := args["flags"].(int)

			c.AddRectFilled(min, max, col, float32(rounding), g.DrawFlags(flags))
			return 0
		})

	/// @func canvas_text(pos, color, text)
	/// @arg pos {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg text {string}
	lib.CreateFunction(tab, "canvas_text",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "pos"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			pos := imageutil.TableToPoint(args["pos"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			text := args["text"].(string)

			c.AddText(pos, col, text)
			return 0
		})

	/// @func canvas_triangle(p1, p2, p3, color, thickness)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg p3 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_triangle",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "p3"},
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(args["p3"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddTriangle(p1, p2, p3, col, float32(thickness))
			return 0
		})

	/// @func canvas_triangle_filled(p1, p2, p3, color)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg p3 {struct<image.Point>}
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "canvas_triangle_filled",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "p3"},
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(args["p3"].(*golua.LTable))
			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))

			c.AddTriangleFilled(p1, p2, p3, col)
			return 0
		})

	/// @func canvas_path_arc_to(center, radius, min, max, segments)
	/// @arg center {struct<image.Point>}
	/// @arg radius {float}
	/// @arg min {float}
	/// @arg max {float}
	/// @arg segments {int}
	lib.CreateFunction(tab, "canvas_path_arc_to",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.FLOAT, Name: "min"},
			{Type: lua.FLOAT, Name: "max"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			min := args["min"].(float64)
			max := args["max"].(float64)
			segments := args["segments"].(int)

			c.PathArcTo(center, float32(radius), float32(min), float32(max), int32(segments))
			return 0
		})

	/// @func canvas_path_arc_to_fast(center, radius, min, max, segments)
	/// @arg center {struct<image.Point>}
	/// @arg radius {float}
	/// @arg min {int}
	/// @arg max {int}
	lib.CreateFunction(tab, "canvas_path_arc_to_fast",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.INT, Name: "min"},
			{Type: lua.INT, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			min := args["min"].(int)
			max := args["max"].(int)

			c.PathArcToFast(center, float32(radius), int32(min), int32(max))
			return 0
		})

	/// @func canvas_path_bezier_cubic_to(p1, p2, p3, segments)
	/// @arg p1 {struct<image.Point>}
	/// @arg p2 {struct<image.Point>}
	/// @arg p3 {struct<image.Point>}
	/// @arg segments {int}
	lib.CreateFunction(tab, "canvas_path_bezier_cubic_to",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
			{Type: lua.RAW_TABLE, Name: "p2"},
			{Type: lua.RAW_TABLE, Name: "p3"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(args["p3"].(*golua.LTable))
			segments := args["segments"].(int)

			c.PathBezierCubicCurveTo(p1, p2, p3, int32(segments))
			return 0
		})

	/// @func canvas_path_clear()
	lib.CreateFunction(tab, "canvas_path_clear",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()
			c.PathClear()
			return 0
		})

	/// @func canvas_fill_convex(color)
	/// @arg color {struct<image.Color>}
	lib.CreateFunction(tab, "canvas_path_fill_convex",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))

			c.PathFillConvex(col)
			return 0
		})

	/// @func canvas_path_line_to(p1, segments)
	/// @arg p1 {struct<image.Point>}
	lib.CreateFunction(tab, "canvas_path_line_to",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))

			c.PathLineTo(p1)
			return 0
		})

	/// @func canvas_path_line_to_merge_duplicate(p1)
	/// @arg p1 {struct<image.Point>}
	lib.CreateFunction(tab, "canvas_path_line_to_merge_duplicate",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "p1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(args["p1"].(*golua.LTable))

			c.PathLineToMergeDuplicate(p1)
			return 0
		})

	/// @func canvas_path_stroke(color, flags, thickness)
	/// @arg color {struct<image.Color>}
	/// @arg flags {int<gui.DrawFlags>}
	/// @arg thickness {float}
	lib.CreateFunction(tab, "canvas_path_stroke",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
			{Type: lua.INT, Name: "flags"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			col := imageutil.ColorTableToRGBAColor(args["color"].(*golua.LTable))
			flags := args["flags"].(int)
			thickness := args["thickness"].(float64)

			c.PathStroke(col, g.DrawFlags(flags), float32(thickness))
			return 0
		})

	/// @func fontatlas_add_font(name, size) -> int<ref.FONT>, bool
	/// @arg name {string}
	/// @arg size {float}
	/// @returns {int<ref.FONT>}
	/// @returns {bool}
	/// @desc
	/// The returned font ref will be nil if ok is false.
	lib.CreateFunction(tab, "fontatlas_add_font",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.FLOAT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			fi := atlas.AddFont(args["name"].(string), float32(args["size"].(float64)))
			if fi == nil {
				lg.Append(fmt.Sprintf("failed to add font: %s", args["name"]), log.LEVEL_WARN)
				state.Push(golua.LNil)
				state.Push(golua.LFalse)
				return 2
			}

			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: fi,
			})

			state.Push(golua.LNumber(ref))
			state.Push(golua.LTrue)
			return 2
		})

	/// @func fontatlas_default_font_strings() -> []string
	/// @returns {[]string}
	lib.CreateFunction(tab, "fontatlas_default_font_strings",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			fonts := atlas.GetDefaultFonts()
			t := state.NewTable()

			for _, f := range fonts {
				t.Append(golua.LString(f.String()))
			}

			state.Push(t)
			return 1
		})

	/// @func fontatlas_default_fonts() -> []int<ref.FONT>
	/// @returns {[]int<ref.FONT>}
	/// @desc
	/// Take note that this creates an array of refs,
	/// refs are only cleared when the workflow ends,
	/// or manually with 'ref.del' or 'ref.del_many'.
	lib.CreateFunction(tab, "fontatlas_default_fonts",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			fonts := atlas.GetDefaultFonts()
			t := state.NewTable()

			for _, f := range fonts {
				ref := r.CR_REF.Add(&collection.RefItem[any]{
					Value: &f,
				})
				t.Append(golua.LNumber(ref))
			}

			state.Push(t)
			return 1
		})

	/// @func fontatlas_register_string(str) -> string
	/// @arg str {string}
	/// @returns {string}
	lib.CreateFunction(tab, "fontatlas_register_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			str := atlas.RegisterString(args["str"].(string))

			state.Push(golua.LString(str))
			return 1
		})

	/// @func fontatlas_register_string_ref(stringref) -> int<ref.STRING>
	/// @arg stringref {int<ref.STRING>}
	/// @returns {int<ref.STRING>}
	lib.CreateFunction(tab, "fontatlas_register_string_ref",
		[]lua.Arg{
			{Type: lua.INT, Name: "stringref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas

			ref := args["stringref"].(int)
			sref, err := r.CR_REF.Item(ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			stringref := sref.Value.(*string)
			atlas.RegisterStringPointer(stringref)

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func fontatlas_register_string_many(str) -> []string
	/// @arg str {[]string}
	/// @returns {[]string}
	lib.CreateFunction(tab, "fontatlas_register_string_many",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			strList := args["str"].(*golua.LTable)
			strSlice := []string{}

			for i := range strList.Len() {
				v := strList.RawGetInt(i + 1).(golua.LString)
				strSlice = append(strSlice, string(v))
			}

			atlas.RegisterStringSlice(strSlice)

			state.Push(strList)
			return 1
		})

	/// @func fontatlas_set_default_font(name, size)
	/// @arg name {string}
	/// @arg size {float}
	lib.CreateFunction(tab, "fontatlas_set_default_font",
		[]lua.Arg{
			{Type: lua.STRING, Name: "name"},
			{Type: lua.FLOAT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			atlas.SetDefaultFont(args["name"].(string), float32(args["size"].(float64)))
			return 0
		})

	/// @func fontatlas_set_default_font_size(size)
	/// @arg size {float}
	lib.CreateFunction(tab, "fontatlas_set_default_font_size",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			atlas.SetDefaultFontSize(float32(args["size"].(float64)))
			return 0
		})

	/// @func font_set_size(fontref, size) -> int<ref.FONT>
	/// @arg fontref {int<ref.FONT>}
	/// @arg size {float}
	/// @returns {int<ref.FLOAT>}
	lib.CreateFunction(tab, "font_set_size",
		[]lua.Arg{
			{Type: lua.INT, Name: "fontref"},
			{Type: lua.FLOAT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ref := args["fontref"].(int)
			sref, err := r.CR_REF.Item(ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			font := sref.Value.(*g.FontInfo)

			font.SetSize(float32(args["size"].(float64)))

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func font_string(fontref) -> string
	/// @arg fontref {int<ref.FONT>}
	/// @returns {string}
	lib.CreateFunction(tab, "font_string",
		[]lua.Arg{
			{Type: lua.INT, Name: "fontref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			ref := args["fontref"].(int)
			sref, err := r.CR_REF.Item(ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			font := sref.Value.(*g.FontInfo)

			state.Push(golua.LString(font.String()))
			return 1
		})

	/// @constants ColorEditFlags {int}
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
	tab.RawSetString("FLAGCOLOREDIT_NONE", golua.LNumber(FLAGCOLOREDIT_NONE))
	tab.RawSetString("FLAGCOLOREDIT_NOALPHA", golua.LNumber(FLAGCOLOREDIT_NOALPHA))
	tab.RawSetString("FLAGCOLOREDIT_NOPICKER", golua.LNumber(FLAGCOLOREDIT_NOPICKER))
	tab.RawSetString("FLAGCOLOREDIT_NOOPTIONS", golua.LNumber(FLAGCOLOREDIT_NOOPTIONS))
	tab.RawSetString("FLAGCOLOREDIT_NOSMALLPREVIEW", golua.LNumber(FLAGCOLOREDIT_NOSMALLPREVIEW))
	tab.RawSetString("FLAGCOLOREDIT_NOINPUTS", golua.LNumber(FLAGCOLOREDIT_NOINPUTS))
	tab.RawSetString("FLAGCOLOREDIT_NOTOOLTIP", golua.LNumber(FLAGCOLOREDIT_NOTOOLTIP))
	tab.RawSetString("FLAGCOLOREDIT_NOLABEL", golua.LNumber(FLAGCOLOREDIT_NOLABEL))
	tab.RawSetString("FLAGCOLOREDIT_NOSIDEPREVIEW", golua.LNumber(FLAGCOLOREDIT_NOSIDEPREVIEW))
	tab.RawSetString("FLAGCOLOREDIT_NODRAGDROP", golua.LNumber(FLAGCOLOREDIT_NODRAGDROP))
	tab.RawSetString("FLAGCOLOREDIT_NOBORDER", golua.LNumber(FLAGCOLOREDIT_NOBORDER))
	tab.RawSetString("FLAGCOLOREDIT_ALPHABAR", golua.LNumber(FLAGCOLOREDIT_ALPHABAR))
	tab.RawSetString("FLAGCOLOREDIT_ALPHAPREVIEW", golua.LNumber(FLAGCOLOREDIT_ALPHAPREVIEW))
	tab.RawSetString("FLAGCOLOREDIT_ALPHAPREVIEWHALF", golua.LNumber(FLAGCOLOREDIT_ALPHAPREVIEWHALF))
	tab.RawSetString("FLAGCOLOREDIT_HDR", golua.LNumber(FLAGCOLOREDIT_HDR))
	tab.RawSetString("FLAGCOLOREDIT_DISPLAYRGB", golua.LNumber(FLAGCOLOREDIT_DISPLAYRGB))
	tab.RawSetString("FLAGCOLOREDIT_DISPLAYHSV", golua.LNumber(FLAGCOLOREDIT_DISPLAYHSV))
	tab.RawSetString("FLAGCOLOREDIT_DISPLAYHEX", golua.LNumber(FLAGCOLOREDIT_DISPLAYHEX))
	tab.RawSetString("FLAGCOLOREDIT_UINT8", golua.LNumber(FLAGCOLOREDIT_UINT8))
	tab.RawSetString("FLAGCOLOREDIT_FLOAT", golua.LNumber(FLAGCOLOREDIT_FLOAT))
	tab.RawSetString("FLAGCOLOREDIT_HUEBAR", golua.LNumber(FLAGCOLOREDIT_HUEBAR))
	tab.RawSetString("FLAGCOLOREDIT_HUEWHEEL", golua.LNumber(FLAGCOLOREDIT_HUEWHEEL))
	tab.RawSetString("FLAGCOLOREDIT_INPUTRGB", golua.LNumber(FLAGCOLOREDIT_INPUTRGB))
	tab.RawSetString("FLAGCOLOREDIT_INPUTHSV", golua.LNumber(FLAGCOLOREDIT_INPUTHSV))
	tab.RawSetString("FLAGCOLOREDIT_DEFAULTOPTIONS", golua.LNumber(FLAGCOLOREDIT_DEFAULTOPTIONS))
	tab.RawSetString("FLAGCOLOREDIT_DISPLAYMASK", golua.LNumber(FLAGCOLOREDIT_DISPLAYMASK))
	tab.RawSetString("FLAGCOLOREDIT_DATATYPEMASK", golua.LNumber(FLAGCOLOREDIT_DATATYPEMASK))
	tab.RawSetString("FLAGCOLOREDIT_PICKERMASK", golua.LNumber(FLAGCOLOREDIT_PICKERMASK))
	tab.RawSetString("FLAGCOLOREDIT_INPUTMASK", golua.LNumber(FLAGCOLOREDIT_INPUTMASK))

	/// @constants ComboFlags {int}
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
	tab.RawSetString("FLAGCOMBO_NONE", golua.LNumber(FLAGCOMBO_NONE))
	tab.RawSetString("FLAGCOMBO_POPUPALIGNLEFT", golua.LNumber(FLAGCOMBO_POPUPALIGNLEFT))
	tab.RawSetString("FLAGCOMBO_HEIGHTSMALL", golua.LNumber(FLAGCOMBO_HEIGHTSMALL))
	tab.RawSetString("FLAGCOMBO_HEIGHTREGULAR", golua.LNumber(FLAGCOMBO_HEIGHTREGULAR))
	tab.RawSetString("FLAGCOMBO_HEIGHTLARGEST", golua.LNumber(FLAGCOMBO_HEIGHTLARGEST))
	tab.RawSetString("FLAGCOMBO_NOARROWBUTTON", golua.LNumber(FLAGCOMBO_NOARROWBUTTON))
	tab.RawSetString("FLAGCOMBO_NOARROWBUTTON", golua.LNumber(FLAGCOMBO_NOARROWBUTTON))
	tab.RawSetString("FLAGCOMBO_NOPREVIEW", golua.LNumber(FLAGCOMBO_NOPREVIEW))
	tab.RawSetString("FLAGCOMBO_WIDTHFITPREVIEW", golua.LNumber(FLAGCOMBO_WIDTHFITPREVIEW))
	tab.RawSetString("FLAGCOMBO_HEIGHTMASK", golua.LNumber(FLAGCOMBO_HEIGHTMASK))

	/// @constants MouseButton {int}
	/// @const MOUSEBUTTON_LEFT
	/// @const MOUSEBUTTON_RIGHT
	/// @const MOUSEBUTTON_MIDDLE
	tab.RawSetString("MOUSEBUTTON_LEFT", golua.LNumber(MOUSEBUTTON_LEFT))
	tab.RawSetString("MOUSEBUTTON_RIGHT", golua.LNumber(MOUSEBUTTON_RIGHT))
	tab.RawSetString("MOUSEBUTTON_MIDDLE", golua.LNumber(MOUSEBUTTON_MIDDLE))

	/// @constants DatePickerLabel {string}
	/// @const DATEPICKERLABEL_MONTH
	/// @const DATEPICKERLABEL_YEAR
	tab.RawSetString("DATEPICKERLABEL_MONTH", golua.LString(DATEPICKERLABEL_MONTH))
	tab.RawSetString("DATEPICKERLABEL_YEAR", golua.LString(DATEPICKERLABEL_YEAR))

	/// @constants InputFlag {int}
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
	tab.RawSetString("FLAGINPUTTEXT_NONE", golua.LNumber(FLAGINPUTTEXT_NONE))
	tab.RawSetString("FLAGINPUTTEXT_CHARSDECIMAL", golua.LNumber(FLAGINPUTTEXT_CHARSDECIMAL))
	tab.RawSetString("FLAGINPUTTEXT_CHARSHEXADECIMAL", golua.LNumber(FLAGINPUTTEXT_CHARSHEXADECIMAL))
	tab.RawSetString("FLAGINPUTTEXT_CHARSUPPERCASE", golua.LNumber(FLAGINPUTTEXT_CHARSUPPERCASE))
	tab.RawSetString("FLAGINPUTTEXT_CHARSNOBLANK", golua.LNumber(FLAGINPUTTEXT_CHARSNOBLANK))
	tab.RawSetString("FLAGINPUTTEXT_AUTOSELECTALL", golua.LNumber(FLAGINPUTTEXT_AUTOSELECTALL))
	tab.RawSetString("FLAGINPUTTEXT_ENTERRETURNSTRUE", golua.LNumber(FLAGINPUTTEXT_ENTERRETURNSTRUE))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKCOMPLETION", golua.LNumber(FLAGINPUTTEXT_CALLBACKCOMPLETION))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKHISTORY", golua.LNumber(FLAGINPUTTEXT_CALLBACKHISTORY))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKALWAYS", golua.LNumber(FLAGINPUTTEXT_CALLBACKALWAYS))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKCHARFILTER", golua.LNumber(FLAGINPUTTEXT_CALLBACKCHARFILTER))
	tab.RawSetString("FLAGINPUTTEXT_ALLOWTABINPUT", golua.LNumber(FLAGINPUTTEXT_ALLOWTABINPUT))
	tab.RawSetString("FLAGINPUTTEXT_CTRLENTERFORNEWLINE", golua.LNumber(FLAGINPUTTEXT_CTRLENTERFORNEWLINE))
	tab.RawSetString("FLAGINPUTTEXT_NOHORIZONTALSCROLL", golua.LNumber(FLAGINPUTTEXT_NOHORIZONTALSCROLL))
	tab.RawSetString("FLAGINPUTTEXT_ALWAYSOVERWRITE", golua.LNumber(FLAGINPUTTEXT_ALWAYSOVERWRITE))
	tab.RawSetString("FLAGINPUTTEXT_READONLY", golua.LNumber(FLAGINPUTTEXT_READONLY))
	tab.RawSetString("FLAGINPUTTEXT_PASSWORD", golua.LNumber(FLAGINPUTTEXT_PASSWORD))
	tab.RawSetString("FLAGINPUTTEXT_NOUNDOREDO", golua.LNumber(FLAGINPUTTEXT_NOUNDOREDO))
	tab.RawSetString("FLAGINPUTTEXT_CHARSSCIENTIFIC", golua.LNumber(FLAGINPUTTEXT_CHARSSCIENTIFIC))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKRESIZE", golua.LNumber(FLAGINPUTTEXT_CALLBACKRESIZE))
	tab.RawSetString("FLAGINPUTTEXT_CALLBACKEDIT", golua.LNumber(FLAGINPUTTEXT_CALLBACKEDIT))
	tab.RawSetString("FLAGINPUTTEXT_ESCAPECLEARSALL", golua.LNumber(FLAGINPUTTEXT_ESCAPECLEARSALL))

	/// @constants SelectableFlags {int}
	/// @const FLAGSELECTABLE_NONE
	/// @const FLAGSELECTABLE_DONTCLOSEPOPUPS
	/// @const FLAGSELECTABLE_SPANALLCOLUMNS
	/// @const FLAGSELECTABLE_ALLOWDOUBLECLICK
	/// @const FLAGSELECTABLE_DISABLED
	/// @const FLAGSELECTABLE_ALLOWOVERLAP
	tab.RawSetString("FLAGSELECTABLE_NONE", golua.LNumber(FLAGSELECTABLE_NONE))
	tab.RawSetString("FLAGSELECTABLE_DONTCLOSEPOPUPS", golua.LNumber(FLAGSELECTABLE_DONTCLOSEPOPUPS))
	tab.RawSetString("FLAGSELECTABLE_SPANALLCOLUMNS", golua.LNumber(FLAGSELECTABLE_SPANALLCOLUMNS))
	tab.RawSetString("FLAGSELECTABLE_ALLOWDOUBLECLICK", golua.LNumber(FLAGSELECTABLE_ALLOWDOUBLECLICK))
	tab.RawSetString("FLAGSELECTABLE_DISABLED", golua.LNumber(FLAGSELECTABLE_DISABLED))
	tab.RawSetString("FLAGSELECTABLE_ALLOWOVERLAP", golua.LNumber(FLAGSELECTABLE_ALLOWOVERLAP))

	/// @constants SliderFlags {int}
	/// @const FLAGSLIDER_NONE
	/// @const FLAGSLIDER_ALWAYSCLAMP
	/// @const FLAGSLIDER_LOGARITHMIC
	/// @const FLAGSLIDER_NOROUNDTOFORMAT
	/// @const FLAGSLIDER_NOINPUT
	/// @const FLAGSLIDER_INVALIDMASK
	tab.RawSetString("FLAGSLIDER_NONE", golua.LNumber(FLAGSLIDER_NONE))
	tab.RawSetString("FLAGSLIDER_ALWAYSCLAMP", golua.LNumber(FLAGSLIDER_ALWAYSCLAMP))
	tab.RawSetString("FLAGSLIDER_LOGARITHMIC", golua.LNumber(FLAGSLIDER_LOGARITHMIC))
	tab.RawSetString("FLAGSLIDER_NOROUNDTOFORMAT", golua.LNumber(FLAGSLIDER_NOROUNDTOFORMAT))
	tab.RawSetString("FLAGSLIDER_NOINPUT", golua.LNumber(FLAGSLIDER_NOINPUT))
	tab.RawSetString("FLAGSLIDER_INVALIDMASK", golua.LNumber(FLAGSLIDER_INVALIDMASK))

	/// @constants TabBarFlags {int}
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
	tab.RawSetString("FLAGTABBAR_NONE", golua.LNumber(FLAGTABBAR_NONE))
	tab.RawSetString("FLAGTABBAR_REORDERABLE", golua.LNumber(FLAGTABBAR_REORDERABLE))
	tab.RawSetString("FLAGTABBAR_AUTOSELECTNEWTABS", golua.LNumber(FLAGTABBAR_AUTOSELECTNEWTABS))
	tab.RawSetString("FLAGTABBAR_TABLLISTPOPUPBUTTON", golua.LNumber(FLAGTABBAR_TABLLISTPOPUPBUTTON))
	tab.RawSetString("FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON", golua.LNumber(FLAGTABBAR_NOCLOSEWITHMIDDLEMOUSEBUTTON))
	tab.RawSetString("FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS", golua.LNumber(FLAGTABBAR_NOTABLISTSCROLLINGBUTTONS))
	tab.RawSetString("FLAGTABBAR_NOTOOLTIP", golua.LNumber(FLAGTABBAR_NOTOOLTIP))
	tab.RawSetString("FLAGTABBAR_FITTINGPOLICYRESIZEDOWN", golua.LNumber(FLAGTABBAR_FITTINGPOLICYRESIZEDOWN))
	tab.RawSetString("FLAGTABBAR_FITTINGPOLICYSCROLL", golua.LNumber(FLAGTABBAR_FITTINGPOLICYSCROLL))
	tab.RawSetString("FLAGTABBAR_FITTINGPOLICYMASK", golua.LNumber(FLAGTABBAR_FITTINGPOLICYMASK))
	tab.RawSetString("FLAGTABBAR_FITTINGPOLICYDEFAULT", golua.LNumber(FLAGTABBAR_FITTINGPOLICYDEFAULT))

	/// @constants TabItemFlags {int}
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
	tab.RawSetString("FLAGTABITEM_NONE", golua.LNumber(FLAGTABITEM_NONE))
	tab.RawSetString("FLAGTABITEM_UNSAVEDOCUMENT", golua.LNumber(FLAGTABITEM_UNSAVEDOCUMENT))
	tab.RawSetString("FLAGTABITEM_SETSELECTED", golua.LNumber(FLAGTABITEM_SETSELECTED))
	tab.RawSetString("FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON", golua.LNumber(FLAGTABITEM_NOCLOSEWITHMIDDLEMOUSEBUTTON))
	tab.RawSetString("FLAGTABITEM_NOPUSHID", golua.LNumber(FLAGTABITEM_NOPUSHID))
	tab.RawSetString("FLAGTABITEM_NOTOOLTIP", golua.LNumber(FLAGTABITEM_NOTOOLTIP))
	tab.RawSetString("FLAGTABITEM_NOREORDER", golua.LNumber(FLAGTABITEM_NOREORDER))
	tab.RawSetString("FLAGTABITEM_LEADING", golua.LNumber(FLAGTABITEM_LEADING))
	tab.RawSetString("FLAGTABITEM_TRAILING", golua.LNumber(FLAGTABITEM_TRAILING))
	tab.RawSetString("FLAGTABITEM_NOASSUMEDCLOSURE", golua.LNumber(FLAGTABITEM_NOASSUMEDCLOSURE))

	/// @constants TableColumnFlags {int}
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
	tab.RawSetString("FLAGTABLECOLUMN_NONE", golua.LNumber(FLAGTABLECOLUMN_NONE))
	tab.RawSetString("FLAGTABLECOLUMN_DEFAULTHIDE", golua.LNumber(FLAGTABLECOLUMN_DEFAULTHIDE))
	tab.RawSetString("FLAGTABLECOLUMN_DEFAULTSORT", golua.LNumber(FLAGTABLECOLUMN_DEFAULTSORT))
	tab.RawSetString("FLAGTABLECOLUMN_WIDTHSTRETCH", golua.LNumber(FLAGTABLECOLUMN_WIDTHSTRETCH))
	tab.RawSetString("FLAGTABLECOLUMN_WIDTHFIXED", golua.LNumber(FLAGTABLECOLUMN_WIDTHFIXED))
	tab.RawSetString("FLAGTABLECOLUMN_NORESIZE", golua.LNumber(FLAGTABLECOLUMN_NORESIZE))
	tab.RawSetString("FLAGTABLECOLUMN_NOREORDER", golua.LNumber(FLAGTABLECOLUMN_NOREORDER))
	tab.RawSetString("FLAGTABLECOLUMN_NOHIDE", golua.LNumber(FLAGTABLECOLUMN_NOHIDE))
	tab.RawSetString("FLAGTABLECOLUMN_NOCLIP", golua.LNumber(FLAGTABLECOLUMN_NOCLIP))
	tab.RawSetString("FLAGTABLECOLUMN_NOSORT", golua.LNumber(FLAGTABLECOLUMN_NOSORT))
	tab.RawSetString("FLAGTABLECOLUMN_NOSORTASCENDING", golua.LNumber(FLAGTABLECOLUMN_NOSORTASCENDING))
	tab.RawSetString("FLAGTABLECOLUMN_NOSORTDESCENDING", golua.LNumber(FLAGTABLECOLUMN_NOSORTDESCENDING))
	tab.RawSetString("FLAGTABLECOLUMN_NOHEADERWIDTH", golua.LNumber(FLAGTABLECOLUMN_NOHEADERWIDTH))
	tab.RawSetString("FLAGTABLECOLUMN_PREFERSORTASCENDING", golua.LNumber(FLAGTABLECOLUMN_PREFERSORTASCENDING))
	tab.RawSetString("FLAGTABLECOLUMN_PREFERSORTDESCENDING", golua.LNumber(FLAGTABLECOLUMN_PREFERSORTDESCENDING))
	tab.RawSetString("FLAGTABLECOLUMN_INDENTENABLE", golua.LNumber(FLAGTABLECOLUMN_INDENTENABLE))
	tab.RawSetString("FLAGTABLECOLUMN_INDENTDISABLE", golua.LNumber(FLAGTABLECOLUMN_INDENTDISABLE))
	tab.RawSetString("FLAGTABLECOLUMN_ISENABLED", golua.LNumber(FLAGTABLECOLUMN_ISENABLED))
	tab.RawSetString("FLAGTABLECOLUMN_ISVISIBLE", golua.LNumber(FLAGTABLECOLUMN_ISVISIBLE))
	tab.RawSetString("FLAGTABLECOLUMN_ISSORTED", golua.LNumber(FLAGTABLECOLUMN_ISSORTED))
	tab.RawSetString("FLAGTABLECOLUMN_ISHOVERED", golua.LNumber(FLAGTABLECOLUMN_ISHOVERED))
	tab.RawSetString("FLAGTABLECOLUMN_WIDTHMASK", golua.LNumber(FLAGTABLECOLUMN_WIDTHMASK))
	tab.RawSetString("FLAGTABLECOLUMN_INDENTMASK", golua.LNumber(FLAGTABLECOLUMN_INDENTMASK))
	tab.RawSetString("FLAGTABLECOLUMN_STATUSMASK", golua.LNumber(FLAGTABLECOLUMN_STATUSMASK))
	tab.RawSetString("FLAGTABLECOLUMN_NODIRECTRESIZE", golua.LNumber(FLAGTABLECOLUMN_NODIRECTRESIZE))

	/// @constants TableRowFlags {int}
	/// @const FLAGTABLEROW_NONE
	/// @const FLAGTABLEROW_HEADERS
	tab.RawSetString("FLAGTABLEROW_NONE", golua.LNumber(FLAGTABLEROW_NONE))
	tab.RawSetString("FLAGTABLEROW_HEADERS", golua.LNumber(FLAGTABLEROW_HEADERS))

	/// @constants TableFlags {int}
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
	tab.RawSetString("FLAGTABLE_NONE", golua.LNumber(FLAGTABLE_NONE))
	tab.RawSetString("FLAGTABLE_RESIZEABLE", golua.LNumber(FLAGTABLE_RESIZEABLE))
	tab.RawSetString("FLAGTABLE_REORDERABLE", golua.LNumber(FLAGTABLE_REORDERABLE))
	tab.RawSetString("FLAGTABLE_HIDEABLE", golua.LNumber(FLAGTABLE_HIDEABLE))
	tab.RawSetString("FLAGTABLE_SORTABLE", golua.LNumber(FLAGTABLE_SORTABLE))
	tab.RawSetString("FLAGTABLE_NOSAVEDSETTINGS", golua.LNumber(FLAGTABLE_NOSAVEDSETTINGS))
	tab.RawSetString("FLAGTABLE_CONTEXTMENUINBODY", golua.LNumber(FLAGTABLE_CONTEXTMENUINBODY))
	tab.RawSetString("FLAGTABLE_ROWBG", golua.LNumber(FLAGTABLE_ROWBG))
	tab.RawSetString("FLAGTABLE_BORDERSINNERH", golua.LNumber(FLAGTABLE_BORDERSINNERH))
	tab.RawSetString("FLAGTABLE_BORDERSOUTERH", golua.LNumber(FLAGTABLE_BORDERSOUTERH))
	tab.RawSetString("FLAGTABLE_BORDERSINNERV", golua.LNumber(FLAGTABLE_BORDERSINNERV))
	tab.RawSetString("FLAGTABLE_BORDERSOUTERV", golua.LNumber(FLAGTABLE_BORDERSOUTERV))
	tab.RawSetString("FLAGTABLE_BORDERSH", golua.LNumber(FLAGTABLE_BORDERSH))
	tab.RawSetString("FLAGTABLE_BORDERSV", golua.LNumber(FLAGTABLE_BORDERSV))
	tab.RawSetString("FLAGTABLE_BORDERSINNER", golua.LNumber(FLAGTABLE_BORDERSINNER))
	tab.RawSetString("FLAGTABLE_BORDERSOUTER", golua.LNumber(FLAGTABLE_BORDERSOUTER))
	tab.RawSetString("FLAGTABLE_BORDERS", golua.LNumber(FLAGTABLE_BORDERS))
	tab.RawSetString("FLAGTABLE_NOBORDERSINBODY", golua.LNumber(FLAGTABLE_NOBORDERSINBODY))
	tab.RawSetString("FLAGTABLE_NOBORDERSINBODYUNTILRESIZE", golua.LNumber(FLAGTABLE_NOBORDERSINBODYUNTILRESIZE))
	tab.RawSetString("FLAGTABLE_SIZINGFIXEDFIT", golua.LNumber(FLAGTABLE_SIZINGFIXEDFIT))
	tab.RawSetString("FLAGTABLE_SIZINGFIXEDSAME", golua.LNumber(FLAGTABLE_SIZINGFIXEDSAME))
	tab.RawSetString("FLAGTABLE_SIZINGSTRETCHPROP", golua.LNumber(FLAGTABLE_SIZINGSTRETCHPROP))
	tab.RawSetString("FLAGTABLE_SIZINGSTRETCHSAME", golua.LNumber(FLAGTABLE_SIZINGSTRETCHSAME))
	tab.RawSetString("FLAGTABLE_NOHOSTEXTENDX", golua.LNumber(FLAGTABLE_NOHOSTEXTENDX))
	tab.RawSetString("FLAGTABLE_NOHOSTEXTENDY", golua.LNumber(FLAGTABLE_NOHOSTEXTENDY))
	tab.RawSetString("FLAGTABLE_NOKEEPCOLUMNSVISIBLE", golua.LNumber(FLAGTABLE_NOKEEPCOLUMNSVISIBLE))
	tab.RawSetString("FLAGTABLE_PRECISEWIDTHS", golua.LNumber(FLAGTABLE_PRECISEWIDTHS))
	tab.RawSetString("FLAGTABLE_NOCLIP", golua.LNumber(FLAGTABLE_NOCLIP))
	tab.RawSetString("FLAGTABLE_PADOUTERX", golua.LNumber(FLAGTABLE_PADOUTERX))
	tab.RawSetString("FLAGTABLE_NOPADOUTERX", golua.LNumber(FLAGTABLE_NOPADOUTERX))
	tab.RawSetString("FLAGTABLE_NOPADINNERX", golua.LNumber(FLAGTABLE_NOPADINNERX))
	tab.RawSetString("FLAGTABLE_SCROLLX", golua.LNumber(FLAGTABLE_SCROLLX))
	tab.RawSetString("FLAGTABLE_SCROLLY", golua.LNumber(FLAGTABLE_SCROLLY))
	tab.RawSetString("FLAGTABLE_SORTMULTI", golua.LNumber(FLAGTABLE_SORTMULTI))
	tab.RawSetString("FLAGTABLE_SORTTRISTATE", golua.LNumber(FLAGTABLE_SORTTRISTATE))
	tab.RawSetString("FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN", golua.LNumber(FLAGTABLE_HIGHLIGHTHOVEREDCOLUMN))
	tab.RawSetString("FLAGTABLE_SIZINGMASK", golua.LNumber(FLAGTABLE_SIZINGMASK))

	/// @constants Direction {int}
	/// @const DIR_NONE
	/// @const DIR_LEFT
	/// @const DIR_RIGHT
	/// @const DIR_UP
	/// @const DIR_DOWN
	/// @const DIR_COUNT
	tab.RawSetString("DIR_NONE", golua.LNumber(DIR_NONE))
	tab.RawSetString("DIR_LEFT", golua.LNumber(DIR_LEFT))
	tab.RawSetString("DIR_RIGHT", golua.LNumber(DIR_RIGHT))
	tab.RawSetString("DIR_UP", golua.LNumber(DIR_UP))
	tab.RawSetString("DIR_DOWN", golua.LNumber(DIR_DOWN))
	tab.RawSetString("DIR_COUNT", golua.LNumber(DIR_COUNT))

	/// @constants TreeNodeFlags {int}
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
	tab.RawSetString("FLAGTREENODE_NONE", golua.LNumber(FLAGTREENODE_NONE))
	tab.RawSetString("FLAGTREENODE_SELECTED", golua.LNumber(FLAGTREENODE_SELECTED))
	tab.RawSetString("FLAGTREENODE_FRAMED", golua.LNumber(FLAGTREENODE_FRAMED))
	tab.RawSetString("FLAGTREENODE_ALLOWOVERLAP", golua.LNumber(FLAGTREENODE_ALLOWOVERLAP))
	tab.RawSetString("FLAGTREENODE_NOTREEPUSHONOPEN", golua.LNumber(FLAGTREENODE_NOTREEPUSHONOPEN))
	tab.RawSetString("FLAGTREENODE_NOAUTOOPENONLOG", golua.LNumber(FLAGTREENODE_NOAUTOOPENONLOG))
	tab.RawSetString("FLAGTREENODE_DEFAULTOPEN", golua.LNumber(FLAGTREENODE_DEFAULTOPEN))
	tab.RawSetString("FLAGTREENODE_OPENONDOUBLECLICK", golua.LNumber(FLAGTREENODE_OPENONDOUBLECLICK))
	tab.RawSetString("FLAGTREENODE_OPENONARROW", golua.LNumber(FLAGTREENODE_OPENONARROW))
	tab.RawSetString("FLAGTREENODE_LEAF", golua.LNumber(FLAGTREENODE_LEAF))
	tab.RawSetString("FLAGTREENODE_BULLET", golua.LNumber(FLAGTREENODE_BULLET))
	tab.RawSetString("FLAGTREENODE_FRAMEPADDING", golua.LNumber(FLAGTREENODE_FRAMEPADDING))
	tab.RawSetString("FLAGTREENODE_SPANAVAILWIDTH", golua.LNumber(FLAGTREENODE_SPANAVAILWIDTH))
	tab.RawSetString("FLAGTREENODE_SPANFULLWIDTH", golua.LNumber(FLAGTREENODE_SPANFULLWIDTH))
	tab.RawSetString("FLAGTREENODE_SPANALLCOLUMNS", golua.LNumber(FLAGTREENODE_SPANALLCOLUMNS))
	tab.RawSetString("FLAGTREENODE_NAVLEFTJUMPSBACKHERE", golua.LNumber(FLAGTREENODE_NAVLEFTJUMPSBACKHERE))
	tab.RawSetString("FLAGTREENODE_COLLAPSINGHEADER", golua.LNumber(FLAGTREENODE_COLLAPSINGHEADER))

	/// @constants MasterWindowFlags {int}
	/// @const FLAGMASTERWINDOW_NOTRESIZABLE
	/// @const FLAGMASTERWINDOW_MAXIMIZED
	/// @const FLAGMASTERWINDOW_FLOATING
	/// @const FLAGMASTERWINDOW_FRAMELESS
	/// @const FLAGMASTERWINDOW_TRANSPARENT
	tab.RawSetString("FLAGMASTERWINDOW_NOTRESIZABLE", golua.LNumber(FLAGMASTERWINDOW_NOTRESIZABLE))
	tab.RawSetString("FLAGMASTERWINDOW_MAXIMIZED", golua.LNumber(FLAGMASTERWINDOW_MAXIMIZED))
	tab.RawSetString("FLAGMASTERWINDOW_FLOATING", golua.LNumber(FLAGMASTERWINDOW_FLOATING))
	tab.RawSetString("FLAGMASTERWINDOW_FRAMELESS", golua.LNumber(FLAGMASTERWINDOW_FRAMELESS))
	tab.RawSetString("FLAGMASTERWINDOW_TRANSPARENT", golua.LNumber(FLAGMASTERWINDOW_TRANSPARENT))

	/// @constants WindowFlags {int}
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
	tab.RawSetString("FLAGWINDOW_NONE", golua.LNumber(FLAGWINDOW_NONE))
	tab.RawSetString("FLAGWINDOW_NOTITLEBAR", golua.LNumber(FLAGWINDOW_NOTITLEBAR))
	tab.RawSetString("FLAGWINDOW_NORESIZE", golua.LNumber(FLAGWINDOW_NORESIZE))
	tab.RawSetString("FLAGWINDOW_NOMOVE", golua.LNumber(FLAGWINDOW_NOMOVE))
	tab.RawSetString("FLAGWINDOW_NOSCROLLBAR", golua.LNumber(FLAGWINDOW_NOSCROLLBAR))
	tab.RawSetString("FLAGWINDOW_NOSCROLLWITHMOUSE", golua.LNumber(FLAGWINDOW_NOSCROLLWITHMOUSE))
	tab.RawSetString("FLAGWINDOW_NOCOLLAPSE", golua.LNumber(FLAGWINDOW_NOCOLLAPSE))
	tab.RawSetString("FLAGWINDOW_ALWAYSAUTORESIZE", golua.LNumber(FLAGWINDOW_ALWAYSAUTORESIZE))
	tab.RawSetString("FLAGWINDOW_NOBACKGROUND", golua.LNumber(FLAGWINDOW_NOBACKGROUND))
	tab.RawSetString("FLAGWINDOW_NOSAVEDSETTINGS", golua.LNumber(FLAGWINDOW_NOSAVEDSETTINGS))
	tab.RawSetString("FLAGWINDOW_NOMOUSEINPUTS", golua.LNumber(FLAGWINDOW_NOMOUSEINPUTS))
	tab.RawSetString("FLAGWINDOW_MENUBAR", golua.LNumber(FLAGWINDOW_MENUBAR))
	tab.RawSetString("FLAGWINDOW_HORIZONTALSCROLLBAR", golua.LNumber(FLAGWINDOW_HORIZONTALSCROLLBAR))
	tab.RawSetString("FLAGWINDOW_NOFOCUSONAPPEARING", golua.LNumber(FLAGWINDOW_NOFOCUSONAPPEARING))
	tab.RawSetString("FLAGWINDOW_NOBRINGTOFRONTONFOCUS", golua.LNumber(FLAGWINDOW_NOBRINGTOFRONTONFOCUS))
	tab.RawSetString("FLAGWINDOW_ALWAYSVERTICALSCROLLBAR", golua.LNumber(FLAGWINDOW_ALWAYSVERTICALSCROLLBAR))
	tab.RawSetString("FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR", golua.LNumber(FLAGWINDOW_ALWAYSHORIZONTALSCROLLBAR))
	tab.RawSetString("FLAGWINDOW_NONAVINPUTS", golua.LNumber(FLAGWINDOW_NONAVINPUTS))
	tab.RawSetString("FLAGWINDOW_NONAVFOCUS", golua.LNumber(FLAGWINDOW_NONAVFOCUS))
	tab.RawSetString("FLAGWINDOW_UNSAVEDDOCUMENT", golua.LNumber(FLAGWINDOW_UNSAVEDDOCUMENT))
	tab.RawSetString("FLAGWINDOW_NONAV", golua.LNumber(FLAGWINDOW_NONAV))
	tab.RawSetString("FLAGWINDOW_NODECORATION", golua.LNumber(FLAGWINDOW_NODECORATION))
	tab.RawSetString("FLAGWINDOW_NOINPUTS", golua.LNumber(FLAGWINDOW_NOINPUTS))

	/// @constants SplitDirection {int}
	/// @const SPLITDIRECTION_HORIZONTAL
	/// @const SPLITDIRECTION_VERTICAL
	tab.RawSetString("SPLITDIRECTION_HORIZONTAL", golua.LNumber(SPLITDIRECTION_HORIZONTAL))
	tab.RawSetString("SPLITDIRECTION_VERTICAL", golua.LNumber(SPLITDIRECTION_VERTICAL))

	/// @constants Alignment {int}
	/// @const ALIGN_LEFT
	/// @const ALIGN_CENTER
	/// @const ALIGN_RIGHT
	tab.RawSetString("ALIGN_LEFT", golua.LNumber(ALIGN_LEFT))
	tab.RawSetString("ALIGN_CENTER", golua.LNumber(ALIGN_CENTER))
	tab.RawSetString("ALIGN_RIGHT", golua.LNumber(ALIGN_RIGHT))

	/// @constants MSGBoxButtons {int}
	/// @const MSGBOXBUTTONS_YESNO
	/// @const MSGBOXBUTTONS_OKCANCEL
	/// @const MSGBOXBUTTONS_OK
	tab.RawSetString("MSGBOXBUTTONS_YESNO", golua.LNumber(MSGBOXBUTTONS_YESNO))
	tab.RawSetString("MSGBOXBUTTONS_OKCANCEL", golua.LNumber(MSGBOXBUTTONS_OKCANCEL))
	tab.RawSetString("MSGBOXBUTTONS_OK", golua.LNumber(MSGBOXBUTTONS_OK))

	/// @constants StyleColorID {int}
	/// @const COLID_TEXT
	/// @const COLID_TEXTDISABLED
	/// @const COLID_WINDOWBG
	/// @const COLID_CHILDBG
	/// @const COLID_POPUPBG
	/// @const COLID_BORDER
	/// @const COLID_BORDERSHADOW
	/// @const COLID_FRAMEBG
	/// @const COLID_FRAMEBGHOVERED
	/// @const COLID_FRAMEBGACTIVE
	/// @const COLID_TITLEBG
	/// @const COLID_TITLEBGACTIVE
	/// @const COLID_TITLEBGCOLLAPSED
	/// @const COLID_MENUBARBG
	/// @const COLID_SCROLLBARBG
	/// @const COLID_SCROLLBARGRAB
	/// @const COLID_SCROLLBARGRABHOVERED
	/// @const COLID_SCROLLBARGRABACTIVE
	/// @const COLID_CHECKMARK
	/// @const COLID_SLIDERGRAB
	/// @const COLID_SLIDERGRABACTIVE
	/// @const COLID_BUTTON
	/// @const COLID_BUTTONHOVERED
	/// @const COLID_BUTTONACTIVE
	/// @const COLID_HEADER
	/// @const COLID_HEADERHOVERED
	/// @const COLID_HEADERACTIVE
	/// @const COLID_SEPARATOR
	/// @const COLID_SEPARATORHOVERED
	/// @const COLID_SEPARATORACTIVE
	/// @const COLID_RESIZEGRIP
	/// @const COLID_RESIZEGRIPHOVERED
	/// @const COLID_RESIZEGRIPACTIVE
	/// @const COLID_TAB
	/// @const COLID_TABHOVERED
	/// @const COLID_TABACTIVE
	/// @const COLID_TABUNFOCUSED
	/// @const COLID_TABUNFOCUSEDACTIVE
	/// @const COLID_DOCKINGPREVIEW
	/// @const COLID_DOCKINGEMPTYBG
	/// @const COLID_PLOTLINES
	/// @const COLID_PLOTLINESHOVERED
	/// @const COLID_PLOTHISTOGRAM
	/// @const COLID_PLOTHISTOGRAMHOVERED
	/// @const COLID_TABLEHEADERBG
	/// @const COLID_TABLEBORDERSTRONG
	/// @const COLID_TABLEBORDERLIGHT
	/// @const COLID_TABLEROWBG
	/// @const COLID_TABLEROWBGALT
	/// @const COLID_TEXTSELECTEDBG
	/// @const COLID_DRAGDROPTARGET
	/// @const COLID_NAVHIGHLIGHT
	/// @const COLID_NAVWINDOWINGHIGHLIGHT
	/// @const COLID_NAVWINDOWINGDIMBG
	/// @const COLID_MODALWINDOWDIMBG
	/// @const COLID_COUNT
	tab.RawSetString("COLID_TEXT", golua.LNumber(COLID_TEXT))
	tab.RawSetString("COLID_TEXTDISABLED", golua.LNumber(COLID_TEXTDISABLED))
	tab.RawSetString("COLID_WINDOWBG", golua.LNumber(COLID_WINDOWBG))
	tab.RawSetString("COLID_CHILDBG", golua.LNumber(COLID_CHILDBG))
	tab.RawSetString("COLID_POPUPBG", golua.LNumber(COLID_POPUPBG))
	tab.RawSetString("COLID_BORDER", golua.LNumber(COLID_BORDER))
	tab.RawSetString("COLID_BORDERSHADOW", golua.LNumber(COLID_BORDERSHADOW))
	tab.RawSetString("COLID_FRAMEBG", golua.LNumber(COLID_FRAMEBG))
	tab.RawSetString("COLID_FRAMEBGHOVERED", golua.LNumber(COLID_FRAMEBGHOVERED))
	tab.RawSetString("COLID_FRAMEBGACTIVE", golua.LNumber(COLID_FRAMEBGACTIVE))
	tab.RawSetString("COLID_TITLEBG", golua.LNumber(COLID_TITLEBG))
	tab.RawSetString("COLID_TITLEBGACTIVE", golua.LNumber(COLID_TITLEBGACTIVE))
	tab.RawSetString("COLID_TITLEBGCOLLAPSED", golua.LNumber(COLID_TITLEBGCOLLAPSED))
	tab.RawSetString("COLID_MENUBARBG", golua.LNumber(COLID_MENUBARBG))
	tab.RawSetString("COLID_SCROLLBARBG", golua.LNumber(COLID_SCROLLBARBG))
	tab.RawSetString("COLID_SCROLLBARGRAB", golua.LNumber(COLID_SCROLLBARGRAB))
	tab.RawSetString("COLID_SCROLLBARGRABHOVERED", golua.LNumber(COLID_SCROLLBARGRABHOVERED))
	tab.RawSetString("COLID_SCROLLBARGRABACTIVE", golua.LNumber(COLID_SCROLLBARGRABACTIVE))
	tab.RawSetString("COLID_CHECKMARK", golua.LNumber(COLID_CHECKMARK))
	tab.RawSetString("COLID_SLIDERGRAB", golua.LNumber(COLID_SLIDERGRAB))
	tab.RawSetString("COLID_SLIDERGRABACTIVE", golua.LNumber(COLID_SLIDERGRABACTIVE))
	tab.RawSetString("COLID_BUTTON", golua.LNumber(COLID_BUTTON))
	tab.RawSetString("COLID_BUTTONHOVERED", golua.LNumber(COLID_BUTTONHOVERED))
	tab.RawSetString("COLID_BUTTONACTIVE", golua.LNumber(COLID_BUTTONACTIVE))
	tab.RawSetString("COLID_HEADER", golua.LNumber(COLID_HEADER))
	tab.RawSetString("COLID_HEADERHOVERED", golua.LNumber(COLID_HEADERHOVERED))
	tab.RawSetString("COLID_HEADERACTIVE", golua.LNumber(COLID_HEADERACTIVE))
	tab.RawSetString("COLID_SEPARATOR", golua.LNumber(COLID_SEPARATOR))
	tab.RawSetString("COLID_SEPARATORHOVERED", golua.LNumber(COLID_SEPARATORHOVERED))
	tab.RawSetString("COLID_SEPARATORACTIVE", golua.LNumber(COLID_SEPARATORACTIVE))
	tab.RawSetString("COLID_RESIZEGRIP", golua.LNumber(COLID_RESIZEGRIP))
	tab.RawSetString("COLID_RESIZEGRIPHOVERED", golua.LNumber(COLID_RESIZEGRIPHOVERED))
	tab.RawSetString("COLID_RESIZEGRIPACTIVE", golua.LNumber(COLID_RESIZEGRIPACTIVE))
	tab.RawSetString("COLID_TAB", golua.LNumber(COLID_TAB))
	tab.RawSetString("COLID_TABHOVERED", golua.LNumber(COLID_TABHOVERED))
	tab.RawSetString("COLID_TABACTIVE", golua.LNumber(COLID_TABACTIVE))
	tab.RawSetString("COLID_TABUNFOCUSED", golua.LNumber(COLID_TABUNFOCUSED))
	tab.RawSetString("COLID_TABUNFOCUSEDACTIVE", golua.LNumber(COLID_TABUNFOCUSEDACTIVE))
	tab.RawSetString("COLID_DOCKINGPREVIEW", golua.LNumber(COLID_DOCKINGPREVIEW))
	tab.RawSetString("COLID_DOCKINGEMPTYBG", golua.LNumber(COLID_DOCKINGEMPTYBG))
	tab.RawSetString("COLID_PLOTLINES", golua.LNumber(COLID_PLOTLINES))
	tab.RawSetString("COLID_PLOTLINESHOVERED", golua.LNumber(COLID_PLOTLINESHOVERED))
	tab.RawSetString("COLID_PLOTHISTOGRAM", golua.LNumber(COLID_PLOTHISTOGRAM))
	tab.RawSetString("COLID_PLOTHISTOGRAMHOVERED", golua.LNumber(COLID_PLOTHISTOGRAMHOVERED))
	tab.RawSetString("COLID_TABLEHEADERBG", golua.LNumber(COLID_TABLEHEADERBG))
	tab.RawSetString("COLID_TABLEBORDERSTRONG", golua.LNumber(COLID_TABLEBORDERSTRONG))
	tab.RawSetString("COLID_TABLEBORDERLIGHT", golua.LNumber(COLID_TABLEBORDERLIGHT))
	tab.RawSetString("COLID_TABLEROWBG", golua.LNumber(COLID_TABLEROWBG))
	tab.RawSetString("COLID_TABLEROWBGALT", golua.LNumber(COLID_TABLEROWBGALT))
	tab.RawSetString("COLID_TEXTSELECTEDBG", golua.LNumber(COLID_TEXTSELECTEDBG))
	tab.RawSetString("COLID_DRAGDROPTARGET", golua.LNumber(COLID_DRAGDROPTARGET))
	tab.RawSetString("COLID_NAVHIGHLIGHT", golua.LNumber(COLID_NAVHIGHLIGHT))
	tab.RawSetString("COLID_NAVWINDOWINGHIGHLIGHT", golua.LNumber(COLID_NAVWINDOWINGHIGHLIGHT))
	tab.RawSetString("COLID_NAVWINDOWINGDIMBG", golua.LNumber(COLID_NAVWINDOWINGDIMBG))
	tab.RawSetString("COLID_MODALWINDOWDIMBG", golua.LNumber(COLID_MODALWINDOWDIMBG))
	tab.RawSetString("COLID_COUNT", golua.LNumber(COLID_COUNT))

	/// @constants StyleVarID {int}
	/// @const STYLEVAR_ALPHA
	/// @const STYLEVAR_DISABLEDALPHA
	/// @const STYLEVAR_WINDOWPADDING
	/// @const STYLEVAR_WINDOWROUNDING
	/// @const STYLEVAR_WINDOWBORDERSIZE
	/// @const STYLEVAR_WINDOWMINSIZE
	/// @const STYLEVAR_WINDOWTITLEALIGN
	/// @const STYLEVAR_CHILDROUNDING
	/// @const STYLEVAR_CHILDBORDERSIZE
	/// @const STYLEVAR_POPUPROUNDING
	/// @const STYLEVAR_POPUPBORDERSIZE
	/// @const STYLEVAR_FRAMEPADDING
	/// @const STYLEVAR_FRAMEROUNDING
	/// @const STYLEVAR_FRAMEBORDERSIZE
	/// @const STYLEVAR_ITEMSPACING
	/// @const STYLEVAR_ITEMINNERSPACING
	/// @const STYLEVAR_INDENTSPACING
	/// @const STYLEVAR_CELLPADDING
	/// @const STYLEVAR_SCROLLBARSIZE
	/// @const STYLEVAR_SCROLLBARROUNDING
	/// @const STYLEVAR_GRABMINSIZE
	/// @const STYLEVAR_GRABROUNDING
	/// @const STYLEVAR_TABROUNDING
	/// @const STYLEVAR_TABBARBORDERSIZE
	/// @const STYLEVAR_BUTTONTEXTALIGN
	/// @const STYLEVAR_SELECTABLETEXTALIGN
	/// @const STYLEVAR_SEPARATORTEXTBORDERSIZE
	/// @const STYLEVAR_SEPARATORTEXTALIGN
	/// @const STYLEVAR_SEPARATORTEXTPADDING
	/// @const STYLEVAR_DOCKINGSEPARATORSIZE
	/// @const STYLEVAR_COUNT
	tab.RawSetString("STYLEVAR_ALPHA", golua.LNumber(STYLEVAR_ALPHA))
	tab.RawSetString("STYLEVAR_DISABLEDALPHA", golua.LNumber(STYLEVAR_DISABLEDALPHA))
	tab.RawSetString("STYLEVAR_WINDOWPADDING", golua.LNumber(STYLEVAR_WINDOWPADDING))
	tab.RawSetString("STYLEVAR_WINDOWROUNDING", golua.LNumber(STYLEVAR_WINDOWROUNDING))
	tab.RawSetString("STYLEVAR_WINDOWBORDERSIZE", golua.LNumber(STYLEVAR_WINDOWBORDERSIZE))
	tab.RawSetString("STYLEVAR_WINDOWMINSIZE", golua.LNumber(STYLEVAR_WINDOWMINSIZE))
	tab.RawSetString("STYLEVAR_WINDOWTITLEALIGN", golua.LNumber(STYLEVAR_WINDOWTITLEALIGN))
	tab.RawSetString("STYLEVAR_CHILDROUNDING", golua.LNumber(STYLEVAR_CHILDROUNDING))
	tab.RawSetString("STYLEVAR_CHILDBORDERSIZE", golua.LNumber(STYLEVAR_CHILDBORDERSIZE))
	tab.RawSetString("STYLEVAR_POPUPROUNDING", golua.LNumber(STYLEVAR_POPUPROUNDING))
	tab.RawSetString("STYLEVAR_POPUPBORDERSIZE", golua.LNumber(STYLEVAR_POPUPBORDERSIZE))
	tab.RawSetString("STYLEVAR_FRAMEPADDING", golua.LNumber(STYLEVAR_FRAMEPADDING))
	tab.RawSetString("STYLEVAR_FRAMEROUNDING", golua.LNumber(STYLEVAR_FRAMEROUNDING))
	tab.RawSetString("STYLEVAR_FRAMEBORDERSIZE", golua.LNumber(STYLEVAR_FRAMEBORDERSIZE))
	tab.RawSetString("STYLEVAR_ITEMSPACING", golua.LNumber(STYLEVAR_ITEMSPACING))
	tab.RawSetString("STYLEVAR_ITEMINNERSPACING", golua.LNumber(STYLEVAR_ITEMINNERSPACING))
	tab.RawSetString("STYLEVAR_INDENTSPACING", golua.LNumber(STYLEVAR_INDENTSPACING))
	tab.RawSetString("STYLEVAR_CELLPADDING", golua.LNumber(STYLEVAR_CELLPADDING))
	tab.RawSetString("STYLEVAR_SCROLLBARSIZE", golua.LNumber(STYLEVAR_SCROLLBARSIZE))
	tab.RawSetString("STYLEVAR_SCROLLBARROUNDING", golua.LNumber(STYLEVAR_SCROLLBARROUNDING))
	tab.RawSetString("STYLEVAR_GRABMINSIZE", golua.LNumber(STYLEVAR_GRABMINSIZE))
	tab.RawSetString("STYLEVAR_GRABROUNDING", golua.LNumber(STYLEVAR_GRABROUNDING))
	tab.RawSetString("STYLEVAR_TABROUNDING", golua.LNumber(STYLEVAR_TABROUNDING))
	tab.RawSetString("STYLEVAR_TABBARBORDERSIZE", golua.LNumber(STYLEVAR_TABBARBORDERSIZE))
	tab.RawSetString("STYLEVAR_BUTTONTEXTALIGN", golua.LNumber(STYLEVAR_BUTTONTEXTALIGN))
	tab.RawSetString("STYLEVAR_SELECTABLETEXTALIGN", golua.LNumber(STYLEVAR_SELECTABLETEXTALIGN))
	tab.RawSetString("STYLEVAR_SEPARATORTEXTBORDERSIZE", golua.LNumber(STYLEVAR_SEPARATORTEXTBORDERSIZE))
	tab.RawSetString("STYLEVAR_SEPARATORTEXTALIGN", golua.LNumber(STYLEVAR_SEPARATORTEXTALIGN))
	tab.RawSetString("STYLEVAR_SEPARATORTEXTPADDING", golua.LNumber(STYLEVAR_SEPARATORTEXTPADDING))
	tab.RawSetString("STYLEVAR_DOCKINGSEPARATORSIZE", golua.LNumber(STYLEVAR_DOCKINGSEPARATORSIZE))
	tab.RawSetString("STYLEVAR_COUNT", golua.LNumber(STYLEVAR_COUNT))

	/// @constants Key {int}
	/// @const KEY_NONE
	/// @const KEY_TAB
	/// @const KEY_LEFTARROW
	/// @const KEY_RIGHTARROW
	/// @const KEY_UPARROW
	/// @const KEY_DOWNARROW
	/// @const KEY_PAGEUP
	/// @const KEY_PAGEDOWN
	/// @const KEY_HOME
	/// @const KEY_END
	/// @const KEY_INSERT
	/// @const KEY_DELETE
	/// @const KEY_BACKSPACE
	/// @const KEY_SPACE
	/// @const KEY_ENTER
	/// @const KEY_ESCAPE
	/// @const KEY_LEFTCTRL
	/// @const KEY_LEFTSHIFT
	/// @const KEY_LEFTALT
	/// @const KEY_LEFTSUPER
	/// @const KEY_RIGHTCTRL
	/// @const KEY_RIGHTSHIFT
	/// @const KEY_RIGHTALT
	/// @const KEY_RIGHTSUPER
	/// @const KEY_MENU
	/// @const KEY_0
	/// @const KEY_1
	/// @const KEY_2
	/// @const KEY_3
	/// @const KEY_4
	/// @const KEY_5
	/// @const KEY_6
	/// @const KEY_7
	/// @const KEY_8
	/// @const KEY_9
	/// @const KEY_A
	/// @const KEY_B
	/// @const KEY_C
	/// @const KEY_D
	/// @const KEY_E
	/// @const KEY_F
	/// @const KEY_G
	/// @const KEY_H
	/// @const KEY_I
	/// @const KEY_J
	/// @const KEY_K
	/// @const KEY_L
	/// @const KEY_M
	/// @const KEY_N
	/// @const KEY_O
	/// @const KEY_P
	/// @const KEY_Q
	/// @const KEY_R
	/// @const KEY_S
	/// @const KEY_T
	/// @const KEY_U
	/// @const KEY_V
	/// @const KEY_W
	/// @const KEY_X
	/// @const KEY_Y
	/// @const KEY_Z
	/// @const KEY_F1
	/// @const KEY_F2
	/// @const KEY_F3
	/// @const KEY_F4
	/// @const KEY_F5
	/// @const KEY_F6
	/// @const KEY_F7
	/// @const KEY_F8
	/// @const KEY_F9
	/// @const KEY_F10
	/// @const KEY_F11
	/// @const KEY_F12
	/// @const KEY_F13
	/// @const KEY_F14
	/// @const KEY_F15
	/// @const KEY_F16
	/// @const KEY_F17
	/// @const KEY_F18
	/// @const KEY_F19
	/// @const KEY_F20
	/// @const KEY_F21
	/// @const KEY_F22
	/// @const KEY_F23
	/// @const KEY_F24
	/// @const KEY_APOSTROPHE
	/// @const KEY_COMMA
	/// @const KEY_MINUS
	/// @const KEY_PERIOD
	/// @const KEY_SLASH
	/// @const KEY_SEMICOLON
	/// @const KEY_EQUAL
	/// @const KEY_LEFTBRACKET
	/// @const KEY_BACKSLASH
	/// @const KEY_RIGHTBRACKET
	/// @const KEY_GRAVEACCENT
	/// @const KEY_CAPSLOCK
	/// @const KEY_SCROLLLOCK
	/// @const KEY_NUMLOCK
	/// @const KEY_PRINTSCREEN
	/// @const KEY_PAUSE
	/// @const KEY_KEYPAD0
	/// @const KEY_KEYPAD1
	/// @const KEY_KEYPAD2
	/// @const KEY_KEYPAD3
	/// @const KEY_KEYPAD4
	/// @const KEY_KEYPAD5
	/// @const KEY_KEYPAD6
	/// @const KEY_KEYPAD7
	/// @const KEY_KEYPAD8
	/// @const KEY_KEYPAD9
	/// @const KEY_KEYPADDECIMAL
	/// @const KEY_KEYPADDIVIDE
	/// @const KEY_KEYPADMULTIPLY
	/// @const KEY_KEYPADSUBTRACT
	/// @const KEY_KEYPADADD
	/// @const KEY_KEYPADENTER
	/// @const KEY_KEYPADEQUAL
	/// @const KEY_APPBACK
	/// @const KEY_APPFORWARD
	/// @const KEY_GAMEPADSTART
	/// @const KEY_GAMEPADBACK
	/// @const KEY_GAMEPADFACELEFT
	/// @const KEY_GAMEPADFACERIGHT
	/// @const KEY_GAMEPADFACEUP
	/// @const KEY_GAMEPADFACEDOWN
	/// @const KEY_GAMEPADDPADLEFT
	/// @const KEY_GAMEPADDPADRIGHT
	/// @const KEY_GAMEPADDPADUP
	/// @const KEY_GAMEPADDPADDOWN
	/// @const KEY_GAMEPADL1
	/// @const KEY_GAMEPADR1
	/// @const KEY_GAMEPADL2
	/// @const KEY_GAMEPADR2
	/// @const KEY_GAMEPADL3
	/// @const KEY_GAMEPADR3
	/// @const KEY_GAMEPADLSTICKLEFT
	/// @const KEY_GAMEPADLSTICKRIGHT
	/// @const KEY_GAMEPADLSTICKUP
	/// @const KEY_GAMEPADLSTICKDOWN
	/// @const KEY_GAMEPADRSTICKLEFT
	/// @const KEY_GAMEPADRSTICKRIGHT
	/// @const KEY_GAMEPADRSTICKUP
	/// @const KEY_GAMEPADRSTICKDOWN
	/// @const KEY_MOUSELEFT
	/// @const KEY_MOUSERIGHT
	/// @const KEY_MOUSEMIDDLE
	/// @const KEY_MOUSEX1
	/// @const KEY_MOUSEX2
	/// @const KEY_MOUSEWHEELX
	/// @const KEY_MOUSEWHEELY
	/// @const KEY_RESERVEDFORMODCTRL
	/// @const KEY_RESERVEDFORMODSHIFT
	/// @const KEY_RESERVEDFORMODALT
	/// @const KEY_RESERVEDFORMODSUPER
	/// @const KEY_COUNT
	/// @const KEY_MODNONE
	/// @const KEY_MODCTRL
	/// @const KEY_MODSHIFT
	/// @const KEY_MODALT
	/// @const KEY_MODSUPER
	/// @const KEY_MODSHORTCUT
	/// @const KEY_MODMASK
	/// @const KEY_NAMEDKEYBEGIN
	/// @const KEY_NAMEDKEYEND
	/// @const KEY_NAMEDKEYCOUNT
	/// @const KEY_KEYSDATASIZE
	/// @const KEY_KEYSDATAOFFSET
	tab.RawSetString("KEY_NONE", golua.LNumber(KEY_NONE))
	tab.RawSetString("KEY_TAB", golua.LNumber(KEY_TAB))
	tab.RawSetString("KEY_LEFTARROW", golua.LNumber(KEY_LEFTARROW))
	tab.RawSetString("KEY_RIGHTARROW", golua.LNumber(KEY_RIGHTARROW))
	tab.RawSetString("KEY_UPARROW", golua.LNumber(KEY_UPARROW))
	tab.RawSetString("KEY_DOWNARROW", golua.LNumber(KEY_DOWNARROW))
	tab.RawSetString("KEY_PAGEUP", golua.LNumber(KEY_PAGEUP))
	tab.RawSetString("KEY_PAGEDOWN", golua.LNumber(KEY_PAGEDOWN))
	tab.RawSetString("KEY_HOME", golua.LNumber(KEY_HOME))
	tab.RawSetString("KEY_END", golua.LNumber(KEY_END))
	tab.RawSetString("KEY_INSERT", golua.LNumber(KEY_INSERT))
	tab.RawSetString("KEY_DELETE", golua.LNumber(KEY_DELETE))
	tab.RawSetString("KEY_BACKSPACE", golua.LNumber(KEY_BACKSPACE))
	tab.RawSetString("KEY_SPACE", golua.LNumber(KEY_SPACE))
	tab.RawSetString("KEY_ENTER", golua.LNumber(KEY_ENTER))
	tab.RawSetString("KEY_ESCAPE", golua.LNumber(KEY_ESCAPE))
	tab.RawSetString("KEY_LEFTCTRL", golua.LNumber(KEY_LEFTCTRL))
	tab.RawSetString("KEY_LEFTSHIFT", golua.LNumber(KEY_LEFTSHIFT))
	tab.RawSetString("KEY_LEFTALT", golua.LNumber(KEY_LEFTALT))
	tab.RawSetString("KEY_LEFTSUPER", golua.LNumber(KEY_LEFTSUPER))
	tab.RawSetString("KEY_RIGHTCTRL", golua.LNumber(KEY_RIGHTCTRL))
	tab.RawSetString("KEY_RIGHTSHIFT", golua.LNumber(KEY_RIGHTSHIFT))
	tab.RawSetString("KEY_RIGHTALT", golua.LNumber(KEY_RIGHTALT))
	tab.RawSetString("KEY_RIGHTSUPER", golua.LNumber(KEY_RIGHTSUPER))
	tab.RawSetString("KEY_MENU", golua.LNumber(KEY_MENU))
	tab.RawSetString("KEY_0", golua.LNumber(KEY_0))
	tab.RawSetString("KEY_1", golua.LNumber(KEY_1))
	tab.RawSetString("KEY_2", golua.LNumber(KEY_2))
	tab.RawSetString("KEY_3", golua.LNumber(KEY_3))
	tab.RawSetString("KEY_4", golua.LNumber(KEY_4))
	tab.RawSetString("KEY_5", golua.LNumber(KEY_5))
	tab.RawSetString("KEY_6", golua.LNumber(KEY_6))
	tab.RawSetString("KEY_7", golua.LNumber(KEY_7))
	tab.RawSetString("KEY_8", golua.LNumber(KEY_8))
	tab.RawSetString("KEY_9", golua.LNumber(KEY_9))
	tab.RawSetString("KEY_A", golua.LNumber(KEY_A))
	tab.RawSetString("KEY_B", golua.LNumber(KEY_B))
	tab.RawSetString("KEY_C", golua.LNumber(KEY_C))
	tab.RawSetString("KEY_D", golua.LNumber(KEY_D))
	tab.RawSetString("KEY_E", golua.LNumber(KEY_E))
	tab.RawSetString("KEY_F", golua.LNumber(KEY_F))
	tab.RawSetString("KEY_G", golua.LNumber(KEY_G))
	tab.RawSetString("KEY_H", golua.LNumber(KEY_H))
	tab.RawSetString("KEY_I", golua.LNumber(KEY_I))
	tab.RawSetString("KEY_J", golua.LNumber(KEY_J))
	tab.RawSetString("KEY_K", golua.LNumber(KEY_K))
	tab.RawSetString("KEY_L", golua.LNumber(KEY_L))
	tab.RawSetString("KEY_M", golua.LNumber(KEY_M))
	tab.RawSetString("KEY_N", golua.LNumber(KEY_N))
	tab.RawSetString("KEY_O", golua.LNumber(KEY_O))
	tab.RawSetString("KEY_P", golua.LNumber(KEY_P))
	tab.RawSetString("KEY_Q", golua.LNumber(KEY_Q))
	tab.RawSetString("KEY_R", golua.LNumber(KEY_R))
	tab.RawSetString("KEY_S", golua.LNumber(KEY_S))
	tab.RawSetString("KEY_T", golua.LNumber(KEY_T))
	tab.RawSetString("KEY_U", golua.LNumber(KEY_U))
	tab.RawSetString("KEY_V", golua.LNumber(KEY_V))
	tab.RawSetString("KEY_W", golua.LNumber(KEY_W))
	tab.RawSetString("KEY_X", golua.LNumber(KEY_X))
	tab.RawSetString("KEY_Y", golua.LNumber(KEY_Y))
	tab.RawSetString("KEY_Z", golua.LNumber(KEY_Z))
	tab.RawSetString("KEY_F1", golua.LNumber(KEY_F1))
	tab.RawSetString("KEY_F2", golua.LNumber(KEY_F2))
	tab.RawSetString("KEY_F3", golua.LNumber(KEY_F3))
	tab.RawSetString("KEY_F4", golua.LNumber(KEY_F4))
	tab.RawSetString("KEY_F5", golua.LNumber(KEY_F5))
	tab.RawSetString("KEY_F6", golua.LNumber(KEY_F6))
	tab.RawSetString("KEY_F7", golua.LNumber(KEY_F7))
	tab.RawSetString("KEY_F8", golua.LNumber(KEY_F8))
	tab.RawSetString("KEY_F9", golua.LNumber(KEY_F9))
	tab.RawSetString("KEY_F10", golua.LNumber(KEY_F10))
	tab.RawSetString("KEY_F11", golua.LNumber(KEY_F11))
	tab.RawSetString("KEY_F12", golua.LNumber(KEY_F12))
	tab.RawSetString("KEY_F13", golua.LNumber(KEY_F13))
	tab.RawSetString("KEY_F14", golua.LNumber(KEY_F14))
	tab.RawSetString("KEY_F15", golua.LNumber(KEY_F15))
	tab.RawSetString("KEY_F16", golua.LNumber(KEY_F16))
	tab.RawSetString("KEY_F17", golua.LNumber(KEY_F17))
	tab.RawSetString("KEY_F18", golua.LNumber(KEY_F18))
	tab.RawSetString("KEY_F19", golua.LNumber(KEY_F19))
	tab.RawSetString("KEY_F20", golua.LNumber(KEY_F20))
	tab.RawSetString("KEY_F21", golua.LNumber(KEY_F21))
	tab.RawSetString("KEY_F22", golua.LNumber(KEY_F22))
	tab.RawSetString("KEY_F23", golua.LNumber(KEY_F23))
	tab.RawSetString("KEY_F24", golua.LNumber(KEY_F24))
	tab.RawSetString("KEY_APOSTROPHE", golua.LNumber(KEY_APOSTROPHE))
	tab.RawSetString("KEY_COMMA", golua.LNumber(KEY_COMMA))
	tab.RawSetString("KEY_MINUS", golua.LNumber(KEY_MINUS))
	tab.RawSetString("KEY_PERIOD", golua.LNumber(KEY_PERIOD))
	tab.RawSetString("KEY_SLASH", golua.LNumber(KEY_SLASH))
	tab.RawSetString("KEY_SEMICOLON", golua.LNumber(KEY_SEMICOLON))
	tab.RawSetString("KEY_EQUAL", golua.LNumber(KEY_EQUAL))
	tab.RawSetString("KEY_LEFTBRACKET", golua.LNumber(KEY_LEFTBRACKET))
	tab.RawSetString("KEY_BACKSLASH", golua.LNumber(KEY_BACKSLASH))
	tab.RawSetString("KEY_RIGHTBRACKET", golua.LNumber(KEY_RIGHTBRACKET))
	tab.RawSetString("KEY_GRAVEACCENT", golua.LNumber(KEY_GRAVEACCENT))
	tab.RawSetString("KEY_CAPSLOCK", golua.LNumber(KEY_CAPSLOCK))
	tab.RawSetString("KEY_SCROLLLOCK", golua.LNumber(KEY_SCROLLLOCK))
	tab.RawSetString("KEY_NUMLOCK", golua.LNumber(KEY_NUMLOCK))
	tab.RawSetString("KEY_PRINTSCREEN", golua.LNumber(KEY_PRINTSCREEN))
	tab.RawSetString("KEY_PAUSE", golua.LNumber(KEY_PAUSE))
	tab.RawSetString("KEY_KEYPAD0", golua.LNumber(KEY_KEYPAD0))
	tab.RawSetString("KEY_KEYPAD1", golua.LNumber(KEY_KEYPAD1))
	tab.RawSetString("KEY_KEYPAD2", golua.LNumber(KEY_KEYPAD2))
	tab.RawSetString("KEY_KEYPAD3", golua.LNumber(KEY_KEYPAD3))
	tab.RawSetString("KEY_KEYPAD4", golua.LNumber(KEY_KEYPAD4))
	tab.RawSetString("KEY_KEYPAD5", golua.LNumber(KEY_KEYPAD5))
	tab.RawSetString("KEY_KEYPAD6", golua.LNumber(KEY_KEYPAD6))
	tab.RawSetString("KEY_KEYPAD7", golua.LNumber(KEY_KEYPAD7))
	tab.RawSetString("KEY_KEYPAD8", golua.LNumber(KEY_KEYPAD8))
	tab.RawSetString("KEY_KEYPAD9", golua.LNumber(KEY_KEYPAD9))
	tab.RawSetString("KEY_KEYPADDECIMAL", golua.LNumber(KEY_KEYPADDECIMAL))
	tab.RawSetString("KEY_KEYPADDIVIDE", golua.LNumber(KEY_KEYPADDIVIDE))
	tab.RawSetString("KEY_KEYPADMULTIPLY", golua.LNumber(KEY_KEYPADMULTIPLY))
	tab.RawSetString("KEY_KEYPADSUBTRACT", golua.LNumber(KEY_KEYPADSUBTRACT))
	tab.RawSetString("KEY_KEYPADADD", golua.LNumber(KEY_KEYPADADD))
	tab.RawSetString("KEY_KEYPADENTER", golua.LNumber(KEY_KEYPADENTER))
	tab.RawSetString("KEY_KEYPADEQUAL", golua.LNumber(KEY_KEYPADEQUAL))
	tab.RawSetString("KEY_APPBACK", golua.LNumber(KEY_APPBACK))
	tab.RawSetString("KEY_APPFORWARD", golua.LNumber(KEY_APPFORWARD))
	tab.RawSetString("KEY_GAMEPADSTART", golua.LNumber(KEY_GAMEPADSTART))
	tab.RawSetString("KEY_GAMEPADBACK", golua.LNumber(KEY_GAMEPADBACK))
	tab.RawSetString("KEY_GAMEPADFACELEFT", golua.LNumber(KEY_GAMEPADFACELEFT))
	tab.RawSetString("KEY_GAMEPADFACERIGHT", golua.LNumber(KEY_GAMEPADFACERIGHT))
	tab.RawSetString("KEY_GAMEPADFACEUP", golua.LNumber(KEY_GAMEPADFACEUP))
	tab.RawSetString("KEY_GAMEPADFACEDOWN", golua.LNumber(KEY_GAMEPADFACEDOWN))
	tab.RawSetString("KEY_GAMEPADDPADLEFT", golua.LNumber(KEY_GAMEPADDPADLEFT))
	tab.RawSetString("KEY_GAMEPADDPADRIGHT", golua.LNumber(KEY_GAMEPADDPADRIGHT))
	tab.RawSetString("KEY_GAMEPADDPADUP", golua.LNumber(KEY_GAMEPADDPADUP))
	tab.RawSetString("KEY_GAMEPADDPADDOWN", golua.LNumber(KEY_GAMEPADDPADDOWN))
	tab.RawSetString("KEY_GAMEPADL1", golua.LNumber(KEY_GAMEPADL1))
	tab.RawSetString("KEY_GAMEPADR1", golua.LNumber(KEY_GAMEPADR1))
	tab.RawSetString("KEY_GAMEPADL2", golua.LNumber(KEY_GAMEPADL2))
	tab.RawSetString("KEY_GAMEPADR2", golua.LNumber(KEY_GAMEPADR2))
	tab.RawSetString("KEY_GAMEPADL3", golua.LNumber(KEY_GAMEPADL3))
	tab.RawSetString("KEY_GAMEPADR3", golua.LNumber(KEY_GAMEPADR3))
	tab.RawSetString("KEY_GAMEPADLSTICKLEFT", golua.LNumber(KEY_GAMEPADLSTICKLEFT))
	tab.RawSetString("KEY_GAMEPADLSTICKRIGHT", golua.LNumber(KEY_GAMEPADLSTICKRIGHT))
	tab.RawSetString("KEY_GAMEPADLSTICKUP", golua.LNumber(KEY_GAMEPADLSTICKUP))
	tab.RawSetString("KEY_GAMEPADLSTICKDOWN", golua.LNumber(KEY_GAMEPADLSTICKDOWN))
	tab.RawSetString("KEY_GAMEPADRSTICKLEFT", golua.LNumber(KEY_GAMEPADRSTICKLEFT))
	tab.RawSetString("KEY_GAMEPADRSTICKRIGHT", golua.LNumber(KEY_GAMEPADRSTICKRIGHT))
	tab.RawSetString("KEY_GAMEPADRSTICKUP", golua.LNumber(KEY_GAMEPADRSTICKUP))
	tab.RawSetString("KEY_GAMEPADRSTICKDOWN", golua.LNumber(KEY_GAMEPADRSTICKDOWN))
	tab.RawSetString("KEY_MOUSELEFT", golua.LNumber(KEY_MOUSELEFT))
	tab.RawSetString("KEY_MOUSERIGHT", golua.LNumber(KEY_MOUSERIGHT))
	tab.RawSetString("KEY_MOUSEMIDDLE", golua.LNumber(KEY_MOUSEMIDDLE))
	tab.RawSetString("KEY_MOUSEX1", golua.LNumber(KEY_MOUSEX1))
	tab.RawSetString("KEY_MOUSEX2", golua.LNumber(KEY_MOUSEX2))
	tab.RawSetString("KEY_MOUSEWHEELX", golua.LNumber(KEY_MOUSEWHEELX))
	tab.RawSetString("KEY_MOUSEWHEELY", golua.LNumber(KEY_MOUSEWHEELY))
	tab.RawSetString("KEY_RESERVEDFORMODCTRL", golua.LNumber(KEY_RESERVEDFORMODCTRL))
	tab.RawSetString("KEY_RESERVEDFORMODSHIFT", golua.LNumber(KEY_RESERVEDFORMODSHIFT))
	tab.RawSetString("KEY_RESERVEDFORMODALT", golua.LNumber(KEY_RESERVEDFORMODALT))
	tab.RawSetString("KEY_RESERVEDFORMODSUPER", golua.LNumber(KEY_RESERVEDFORMODSUPER))
	tab.RawSetString("KEY_COUNT", golua.LNumber(KEY_COUNT))
	tab.RawSetString("KEY_MODNONE", golua.LNumber(KEY_MODNONE))
	tab.RawSetString("KEY_MODCTRL", golua.LNumber(KEY_MODCTRL))
	tab.RawSetString("KEY_MODSHIFT", golua.LNumber(KEY_MODSHIFT))
	tab.RawSetString("KEY_MODALT", golua.LNumber(KEY_MODALT))
	tab.RawSetString("KEY_MODSUPER", golua.LNumber(KEY_MODSUPER))
	tab.RawSetString("KEY_MODSHORTCUT", golua.LNumber(KEY_MODSHORTCUT))
	tab.RawSetString("KEY_MODMASK", golua.LNumber(KEY_MODMASK))
	tab.RawSetString("KEY_NAMEDKEYBEGIN", golua.LNumber(KEY_NAMEDKEYBEGIN))
	tab.RawSetString("KEY_NAMEDKEYEND", golua.LNumber(KEY_NAMEDKEYEND))
	tab.RawSetString("KEY_NAMEDKEYCOUNT", golua.LNumber(KEY_NAMEDKEYCOUNT))
	tab.RawSetString("KEY_KEYSDATASIZE", golua.LNumber(KEY_KEYSDATASIZE))
	tab.RawSetString("KEY_KEYSDATAOFFSET", golua.LNumber(KEY_KEYSDATAOFFSET))

	/// @constants Condition {int}
	/// @const COND_NONE
	/// @const COND_ALWAYS
	/// @const COND_ONCE
	/// @const COND_FIRSTUSEEVER
	/// @const COND_APPEARING
	tab.RawSetString("COND_NONE", golua.LNumber(COND_NONE))
	tab.RawSetString("COND_ALWAYS", golua.LNumber(COND_ALWAYS))
	tab.RawSetString("COND_ONCE", golua.LNumber(COND_ONCE))
	tab.RawSetString("COND_FIRSTUSEEVER", golua.LNumber(COND_FIRSTUSEEVER))
	tab.RawSetString("COND_APPEARING", golua.LNumber(COND_APPEARING))

	/// @constants PlotFlags {int}
	/// @const FLAGPLOT_NONE
	/// @const FLAGPLOT_NOTITLE
	/// @const FLAGPLOT_NOLEGEND
	/// @const FLAGPLOT_NOMOUSETEXT
	/// @const FLAGPLOT_NOINPUTS
	/// @const FLAGPLOT_NOMENUS
	/// @const FLAGPLOT_NOBOXSELECT
	/// @const FLAGPLOT_NOFRAME
	/// @const FLAGPLOT_EQUAL
	/// @const FLAGPLOT_CROSSHAIRS
	/// @const FLAGPLOT_CANVASONLY
	tab.RawSetString("FLAGPLOT_NONE", golua.LNumber(FLAGPLOT_NONE))
	tab.RawSetString("FLAGPLOT_NOTITLE", golua.LNumber(FLAGPLOT_NOTITLE))
	tab.RawSetString("FLAGPLOT_NOLEGEND", golua.LNumber(FLAGPLOT_NOLEGEND))
	tab.RawSetString("FLAGPLOT_NOMOUSETEXT", golua.LNumber(FLAGPLOT_NOMOUSETEXT))
	tab.RawSetString("FLAGPLOT_NOINPUTS", golua.LNumber(FLAGPLOT_NOINPUTS))
	tab.RawSetString("FLAGPLOT_NOMENUS", golua.LNumber(FLAGPLOT_NOMENUS))
	tab.RawSetString("FLAGPLOT_NOBOXSELECT", golua.LNumber(FLAGPLOT_NOBOXSELECT))
	tab.RawSetString("FLAGPLOT_NOFRAME", golua.LNumber(FLAGPLOT_NOFRAME))
	tab.RawSetString("FLAGPLOT_EQUAL", golua.LNumber(FLAGPLOT_EQUAL))
	tab.RawSetString("FLAGPLOT_CROSSHAIRS", golua.LNumber(FLAGPLOT_CROSSHAIRS))
	tab.RawSetString("FLAGPLOT_CANVASONLY", golua.LNumber(FLAGPLOT_CANVASONLY))

	/// @constants PlotAxis {int}
	/// @const PLOTAXIS_X1
	/// @const PLOTAXIS_X2
	/// @const PLOTAXIS_X3
	/// @const PLOTAXIS_Y1
	/// @const PLOTAXIS_Y2
	/// @const PLOTAXIS_Y3
	/// @const PLOTAXIS_COUNT
	tab.RawSetString("PLOTAXIS_X1", golua.LNumber(PLOTAXIS_X1))
	tab.RawSetString("PLOTAXIS_X2", golua.LNumber(PLOTAXIS_X2))
	tab.RawSetString("PLOTAXIS_X3", golua.LNumber(PLOTAXIS_X3))
	tab.RawSetString("PLOTAXIS_Y1", golua.LNumber(PLOTAXIS_Y1))
	tab.RawSetString("PLOTAXIS_Y2", golua.LNumber(PLOTAXIS_Y2))
	tab.RawSetString("PLOTAXIS_Y3", golua.LNumber(PLOTAXIS_Y3))
	tab.RawSetString("PLOTAXIS_COUNT", golua.LNumber(PLOTAXIS_COUNT))

	/// @constants PlotAxisFlags {int}
	/// @const FLAGPLOTAXIS_NONE
	/// @const FLAGPLOTAXIS_NOLABEL
	/// @const FLAGPLOTAXIS_NOGRIDLINES
	/// @const FLAGPLOTAXIS_NOTICKMARKS
	/// @const FLAGPLOTAXIS_NOTICKLABELS
	/// @const FLAGPLOTAXIS_NOINITIALFIT
	/// @const FLAGPLOTAXIS_NOMENUS
	/// @const FLAGPLOTAXIS_NOSIDESWITCH
	/// @const FLAGPLOTAXIS_NOHIGHLIGHT
	/// @const FLAGPLOTAXIS_OPPOSITE
	/// @const FLAGPLOTAXIS_FOREGROUND
	/// @const FLAGPLOTAXIS_INVERT
	/// @const FLAGPLOTAXIS_AUTOFIT
	/// @const FLAGPLOTAXIS_RANGEFIT
	/// @const FLAGPLOTAXIS_PANSTRETCH
	/// @const FLAGPLOTAXIS_LOCKMIN
	/// @const FLAGPLOTAXIS_LOCKMAX
	/// @const FLAGPLOTAXIS_LOCK
	/// @const FLAGPLOTAXIS_NODECORATIONS
	/// @const FLAGPLOTAXIS_AUXDEFAULT
	tab.RawSetString("FLAGPLOTAXIS_NONE", golua.LNumber(FLAGPLOTAXIS_NONE))
	tab.RawSetString("FLAGPLOTAXIS_NOLABEL", golua.LNumber(FLAGPLOTAXIS_NOLABEL))
	tab.RawSetString("FLAGPLOTAXIS_NOGRIDLINES", golua.LNumber(FLAGPLOTAXIS_NOGRIDLINES))
	tab.RawSetString("FLAGPLOTAXIS_NOTICKMARKS", golua.LNumber(FLAGPLOTAXIS_NOTICKMARKS))
	tab.RawSetString("FLAGPLOTAXIS_NOTICKLABELS", golua.LNumber(FLAGPLOTAXIS_NOTICKLABELS))
	tab.RawSetString("FLAGPLOTAXIS_NOINITIALFIT", golua.LNumber(FLAGPLOTAXIS_NOINITIALFIT))
	tab.RawSetString("FLAGPLOTAXIS_NOMENUS", golua.LNumber(FLAGPLOTAXIS_NOMENUS))
	tab.RawSetString("FLAGPLOTAXIS_NOSIDESWITCH", golua.LNumber(FLAGPLOTAXIS_NOSIDESWITCH))
	tab.RawSetString("FLAGPLOTAXIS_NOHIGHLIGHT", golua.LNumber(FLAGPLOTAXIS_NOHIGHLIGHT))
	tab.RawSetString("FLAGPLOTAXIS_OPPOSITE", golua.LNumber(FLAGPLOTAXIS_OPPOSITE))
	tab.RawSetString("FLAGPLOTAXIS_FOREGROUND", golua.LNumber(FLAGPLOTAXIS_FOREGROUND))
	tab.RawSetString("FLAGPLOTAXIS_INVERT", golua.LNumber(FLAGPLOTAXIS_INVERT))
	tab.RawSetString("FLAGPLOTAXIS_AUTOFIT", golua.LNumber(FLAGPLOTAXIS_AUTOFIT))
	tab.RawSetString("FLAGPLOTAXIS_RANGEFIT", golua.LNumber(FLAGPLOTAXIS_RANGEFIT))
	tab.RawSetString("FLAGPLOTAXIS_PANSTRETCH", golua.LNumber(FLAGPLOTAXIS_PANSTRETCH))
	tab.RawSetString("FLAGPLOTAXIS_LOCKMIN", golua.LNumber(FLAGPLOTAXIS_LOCKMIN))
	tab.RawSetString("FLAGPLOTAXIS_LOCKMAX", golua.LNumber(FLAGPLOTAXIS_LOCKMAX))
	tab.RawSetString("FLAGPLOTAXIS_LOCK", golua.LNumber(FLAGPLOTAXIS_LOCK))
	tab.RawSetString("FLAGPLOTAXIS_NODECORATIONS", golua.LNumber(FLAGPLOTAXIS_NODECORATIONS))
	tab.RawSetString("FLAGPLOTAXIS_AUXDEFAULT", golua.LNumber(FLAGPLOTAXIS_AUXDEFAULT))

	/// @constants PlotYAxis {int}
	/// @const PLOTYAXIS_LEFT
	/// @const PLOTYAXIS_FIRSTONRIGHT
	/// @const PLOTYAXIS_SECONDONRIGHT
	tab.RawSetString("PLOTYAXIS_LEFT", golua.LNumber(PLOTYAXIS_LEFT))
	tab.RawSetString("PLOTYAXIS_FIRSTONRIGHT", golua.LNumber(PLOTYAXIS_FIRSTONRIGHT))
	tab.RawSetString("PLOTYAXIS_SECONDONRIGHT", golua.LNumber(PLOTYAXIS_SECONDONRIGHT))

	/// @constants DrawFlags {int}
	/// @const FLAGDRAW_NONE
	/// @const FLAGDRAW_CLOSED
	/// @const FLAGDRAW_ROUNDCORNERSTOPLEFT
	/// @const FLAGDRAW_ROUNDCORNERSTOPRIGHT
	/// @const FLAGDRAW_ROUNDCORNERSBOTTOMLEFT
	/// @const FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT
	/// @const FLAGDRAW_ROUNDCORNERSNONE
	/// @const FLAGDRAW_ROUNDCORNERSTOP
	/// @const FLAGDRAW_ROUNDCORNERSBOTTOM
	/// @const FLAGDRAW_ROUNDCORNERSLEFT
	/// @const FLAGDRAW_ROUNDCORNERSRIGHT
	/// @const FLAGDRAW_ROUNDCORNERSALL
	/// @const FLAGDRAW_ROUNDCORNERSDEFAULT
	/// @const FLAGDRAW_ROUNDCORNERSMASK
	tab.RawSetString("FLAGDRAW_NONE", golua.LNumber(FLAGDRAW_NONE))
	tab.RawSetString("FLAGDRAW_CLOSED", golua.LNumber(FLAGDRAW_CLOSED))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSTOPLEFT", golua.LNumber(FLAGDRAW_ROUNDCORNERSTOPLEFT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSTOPRIGHT", golua.LNumber(FLAGDRAW_ROUNDCORNERSTOPRIGHT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSBOTTOMLEFT", golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOMLEFT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT", golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSNONE", golua.LNumber(FLAGDRAW_ROUNDCORNERSNONE))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSTOP", golua.LNumber(FLAGDRAW_ROUNDCORNERSTOP))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSBOTTOM", golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOM))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSLEFT", golua.LNumber(FLAGDRAW_ROUNDCORNERSLEFT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSRIGHT", golua.LNumber(FLAGDRAW_ROUNDCORNERSRIGHT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSALL", golua.LNumber(FLAGDRAW_ROUNDCORNERSALL))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSDEFAULT", golua.LNumber(FLAGDRAW_ROUNDCORNERSDEFAULT))
	tab.RawSetString("FLAGDRAW_ROUNDCORNERSMASK", golua.LNumber(FLAGDRAW_ROUNDCORNERSMASK))

	/// @constants FocusedFlags {int}
	/// @const FLAGFOCUS_NONE
	/// @const FLAGFOCUS_CHILDWINDOWS
	/// @const FLAGFOCUS_ROOTWINDOW
	/// @const FLAGFOCUS_ANYWINDOW
	/// @const FLAGFOCUS_NOPOPUPHIERARCHY
	/// @const FLAGFOCUS_DOCKHIERARCHY
	/// @const FLAGFOCUS_ROOTANDCHILDWINDOWS
	tab.RawSetString("FLAGFOCUS_NONE", golua.LNumber(FLAGFOCUS_NONE))
	tab.RawSetString("FLAGFOCUS_CHILDWINDOWS", golua.LNumber(FLAGFOCUS_CHILDWINDOWS))
	tab.RawSetString("FLAGFOCUS_ROOTWINDOW", golua.LNumber(FLAGFOCUS_ROOTWINDOW))
	tab.RawSetString("FLAGFOCUS_ANYWINDOW", golua.LNumber(FLAGFOCUS_ANYWINDOW))
	tab.RawSetString("FLAGFOCUS_NOPOPUPHIERARCHY", golua.LNumber(FLAGFOCUS_NOPOPUPHIERARCHY))
	tab.RawSetString("FLAGFOCUS_DOCKHIERARCHY", golua.LNumber(FLAGFOCUS_DOCKHIERARCHY))
	tab.RawSetString("FLAGFOCUS_ROOTANDCHILDWINDOWS", golua.LNumber(FLAGFOCUS_ROOTANDCHILDWINDOWS))

	/// @constants HoveredFlags {int}
	/// @const FLAGHOVERED_NONE
	/// @const FLAGHOVERED_CHILDWINDOWS
	/// @const FLAGHOVERED_ROOTWINDOW
	/// @const FLAGHOVERED_ANYWINDOW
	/// @const FLAGHOVERED_NOPOPUPHIERARCHY
	/// @const FLAGHOVERED_DOCKHIERARCHY
	/// @const FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP
	/// @const FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM
	/// @const FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM
	/// @const FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW
	/// @const FLAGHOVERED_ALLOWWHENDISABLED
	/// @const FLAGHOVERED_NONAVOVERRIDE
	/// @const FLAGHOVERED_ALLOWWHENOVERLAPPED
	/// @const FLAGHOVERED_RECTONLY
	/// @const FLAGHOVERED_ROOTANDCHILDWINDOWS
	/// @const FLAGHOVERED_FORTOOLTIP
	/// @const FLAGHOVERED_STATIONARY
	/// @const FLAGHOVERED_DELAYNONE
	/// @const FLAGHOVERED_DELAYSHORT
	/// @const FLAGHOVERED_DELAYNORMAL
	/// @const FLAGHOVERED_NOSHAREDDELAY
	tab.RawSetString("FLAGHOVERED_NONE", golua.LNumber(FLAGHOVERED_NONE))
	tab.RawSetString("FLAGHOVERED_CHILDWINDOWS", golua.LNumber(FLAGHOVERED_CHILDWINDOWS))
	tab.RawSetString("FLAGHOVERED_ROOTWINDOW", golua.LNumber(FLAGHOVERED_ROOTWINDOW))
	tab.RawSetString("FLAGHOVERED_ANYWINDOW", golua.LNumber(FLAGHOVERED_ANYWINDOW))
	tab.RawSetString("FLAGHOVERED_NOPOPUPHIERARCHY", golua.LNumber(FLAGHOVERED_NOPOPUPHIERARCHY))
	tab.RawSetString("FLAGHOVERED_DOCKHIERARCHY", golua.LNumber(FLAGHOVERED_DOCKHIERARCHY))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP", golua.LNumber(FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM", golua.LNumber(FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM", golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW", golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENDISABLED", golua.LNumber(FLAGHOVERED_ALLOWWHENDISABLED))
	tab.RawSetString("FLAGHOVERED_NONAVOVERRIDE", golua.LNumber(FLAGHOVERED_NONAVOVERRIDE))
	tab.RawSetString("FLAGHOVERED_ALLOWWHENOVERLAPPED", golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPED))
	tab.RawSetString("FLAGHOVERED_RECTONLY", golua.LNumber(FLAGHOVERED_RECTONLY))
	tab.RawSetString("FLAGHOVERED_ROOTANDCHILDWINDOWS", golua.LNumber(FLAGHOVERED_ROOTANDCHILDWINDOWS))
	tab.RawSetString("FLAGHOVERED_FORTOOLTIP", golua.LNumber(FLAGHOVERED_FORTOOLTIP))
	tab.RawSetString("FLAGHOVERED_STATIONARY", golua.LNumber(FLAGHOVERED_STATIONARY))
	tab.RawSetString("FLAGHOVERED_DELAYNONE", golua.LNumber(FLAGHOVERED_DELAYNONE))
	tab.RawSetString("FLAGHOVERED_DELAYSHORT", golua.LNumber(FLAGHOVERED_DELAYSHORT))
	tab.RawSetString("FLAGHOVERED_DELAYNORMAL", golua.LNumber(FLAGHOVERED_DELAYNORMAL))
	tab.RawSetString("FLAGHOVERED_NOSHAREDDELAY", golua.LNumber(FLAGHOVERED_NOSHAREDDELAY))

	/// @constants MouseCursor {int}
	/// @const MOUSECURSOR_NONE
	/// @const MOUSECURSOR_ARROW
	/// @const MOUSECURSOR_TEXTINPUT
	/// @const MOUSECURSOR_RESIZEALL
	/// @const MOUSECURSOR_RESIZENS
	/// @const MOUSECURSOR_RESIZEEW
	/// @const MOUSECURSOR_RESIZENESW
	/// @const MOUSECURSOR_RESIZENWSE
	/// @const MOUSECURSOR_HAND
	/// @const MOUSECURSOR_NOTALLOWED
	/// @const MOUSECURSOR_COUNT
	tab.RawSetString("MOUSECURSOR_NONE", golua.LNumber(MOUSECURSOR_NONE))
	tab.RawSetString("MOUSECURSOR_ARROW", golua.LNumber(MOUSECURSOR_ARROW))
	tab.RawSetString("MOUSECURSOR_TEXTINPUT", golua.LNumber(MOUSECURSOR_TEXTINPUT))
	tab.RawSetString("MOUSECURSOR_RESIZEALL", golua.LNumber(MOUSECURSOR_RESIZEALL))
	tab.RawSetString("MOUSECURSOR_RESIZENS", golua.LNumber(MOUSECURSOR_RESIZENS))
	tab.RawSetString("MOUSECURSOR_RESIZEEW", golua.LNumber(MOUSECURSOR_RESIZEEW))
	tab.RawSetString("MOUSECURSOR_RESIZENESW", golua.LNumber(MOUSECURSOR_RESIZENESW))
	tab.RawSetString("MOUSECURSOR_RESIZENWSE", golua.LNumber(MOUSECURSOR_RESIZENWSE))
	tab.RawSetString("MOUSECURSOR_HAND", golua.LNumber(MOUSECURSOR_HAND))
	tab.RawSetString("MOUSECURSOR_NOTALLOWED", golua.LNumber(MOUSECURSOR_NOTALLOWED))
	tab.RawSetString("MOUSECURSOR_COUNT", golua.LNumber(MOUSECURSOR_COUNT))

	/// @constants Action {int}
	/// @const ACTION_RELEASE
	/// @const ACTION_PRESS
	/// @const ACTION_REPEAT
	tab.RawSetString("ACTION_RELEASE", golua.LNumber(ACTION_RELEASE))
	tab.RawSetString("ACTION_PRESS", golua.LNumber(ACTION_PRESS))
	tab.RawSetString("ACTION_REPEAT", golua.LNumber(ACTION_REPEAT))

	/// @constants WidgetType {string}
	/// @const WIDGET_LABEL
	/// @const WIDGET_BUTTON
	/// @const WIDGET_DUMMY
	/// @const WIDGET_SEPARATOR
	/// @const WIDGET_BULLET_TEXT
	/// @const WIDGET_BULLET
	/// @const WIDGET_CHECKBOX
	/// @const WIDGET_CHILD
	/// @const WIDGET_COLOR_EDIT
	/// @const WIDGET_COLUMN
	/// @const WIDGET_ROW
	/// @const WIDGET_COMBO_CUSTOM
	/// @const WIDGET_COMBO
	/// @const WIDGET_CONDITION
	/// @const WIDGET_CONTEXT_MENU
	/// @const WIDGET_DATE_PICKER
	/// @const WIDGET_DRAG_INT
	/// @const WIDGET_INPUT_FLOAT
	/// @const WIDGET_INPUT_INT
	/// @const WIDGET_INPUT_TEXT
	/// @const WIDGET_INPUT_TEXT_MULTILINE
	/// @const WIDGET_PROGRESS_BAR
	/// @const WIDGET_PROGRESS_INDICATOR
	/// @const WIDGET_SPACING
	/// @const WIDGET_BUTTON_SMALL
	/// @const WIDGET_BUTTON_RADIO
	/// @const WIDGET_IMAGE_URL
	/// @const WIDGET_IMAGE
	/// @const WIDGET_LIST_BOX
	/// @const WIDGET_LIST_CLIPPER
	/// @const WIDGET_MENU_BAR_MAIN
	/// @const WIDGET_MENU_BAR
	/// @const WIDGET_MENU_ITEM
	/// @const WIDGET_MENU
	/// @const WIDGET_SELECTABLE
	/// @const WIDGET_SLIDER_FLOAT
	/// @const WIDGET_SLIDER_INT
	/// @const WIDGET_VSLIDER_INT
	/// @const WIDGET_TAB_BAR
	/// @const WIDGET_TAB_ITEM
	/// @const WIDGET_TOOLTIP
	/// @const WIDGET_TABLE_COLUMN
	/// @const WIDGET_TABLE_ROW
	/// @const WIDGET_TABLE
	/// @const WIDGET_BUTTON_ARROW
	/// @const WIDGET_TREE_NODE
	/// @const WIDGET_TREE_TABLE_ROW
	/// @const WIDGET_TREE_TABLE
	/// @const WIDGET_WINDOW_SINGLE
	/// @const WIDGET_POPUP_MODAL
	/// @const WIDGET_POPUP
	/// @const WIDGET_LAYOUT_SPLIT
	/// @const WIDGET_SPLITTER
	/// @const WIDGET_STACK
	/// @const WIDGET_ALIGN
	/// @const WIDGET_MSG_BOX
	/// @const WIDGET_MSG_BOX_PREPARE
	/// @const WIDGET_BUTTON_INVISIBLE
	/// @const WIDGET_BUTTON_IMAGE
	/// @const WIDGET_STYLE
	/// @const WIDGET_CUSTOM
	/// @const WIDGET_EVENT_HANDLER
	/// @const WIDGET_PLOT
	/// @const WIDGET_CSS_TAG
	tab.RawSetString("WIDGET_LABEL", golua.LString(WIDGET_LABEL))
	tab.RawSetString("WIDGET_BUTTON", golua.LString(WIDGET_BUTTON))
	tab.RawSetString("WIDGET_DUMMY", golua.LString(WIDGET_DUMMY))
	tab.RawSetString("WIDGET_SEPARATOR", golua.LString(WIDGET_SEPARATOR))
	tab.RawSetString("WIDGET_BULLET_TEXT", golua.LString(WIDGET_BULLET_TEXT))
	tab.RawSetString("WIDGET_BULLET", golua.LString(WIDGET_BULLET))
	tab.RawSetString("WIDGET_CHECKBOX", golua.LString(WIDGET_CHECKBOX))
	tab.RawSetString("WIDGET_CHILD", golua.LString(WIDGET_CHILD))
	tab.RawSetString("WIDGET_COLOR_EDIT", golua.LString(WIDGET_COLOR_EDIT))
	tab.RawSetString("WIDGET_COLUMN", golua.LString(WIDGET_COLUMN))
	tab.RawSetString("WIDGET_ROW", golua.LString(WIDGET_ROW))
	tab.RawSetString("WIDGET_COMBO_CUSTOM", golua.LString(WIDGET_COMBO_CUSTOM))
	tab.RawSetString("WIDGET_COMBO", golua.LString(WIDGET_COMBO))
	tab.RawSetString("WIDGET_CONDITION", golua.LString(WIDGET_CONDITION))
	tab.RawSetString("WIDGET_CONTEXT_MENU", golua.LString(WIDGET_CONTEXT_MENU))
	tab.RawSetString("WIDGET_DATE_PICKER", golua.LString(WIDGET_DATE_PICKER))
	tab.RawSetString("WIDGET_DRAG_INT", golua.LString(WIDGET_DRAG_INT))
	tab.RawSetString("WIDGET_INPUT_FLOAT", golua.LString(WIDGET_INPUT_FLOAT))
	tab.RawSetString("WIDGET_INPUT_INT", golua.LString(WIDGET_INPUT_INT))
	tab.RawSetString("WIDGET_INPUT_TEXT", golua.LString(WIDGET_INPUT_TEXT))
	tab.RawSetString("WIDGET_INPUT_TEXT_MULTILINE", golua.LString(WIDGET_INPUT_MULTILINE_TEXT))
	tab.RawSetString("WIDGET_PROGRESS_BAR", golua.LString(WIDGET_PROGRESS_BAR))
	tab.RawSetString("WIDGET_PROGRESS_INDICATOR", golua.LString(WIDGET_PROGRESS_INDICATOR))
	tab.RawSetString("WIDGET_SPACING", golua.LString(WIDGET_SPACING))
	tab.RawSetString("WIDGET_BUTTON_SMALL", golua.LString(WIDGET_BUTTON_SMALL))
	tab.RawSetString("WIDGET_BUTTON_RADIO", golua.LString(WIDGET_BUTTON_RADIO))
	tab.RawSetString("WIDGET_IMAGE_URL", golua.LString(WIDGET_IMAGE_URL))
	tab.RawSetString("WIDGET_IMAGE", golua.LString(WIDGET_IMAGE))
	tab.RawSetString("WIDGET_LIST_BOX", golua.LString(WIDGET_LIST_BOX))
	tab.RawSetString("WIDGET_LIST_CLIPPER", golua.LString(WIDGET_LIST_CLIPPER))
	tab.RawSetString("WIDGET_MENU_BAR_MAIN", golua.LString(WIDGET_MENU_BAR_MAIN))
	tab.RawSetString("WIDGET_MENU_BAR", golua.LString(WIDGET_MENU_BAR))
	tab.RawSetString("WIDGET_MENU_ITEM", golua.LString(WIDGET_MENU_ITEM))
	tab.RawSetString("WIDGET_MENU", golua.LString(WIDGET_MENU))
	tab.RawSetString("WIDGET_SELECTABLE", golua.LString(WIDGET_SELECTABLE))
	tab.RawSetString("WIDGET_SLIDER_FLOAT", golua.LString(WIDGET_SLIDER_FLOAT))
	tab.RawSetString("WIDGET_SLIDER_INT", golua.LString(WIDGET_SLIDER_INT))
	tab.RawSetString("WIDGET_VSLIDER_INT", golua.LString(WIDGET_VSLIDER_INT))
	tab.RawSetString("WIDGET_TAB_BAR", golua.LString(WIDGET_TAB_BAR))
	tab.RawSetString("WIDGET_TAB_ITEM", golua.LString(WIDGET_TAB_ITEM))
	tab.RawSetString("WIDGET_TOOLTIP", golua.LString(WIDGET_TOOLTIP))
	tab.RawSetString("WIDGET_TABLE_COLUMN", golua.LString(WIDGET_TABLE_COLUMN))
	tab.RawSetString("WIDGET_TABLE_ROW", golua.LString(WIDGET_TABLE_ROW))
	tab.RawSetString("WIDGET_TABLE", golua.LString(WIDGET_TABLE))
	tab.RawSetString("WIDGET_BUTTON_ARROW", golua.LString(WIDGET_BUTTON_ARROW))
	tab.RawSetString("WIDGET_TREE_NODE", golua.LString(WIDGET_TREE_NODE))
	tab.RawSetString("WIDGET_TREE_TABLE_ROW", golua.LString(WIDGET_TREE_TABLE_ROW))
	tab.RawSetString("WIDGET_TREE_TABLE", golua.LString(WIDGET_TREE_TABLE))
	tab.RawSetString("WIDGET_WINDOW_SINGLE", golua.LString(WIDGET_WINDOW_SINGLE))
	tab.RawSetString("WIDGET_POPUP_MODAL", golua.LString(WIDGET_POPUP_MODAL))
	tab.RawSetString("WIDGET_POPUP", golua.LString(WIDGET_POPUP))
	tab.RawSetString("WIDGET_LAYOUT_SPLIT", golua.LString(WIDGET_LAYOUT_SPLIT))
	tab.RawSetString("WIDGET_SPLITTER", golua.LString(WIDGET_SPLITTER))
	tab.RawSetString("WIDGET_STACK", golua.LString(WIDGET_STACK))
	tab.RawSetString("WIDGET_ALIGN", golua.LString(WIDGET_ALIGN))
	tab.RawSetString("WIDGET_MSG_BOX", golua.LString(WIDGET_MSG_BOX))
	tab.RawSetString("WIDGET_MSG_BOX_PREPARE", golua.LString(WIDGET_MSG_BOX_PREPARE))
	tab.RawSetString("WIDGET_BUTTON_INVISIBLE", golua.LString(WIDGET_BUTTON_INVISIBLE))
	tab.RawSetString("WIDGET_BUTTON_IMAGE", golua.LString(WIDGET_BUTTON_IMAGE))
	tab.RawSetString("WIDGET_STYLE", golua.LString(WIDGET_STYLE))
	tab.RawSetString("WIDGET_CUSTOM", golua.LString(WIDGET_CUSTOM))
	tab.RawSetString("WIDGET_EVENT_HANDLER", golua.LString(WIDGET_EVENT_HANDLER))
	tab.RawSetString("WIDGET_PLOT", golua.LString(WIDGET_PLOT))
	tab.RawSetString("WIDGET_CSS_TAG", golua.LString(WIDGET_CSS_TAG))

	/// @constants PlotType {string}
	/// @const PLOT_BAR_H
	/// @const PLOT_BAR
	/// @const PLOT_LINE
	/// @const PLOT_LINE_XY
	/// @const PLOT_PIE_CHART
	/// @const PLOT_SCATTER
	/// @const PLOT_SCATTER_XY
	/// @const PLOT_CUSTOM
	tab.RawSetString("PLOT_BAR_H", golua.LString(PLOT_BAR_H))
	tab.RawSetString("PLOT_BAR", golua.LString(PLOT_BAR))
	tab.RawSetString("PLOT_LINE", golua.LString(PLOT_LINE))
	tab.RawSetString("PLOT_LINE_XY", golua.LString(PLOT_LINE_XY))
	tab.RawSetString("PLOT_PIE_CHART", golua.LString(PLOT_PIE_CHART))
	tab.RawSetString("PLOT_SCATTER", golua.LString(PLOT_SCATTER))
	tab.RawSetString("PLOT_SCATTER_XY", golua.LString(PLOT_SCATTER_XY))
	tab.RawSetString("PLOT_CUSTOM", golua.LString(PLOT_CUSTOM))
}

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
	ALIGN_LEFT int = iota
	ALIGN_CENTER
	ALIGN_RIGHT
)

const (
	MSGBOXBUTTONS_YESNO = 1 << iota
	MSGBOXBUTTONS_OKCANCEL
	MSGBOXBUTTONS_OK
)

const (
	COLID_TEXT int = iota
	COLID_TEXTDISABLED
	COLID_WINDOWBG
	COLID_CHILDBG
	COLID_POPUPBG
	COLID_BORDER
	COLID_BORDERSHADOW
	COLID_FRAMEBG
	COLID_FRAMEBGHOVERED
	COLID_FRAMEBGACTIVE
	COLID_TITLEBG
	COLID_TITLEBGACTIVE
	COLID_TITLEBGCOLLAPSED
	COLID_MENUBARBG
	COLID_SCROLLBARBG
	COLID_SCROLLBARGRAB
	COLID_SCROLLBARGRABHOVERED
	COLID_SCROLLBARGRABACTIVE
	COLID_CHECKMARK
	COLID_SLIDERGRAB
	COLID_SLIDERGRABACTIVE
	COLID_BUTTON
	COLID_BUTTONHOVERED
	COLID_BUTTONACTIVE
	COLID_HEADER
	COLID_HEADERHOVERED
	COLID_HEADERACTIVE
	COLID_SEPARATOR
	COLID_SEPARATORHOVERED
	COLID_SEPARATORACTIVE
	COLID_RESIZEGRIP
	COLID_RESIZEGRIPHOVERED
	COLID_RESIZEGRIPACTIVE
	COLID_TAB
	COLID_TABHOVERED
	COLID_TABACTIVE
	COLID_TABUNFOCUSED
	COLID_TABUNFOCUSEDACTIVE
	COLID_DOCKINGPREVIEW
	COLID_DOCKINGEMPTYBG
	COLID_PLOTLINES
	COLID_PLOTLINESHOVERED
	COLID_PLOTHISTOGRAM
	COLID_PLOTHISTOGRAMHOVERED
	COLID_TABLEHEADERBG
	COLID_TABLEBORDERSTRONG
	COLID_TABLEBORDERLIGHT
	COLID_TABLEROWBG
	COLID_TABLEROWBGALT
	COLID_TEXTSELECTEDBG
	COLID_DRAGDROPTARGET
	COLID_NAVHIGHLIGHT
	COLID_NAVWINDOWINGHIGHLIGHT
	COLID_NAVWINDOWINGDIMBG
	COLID_MODALWINDOWDIMBG
	COLID_COUNT
)

const (
	STYLEVAR_ALPHA int = iota
	STYLEVAR_DISABLEDALPHA
	STYLEVAR_WINDOWPADDING
	STYLEVAR_WINDOWROUNDING
	STYLEVAR_WINDOWBORDERSIZE
	STYLEVAR_WINDOWMINSIZE
	STYLEVAR_WINDOWTITLEALIGN
	STYLEVAR_CHILDROUNDING
	STYLEVAR_CHILDBORDERSIZE
	STYLEVAR_POPUPROUNDING
	STYLEVAR_POPUPBORDERSIZE
	STYLEVAR_FRAMEPADDING
	STYLEVAR_FRAMEROUNDING
	STYLEVAR_FRAMEBORDERSIZE
	STYLEVAR_ITEMSPACING
	STYLEVAR_ITEMINNERSPACING
	STYLEVAR_INDENTSPACING
	STYLEVAR_CELLPADDING
	STYLEVAR_SCROLLBARSIZE
	STYLEVAR_SCROLLBARROUNDING
	STYLEVAR_GRABMINSIZE
	STYLEVAR_GRABROUNDING
	STYLEVAR_TABROUNDING
	STYLEVAR_TABBARBORDERSIZE
	STYLEVAR_BUTTONTEXTALIGN
	STYLEVAR_SELECTABLETEXTALIGN
	STYLEVAR_SEPARATORTEXTBORDERSIZE
	STYLEVAR_SEPARATORTEXTALIGN
	STYLEVAR_SEPARATORTEXTPADDING
	STYLEVAR_DOCKINGSEPARATORSIZE
	STYLEVAR_COUNT
)

type Key int

const (
	KEY_NONE                Key = 0
	KEY_TAB                 Key = 512
	KEY_LEFTARROW           Key = 513
	KEY_RIGHTARROW          Key = 514
	KEY_UPARROW             Key = 515
	KEY_DOWNARROW           Key = 516
	KEY_PAGEUP              Key = 517
	KEY_PAGEDOWN            Key = 518
	KEY_HOME                Key = 519
	KEY_END                 Key = 520
	KEY_INSERT              Key = 521
	KEY_DELETE              Key = 522
	KEY_BACKSPACE           Key = 523
	KEY_SPACE               Key = 524
	KEY_ENTER               Key = 525
	KEY_ESCAPE              Key = 526
	KEY_LEFTCTRL            Key = 527
	KEY_LEFTSHIFT           Key = 528
	KEY_LEFTALT             Key = 529
	KEY_LEFTSUPER           Key = 530
	KEY_RIGHTCTRL           Key = 531
	KEY_RIGHTSHIFT          Key = 532
	KEY_RIGHTALT            Key = 533
	KEY_RIGHTSUPER          Key = 534
	KEY_MENU                Key = 535
	KEY_0                   Key = 536
	KEY_1                   Key = 537
	KEY_2                   Key = 538
	KEY_3                   Key = 539
	KEY_4                   Key = 540
	KEY_5                   Key = 541
	KEY_6                   Key = 542
	KEY_7                   Key = 543
	KEY_8                   Key = 544
	KEY_9                   Key = 545
	KEY_A                   Key = 546
	KEY_B                   Key = 547
	KEY_C                   Key = 548
	KEY_D                   Key = 549
	KEY_E                   Key = 550
	KEY_F                   Key = 551
	KEY_G                   Key = 552
	KEY_H                   Key = 553
	KEY_I                   Key = 554
	KEY_J                   Key = 555
	KEY_K                   Key = 556
	KEY_L                   Key = 557
	KEY_M                   Key = 558
	KEY_N                   Key = 559
	KEY_O                   Key = 560
	KEY_P                   Key = 561
	KEY_Q                   Key = 562
	KEY_R                   Key = 563
	KEY_S                   Key = 564
	KEY_T                   Key = 565
	KEY_U                   Key = 566
	KEY_V                   Key = 567
	KEY_W                   Key = 568
	KEY_X                   Key = 569
	KEY_Y                   Key = 570
	KEY_Z                   Key = 571
	KEY_F1                  Key = 572
	KEY_F2                  Key = 573
	KEY_F3                  Key = 574
	KEY_F4                  Key = 575
	KEY_F5                  Key = 576
	KEY_F6                  Key = 577
	KEY_F7                  Key = 578
	KEY_F8                  Key = 579
	KEY_F9                  Key = 580
	KEY_F10                 Key = 581
	KEY_F11                 Key = 582
	KEY_F12                 Key = 583
	KEY_F13                 Key = 584
	KEY_F14                 Key = 585
	KEY_F15                 Key = 586
	KEY_F16                 Key = 587
	KEY_F17                 Key = 588
	KEY_F18                 Key = 589
	KEY_F19                 Key = 590
	KEY_F20                 Key = 591
	KEY_F21                 Key = 592
	KEY_F22                 Key = 593
	KEY_F23                 Key = 594
	KEY_F24                 Key = 595
	KEY_APOSTROPHE          Key = 596
	KEY_COMMA               Key = 597
	KEY_MINUS               Key = 598
	KEY_PERIOD              Key = 599
	KEY_SLASH               Key = 600
	KEY_SEMICOLON           Key = 601
	KEY_EQUAL               Key = 602
	KEY_LEFTBRACKET         Key = 603
	KEY_BACKSLASH           Key = 604
	KEY_RIGHTBRACKET        Key = 605
	KEY_GRAVEACCENT         Key = 606
	KEY_CAPSLOCK            Key = 607
	KEY_SCROLLLOCK          Key = 608
	KEY_NUMLOCK             Key = 609
	KEY_PRINTSCREEN         Key = 610
	KEY_PAUSE               Key = 611
	KEY_KEYPAD0             Key = 612
	KEY_KEYPAD1             Key = 613
	KEY_KEYPAD2             Key = 614
	KEY_KEYPAD3             Key = 615
	KEY_KEYPAD4             Key = 616
	KEY_KEYPAD5             Key = 617
	KEY_KEYPAD6             Key = 618
	KEY_KEYPAD7             Key = 619
	KEY_KEYPAD8             Key = 620
	KEY_KEYPAD9             Key = 621
	KEY_KEYPADDECIMAL       Key = 622
	KEY_KEYPADDIVIDE        Key = 623
	KEY_KEYPADMULTIPLY      Key = 624
	KEY_KEYPADSUBTRACT      Key = 625
	KEY_KEYPADADD           Key = 626
	KEY_KEYPADENTER         Key = 627
	KEY_KEYPADEQUAL         Key = 628
	KEY_APPBACK             Key = 629
	KEY_APPFORWARD          Key = 630
	KEY_GAMEPADSTART        Key = 631
	KEY_GAMEPADBACK         Key = 632
	KEY_GAMEPADFACELEFT     Key = 633
	KEY_GAMEPADFACERIGHT    Key = 634
	KEY_GAMEPADFACEUP       Key = 635
	KEY_GAMEPADFACEDOWN     Key = 636
	KEY_GAMEPADDPADLEFT     Key = 637
	KEY_GAMEPADDPADRIGHT    Key = 638
	KEY_GAMEPADDPADUP       Key = 639
	KEY_GAMEPADDPADDOWN     Key = 640
	KEY_GAMEPADL1           Key = 641
	KEY_GAMEPADR1           Key = 642
	KEY_GAMEPADL2           Key = 643
	KEY_GAMEPADR2           Key = 644
	KEY_GAMEPADL3           Key = 645
	KEY_GAMEPADR3           Key = 646
	KEY_GAMEPADLSTICKLEFT   Key = 647
	KEY_GAMEPADLSTICKRIGHT  Key = 648
	KEY_GAMEPADLSTICKUP     Key = 649
	KEY_GAMEPADLSTICKDOWN   Key = 650
	KEY_GAMEPADRSTICKLEFT   Key = 651
	KEY_GAMEPADRSTICKRIGHT  Key = 652
	KEY_GAMEPADRSTICKUP     Key = 653
	KEY_GAMEPADRSTICKDOWN   Key = 654
	KEY_MOUSELEFT           Key = 655
	KEY_MOUSERIGHT          Key = 656
	KEY_MOUSEMIDDLE         Key = 657
	KEY_MOUSEX1             Key = 658
	KEY_MOUSEX2             Key = 659
	KEY_MOUSEWHEELX         Key = 660
	KEY_MOUSEWHEELY         Key = 661
	KEY_RESERVEDFORMODCTRL  Key = 662
	KEY_RESERVEDFORMODSHIFT Key = 663
	KEY_RESERVEDFORMODALT   Key = 664
	KEY_RESERVEDFORMODSUPER Key = 665
	KEY_COUNT               Key = 666
	KEY_MODNONE             Key = 0
	KEY_MODCTRL             Key = 4096
	KEY_MODSHIFT            Key = 8192
	KEY_MODALT              Key = 16384
	KEY_MODSUPER            Key = 32768
	KEY_MODSHORTCUT         Key = 2048
	KEY_MODMASK             Key = 63488
	KEY_NAMEDKEYBEGIN       Key = 512
	KEY_NAMEDKEYEND         Key = 666
	KEY_NAMEDKEYCOUNT       Key = 154
	KEY_KEYSDATASIZE        Key = 154
	KEY_KEYSDATAOFFSET      Key = 512
)

const (
	COND_NONE         int = 0b0000
	COND_ALWAYS       int = 0b0001
	COND_ONCE         int = 0b0010
	COND_FIRSTUSEEVER int = 0b0100
	COND_APPEARING    int = 0b1000
)

const (
	FLAGPLOT_NONE        int = 0b0000_0000_0000
	FLAGPLOT_NOTITLE     int = 0b0000_0000_0001
	FLAGPLOT_NOLEGEND    int = 0b0000_0000_0010
	FLAGPLOT_NOMOUSETEXT int = 0b0000_0000_0100
	FLAGPLOT_NOINPUTS    int = 0b0000_0000_1000
	FLAGPLOT_NOMENUS     int = 0b0000_0001_0000
	FLAGPLOT_NOBOXSELECT int = 0b0000_0010_0000
	FLAGPLOT_NOFRAME     int = 0b0000_0100_0000
	FLAGPLOT_EQUAL       int = 0b0000_1000_0000
	FLAGPLOT_CROSSHAIRS  int = 0b0001_0000_0000
	FLAGPLOT_CANVASONLY  int = 0b0000_0011_0111
)

const (
	PLOTAXIS_X1 int = iota
	PLOTAXIS_X2
	PLOTAXIS_X3
	PLOTAXIS_Y1
	PLOTAXIS_Y2
	PLOTAXIS_Y3
	PLOTAXIS_COUNT
)

const (
	FLAGPLOTAXIS_NONE          int = 0b0000_0000_0000_0000
	FLAGPLOTAXIS_NOLABEL       int = 0b0000_0000_0000_0001
	FLAGPLOTAXIS_NOGRIDLINES   int = 0b0000_0000_0000_0010
	FLAGPLOTAXIS_NOTICKMARKS   int = 0b0000_0000_0000_0100
	FLAGPLOTAXIS_NOTICKLABELS  int = 0b0000_0000_0000_1000
	FLAGPLOTAXIS_NOINITIALFIT  int = 0b0000_0000_0001_0000
	FLAGPLOTAXIS_NOMENUS       int = 0b0000_0000_0010_0000
	FLAGPLOTAXIS_NOSIDESWITCH  int = 0b0000_0000_0100_0000
	FLAGPLOTAXIS_NOHIGHLIGHT   int = 0b0000_0000_1000_0000
	FLAGPLOTAXIS_OPPOSITE      int = 0b0000_0001_0000_0000
	FLAGPLOTAXIS_FOREGROUND    int = 0b0000_0010_0000_0000
	FLAGPLOTAXIS_INVERT        int = 0b0000_0100_0000_0000
	FLAGPLOTAXIS_AUTOFIT       int = 0b0000_1000_0000_0000
	FLAGPLOTAXIS_RANGEFIT      int = 0b0001_0000_0000_0000
	FLAGPLOTAXIS_PANSTRETCH    int = 0b0010_0000_0000_0000
	FLAGPLOTAXIS_LOCKMIN       int = 0b0100_0000_0000_0000
	FLAGPLOTAXIS_LOCKMAX       int = 0b1000_0000_0000_0000
	FLAGPLOTAXIS_LOCK          int = 0b1100_0000_0000_0000
	FLAGPLOTAXIS_NODECORATIONS int = 0b0000_0000_0000_1111
	FLAGPLOTAXIS_AUXDEFAULT    int = 0b0000_0001_0000_0010
)

const (
	PLOTYAXIS_LEFT          int = 0
	PLOTYAXIS_FIRSTONRIGHT  int = 1
	PLOTYAXIS_SECONDONRIGHT int = 2
)

const (
	FLAGDRAW_NONE                    int = 0b0000_0000_0000
	FLAGDRAW_CLOSED                  int = 0b0000_0000_0001
	FLAGDRAW_ROUNDCORNERSTOPLEFT     int = 0b0000_0001_0000
	FLAGDRAW_ROUNDCORNERSTOPRIGHT    int = 0b0000_0010_0000
	FLAGDRAW_ROUNDCORNERSBOTTOMLEFT  int = 0b0000_0100_0000
	FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT int = 0b0000_1000_0000
	FLAGDRAW_ROUNDCORNERSNONE        int = 0b0001_0000_0000
	FLAGDRAW_ROUNDCORNERSTOP         int = 0b0000_0011_0000
	FLAGDRAW_ROUNDCORNERSBOTTOM      int = 0b0000_1100_0000
	FLAGDRAW_ROUNDCORNERSLEFT        int = 0b0000_0101_0000
	FLAGDRAW_ROUNDCORNERSRIGHT       int = 0b0000_1010_0000
	FLAGDRAW_ROUNDCORNERSALL         int = 0b0000_1111_0000
	FLAGDRAW_ROUNDCORNERSDEFAULT     int = 0b0000_1111_0000
	FLAGDRAW_ROUNDCORNERSMASK        int = 0b0001_1111_0000
)

const (
	FLAGFOCUS_NONE                int = 0b0000_0000
	FLAGFOCUS_CHILDWINDOWS        int = 0b0000_0001
	FLAGFOCUS_ROOTWINDOW          int = 0b0000_0010
	FLAGFOCUS_ANYWINDOW           int = 0b0000_0100
	FLAGFOCUS_NOPOPUPHIERARCHY    int = 0b0000_1000
	FLAGFOCUS_DOCKHIERARCHY       int = 0b0001_0000
	FLAGFOCUS_ROOTANDCHILDWINDOWS int = 0b0000_0011
)

const (
	FLAGHOVERED_NONE                         int = 0b0000_0000_0000_0000_0000
	FLAGHOVERED_CHILDWINDOWS                 int = 0b0000_0000_0000_0000_0001
	FLAGHOVERED_ROOTWINDOW                   int = 0b0000_0000_0000_0000_0010
	FLAGHOVERED_ANYWINDOW                    int = 0b0000_0000_0000_0000_0100
	FLAGHOVERED_NOPOPUPHIERARCHY             int = 0b0000_0000_0000_0000_1000
	FLAGHOVERED_DOCKHIERARCHY                int = 0b0000_0000_0000_0001_0000
	FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP      int = 0b0000_0000_0000_0010_0000
	FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM int = 0b0000_0000_0000_1000_0000
	FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM    int = 0b0000_0000_0001_0000_0000
	FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW  int = 0b0000_0000_0010_0000_0000
	FLAGHOVERED_ALLOWWHENDISABLED            int = 0b0000_0000_0100_0000_0000
	FLAGHOVERED_NONAVOVERRIDE                int = 0b0000_0000_1000_0000_0000
	FLAGHOVERED_ALLOWWHENOVERLAPPED          int = 0b0000_0000_0011_0000_0000
	FLAGHOVERED_RECTONLY                     int = 0b0000_0000_0011_1010_0000
	FLAGHOVERED_ROOTANDCHILDWINDOWS          int = 0b0000_0000_0000_0000_0011
	FLAGHOVERED_FORTOOLTIP                   int = 0b0000_0001_0000_0000_0000
	FLAGHOVERED_STATIONARY                   int = 0b0000_0010_0000_0000_0000
	FLAGHOVERED_DELAYNONE                    int = 0b0000_0100_0000_0000_0000
	FLAGHOVERED_DELAYSHORT                   int = 0b0000_1000_0000_0000_0000
	FLAGHOVERED_DELAYNORMAL                  int = 0b0001_0000_0000_0000_0000
	FLAGHOVERED_NOSHAREDDELAY                int = 0b0010_0000_0000_0000_0000
)

const (
	MOUSECURSOR_NONE int = iota - 1
	MOUSECURSOR_ARROW
	MOUSECURSOR_TEXTINPUT
	MOUSECURSOR_RESIZEALL
	MOUSECURSOR_RESIZENS
	MOUSECURSOR_RESIZEEW
	MOUSECURSOR_RESIZENESW
	MOUSECURSOR_RESIZENWSE
	MOUSECURSOR_HAND
	MOUSECURSOR_NOTALLOWED
	MOUSECURSOR_COUNT
)

const (
	ACTION_RELEASE int = iota
	ACTION_PRESS
	ACTION_REPEAT
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
	WIDGET_TREE_NODE            = "tree_node"
	WIDGET_TREE_TABLE_ROW       = "tree_table_row"
	WIDGET_TREE_TABLE           = "tree_table"
	WIDGET_WINDOW_SINGLE        = "window_single"
	WIDGET_POPUP_MODAL          = "popup_modal"
	WIDGET_POPUP                = "popup"
	WIDGET_LAYOUT_SPLIT         = "layout_split"
	WIDGET_SPLITTER             = "splitter"
	WIDGET_STACK                = "stack"
	WIDGET_ALIGN                = "align"
	WIDGET_MSG_BOX              = "msg_box"
	WIDGET_MSG_BOX_PREPARE      = "msg_box_prepare"
	WIDGET_BUTTON_INVISIBLE     = "button_invisible"
	WIDGET_BUTTON_IMAGE         = "button_image"
	WIDGET_STYLE                = "style"
	WIDGET_CUSTOM               = "custom"
	WIDGET_EVENT_HANDLER        = "event_handler"
	WIDGET_PLOT                 = "plot"
	WIDGET_CSS_TAG              = "css_tag"
)

const (
	PLOT_BAR_H      = "plot_bar_h"
	PLOT_BAR        = "plot_bar"
	PLOT_LINE       = "plot_line"
	PLOT_LINE_XY    = "plot_line_xy"
	PLOT_PIE_CHART  = "plot_pie_chart"
	PLOT_SCATTER    = "plot_scatter"
	PLOT_SCATTER_XY = "plot_scatter_xy"
	PLOT_CUSTOM     = "plot_custom"
)

var (
	buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{}
	plotList  = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget{}
)

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
		WIDGET_TREE_NODE:            treeNodeBuild,
		WIDGET_TREE_TABLE:           treeTableBuild,
		WIDGET_POPUP_MODAL:          popupModalBuild,
		WIDGET_POPUP:                popupBuild,
		WIDGET_LAYOUT_SPLIT:         splitLayoutBuild,
		WIDGET_SPLITTER:             splitterBuild,
		WIDGET_STACK:                stackBuild,
		WIDGET_ALIGN:                alignBuild,
		WIDGET_MSG_BOX_PREPARE:      msgBoxPrepareBuild,
		WIDGET_BUTTON_INVISIBLE:     buttonInvisibleBuild,
		WIDGET_BUTTON_IMAGE:         buttonImageBuild,
		WIDGET_STYLE:                styleBuild,
		WIDGET_CUSTOM:               customBuild,
		WIDGET_EVENT_HANDLER:        eventHandlerBuild,
		WIDGET_PLOT:                 plotBuild,
		WIDGET_CSS_TAG:              cssTagBuild,
	}

	plotList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget{
		PLOT_BAR_H:      plotBarHBuild,
		PLOT_BAR:        plotBarBuild,
		PLOT_LINE:       plotLineBuild,
		PLOT_LINE_XY:    plotLineXYBuild,
		PLOT_PIE_CHART:  plotPieBuild,
		PLOT_SCATTER:    plotScatterBuild,
		PLOT_SCATTER_XY: plotScatterXYBuild,
		PLOT_CUSTOM:     plotCustomBuild,
	}
}

func parseWidgets(widgetTable map[string]any, state *golua.LState, lg *log.Logger) []*golua.LTable {
	/// @interface Widget
	/// @prop type {string<gui.WidgetType>}

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

func parseTable(t *golua.LTable) map[string]any {
	v := map[string]any{}

	ln := t.Len()
	for i := range ln {
		v[string(i+1)] = t.RawGetInt(i + 1)
	}

	return v
}

func layoutBuild(r *lua.Runner, state *golua.LState, widgets []*golua.LTable, lg *log.Logger) []g.Widget {
	w := []g.Widget{}

	for _, wt := range widgets {
		t := string(wt.RawGetString("type").(golua.LString))

		build, ok := buildList[t]
		if !ok {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unknown widget: %s", t), log.LEVEL_ERROR)), 0)
		}

		w = append(w, build(r, lg, state, wt))
	}

	return w
}

func labelTable(state *golua.LState, text string) *golua.LTable {
	/// @struct WidgetLabel
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method wrapped(self, bool) -> self
	/// @method font(self, int<ref.FONT>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_LABEL))
	t.RawSetString("label", golua.LString(text))
	t.RawSetString("__wrapped", golua.LNil)
	t.RawSetString("__font", golua.LNil)

	tableBuilderFunc(state, t, "wrapped", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		t.RawSetString("__wrapped", golua.LBool(v))
	})

	tableBuilderFunc(state, t, "font", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		t.RawSetString("__font", v)
	})

	return t
}

func labelBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	l := g.Label(t.RawGetString("label").String())

	wrapped := t.RawGetString("__wrapped")
	if wrapped.Type() == golua.LTBool {
		l.Wrapped(bool(wrapped.(golua.LBool)))
	}

	fontref := t.RawGetString("__font")
	if fontref.Type() == golua.LTNumber {
		ref := int(fontref.(golua.LNumber))
		sref, err := r.CR_REF.Item(ref)
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		font := sref.Value.(*g.FontInfo)

		l.Font(font)
	}

	return l
}

func buttonTable(state *golua.LState, text string) *golua.LTable {
	/// @struct WidgetButton
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method disabled(self, bool) -> self
	/// @method size(self, width float, height float) -> self
	/// @method on_click(self, callback {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON))
	t.RawSetString("label", golua.LString(text))
	t.RawSetString("__disabled", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__click", golua.LNil)

	tableBuilderFunc(state, t, "disabled", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		t.RawSetString("__disabled", golua.LBool(v))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	return t
}

func buttonBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Button(t.RawGetString("label").String())

	disabled := t.RawGetString("__disabled")
	if disabled.Type() == golua.LTBool {
		b.Disabled(bool(disabled.(golua.LBool)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func dummyTable(state *golua.LState, width, height float64) *golua.LTable {
	/// @struct WidgetDummy
	/// @prop type {string<gui.WidgetType>}
	/// @prop width {float}
	/// @prop height {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_DUMMY))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))

	return t
}

func dummyBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	w := t.RawGetString("width").(golua.LNumber)
	h := t.RawGetString("height").(golua.LNumber)
	d := g.Dummy(float32(w), float32(h))

	return d
}

func separatorTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetSeparator
	/// @prop type {string<gui.WidgetType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SEPARATOR))

	return t
}

func separatorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	s := g.Separator()

	return s
}

func bulletTextTable(state *golua.LState, text string) *golua.LTable {
	/// @struct WidgetBulletText
	/// @prop type {string<gui.WidgetType>}
	/// @prop text {string}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BULLET_TEXT))
	t.RawSetString("text", golua.LString(text))

	return t
}

func bulletTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.BulletText(t.RawGetString("text").String())

	return b
}

func bulletTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetBullet
	/// @prop type {string<gui.WidgetType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BULLET))

	return t
}

func bulletBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Bullet()

	return b
}

func checkboxTable(state *golua.LState, text string, boolref int) *golua.LTable {
	/// @struct WidgetCheckbox
	/// @prop type {string<gui.WidgetType>}
	/// @prop text {string}
	/// @prop boolref {int<ref.BOOL>}
	/// @method on_change(self, {function(bool, int<ref.BOOL>)}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CHECKBOX))
	t.RawSetString("text", golua.LString(text))
	t.RawSetString("boolref", golua.LNumber(boolref))
	t.RawSetString("__change", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	return t
}

func checkboxBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("boolref").(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	selected := sref.Value.(*bool)
	c := g.Checkbox(t.RawGetString("text").String(), selected)

	change := t.RawGetString("__change")
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
	/// @struct WidgetChild
	/// @prop type {string<gui.WidgetType>}
	/// @method border(self, bool) -> self
	/// @method size(self, width float, height float) -> self
	/// @method layout(self, widgets []struct<gui.Widget>) -> self
	/// @method flags(self, flags int<gui.WindowFlags>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CHILD))
	t.RawSetString("__border", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckBool(-1)
		t.RawSetString("__border", golua.LBool(b))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func childBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.Child()

	border := t.RawGetString("__border")
	if border.Type() == golua.LTBool {
		c.Border(bool(border.(golua.LBool)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func colorEditTable(state *golua.LState, text string, colorref int) *golua.LTable {
	/// @struct WidgetColorEdit
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @prop colorref {int<ref.RGBA>}
	/// @method size(self, width float) -> self
	/// @method on_change(self, callback {function(color struct<image.Color>, int<ref.RGBA>)}) -> self
	/// @method flags(self, flags int<gui.ColorEditFlags>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_COLOR_EDIT))
	t.RawSetString("label", golua.LString(text))
	t.RawSetString("colorref", golua.LNumber(colorref))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func colorEditBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("colorref").(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	selected := sref.Value.(*color.RGBA)
	c := g.ColorEdit(t.RawGetString("label").String(), selected)

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			ct := imageutil.RGBAColorToColorTable(state, selected)

			state.Push(change)
			state.Push(ct)
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ColorEditFlags(flags.(golua.LNumber)))
	}

	return c
}

func columnTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	/// @struct WidgetColumn
	/// @prop type {string<gui.WidgetType>}
	/// @prop widgets {[]struct<gui.Widget>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_COLUMN))
	t.RawSetString("widgets", widgets)

	return t
}

func columnBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	var widgets []g.Widget

	wid := t.RawGetString("widgets")
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable)), state, lg), lg)
	}

	s := g.Column(widgets...)

	return s
}

func rowTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	/// @struct WidgetRow
	/// @prop type {string<gui.WidgetType>}
	/// @prop widgets {[]struct<gui.Widget>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_ROW))
	t.RawSetString("widgets", widgets)

	return t
}

func rowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	var widgets []g.Widget

	wid := t.RawGetString("widgets")
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable)), state, lg), lg)
	}

	s := g.Row(widgets...)

	return s
}

func comboCustomTable(state *golua.LState, text, preview string) *golua.LTable {
	/// @struct WidgetComboCustom
	/// @prop type {string<gui.WidgetType>}
	/// @prop text {string}
	/// @prop preview {string}
	/// @method size(self, width float) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self
	/// @method flags(self, flags int<gui.ComboFlags>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_COMBO_CUSTOM))
	t.RawSetString("text", golua.LString(text))
	t.RawSetString("preview", golua.LString(preview))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func comboCustomBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := t.RawGetString("text").String()
	preview := t.RawGetString("preview").String()
	c := g.ComboCustom(text, preview)

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ComboFlags(flags.(golua.LNumber)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func comboTable(state *golua.LState, text, preview string, items golua.LValue, i32Ref int) *golua.LTable {
	/// @struct WidgetCombo
	/// @prop type {string<gui.WidgetType>}
	/// @prop text {string}
	/// @prop preview {string}
	/// @prop items {[]string}
	/// @prop i32ref {int<ref.INT32>}
	/// @method size(self, width float) -> self
	/// @method on_change(self, {function(int, int<ref.INT32>)}) -> self
	/// @method flags(self, flags int<gui.ComboFlags) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_COMBO))
	t.RawSetString("text", golua.LString(text))
	t.RawSetString("preview", golua.LString(preview))
	t.RawSetString("items", items)
	t.RawSetString("i32ref", golua.LNumber(i32Ref))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func comboBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := t.RawGetString("text").String()
	preview := t.RawGetString("preview").String()

	items := []string{}
	it := t.RawGetString("items").(*golua.LTable)
	for i := range it.Len() {
		v := it.RawGetInt(i + 1).(golua.LString)
		items = append(items, string(v))
	}

	ref := int(t.RawGetString("i32ref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.Combo(text, preview, items, selected)

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.ComboFlags(flags.(golua.LNumber)))
	}

	change := t.RawGetString("__change")
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
	/// @struct WidgetCondition
	/// @prop type {string<gui.WidgetType>}
	/// @prop condition {bool}
	/// @prop layoutIf {struct<gui.Widget>}
	/// @prop layoutElse {struct<gui.Widget>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CONDITION))
	t.RawSetString("condition", golua.LBool(condition))
	t.RawSetString("layoutIf", layoutIf)
	t.RawSetString("layoutElse", layoutElse)

	return t
}

func conditionBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	condition := t.RawGetString("condition").(golua.LBool)

	widIf := t.RawGetString("layoutIf").(*golua.LTable)
	widgetsIf := layoutBuild(r, state, []*golua.LTable{widIf}, lg)
	widElse := t.RawGetString("layoutElse").(*golua.LTable)
	widgetsElse := layoutBuild(r, state, []*golua.LTable{widElse}, lg)

	s := g.Condition(bool(condition), widgetsIf[0], widgetsElse[0])

	return s
}

func contextMenuTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetContextMenu
	/// @prop type {string<gui.WidgetType>}
	/// @method mouse_button(self, button int<gui.MouseButton>) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CONTEXT_MENU))
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__button", golua.LNil)

	tableBuilderFunc(state, t, "mouse_button", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckNumber(-1)
		t.RawSetString("__button", lt)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func contextMenuBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.ContextMenu()

	button := t.RawGetString("__button")
	if button.Type() == golua.LTNumber {
		c.MouseButton(g.MouseButton(button.(golua.LNumber)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func datePickerTable(state *golua.LState, id string, timeref int) *golua.LTable {
	/// @struct WidgetDatePicker
	/// @prop type {string<gui.WidgetType>}
	/// @prop id {string}
	/// @prop timeref {int<ref.TIME>}
	/// @method on_change(self, {function(string, int<ref.TIME>)}) -> self
	/// @method format(self, format string) -> self
	/// @method size(self, width float) -> self
	/// @method start_of_week(self, day int<time.Weekday>) -> self
	/// @method translation(self, label string<gui.DatePickerLabel>, value string) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_DATE_PICKER))
	t.RawSetString("id", golua.LString(id))
	t.RawSetString("timeref", golua.LNumber(timeref))
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__startofweek", golua.LNil)
	t.RawSetString("__translationMonth", golua.LNil)
	t.RawSetString("__translationYear", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "start_of_week", func(state *golua.LState, t *golua.LTable) {
		sow := state.CheckNumber(-1)
		t.RawSetString("__startofweek", sow)
	})

	tableBuilderFunc(state, t, "translation", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-2)
		value := state.CheckString(-1)

		if label == string(DATEPICKERLABEL_MONTH) {
			t.RawSetString("__translationMonth", golua.LString(value))
		} else if label == string(DATEPICKERLABEL_YEAR) {
			t.RawSetString("__translationYear", golua.LString(value))
		}
	})

	return t
}

func datePickerBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("timeref").(golua.LNumber))

	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}

	date := sref.Value.(*time.Time)
	c := g.DatePicker(t.RawGetString("id").String(), date)

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(date.UnixMilli()))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	sow := t.RawGetString("__startofweek")
	if sow.Type() == golua.LTNumber {
		c.StartOfWeek(time.Weekday((sow.(golua.LNumber))))
	}

	translationMonth := t.RawGetString("__translationMonth")
	if translationMonth.Type() == golua.LTString {
		c.Translation(DATEPICKERLABEL_MONTH, translationMonth.String())
	}
	translationYear := t.RawGetString("__translationYear")
	if translationYear.Type() == golua.LTString {
		c.Translation(DATEPICKERLABEL_YEAR, translationYear.String())
	}

	return c
}

func dragIntTable(state *golua.LState, text string, i32Ref, minValue, maxValue int) *golua.LTable {
	/// @struct WidgetDragInt
	/// @prop type {string<gui.WidgetType>}
	/// @prop text {string}
	/// @prop i32ref {int<ref.INT32>}
	/// @prop minvalue {int}
	/// @prop maxvalue {int}
	/// @method speed(self, speed float) -> self
	/// @method format(self, format string) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_DRAG_INT))
	t.RawSetString("text", golua.LString(text))
	t.RawSetString("i32ref", golua.LNumber(i32Ref))
	t.RawSetString("minvalue", golua.LNumber(minValue))
	t.RawSetString("maxvalue", golua.LNumber(maxValue))
	t.RawSetString("__speed", golua.LNil)
	t.RawSetString("__format", golua.LNil)

	tableBuilderFunc(state, t, "speed", func(state *golua.LState, t *golua.LTable) {
		speed := state.CheckNumber(-1)
		t.RawSetString("__speed", speed)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	return t
}

func dragIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := t.RawGetString("text").String()
	min := t.RawGetString("minvalue").(golua.LNumber)
	max := t.RawGetString("maxvalue").(golua.LNumber)

	ref := int(t.RawGetString("i32ref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.DragInt(text, selected, int32(min), int32(max))

	speed := t.RawGetString("__speed")
	if speed.Type() == golua.LTNumber {
		c.Speed(float32(speed.(golua.LNumber)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	return c
}

func inputFloatTable(state *golua.LState, floatref int) *golua.LTable {
	/// @struct WidgetInputFloat
	/// @prop type {string<gui.WidgetType>}
	/// @prop f32ref {int<ref.FLOAT32>}
	/// @method size(self, width float) -> self
	/// @method on_change(self, {function(float, int<ref.FLOAT32>)}) -> self
	/// @method format(self, format string) -> self
	/// @method flags(self, flags int<gui.InputFlags>) -> self
	/// @method label(self, label string) -> self
	/// @method step_size(self, stepsize float) -> self
	/// @method step_size_fast(self, stepsize float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_INPUT_FLOAT))
	t.RawSetString("f32ref", golua.LNumber(floatref))
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__stepsize", golua.LNil)
	t.RawSetString("__stepsizefast", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "step_size", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		t.RawSetString("__stepsize", stepsize)
	})

	tableBuilderFunc(state, t, "step_size_fast", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		t.RawSetString("__stepsizefast", stepsize)
	})

	return t
}

func inputFloatBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("f32ref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*float32)

	c := g.InputFloat(selected)

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		c.Format(string(format.(golua.LString)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	stepsize := t.RawGetString("__stepsize")
	if stepsize.Type() == golua.LTNumber {
		c.StepSize(float32(stepsize.(golua.LNumber)))
	}

	stepsizefast := t.RawGetString("__stepsizefast")
	if stepsizefast.Type() == golua.LTNumber {
		c.StepSizeFast(float32(stepsizefast.(golua.LNumber)))
	}

	return c
}

func inputIntTable(state *golua.LState, intref int) *golua.LTable {
	/// @struct WidgetInputInt
	/// @prop type {string<gui.WidgetType>}
	/// @prop i32ref {int<ref.INT32>}
	/// @method size(self, width float) -> self
	/// @method on_change({function(int, int<ref.INT32>)}) -> self
	/// @method flags(self, flags int<gui.InputFlags>) -> self
	/// @method label(self, label string) -> self
	/// @method step_size(self, stepsize int) -> self
	/// @method step_size_fast(self, stepsize int) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_INPUT_INT))
	t.RawSetString("i32ref", golua.LNumber(intref))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__stepsize", golua.LNil)
	t.RawSetString("__stepsizefast", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "step_size", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		t.RawSetString("__stepsize", stepsize)
	})

	tableBuilderFunc(state, t, "step_size_fast", func(state *golua.LState, t *golua.LTable) {
		stepsize := state.CheckNumber(-1)
		t.RawSetString("__stepsizefast", stepsize)
	})

	return t
}

func inputIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("i32ref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*int32)

	c := g.InputInt(selected)

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	stepsize := t.RawGetString("__stepsize")
	if stepsize.Type() == golua.LTNumber {
		c.StepSize(int(stepsize.(golua.LNumber)))
	}

	stepsizefast := t.RawGetString("__stepsizefast")
	if stepsizefast.Type() == golua.LTNumber {
		c.StepSizeFast(int(stepsizefast.(golua.LNumber)))
	}

	return c
}

func inputTextTable(state *golua.LState, strref int) *golua.LTable {
	/// @struct WidgetInputText
	/// @prop type {string<gui.WidgetType>}
	/// @prop strref {int<ref.STRING>}
	/// @method size(self, width float) -> self
	/// @method flags(self, flags int<gui.InputFlags>) -> self
	/// @method label(self, label strings) -> self
	/// @method autocomplete(self, []string) -> self
	/// @method callback(self, {function(string, int<ref.STRING>)}) -> self
	/// @method hint(self, hint string) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_INPUT_TEXT))
	t.RawSetString("strref", golua.LNumber(strref))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__autocomplete", golua.LNil)
	t.RawSetString("__callback", golua.LNil)
	t.RawSetString("__hint", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "autocomplete", func(state *golua.LState, t *golua.LTable) {
		ac := state.CheckTable(-1)
		t.RawSetString("__autocomplete", ac)
	})

	tableBuilderFunc(state, t, "callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__callback", fn)
	})

	tableBuilderFunc(state, t, "hint", func(state *golua.LState, t *golua.LTable) {
		hint := state.CheckString(-1)
		t.RawSetString("__hint", golua.LString(hint))
	})

	return t
}

func inputTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("strref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*string)

	c := g.InputText(selected)

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	hint := t.RawGetString("__hint")
	if hint.Type() == golua.LTString {
		c.Hint(string(hint.(golua.LString)))
	}

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	callback := t.RawGetString("__callback")
	if callback.Type() == golua.LTFunction {
		c.Callback(func(data imgui.InputTextCallbackData) int {
			state.Push(callback)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
			return 0
		})
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)))
	}

	ac := t.RawGetString("__autocomplete")
	if ac.Type() == golua.LTTable {
		acList := []string{}
		at := ac.(*golua.LTable)
		for i := range at.Len() {
			ai := at.RawGetInt(i + 1).(golua.LString)
			acList = append(acList, string(ai))
		}

		c.AutoComplete(acList)
	}

	return c
}

func inputMultilineTextTable(state *golua.LState, strref int) *golua.LTable {
	/// @struct WidgetInputTextMultiline
	/// @prop type {string<gui.WidgetType>}
	/// @prop strref {int<ref.STRING>}
	/// @method size(self, width float, height float) -> self
	/// @method on_change(self, {function(string, int<ref.STRING>)}) -> self
	/// @method flags(self, flags int<gui.InputFlags>) -> self
	/// @method label(self, label string) -> self
	/// @method callback(self, {function(string, int<ref.STRING>)}) -> self
	/// @method autoscroll_to_bottom(self, bool) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_INPUT_MULTILINE_TEXT))
	t.RawSetString("strref", golua.LNumber(strref))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__callback", golua.LNil)
	t.RawSetString("__autoscroll", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__callback", fn)
	})

	tableBuilderFunc(state, t, "autoscroll_to_bottom", func(state *golua.LState, t *golua.LTable) {
		as := state.CheckBool(-1)
		t.RawSetString("__autoscroll", golua.LBool(as))
	})

	return t
}

func inputMultilineTextBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ref := int(t.RawGetString("strref").(golua.LNumber))
	sref, err := r.CR_REF.Item(ref)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	selected := sref.Value.(*string)

	c := g.InputTextMultiline(selected)

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.InputTextFlags(flags.(golua.LNumber)))
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		c.Label(string(label.(golua.LString)))
	}

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		c.OnChange(func() {
			state.Push(change)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
		})
	}

	callback := t.RawGetString("__callback")
	if callback.Type() == golua.LTFunction {
		c.Callback(func(data imgui.InputTextCallbackData) int {
			state.Push(callback)
			state.Push(golua.LString(*selected))
			state.Push(golua.LNumber(ref))
			state.Call(2, 0)
			return 0
		})
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		c.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	autoscroll := t.RawGetString("__autoscroll")
	if autoscroll.Type() == golua.LTBool {
		c.AutoScrollToBottom(bool(autoscroll.(golua.LBool)))
	}

	return c
}

func progressBarTable(state *golua.LState, fraction float64) *golua.LTable {
	/// @struct WidgetProgressBar
	/// @prop type {string<gui.WidgetType>}
	/// @prop fraction {float}
	/// @method overlay(self, label string) -> self
	/// @method size(self, width float, height float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_PROGRESS_BAR))
	t.RawSetString("fraction", golua.LNumber(fraction))
	t.RawSetString("__overlay", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "overlay", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__overlay", golua.LString(label))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	return t
}

func progressBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	fraction := t.RawGetString("fraction").(golua.LNumber)
	p := g.ProgressBar(float32(fraction))

	overlay := t.RawGetString("__overlay")
	if overlay.Type() == golua.LTString {
		p.Overlay(string(overlay.(golua.LString)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		p.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return p
}

func progressIndicatorTable(state *golua.LState, label string, width, height, radius float64) *golua.LTable {
	/// @struct WidgetProgressIndicator
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @prop width {float}
	/// @prop height {float}
	/// @prop radius {float}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_PROGRESS_INDICATOR))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("width", golua.LNumber(width))
	t.RawSetString("height", golua.LNumber(height))
	t.RawSetString("radius", golua.LNumber(radius))

	return t
}

func progressIndicatorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	width := t.RawGetString("width").(golua.LNumber)
	height := t.RawGetString("height").(golua.LNumber)
	radius := t.RawGetString("radius").(golua.LNumber)
	p := g.ProgressIndicator(string(label), float32(width), float32(height), float32(radius))

	return p
}

func spacingTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetSpacing
	/// @prop type {string<gui.WidgetType>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SPACING))

	return t
}

func spacingBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Spacing()

	return b
}

func buttonSmallTable(state *golua.LState, text string) *golua.LTable {
	/// @struct WidgetButtonSmall
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method on_click(self, callback {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON_SMALL))
	t.RawSetString("label", golua.LString(text))
	t.RawSetString("__click", golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	return t
}

func buttonSmallBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := t.RawGetString("label").(golua.LString)
	b := g.SmallButton(string(text))

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func buttonRadioTable(state *golua.LState, text string, active bool) *golua.LTable {
	/// @struct WidgetButtonRadio
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @prop active {bool}
	/// @method on_change(self, {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON_RADIO))
	t.RawSetString("label", golua.LString(text))
	t.RawSetString("active", golua.LBool(active))
	t.RawSetString("__change", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	return t
}

func buttonRadioBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	text := t.RawGetString("label").(golua.LString)
	active := t.RawGetString("active").(golua.LBool)
	b := g.RadioButton(string(text), bool(active))

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Call(0, 0)
		})
	}

	return b
}

func imageUrlTable(state *golua.LState, url string) *golua.LTable {
	/// @struct WidgetImageURL
	/// @prop type {string<gui.WidgetType>}
	/// @prop url {string}
	/// @method on_click(self, {function()}) -> self
	/// @method size(self, width float, height float) -> self
	/// @method timeout(self, timeout int) -> self
	/// @method layout_for_failure(self, []struct<gui.Widget>) -> self
	/// @method layout_for_loading(self, []struct<gui.Widget>) -> self
	/// @method on_failure(self, {function()}) -> self
	/// @method on_ready(self, {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_IMAGE_URL))
	t.RawSetString("url", golua.LString(url))
	t.RawSetString("__click", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__timeout", golua.LNil)
	t.RawSetString("__failwidgets", golua.LNil)
	t.RawSetString("__loadwidgets", golua.LNil)
	t.RawSetString("__failure", golua.LNil)
	t.RawSetString("__ready", golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "timeout", func(state *golua.LState, t *golua.LTable) {
		timeout := state.CheckNumber(-1)
		t.RawSetString("__timeout", timeout)
	})

	tableBuilderFunc(state, t, "layout_for_failure", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__failwidgets", lt)
	})

	tableBuilderFunc(state, t, "layout_for_loading", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__loadwidgets", lt)
	})

	tableBuilderFunc(state, t, "on_failure", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__failure", fn)
	})

	tableBuilderFunc(state, t, "on_ready", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__ready", fn)
	})

	return t
}

func imageUrlBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	url := t.RawGetString("url").(golua.LString)
	i := g.ImageWithURL(string(url))

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		i.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	timeout := t.RawGetString("__timeout")
	if timeout.Type() == golua.LTNumber {
		i.Timeout(time.Duration(timeout.(golua.LNumber)))
	}

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		i.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	lfail := t.RawGetString("__failwidgets")
	if lfail.Type() == golua.LTTable {
		i.LayoutForFailure(layoutBuild(r, state, parseWidgets(parseTable(lfail.(*golua.LTable)), state, lg), lg)...)
	}

	lload := t.RawGetString("__loadwidgets")
	if lload.Type() == golua.LTTable {
		i.LayoutForFailure(layoutBuild(r, state, parseWidgets(parseTable(lload.(*golua.LTable)), state, lg), lg)...)
	}

	failure := t.RawGetString("__failure")
	if failure.Type() == golua.LTFunction {
		i.OnFailure(func(err error) {
			lg.Append(fmt.Sprintf("error occured while loading image url: %s", err), log.LEVEL_WARN)
			state.Push(failure)
			state.Call(0, 0)
		})
	}

	ready := t.RawGetString("__ready")
	if ready.Type() == golua.LTFunction {
		i.OnReady(func() {
			state.Push(ready)
			state.Call(0, 0)
		})
	}

	return i
}

func imageTable(state *golua.LState, image int, sync, cache bool) *golua.LTable {
	/// @struct WidgetImage
	/// @prop type {string<gui.WidgetType>}
	/// @prop image {int<collection.IMAGE>}
	/// @prop imagecached {int<collection.CRATE_CACHEDIMAGE>}
	/// @prop sync {bool}
	/// @prop cached {bool}
	/// @method on_click(self, {function()}) -> self
	/// @method size(self, width float, height float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_IMAGE))
	t.RawSetString("image", golua.LNumber(image))
	t.RawSetString("imagecached", golua.LNumber(image))
	t.RawSetString("sync", golua.LBool(sync))
	t.RawSetString("cached", golua.LBool(cache))
	t.RawSetString("__click", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	return t
}

func imageBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ig := t.RawGetString("image").(golua.LNumber)
	var img image.Image

	sync := t.RawGetString("sync").(golua.LBool)
	cache := t.RawGetString("cached").(golua.LBool)

	if !sync {
		if cache {
			ci, err := r.CR_CIM.Item(int(ig))
			if err == nil {
				img = ci.Image
			}
		} else {
			<-r.IC.Schedule(state, int(ig), &collection.Task[collection.ItemImage]{
				Lib:  LIB_GUI,
				Name: "wg_image",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})
		}
	} else {
		item := r.IC.Item(int(ig))
		if item != nil && item.Self != nil {
			img = item.Self.Image
		}
	}

	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	i := g.ImageWithRgba(img)

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		i.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		i.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return i
}

func listBoxTable(state *golua.LState, items golua.LValue) *golua.LTable {
	/// @struct WidgetListBox
	/// @prop type {string<gui.WidgetType>}
	/// @prop items {[]string}
	/// @method on_change(self, {function(int)}) -> self
	/// @method border(self, bool) -> self
	/// @method context_menu(self, []struct<gui.Widget>) -> self
	/// @method on_double_click(self, {function(int)}) -> self
	/// @method on_menu(self, {function(int, string)}) -> self
	/// @method selected_index(self, int) -> self
	/// @method size(self, width float, height float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_LIST_BOX))
	t.RawSetString("items", items)
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__border", golua.LNil)
	t.RawSetString("__context", golua.LNil)
	t.RawSetString("__dclick", golua.LNil)
	t.RawSetString("__menu", golua.LNil)
	t.RawSetString("__sel", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckBool(-1)
		t.RawSetString("__border", golua.LBool(b))
	})

	tableBuilderFunc(state, t, "context_menu", func(state *golua.LState, t *golua.LTable) {
		cmt := state.CheckTable(-1)
		t.RawSetString("__context", cmt)
	})

	tableBuilderFunc(state, t, "on_double_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__dclick", fn)
	})

	tableBuilderFunc(state, t, "on_menu", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__menu", fn)
	})

	tableBuilderFunc(state, t, "selected_index", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckNumber(-1)
		t.RawSetString("__sel", sel)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	return t
}

func listBoxBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	it := t.RawGetString("items").(*golua.LTable)
	items := []string{}
	for i := range it.Len() {
		is := it.RawGetInt(i + 1).(golua.LString)
		items = append(items, string(is))
	}

	b := g.ListBox(items)

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		b.OnChange(func(index int) {
			state.Push(change)
			state.Push(golua.LNumber(index))
			state.Call(1, 0)
		})
	}

	dclick := t.RawGetString("__dclick")
	if dclick.Type() == golua.LTFunction {
		b.OnDClick(func(index int) {
			state.Push(dclick)
			state.Push(golua.LNumber(index))
			state.Call(1, 0)
		})
	}

	selmenu := t.RawGetString("__menu")
	if selmenu.Type() == golua.LTFunction {
		b.OnMenu(func(index int, menu string) {
			state.Push(selmenu)
			state.Push(golua.LNumber(index))
			state.Push(golua.LString(menu))
			state.Call(2, 0)
		})
	}

	sel := t.RawGetString("__sel")
	if sel.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(sel.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		b.SelectedIndex(ref.Value.(*int32))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return b
}

func listClipperTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetListClipper
	/// @prop type {string<gui.WidgetType>}
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_LIST_CLIPPER))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func listClipperBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.ListClipper()

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func mainMenuBarTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetMainMenuBar
	/// @prop type {string<gui.WidgetType>}
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MENU_BAR_MAIN))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func mainMenuBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.MainMenuBar()

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func menuBarTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetMenuBar
	/// @prop type {string<gui.WidgetType>}
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MENU_BAR))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func menuBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c := g.MenuBar()

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func menuItemTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetMenuItem
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method enabled(self, bool) -> self
	/// @method on_click(self, {function()}) -> self
	/// @method selected(self, bool) -> self
	/// @method shortcut(self, string) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MENU_ITEM))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__enabled", golua.LNil)
	t.RawSetString("__click", golua.LNil)
	t.RawSetString("__sel", golua.LNil)
	t.RawSetString("__shortcut", golua.LNil)

	tableBuilderFunc(state, t, "enabled", func(state *golua.LState, t *golua.LTable) {
		en := state.CheckBool(-1)
		t.RawSetString("__enabled", golua.LBool(en))
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	tableBuilderFunc(state, t, "selected", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckBool(-1)
		t.RawSetString("__sel", golua.LBool(sel))
	})

	tableBuilderFunc(state, t, "shortcut", func(state *golua.LState, t *golua.LTable) {
		sc := state.CheckString(-1)
		t.RawSetString("__shortcut", golua.LString(sc))
	})

	return t
}

func menuItemBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	m := g.MenuItem(string(label))

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		m.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	enabled := t.RawGetString("__enabled")
	if enabled.Type() == golua.LTBool {
		m.Enabled(bool(enabled.(golua.LBool)))
	}

	sel := t.RawGetString("__sel")
	if sel.Type() == golua.LTBool {
		m.Selected(bool(sel.(golua.LBool)))
	}

	shortcut := t.RawGetString("__shortcut")
	if shortcut.Type() == golua.LTString {
		m.Shortcut(string(shortcut.(golua.LString)))
	}

	return m
}

func menuTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetMenu
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method enabled(self, bool) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MENU))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__enabled", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "enabled", func(state *golua.LState, t *golua.LTable) {
		en := state.CheckBool(-1)
		t.RawSetString("__enabled", golua.LBool(en))
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func menuBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	m := g.Menu(string(label))

	enabled := t.RawGetString("__enabled")
	if enabled.Type() == golua.LTBool {
		m.Enabled(bool(enabled.(golua.LBool)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return m
}

func selectableTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetSelectable
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method on_click(self, {function()}) -> self
	/// @method on_double_click(self, {function()}) -> self
	/// @method selected(self, bool) -> self
	/// @method size(self, width float, height float) -> self
	/// @method flags(self, flags int<gui.SelectableFlags>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SELECTABLE))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__click", golua.LNil)
	t.RawSetString("__dclick", golua.LNil)
	t.RawSetString("__sel", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	tableBuilderFunc(state, t, "on_double_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__dclick", fn)
	})

	tableBuilderFunc(state, t, "selected", func(state *golua.LState, t *golua.LTable) {
		sel := state.CheckBool(-1)
		t.RawSetString("__sel", golua.LBool(sel))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func selectableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	b := g.Selectable(string(label))

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	dclick := t.RawGetString("__dclick")
	if dclick.Type() == golua.LTFunction {
		b.OnDClick(func() {
			state.Push(dclick)
			state.Call(0, 0)
		})
	}

	sel := t.RawGetString("__sel")
	if sel.Type() == golua.LTBool {
		b.Selected(bool(sel.(golua.LBool)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		b.Flags(g.SelectableFlags(flags.(golua.LNumber)))
	}

	return b
}

func sliderFloatTable(state *golua.LState, f32ref int, min, max float64) *golua.LTable {
	/// @struct WidgetSliderFloat
	/// @prop type {string<gui.WidgetType>}
	/// @prop f32ref {int<ref.FLOAT32>}
	/// @prop min {float}
	/// @prop max {float}
	/// @method on_change(self, {function(float, int<ref.FLOAT32>)}) -> self
	/// @method label(self, string) -> self
	/// @method format(self, string) -> self
	/// @method size(self, width float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SLIDER_FLOAT))
	t.RawSetString("f32ref", golua.LNumber(f32ref))
	t.RawSetString("min", golua.LNumber(min))
	t.RawSetString("max", golua.LNumber(max))
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__width", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	return t
}

func sliderFloatBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	floatref := t.RawGetString("f32ref").(golua.LNumber)
	ref, err := r.CR_REF.Item(int(floatref))
	value := ref.Value.(*float32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := t.RawGetString("min").(golua.LNumber)
	max := t.RawGetString("max").(golua.LNumber)
	b := g.SliderFloat(value, float32(min), float32(max))

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(floatref)
			state.Call(2, 0)
		})
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)))
	}

	return b
}

func sliderIntTable(state *golua.LState, i32ref int, min, max int) *golua.LTable {
	/// @struct WidgetSliderInt
	/// @prop type {string<gui.WidgetType>}
	/// @prop i32ref {int<ref.INT32>}
	/// @prop min {int}
	/// @prop max {int}
	/// @method on_change(self, {function(int, int<ref.INT32>)}) -> self
	/// @method label(self, string) -> self
	/// @method format(self, string) -> self
	/// @method size(self, width float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SLIDER_INT))
	t.RawSetString("i32ref", golua.LNumber(i32ref))
	t.RawSetString("min", golua.LNumber(min))
	t.RawSetString("max", golua.LNumber(max))
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__width", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	return t
}

func sliderIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	intref := t.RawGetString("i32ref").(golua.LNumber)
	ref, err := r.CR_REF.Item(int(intref))
	value := ref.Value.(*int32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := t.RawGetString("min").(golua.LNumber)
	max := t.RawGetString("max").(golua.LNumber)
	b := g.SliderInt(value, int32(min), int32(max))

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(intref)
			state.Call(2, 0)
		})
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)))
	}

	return b
}

func vsliderIntTable(state *golua.LState, i32ref int, min, max int) *golua.LTable {
	/// @struct WidgetVSliderInt
	/// @prop type {string<gui.WidgetType>}
	/// @prop i32ref {int<ref.INT32>}
	/// @prop min {int}
	/// @prop max {int}
	/// @method on_change(self, {function(int, int<ref.INT32>)}) -> self
	/// @method label(self, string) -> self
	/// @method format(self, string) -> self
	/// @method size(self, width float, height float) -> self
	/// @method flags(self, flags int<gui.SliderFlags>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_VSLIDER_INT))
	t.RawSetString("i32ref", golua.LNumber(i32ref))
	t.RawSetString("min", golua.LNumber(min))
	t.RawSetString("max", golua.LNumber(max))
	t.RawSetString("__change", golua.LNil)
	t.RawSetString("__label", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__flags", golua.LNil)

	tableBuilderFunc(state, t, "on_change", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__change", fn)
	})

	tableBuilderFunc(state, t, "label", func(state *golua.LState, t *golua.LTable) {
		label := state.CheckString(-1)
		t.RawSetString("__label", golua.LString(label))
	})

	tableBuilderFunc(state, t, "format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	return t
}

func vsliderIntBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	intref := t.RawGetString("i32ref").(golua.LNumber)
	ref, err := r.CR_REF.Item(int(intref))
	value := ref.Value.(*int32)
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	min := t.RawGetString("min").(golua.LNumber)
	max := t.RawGetString("max").(golua.LNumber)
	b := g.VSliderInt(value, int32(min), int32(max))

	change := t.RawGetString("__change")
	if change.Type() == golua.LTFunction {
		b.OnChange(func() {
			state.Push(change)
			state.Push(golua.LNumber(*value))
			state.Push(intref)
			state.Call(2, 0)
		})
	}

	label := t.RawGetString("__label")
	if label.Type() == golua.LTString {
		b.Label(string(label.(golua.LString)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		b.Format(string(format.(golua.LString)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		b.Flags(g.SliderFlags(flags.(golua.LNumber)))
	}

	return b
}

func tabbarTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetTabBar
	/// @prop type {string<gui.WidgetType>}
	/// @method flags(self, flags int<gui.TabBarFlags>) -> self
	/// @method tab_items(self, []struct<gui.TabItem>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TAB_BAR))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "tab_items", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func tabbarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.TabBar()

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TabBarFlags(flags.(golua.LNumber)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(layout.(*golua.LTable)), state, lg)
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
	/// @struct TabItem
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.TabItemFlags>) -> self
	/// @method is_open(self, bool) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TAB_ITEM))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__open", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__open", flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func tabitemBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TabItemWidget {
	label := t.RawGetString("label").(golua.LString)
	i := g.TabItem(string(label))

	open := t.RawGetString("__open")
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		i.IsOpen(ref.Value.(*bool))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		i.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		i.Flags(g.TabItemFlags(flags.(golua.LNumber)))
	}

	return i
}

func tooltipTable(state *golua.LState, tip string) *golua.LTable {
	/// @struct WidgetTooltip
	/// @prop type {string<gui.WidgetType>}
	/// @prop tip {string}
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TOOLTIP))
	t.RawSetString("tip", golua.LString(tip))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func tooltipBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tip := t.RawGetString("tip").(golua.LString)
	i := g.Tooltip(string(tip))

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		i.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return i
}

func tableColumnTable(state *golua.LState, label string) *golua.LTable {
	/// @struct TableColumn
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.TableColumnFlags>) -> self
	/// @method inner_width_or_weight(self, width float) -> self
	/// @desc
	/// Only used in table widget columns.

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TABLE_COLUMN))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__width", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "inner_width_or_weight", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__width", flags)
	})

	return t
}

func tableColumnBuild(t *golua.LTable) *g.TableColumnWidget {
	label := t.RawGetString("label").(golua.LString)
	c := g.TableColumn(string(label))

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		c.Flags(g.TableColumnFlags(flags.(golua.LNumber)))
	}

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		c.InnerWidthOrWeight(float32(width.(golua.LNumber)))
	}

	return c
}

func tableRowTable(state *golua.LState, widgets golua.LValue) *golua.LTable {
	/// @struct TableRow
	/// @prop type {string<gui.WidgetType>}
	/// @prop widgets {[]struct<gui.Widget>}
	/// @method flags(self, flags int<gui.TableRowFlags>) -> self
	/// @method bg_color(self, color struct<image.Color>) -> self
	/// @method min_height(self, height float) -> self
	/// @desc
	/// Only used in table widget rows.

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TABLE_ROW))
	t.RawSetString("widgets", widgets)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__color", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "bg_color", func(state *golua.LState, t *golua.LTable) {
		clr := state.CheckTable(-1)
		t.RawSetString("__color", clr)
	})

	tableBuilderFunc(state, t, "min_height", func(state *golua.LState, t *golua.LTable) {
		height := state.CheckNumber(-1)
		t.RawSetString("__height", height)
	})

	return t
}

func tableRowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TableRowWidget {
	var widgets []g.Widget

	wid := t.RawGetString("widgets")
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable)), state, lg), lg)
	}

	s := g.TableRow(widgets...)

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		s.Flags(g.TableRowFlags(flags.(golua.LNumber)))
	}

	height := t.RawGetString("__height")
	if height.Type() == golua.LTNumber {
		s.MinHeight(float64(height.(golua.LNumber)))
	}

	clr := t.RawGetString("__color")
	if clr.Type() == golua.LTTable {
		rgba := imageutil.ColorTableToRGBAColor(clr.(*golua.LTable))
		s.BgColor(rgba)
	}

	return s
}

func tableTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetTable
	/// @prop type {string<gui.WidgetType>}
	/// @method flags(self, flags int<gui.TableFlags>) -> self
	/// @method fast_mode(self, bool) -> self
	/// @method size(self, width float, height float) -> self
	/// @method columns(self, []struct<gui.TableColumn>) -> self
	/// @method rows(self, []struct<gui.TableRow>) -> self
	/// @method inner_width(self, width float) -> self
	/// @method freeze(self, col int, row int) -> self - Can be called multiple times.

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TABLE))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__columns", golua.LNil)
	t.RawSetString("__rows", golua.LNil)
	t.RawSetString("__fast", golua.LNil)
	t.RawSetString("__freeze", state.NewTable())
	t.RawSetString("__innerwidth", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "fast_mode", func(state *golua.LState, t *golua.LTable) {
		fast := state.CheckBool(-1)
		t.RawSetString("__fast", golua.LBool(fast))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "columns", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__columns", lt)
	})

	tableBuilderFunc(state, t, "rows", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__rows", lt)
	})

	tableBuilderFunc(state, t, "inner_width", func(state *golua.LState, t *golua.LTable) {
		innerwidth := state.CheckNumber(-1)
		t.RawSetString("__innerwidth", innerwidth)
	})

	tableBuilderFunc(state, t, "freeze", func(state *golua.LState, t *golua.LTable) {
		col := state.CheckNumber(-2)
		row := state.CheckNumber(-1)
		pt := state.NewTable()
		pt.RawSetString("col", col)
		pt.RawSetString("row", row)

		ft := t.RawGetString("__freeze").(*golua.LTable)
		ft.Append(pt)
	})

	return t
}

func tableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.Table()

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TableFlags(flags.(golua.LNumber)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		tb.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	innerwidth := t.RawGetString("__innerwidth")
	if innerwidth.Type() == golua.LTNumber {
		tb.InnerWidth(float64(innerwidth.(golua.LNumber)))
	}

	fast := t.RawGetString("__fast")
	if fast.Type() == golua.LTBool {
		tb.FastMode(bool(fast.(golua.LBool)))
	}

	freeze := t.RawGetString("__freeze").(*golua.LTable)
	for i := range freeze.Len() {
		pt := freeze.RawGetInt(i + 1).(*golua.LTable)
		col := pt.RawGetString("col").(golua.LNumber)
		row := pt.RawGetString("row").(golua.LNumber)

		tb.Freeze(int(col), int(row))
	}

	columns := t.RawGetString("__columns")
	if columns.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(columns.(*golua.LTable)), state, lg)
		wdi := []*g.TableColumnWidget{}
		for _, w := range wd {
			i := tableColumnBuild(w)
			wdi = append(wdi, i)
		}
		tb.Columns(wdi...)
	}

	rows := t.RawGetString("__rows")
	if rows.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(rows.(*golua.LTable)), state, lg)
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
	/// @struct WidgetButtonArrow
	/// @prop type {string<gui.WidgetType>}
	/// @prop dir {int<gui.Direction>}
	/// @method on_click(self, {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON_ARROW))
	t.RawSetString("dir", golua.LNumber(dir))
	t.RawSetString("__click", golua.LNil)

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	return t
}

func buttonArrowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	dir := t.RawGetString("dir").(golua.LNumber)
	b := g.ArrowButton(g.Direction(dir))

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func treeNodeTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetTreeNode
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.TreeNodeFlags>) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TREE_NODE))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func treeNodeBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	n := g.TreeNode(string(label))

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		n.Flags(g.TreeNodeFlags(flags.(golua.LNumber)))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		n.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return n
}

func treeTableRowTable(state *golua.LState, label string, widgets golua.LValue) *golua.LTable {
	/// @struct TreeTableRow
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @prop widgets {[]struct<gui.Widget>}
	/// @method flags(self, flags int<gui.TreeNodeFlags>) -> self
	/// @method children(self, []struct<gui.TreeTableRow>) -> self
	/// @desc
	/// Only used in tree table widget rows.

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TREE_TABLE_ROW))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("widgets", widgets)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__children", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "children", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__children", lt)
	})

	return t
}

func treeTableRowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.TreeTableRowWidget {
	label := t.RawGetString("label").(golua.LString)
	var widgets []g.Widget

	wid := t.RawGetString("widgets")
	if wid.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid.(*golua.LTable)), state, lg), lg)
	}

	n := g.TreeTableRow(string(label), widgets...)

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		n.Flags(g.TreeNodeFlags(flags.(golua.LNumber)))
	}

	children := t.RawGetString("__children")
	if children.Type() == golua.LTTable {
		rwid := parseWidgets(parseTable(children.(*golua.LTable)), state, lg)
		childs := []*g.TreeTableRowWidget{}
		for _, c := range rwid {
			childs = append(childs, treeTableRowBuild(r, lg, state, c))
		}
		n.Children(childs...)
	}

	return n
}

func treeTableTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetTreeTable
	/// @prop type {string<gui.WidgetType>}
	/// @method flags(self, flags int<gui.TableFlags>) -> self
	/// @method size(self, width float, height float) -> self
	/// @method columns(self, []struct<TableColumn>) -> self
	/// @method rows(self, []struct<TreeTableRow>) -> self
	/// @method freeze(self, col int, row int) -> self - Can be called multiple times.

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_TREE_TABLE))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__columns", golua.LNil)
	t.RawSetString("__rows", golua.LNil)
	t.RawSetString("__freeze", state.NewTable())
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "columns", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__columns", lt)
	})

	tableBuilderFunc(state, t, "rows", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__rows", lt)
	})

	tableBuilderFunc(state, t, "freeze", func(state *golua.LState, t *golua.LTable) {
		col := state.CheckNumber(-2)
		row := state.CheckNumber(-1)
		pt := state.NewTable()
		pt.RawSetString("col", col)
		pt.RawSetString("row", row)

		ft := t.RawGetString("__freeze").(*golua.LTable)
		ft.Append(pt)
	})

	return t
}

func treeTableBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tb := g.TreeTable()

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		tb.Flags(g.TableFlags(flags.(golua.LNumber)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		tb.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	freeze := t.RawGetString("__freeze").(*golua.LTable)
	for i := range freeze.Len() {
		pt := freeze.RawGetInt(i + 1).(*golua.LTable)
		col := pt.RawGetString("col").(golua.LNumber)
		row := pt.RawGetString("row").(golua.LNumber)

		tb.Freeze(int(col), int(row))
	}

	columns := t.RawGetString("__columns")
	if columns.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(columns.(*golua.LTable)), state, lg)
		wdi := []*g.TableColumnWidget{}
		for _, w := range wd {
			i := tableColumnBuild(w)
			wdi = append(wdi, i)
		}
		tb.Columns(wdi...)
	}

	rows := t.RawGetString("__rows")
	if rows.Type() == golua.LTTable {
		wd := parseWidgets(parseTable(rows.(*golua.LTable)), state, lg)
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
	/// @struct WidgetWindow
	/// @prop type {string<gui.WidgetType>}
	/// @prop single {bool}
	/// @prop menubar {bool}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.WindowFlags>) -> self
	/// @method size(self, width float, height float) -> self
	/// @method pos(self, x float, y float) -> self
	/// @method is_open(self, bool) -> self
	/// @method bring_to_front(self) -> self
	/// @method ready(self, {function(struct<gui.StateWindow>)}) -> self
	/// @method register_keyboard_shortcuts(self, []struct<gui.Shortcut>) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_WINDOW_SINGLE))
	t.RawSetString("single", golua.LBool(single))
	t.RawSetString("menubar", golua.LBool(menubar))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__front", golua.LNil)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__open", golua.LNil)
	t.RawSetString("__posx", golua.LNil)
	t.RawSetString("__posy", golua.LNil)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__ready", golua.LNil)
	t.RawSetString("__shortcuts", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "pos", func(state *golua.LState, t *golua.LTable) {
		posx := state.CheckNumber(-2)
		posy := state.CheckNumber(-1)
		t.RawSetString("__posx", posx)
		t.RawSetString("__posy", posy)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		open := state.CheckNumber(-1)
		t.RawSetString("__open", open)
	})

	tableBuilderFunc(state, t, "bring_to_front", func(state *golua.LState, t *golua.LTable) {
		t.RawSetString("__front", golua.LTrue)
	})

	tableBuilderFunc(state, t, "ready", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__ready", fn)
	})

	tableBuilderFunc(state, t, "register_keyboard_shortcuts", func(state *golua.LState, t *golua.LTable) {
		st := state.CheckTable(-1)
		t.RawSetString("__shortcuts", st)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
		windowBuild(r, lg, state, t)
	})

	return t
}

func windowBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) *g.WindowWidget {
	var w *g.WindowWidget

	single := t.RawGetString("single").(golua.LBool)
	if single {
		menubar := t.RawGetString("menubar").(golua.LBool)
		if menubar {
			w = g.SingleWindowWithMenuBar()
		} else {
			w = g.SingleWindow()
		}
	} else {
		label := t.RawGetString("label").(golua.LString)
		w = g.Window(string(label))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		w.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		w.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	posx := t.RawGetString("__posx")
	posy := t.RawGetString("__posy")
	if posx.Type() == golua.LTNumber && posy.Type() == golua.LTNumber {
		w.Pos(float32(posx.(golua.LNumber)), float32(posy.(golua.LNumber)))
	}

	front := t.RawGetString("__front")
	if front.Type() == golua.LTBool {
		if front.(golua.LBool) {
			w.BringToFront()
		}
	}

	open := t.RawGetString("__open")
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		w.IsOpen(ref.Value.(*bool))
	}

	/// @struct StateWindow
	/// @method current_position() -> float, float
	/// @method current_size() -> float, float
	/// @method has_focus() -> bool
	ready := t.RawGetString("__ready")
	if ready.Type() == golua.LTFunction {
		fnt := state.NewTable()

		fnt.RawSetString("current_position", state.NewFunction(func(state *golua.LState) int {
			x, y := w.CurrentPosition()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		}))

		fnt.RawSetString("current_size", state.NewFunction(func(state *golua.LState) int {
			w, h := w.CurrentSize()

			state.Push(golua.LNumber(w))
			state.Push(golua.LNumber(h))
			return 2
		}))

		fnt.RawSetString("has_focus", state.NewFunction(func(state *golua.LState) int {
			f := w.HasFocus()

			state.Push(golua.LBool(f))
			return 1
		}))

		state.Push(ready)
		state.Push(fnt)
		state.Call(1, 0)
	}

	shortcuts := t.RawGetString("__shortcuts")
	if shortcuts.Type() == golua.LTTable {
		stList := []g.WindowShortcut{}
		st := shortcuts.(*golua.LTable)
		for i := range st.Len() {
			s := st.RawGetInt(i + 1).(*golua.LTable)

			key := s.RawGetString("key").(golua.LNumber)
			mod := s.RawGetString("mod").(golua.LNumber)
			callback := s.RawGetString("callback")

			shortcut := g.WindowShortcut{
				Key:      g.Key(key),
				Modifier: g.Modifier(mod),
				Callback: func() {
					state.Push(callback)
					state.Call(0, 0)
				},
			}

			stList = append(stList, shortcut)
		}

		w.RegisterKeyboardShortcuts(stList...)
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		w.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return w
}

func popupModalTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetPopupModel
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.WindowFlags>) -> self
	/// @method is_open(self, bool) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_POPUP_MODAL))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__open", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "is_open", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__open", flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func popupModalBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	m := g.PopupModal(string(label))

	open := t.RawGetString("__open")
	if open.Type() == golua.LTNumber {
		ref, err := r.CR_REF.Item(int(open.(golua.LNumber)))
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		m.IsOpen(ref.Value.(*bool))
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		m.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	return m
}

func popupTable(state *golua.LState, label string) *golua.LTable {
	/// @struct WidgetPopup
	/// @prop type {string<gui.WidgetType>}
	/// @prop label {string}
	/// @method flags(self, flags int<gui.WindowFlags>) -> self
	/// @method layout(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_POPUP))
	t.RawSetString("label", golua.LString(label))
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "layout", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func popupBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	label := t.RawGetString("label").(golua.LString)
	m := g.Popup(string(label))

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		m.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		m.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	return m
}

func splitLayoutTable(state *golua.LState, direction, floatref int, layout1 golua.LValue, layout2 golua.LValue) *golua.LTable {
	/// @struct WidgetSplitLayout
	/// @prop type {string<gui.WidgetType>}
	/// @prop direction {int<gui.SplitDirection>}
	/// @prop floatref {int<ref.FLOAT32>}
	/// @prop layout1 {[]struct<gui.Widget>}
	/// @prop layout2 {[]struct<gui.Widget>}
	/// @method border(self, bool) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_LAYOUT_SPLIT))
	t.RawSetString("direction", golua.LNumber(direction))
	t.RawSetString("floatref", golua.LNumber(floatref))
	t.RawSetString("layout1", layout1)
	t.RawSetString("layout2", layout2)
	t.RawSetString("__border", golua.LNil)

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		border := state.CheckBool(-1)
		t.RawSetString("__border", golua.LBool(border))
	})

	return t
}

func splitLayoutBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	direction := t.RawGetString("direction").(golua.LNumber)

	floatref := t.RawGetString("floatref")
	ref, err := r.CR_REF.Item(int(floatref.(golua.LNumber)))
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	pos := ref.Value.(*float32)

	var widgets1 []g.Widget
	wid1 := t.RawGetString("layout1")
	if wid1.Type() == golua.LTTable {
		widgets1 = layoutBuild(r, state, parseWidgets(parseTable(wid1.(*golua.LTable)), state, lg), lg)
	}

	var widgets2 []g.Widget
	wid2 := t.RawGetString("layout2")
	if wid2.Type() == golua.LTTable {
		widgets2 = layoutBuild(r, state, parseWidgets(parseTable(wid2.(*golua.LTable)), state, lg), lg)
	}

	s := g.SplitLayout(g.SplitDirection(direction), pos, g.Layout(widgets1), g.Layout(widgets2))

	border := t.RawGetString("__border")
	if border.Type() == golua.LTBool {
		s.Border(bool(border.(golua.LBool)))
	}

	return s
}

func splitterTable(state *golua.LState, direction, floatref int) *golua.LTable {
	/// @struct WidgetSplitter
	/// @prop type {string<gui.WidgetType>}
	/// @prop direction {int<gui.SplitDirection>}
	/// @prop floatref {int<ref.FLOAT32>}
	/// @method size(self, width float, height float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_SPLITTER))
	t.RawSetString("direction", golua.LNumber(direction))
	t.RawSetString("floatref", golua.LNumber(floatref))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	return t
}

func splitterBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	direction := t.RawGetString("direction").(golua.LNumber)

	floatref := t.RawGetString("floatref")
	ref, err := r.CR_REF.Item(int(floatref.(golua.LNumber)))
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	pos := ref.Value.(*float32)

	s := g.Splitter(g.SplitDirection(direction), pos)

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		s.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return s
}

func stackTable(state *golua.LState, visible int, widgets golua.LValue) *golua.LTable {
	/// @struct WidgetStack
	/// @prop type {string<gui.WidgetType>}
	/// @prop visible {int}
	/// @prop widgets {[]struct<gui.Widget>}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_STACK))
	t.RawSetString("visible", golua.LNumber(visible))
	t.RawSetString("widgets", widgets)

	return t
}

func stackBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	visible := t.RawGetString("visible").(golua.LNumber)

	var widgets []g.Widget
	wid1 := t.RawGetString("widgets")
	if wid1.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid1.(*golua.LTable)), state, lg), lg)
	}

	s := g.Stack(int32(visible), widgets...)

	return s
}

func alignTable(state *golua.LState, at int) *golua.LTable {
	/// @struct WidgetAlign
	/// @prop type {string<gui.WidgetType>}
	/// @prop at {int<gui.Alignment>}
	/// @method to(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_ALIGN))
	t.RawSetString("at", golua.LNumber(at))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func alignBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	at := t.RawGetString("at").(golua.LNumber)
	a := g.Align(g.AlignmentType(at))

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		a.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return a
}

func msgBoxTable(state *golua.LState, title, content string) *golua.LTable {
	/// @struct WidgetMSGBox
	/// @prop type {string<gui.WidgetType>}
	/// @prop title {string}
	/// @prop content {string}
	/// @method buttons(self, int<gui.MSGBoxButtons>) -> self
	/// @method result_callback(self, {function(bool)}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MSG_BOX))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("content", golua.LString(content))
	t.RawSetString("__buttons", golua.LNil)
	t.RawSetString("__callback", golua.LNil)

	tableBuilderFunc(state, t, "buttons", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckNumber(-1)
		t.RawSetString("__buttons", b)
	})

	tableBuilderFunc(state, t, "result_callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__callback", fn)
	})

	tableBuilderFunc(state, t, "build", func(state *golua.LState, t *golua.LTable) {
		msgBoxBuild(state, t)
	})

	return t
}

func msgBoxBuild(state *golua.LState, t *golua.LTable) *g.MsgboxWidget {
	title := t.RawGetString("title").(golua.LString)
	content := t.RawGetString("content").(golua.LString)
	m := g.Msgbox(string(title), string(content))

	buttons := t.RawGetString("__buttons")
	if buttons.Type() == golua.LTNumber {
		m.Buttons(g.MsgboxButtons(buttons.(golua.LNumber)))
	}

	callback := t.RawGetString("__callback")
	if callback.Type() == golua.LTFunction {
		m.ResultCallback(func(dr g.DialogResult) {
			state.Push(callback)
			state.Push(golua.LBool(dr))
			state.Call(1, 0)
		})
	}

	return m
}

func msgBoxPrepareTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetMSGBoxPrepare
	/// @prop type {string<gui.WidgetType>}
	/// @desc
	/// This is used internally with gui.prepare_msg_box().

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_MSG_BOX_PREPARE))

	return t
}

func msgBoxPrepareBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	return g.PrepareMsgbox()
}

func buttonInvisibleTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetButtonInvisible
	/// @prop type {string<gui.WidgetType>}
	/// @method size(self, width float, height float) -> self
	/// @method on_click(self, {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON_INVISIBLE))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__click", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	return t
}

func buttonInvisibleBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.InvisibleButton()

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	return b
}

func buttonImageTable(state *golua.LState, id int, sync, cache bool) *golua.LTable {
	/// @struct WidgetButtonImage
	/// @prop type {string<gui.WidgetType>}
	/// @prop id {int<collection.IMAGE>}
	/// @prop cachedid {int<collection.CRATE_CACHEDIMAGE>}
	/// @prop sync {bool}
	/// @prop cached {bool}
	/// @method size(self, width float, height float) -> self
	/// @method on_click(self, {function()}) -> self
	/// @method bg_color(self, struct<image.Color>) -> self
	/// @method tint_color(self, struct<image.Color>) -> self
	/// @method frame_padding(self, padding float) -> self
	/// @method uv(self, uv0 struct<image.Point>, uv1 struct<image.Point>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_BUTTON_IMAGE))
	t.RawSetString("id", golua.LNumber(id))
	t.RawSetString("cachedid", golua.LNumber(id))
	t.RawSetString("sync", golua.LBool(sync))
	t.RawSetString("cached", golua.LBool(cache))
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__click", golua.LNil)
	t.RawSetString("__bgcolor", golua.LNil)
	t.RawSetString("__padding", golua.LNil)
	t.RawSetString("__tint", golua.LNil)
	t.RawSetString("__uv0", golua.LNil)
	t.RawSetString("__uv1", golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__click", fn)
	})

	tableBuilderFunc(state, t, "bg_color", func(state *golua.LState, t *golua.LTable) {
		tc := state.CheckTable(-1)
		t.RawSetString("__bgcolor", tc)
	})

	tableBuilderFunc(state, t, "tint_color", func(state *golua.LState, t *golua.LTable) {
		tc := state.CheckTable(-1)
		t.RawSetString("__tint", tc)
	})

	tableBuilderFunc(state, t, "frame_padding", func(state *golua.LState, t *golua.LTable) {
		n := state.CheckNumber(-1)
		t.RawSetString("__padding", n)
	})

	tableBuilderFunc(state, t, "uv", func(state *golua.LState, t *golua.LTable) {
		uv0 := state.CheckTable(-2)
		uv1 := state.CheckTable(-1)
		t.RawSetString("__uv0", uv0)
		t.RawSetString("__uv1", uv1)
	})

	return t
}

func buttonImageBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ig := t.RawGetString("id").(golua.LNumber)
	var img image.Image

	sync := t.RawGetString("sync").(golua.LBool)
	cache := t.RawGetString("cached").(golua.LBool)

	if !sync {
		if cache {
			ci, err := r.CR_CIM.Item(int(ig))
			if err == nil {
				img = ci.Image
			}
		} else {
			<-r.IC.Schedule(state, int(ig), &collection.Task[collection.ItemImage]{
				Lib:  LIB_GUI,
				Name: "wg_button_image",
				Fn: func(i *collection.Item[collection.ItemImage]) {
					img = i.Self.Image
				},
			})
		}
	} else {
		item := r.IC.Item(int(ig))
		if item.Self.Image != nil {
			img = item.Self.Image
		}
	}

	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	b := g.ImageButtonWithRgba(img)

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		b.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	click := t.RawGetString("__click")
	if click.Type() == golua.LTFunction {
		b.OnClick(func() {
			state.Push(click)
			state.Call(0, 0)
		})
	}

	bgcolor := t.RawGetString("__bgcolor")
	if bgcolor.Type() == golua.LTTable {
		rgba := imageutil.ColorTableToRGBAColor(bgcolor.(*golua.LTable))
		b.BgColor(rgba)
	}

	tint := t.RawGetString("__tint")
	if tint.Type() == golua.LTTable {
		rgba := imageutil.ColorTableToRGBAColor(tint.(*golua.LTable))
		b.TintColor(rgba)
	}

	padding := t.RawGetString("__padding")
	if padding.Type() == golua.LTNumber {
		b.FramePadding(int(padding.(golua.LNumber)))
	}

	uv0 := t.RawGetString("__uv0")
	uv1 := t.RawGetString("__uv1")
	if uv0.Type() == golua.LTTable && uv1.Type() == golua.LTTable {
		p1 := imageutil.TableToPoint(uv0.(*golua.LTable))
		p2 := imageutil.TableToPoint(uv1.(*golua.LTable))

		b.UV(p1, p2)
	}

	return b
}

func styleTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetStyle
	/// @prop type {string<gui.WidgetType>}
	/// @method set_disabled(self, bool) -> self
	/// @method to(self, []struct<gui.Widget>) -> self
	/// @method set_font_size(self, float) -> self
	/// @method set_color(self, int<gui.StyleColorID>, struct<image.Color>) -> self
	/// @method set_style(self, int<gui.StyleVarID>, width float, height float) -> self
	/// @method set_style_float(self, int<gui.StyleVarID>, float) -> self
	/// @method font(self, int<ref.FONT>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_STYLE))
	t.RawSetString("__disabled", golua.LNil)
	t.RawSetString("__widgets", golua.LNil)
	t.RawSetString("__fontsize", golua.LNil)
	t.RawSetString("__colors", state.NewTable())
	t.RawSetString("__styles", state.NewTable())
	t.RawSetString("__stylesfloat", state.NewTable())
	t.RawSetString("__font", golua.LNil)

	tableBuilderFunc(state, t, "set_disabled", func(state *golua.LState, t *golua.LTable) {
		d := state.CheckBool(-1)
		t.RawSetString("__disabled", golua.LBool(d))
	})

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	tableBuilderFunc(state, t, "set_font_size", func(state *golua.LState, t *golua.LTable) {
		fnt := state.CheckNumber(-1)
		t.RawSetString("__fontsize", fnt)
	})

	tableBuilderFunc(state, t, "set_color", func(state *golua.LState, t *golua.LTable) {
		cid := state.CheckNumber(-2)
		col := state.CheckTable(-1)
		ct := state.NewTable()
		ct.RawSetString("colorid", cid)
		ct.RawSetString("color", col)

		ft := t.RawGetString("__colors").(*golua.LTable)
		ft.Append(ct)
	})

	tableBuilderFunc(state, t, "set_style", func(state *golua.LState, t *golua.LTable) {
		sid := state.CheckNumber(-3)
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		st := state.NewTable()
		st.RawSetString("styleid", sid)
		st.RawSetString("width", width)
		st.RawSetString("height", height)

		ft := t.RawGetString("__styles").(*golua.LTable)
		ft.Append(st)
	})

	tableBuilderFunc(state, t, "set_style_float", func(state *golua.LState, t *golua.LTable) {
		sid := state.CheckNumber(-2)
		float := state.CheckNumber(-1)
		st := state.NewTable()
		st.RawSetString("styleid", sid)
		st.RawSetString("float", float)

		ft := t.RawGetString("__stylesfloat").(*golua.LTable)
		ft.Append(st)
	})

	tableBuilderFunc(state, t, "font", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		t.RawSetString("__font", v)
	})

	return t
}

func styleBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	s := g.Style()

	disabled := t.RawGetString("__disabled")
	if disabled.Type() == golua.LTBool {
		s.SetDisabled(bool(disabled.(golua.LBool)))
	}

	fontsize := t.RawGetString("__fontsize")
	if fontsize.Type() == golua.LTNumber {
		s.SetFontSize(float32(fontsize.(golua.LNumber)))
	}

	colors := t.RawGetString("__colors").(*golua.LTable)
	for i := range colors.Len() {
		ct := colors.RawGetInt(i + 1).(*golua.LTable)
		cid := ct.RawGetString("colorid").(golua.LNumber)
		col := ct.RawGetString("color").(*golua.LTable)

		rgba := imageutil.ColorTableToRGBAColor(col)
		s.SetColor(g.StyleColorID(cid), rgba)
	}

	styles := t.RawGetString("__styles").(*golua.LTable)
	for i := range styles.Len() {
		st := styles.RawGetInt(i + 1).(*golua.LTable)
		sid := st.RawGetString("styleid").(golua.LNumber)
		width := st.RawGetString("width").(golua.LNumber)
		height := st.RawGetString("height").(golua.LNumber)

		s.SetStyle(g.StyleVarID(sid), float32(width), float32(height))
	}

	stylesfloat := t.RawGetString("__stylesfloat").(*golua.LTable)
	for i := range stylesfloat.Len() {
		st := stylesfloat.RawGetInt(i + 1).(*golua.LTable)
		sid := st.RawGetString("styleid").(golua.LNumber)
		float := st.RawGetString("float").(golua.LNumber)

		s.SetStyleFloat(g.StyleVarID(sid), float32(float))
	}

	fontref := t.RawGetString("__font")
	if fontref.Type() == golua.LTNumber {
		ref := int(fontref.(golua.LNumber))
		sref, err := r.CR_REF.Item(ref)
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		font := sref.Value.(*g.FontInfo)

		s.SetFont(font)
	}

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		s.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return s
}

func customTable(state *golua.LState, builder *golua.LFunction) *golua.LTable {
	/// @struct WidgetCustom
	/// @prop type {string<gui.WidgetType>}
	/// @prop builder {function()}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CUSTOM))
	t.RawSetString("builder", builder)

	return t
}

func customBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	builder := t.RawGetString("builder").(*golua.LFunction)

	c := g.Custom(func() {
		state.Push(builder)
		state.Call(0, 0)
	})

	return c
}

func eventHandlerTable(state *golua.LState) *golua.LTable {
	/// @struct WidgetEvent
	/// @prop type {string<gui.WidgetType>}
	/// @method on_activate(self, {function()}) -> self
	/// @method on_active(self, {function()}) -> self
	/// @method on_deactivate(self, {function()}) -> self
	/// @method on_hover(self, {function()}) -> self
	/// @method on_click(self, int<gui.MouseButton>, {function()}) -> self
	/// @method on_dclick(self, int<gui.MouseButton>, {function()}) -> self
	/// @method on_key_down(self, int<gui.Key>, {function()}) -> self
	/// @method on_key_pressed(self, int<gui.Key>, {function()}) -> self
	/// @method on_key_released(self, int<gui.Key>, {function()}) -> self
	/// @method on_mouse_down(self, int<gui.MouseButton>, {function()}) -> self
	/// @method on_mouse_released(self, int<gui.MouseButton>, {function()}) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_EVENT_HANDLER))
	t.RawSetString("__activate", golua.LNil)
	t.RawSetString("__active", golua.LNil)
	t.RawSetString("__deactivate", golua.LNil)
	t.RawSetString("__hover", golua.LNil)
	t.RawSetString("__click", state.NewTable())
	t.RawSetString("__dclick", state.NewTable())
	t.RawSetString("__keydown", state.NewTable())
	t.RawSetString("__keypressed", state.NewTable())
	t.RawSetString("__keyreleased", state.NewTable())
	t.RawSetString("__mousedown", state.NewTable())
	t.RawSetString("__mousereleased", state.NewTable())

	tableBuilderFunc(state, t, "on_activate", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__activate", fn)
	})

	tableBuilderFunc(state, t, "on_active", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__active", fn)
	})

	tableBuilderFunc(state, t, "on_deactivate", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__deactivate", fn)
	})

	tableBuilderFunc(state, t, "on_hover", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		t.RawSetString("__hover", fn)
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__click").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_dclick", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__dclick").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_down", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__keydown").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_pressed", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__keypressed").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_released", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__keyreleased").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_mouse_down", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__mousedown").(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_mouse_released", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		ev.RawSetString("key", key)
		ev.RawSetString("callback", cb)

		ft := t.RawGetString("__mousereleased").(*golua.LTable)
		ft.Append(ev)
	})

	return t
}

func eventHandlerBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	e := g.Event()

	activate := t.RawGetString("__activate")
	if activate.Type() == golua.LTFunction {
		e.OnActivate(func() {
			state.Push(activate)
			state.Call(0, 0)
		})
	}

	active := t.RawGetString("__active")
	if active.Type() == golua.LTFunction {
		e.OnActive(func() {
			state.Push(active)
			state.Call(0, 0)
		})
	}

	deactivate := t.RawGetString("__deactivate")
	if deactivate.Type() == golua.LTFunction {
		e.OnDeactivate(func() {
			state.Push(deactivate)
			state.Call(0, 0)
		})
	}

	hover := t.RawGetString("__hover")
	if hover.Type() == golua.LTFunction {
		e.OnHover(func() {
			state.Push(hover)
			state.Call(0, 0)
		})
	}

	click := t.RawGetString("__click").(*golua.LTable)
	for i := range click.Len() {
		events := click.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnClick(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	dclick := t.RawGetString("__dclick").(*golua.LTable)
	for i := range dclick.Len() {
		events := dclick.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnDClick(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keydown := t.RawGetString("__keydown").(*golua.LTable)
	for i := range keydown.Len() {
		events := keydown.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnKeyDown(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keypressed := t.RawGetString("__keypressed").(*golua.LTable)
	for i := range keypressed.Len() {
		events := keypressed.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnKeyPressed(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keyreleased := t.RawGetString("__keyreleased").(*golua.LTable)
	for i := range keyreleased.Len() {
		events := keyreleased.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnKeyReleased(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	mousedown := t.RawGetString("__mousedown").(*golua.LTable)
	for i := range mousedown.Len() {
		events := mousedown.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnMouseDown(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	mousereleased := t.RawGetString("__mousereleased").(*golua.LTable)
	for i := range mousereleased.Len() {
		events := mousereleased.RawGetInt(i + 1).(*golua.LTable)
		key := events.RawGetString("key").(golua.LNumber)
		callback := events.RawGetString("callback").(*golua.LFunction)

		e.OnMouseReleased(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	return e
}

func cssTagTable(state *golua.LState, tag string) *golua.LTable {
	/// @struct WidgetCSSTag
	/// @prop type {string<gui.WidgetType>}
	/// @prop tag {string}
	/// @method to(self, []struct<gui.Widget>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CSS_TAG))
	t.RawSetString("tag", golua.LString(tag))
	t.RawSetString("__widgets", golua.LNil)

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		t.RawSetString("__widgets", lt)
	})

	return t
}

func cssTagBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tag := t.RawGetString("tag").(golua.LString)
	c := g.CSSTag(string(tag))

	layout := t.RawGetString("__widgets")
	if layout.Type() == golua.LTTable {
		c.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable)), state, lg), lg)...)
	}

	return c
}

func plotTable(state *golua.LState, title string) *golua.LTable {
	/// @struct WidgetPlot
	/// @prop type {string<gui.WidgetType>}
	/// @prop title {string}
	/// @method axis_limits(self, xmin float, xmax float, ymin float, ymax float, cond int<gui.Condition>) -> self
	/// @method flags(self, flags int<gui.PlotFlags>) -> self
	/// @method set_xaxis_label(self, axis int<gui.PlotAxis>, label string) -> self
	/// @method set_yaxis_label(self, axis int<gui.PlotAxis>, label string) -> self
	/// @method size(self, width float, height float) -> self
	/// @method x_axeflags(self, flags int<gui.PlotAxisFlags>) -> self
	/// @method xticks(self, ticks []struct<gui.PlotTicker>, default bool) -> self
	/// @method y_axeflags(self, flags1 int<gui.PlotAxisFlags>, flags2 int<gui.PlotAxisFlags>, flags3 int<gui.PlotAxisFlags>) -> self
	/// @method yticks(self, ticks []struct<gui.PlotTicker>) -> self
	/// @method plots(self, []struct<gui.Plot>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_PLOT))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("__xmin", golua.LNil)
	t.RawSetString("__xmax", golua.LNil)
	t.RawSetString("__ymin", golua.LNil)
	t.RawSetString("__ymax", golua.LNil)
	t.RawSetString("__cond", golua.LNil)
	t.RawSetString("__flags", golua.LNil)
	t.RawSetString("__xlabels", state.NewTable())
	t.RawSetString("__ylabels", state.NewTable())
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__xaxeflags", golua.LNil)
	t.RawSetString("__xticks", golua.LNil)
	t.RawSetString("__xaticksdefault", golua.LNil)
	t.RawSetString("__yaxeflags1", golua.LNil)
	t.RawSetString("__yaxeflags2", golua.LNil)
	t.RawSetString("__yaxeflags3", golua.LNil)
	t.RawSetString("__yticks", state.NewTable())
	t.RawSetString("__plots", golua.LNil)

	tableBuilderFunc(state, t, "axis_limits", func(state *golua.LState, t *golua.LTable) {
		xmin := state.CheckNumber(-5)
		xmax := state.CheckNumber(-4)
		ymin := state.CheckNumber(-3)
		ymax := state.CheckNumber(-2)
		cond := state.CheckNumber(-1)
		t.RawSetString("__xmin", xmin)
		t.RawSetString("__xmax", xmax)
		t.RawSetString("__ymin", ymin)
		t.RawSetString("__ymax", ymax)
		t.RawSetString("__cond", cond)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__flags", flags)
	})

	tableBuilderFunc(state, t, "set_xaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		lt.RawSetString("axis", axis)
		lt.RawSetString("label", golua.LString(label))

		ft := t.RawGetString("__xlabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "set_yaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		lt.RawSetString("axis", axis)
		lt.RawSetString("label", golua.LString(label))

		ft := t.RawGetString("__ylabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		t.RawSetString("__width", width)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "x_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		t.RawSetString("__xaxeflags", flags)
	})

	tableBuilderFunc(state, t, "xticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-2)
		dflt := state.CheckBool(-1)
		t.RawSetString("__xticks", ticks)
		t.RawSetString("__xaticksdefault", golua.LBool(dflt))
	})

	tableBuilderFunc(state, t, "y_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags1 := state.CheckNumber(-3)
		flags2 := state.CheckNumber(-2)
		flags3 := state.CheckNumber(-1)
		t.RawSetString("__yaxeflags1", flags1)
		t.RawSetString("__yaxeflags2", flags2)
		t.RawSetString("__yaxeflags3", flags3)
	})

	tableBuilderFunc(state, t, "yticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-3)
		dflt := state.CheckBool(-2)
		axis := state.CheckNumber(-1)
		lt := state.NewTable()
		lt.RawSetString("ticks", ticks)
		lt.RawSetString("dflt", golua.LBool(dflt))
		lt.RawSetString("axis", axis)

		ft := t.RawGetString("__ylabels").(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "plots", func(state *golua.LState, t *golua.LTable) {
		plots := state.CheckTable(-1)
		t.RawSetString("__plots", plots)
	})

	return t
}

func plotBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	/// @interface Plot
	/// @prop type {string<gui.PlotType>}

	title := t.RawGetString("title").(golua.LString)
	p := g.Plot(string(title))

	width := t.RawGetString("__width")
	height := t.RawGetString("__height")
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		p.Size(int(width.(golua.LNumber)), int(height.(golua.LNumber)))
	}

	flags := t.RawGetString("__flags")
	if flags.Type() == golua.LTNumber {
		p.Flags(g.PlotFlags(flags.(golua.LNumber)))
	}

	xaxeflags := t.RawGetString("__xaxeflags")
	if xaxeflags.Type() == golua.LTNumber {
		p.XAxeFlags(g.PlotAxisFlags(xaxeflags.(golua.LNumber)))
	}

	yaxeflags1 := t.RawGetString("__yaxeflags1")
	yaxeflags2 := t.RawGetString("__yaxeflags2")
	yaxeflags3 := t.RawGetString("__yaxeflags3")
	if yaxeflags1.Type() == golua.LTNumber && yaxeflags2.Type() == golua.LTNumber && yaxeflags3.Type() == golua.LTNumber {
		p.YAxeFlags(g.PlotAxisFlags(yaxeflags1.(golua.LNumber)), g.PlotAxisFlags(yaxeflags2.(golua.LNumber)), g.PlotAxisFlags(yaxeflags3.(golua.LNumber)))
	}

	xmin := t.RawGetString("__xmin")
	xmax := t.RawGetString("__xmax")
	ymin := t.RawGetString("__ymin")
	ymax := t.RawGetString("__ymax")
	cond := t.RawGetString("__cond")
	if xmin.Type() == golua.LTNumber && xmax.Type() == golua.LTNumber && ymin.Type() == golua.LTNumber && ymax.Type() == golua.LTNumber && cond.Type() == golua.LTNumber {
		p.AxisLimits(
			float64(xmin.(golua.LNumber)), float64(xmax.(golua.LNumber)),
			float64(ymin.(golua.LNumber)), float64(ymax.(golua.LNumber)),
			g.ExecCondition(cond.(golua.LNumber)),
		)
	}

	xlabels := t.RawGetString("__xlabels").(*golua.LTable)
	for i := range xlabels.Len() {
		lt := xlabels.RawGetInt(i + 1).(*golua.LTable)
		axis := lt.RawGetString("axis").(golua.LNumber)
		label := lt.RawGetString("label").(golua.LString)

		p.SetXAxisLabel(g.PlotXAxis(axis), string(label))
	}

	ylabels := t.RawGetString("__ylabels").(*golua.LTable)
	for i := range ylabels.Len() {
		lt := ylabels.RawGetInt(i + 1).(*golua.LTable)
		axis := lt.RawGetString("axis").(golua.LNumber)
		label := lt.RawGetString("label").(golua.LString)

		p.SetYAxisLabel(g.PlotYAxis(axis), string(label))
	}

	xticks := t.RawGetString("__xticks")
	xticksdefault := t.RawGetString("__xticksdefault")
	if xticks.Type() == golua.LTTable && xticksdefault.Type() == golua.LTBool {
		ticks := []g.PlotTicker{}

		xtickst := xticks.(*golua.LTable)
		for i := range xtickst.Len() {
			tick := xtickst.RawGetInt(i + 1).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick))
		}

		p.XTicks(ticks, bool(xticksdefault.(golua.LBool)))
	}

	yticks := t.RawGetString("__yticks").(*golua.LTable)
	for z := range yticks.Len() {
		yticksaxis := yticks.RawGetInt(z + 1).(*golua.LTable)

		ytickaxis := yticksaxis.RawGetString("ticks").(*golua.LTable)
		dflt := yticksaxis.RawGetString("dflt").(golua.LBool)
		axis := yticksaxis.RawGetString("axis").(golua.LNumber)

		ticks := []g.PlotTicker{}

		for i := range ytickaxis.Len() {
			tick := ytickaxis.RawGetInt(i + 1).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick))
		}

		p.YTicks(ticks, bool(dflt), g.ImPlotYAxis(axis))
	}

	plots := t.RawGetString("__plots")
	if plots.Type() == golua.LTTable {
		plist := []g.PlotWidget{}

		for i := range (plots.(*golua.LTable)).Len() {
			pt := plots.(*golua.LTable).RawGetInt(i + 1).(*golua.LTable)
			plottype := pt.RawGetString("type").(golua.LString)

			build := plotList[string(plottype)]
			plist = append(plist, build(r, lg, state, pt))
		}

		p.Plots(plist...)
	}

	return p
}

func plotTickerBuild(t *golua.LTable) g.PlotTicker {
	/// @struct PlotTicker
	/// @prop position {float}
	/// @prop label {string}

	position := t.RawGetString("position").(golua.LNumber)
	label := t.RawGetString("label").(golua.LString)

	return g.PlotTicker{
		Position: float64(position),
		Label:    string(label),
	}
}

func plotBarHTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotBarH
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method height(self, height float) -> self
	/// @method offset(self, offset float) -> self
	/// @method shift(self, shift float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_BAR_H))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__height", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__shift", golua.LNil)

	tableBuilderFunc(state, t, "height", func(state *golua.LState, t *golua.LTable) {
		height := state.CheckNumber(-1)
		t.RawSetString("__height", height)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		t.RawSetString("__shift", shift)
	})

	return t
}

func plotBarHBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.BarH(string(title), dataPoints)

	height := t.RawGetString("__height")
	if height.Type() == golua.LTNumber {
		p.Height(float64(height.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := t.RawGetString("__shift")
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotBarTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotBar
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method width(self, width float) -> self
	/// @method offset(self, offset float) -> self
	/// @method shift(self, shift float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_BAR))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__width", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__shift", golua.LNil)

	tableBuilderFunc(state, t, "width", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		t.RawSetString("__width", width)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		t.RawSetString("__shift", shift)
	})

	return t
}

func plotBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Bar(string(title), dataPoints)

	width := t.RawGetString("__width")
	if width.Type() == golua.LTNumber {
		p.Width(float64(width.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := t.RawGetString("__shift")
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotLineTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotLine
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method set_plot_y_axis(self, axis int<gui.PlotYAxis>) -> self
	/// @method offset(self, offset float) -> self
	/// @method x0(self, x0 float) -> self
	/// @method xscale(self, xscale float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_LINE))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__yaxis", golua.LNil)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__x0", golua.LNil)
	t.RawSetString("__xscale", golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		t.RawSetString("__yaxis", axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		t.RawSetString("__x0", x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		t.RawSetString("__xscale", xscale)
	})

	return t
}

func plotLineBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Line(string(title), dataPoints)

	yaxis := t.RawGetString("__yaxis")
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := t.RawGetString("__x0")
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := t.RawGetString("__xscale")
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotLineXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct PlotLineXY
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop xdata {[]float}
	/// @prop ydata {[]float}
	/// @method set_plot_y_axis(self, axis int<gui.PlotYAxis>) -> self
	/// @method offset(self, offset float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_LINE_XY))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("xdata", xdata)
	t.RawSetString("ydata", ydata)
	t.RawSetString("__yaxis", golua.LNil)
	t.RawSetString("__offset", golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		t.RawSetString("__yaxis", axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	return t
}

func plotLineXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	xdata := t.RawGetString("xdata").(*golua.LTable)
	ydata := t.RawGetString("ydata").(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := xdata.RawGetInt(i + 1).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := ydata.RawGetInt(i + 1).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.LineXY(string(title), xdataPoints, ydataPoints)

	yaxis := t.RawGetString("__yaxis")
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotPieTable(state *golua.LState, labels golua.LValue, data golua.LValue, x, y, radius float64) *golua.LTable {
	/// @struct PlotPieChart
	/// @prop type {string<gui.PlotType>}
	/// @prop labels {[]string}
	/// @prop data {[]float}
	/// @prop x {float}
	/// @prop y {float}
	/// @prop radius {float}
	/// @method angle0(self, angle0 float) -> self
	/// @method label_format(self, format string) -> self
	/// @method normalize(self, bool) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_PIE_CHART))
	t.RawSetString("labels", labels)
	t.RawSetString("data", data)
	t.RawSetString("x", golua.LNumber(x))
	t.RawSetString("y", golua.LNumber(y))
	t.RawSetString("radius", golua.LNumber(radius))
	t.RawSetString("__angle0", golua.LNil)
	t.RawSetString("__format", golua.LNil)
	t.RawSetString("__normalize", golua.LNil)

	tableBuilderFunc(state, t, "angle0", func(state *golua.LState, t *golua.LTable) {
		angle0 := state.CheckNumber(-1)
		t.RawSetString("__angle0", angle0)
	})

	tableBuilderFunc(state, t, "label_format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		t.RawSetString("__format", golua.LString(format))
	})

	tableBuilderFunc(state, t, "normalize", func(state *golua.LState, t *golua.LTable) {
		normalize := state.CheckBool(-1)
		t.RawSetString("__normalize", golua.LBool(normalize))
	})

	return t
}

func plotPieBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	labels := t.RawGetString("labels").(*golua.LTable)
	data := t.RawGetString("data").(*golua.LTable)
	x := t.RawGetString("x").(golua.LNumber)
	y := t.RawGetString("y").(golua.LNumber)
	radius := t.RawGetString("radius").(golua.LNumber)

	labelPoints := []string{}
	for i := range labels.Len() {
		point := labels.RawGetInt(i + 1).(golua.LString)
		labelPoints = append(labelPoints, string(point))
	}

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.PieChart(labelPoints, dataPoints, float64(x), float64(y), float64(radius))

	angle0 := t.RawGetString("__angle0")
	if angle0.Type() == golua.LTNumber {
		p.Angle0(float64(angle0.(golua.LNumber)))
	}

	format := t.RawGetString("__format")
	if format.Type() == golua.LTString {
		p.LabelFormat(string(format.(golua.LString)))
	}

	normalize := t.RawGetString("__normalize")
	if normalize.Type() == golua.LTBool {
		p.Normalize(bool(normalize.(golua.LBool)))
	}

	return p
}

func plotScatterTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct PlotScatter
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop data {[]float}
	/// @method offset(self, offset float) -> self
	/// @method x0(self, x0 float) -> self
	/// @method xscale(self, xscale float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_SCATTER))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("data", data)
	t.RawSetString("__offset", golua.LNil)
	t.RawSetString("__x0", golua.LNil)
	t.RawSetString("__xscale", golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		t.RawSetString("__x0", x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		t.RawSetString("__xscale", xscale)
	})

	return t
}

func plotScatterBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	data := t.RawGetString("data").(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := data.RawGetInt(i + 1).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Scatter(string(title), dataPoints)

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := t.RawGetString("__x0")
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := t.RawGetString("__xscale")
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotScatterXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct PlotScatterXY
	/// @prop type {string<gui.PlotType>}
	/// @prop title {string}
	/// @prop xdata {[]float}
	/// @prop ydata {[]float}
	/// @method offset(self, offset float) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_SCATTER_XY))
	t.RawSetString("title", golua.LString(title))
	t.RawSetString("xdata", xdata)
	t.RawSetString("ydata", ydata)
	t.RawSetString("__offset", golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		t.RawSetString("__offset", offset)
	})

	return t
}

func plotScatterXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := t.RawGetString("title").(golua.LString)
	xdata := t.RawGetString("xdata").(*golua.LTable)
	ydata := t.RawGetString("ydata").(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := xdata.RawGetInt(i + 1).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := ydata.RawGetInt(i + 1).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.ScatterXY(string(title), xdataPoints, ydataPoints)

	offset := t.RawGetString("__offset")
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotCustomTable(state *golua.LState, builder *golua.LFunction) *golua.LTable {
	/// @struct PlotCustom
	/// @prop type {string<gui.PlotType>}
	/// @prop builder {function()}

	t := state.NewTable()
	t.RawSetString("type", golua.LString(PLOT_CUSTOM))
	t.RawSetString("builder", builder)

	return t
}

func plotCustomBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	builder := t.RawGetString("builder").(*golua.LFunction)

	c := g.Custom(func() {
		state.Push(builder)
		state.Call(0, 0)
	})

	return c
}
