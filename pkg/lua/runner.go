package lua

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/AllenDang/giu"

	"github.com/ArtificialLegacy/gm-proj-tool/yyp"
	"github.com/ArtificialLegacy/imgscal/pkg/collection"
	"github.com/ArtificialLegacy/imgscal/pkg/config"
	teamodels "github.com/ArtificialLegacy/imgscal/pkg/custom_tea/models"
	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/workflow"
	"github.com/akamensky/argparse"
	lua "github.com/yuin/gopher-lua"
)

type Runner struct {
	State            *lua.LState
	lg               *log.Logger
	Dir              string
	Config           *config.Config
	Entry            string
	ConfigData       map[string]any
	SecretData       map[string]any
	UseDefaultInput  bool
	UseDefaultOutput bool

	Failed string

	Wg *sync.WaitGroup

	CMDParser *argparse.Parser
	CLIMode   bool

	// -- collections
	TC *collection.Collection[collection.ItemTask]
	IC *collection.Collection[collection.ItemImage]
	CC *collection.Collection[collection.ItemContext]
	QR *collection.Collection[collection.ItemQR]

	// -- crates
	CR_WIN *collection.Crate[giu.MasterWindow]
	CR_REF *collection.Crate[collection.RefItem[any]]
	CR_GMP *collection.Crate[yyp.Project]
	CR_TEA *collection.Crate[teamodels.TeaItem]
	CR_LIP *collection.Crate[collection.StyleItem]
	CR_CIM *collection.Crate[collection.CachedImageItem]
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
		IC: collection.NewCollection[collection.ItemImage](lg, wg, collection.TYPE_IMAGE),
		CC: collection.NewCollection[collection.ItemContext](lg, wg, collection.TYPE_CONTEXT),
		QR: collection.NewCollection[collection.ItemQR](lg, wg, collection.TYPE_QR),
		TC: collection.NewCollection[collection.ItemTask](lg, wg, collection.TYPE_TASK),

		// -- crates
		CR_WIN: collection.NewCrate[giu.MasterWindow](),
		CR_REF: collection.NewCrate[collection.RefItem[any]](),
		CR_GMP: collection.NewCrate[yyp.Project](),
		CR_TEA: collection.NewCrate[teamodels.TeaItem](),
		CR_LIP: collection.NewCrate[collection.StyleItem](),
		CR_CIM: collection.NewCrate[collection.CachedImageItem](),
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
	r.State.SetField(pkg, "path", lua.LString(fmt.Sprintf("%s/?.lua;%s/?.lua;%s/?/?.lua", r.Dir, r.Config.PluginDirectory, r.Config.PluginDirectory)))

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

	t.RawSetString("config", r.State.NewFunction(func(l *lua.LState) int {
		r.ConfigData = r.WorkflowConfig(l, ".json")
		return 0
	}))

	t.RawSetString("secrets", r.State.NewFunction(func(l *lua.LState) int {
		r.SecretData = r.WorkflowConfig(l, ".secrets.json")
		return 0
	}))

	t.RawSetString("use_default_input", r.State.NewFunction(func(l *lua.LState) int {
		pth := path.Join(r.Config.InputDirectory, r.Entry)
		err := os.MkdirAll(pth, 0o777)
		if err != nil {
			Error(r.State, lg.Appendf("failed to create default input directory %s, with error (%s)", log.LEVEL_ERROR, pth, err))
		}
		r.UseDefaultInput = true

		return 0
	}))

	t.RawSetString("use_default_output", r.State.NewFunction(func(l *lua.LState) int {
		pth := path.Join(r.Config.OutputDirectory, r.Entry)
		err := os.MkdirAll(pth, 0o777)
		if err != nil {
			Error(r.State, lg.Appendf("failed to create default output directory %s, with error (%s)", log.LEVEL_ERROR, pth, err))
		}
		r.UseDefaultOutput = true

		return 0
	}))

	return t
}

func MapSchema(schema, data map[string]any) map[string]any {
	result := map[string]any{}

	for k, v := range schema {
		if d, ok := data[k]; ok {
			if v1, ok := v.([]any); ok {
				if d1, ok := d.([]any); ok {
					result[k] = d1
				} else {
					result[k] = v1
				}
			} else if v1, ok := v.(map[string]any); ok {
				if d1, ok := d.(map[string]any); ok {
					result[k] = MapSchema(v1, d1)
				} else {
					result[k] = d
				}
			} else {
				result[k] = d
			}
		} else {
			result[k] = v
		}
	}

	return result
}

