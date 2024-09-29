package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/charmbracelet/lipgloss"
	golua "github.com/yuin/gopher-lua"
)

const LIB_LIPGLOSS = "lipgloss"

/// @lib LipGloss
/// @import lipgloss
/// @desc
/// Wrapper for lipgloss library.

func RegisterLipGloss(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_LIPGLOSS, r, r.State, lg)

	/// @func has_dark_background() -> bool
	/// @returns {bool}
	lib.CreateFunction(tab, "has_dark_background",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			isDark := lipgloss.HasDarkBackground()

			state.Push(golua.LBool(isDark))
			return 1
		})

	/// @func size(str) -> int, int
	/// @arg str {string} - The string to measure.
	/// @returns {int} - The width of the string.
	/// @returns {int} - The height of the string.
	lib.CreateFunction(tab, "size",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width, height := lipgloss.Size(args["str"].(string))

			state.Push(golua.LNumber(width))
			state.Push(golua.LNumber(height))
			return 2
		})

	/// @func width(str) -> int
	/// @arg str {string} - The string to measure.
	/// @returns {int} - The width of the string.
	lib.CreateFunction(tab, "width",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := lipgloss.Width(args["str"].(string))

			state.Push(golua.LNumber(width))
			return 1
		})

	/// @func height(str) -> int
	/// @arg str {string} - The string to measure.
	/// @returns {int} - The height of the string.
	lib.CreateFunction(tab, "height",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			height := lipgloss.Height(args["str"].(string))

			state.Push(golua.LNumber(height))
			return 1
		})

	/// @func join_horizontal(pos, str...) -> string
	/// @arg pos {float<lipgloss.Position>} - Percentage between 0 and 1.
	/// @arg str {string...} - The strings to join.
	/// @returns {string} - The joined string.
	lib.CreateFunction(tab, "join_horizontal",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "pos"},
			lua.ArgVariadic("str", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pos := args["pos"].(float64)

			str := args["str"].([]any)
			strList := make([]string, len(str))
			for i, s := range str {
				strList[i] = s.(string)
			}

			result := lipgloss.JoinHorizontal(lipgloss.Position(pos), strList...)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func join_vertical(pos, str...) -> string
	/// @arg pos {float<lipgloss.Position>} - Percentage between 0 and 1.
	/// @arg str {string...} - The strings to join.
	/// @returns {string} - The joined string.
	lib.CreateFunction(tab, "join_vertical",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "pos"},
			lua.ArgVariadic("str", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pos := args["pos"].(float64)

			str := args["str"].([]any)
			strList := make([]string, len(str))
			for i, s := range str {
				strList[i] = s.(string)
			}

			result := lipgloss.JoinVertical(lipgloss.Position(pos), strList...)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func color(col) -> struct<lipgloss.Color>
	/// @arg col {string}
	/// @returns {struct<lipgloss.Color>}
	lib.CreateFunction(tab, "color",
		[]lua.Arg{
			{Type: lua.STRING, Name: "col"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorTable(state, args["col"].(string))

			state.Push(col)
			return 1
		})

	/// @func color_none() -> struct<lipgloss.ColorNone>
	/// @returns {struct<lipgloss.ColorNone>}
	lib.CreateFunction(tab, "color_none",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorNoneTable(state)

			state.Push(col)
			return 1
		})

	/// @func color_ansi(col) -> struct<lipgloss.ColorAnsi>
	/// @arg col {int}
	/// @returns {struct<lipgloss.Color>}
	lib.CreateFunction(tab, "color_ansi",
		[]lua.Arg{
			{Type: lua.INT, Name: "col"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorAnsiTable(state, args["col"].(int))

			state.Push(col)
			return 1
		})

	/// @func color_adaptive(light, dark) -> struct<lipgloss.ColorAdaptive>
	/// @arg light {string}
	/// @arg dark {string}
	/// @returns {struct<lipgloss.ColorAdaptive>}
	lib.CreateFunction(tab, "color_adaptive",
		[]lua.Arg{
			{Type: lua.STRING, Name: "light"},
			{Type: lua.STRING, Name: "dark"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorAdaptiveTable(state, args["light"].(string), args["dark"].(string))

			state.Push(col)
			return 1
		})

	/// @func color_complete(truecolor, ansi256, ansi) -> struct<lipgloss.ColorComplete>
	/// @arg truecolor {string}
	/// @arg ansi256 {string}
	/// @arg ansi {string}
	/// @returns {struct<lipgloss.ColorComplete>}
	lib.CreateFunction(tab, "color_complete",
		[]lua.Arg{
			{Type: lua.STRING, Name: "truecolor"},
			{Type: lua.STRING, Name: "ansi256"},
			{Type: lua.STRING, Name: "ansi"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorCompleteTable(state, args["truecolor"].(string), args["ansi256"].(string), args["ansi"].(string))

			state.Push(col)
			return 1
		})

	/// @func color_complete_adaptive(light, dark) -> struct<lipgloss.ColorCompleteAdaptive>
	/// @arg light {struct<lipgloss.ColorComplete>}
	/// @arg dark {struct<lipgloss.ColorComplete>}
	/// @returns {struct<lipgloss.ColorCompleteAdaptive>}
	lib.CreateFunction(tab, "color_complete_adaptive",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "light"},
			{Type: lua.RAW_TABLE, Name: "dark"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			col := lgColorCompleteAdaptiveTable(state, args["light"].(*golua.LTable), args["dark"].(*golua.LTable))

			state.Push(col)
			return 1
		})

	/// @func border(top, bottom, left, right, topleft, topright, bottomleft, bottomright, middleleft, middleright, middle, middletop, middlebottom) -> struct<lipgloss.Border>
	/// @arg top {string}
	/// @arg bottom {string}
	/// @arg left {string}
	/// @arg right {string}
	/// @arg topleft {string}
	/// @arg topright {string}
	/// @arg bottomleft {string}
	/// @arg bottomright {string}
	/// @arg middleleft {string}
	/// @arg middleright {string}
	/// @arg middle {string}
	/// @arg middletop {string}
	/// @arg middlebottom {string}
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border",
		[]lua.Arg{
			{Type: lua.STRING, Name: "top"},
			{Type: lua.STRING, Name: "bottom"},
			{Type: lua.STRING, Name: "left"},
			{Type: lua.STRING, Name: "right"},
			{Type: lua.STRING, Name: "topleft"},
			{Type: lua.STRING, Name: "topright"},
			{Type: lua.STRING, Name: "bottomleft"},
			{Type: lua.STRING, Name: "bottomright"},
			{Type: lua.STRING, Name: "middleleft"},
			{Type: lua.STRING, Name: "middleright"},
			{Type: lua.STRING, Name: "middle"},
			{Type: lua.STRING, Name: "middletop"},
			{Type: lua.STRING, Name: "middlebottom"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTable(state,
				args["top"].(string),
				args["bottom"].(string),
				args["left"].(string),
				args["right"].(string),
				args["topleft"].(string),
				args["topright"].(string),
				args["bottomleft"].(string),
				args["bottomright"].(string),
				args["middleleft"].(string),
				args["middleright"].(string),
				args["middle"].(string),
				args["middletop"].(string),
				args["middlebottom"].(string),
			)

			state.Push(border)
			return 1
		})

	/// @func border_block() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_block",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.BlockBorder())

			state.Push(border)
			return 1
		})

	/// @func border_double() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_double",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.DoubleBorder())

			state.Push(border)
			return 1
		})

	/// @func border_hidden() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_hidden",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.HiddenBorder())

			state.Push(border)
			return 1
		})

	/// @func border_block_inner_half() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_block_inner_half",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.InnerHalfBlockBorder())

			state.Push(border)
			return 1
		})

	/// @func border_normal() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_normal",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.NormalBorder())

			state.Push(border)
			return 1
		})

	/// @func border_block_outer_half() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_block_outer_half",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.OuterHalfBlockBorder())

			state.Push(border)
			return 1
		})

	/// @func border_rounded() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_rounded",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.RoundedBorder())

			state.Push(border)
			return 1
		})

	/// @func border_thick() -> struct<lipgloss.Border>
	/// @returns {struct<lipgloss.Border>}
	lib.CreateFunction(tab, "border_thick",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			border := lgBorderTableFrom(state, lipgloss.ThickBorder())

			state.Push(border)
			return 1
		})

	/// @func whitespace_option() -> struct<lipgloss.WhitespaceOption>
	/// @returns {lipgloss.WhitespaceOption}
	lib.CreateFunction(tab, "whitespace_option",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := whitespaceOptionTable(state, lib)

			state.Push(t)
			return 1
		})

	/// @func place(width, height, hpos, vpos, str, opts?) -> string
	/// @arg width {int}
	/// @arg height {int}
	/// @arg hpos {float<lipgloss.Position>}
	/// @arg vpos {float<lipgloss.Position>}
	/// @arg str {string}
	/// @arg? opts {struct<lipgloss.WhitespaceOption>}
	/// @returns {string}
	lib.CreateFunction(tab, "place",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.INT, Name: "height"},
			{Type: lua.FLOAT, Name: "hpos"},
			{Type: lua.FLOAT, Name: "vpos"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "opts", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := args["width"].(int)
			height := args["height"].(int)
			hpos := args["hpos"].(float64)
			vpos := args["vpos"].(float64)
			str := args["str"].(string)
			opts := args["opts"].(*golua.LTable)
			wopts := whitespaceOptionBuild(opts)

			result := lipgloss.Place(width, height, lipgloss.Position(hpos), lipgloss.Position(vpos), str, wopts...)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func place_horizontal(width, hpos, str, opts?) -> string
	/// @arg width {int}
	/// @arg hpos {float<lipgloss.Position>}
	/// @arg str {string}
	/// @arg? opts {struct<lipgloss.WhitespaceOption>}
	/// @returns {string}
	lib.CreateFunction(tab, "place_horizontal",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
			{Type: lua.FLOAT, Name: "hpos"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "opts", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			width := args["width"].(int)
			hpos := args["hpos"].(float64)
			str := args["str"].(string)
			opts := args["opts"].(*golua.LTable)
			wopts := whitespaceOptionBuild(opts)

			result := lipgloss.PlaceHorizontal(width, lipgloss.Position(hpos), str, wopts...)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func place_vertical(height, vpos, str, opts?) -> string
	/// @arg height {int}
	/// @arg vpos {float<lipgloss.Position>}
	/// @arg str {string}
	/// @arg? opts {struct<lipgloss.WhitespaceOption>}
	/// @returns {string}
	lib.CreateFunction(tab, "place_vertical",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
			{Type: lua.FLOAT, Name: "vpos"},
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "opts", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			height := args["height"].(int)
			vpos := args["vpos"].(float64)
			str := args["str"].(string)
			opts := args["opts"].(*golua.LTable)
			wopts := whitespaceOptionBuild(opts)

			result := lipgloss.PlaceVertical(height, lipgloss.Position(vpos), str, wopts...)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func style() -> struct<lipgloss.Style>
	/// @returns {struct<lipgloss.Style>}
	lib.CreateFunction(tab, "style",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			s := lipgloss.NewStyle()
			id := r.CR_LIP.Add(&collection.StyleItem{
				Style: &s,
			})

			style := lipglossStyleTable(state, lib, r, id)
			state.Push(style)
			return 1
		})

	/// @func style_runes(str, indices, match, unmatch) -> string
	/// @arg str {string}
	/// @arg indices {[]int}
	/// @arg match {struct<lipgloss.Style>}
	/// @arg unmatch {struct<lipgloss.Style>}
	/// @returns {string}
	lib.CreateFunction(tab, "style_runes",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			lua.ArgArray("indices", lua.ArrayType{Type: lua.INT}, false),
			{Type: lua.RAW_TABLE, Name: "match"},
			{Type: lua.RAW_TABLE, Name: "unmatch"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := args["str"].(string)

			indices := args["indices"].([]any)
			indList := make([]int, len(indices))
			for i, v := range indices {
				indList[i] = v.(int)
			}

			match := args["match"].(*golua.LTable)
			unmatch := args["unmatch"].(*golua.LTable)

			matchid := match.RawGetString("id").(golua.LNumber)
			unmatchid := unmatch.RawGetString("id").(golua.LNumber)

			ms, _ := r.CR_LIP.Item(int(matchid))
			us, _ := r.CR_LIP.Item(int(unmatchid))

			result := lipgloss.StyleRunes(str, indList, *ms.Style, *us.Style)

			state.Push(golua.LString(result))
			return 1
		})

	/// @func style_string(str, style) -> string
	/// @arg str {string}
	/// @arg style {struct<lipgloss.Style>}
	/// @returns {string}
	lib.CreateFunction(tab, "style_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
			{Type: lua.RAW_TABLE, Name: "style"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			str := args["str"].(string)

			style := args["style"].(*golua.LTable)

			sid := style.RawGetString("id").(golua.LNumber)

			ss, _ := r.CR_LIP.Item(int(sid))

			result := ss.Style.Render(str)

			state.Push(golua.LString(result))
			return 1
		})

	/// @constants Position
	/// @const POSITION_TOP
	/// @const POSITION_BOTTOM
	/// @const POSITION_CENTER
	/// @const POSITION_LEFT
	/// @const POSITION_RIGHT
	tab.RawSetString("POSITION_TOP", golua.LNumber(0.0))
	tab.RawSetString("POSITION_BOTTOM", golua.LNumber(1.0))
	tab.RawSetString("POSITION_CENTER", golua.LNumber(0.5))
	tab.RawSetString("POSITION_LEFT", golua.LNumber(0.0))
	tab.RawSetString("POSITION_RIGHT", golua.LNumber(1.0))

	/// @constants ColorType
	/// @const COLOR
	/// @const COLOR_NONE
	/// @const COLOR_ANSI
	/// @const COLOR_ADAPTIVE
	/// @const COLOR_COMPLETE
	/// @const COLOR_COMPLETEADAPTIVE
	tab.RawSetString("COLOR", golua.LNumber(LG_COLOR))
	tab.RawSetString("COLOR_NONE", golua.LNumber(LG_COLOR_NONE))
	tab.RawSetString("COLOR_ANSI", golua.LNumber(LG_COLOR_ANSI))
	tab.RawSetString("COLOR_ADAPTIVE", golua.LNumber(LG_COLOR_ADAPTIVE))
	tab.RawSetString("COLOR_COMPLETE", golua.LNumber(LG_COLOR_COMPLETE))
	tab.RawSetString("COLOR_COMPLETEADAPTIVE", golua.LNumber(LG_COLOR_COMPLETEADAPTIVE))

	/// @constants Special
	/// @const NOTAB_CONVERSION
	tab.RawSetString("NOTAB_CONVERSION", golua.LNumber(lipgloss.NoTabConversion))
}

type LipGlossColor int

const (
	LG_COLOR LipGlossColor = iota
	LG_COLOR_NONE
	LG_COLOR_ANSI
	LG_COLOR_ADAPTIVE
	LG_COLOR_COMPLETE
	LG_COLOR_COMPLETEADAPTIVE
)

type ColorBuilder func(t *golua.LTable) lipgloss.TerminalColor

var colorList = []ColorBuilder{
	lgColorBuild,
	lgColorNoneBuild,
	lgColorAnsiBuild,
	lgColorAdaptiveBuild,
	lgColorCompleteBuild,
	lgColorCompleteAdaptiveBuild,
}

func lgColorGenericTable(state *golua.LState, c lipgloss.TerminalColor) *golua.LTable {
	var t *golua.LTable

	switch c := c.(type) {
	case lipgloss.Color:
		t = lgColorTable(state, string(c))
	case lipgloss.NoColor:
		t = lgColorNoneTable(state)
	case lipgloss.ANSIColor:
		t = lgColorAnsiTable(state, int(c))
	case lipgloss.AdaptiveColor:
		t = lgColorAdaptiveTable(state, c.Light, c.Dark)
	case lipgloss.CompleteColor:
		t = lgColorCompleteTable(state, c.TrueColor, c.ANSI256, c.ANSI)
	case lipgloss.CompleteAdaptiveColor:
		t = lgColorCompleteAdaptiveTable(state,
			lgColorCompleteTable(state, c.Light.TrueColor, c.Light.ANSI256, c.Light.ANSI),
			lgColorCompleteTable(state, c.Dark.TrueColor, c.Dark.ANSI256, c.Dark.ANSI),
		)
	}

	return t
}

func lgColorGenericBuild(t *golua.LTable) lipgloss.TerminalColor {
	/// @struct ColorAny
	/// @prop type {int<lipgloss.ColorType>} - The type of color.

	typ := t.RawGetString("type")
	if typ.Type() != golua.LTNumber {
		return lipgloss.NoColor{}
	}

	col := colorList[int(typ.(golua.LNumber))](t)
	return col
}

func lgColorTable(state *golua.LState, value string) *golua.LTable {
	/// @struct Color
	/// @prop type {int<lipgloss.ColorType>} - The type of color.
	/// @prop value {string} - The color value.

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR))
	t.RawSetString("value", golua.LString(value))

	return t
}

