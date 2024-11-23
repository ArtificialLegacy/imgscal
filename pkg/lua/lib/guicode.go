package lib

import (
	g "github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_GUICODE = "guicode"

/// @lib GUI Code
/// @import guicode
/// @desc
/// Extension of the GUI library for code editor widgets.

func RegisterGUICode(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_GUICODE, r, r.State, lg)

	/// @func wg_code_editor() -> struct<guicode.WidgetCodeEditor>
	/// @returns {struct<guicode.WidgetCodeEditor>}
	/// @desc
	/// Note this should be treated as persistent, do not recreate on each render.
	lib.CreateFunction(tab, "wg_code_editor",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			editor := g.CodeEditor()
			id := r.CR_CED.Add(editor)
			t := codeeditorTable(state, id, editor)

			state.Push(t)
			return 1
		})

	/// @constants LanguageDefinition {int}
	/// @const LANG_NONE
	/// @const LANG_CPP
	/// @const LANG_C
	/// @const LANG_CS
	/// @const LANG_PYTHON
	/// @const LANG_LUA
	/// @const LANG_JSON
	/// @const LANG_SQL
	/// @const LANG_ANGELSCRIPT
	/// @const LANG_GLSL
	/// @const LANG_HLSL
	tab.RawSetString("LANG_NONE", golua.LNumber(g.LanguageDefinitionNone))
	tab.RawSetString("LANG_CPP", golua.LNumber(g.LanguageDefinitionCPP))
	tab.RawSetString("LANG_C", golua.LNumber(g.LanguageDefinitionC))
	tab.RawSetString("LANG_CS", golua.LNumber(g.LanguageDefinitionCs))
	tab.RawSetString("LANG_PYTHON", golua.LNumber(g.LanguageDefinitionPython))
	tab.RawSetString("LANG_LUA", golua.LNumber(g.LanguageDefinitionLua))
	tab.RawSetString("LANG_JSON", golua.LNumber(g.LanguageDefinitionJSON))
	tab.RawSetString("LANG_SQL", golua.LNumber(g.LanguageDefinitionSQL))
	tab.RawSetString("LANG_ANGELSCRIPT", golua.LNumber(g.LanguageDefinitionAngelScript))
	tab.RawSetString("LANG_GLSL", golua.LNumber(g.LanguageDefinitionGlsl))
	tab.RawSetString("LANG_HLSL", golua.LNumber(g.LanguageDefinitionHlsl))
}

func codeeditorTable(state *golua.LState, id int, editor *g.CodeEditorWidget) *golua.LTable {
	/// @struct WidgetCodeEditor
	/// @prop type {string<gui.WidgetType>}
	/// @prop id {int<collection.CRATE_CODEEDITOR>}
	/// @method border(self, v bool) -> self
	/// @method copy()
	/// @method cut()
	/// @method delete()
	/// @method paste()
	/// @method text() -> string
	/// @method current_line_text() -> string
	/// @method selected_text() -> string
	/// @method word_under_cursor() -> string
	/// @method cursor_pos() -> int, int
	/// @method cursor_screen_pos() -> int, int
	/// @method selection_start() -> int, int
	/// @method handle_keyboard_inputs(self, v bool) -> self
	/// @method has_selection() -> bool
	/// @method is_text_changed() -> bool
	/// @method insert_text(self, v string) -> self
	/// @method text(self, v string) -> self
	/// @method show_whitespace(self, v bool) -> self
	/// @method tab_size(self, v int) -> self
	/// @method size(self, x int, y int) -> self
	/// @method language_definition(self, v int<guicode.LanguageDefinition>) -> self

	t := state.NewTable()
	t.RawSetString("type", golua.LString(WIDGET_CODEEDITOR))
	t.RawSetString("id", golua.LNumber(id))

	tableBuilderFunc(state, t, "border", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		editor.Border(v)
	})

	t.RawSetString("copy", state.NewFunction(func(state *golua.LState) int {
		editor.Copy()
		return 0
	}))

	t.RawSetString("cut", state.NewFunction(func(state *golua.LState) int {
		editor.Cut()
		return 0
	}))

	t.RawSetString("delete", state.NewFunction(func(state *golua.LState) int {
		editor.Delete()
		return 0
	}))

	t.RawSetString("paste", state.NewFunction(func(state *golua.LState) int {
		editor.Paste()
		return 0
	}))

	t.RawSetString("text", state.NewFunction(func(state *golua.LState) int {
		s := editor.GetText()

		state.Push(golua.LString(s))
		return 1
	}))

	t.RawSetString("current_line_text", state.NewFunction(func(state *golua.LState) int {
		s := editor.GetCurrentLineText()

		state.Push(golua.LString(s))
		return 1
	}))

	t.RawSetString("selected_text", state.NewFunction(func(state *golua.LState) int {
		s := editor.GetSelectedText()

		state.Push(golua.LString(s))
		return 1
	}))

	t.RawSetString("word_under_cursor", state.NewFunction(func(state *golua.LState) int {
		s := editor.GetWordUnderCursor()

		state.Push(golua.LString(s))
		return 1
	}))

	t.RawSetString("cursor_pos", state.NewFunction(func(state *golua.LState) int {
		x, y := editor.GetCursorPos()

		state.Push(golua.LNumber(x))
		state.Push(golua.LNumber(y))
		return 2
	}))

	t.RawSetString("cursor_screen_pos", state.NewFunction(func(state *golua.LState) int {
		x, y := editor.GetScreenCursorPos()

		state.Push(golua.LNumber(x))
		state.Push(golua.LNumber(y))
		return 2
	}))

	t.RawSetString("selection_start", state.NewFunction(func(state *golua.LState) int {
		x, y := editor.GetSelectionStart()

		state.Push(golua.LNumber(x))
		state.Push(golua.LNumber(y))
		return 2
	}))

	tableBuilderFunc(state, t, "handle_keyboard_inputs", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		editor.HandleKeyboardInputs(v)
	})

	t.RawSetString("has_selection", state.NewFunction(func(state *golua.LState) int {
		b := editor.HasSelection()

		state.Push(golua.LBool(b))
		return 1
	}))

	t.RawSetString("is_text_changed", state.NewFunction(func(state *golua.LState) int {
		b := editor.IsTextChanged()

		state.Push(golua.LBool(b))
		return 1
	}))

	tableBuilderFunc(state, t, "insert_text", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckString(-1)
		editor.InsertText(v)
	})

	tableBuilderFunc(state, t, "text", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckString(-1)
		editor.Text(v)
	})

	tableBuilderFunc(state, t, "show_whitespace", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckBool(-1)
		editor.ShowWhitespaces(v)
	})

	tableBuilderFunc(state, t, "tab_size", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		editor.TabSize(int(v))
	})

	tableBuilderFunc(state, t, "size", func(state *golua.LState, t *golua.LTable) {
		x := state.CheckNumber(-2)
		y := state.CheckNumber(-1)
		editor.Size(float32(x), float32(y))
	})

	tableBuilderFunc(state, t, "language_definition", func(state *golua.LState, t *golua.LTable) {
		v := state.CheckNumber(-1)
		editor.LanguageDefinition(g.LanguageDefinition(v))
	})

	return t
}

func codeeditorBuild(r *lua.Runner, lg *log.Logger, state *golua.LState, t *golua.LTable) g.Widget {
	c, err := r.CR_CED.Item(int(t.RawGetString("id").(golua.LNumber)))
	if err != nil {
		lua.Error(state, lg.Append("failed to find code editor", log.LEVEL_ERROR))
	}

	return c
}
