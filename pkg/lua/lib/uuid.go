package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/google/uuid"
	golua "github.com/yuin/gopher-lua"
)

const LIB_UUID = "uuid"

/// @lib UUID
/// @import uuid
/// @desc
/// Small library for generating UUID strings.

func RegisterUUID(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_UUID, r, r.State, lg)

	/// @func string() -> string
	/// @returns {string} - The generated uuid.
	lib.CreateFunction(tab, "string", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			uuid := uuid.NewString()
			lg.Append(fmt.Sprintf("got uuid: %s", uuid), log.LEVEL_INFO)

			state.Push(golua.LString(uuid))
			return 1
		})
}
