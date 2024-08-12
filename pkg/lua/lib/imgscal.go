package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

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
}
