package lua

import (
	"fmt"
	"math"
	"path"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/AllenDang/giu"

	"github.com/ArtificialLegacy/gm-proj-tool/yyp"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/akamensky/argparse"
	lua "github.com/yuin/gopher-lua"
)

type Runner struct {
	State  *lua.LState
	lg     *log.Logger
	Dir    string
	Output string

	Failed string

	Wg *sync.WaitGroup

	CMDParser *argparse.Parser
	CLIMode   bool

	// -- collections
	TC *collection.Collection[collection.ItemTask]
	IC *collection.Collection[collection.ItemImage]
	FC *collection.Collection[collection.ItemFile]
	CC *collection.Collection[collection.ItemContext]
	QR *collection.Collection[collection.ItemQR]

	// -- crates
	CR_WIN *collection.Crate[giu.MasterWindow]
	CR_REF *collection.Crate[collection.RefItem[any]]
	CR_GMP *collection.Crate[yyp.Project]
}

func NewRunner(state *lua.LState, lg *log.Logger, cliMode bool) Runner {
	wg := &sync.WaitGroup{}
	return Runner{
		State: state,
		lg:    lg,

		Wg: wg,

		CMDParser: argparse.NewParser("imgscal", ""),
		CLIMode:   cliMode,

		// -- collections
		IC: collection.NewCollection[collection.ItemImage](lg, wg),
		FC: collection.NewCollection[collection.ItemFile](lg, wg).OnCollect(
			func(i *collection.Item[collection.ItemFile]) {
				if i.Self != nil {
					i.Self.File.Close()
				}
			}),
		CC: collection.NewCollection[collection.ItemContext](lg, wg),
		QR: collection.NewCollection[collection.ItemQR](lg, wg),
		TC: collection.NewCollection[collection.ItemTask](lg, wg),

		// -- crates
		CR_WIN: collection.NewCrate[giu.MasterWindow](),
		CR_REF: collection.NewCrate[collection.RefItem[any]](),
		CR_GMP: collection.NewCrate[yyp.Project](),
	}
}

func (r *Runner) Run(file string, plugins PluginMap) error {
	defer func() {
		if p := recover(); p != nil {
			r.lg.Append("recovered from panic during lua runtime", log.LEVEL_ERROR)
			r.lg.Append(string(debug.Stack()), log.LEVEL_ERROR)
			r.Failed = fmt.Sprintf("%s", p)
		}
	}()

	r.Dir = path.Dir(file)

	pkg := r.State.GetField(r.State.Get(lua.EnvironIndex), "package")
	r.State.SetField(pkg, "path", lua.LString(r.Dir+"/?.lua"))

	lua.OpenBase(r.State)
	lua.OpenMath(r.State)
	lua.OpenString(r.State)
	lua.OpenTable(r.State)

	err := r.State.DoFile(file)
	if err != nil {
		return err
	}

	initFunc := r.State.GetGlobal("init")
	if initFunc.Type() != lua.LTFunction {
		return fmt.Errorf("failed to run init function, it is not a function: %s", initFunc.Type())
	}
	r.State.Push(initFunc)
	r.State.Push(r.WorkflowInit(file, r.lg, plugins))
	r.State.Call(1, 0)

	mainFunc := r.State.GetGlobal("main")
	if mainFunc.Type() != lua.LTFunction {
		return fmt.Errorf("failed to run main function, it is not a function: %s", mainFunc.Type())
	}
	r.State.Push(mainFunc)
	r.State.Call(0, 0)

	return nil
}

func (r *Runner) Help(file string, wf *workflow.Workflow) (string, error) {
	defer func() {
		if p := recover(); p != nil {
			r.lg.Append("recovered from panic during lua runtime.", log.LEVEL_ERROR)
			r.Failed = fmt.Sprintf("%s", p)
		}
	}()

	r.Dir = path.Dir(file)
	err := r.State.DoFile(file)
	if err != nil {
		return "", err
	}

	helpFunc := r.State.GetGlobal("help")
	if helpFunc.Type() != lua.LTFunction {
		return "", fmt.Errorf("failed to run help function, it is not a function: %s", helpFunc.Type())
	}
	r.State.Push(helpFunc)
	r.State.Push(r.WorkflowInfo(file, wf, r.lg))
	r.State.Call(1, 1)
	str := r.State.CheckString(-1)
	r.State.Pop(1)

	return str, nil
}

func (r *Runner) WorkflowInit(name string, lg *log.Logger, plugins PluginMap) *lua.LTable {
	t := r.State.NewTable()

	t.RawSetString("is_cli", lua.LBool(r.CLIMode))

	t.RawSetString("debug", r.State.NewFunction(func(l *lua.LState) int {
		lua.OpenDebug(r.State)
		return 0
	}))

	t.RawSetString("verbose", r.State.NewFunction(func(l *lua.LState) int {
		r.lg.EnableVerbose()
		return 0
	}))

	t.RawSetString("import", r.State.NewFunction(func(l *lua.LState) int {
		pt := l.CheckTable(-1)
		reqs := []string{}

		pt.ForEach(func(l1, l2 lua.LValue) {
			reqs = append(reqs, string(l2.(lua.LString)))
		})

		LoadPlugins(name, r, lg, plugins, reqs, r.State)

		return 0
	}))

	return t
}

func (r *Runner) WorkflowInfo(name string, wf *workflow.Workflow, lg *log.Logger) *lua.LTable {
	t := r.State.NewTable()

	t.RawSetString("is_cli", lua.LBool(r.CLIMode))

	t.RawSetString("name", lua.LString(wf.Name))
	t.RawSetString("author", lua.LString(wf.Author))
	t.RawSetString("version", lua.LString(wf.Version))
	t.RawSetString("desc", lua.LString(wf.Desc))

	return t
}

