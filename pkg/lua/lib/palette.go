package lib

import (
	imageutil "github.com/ArtificialLegacy/imgscal/pkg/image_util"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_PALETTE = "palette"

/// @lib Palette
/// @import palette
/// @desc
/// A collection of common color palettes.

func RegisterPalette(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BIT, r, r.State, lg)

	/// @func dracula() -> struct<palette.Dracula>
	/// @returns {struct<palette.Dracula>} - The Dracula color palette.
	lib.CreateFunction(tab, "dracula",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := draculaTable(state)

			state.Push(t)
			return 1
		})

	/// @func tokyo_night() -> struct<palette.TokyoNight>
	/// @returns {struct<palette.TokyoNight>} - The Tokyo Night color palette.
	lib.CreateFunction(tab, "tokyo_night",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := tokyoNightTable(state)

			state.Push(t)
			return 1
		})

	/// @func gamemaker() -> struct<palette.GameMaker>
	/// @returns {struct<palette.GameMaker>} - The GameMaker color palette.
	lib.CreateFunction(tab, "gamemaker",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := gamemakerTable(state)

			state.Push(t)
			return 1
		})
}

func draculaTable(state *golua.LState) *golua.LTable {
	/// @struct Dracula
	/// @prop background {struct<image.ColorRGBA>} - The background color.
	/// @prop current_line {struct<image.ColorRGBA>} - The current line color.
	/// @prop foreground {struct<image.ColorRGBA>} - The foreground color.
	/// @prop comment {struct<image.ColorRGBA>} - The comment color.
	/// @prop cyan {struct<image.ColorRGBA>} - The cyan color.
	/// @prop green {struct<image.ColorRGBA>} - The green color.
	/// @prop orange {struct<image.ColorRGBA>} - The orange color.
	/// @prop pink {struct<image.ColorRGBA>} - The pink color.
	/// @prop purple {struct<image.ColorRGBA>} - The purple color.
	/// @prop red {struct<image.ColorRGBA>} - The red color.
	/// @prop yellow {struct<image.ColorRGBA>} - The yellow color.

	t := state.NewTable()

	t.RawSetString("background", imageutil.RGBAToColorTable(state, 40, 42, 54, 255))
	t.RawSetString("current_line", imageutil.RGBAToColorTable(state, 68, 71, 90, 255))
	t.RawSetString("foreground", imageutil.RGBAToColorTable(state, 248, 248, 242, 255))
	t.RawSetString("comment", imageutil.RGBAToColorTable(state, 98, 114, 164, 255))
	t.RawSetString("cyan", imageutil.RGBAToColorTable(state, 139, 233, 253, 255))
	t.RawSetString("green", imageutil.RGBAToColorTable(state, 80, 250, 123, 255))
	t.RawSetString("orange", imageutil.RGBAToColorTable(state, 255, 184, 108, 255))
	t.RawSetString("pink", imageutil.RGBAToColorTable(state, 255, 121, 198, 255))
	t.RawSetString("purple", imageutil.RGBAToColorTable(state, 189, 147, 249, 255))
	t.RawSetString("red", imageutil.RGBAToColorTable(state, 255, 85, 85, 255))
	t.RawSetString("yellow", imageutil.RGBAToColorTable(state, 241, 250, 140, 255))

	return t
}

func tokyoNightTable(state *golua.LState) *golua.LTable {
	/// @struct TokyoNight
	/// @prop background_night {struct<image.ColorRGBA>} - The background color for night.
	/// @prop background_storm {struct<image.ColorRGBA>} - The background color for storm.
	/// @prop terminal_black {struct<image.ColorRGBA>} - The terminal black color.
	/// @prop comment {struct<image.ColorRGBA>} - The comment color.
	/// @prop text {struct<image.ColorRGBA>} - The text color.
	/// @prop foreground {struct<image.ColorRGBA>} - The foreground color.
	/// @prop terminal_white {struct<image.ColorRGBA>} - The terminal white color.
	/// @prop purple {struct<image.ColorRGBA>} - The purple color.
	/// @prop blue {struct<image.ColorRGBA>} - The blue color.
	/// @prop light_blue {struct<image.ColorRGBA>} - The light blue color.
	/// @prop cyan {struct<image.ColorRGBA>} - The cyan color.
	/// @prop light_cyan {struct<image.ColorRGBA>} - The light cyan color.
	/// @prop blue_green {struct<image.ColorRGBA>} - The blue green color.
	/// @prop green {struct<image.ColorRGBA>} - The green color.
	/// @prop white {struct<image.ColorRGBA>} - The white color.
	/// @prop light_orange {struct<image.ColorRGBA>} - The light orange color.
	/// @prop orange {struct<image.ColorRGBA>} - The orange color.
	/// @prop red {struct<image.ColorRGBA>} - The red color.

	t := state.NewTable()

	t.RawSetString("background_night", imageutil.RGBAToColorTable(state, 26, 27, 38, 255))
	t.RawSetString("background_storm", imageutil.RGBAToColorTable(state, 36, 40, 59, 255))
	t.RawSetString("terminal_black", imageutil.RGBAToColorTable(state, 65, 72, 104, 255))
	t.RawSetString("comment", imageutil.RGBAToColorTable(state, 86, 95, 137, 255))
	t.RawSetString("text", imageutil.RGBAToColorTable(state, 154, 165, 206, 255))
	t.RawSetString("foreground", imageutil.RGBAToColorTable(state, 169, 177, 214, 255))
	t.RawSetString("terminal_white", imageutil.RGBAToColorTable(state, 192, 202, 245, 255))
	t.RawSetString("purple", imageutil.RGBAToColorTable(state, 187, 154, 247, 255))
	t.RawSetString("blue", imageutil.RGBAToColorTable(state, 122, 162, 247, 255))
	t.RawSetString("light_blue", imageutil.RGBAToColorTable(state, 125, 207, 255, 255))
	t.RawSetString("cyan", imageutil.RGBAToColorTable(state, 42, 195, 222, 255))
	t.RawSetString("light_cyan", imageutil.RGBAToColorTable(state, 180, 249, 248, 255))
	t.RawSetString("blue_green", imageutil.RGBAToColorTable(state, 115, 218, 202, 255))
	t.RawSetString("green", imageutil.RGBAToColorTable(state, 158, 206, 106, 255))
	t.RawSetString("white", imageutil.RGBAToColorTable(state, 207, 201, 194, 255))
	t.RawSetString("light_orange", imageutil.RGBAToColorTable(state, 224, 175, 104, 255))
	t.RawSetString("orange", imageutil.RGBAToColorTable(state, 255, 158, 100, 255))
	t.RawSetString("red", imageutil.RGBAToColorTable(state, 247, 118, 142, 255))

	return t
}