func lgColorBuild(t *golua.LTable) lipgloss.TerminalColor {
	value := t.RawGetString("value").(golua.LString)
	return lipgloss.Color(value)
}

func lgColorNoneTable(state *golua.LState) *golua.LTable {
	/// @struct ColorNone
	/// @prop type {int<lipgloss.ColorType>} - The type of color.

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR_NONE))

	return t
}

func lgColorNoneBuild(t *golua.LTable) lipgloss.TerminalColor {
	return lipgloss.NoColor{}
}

func lgColorAnsiTable(state *golua.LState, value int) *golua.LTable {
	/// @struct ColorAnsi
	/// @prop type {int<lipgloss.ColorType>} - The type of color.
	/// @prop value {int} - The color value.

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR_ANSI))
	t.RawSetString("value", golua.LNumber(value))

	return t
}

func lgColorAnsiBuild(t *golua.LTable) lipgloss.TerminalColor {
	value := t.RawGetString("value").(golua.LNumber)
	return lipgloss.ANSIColor(value)
}

func lgColorAdaptiveTable(state *golua.LState, light, dark string) *golua.LTable {
	/// @struct ColorAdaptive
	/// @prop type {int<lipgloss.ColorType>} - The type of color.
	/// @prop light {string} - The color value for light backgrounds.
	/// @prop dark {string} - The color value for dark backgrounds.

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR_ADAPTIVE))
	t.RawSetString("light", golua.LString(light))
	t.RawSetString("dark", golua.LString(dark))

	return t
}

