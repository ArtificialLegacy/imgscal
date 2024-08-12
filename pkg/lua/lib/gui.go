package lib

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"sync"
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

	/// @func window_set_icon_imgscal()
	/// @arg id
	/// @arg? circled
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

			iconPaths := []string{
				"favicon-16x16.png",
				"favicon-32x32.png",
			}
			if args["circled"].(bool) {
				iconPaths = []string{
					"favicon-16x16-circle.png",
					"favicon-32x32-circle.png",
				}
			}

			icons := []image.Image{}

			wd, _ := os.Getwd()
			for _, p := range iconPaths {
				f, err := os.Open(path.Join(wd, "assets", p))
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("cannot open %s", p), log.LEVEL_ERROR)), 0)
				}
				defer f.Close()

				ic, err := imageutil.Decode(f, imageutil.ENCODING_PNG)
				if err != nil {
					state.Error(golua.LString(lg.Append(fmt.Sprintf("%s is an invalid image: %s", p, err), log.LEVEL_ERROR)), 0)
				}

				icons = append(icons, ic)
			}

			w.SetIcon(icons...)
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

	/// @func window_set_icon_many()
	/// @arg id
	/// @arg icon_ids
	/// @blocking
	/// @desc
	/// setting multiple icons allows it select the closest to the system's desired size.
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

			imgids := args["icon_ids"].(map[string]any)
			imgList := []image.Image{}
			wg := sync.WaitGroup{}

			for _, id := range imgids {
				wg.Add(1)
				r.IC.Schedule(id.(int), &collection.Task[collection.ItemImage]{
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

	/// @func window_clear_icon()
	/// @arg id
	/// @desc
	/// resets window icon to default, same as window_set_icon_many(id, {})
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

	/// @func window_register_keyboard_shortcuts()
	/// @arg id
	/// @arg []shortcuts
	lib.CreateFunction(tab, "window_register_keyboard_shortcuts",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "shortcuts"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			w, err := r.CR_WIN.Item(args["id"].(int))
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("error getting window: %s", err), log.LEVEL_ERROR)), 0)
			}

			st := args["shortcuts"].(*golua.LTable)
			stList := []g.WindowShortcut{}
			for i := range st.Len() {
				s := state.GetTable(st, golua.LNumber(i+1)).(*golua.LTable)

				key := state.GetTable(s, golua.LString("key")).(golua.LNumber)
				mod := state.GetTable(s, golua.LString("mod")).(golua.LNumber)
				callback := state.GetTable(s, golua.LString("callback"))

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

	/// @func window_set_close_callback()
	/// @arg id
	/// @arg callback - returns bool
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
				return bool(state.ToBool(-1))
			})
			return 0
		})

	/// @func window_set_drop_callback()
	/// @arg id
	/// @arg callback([]string)
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

	/// @func window_additional_input_handler_callback()
	/// @arg id
	/// @arg callback(key, mod, action)
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

	/// @func window_close()
	/// @arg id
	/// @desc
	/// same as window_should_close(id, true)
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

	/// @func prepare_msg_box()
	/// @returns widget
	lib.CreateFunction(tab, "prepare_msg_box",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := msgBoxPrepareTable(state)

			state.Push(t)
			return 1
		})

	/// @func style_var_is_vec2()
	/// @arg var
	/// @returns bool
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

	/// @func style_var_string()
	/// @arg var
	/// @returns string
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

	/// @func style_var_from_string()
	/// @arg string
	/// @returns int
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

	/// @func shortcut
	/// @arg key
	/// @arg mod
	/// @arg callback
	/// @returns shortcut table
	lib.CreateFunction(tab, "shortcut",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
			{Type: lua.INT, Name: "mod"},
			{Type: lua.FUNC, Name: "callback"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct shortcut
			/// @prop key
			/// @prop mod
			/// @prop callback()

			key := args["key"].(int)
			mod := args["mod"].(int)
			callback := args["callback"].(*golua.LFunction)

			t := state.NewTable()
			state.SetTable(t, golua.LString("key"), golua.LNumber(key))
			state.SetTable(t, golua.LString("mod"), golua.LNumber(mod))
			state.SetTable(t, golua.LString("callback"), callback)

			state.Push(t)
			return 1
		})

	/// @func plot_ticker
	/// @arg position
	/// @arg label
	/// @returns plot ticker
	lib.CreateFunction(tab, "plot_ticker",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "position"},
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pos := args["position"].(float64)
			label := args["label"].(string)

			t := state.NewTable()
			state.SetTable(t, golua.LString("position"), golua.LNumber(pos))
			state.SetTable(t, golua.LString("label"), golua.LString(label))

			state.Push(t)
			return 1
		})

	/// @func css_parse
	/// @arg path
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

	/// @func calc_text_size()
	/// @arg text
	/// @returns width, height
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

	/// @func calc_text_size_width()
	/// @arg text
	/// @returns width
	lib.CreateFunction(tab, "calc_text_size_width",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, _ := g.CalcTextSize(args["text"].(string))

			state.Push(golua.LNumber(width))
			return 1
		})

	/// @func calc_text_size_height()
	/// @arg text
	/// @returns height
	lib.CreateFunction(tab, "calc_text_size_height",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			_, height := g.CalcTextSize(args["text"].(string))

			state.Push(golua.LNumber(height))
			return 1
		})

	/// @func calc_text_size_v()
	/// @arg text
	/// @arg hideAfterDoubleHash
	/// @arg wrapWidth
	/// @returns width, height
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

	/// @func available_region()
	/// @returns width, height
	lib.CreateFunction(tab, "available_region",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetAvailableRegion()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func frame_padding()
	/// @returns x, y
	lib.CreateFunction(tab, "frame_padding",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x, y := g.GetFramePadding()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func item_inner_spacing()
	/// @returns width, height
	lib.CreateFunction(tab, "item_inner_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetItemInnerSpacing()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func item_spacing()
	/// @returns width, height
	lib.CreateFunction(tab, "item_spacing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := g.GetItemSpacing()

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func mouse_pos_xy()
	/// @returns x, y
	lib.CreateFunction(tab, "mouse_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetMousePos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func mouse_pos()
	/// @returns image.point
	lib.CreateFunction(tab, "mouse_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetMousePos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func window_padding()
	/// @returns x, y
	lib.CreateFunction(tab, "window_padding",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			x, y := g.GetWindowPadding()

			state.Push(golua.LNumber(x))
			state.Push(golua.LNumber(y))
			return 2
		})

	/// @func is_item_active()
	/// @returns bool
	lib.CreateFunction(tab, "is_item_active",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemActive()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_item_clicked()
	/// @arg button
	/// @returns bool
	lib.CreateFunction(tab, "is_item_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_item_hovered()
	/// @returns bool
	lib.CreateFunction(tab, "is_item_hovered",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsItemHovered()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_down()
	/// @arg key
	/// @returns bool
	lib.CreateFunction(tab, "is_key_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyDown(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_pressed()
	/// @arg key
	/// @returns bool
	lib.CreateFunction(tab, "is_key_pressed",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyPressed(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_key_released()
	/// @arg key
	/// @returns bool
	lib.CreateFunction(tab, "is_key_released",
		[]lua.Arg{
			{Type: lua.INT, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsKeyReleased(g.Key(args["key"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_clicked()
	/// @arg button
	/// @returns bool
	lib.CreateFunction(tab, "is_mouse_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_double_clicked()
	/// @arg button
	/// @returns bool
	lib.CreateFunction(tab, "is_mouse_double_clicked",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseDoubleClicked(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_down()
	/// @arg button
	/// @returns bool
	lib.CreateFunction(tab, "is_mouse_down",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseDown(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_mouse_released()
	/// @arg button
	/// @returns bool
	lib.CreateFunction(tab, "is_mouse_released",
		[]lua.Arg{
			{Type: lua.INT, Name: "button"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsMouseReleased(g.MouseButton(args["button"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_appearing()
	/// @returns bool
	lib.CreateFunction(tab, "is_window_appearing",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowAppearing()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_collapsed()
	/// @returns bool
	lib.CreateFunction(tab, "is_window_collapsed",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowCollapsed()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_collapsed()
	/// @returns bool
	lib.CreateFunction(tab, "is_window_collapsed",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowCollapsed()

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_focused()
	/// @arg flags
	/// @returns bool
	lib.CreateFunction(tab, "is_window_focused",
		[]lua.Arg{
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowFocused(g.FocusedFlags(args["flags"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func is_window_hovered()
	/// @arg flags
	/// @returns bool
	lib.CreateFunction(tab, "is_window_hovered",
		[]lua.Arg{
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			active := g.IsWindowHovered(g.HoveredFlags(args["flags"].(int)))

			state.Push(golua.LBool(active))
			return 1
		})

	/// @func open_url()
	/// @arg url
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

	/// @func pop_style_color_v()
	/// @arg count
	lib.CreateFunction(tab, "pop_style_color_v",
		[]lua.Arg{
			{Type: lua.INT, Name: "count"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PopStyleColorV(args["count"].(int))
			return 0
		})

	/// @func pop_style_v()
	/// @arg count
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

	/// @func push_button_text_align()
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "push_button_text_align",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushButtonTextAlign(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_clip_rect()
	/// @arg min
	/// @arg max
	/// @rag intersect
	lib.CreateFunction(tab, "push_clip_rect",
		[]lua.Arg{
			{Type: lua.ANY, Name: "min"},
			{Type: lua.ANY, Name: "max"},
			{Type: lua.BOOL, Name: "intersect"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			min := imageutil.TableToPoint(state, args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(state, args["max"].(*golua.LTable))
			g.PushClipRect(min, max, args["intersect"].(bool))
			return 0
		})

	/// @func push_color_button()
	/// @arg color
	lib.CreateFunction(tab, "push_color_button",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButton(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_button_active()
	/// @arg color
	lib.CreateFunction(tab, "push_color_button_active",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButtonActive(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_button_hovered()
	/// @arg color
	lib.CreateFunction(tab, "push_color_button_hovered",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorButtonHovered(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_frame_bg()
	/// @arg color
	lib.CreateFunction(tab, "push_color_frame_bg",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorFrameBg(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_text()
	/// @arg color
	lib.CreateFunction(tab, "push_color_text",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorText(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_text_disabled()
	/// @arg color
	lib.CreateFunction(tab, "push_color_text_disabled",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorTextDisabled(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_color_window_bg()
	/// @arg color
	lib.CreateFunction(tab, "push_color_window_bg",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushColorWindowBg(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_font()
	/// @arg fontref
	/// @returns bool
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

	/// @func push_frame_padding()
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "push_frame_padding",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushFramePadding(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_item_spacing()
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "push_item_spacing",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushItemSpacing(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_item_width()
	/// @arg width
	lib.CreateFunction(tab, "push_item_width",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushItemWidth(float32(args["width"].(float64)))
			return 0
		})

	/// @func push_selectable_text_align()
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "push_selectable_text_align",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushSelectableTextAlign(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func push_style_color()
	/// @arg id
	/// @arg color
	lib.CreateFunction(tab, "push_style_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushStyleColor(g.StyleColorID(args["id"].(int)), imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))
			return 0
		})

	/// @func push_text_wrap_pos()
	lib.CreateFunction(tab, "push_text_wrap_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.PushTextWrapPos()
			return 0
		})

	/// @func push_window_padding()
	/// @arg width
	/// @arg height
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

	/// @func cursor_pos_set_xy()
	/// @arg x
	/// @arg y
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

	/// @func cursor_pos_set()
	/// @arg point
	lib.CreateFunction(tab, "cursor_pos_set",
		[]lua.Arg{
			{Type: lua.ANY, Name: "point"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorPos(imageutil.TableToPoint(state, args["point"].(*golua.LTable)))
			return 0
		})

	/// @func cursor_screeen_pos_set_xy()
	/// @arg x
	/// @arg y
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

	/// @func cursor_screen_pos_set()
	/// @arg point
	lib.CreateFunction(tab, "cursor_screen_pos_set",
		[]lua.Arg{
			{Type: lua.ANY, Name: "point"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetCursorScreenPos(imageutil.TableToPoint(state, args["point"].(*golua.LTable)))
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

	/// @func keyboard_focus_here_v()
	/// @arg i
	lib.CreateFunction(tab, "keyboard_focus_here_v",
		[]lua.Arg{
			{Type: lua.INT, Name: "i"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetKeyboardFocusHereV(args["i"].(int))
			return 0
		})

	/// @func mouse_cursor_set()
	/// @arg cursor
	lib.CreateFunction(tab, "mouse_cursor_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "cursor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetMouseCursor(g.MouseCursorType(args["cursor"].(int)))
			return 0
		})

	/// @func next_window_pos_set()
	/// @arg x
	/// @arg y
	lib.CreateFunction(tab, "next_window_pos_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "x"},
			{Type: lua.INT, Name: "y"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetNextWindowPos(float32(args["x"].(float64)), float32(args["y"].(float64)))
			return 0
		})

	/// @func next_window_size_set()
	/// @arg width
	/// @arg height
	lib.CreateFunction(tab, "next_window_size_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "width"},
			{Type: lua.FLOAT, Name: "height"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			g.SetNextWindowSize(float32(args["width"].(float64)), float32(args["height"].(float64)))
			return 0
		})

	/// @func next_window_size_v_set()
	/// @arg width
	/// @arg height
	/// @arg cond
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

	/// @func color_to_uint32()
	/// @arg color
	/// @returns number representation of color
	/// @desc
	/// returns the uint32 to lua as a float64
	lib.CreateFunction(tab, "color_to_uint32",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cuint := g.ColorToUint(imageutil.TableToRGBA(state, args["color"].(*golua.LTable)))

			state.Push(golua.LNumber(cuint))
			return 1
		})

	/// @func uint32_to_color()
	/// @arg ucolor
	/// @returns color
	lib.CreateFunction(tab, "uint32_to_color",
		[]lua.Arg{
			{Type: lua.INT, Name: "ucolor"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.UintToColor(uint32(args["ucolor"].(int)))
			state.Push(imageutil.RGBAToTable(state, c))
			return 1
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

	/// @func wg_combo_preview()
	/// @arg text
	/// @arg items - []string
	/// @arg i32ref
	/// @returns widget
	/// @desc
	/// Same as wg_combo but sets preview to the selected value in items.
	lib.CreateFunction(tab, "wg_combo_preview",
		[]lua.Arg{
			{Type: lua.STRING, Name: "text"},
			{Type: lua.ANY, Name: "items"},
			{Type: lua.INT, Name: "i32ref"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			text := args["text"].(string)
			items := args["items"].(golua.LValue)
			i32ref := args["i32ref"].(int)

			sref, err := r.CR_REF.Item(i32ref)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
			}
			selected := sref.Value.(*int32)
			preview := state.GetTable(items.(*golua.LTable), golua.LNumber(*selected+1)).(golua.LString)

			t := comboTable(state, text, string(preview), items, i32ref)

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

	/// @func wg_tree_node()
	/// @arg label
	/// @returns widget
	lib.CreateFunction(tab, "wg_tree_node",
		[]lua.Arg{
			{Type: lua.STRING, Name: "label"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := treeNodeTable(state, args["label"].(string))

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

	/// @func wg_splitter()
	/// @arg direction
	/// @arg f32ref
	/// @returns widget
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

	/// @func wg_stack()
	/// @arg visible
	/// @arg widgets
	/// @returns widget
	lib.CreateFunction(tab, "wg_stack",
		[]lua.Arg{
			{Type: lua.INT, Name: "visible"},
			{Type: lua.ANY, Name: "widgets"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := stackTable(state, args["visible"].(int), args["widgets"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func wg_align()
	/// @arg at
	/// @returns widget
	lib.CreateFunction(tab, "wg_align",
		[]lua.Arg{
			{Type: lua.INT, Name: "at"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := alignTable(state, args["at"].(int))

			state.Push(t)
			return 1
		})

	/// @func wg_msg_box()
	/// @arg title
	/// @arg content
	/// @returns msg box widget
	/// @desc
	/// prepare_msg_box() must be called once a loop when using a msg box.
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

	/// @func wg_button_invisible()
	/// @returns widget
	lib.CreateFunction(tab, "wg_button_invisible",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonInvisibleTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_button_image()
	/// @arg id
	/// @returns widget
	lib.CreateFunction(tab, "wg_button_image",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonImageTable(state, args["id"].(int), false)

			state.Push(t)
			return 1
		})

	/// @func wg_button_image_sync()
	/// @arg id
	/// @returns widget
	/// @desc
	/// Note: this does not wait for the image to be ready or idle,
	/// if the image is not loaded it will dislay an empy image
	/// May look weird if the image is also being processed while displayed here.
	lib.CreateFunction(tab, "wg_button_image_sync",
		[]lua.Arg{
			{Type: lua.INT, Name: "id"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := buttonImageTable(state, args["id"].(int), true)

			state.Push(t)
			return 1
		})

	/// @func wg_style()
	/// @returns widget
	lib.CreateFunction(tab, "wg_style",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := styleTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_custom()
	/// @arg builder
	/// @returns widget
	lib.CreateFunction(tab, "wg_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "builder"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := customTable(state, args["builder"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func wg_event()
	/// @returns widget
	lib.CreateFunction(tab, "wg_event",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := eventHandlerTable(state)

			state.Push(t)
			return 1
		})

	/// @func wg_plot()
	/// @arg title
	/// @returns widget
	lib.CreateFunction(tab, "wg_plot",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotTable(state, args["title"].(string))

			state.Push(t)
			return 1
		})

	/// @func pt_bar_h()
	/// @arg title
	/// @arg data
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_bar_h",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarHTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_bar()
	/// @arg title
	/// @arg data
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_bar",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotBarTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line()
	/// @arg title
	/// @arg data
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_line",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_line_xy()
	/// @arg title
	/// @arg xdata
	/// @arg ydata
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_line_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "xdata"},
			{Type: lua.ANY, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotLineXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_pie_chart()
	/// @arg labels
	/// @arg data
	/// @arg x
	/// @arg y
	/// @arg radius
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_pie_chart",
		[]lua.Arg{
			{Type: lua.ANY, Name: "labels"},
			{Type: lua.ANY, Name: "data"},
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

	/// @func pt_scatter()
	/// @arg title
	/// @arg data
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_scatter",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "data"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterTable(state, args["title"].(string), args["data"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_scatter_xy()
	/// @arg title
	/// @arg xdata
	/// @arg ydata
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_scatter_xy",
		[]lua.Arg{
			{Type: lua.STRING, Name: "title"},
			{Type: lua.ANY, Name: "xdata"},
			{Type: lua.ANY, Name: "ydata"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotScatterXYTable(state, args["title"].(string), args["xdata"].(golua.LValue), args["ydata"].(golua.LValue))

			state.Push(t)
			return 1
		})

	/// @func pt_custom()
	/// @arg builder
	/// @returns plot widget
	lib.CreateFunction(tab, "pt_custom",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "builder"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := plotCustomTable(state, args["builder"].(*golua.LFunction))

			state.Push(t)
			return 1
		})

	/// @func wg_css_tag()
	/// @arg tag
	/// @returns widget
	lib.CreateFunction(tab, "wg_css_tag",
		[]lua.Arg{
			{Type: lua.STRING, Name: "tag"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := cssTagTable(state, args["tag"].(string))

			state.Push(t)
			return 1
		})

	/// @func cursor_screen_pos_xy()
	/// @returns x, y
	lib.CreateFunction(tab, "cursor_screen_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorScreenPos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func cursor_screen_pos()
	/// @returns image point
	lib.CreateFunction(tab, "cursor_screen_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorScreenPos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func cursor_pos_xy()
	/// @returns x, y
	lib.CreateFunction(tab, "cursor_pos_xy",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorPos()

			state.Push(golua.LNumber(p.X))
			state.Push(golua.LNumber(p.Y))
			return 2
		})

	/// @func cursor_pos()
	/// @returns image point
	lib.CreateFunction(tab, "cursor_pos",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			p := g.GetCursorPos()

			state.Push(imageutil.PointToTable(state, p))
			return 1
		})

	/// @func canvas_bezier_cubic()
	/// @arg pos0
	/// @arg cp0
	/// @arg cp1
	/// @arg pos1
	/// @arg color
	/// @arg thickness
	/// @arg segments
	lib.CreateFunction(tab, "canvas_bezier_cubic",
		[]lua.Arg{
			{Type: lua.ANY, Name: "pos0"},
			{Type: lua.ANY, Name: "cp0"},
			{Type: lua.ANY, Name: "cp1"},
			{Type: lua.ANY, Name: "pos1"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			pos0 := imageutil.TableToPoint(state, args["pos0"].(*golua.LTable))
			cp0 := imageutil.TableToPoint(state, args["cp0"].(*golua.LTable))
			cp1 := imageutil.TableToPoint(state, args["cp1"].(*golua.LTable))
			pos1 := imageutil.TableToPoint(state, args["pos1"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)
			segments := args["segments"].(int)

			c.AddBezierCubic(pos0, cp0, cp1, pos1, col, float32(thickness), int32(segments))
			return 0
		})

	/// @func canvas_circle()
	/// @arg center
	/// @arg radius
	/// @arg color
	/// @arg segments
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_circle",
		[]lua.Arg{
			{Type: lua.ANY, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.INT, Name: "segments"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(state, args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			segments := args["segments"].(int)
			thickness := args["thickness"].(float64)

			c.AddCircle(center, float32(radius), col, int32(segments), float32(thickness))
			return 0
		})

	/// @func canvas_circle_filled()
	/// @arg center
	/// @arg radius
	/// @arg color
	lib.CreateFunction(tab, "canvas_circle_filled",
		[]lua.Arg{
			{Type: lua.ANY, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(state, args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))

			c.AddCircleFilled(center, float32(radius), col)
			return 0
		})

	/// @func canvas_line()
	/// @arg p1
	/// @arg p2
	/// @arg color
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_line",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddLine(p1, p2, col, float32(thickness))
			return 0
		})

	/// @func canvas_quad()
	/// @arg p1
	/// @arg p2
	/// @arg p3
	/// @arg p4
	/// @arg color
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_quad",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "p3"},
			{Type: lua.ANY, Name: "p4"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(state, args["p3"].(*golua.LTable))
			p4 := imageutil.TableToPoint(state, args["p4"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddQuad(p1, p2, p3, p4, col, float32(thickness))
			return 0
		})

	/// @func canvas_quad_filled()
	/// @arg p1
	/// @arg p2
	/// @arg p3
	/// @arg p4
	/// @arg color
	lib.CreateFunction(tab, "canvas_quad_filled",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "p3"},
			{Type: lua.ANY, Name: "p4"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(state, args["p3"].(*golua.LTable))
			p4 := imageutil.TableToPoint(state, args["p4"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))

			c.AddQuadFilled(p1, p2, p3, p4, col)
			return 0
		})

	/// @func canvas_rect()
	/// @arg min
	/// @arg max
	/// @arg color
	/// @arg rounding
	/// @arg flags
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_rect",
		[]lua.Arg{
			{Type: lua.ANY, Name: "min"},
			{Type: lua.ANY, Name: "max"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "rounding"},
			{Type: lua.INT, Name: "flags"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			min := imageutil.TableToPoint(state, args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(state, args["max"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			rounding := args["rounding"].(float64)
			flags := args["flags"].(int)
			thickness := args["thickness"].(float64)

			c.AddRect(min, max, col, float32(rounding), g.DrawFlags(flags), float32(thickness))
			return 0
		})

	/// @func canvas_rect_filled()
	/// @arg min
	/// @arg max
	/// @arg color
	/// @arg rounding
	/// @arg flags
	lib.CreateFunction(tab, "canvas_rect_filled",
		[]lua.Arg{
			{Type: lua.ANY, Name: "min"},
			{Type: lua.ANY, Name: "max"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "rounding"},
			{Type: lua.INT, Name: "flags"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			min := imageutil.TableToPoint(state, args["min"].(*golua.LTable))
			max := imageutil.TableToPoint(state, args["max"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			rounding := args["rounding"].(float64)
			flags := args["flags"].(int)

			c.AddRectFilled(min, max, col, float32(rounding), g.DrawFlags(flags))
			return 0
		})

	/// @func canvas_text()
	/// @arg pos
	/// @arg color
	/// @arg text
	lib.CreateFunction(tab, "canvas_text",
		[]lua.Arg{
			{Type: lua.ANY, Name: "pos"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.STRING, Name: "text"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			pos := imageutil.TableToPoint(state, args["pos"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			text := args["text"].(string)

			c.AddText(pos, col, text)
			return 0
		})

	/// @func canvas_triangle()
	/// @arg p1
	/// @arg p2
	/// @arg p3
	/// @arg color
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_triangle",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "p3"},
			{Type: lua.ANY, Name: "color"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(state, args["p3"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			thickness := args["thickness"].(float64)

			c.AddTriangle(p1, p2, p3, col, float32(thickness))
			return 0
		})

	/// @func canvas_triangle_filled()
	/// @arg p1
	/// @arg p2
	/// @arg p3
	/// @arg color
	lib.CreateFunction(tab, "canvas_triangle_filled",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "p3"},
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(state, args["p3"].(*golua.LTable))
			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))

			c.AddTriangleFilled(p1, p2, p3, col)
			return 0
		})

	/// @func canvas_path_arc_to()
	/// @arg center
	/// @arg radius
	/// @arg min
	/// @arg max
	/// @arg segments
	lib.CreateFunction(tab, "canvas_path_arc_to",
		[]lua.Arg{
			{Type: lua.ANY, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.FLOAT, Name: "min"},
			{Type: lua.FLOAT, Name: "max"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(state, args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			min := args["min"].(float64)
			max := args["max"].(float64)
			segments := args["segments"].(int)

			c.PathArcTo(center, float32(radius), float32(min), float32(max), int32(segments))
			return 0
		})

	/// @func canvas_path_arc_to_fast()
	/// @arg center
	/// @arg radius
	/// @arg min
	/// @arg max
	/// @arg segments
	lib.CreateFunction(tab, "canvas_path_arc_to_fast",
		[]lua.Arg{
			{Type: lua.ANY, Name: "center"},
			{Type: lua.FLOAT, Name: "radius"},
			{Type: lua.INT, Name: "min"},
			{Type: lua.INT, Name: "max"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			center := imageutil.TableToPoint(state, args["center"].(*golua.LTable))
			radius := args["radius"].(float64)
			min := args["min"].(int)
			max := args["max"].(int)

			c.PathArcToFast(center, float32(radius), int32(min), int32(max))
			return 0
		})

	/// @func canvas_path_bezier_cubic_to()
	/// @arg p1
	/// @arg p2
	/// @arg p3
	/// @arg segments
	lib.CreateFunction(tab, "canvas_path_bezier_cubic_to",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
			{Type: lua.ANY, Name: "p2"},
			{Type: lua.ANY, Name: "p3"},
			{Type: lua.INT, Name: "segments"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))
			p2 := imageutil.TableToPoint(state, args["p2"].(*golua.LTable))
			p3 := imageutil.TableToPoint(state, args["p3"].(*golua.LTable))
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

	/// @func canvas_fill_convex()
	/// @arg color
	lib.CreateFunction(tab, "canvas_path_fill_convex",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))

			c.PathFillConvex(col)
			return 0
		})

	/// @func canvas_path_line_to()
	/// @arg p1
	/// @arg segments
	lib.CreateFunction(tab, "canvas_path_line_to",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))

			c.PathLineTo(p1)
			return 0
		})

	/// @func canvas_path_line_to_merge_duplicate()
	/// @arg p1
	/// @arg segments
	lib.CreateFunction(tab, "canvas_path_line_to_merge_duplicate",
		[]lua.Arg{
			{Type: lua.ANY, Name: "p1"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			p1 := imageutil.TableToPoint(state, args["p1"].(*golua.LTable))

			c.PathLineToMergeDuplicate(p1)
			return 0
		})

	/// @func canvas_path_stroke()
	/// @arg color
	/// @arg flags
	/// @arg thickness
	lib.CreateFunction(tab, "canvas_path_stroke",
		[]lua.Arg{
			{Type: lua.ANY, Name: "color"},
			{Type: lua.INT, Name: "flags"},
			{Type: lua.FLOAT, Name: "thickness"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			c := g.GetCanvas()

			col := imageutil.TableToRGBA(state, args["color"].(*golua.LTable))
			flags := args["flags"].(int)
			thickness := args["thickness"].(float64)

			c.PathStroke(col, g.DrawFlags(flags), float32(thickness))
			return 0
		})

	/// @func fontatlas_add_font()
	/// @arg name
	/// @arg size
	/// @returns fontref, ok
	/// @desc
	/// fontref will be nil if ok is false
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

	/// @func fontatlas_default_font_strings()
	/// @returns []string
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

	/// @func fontatlas_default_fonts()
	/// @returns []fontref
	/// @desc
	/// Take note that this creates an array of refs,
	/// refs are only cleared at the end of a workflow,
	/// or with ref.del / ref.del_many
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

	/// @func fontatlas_register_string()
	/// @arg str
	/// @returns str
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

	/// @func fontatlas_register_string_ref()
	/// @arg stringref
	/// @returns stringref
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

	/// @func fontatlas_register_string_many()
	/// @arg []str
	/// @returns []str
	lib.CreateFunction(tab, "fontatlas_register_string_many",
		[]lua.Arg{
			{Type: lua.ANY, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			strList := args["str"].(*golua.LTable)
			strSlice := []string{}

			for i := range strList.Len() {
				v := state.GetTable(strList, golua.LNumber(i+1)).(golua.LString)
				strSlice = append(strSlice, string(v))
			}

			atlas.RegisterStringSlice(strSlice)

			state.Push(strList)
			return 1
		})

	/// @func fontatlas_set_default_font()
	/// @arg name
	/// @arg size
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

	/// @func fontatlas_set_default_font_size()
	/// @arg size
	lib.CreateFunction(tab, "fontatlas_set_default_font_size",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "size"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			atlas := g.Context.FontAtlas
			atlas.SetDefaultFontSize(float32(args["size"].(float64)))
			return 0
		})

	/// @func font_set_size()
	/// @arg fontref
	/// @arg size
	/// @returns fontref
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

	/// @func font_string()
	/// @arg fontref
	/// @returns string
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

	/// @constants Alignment
	/// @const ALIGN_LEFT
	/// @const ALIGN_CENTER
	/// @const ALIGN_RIGHT
	r.State.SetTable(tab, golua.LString("ALIGN_LEFT"), golua.LNumber(ALIGN_LEFT))
	r.State.SetTable(tab, golua.LString("ALIGN_CENTER"), golua.LNumber(ALIGN_CENTER))
	r.State.SetTable(tab, golua.LString("ALIGN_RIGHT"), golua.LNumber(ALIGN_RIGHT))

	/// @constants MSG Box Buttons
	/// @const MSGBOXBUTTONS_YESNO
	/// @const MSGBOXBUTTONS_OKCANCEL
	/// @const MSGBOXBUTTONS_OK
	r.State.SetTable(tab, golua.LString("MSGBOXBUTTONS_YESNO"), golua.LNumber(MSGBOXBUTTONS_YESNO))
	r.State.SetTable(tab, golua.LString("MSGBOXBUTTONS_OKCANCEL"), golua.LNumber(MSGBOXBUTTONS_OKCANCEL))
	r.State.SetTable(tab, golua.LString("MSGBOXBUTTONS_OK"), golua.LNumber(MSGBOXBUTTONS_OK))

	/// @constants Color IDs
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
	r.State.SetTable(tab, golua.LString("COLID_TEXT"), golua.LNumber(COLID_TEXT))
	r.State.SetTable(tab, golua.LString("COLID_TEXTDISABLED"), golua.LNumber(COLID_TEXTDISABLED))
	r.State.SetTable(tab, golua.LString("COLID_WINDOWBG"), golua.LNumber(COLID_WINDOWBG))
	r.State.SetTable(tab, golua.LString("COLID_CHILDBG"), golua.LNumber(COLID_CHILDBG))
	r.State.SetTable(tab, golua.LString("COLID_POPUPBG"), golua.LNumber(COLID_POPUPBG))
	r.State.SetTable(tab, golua.LString("COLID_BORDER"), golua.LNumber(COLID_BORDER))
	r.State.SetTable(tab, golua.LString("COLID_BORDERSHADOW"), golua.LNumber(COLID_BORDERSHADOW))
	r.State.SetTable(tab, golua.LString("COLID_FRAMEBG"), golua.LNumber(COLID_FRAMEBG))
	r.State.SetTable(tab, golua.LString("COLID_FRAMEBGHOVERED"), golua.LNumber(COLID_FRAMEBGHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_FRAMEBGACTIVE"), golua.LNumber(COLID_FRAMEBGACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_TITLEBG"), golua.LNumber(COLID_TITLEBG))
	r.State.SetTable(tab, golua.LString("COLID_TITLEBGACTIVE"), golua.LNumber(COLID_TITLEBGACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_TITLEBGCOLLAPSED"), golua.LNumber(COLID_TITLEBGCOLLAPSED))
	r.State.SetTable(tab, golua.LString("COLID_MENUBARBG"), golua.LNumber(COLID_MENUBARBG))
	r.State.SetTable(tab, golua.LString("COLID_SCROLLBARBG"), golua.LNumber(COLID_SCROLLBARBG))
	r.State.SetTable(tab, golua.LString("COLID_SCROLLBARGRAB"), golua.LNumber(COLID_SCROLLBARGRAB))
	r.State.SetTable(tab, golua.LString("COLID_SCROLLBARGRABHOVERED"), golua.LNumber(COLID_SCROLLBARGRABHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_SCROLLBARGRABACTIVE"), golua.LNumber(COLID_SCROLLBARGRABACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_CHECKMARK"), golua.LNumber(COLID_CHECKMARK))
	r.State.SetTable(tab, golua.LString("COLID_SLIDERGRAB"), golua.LNumber(COLID_SLIDERGRAB))
	r.State.SetTable(tab, golua.LString("COLID_SLIDERGRABACTIVE"), golua.LNumber(COLID_SLIDERGRABACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_BUTTON"), golua.LNumber(COLID_BUTTON))
	r.State.SetTable(tab, golua.LString("COLID_BUTTONHOVERED"), golua.LNumber(COLID_BUTTONHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_BUTTONACTIVE"), golua.LNumber(COLID_BUTTONACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_HEADER"), golua.LNumber(COLID_HEADER))
	r.State.SetTable(tab, golua.LString("COLID_HEADERHOVERED"), golua.LNumber(COLID_HEADERHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_HEADERACTIVE"), golua.LNumber(COLID_HEADERACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_SEPARATOR"), golua.LNumber(COLID_SEPARATOR))
	r.State.SetTable(tab, golua.LString("COLID_SEPARATORHOVERED"), golua.LNumber(COLID_SEPARATORHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_SEPARATORACTIVE"), golua.LNumber(COLID_SEPARATORACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_RESIZEGRIP"), golua.LNumber(COLID_RESIZEGRIP))
	r.State.SetTable(tab, golua.LString("COLID_RESIZEGRIPHOVERED"), golua.LNumber(COLID_RESIZEGRIPHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_RESIZEGRIPACTIVE"), golua.LNumber(COLID_RESIZEGRIPACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_TAB"), golua.LNumber(COLID_TAB))
	r.State.SetTable(tab, golua.LString("COLID_TABHOVERED"), golua.LNumber(COLID_TABHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_TABACTIVE"), golua.LNumber(COLID_TABACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_TABUNFOCUSED"), golua.LNumber(COLID_TABUNFOCUSED))
	r.State.SetTable(tab, golua.LString("COLID_TABUNFOCUSEDACTIVE"), golua.LNumber(COLID_TABUNFOCUSEDACTIVE))
	r.State.SetTable(tab, golua.LString("COLID_DOCKINGPREVIEW"), golua.LNumber(COLID_DOCKINGPREVIEW))
	r.State.SetTable(tab, golua.LString("COLID_DOCKINGEMPTYBG"), golua.LNumber(COLID_DOCKINGEMPTYBG))
	r.State.SetTable(tab, golua.LString("COLID_PLOTLINES"), golua.LNumber(COLID_PLOTLINES))
	r.State.SetTable(tab, golua.LString("COLID_PLOTLINESHOVERED"), golua.LNumber(COLID_PLOTLINESHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_PLOTHISTOGRAM"), golua.LNumber(COLID_PLOTHISTOGRAM))
	r.State.SetTable(tab, golua.LString("COLID_PLOTHISTOGRAMHOVERED"), golua.LNumber(COLID_PLOTHISTOGRAMHOVERED))
	r.State.SetTable(tab, golua.LString("COLID_TABLEHEADERBG"), golua.LNumber(COLID_TABLEHEADERBG))
	r.State.SetTable(tab, golua.LString("COLID_TABLEBORDERSTRONG"), golua.LNumber(COLID_TABLEBORDERSTRONG))
	r.State.SetTable(tab, golua.LString("COLID_TABLEBORDERLIGHT"), golua.LNumber(COLID_TABLEBORDERLIGHT))
	r.State.SetTable(tab, golua.LString("COLID_TABLEROWBG"), golua.LNumber(COLID_TABLEROWBG))
	r.State.SetTable(tab, golua.LString("COLID_TABLEROWBGALT"), golua.LNumber(COLID_TABLEROWBGALT))
	r.State.SetTable(tab, golua.LString("COLID_TEXTSELECTEDBG"), golua.LNumber(COLID_TEXTSELECTEDBG))
	r.State.SetTable(tab, golua.LString("COLID_DRAGDROPTARGET"), golua.LNumber(COLID_DRAGDROPTARGET))
	r.State.SetTable(tab, golua.LString("COLID_NAVHIGHLIGHT"), golua.LNumber(COLID_NAVHIGHLIGHT))
	r.State.SetTable(tab, golua.LString("COLID_NAVWINDOWINGHIGHLIGHT"), golua.LNumber(COLID_NAVWINDOWINGHIGHLIGHT))
	r.State.SetTable(tab, golua.LString("COLID_NAVWINDOWINGDIMBG"), golua.LNumber(COLID_NAVWINDOWINGDIMBG))
	r.State.SetTable(tab, golua.LString("COLID_MODALWINDOWDIMBG"), golua.LNumber(COLID_MODALWINDOWDIMBG))
	r.State.SetTable(tab, golua.LString("COLID_COUNT"), golua.LNumber(COLID_COUNT))

	/// @constants Style Var
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
	r.State.SetTable(tab, golua.LString("STYLEVAR_ALPHA"), golua.LNumber(STYLEVAR_ALPHA))
	r.State.SetTable(tab, golua.LString("STYLEVAR_DISABLEDALPHA"), golua.LNumber(STYLEVAR_DISABLEDALPHA))
	r.State.SetTable(tab, golua.LString("STYLEVAR_WINDOWPADDING"), golua.LNumber(STYLEVAR_WINDOWPADDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_WINDOWROUNDING"), golua.LNumber(STYLEVAR_WINDOWROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_WINDOWBORDERSIZE"), golua.LNumber(STYLEVAR_WINDOWBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_WINDOWMINSIZE"), golua.LNumber(STYLEVAR_WINDOWMINSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_WINDOWTITLEALIGN"), golua.LNumber(STYLEVAR_WINDOWTITLEALIGN))
	r.State.SetTable(tab, golua.LString("STYLEVAR_CHILDROUNDING"), golua.LNumber(STYLEVAR_CHILDROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_CHILDBORDERSIZE"), golua.LNumber(STYLEVAR_CHILDBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_POPUPROUNDING"), golua.LNumber(STYLEVAR_POPUPROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_POPUPBORDERSIZE"), golua.LNumber(STYLEVAR_POPUPBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_FRAMEPADDING"), golua.LNumber(STYLEVAR_FRAMEPADDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_FRAMEROUNDING"), golua.LNumber(STYLEVAR_FRAMEROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_FRAMEBORDERSIZE"), golua.LNumber(STYLEVAR_FRAMEBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_ITEMSPACING"), golua.LNumber(STYLEVAR_ITEMSPACING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_ITEMINNERSPACING"), golua.LNumber(STYLEVAR_ITEMINNERSPACING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_INDENTSPACING"), golua.LNumber(STYLEVAR_INDENTSPACING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_CELLPADDING"), golua.LNumber(STYLEVAR_CELLPADDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SCROLLBARSIZE"), golua.LNumber(STYLEVAR_SCROLLBARSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SCROLLBARROUNDING"), golua.LNumber(STYLEVAR_SCROLLBARROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_GRABMINSIZE"), golua.LNumber(STYLEVAR_GRABMINSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_GRABROUNDING"), golua.LNumber(STYLEVAR_GRABROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_TABROUNDING"), golua.LNumber(STYLEVAR_TABROUNDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_TABBARBORDERSIZE"), golua.LNumber(STYLEVAR_TABBARBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_BUTTONTEXTALIGN"), golua.LNumber(STYLEVAR_BUTTONTEXTALIGN))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SELECTABLETEXTALIGN"), golua.LNumber(STYLEVAR_SELECTABLETEXTALIGN))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SEPARATORTEXTBORDERSIZE"), golua.LNumber(STYLEVAR_SEPARATORTEXTBORDERSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SEPARATORTEXTALIGN"), golua.LNumber(STYLEVAR_SEPARATORTEXTALIGN))
	r.State.SetTable(tab, golua.LString("STYLEVAR_SEPARATORTEXTPADDING"), golua.LNumber(STYLEVAR_SEPARATORTEXTPADDING))
	r.State.SetTable(tab, golua.LString("STYLEVAR_DOCKINGSEPARATORSIZE"), golua.LNumber(STYLEVAR_DOCKINGSEPARATORSIZE))
	r.State.SetTable(tab, golua.LString("STYLEVAR_COUNT"), golua.LNumber(STYLEVAR_COUNT))

	/// @constants Keys
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
	r.State.SetTable(tab, golua.LString("KEY_NONE"), golua.LNumber(KEY_NONE))
	r.State.SetTable(tab, golua.LString("KEY_TAB"), golua.LNumber(KEY_TAB))
	r.State.SetTable(tab, golua.LString("KEY_LEFTARROW"), golua.LNumber(KEY_LEFTARROW))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTARROW"), golua.LNumber(KEY_RIGHTARROW))
	r.State.SetTable(tab, golua.LString("KEY_UPARROW"), golua.LNumber(KEY_UPARROW))
	r.State.SetTable(tab, golua.LString("KEY_DOWNARROW"), golua.LNumber(KEY_DOWNARROW))
	r.State.SetTable(tab, golua.LString("KEY_PAGEUP"), golua.LNumber(KEY_PAGEUP))
	r.State.SetTable(tab, golua.LString("KEY_PAGEDOWN"), golua.LNumber(KEY_PAGEDOWN))
	r.State.SetTable(tab, golua.LString("KEY_HOME"), golua.LNumber(KEY_HOME))
	r.State.SetTable(tab, golua.LString("KEY_END"), golua.LNumber(KEY_END))
	r.State.SetTable(tab, golua.LString("KEY_INSERT"), golua.LNumber(KEY_INSERT))
	r.State.SetTable(tab, golua.LString("KEY_DELETE"), golua.LNumber(KEY_DELETE))
	r.State.SetTable(tab, golua.LString("KEY_BACKSPACE"), golua.LNumber(KEY_BACKSPACE))
	r.State.SetTable(tab, golua.LString("KEY_SPACE"), golua.LNumber(KEY_SPACE))
	r.State.SetTable(tab, golua.LString("KEY_ENTER"), golua.LNumber(KEY_ENTER))
	r.State.SetTable(tab, golua.LString("KEY_ESCAPE"), golua.LNumber(KEY_ESCAPE))
	r.State.SetTable(tab, golua.LString("KEY_LEFTCTRL"), golua.LNumber(KEY_LEFTCTRL))
	r.State.SetTable(tab, golua.LString("KEY_LEFTSHIFT"), golua.LNumber(KEY_LEFTSHIFT))
	r.State.SetTable(tab, golua.LString("KEY_LEFTALT"), golua.LNumber(KEY_LEFTALT))
	r.State.SetTable(tab, golua.LString("KEY_LEFTSUPER"), golua.LNumber(KEY_LEFTSUPER))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTCTRL"), golua.LNumber(KEY_RIGHTCTRL))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTSHIFT"), golua.LNumber(KEY_RIGHTSHIFT))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTALT"), golua.LNumber(KEY_RIGHTALT))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTSUPER"), golua.LNumber(KEY_RIGHTSUPER))
	r.State.SetTable(tab, golua.LString("KEY_MENU"), golua.LNumber(KEY_MENU))
	r.State.SetTable(tab, golua.LString("KEY_0"), golua.LNumber(KEY_0))
	r.State.SetTable(tab, golua.LString("KEY_1"), golua.LNumber(KEY_1))
	r.State.SetTable(tab, golua.LString("KEY_2"), golua.LNumber(KEY_2))
	r.State.SetTable(tab, golua.LString("KEY_3"), golua.LNumber(KEY_3))
	r.State.SetTable(tab, golua.LString("KEY_4"), golua.LNumber(KEY_4))
	r.State.SetTable(tab, golua.LString("KEY_5"), golua.LNumber(KEY_5))
	r.State.SetTable(tab, golua.LString("KEY_6"), golua.LNumber(KEY_6))
	r.State.SetTable(tab, golua.LString("KEY_7"), golua.LNumber(KEY_7))
	r.State.SetTable(tab, golua.LString("KEY_8"), golua.LNumber(KEY_8))
	r.State.SetTable(tab, golua.LString("KEY_9"), golua.LNumber(KEY_9))
	r.State.SetTable(tab, golua.LString("KEY_A"), golua.LNumber(KEY_A))
	r.State.SetTable(tab, golua.LString("KEY_B"), golua.LNumber(KEY_B))
	r.State.SetTable(tab, golua.LString("KEY_C"), golua.LNumber(KEY_C))
	r.State.SetTable(tab, golua.LString("KEY_D"), golua.LNumber(KEY_D))
	r.State.SetTable(tab, golua.LString("KEY_E"), golua.LNumber(KEY_E))
	r.State.SetTable(tab, golua.LString("KEY_F"), golua.LNumber(KEY_F))
	r.State.SetTable(tab, golua.LString("KEY_G"), golua.LNumber(KEY_G))
	r.State.SetTable(tab, golua.LString("KEY_H"), golua.LNumber(KEY_H))
	r.State.SetTable(tab, golua.LString("KEY_I"), golua.LNumber(KEY_I))
	r.State.SetTable(tab, golua.LString("KEY_J"), golua.LNumber(KEY_J))
	r.State.SetTable(tab, golua.LString("KEY_K"), golua.LNumber(KEY_K))
	r.State.SetTable(tab, golua.LString("KEY_L"), golua.LNumber(KEY_L))
	r.State.SetTable(tab, golua.LString("KEY_M"), golua.LNumber(KEY_M))
	r.State.SetTable(tab, golua.LString("KEY_N"), golua.LNumber(KEY_N))
	r.State.SetTable(tab, golua.LString("KEY_O"), golua.LNumber(KEY_O))
	r.State.SetTable(tab, golua.LString("KEY_P"), golua.LNumber(KEY_P))
	r.State.SetTable(tab, golua.LString("KEY_Q"), golua.LNumber(KEY_Q))
	r.State.SetTable(tab, golua.LString("KEY_R"), golua.LNumber(KEY_R))
	r.State.SetTable(tab, golua.LString("KEY_S"), golua.LNumber(KEY_S))
	r.State.SetTable(tab, golua.LString("KEY_T"), golua.LNumber(KEY_T))
	r.State.SetTable(tab, golua.LString("KEY_U"), golua.LNumber(KEY_U))
	r.State.SetTable(tab, golua.LString("KEY_V"), golua.LNumber(KEY_V))
	r.State.SetTable(tab, golua.LString("KEY_W"), golua.LNumber(KEY_W))
	r.State.SetTable(tab, golua.LString("KEY_X"), golua.LNumber(KEY_X))
	r.State.SetTable(tab, golua.LString("KEY_Y"), golua.LNumber(KEY_Y))
	r.State.SetTable(tab, golua.LString("KEY_Z"), golua.LNumber(KEY_Z))
	r.State.SetTable(tab, golua.LString("KEY_F1"), golua.LNumber(KEY_F1))
	r.State.SetTable(tab, golua.LString("KEY_F2"), golua.LNumber(KEY_F2))
	r.State.SetTable(tab, golua.LString("KEY_F3"), golua.LNumber(KEY_F3))
	r.State.SetTable(tab, golua.LString("KEY_F4"), golua.LNumber(KEY_F4))
	r.State.SetTable(tab, golua.LString("KEY_F5"), golua.LNumber(KEY_F5))
	r.State.SetTable(tab, golua.LString("KEY_F6"), golua.LNumber(KEY_F6))
	r.State.SetTable(tab, golua.LString("KEY_F7"), golua.LNumber(KEY_F7))
	r.State.SetTable(tab, golua.LString("KEY_F8"), golua.LNumber(KEY_F8))
	r.State.SetTable(tab, golua.LString("KEY_F9"), golua.LNumber(KEY_F9))
	r.State.SetTable(tab, golua.LString("KEY_F10"), golua.LNumber(KEY_F10))
	r.State.SetTable(tab, golua.LString("KEY_F11"), golua.LNumber(KEY_F11))
	r.State.SetTable(tab, golua.LString("KEY_F12"), golua.LNumber(KEY_F12))
	r.State.SetTable(tab, golua.LString("KEY_F13"), golua.LNumber(KEY_F13))
	r.State.SetTable(tab, golua.LString("KEY_F14"), golua.LNumber(KEY_F14))
	r.State.SetTable(tab, golua.LString("KEY_F15"), golua.LNumber(KEY_F15))
	r.State.SetTable(tab, golua.LString("KEY_F16"), golua.LNumber(KEY_F16))
	r.State.SetTable(tab, golua.LString("KEY_F17"), golua.LNumber(KEY_F17))
	r.State.SetTable(tab, golua.LString("KEY_F18"), golua.LNumber(KEY_F18))
	r.State.SetTable(tab, golua.LString("KEY_F19"), golua.LNumber(KEY_F19))
	r.State.SetTable(tab, golua.LString("KEY_F20"), golua.LNumber(KEY_F20))
	r.State.SetTable(tab, golua.LString("KEY_F21"), golua.LNumber(KEY_F21))
	r.State.SetTable(tab, golua.LString("KEY_F22"), golua.LNumber(KEY_F22))
	r.State.SetTable(tab, golua.LString("KEY_F23"), golua.LNumber(KEY_F23))
	r.State.SetTable(tab, golua.LString("KEY_F24"), golua.LNumber(KEY_F24))
	r.State.SetTable(tab, golua.LString("KEY_APOSTROPHE"), golua.LNumber(KEY_APOSTROPHE))
	r.State.SetTable(tab, golua.LString("KEY_COMMA"), golua.LNumber(KEY_COMMA))
	r.State.SetTable(tab, golua.LString("KEY_MINUS"), golua.LNumber(KEY_MINUS))
	r.State.SetTable(tab, golua.LString("KEY_PERIOD"), golua.LNumber(KEY_PERIOD))
	r.State.SetTable(tab, golua.LString("KEY_SLASH"), golua.LNumber(KEY_SLASH))
	r.State.SetTable(tab, golua.LString("KEY_SEMICOLON"), golua.LNumber(KEY_SEMICOLON))
	r.State.SetTable(tab, golua.LString("KEY_EQUAL"), golua.LNumber(KEY_EQUAL))
	r.State.SetTable(tab, golua.LString("KEY_LEFTBRACKET"), golua.LNumber(KEY_LEFTBRACKET))
	r.State.SetTable(tab, golua.LString("KEY_BACKSLASH"), golua.LNumber(KEY_BACKSLASH))
	r.State.SetTable(tab, golua.LString("KEY_RIGHTBRACKET"), golua.LNumber(KEY_RIGHTBRACKET))
	r.State.SetTable(tab, golua.LString("KEY_GRAVEACCENT"), golua.LNumber(KEY_GRAVEACCENT))
	r.State.SetTable(tab, golua.LString("KEY_CAPSLOCK"), golua.LNumber(KEY_CAPSLOCK))
	r.State.SetTable(tab, golua.LString("KEY_SCROLLLOCK"), golua.LNumber(KEY_SCROLLLOCK))
	r.State.SetTable(tab, golua.LString("KEY_NUMLOCK"), golua.LNumber(KEY_NUMLOCK))
	r.State.SetTable(tab, golua.LString("KEY_PRINTSCREEN"), golua.LNumber(KEY_PRINTSCREEN))
	r.State.SetTable(tab, golua.LString("KEY_PAUSE"), golua.LNumber(KEY_PAUSE))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD0"), golua.LNumber(KEY_KEYPAD0))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD1"), golua.LNumber(KEY_KEYPAD1))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD2"), golua.LNumber(KEY_KEYPAD2))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD3"), golua.LNumber(KEY_KEYPAD3))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD4"), golua.LNumber(KEY_KEYPAD4))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD5"), golua.LNumber(KEY_KEYPAD5))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD6"), golua.LNumber(KEY_KEYPAD6))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD7"), golua.LNumber(KEY_KEYPAD7))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD8"), golua.LNumber(KEY_KEYPAD8))
	r.State.SetTable(tab, golua.LString("KEY_KEYPAD9"), golua.LNumber(KEY_KEYPAD9))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADDECIMAL"), golua.LNumber(KEY_KEYPADDECIMAL))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADDIVIDE"), golua.LNumber(KEY_KEYPADDIVIDE))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADMULTIPLY"), golua.LNumber(KEY_KEYPADMULTIPLY))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADSUBTRACT"), golua.LNumber(KEY_KEYPADSUBTRACT))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADADD"), golua.LNumber(KEY_KEYPADADD))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADENTER"), golua.LNumber(KEY_KEYPADENTER))
	r.State.SetTable(tab, golua.LString("KEY_KEYPADEQUAL"), golua.LNumber(KEY_KEYPADEQUAL))
	r.State.SetTable(tab, golua.LString("KEY_APPBACK"), golua.LNumber(KEY_APPBACK))
	r.State.SetTable(tab, golua.LString("KEY_APPFORWARD"), golua.LNumber(KEY_APPFORWARD))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADSTART"), golua.LNumber(KEY_GAMEPADSTART))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADBACK"), golua.LNumber(KEY_GAMEPADBACK))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADFACELEFT"), golua.LNumber(KEY_GAMEPADFACELEFT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADFACERIGHT"), golua.LNumber(KEY_GAMEPADFACERIGHT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADFACEUP"), golua.LNumber(KEY_GAMEPADFACEUP))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADFACEDOWN"), golua.LNumber(KEY_GAMEPADFACEDOWN))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADDPADLEFT"), golua.LNumber(KEY_GAMEPADDPADLEFT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADDPADRIGHT"), golua.LNumber(KEY_GAMEPADDPADRIGHT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADDPADUP"), golua.LNumber(KEY_GAMEPADDPADUP))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADDPADDOWN"), golua.LNumber(KEY_GAMEPADDPADDOWN))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADL1"), golua.LNumber(KEY_GAMEPADL1))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADR1"), golua.LNumber(KEY_GAMEPADR1))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADL2"), golua.LNumber(KEY_GAMEPADL2))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADR2"), golua.LNumber(KEY_GAMEPADR2))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADL3"), golua.LNumber(KEY_GAMEPADL3))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADR3"), golua.LNumber(KEY_GAMEPADR3))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADLSTICKLEFT"), golua.LNumber(KEY_GAMEPADLSTICKLEFT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADLSTICKRIGHT"), golua.LNumber(KEY_GAMEPADLSTICKRIGHT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADLSTICKUP"), golua.LNumber(KEY_GAMEPADLSTICKUP))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADLSTICKDOWN"), golua.LNumber(KEY_GAMEPADLSTICKDOWN))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADRSTICKLEFT"), golua.LNumber(KEY_GAMEPADRSTICKLEFT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADRSTICKRIGHT"), golua.LNumber(KEY_GAMEPADRSTICKRIGHT))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADRSTICKUP"), golua.LNumber(KEY_GAMEPADRSTICKUP))
	r.State.SetTable(tab, golua.LString("KEY_GAMEPADRSTICKDOWN"), golua.LNumber(KEY_GAMEPADRSTICKDOWN))
	r.State.SetTable(tab, golua.LString("KEY_MOUSELEFT"), golua.LNumber(KEY_MOUSELEFT))
	r.State.SetTable(tab, golua.LString("KEY_MOUSERIGHT"), golua.LNumber(KEY_MOUSERIGHT))
	r.State.SetTable(tab, golua.LString("KEY_MOUSEMIDDLE"), golua.LNumber(KEY_MOUSEMIDDLE))
	r.State.SetTable(tab, golua.LString("KEY_MOUSEX1"), golua.LNumber(KEY_MOUSEX1))
	r.State.SetTable(tab, golua.LString("KEY_MOUSEX2"), golua.LNumber(KEY_MOUSEX2))
	r.State.SetTable(tab, golua.LString("KEY_MOUSEWHEELX"), golua.LNumber(KEY_MOUSEWHEELX))
	r.State.SetTable(tab, golua.LString("KEY_MOUSEWHEELY"), golua.LNumber(KEY_MOUSEWHEELY))
	r.State.SetTable(tab, golua.LString("KEY_RESERVEDFORMODCTRL"), golua.LNumber(KEY_RESERVEDFORMODCTRL))
	r.State.SetTable(tab, golua.LString("KEY_RESERVEDFORMODSHIFT"), golua.LNumber(KEY_RESERVEDFORMODSHIFT))
	r.State.SetTable(tab, golua.LString("KEY_RESERVEDFORMODALT"), golua.LNumber(KEY_RESERVEDFORMODALT))
	r.State.SetTable(tab, golua.LString("KEY_RESERVEDFORMODSUPER"), golua.LNumber(KEY_RESERVEDFORMODSUPER))
	r.State.SetTable(tab, golua.LString("KEY_COUNT"), golua.LNumber(KEY_COUNT))
	r.State.SetTable(tab, golua.LString("KEY_MODNONE"), golua.LNumber(KEY_MODNONE))
	r.State.SetTable(tab, golua.LString("KEY_MODCTRL"), golua.LNumber(KEY_MODCTRL))
	r.State.SetTable(tab, golua.LString("KEY_MODSHIFT"), golua.LNumber(KEY_MODSHIFT))
	r.State.SetTable(tab, golua.LString("KEY_MODALT"), golua.LNumber(KEY_MODALT))
	r.State.SetTable(tab, golua.LString("KEY_MODSUPER"), golua.LNumber(KEY_MODSUPER))
	r.State.SetTable(tab, golua.LString("KEY_MODSHORTCUT"), golua.LNumber(KEY_MODSHORTCUT))
	r.State.SetTable(tab, golua.LString("KEY_MODMASK"), golua.LNumber(KEY_MODMASK))
	r.State.SetTable(tab, golua.LString("KEY_NAMEDKEYBEGIN"), golua.LNumber(KEY_NAMEDKEYBEGIN))
	r.State.SetTable(tab, golua.LString("KEY_NAMEDKEYEND"), golua.LNumber(KEY_NAMEDKEYEND))
	r.State.SetTable(tab, golua.LString("KEY_NAMEDKEYCOUNT"), golua.LNumber(KEY_NAMEDKEYCOUNT))
	r.State.SetTable(tab, golua.LString("KEY_KEYSDATASIZE"), golua.LNumber(KEY_KEYSDATASIZE))
	r.State.SetTable(tab, golua.LString("KEY_KEYSDATAOFFSET"), golua.LNumber(KEY_KEYSDATAOFFSET))

	/// @constants Exec Conditions
	/// @const COND_NONE
	/// @const COND_ALWAYS
	/// @const COND_ONCE
	/// @const COND_FIRSTUSEEVER
	/// @const COND_APPEARING
	r.State.SetTable(tab, golua.LString("COND_NONE"), golua.LNumber(COND_NONE))
	r.State.SetTable(tab, golua.LString("COND_ALWAYS"), golua.LNumber(COND_ALWAYS))
	r.State.SetTable(tab, golua.LString("COND_ONCE"), golua.LNumber(COND_ONCE))
	r.State.SetTable(tab, golua.LString("COND_FIRSTUSEEVER"), golua.LNumber(COND_FIRSTUSEEVER))
	r.State.SetTable(tab, golua.LString("COND_APPEARING"), golua.LNumber(COND_APPEARING))

	/// @constants Plot Flags
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
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NONE"), golua.LNumber(FLAGPLOT_NONE))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOTITLE"), golua.LNumber(FLAGPLOT_NOTITLE))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOLEGEND"), golua.LNumber(FLAGPLOT_NOLEGEND))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOMOUSETEXT"), golua.LNumber(FLAGPLOT_NOMOUSETEXT))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOINPUTS"), golua.LNumber(FLAGPLOT_NOINPUTS))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOMENUS"), golua.LNumber(FLAGPLOT_NOMENUS))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOBOXSELECT"), golua.LNumber(FLAGPLOT_NOBOXSELECT))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_NOFRAME"), golua.LNumber(FLAGPLOT_NOFRAME))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_EQUAL"), golua.LNumber(FLAGPLOT_EQUAL))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_CROSSHAIRS"), golua.LNumber(FLAGPLOT_CROSSHAIRS))
	r.State.SetTable(tab, golua.LString("FLAGPLOT_CANVASONLY"), golua.LNumber(FLAGPLOT_CANVASONLY))

	/// @constants Plot Axis
	/// @const PLOTAXIS_X1
	/// @const PLOTAXIS_X2
	/// @const PLOTAXIS_X3
	/// @const PLOTAXIS_Y1
	/// @const PLOTAXIS_Y2
	/// @const PLOTAXIS_Y3
	/// @const PLOTAXIS_COUNT
	r.State.SetTable(tab, golua.LString("PLOTAXIS_X1"), golua.LNumber(PLOTAXIS_X1))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_X2"), golua.LNumber(PLOTAXIS_X2))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_X3"), golua.LNumber(PLOTAXIS_X3))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_Y1"), golua.LNumber(PLOTAXIS_Y1))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_Y2"), golua.LNumber(PLOTAXIS_Y2))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_Y3"), golua.LNumber(PLOTAXIS_Y3))
	r.State.SetTable(tab, golua.LString("PLOTAXIS_COUNT"), golua.LNumber(PLOTAXIS_COUNT))

	/// @constants Plot Axis Flags
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
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NONE"), golua.LNumber(FLAGPLOTAXIS_NONE))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOLABEL"), golua.LNumber(FLAGPLOTAXIS_NOLABEL))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOGRIDLINES"), golua.LNumber(FLAGPLOTAXIS_NOGRIDLINES))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOTICKMARKS"), golua.LNumber(FLAGPLOTAXIS_NOTICKMARKS))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOTICKLABELS"), golua.LNumber(FLAGPLOTAXIS_NOTICKLABELS))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOINITIALFIT"), golua.LNumber(FLAGPLOTAXIS_NOINITIALFIT))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOMENUS"), golua.LNumber(FLAGPLOTAXIS_NOMENUS))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOSIDESWITCH"), golua.LNumber(FLAGPLOTAXIS_NOSIDESWITCH))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NOHIGHLIGHT"), golua.LNumber(FLAGPLOTAXIS_NOHIGHLIGHT))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_OPPOSITE"), golua.LNumber(FLAGPLOTAXIS_OPPOSITE))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_FOREGROUND"), golua.LNumber(FLAGPLOTAXIS_FOREGROUND))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_INVERT"), golua.LNumber(FLAGPLOTAXIS_INVERT))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_AUTOFIT"), golua.LNumber(FLAGPLOTAXIS_AUTOFIT))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_RANGEFIT"), golua.LNumber(FLAGPLOTAXIS_RANGEFIT))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_PANSTRETCH"), golua.LNumber(FLAGPLOTAXIS_PANSTRETCH))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_LOCKMIN"), golua.LNumber(FLAGPLOTAXIS_LOCKMIN))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_LOCKMAX"), golua.LNumber(FLAGPLOTAXIS_LOCKMAX))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_LOCK"), golua.LNumber(FLAGPLOTAXIS_LOCK))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_NODECORATIONS"), golua.LNumber(FLAGPLOTAXIS_NODECORATIONS))
	r.State.SetTable(tab, golua.LString("FLAGPLOTAXIS_AUXDEFAULT"), golua.LNumber(FLAGPLOTAXIS_AUXDEFAULT))

	/// @constants Plot Y Axis
	/// @const PLOTYAXIS_LEFT
	/// @const PLOTYAXIS_FIRSTONRIGHT
	/// @const PLOTYAXIS_SECONDONRIGHT
	r.State.SetTable(tab, golua.LString("PLOTYAXIS_LEFT"), golua.LNumber(PLOTYAXIS_LEFT))
	r.State.SetTable(tab, golua.LString("PLOTYAXIS_FIRSTONRIGHT"), golua.LNumber(PLOTYAXIS_FIRSTONRIGHT))
	r.State.SetTable(tab, golua.LString("PLOTYAXIS_SECONDONRIGHT"), golua.LNumber(PLOTYAXIS_SECONDONRIGHT))

	/// @constants Draw Flags
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
	r.State.SetTable(tab, golua.LString("FLAGDRAW_NONE"), golua.LNumber(FLAGDRAW_NONE))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_CLOSED"), golua.LNumber(FLAGDRAW_CLOSED))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSTOPLEFT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSTOPLEFT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSTOPRIGHT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSTOPRIGHT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSBOTTOMLEFT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOMLEFT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSNONE"), golua.LNumber(FLAGDRAW_ROUNDCORNERSNONE))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSTOP"), golua.LNumber(FLAGDRAW_ROUNDCORNERSTOP))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSBOTTOM"), golua.LNumber(FLAGDRAW_ROUNDCORNERSBOTTOM))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSLEFT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSLEFT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSRIGHT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSRIGHT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSALL"), golua.LNumber(FLAGDRAW_ROUNDCORNERSALL))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSDEFAULT"), golua.LNumber(FLAGDRAW_ROUNDCORNERSDEFAULT))
	r.State.SetTable(tab, golua.LString("FLAGDRAW_ROUNDCORNERSMASK"), golua.LNumber(FLAGDRAW_ROUNDCORNERSMASK))

	/// @constants Focus Flags
	/// @const FLAGFOCUS_NONE
	/// @const FLAGFOCUS_CHILDWINDOWS
	/// @const FLAGFOCUS_ROOTWINDOW
	/// @const FLAGFOCUS_ANYWINDOW
	/// @const FLAGFOCUS_NOPOPUPHIERARCHY
	/// @const FLAGFOCUS_DOCKHIERARCHY
	/// @const FLAGFOCUS_ROOTANDCHILDWINDOWS
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_NONE"), golua.LNumber(FLAGFOCUS_NONE))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_CHILDWINDOWS"), golua.LNumber(FLAGFOCUS_CHILDWINDOWS))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_ROOTWINDOW"), golua.LNumber(FLAGFOCUS_ROOTWINDOW))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_ANYWINDOW"), golua.LNumber(FLAGFOCUS_ANYWINDOW))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_NOPOPUPHIERARCHY"), golua.LNumber(FLAGFOCUS_NOPOPUPHIERARCHY))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_DOCKHIERARCHY"), golua.LNumber(FLAGFOCUS_DOCKHIERARCHY))
	r.State.SetTable(tab, golua.LString("FLAGFOCUS_ROOTANDCHILDWINDOWS"), golua.LNumber(FLAGFOCUS_ROOTANDCHILDWINDOWS))

	/// @constants Hover Flags
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
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_NONE"), golua.LNumber(FLAGHOVERED_NONE))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_CHILDWINDOWS"), golua.LNumber(FLAGHOVERED_CHILDWINDOWS))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ROOTWINDOW"), golua.LNumber(FLAGHOVERED_ROOTWINDOW))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ANYWINDOW"), golua.LNumber(FLAGHOVERED_ANYWINDOW))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_NOPOPUPHIERARCHY"), golua.LNumber(FLAGHOVERED_NOPOPUPHIERARCHY))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_DOCKHIERARCHY"), golua.LNumber(FLAGHOVERED_DOCKHIERARCHY))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP"), golua.LNumber(FLAGHOVERED_ALLOWWHENBLOCKEDBYPOPUP))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM"), golua.LNumber(FLAGHOVERED_ALLOWWHENBLOCKEDBYACTIVEITEM))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM"), golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPEDBYITEM))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW"), golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPEDBYWINDOW))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENDISABLED"), golua.LNumber(FLAGHOVERED_ALLOWWHENDISABLED))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_NONAVOVERRIDE"), golua.LNumber(FLAGHOVERED_NONAVOVERRIDE))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ALLOWWHENOVERLAPPED"), golua.LNumber(FLAGHOVERED_ALLOWWHENOVERLAPPED))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_RECTONLY"), golua.LNumber(FLAGHOVERED_RECTONLY))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_ROOTANDCHILDWINDOWS"), golua.LNumber(FLAGHOVERED_ROOTANDCHILDWINDOWS))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_FORTOOLTIP"), golua.LNumber(FLAGHOVERED_FORTOOLTIP))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_STATIONARY"), golua.LNumber(FLAGHOVERED_STATIONARY))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_DELAYNONE"), golua.LNumber(FLAGHOVERED_DELAYNONE))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_DELAYSHORT"), golua.LNumber(FLAGHOVERED_DELAYSHORT))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_DELAYNORMAL"), golua.LNumber(FLAGHOVERED_DELAYNORMAL))
	r.State.SetTable(tab, golua.LString("FLAGHOVERED_NOSHAREDDELAY"), golua.LNumber(FLAGHOVERED_NOSHAREDDELAY))

	/// @constants Mouse Cursors
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
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_NONE"), golua.LNumber(MOUSECURSOR_NONE))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_ARROW"), golua.LNumber(MOUSECURSOR_ARROW))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_TEXTINPUT"), golua.LNumber(MOUSECURSOR_TEXTINPUT))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_RESIZEALL"), golua.LNumber(MOUSECURSOR_RESIZEALL))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_RESIZENS"), golua.LNumber(MOUSECURSOR_RESIZENS))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_RESIZEEW"), golua.LNumber(MOUSECURSOR_RESIZEEW))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_RESIZENESW"), golua.LNumber(MOUSECURSOR_RESIZENESW))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_RESIZENWSE"), golua.LNumber(MOUSECURSOR_RESIZENWSE))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_HAND"), golua.LNumber(MOUSECURSOR_HAND))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_NOTALLOWED"), golua.LNumber(MOUSECURSOR_NOTALLOWED))
	r.State.SetTable(tab, golua.LString("MOUSECURSOR_COUNT"), golua.LNumber(MOUSECURSOR_COUNT))

	/// @constants Actions
	/// @const ACTION_RELEASE
	/// @const ACTION_PRESS
	/// @const ACTION_REPEAT
	r.State.SetTable(tab, golua.LString("ACTION_RELEASE"), golua.LNumber(ACTION_RELEASE))
	r.State.SetTable(tab, golua.LString("ACTION_PRESS"), golua.LNumber(ACTION_PRESS))
	r.State.SetTable(tab, golua.LString("ACTION_REPEAT"), golua.LNumber(ACTION_REPEAT))
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

