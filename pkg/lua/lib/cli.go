package lib

import (
	"fmt"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/cli"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CLI = "cli"

/// @lib CLI
/// @import cli
/// @desc
/// Library for interacting with the command-line.

func RegisterCli(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CLI, r, r.State, lg)

	/// @func print(msg...)
	/// @arg msg {string...} - The messages to print to the console.
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "print",
		[]lua.Arg{
			lua.ArgVariadic("msg", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			msgs := args["msg"].([]any)
			msg := make([]string, len(msgs))
			for i, v := range msgs {
				msg[i] = v.(string)
			}

			str := strings.Join(msg, ", ")

			fmt.Println(str)
			lg.Append(fmt.Sprintf("lua msg printed: %s", str), log.LEVEL_INFO)
			return 0
		})

	/// @func println(msg)
	/// @arg msg {string} - The message to print to the console.
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "println",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			fmt.Println(args["msg"].(string))
			lg.Append(fmt.Sprintf("lua msg printed: %s", args["msg"]), log.LEVEL_INFO)
			return 0
		})

	/// @func print_float(numbers...)
	/// @arg numbers {float64...} - The numbers to print to the console.
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "print_float",
		[]lua.Arg{
			lua.ArgVariadic("numbers", lua.ArrayType{Type: lua.FLOAT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			nums := args["numbers"].([]any)
			msg := make([]string, len(nums))
			for i, v := range nums {
				msg[i] = fmt.Sprintf("%f", v)
			}

			str := strings.Join(msg, ", ")

			fmt.Println(str)
			lg.Append(fmt.Sprintf("lua msg printed: %s", str), log.LEVEL_INFO)
			return 0
		})

	/// @func print_int(numbers...)
	/// @arg numbers {int...} - The numbers to print to the console.
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "print_int",
		[]lua.Arg{
			lua.ArgVariadic("numbers", lua.ArrayType{Type: lua.INT}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			nums := args["numbers"].([]any)
			msg := make([]string, len(nums))
			for i, v := range nums {
				msg[i] = fmt.Sprintf("%d", v)
			}

			str := strings.Join(msg, ", ")

			fmt.Println(str)
			lg.Append(fmt.Sprintf("lua msg printed: %s", str), log.LEVEL_INFO)
			return 0
		})

	/// @func println_number(number, trunc?)
	/// @arg number {float64} - The number to print to the console.
	/// @arg? trunc {bool}
	/// @desc
	/// This is also including in the log similar to std.log.
	lib.CreateFunction(tab, "println_number",
		[]lua.Arg{
			{Type: lua.FLOAT, Name: "number"},
			{Type: lua.BOOL, Name: "trunc", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			num := args["number"].(float64)

			var msg string
			if args["trunc"].(bool) {
				msg = fmt.Sprintf("%d", int(num))
			} else {
				msg = fmt.Sprintf("%f", num)
			}

			fmt.Println(msg)
			lg.Append(fmt.Sprintf("lua msg printed: %s", msg), log.LEVEL_INFO)
			return 0
		})

	/// @func print_table(table)
	/// @arg table {table<any>} - The table to print to the console.
	lib.CreateFunction(tab, "print_table",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "table"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			printValue(args["table"].(golua.LValue), "", "")
			return 0
		})

	/// @func question(question) -> string
	/// @arg question {string} - The message to be displayed.
	/// @returns {string} - The answer given by the user.
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

	/// @func question_ext(question, options) -> string
	/// @arg question {string} - The message to be displayed.
	/// @arg options {struct<cli.Options>} - Options used for processing the response.
	/// @returns {string} - The answer given by the user.
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
			/// @struct Options
			/// @prop normalize {bool} - Set to lowercase the recieved answer.
			/// @prop accepts {[]string} - List of accepted responses.
			/// @prop fallback {string} - A default response to return when the one entered by the user is not in 'accepts'.

			acc := args["options"].(map[string]any)["accepts"].([]any)
			accepts := make([]string, len(acc))
			for i, v := range acc {
				accepts[i] = v.(string)
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

	/// @func confirm(msg)
	/// @arg msg {string} - The message to be displayed.
	/// @desc
	/// Waits for the user to press enter before continuing.
	lib.CreateFunction(tab, "confirm",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			cli.Question(fmt.Sprintf("%s [ENTER]", args["msg"].(string)), cli.QuestionOptions{})
			return 0
		})

	/// @func select(msg, options) -> int
	/// @arg msg {string} - The message to be displayed.
	/// @arg options {[]string} - List of options for the user to select from.
	/// @returns {int} - The index of selected option, or 0 if none were picked.
	lib.CreateFunction(tab, "select",
		[]lua.Arg{
			{Type: lua.STRING, Name: "msg"},
			lua.ArgArray("options", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			options := args["options"].([]any)
			opts := make([]string, len(options))
			for i, v := range options {
				opts[i] = v.(string)
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

	/// @func color_256(code) -> string
	/// @arg code {int}
	/// @returns {string} - The color control code.
	lib.CreateFunction(tab, "color_256",
		[]lua.Arg{
			{Type: lua.INT, Name: "code"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			code := args["code"].(int)
			state.Push(golua.LString(fmt.Sprintf("\u001b[38;5;%dm", code)))
			return 1
		})

	/// @func color_bg_256(code) -> string
	/// @arg code {int}
	/// @returns {string} - The color control code.
	lib.CreateFunction(tab, "color_bg_256",
		[]lua.Arg{
			{Type: lua.INT, Name: "code"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			code := args["code"].(int)
			state.Push(golua.LString(fmt.Sprintf("\u001b[48;5;%dm", code)))
			return 1
		})

	/// @constants Control
	/// @const RESET
	tab.RawSetString("RESET", golua.LString(cli.COLOR_RESET))

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
	tab.RawSetString("BLACK", golua.LString(cli.COLOR_BLACK))
	tab.RawSetString("RED", golua.LString(cli.COLOR_RED))
	tab.RawSetString("GREEN", golua.LString(cli.COLOR_GREEN))
	tab.RawSetString("YELLOW", golua.LString(cli.COLOR_YELLOW))
	tab.RawSetString("BLUE", golua.LString(cli.COLOR_BLUE))
	tab.RawSetString("MAGENTA", golua.LString(cli.COLOR_MAGENTA))
	tab.RawSetString("CYAN", golua.LString(cli.COLOR_CYAN))
	tab.RawSetString("WHITE", golua.LString(cli.COLOR_WHITE))

	tab.RawSetString("BRIGHT_BLACK", golua.LString(cli.COLOR_BRIGHT_BLACK))
	tab.RawSetString("BRIGHT_RED", golua.LString(cli.COLOR_BRIGHT_RED))
	tab.RawSetString("BRIGHT_GREEN", golua.LString(cli.COLOR_BRIGHT_GREEN))
	tab.RawSetString("BRIGHT_YELLOW", golua.LString(cli.COLOR_BRIGHT_YELLOW))
	tab.RawSetString("BRIGHT_BLUE", golua.LString(cli.COLOR_BRIGHT_BLUE))
	tab.RawSetString("BRIGHT_MAGENTA", golua.LString(cli.COLOR_BRIGHT_MAGENTA))
	tab.RawSetString("BRIGHT_CYAN", golua.LString(cli.COLOR_BRIGHT_CYAN))
	tab.RawSetString("BRIGHT_WHITE", golua.LString(cli.COLOR_BRIGHT_WHITE))

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
	tab.RawSetString("BACKGROUND_BLACK", golua.LString(cli.COLOR_BACKGROUND_BLACK))
	tab.RawSetString("BACKGROUND_RED", golua.LString(cli.COLOR_BACKGROUND_RED))
	tab.RawSetString("BACKGROUND_GREEN", golua.LString(cli.COLOR_BACKGROUND_GREEN))
	tab.RawSetString("BACKGROUND_YELLOW", golua.LString(cli.COLOR_BACKGROUND_YELLOW))
	tab.RawSetString("BACKGROUND_BLUE", golua.LString(cli.COLOR_BACKGROUND_BLUE))
	tab.RawSetString("BACKGROUND_MAGENTA", golua.LString(cli.COLOR_BACKGROUND_MAGENTA))
	tab.RawSetString("BACKGROUND_CYAN", golua.LString(cli.COLOR_BACKGROUND_CYAN))
	tab.RawSetString("BACKGROUND_WHITE", golua.LString(cli.COLOR_BACKGROUND_WHITE))

	tab.RawSetString("BRIGHT_BACKGROUND_BLACK", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_BLACK))
	tab.RawSetString("BRIGHT_BACKGROUND_RED", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_RED))
	tab.RawSetString("BRIGHT_BACKGROUND_GREEN", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_GREEN))
	tab.RawSetString("BRIGHT_BACKGROUND_YELLOW", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_YELLOW))
	tab.RawSetString("BRIGHT_BACKGROUND_BLUE", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_BLUE))
	tab.RawSetString("BRIGHT_BACKGROUND_MAGENTA", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_MAGENTA))
	tab.RawSetString("BRIGHT_BACKGROUND_CYAN", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_CYAN))
	tab.RawSetString("BRIGHT_BACKGROUND_WHITE", golua.LString(cli.COLOR_BRIGHT_BACKGROUND_WHITE))

	/// @constants Styles
	/// @const BOLD
	/// @const UNDERLINE
	/// @const REVERSED
	tab.RawSetString("BOLD", golua.LString(cli.COLOR_BOLD))
	tab.RawSetString("UNDERLINE", golua.LString(cli.COLOR_UNDERLINE))
	tab.RawSetString("REVERSED", golua.LString(cli.COLOR_REVERSED))
}

