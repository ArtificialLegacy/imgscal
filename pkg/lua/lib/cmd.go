package lib

import (
	"fmt"
	"os"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/akamensky/argparse"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CMD = "cmd"

/// @lib Command
/// @import cmd
/// @desc
/// Library for parsing command-line arguments.
/// @section
/// Cannot be used when the workflow was called from the workflow selection menu.

func RegisterCmd(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CMD, r, r.State, lg)

	/// @func parse() -> bool, string
	/// @returns {bool} - If the results are valid.
	/// @returns {string} - A string of the error if on occured.
	lib.CreateFunction(tab, "parse",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			if !r.CLIMode {
				state.Push(golua.LFalse)
				state.Push(golua.LString("workflow was not called directly, cannot parse args."))
				return 2
			}

			err := r.CMDParser.Parse(os.Args[1:])
			if err != nil {
				state.Push(golua.LFalse)
				state.Push(golua.LString(err.Error()))
			} else {
				state.Push(golua.LTrue)
				state.Push(golua.LString(""))
			}

			return 2
		})

	/// @func called() -> bool
	/// @returns {bool} - If the workflow was called directly from the command-line.
	lib.CreateFunction(tab, "called",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(golua.LBool(r.CLIMode))
			return 1
		})

	/// @func options() -> struct<cmd.Options>
	/// @returns {struct<cmd.Options>}
	/// @desc
	/// Creates a table for passing options into arguments.
	lib.CreateFunction(tab, "options",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			/// @struct Options
			/// @method required(self, bool) -> self
			/// @method validate(self, {function([]string) -> string}) -> self - Takes in all values for the argument, and returns an error if there was one. Use an empty string for when there is no error.
			/// @method default(self, any, int<cmd.DefaultType>) -> self - The first value should match the default type given.

			t := state.NewTable()

			t.RawSetString("__required", golua.LFalse)
			t.RawSetString("__validate", golua.LNil)
			t.RawSetString("__default", golua.LNil)
			t.RawSetString("__defaultType", golua.LNil)

			tableBuilderFunc(state, t, "required", func(state *golua.LState, t *golua.LTable) {
				b := state.CheckBool(-1)
				t.RawSetString("__required", golua.LBool(b))
			})

			tableBuilderFunc(state, t, "validate", func(state *golua.LState, t *golua.LTable) {
				fn := state.CheckFunction(-1)
				t.RawSetString("__validate", fn)
			})

			tableBuilderFunc(state, t, "default", func(state *golua.LState, t *golua.LTable) {
				v := state.CheckAny(-2)
				dt := state.CheckNumber(-1)
				t.RawSetString("__default", v)
				t.RawSetString("__defaultType", dt)
			})

			state.Push(t)
			return 1
		})

	/// @func arg_flag(short, long, options?) -> int<ref.BOOL>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.BOOL>}
	lib.CreateFunction(tab, "arg_flag",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			b := r.CMDParser.Command.Flag(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: b,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_flag_count(short, long, options?) -> int<ref.INT>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.INT>}
	lib.CreateFunction(tab, "arg_flag_count",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			i := r.CMDParser.Command.FlagCounter(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: i,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_string(short, long, options?) -> int<ref.STRING>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.STRING>}
	lib.CreateFunction(tab, "arg_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			s := r.CMDParser.Command.String(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: s,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_string_list(short, long, options?) -> int<[]ref.STRING>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<[]ref.STRING>}
	lib.CreateFunction(tab, "arg_string_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			s := r.CMDParser.Command.StringList(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: s,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_string_pos(options?) -> int<ref.STRING>
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.STRING>}
	lib.CreateFunction(tab, "arg_string_pos",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			s := r.CMDParser.Command.StringPositional(opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: s,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_int(short, long, options?) -> int<ref.INT>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.INT>}
	lib.CreateFunction(tab, "arg_int",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			i := r.CMDParser.Command.Int(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: i,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_int_list(short, long, options?) -> int<[]ref.INT>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<[]ref.INT>}
	lib.CreateFunction(tab, "arg_int_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			i := r.CMDParser.Command.IntList(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: i,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_int_pos(options?) -> int<ref.INT>
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.INT>}
	lib.CreateFunction(tab, "arg_int_pos",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			i := r.CMDParser.Command.IntPositional(opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: i,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_float(short, long, options?) -> int<ref.FLOAT>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.FLOAT>}
	lib.CreateFunction(tab, "arg_float",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			f := r.CMDParser.Command.Float(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_float_list(short, long, options?) -> int<[]ref.FLOAT>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<[]ref.FLOAT>}
	lib.CreateFunction(tab, "arg_float_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			f := r.CMDParser.Command.FloatList(args["short"].(string), args["long"].(string), opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_float_pos(options?) -> int<[]ref.FLOAT>
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.FLOAT>}
	lib.CreateFunction(tab, "arg_float_pos",
		[]lua.Arg{
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			f := r.CMDParser.Command.FloatPositional(opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_selector(short, long, choices, options?) -> int<ref.STRING>
	/// @arg short {string}
	/// @arg long {string}
	/// @arg choices {[]string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.STRING>}
	lib.CreateFunction(tab, "arg_selector",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			lua.ArgArray("choices", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			chv := args["choices"].([]any)
			choices := make([]string, len(chv))
			for i, v := range chv {
				choices[i] = v.(string)
			}

			f := r.CMDParser.Command.Selector(args["short"].(string), args["long"].(string), choices, opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_selector_pos(choices, options?) -> int<ref.STRING>
	/// @arg choices {[]string}
	/// @arg? options {struct<cmd.Options>}
	/// @returns {int<ref.STRING>}
	lib.CreateFunction(tab, "arg_selector_pos",
		[]lua.Arg{
			lua.ArgArray("choices", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.RAW_TABLE, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			chv := args["choices"].([]any)
			choices := make([]string, len(chv))
			for i, v := range chv {
				choices[i] = v.(string)
			}

			f := r.CMDParser.Command.SelectorPositional(choices, opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @constants DefaultType {int}
	/// @const DEFAULT_INT
	/// @const DEFAULT_FLOAT
	/// @const DEFAULT_STRING
	/// @const DEFAULT_BOOL
	tab.RawSetString("DEFAULT_INT", golua.LString(DEFAULT_INT))
	tab.RawSetString("DEFAULT_FLOAT", golua.LString(DEFAULT_FLOAT))
	tab.RawSetString("DEFAULT_STRING", golua.LString(DEFAULT_STRING))
	tab.RawSetString("DEFAULT_BOOL", golua.LString(DEFAULT_BOOL))
}

const (
	DEFAULT_INT int = iota
	DEFAULT_FLOAT
	DEFAULT_STRING
	DEFAULT_BOOL
)

func buildOptions(state *golua.LState, t *golua.LTable) *argparse.Options {
	reqv := t.RawGetString("__required")
	required := false
	if reqv.Type() == golua.LTBool {
		required = bool(reqv.(golua.LBool))
	}

	var dflt any = nil
	defv := t.RawGetString("__default")
	defvt := t.RawGetString("__defaultType")

	if defvt.Type() == golua.LTNumber {
		switch int(defvt.(golua.LNumber)) {
		case DEFAULT_INT:
			dflt = int(defv.(golua.LNumber))
		case DEFAULT_FLOAT:
			dflt = float64(defv.(golua.LNumber))
		case DEFAULT_STRING:
			dflt = string(defv.(golua.LString))
		case DEFAULT_BOOL:
			dflt = bool(defv.(golua.LBool))
		}
	}

	opt := &argparse.Options{
		Required: required,
		Default:  dflt,
	}

	valv := t.RawGetString("__validate")
	if valv.Type() == golua.LTFunction {
		opt.Validate = func(args []string) error {
			state.Push(valv)
			at := state.NewTable()
			for _, a := range args {
				at.Append(golua.LString(a))
			}
			state.Push(at)
			state.Call(1, 1)
			err := state.CheckString(-1)
			state.Pop(1)
			return fmt.Errorf(err)
		}
	}

	return opt
}
