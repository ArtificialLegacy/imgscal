package lua

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"golang.org/x/exp/constraints"
)

func ParseEnum[T constraints.Integer](enum int, enums []T, lib *Lib) T {
	if enum < 0 || enum >= len(enums) {
		lib.State.PushString(lib.Lg.Append(fmt.Sprintf("invalid enum value for %T: %d", enums, enum), log.LEVEL_ERROR))
		lib.State.Error()

		return enums[0]
	}

	return enums[enum]
}