func (r *Runner) WorkflowConfig(l *lua.LState, ext string) map[string]any {
	cfg := l.CheckTable(-1)
	v := GetValue(cfg).(map[string]any)

	fpath := path.Join(r.Config.ConfigDirectory, strings.ReplaceAll(r.Entry, "/", "_")+ext)

	_, err := os.Stat(fpath)
	if err != nil {
		err := jsonWrite(fpath, v)
		if err != nil {
			l.Error(lua.LString(r.lg.Append(err.Error(), log.LEVEL_ERROR)), 0)
		}

		return v
	}

	b, err := os.ReadFile(fpath)
	if err != nil {
		l.Error(lua.LString(r.lg.Append(fmt.Sprintf("failed to read config file: %s", err), log.LEVEL_ERROR)), 0)
	}

	vs := map[string]any{}
	err = json.Unmarshal(b, &vs)
	if err != nil {
		l.Error(lua.LString(r.lg.Append(fmt.Sprintf("failed to unmarshal config: %s", err), log.LEVEL_ERROR)), 0)
	}

	result := MapSchema(v, vs)
	err = jsonWrite(fpath, result)
	if err != nil {
		l.Error(lua.LString(r.lg.Append(err.Error(), log.LEVEL_ERROR)), 0)
	}

	return result
}

func jsonWrite(fpath string, data map[string]any) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %s", err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %s", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write config file: %s", err)
	}

	return nil
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
	VARIADIC
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

func ArgVariadic(name string, arrType ArrayType, optional bool) Arg {
	table := []Arg{
		{Type: arrType.Type, Name: name, Table: arrType.Table},
	}

	return Arg{
		Type:     VARIADIC,
		Name:     name,
		Optional: optional,
		Table:    &table,
	}
}

