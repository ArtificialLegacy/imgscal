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
	LIB_GUIPLOT:     RegisterGUIPlot,
	LIB_GUICODE:     RegisterGUICode,
	LIB_BIT:         RegisterBit,
	LIB_REF:         RegisterRef,
	LIB_NOISE:       RegisterNoise,
	LIB_FILTER:      RegisterFilter,
	LIB_CMD:         RegisterCmd,
	LIB_TEMPLATE:    RegisterTemplate,
	LIB_GAMEMAKER:   RegisterGamemaker,
	LIB_TEST:        RegisterTest,
	LIB_IMGER:       RegisterImger,
	LIB_BASE64:      RegisterBase64,
	LIB_TUI:         RegisterTUI,
	LIB_LIPGLOSS:    RegisterLipGloss,
	LIB_PALETTE:     RegisterPalette,
	LIB_STRINGS:     RegisterStrings,
	LIB_BLEND:       RegisterBlend,
	LIB_CACHE:       RegisterCache,
	LIB_SHADER:      RegisterShader,
	LIB_NET:         RegisterNet,
	LIB_PIPE:        RegisterPipe,
}

func tableBuilderFunc(state *golua.LState, t *golua.LTable, name string, fn func(state *golua.LState, t *golua.LTable)) {
	t.RawSetString(name, state.NewFunction(func(state *golua.LState) int {
		self := state.CheckTable(1)

		fn(state, self)

		state.Push(self)
		return 1
	}))
}

/// @struct WorkflowInit
/// @prop is_cli {bool}
/// @method debug() - Open the lua debug library to the workflow.
/// @method verbose() - Enable verbose logging when running the workflow.
/// @method import([]string<imgscal_Imports>) - Array containing the import names for built-in libraries.
/// @method config(table<any>) - Set the default configuration for the workflow.
/// @method secrets(table<any>) - Set the default secrets for the workflow.
/// @method finish_bell() - Causes the bell control character to be printed when the workflow either finishes or fails.
/// @method use_default_input() - Enable using io.default_input().
/// @method use_default_output() - Enable using io.default_output().

/// @struct WorkflowInfo
/// @prop is_cli {bool}
/// @prop name {string}
/// @prop author {string}
/// @prop version {string}
/// @prop desc {string}