func lgColorAdaptiveBuild(t *golua.LTable) lipgloss.TerminalColor {
	light := t.RawGetString("light").(golua.LString)
	dark := t.RawGetString("dark").(golua.LString)
	return lipgloss.AdaptiveColor{
		Light: string(light),
		Dark:  string(dark),
	}
}

func lgColorCompleteTable(state *golua.LState, truecolor, ansi256, ansi string) *golua.LTable {
	/// @struct ColorComplete
	/// @prop type {int<lipgloss.ColorType>} - The type of color.
	/// @prop truecolor {string}
	/// @prop ansi256 {string}
	/// @prop ansi {string}

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR_COMPLETE))
	t.RawSetString("truecolor", golua.LString(truecolor))
	t.RawSetString("ansi256", golua.LString(ansi256))
	t.RawSetString("ansi", golua.LString(ansi))

	return t
}

func lgColorCompleteBuild(t *golua.LTable) lipgloss.TerminalColor {
	truecolor := t.RawGetString("truecolor").(golua.LString)
	ansi256 := t.RawGetString("ansi256").(golua.LString)
	ansi := t.RawGetString("ansi").(golua.LString)
	return lipgloss.CompleteColor{
		TrueColor: string(truecolor),
		ANSI256:   string(ansi256),
		ANSI:      string(ansi),
	}
}

