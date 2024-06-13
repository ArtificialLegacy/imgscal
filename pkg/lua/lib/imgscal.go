package lib

import (
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
)

var Builtins = map[string]func(r *lua.Runner, lg *log.Logger){
	LIB_CLI:   RegisterCli,
	LIB_IMAGE: RegisterImage,
	LIB_IO:    RegisterIO,
	LIB_STD:   RegisterStd,
	LIB_NSFW:  RegisterNSFW,
	LIB_UUID:  RegisterUUID,
}