func (l *Lib) ParseValue(state *lua.LState, pos int, value lua.LValue, arg Arg) any {
	switch arg.Type {
	case INT:
		if value.Type() == lua.LTNumber {
			return int(math.Round(float64(value.(lua.LNumber))))
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case FLOAT:
		if value.Type() == lua.LTNumber {
			return float64(value.(lua.LNumber))
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid number provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case BOOL:
		if value.Type() == lua.LTBool {
			return bool(value.(lua.LBool))
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid bool provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case STRING:
		if value.Type() == lua.LTString {
			return string(value.(lua.LString))
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid string provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case ANY:
		return value

	case FUNC:
		if value.Type() == lua.LTFunction {
			return value
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid function provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case RAW_TABLE:
		if value.Type() == lua.LTTable {
			return value
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid table provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case ARRAY:
		if value.Type() == lua.LTTable {
			v := value.(*lua.LTable)

			m := make([]any, v.Len())

			for i := range v.Len() {
				m[i] = l.ParseValue(state, pos, v.RawGetInt(i+1), Arg{
					Type:  (*arg.Table)[0].Type,
					Name:  strconv.Itoa(i + 1),
					Table: (*arg.Table)[0].Table,
				})
			}
			return m
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid array provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	case TABLE:
		if value.Type() == lua.LTTable {
			v := value.(*lua.LTable)
			m := map[string]any{}

			for _, a := range *arg.Table {
				m[a.Name] = l.ParseValue(state, pos, v.RawGetString(a.Name), a)
			}

			return m
		} else if value.Type() == lua.LTNil && arg.Optional {
			return l.getDefault(arg, state)
		} else {
			state.ArgError(pos, l.Lg.Append(fmt.Sprintf("invalid table provided to %s: %+v", arg.Name, value), log.LEVEL_ERROR))
		}

	default:
		state.ArgError(pos, l.Lg.Append(fmt.Sprintf("unsupport arg type provided: %T", value), log.LEVEL_ERROR))
	}

	return lua.LNil
}

func (l *Lib) ParseArgs(state *lua.LState, name string, args []Arg, ln, level int) (map[string]any, int) {
	argMap := map[string]any{}
	count := 0

	for i, a := range args {
		ind := -ln + i
		v := state.CheckAny(ind)

		if a.Type == VARIADIC {
			if ind < 0 {
				m := make([]any, ind*-1)

				vi := 0
				for ; ind < 0; ind++ {
					m[vi] = l.ParseValue(state, i, state.CheckAny(ind), Arg{
						Type:  (*a.Table)[0].Type,
						Name:  strconv.Itoa(vi + 1),
						Table: (*a.Table)[0].Table,
					})
					vi++
				}

				argMap[a.Name] = m
				count++
			} else if a.Optional {
				argMap[a.Name] = l.getDefault(a, state)
			} else {
				state.ArgError(i, l.Lg.Append(fmt.Sprintf("invalid variadic provided to %s (non-optional variadics require at least 1 value): %+v", a.Name, v), log.LEVEL_ERROR))
			}

			return argMap, count
		}

		if i >= ln {
			if !a.Optional {
				state.ArgError(i, l.Lg.Append(fmt.Sprintf("required arg not provided: %d", i), log.LEVEL_ERROR))
			} else {
				argMap[a.Name] = l.getDefault(a, state)
			}
		} else {
			argMap[a.Name] = l.ParseValue(state, i, v, a)
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
		return []any{}
	case VARIADIC:
		return []any{}

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

func Error(state *lua.LState, err string) {
	state.Error(lua.LString(err), 0)
}

func (l *Lib) CreateFunction(lib *lua.LTable, name string, args []Arg, fn func(state *lua.LState, d TaskData, args map[string]any) int) {
	lib.RawSetString(name, l.State.NewFunction(func(state *lua.LState) int {
		l.Lg.Append(fmt.Sprintf("%s.%s called.", l.Lib, name), log.LEVEL_VERBOSE)

		argMap, c := l.ParseArgs(state, name, args, state.GetTop(), 0)
		state.Pop(c)

		ret := fn(state, TaskData{Lib: l.Lib, Name: name}, argMap)

		l.Lg.Append(fmt.Sprintf("%s.%s finished.", l.Lib, name), log.LEVEL_VERBOSE)
		return ret
	}))
}

func (l *Lib) TableFunction(state *lua.LState, t *lua.LTable, name string, args []Arg, fn func(state *lua.LState, args map[string]any) int) {
	t.RawSetString(name, state.NewFunction(func(state *lua.LState) int {
		argMap, c := l.ParseArgs(state, name, args, state.GetTop(), 0)
		state.Pop(c)

		ret := fn(state, argMap)

		return ret
	}))
}

func (l *Lib) BuilderFunction(state *lua.LState, t *lua.LTable, name string, args []Arg, fn func(state *lua.LState, t *lua.LTable, args map[string]any)) {
	t.RawSetString(name, state.NewFunction(func(state *lua.LState) int {
		self := state.CheckTable(1)
		state.Remove(1)

		argMap, c := l.ParseArgs(state, name, args, state.GetTop(), 0)
		state.Pop(c)

		fn(state, self, argMap)

		state.Push(self)
		return 1
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

func CreateValue(value any, state *lua.LState) lua.LValue {
	switch v := value.(type) {
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case string:
		return lua.LString(v)

	case []any:
		t := state.NewTable()
		for _, va := range v {
			t.Append(CreateValue(va, state))
		}
		return t

	case map[string]any:
		t := state.NewTable()
		for k, va := range v {
			t.RawSetString(k, CreateValue(va, state))
		}
		return t

	default:
		return lua.LNil
	}
}

func GetValue(value lua.LValue) any {
	switch v := value.(type) {
	case lua.LNumber:
		if float64(v) == float64(int(v)) {
			return int(v)
		}
		return float64(v)
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		isNumeric := true
		v.ForEach(func(l1, l2 lua.LValue) {
			if l1.Type() != lua.LTNumber {
				isNumeric = false
			} else if float64(l1.(lua.LNumber)) != float64(int(l1.(lua.LNumber))) {
				isNumeric = false
			}
		})

		if isNumeric {
			t := []any{}
			v.ForEach(func(l1, l2 lua.LValue) {
				t = append(t, GetValue(l2))
			})
			return t
		} else {
			t := map[string]any{}
			v.ForEach(func(l1, l2 lua.LValue) {
				t[l1.String()] = GetValue(l2)
			})
			return t
		}

	default:
		return nil
	}
}