var buildList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget{}
var plotList = map[string]func(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget{}

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
	/// @struct wg_label
	/// @prop type
	/// @prop label
	/// @method wrapped(bool)
	/// @method font(fontref)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_LABEL))
	state.SetTable(t, golua.LString("label"), golua.LString(text))
	state.SetTable(t, golua.LString("__wrapped"), golua.LNil)
	state.SetTable(t, golua.LString("__font"), golua.LNil)

	tableBuilderFunc(state, t, "wrapped", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__wrapped"), golua.LBool(v))
	})

	tableBuilderFunc(state, t, "font", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__font"), v)
	})

	return t
}

func labelBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	l := g.Label(state.GetTable(t, golua.LString("label")).String())

	wrapped := state.GetTable(t, golua.LString("__wrapped"))
	if wrapped.Type() == golua.LTBool {
		l.Wrapped(bool(wrapped.(golua.LBool)))
	}

	fontref := state.GetTable(t, golua.LString("__font"))
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
	/// @struct wg_button
	/// @prop type
	/// @prop label
	/// @method disabled(bool)
	/// @method size(width, height)
	/// @method on_click(callback)

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
	/// @struct wg_dummy
	/// @prop type
	/// @prop width
	/// @prop height

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
	/// @struct wg_separator
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SEPARATOR))

	return t
}

func separatorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	s := g.Separator()

	return s
}

func bulletTextTable(state *golua.LState, text string) *golua.LTable {
	/// @struct wg_bullet_text
	/// @prop type
	/// @prop text

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
	/// @struct wg_bullet
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BULLET))

	return t
}

func bulletBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Bullet()

	return b
}

func checkboxTable(state *golua.LState, text string, boolref int) *golua.LTable {
	/// @struct wg_checkbox
	/// @prop type
	/// @prop text
	/// @prop boolref
	/// @method on_change(callback(bool, boolref))

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
	/// @struct wg_child
	/// @prop type
	/// @method border(bool)
	/// @method size(width, height)
	/// @method layout([]widgets)
	/// @method flags(flags)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CHILD))
	state.SetTable(t, golua.LString("__border"), golua.LNil)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)

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

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
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

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		c.Flags(g.WindowFlags(flags.(golua.LNumber)))
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}

func colorEditTable(state *golua.LState, text string, colorref int) *golua.LTable {
	/// @struct wg_color_edit
	/// @prop type
	/// @prop label
	/// @prop colorref
	/// @method size(width)
	/// @method on_change(callback(color, colorref))
	/// @method flags(flags)

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
	/// @struct wg_column
	/// @prop type
	/// @prop widgets

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
	/// @struct wg_row
	/// @prop type
	/// @prop widgets

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
	/// @struct wg_combo_custom
	/// @prop type
	/// @prop text
	/// @prop preview
	/// @method size(width)
	/// @method layout([]widgets)
	/// @method flags(flags)

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
	/// @struct wg_combo
	/// @prop type
	/// @prop text
	/// @prop preview
	/// @prop items
	/// @prop i32ref
	/// @method size(width)
	/// @method on_change(callback(int, i32ref))
	/// @method flags(flags)

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
	/// @struct wg_condition
	/// @prop type
	/// @prop condition
	/// @prop layoutIf
	/// @prop layoutElse

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
	/// @struct wg_context_menu
	/// @prop type
	/// @method mouse_button(button)
	/// @method layout([]widgets)

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
	/// @struct wg_date_picker
	/// @prop type
	/// @prop id
	/// @prop timeref
	/// @method on_change(callback(time, timeref))
	/// @method format(format)
	/// @method size(width)
	/// @method start_of_week(day)
	/// @method translation(label, value)

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
	/// @struct wg_drag_int
	/// @prop type
	/// @prop text
	/// @prop i32ref
	/// @prop minvalue
	/// @prop maxvalue
	/// @method speed(speed)
	/// @method format(format)

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
	/// @struct wg_input_float
	/// @prop type
	/// @prop f32ref
	/// @method size(width)
	/// @method on_change(callback(float, f32ref))
	/// @method format(format)
	/// @method flags(flags)
	/// @method label(label)
	/// @method step_size(stepsize)
	/// @method step_size_fast(stepsize)

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
	/// @struct wg_input_int
	/// @prop type
	/// @prop i32ref
	/// @method size(width)
	/// @method on_change(callback(int, i32ref))
	/// @method flags(flags)
	/// @method label(label)
	/// @method step_size(stepsize)
	/// @method step_size_fast(stepsize)

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
	/// @struct wg_input_text
	/// @prop type
	/// @prop strref
	/// @method size(width)
	/// @method flags(flags)
	/// @method label(label)
	/// @method autocomplete([]string)
	/// @method callback(callback(string, strref))
	/// @method hint(hint)

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
	/// @struct wg_input_multiline_text
	/// @prop type
	/// @prop strref
	/// @method size(width, height)
	/// @method on_change(callback(string, strref))
	/// @method flags(flags)
	/// @method label(label)
	/// @method callback(callback(string, strref))
	/// @method autoscroll_to_bottom(bool)

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
	/// @struct wg_progress_bar
	/// @prop type
	/// @prop fraction
	/// @method overlay(label)
	/// @method size(width, height)

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
	/// @struct wg_progress_indicator
	/// @prop type
	/// @prop label
	/// @prop width
	/// @prop height
	/// @prop radius

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
	/// @struct wg_spacing
	/// @prop type

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SPACING))

	return t
}

func spacingBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.Spacing()

	return b
}

func buttonSmallTable(state *golua.LState, text string) *golua.LTable {
	/// @struct wg_button_small
	/// @prop type
	/// @prop label
	/// @method on_click(callback())

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
	/// @struct wg_button_radio
	/// @prop type
	/// @prop label
	/// @prop active
	/// @method on_change(callback())

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
	/// @struct wg_image_url
	/// @prop type
	/// @prop url
	/// @method on_click(callback())
	/// @method size(width, height)
	/// @method timeout(timeout)
	/// @method layout_for_failure([]widgets)
	/// @method layout_for_loading([]widgets)
	/// @method on_failure(callback())
	/// @method on_ready(callback())

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
	/// @struct wg_image
	/// @prop type
	/// @prop image
	/// @prop sync
	/// @method on_click(callback())
	/// @method size(width, height)

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
	/// @struct wg_list_box
	/// @prop type
	/// @prop items
	/// @method on_change(callback(index))
	/// @method border(bool)
	/// @method context_menu([]widgets)
	/// @method on_double_click(callback(index))
	/// @method on_menu(callback(index, menu))
	/// @method selected_index(index)
	/// @method size(width, height)

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
	/// @struct wg_list_clipper
	/// @prop type
	/// @method layout([]widgets)

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
	/// @struct wg_main_menu_bar
	/// @prop type
	/// @method layout([]widgets)

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
	/// @struct wg_menu_bar
	/// @prop type
	/// @method layout([]widgets)

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
	/// @struct wg_menu_item
	/// @prop type
	/// @prop label
	/// @method enabled(bool)
	/// @method on_click(callback())
	/// @method selected(bool)
	/// @method shortcut(string)

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
	/// @struct wg_menu
	/// @prop type
	/// @prop label
	/// @method enabled(bool)
	/// @method layout([]widgets)

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
	/// @struct wg_selectable
	/// @prop type
	/// @prop label
	/// @method on_click(callback())
	/// @method on_double_click(callback())
	/// @method selected(bool)
	/// @method size(width, height)
	/// @method flags(flags)

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
	/// @struct wg_slider_float
	/// @prop type
	/// @prop f32ref
	/// @prop min
	/// @prop max
	/// @method on_change(callback(value, f32ref))
	/// @method label(string)
	/// @method format(string)
	/// @method size(width)

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
	/// @struct wg_slider_int
	/// @prop type
	/// @prop i32ref
	/// @prop min
	/// @prop max
	/// @method on_change(callback(value, i32ref))
	/// @method label(string)
	/// @method format(string)
	/// @method size(width)

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
	/// @struct wg_vslider_int
	/// @prop type
	/// @prop i32ref
	/// @prop min
	/// @prop max
	/// @method on_change(callback(value, i32ref))
	/// @method label(string)
	/// @method format(string)
	/// @method size(width, height)
	/// @method flags(flags)

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
	/// @struct wg_tab_bar
	/// @prop type
	/// @method flags(flags)
	/// @method tab_items([]wg_tab_item)

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
	/// @struct wg_tab_item
	/// @prop type
	/// @prop label
	/// @method flags(flags)
	/// @method is_open(bool)
	/// @method layout([]widgets)

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
	/// @struct wg_tooltip
	/// @prop type
	/// @prop tip
	/// @method layout([]widgets)

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
	/// @struct wg_table_column
	/// @prop type
	/// @prop label
	/// @method flags(flags)
	/// @method inner_width_or_weight(width)
	/// @desc
	/// only used in table widget columns

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

func tableColumnBuild(state *golua.LState, t *golua.LTable) *g.TableColumnWidget {
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
	/// @struct wg_table_row
	/// @prop type
	/// @prop widgets
	/// @method flags(flags)
	/// @method bg_color(color)
	/// @method min_height(height)
	/// @desc
	/// only used in table widget rows

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
	/// @struct wg_table
	/// @prop type
	/// @method flags(flags)
	/// @method fast_mode(bool)
	/// @method size(width, height)
	/// @method columns([]wg_table_column)
	/// @method rows([]wg_table_row)
	/// @method inner_width(width)
	/// @method freeze(col, row) - can be called multiple times

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
			i := tableColumnBuild(state, w)
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
	/// @struct wg_button_arrow
	/// @prop type
	/// @prop dir
	/// @method on_click(callback)

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

func treeNodeTable(state *golua.LState, label string) *golua.LTable {
	/// @struct wg_tree_node
	/// @prop type
	/// @prop label
	/// @method flags(flags)
	/// @method layout([]widgets)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_TREE_NODE))
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

func treeNodeBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
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
	/// @struct wg_tree_table_row
	/// @prop type
	/// @prop label
	/// @prop widgets
	/// @method flags(flags)
	/// @method children([]wg_tree_table_row)
	/// @desc
	/// only used in tree table widget rows

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
	/// @struct wg_tree_table
	/// @prop type
	/// @method flags(flags)
	/// @method size(width, height)
	/// @method columns([]wg_table_column)
	/// @method rows([]wg_tree_table_row)
	/// @method freeze(col, row) - can be called multiple times

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
			i := tableColumnBuild(state, w)
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
	/// @struct wg_window
	/// @prop type
	/// @prop single
	/// @prop menubar
	/// @prop label
	/// @method flags(flags)
	/// @method size(width, height)
	/// @method pos(x, y)
	/// @method is_open(bool)
	/// @method bring_to_front()
	/// @method ready(callback(state_window))
	/// @method register_keyboard_shortcuts([]shortcut)
	/// @method layout([]widgets)

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
	state.SetTable(t, golua.LString("__shortcuts"), golua.LNil)

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

	tableBuilderFunc(state, t, "register_keyboard_shortcuts", func(state *golua.LState, t *golua.LTable) {
		st := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__shortcuts"), st)
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

	/// @struct state_window
	/// @method current_position() x, y
	/// @method current_size() width, height
	/// @method has_focus() bool
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

	shortcuts := state.GetTable(t, golua.LString("__shortcuts"))
	if shortcuts.Type() == golua.LTTable {
		stList := []g.WindowShortcut{}
		st := shortcuts.(*golua.LTable)
		for i := range st.Len() {
			s := state.GetTable(st, golua.LNumber(i+1)).(*golua.LTable)

			key := state.GetTable(s, golua.LString("key")).(golua.LNumber)
			mod := state.GetTable(s, golua.LString("mod")).(golua.LNumber)
			callback := state.GetTable(s, golua.LString("callback"))

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

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		w.Layout(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return w
}

func popupModalTable(state *golua.LState, label string) *golua.LTable {
	/// @struct wg_popup_modal
	/// @prop type
	/// @prop label
	/// @method flags(flags)
	/// @method is_open(bool)
	/// @method layout([]widgets)

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
	/// @struct wg_popup
	/// @prop type
	/// @prop label
	/// @method flags(flags)
	/// @method layout([]widgets)

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
	/// @struct wg_split_layout
	/// @prop type
	/// @prop direction
	/// @prop floatref
	/// @prop layout1
	/// @prop layout2
	/// @method border(bool)

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

func splitterTable(state *golua.LState, direction, floatref int) *golua.LTable {
	/// @struct wg_splitter
	/// @prop type
	/// @prop direction
	/// @prop floatref
	/// @method size(width, height)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_SPLITTER))
	state.SetTable(t, golua.LString("direction"), golua.LNumber(direction))
	state.SetTable(t, golua.LString("floatref"), golua.LNumber(floatref))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	return t
}

func splitterBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	direction := state.GetTable(t, golua.LString("direction")).(golua.LNumber)

	floatref := state.GetTable(t, golua.LString("floatref"))
	ref, err := r.CR_REF.Item(int(floatref.(golua.LNumber)))
	if err != nil {
		state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
	}
	pos := ref.Value.(*float32)

	s := g.Splitter(g.SplitDirection(direction), pos)

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		s.Size(float32(width.(golua.LNumber)), float32(height.(golua.LNumber)))
	}

	return s
}

func stackTable(state *golua.LState, visible int, widgets golua.LValue) *golua.LTable {
	/// @struct wg_stack
	/// @prop type
	/// @prop visible
	/// @prop widgets

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_STACK))
	state.SetTable(t, golua.LString("visible"), golua.LNumber(visible))
	state.SetTable(t, golua.LString("widgets"), widgets)

	return t
}

func stackBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	visible := state.GetTable(t, golua.LString("visible")).(golua.LNumber)

	var widgets []g.Widget
	wid1 := state.GetTable(t, golua.LString("widgets"))
	if wid1.Type() == golua.LTTable {
		widgets = layoutBuild(r, state, parseWidgets(parseTable(wid1.(*golua.LTable), state), state, lg), lg)
	}

	s := g.Stack(int32(visible), widgets...)

	return s
}

func alignTable(state *golua.LState, at int) *golua.LTable {
	/// @struct wg_align
	/// @prop type
	/// @prop at
	/// @method to([]widgets)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_ALIGN))
	state.SetTable(t, golua.LString("at"), golua.LNumber(at))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func alignBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	at := state.GetTable(t, golua.LString("at")).(golua.LNumber)
	a := g.Align(g.AlignmentType(at))

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		a.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return a
}

func msgBoxTable(state *golua.LState, title, content string) *golua.LTable {
	/// @struct wg_msg_box
	/// @prop type
	/// @prop title
	/// @prop content
	/// @method buttons(int)
	/// @method result_callback(callback(bool))

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MSG_BOX))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("content"), golua.LString(content))
	state.SetTable(t, golua.LString("__buttons"), golua.LNil)
	state.SetTable(t, golua.LString("__callback"), golua.LNil)

	tableBuilderFunc(state, t, "buttons", func(state *golua.LState, t *golua.LTable) {
		b := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__buttons"), b)
	})

	tableBuilderFunc(state, t, "result_callback", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__callback"), fn)
	})

	tableBuilderFunc(state, t, "build", func(state *golua.LState, t *golua.LTable) {
		msgBoxBuild(state, t)
	})

	return t
}

