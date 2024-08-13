package lua

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"

	"github.com/AllenDang/giu"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	lua "github.com/yuin/gopher-lua"
)

type Runner struct {
	State   *lua.LState
	lg      *log.Logger
	Plugins []string
	Dir     string

	// -- collections
	TC *collection.Collection[collection.ItemTask]
	IC *collection.Collection[collection.ItemImage]
	FC *collection.Collection[collection.ItemFile]
	CC *collection.Collection[collection.ItemContext]
	QR *collection.Collection[collection.ItemQR]

	// -- crates
	CR_WIN *collection.Crate[giu.MasterWindow]
	CR_REF *collection.Crate[collection.RefItem[any]]
}

func NewRunner(plugins []string, state *lua.LState, lg *log.Logger) Runner {
	return Runner{
		State:   state,
		lg:      lg,
		Plugins: plugins,

		// -- collections
		IC: collection.NewCollection[collection.ItemImage](lg),
		FC: collection.NewCollection[collection.ItemFile](lg).OnCollect(
			func(i *collection.Item[collection.ItemFile]) {
				if i.Self != nil {
					i.Self.File.Close()
				}
			}),
		CC: collection.NewCollection[collection.ItemContext](lg),
		QR: collection.NewCollection[collection.ItemQR](lg),
		TC: collection.NewCollection[collection.ItemTask](lg),

		// -- crates
		CR_WIN: collection.NewCrate[giu.MasterWindow](),
		CR_REF: collection.NewCrate[collection.RefItem[any]](),
	}
}

func (r *Runner) Run(file string) error {
	pwd, _ := os.Getwd()

	defer func() {
		if p := recover(); p != nil {
			r.lg.Append("recovered from panic during lua runtime.", log.LEVEL_ERROR)
		}
	}()

	pth := path.Join(pwd, file)
	r.Dir = path.Dir(pth)

	err := r.State.DoFile(pth)
	if err != nil {
		return err
	}

	return nil
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
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case FLOAT:
		if value.Type() == lua.LTNumber {
			argMap[arg.Name] = float64(value.(lua.LNumber))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case BOOL:
		if value.Type() == lua.LTBool {
			argMap[arg.Name] = bool(value.(lua.LBool))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid bool provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case STRING:
		if value.Type() == lua.LTString {
			argMap[arg.Name] = string(value.(lua.LString))
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid string provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case ANY:
		argMap[arg.Name] = value

	case FUNC:
		if value.Type() == lua.LTFunction {
			argMap[arg.Name] = value
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid function provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
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
			argMap[arg.Name] = l.getDefault(arg)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid array provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case TABLE:
		if value.Type() == lua.LTTable {
			v := value.(*lua.LTable)

			m := map[string]any{}

			i := 0
			v.ForEach(func(l1, l2 lua.LValue) {
				m = l.ParseValue(state, pos, l2, (*arg.Table)[i], m)
				i++
			})
			argMap[arg.Name] = m
		} else if value.Type() == lua.LTNil && arg.Optional {
			argMap[arg.Name] = l.getDefault(arg)
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
				argMap[a.Name] = l.getDefault(a)
			}
		} else {
			argMap = l.ParseValue(state, i, v, a, argMap)
			count++
		}
	}

	return argMap, count
}

func (l *Lib) getDefault(a Arg) any {
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
			tab[v.Name] = l.getDefault(v)
		}
		return tab

	case ARRAY:
		return map[string]any{}

	case ANY:
		fallthrough
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

func LoadPlugins(to string, r *Runner, lg *log.Logger, plugins map[string]func(r *Runner, lg *log.Logger), reqs []string, state *lua.LState) {
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