func printValue(val golua.LValue, prefix, indent string) {
	switch val.Type() {
	case golua.LTNil:
		fmt.Printf("%s%snil\n", indent, prefix)
	case golua.LTBool:
		fmt.Printf("%s%s%t\n", indent, prefix, val.(golua.LBool))
	case golua.LTNumber:
		fmt.Printf("%s%s%f\n", indent, prefix, val.(golua.LNumber))
	case golua.LTString:
		fmt.Printf("%s%s%s\n", indent, prefix, val.(golua.LString))
	case golua.LTTable:
		tbl := val.(*golua.LTable)
		fmt.Printf("%s%s{\n", indent, prefix)
		tbl.ForEach(func(k, v golua.LValue) {
			prefix := ""
			if k.Type() == golua.LTString {
				prefix = fmt.Sprintf("%s: ", k.(golua.LString))
			} else if k.Type() == golua.LTNumber {
				prefix = fmt.Sprintf("%d: ", k.(golua.LNumber))
			} else {
				prefix = fmt.Sprintf("%s: ", k.String())
			}
			printValue(v, prefix, indent+"  ")
		})
		fmt.Printf("%s}\n", indent)
	case golua.LTFunction:
		name := val.(*golua.LFunction).String()
		fmt.Printf("%s%s%s\n", indent, prefix, name)
	}
}