func msgBoxBuild(state *golua.LState, t *golua.LTable) *g.MsgboxWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	content := state.GetTable(t, golua.LString("content")).(golua.LString)
	m := g.Msgbox(string(title), string(content))

	buttons := state.GetTable(t, golua.LString("__buttons"))
	if buttons.Type() == golua.LTNumber {
		m.Buttons(g.MsgboxButtons(buttons.(golua.LNumber)))
	}

	callback := state.GetTable(t, golua.LString("__callback"))
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
	/// @struct wg_msg_box_prepare
	/// @prop type
	/// @desc
	/// this is used internally with gui.prepare_msg_box()

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_MSG_BOX_PREPARE))

	return t
}

func msgBoxPrepareBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	return g.PrepareMsgbox()
}

func buttonInvisibleTable(state *golua.LState) *golua.LTable {
	/// @struct wg_button_invisible
	/// @prop type
	/// @method size(width, height)
	/// @method on_click(callback())

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON_INVISIBLE))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__click"), golua.LNil)

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

func buttonInvisibleBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	b := g.InvisibleButton()

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

func buttonImageTable(state *golua.LState, id int, sync bool) *golua.LTable {
	/// @struct wg_button_image
	/// @prop type
	/// @prop id
	/// @prop sync
	/// @method size(width, height)
	/// @method on_click(callback())
	/// @method bg_color(color)
	/// @method tint_color(color)
	/// @method frame_padding(padding)
	/// @method uv(uv0 point, uv1 point)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_BUTTON_IMAGE))
	state.SetTable(t, golua.LString("id"), golua.LNumber(id))
	state.SetTable(t, golua.LString("sync"), golua.LBool(sync))
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__click"), golua.LNil)
	state.SetTable(t, golua.LString("__bgcolor"), golua.LNil)
	state.SetTable(t, golua.LString("__padding"), golua.LNil)
	state.SetTable(t, golua.LString("__tint"), golua.LNil)
	state.SetTable(t, golua.LString("__uv0"), golua.LNil)
	state.SetTable(t, golua.LString("__uv1"), golua.LNil)

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

	tableBuilderFunc(state, t, "bg_color", func(state *golua.LState, t *golua.LTable) {
		tc := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__bgcolor"), tc)
	})

	tableBuilderFunc(state, t, "tint_color", func(state *golua.LState, t *golua.LTable) {
		tc := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__tint"), tc)
	})

	tableBuilderFunc(state, t, "frame_padding", func(state *golua.LState, t *golua.LTable) {
		n := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__padding"), n)
	})

	tableBuilderFunc(state, t, "uv", func(state *golua.LState, t *golua.LTable) {
		uv0 := state.CheckTable(-2)
		uv1 := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__uv0"), uv0)
		state.SetTable(t, golua.LString("__uv1"), uv1)
	})

	return t
}

func buttonImageBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	ig := state.GetTable(t, golua.LString("id")).(golua.LNumber)
	var img image.Image

	sync := state.GetTable(t, golua.LString("sync")).(golua.LBool)

	if !sync {
		<-r.IC.Schedule(int(ig), &collection.Task[collection.ItemImage]{
			Lib:  LIB_GUI,
			Name: "wg_button_image",
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

	b := g.ImageButtonWithRgba(img)

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

	bgcolor := state.GetTable(t, golua.LString("__bgcolor"))
	if bgcolor.Type() == golua.LTTable {
		rgba := imageutil.TableToRGBA(state, bgcolor.(*golua.LTable))
		b.BgColor(rgba)
	}

	tint := state.GetTable(t, golua.LString("__tint"))
	if tint.Type() == golua.LTTable {
		rgba := imageutil.TableToRGBA(state, tint.(*golua.LTable))
		b.TintColor(rgba)
	}

	padding := state.GetTable(t, golua.LString("__padding"))
	if padding.Type() == golua.LTNumber {
		b.FramePadding(int(padding.(golua.LNumber)))
	}

	uv0 := state.GetTable(t, golua.LString("__uv0"))
	uv1 := state.GetTable(t, golua.LString("__uv1"))
	if uv0.Type() == golua.LTTable && uv1.Type() == golua.LTTable {
		p1 := imageutil.TableToPoint(state, uv0.(*golua.LTable))
		p2 := imageutil.TableToPoint(state, uv1.(*golua.LTable))

		b.UV(p1, p2)
	}

	return b
}

func styleTable(state *golua.LState) *golua.LTable {
	/// @struct wg_style
	/// @prop type
	/// @method set_disabled(bool)
	/// @method to([]widgets)
	/// @method set_font_size(float)
	/// @method set_color(colorid, color)
	/// @method set_style(styleid, width, height)
	/// @method set_style_float(styleid, float)
	/// @method font(fontref)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_STYLE))
	state.SetTable(t, golua.LString("__disabled"), golua.LNil)
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)
	state.SetTable(t, golua.LString("__fontsize"), golua.LNil)
	state.SetTable(t, golua.LString("__colors"), state.NewTable())
	state.SetTable(t, golua.LString("__styles"), state.NewTable())
	state.SetTable(t, golua.LString("__stylesfloat"), state.NewTable())
	state.SetTable(t, golua.LString("__font"), golua.LNil)

	tableBuilderFunc(state, t, "set_disabled", func(state *golua.LState, t *golua.LTable) {
		d := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__disabled"), golua.LBool(d))
	})

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	tableBuilderFunc(state, t, "set_font_size", func(state *golua.LState, t *golua.LTable) {
		fnt := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__fontsize"), fnt)
	})

	tableBuilderFunc(state, t, "set_color", func(state *golua.LState, t *golua.LTable) {
		cid := state.CheckNumber(-2)
		col := state.CheckTable(-1)
		ct := state.NewTable()
		state.SetTable(ct, golua.LString("colorid"), cid)
		state.SetTable(ct, golua.LString("color"), col)

		ft := state.GetTable(t, golua.LString("__colors")).(*golua.LTable)
		ft.Append(ct)
	})

	tableBuilderFunc(state, t, "set_style", func(state *golua.LState, t *golua.LTable) {
		sid := state.CheckNumber(-3)
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		st := state.NewTable()
		state.SetTable(st, golua.LString("styleid"), sid)
		state.SetTable(st, golua.LString("width"), width)
		state.SetTable(st, golua.LString("height"), height)

		ft := state.GetTable(t, golua.LString("__styles")).(*golua.LTable)
		ft.Append(st)
	})

	tableBuilderFunc(state, t, "set_style_float", func(state *golua.LState, t *golua.LTable) {
		sid := state.CheckNumber(-2)
		float := state.CheckNumber(-1)
		st := state.NewTable()
		state.SetTable(st, golua.LString("styleid"), sid)
		state.SetTable(st, golua.LString("float"), float)

		ft := state.GetTable(t, golua.LString("__stylesfloat")).(*golua.LTable)
		ft.Append(st)
	})

	tableBuilderFunc(state, t, "font", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__font"), v)
	})

	return t
}

func styleBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	s := g.Style()

	disabled := state.GetTable(t, golua.LString("__disabled"))
	if disabled.Type() == golua.LTBool {
		s.SetDisabled(bool(disabled.(golua.LBool)))
	}

	fontsize := state.GetTable(t, golua.LString("__fontsize"))
	if fontsize.Type() == golua.LTNumber {
		s.SetFontSize(float32(fontsize.(golua.LNumber)))
	}

	colors := state.GetTable(t, golua.LString("__colors")).(*golua.LTable)
	for i := range colors.Len() {
		ct := state.GetTable(colors, golua.LNumber(i+1)).(*golua.LTable)
		cid := state.GetTable(ct, golua.LString("colorid")).(golua.LNumber)
		col := state.GetTable(ct, golua.LString("color")).(*golua.LTable)

		rgba := imageutil.TableToRGBA(state, col)
		s.SetColor(g.StyleColorID(cid), rgba)
	}

	styles := state.GetTable(t, golua.LString("__styles")).(*golua.LTable)
	for i := range styles.Len() {
		st := state.GetTable(styles, golua.LNumber(i+1)).(*golua.LTable)
		sid := state.GetTable(st, golua.LString("styleid")).(golua.LNumber)
		width := state.GetTable(st, golua.LString("width")).(golua.LNumber)
		height := state.GetTable(st, golua.LString("height")).(golua.LNumber)

		s.SetStyle(g.StyleVarID(sid), float32(width), float32(height))
	}

	stylesfloat := state.GetTable(t, golua.LString("__stylesfloat")).(*golua.LTable)
	for i := range stylesfloat.Len() {
		st := state.GetTable(stylesfloat, golua.LNumber(i+1)).(*golua.LTable)
		sid := state.GetTable(st, golua.LString("styleid")).(golua.LNumber)
		float := state.GetTable(st, golua.LString("float")).(golua.LNumber)

		s.SetStyleFloat(g.StyleVarID(sid), float32(float))
	}

	fontref := state.GetTable(t, golua.LString("__font"))
	if fontref.Type() == golua.LTNumber {
		ref := int(fontref.(golua.LNumber))
		sref, err := r.CR_REF.Item(ref)
		if err != nil {
			state.Error(golua.LString(lg.Append(fmt.Sprintf("unable to find ref: %s", err), log.LEVEL_ERROR)), 0)
		}
		font := sref.Value.(*g.FontInfo)

		s.SetFont(font)
	}

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		s.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return s
}

func customTable(state *golua.LState, builder *golua.LFunction) *golua.LTable {
	/// @struct wg_custom
	/// @prop type
	/// @prop builder

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CUSTOM))
	state.SetTable(t, golua.LString("builder"), builder)

	return t
}

func customBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	builder := state.GetTable(t, golua.LString("builder")).(*golua.LFunction)

	c := g.Custom(func() {
		state.Push(builder)
		state.Call(0, 0)
	})

	return c
}

func eventHandlerTable(state *golua.LState) *golua.LTable {
	/// @struct wg_event_handler
	/// @prop type
	/// @method on_activate(callback())
	/// @method on_active(callback())
	/// @method on_deactivate(callback())
	/// @method on_hover(callback())
	/// @method on_click(key, callback())
	/// @method on_dclick(key, callback())
	/// @method on_key_down(key, callback())
	/// @method on_key_pressed(key, callback())
	/// @method on_key_released(key, callback())
	/// @method on_mouse_down(key, callback())
	/// @method on_mouse_released(key, callback())

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_EVENT_HANDLER))
	state.SetTable(t, golua.LString("__activate"), golua.LNil)
	state.SetTable(t, golua.LString("__active"), golua.LNil)
	state.SetTable(t, golua.LString("__deactivate"), golua.LNil)
	state.SetTable(t, golua.LString("__hover"), golua.LNil)
	state.SetTable(t, golua.LString("__click"), state.NewTable())
	state.SetTable(t, golua.LString("__dclick"), state.NewTable())
	state.SetTable(t, golua.LString("__keydown"), state.NewTable())
	state.SetTable(t, golua.LString("__keypressed"), state.NewTable())
	state.SetTable(t, golua.LString("__keyreleased"), state.NewTable())
	state.SetTable(t, golua.LString("__mousedown"), state.NewTable())
	state.SetTable(t, golua.LString("__mousereleased"), state.NewTable())

	tableBuilderFunc(state, t, "on_activate", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__activate"), fn)
	})

	tableBuilderFunc(state, t, "on_active", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__active"), fn)
	})

	tableBuilderFunc(state, t, "on_deactivate", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__deactivate"), fn)
	})

	tableBuilderFunc(state, t, "on_hover", func(state *golua.LState, t *golua.LTable) {
		fn := state.CheckFunction(-1)
		state.SetTable(t, golua.LString("__hover"), fn)
	})

	tableBuilderFunc(state, t, "on_click", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__click")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_dclick", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__dclick")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_down", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__keydown")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_pressed", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__keypressed")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_key_released", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__keyreleased")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_mouse_down", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__mousedown")).(*golua.LTable)
		ft.Append(ev)
	})

	tableBuilderFunc(state, t, "on_mouse_released", func(state *golua.LState, t *golua.LTable) {
		key := state.CheckNumber(-2)
		cb := state.CheckFunction(-1)
		ev := state.NewTable()
		state.SetTable(ev, golua.LString("key"), key)
		state.SetTable(ev, golua.LString("callback"), cb)

		ft := state.GetTable(t, golua.LString("__mousereleased")).(*golua.LTable)
		ft.Append(ev)
	})

	return t
}

func eventHandlerBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	e := g.Event()

	activate := state.GetTable(t, golua.LString("__activate"))
	if activate.Type() == golua.LTFunction {
		e.OnActivate(func() {
			state.Push(activate)
			state.Call(0, 0)
		})
	}

	active := state.GetTable(t, golua.LString("__active"))
	if active.Type() == golua.LTFunction {
		e.OnActive(func() {
			state.Push(active)
			state.Call(0, 0)
		})
	}

	deactivate := state.GetTable(t, golua.LString("__deactivate"))
	if deactivate.Type() == golua.LTFunction {
		e.OnDeactivate(func() {
			state.Push(deactivate)
			state.Call(0, 0)
		})
	}

	hover := state.GetTable(t, golua.LString("__hover"))
	if hover.Type() == golua.LTFunction {
		e.OnHover(func() {
			state.Push(hover)
			state.Call(0, 0)
		})
	}

	click := state.GetTable(t, golua.LString("__click")).(*golua.LTable)
	for i := range click.Len() {
		events := state.GetTable(click, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnClick(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	dclick := state.GetTable(t, golua.LString("__dclick")).(*golua.LTable)
	for i := range dclick.Len() {
		events := state.GetTable(dclick, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnDClick(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keydown := state.GetTable(t, golua.LString("__keydown")).(*golua.LTable)
	for i := range keydown.Len() {
		events := state.GetTable(keydown, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnKeyDown(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keypressed := state.GetTable(t, golua.LString("__keypressed")).(*golua.LTable)
	for i := range keypressed.Len() {
		events := state.GetTable(keypressed, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnKeyPressed(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	keyreleased := state.GetTable(t, golua.LString("__keyreleased")).(*golua.LTable)
	for i := range keyreleased.Len() {
		events := state.GetTable(keyreleased, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnKeyReleased(g.Key(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	mousedown := state.GetTable(t, golua.LString("__mousedown")).(*golua.LTable)
	for i := range mousedown.Len() {
		events := state.GetTable(mousedown, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnMouseDown(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	mousereleased := state.GetTable(t, golua.LString("__mousereleased")).(*golua.LTable)
	for i := range mousereleased.Len() {
		events := state.GetTable(mousereleased, golua.LNumber(i+1)).(*golua.LTable)
		key := state.GetTable(events, golua.LString("key")).(golua.LNumber)
		callback := state.GetTable(events, golua.LString("callback")).(*golua.LFunction)

		e.OnMouseReleased(g.MouseButton(key), func() {
			state.Push(callback)
			state.Call(0, 0)
		})
	}

	return e
}

func plotTable(state *golua.LState, title string) *golua.LTable {
	/// @struct wg_plot
	/// @prop type
	/// @prop title
	/// @method axis_limits(xmin, xmax, ymin, ymax, cond)
	/// @method flags(flags)
	/// @method set_xaxis_label(axis, label)
	/// @method set_yaxis_label(axis, label)
	/// @method size(width, height)
	/// @method x_axeflags(flags)
	/// @method xticks(ticks, default)
	/// @method y_axeflags(flags1, flags2, flags3)
	/// @method yticks(ticks)
	/// @method plots([]plots)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_PLOT))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("__xmin"), golua.LNil)
	state.SetTable(t, golua.LString("__xmax"), golua.LNil)
	state.SetTable(t, golua.LString("__ymin"), golua.LNil)
	state.SetTable(t, golua.LString("__ymax"), golua.LNil)
	state.SetTable(t, golua.LString("__cond"), golua.LNil)
	state.SetTable(t, golua.LString("__flags"), golua.LNil)
	state.SetTable(t, golua.LString("__xlabels"), state.NewTable())
	state.SetTable(t, golua.LString("__ylabels"), state.NewTable())
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__xaxeflags"), golua.LNil)
	state.SetTable(t, golua.LString("__xticks"), golua.LNil)
	state.SetTable(t, golua.LString("__xaticksdefault"), golua.LNil)
	state.SetTable(t, golua.LString("__yaxeflags1"), golua.LNil)
	state.SetTable(t, golua.LString("__yaxeflags2"), golua.LNil)
	state.SetTable(t, golua.LString("__yaxeflags3"), golua.LNil)
	state.SetTable(t, golua.LString("__yticks"), state.NewTable())
	state.SetTable(t, golua.LString("__plots"), golua.LNil)

	tableBuilderFunc(state, t, "axis_limits", func(state *golua.LState, t *golua.LTable) {
		xmin := state.CheckNumber(-5)
		xmax := state.CheckNumber(-4)
		ymin := state.CheckNumber(-3)
		ymax := state.CheckNumber(-2)
		cond := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__xmin"), xmin)
		state.SetTable(t, golua.LString("__xmax"), xmax)
		state.SetTable(t, golua.LString("__ymin"), ymin)
		state.SetTable(t, golua.LString("__ymax"), ymax)
		state.SetTable(t, golua.LString("__cond"), cond)
	})

	tableBuilderFunc(state, t, "flags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__flags"), flags)
	})

	tableBuilderFunc(state, t, "set_xaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		state.SetTable(lt, golua.LString("axis"), axis)
		state.SetTable(lt, golua.LString("label"), golua.LString(label))

		ft := state.GetTable(t, golua.LString("__xlabels")).(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "set_yaxis_label", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-2)
		label := state.CheckString(-1)
		lt := state.NewTable()
		state.SetTable(lt, golua.LString("axis"), axis)
		state.SetTable(lt, golua.LString("label"), golua.LString(label))

		ft := state.GetTable(t, golua.LString("__ylabels")).(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-2)
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "x_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__xaxeflags"), flags)
	})

	tableBuilderFunc(state, t, "xticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-2)
		dflt := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__xticks"), ticks)
		state.SetTable(t, golua.LString("__xaticksdefault"), golua.LBool(dflt))
	})

	tableBuilderFunc(state, t, "y_axeflags", func(state *golua.LState, t *golua.LTable) {
		flags1 := state.CheckNumber(-3)
		flags2 := state.CheckNumber(-2)
		flags3 := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__yaxeflags1"), flags1)
		state.SetTable(t, golua.LString("__yaxeflags2"), flags2)
		state.SetTable(t, golua.LString("__yaxeflags3"), flags3)
	})

	tableBuilderFunc(state, t, "yticks", func(state *golua.LState, t *golua.LTable) {
		ticks := state.CheckTable(-3)
		dflt := state.CheckBool(-2)
		axis := state.CheckNumber(-1)
		lt := state.NewTable()
		state.SetTable(lt, golua.LString("ticks"), ticks)
		state.SetTable(lt, golua.LString("dflt"), golua.LBool(dflt))
		state.SetTable(lt, golua.LString("axis"), axis)

		ft := state.GetTable(t, golua.LString("__ylabels")).(*golua.LTable)
		ft.Append(lt)
	})

	tableBuilderFunc(state, t, "plots", func(state *golua.LState, t *golua.LTable) {
		plots := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__plots"), plots)
	})

	return t
}

func plotBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	p := g.Plot(string(title))

	width := state.GetTable(t, golua.LString("__width"))
	height := state.GetTable(t, golua.LString("__height"))
	if width.Type() == golua.LTNumber && height.Type() == golua.LTNumber {
		p.Size(int(width.(golua.LNumber)), int(height.(golua.LNumber)))
	}

	flags := state.GetTable(t, golua.LString("__flags"))
	if flags.Type() == golua.LTNumber {
		p.Flags(g.PlotFlags(flags.(golua.LNumber)))
	}

	xaxeflags := state.GetTable(t, golua.LString("__xaxeflags"))
	if xaxeflags.Type() == golua.LTNumber {
		p.XAxeFlags(g.PlotAxisFlags(xaxeflags.(golua.LNumber)))
	}

	yaxeflags1 := state.GetTable(t, golua.LString("__yaxeflags1"))
	yaxeflags2 := state.GetTable(t, golua.LString("__yaxeflags2"))
	yaxeflags3 := state.GetTable(t, golua.LString("__yaxeflags3"))
	if yaxeflags1.Type() == golua.LTNumber && yaxeflags2.Type() == golua.LTNumber && yaxeflags3.Type() == golua.LTNumber {
		p.YAxeFlags(g.PlotAxisFlags(yaxeflags1.(golua.LNumber)), g.PlotAxisFlags(yaxeflags2.(golua.LNumber)), g.PlotAxisFlags(yaxeflags3.(golua.LNumber)))
	}

	xmin := state.GetTable(t, golua.LString("__xmin"))
	xmax := state.GetTable(t, golua.LString("__xmax"))
	ymin := state.GetTable(t, golua.LString("__ymin"))
	ymax := state.GetTable(t, golua.LString("__ymax"))
	cond := state.GetTable(t, golua.LString("__cond"))
	if xmin.Type() == golua.LTNumber && xmax.Type() == golua.LTNumber && ymin.Type() == golua.LTNumber && ymax.Type() == golua.LTNumber && cond.Type() == golua.LTNumber {
		p.AxisLimits(
			float64(xmin.(golua.LNumber)), float64(xmax.(golua.LNumber)),
			float64(ymin.(golua.LNumber)), float64(ymax.(golua.LNumber)),
			g.ExecCondition(cond.(golua.LNumber)),
		)
	}

	xlabels := state.GetTable(t, golua.LString("__xlabels")).(*golua.LTable)
	for i := range xlabels.Len() {
		lt := state.GetTable(xlabels, golua.LNumber(i+1)).(*golua.LTable)
		axis := state.GetTable(lt, golua.LString("axis")).(golua.LNumber)
		label := state.GetTable(lt, golua.LString("label")).(golua.LString)

		p.SetXAxisLabel(g.PlotXAxis(axis), string(label))
	}

	ylabels := state.GetTable(t, golua.LString("__ylabels")).(*golua.LTable)
	for i := range ylabels.Len() {
		lt := state.GetTable(ylabels, golua.LNumber(i+1)).(*golua.LTable)
		axis := state.GetTable(lt, golua.LString("axis")).(golua.LNumber)
		label := state.GetTable(lt, golua.LString("label")).(golua.LString)

		p.SetYAxisLabel(g.PlotYAxis(axis), string(label))
	}

	xticks := state.GetTable(t, golua.LString("__xticks"))
	xticksdefault := state.GetTable(t, golua.LString("__xticksdefault"))
	if xticks.Type() == golua.LTTable && xticksdefault.Type() == golua.LTBool {
		ticks := []g.PlotTicker{}

		for i := range (xticks.(*golua.LTable)).Len() {
			tick := state.GetTable(xticks, golua.LNumber(i+1)).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick, state))
		}

		p.XTicks(ticks, bool(xticksdefault.(golua.LBool)))
	}

	yticks := state.GetTable(t, golua.LString("__yticks")).(*golua.LTable)
	for z := range yticks.Len() {
		yticksaxis := state.GetTable(yticks, golua.LNumber(z+1)).(*golua.LTable)

		ytickaxis := state.GetTable(yticksaxis, golua.LString("ticks")).(*golua.LTable)
		dflt := state.GetTable(yticksaxis, golua.LString("dflt")).(golua.LBool)
		axis := state.GetTable(yticksaxis, golua.LString("axis")).(golua.LNumber)

		ticks := []g.PlotTicker{}

		for i := range ytickaxis.Len() {
			tick := state.GetTable(ytickaxis, golua.LNumber(i+1)).(*golua.LTable)
			ticks = append(ticks, plotTickerBuild(tick, state))
		}

		p.YTicks(ticks, bool(dflt), g.ImPlotYAxis(axis))
	}

	plots := state.GetTable(t, golua.LString("__plots"))
	if plots.Type() == golua.LTTable {
		plist := []g.PlotWidget{}

		for i := range (plots.(*golua.LTable)).Len() {
			pt := state.GetTable(plots.(*golua.LTable), golua.LNumber(i+1)).(*golua.LTable)
			plottype := state.GetTable(pt, golua.LString("type")).(golua.LString)

			build := plotList[string(plottype)]
			plist = append(plist, build(r, lg, state, pt))
		}

		p.Plots(plist...)
	}

	return p
}

func plotTickerBuild(t *golua.LTable, state *golua.LState) g.PlotTicker {
	/// @struct plot_ticker
	/// @prop position
	/// @prop label

	position := state.GetTable(t, golua.LString("position")).(golua.LNumber)
	label := state.GetTable(t, golua.LString("label")).(golua.LString)

	return g.PlotTicker{
		Position: float64(position),
		Label:    string(label),
	}
}

func plotBarHTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct pt_bar_h
	/// @prop type
	/// @prop title
	/// @prop data
	/// @method height(height)
	/// @method offset(offset)
	/// @method shift(shift)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_BAR_H))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("data"), data)
	state.SetTable(t, golua.LString("__height"), golua.LNil)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)
	state.SetTable(t, golua.LString("__shift"), golua.LNil)

	tableBuilderFunc(state, t, "height", func(state *golua.LState, t *golua.LTable) {
		height := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__height"), height)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__shift"), shift)
	})

	return t
}

func plotBarHBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	data := state.GetTable(t, golua.LString("data")).(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := state.GetTable(data, golua.LNumber(i+1)).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.BarH(string(title), dataPoints)

	height := state.GetTable(t, golua.LString("__height"))
	if height.Type() == golua.LTNumber {
		p.Height(float64(height.(golua.LNumber)))
	}

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := state.GetTable(t, golua.LString("__shift"))
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotBarTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct pt_bar
	/// @prop type
	/// @prop title
	/// @prop data
	/// @method width(width)
	/// @method offset(offset)
	/// @method shift(shift)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_BAR))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("data"), data)
	state.SetTable(t, golua.LString("__width"), golua.LNil)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)
	state.SetTable(t, golua.LString("__shift"), golua.LNil)

	tableBuilderFunc(state, t, "width", func(state *golua.LState, t *golua.LTable) {
		width := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__width"), width)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	tableBuilderFunc(state, t, "shift", func(state *golua.LState, t *golua.LTable) {
		shift := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__shift"), shift)
	})

	return t
}

func plotBarBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	data := state.GetTable(t, golua.LString("data")).(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := state.GetTable(data, golua.LNumber(i+1)).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Bar(string(title), dataPoints)

	width := state.GetTable(t, golua.LString("__width"))
	if width.Type() == golua.LTNumber {
		p.Width(float64(width.(golua.LNumber)))
	}

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	shift := state.GetTable(t, golua.LString("__shift"))
	if shift.Type() == golua.LTNumber {
		p.Shift(float64(shift.(golua.LNumber)))
	}

	return p
}

func plotLineTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct pt_line
	/// @prop type
	/// @prop title
	/// @prop data
	/// @method set_plot_y_axis(axis)
	/// @method offset(offset)
	/// @method x0(x0)
	/// @method xscale(xscale)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_LINE))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("data"), data)
	state.SetTable(t, golua.LString("__yaxis"), golua.LNil)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)
	state.SetTable(t, golua.LString("__x0"), golua.LNil)
	state.SetTable(t, golua.LString("__xscale"), golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__yaxis"), axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__x0"), x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__xscale"), xscale)
	})

	return t
}

func plotLineBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	data := state.GetTable(t, golua.LString("data")).(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := state.GetTable(data, golua.LNumber(i+1)).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Line(string(title), dataPoints)

	yaxis := state.GetTable(t, golua.LString("__yaxis"))
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := state.GetTable(t, golua.LString("__x0"))
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := state.GetTable(t, golua.LString("__xscale"))
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotLineXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct pt_line_xy
	/// @prop type
	/// @prop title
	/// @prop xdata
	/// @prop ydata
	/// @method set_plot_y_axis(axis)
	/// @method offset(offset)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_LINE_XY))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("xdata"), xdata)
	state.SetTable(t, golua.LString("ydata"), ydata)
	state.SetTable(t, golua.LString("__yaxis"), golua.LNil)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)

	tableBuilderFunc(state, t, "set_plot_y_axis", func(state *golua.LState, t *golua.LTable) {
		axis := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__yaxis"), axis)
	})

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	return t
}

func plotLineXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	xdata := state.GetTable(t, golua.LString("xdata")).(*golua.LTable)
	ydata := state.GetTable(t, golua.LString("ydata")).(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := state.GetTable(xdata, golua.LNumber(i+1)).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := state.GetTable(ydata, golua.LNumber(i+1)).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.LineXY(string(title), xdataPoints, ydataPoints)

	yaxis := state.GetTable(t, golua.LString("__yaxis"))
	if yaxis.Type() == golua.LTNumber {
		p.SetPlotYAxis(g.ImPlotYAxis(yaxis.(golua.LNumber)))
	}

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotPieTable(state *golua.LState, labels golua.LValue, data golua.LValue, x, y, radius float64) *golua.LTable {
	/// @struct pt_pie
	/// @prop type
	/// @prop labels
	/// @prop data
	/// @prop x
	/// @prop y
	/// @prop radius
	/// @method angle0(angle0)
	/// @method label_format(format)
	/// @method normalize(bool)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_PIE_CHART))
	state.SetTable(t, golua.LString("labels"), labels)
	state.SetTable(t, golua.LString("data"), data)
	state.SetTable(t, golua.LString("x"), golua.LNumber(x))
	state.SetTable(t, golua.LString("y"), golua.LNumber(y))
	state.SetTable(t, golua.LString("radius"), golua.LNumber(radius))
	state.SetTable(t, golua.LString("__angle0"), golua.LNil)
	state.SetTable(t, golua.LString("__format"), golua.LNil)
	state.SetTable(t, golua.LString("__normalize"), golua.LNil)

	tableBuilderFunc(state, t, "angle0", func(state *golua.LState, t *golua.LTable) {
		angle0 := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__angle0"), angle0)
	})

	tableBuilderFunc(state, t, "label_format", func(state *golua.LState, t *golua.LTable) {
		format := state.CheckString(-1)
		state.SetTable(t, golua.LString("__format"), golua.LString(format))
	})

	tableBuilderFunc(state, t, "normalize", func(state *golua.LState, t *golua.LTable) {
		normalize := state.CheckBool(-1)
		state.SetTable(t, golua.LString("__normalize"), golua.LBool(normalize))
	})

	return t
}

func plotPieBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	labels := state.GetTable(t, golua.LString("labels")).(*golua.LTable)
	data := state.GetTable(t, golua.LString("data")).(*golua.LTable)
	x := state.GetTable(t, golua.LString("x")).(golua.LNumber)
	y := state.GetTable(t, golua.LString("y")).(golua.LNumber)
	radius := state.GetTable(t, golua.LString("radius")).(golua.LNumber)

	labelPoints := []string{}
	for i := range labels.Len() {
		point := state.GetTable(labels, golua.LNumber(i+1)).(golua.LString)
		labelPoints = append(labelPoints, string(point))
	}

	dataPoints := []float64{}
	for i := range data.Len() {
		point := state.GetTable(data, golua.LNumber(i+1)).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.PieChart(labelPoints, dataPoints, float64(x), float64(y), float64(radius))

	angle0 := state.GetTable(t, golua.LString("__angle0"))
	if angle0.Type() == golua.LTNumber {
		p.Angle0(float64(angle0.(golua.LNumber)))
	}

	format := state.GetTable(t, golua.LString("__format"))
	if format.Type() == golua.LTString {
		p.LabelFormat(string(format.(golua.LString)))
	}

	normalize := state.GetTable(t, golua.LString("__normalize"))
	if normalize.Type() == golua.LTBool {
		p.Normalize(bool(normalize.(golua.LBool)))
	}

	return p
}

func plotScatterTable(state *golua.LState, title string, data golua.LValue) *golua.LTable {
	/// @struct pt_scatter
	/// @prop type
	/// @prop title
	/// @prop data
	/// @method offset(offset)
	/// @method x0(x0)
	/// @method xscale(xscale)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_SCATTER))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("data"), data)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)
	state.SetTable(t, golua.LString("__x0"), golua.LNil)
	state.SetTable(t, golua.LString("__xscale"), golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	tableBuilderFunc(state, t, "x0", func(state *golua.LState, t *golua.LTable) {
		x0 := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__x0"), x0)
	})

	tableBuilderFunc(state, t, "xscale", func(state *golua.LState, t *golua.LTable) {
		xscale := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__xscale"), xscale)
	})

	return t
}

func plotScatterBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	data := state.GetTable(t, golua.LString("data")).(*golua.LTable)

	dataPoints := []float64{}
	for i := range data.Len() {
		point := state.GetTable(data, golua.LNumber(i+1)).(golua.LNumber)
		dataPoints = append(dataPoints, float64(point))
	}

	p := g.Scatter(string(title), dataPoints)

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	x0 := state.GetTable(t, golua.LString("__x0"))
	if x0.Type() == golua.LTNumber {
		p.X0(float64(x0.(golua.LNumber)))
	}

	xscale := state.GetTable(t, golua.LString("__xscale"))
	if xscale.Type() == golua.LTNumber {
		p.XScale(float64(xscale.(golua.LNumber)))
	}

	return p
}

func plotScatterXYTable(state *golua.LState, title string, xdata, ydata golua.LValue) *golua.LTable {
	/// @struct pt_scatter_xy
	/// @prop type
	/// @prop title
	/// @prop xdata
	/// @prop ydata
	/// @method offset(offset)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_SCATTER_XY))
	state.SetTable(t, golua.LString("title"), golua.LString(title))
	state.SetTable(t, golua.LString("xdata"), xdata)
	state.SetTable(t, golua.LString("ydata"), ydata)
	state.SetTable(t, golua.LString("__offset"), golua.LNil)

	tableBuilderFunc(state, t, "offset", func(state *golua.LState, t *golua.LTable) {
		offset := state.CheckNumber(-1)
		state.SetTable(t, golua.LString("__offset"), offset)
	})

	return t
}

func plotScatterXYBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	title := state.GetTable(t, golua.LString("title")).(golua.LString)
	xdata := state.GetTable(t, golua.LString("xdata")).(*golua.LTable)
	ydata := state.GetTable(t, golua.LString("ydata")).(*golua.LTable)

	xdataPoints := []float64{}
	for i := range xdata.Len() {
		point := state.GetTable(xdata, golua.LNumber(i+1)).(golua.LNumber)
		xdataPoints = append(xdataPoints, float64(point))
	}

	ydataPoints := []float64{}
	for i := range ydata.Len() {
		point := state.GetTable(ydata, golua.LNumber(i+1)).(golua.LNumber)
		ydataPoints = append(ydataPoints, float64(point))
	}

	p := g.ScatterXY(string(title), xdataPoints, ydataPoints)

	offset := state.GetTable(t, golua.LString("__offset"))
	if offset.Type() == golua.LTNumber {
		p.Offset(int(offset.(golua.LNumber)))
	}

	return p
}

func plotCustomTable(state *golua.LState, builder *golua.LFunction) *golua.LTable {
	/// @struct pt_custom
	/// @prop type
	/// @prop builder

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(PLOT_CUSTOM))
	state.SetTable(t, golua.LString("builder"), builder)

	return t
}

func plotCustomBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.PlotWidget {
	builder := state.GetTable(t, golua.LString("builder")).(*golua.LFunction)

	c := g.Custom(func() {
		state.Push(builder)
		state.Call(0, 0)
	})

	return c
}

func cssTagTable(state *golua.LState, tag string) *golua.LTable {
	/// @struct wg_css_tag
	/// @prop type
	/// @prop tag
	/// @method to([]widgets)

	t := state.NewTable()
	state.SetTable(t, golua.LString("type"), golua.LString(WIDGET_CSS_TAG))
	state.SetTable(t, golua.LString("tag"), golua.LString(tag))
	state.SetTable(t, golua.LString("__widgets"), golua.LNil)

	tableBuilderFunc(state, t, "to", func(state *golua.LState, t *golua.LTable) {
		lt := state.CheckTable(-1)
		state.SetTable(t, golua.LString("__widgets"), lt)
	})

	return t
}

func cssTagBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	tag := state.GetTable(t, golua.LString("tag")).(golua.LString)
	c := g.CSSTag(string(tag))

	layout := state.GetTable(t, golua.LString("__widgets"))
	if layout.Type() == golua.LTTable {
		c.To(layoutBuild(r, state, parseWidgets(parseTable(layout.(*golua.LTable), state), state, lg), lg)...)
	}

	return c
}