type Lib struct {
	Lib    string
	Runner *Runner
	State  *lua.LState
	Lg     *log.Logger
}

func NewLib(lib string, runner *Runner, state *lua.LState, lg *log.Logger) (*Lib, *lua.LTable) {
	t := state.NewTable()
	state.SetGlobal(lib, t)

	return &Lib{
		Lib:    lib,
		State:  state,
		Lg:     lg,
		Runner: runner,
	}, t
}

type ArgType int

const (
	INT ArgType = iota
	FLOAT
	BOOL
	STRING
	TABLE
	ARRAY
	ANY
	FUNC
	RAW_TABLE
)

type Arg struct {
	Type     ArgType
	Name     string
	Optional bool
	Table    *[]Arg
}

type ArrayType struct {
	Type  ArgType
	Table *[]Arg
}

func ArgArray(name string, arrType ArrayType, optional bool) Arg {
	table := []Arg{
		{Type: arrType.Type, Name: name, Table: arrType.Table},
	}

	return Arg{
		Type:     ARRAY,
		Name:     name,
		Optional: optional,
		Table:    &table,
	}
}

func (l *Lib) ParseValue(state *lua.LState, pos int, value lua.LValue, arg Arg, argMap map[string]any) map[string]any {
	switch arg.Type {
	case INT:
		if value.Type() == lua.LTNumber {
			argMap[arg.Name] = int(math.Round(float64(value.(lua.LNumber))))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case FLOAT:
		if value.Type() == lua.LTNumber {
			argMap[arg.Name] = float64(value.(lua.LNumber))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case BOOL:
		if value.Type() == lua.LTBool {
			argMap[arg.Name] = bool(value.(lua.LBool))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid bool provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case STRING:
		if value.Type() == lua.LTString {
			argMap[arg.Name] = string(value.(lua.LString))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid string provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case ANY:
		argMap[arg.Name] = value

	case FUNC:
		if value.Type() == lua.LTFunction {
			argMap[arg.Name] = value
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid function provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case RAW_TABLE:
		if value.Type() == lua.LTTable {
			argMap[arg.Name] = value
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid table provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case ARRAY:
		if value.Type() == lua.LTTable {
			v := value.(*lua.LTable)

			m := map[string]any{}

			for i := range v.Len() {
				m = l.ParseValue(state, pos, v.RawGetInt(i+1), Arg{
					Type:  (*arg.Table)[0].Type,
					Name:  strconv.Itoa(i + 1),
					Table: (*arg.Table)[0].Table,
				}, m)
			}
			argMap[arg.Name] = m
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid array provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case TABLE:
		if value.Type() == lua.LTTable {
			v := value.(*lua.LTable)
			m := map[string]any{}

			for _, a := range *arg.Table {
				m = l.ParseValue(state, pos, v.RawGetString(a.Name), a, m)
			}

			argMap[arg.Name] = m
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid table provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	default:
		state.ArgError(pos, l.Lg.Append(fmt.Sprintf("unsupport arg type provided: %T", value), log.LEVEL_ERROR))
	}

	return argMap
}

func (l *Lib) ParseArgs(state *lua.LState, name string, args []Arg, ln, level int) (map[string]any, int) {
	argMap := map[string]any{}
	count := 0

	for i, a := range args {
		ind := -ln + i
		v := state.CheckAny(ind)
		if i >= ln {
			if !a.Optional {
				state.ArgError(i, l.Lg.Append(fmt.Sprintf("required arg not provided: %d", i), log.LEVEL_ERROR))
			} else {
				argMap[a.Name] = l.getDefault(a, state)
			}
		} else {
			argMap = l.ParseValue(state, i, v, a, argMap)
			count++
		}
	}

	return argMap, count
}

func (l *Lib) getDefault(a Arg, state *lua.LState) any {
	switch a.Type {
	case INT:
		return 0

	case FLOAT:
		return 0.0

	case BOOL:
		return false

	case STRING:
		return ""

	case TABLE:
		tab := map[string]any{}
		for _, v := range *a.Table {
			tab[v.Name] = l.getDefault(v, state)
		}
		return tab

	case ARRAY:
		return map[string]any{}

	case RAW_TABLE:
		return state.NewTable()

	case ANY:
		return lua.LNil
	default:
		return nil
	}
}

type TaskData struct {
	Lib  string
	Name string
}

func (l *Lib) CreateFunction(lib lua.LValue, name string, args []Arg, fn func(state *lua.LState, d TaskData, args map[string]any) int) {
	l.State.SetField(lib, name, l.State.NewFunction(func(state *lua.LState) int {
		l.Lg.Append(fmt.Sprintf("%s.%s called.", l.Lib, name), log.LEVEL_VERBOSE)

		argMap, c := l.ParseArgs(state, name, args, state.GetTop(), 0)
		state.Pop(c)

		ret := fn(state, TaskData{Lib: l.Lib, Name: name}, argMap)

		l.Lg.Append(fmt.Sprintf("%s.%s finished.", l.Lib, name), log.LEVEL_VERBOSE)
		return ret
	}))
}

type PluginMap map[string]func(r *Runner, lg *log.Logger)

func LoadPlugins(to string, r *Runner, lg *log.Logger, plugins PluginMap, reqs []string, state *lua.LState) {
	for _, req := range reqs {
		builtin, ok := plugins[req]
		if !ok {
			lg.Append(fmt.Sprintf("%s: plugin %s does not exist", to, req), log.LEVEL_WARN)
		} else {
			builtin(r, lg)
			lg.Append(fmt.Sprintf("%s: registered plugin %s", to, req), log.LEVEL_SYSTEM)
		}
	}
}
