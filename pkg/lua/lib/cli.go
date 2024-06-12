package lib

import (
	"fmt"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/Shopify/go-lua"
)

const LIB_CLI = "cli"

func RegisterCli(r *lua.Runner, lg *log.Logger) {
	r.State.NewTable()

	/// @func print()
	/// @arg msg - the message to print to the console.
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("cli.print called", log.LEVEL_INFO)

		msg, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid question provided to cli.question", log.LEVEL_ERROR))
			state.Error()
		}

		println(msg)
		lg.Append(fmt.Sprintf("lua msg printed: %s", msg), log.LEVEL_INFO)

		return 0
	})
	r.State.SetField(-2, "print")

	/// @func question()
	/// @arg question - the message to be displayed.
	/// @returns string - the answer given by the user
	r.State.PushGoFunction(func(state *golua.State) int {
		lg.Append("cli.question called", log.LEVEL_INFO)

		question, ok := state.ToString(-1)
		if !ok {
			state.PushString(lg.Append("invalid question provided to cli.question", log.LEVEL_ERROR))
			state.Error()
		}

		result, err := cli.Question(question, cli.QuestionOptions{})
		if err != nil {
			state.PushString(lg.Append("invalid answer provided to cli.question", log.LEVEL_ERROR))
			state.Error()
		}

		state.PushString(result)
		return 1
	})
	r.State.SetField(-2, "question")

	/// @constants Control
	/// @const RESET
	r.State.PushString(string(cli.COLOR_RESET))
	r.State.SetField(-2, "RESET")

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
	r.State.PushString(string(cli.COLOR_BLACK))
	r.State.SetField(-2, "BLACK")
	r.State.PushString(string(cli.COLOR_RED))
	r.State.SetField(-2, "RED")
	r.State.PushString(string(cli.COLOR_GREEN))
	r.State.SetField(-2, "GREEN")
	r.State.PushString(string(cli.COLOR_YELLOW))
	r.State.SetField(-2, "YELLOW")
	r.State.PushString(string(cli.COLOR_BLUE))
	r.State.SetField(-2, "BLUE")
	r.State.PushString(string(cli.COLOR_MAGENTA))
	r.State.SetField(-2, "MAGENTA")
	r.State.PushString(string(cli.COLOR_CYAN))
	r.State.SetField(-2, "CYAN")
	r.State.PushString(string(cli.COLOR_WHITE))
	r.State.SetField(-2, "WHITE")

	r.State.PushString(string(cli.COLOR_BRIGHT_BLACK))
	r.State.SetField(-2, "BRIGHT_BLACK")
	r.State.PushString(string(cli.COLOR_BRIGHT_RED))
	r.State.SetField(-2, "BRIGHT_RED")
	r.State.PushString(string(cli.COLOR_BRIGHT_GREEN))
	r.State.SetField(-2, "BRIGHT_GREEN")
	r.State.PushString(string(cli.COLOR_BRIGHT_YELLOW))
	r.State.SetField(-2, "BRIGHT_YELLOW")
	r.State.PushString(string(cli.COLOR_BRIGHT_BLUE))
	r.State.SetField(-2, "BRIGHT_BLUE")
	r.State.PushString(string(cli.COLOR_BRIGHT_MAGENTA))
	r.State.SetField(-2, "BRIGHT_MAGENTA")
	r.State.PushString(string(cli.COLOR_BRIGHT_CYAN))
	r.State.SetField(-2, "BRIGHT_CYAN")
	r.State.PushString(string(cli.COLOR_BRIGHT_WHITE))
	r.State.SetField(-2, "BRIGHT_WHITE")

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
	r.State.PushString(string(cli.COLOR_BACKGROUND_BLACK))
	r.State.SetField(-2, "BACKGROUND_BLACK")
	r.State.PushString(string(cli.COLOR_BACKGROUND_RED))
	r.State.SetField(-2, "BACKGROUND_RED")
	r.State.PushString(string(cli.COLOR_BACKGROUND_GREEN))
	r.State.SetField(-2, "BACKGROUND_GREEN")
	r.State.PushString(string(cli.COLOR_BACKGROUND_YELLOW))
	r.State.SetField(-2, "BACKGROUND_YELLOW")
	r.State.PushString(string(cli.COLOR_BACKGROUND_BLUE))
	r.State.SetField(-2, "BACKGROUND_BLUE")
	r.State.PushString(string(cli.COLOR_BACKGROUND_MAGENTA))
	r.State.SetField(-2, "BACKGROUND_MAGENTA")
	r.State.PushString(string(cli.COLOR_BACKGROUND_CYAN))
	r.State.SetField(-2, "BACKGROUND_CYAN")
	r.State.PushString(string(cli.COLOR_BACKGROUND_WHITE))
	r.State.SetField(-2, "BACKGROUND_WHITE")

	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_BLACK))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_BLACK")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_RED))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_RED")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_GREEN))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_GREEN")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_YELLOW))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_YELLOW")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_BLUE))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_BLUE")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_MAGENTA))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_MAGENTA")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_CYAN))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_CYAN")
	r.State.PushString(string(cli.COLOR_BRIGHT_BACKGROUND_WHITE))
	r.State.SetField(-2, "BRIGHT_BACKGROUND_WHITE")

	/// @constants Styles
	/// @const BOLD
	/// @const UNDERLINE
	/// @const REVERSED
	r.State.PushString(string(cli.COLOR_BOLD))
	r.State.SetField(-2, "BOLD")
	r.State.PushString(string(cli.COLOR_UNDERLINE))
	r.State.SetField(-2, "UNDERLINE")
	r.State.PushString(string(cli.COLOR_REVERSED))
	r.State.SetField(-2, "REVERSED")

	r.State.SetGlobal(LIB_CLI)
}