func gamemakerTable(state *golua.LState) *golua.LTable {
	/// @struct GameMaker
	/// @prop aqua {struct<image.ColorRGBA>} - The aqua color.
	/// @prop black {struct<image.ColorRGBA>} - The black color.
	/// @prop blue {struct<image.ColorRGBA>} - The blue color.
	/// @prop dkgray {struct<image.ColorRGBA>} - The dark gray color.
	/// @prop fuchsia {struct<image.ColorRGBA>} - The fuchsia color.
	/// @prop gray {struct<image.ColorRGBA>} - The gray color.
	/// @prop green {struct<image.ColorRGBA>} - The green color.
	/// @prop lime {struct<image.ColorRGBA>} - The lime color.
	/// @prop ltgray {struct<image.ColorRGBA>} - The light gray color.
	/// @prop maroon {struct<image.ColorRGBA>} - The maroon color.
	/// @prop navy {struct<image.ColorRGBA>} - The navy color.
	/// @prop olive {struct<image.ColorRGBA>} - The olive color.
	/// @prop orange {struct<image.ColorRGBA>} - The orange color.
	/// @prop purple {struct<image.ColorRGBA>} - The purple color.
	/// @prop red {struct<image.ColorRGBA>} - The red color.
	/// @prop silver {struct<image.ColorRGBA>} - The silver color.
	/// @prop teal {struct<image.ColorRGBA>} - The teal color.
	/// @prop white {struct<image.ColorRGBA>} - The white color.
	/// @prop yellow {struct<image.ColorRGBA>} - The yellow color.

	t := state.NewTable()

	t.RawSetString("aqua", imageutil.RGBAToColorTable(state, 0, 255, 255, 255))
	t.RawSetString("black", imageutil.RGBAToColorTable(state, 0, 0, 0, 255))
	t.RawSetString("blue", imageutil.RGBAToColorTable(state, 0, 0, 255, 255))
	t.RawSetString("dkgray", imageutil.RGBAToColorTable(state, 64, 64, 64, 255))
	t.RawSetString("fuchsia", imageutil.RGBAToColorTable(state, 255, 0, 255, 255))
	t.RawSetString("gray", imageutil.RGBAToColorTable(state, 128, 128, 128, 255))
	t.RawSetString("green", imageutil.RGBAToColorTable(state, 0, 128, 0, 255))
	t.RawSetString("lime", imageutil.RGBAToColorTable(state, 0, 255, 0, 255))
	t.RawSetString("ltgray", imageutil.RGBAToColorTable(state, 192, 192, 192, 255))
	t.RawSetString("maroon", imageutil.RGBAToColorTable(state, 128, 0, 0, 255))
	t.RawSetString("navy", imageutil.RGBAToColorTable(state, 0, 0, 128, 255))
	t.RawSetString("olive", imageutil.RGBAToColorTable(state, 128, 128, 0, 255))
	t.RawSetString("orange", imageutil.RGBAToColorTable(state, 255, 160, 64, 255))
	t.RawSetString("purple", imageutil.RGBAToColorTable(state, 128, 0, 128, 255))
	t.RawSetString("red", imageutil.RGBAToColorTable(state, 255, 0, 0, 255))
	t.RawSetString("silver", imageutil.RGBAToColorTable(state, 192, 192, 192, 255))
	t.RawSetString("teal", imageutil.RGBAToColorTable(state, 0, 128, 128, 255))
	t.RawSetString("white", imageutil.RGBAToColorTable(state, 255, 255, 255, 255))
	t.RawSetString("yellow", imageutil.RGBAToColorTable(state, 255, 255, 0, 255))

	return t
}
