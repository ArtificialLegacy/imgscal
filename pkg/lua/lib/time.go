package lib

import (
	"time"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_TIME = "time"

/// @lib Time
/// @import time
/// @desc
/// Library for getting basic information about the time.

func RegisterTime(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_TIME, r, r.State, lg)

	/// @func now_ms() -> int
	/// @returns {int} - The current time in ms.
	lib.CreateFunction(tab, "now_ms", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano() / int64(time.Millisecond)

			state.Push(golua.LNumber(t))
			return 1
		})

	/// @func now_mc() -> int
	/// @returns {int} - The current time in mc.
	lib.CreateFunction(tab, "now_mc", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().UnixNano() / int64(time.Microsecond)

			state.Push(golua.LNumber(t))
			return 1
		})

	/// @func now_date() -> string
	/// @returns {string} - The current date in MM-DD-YEAR format.
	lib.CreateFunction(tab, "now_date", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format("01-02-2006")

			state.Push(golua.LString(t))
			return 1
		})

	/// @func now_timestamp() -> string
	/// @returns {string} - The current time in the default format.
	lib.CreateFunction(tab, "now_timestamp", []lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().String()

			state.Push(golua.LString(t))
			return 1
		})

	/// @func now_format(format) -> string
	/// @arg format {string} - Go time format string.
	/// @returns {string} - The current time in given format.
	lib.CreateFunction(tab, "now_format",
		[]lua.Arg{
			{Type: lua.STRING, Name: "format"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			t := time.Now().Local().Format(args["format"].(string))

			state.Push(golua.LString(t))
			return 1
		})

	/// @constants Weekdays
	/// @const WEEKDAY_SUNDAY
	/// @const WEEKDAY_MONDAY
	/// @const WEEKDAY_TUESDAY
	/// @const WEEKDAY_WEDNESDAY
	/// @const WEEKDAY_THURSDAY
	/// @const WEEKDAY_FRIDAY
	/// @const WEEKDAY_SATURDAY
	r.State.SetTable(tab, golua.LString("WEEKDAY_SUNDAY"), golua.LNumber(WEEKDAY_SUNDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_MONDAY"), golua.LNumber(WEEKDAY_MONDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_TUESDAY"), golua.LNumber(WEEKDAY_TUESDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_WEDNESDAY"), golua.LNumber(WEEKDAY_WEDNESDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_THURSDAY"), golua.LNumber(WEEKDAY_THURSDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_FRIDAY"), golua.LNumber(WEEKDAY_FRIDAY))
	r.State.SetTable(tab, golua.LString("WEEKDAY_SATURDAY"), golua.LNumber(WEEKDAY_SATURDAY))
}

const (
	WEEKDAY_SUNDAY int = iota
	WEEKDAY_MONDAY
	WEEKDAY_TUESDAY
	WEEKDAY_WEDNESDAY
	WEEKDAY_THURSDAY
	WEEKDAY_FRIDAY
	WEEKDAY_SATURDAY
)
