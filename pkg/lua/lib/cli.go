package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CLI = "cli"

func RegisterCli(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CLI, r, r.State, lg)

	/// @func print()
	/// @arg msg - the message to print to the console.
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "print",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			println(args["msg"])
			lg.Append(fmt.Sprintf("lua msg printed: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func question()
	/// @arg question - the message to be displayed.
	/// @returns string - the answer given by the user
	lib.CreateFunction(tab, "question",
		[]lua.Arg{
			{Type: lua.STRING, Name: "question"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			result, err := cli.Question(args["question"].(string), cli.QuestionOptions{})
			if err != nil {
				state.Error(golua.LString(lg.Append("invalid answer provided to cli.question", log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(result))
			return 1
		})

	/// @func question_ext()
	/// @arg question - the message to be displayed.
	/// @arg options - the options to use for processing the answer. [normalize, accepts, fallback]
	/// @returns string - the answer given by the user
	lib.CreateFunction(tab, "question_ext",
		[]lua.Arg{
			{Type: lua.STRING, Name: "question"},
			{Type: lua.TABLE, Name: "options", Table: &[]lua.Arg{
				{Type: lua.BOOL, Name: "normalize", Optional: true},
				lua.ArgArray("accepts", lua.ArrayType{Type: lua.STRING}, true),
				{Type: lua.STRING, Name: "fallback", Optional: true},
			}},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			acc := args["options"].(map[string]any)["accepts"]
			accepts := []string{}
			if str, ok := acc.([]string); ok {
				accepts = str
			}

			opts := cli.QuestionOptions{
				Normalize: args["options"].(map[string]any)["normalize"].(bool),
				Accepts:   accepts,
				Fallback:  args["options"].(map[string]any)["fallback"].(string),
			}

			result, err := cli.Question(args["question"].(string), opts)
			if err != nil {
				state.Error(golua.LString(lg.Append(fmt.Sprintf("invalid answer provided to cli.question_ext: %s", err), log.LEVEL_ERROR)), 0)
			}

			state.Push(golua.LString(result))
			return 1
		})

	/// @func confirm()
	/// @arg msg
	/// @desc
	/// waits for enter to be pressed before continuing.
	lib.CreateFunction(tab, "confirm",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cli.Question(fmt.Sprintf("%s [ENTER]", args["msg"].(string)), cli.QuestionOptions{})
			return 0
		})

	/// @func select()
	/// @arg msg
	/// @arg options - array of strings
	/// @returns index of selected option, or 0.
	lib.CreateFunction(tab, "select",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
			lua.ArgArray("options", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			opts := []string{}
			for _, v := range args["options"].(map[string]any) {
				opts = append(opts, v.(string))
			}

			ind, err := cli.SelectMenu(
				args["msg"].(string),
				opts,
			)
			if err != nil {
				lg.Append("selection failed", log.LEVEL_WARN)
			}

			lg.Append(fmt.Sprintf("selection option picked: %d", ind+1), log.LEVEL_INFO)

			state.Push(golua.LNumber(ind + 1))
			return 1
		})

	/// @constants Control
	/// @const RESET
	r.State.SetField(tab, "RESET", golua.LString(cli.COLOR_RESET))

	/// @constants Text Colors
	/// @const BLACK
	/// @const RED
	/// @const GREEN
	/// @const YELLOW
	/// @const BLUE
	/// @const MAGENTA
	/// @const CYAN
	/// @const WHITE
	/// @const BRIGHT_BLACK
	/// @const BRIGHT_RED
	/// @const BRIGHT_GREEN
	/// @const BRIGHT_YELLOW
	/// @const BRIGHT_BLUE
	/// @const BRIGHT_MAGENTA
	/// @const BRIGHT_CYAN
	/// @const BRIGHT_WHITE
	r.State.SetField(tab, "BLACK", golua.LString(cli.COLOR_BLACK))
	r.State.SetField(tab, "RED", golua.LString(cli.COLOR_RED))
	r.State.SetField(tab, "GREEN", golua.LString(cli.COLOR_GREEN))
	r.State.SetField(tab, "YELLOW", golua.LString(cli.COLOR_YELLOW))
	r.State.SetField(tab, "BLUE", golua.LString(cli.COLOR_BLUE))
	r.State.SetField(tab, "MAGENTA", golua.LString(cli.COLOR_MAGENTA))
	r.State.SetField(tab, "CYAN", golua.LString(cli.COLOR_CYAN))
	r.State.SetField(tab, "WHITE", golua.LString(cli.COLOR_WHITE))

	r.State.SetField(tab, "BRIGHT_BLACK", golua.LString(cli.COLOR_BRIGHT_BLACK))
	r.State.SetField(tab, "BRIGHT_RED", golua.LString(cli.COLOR_BRIGHT_RED))
	r.State.SetField(tab, "BRIGHT_GREEN", golua.LString(cli.COLOR_BRIGHT_GREEN))
	r.State.SetField(tab, "BRIGHT_YELLOW", golua.LString(cli.COLOR_BRIGHT_YELLOW))
	r.State.SetField(tab, "BRIGHT_BLUE", golua.LString(cli.COLOR_BRIGHT_BLUE))
	r.State.SetField(tab, "BRIGHT_MAGENTA", golua.LString(cli.COLOR_BRIGHT_MAGENTA))
	r.State.SetField(tab, "BRIGHT_CYAN", golua.LString(cli.COLOR_BRIGHT_CYAN))
	r.State.SetField(tab, "BRIGHT_WHITE", golua.LString(cli.COLOR_BRIGHT_WHITE))

	/// @constants Background Colors
	/// @const BACKGROUND_BLACK
	/// @const BACKGROUND_RED
	/// @const BACKGROUND_GREEN
	/// @const BACKGROUND_YELLOW
	/// @const BACKGROUND_BLUE
	/// @const BACKGROUND_MAGENTA
	/// @const BACKGROUND_CYAN
	/// @const BACKGROUND_WHITE
	/// @const BRIGHT_BACKGROUND_BLACK
	/// @const BRIGHT_BACKGROUND_RED
	/// @const BRIGHT_BACKGROUND_GREEN
	/// @const BRIGHT_BACKGROUND_YELLOW
	/// @const BRIGHT_BACKGROUND_BLUE
	/// @const BRIGHT_BACKGROUND_MAGENTA
	/// @const BRIGHT_BACKGROUND_CYAN
	/// @const BRIGHT_BACKGROUND_WHITE
	r.State.SetField(tab, "BACKGROUND_BLACK", golua.LString(cli.COLOR_BACKGROUND_BLACK))
	r.State.SetField(tab, "BACKGROUND_RED", golua.LString(cli.COLOR_BACKGROUND_RED))
	r.State.SetField(tab, "BACKGROUND_GREEN", golua.LString(cli.COLOR_BACKGROUND_GREEN))
	r.State.SetField(tab, "BACKGROUND_YELLOW", golua.LString(cli.COLOR_BACKGROUND_YELLOW))
	r.State.SetField(tab, "BACKGROUND_BLUE", golua.LString(cli.COLOR_BACKGROUND_BLUE))
	r.State.SetField(tab, "BACKGROUND_MAGENTA", golua.LString(cli.COLOR_BACKGROUND_MAGENTA))
	r.State.SetField(tab, "BACKGROUND_CYAN", golua.LString(cli.COLOR_BACKGROUND_CYAN))
	r.State.SetField(tab, "BACKGROUND_WHITE", golua.LString(cli.COLOR_BACKGROUND_WHITE))

	r.State.SetField(tab, "BRIGHT_BACKGROUND_BLACK", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_BLACK))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_RED", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_RED))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_GREEN", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_GREEN))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_YELLOW", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_YELLOW))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_BLUE", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_BLUE))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_MAGENTA", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_MAGENTA))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_CYAN", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_CYAN))
	r.State.SetField(tab, "BRIGHT_BACKGROUND_WHITE", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_WHITE))

	/// @constants Styles
	/// @const BOLD
	/// @const UNDERLINE
	/// @const REVERSED
	r.State.SetField(tab, "BOLD", golua.LString(cli.COLOR_BOLD))
	r.State.SetField(tab, "UNDERLINE", golua.LString(cli.COLOR_UNDERLINE))
	r.State.SetField(tab, "REVERSED", golua.LString(cli.COLOR_REVERSED))
}