func lgColorCompleteAdaptiveTable(state *golua.LState, light, dark *golua.LTable) *golua.LTable {
	/// @struct ColorCompleteAdaptive
	/// @prop type {int<lipgloss.ColorType>} - The type of color.
	/// @prop light {struct<lipgloss.ColorComplete>}
	/// @prop dark {struct<lipgloss.ColorComplete>}

	t := state.NewTable()

	t.RawSetString("type", golua.LNumber(LG_COLOR_COMPLETEADAPTIVE))
	t.RawSetString("light", light)
	t.RawSetString("dark", dark)

	return t
}

func lgColorCompleteAdaptiveBuild(t *golua.LTable) lipgloss.TerminalColor {
	light := t.RawGetString("light").(*golua.LTable)
	dark := t.RawGetString("dark").(*golua.LTable)
	return lipgloss.CompleteAdaptiveColor{
		Light: lgColorCompleteBuild(light).(lipgloss.CompleteColor),
		Dark:  lgColorCompleteBuild(dark).(lipgloss.CompleteColor),
	}
}

func lgBorderTableFrom(state *golua.LState, border lipgloss.Border) *golua.LTable {
	b := lgBorderTable(state,
		border.Top,
		border.Bottom,
		border.Left,
		border.Right,
		border.TopLeft,
		border.TopRight,
		border.BottomLeft,
		border.BottomRight,
		border.MiddleLeft,
		border.MiddleRight,
		border.Middle,
		border.MiddleTop,
		border.MiddleBottom,
	)

	return b
}

