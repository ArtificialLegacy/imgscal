package lib

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	"github.com/akamensky/argparse"
	golua "github.com/yuin/gopher-lua"
)

const LIB_CMD = "cmd"

func RegisterCmd(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_CMD, r, r.State, lg)

	/// @func parse()
	/// @returns ok
	/// @returns error
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

	/// @func called()
	/// @returns bool - if the workflow was called directly from the cmd line
	lib.CreateFunction(tab, "called",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			state.Push(golua.LBool(r.CLIMode))
			return 1
		})

	/// @func options()
	/// @returns options struct
	lib.CreateFunction(tab, "options",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
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

	/// @func arg_flag()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns boolref
	lib.CreateFunction(tab, "arg_flag",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_flag_count()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns intref
	lib.CreateFunction(tab, "arg_flag_count",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_string()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns stringref
	lib.CreateFunction(tab, "arg_string",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_string_list()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns sliceref of strings
	lib.CreateFunction(tab, "arg_string_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_string_pos()
	/// @arg options
	/// @returns stringref
	lib.CreateFunction(tab, "arg_string_pos",
		[]lua.Arg{
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_int()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns intref
	lib.CreateFunction(tab, "arg_int",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_int_list()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns sliceref of ints
	lib.CreateFunction(tab, "arg_int_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_int_pos()
	/// @arg options
	/// @returns intref
	lib.CreateFunction(tab, "arg_int_pos",
		[]lua.Arg{
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_float()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns floatref
	lib.CreateFunction(tab, "arg_float",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_float_list()
	/// @arg short
	/// @arg long
	/// @arg options
	/// @returns sliceref of floats
	lib.CreateFunction(tab, "arg_float_list",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_float_pos()
	/// @arg options
	/// @returns floatref
	lib.CreateFunction(tab, "arg_float_pos",
		[]lua.Arg{
			{Type: lua.ANY, Name: "options", Optional: true},
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

	/// @func arg_selector()
	/// @arg short
	/// @arg long
	/// @arg choices - []string
	/// @arg options
	/// @returns stringref
	lib.CreateFunction(tab, "arg_selector",
		[]lua.Arg{
			{Type: lua.STRING, Name: "short"},
			{Type: lua.STRING, Name: "long"},
			lua.ArgArray("choices", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.ANY, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			choices := []string{}
			chv := args["choices"].(map[string]any)
			for i := range len(chv) {
				choices = append(choices, chv[strconv.Itoa(i+1)].(string))
			}

			f := r.CMDParser.Command.Selector(args["short"].(string), args["long"].(string), choices, opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})

	/// @func arg_selector_pos()
	/// @arg choices - []string
	/// @arg options
	/// @returns stringref
	lib.CreateFunction(tab, "arg_selector_pos",
		[]lua.Arg{
			lua.ArgArray("choices", lua.ArrayType{Type: lua.STRING}, false),
			{Type: lua.ANY, Name: "options", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			var opt *argparse.Options = nil

			optv := args["options"].(golua.LValue)
			if optv.Type() == golua.LTTable {
				opt = buildOptions(state, optv.(*golua.LTable))
			}

			choices := []string{}
			chv := args["choices"].(map[string]any)
			for i := range len(chv) {
				choices = append(choices, chv[strconv.Itoa(i+1)].(string))
			}

			f := r.CMDParser.Command.SelectorPositional(choices, opt)
			ref := r.CR_REF.Add(&collection.RefItem[any]{
				Value: f,
			})

			state.Push(golua.LNumber(ref))
			return 1
		})
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
