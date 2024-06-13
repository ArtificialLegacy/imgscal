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
	r.State.NewTable()

	/// @func string()
	/// @returns string - the generated uuid.
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("uuid.string called", log.LEVEL_INFO)

		uuid := uuid.NewString()
		lg.Append(fmt.Sprintf("got uuid: %s", uuid), log.LEVEL_INFO)

		r.State.PushString(uuid)
		return 1
	})
	r.State.SetField(-2, "string")

	r.State.SetGlobal(LIB_UUID)
}