func lgBorderTable(state *golua.LState, top, bottom, left, right, topleft, topright, bottomleft, bottomright, middleleft, middleright, middle, middletop, middlebottom string) *golua.LTable {
	/// @struct Border
	/// @prop top {string}
	/// @prop bottom {string}
	/// @prop left {string}
	/// @prop right {string}
	/// @prop topleft {string}
	/// @prop topright {string}
	/// @prop bottomleft {string}
	/// @prop bottomright {string}
	/// @prop middleleft {string}
	/// @prop middleright {string}
	/// @prop middle {string}
	/// @prop middletop {string}
	/// @prop middlebottom {string}
	/// @method size_top() -> int
	/// @method size_bottom() -> int
	/// @method size_left() -> int
	/// @method size_right() -> int

	t := state.NewTable()

	t.RawSetString("top", golua.LString(top))
	t.RawSetString("bottom", golua.LString(bottom))
	t.RawSetString("left", golua.LString(left))
	t.RawSetString("right", golua.LString(right))
	t.RawSetString("topleft", golua.LString(topleft))
	t.RawSetString("topright", golua.LString(topright))
	t.RawSetString("bottomleft", golua.LString(bottomleft))
	t.RawSetString("bottomright", golua.LString(bottomright))
	t.RawSetString("middleleft", golua.LString(middleleft))
	t.RawSetString("middleright", golua.LString(middleright))
	t.RawSetString("middle", golua.LString(middle))
	t.RawSetString("middletop", golua.LString(middletop))
	t.RawSetString("middlebottom", golua.LString(middlebottom))

	t.RawSetString("size_top", state.NewFunction(func(state *golua.LState) int {
		b := lgBorderBuild(t)
		size := b.GetTopSize()

		state.Push(golua.LNumber(size))
		return 1
	}))

	t.RawSetString("size_bottom", state.NewFunction(func(state *golua.LState) int {
		b := lgBorderBuild(t)
		size := b.GetBottomSize()

		state.Push(golua.LNumber(size))
		return 1
	}))

	t.RawSetString("size_left", state.NewFunction(func(state *golua.LState) int {
		b := lgBorderBuild(t)
		size := b.GetLeftSize()

		state.Push(golua.LNumber(size))
		return 1
	}))

	t.RawSetString("size_right", state.NewFunction(func(state *golua.LState) int {
		b := lgBorderBuild(t)
		size := b.GetRightSize()

		state.Push(golua.LNumber(size))
		return 1
	}))

	return t
}

func lgBorderBuild(t *golua.LTable) lipgloss.Border {
	top := t.RawGetString("top").(golua.LString)
	bottom := t.RawGetString("bottom").(golua.LString)
	left := t.RawGetString("left").(golua.LString)
	right := t.RawGetString("right").(golua.LString)
	topleft := t.RawGetString("topleft").(golua.LString)
	topright := t.RawGetString("topright").(golua.LString)
	bottomleft := t.RawGetString("bottomleft").(golua.LString)
	bottomright := t.RawGetString("bottomright").(golua.LString)
	middleleft := t.RawGetString("middleleft").(golua.LString)
	middleright := t.RawGetString("middleright").(golua.LString)
	middle := t.RawGetString("middle").(golua.LString)
	middletop := t.RawGetString("middletop").(golua.LString)
	middlebottom := t.RawGetString("middlebottom").(golua.LString)
	return lipgloss.Border{
		Top:          string(top),
		Bottom:       string(bottom),
		Left:         string(left),
		Right:        string(right),
		TopLeft:      string(topleft),
		TopRight:     string(topright),
		BottomLeft:   string(bottomleft),
		BottomRight:  string(bottomright),
		MiddleLeft:   string(middleleft),
		MiddleRight:  string(middleright),
		Middle:       string(middle),
		MiddleTop:    string(middletop),
		MiddleBottom: string(middlebottom),
	}
}

func whitespaceOptionTable(state *golua.LState, lib *lua.Lib) *golua.LTable {
	/// @struct WhitespaceOption
	/// @method foreground(color struct<lipgloss.ColorAny>) -> self
	/// @method background(color struct<lipgloss.ColorAny>) -> self
	/// @method chars(string) -> self

	t := state.NewTable()

	t.RawSetString("__foreground", golua.LNil)
	t.RawSetString("__background", golua.LNil)
	t.RawSetString("__chars", golua.LNil)

	lib.BuilderFunction(state, t, "foreground",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__foreground", args["color"].(*golua.LTable))
		})

	lib.BuilderFunction(state, t, "background",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "color"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__background", args["color"].(*golua.LTable))
		})

	lib.BuilderFunction(state, t, "chars",
		[]lua.Arg{
			{Type: lua.STRING, Name: "str"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			t.RawSetString("__chars", golua.LString(args["str"].(string)))
		})

	return t
}

func whitespaceOptionBuild(t *golua.LTable) []lipgloss.WhitespaceOption {
	opts := []lipgloss.WhitespaceOption{}

	foreground := t.RawGetString("__foreground")
	if foreground.Type() == golua.LTTable {
		opts = append(opts, lipgloss.WithWhitespaceForeground(lgColorGenericBuild(foreground.(*golua.LTable))))
	}

	background := t.RawGetString("__background")
	if background.Type() == golua.LTTable {
		opts = append(opts, lipgloss.WithWhitespaceBackground(lgColorGenericBuild(background.(*golua.LTable))))
	}

	chars := t.RawGetString("__chars")
	if chars.Type() == golua.LTString {
		opts = append(opts, lipgloss.WithWhitespaceChars(string(chars.(golua.LString))))
	}

	return opts
}

