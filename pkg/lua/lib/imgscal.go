package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

/// @lib ImgScal
/// @import ~
/// @desc
/// Automate image processing programmatically.
/// @section
/// Built around concurrency.
/// Workflows writteng in lua.
/// Builtin handling for image encodings and color models.
/// Included ImGui wrapper for building custom GUI tools.
/// Spritesheet support.
/// Command-line support, e.g. imgscal resize ./image.png 100 100.

var Builtins = map[string]func(r *lua.Runner, lg *log.Logger){
	LIB_CLI:         RegisterCli,
	LIB_IMAGE:       RegisterImage,
	LIB_IO:          RegisterIO,
	LIB_STD:         RegisterStd,
	LIB_NSFW:        RegisterNSFW,
	LIB_UUID:        RegisterUUID,
	LIB_ASCII:       RegisterASCII,
	LIB_TXT:         RegisterTXT,
	LIB_COLLECTION:  RegisterCollection,
	LIB_CONTEXT:     RegisterContext,
	LIB_SPRITESHEET: RegisterSpritesheet,
	LIB_QRCODE:      RegisterQRCode,
	LIB_TIME:        RegisterTime,
	LIB_JSON:        RegisterJSON,
	LIB_GUI:         RegisterGUI,
	LIB_BIT:         RegisterBit,
	LIB_REF:         RegisterRef,
	LIB_NOISE:       RegisterNoise,
	LIB_FILTER:      RegisterFilter,
	LIB_CMD:         RegisterCmd,
}

func tableBuilderFunc(state *golua.LState, t *golua.LTable, name string, fn func(state *golua.LState, t *golua.LTable)) {
	t.RawSetString(name, state.NewFunction(func(state *golua.LState) int {
		self := state.CheckTable(1)

		fn(state, self)

		state.Push(self)
		return 1
	}))
}
