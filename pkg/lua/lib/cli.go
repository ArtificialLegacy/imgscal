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

	r.State.SetGlobal(LIB_CLI)
}