func lipglossStyleTable(state *golua.LState, lib *lua.Lib, r *lua.Runner, id int) *golua.LTable {
	/// @struct Style
	/// @prop id {int} - The ID of the style.
	/// @method render(string...) -> string
	/// @method align() -> float<lipgloss.Position>, float<lipgloss.Position>
	/// @method align_set(hpos float<lipgloss.Position>, vpos float<lipgloss.Position>) -> self
	/// @method align_unset() -> self
	/// @method align_horizontal() -> float<lipgloss.Position>
	/// @method align_horizontal_set(hpos float<lipgloss.Position>) -> self
	/// @method align_horizontal_unset() -> self
	/// @method align_vertical() -> float<lipgloss.Position>
	/// @method align_vertical_set(vpos float<lipgloss.Position>) -> self
	/// @method align_vertical_unset() -> self
	/// @method blink() -> bool
	/// @method blink_set(enabled bool) -> self
	/// @method blink_unset() -> self
	/// @method bold() -> bool
	/// @method bold_set(enabled bool) -> self
	/// @method bold_unset() -> self
	/// @method faint() -> bool
	/// @method faint_set(enabled bool) -> self
	/// @method faint_unset() -> self
	/// @method italic() -> bool
	/// @method italic_set(enabled bool) -> self
	/// @method italic_unset() -> self
	/// @method underline() -> bool
	/// @method underline_set(enabled bool) -> self
	/// @method underline_unset() -> self
	/// @method underline_spaces() -> bool
	/// @method underline_spaces_set(enabled bool) -> self
	/// @method underline_spaces_unset() -> self
	/// @method strikethrough() -> bool
	/// @method strikethrough_set(enabled bool) -> self
	/// @method strikethrough_unset() -> self
	/// @method strikethrough_spaces() -> bool
	/// @method strikethrough_spaces_set(enabled bool) -> self
	/// @method strikethrough_spaces_unset() -> self
	/// @method reverse() -> bool
	/// @method reverse_set(enabled bool) -> self
	/// @method reverse_unset() -> self
	/// @method foreground() -> struct<lipgloss.ColorAny>
	/// @method foreground_set(color struct<lipgloss.ColorAny>) -> self
	/// @method foreground_unset() -> self
	/// @method background() -> struct<lipgloss.ColorAny>
	/// @method background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method background_unset() -> self
	/// @method inline() -> bool
	/// @method inline_set(enabled bool) -> self
	/// @method inline_unset() -> self
	/// @method width() -> int
	/// @method width_set(width int) -> self
	/// @method width_unset() -> self
	/// @method height() -> int
	/// @method height_set(height int) -> self
	/// @method height_unset() -> self
	/// @method width_max() -> int
	/// @method width_max_set(width int) -> self
	/// @method width_max_unset() -> self
	/// @method height_max() -> int
	/// @method height_max_set(height int) -> self
	/// @method height_max_unset() -> self
	/// @method tab_width() -> int
	/// @method tab_width_set(width int) -> self
	/// @method tab_width_unset() -> self
	/// @method border() -> struct<lipgloss.Border>, bool, bool, bool, bool
	/// @method border_set(border struct<lipgloss.Border>, sides bool...) -> self
	/// @method border_unset() -> self
	/// @method border_foreground_set(color struct<lipgloss.ColorAny>...) -> self
	/// @method border_background_set(color struct<lipgloss.ColorAny>...) -> self
	/// @method border_style() -> struct<lipgloss.Border>
	/// @method border_style_set(border struct<lipgloss.Border>) -> self
	/// @method border_style_unset() -> self
	/// @method border_top() -> bool
	/// @method border_top_set(enabled bool) -> self
	/// @method border_top_unset() -> self
	/// @method border_top_foreground() -> struct<lipgloss.ColorAny>
	/// @method border_top_foreground_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_top_foreground_unset() -> self
	/// @method border_top_background() -> struct<lipgloss.ColorAny>
	/// @method border_top_background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_top_background_unset() -> self
	/// @method border_bottom() -> bool
	/// @method border_bottom_set(enabled bool) -> self
	/// @method border_bottom_unset() -> self
	/// @method border_bottom_foreground() -> struct<lipgloss.ColorAny>
	/// @method border_bottom_foreground_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_bottom_foreground_unset() -> self
	/// @method border_bottom_background() -> struct<lipgloss.ColorAny>
	/// @method border_bottom_background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_bottom_background_unset() -> self
	/// @method border_left() -> bool
	/// @method border_left_set(enabled bool) -> self
	/// @method border_left_unset() -> self
	/// @method border_left_foreground() -> struct<lipgloss.ColorAny>
	/// @method border_left_foreground_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_left_foreground_unset() -> self
	/// @method border_left_background() -> struct<lipgloss.ColorAny>
	/// @method border_left_background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_left_background_unset() -> self
	/// @method border_right() -> bool
	/// @method border_right_set(enabled bool) -> self
	/// @method border_right_unset() -> self
	/// @method border_right_foreground() -> struct<lipgloss.ColorAny>
	/// @method border_right_foreground_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_right_foreground_unset() -> self
	/// @method border_right_background() -> struct<lipgloss.ColorAny>
	/// @method border_right_background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method border_right_background_unset() -> self
	/// @method margin() -> int, int, int, int
	/// @method margin_set(sides int...) -> self
	/// @method margin_unset() -> self
	/// @method margin_top() -> int
	/// @method margin_top_set(size int) -> self
	/// @method margin_top_unset() -> self
	/// @method margin_bottom() -> int
	/// @method margin_bottom_set(size int) -> self
	/// @method margin_bottom_unset() -> self
	/// @method margin_left() -> int
	/// @method margin_left_set(size int) -> self
	/// @method margin_left_unset() -> self
	/// @method margin_right() -> int
	/// @method margin_right_set(size int) -> self
	/// @method margin_right_unset() -> self
	/// @method margin_background_set(color struct<lipgloss.ColorAny>) -> self
	/// @method margin_background_unset() -> self
	/// @method padding() -> int, int, int, int
	/// @method padding_set(sides int...) -> self
	/// @method padding_unset() -> self
	/// @method padding_top() -> int
	/// @method padding_top_set(size int) -> self
	/// @method padding_top_unset() -> self
	/// @method padding_bottom() -> int
	/// @method padding_bottom_set(size int) -> self
	/// @method padding_bottom_unset() -> self
	/// @method padding_left() -> int
	/// @method padding_left_set(size int) -> self
	/// @method padding_left_unset() -> self
	/// @method padding_right() -> int
	/// @method padding_right_set(size int) -> self
	/// @method padding_right_unset() -> self
	/// @method transform_set(transform func(string) -> string) -> self
	/// @method transform_unset() -> self

	t := state.NewTable()

	t.RawSetString("id", golua.LNumber(id))

	lib.TableFunction(state, t, "render",
		[]lua.Arg{
			lua.ArgVariadic("strs", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			strs := args["strs"].([]any)
			strList := make([]string, len(strs))
			for i, v := range strs {
				strList[i] = v.(string)
			}

			state.Push(golua.LString(item.Style.Render(strList...)))
			return 1
		})

	lib.TableFunction(state, t, "align",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			hpos := item.Style.GetAlignHorizontal()
			vpos := item.Style.GetAlignVertical()

			state.Push(golua.LNumber(hpos))
			state.Push(golua.LNumber(vpos))
			return 2
		})

	lib.BuilderFunction(state, t, "align_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "hpos"},
			{Type: lua.FLOAT, Name: "vpos"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Align(lipgloss.Position(args["hpos"].(float64)), lipgloss.Position(args["vpos"].(float64)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "align_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetAlign()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "align_horizontal",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			hpos := item.Style.GetAlignHorizontal()

			state.Push(golua.LNumber(hpos))
			return 1
		})

	lib.BuilderFunction(state, t, "align_horizontal_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "hpos"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.AlignHorizontal(lipgloss.Position(args["hpos"].(float64)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "align_horizontal_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetAlignHorizontal()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "align_vertical",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			vpos := item.Style.GetAlignVertical()

			state.Push(golua.LNumber(vpos))
			return 1
		})

	lib.BuilderFunction(state, t, "align_vertical_set",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "vpos"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.AlignVertical(lipgloss.Position(args["vpos"].(float64)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "align_vertical_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetAlignVertical()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "blink",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBlink()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "blink_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Blink(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "blink_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBlink()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "bold",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBold()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "bold_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Bold(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "bold_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBold()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "faint",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetFaint()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "faint_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Faint(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "faint_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetFaint()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "italic",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetItalic()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "italic_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Italic(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "italic_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetItalic()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "underline",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetUnderline()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "underline_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Underline(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "underline_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetUnderline()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "underline_spaces",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetUnderlineSpaces()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "underline_spaces_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnderlineSpaces(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "underline_spaces_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetUnderlineSpaces()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "strikethrough",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetStrikethrough()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "strikethrough_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Strikethrough(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "strikethrough_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetStrikethrough()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "strikethrough_spaces",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetStrikethroughSpaces()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "strikethrough_spaces_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.StrikethroughSpaces(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "strikethrough_spaces_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetStrikethroughSpaces()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "reverse",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetReverse()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "reverse_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Reverse(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "reverse_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetReverse()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "foreground",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetForeground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "foreground_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Foreground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "foreground_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetForeground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "background",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBackground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Background(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "inline",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetInline()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "inline_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Inline(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "inline_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetInline()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "width",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetWidth()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "width_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Width(args["width"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "width_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetWidth()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "height",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetHeight()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "height_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Height(args["height"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "height_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetHeight()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "width_max",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMaxWidth()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "width_max_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MaxWidth(args["width"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "width_max_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMaxWidth()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "height_max",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMaxHeight()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "height_max_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "height"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MaxHeight(args["height"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "height_max_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMaxHeight()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "tab_width",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetTabWidth()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "tab_width_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "width"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.TabWidth(args["width"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "tab_width_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetTabWidth()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			border, top, right, bottom, left := item.Style.GetBorder()

			state.Push(lgBorderTableFrom(state, border))
			state.Push(golua.LBool(top))
			state.Push(golua.LBool(right))
			state.Push(golua.LBool(bottom))
			state.Push(golua.LBool(left))
			return 5
		})

	lib.BuilderFunction(state, t, "border_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "border"},
			lua.ArgVariadic("sides", lua.ArrayType{Type: lua.BOOL}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			sides := args["sides"].([]any)
			sideList := make([]bool, len(sides))
			for i, v := range sides {
				sideList[i] = v.(bool)
			}

			border := lgBorderBuild(args["border"].(*golua.LTable))

			newStyle := item.Style.Border(border, sideList...)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_foreground_set",
		[]lua.Arg{
			lua.ArgVariadic("col", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			col := args["col"].([]any)
			colList := make([]lipgloss.TerminalColor, len(col))
			for i, v := range col {
				colList[i] = lgColorGenericBuild(v.(*golua.LTable))
			}

			newStyle := item.Style.BorderForeground(colList...)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_background_set",
		[]lua.Arg{
			lua.ArgVariadic("col", lua.ArrayType{Type: lua.RAW_TABLE}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			col := args["col"].([]any)
			colList := make([]lipgloss.TerminalColor, len(col))
			for i, v := range col {
				colList[i] = lgColorGenericBuild(v.(*golua.LTable))
			}

			newStyle := item.Style.BorderBackground(colList...)
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_style",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			border := item.Style.GetBorderStyle()

			state.Push(lgBorderTableFrom(state, border))
			return 1
		})

	lib.BuilderFunction(state, t, "border_style_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "border"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			border := lgBorderBuild(args["border"].(*golua.LTable))

			newStyle := item.Style.BorderStyle(border)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_style_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderStyle()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_top",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderTop()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_top_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderTop(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_top_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderTop()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_top_foreground",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderTopForeground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_top_foreground_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderTopForeground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_top_foreground_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderTopForeground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_top_background",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderTopBackground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_top_background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderTopBackground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_top_background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderTopBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderBottom()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_bottom_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderBottom(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_bottom_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderBottom()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_bottom_foreground",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderBottomForeground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_bottom_foreground_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderBottomForeground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_bottom_foreground_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderBottomForeground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_bottom_background",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderBottomBackground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_bottom_background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderBottomBackground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_bottom_background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderBottomBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_left",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderLeft()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_left_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderLeft(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_left_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderLeft()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_left_foreground",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderLeftForeground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_left_foreground_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderLeftForeground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_left_foreground_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderLeftForeground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_left_background",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderLeftBackground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_left_background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderLeftBackground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_left_background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderLeftBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_right",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderRight()

			state.Push(golua.LBool(value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_right_set",
		[]lua.Arg{
			{Type: lua.BOOL, Name: "enabled"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderRight(args["enabled"].(bool))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_right_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderRight()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_right_foreground",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderRightForeground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_right_foreground_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderRightForeground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_right_foreground_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderRightForeground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "border_right_background",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetBorderRightBackground()

			state.Push(lgColorGenericTable(state, value))
			return 1
		})

	lib.BuilderFunction(state, t, "border_right_background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.BorderRightBackground(lgColorGenericBuild(args["col"].(*golua.LTable)))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "border_right_background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetBorderRightBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "margin",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			top, right, bottom, left := item.Style.GetMargin()

			state.Push(golua.LNumber(top))
			state.Push(golua.LNumber(right))
			state.Push(golua.LNumber(bottom))
			state.Push(golua.LNumber(left))
			return 4
		})

	lib.BuilderFunction(state, t, "margin_set",
		[]lua.Arg{
			lua.ArgVariadic("margins", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			margins := args["margins"].([]any)
			marginsList := make([]int, len(margins))
			for i, v := range margins {
				marginsList[i] = v.(int)
			}

			newStyle := item.Style.Margin(marginsList...)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMargins()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "margin_top",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMarginTop()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "margin_top_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "margin"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MarginTop(args["margin"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_top_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMarginTop()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "margin_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMarginBottom()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "margin_bottom_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "margin"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MarginBottom(args["margin"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_bottom_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMarginBottom()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "margin_left",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMarginLeft()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "margin_left_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "margin"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MarginLeft(args["margin"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_left_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMarginLeft()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "margin_right",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetMarginRight()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "margin_right_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "margin"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.MarginRight(args["margin"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_right_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMarginRight()
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_background_set",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "col"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			col := lgColorGenericBuild(args["col"].(*golua.LTable))

			newStyle := item.Style.MarginBackground(col)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "margin_background_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetMarginBackground()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "padding",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			top, right, bottom, left := item.Style.GetPadding()

			state.Push(golua.LNumber(top))
			state.Push(golua.LNumber(right))
			state.Push(golua.LNumber(bottom))
			state.Push(golua.LNumber(left))
			return 4
		})

	lib.BuilderFunction(state, t, "padding_set",
		[]lua.Arg{
			lua.ArgVariadic("padding", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			padding := args["padding"].([]any)
			paddingList := make([]int, len(padding))
			for i, v := range padding {
				paddingList[i] = v.(int)
			}

			newStyle := item.Style.Padding(paddingList...)
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "padding_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetPadding()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "padding_top",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetPaddingTop()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "padding_top_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "padding"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.PaddingTop(args["padding"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "padding_top_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetPaddingTop()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "padding_bottom",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetPaddingBottom()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "padding_bottom_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "padding"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.PaddingBottom(args["padding"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "padding_bottom_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetPaddingBottom()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "padding_left",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetPaddingLeft()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "padding_left_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "padding"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.PaddingLeft(args["padding"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "padding_left_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetPaddingLeft()
			item.Style = &newStyle
		})

	lib.TableFunction(state, t, "padding_right",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			value := item.Style.GetPaddingRight()

			state.Push(golua.LNumber(value))
			return 1
		})

	lib.BuilderFunction(state, t, "padding_right_set",
		[]lua.Arg{
			{Type: lua.INT, Name: "padding"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.PaddingRight(args["padding"].(int))
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "padding_right_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetPaddingRight()
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "transform_set",
		[]lua.Arg{
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.Transform(func(s string) string {
				state.Push(args["fn"].(*golua.LFunction))
				state.Push(golua.LString(s))
				state.Call(1, 1)
				str := state.CheckString(-1)
				state.Pop(1)

				return str
			})
			item.Style = &newStyle
		})

	lib.BuilderFunction(state, t, "transform_unset",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			item, _ := r.CR_LIP.Item(int(t.RawGetString("id").(golua.LNumber)))

			newStyle := item.Style.UnsetTransform()
			item.Style = &newStyle
		})

	return t
}
