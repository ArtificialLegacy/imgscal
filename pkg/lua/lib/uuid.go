package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
	"github.com/google/uuid"
)

const LIB_UUID = "uuid"

func RegisterUUID(r *lua.Runner, lg *log.Logger) {
	lib := lua.NewLib(LIB_UUID, r.State, lg)

	/// @func string()
	/// @returns string - the generated uuid.
	lib.CreateFunction("string", []lua.Arg{},
		func(state *golua.State, args map[string]any) int {
			uuid := uuid.NewString()
			lg.Append(fmt.Sprintf("got uuid: %s", uuid), log.LEVEL_INFO)

			r.State.PushString(uuid)
			return 1
		})
}
